package main

import (
	"dolittle2/database"
	"dolittle2/database/migrations"
	"dolittle2/internal/controllers"
	"dolittle2/internal/domain/repos"
	"dolittle2/internal/domain/services"
	"dolittle2/internal/server"
	"github.com/labstack/echo/v4"

	"dolittle2/internal/logger"
	"errors"
	"log/slog"
	"net/http"
)

func main() {
	logger.InitLogger()

	slog.Info("Сервер запущен")

	db, err := database.InitDB()
	if err != nil {
		slog.Error("Ошибка подключения к бд", "error: ", err)
	}

	err = migrations.Migration(db)
	if err != nil {
		slog.Error("Ошибка миграции", "error:", err)
	}

	repo := repos.NewScheduleRepo(db)
	service := services.NewService(repo)
	controller := server.NewScheduleController(service)

	go func() {
		e := echo.New()
		server.InitRoutes(e, controller)

		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start echo server", "error", err)
		}
	}()

	go controllers.StartGRPCServer(service)

	select {}
}
