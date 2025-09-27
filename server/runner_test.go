package main

import (
	"log/slog"
	"sync"
	"testing"
)

type MockDB struct{}

func (db *MockDB) NewCommand(command string) (*Command, error) {
	commandInfo := Command{
		Id:        IdGen(),
		CreatedAt: FormatNow(),
		Status:    COMMAND_IDLE,
		Command:   command,
	}
	slog.Info("[TEST] new command", "command", command, "info", commandInfo)
	return &commandInfo, nil
}

func (db *MockDB) GetCommand(id string) (*Command, error) {
	return nil, nil
}

func (db *MockDB) ListCommands(before string, n uint) ([]Command, error) {
	return []Command{}, nil
}

func (db *MockDB) AddLog(log *Log) error {
	slog.Info("[DB] ", "content", log.Content, "time", log.CreatedAt)
	return nil
}

func (db *MockDB) UpdateStatus(id string, status CommandStatus) error {
	slog.Info("[TEST] update status", "id", id, "status", status)
	return nil
}

func (db *MockDB) GetLogs(commandId string, before string, n uint) ([]Log, error) {
	return nil, nil
}

func TestRunner(t *testing.T) {
	bus := NewEventBus()
	runner := NewRunner(bus)
	db := &MockDB{}

	// commandInfo, err := db.NewCommand("tail -f /mnt/d/vod/memo.dat")
	// commandInfo, err := db.NewCommand("ls -alh /mnt/d/vod")
	// commandInfo, err := db.NewCommand("ls -alh")
	command, err := db.NewCommand("while true; do date; sleep 1; done")
	if err != nil {
		t.Errorf("db NewCommand error: %s\n", err)
	}

	runner.Run(db, command)

	// id := commandInfo.Id
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		rc, _, err := SubChan[*Log](bus, command.Id)
		if err != nil {
			t.Errorf("runner.AddConsumer error: %s\n", err)
			return
		}
		for log := range rc {
			slog.Info("[OUT]", "content", log.Content, "time", log.CreatedAt)
		}
	}()

	go func() {
		defer wg.Done()
		rc, unsub, err := SubChan[*Log](bus, command.Id)
		if err != nil {
			t.Errorf("runner.AddConsumer error: %s\n", err)
			return
		}
		cnt := 0
		for log := range rc {
			cnt += 1
			slog.Info("[OUT]", "content", log.Content, "time", log.CreatedAt)
			if cnt > 5 {
				unsub()
				return
			}
		}
	}()

	wg.Wait()
}

