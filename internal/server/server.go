package server

import (
	"restaurant_reservation/pkg/logger"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
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
