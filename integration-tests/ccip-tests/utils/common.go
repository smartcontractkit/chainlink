package utils

import (
	"path/filepath"
	"runtime"
)

func ProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "/..")
}
