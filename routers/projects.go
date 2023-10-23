package routers

import (
	"os"
	"osp-allure/utils"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ProjectsRouters(router *gin.RouterGroup) {
	projectsRouter := router.Group("/projects")
	projectsRouter.GET("/", getAllProjects)
	projectsRouter.GET("/:project_id", getProject)
	projectsRouter.GET("/:project_id/reports/*path", getReport)
}

func getAllProjects(c *gin.Context) {
	projectsLink := []string{}

	listDir, _ := os.ReadDir(utils.ProjectsPath())

	if len(listDir) > 0 {
		for _, dir := range listDir {
			projectsLink = append(projectsLink, dir.Name())
		}
	}

	c.JSON(200, gin.H{
		"projects": projectsLink,
	})
}

func getProject(c *gin.Context) {
	var err error
	projectId := c.Param("project_id")
	redirect, redirectErr := strconv.ParseBool(c.Query("redirect"))
	if redirectErr != nil {
		redirect = false
	}

	currentProjectDir, err := utils.GetProjectPath(projectId)

	if err != nil {
		c.Error(err)
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
	currentProjectDir, err := utils.GetProjectPath(projectId)
	if err != nil {
		c.Error(err)
		return
	}

	reportDir := filepath.Join(currentProjectDir, "reports", path)

	c.File(reportDir)
}
