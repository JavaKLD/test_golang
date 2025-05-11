package tests

import (
	"dolittle2/internal/domain/repos"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log/slog"
)

func TestRepo() *repos.ScheduleRepo {
	dsn := "root:strong_password@tcp(mysql:3306)/app_db?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect db", slog.Any("error", err))
	}

	return repos.NewScheduleRepo(db)
}
