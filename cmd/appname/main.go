package main

import (
	"dolittle2/internal/config"
	"dolittle2/internal/domain/repos"
	"dolittle2/internal/domain/services"
	"dolittle2/internal/server"
	"dolittle2/pkg/connectors"
	"github.com/labstack/echo/v4"

	"errors"
	"log/slog"
	"net/http"
)

func main() {
	connectors.InitLogger()

	slog.Info("Сервер запущен")

	db, err := connectors.InitDB()
	if err != nil {
		slog.Error("Ошибка подключения к бд", "error: ", err)
	}

	cfg := config.LoadConfig()
	repo := repos.NewScheduleRepo(db)
	service := services.NewService(repo, cfg)
	controller := server.NewScheduleController(service)

	go func() {
		e := echo.New()
		server.InitRoutes(e, controller)

		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start echo server", "error", err)
		}
	}()

	go server.StartGRPCServer(service)

	select {}
}
