package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

type Runner struct {
	hooks RunnerHooks
}

type RunnerConfig struct {
	EnableBeholder             bool
	EnableBilling              bool
	EnableStandardCapabilities bool
	Lggr                       logger.Logger
	LifecycleHooks             v2.LifecycleHooks
}

type RunnerHooks struct {
	// Initialize hook sets up resources used by the Runner
	Initialize func(context.Context, RunnerConfig) (*capabilities.Registry, []services.Service)
	// BeforeStart hook is a testing hook that can be used to check that resources were set up
	BeforeStart func(context.Context, RunnerConfig, *capabilities.Registry, []services.Service, []*sdk.TriggerSubscription)
	// Wait hook handles blocking for the runner to keep the standalone engine running
	Wait func(context.Context, RunnerConfig, *capabilities.Registry, []services.Service)
	// AfterRun hook is a testing hook that can be used for checking engine and capability state directly after waiting
	AfterRun func(context.Context, RunnerConfig, *capabilities.Registry, []services.Service)
	// Cleanup hook shuts down the services that were started in the Initialize hook
	Cleanup func(context.Context, RunnerConfig, *capabilities.Registry, []services.Service)
	// Finally hook is a testing hook that can be used to check that resources were cleaned up
	Finally func(context.Context, RunnerConfig, *capabilities.Registry, []services.Service)
}

var emptyHook = func(context.Context, RunnerConfig, *capabilities.Registry, []services.Service) {}
var emptyBeforeStart = func(context.Context, RunnerConfig, *capabilities.Registry, []services.Service, []*sdk.TriggerSubscription) {
}

var defaultInitialize = func(ctx context.Context, cfg RunnerConfig) (*capabilities.Registry, []services.Service) {
	registry := capabilities.NewRegistry(cfg.Lggr)
	registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})

	srvcs := []services.Service{}
	if cfg.EnableBilling {
		bs := NewBillingService(logger.Named(cfg.Lggr, "Fake_Billing_Client"))
		err := bs.Start(ctx)
		if err != nil {
			fmt.Printf("Failed to start billing service: %v\n", err)
			os.Exit(1)
		}

		srvcs = append(srvcs, bs)
	}

	var caps []services.Service
	var err error

	if cfg.EnableStandardCapabilities {
		caps, err = NewCapabilities(ctx, cfg.Lggr, registry)
	} else {
		caps, err = NewFakeCapabilities(ctx, cfg.Lggr, registry)
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

	if cfg.EnableBeholder {
		_ = SetupBeholder(logger.Named(cfg.Lggr, "Fake_Stdlog_Beholder"))
	}

	return registry, srvcs
}

var defaultWait = func(ctx context.Context, cfg RunnerConfig, registry *capabilities.Registry, services []services.Service) {
	<-ctx.Done()
}

var defaultCleanup = func(ctx context.Context, cfg RunnerConfig, registry *capabilities.Registry, services []services.Service) {
	for _, service := range services {
		cfg.Lggr.Infow("Shutting down", "id", service.Name())
		_ = service.Close()
	}

	_ = cleanupBeholder()
}

func DefaultHooks() *RunnerHooks {
	return &RunnerHooks{
		Initialize:  defaultInitialize,
		BeforeStart: emptyBeforeStart,
		AfterRun:    emptyHook,
		Wait:        defaultWait,
		Cleanup:     defaultCleanup,
		Finally:     emptyHook,
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
func (r *Runner) Run(
	ctx context.Context,
	workflowName string,
	binary, config, secrets []byte,
	cfg RunnerConfig,
) {
	cfg.Lggr.Infof("executing engine in process: %d", os.Getpid())

	registry, services := r.hooks.Initialize(ctx, cfg)

	billingAddress := ""
	if cfg.EnableBilling {
		billingAddress = "localhost:4319"
	}

	engine, triggerSub, err := NewStandaloneEngine(ctx, cfg.Lggr, registry, binary, config, secrets, billingAddress, cfg.LifecycleHooks, workflowName)
	if err != nil {
		fmt.Printf("Failed to create engine: %v\n", err)
		os.Exit(1)
	}

	services = append(services, engine)

	r.hooks.BeforeStart(ctx, cfg, registry, services, triggerSub)

	err = engine.Start(ctx)
	if err != nil {
		fmt.Printf("Failed to start engine: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err2 := engine.Close(); err2 != nil {
			fmt.Printf("Failed to close engine: %v\n", err2)
		}
	}()

	r.hooks.Wait(ctx, cfg, registry, services)

	r.hooks.AfterRun(ctx, cfg, registry, services)

	r.hooks.Cleanup(ctx, cfg, registry, services)

	r.hooks.Finally(ctx, cfg, registry, services)
}
