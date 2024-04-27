package server

import (
	"context"
	"github.com/Atvit/assessment-tax/config"
	"github.com/Atvit/assessment-tax/internals/tax"
	mw "github.com/Atvit/assessment-tax/middleware"
	"github.com/go-playground/validator/v10"
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
}

func New(e *echo.Echo, cfg *config.Configuration, logger *zap.Logger) Server {
	return &server{
		e:      e,
		cfg:    cfg,
		logger: logger,
	}
}

func (s server) registerRoutes() {
	validate := validator.New()
	e := s.e

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	admin := e.Group("/admin")
	admin.Use(middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		return mw.Authenticate(username, password, s.cfg)
	}))

	taxHandler := tax.NewHandler(s.logger, validate)
	tax := e.Group("/tax")
	tax.POST("/calculations", taxHandler.CalculateTax)
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
