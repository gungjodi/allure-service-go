package main

import (
	"os"
	"osp-allure/routers"
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
}

// @title Gin Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	gin.SetMode(os.Getenv("GIN_MODE"))
	app := gin.Default()
	app.MaxMultipartMemory = 512 << 20

	base := app.Group("/")
	{
		BindAllRouters(base)
	}

	sub := app.Group(os.Getenv("BASE_PATH"))
	{
		BindAllRouters(sub)
	}

	app.Run("0.0.0.0:5050")
}

func BindAllRouters(group *gin.RouterGroup) {
	routers.GeneralRouters(group)
	routers.ProjectsRouters(group)
	routers.ReportsRouters(group)
}
