package server

import (
	"context"
	"github.com/Atvit/assessment-tax/config"
	"github.com/Atvit/assessment-tax/internals/setting"
	"github.com/Atvit/assessment-tax/internals/tax"
	mw "github.com/Atvit/assessment-tax/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e      *echo.Echo
	cfg    *config.Configuration
	logger *zap.Logger

	settingHandler setting.Handler
	taxHandler     tax.Handler
}

func New(
	e *echo.Echo,
	cfg *config.Configuration,
	logger *zap.Logger,

	taxHandler tax.Handler,
	settingHandler setting.Handler,
) Server {
	return &server{
		e:      e,
		cfg:    cfg,
		logger: logger,

		taxHandler:     taxHandler,
		settingHandler: settingHandler,
	}
}

func (s server) registerRoutes() {
	e := s.e

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	admin := e.Group("/admin")
	admin.Use(middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		return mw.Authenticate(username, password, s.cfg)
	}))
	admin.POST("/deductions/personal", s.settingHandler.UpdatePersonalDeduction)

	tax := e.Group("/tax")
	tax.POST("/calculations", s.taxHandler.CalculateTax)
}

func (s server) Start() {
	s.registerRoutes()

	go func() {
		err := s.e.Start(":" + strconv.Itoa(s.cfg.Port))
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
	if err := s.e.Shutdown(ctx); err != nil {
		s.logger.Fatal("unexpected shutdown the server", zap.Error(err))
	}
}
