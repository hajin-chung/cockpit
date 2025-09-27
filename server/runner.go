package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"sync"
)

type Runner interface {
	Run(db DB, command *Command) error
}

type Session struct {
	*Command
	Consumers *Log
}

type CockpitRunner struct {
	Bus      *EventBus
	Sessions map[string]*Session
}

func NewRunner(bus *EventBus) Runner {
	sessions := make(map[string]*Session)
	runner := CockpitRunner{
		Sessions: sessions,
		Bus:      bus,
	}
	return &runner
}

func (r *CockpitRunner) Run(db DB, command *Command) error {
	cmd := exec.Command("bash", "-c", command.Command)

	session := &Session{
		Command: command,
	}
	r.Sessions[command.Id] = session

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("cannot get stdout", "command", command.Command, "error", err)
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("cannot get stderr", "command", command.Command, "error", err)
		return err
	}

	_, err = CreateTopic[*Log](r.Bus, command.Id)
	if err != nil {
		slog.Error("CockpitRunner.Run", "error", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go session.Drainer(&wg, r.Bus, command, stdout, LOG_STDOUT)
	go session.Drainer(&wg, r.Bus, command, stderr, LOG_STDERR)
	go session.Logger(db, r.Bus, command)
	go session.Waiter(&wg, db, r.Bus, command, cmd)

	return nil
}

func SplitLines(buf []byte) ([]string, int) {
	lines := []string{}
	idx := 0
	for i, b := range buf {
		if b == '\n' {
			bufline := make([]byte, i-idx)
			copy(bufline, buf[idx:i])
			lines = append(lines, string(bufline))
			idx = i + 1
		}
	}
	return lines, idx
}

// read pipe and write to channel
func (s *Session) Drainer(wg *sync.WaitGroup, bus *EventBus, command *Command, reader io.ReadCloser, fd LogFD) {
	defer wg.Done()
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		log := &Log{
			Id:        IdGen(),
			CommandId: s.Id,
			CreatedAt: FormatNow(),
			Content:   line,
			FD:        fd,
		}
		slog.Info("[IN] ", "content", line, "time", log.CreatedAt)
		Pub(bus, log.CommandId, log)
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Drainer", "error", err)
	}
}

// write log to db
func (s *Session) Logger(db DB, bus *EventBus, command *Command) {
	_, err := Sub(bus, command.Id, func(log *Log) {
		db.AddLog(log)
	})
	if err != nil {
		slog.Error("Logger", "error", err)
	}
}

type UpdateCommandMessage struct {
	Id     string        `json:"id"`
	Status CommandStatus `json:"status"`
}

// resposible for startup and cleanup
func (s *Session) Waiter(wg *sync.WaitGroup, db DB, bus *EventBus, command *Command, cmd *exec.Cmd) {
	if err := cmd.Start(); err != nil {
		slog.Error("failed to start command", "command", s.Command, "error", err)

		db.UpdateStatus(s.Id, COMMAND_ERROR)
		db.AddLog(&Log{
			IdGen(),
			s.Id,
			FormatNow(),
			fmt.Sprintf("failed to start command %s error: %s", s.Command, err),
			-1,
		})

		msg := UpdateCommandMessage{command.Id, COMMAND_ERROR}
		if err := Pub[any](bus, "command", msg); err != nil {
			slog.Error("failed to send update command message", "message", msg, "error", err)
		}

		return
	}
	db.UpdateStatus(s.Id, COMMAND_RUNNING)
	msg := UpdateCommandMessage{command.Id, COMMAND_RUNNING}
	if err := Pub[any](bus, "command", msg); err != nil {
		slog.Error("failed to send update command message", "message", msg, "error", err)
	}
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		slog.Error("failed to wait command", "command", s.Command, "error", err)

		db.UpdateStatus(s.Id, COMMAND_ERROR)
		db.AddLog(&Log{
			IdGen(),
			s.Id,
			FormatNow(),
			fmt.Sprintf("failed to wait command %s error: %s", s.Command, err),
			-1,
		})

		msg := UpdateCommandMessage{command.Id, COMMAND_ERROR}
		if err := Pub[any](bus, "command", msg); err != nil {
			slog.Error("failed to send update command message", "message", msg, "error", err)
		}

		return
	}
	db.UpdateStatus(s.Id, COMMAND_EXITED)
	msg = UpdateCommandMessage{command.Id, COMMAND_EXITED}
	if err := Pub[any](bus, "command", msg); err != nil {
		slog.Error("failed to send update command message", "message", msg, "error", err)
	}

	CloseTopic[*Log](bus, command.Id)
}

