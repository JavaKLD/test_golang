package main

import (
	"dolittle2/internal/controllers"
	"dolittle2/internal/database"
	"dolittle2/internal/repos"
	"dolittle2/internal/services"
	"dolittle2/migrations"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed %v", err)
	}

	err = migrations.Migration(db)
	if err != nil {
		log.Fatal("Ошибка миграции", err)
	}

	repo := repos.NewScheduleRepo(db)
	service := services.NewService(repo)
	controller := controllers.NewScheduleController(service)

	e := echo.New()

	e.POST("/schedule", controller.CreateSchedule)
	e.GET("/schedules", controller.GetUserSchedule)
	e.GET("/schedule", controller.GetSchedule)
	e.GET("/next_takings", controller.GetNextTakings)

	err = e.Start(":8080")

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
