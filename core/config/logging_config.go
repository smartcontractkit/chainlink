package config

import (
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Logging interface {
	DefaultLogLevel() zapcore.Level
	LogFileDir() string
	LogLevel() zapcore.Level
	LogFileMaxSize() utils.FileSize
	LogFileMaxAge() int64
	LogFileMaxBackups() int64
	LogUnixTimestamps() bool
	JSONConsole() bool
}
