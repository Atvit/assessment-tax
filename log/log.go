package log

import (
	"go.uber.org/zap"
	"log"
)

func New() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	return logger
}
