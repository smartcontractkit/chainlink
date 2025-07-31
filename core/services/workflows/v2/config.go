package v2

import (
	"errors"
	"fmt"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

type EngineConfig struct {
	Lggr            logger.Logger
	Module          host.ModuleV2
	WorkflowConfig  []byte // workflow author provided config
	CapRegistry     core.CapabilitiesRegistry
	ExecutionsStore store.Store
	Clock           clockwork.Clock
	SecretsFetcher  SecretsFetcher

	WorkflowID    string // hex-encoded [32]byte, no "0x" prefix
	WorkflowOwner string // hex-encoded [20]byte, no "0x" prefix
	WorkflowName  types.WorkflowName

	LocalLimits          EngineLimits                // local to a single workflow
	GlobalLimits         limits.ResourceLimiter[int] // global to all workflows
	ExecutionRateLimiter limits.RateLimiter          // global + per owner

	BeholderEmitter custmsg.MessageEmitter

	Hooks         LifecycleHooks
	BillingClient metering.BillingClient

	// WorkflowRegistryAddress is the address of the workflow registry contract
	WorkflowRegistryAddress string
	// WorkflowRegistryChainSelector is the chain selector for the workflow registry
	WorkflowRegistryChainSelector string

	// includes additional logging of events internal to user workflows
	DebugMode bool
}

const (
	defaultModuleExecuteMaxResponseSizeBytes   = 100000
	defaultTriggerSubscriptionRequestTimeoutMs = 500
	defaultTriggerAllRegistrationsTimeoutMs    = 1000
	defaultMaxTriggerSubscriptions             = 10
	defaultTriggerEventQueueSize               = 1000
	defaultTriggerEventMaxAgeMs                = 1000 * 60 * 10 // 10 minutes

	defaultMaxConcurrentWorkflowExecutions         = 100
	defaultMaxConcurrentCapabilityCallsPerWorkflow = 10
	defaultMaxConcurrentSecretsCallsPerWorkflow    = 3
	defaultWorkflowExecutionTimeoutMs              = 1000 * 60 * 10 // 10 minutes
	defaultCapabilityCallTimeoutMs                 = 1000 * 60 * 8  // 8 minutes
	defaultMaxUserLogEventsPerExecution            = 1000
	defaultMaxUserLogLineLength                    = 1000

	defaultHeartbeatFrequencyMs = 1000 * 60 // 1 minute
	defaultShutdownTimeoutMs    = 5000
)

type EngineLimits struct {
	ModuleExecuteMaxResponseSizeBytes   uint32
	TriggerSubscriptionRequestTimeoutMs uint32
	TriggerAllRegistrationsTimeoutMs    uint32
	MaxTriggerSubscriptions             uint16
	TriggerEventQueueSize               uint16
	TriggerEventMaxAgeMs                uint32

	MaxConcurrentWorkflowExecutions         uint16
	MaxConcurrentCapabilityCallsPerWorkflow uint16
	MaxConcurrentSecretsCallsPerWorkflow    uint16
	WorkflowExecutionTimeoutMs              uint32
	CapabilityCallTimeoutMs                 uint32
	MaxUserLogEventsPerExecution            uint32
	MaxUserLogLineLength                    uint32

	HeartbeatFrequencyMs uint32
	ShutdownTimeoutMs    uint32
}

type LifecycleHooks struct {
	OnInitialized          func(err error)
	OnSubscribedToTriggers func(triggerIDs []string)
	OnExecutionFinished    func(executionID string, status string)
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
	if c.GlobalLimits == nil {
		return errors.New("global limits not set")
	}
	if c.ExecutionRateLimiter == nil {
		return errors.New("execution rate limiter not set")
	}

	if c.BeholderEmitter == nil {
		return errors.New("beholder emitter not set")
	}

	c.Hooks.setDefaultHooks()
	return nil
}

func (l *EngineLimits) setDefaultLimits() {
	if l.ModuleExecuteMaxResponseSizeBytes == 0 {
		l.ModuleExecuteMaxResponseSizeBytes = defaultModuleExecuteMaxResponseSizeBytes
	}
	if l.TriggerSubscriptionRequestTimeoutMs == 0 {
		l.TriggerSubscriptionRequestTimeoutMs = defaultTriggerSubscriptionRequestTimeoutMs
	}
	if l.TriggerAllRegistrationsTimeoutMs == 0 {
		l.TriggerAllRegistrationsTimeoutMs = defaultTriggerAllRegistrationsTimeoutMs
	}
	if l.MaxTriggerSubscriptions == 0 {
		l.MaxTriggerSubscriptions = defaultMaxTriggerSubscriptions
	}
	if l.TriggerEventQueueSize == 0 {
		l.TriggerEventQueueSize = defaultTriggerEventQueueSize
	}
	if l.TriggerEventMaxAgeMs == 0 {
		l.TriggerEventMaxAgeMs = defaultTriggerEventMaxAgeMs
	}
	if l.MaxConcurrentWorkflowExecutions == 0 {
		l.MaxConcurrentWorkflowExecutions = defaultMaxConcurrentWorkflowExecutions
	}
	if l.MaxConcurrentCapabilityCallsPerWorkflow == 0 {
		l.MaxConcurrentCapabilityCallsPerWorkflow = defaultMaxConcurrentCapabilityCallsPerWorkflow
	}
	if l.MaxConcurrentSecretsCallsPerWorkflow == 0 {
		l.MaxConcurrentSecretsCallsPerWorkflow = defaultMaxConcurrentSecretsCallsPerWorkflow
	}
	if l.WorkflowExecutionTimeoutMs == 0 {
		l.WorkflowExecutionTimeoutMs = defaultWorkflowExecutionTimeoutMs
	}
	if l.CapabilityCallTimeoutMs == 0 {
		l.CapabilityCallTimeoutMs = defaultCapabilityCallTimeoutMs
	}
	if l.MaxUserLogEventsPerExecution == 0 {
		l.MaxUserLogEventsPerExecution = defaultMaxUserLogEventsPerExecution
	}
	if l.MaxUserLogLineLength == 0 {
		l.MaxUserLogLineLength = defaultMaxUserLogLineLength
	}
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
	if h.OnExecutionFinished == nil {
		h.OnExecutionFinished = func(executionID string, status string) {}
	}
	if h.OnRateLimited == nil {
		h.OnRateLimited = func(executionID string) {}
	}
}
