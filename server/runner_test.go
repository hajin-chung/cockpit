package main

import (
	"log/slog"
	"testing"
)

type MockDB struct{}

func (db *MockDB) NewCommand(command string) (*CommandInfo, error) {
	commandInfo := CommandInfo{
		Id:        IdGen(),
		CreatedAt: FormatNow(), 
		Status:    COMMAND_IDLE,
		Command:   command,
	}
	slog.Info("[TEST] new command", "command", command, "info", commandInfo)
	return &commandInfo, nil
}

func (db *MockDB) GetCommand(id string) (*CommandInfo, error) {
	return nil, nil
}

func (db *MockDB) ListCommands(before string, n uint) ([]CommandInfo, error) {
	return []CommandInfo{}, nil
}

func (db *MockDB) AddLog(log *CommandLog) error {
	slog.Info("[DB] ", "content", log.Content, "time", log.CreatedAt)
	return nil
}

func (db *MockDB) UpdateStatus(id string, status CommandStatus) error {
	slog.Info("[TEST] update status", "id", id, "status", status)
	return nil
}

func (db *MockDB) GetLogs(commandId string, before string, n uint) ([]CommandLog, error) {
	return nil, nil
}

func TestRunner(t *testing.T) {
	runner := NewRunner()
	db := &MockDB{}

	// commandInfo, err := db.NewCommand("tail -f /mnt/d/vod/memo.dat")
	// commandInfo, err := db.NewCommand("ls -alh /mnt/d/vod")
	commandInfo, err := db.NewCommand("ls -alh")
	if err != nil {
		t.Errorf("db NewCommand error: %s\n", err)
	}

	runner.Run(db, commandInfo)

	id := commandInfo.Id
	pipeOut := runner.AddConsumer(id)
	for log := range pipeOut {
		slog.Info("[OUT]", "content", log.Content, "time", log.CreatedAt)
	}
}

