package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func ProjectsPath() string {
	return filepath.Join(GetAppDataPath(), "projects")
}

func GetExistentsProjects(projectId string) bool {
	if _, err := os.Stat(GetProjectPath(projectId)); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetProjectPath(projectId string) string {
	return filepath.Join(ProjectsPath(), projectId)
}

func GetLatestProjectReport(projectId string) string {
	return filepath.Join(GetProjectPath(projectId), "reports", "latest")
}

func CreateProject(projectId string) error {
	if !GetExistentsProjects(projectId) {
		os.MkdirAll(GetProjectPath(projectId), 0755)
		os.MkdirAll(filepath.Join(GetProjectPath(projectId), "results"), 0755)
		os.MkdirAll(filepath.Join(GetProjectPath(projectId), "reports"), 0755)
	} else {
		return fmt.Errorf("project %s already exists", projectId)
	}
	return nil
}
