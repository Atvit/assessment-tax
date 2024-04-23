package main

import (
	"github.com/Atvit/assessment-tax/config"
	"github.com/Atvit/assessment-tax/log"
	"github.com/Atvit/assessment-tax/server"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	logger := log.New()
	cfg := config.New()
	server := server.New(e, cfg, logger)

	server.Start()
}
