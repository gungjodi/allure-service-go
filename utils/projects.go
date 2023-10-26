package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"syscall"

	"github.com/go-cmd/cmd"
	"github.com/rs/zerolog/log"
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

func GetLatestProjectReportId(projectId string) string {
	var latestReport string
	sortedReportDir := GetSortedReportsDir(projectId)
	if len(sortedReportDir) > 0 {
		sort.Sort(sort.Reverse(sort.IntSlice(sortedReportDir)))
		latestReport = strconv.Itoa(sortedReportDir[0])
	}

	if latestReport == "" {
		latestReport = "1"
	}
	return latestReport
}

func GetLatestProjectReport(projectId string) string {
	return filepath.Join(GetProjectPath(projectId), "reports", GetLatestProjectReportId(projectId))
}

func CreateProject(projectId string) error {
	if !GetExistentsProjects(projectId) {
		syscall.Umask(0)
		if err := os.MkdirAll(GetProjectPath(projectId), os.ModeSticky|os.ModePerm); err != nil {
			log.Err(err)
			return err
		}
		if err := os.MkdirAll(filepath.Join(GetProjectPath(projectId), "results"), os.ModeSticky|os.ModePerm); err != nil {
			log.Err(err)
			return err
		}
		if err := os.MkdirAll(filepath.Join(GetProjectPath(projectId), "reports"), os.ModeSticky|os.ModePerm); err != nil {
			log.Err(err)
			return err
		}
	} else {
		log.Error().Msgf("project %s already exists", projectId)
		return fmt.Errorf("project %s already exists", projectId)
	}

	return nil
}

func DeleteProject(projectId string) error {
	if GetExistentsProjects(projectId) {
		if err := os.RemoveAll(GetProjectPath(projectId)); err != nil {
			log.Err(err)
			return err
		}
	} else {
		return fmt.Errorf("project %s not exists", projectId)
	}
	return nil
}

func CreateReportLatestSymlink(projectId string, latestBuildOrder int) error {
	latestProjectPath := GetLatestProjectReport(projectId)
	latestReportSymlink := filepath.Join(GetProjectPath(projectId), "reports", "latest")

	log.Debug().Msgf("%s %s %s %s", "ln", "-sf", filepath.Join(GetAppDataPath(), "resources", "app.js"), filepath.Join(latestProjectPath, "app.js"))
	log.Debug().Msgf("%s %s %s %s", "ln", "-sf", filepath.Join(GetAppDataPath(), "resources", "styles.css"), filepath.Join(latestProjectPath, "styles.css"))
	status1 := <-cmd.NewCmd("ln", "-sf", filepath.Join(GetAppDataPath(), "resources", "app.js"), filepath.Join(latestProjectPath, "app.js")).Start()
	status2 := <-cmd.NewCmd("ln", "-sf", filepath.Join(GetAppDataPath(), "resources", "styles.css"), filepath.Join(latestProjectPath, "styles.css")).Start()

	if len(status1.Stderr) > 0 || len(status2.Stderr) > 0 {
		log.Err(fmt.Errorf("%v %v", status1.Stderr, status2.Stderr))
	}

	if _, err := os.Lstat(latestReportSymlink); err == nil {
		os.Remove(latestReportSymlink)
	}

	if err := os.Symlink(latestProjectPath, latestReportSymlink); err != nil {
		log.Err(err)
		return err
	}
	log.Info().Msgf("Created report symlink %s -> %s", latestProjectPath, latestReportSymlink)

	return nil
}
