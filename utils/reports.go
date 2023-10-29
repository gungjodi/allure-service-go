package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"osp-allure/models"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/go-cmd/cmd"
	"github.com/otiai10/copy"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

func GetSortedReportsDir(projectId string) (reportDirsName []int) {
	reportsDir, err := os.ReadDir(GetProjectReportsPath(projectId))

	if err != nil {
		log.Logger.Error().Msgf("Error reading reports directory: %s", err)
	}

	if len(reportsDir) > 0 {
		for _, dir := range reportsDir {
			if dir.Name() != "latest" {
				dirInt, _ := strconv.Atoi(dir.Name())
				reportDirsName = append(reportDirsName, dirInt)
			}
		}
		slices.Sort(reportDirsName)
	}
	return reportDirsName
}

func GetLatestProjectBuildOrder(projectId string) int {
	latestReport := GetLatestProjectReport(projectId)

	if !CheckFileOrDirExist(latestReport) {
		log.Logger.Info().Msgf("No latest report for projectId %s, latest buildOrder is 0", projectId)
		return 0
	}

	executorFile := filepath.Join(latestReport, "widgets", "executors.json")
	content, err := os.ReadFile(executorFile)
	if err != nil {
		log.Logger.Error().Msgf("Error reading executor file for projectId %s: %s, latest buildOrder is 0", projectId, err)
		return 0
	}

	executorJSON := []models.ExecutorInfo{}
	if err := json.Unmarshal(content, &executorJSON); err != nil {
		log.Logger.Error().Msgf("Error unmarshal executor file for projectId %s: %s", projectId, err)
		return 0
	}

	if len(executorJSON) == 0 {
		log.Logger.Info().Msgf("No executor file for projectId %s, latest buildOrder is 0", projectId)
		return 0
	}
	return executorJSON[0].BuildOrder
}

func KeepReportHistory(projectId string) {
	keepHistoryLatest, _ := strconv.Atoi(GetEnv("KEEP_HISTORY_LATEST", "0"))
	projectReportPath := GetProjectReportsPath(projectId)

	reportDirsName := GetSortedReportsDir(projectId)
	log.Info().Msgf("Keeping latest history, max = %d", keepHistoryLatest)
	if keepHistoryLatest > 0 && (len(reportDirsName) > keepHistoryLatest) {
		sizeToRemove := len(reportDirsName) - keepHistoryLatest
		for i := 0; i < sizeToRemove; i++ {
			os.RemoveAll(filepath.Join(projectReportPath, strconv.Itoa(reportDirsName[i])))
			log.Info().Msgf("Removed report history %d for PROJECT_ID: %s", reportDirsName[i], projectId)
		}
	}
}

func GenerateExecutorJson(projectId string, buildOrder int, executionName string, executionFrom string, executionType string) error {
	if executionName == "" {
		executionName = "Execution On Demand"
	}
	if executionType == "" {
		executionType = "another"
	}
	executorMap := models.ExecutorInfo{
		BuildName:  fmt.Sprintf("%s #%d", projectId, buildOrder),
		ReportName: projectId,
		BuildOrder: buildOrder,
		Name:       executionName,
		ReportURL:  fmt.Sprintf("../%d/index.html", buildOrder),
		BuildURL:   executionFrom,
		Type:       executionType,
	}

	file, err := json.MarshalIndent(executorMap, "", "    ")
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("Generating executor.json for projectId %s on build order %d", projectId, buildOrder)

	if err := os.WriteFile(filepath.Join(GetProjectResultsPath(projectId), "executor.json"), file, os.ModePerm); err != nil {
		log.Error().Err(err)
		return err
	}
	return nil
}

func GenerateReportCmd(projectId string, newBuildOrder int) error {
	resultPath := filepath.Join(GetProjectResultsPath(projectId))
	allureCmd := GetEnv("LOCAL_ALLURE_EXECUTABLE", "allure")

	generateAllureCmd := cmd.NewCmd(
		allureCmd, "generate", "--clean", resultPath,
		"-o", filepath.Join(GetProjectReportsPath(projectId), strconv.Itoa(newBuildOrder)),
	)
	log.Info().Msgf("Generating report for project %s", projectId)
	log.Debug().Msgf("[allure] %v %v", allureCmd, generateAllureCmd.Args)
	status := <-generateAllureCmd.Start()

	for _, line := range status.Stderr {
		if strings.Contains(line, "Caused by") {
			log.Error().Msgf("[allure] %v", line)
			return fmt.Errorf("error generating report for project %s : %v", projectId, line)
		}
	}

	for _, line := range status.Stdout {
		log.Info().Msgf("[allure] %v", line)
		if strings.Contains(line, "successfully generated") {
			return nil
		}
	}

	return fmt.Errorf("error generating report for project %s", projectId)
}

