package configuration

import (
	"context"

	"restaurant_reservation/pkg/logger"

	"github.com/sethvargo/go-envconfig"
)

type Configuration struct {
	Environment string `env:"APP_ENV" required:"true"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"DEBUG"`
	Port        string `env:"PORT" envDefault:"8080"`
	Version     string `env:"VERSION" envDefault:"0.0.1"`
}

// NewConfiguration :
func NewConfiguration() *Configuration {
	cfg := new(Configuration)

	readConfigEnv(cfg)
	return cfg
}

func readConfigEnv(configuration *Configuration) {
	ctx := context.Background()
	if err := envconfig.Process(ctx, configuration); err != nil {
		logger.Get().FatalIfError("config", "readConfigEnv", err)
	}
}
