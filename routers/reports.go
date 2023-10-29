package routers

import (
	"fmt"
	"os"
	"osp-allure/models/response"
	"osp-allure/utils"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func ReportsRouters(router fiber.Router) {
	router.Post("/send-results", sendResults)
	router.Get("/generate-report", generateReport)
	router.Get("/clean-results", cleanResults)
	router.Get("/clean-history", cleanHistory)

	reportRouter := router.Group("/report")
	reportRouter.Get("/backup/:project_id/:report_id", backupReport)
	reportRouter.Get("/download", downloadLatestReportCSV)
	reportRouter.Get("/copy", backupReport)
}

// @Summary Send results
// @Description Send allure result files to server
// @Tags Reports
// @Accept mpfd
// @Produce json
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /send-results [post]
// @Param project_id query string true "projectId" (default)
// @Param force_project_creation query boolean false "create project if not exists" (false)
// @Param files[] formData []file true "result files"
func sendResults(c *fiber.Ctx) error {
	time1 := time.Now()
	projectId := c.Query("project_id")
	projectResultsDir := utils.GetProjectResultsPath(projectId)
	projectExist := utils.GetExistentsProjects(projectId)

	if !projectExist {
		utils.CreateProject(projectId)
	}

	form, err := utils.MultipartForm(c)

	if err != nil {
		log.Error().Msgf("Error: %v", err)
	}

	files := form.File["files[]"]
	log.Info().Msgf("Received %v files", files[0].Header.Get("Content-Type"))

	if !utils.CheckFileOrDirExist(projectResultsDir) {
		if err := os.MkdirAll(projectResultsDir, os.ModePerm); err != nil {
			log.Error().Msgf("Error: %v", err)
		}
	}

	log.Info().Msgf("Saving results for project_id %s", projectId)

	// if only one file received and it's a tar file, extract it
	if len(files) == 1 && strings.HasSuffix(files[0].Filename, ".tar") {
		tarFileName := filepath.Join(projectResultsDir, files[0].Filename)
		if err := c.SaveFile(files[0], tarFileName); err != nil {
			log.Error().Msgf("Error: %v", err)

			return response.ResponseError(c, fiber.Map{
				"message": fmt.Sprintf("%v", err),
			})
		}
		utils.ExtractResultsArchive(tarFileName, projectResultsDir)
	} else {
		for _, file := range files {
			c.SaveFile(file, filepath.Join(projectResultsDir, file.Filename))
		}
	}

	time2 := time.Now()
	return response.ResponseSuccess(c, fiber.Map{
		"count":    len(files),
		"duration": time2.Sub(time1).Seconds(),
		"message":  fmt.Sprintf("Results successfully sent for project_id %s", projectId),
	})
}

// @Summary Generate report from sent results
// @Description
// @Tags Reports
// @Produce json
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /generate-report [get]
// @Param project_id query string true "projectId" (default)
// @Param execution_name query string 	false "executionName"
// @Param execution_from query string 	false "executionName"
// @Param execution_type query string 	false "executionName"
// @Param backup_latest	 query boolean 	false "executionName"	default(false)
func generateReport(c *fiber.Ctx) error {
	projectId := c.Query("project_id")
	executionName := c.Query("execution_name")
	executionFrom := c.Query("execution_from")
	executionType := c.Query("execution_type")
	backupLatest := c.QueryBool("backup_latest", false)

	projectExist := utils.GetExistentsProjects(projectId)
	if !projectExist {
		return response.ResponseNotFound(c, "Project not found")
	}

	latestBuildOrder := utils.GetLatestProjectBuildOrder(projectId)
	latestBackupReportId := utils.GetLatestProjectReportIdBackup(projectId)

	// if backup reportId is newer than latestBuildOrder, use it
	if latestBackupReportId > latestBuildOrder {
		latestBuildOrder = latestBackupReportId
		utils.CopyResultHistoryFromBackup(projectId, strconv.Itoa(latestBackupReportId))
	} else {
		utils.KeepResultHistory(projectId)
	}

	newBuildOrder := latestBuildOrder + 1
	utils.GenerateExecutorJson(projectId, newBuildOrder, executionName, executionFrom, executionType)

	if err := utils.GenerateReportCmd(projectId, newBuildOrder); err != nil {
		return response.ResponseError(c, fiber.Map{
			"message": err.Error(),
		})
	}

	if err := utils.CreateReportLatestSymlink(projectId, newBuildOrder); err != nil {
		return response.ResponseError(c, fiber.Map{
			"message": err.Error(),
		})
	}

	reportUrl := fmt.Sprintf("%s%s/projects/%s/reports/%d/index.html", c.BaseURL(), os.Getenv("BASE_PATH"), projectId, newBuildOrder)

	go utils.KeepReportHistory(projectId)

	if backupLatest {
		log.Info().Msgf("Async backup project_id %s, buildOrder %d", projectId, newBuildOrder)
		go utils.BackupReport(projectId, strconv.Itoa(newBuildOrder), true)
	}

	return response.ResponseSuccess(c, fiber.Map{
		"reportUrl":    reportUrl,
		"projectId":    projectId,
		"reportId":     newBuildOrder,
		"backupLatest": backupLatest,
		"message":      "Report successfully generated",
	})
}

// @Summary Clean results
// @Description Clean allure result files on server
// @Tags Reports
// @Accept mpfd
// @Produce json
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /clean-results [get]
// @Param project_id query string true "projectId" (default)
func cleanResults(c *fiber.Ctx) error {
	projectId := c.Query("project_id")
	projectExist := utils.GetExistentsProjects(projectId)
	if !projectExist {
		return response.ResponseNotFound(c, "Project not found")
	}

	if err := utils.CleanResults(projectId); err != nil {
		return response.ResponseError(c, fiber.Map{
			"message": err.Error(),
		})
	}

	return response.ResponseSuccess(c, fiber.Map{
		"message": "Results successfully cleaned",
	})
}

// @Summary Clean history
// @Description Clean history project
// @Tags Reports
// @Accept mpfd
// @Produce json
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /clean-history [get]
// @Param project_id query string true "projectId" (default)
func cleanHistory(c *fiber.Ctx) error {
	projectId := c.Query("project_id")
	projectExist := utils.GetExistentsProjects(projectId)
	if !projectExist {
		return response.ResponseNotFound(c, "Project not found")
	}

	if err := utils.CleanHistory(projectId); err != nil {
		return response.ResponseError(c, fiber.Map{
			"message": err.Error(),
		})
	}

	return response.ResponseSuccess(c, fiber.Map{
		"message": "History successfully cleaned",
	})
}

// Backup Report By ID godoc
// @Summary Backup Report By ID
// @Description Backup Report By ID
// @Tags Reports
// @Accept */*
// @Produce json
// @Success 200
// @Param   project_id		path      string     true  "default"	default(default)
// @Param   report_id		path      string     true  "default"	default(latest)
// @Param   should_delete	query     boolean    false  "default"	default(true)
// @Param   async			query     boolean    false  "default"	default(true)
// @Router /backup/{project_id}/{report_id} [get]
func backupReport(c *fiber.Ctx) error {
	projectId := c.Params("project_id")
	reportId := c.Params("report_id")
	shouldDelete := c.QueryBool("should_delete", true)
	async := c.QueryBool("async", true)

	if isExists := utils.GetExistentsProjects(projectId); !isExists {
		return response.ResponseNotFound(c, "Project not found")
	}

	if utils.GetLatestProjectBuildOrder(projectId) == 0 {
		return response.ResponseNotFound(c, "No report found")
	}

	if reportId == "latest" {
		reportId = utils.GetLatestProjectReportId(projectId)
	}

	if async {
		go func() {
			utils.BackupReport(projectId, reportId, shouldDelete)
		}()
	} else {
		if err := utils.BackupReport(projectId, reportId, shouldDelete); err != nil {
			return response.ResponseError(c, fiber.Map{
				"message": err.Error(),
			})
		}
	}

	return response.ResponseSuccess(c, fiber.Map{
		"message":      "Successfully backed up",
		"async":        async,
		"shouldDelete": shouldDelete,
		"projectId":    projectId,
		"reportId":     reportId,
	})
}

// Download latest report csv godoc
// @Summary Download latest report csv
// @Description Download latest report csv
// @Tags Reports
// @Accept */*
// @Produce application/octet-stream
// @Success 200
// @Param   project_id		query      string     true  "default"	default(default)
// @Router /report/download [get]
func downloadLatestReportCSV(c *fiber.Ctx) error {
	projectId := c.Query("project_id")

	if isExists := utils.GetExistentsProjects(projectId); !isExists {
		return response.ResponseNotFound(c, "Project not found")
	}
	csvPath := utils.StoreCSVLatestReport(projectId, "latest")

	return c.Download(csvPath)
}
