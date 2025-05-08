package server

import (
	"dolittle2/pkg/middleware"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, s *ScheduleRestServer) {
	e.Use(middleware.ReqLogger)

	e.POST("/schedule", s.PostSchedule)
	e.GET("/schedules", s.GetUserSchedule)
	e.GET("/schedule", s.GetSchedule)
	e.GET("/next_takings", s.GetNextTakings)
}
