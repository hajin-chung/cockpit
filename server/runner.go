package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"slices"
	"sync"
)

type Runner interface {
	Run(db DB, commandInfo *CommandInfo) error
	AddConsumer(commandId string) (chan *CommandLog, error)
	CloseConsumer(commandId string, c chan *CommandLog) error 
}

type Session struct {
	*CommandInfo
	Consumers []chan *CommandLog
}

type CockpitRunner struct {
	Sessions map[string]*Session
}

func NewRunner() Runner {
	sessions := make(map[string]*Session)
	runner := CockpitRunner{
		Sessions: sessions,
	}
	return &runner
}

func (r *CockpitRunner) Run(db DB, commandInfo *CommandInfo) error {
	cmd := exec.Command("bash", "-c", commandInfo.Command)

	consumers := make([]chan *CommandLog, 0)
	session := &Session{
		CommandInfo: commandInfo,
		Consumers:   consumers,
	}
	r.Sessions[commandInfo.Id] = session

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("cannot get stdout", "command", commandInfo.Command, "error", err)
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("cannot get stderr", "command", commandInfo.Command, "error", err)
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go session.Drainer(&wg, commandInfo, stdout, LOG_STDOUT)
	go session.Drainer(&wg, commandInfo, stderr, LOG_STDERR)
	go session.Logger(db, commandInfo)
	go session.Waiter(&wg, db, commandInfo, cmd)

	return nil
}

func (r *CockpitRunner) AddConsumer(commandId string) (chan *CommandLog, error) {
	session := r.Sessions[commandId]
	if session == nil {
		return nil, fmt.Errorf("no command with commandId %s", commandId)
	}
	c, err := session.AddConsumer()
	return c, err
}

func (r *CockpitRunner) CloseConsumer(commandId string, c chan *CommandLog) error {
	session := r.Sessions[commandId]
	if session == nil {
		return fmt.Errorf("no command with commandId %s", commandId)
	}
	err := session.CloseConsumer(c)
	return err
}

func (s *Session) AddConsumer() (chan *CommandLog, error) {
	if s.Status != COMMAND_RUNNING && s.Status != COMMAND_IDLE {
		return nil, errors.New("command closed")
	}
	c := make(chan *CommandLog)
	s.Consumers = append(s.Consumers, c)
	return c, nil
}

func (s *Session) CloseConsumer(c chan *CommandLog) error {
	s.Consumers = slices.DeleteFunc(s.Consumers, func(cc chan *CommandLog) bool {
		return cc == c
	})
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
func (s *Session) Drainer(wg *sync.WaitGroup, commandInfo *CommandInfo, reader io.ReadCloser, fd LogFD) {
	defer wg.Done()
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		commandLog := &CommandLog{
			Id:        IdGen(),
			CommandId: s.Id,
			CreatedAt: FormatNow(),
			Content:   line,
			FD:        fd,
		}
		slog.Info("[IN] ", "content", line, "time", commandLog.CreatedAt)
		for _, c := range s.Consumers {
			c <- commandLog
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Drainer", "error", err)
	}
}

// write log to db
func (s *Session) Logger(db DB, commandInfo *CommandInfo) {
	c, err := s.AddConsumer()
	if err != nil {
		slog.Error("Logger", "error", err)
	}
	for log := range c {
		db.AddLog(log)
	}
}

// resposible for startup and cleanup
func (s *Session) Waiter(wg *sync.WaitGroup, db DB, commandInfo *CommandInfo, cmd *exec.Cmd) {
	if err := cmd.Start(); err != nil {
		slog.Error("failed to start command", "command", s.Command, "error", err)
		db.UpdateStatus(s.Id, COMMAND_ERROR)
		db.AddLog(&CommandLog{
			IdGen(),
			s.Id,
			FormatNow(),
			fmt.Sprintf("failed to start command %s error: %s", s.Command, err),
			-1,
		})
		return
	}
	db.UpdateStatus(s.Id, COMMAND_RUNNING)
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		slog.Error("failed to wait command", "command", s.Command, "error", err)
		db.UpdateStatus(s.Id, COMMAND_ERROR)
		db.AddLog(&CommandLog{
			IdGen(),
			s.Id,
			FormatNow(),
			fmt.Sprintf("failed to wait command %s error: %s", s.Command, err),
			-1,
		})
		return
	}
	db.UpdateStatus(s.Id, COMMAND_EXITED)

	// TODO: more cleanup
	for _, pipeIn := range s.Consumers {
		close(pipeIn)
	}
}

