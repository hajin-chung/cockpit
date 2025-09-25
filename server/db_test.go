package main

import (
	"os"
	"testing"
)

func TestDB(t *testing.T) {
	db, err := NewDB("file:test.db")
	if err != nil {
		t.Fatalf("NewDB error: %s\n", err)
	}

	defer func() {
		err := os.Remove("test.db")
		if err != nil {
			t.Errorf("os.Remove test db error: %s\n", err)
		}
	}()

	var info *CommandInfo
	t.Run("db command", func(t *testing.T) {
		info = testDBCommand(t, db)	
	})

	t.Run("db log", func(t *testing.T) {
		testDBLog(t, db, info)
	})
}

func testDBCommand(t *testing.T, db DB) *CommandInfo {
	info, err := db.NewCommand("ls -alh")
	if err != nil {
		t.Fatalf("NewCommand error: %s\n", err)
	}
	t.Logf("info: %v\n", info)

	infos, err := db.ListCommands("", 10)
	if err != nil {
		t.Fatalf("ListCommands error: %s\n", err)
	}

	t.Logf("total %d infos retrieved\n", len(infos))
	for i, info := range infos {
		t.Logf("infos[%d]: %v\n", i, info)
	}

	if len(infos) == 0 {
		t.Fatalf("no infos retrieved\n")
	}

	if infos[0].Id != info.Id {
		t.Errorf("commands differ!\n")
	}
	return info
}

func testDBLog(t *testing.T, db DB, info *CommandInfo) {
	log := &CommandLog {
		Id: IdGen(),
		CommandId: info.Id,
		CreatedAt: FormatNow(),
		Content: "hi",
		FD: LOG_STDOUT,
	}
	if err := db.AddLog(log); err != nil {
		t.Fatalf("AddLog error: %s\n", err)
	}
	t.Logf("log: %v\n", log)

	logs, err := db.GetLogs(info.Id, "", 10)
	if err != nil {
		t.Fatalf("GetLogs error: %s\n", err)
	}

	t.Logf("total %d logs retrieved\n", len(logs))
	for i, log := range logs {
		t.Logf("logs[%d]: %v\n", i, log)
	}

	if len(logs) == 0 {
		t.Fatalf("no logs retrieved\n")
	}

	if logs[0].Id != log.Id {
		t.Errorf("logs differ\n")
	}
}
