package main

import (
	"os"
	"osp-allure/utils"
	"path/filepath"
	"testing"
	"time"

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

func TestStoreAllureReport(t *testing.T) {

	// projectId := "test123"

	// utils.CreateProject(projectId)

	// utils.StoreAllureReport(projectId, 1)

	// os.RemoveAll(utils.GetProjectPath(projectId))
}

func TestGenerateExecutorJSON(t *testing.T) {
	projectId := "test123"
	// create a project
	utils.CreateProject(projectId)
	buildOrder := utils.GetLatestProjectBuildOrder(projectId)
	utils.GenerateExecutorJson(projectId, buildOrder, "", "", "")
}

func TestGenerateAllureReport(t *testing.T) {
	projectId := "test1600"
	utils.CreateProject(projectId)
	utils.KeepResultHistory(projectId)
	latestBuildOrder := utils.GetLatestProjectBuildOrder(projectId)
	utils.StoreAllureReport(projectId, latestBuildOrder)

	newBuildOrder := latestBuildOrder + 1
	utils.GenerateExecutorJson(projectId, newBuildOrder, "", "", "")
	status := utils.GenerateReportCmd(projectId)
	assert.True(t, status.Complete)
}
