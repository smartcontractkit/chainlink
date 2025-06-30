package main

import (
	"context"
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

type Runner struct {
	hooks RunnerHooks
}

type RunnerConfig struct {
	enableBeholder             bool
	enableBilling              bool
	enableStandardCapabilities bool
	lggr                       logger.Logger
	engineHooks                v2.LifecycleHooks
}

type RunnerHooks struct {
	Initialize func(context.Context, RunnerConfig) (*capabilities.Registry, []services.Service)
	BeforeRun  func(context.Context, RunnerConfig, *capabilities.Registry, []services.Service)
	AfterRun   func(context.Context, RunnerConfig, *capabilities.Registry, []services.Service)
	Cleanup    func(context.Context, RunnerConfig, *capabilities.Registry, []services.Service)
	Finally    func(context.Context, RunnerConfig, *capabilities.Registry, []services.Service)
}

var emptyHook = func(context.Context, RunnerConfig, *capabilities.Registry, []services.Service) {}

var defaultInitialize = func(ctx context.Context, cfg RunnerConfig) (*capabilities.Registry, []services.Service) {
	registry := capabilities.NewRegistry(cfg.lggr)
	registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})

	srvcs := []services.Service{}

	if cfg.enableBilling {
		bs := NewBillingService(cfg.lggr.Named("Fake_Billing_Client"))
		err := bs.Start(ctx)
		if err != nil {
			fmt.Printf("Failed to start billing service: %v\n", err)
			os.Exit(1)
		}

		srvcs = append(srvcs, bs)
	}

	var caps []services.Service
	var err error

	if cfg.enableStandardCapabilities {
		caps, err = NewCapabilities(ctx, cfg.lggr, registry)
	} else {
		caps, err = NewFakeCapabilities(ctx, cfg.lggr, registry)
	}
	if err != nil {
		fmt.Printf("Failed to create capabilities: %v\n", err)
		os.Exit(1)
	}

	for _, cap := range caps {
		err = cap.Start(ctx)
		if err != nil {
			fmt.Printf("Failed to start capability: %v\n", err)
			os.Exit(1)
		}

		// await the capability to be initialized if using a loop plugin
		if standardcap, ok := cap.(*standaloneLoopWrapper); ok {
			err = standardcap.Await(ctx)
			if err != nil {
				fmt.Printf("Failed to await capability: %v\n", err)
				os.Exit(1)
			}
		}

		srvcs = append(srvcs, cap)
	}

	if cfg.enableBeholder {
		_ = setupBeholder(cfg.lggr.Named("Fake_Stdlog_Beholder"))
	}

	return registry, srvcs
}

var defaultCleanup = func(ctx context.Context, cfg RunnerConfig, registry *capabilities.Registry, services []services.Service) {
	for _, service := range services {
		cfg.lggr.Infow("Shutting down", "id", service.Name())
		_ = service.Close()
	}

	_ = cleanupBeholder()
}

func DefaultHooks() *RunnerHooks {
	return &RunnerHooks{
		Initialize: defaultInitialize,
		BeforeRun:  emptyHook,
		AfterRun:   emptyHook,
		Cleanup:    defaultCleanup,
		Finally:    emptyHook,
	}
}

func NewRunner(hooks *RunnerHooks) *Runner {
	if hooks == nil {
		hooks = DefaultHooks()
	}

	return &Runner{
		hooks: *hooks,
	}
}

// run instantiates the engine, starts it and blocks until the context is canceled.
func (r *Runner) run(
	ctx context.Context,
	binary, config, secrets []byte,
	cfg RunnerConfig,
) {
	cfg.lggr.Infof("executing engine in process: %d", os.Getpid())

	registry, services := r.hooks.Initialize(ctx, cfg)

	billingAddress := ""
	if cfg.enableBilling {
		billingAddress = "localhost:4319"
	}

	engine, err := NewStandaloneEngine(ctx, cfg.lggr, registry, binary, config, secrets, billingAddress, cfg.engineHooks)
	if err != nil {
		fmt.Printf("Failed to create engine: %v\n", err)
		os.Exit(1)
	}

	services = append(services, engine)

	r.hooks.BeforeRun(ctx, cfg, registry, services)

	err = engine.Start(ctx)
	if err != nil {
		fmt.Printf("Failed to start engine: %v\n", err)
		os.Exit(1)
	}

	<-ctx.Done()

	r.hooks.AfterRun(ctx, cfg, registry, services)

	r.hooks.Cleanup(ctx, cfg, registry, services)

	r.hooks.Finally(ctx, cfg, registry, services)
}
