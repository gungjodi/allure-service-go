package routers

import (
	"fmt"
	"os"
	"runtime"

	docs "osp-allure/docs"
	"osp-allure/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog/log"
)

func GeneralRouters(router fiber.Router) {
	docs.SwaggerInfo.BasePath = os.Getenv("BASE_PATH")
	router.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("swagger/index.html")
	})
	router.Get("/swagger/*", swagger.HandlerDefault)
	router.Get("/config", config)
}

// Get Config godoc
// @Summary Get App Config
// @Description Get app config
// @Tags General
// @Accept */*
// @Produce json
// @Success 200
// @Router /config [get]
func config(c *fiber.Ctx) error {
	allureVersion, error := utils.GetAllureVersion()

	if error != nil {
		log.Err(error)
		allureVersion = error.Error()
	}

	return c.JSON(map[string]any{
		"APP_MODE":                             utils.GetEnv("APP_MODE", "release"),
		"GO_VERSION":                           runtime.Version(),
		"APP_VERSION":                          utils.GetEnv("APP_VERSION", ""),
		"BASE_PATH":                            fmt.Sprintf("%s%s", c.BaseURL(), os.Getenv("BASE_PATH")),
		"PROJECTS_PATH":                        utils.ProjectsPath(),
		"PROJECTS_BACKUP_PATH":                 utils.BackupProjectsPath(),
		"ALLURE_VERSION":                       allureVersion,
		"KEEP_RESULTS_HISTORY":                 os.Getenv("KEEP_RESULTS_HISTORY"),
		"KEEP_HISTORY_LATEST":                  os.Getenv("KEEP_HISTORY_LATEST"),
		"DOWNLOAD_REPORT_CSV_DESTINATION_PATH": utils.GetBackupReportCSVPath(),
		"GOROUTINE_COUNT":                      runtime.NumGoroutine(),
	})
}
