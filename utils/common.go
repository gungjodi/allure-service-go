package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-cmd/cmd"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func GetAppDataPath() string {
	appDataDir := strings.TrimSpace(os.Getenv("APP_DATA_DIR"))

	if appDataDir == "" {
		log.Fatal().Msg("APP_DATA_DIR IS NOT SET !!!")
	}

	// If the path is absolute, return it
	if strings.HasPrefix(appDataDir, "/") {
		return appDataDir
	}

	// If the path is relative, return it relative to the current directory
	currentDir, _ := os.Getwd()
	return filepath.Join(currentDir, appDataDir)
}

func GetBackupAppDataPath() string {
	backupAppDataDir := strings.TrimSpace(os.Getenv("BACKUP_DATA_DIR"))

	if backupAppDataDir == "" {
		log.Fatal().Msg("BACKUP_DATA_DIR IS NOT SET !!!")
	}
	// If the path is absolute, return it
	if strings.HasPrefix(backupAppDataDir, "/") {
		return backupAppDataDir
	}

	// If the path is relative, return it relative to the current directory
	currentDir, _ := os.Getwd()
	return filepath.Join(currentDir, backupAppDataDir)
}

func GetBackupReportCSVPath() string {
	reportCsvPath := strings.TrimSpace(os.Getenv("DOWNLOAD_REPORT_CSV_DESTINATION_PATH"))

	if reportCsvPath == "" {
		log.Fatal().Msg("BACKUP_DATA_DIR IS NOT SET !!!")
	}
	// If the path is absolute, return it
	if strings.HasPrefix(reportCsvPath, "/") {
		return reportCsvPath
	}

	// If the path is relative, return it relative to the current directory
	currentDir, _ := os.Getwd()
	return filepath.Join(currentDir, reportCsvPath)
}

func GetAllureResourcesPath() string {
	allureResourcesPath := strings.TrimSpace(os.Getenv("ALLURE_RESOURCES"))

	if allureResourcesPath == "" {
		log.Fatal().Msg("ALLURE_RESOURCES IS NOT SET !!!")
	}

	// If the path is absolute, return it
	if strings.HasPrefix(allureResourcesPath, "/") {
		return allureResourcesPath
	}

	// If the path is relative, return it relative to the current directory
	currentDir, _ := os.Getwd()
	return filepath.Join(currentDir, allureResourcesPath)
}

func CheckFileOrDirExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			log.Err(err)
			return false
		}
	}

	return true
}

func CheckFileExist(path string) bool {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if stat.IsDir() {
		return false
	}

	return true
}

func MultipartForm(c *fiber.Ctx) (*multipart.Form, error) {
	buff := c.Context().PostBody()
	reader := bytes.NewReader(buff)

	partForm, err := readMultipartForm(reader, string(c.Request().Header.MultipartFormBoundary()), reader.Size(), reader.Size())
	defer partForm.RemoveAll()
	return partForm, err
}

func readMultipartForm(r io.Reader, boundary string, size int64, maxInMemoryFileSize int64) (*multipart.Form, error) {
	// Do not care about memory allocations here, since they are tiny
	// compared to multipart data (aka multi-MB files) usually sent
	// in multipart/form-data requests.

	if size <= 0 {
		return nil, fmt.Errorf("form size must be greater than 0. Given %d", size)
	}
	lr := io.LimitReader(r, size)
	mr := multipart.NewReader(lr, boundary)
	f, err := mr.ReadForm(maxInMemoryFileSize + (10 << 20))
	if err != nil {
		return nil, fmt.Errorf("cannot read multipart/form-data body: %w", err)
	}
	return f, nil
}

func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func GetAllureVersion() (string, error) {
	allureCmd := GetEnv("LOCAL_ALLURE_EXECUTABLE", "allure")

	getAllureVersion := <-cmd.NewCmd(allureCmd, "--version").Start()

	if getAllureVersion.Error != nil {
		return "", fmt.Errorf("%v - make sure Allure executable configured properly!", getAllureVersion.Error.Error())
	}

	return getAllureVersion.Stdout[0], nil
}