func KeepResultHistory(projectId string) {
	keepResultHistory, _ := strconv.ParseBool(os.Getenv("KEEP_RESULTS_HISTORY"))
	projectResultsHistory := filepath.Join(GetProjectResultsPath(projectId), "history")
	projectLatestReportHistory := filepath.Join(GetLatestProjectReport(projectId), "history")

	if keepResultHistory {
		log.Info().Msgf("Creating history on results directory for PROJECT_ID: %s ...", projectId)
		if err := os.MkdirAll(projectResultsHistory, os.ModePerm); err != nil {
			log.Error().Msgf("Error creating history directory on results for PROJECT_ID: %s", projectId)
		}
		if CheckFileOrDirExist(projectLatestReportHistory) {
			log.Info().Msgf("Copying history from previous results on projectId %s ...", projectId)
			if err := copy.Copy(projectLatestReportHistory, projectResultsHistory, copy.Options{
				PreserveTimes: true,
				PreserveOwner: true,
			}); err != nil {
				log.Error().Msgf("Error copying history from previous results for PROJECT_ID: %s", projectId)
			}
		}
	} else {
		if CheckFileOrDirExist(projectResultsHistory) {
			log.Info().Msgf("Removing history directory from results for PROJECT_ID: %s ...", projectId)
			if err := os.RemoveAll(projectResultsHistory); err != nil {
				log.Error().Msgf("Error removing history directory from results for PROJECT_ID: %s", projectId)
			}
		}
	}
}

func ExtractResultsArchive(tarFileName string, projectResultsDir string) error {
	log.Debug().Msgf("%v %v %v %v %v", "tar", "-xf", tarFileName, "-C", projectResultsDir)
	extractCmd := cmd.NewCmd("tar", "-xpf", tarFileName, "-C", projectResultsDir)
	status := <-extractCmd.Start()
	if len(status.Stderr) > 0 {
		log.Error().Msgf("Extracting tar file error, %v", status.Stderr)
		return fmt.Errorf("%v", status.Stderr)
	}
	go func() {
		os.Remove(tarFileName)
	}()

	return nil
}

func CleanResults(projectId string) error {
	projectResultsDir := GetProjectResultsPath(projectId)
	if CheckFileOrDirExist(projectResultsDir) {
		log.Info().Msgf("Cleaning results directory for PROJECT_ID: %s ...", projectId)
		if err := os.RemoveAll(projectResultsDir); err != nil {
			log.Error().Msgf("Error cleaning results directory for PROJECT_ID: %s", projectId)
			return err
		}
	} else {
		log.Info().Msgf("Results directory for PROJECT_ID: %s does not exist", projectId)
	}
	return nil
}

func CleanHistory(projectId string) error {
	projectResultsHistory := filepath.Join(GetProjectResultsPath(projectId), "history")
	executorFile := filepath.Join(GetProjectResultsPath(projectId), "executor.json")

	if CheckFileOrDirExist(projectResultsHistory) {
		log.Info().Msgf("Cleaning history directory for PROJECT_ID: %s ...", projectId)
		if err := os.RemoveAll(projectResultsHistory); err != nil {
			log.Error().Msgf("Error cleaning history directory for PROJECT_ID: %s", projectId)
		}
		if err := os.Remove(executorFile); err != nil {
			log.Error().Msgf("Error removing executor.json file for PROJECT_ID: %s", projectId)
		}
	} else {
		log.Info().Msgf("History directory for PROJECT_ID: %s does not exist", projectId)
	}
	return nil
}

