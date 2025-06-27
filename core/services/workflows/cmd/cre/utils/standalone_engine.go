package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/billing"
	consensusserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/consensus/server"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/fakes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

const (
	defaultMaxUncompressedBinarySize = 1000000000
	defaultRPS                       = 1000.0
	defaultBurst                     = 1000
	defaultWorkflowID                = "1111111111111111111111111111111111111111111111111111111111111111"
	defaultOwner                     = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	defaultName                      = "myworkflow"
)

func NewStandaloneEngine(
	ctx context.Context,
	lggr logger.Logger,
	registry *capabilities.Registry,
	binary, config []byte,
	billingClientAddr string,
	lifecycleHooks v2.LifecycleHooks,
) (services.Service, *sdkpb.TriggerSubscriptionRequest, error) {
	labeler := custmsg.NewLabeler()
	moduleConfig := &host.ModuleConfig{
		Logger:                  lggr,
		Labeler:                 labeler,
		MaxCompressedBinarySize: defaultMaxUncompressedBinarySize,
		IsUncompressed:          true,
	}

	module, err := host.NewModule(moduleConfig, binary, host.WithDeterminism())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create module from config: %w", err)
	}

	result, err := module.Execute(ctx, &sdkpb.ExecuteRequest{
		Request:         &sdkpb.ExecuteRequest_Subscribe{},
		MaxResponseSize: uint64(1000000000),
		Config:          config,
	}, &v2.DisallowedExecutionHelper{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute subscribe: %w", err)
	}
	if result.GetError() != "" {
		return nil, nil, fmt.Errorf("failed to execute subscribe: %s", result.GetError())
	}
	triggerSubscriptions := result.GetTriggerSubscriptions()

	name, err := types.NewWorkflowName(defaultName)
	if err != nil {
		return nil, nil, err
	}

	rl, err := ratelimiter.NewRateLimiter(ratelimiter.Config{
		GlobalRPS:      defaultRPS,
		GlobalBurst:    defaultBurst,
		PerSenderRPS:   defaultRPS,
		PerSenderBurst: defaultBurst,
	})
	if err != nil {
		return nil, nil, err
	}

	workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{
		Global:   1000000000,
		PerOwner: 1000000000,
	})
	if err != nil {
		return nil, nil, err
	}

	var billingClient billing.WorkflowClient
	if billingClientAddr != "" {
		billingClient, _ = billing.NewWorkflowClient(billingClientAddr)
	}

	if module.IsLegacyDAG() {
		sdkSpec, err := host.GetWorkflowSpec(ctx, moduleConfig, binary, config)
		if err != nil {
			return nil, nil, err
		}

		cfg := workflows.Config{
			Lggr:                 lggr,
			Workflow:             *sdkSpec,
			WorkflowID:           defaultWorkflowID,
			WorkflowOwner:        defaultOwner,
			WorkflowName:         name,
			Registry:             registry,
			Store:                store.NewInMemoryStore(lggr, clockwork.NewRealClock()),
			Config:               config,
			Binary:               binary,
			SecretsFetcher:       SecretsFor,
			RateLimiter:          rl,
			WorkflowLimits:       workflowLimits,
			NewWorkerTimeout:     time.Minute,
			StepTimeout:          time.Minute,
			MaxExecutionDuration: time.Minute,
			BillingClient:        billingClient,
		}

		engine, err := workflows.NewEngine(ctx, cfg)
		if err != nil {
			return nil, nil, err
		}
		return engine, triggerSubscriptions, nil
	}

	cfg := &v2.EngineConfig{
		Lggr:            lggr,
		Module:          module,
		WorkflowConfig:  config,
		CapRegistry:     registry,
		ExecutionsStore: store.NewInMemoryStore(lggr, clockwork.NewRealClock()),

		WorkflowID:    defaultWorkflowID,
		WorkflowOwner: defaultOwner,
		WorkflowName:  name,

		LocalLimits:          v2.EngineLimits{},
		GlobalLimits:         workflowLimits,
		ExecutionRateLimiter: rl,

		BeholderEmitter: custmsg.NewLabeler(),

		BillingClient: billingClient,
		Hooks:         lifecycleHooks,

		DebugMode: true,
	}

	engine, err := v2.NewEngine(cfg)
	if err != nil {
		return nil, nil, err
	}
	return engine, triggerSubscriptions, nil
}

// TODO support fetching secrets (from a local file)
func SecretsFor(ctx context.Context, workflowOwner, hexWorkflowName, decodedWorkflowName, workflowID string) (map[string]string, error) {
	return map[string]string{}, nil
}

// NewCapabilities builds capabilities using latest standard capabilities where possible, otherwise filled in with faked capabilities.
// Capabilities are then registered with the capability registry.
func NewCapabilities(ctx context.Context, lggr logger.Logger, registry *capabilities.Registry) ([]services.Service, error) {
	caps, err := NewFakeCapabilities(ctx, lggr, registry)
	if err != nil {
		return nil, err
	}

	caps = append(caps, newStandardCapabilities(standardCapabilities, lggr, registry)...)

	return caps, nil
}

func NewFakeCapabilities(ctx context.Context, lggr logger.Logger, registry *capabilities.Registry) ([]services.Service, error) {
	caps := make([]services.Service, 0)

	fakeConsensus, err := fakes.NewFakeConsensus(lggr, fakes.DefaultFakeConsensusConfig())
	if err != nil {
		return nil, err
	}
	if err := registry.Add(ctx, fakeConsensus); err != nil {
		return nil, err
	}
	caps = append(caps, fakeConsensus)

	fakeConsensusNoDAG := fakes.NewFakeConsensusNoDAG(lggr)
	if err := registry.Add(ctx, consensusserver.NewConsensusServer(fakeConsensusNoDAG)); err != nil {
		return nil, err
	}
	caps = append(caps, fakeConsensusNoDAG)

	writers := []string{"write_aptos-testnet@1.0.0"}
	for _, writer := range writers {
		writeCap := fakes.NewFakeWriteChain(lggr, writer)
		if err := registry.Add(ctx, writeCap); err != nil {
			return nil, err
		}
		caps = append(caps, writeCap)
	}

	return caps, nil
}
