package main

import "github.com/labstack/echo/v4"

type CockpitContext struct {
	echo.Context

	Runner Runner
	DB     DB
	Bus    *EventBus
}

func CockpitContextMiddleware(runner Runner, db DB, bus *EventBus) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CockpitContext{
				Context: c,
				Runner:  runner,
				DB:      db,
				Bus:     bus,
			}
			return next(cc)
		}
	}
}

