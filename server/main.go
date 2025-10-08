package main

import (
	_ "embed"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed build/index.html
var IndexHTML string

func main() {
	bus := NewEventBus()
	CreateTopic[any](bus, "command")
	runner := NewRunner(bus)
	db, err := NewDB("file:cockpit.db", bus)
	if err != nil {
		slog.Error("failed to init db", "error", err)
		return
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.Use(CockpitContextMiddleware(runner, db, bus))

	e.GET("/test/sse", TestSSE)
	e.POST("/api/v1/command/new", NewCommandHandler)
	e.GET("/api/v1/command/:id", GetCommandHandler)
	e.GET("/api/v1/command/list", ListCommandHandler)
	e.POST("/api/v1/command/:id/stop", StopCommandHandler)
	e.DELETE("/api/v1/command/:id", DeleteCommandHandler)
	e.GET("/api/v1/command/stream", CommandStreamHandler)
	e.GET("/api/v1/command/:id/log/stream", LogStreamHandler)
	e.GET("/api/v1/command/:id/log", LogHandler)

	e.GET("/*", func(c echo.Context) error {
		return c.HTML(http.StatusOK, IndexHTML)
	})

	if err := e.Start(":4000"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

