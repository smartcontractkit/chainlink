// Copyright 2019 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var AutoStartServer = true

// Deprecated: You should use MustCreateLoggerWithServiceName
func MustCreateLogger(opts ...zap.Option) *zap.Logger {
	return MustCreateLoggerWithServiceName("Unknown service", opts...)
}

// MustCreateLogger has the same behavior as `CreateLogger` function. However, it
// automatically panic if the logger was not created successfully.
func MustCreateLoggerWithServiceName(serviceName string, opts ...zap.Option) *zap.Logger {
	return MustCreateLoggerWithLevel(serviceName, LevelFromEnvironment(), opts...)
}

// MustCreateLoggerWithLevel behaves exactly like `MustCreateLogger`, but you can pass the atomic level
// that should be used for the logger that will use this atomic level. By keeping a reference to it,
// later on, you will be able to change the level at runtime by calling `atomicLevel.SetLevel`
// on your reference and logger level will be changed.
func MustCreateLoggerWithLevel(serviceName string, atomicLevel zap.AtomicLevel, opts ...zap.Option) *zap.Logger {
	logger, err := CreateLoggerWithLevel(serviceName, atomicLevel, opts...)
	if err != nil {
		panic(fmt.Errorf("unable to create logger (in production: %t): %s", IsProductionEnvironment(), err))
	}

	return logger
}

// CreateLogger can be used to create the correct zap logger based on the environment.
//
// First, if an environment variable `ZAP_PRETTY` pretty is present, a `zapdriver.NewProduction`
// is used but logging all levels (`Debug` level and more). Furthermore, this logger will
// print everything into the standard output of the process (opposed to standard error by
// default). If the env is set, it overrides everything. If the value of the `ZAP_PRETTY`
// environment variable is a valid Zap level (`debug`, `info`, `warn`, `error`), the logger
// level will be configured using the value. In all other cases, `debug` level is used as
// the default.
//
// Then, if in production, automatically a `zap.NewProduction()` is returned. The production
// environment is determined based on the presence of the `/.dockerenv` file.
//
// In all other cases, return a `zap.NewDevelopment()` logger.
func CreateLogger(serviceName string, opts ...zap.Option) (*zap.Logger, error) {
	return CreateLoggerWithLevel(serviceName, LevelFromEnvironment(), opts...)
}

// CreateLoggerWithLevel behaves exactly like `CreateLogger`, but you can pass the atomic level
// that should be used for the logger that will use this atomic level. By keeping a reference to it,
// later on, you will be able to change the level at runtime by calling `atomicLevel.SetLevel`
// on your reference and logger level will be changed.
func CreateLoggerWithLevel(serviceName string, atomicLevel zap.AtomicLevel, opts ...zap.Option) (*zap.Logger, error) {
	config := BasicLoggingConfig(serviceName, atomicLevel, opts...)

	zlog, err := config.Build(opts...)
	if err != nil {
		return nil, err
	}

	if AutoStartServer {
		go func() {
			zlog.Info("starting atomic level switcher, port :1065")
			if err := http.ListenAndServe(":1065", atomicLevel); err != nil {
				zlog.Info("failed listening on :1065 to switch log level:", zap.Error(err))
			}
		}()
	}

	return zlog, nil
}

func LevelFromEnvironment() zap.AtomicLevel {
	zapPrettyValue := os.Getenv("ZAP_PRETTY")
	if zapPrettyValue != "" {
		return zap.NewAtomicLevelAt(zapLevelFromString(zapPrettyValue))
	}

	if IsProductionEnvironment() {
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return zap.NewAtomicLevelAt(zap.DebugLevel)
}

func IsProductionEnvironment() bool {
	_, err := os.Stat("/.dockerenv")

	return !os.IsNotExist(err)
}

// Deprecated: Will be removed in a future version, use `InstantiateLoggers` and configure it the way you want
// instead.
func BasicLoggingConfig(serviceName string, atomicLevel zap.AtomicLevel, opts ...zap.Option) *zap.Config {
	var config zap.Config

	if IsProductionEnvironment() || os.Getenv("ZAP_PRETTY") != "" {
		config = zapdriver.NewProductionConfig()
		opts = append(opts, zapdriver.WrapCore(
			zapdriver.ReportAllErrors(true),
			zapdriver.ServiceName(serviceName),
		))
	} else {
		config = zap.NewDevelopmentConfig()
	}

	if os.Getenv("ZAP_PRETTY") != "" {
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stdout"}
	}

	config.Level = atomicLevel
	return &config
}

func zapLevelFromString(input string) zapcore.Level {
	switch strings.ToLower(input) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warning", "warn":
		return zap.WarnLevel
	case "error", "err":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	case "panic":
		return zap.PanicLevel
	default:
		return zap.DebugLevel
	}
}
