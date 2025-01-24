package middleware

import (
	"github.com/labstack/echo/v4"
	"restaurant_reservation/pkg/logger"
)

func LoggerMiddleware(log *logger.Log) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Infof("Request: %s %s", c.Request().Method, c.Request().URL.Path)
			return next(c)
		}
	}
}
