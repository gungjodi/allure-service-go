package utils

import (
	"os"
	"path/filepath"
)

func ProjectsPath() string {
	rootDir := os.Getenv("ROOT_DIR")
	return filepath.Join(rootDir, "projects")
}

func GetProjectPath(projectId string) (string, error) {
	projectPath := filepath.Join(ProjectsPath(), projectId)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return "", err
	}
	return projectPath, nil
}

func GetLatestProjectReport(projectId string) string {
	latestReportPath, _ := GetProjectPath(projectId)
	return filepath.Join(latestReportPath, "reports", "latest")
}
