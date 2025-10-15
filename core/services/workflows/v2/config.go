package v2

import (
	"errors"
	"fmt"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/services"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

type EngineConfig struct {
	Lggr                 logger.Logger
	Module               host.ModuleV2
	WorkflowConfig       []byte // workflow author provided config
	CapRegistry          core.CapabilitiesRegistry
	DonTimeStore         *dontime.Store
	UseLocalTimeProvider bool // Set true when DON Time Plugin is not running
	ExecutionsStore      store.Store
	Clock                clockwork.Clock
	SecretsFetcher       SecretsFetcher

	WorkflowID            string // hex-encoded [32]byte, no "0x" prefix
	WorkflowOwner         string // hex-encoded [20]byte, no "0x" prefix
	WorkflowName          types.WorkflowName
	WorkflowTag           string // workflow tag is required during workflow registration. owner + name + tag uniquely identifies a workflow.
	WorkflowEncryptionKey workflowkey.Key

	LocalLimits                       EngineLimits
	LocalLimiters                     *EngineLimiters
	GlobalExecutionConcurrencyLimiter limits.ResourceLimiter[int] // global + per owner
	GlobalExecutionRateLimiter        limits.RateLimiter          // global + per owner

	BeholderEmitter custmsg.MessageEmitter

	Hooks         LifecycleHooks
	BillingClient metering.BillingClient

	// WorkflowRegistryAddress is the address of the workflow registry contract
	WorkflowRegistryAddress string
	// WorkflowRegistryChainSelector is the chain selector for the workflow registry
	WorkflowRegistryChainSelector string

	// OrgResolver is used to resolve organization IDs from workflow owners
	OrgResolver orgresolver.OrgResolver

	// includes additional logging of events internal to user workflows
	DebugMode bool
}

type EngineLimiters struct {
	ExecutionResponse        limits.BoundLimiter[config.Size]
	TriggerSubscriptionTime  limits.TimeLimiter
	TriggerRegistrationsTime limits.TimeLimiter
	TriggerSubscription      limits.BoundLimiter[int]
	TriggerEventQueue        limits.QueueLimiter[enqueuedTriggerEvent]
	TriggerEventQueueTime    limits.TimeLimiter
	ExecutionConcurrency     limits.ResourcePoolLimiter[int]

	CapabilityConcurrency limits.ResourcePoolLimiter[int]
	SecretsConcurrency    limits.ResourcePoolLimiter[int]
	ExecutionTime         limits.TimeLimiter
	CapabilityCallTime    limits.TimeLimiter
	LogEvent              limits.BoundLimiter[int]
	LogLine               limits.BoundLimiter[config.Size]

	ChainWriteTargets limits.BoundLimiter[int]
	ChainReadCalls    limits.BoundLimiter[int]
	HTTPActionCalls   limits.BoundLimiter[int]
}

// NewLimiters returns a new set of EngineLimiters based on the default configuration, and optionally modified by cfgFn.
func NewLimiters(lf limits.Factory, cfgFn func(*cresettings.Workflows)) (*EngineLimiters, error) {
	l := &EngineLimiters{}
	err := l.init(lf, cfgFn)
	return l, err
}

func (l *EngineLimiters) init(lf limits.Factory, cfgFn func(*cresettings.Workflows)) (err error) {
	cfg := cresettings.Default.PerWorkflow // make copy
	if cfgFn != nil {
		cfgFn(&cfg)
	}
	l.ExecutionResponse, err = limits.MakeBoundLimiter(lf, cfg.ExecutionResponseLimit)
	if err != nil {
		return
	}
	l.TriggerSubscriptionTime, err = lf.MakeTimeLimiter(cfg.TriggerSubscriptionTimeout)
	if err != nil {
		return
	}
	l.TriggerRegistrationsTime, err = lf.MakeTimeLimiter(cfg.TriggerRegistrationsTimeout)
	if err != nil {
		return
	}
	l.TriggerSubscription, err = limits.MakeBoundLimiter(lf, cfg.TriggerSubscriptionLimit)
	if err != nil {
		return
	}
	l.TriggerEventQueue, err = limits.MakeQueueLimiter[enqueuedTriggerEvent](lf, cfg.TriggerEventQueueLimit)
	if err != nil {
		return
	}
	l.TriggerEventQueueTime, err = lf.MakeTimeLimiter(cfg.TriggerEventQueueTimeout)
	if err != nil {
		return
	}
	l.ExecutionConcurrency, err = limits.MakeResourcePoolLimiter(lf, cfg.ExecutionConcurrencyLimit)
	if err != nil {
		return
	}

	l.CapabilityConcurrency, err = limits.MakeResourcePoolLimiter(lf, cfg.CapabilityConcurrencyLimit)
	if err != nil {
		return
	}
	l.SecretsConcurrency, err = limits.MakeResourcePoolLimiter(lf, cfg.SecretsConcurrencyLimit)
	if err != nil {
		return
	}
	l.ExecutionTime, err = lf.MakeTimeLimiter(cfg.ExecutionTimeout)
	if err != nil {
		return
	}
	l.CapabilityCallTime, err = lf.MakeTimeLimiter(cfg.CapabilityCallTimeout)
	if err != nil {
		return
	}
	l.LogEvent, err = limits.MakeBoundLimiter(lf, cfg.LogEventLimit)
	if err != nil {
		return
	}
	l.LogLine, err = limits.MakeBoundLimiter(lf, cfg.LogLineLimit)
	if err != nil {
		return
	}
	l.ChainWriteTargets, err = limits.MakeBoundLimiter(lf, cfg.ChainWrite.TargetsLimit)
	if err != nil {
		return
	}
	l.ChainReadCalls, err = limits.MakeBoundLimiter(lf, cfg.ChainRead.CallLimit)
	if err != nil {
		return
	}
	l.HTTPActionCalls, err = limits.MakeBoundLimiter(lf, cfg.HTTPAction.CallLimit)
	return
}

func (l *EngineLimiters) Close() error {
	return services.CloseAll(
		l.ExecutionResponse,
		l.TriggerSubscriptionTime,
		l.TriggerRegistrationsTime,
		l.TriggerSubscription,
		l.TriggerEventQueue,
		l.TriggerEventQueueTime,
		l.ExecutionConcurrency,
		l.CapabilityConcurrency,
		l.SecretsConcurrency,
		l.ExecutionTime,
		l.CapabilityCallTime,
		l.LogEvent,
		l.LogLine,
		l.ChainWriteTargets,
		l.ChainReadCalls,
		l.HTTPActionCalls,
	)
}

const (
	defaultHeartbeatFrequencyMs = 1000 * 60 // 1 minute
	defaultShutdownTimeoutMs    = 5000
)

type EngineLimits struct {
	HeartbeatFrequencyMs uint32
	ShutdownTimeoutMs    uint32
}

type LifecycleHooks struct {
	OnInitialized          func(err error)
	OnSubscribedToTriggers func(triggerIDs []string)
	OnExecutionFinished    func(executionID string, status string)
	OnExecutionError       func(msg string)
	OnResultReceived       func(*sdkpb.ExecutionResult)
	OnRateLimited          func(executionID string)
}

func (c *EngineConfig) Validate() error {
	if c.Lggr == nil {
		return errors.New("logger not set")
	}
	if c.Module == nil {
		return errors.New("module not set")
	}
	if c.CapRegistry == nil {
		return errors.New("capabilities registry not set")
	}
	if c.DonTimeStore == nil && !c.UseLocalTimeProvider {
		return errors.New("dontime store not set")
	}
	if c.ExecutionsStore == nil {
		return errors.New("executions store not set")
	}
	if c.Clock == nil {
		c.Clock = clockwork.NewRealClock()
	}

	_, err := types.WorkflowIDFromHex(c.WorkflowID)
	if err != nil {
		return fmt.Errorf("invalid workflowID: %w", err)
	}
	err = types.ValidateWorkflowOwner(c.WorkflowOwner)
	if err != nil {
		return fmt.Errorf("invalid workflowOwner: %w", err)
	}
	if c.WorkflowName == nil {
		return errors.New("workflowName not set")
	}

	c.LocalLimits.setDefaultLimits()
	if c.GlobalExecutionConcurrencyLimiter == nil {
		return errors.New("execution concurrency limiter not set")
	}
	if c.GlobalExecutionRateLimiter == nil {
		return errors.New("execution rate limiter not set")
	}

	if c.BeholderEmitter == nil {
		return errors.New("beholder emitter not set")
	}

	c.Hooks.setDefaultHooks()
	return nil
}

func (l *EngineLimits) setDefaultLimits() {
	if l.HeartbeatFrequencyMs == 0 {
		l.HeartbeatFrequencyMs = defaultHeartbeatFrequencyMs
	}
	if l.ShutdownTimeoutMs == 0 {
		l.ShutdownTimeoutMs = defaultShutdownTimeoutMs
	}
}

// set all to non-nil so the Engine doesn't have to check before each call
func (h *LifecycleHooks) setDefaultHooks() {
	if h.OnInitialized == nil {
		h.OnInitialized = func(err error) {}
	}
	if h.OnSubscribedToTriggers == nil {
		h.OnSubscribedToTriggers = func(triggerIDs []string) {}
	}
	if h.OnResultReceived == nil {
		h.OnResultReceived = func(res *sdkpb.ExecutionResult) {}
	}
	if h.OnExecutionError == nil {
		h.OnExecutionError = func(msg string) {}
	}
	if h.OnExecutionFinished == nil {
		h.OnExecutionFinished = func(executionID string, status string) {}
	}
	if h.OnRateLimited == nil {
		h.OnRateLimited = func(executionID string) {}
	}
}
