package controllers

import (
	"dolittle2/internal/models"
	"dolittle2/internal/services"
	"dolittle2/internal/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type ScheduleController struct {
	Service *services.ScheduleService
}

func NewScheduleController(service *services.ScheduleService) *ScheduleController {
	return &ScheduleController{Service: service}
}

func (c *ScheduleController) CreateSchedule(ctx echo.Context) error {
	var schedule models.Schedule
	err := ctx.Bind(&schedule)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "неправильный формат запроса"})
	}

	schedule.Create_at = utils.RoundTime(time.Now())

	id, err := c.Service.CreateSchedule(&schedule)
	if err != nil {
		if err.Error() == "Запись с таким именем для пользователя уже существует" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Запись с таким aid_name для данного пользователя уже существует"})
		} else {
			return ctx.JSON(http.StatusInternalServerError, "Лекарства принимаются с 8 до 22")
		}
	}
	return ctx.JSON(http.StatusOK, map[string]uint{"id": id})
}

func (c *ScheduleController) UserSchedule(ctx echo.Context) error {
	queryParam := strings.TrimSpace(ctx.QueryParam("user_id"))
	if queryParam == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "не указан user_id"})
	}

	userID, err := strconv.ParseUint(queryParam, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "неверный формат user_id"})
	}

	scheduleID, err := c.Service.FindByUserID(uint(userID))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось получить данные"})
	}

	if len(scheduleID) == 0 {
		return ctx.JSON(http.StatusOK, []uint{})
	}

	return ctx.JSON(http.StatusOK, map[string][]uint{"schedules": scheduleID})
}

func (c *ScheduleController) GetSchedule(ctx echo.Context) error {
	queryParamID := strings.TrimSpace(ctx.QueryParam("user_id"))
	queryParamSchedule := strings.TrimSpace(ctx.QueryParam("schedule_id"))

	if queryParamID == "" || queryParamSchedule == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Не указан user_id или schedule_id"})
	}

	userID, err := strconv.ParseUint(queryParamID, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат user_id"})
	}

	scheduleID, err := strconv.ParseUint(queryParamSchedule, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат schedule_id"})
	}

	scheduleTimes, err := c.Service.GetDailySchedule(uint(userID), uint(scheduleID))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка вывода графика приема лекарств"})
	}

	var formattedTimes []string
	for _, t := range scheduleTimes {
		formattedTimes = append(formattedTimes, t.Format("15:04"))
	}

	return ctx.JSON(http.StatusOK, map[string][]string{"schedule": formattedTimes})

}

func (c *ScheduleController) GetNextTakings(ctx echo.Context) error {
	queryParam := strings.TrimSpace(ctx.QueryParam("user_id"))

	if queryParam == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Отсутствует user_id"})
	}

	userID, err := strconv.ParseUint(queryParam, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат user_id"})
	}

	nextTakings, err := c.Service.GetNextTakings(uint(userID))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения данных"})
	}

	return ctx.JSON(http.StatusOK, map[string]map[string][]string{"schedule": nextTakings})
}
