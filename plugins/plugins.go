package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"go.uber.org/zap/zapcore"
)

type EnvConfig interface {
	LogLevel() zapcore.Level
	JSONConsole() bool
	LogUnixTimestamps() bool
}

func SetEnvConfig(cmd *exec.Cmd, cfg EnvConfig) {
	forward := func(name string) {
		if v, ok := os.LookupEnv(name); ok {
			cmd.Env = append(cmd.Env, name+"="+v)
		}
	}
	forward("CL_DEV")
	forward("CL_LOG_SQL_MIGRATIONS")
	forward("CL_LOG_COLOR")
	cmd.Env = append(cmd.Env,
		"CL_LOG_LEVEL="+cfg.LogLevel().String(),
		"CL_JSON_CONSOLE="+strconv.FormatBool(cfg.JSONConsole()),
		"CL_UNIX_TS="+strconv.FormatBool(cfg.LogUnixTimestamps()),
	)
}

func GetEnvConfig() (EnvConfig, error) {
	logLevelStr := os.Getenv("CL_LOG_LEVEL")
	logLevel, err := zapcore.ParseLevel(logLevelStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CL_LOG_LEVEL = %q: %w", logLevelStr, err)
	}
	return &envConfig{
		logLevel:       logLevel,
		jsonConsole:    strings.EqualFold("true", os.Getenv("CL_JSON_CONSOLE")),
		unixTimestamps: strings.EqualFold("true", os.Getenv("CL_UNIX_TS")),
	}, nil
}

type envConfig struct {
	logLevel       zapcore.Level
	jsonConsole    bool
	unixTimestamps bool
}

func (e *envConfig) LogLevel() zapcore.Level {
	return e.logLevel
}

func (e *envConfig) JSONConsole() bool {
	return e.jsonConsole
}

func (e *envConfig) LogUnixTimestamps() bool {
	return e.unixTimestamps
}
