package utils

import (
	"errors"
	"os"
	"path/filepath"
)

func ExpandHomeDir(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}
	if path[0] != '~' {
		return path, nil
	}
	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, path[1:]), nil
}
