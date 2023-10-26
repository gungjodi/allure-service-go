package main

import (
	"os"
	"osp-allure/routers"
	"syscall"
	"time"

	_ "osp-allure/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Error().AnErr("Error loading .env file", err)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	syscall.Umask(0)
}

func main() {
	gin.SetMode(os.Getenv("GIN_MODE"))
	app := gin.Default()
	app.RedirectTrailingSlash = false
	app.UseRawPath = true
	app.UnescapePathValues = false
	app.RedirectFixedPath = false
	app.RemoveExtraSlash = true
	app.MaxMultipartMemory = 512 << 20

	base := app.Group(os.Getenv("BASE_PATH"))
	BindAllRouters(base)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5050"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	log.Fatal().Err(app.Run(host + ":" + port))
}

func BindAllRouters(group *gin.RouterGroup) {
	routers.GeneralRouters(group)
	routers.ProjectsRouters(group)
	routers.ReportsRouters(group)
}
