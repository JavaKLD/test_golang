package routers

import (
	"dolittle2/internal/controllers"
	"dolittle2/internal/middleware"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, controller *controllers.ScheduleController) {
	e.Use(middleware.ReqLogger)

	e.POST("/schedule", controller.CreateSchedule)
	e.GET("/schedules", controller.GetUserSchedule)
	e.GET("/schedule", controller.GetSchedule)
	e.GET("/next_takings", controller.GetNextTakings)
}
