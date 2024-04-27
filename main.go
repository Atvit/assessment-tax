package main

import (
	"github.com/Atvit/assessment-tax/config"
	"github.com/Atvit/assessment-tax/db"
	"github.com/Atvit/assessment-tax/log"
	"github.com/Atvit/assessment-tax/server"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	logger := log.New()
	cfg := config.New(logger)
	sv := server.New(e, cfg, logger)
	db := db.New(cfg, logger)

	sv.Start()
	conn, err := db.Connect()
	if err != nil {
		panic(err)
	}

	_ = conn
}
