package routers

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	docs "osp-allure/docs"
	"osp-allure/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GeneralRouters(router *gin.RouterGroup) {
	docs.SwaggerInfo.BasePath = os.Getenv("BASE_PATH")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/config", config)
}

// Get Config godoc
// @Summary Get App Config
// @Description Get app config
// @Tags General
// @Accept */*
// @Produce json
// @Success 200
// @Router /config [get]
func config(c *gin.Context) {
	dat, err := os.ReadFile(os.Getenv("ALLURE_VERSION"))

	if err != nil {
		log.Err(err)
	}

	c.JSON(http.StatusOK, map[string]any{
		"GO_VERSION":           runtime.Version(),
		"BASE_PATH":            fmt.Sprintf("http://%v%s", c.Request.Host, os.Getenv("BASE_PATH")),
		"PROJECTS_PATH":        utils.ProjectsPath(),
		"ALLURE_VERSION":       string(dat),
		"KEEP_RESULTS_HISTORY": os.Getenv("KEEP_RESULTS_HISTORY"),
		"KEEP_HISTORY_LATEST":  os.Getenv("KEEP_HISTORY_LATEST"),
	})
}
