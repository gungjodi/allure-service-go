package main

import (
	"os"
	"osp-allure/utils"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	currentPath, _ := os.Getwd()
	err := godotenv.Load(filepath.Join(currentPath, ".env"))

	if err != nil {
		log.Fatal().Any("Error loading .env file", err)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
}

var projectId = "test123"

func TestCreateProject(t *testing.T) {
	if isExist := utils.GetExistentsProjects(projectId); isExist {
		utils.DeleteProject(projectId)
	}
	errCreate := utils.CreateProject(projectId)
	isExist := utils.GetExistentsProjects(projectId)
	projectPath := utils.GetProjectPath(projectId)
	assert.Nil(t, errCreate)
	assert.True(t, isExist)
	assert.NotEqual(t, "", projectPath)

	t.Cleanup(func() {
		utils.DeleteProject(projectId)
	})
}

func TestGenerateExecutorJSON(t *testing.T) {
	// create a project
	utils.CreateProject(projectId)
	latestBuildOrder := utils.GetLatestProjectBuildOrder(projectId)
	newBuildOrder := latestBuildOrder + 1
	err := utils.GenerateExecutorJson(projectId, newBuildOrder, "", "", "")
	assert.Nil(t, err)

	t.Cleanup(func() {
		utils.DeleteProject(projectId)
	})
}

func TestGenerateAllureReport(t *testing.T) {
	utils.CreateProject(projectId)
	utils.KeepResultHistory(projectId)
	latestBuildOrder := utils.GetLatestProjectBuildOrder(projectId)
	utils.KeepReportHistory(projectId)

	newBuildOrder := latestBuildOrder + 1
	utils.GenerateExecutorJson(projectId, newBuildOrder, "", "", "")
	status := utils.GenerateReportCmd(projectId, newBuildOrder)
	assert.Nil(t, status)
	err := utils.CreateReportLatestSymlink(projectId, latestBuildOrder)
	assert.Nil(t, err)

	// t.Cleanup(func() {
	// 	utils.DeleteProject(projectId)
	// })
}

func TestBackupProject(t *testing.T) {
	utils.CreateProject(projectId)
	utils.GenerateExecutorJson(projectId, 1, "a", "b", "c")
	utils.GenerateReportCmd(projectId, 1)
	utils.CreateReportLatestSymlink(projectId, 1)

	err := utils.BackupReport(projectId, "1", false)
	assert.Nil(t, err)

	t.Cleanup(func() {
		utils.DeleteProject(projectId)
	})
}

func TestBackupLatestReport(t *testing.T) {
	utils.CreateProject(projectId)
	latestBuildOrder := utils.GetLatestProjectBuildOrder(projectId)
	utils.GenerateExecutorJson(projectId, latestBuildOrder+1, "a", "b", "c")
	utils.GenerateReportCmd(projectId, latestBuildOrder+1)
	utils.CreateReportLatestSymlink(projectId, latestBuildOrder+1)

	latestReportId := "latest"
	if latestReportId == "latest" {
		latestReportId = utils.GetLatestProjectReportId(projectId)
	}

	assert.Equal(t, "1", latestReportId)

	latestBuildOrder = utils.GetLatestProjectBuildOrder(projectId)
	utils.GenerateExecutorJson(projectId, latestBuildOrder+1, "a", "b", "c")
	utils.GenerateReportCmd(projectId, latestBuildOrder+1)
	utils.CreateReportLatestSymlink(projectId, latestBuildOrder+1)

	latestReportId = "latest"
	if latestReportId == "latest" {
		latestReportId = utils.GetLatestProjectReportId(projectId)
	}
	assert.Equal(t, "2", latestReportId)

	utils.BackupReport(projectId, latestReportId, false)
	dir, _ := os.ReadDir(filepath.Join(utils.GetProjectReportsBackupPath(projectId)))
	assert.Equal(t, 2, len(dir))

	for _, d := range dir {
		assert.True(t, d.Name() == "latest" || d.Name() == "2")
	}

	t.Cleanup(func() {
		utils.DeleteProject(projectId)
		<-cmd.NewCmd("rm", "-rf", utils.GetProjectBackupPath(projectId)).Start()
	})
}
