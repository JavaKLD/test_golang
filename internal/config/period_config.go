package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig() time.Duration {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки env файла", err)
	}

	periodStr := os.Getenv("NEXT_TAKING_PERIOD")
	period, err := time.ParseDuration(periodStr)
	if err != nil {
		log.Fatal("Ошибка прасинга периода", err)
	}

	return period
}
