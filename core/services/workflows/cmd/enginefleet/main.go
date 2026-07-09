// Command enginefleet runs one or more instances of the v2 workflow engine
// against a shared capabilities registry populated with in-process capability
// fakes (see core/capabilities/fakes).
//
// It is intended for local experimentation and load/soak testing of the engine:
// point it at a compiled workflow binary and a config file, and it will start N
// engine instances all executing that workflow.
//
// Usage:
//
//	enginefleet -n 5 -w ./testdata/output.wasm -c ./config.yaml
//
// A cron trigger, HTTP action, consensus, and chain-write capability are wired
// up. By default the cron triggers are driven on their schedule so workflows
// actually execute; pass -fire=false to leave the engines subscribed but idle
// (useful for measuring the footprint of suspended engines).
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof" // registers /debug/pprof handlers on the default mux
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jonboulle/clockwork"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	httpserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/http/server"
	evmserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm/server"
	consensusserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/consensus/server"
	cronserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron/server"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	generichost "github.com/smartcontractkit/chainlink-common/pkg/workflows/host"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/fakes"
	cllogger "github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

const (
	// chainWriteTargetID is the capability ID advertised by the fake chain-write
	// capability. Workflows that write reports must target this ID.
	chainWriteTargetID = "write_aptos-testnet@1.0.0"

	// defaultOwner is the (shared) workflow owner. Instances are distinguished by
	// a unique workflow ID, not owner.
	defaultOwner = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	maxUncompressedBinarySize = 1_000_000_000
)

var moduleTimeout = 10 * time.Minute

