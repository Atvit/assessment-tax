package db

import (
	"database/sql"
	"github.com/Atvit/assessment-tax/config"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type DB interface {
	Connect() (*sql.DB, error)
}

type db struct {
	config config.Configuration
	logger *zap.Logger
}

func New(config config.Configuration, logger *zap.Logger) DB {
	return &db{
		config: config,
		logger: logger,
	}
}

func (d db) Connect() (*sql.DB, error) {
	sql, err := sql.Open("postgres", d.config.DatabaseURL)
	if err != nil {
		d.logger.Fatal("unable to configure database", zap.Error(err))
		return nil, err
	}
	err = sql.Ping()
	if err != nil {
		d.logger.Fatal("unable to connect database", zap.Error(err))
		return nil, err
	}

	return sql, nil
}
