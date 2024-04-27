package main

import (
	"github.com/Atvit/assessment-tax/config"
	"github.com/Atvit/assessment-tax/db"
	"github.com/Atvit/assessment-tax/internals/setting"
	"github.com/Atvit/assessment-tax/internals/tax"
	"github.com/Atvit/assessment-tax/log"
	"github.com/Atvit/assessment-tax/server"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	logger := log.New()
	cfg := config.New(logger)
	validate := validator.New()
	db := db.New(cfg, logger)

	conn, err := db.Connect()
	if err != nil {
		panic(err)
	}

	settingRepo := setting.NewRepository(conn)
	settingHandler := setting.NewHandler(logger, validate, settingRepo)

	taxHandler := tax.NewHandler(logger, validate, settingRepo)

	sv := server.New(e, cfg, logger, taxHandler, settingHandler)
	sv.Start()
}
