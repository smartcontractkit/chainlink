package utils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func DefaultDataDir() (string, error) {
	userDataDir := os.Getenv("XDG_DATA_HOME")
	if runtime.GOOS == "windows" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		userDataDir = configDir
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dataDir := DataDir(runtime.GOOS, userDataDir, userHomeDir)
	if dataDir == "" {
		return "", errors.New("data directory not found")
	}
	return dataDir, nil
}

func DataDir(operatingSystem, userDataDir, homeDir string) string {
	local, share, juno := ".local", "share", "juno"

	if operatingSystem == "" || (userDataDir == "" && homeDir == "") {
		return ""
	}

	if operatingSystem == "windows" {
		if userDataDir == "" {
			return ""
		}
		return filepath.Join(userDataDir, juno)
	}

	if userDataDir != "" {
		return filepath.Join(userDataDir, juno)
	}
	return filepath.Join(homeDir, local, share, juno)
}
