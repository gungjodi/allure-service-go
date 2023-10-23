package routers

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BindProjectsRouters(router *gin.RouterGroup) {
	router.GET("/:project_id", getProject)
	router.GET("/:project_id/reports/*path", getReport)
}

func getProject(c *gin.Context) {
	projectId := c.Param("project_id")
	redirect, redirectErr := strconv.ParseBool(c.Query("redirect"))
	if redirectErr != nil {
		redirect = false
	}

	currentProjectDir := filepath.Join(projectDir, projectId)
	latestProjectReportsDir := filepath.Join(currentProjectDir, "reports", "latest")

	if _, err := os.Stat(latestProjectReportsDir); os.IsNotExist(err) {
		c.JSON(404, gin.H{
			"error": fmt.Sprintf("Project %s not found", projectId),
		})
		return
	}

	var projectsLink []string
	listDir, _ := os.ReadDir(filepath.Join(currentProjectDir, "reports"))
	for _, dir := range listDir {
		projectsLink = append(projectsLink, path.Join(c.Request.URL.String(), "reports", dir.Name()))
	}

	if redirect {
		return
	}

	c.JSON(200, gin.H{
		"projects": projectsLink,
	})
}

func getReport(c *gin.Context) {
	projectId := c.Param("project_id")
	path := c.Param("path")
	currentProjectDir := filepath.Join(projectDir, projectId)
	reportDir := filepath.Join(currentProjectDir, "reports", path)
	reportPath := filepath.Join(projectDir, projectId, "reports", path)
	if _, err := os.Stat(reportDir); os.IsNotExist(err) {
		c.JSON(404, gin.H{
			"error": fmt.Sprintf("Report not found for projectId %s", projectId),
		})
		return
	}

	c.File(reportPath)
}
