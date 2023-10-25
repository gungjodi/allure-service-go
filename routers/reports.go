package routers

import (
	"fmt"
	"net/http"
	"os"
	"osp-allure/utils"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ReportsRouters(router *gin.RouterGroup) {
	router.POST("/send-results", sendResults)
	router.POST("/generate-report", generateReport)
}

// @Summary Send results
// @Description Send allure result files to server
// @Produce json
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /send-results [post]
// @Param project_id query string true "projectId" (default)
// @Param force_project_creation query boolean false "create project if not exists" (false)
// @Param files formData []file true "result files"
func sendResults(c *gin.Context) {
	time1 := time.Now()
	projectId := c.Query("project_id")
	projectDir := utils.GetProjectPath(projectId)
	projectExist := utils.GetExistentsProjects(projectId)
	form, err := c.MultipartForm()

	if !projectExist {
		utils.CreateProject(projectId)
	}

	if err != nil {
		log.Error().Msgf("Error: %v", err)
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}

	files := form.File["files"]
	var size int64 = 0
	for _, file := range files {
		size = size + file.Size
		c.SaveUploadedFile(file, filepath.Join(projectDir, "results", file.Filename))
	}

	time2 := time.Now()
	c.JSON(http.StatusOK, map[string]any{
		"count":    len(files),
		"size":     size,
		"duration": time2.Sub(time1).Seconds(),
		"message":  fmt.Sprintf("Results successfully sent for project_id %s", projectId),
	})
}

// @Summary <summary>
// @Description <API Description>
// @Produce json
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /generate-report [get]
// @Param project_id query string true "projectId" (default)
// @Param execution_name query string false "executionName"
// @Param execution_from query string false "executionName"
// @Param execution_type query string false "executionName"
func generateReport(c *gin.Context) {
	projectId := c.Query("project_id")
	executionName := c.Query("execution_name")
	executionFrom := c.Query("execution_from")
	executionType := c.Query("execution_type")
	projectExist := utils.GetExistentsProjects(projectId)
	if !projectExist {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Project not found",
		})
		return
	}

	utils.KeepResultHistory(projectId)
	latestBuildOrder := utils.GetLatestProjectBuildOrder(projectId)
	utils.StoreAllureReport(projectId, latestBuildOrder)

	newBuildOrder := latestBuildOrder + 1
	utils.GenerateExecutorJson(projectId, newBuildOrder, executionName, executionFrom, executionType)

	utils.GenerateReportCmd(projectId)

	reportUrl := fmt.Sprintf("http://%s%s/projects/%s/reports/latest/index.html", c.Request.Host, os.Getenv("BASE_PATH"), projectId)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"reportUrl": reportUrl,
			"projectId": projectId,
		},
		"message": "Report successfully generated",
	})
}
