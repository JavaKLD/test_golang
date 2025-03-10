package repos

import (
	"dolittle2/internal/models"
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
	if err := r.db.Model(&models.Schedule{}).Where("user_id = ?", userID).Pluck("id",&scheduleID ).Error; err != nil {
		return nil, err
	}
	return scheduleID, nil
}

func (r * ScheduleRepo) FindSchedule(userID, scheduleID uint) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := r.db.Where("user_id = ? AND id = ?", userID,scheduleID).First(&schedule).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}
