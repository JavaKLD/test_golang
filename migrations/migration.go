package migrations

import (
	"dolittle2/internal/models"
	"gorm.io/gorm"
	"log"
)

func Migration(db *gorm.DB) error {
	return db.AutoMigrate(&models.Schedule{})
}

func Rollback(db *gorm.DB) {
	err := db.Migrator().DropTable(&models.Schedule{})
	if err != nil {
		log.Fatal("Ошибка ", err)
	}
}
