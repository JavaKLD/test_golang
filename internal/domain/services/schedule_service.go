package services

import (
	"dolittle2/internal/config"
	"dolittle2/internal/domain/models"
	"dolittle2/internal/utils"
	"errors"
	"log/slog"
	"time"
)

type scheduleRepo interface {
	CreateSchedule(schedule *models.Schedule) (uint64, error)
	AidNameExists(aidName string, userID uint64) (bool, error)
	UserIdExists(userID uint64) (bool, error)
	FindByUserID(userID uint64) ([]uint64, error)
	FindSchedule(userID, scheduleID uint64) (*models.Schedule, error)
	NextTakings(userID uint64) ([]models.Schedule, error)
}

type ScheduleService struct {
	scheduleRepo scheduleRepo
}

func NewService(scheduleRepo scheduleRepo) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
	}
}

func (s *ScheduleService) CreateSchedule(schedule *models.Schedule) (uint64, error) {
	slog.Info(
		"Создание расписания",
		slog.String("aid_name", schedule.Aid_name),
		slog.Uint64("user_id", schedule.UserID),
	)

	exists, err := s.scheduleRepo.AidNameExists(schedule.Aid_name, schedule.UserID)
	if err != nil {
		slog.Error(
			"Ошибка при проверке существования имени лекарства",
			slog.String("error", err.Error()),
		)
		return 0, err
	}

	if exists {
		slog.Error(
			"Расписание с таким именем уже существует",
			slog.String("aid_name", schedule.Aid_name),
		)
		return 0, errors.New("Запись с таким именем для пользователя уже существует")
	}

	hpur := time.Now().Hour()

	if hpur < 8 || hpur > 22 {
		slog.Error(
			"Попытка создания расписания вне допустимого времени",
			slog.Int("current_hour", hpur),
		)
		return 0, errors.New("Лекарства можно принимать только с 8:00 до 22:00")
	}

	id, err := s.scheduleRepo.CreateSchedule(schedule)
	if err != nil {
		slog.Error(
			"Ошибка при создании расписания",
			slog.String("error", err.Error()),
		)
		return 0, err
	}

	slog.Info(
		"Расписание успешно создано",
		slog.Uint64("schedule_id", id),
	)
	return id, nil
}

func (s *ScheduleService) FindByUserID(userID uint64) ([]uint64, error) {
	slog.Info(
		"Поиск всех расписаний для пользователя",
		slog.Uint64("user_id", userID),
	)

	schedules, err := s.scheduleRepo.FindByUserID(userID)
	if err != nil {
		slog.Error(
			"Ошибка при поиске расписаний для пользователя",
			slog.Uint64("user_id", userID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	slog.Info(
		"Найдено расписаний",
		slog.Uint64("user_id", userID),
		slog.Int("count", len(schedules)),
	)

	return schedules, nil
}

func (s *ScheduleService) CheckUserExists(userID uint64) (bool, error) {
	slog.Info(
		"Проверка существования пользователя",
		slog.Uint64("user_id", userID),
	)

	exists, err := s.scheduleRepo.UserIdExists(userID)
	if err != nil {
		slog.Error(
			"Ошибка при проверке существования пользователя",
			slog.Uint64("user_id", userID),
			slog.String("error", err.Error()),
		)
		return false, err
	}

	if exists {
		slog.Info(
			"Пользователь существует",
			slog.Uint64("user_id", userID),
		)
	} else {
		slog.Warn(
			"Пользователь не существует",
			slog.Uint64("user_id", userID),
		)
	}

	return exists, nil
}

func (s *ScheduleService) GetDailySchedule(userID, scheduleID uint64) ([]time.Time, error) {
	slog.Info(
		"Получение расписания на день для user для конкретного расписания",
		slog.Uint64("user_id", userID),
		slog.Uint64("schedule_id", scheduleID),
	)

	schedule, err := s.scheduleRepo.FindSchedule(userID, scheduleID)
	if err != nil {
		slog.Error("Ошибка нахождения расписания",
			slog.Uint64("user_id", userID),
			slog.Uint64("schedule_id", scheduleID),
			slog.String("error", err.Error()))
		return nil, err
	}

	return utils.GenerateScheduleTimes(time.Now(), schedule.Aid_per_day)
}

func (s *ScheduleService) GetNextTakings(userID uint64) (map[string][]string, error) {
	slog.Info(
		"Запуск сервиса следующих приемов для пользователя",
		slog.Uint64("user_id", userID),
	)
	now := time.Now()
	end := now.Add(config.LoadConfig())

	schedules, err := s.scheduleRepo.NextTakings(userID)
	if err != nil {
		slog.Error(
			"Ошибка при получении расписания для след приемов",
			slog.Uint64("user_id", userID),
		)
		return nil, err
	}

	nextTakings := make(map[string][]string)

	for _, schedule := range schedules {
		times, err := utils.GenerateScheduleTimes(now, schedule.Aid_per_day)
		if err != nil {
			slog.Error(
				"Ошибка генерации времени расписания",
				slog.String("error", err.Error()),
				slog.String("aid_name", schedule.Aid_name),
				slog.Uint64("user_id", userID),
			)
			return nil, err
		}

		var nextPer []time.Time
		createEndTime := schedule.Create_at.Add(time.Duration(schedule.Duration*24) * time.Hour)

		for _, t := range times {
			if t.After(now) && t.Before(end) && t.Before(createEndTime) {
				slog.Info(
					"Добавление подходящего времени",
					slog.Time("time", t),
				)
				nextPer = append(nextPer, t)
			}
		}

		var formattedTimes []string
		for _, t := range nextPer {
			formattedTimes = append(formattedTimes, t.Format("15:04"))
		}

		if formattedTimes != nil {
			slog.Info(
				"Добавление времени для лекарства",
				slog.String("Название лекарства", schedule.Aid_name),
				slog.Any("Время", formattedTimes),
			)
			nextTakings[schedule.Aid_name] = formattedTimes
		}
	}

	if len(nextTakings) == 0 {
		return nil, errors.New("Нет ближайших приемов")
	}

	return nextTakings, nil
}
