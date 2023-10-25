package utils

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func GetEnvValue(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func GetAppDataPath() string {
	workpath, _ := os.Getwd()
	rootDir := filepath.Join(workpath, os.Getenv("APP_DATA_DIR"))
	return rootDir
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
