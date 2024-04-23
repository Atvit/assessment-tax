package server

import (
	"context"
	"github.com/Atvit/assessment-tax/config"
	"github.com/Atvit/assessment-tax/router"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type Server interface {
	Start()
}

type server struct {
	echo   *echo.Echo
	cfg    *config.Configuration
	logger *zap.Logger
}

func New(echo *echo.Echo, cfg *config.Configuration, logger *zap.Logger) Server {
	return &server{
		echo:   echo,
		cfg:    cfg,
		logger: logger,
	}
}

func (s server) Start() {
	router.Register(s.echo, s.cfg)

	go func() {
		err := s.echo.Start(":" + strconv.Itoa(s.cfg.Port))
		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("unexpected shutdown the server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	gCtx := context.Background()
	ctx, cancel := context.WithTimeout(gCtx, 10*time.Second)
	defer cancel()

	s.logger.Info("shutting down the server")
	if err := s.echo.Shutdown(ctx); err != nil {
		s.logger.Fatal("unexpected shutdown the server", zap.Error(err))
	}
}
