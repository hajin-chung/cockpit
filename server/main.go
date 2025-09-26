package main

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	bus := NewEventBus()
	runner := NewRunner()
	db, err := NewDB("file:cockpit.db")
	if err != nil {
		slog.Error("failed to init db", "error", err)
		return
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.Use(CockpitContextMiddleware(runner, db, bus))
	e.Static("/", "static")

	e.GET("/test/sse", TestSSE)
	e.POST("/api/v1/command/new", NewCommandHandler)
	e.GET("/api/v1/command/:id", GetCommandHandler)
	e.GET("/api/v1/command/list", ListCommandHandler)
	e.GET("/api/v1/command/:id/log/stream", LogStreamHandler)
	e.GET("/api/v1/command/:id/log", LogHandler)

	if err := e.Start(":4000"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
