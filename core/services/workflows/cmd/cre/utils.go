package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/billing"
	evmserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm/server"
	consensusserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/consensus/server"
	crontrigger "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron/server"
	httptrigger "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/http/server"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/fakes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

const (
	defaultMaxUncompressedBinarySize = 1000000000
	defaultRPS                       = 1000.0
	defaultBurst                     = 1000
	defaultWorkflowID                = "1111111111111111111111111111111111111111111111111111111111111111"
	defaultOwner                     = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	defaultName                      = "myworkflow"
)

type standardCapConfig struct {
	Config string

	// Set enabled to true to run the loop plugin.  Requires the plugin be installed.
	// Config will be passed to Initialise method of plugin.
	Enabled bool
}

var (
	goBinPath            = os.Getenv("GOBIN")
	standardCapabilities = map[string]standardCapConfig{
		"cron": {
			Config:  `{"fastestScheduleIntervalSeconds": 1}`,
			Enabled: true,
		},
		"readcontract":  {},
		"kvstore":       {},
		"workflowevent": {},
	}
)

func NewStandaloneEngine(
	ctx context.Context,
	lggr logger.Logger,
	registry *capabilities.Registry,
	binary []byte, config []byte,
	billingClientAddr string,
	lifecycleHooks v2.LifecycleHooks,
) (services.Service, error) {
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

	name, err := types.NewWorkflowName(defaultName)
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

	var billingClient billing.WorkflowClient
	if billingClientAddr != "" {
		billingClient, _ = billing.NewWorkflowClient(billingClientAddr)
	}

	if module.IsLegacyDAG() {
		sdkSpec, err := host.GetWorkflowSpec(ctx, moduleConfig, binary, config)
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
			BillingClient:        billingClient,
		}
		return workflows.NewEngine(ctx, cfg)
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
	}

	return v2.NewEngine(ctx, cfg)
}

// TODO support fetching secrets (from a local file)
func SecretsFor(ctx context.Context, workflowOwner, hexWorkflowName, decodedWorkflowName, workflowID string) (map[string]string, error) {
	return map[string]string{}, nil
}

func NewFakeManualTriggerCapabilities(ctx context.Context, lggr logger.Logger, registry *capabilities.Registry) ([]fakes.ManualTriggerCapability, error) {
	caps := make([]fakes.ManualTriggerCapability, 0)

	// Cron
	manualCronTrigger := fakes.NewFakeCronTriggerService(lggr)
	manualCronTriggerServer := crontrigger.NewCronServer(manualCronTrigger)
	if err := registry.Add(ctx, manualCronTriggerServer); err != nil {
		return nil, err
	}
	caps = append(caps, manualCronTrigger)

	// HTTP
	manualHttpTrigger := fakes.NewFakeManualHttpTriggerService(lggr)
	manualHttpTriggerServer := httptrigger.NewHTTPServer(manualHttpTrigger)
	if err := registry.Add(ctx, manualHttpTriggerServer); err != nil {
		return nil, err
	}
	caps = append(caps, manualHttpTrigger)

	return caps, nil
}

func NewFakeComputeCapabilities(ctx context.Context, lggr logger.Logger, registry *capabilities.Registry) ([]services.Service, error) {
	caps := make([]services.Service, 0)

	// EVM
	evm := fakes.NewFakeEvmChain(lggr, nil)
	evmServer := evmserver.NewClientServer(evm)
	if err := registry.Add(ctx, evmServer); err != nil {
		return nil, err
	}
	caps = append(caps, evm)

	// Consensus
	consensus, err := fakes.NewFakeConsensus(lggr, fakes.DefaultFakeConsensusConfig())
	consensusServer := consensusserver.NewConsensusServer(consensus)
	if err != nil {
		return nil, err
	}

	if err := registry.Add(ctx, consensusServer); err != nil {
		return nil, err
	}
	caps = append(caps, consensus)

	return caps, nil
}

func NewFakeCapabilities(ctx context.Context, lggr logger.Logger, registry *capabilities.Registry) ([]services.Service, error) {
	caps := make([]services.Service, 0)

	caps = append(caps, newStandardCapabilities(standardCapabilities, lggr, registry)...)

	fakeConsensus, err := fakes.NewFakeConsensus(lggr, fakes.DefaultFakeConsensusConfig())
	if err != nil {
		return nil, err
	}
	if err := registry.Add(ctx, fakeConsensus); err != nil {
		return nil, err
	}
	caps = append(caps, fakeConsensus)

	return caps, nil
}

// standaloneLoopWrapper wraps a StandardCapabilities to implement services.Service
type standaloneLoopWrapper struct {
	*standardcapabilities.StandardCapabilities
}

func (l *standaloneLoopWrapper) Ready() error { return l.StandardCapabilities.Ready() }

func (l *standaloneLoopWrapper) HealthReport() map[string]error { return make(map[string]error) }

func (l *standaloneLoopWrapper) Name() string { return "wrapped" }

func newStandardCapabilities(
	standardCapabilities map[string]standardCapConfig,
	lggr logger.Logger,
	registry *capabilities.Registry,
) []services.Service {
	caps := make([]services.Service, 0)

	pluginRegistrar := plugins.NewRegistrarConfig(
		loop.GRPCOpts{},
		func(name string) (*plugins.RegisteredLoop, error) { return &plugins.RegisteredLoop{}, nil },
		func(loopId string) {})

	for name, config := range standardCapabilities {
		if !config.Enabled {
			continue
		}

		spec := &job.StandardCapabilitiesSpec{
			Command: path.Join(goBinPath, name),
			Config:  config.Config,
		}

		loop := standardcapabilities.NewStandardCapabilities(lggr, spec,
			pluginRegistrar, &fakes.TelemetryServiceMock{}, &fakes.KVStoreMock{},
			registry, &fakes.ErrorLogMock{}, &fakes.PipelineRunnerServiceMock{},
			&fakes.RelayerSetMock{}, &fakes.OracleFactoryMock{})

		service := &standaloneLoopWrapper{
			StandardCapabilities: loop,
		}
		caps = append(caps, service)
	}

	return caps
}
