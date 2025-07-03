package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/cmd/cre/utils"
)

func main() {
	var (
		wasmPath                   string
		configPath                 string
		secretsPath                string
		debugMode                  bool
		enableBeholder             bool
		enableBilling              bool
		enableStandardCapabilities bool
	)

	flag.StringVar(&wasmPath, "wasm", "", "Path to the WASM binary file")
	flag.StringVar(&configPath, "config", "", "Path to the Config file")
	flag.StringVar(&secretsPath, "secrets", "", "Path to the secrets file")
	flag.BoolVar(&debugMode, "debug", false, "Enable debug-level logging")
	flag.BoolVar(&enableBeholder, "beholder", false, "Enable printing beholder messages to standard log")
	flag.BoolVar(&enableBilling, "billing", false, "Enable to run a faked billing service that prints to the standard log.")
	flag.BoolVar(&enableStandardCapabilities, "standardCapabilities", true, "Enable to use the latest production standard capability binaries for capabilities. The binaries must be available in local GOBIN.")
	flag.Parse()

	if wasmPath == "" {
		fmt.Println("--wasm must be set")
		os.Exit(1)
	}

	binary, err := os.ReadFile(wasmPath)
	if err != nil {
		fmt.Printf("Failed to read WASM binary file: %v\n", err)
		os.Exit(1)
	}

	var config []byte
	if configPath != "" {
		config, err = os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Failed to read config file: %v\n", err)
			os.Exit(1)
		}
	}

	var secrets []byte
	if secretsPath != "" {
		secrets, err = os.ReadFile(secretsPath)
		if err != nil {
			fmt.Printf("Failed to read secrets file: %v\n", err)
			os.Exit(1)
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Set log level based on debug flag
	logLevel := zapcore.InfoLevel
	if debugMode {
		logLevel = zapcore.DebugLevel
	}

	logCfg := logger.Config{LogLevel: logLevel}
	lggr, _ := logCfg.New()

	runner := utils.NewRunner(nil)
	runner.Run(ctx, "", binary, config, secrets, utils.RunnerConfig{
		EnableBilling:              enableBilling,
		EnableBeholder:             enableBeholder,
		EnableStandardCapabilities: enableStandardCapabilities,
		Lggr:                       lggr,
	})
}
