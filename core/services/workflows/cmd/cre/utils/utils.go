package utils

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/billing"
	httpserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/http/server"
	evmserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm/server"
	consensusserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/consensus/server"
	crontrigger "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron/server"
	httptrigger "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/http/server"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
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

type ManualTriggers struct {
	ManualCronTrigger *fakes.ManualCronTriggerService
	ManualHTTPTrigger *fakes.ManualHTTPTriggerService
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

type StandaloneEngineLoggerConfig struct {
	ModuleLogger         logger.Logger
	WorkflowLimitsLogger logger.Logger
	EngineLogger         logger.Logger
}

func NewStandaloneEngine(
	ctx context.Context,
	lggrCfg StandaloneEngineLoggerConfig,
	registry *capabilities.Registry,
	binary, config []byte,
	billingClientAddr string,
	lifecycleHooks v2.LifecycleHooks,
) (services.Service, *sdkpb.TriggerSubscriptionRequest, error) {
	labeler := custmsg.NewLabeler()
	moduleConfig := &host.ModuleConfig{
		Logger:                  lggrCfg.ModuleLogger,
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
	// for _, sub := range result.GetTriggerSubscriptions().Subscriptions {
	// 	fmt.Println("Registered Trigger: ", sub.Id)
	// 	fmt.Println("Method: ", sub.Method)
	// }

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

	workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggrCfg.WorkflowLimitsLogger, syncerlimiter.Config{
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
			Lggr:                 lggrCfg.EngineLogger,
			Workflow:             *sdkSpec,
			WorkflowID:           defaultWorkflowID,
			WorkflowOwner:        defaultOwner,
			WorkflowName:         name,
			Registry:             registry,
			Store:                store.NewInMemoryStore(lggrCfg.EngineLogger, clockwork.NewRealClock()),
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
		Lggr:            lggrCfg.EngineLogger,
		Module:          module,
		WorkflowConfig:  config,
		CapRegistry:     registry,
		ExecutionsStore: store.NewInMemoryStore(lggrCfg.EngineLogger, clockwork.NewRealClock()),

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

// NewFakeCapabilities builds faked capabilities, then registers them with the capability registry.
func NewFakeComputeCapabilities(ctx context.Context, lggr logger.Logger, registry *capabilities.Registry) ([]services.Service, error) {
	caps := make([]services.Service, 0)

	// EVM
	// evmClient, err := ethclient.Dial("https://sepolia.infura.io/v3/dbe1bfd45172477084dfe080e0754c1e")
	evmClient, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		return nil, err
	}

	// TODO: get private key from env var
	privateKey, err := crypto.HexToECDSA("bc2c4e2ed2af93035d2e616de9d8dfb6fc35117cf2bf00bf7f7a7ab54981536d")
	if err != nil {
		return nil, err
	}

	evm := fakes.NewFakeEvmChain(lggr, evmClient, privateKey)
	evmServer := evmserver.NewClientServer(evm)
	if err := registry.Add(ctx, evmServer); err != nil {
		return nil, err
	}
	caps = append(caps, evm)

	// Consensus
	consensus, err := fakes.NewFakeConsensus(lggr, fakes.DefaultFakeConsensusConfig())
	consensusServer := consensusserver.NewConsensusServer(consensus)
	if err := registry.Add(ctx, consensusServer); err != nil {
		return nil, err
	}
	caps = append(caps, consensusServer)

	// HTTP Action
	httpAction := fakes.NewDirectHTTPAction(lggr)
	httpActionServer := httpserver.NewClientServer(httpAction)
	if err := registry.Add(ctx, httpActionServer); err != nil {
		return nil, err
	}
	caps = append(caps, httpActionServer)

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

func NewManualTriggerCapabilities(ctx context.Context, lggr logger.Logger, registry *capabilities.Registry) (ManualTriggers, error) {
	// Cron
	manualCronTrigger := fakes.NewManualCronTriggerService(lggr)
	manualCronTriggerServer := crontrigger.NewCronServer(manualCronTrigger)
	if err := registry.Add(ctx, manualCronTriggerServer); err != nil {
		return ManualTriggers{}, err
	}

	// HTTP
	manualHTTPTrigger := fakes.NewManualHTTPTriggerService(lggr)
	manualHTTPTriggerServer := httptrigger.NewHTTPServer(manualHTTPTrigger)
	if err := registry.Add(ctx, manualHTTPTriggerServer); err != nil {
		return ManualTriggers{}, err
	}

	return ManualTriggers{
		ManualCronTrigger: manualCronTrigger,
		ManualHTTPTrigger: manualHTTPTrigger,
	}, nil
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
			&fakes.RelayerSetMock{}, &fakes.OracleFactoryMock{}, &fakes.GatewayConnectorMock{})

		service := &standaloneLoopWrapper{
			StandardCapabilities: loop,
		}
		caps = append(caps, service)
	}

	return caps
}
