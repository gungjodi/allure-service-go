package routers

import (
	"net/url"
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

// Get All projects godoc
// @Summary Get All projects
// @Description Get All projects
// @Tags Projects
// @Accept */*
// @Produce json
// @Success 200
// @Router /projects [get]
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

// Get Project By ID godoc
// @Summary Get Project By ID
// @Description Get Project By ID
// @Tags Projects
// @Accept */*
// @Produce json
// @Success 200
// @Param   project_id     path     string     true  "default"     default(default)
// @Router /projects/{project_id} [get]
func getProject(c *gin.Context) {
	var err error
	projectId := c.Param("project_id")
	redirect, redirectErr := strconv.ParseBool(c.Query("redirect"))
	if redirectErr != nil {
		redirect = false
	}

	currentProjectDir := utils.GetProjectPath(projectId)

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

// Get Project By ID godoc
// @Summary Get Project By ID
// @Description Get Project By ID
// @Tags Projects
// @Accept */*
// @Produce json
// @Success 200
// @Param   project_id     path     string     true  "default"     default(default)
// @Param   path     path     string     true  "default"     default(latest/widgets/summary.json)
// @Router /projects/{project_id}/reports/{path} [get]
func getReport(c *gin.Context) {
	projectId := c.Param("project_id")
	path, _ := url.QueryUnescape(c.Param("path"))
	currentProjectDir := utils.GetProjectPath(projectId)
	reportDir := filepath.Join(currentProjectDir, "reports", path)

	c.File(reportDir)
}
