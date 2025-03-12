package services

import (
	"dolittle2/internal/models"
	"dolittle2/internal/repos"
	"dolittle2/internal/utils"
	"errors"
	"time"
)

type ScheduleService struct {
	Repo *repos.ScheduleRepo
}

func NewService(repo *repos.ScheduleRepo) *ScheduleService {
	return &ScheduleService{Repo: repo}
}

func (s *ScheduleService) CreateSchedule(schedule *models.Schedule) (uint, error) {
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

	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 8, 0, 0, 0, time.Local)
	endTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 22, 0, 0, 0, time.Local)

	timesPerDay := schedule.Aid_per_day
	if timesPerDay <= 0 {
		return nil, errors.New("Неверное количество приемов в день")
	}

	interval := (endTime.Sub(startTime)) / time.Duration(timesPerDay)

	var scheduleTimes []time.Time
	for i := 0; i < timesPerDay; i++ {
		appointmentTime := startTime.Add(time.Duration(i) * interval)
		roundedTime := utils.RoundTime(appointmentTime)
		scheduleTimes = append(scheduleTimes, roundedTime)
	}

	return scheduleTimes, nil
}

func (s *ScheduleService) GetNextTakings(userID uint, period time.Duration) ([]models.Schedule, error) {
	now := time.Now()
	end := now.Add(period)

	schedule, err := s.Repo.NextTakings(userID, now, end)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}
