package server

import (
	"go.uber.org/dig"
	"net/http"
	"restaurant_reservation/internal/restaurantreservation/handlers"
	"restaurant_reservation/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes(container *dig.Container, e *echo.Echo) http.Handler {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	errInvoke := container.Invoke(func(
		helloHandler *handlers.HelloWorldHandler,
		healthHandler *handlers.HealthHandler,
	) {
		e.Add("GET", "/hello", helloHandler.HelloWorldHandler)
		e.Add("GET", "/health", healthHandler.HealthHandler)
	})

	logger.Get().FatalIfError("server", "RegisterRoutes", errInvoke)

	return e
}
