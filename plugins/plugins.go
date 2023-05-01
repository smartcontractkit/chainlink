package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"go.uber.org/zap/zapcore"
)

type LoggingConfigurer interface {
	LogLevel() zapcore.Level
	JSONConsole() bool
	LogUnixTimestamps() bool
}

type EnvConfigurer interface {
	LoggingConfigurer
	PrometheusPort() int
}

func SetEnvConfig(cmd *exec.Cmd, cfg EnvConfigurer) {
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
		"CL_PROMETHEUS_PORT="+strconv.FormatInt(int64(cfg.PrometheusPort()), 10),
	)
}

func GetEnvConfig() (*EnvConfig, error) {
	logLevelStr := os.Getenv("CL_LOG_LEVEL")
	logLevel, err := zapcore.ParseLevel(logLevelStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CL_LOG_LEVEL = %q: %w", logLevelStr, err)
	}
	promPortStr := os.Getenv("CL_PROMETHEUS_PORT")
	promPort, err := strconv.Atoi(promPortStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CL_PROMETHEUS_PORT = %q: %w", promPortStr, err)
	}
	return &EnvConfig{
		logLevel:       logLevel,
		jsonConsole:    strings.EqualFold("true", os.Getenv("CL_JSON_CONSOLE")),
		unixTimestamps: strings.EqualFold("true", os.Getenv("CL_UNIX_TS")),
		prometheusPort: promPort,
	}, nil
}

type EnvConfig struct {
	logLevel       zapcore.Level
	jsonConsole    bool
	unixTimestamps bool
	prometheusPort int
}

func NewEnvConfig(logLevel zapcore.Level, jsonConsole bool, unixTimestamps bool, prometheusPort int) *EnvConfig {
	return &EnvConfig{
		logLevel:       logLevel,
		jsonConsole:    jsonConsole,
		unixTimestamps: unixTimestamps,
		prometheusPort: prometheusPort,
	}
}

func (e *EnvConfig) LogLevel() zapcore.Level {
	return e.logLevel
}

func (e *EnvConfig) JSONConsole() bool {
	return e.jsonConsole
}

func (e *EnvConfig) LogUnixTimestamps() bool {
	return e.unixTimestamps
}

func (e *EnvConfig) PrometheusPort() int {
	return e.prometheusPort
}
