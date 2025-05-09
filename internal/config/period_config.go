package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig() time.Duration {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Ошибка загрузки env файла", slog.Any("error", err))
	}

	periodStr := os.Getenv("NEXT_TAKING_PERIOD")
	period, err := time.ParseDuration(periodStr)

	if err != nil {
		slog.Error("Ошибка прасинга периода", slog.Any("error", err))
	}

	return period
}
