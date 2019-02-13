// +build windows

package logger

import (
	"net/url"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// ProductionLoggerFilepath returns the full path to the file the
// ProductionLogger logs to, with a custom scheme specifically tailored for
// Windows to get around their handling of the file:// schema.
// https://github.com/uber-go/zap/issues/621
func ProductionLoggerFilepath(configRootDir string) string {
	return "winfile:///" + filepath.ToSlash(filepath.Join(configRootDir, "log.jsonl"))
}

func registerOSSinks() error {
	return zap.RegisterSink("winfile", newWinFileSink)
}

func newWinFileSink(u *url.URL) (zap.Sink, error) {
	// https://github.com/uber-go/zap/issues/621
	// Remove leading slash left by url.Parse()
	return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}
