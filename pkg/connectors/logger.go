package connectors

import (
	"log/slog"
	"os"
)

func InitLogger() {
	logger := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(logger))
}
