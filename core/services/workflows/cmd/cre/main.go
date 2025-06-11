package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

func main() {
	var (
		wasmPath          string
		configPath        string
		debugMode         bool
		billingClientAddr string
		enableBeholder    bool
	)

	flag.StringVar(&wasmPath, "wasm", "", "Path to the WASM binary file")
	flag.StringVar(&configPath, "config", "", "Path to the Config file")
	flag.BoolVar(&debugMode, "debug", false, "Enable debug-level logging")
	flag.StringVar(&billingClientAddr, "billing-client-address", "", "Billing client address; Leave empty to run a local client that prints to the standard log.")
	flag.BoolVar(&enableBeholder, "beholder", false, "Enable printing beholder messages to standard log")
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

	run(ctx, lggr, binary, config, billingClientAddr, enableBeholder)
}

// run instantiates the engine, starts it and blocks until the context is canceled.
func run(
	ctx context.Context,
	lggr logger.Logger,
	binary, config []byte,
	billingClientAddr string,
	enableBeholder bool,
) {
	lggr.Infof("executing engine in process: %d", os.Getpid())

	// Create the registry and fake capabilities
	registry := capabilities.NewRegistry(lggr)
	registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})
	capabilities, err := NewFakeCapabilities(ctx, lggr, registry)
	if err != nil {
		fmt.Printf("Failed to create capabilities: %v\n", err)
		os.Exit(1)
	}

	if enableBeholder {
		_ = setupBeholder(lggr.Named("Fake_Stdlog_Beholder"))
	}

	if billingClientAddr == "" {
		billingClientAddr = "localhost:4319"
	}
	bs := NewBillingService(lggr.Named("Fake_Billing_Client"))
	err = bs.Start(ctx)
	if err != nil {
		fmt.Printf("Failed to start billing service: %v\n", err)
		os.Exit(1)
	}

	for _, cap := range capabilities {
		if err2 := cap.Start(ctx); err2 != nil {
			fmt.Printf("Failed to start capability: %v\n", err2)
			os.Exit(1)
		}

		// await the capability to be initialized if using a loop plugin
		if standardcap, ok := cap.(*standaloneLoopWrapper); ok {
			if err = standardcap.Await(ctx); err != nil {
				fmt.Printf("Failed to await capability: %v\n", err)
				os.Exit(1)
			}
		}
	}

	engine, err := NewStandaloneEngine(ctx, lggr, registry, binary, config, billingClientAddr, v2.LifecycleHooks{})
	if err != nil {
		fmt.Printf("Failed to create engine: %v\n", err)
		os.Exit(1)
	}

	err = engine.Start(ctx)
	if err != nil {
		fmt.Printf("Failed to start engine: %v\n", err)
		os.Exit(1)
	}

	<-ctx.Done()

	lggr.Info("Shutting down the Engine")
	_ = engine.Close()
	for _, cap := range capabilities {
		lggr.Infow("Shutting down capability", "id", cap.Name())
		_ = cap.Close()
	}
	_ = bs.Close()
}
