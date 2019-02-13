// +build !windows

package logger

import "path/filepath"

func registerOSSinks() error {
	return nil
}

// ProductionLoggerFilepath returns the full path to the file the
// ProductionLogger logs to, and uses zap's built in default file sink.
func ProductionLoggerFilepath(configRootDir string) string {
	return filepath.ToSlash(filepath.Join(configRootDir, "log.jsonl"))
}
