package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"osp-allure/models"
	"path/filepath"
	"strconv"

	"github.com/go-cmd/cmd"
	"github.com/otiai10/copy"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

func GetSortedReportsDir(projectId string) (reportDirsName []int) {
	reportsDir, err := os.ReadDir(filepath.Join(GetProjectPath(projectId), "reports"))

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
	keepHistoryLatest, _ := strconv.Atoi(os.Getenv("KEEP_HISTORY_LATEST"))
	projectReportPath := filepath.Join(GetProjectPath(projectId), "reports")

	reportDirsName := GetSortedReportsDir(projectId)
	log.Info().Msgf("Keeping latest history, max = %d", keepHistoryLatest)
	if len(reportDirsName) > keepHistoryLatest {
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
		ReportURL:  fmt.Sprintf("../%d", buildOrder),
		BuildURL:   executionFrom,
		Type:       executionType,
	}

	file, err := json.MarshalIndent(executorMap, "", " ")
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Info().Msgf("Generating executor.json for projectId %s on build order %d", projectId, buildOrder)

	if err := os.WriteFile(filepath.Join(GetProjectPath(projectId), "results", "executor.json"), file, os.ModePerm); err != nil {
		log.Error().Err(err)
		return err
	}
	return nil
}

func GenerateReportCmd(projectId string, newBuildOrder int) cmd.Status {
	resultPath := filepath.Join(GetProjectPath(projectId), "results")
	allureCmd := "allure"
	localAllure := os.Getenv("LOCAL_ALLURE_HOME")
	if localAllure != "" {
		allureCmd = localAllure
	}

	log.Info().Msgf("Generating report for project %s", projectId)
	log.Debug().Msgf("%v %v %v %v %v %v",
		allureCmd, "generate", "--clean", resultPath,
		"-o", filepath.Join(GetProjectPath(projectId), "reports", strconv.Itoa(newBuildOrder)),
	)

	generateAllureCmd := cmd.NewCmdOptions(
		cmd.Options{Streaming: true},
		allureCmd, "generate", "--clean", resultPath,
		"-o", filepath.Join(GetProjectPath(projectId), "reports", strconv.Itoa(newBuildOrder)),
	)
	statusChan := generateAllureCmd.Start()

	if <-generateAllureCmd.Stderr != "" {
		log.Err(fmt.Errorf(<-generateAllureCmd.Stderr))
	}
	log.Info().Msgf("%v", <-generateAllureCmd.Stdout)

	return <-statusChan
}

func KeepResultHistory(projectId string) {
	keepResultHistory, _ := strconv.ParseBool(os.Getenv("KEEP_RESULTS_HISTORY"))
	projectResultsHistory := filepath.Join(GetProjectPath(projectId), "results", "history")
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
