package server

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"

	"restaurant_reservation/internal/middleware"
	"restaurant_reservation/pkg/logger"
)

type Server struct {
	server *echo.Echo
	log    *logger.Log
}

func NewServer(container *dig.Container) *Server {
	var log *logger.Log

	err := container.Invoke(func(l *logger.Log) {
		log = l
	})
	if err != nil {
		log.Fatalf("Server", "NewServer", "Error getting logger: %v", err)
	}

	return &Server{
		server: echo.New(),
		log:    log,
	}
}

func (s *Server) InitServer(container *dig.Container) {
	s.setupMiddlewares()
	s.registerRoutes(container)
}

// middleware Setup
func (s *Server) setupMiddlewares() {
	s.server.Use(middleware.LoggerMiddleware(s.log))
	s.server.Use(echomiddleware.Recover())

	s.server.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

// RegisterRoutes registers the routes for the server.
func (s *Server) registerRoutes(container *dig.Container) {
	s.registerApplicationRoutes(container)
}
