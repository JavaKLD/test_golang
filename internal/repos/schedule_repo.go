package repos

import (
	"dolittle2/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type ScheduleRepo struct {
	db *gorm.DB
}

func NewScheduleRepo(db *gorm.DB) *ScheduleRepo {
	return &ScheduleRepo{db: db}
}

func (r *ScheduleRepo) CreateSchedule(schedule *models.Schedule) (uint, error) {
	if err := r.db.Create(schedule).Error; err != nil {
		return 0, err
	}
	return schedule.ID, nil
}

func (r *ScheduleRepo) FindByUserID(userID uint) ([]uint, error) {
	var scheduleID []uint
	err := r.db.Model(&models.Schedule{}).Where("user_id = ?", userID).Pluck("id", &scheduleID).Error
	if err != nil {
		return nil, err
	}
	return scheduleID, nil
}

func (r *ScheduleRepo) FindSchedule(userID, scheduleID uint) (*models.Schedule, error) {
	var schedule models.Schedule
	err := r.db.Where("user_id = ? AND id = ?", userID, scheduleID).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepo) NextTakings(userID uint) ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.Where("user_id = ?", userID).Find(&schedules).Error
	if err != nil {
		return nil, err
	}

	fmt.Println("Schedules", schedules)
	return schedules, nil
}
