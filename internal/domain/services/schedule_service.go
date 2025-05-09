package services

import (
	"dolittle2/internal/domain/models"
	"dolittle2/internal/utils"
	"errors"
	"log/slog"
	"time"
)

type scheduleRepo interface {
	CreateSchedule(schedule *models.Schedule) (uint64, error)
	AidNameExists(aidName string, userID uint64) (bool, error)
	UserIDExists(userID uint64) (bool, error)
	FindByUserID(userID uint64) ([]uint64, error)
	FindSchedule(userID, scheduleID uint64) (*models.Schedule, error)
	NextTakings(userID uint64) ([]models.Schedule, error)
}

type ScheduleService struct {
	scheduleRepo scheduleRepo
	endTime      time.Duration
}

func NewService(scheduleRepo scheduleRepo, endTime time.Duration) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
		endTime:      endTime,
	}
}

func (s *ScheduleService) CreateSchedule(schedule *models.Schedule) (uint64, error) {
	slog.Info(
		"Создание расписания",
		slog.String("aid_name", schedule.AidName),
		slog.Uint64("user_id", schedule.UserID),
	)

	exists, err := s.scheduleRepo.AidNameExists(schedule.AidName, schedule.UserID)
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
			slog.String("aid_name", schedule.AidName),
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

	exists, err := s.scheduleRepo.UserIDExists(userID)
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
		slog.Error("ошибка нахождения расписания",
			slog.Uint64("user_id", userID),
			slog.Uint64("schedule_id", scheduleID),
			slog.String("error", err.Error()))
		return nil, err
	}

	res, err := utils.GenerateScheduleTimes(time.Now(), schedule.AidPerDay, time.Now)
	if err != nil {
		slog.Info("ошибка гена расписания")
		return nil, err
	}
	return res, nil
}

func (s *ScheduleService) GetNextTakings(userID uint64) (map[string][]string, error) {
	slog.Info(
		"Запуск сервиса следующих приемов для пользователя",
		slog.Uint64("user_id", userID),
	)

	now := time.Now()
	end := now.Add(s.endTime)

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
		endTime := schedule.CreatedAt.Add(time.Duration(schedule.Duration) * 24 * time.Hour)
		if now.After(endTime) {
			slog.Error(
				"Время приема лекарства истекло",
				slog.Uint64("user_id", userID),
				slog.Uint64("schedule_id", schedule.ID),
				slog.Any("end", endTime),
			)
			continue
		}

		times, err := utils.GenerateScheduleTimes(now, schedule.AidPerDay, time.Now)
		if err != nil {
			slog.Error(
				"Ошибка генерации времени расписания",
				slog.String("error", err.Error()),
				slog.String("aid_name", schedule.AidName),
				slog.Uint64("user_id", userID),
			)
			return nil, err
		}

		var nextPer []time.Time
		createEndTime := schedule.CreatedAt.Add(time.Duration(schedule.Duration*24) * time.Hour)

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
				slog.String("Название лекарства", schedule.AidName),
				slog.Any("Время", formattedTimes),
			)

			nextTakings[schedule.AidName] = formattedTimes
		}
	}

	if len(nextTakings) == 0 {
		return nil, errors.New("Нет ближайших приемов")
	}

	slog.Info("nextTakings", slog.Any("&&", nextTakings))
	return nextTakings, nil
}
