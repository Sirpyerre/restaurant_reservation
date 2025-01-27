package server

import (
	"go.uber.org/dig"

	"restaurant_reservation/internal/restaurantreservation/handlers"
	"restaurant_reservation/pkg/logger"
)

// registerApplicationRoutes registers the routes for the server.
func (s *Server) registerApplicationRoutes(container *dig.Container) {
	errInvoke := container.Invoke(func(
		helloHandler *handlers.HelloWorldHandler,
		healthHandler *handlers.HealthHandler,
	) {
		s.registerBaseRouters(helloHandler, healthHandler)
	})

	logger.Get().FatalIfError("server", "RegisterRoutes", errInvoke)
}

// registerBaseRouters registers the routes for the server.
func (s *Server) registerBaseRouters(helloHandler *handlers.HelloWorldHandler,
	healthHandler *handlers.HealthHandler,
) {
	// base routes
	s.server.Add("GET", "/hello", helloHandler.HelloWorldHandler)
	s.server.Add("GET", "/health", healthHandler.HealthHandler)
}