func BackupReport(projectId string, reportId string, shouldDelete bool) error {
	StoreCSVLatestReport(projectId, reportId)

	backupProjectReportPath := filepath.Join(GetProjectReportsBackupPath(projectId), reportId)

	if !strings.HasPrefix(backupProjectReportPath, "/") {
		log.Error().Any("Error backup project", "baseTarget must be an absolute path")
		return fmt.Errorf("baseTarget must be an absolute path")
	}

	projectReportPath := filepath.Join(GetProjectReportsPath(projectId), reportId)

	if CheckFileOrDirExist(projectReportPath) {
		log.Info().Msgf("backuping project %s to %s", projectId, backupProjectReportPath)
		if err := copy.Copy(projectReportPath, backupProjectReportPath, copy.Options{
			PreserveTimes: true,
			PreserveOwner: true,
			OnSymlink: func(src string) copy.SymlinkAction {
				return copy.Skip
			},
		}); err != nil {
			log.Error().Msgf("Error backup project %s to %s, %v", projectId, backupProjectReportPath, err)
			return err
		}
	}

	CreateProjectBackupSymlink(projectId, reportId)

	if shouldDelete {
		DeleteProjectReport(projectId, reportId)
	}

	return nil
}

func CreateProjectBackupSymlink(projectId string, reportId string) {
	backupProjectReportPath := filepath.Join(GetProjectReportsBackupPath(projectId), reportId)
	latestBackupReportSymlink := filepath.Join(GetProjectReportsBackupPath(projectId), "latest")
	allureResourcesPath := GetAllureResourcesPath()

	log.Info().Msgf("Creating backup report symlink %s -> %s", backupProjectReportPath, latestBackupReportSymlink)

	log.Debug().Msgf("%s %s %s %s", "ln", "-sf", filepath.Join(allureResourcesPath, "app.js"), filepath.Join(backupProjectReportPath, "app.js"))
	<-cmd.NewCmd("ln", "-sf", filepath.Join(allureResourcesPath, "app.js"), filepath.Join(backupProjectReportPath, "app.js")).Start()

	log.Debug().Msgf("%s %s %s %s", "ln", "-sf", filepath.Join(allureResourcesPath, "styles.css"), filepath.Join(backupProjectReportPath, "styles.css"))
	<-cmd.NewCmd("ln", "-sf", filepath.Join(allureResourcesPath, "styles.css"), filepath.Join(backupProjectReportPath, "styles.css")).Start()

	log.Info().Msgf("Creating backup report symlink %s -> %s", backupProjectReportPath, latestBackupReportSymlink)
	log.Debug().Msgf("%s %s %s %s", "ln", "-sf", backupProjectReportPath, latestBackupReportSymlink)
	<-cmd.NewCmd("rm", "-f", latestBackupReportSymlink).Start()
	<-cmd.NewCmd("ln", "-sf", backupProjectReportPath, latestBackupReportSymlink).Start()
}

func GetLatestProjectReportIdBackup(projectId string) int {
	projectReportBackupPath := GetProjectReportsBackupPath(projectId)

	if !CheckFileOrDirExist(projectReportBackupPath) {
		log.Info().Msgf("No backup report for projectId %s", projectId)
		return 0
	}
	backupReportsDir, err := os.ReadDir(projectReportBackupPath)

	if err != nil {
		log.Logger.Error().Msgf("Error reading backup reports directory: %s", err)
	}

	var reportDirsName []int

	if len(backupReportsDir) > 0 {
		for _, dir := range backupReportsDir {
			if dir.Name() != "latest" {
				dirInt, _ := strconv.Atoi(dir.Name())
				reportDirsName = append(reportDirsName, dirInt)
			}
		}
		sort.Sort(sort.Reverse(sort.IntSlice(reportDirsName)))
		return reportDirsName[0]
	}

	return 0
}

func CopyResultHistoryFromBackup(projectId string, reportId string) {
	projectResultsHistory := filepath.Join(GetProjectResultsPath(projectId), "history")
	projectReportHistoryBackupPath := filepath.Join(GetProjectReportsBackupPath(projectId), reportId, "history")

	if CheckFileOrDirExist(projectReportHistoryBackupPath) {
		log.Info().Msgf("Copying history from backup results on projectId %s, reportId %s ...", projectId, reportId)
		if err := copy.Copy(projectReportHistoryBackupPath, projectResultsHistory, copy.Options{
			PreserveTimes: true,
			PreserveOwner: true,
		}); err != nil {
			log.Error().Msgf("Error copying history from previous results for PROJECT_ID: %s", projectId)
		}
	}
}