func main() {
	var (
		numInstances int
		wasmPath     string
		configPath   string
		fire         bool
		debugMode    bool
		pprofAddr    string
		suspend      bool
	)

	flag.IntVar(&numInstances, "n", 1, "Number of workflow engine instances to start")
	flag.StringVar(&wasmPath, "w", "", "Path to the workflow (WASM) binary on disk")
	flag.StringVar(&configPath, "c", "", "Path to the workflow config file")
	flag.BoolVar(&fire, "fire", true, "Drive registered cron triggers on their schedule so workflows execute (set false to leave engines idle)")
	flag.BoolVar(&debugMode, "debug", false, "Enable debug-level logging")
	flag.StringVar(&pprofAddr, "pprof", ":6060", "Address for the pprof/HTTP debug server (empty to disable). Inspect memory via /debug/pprof/heap")
	flag.BoolVar(&suspend, "suspend", false, "Set SuspendOnAwait: executions suspend (freeing the WASM instance) while awaiting capability responses instead of blocking")
	flag.Parse()

	if wasmPath == "" {
		fmt.Println("-w (path to workflow binary) must be set")
		os.Exit(1)
	}
	if numInstances < 1 {
		fmt.Println("-n must be at least 1")
		os.Exit(1)
	}

	binary, err := os.ReadFile(wasmPath)
	if err != nil {
		fmt.Printf("Failed to read workflow binary %q: %v\n", wasmPath, err)
		os.Exit(1)
	}

	var config []byte
	if configPath != "" {
		config, err = os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Failed to read config file %q: %v\n", configPath, err)
			os.Exit(1)
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logLevel := zapcore.InfoLevel
	if debugMode {
		logLevel = zapcore.DebugLevel
	}
	logCfg := cllogger.Config{LogLevel: logLevel}
	lggr, closeLggr := logCfg.New()
	defer func() { _ = closeLggr() }()

	// The engine reads SuspendOnAwait from PerOwner.SuspendOnAwaitEnabled; with a
	// nil settings getter it falls back to this default, so override it here.
	cresettings.Default.PerOwner.SuspendOnAwaitEnabled = settings.Bool(suspend)
	lggr.Infow("SuspendOnAwait configured", "enabled", suspend)

	// pprof/HTTP debug server for inspecting live memory (heap, goroutines, etc.).
	if pprofAddr != "" {
		go func() {
			lggr.Infow("Starting pprof debug server", "addr", pprofAddr)
			if err := http.ListenAndServe(pprofAddr, nil); err != nil { //nolint:gosec // debug server, not exposed to untrusted networks
				lggr.Errorw("pprof server exited", "error", err)
			}
		}()
	}

	// A single shared capabilities registry, populated with fakes and consumed by
	// every engine instance.
	registry := capabilities.NewRegistry(lggr)
	registry.SetLocalRegistry(&capabilities.TestMetadataRegistry{})

	// If the config declares an EVM chain selector, register a fake EVM write
	// capability for it so workflows that write on-chain (via
	// evm:ChainSelector:<sel>@1.0.0) actually exercise the chain-write path.
	evmChainSelector := parseChainSelector(config)

	capServices, cronTrigger, err := registerCapabilities(ctx, lggr, registry, evmChainSelector)
	if err != nil {
		fmt.Printf("Failed to register capabilities: %v\n", err)
		os.Exit(1)
	}

	// Track everything we start so we can shut it down cleanly, engines first.
	var engineServices []services.Service
	for i := 0; i < numInstances; i++ {
		instanceLggr := commonlogger.Named(lggr, fmt.Sprintf("engine-%d", i))
		// Each instance gets a unique workflow ID so their trigger registrations
		// and execution IDs don't collide in the shared capabilities.
		workflowID := fmt.Sprintf("%064x", i+1)
		workflowName := fmt.Sprintf("enginefleet-%d", i)

		engine, engErr := buildEngine(ctx, instanceLggr, registry, binary, config, workflowID, workflowName)
		if engErr != nil {
			fmt.Printf("Failed to create engine instance %d: %v\n", i, engErr)
			shutdown(lggr, engineServices, capServices)
			os.Exit(1)
		}

		if startErr := engine.Start(ctx); startErr != nil {
			fmt.Printf("Failed to start engine instance %d: %v\n", i, startErr)
			shutdown(lggr, engineServices, capServices)
			os.Exit(1)
		}
		engineServices = append(engineServices, engine)
		lggr.Infow("Started workflow engine instance", "instance", i, "workflowID", workflowID, "workflowName", workflowName)
	}

	if fire {
		// Give engines a moment to finish registering their triggers, then drive
		// the cron triggers on their schedule.
		startCronDriver(ctx, lggr, cronTrigger)
	}

	lggr.Infow("Workflow engine fleet running", "instances", numInstances, "firingCronTriggers", fire, "pid", os.Getpid())

	<-ctx.Done()

	lggr.Infow("Shutting down workflow engine fleet")
	shutdown(lggr, engineServices, capServices)
}

// registerCapabilities wires the fake capabilities into the shared registry and
// starts them, returning the started services so they can be closed on shutdown
// and the cron trigger so it can be driven.
//
// The DirectHTTPAction and cron trigger are typed "server" capabilities, so they
// are adapted into generic capabilities via httpserver.NewClientServer /
// cronserver.NewCronServer before being added; the consensus and chain-write
// fakes implement ExecutableCapability directly and are added as-is.
func registerCapabilities(ctx context.Context, lggr commonlogger.Logger, registry *capabilities.Registry, evmChainSelector uint64) ([]services.Service, *fakes.ManualCronTriggerService, error) {
	var started []services.Service

	// Cron trigger capability.
	cronTrigger, err := fakes.NewManualCronTriggerService(lggr)
	if err != nil {
		return started, nil, fmt.Errorf("failed to create cron trigger capability: %w", err)
	}
	if err = registry.Add(ctx, cronserver.NewCronServer(cronTrigger)); err != nil {
		return started, nil, fmt.Errorf("failed to add cron trigger capability: %w", err)
	}
	if err = cronTrigger.Start(ctx); err != nil {
		return started, nil, fmt.Errorf("failed to start cron trigger capability: %w", err)
	}
	started = append(started, cronTrigger)

	// HTTP action capability.
	httpAction := fakes.NewDirectHTTPAction(lggr)
	if err = registry.Add(ctx, httpserver.NewClientServer(httpAction)); err != nil {
		return started, nil, fmt.Errorf("failed to add http action capability: %w", err)
	}
	if err = httpAction.Start(ctx); err != nil {
		return started, nil, fmt.Errorf("failed to start http action capability: %w", err)
	}
	started = append(started, httpAction)

	// Consensus capability. v2 (NoDAG) workflows use consensus@1.0.0-alpha,
	// served by the NoDAG consensus fake wrapped in a consensus server.
	nSigners := 4
	signers := make([]ocr2key.KeyBundle, nSigners)
	for i := range nSigners {
		signers[i] = ocr2key.MustNewInsecure(fakes.SeedForKeys(), corekeys.EVM)
	}
	consensus := fakes.NewFakeConsensusNoDAG(signers, lggr)
	if err = registry.Add(ctx, consensusserver.NewConsensusServer(consensus)); err != nil {
		return started, nil, fmt.Errorf("failed to add consensus capability: %w", err)
	}
	if err = consensus.Start(ctx); err != nil {
		return started, nil, fmt.Errorf("failed to start consensus capability: %w", err)
	}
	started = append(started, consensus)

	// Chain-write (target) capability.
	writeChain := fakes.NewFakeWriteChain(lggr, chainWriteTargetID)
	if err = registry.Add(ctx, writeChain); err != nil {
		return started, nil, fmt.Errorf("failed to add chain write capability: %w", err)
	}
	if err = writeChain.Start(ctx); err != nil {
		return started, nil, fmt.Errorf("failed to start chain write capability: %w", err)
	}
	started = append(started, writeChain)

	// EVM write capability (only if the workflow config declares a chain
	// selector). Registered under evm:ChainSelector:<sel>@1.0.0 via the EVM
	// server wrapper.
	if evmChainSelector > 0 {
		evmWrite := newFakeEVMWrite(lggr, evmChainSelector)
		if err = registry.Add(ctx, evmserver.NewClientServer(evmWrite)); err != nil {
			return started, nil, fmt.Errorf("failed to add evm write capability: %w", err)
		}
		started = append(started, evmWrite)
		lggr.Infow("Registered fake EVM write capability", "chainSelector", evmChainSelector)
	}

	return started, cronTrigger, nil
}

// parseChainSelector extracts the "chainSelector" field from a JSON workflow
// config. It decodes with json.Number to preserve the full uint64 range (chain
// selectors exceed float64's exact-integer range). Returns 0 if absent or if
// the config isn't JSON.
func parseChainSelector(config []byte) uint64 {
	if len(config) == 0 {
		return 0
	}
	dec := json.NewDecoder(bytes.NewReader(config))
	dec.UseNumber()
	var m map[string]any
	if err := dec.Decode(&m); err != nil {
		return 0
	}
	num, ok := m["chainSelector"].(json.Number)
	if !ok {
		return 0
	}
	sel, err := strconv.ParseUint(num.String(), 10, 64)
	if err != nil {
		return 0
	}
	return sel
}

// noopSecrets is a SecretsFetcher that returns no secrets. The scenario
// workflows don't use secrets; this satisfies the optional dependency.
type noopSecrets struct{}

func (noopSecrets) GetSecrets(_ context.Context, _ *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	return nil, nil
}

// requirementSelectingModule wraps a module so it satisfies the engine's
// RequirementEnforcingModule expectations without a TEE requirements hook.
type requirementModule struct {
	generichost.Module
}

func (requirementModule) SetRequirements(string, *sdkpb.Requirements) {}

// buildEngine compiles the workflow module and constructs a v2 engine for a
// single instance, identified by workflowID/workflowName. It mirrors the
// standalone-engine wiring (in-memory store, local time provider, no billing)
// but is parameterized so each instance is independent.
func buildEngine(
	ctx context.Context,
	lggr commonlogger.Logger,
	registry *capabilities.Registry,
	binary, config []byte,
	workflowID, workflowName string,
) (services.Service, error) {
	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: defaultOwner, Workflow: workflowID})

	moduleConfig := &host.ModuleConfig{
		Logger:                  lggr,
		Labeler:                 custmsg.NewLabeler(),
		MaxCompressedBinarySize: maxUncompressedBinarySize,
		IsUncompressed:          true,
		Timeout:                 &moduleTimeout,
	}

	mainModule, err := host.NewModule(ctx, moduleConfig, binary, host.WithDeterminism())
	if err != nil {
		return nil, fmt.Errorf("unable to create module: %w", err)
	}

	handler := generichost.RequirementsHandler{Tee: func(context.Context, *sdkpb.Tee) bool { return true }}
	module := generichost.NewRequirementSelectingModule(
		generichost.ModuleAndHandler{Module: requirementModule{Module: mainModule}, RequirementsHandler: handler},
		[]generichost.ModuleAndHandler{},
	)

	name, err := types.NewWorkflowName(workflowName)
	if err != nil {
		return nil, err
	}

	// Allow chain access for all chain selectors so on-chain writes aren't gated.
	allowAllChains := func(cfg *cresettings.Workflows) {
		cfg.ChainAllowed = settings.PerChainSelector(settings.Bool(true), map[string]bool{})
	}

	lf := limits.Factory{Logger: commonlogger.Named(lggr, "Limits")}
	limiters, err := v2.NewLimiters(lf, allowAllChains)
	if err != nil {
		return nil, err
	}
	moduleConfig.EnableUserMetricsLimiter = limiters.UserMetricEnabled
	moduleConfig.MaxUserMetricPayloadLimiter = limiters.UserMetricPayload
	moduleConfig.MaxUserMetricNameLengthLimiter = limiters.UserMetricNameLength
	moduleConfig.MaxUserMetricLabelsPerMetricLimiter = limiters.UserMetricLabelsPerMetric
	moduleConfig.MaxUserMetricLabelValueLengthLimiter = limiters.UserMetricLabelValueLength

	featureFlags, err := v2.NewFeatureFlags(lf, allowAllChains)
	if err != nil {
		return nil, err
	}
	workflowLimits, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{
		Global:   1_000_000_000,
		PerOwner: 1_000_000_000,
	}, lf)
	if err != nil {
		return nil, err
	}

	cfg := &v2.EngineConfig{
		Lggr:                 lggr,
		Module:               module,
		WorkflowConfig:       config,
		CapRegistry:          registry,
		DonSubscriber:        mockSubscriber{},
		UseLocalTimeProvider: true,
		ExecutionsStore:      store.NewInMemoryStore(lggr, clockwork.NewRealClock()),

		WorkflowID:    workflowID,
		WorkflowOwner: defaultOwner,
		WorkflowName:  name,
		WorkflowTag:   "workflowTag",

		LocalLimits:         v2.EngineLimits{},
		LocalLimiters:       limiters,
		FeatureFlags:        featureFlags,
		GlobalWorkflowLimit: workflowLimits,

		BeholderEmitter: custmsg.NewLabeler(),
		SecretsFetcher:  noopSecrets{},
	}

	return v2.NewEngine(cfg)
}

