package main

import (
	"dolittle2/internal/controllers"
	"dolittle2/internal/database"
	"dolittle2/internal/repos"
	"dolittle2/internal/routers"
	"dolittle2/internal/services"
	"dolittle2/migrations"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	slog.Info("Сервер запущен")
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

	go func() {
		e := echo.New()
		routers.InitRoutes(e, controller)

		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start echo server", "error", err)
		}
	}()

	go controllers.StartGRPCServer(service)

	select {}
}
