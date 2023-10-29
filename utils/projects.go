package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/go-cmd/cmd"
	"github.com/gofiber/fiber/v2"
	"github.com/otiai10/copy"
	"github.com/rs/zerolog/log"
)

func ProjectsPath() string {
	return filepath.Join(GetAppDataPath(), "projects")
}

func BackupProjectsPath() string {
	return filepath.Clean(filepath.Join(GetBackupAppDataPath(), "projects"))
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

func GetProjectBackupPath(projectId string) string {
	return filepath.Join(BackupProjectsPath(), projectId)
}

func GetProjectResultsPath(projectId string) string {
	return filepath.Join(GetProjectPath(projectId), "results")
}

func GetProjectReportsPath(projectId string) string {
	return filepath.Join(GetProjectPath(projectId), "reports")
}

func GetProjectReportsBackupPath(projectId string) string {
	return filepath.Join(GetProjectBackupPath(projectId), "reports")
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
	return filepath.Join(GetProjectReportsPath(projectId), GetLatestProjectReportId(projectId))
}

func GetFullReportUrl(c *fiber.Ctx, projectId string, reportId string) string {
	fullUrl, _ := c.GetRouteURL("project_reports_endpoint", fiber.Map{"project_id": projectId, "*1": filepath.Join(reportId, "index.html")})
	return fullUrl
}

func CreateProject(projectId string) error {
	if !GetExistentsProjects(projectId) {
		if err := os.MkdirAll(GetProjectPath(projectId), os.ModeSticky|os.ModePerm); err != nil {
			log.Err(err)
			return err
		}
		if err := os.MkdirAll(GetProjectResultsPath(projectId), os.ModeSticky|os.ModePerm); err != nil {
			log.Err(err)
			return err
		}
		if err := os.MkdirAll(GetProjectReportsPath(projectId), os.ModeSticky|os.ModePerm); err != nil {
			log.Err(err)
			return err
		}
	} else {
		log.Error().Msgf("project %s already exists", projectId)
		return fmt.Errorf("project %s already exists", projectId)
	}

	return nil
}

func copyCsvReport(origin string, projectId string) string {
	projectId = strings.ReplaceAll(projectId, "/", "")

	if CheckFileExist(origin) {
		destination := GetBackupReportCSVPath() + "/" + projectId + ".csv"
		if origin != "" || destination != "" {
			if err := copy.Copy(origin, destination, copy.Options{
				PreserveTimes: true,
				PreserveOwner: true,
			}); err != nil {
				log.Error().Msgf("Skipping report csv copy, %v", err)
			}
			log.Info().Msgf("Project copied from %s to %s", origin, destination)
		} else {
			log.Error().Msg("Skipping report csv copy, origin or destination dir is empty")
		}
		return destination
	}
	return ""
}

func StoreCSVLatestReport(projectId string, reportId string) string {
	LATEST_REPORT_CSV_PATH := filepath.Join(GetProjectReportsPath(projectId), reportId, "data", "suites.csv")
	LATEST_BACKUP_REPORT_CSV_PATH := filepath.Join(GetProjectReportsBackupPath(projectId), reportId, "data", "suites.csv")

	latestCsv := copyCsvReport(LATEST_REPORT_CSV_PATH, projectId)
	latestBackupCsv := copyCsvReport(LATEST_BACKUP_REPORT_CSV_PATH, projectId)

	if latestCsv != "" {
		return latestCsv
	}

	if latestBackupCsv != "" {
		return latestBackupCsv
	}

	return ""
}

func DeleteProject(projectId string) (deletedProject string) {
	removeCmd := <-cmd.NewCmd("rm", "-rf", GetProjectPath(projectId)).Start()
	removeBackupCmd := <-cmd.NewCmd("rm", "-rf", GetProjectBackupPath(projectId)).Start()
	if len(removeCmd.Stderr) > 0 || len(removeBackupCmd.Stderr) > 0 {
		log.Error().Msgf("Error deleting project %s, %v", projectId, removeCmd.Stderr)
		return ""
	}
	log.Info().Msgf("Project %s deleted", projectId)
	return projectId
}

func DeleteProjectReport(projectId string, reportId string) {
	removeCmd := <-cmd.NewCmd("rm", "-rf", filepath.Join(GetProjectReportsPath(projectId), reportId)).Start()
	if len(removeCmd.Stderr) > 0 {
		log.Error().Msgf("Error deleting report %s, %v", reportId, removeCmd.Stderr)
	}

	//if reportId is "latest", delete symlink
	if reportId == "latest" {
		<-cmd.NewCmd("rm", "-f", filepath.Join(GetProjectReportsPath(projectId), "latest")).Start()
	}
	log.Info().Msgf("Report %s deleted", reportId)
}

func CreateReportLatestSymlink(projectId string, latestBuildOrder int) error {
	latestProjectPath := GetLatestProjectReport(projectId)
	latestReportSymlink := filepath.Join(GetProjectReportsPath(projectId), "latest")
	allureResourcesPath := GetAllureResourcesPath()

	log.Debug().Msgf("%s %s %s %s", "ln", "-sf", filepath.Join(allureResourcesPath, "app.js"), filepath.Join(latestProjectPath, "app.js"))
	log.Debug().Msgf("%s %s %s %s", "ln", "-sf", filepath.Join(allureResourcesPath, "styles.css"), filepath.Join(latestProjectPath, "styles.css"))
	<-cmd.NewCmd("ln", "-sf", filepath.Join(allureResourcesPath, "app.js"), filepath.Join(latestProjectPath, "app.js")).Start()
	<-cmd.NewCmd("ln", "-sf", filepath.Join(allureResourcesPath, "styles.css"), filepath.Join(latestProjectPath, "styles.css")).Start()

	log.Info().Msgf("Creating latest report symlink %s -> %s", latestProjectPath, latestReportSymlink)
	<-cmd.NewCmd("rm", "-f", latestReportSymlink).Start()
	<-cmd.NewCmd("ln", "-sf", latestProjectPath, latestReportSymlink).Start()
	return nil
}
