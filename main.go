package main

import (
	"os"
	"strings"
	"syscall"

	app_middlwares "osp-allure/middlewares"
	"osp-allure/routers"
	"osp-allure/utils"
	"time"

	_ "osp-allure/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var appMode string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Error().AnErr("Error loading .env file", err)
	}

	appMode = utils.GetEnv("APP_MODE", "release")

	logLevel := zerolog.InfoLevel
	if strings.ToLower(appMode) == "debug" {
		logLevel = zerolog.DebugLevel
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger().Level(logLevel)
	log.Debug().Msgf("App mode: %s", appMode)
	syscall.Umask(0)
}

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit:                    512 << 20,
		StreamRequestBody:            true,
		DisablePreParseMultipartForm: true,
		UnescapePath:                 true,
	})
	app.Use(app_middlwares.AddTrailingSlashes())
	app.Use(app_middlwares.FavIcon())
	app.Use(cors.New())
	if appMode != "release" {
		app.Use(logger.New())
	}
	app.Use(recover.New())

	basePath := utils.GetEnv("BASE_PATH", "/")
	app.Route(basePath, func(router fiber.Router) {
		bindAllRouters(router)
	})

	if err := app.Listen(utils.GetEnv("HOST", "0.0.0.0") + ":" + utils.GetEnv("PORT", "5050")); err != nil {
		log.Fatal().Msgf("%v", err)
	}
}

func bindAllRouters(group fiber.Router) {
	routers.GeneralRouters(group)
	routers.ProjectsRouters(group)
	routers.ReportsRouters(group)
}
