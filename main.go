package main

import (
	"net/http"
	"os"
	reports "osp-allure/routers"
	"time"

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
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	app := gin.Default()
	app.MaxMultipartMemory = 512 << 20
	router := app.Group(os.Getenv("BASE_PATH"))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{
			"message": "Hello World",
		})
	})

	reports.BindProjectsRouters(router.Group("/projects"))
	reports.BindReportsRouters(router)

	app.Run("0.0.0.0:5050")
}
