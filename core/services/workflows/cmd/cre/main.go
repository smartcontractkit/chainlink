package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func main() {
	var wasmPath string
	var configPath string
	var debugMode bool
	var billingClientAddr string

	flag.StringVar(&wasmPath, "wasm", "", "Path to the WASM binary file")
	flag.StringVar(&configPath, "config", "", "Path to the Config file")
	flag.BoolVar(&debugMode, "debug", false, "Enable debug-level logging")
	flag.StringVar(&billingClientAddr, "billing-client-address", "", "Billing client address; Leave empty for no client.")
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Set log level based on debug flag
	logLevel := zapcore.InfoLevel
	if debugMode {
		logLevel = zapcore.DebugLevel
	}

	logCfg := logger.Config{LogLevel: logLevel}
	lggr, _ := logCfg.New()

	// Create the registry and fake capabilities
	registry := capabilities.NewRegistry(lggr)
	registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
	capabilities, err := NewFakeCapabilities(ctx, lggr, registry)
	if err != nil {
		fmt.Printf("Failed to create capabilities: %v\n", err)
		os.Exit(1)
	}

	run(ctx, lggr, registry, capabilities, binary, config, billingClientAddr)
}

// run instantiates the engine, starts it and blocks until the context is canceled.
func run(
	ctx context.Context,
	lggr logger.Logger,
	registry *capabilities.Registry,
	capabilities []services.Service,
	binary, config []byte,
	billingClientAddr string,
) {
	engine, err := NewStandaloneEngine(ctx, lggr, registry, binary, config, billingClientAddr)
	if err != nil {
		fmt.Printf("Failed to create engine: %v\n", err)
		os.Exit(1)
	}

	for _, cap := range capabilities {
		if err2 := cap.Start(ctx); err2 != nil {
			fmt.Printf("Failed to start capability: %v\n", err2)
			os.Exit(1)
		}
	}
	err = engine.Start(ctx)
	if err != nil {
		fmt.Printf("Failed to start engine: %v\n", err)
		os.Exit(1)
	}

	<-ctx.Done()

	fmt.Println("Shutting down the Engine")
	_ = engine.Close()
	for _, cap := range capabilities {
		lggr.Infow("Shutting down capability", "id", cap.Name())
		_ = cap.Close()
	}
}
