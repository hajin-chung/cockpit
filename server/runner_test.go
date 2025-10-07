package main

import (
	"log/slog"
	"sync"
	"testing"
)

func TestChan(t *testing.T) {
	slog.Info("hi")
	var wgt sync.WaitGroup
	wgt.Add(2)

	ch := make(chan int)
	go func() {
		defer wgt.Done()
		for i := range ch {
			slog.Info("reciever", "i", i)
		}
		slog.Info("closing reciever\n")
	}()

	go func() {
		defer wgt.Done()
		for i := range 10 {
			slog.Info("hi", "i", i)
			ch <- i
		}
		slog.Info("closing ch\n")
		close(ch)
	}()
	wgt.Wait()
}

func TestRunner(t *testing.T) {
	bus := NewEventBus()
	CreateTopic[any](bus, "command")

	runner := NewRunner(bus)
	db, err := NewDB("file:test.db", bus)
	if err != nil {
		t.Errorf("NewDB error: %s\n", err)
	}

	go func() {
		rc, _, err := SubChan[any](bus, "command")
		if err != nil {
			t.Errorf("SubChan command error: %s\n", err)
		}
		for evt := range rc {
			msg := evt.(*CommandEvent)
			slog.Info("[SUB]", "event", msg.Type, "content", msg.Command)
		}
	}()

	// commandInfo, err := db.NewCommand("tail -f /mnt/d/vod/memo.dat")
	// commandInfo, err := db.NewCommand("ls -alh /mnt/d/vod")
	// commandInfo, err := db.NewCommand("ls -alh")
	command, err := db.NewCommand("while true; do date; sleep 1; done")
	if err != nil {
		t.Errorf("db NewCommand error: %s\n", err)
	}

	msg := CommandMessage(command, COMMAND_CREATE)
	err = Pub[any](bus, "command", msg)
	slog.Info("[PUB]", "msg", msg)
	if err != nil {
		slog.Error("[PUB]", "error", err)
	}

	runner.Run(db, command)

	// id := commandInfo.Id
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer slog.Info("[DONE]", "func", "top")
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
		defer slog.Info("[DONE]", "func", "bot")
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
				slog.Info("[STOP]", "id", command.Id)
				runner.Stop(command.Id)
				unsub()
				return
			}
		}
	}()

	wg.Wait()
}

