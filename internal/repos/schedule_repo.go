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

func (r *ScheduleRepo) CreateSchedule(schedule *models.Schedule) (uint64, error) {
	err := r.db.Create(schedule).Error
	if err != nil {
		return 0, err
	}

	return schedule.ID, nil
}

func (r *ScheduleRepo) AidNameExists(aidName string, userID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Schedule{}).
		Where("aid_name = ? AND user_id = ?", aidName, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ScheduleRepo) FindByUserID(userID uint64) ([]uint64, error) {
	var scheduleID []uint64
	err := r.db.Model(&models.Schedule{}).Where("user_id = ?", userID).Pluck("id", &scheduleID).Error
	if err != nil {
		return nil, err
	}
	return scheduleID, nil
}

func (r *ScheduleRepo) FindSchedule(userID, scheduleID uint64) (*models.Schedule, error) {
	var schedule models.Schedule
	err := r.db.Where("user_id = ? AND id = ?", userID, scheduleID).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepo) NextTakings(userID uint64) ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.Where("user_id = ?", userID).Find(&schedules).Error
	if err != nil {
		return nil, err
	}

	return schedules, nil
}
