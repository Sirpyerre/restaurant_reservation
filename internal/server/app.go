package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"restaurant_reservation/internal/dependencies"
)

func RunServer() {
	container := dependencies.GetContainer()
	server := NewServer(container)

	// Initialize the server
	server.InitServer(container)
	e := server.server

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(e, done)

	// Start the server
	if err := e.Start(":" + os.Getenv("PORT")); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	// Wait for the graceful shutdown to complete
	<-done
	server.log.Infof("Server shutdown complete")
}

func gracefulShutdown(apiServer *echo.Echo, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
