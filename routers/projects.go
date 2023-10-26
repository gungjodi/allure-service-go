package routers

import (
	"net/http"
	"net/url"
	"os"
	"osp-allure/utils"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ProjectsRouters(router *gin.RouterGroup) {
	projectsRouter := router.Group("/projects")
	projectsRouter.GET("/", getAllProjects)
	projectsRouter.GET("/:project_id", getProject)
	projectsRouter.GET("/:project_id/reports/*path", getReport)
	projectsRouter.POST("/:project_id", createProject)
	projectsRouter.DELETE("/:project_id", deleteProject)
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
	projectId := c.Param("project_id")
	currentProjectDir := utils.GetProjectPath(projectId)

	var projectsLink []string
	listDir, _ := os.ReadDir(filepath.Join(currentProjectDir, "reports"))
	for _, dir := range listDir {
		projectsLink = append(projectsLink, path.Join(c.Request.URL.String(), "reports", dir.Name()))
	}

	c.JSON(200, gin.H{
		"projects": projectsLink,
	})
}

// Create Project godoc
// @Summary Create Project
// @Description Create Project
// @Tags Projects
// @Accept */*
// @Produce json
// @Success 200
// @Param   project_id     path     string     true  "default"     default(default)
// @Router /projects/{project_id} [post]
func createProject(c *gin.Context) {
	projectId := c.Param("project_id")

	isExists := utils.GetExistentsProjects(projectId)

	if isExists {
		c.JSON(200, gin.H{
			"message": "Project already exists",
		})
		return
	}

	if err := utils.CreateProject(projectId); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Project successfully created",
	})
}

// Delete Project godoc
// @Summary Delete Project
// @Description Delete Project
// @Tags Projects
// @Accept */*
// @Produce json
// @Success 200
// @Param   project_id     path     string     true  "default"     default(default)
// @Router /projects/{project_id} [delete]
func deleteProject(c *gin.Context) {
	projectId := c.Param("project_id")

	isExists := utils.GetExistentsProjects(projectId)

	if !isExists {
		c.JSON(404, gin.H{
			"message": "Project not found",
		})
		return
	}

	if err := utils.DeleteProject(projectId); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Project successfully deleted",
	})
}

// Get Path in a project directory By ID godoc
// @Summary 	Get Path in a project directory
// @Description Get Path in a project directory
// @Tags	Projects
// @Accept	*/*
// @Success 200
// @Param	project_id	path	string	true	"default"	default(default)
// @Param   path		path	string	true	"default"	default(latest/widgets/summary.json)
// @Router	/projects/{project_id}/reports/{path} [get]
func getReport(c *gin.Context) {
	projectId := c.Param("project_id")
	path, err := url.QueryUnescape(c.Param("path"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid path",
		})
		return
	}

	currentProjectDir := utils.GetProjectPath(projectId)
	reportDir := filepath.Join(currentProjectDir, "reports", path)

	c.File(reportDir)
}
