package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ctxKeyTraceID struct{}

func ReqLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		traceID := c.Request().Header.Get("X-TRACE-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		ctx := context.WithValue(c.Request().Context(), ctxKeyTraceID{}, traceID)
		c.SetRequest(c.Request().WithContext(ctx))

		start := time.Now()
		err := next(c)

		slog.Info("incoming request",
			"traceID", traceID,
			"method", c.Request().Method,
			"url", c.Request().URL.String(),
			"headers", c.Request().Header,
			"params", MaskSesitiveHTTP(c),
			"received_at", time.Now().Format(time.RFC3339),
			"ip_address", c.RealIP(),
		)

		slog.Info("outgoing response",
			"traceID", traceID,
			"status", c.Response().Status,
			"duration", time.Since(start).String(),
			"size", c.Response().Size,
			"received_at", time.Now().Format(time.RFC3339),
		)

		return err
	}
}
