package dependencies

import (
	"go.uber.org/dig"
	"restaurant_reservation/cmd/configuration"
	"restaurant_reservation/internal/database"
	"restaurant_reservation/internal/restaurantreservation/handlers"
	"restaurant_reservation/pkg/logger"
	"sync"
)

var (
	container *dig.Container
	once      sync.Once
)

// GetContainer :
func GetContainer() *dig.Container {
	once.Do(func() {
		container = buildContainer()
	})
	return container
}

func buildContainer() *dig.Container {
	c := dig.New()

	logger.Get().FatalIfError("container", "buildContainer",
		c.Provide(configuration.NewConfiguration),
		c.Provide(logger.NewLog),
		c.Provide(database.New),
		c.Provide(handlers.NewHealthHandler),
		c.Provide(handlers.NewHelloWorldHandler),
	)

	return c
}
