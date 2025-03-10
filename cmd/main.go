package main

import (
	"dolittle2/internal/controllers"
	"dolittle2/internal/database"
	"dolittle2/internal/repos"
	"dolittle2/internal/services"
	"errors"
	"github.com/labstack/echo/v4"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed %v", err)
	}
	repo := repos.NewScheduleRepo(db)
	service := services.NewService(repo)
	controller := controllers.NewScheduleController(service)


	e := echo.New()

	e.POST("/schedule", controller.CreateSchedule)

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}


