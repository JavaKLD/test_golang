package server

import (
	"dolittle2/internal/domain/models"
	"dolittle2/internal/utils"
	openapi "dolittle2/openapi/gen/go"
	"time"

	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

type scheduleService interface {
	CreateSchedule(schedule *models.Schedule) (uint64, error)
	FindByUserID(userId uint64) ([]uint64, error)
	CheckUserExists(userId uint64) (bool, error)
	GetDailySchedule(userId, scheduleID uint64) ([]time.Time, error)
	GetNextTakings(userID uint64) (map[string][]string, error)
}

type ScheduleRestServer struct {
	scheduleService scheduleService
}

func NewScheduleController(scheduleService scheduleService) *ScheduleRestServer {
	return &ScheduleRestServer{scheduleService: scheduleService}
}

func (s *ScheduleRestServer) PostSchedule(ctx echo.Context) error {
	var req openapi.ScheduleRequest

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "неправильный формат запроса"})
	}

	schedule := &models.Schedule{
		Aid_name:    req.AidName,
		Aid_per_day: uint64(req.AidPerDay),
		Duration:    int64(req.Duration),
		UserID:      uint64(req.UserId),
		Create_at:   utils.RoundTime(time.Now()),
	}

	id, err := s.scheduleService.CreateSchedule(schedule)
	if err != nil {
		if err.Error() == "Запись с таким именем для пользователя уже существует" {
			return ctx.JSON(http.StatusConflict, map[string]string{"error": "Запись с таким aid_name для данного пользователя уже существует"})
		} else {
			return ctx.JSON(http.StatusUnprocessableEntity, "Лекарства принимаются с 8 до 22")
		}
	}
	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "Запись создана",
	})
}

func (s *ScheduleRestServer) GetUserSchedule(ctx echo.Context) error {
	queryParam := strings.TrimSpace(ctx.QueryParam("user_id"))
	if queryParam == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "не указан user_id"})
	}

	userID, err := strconv.ParseUint(queryParam, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "неверный формат user_id"})
	}

	scheduleID, err := s.scheduleService.FindByUserID(userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось получить данные"})
	}

	if len(scheduleID) == 0 {
		return ctx.JSON(http.StatusOK, []uint{})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"schedules": scheduleID,
		"message":   "Успешный ответ с расписанием",
	})
}

func (s *ScheduleRestServer) GetSchedule(ctx echo.Context) error {
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

	scheduleTimes, err := s.scheduleService.GetDailySchedule(userID, scheduleID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка вывода графика приема лекарств"})
	}

	var formattedTimes []string
	for _, t := range scheduleTimes {
		formattedTimes = append(formattedTimes, t.Format("15:04"))
	}

	return ctx.JSON(http.StatusOK, map[string][]string{"schedule": formattedTimes})

}

func (s *ScheduleRestServer) GetNextTakings(ctx echo.Context) error {
	queryParam := strings.TrimSpace(ctx.QueryParam("user_id"))

	if queryParam == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Отсутствует user_id"})
	}

	userID, err := strconv.ParseUint(queryParam, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат user_id"})
	}

	nextTakings, err := s.scheduleService.GetNextTakings(userID)
	if err != nil {
		if err.Error() == "Нет ближайших приемов" {
			return ctx.JSON(http.StatusOK, map[string]string{"message": "Нет ближайших приемов"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка получения данных"})
	}

	return ctx.JSON(http.StatusOK, map[string]map[string][]string{"schedule": nextTakings})
}
