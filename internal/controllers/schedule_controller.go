package controllers

import (
	"dolittle2/internal/models"
	"dolittle2/internal/services"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ScheduleController struct {
	Service *services.ScheduleService
}

func NewScheduleController(service *services.ScheduleService) *ScheduleController {
	return &ScheduleController{Service: service}
}

func (c *ScheduleController) CreateSchedule(ctx echo.Context) error {
	var schedule models.Schedule
	if err := ctx.Bind(&schedule); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "неправильный формат запроса"})
	}
	id, err := c.Service.CreateSchedule(&schedule)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err})
	}
	return ctx.JSON(http.StatusOK, map[string]uint{"id": id})
}

func (c *ScheduleController) UserSchedule (ctx echo.Context) error {
	userIDstr := ctx.QueryParam("user_id")
	if userIDstr == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "не указан user_id"})
	}

	userID, err := strconv.ParseUint(userIDstr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "неверный формат user_id"})
	}

	scheduleID, err := c.Service.FindByUserID(uint(userID))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "не удаось получить данные"})
	}

	if len(scheduleID) == 0 {
		return ctx.JSON(http.StatusOK, []uint{})
	}

	return ctx.JSON(http.StatusOK, scheduleID)
}
