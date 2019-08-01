package orm

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

// LogLevel determines the verbosity of the events to be logged.
type LogLevel struct {
	zapcore.Level
}

// ForGin keeps Gin's mode at the appropriate level with the LogLevel.
func (ll LogLevel) ForGin() string {
	switch {
	case ll.Level < zapcore.InfoLevel:
		return gin.DebugMode
	default:
		return gin.ReleaseMode
	}
}
