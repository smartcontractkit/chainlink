package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"go.uber.org/zap/zapcore"
)

// LoggingConfig controls static logging related configuration that is inherited from the chainlink application to the
// given LOOP executable.
type LoggingConfig interface {
	LogLevel() zapcore.Level
	JSONConsole() bool
	LogUnixTimestamps() bool
}

// ProcessConfig generates configuration for loop commands
type ProcessConfig interface {
	LoggingConfig
	RegisterLOOP(loopId string) *RegisteredLoop
	//GenerateEnvConfig(loopId string) EnvConfig
}

type processConfig struct {
	LoggingConfig
	//portReservationFn  func(id string) int
	loopEnvGeneratorFn func(loopId string, loopStaticCfg LoggingConfig) *RegisteredLoop
}

// NewProcessConfig loopEnvGeneratorFn must act as a global, idempotent registry function.
func NewProcessConfig(lc LoggingConfig, loopEnvGeneratorFn func(loopId string, loopStaticCfg LoggingConfig) *RegisteredLoop) ProcessConfig {
	return &processConfig{
		LoggingConfig: lc,
		//portReservationFn:  portReservationFn,
		loopEnvGeneratorFn: loopEnvGeneratorFn,
	}
}

func (pc *processConfig) RegisterLOOP(id string) *RegisteredLoop {
	return pc.loopEnvGeneratorFn(id, pc.LoggingConfig)
}

/*
func (pc *processConfig) GenerateEnvConfig(loopId string) EnvConfig {
	return NewEnvConfig(pc.LogLevel(), pc.JSONConsole(), pc.LogUnixTimestamps(), pc.portReservationFn(loopId))
}
*/
// EnvConfig is the configuration interface between the application and the LOOP, which is passed via the environment.
// It separates static and dynamic configuration. Logging configuration can and is inherited statically while the
// port the the LOOP is to use for prometheus, which is created dynamically at run time the chainlink Application.
type EnvConfig interface {
	LoggingConfig
	PrometheusPort() int
}

func SetCmdEnvFromConfig(cmd *exec.Cmd, cfg EnvConfig) {
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

func GetEnvConfig() (EnvConfig, error) {
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
	return &envConfig{
		logLevel:       logLevel,
		jsonConsole:    strings.EqualFold("true", os.Getenv("CL_JSON_CONSOLE")),
		unixTimestamps: strings.EqualFold("true", os.Getenv("CL_UNIX_TS")),
		prometheusPort: promPort,
	}, nil
}

// envConfig is an implementation of EnvConfig.
type envConfig struct {
	logLevel       zapcore.Level
	jsonConsole    bool
	unixTimestamps bool
	prometheusPort int
}

func NewEnvConfig(logLevel zapcore.Level, jsonConsole bool, unixTimestamps bool, prometheusPort int) EnvConfig {
	//prometheusPort := prometheusPortFn(name)
	return &envConfig{
		logLevel:       logLevel,
		jsonConsole:    jsonConsole,
		unixTimestamps: unixTimestamps,
		prometheusPort: prometheusPort,
	}
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

func (e *envConfig) PrometheusPort() int {
	return e.prometheusPort
}
