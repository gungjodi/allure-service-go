package routers

import (
	"fmt"
	"net/url"
	"os"
	"osp-allure/models"
	"osp-allure/models/response"
	"osp-allure/utils"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func ProjectsRouters(router fiber.Router) {
	projectsRouter := router.Group("/projects")
	projectsRouter.Get("/", getAllProjects)
	projectsRouter.Post("/", createProject)
	projectsRouter.Post("/batch-delete", batchDeleteProject)
	projectsRouter.Get("/:project_id", getProject)
	projectsRouter.Get("/:project_id/reports/*", getReport).Name("project_reports_endpoint")
	projectsRouter.Delete("/:project_id", deleteProject)
}

// Get All projects godoc
// @Summary Get All projects
// @Description Get All projects
// @Tags Projects
// @Accept */*
// @Produce json
// @Success 200
// @Router /projects [get]
func getAllProjects(c *fiber.Ctx) error {
	projectsLink := []string{}

	listDir, _ := os.ReadDir(utils.ProjectsPath())

	if len(listDir) > 0 {
		for _, dir := range listDir {
			projectsLink = append(projectsLink, dir.Name())
		}
	}

	return response.ResponseSuccess(c, fiber.Map{
		"projects": projectsLink,
	})
}

// Get Project By ID godoc
// @Summary Get Project By ID
// @Description Get Project By ID
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200
// @Param   project_id     path     string     true  "default"     default(default)
// @Router /projects/{project_id} [get]
func getProject(c *fiber.Ctx) error {
	projectId := c.Params("project_id")

	if isExists := utils.GetExistentsProjects(projectId); !isExists {
		return response.ResponseNotFound(c, "Project not found")
	}

	currentProjectReportsDir := utils.GetProjectReportsPath(projectId)

	var reportsLink []string
	var size int64 = 0

	listDir, _ := os.ReadDir(currentProjectReportsDir)
	for _, dir := range listDir {
		info, _ := dir.Info()
		size = size + info.Size()
		reportsLink = append(reportsLink, filepath.Join(c.BaseURL(), utils.GetFullReportUrl(c, projectId, dir.Name())))
	}

	return response.ResponseSuccess(c, fiber.Map{
		"reports":   reportsLink,
		"totalSize": fmt.Sprintf("%v%s", (size / 1024 / 1024), "MB"),
	})
}

// Create Project godoc
// @Summary Create Project
// @Description Create Project
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200
// @Param   request     body     models.CreateProjectRequest     true  "default"     (default)
// @Router /projects [post]
func createProject(c *fiber.Ctx) error {
	var createProjectRequest models.CreateProjectRequest

	if err := c.BodyParser(&createProjectRequest); err != nil {
		log.Error().Msgf("Error parsing body: %s", err)
		return response.ResponseBadRequest(c, err.Error())
	}

	projectId := createProjectRequest.ID

	isExists := utils.GetExistentsProjects(projectId)

	if isExists {
		return c.JSON(fiber.Map{
			"message": "Project already exists",
		})
	}

	if err := utils.CreateProject(projectId); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Project successfully created",
	})
}

// Batch Delete Project godoc
// @Summary Batch Delete Project
// @Description Batch Delete Project
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200
// @Param   request     body     models.BatchDeleteRequest     true  "default"     default(default)
// @Router /projects/batch-delete [post]
func batchDeleteProject(c *fiber.Ctx) error {
	delete_success := []string{}
	delete_failed := []string{}

	var batchDeleteRequest models.BatchDeleteRequest

	if err := c.BodyParser(&batchDeleteRequest); err != nil {
		log.Error().Msgf("Error parsing body: %s", err)
		return response.ResponseBadRequest(c, err.Error())
	}

	if batchDeleteRequest.Async {
		log.Info().Msgf("Deleting %d projects asynchronously", len(batchDeleteRequest.ProjectIds))
		go func() {
			for _, projectId := range batchDeleteRequest.ProjectIds {
				utils.StoreCSVLatestReport(projectId, "latest")
				utils.DeleteProject(projectId)
			}
		}()
		return response.ResponseSuccess(c, fiber.Map{
			"message": "Projects successfully deleted",
			"total":   len(batchDeleteRequest.ProjectIds),
		})
	} else {
		for _, projectId := range batchDeleteRequest.ProjectIds {
			utils.StoreCSVLatestReport(projectId, "latest")
			deleted := utils.DeleteProject(projectId)
			if deleted != "" {
				delete_success = append(delete_success, deleted)
			} else {
				delete_failed = append(delete_failed, projectId)
			}
		}

		return response.ResponseSuccess(c, fiber.Map{
			"message":            "Projects successfully deleted",
			"deletedProjects":    delete_success,
			"notDeletedProjects": delete_failed,
		})
	}
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
func deleteProject(c *fiber.Ctx) error {
	projectId := c.Params("project_id")

	isExists := utils.GetExistentsProjects(projectId)

	if !isExists {
		return response.ResponseNotFound(c, "Project not found")
	}

	deletedProject := utils.DeleteProject(projectId)

	return response.ResponseSuccess(c, fiber.Map{
		"message":        "Project successfully deleted",
		"deletedProject": deletedProject,
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
func getReport(c *fiber.Ctx) error {
	projectId := c.Params("project_id")
	pathParam, err := url.PathUnescape(c.Params("*"))

	if err != nil {
		log.Error().Msgf("Error parsing path: %s", err)
		return response.ResponseBadRequest(c, err.Error())
	}

	// Check if the file is in the current directory
	currentProjectReportsDir := utils.GetProjectReportsPath(projectId)
	currentPath := filepath.Join(currentProjectReportsDir, pathParam)
	if utils.CheckFileOrDirExist(currentPath) {
		info, _ := os.Stat(currentPath)
		if info.IsDir() {
			currentPath = filepath.Join(currentPath, "index.html")
		}
		return c.SendFile(currentPath)
	}

	// Check if the file is in the backup directory
	currentProjectReportsBackupDir := utils.GetProjectReportsBackupPath(projectId)
	currentBackupPath := filepath.Join(currentProjectReportsBackupDir, pathParam)
	if utils.CheckFileOrDirExist(currentBackupPath) {
		info, _ := os.Stat(currentBackupPath)
		if info.IsDir() {
			currentBackupPath = filepath.Join(currentBackupPath, "index.html")
		}

		return c.SendFile(currentBackupPath)
	}

	return response.ResponseNotFound(c, "Report not found")
}
