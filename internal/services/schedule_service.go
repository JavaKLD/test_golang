package services

import (
	"dolittle2/internal/models"
	"dolittle2/internal/repos"
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

func (s *ScheduleService) FindSchedule(userID, scheduleID uint) (*models.Schedule, error) {
	return s.Repo.FindSchedule(userID, scheduleID)
}
