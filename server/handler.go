package main

import (
	"encoding/json"
	"errors"
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
	if err := c.Bind(newCommand); err != nil {
		return err
	}

	commandInfo, err := cc.DB.NewCommand(newCommand.Command)
	if err != nil {
		return err
	}

	err = cc.Runner.Run(cc.DB, commandInfo)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, commandInfo)
}

func GetCommandHandler(c echo.Context) error {
	cc := c.(*CockpitContext)

	commandId := cc.Param("id")
	info, err := cc.DB.GetCommand(commandId)
	if err != nil {
		return nil
	}
	return c.JSON(http.StatusOK, info)
}

func ListCommandHandler(c echo.Context) error {
	cc := c.(*CockpitContext)

	before := cc.QueryParam("before")
	limit, err := strconv.Atoi(cc.QueryParam("limit"))
	if err != nil{
		return err
	}
	if limit < 0 {
		return errors.New("negative limit")
	}

	commands, err := cc.DB.ListCommands(before, uint(limit))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, commands)
}

func LogHandler(c echo.Context) error {
	cc := c.(*CockpitContext)
	commandId := cc.Param("id")
	before := cc.QueryParam("before")
	limit, err := strconv.Atoi(cc.QueryParam("limit"))
	if err != nil{
		return err
	}
	if limit < 0 {
		return errors.New("negative limit")
	}

	logs, err := cc.DB.GetLogs(commandId, before, uint(limit))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, logs)
}

func LogStreamHandler(c echo.Context) error {
	cc := c.(*CockpitContext)
	commandId := cc.Param("id")

	rc := cc.Runner.AddConsumer(commandId)

	w := cc.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case <-c.Request().Context().Done():
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
