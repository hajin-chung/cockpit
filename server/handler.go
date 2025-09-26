package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type NewCommand struct {
	Command string `json:"command"`
}

func NewCommandHandler(c echo.Context) error {
	cc := c.(*CockpitContext)
	newCommand := new(NewCommand)
	if err := cc.Bind(newCommand); err != nil {
		slog.Error("NewCommandHandler cc.Bind", "error", err)
		return cc.String(http.StatusBadRequest, "invalid json format");
	}

	commandInfo, err := cc.DB.NewCommand(newCommand.Command)
	if err != nil {
		slog.Error("NewCommandHandler cc.DB.NewCommand", "error", err)
		return cc.String(http.StatusInternalServerError, "db fail");
	}

	err = cc.Runner.Run(cc.DB, commandInfo)
	if err != nil {
		slog.Error("NewCommandHandler cc.Runner.Run", "error", err)
		return cc.String(http.StatusInternalServerError, "runner fail");
	}

	return cc.JSON(http.StatusCreated, commandInfo)
}

func GetCommandHandler(c echo.Context) error {
	cc := c.(*CockpitContext)

	commandId := cc.Param("id")
	info, err := cc.DB.GetCommand(commandId)
	if err != nil {
		slog.Error("GetCommandHandler cc.DB.GetCommand", "error", err)
		return cc.String(http.StatusInternalServerError, "db fail");
	}
	return cc.JSON(http.StatusOK, info)
}

func ListCommandHandler(c echo.Context) error {
	cc := c.(*CockpitContext)

	before := cc.QueryParam("before")
	limit, err := strconv.Atoi(cc.QueryParam("limit"))
	if err != nil{
		return cc.String(http.StatusBadRequest, "invalid limit param");
	}
	if limit < 0 {
		return cc.String(http.StatusBadRequest, "negative limit param");
	}

	commands, err := cc.DB.ListCommands(before, uint(limit))
	if err != nil {
		slog.Error("ListCommandHandler cc.DB.ListCommand", "error", err)
		return cc.String(http.StatusInternalServerError, "db fail");
	}

	return cc.JSON(http.StatusOK, commands)
}

func LogHandler(c echo.Context) error {
	cc := c.(*CockpitContext)
	commandId := cc.Param("id")
	before := cc.QueryParam("before")
	limit, err := strconv.Atoi(cc.QueryParam("limit"))
	if err != nil{
		return cc.String(http.StatusBadRequest, "invalid limit param");
	}
	if limit < 0 {
		return cc.String(http.StatusBadRequest, "negative limit param");
	}

	logs, err := cc.DB.GetLogs(commandId, before, uint(limit))
	if err != nil {
		slog.Error("LogHandler cc.DB.GetLogs", "error", err)
		return cc.String(http.StatusInternalServerError, "db fail");
	}

	return cc.JSON(http.StatusOK, logs)
}

func LogStreamHandler(c echo.Context) error {
	cc := c.(*CockpitContext)
	commandId := cc.Param("id")

	rc, err := cc.Runner.AddConsumer(commandId)
	if err != nil {
		slog.Error("LogStreamHandler cc.Runner.AddConsumer", "error", err)
		return cc.String(http.StatusInternalServerError, "runner fail");
	}

	w := cc.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case <-cc.Request().Context().Done():
			cc.Runner.CloseConsumer(commandId, rc)
			return nil
		case log := <-rc:
			data, err := json.Marshal(log)
			if err != nil {
				slog.Error("LogStreamHandler json.Marshal(log)", "error", err)
				continue
			}
			event := Event { Data: data }

			if err := event.MarshalTo(w); err != nil {
				return err
			}
			w.Flush()
		}
	}
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
				Data: []byte("time: " + FormatNow()),
			}
			if err := event.MarshalTo(w); err != nil {
				return err
			}
			w.Flush()
		}
	}
}
