package routers

import (
	"dolittle2/internal/controllers"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, controller *controllers.ScheduleController) {
	e.POST("/schedule", controller.CreateSchedule)
	e.GET("/schedule", controller.GetUserSchedule)
	e.GET("/schedule/:schedule_id", controller.GetSchedule)
	e.GET("/next_takings", controller.GetNextTakings)
}
