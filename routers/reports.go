package routers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-cmd/cmd"
	"github.com/rs/zerolog/log"
)

var projectDir = ".projects"

func BindReportsRouters(router *gin.RouterGroup) {
	router.POST("/send-results", sendResults)
	router.GET("/generate-report", generateReport)
}

func sendResults(c *gin.Context) {
	projectId := c.Query("project_id")
	forceProjectCreation := c.Query("force_project_creation")
	log.Info().Msgf("Project ID: %s, force %s", projectId, forceProjectCreation)
	time1 := time.Now()

	// Multipart form
	form, err := c.MultipartForm()

	if err != nil {
		log.Error().Msgf("Error: %v", err)
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": fmt.Sprintf("%v", err),
		})
	}

	files := form.File["files"]
	var size int64 = 0
	for _, file := range files {
		size = size + file.Size
		c.SaveUploadedFile(file, fmt.Sprintf("%s/%s/results/%s", projectDir, projectId, file.Filename))
	}
	time2 := time.Now()
	diff := time2.Sub(time1)

	c.JSON(http.StatusOK, map[string]any{
		"count":    len(files),
		"size":     size,
		"duration": diff.Seconds(),
	})
}

func generateReport(c *gin.Context) {
	projectId := c.Query("project_id")
	log.Info().Msgf("Generate report for project %s: START", projectId)
	status := generateReportCmd(projectId)
	c.JSON(http.StatusOK, map[string]any{
		"status": status.Stdout,
	})
}

func generateReportCmd(projectId string) cmd.Status {
	// Start a long-running process, capture stdout and stderr
	generateAllureCmd := cmd.NewCmd(".data/allure/bin/allure", "generate", "--clean", "-o", fmt.Sprintf("%s/%s/reports/latest", projectDir, projectId), fmt.Sprintf("%s/%s/results", projectDir, projectId))
	statusChan := <-generateAllureCmd.Start()
	log.Info().Msgf("Generate report for project DONE: %v", statusChan)
	return statusChan
}
