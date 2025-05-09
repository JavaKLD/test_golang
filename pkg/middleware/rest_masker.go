package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func MaskSesitiveHTTP(c echo.Context) map[string]string {
	masked := make(map[string]string)

	for key, val := range c.QueryParams() {
		if strings.EqualFold(key, "user_id") {
			masked[key] = "***"
		} else if len(val) > 0 {
			masked[key] = val[0]
		}
	}
	return masked
}
