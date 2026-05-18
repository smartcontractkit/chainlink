//go:build windows
// +build windows

package logger

import (
	"net/url"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// logFileURI returns the scheme and path to the log file for the passed
// directory, with a custom scheme winfile:/// specifically tailored for
// Windows to get around their handling of the file:// schema in uber.org/zap.
// https://github.com/uber-go/zap/issues/621
func (c Config) logFileURI() string {
	return "winfile:///" + filepath.ToSlash(c.LogsFile())
}

func registerOSSinks() error {
	return zap.RegisterSink("winfile", newWinFileSink)
}

func newWinFileSink(u *url.URL) (zap.Sink, error) {
	// https://github.com/uber-go/zap/issues/621
	// Remove leading slash left by url.Parse()
	return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}