type mockSubscriber struct{}

func (mockSubscriber) Subscribe(context.Context) (<-chan commoncap.DON, func(), error) {
	return make(<-chan commoncap.DON), func() {}, nil
}

// startCronDriver fires every registered cron trigger on its schedule. Engines
// register their triggers asynchronously after Start, so it re-scans
// periodically and spawns a driver goroutine for each newly-seen trigger.
//
// Each ManualTrigger call blocks until the trigger's next scheduled tick and
// then delivers one event, so looping per trigger reproduces the cron cadence.
func startCronDriver(ctx context.Context, lggr commonlogger.Logger, cron *fakes.ManualCronTriggerService) {
	go func() {
		seen := make(map[string]bool)
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			for _, id := range cron.TriggerIDs() {
				if seen[id] {
					continue
				}
				seen[id] = true
				lggr.Infow("Driving cron trigger", "triggerID", id, "totalDriven", len(seen))
				go driveTrigger(ctx, lggr, cron, id)
			}
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
}

func driveTrigger(ctx context.Context, lggr commonlogger.Logger, cron *fakes.ManualCronTriggerService, id string) {
	for ctx.Err() == nil {
		if err := cron.ManualTrigger(ctx, id, nil); err != nil {
			if ctx.Err() != nil {
				return
			}
			lggr.Debugw("cron fire error", "triggerID", id, "err", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
		}
	}
}

// shutdown closes engines before the capabilities they depend on.
func shutdown(lggr commonlogger.Logger, engines, caps []services.Service) {
	for _, e := range engines {
		if err := e.Close(); err != nil {
			lggr.Errorw("Failed to close engine", "error", err)
		}
	}
	for _, c := range caps {
		if err := c.Close(); err != nil {
			lggr.Errorw("Failed to close capability", "name", c.Name(), "error", err)
		}
	}
}
