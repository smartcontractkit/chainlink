package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/fakes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

const (
	defaultMaxUncompressedBinarySize = 1000000000
	defaultRPS                       = 1000.0
	defaultBurst                     = 1000
	defaultWorkflowID                = "1111111111111111111111111111111111111111111111111111111111111111"
	defaultOwner                     = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	defaultName                      = "myworkflow"
)

func NewStandaloneEngine(ctx context.Context, lggr logger.Logger, registry *capabilities.Registry, binary []byte, config []byte) (*workflows.Engine, error) {
	labeler := custmsg.NewLabeler()
	moduleConfig := &host.ModuleConfig{
		Logger:                  lggr,
		Labeler:                 labeler,
		MaxCompressedBinarySize: defaultMaxUncompressedBinarySize,
		IsUncompressed:          true,
	}

	module, err := host.NewModule(moduleConfig, binary, host.WithDeterminism())
	if err != nil {
		return nil, fmt.Errorf("unable to create module from config: %w", err)
	}

	if !module.IsLegacyDAG() {
		return nil, errors.New("no DAG not yet supported")
	}

	sdkSpec, err := host.GetWorkflowSpec(ctx, moduleConfig, binary, config)
	if err != nil {
		return nil, err
	}

	rl, err := ratelimiter.NewRateLimiter(ratelimiter.Config{
		GlobalRPS:      defaultRPS,
		GlobalBurst:    defaultBurst,
		PerSenderRPS:   defaultRPS,
		PerSenderBurst: defaultBurst,
	})
	if err != nil {
		return nil, err
	}

	workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{
		Global:   1000000000,
		PerOwner: 1000000000,
	})
	if err != nil {
		return nil, err
	}

	name, err := types.NewWorkflowName(defaultName)
	if err != nil {
		return nil, err
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
	}
	return workflows.NewEngine(ctx, cfg)
}

// TODO support fetching secrets (from a local file)
func SecretsFor(ctx context.Context, workflowOwner, hexWorkflowName, decodedWorkflowName, workflowID string) (map[string]string, error) {
	return map[string]string{}, nil
}

func NewFakeCapabilities(ctx context.Context, lggr logger.Logger, registry *capabilities.Registry) ([]services.Service, error) {
	caps := make([]services.Service, 0, 5)

	// Streams Trigger
	streamsTrigger := fakes.NewFakeStreamsTrigger(lggr, 6)
	err := registry.Add(ctx, streamsTrigger)
	if err != nil {
		return nil, err
	}
	caps = append(caps, streamsTrigger)

	// Consensus
	fakeConsensus, err := fakes.NewFakeConsensus(lggr, fakes.DefaultFakeConsensusConfig())
	if err != nil {
		return nil, err
	}
	err = registry.Add(ctx, fakeConsensus)
	if err != nil {
		return nil, err
	}
	caps = append(caps, fakeConsensus)

	// Chain Writers
	writers := []string{"write_aptos-testnet@1.0.0"}
	for _, writer := range writers {
		writeCap := fakes.NewFakeWriteChain(lggr, writer)
		err = registry.Add(ctx, writeCap)
		if err != nil {
			return nil, err
		}
		caps = append(caps, writeCap)
	}

	return caps, nil
}
