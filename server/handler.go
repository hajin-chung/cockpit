package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type NewCommand struct {
	Command string `json:"command"`
}

func CommandNewHandler(c echo.Context) error {
	newCommand := new(NewCommand)
	if err := c.Bind(newCommand); err != nil {
		return err
	}
	slog.Info("POST new command", "command", newCommand.Command)

	// TODO: start command

	// TODO: return command info
	return c.JSON(http.StatusCreated, nil)
}

func CommandListHandler(c echo.Context) error {
	return nil
}

func TestSSE(c echo.Context) error {
	slog.Info("SSE client connected, ip: %v", c.RealIP(), "info")
	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.Request().Context().Done():
			slog.Info("SSE client disconnected, ip: %v", c.RealIP(), "info")
			return nil
		case <-ticker.C:
			event := Event{
				Data: []byte("time: " + time.Now().Format(time.RFC3339Nano)),
			}
			if err := event.MarshalTo(w); err != nil {
				return err
			}
			w.Flush()
		}
	}
}
