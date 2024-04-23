package router

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

func Register(e *echo.Echo, logger *zap.Logger) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
}
