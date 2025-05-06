package migrations

import (
	"dolittle2/internal/domain/models"
	"gorm.io/gorm"
)

func Migration(db *gorm.DB) error {
	return db.AutoMigrate(&models.Schedule{})
}
