package plugins

import (
	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

func NewAppConfig(appID uuid.UUID, logLevel zapcore.Level, jsonConsole bool, unixTimestamps bool) AppConfig {
	return &appConfig{
		appID:          appID,
		logLevel:       logLevel,
		jsonConsole:    jsonConsole,
		unixTimestamps: unixTimestamps,
	}
}
