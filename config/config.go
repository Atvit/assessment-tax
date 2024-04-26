package config

import (
	"github.com/caarlos0/env"
	"go.uber.org/zap"
)

type Configuration struct {
	Port          int    `env:"PORT" envDefault:"8080"`
	DatabaseURL   string `env:"DATABASE_URL" envDefault:"host=localhost port=5432 user=postgres password=postgres dbname=ktaxes sslmode=disable"`
	AdminUsername string `env:"ADMIN_USERNAME" envDefault:"adminTax"`
	AdminPassword string `env:"ADMIN_PASSWORD" envDefault:"admin!"`
}

func New(logger *zap.Logger) *Configuration {
	cfg := Configuration{}
	if err := env.Parse(&cfg); err != nil {
		logger.Error("parsing config error", zap.Error(err))
	}

	return &cfg
}
