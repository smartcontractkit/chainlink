//go:build !windows
// +build !windows

package logger

import "path/filepath"

func registerOSSinks() error {
	return nil
}

// logFileURI returns the full path to the file the
// NewLogger logs to, and uses zap's built in default file sink.
func (c Config) logFileURI() string {
	return filepath.ToSlash(c.LogsFile())
}
