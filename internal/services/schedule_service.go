package services

import (
	//"dolittle2/internal/config"
	"dolittle2/internal/config"
	"dolittle2/internal/models"
	"dolittle2/internal/repos"
	"dolittle2/internal/utils"
	"errors"
	"log"
	"time"
)

type ScheduleService struct {
	Repo *repos.ScheduleRepo
}

func NewService(repo *repos.ScheduleRepo) *ScheduleService {
	return &ScheduleService{Repo: repo}
}

func (s *ScheduleService) CreateSchedule(schedule *models.Schedule) (uint, error) {
	exists, err := s.Repo.AidNameExists(schedule.Aid_name, schedule.UserID)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, errors.New("Запись с таким именем для пользователя уже существует")
	}

	hpur := time.Now().Hour()

	if hpur < 8 || hpur > 22 {
		return 0, errors.New("Лекарства можно принимать только с 8:00 до 22:00")
	}

	id, err := s.Repo.CreateSchedule(schedule)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *ScheduleService) FindByUserID(userID uint) ([]uint, error) {
	return s.Repo.FindByUserID(userID)
}

func (s *ScheduleService) GetDailySchedule(userID, scheduleID uint) ([]time.Time, error) {
	schedule, err := s.Repo.FindSchedule(userID, scheduleID)
	if err != nil {
		return nil, err
	}

	return utils.GenerateScheduleTimes(time.Now(), schedule.Aid_per_day)
}

func (s *ScheduleService) GetNextTakings(userID uint) (map[string][]string, error) {
	now := time.Now()
	end := now.Add(config.LoadConfig())

	schedules, err := s.Repo.NextTakings(userID)
	if err != nil {
		return nil, err
	}

	nextTakings := make(map[string][]string)

	for _, schedule := range schedules {
		times, err := utils.GenerateScheduleTimes(now, schedule.Aid_per_day)
		if err != nil {
			log.Fatal("Ошибка генерации расписания", err)
		}

		var nextPer []time.Time
		createEndTime := schedule.Create_at.Add(time.Duration(schedule.Duration*24) * time.Hour)

		for _, t := range times {
			if t.After(now) && t.Before(end) && t.Before(createEndTime) {
				nextPer = append(nextPer, t)
			}
		}

		var formattedTimes []string
		for _, t := range nextPer {
			formattedTimes = append(formattedTimes, t.Format("15:04"))
		}

		if formattedTimes != nil {
			nextTakings[schedule.Aid_name] = formattedTimes
		}
	}

	if len(nextTakings) == 0 {
		return nil, errors.New("Нет ближайших приемов")
	}

	return nextTakings, nil
}
