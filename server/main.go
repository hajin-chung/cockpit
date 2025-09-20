package main

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var runner Runner = NewRunner()
	var db DB = NewDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(CockpitContextMiddleware(&runner, &db))
	e.Static("/", "static")

	e.GET("/test/sse", TestSSE)
	e.POST("/api/v1/command/new", CommandNewHandler)
	e.GET("/api/v1/command/list", CommandListHandler)
	// e.GET("/api/v1/command/log/:id", CommandLogHandler)

	if err := e.Start(":4000"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
