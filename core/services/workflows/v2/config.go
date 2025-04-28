package v2

import (
	"errors"
	"fmt"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

type EngineConfig struct {
	Lggr            logger.Logger
	Module          host.ModuleV2
	CapRegistry     core.CapabilitiesRegistry
	ExecutionsStore store.Store
	Clock           clockwork.Clock

	WorkflowID    string // hex-encoded [32]byte, no "0x" prefix
	WorkflowOwner string // hex-encoded [20]byte, no "0x" prefix
	WorkflowName  types.WorkflowName

	LocalLimits          EngineLimits             // local to a single workflow
	GlobalLimits         *syncerlimiter.Limits    // global to all workflows
	ExecutionRateLimiter *ratelimiter.RateLimiter // global + per owner

	Hooks LifecycleHooks
}

const (
	defaultMaxCapRegistryAccessRetries      = 0 // infinity
	defaultCapRegistryAccessRetryIntervalMs = 5000
	defaultMaxTotalTriggerSubscriptions     = 10
	defaultMaxConcurrentCapabilityCalls     = 10
)

type EngineLimits struct {
	MaxCapRegistryAccessRetries      int
	CapRegistryAccessRetryIntervalMs int

	MaxTotalTriggerSubscriptions int

	MaxConcurrentCapabilityCalls int
}

type LifecycleHooks struct {
	OnInitialized       func(err error)
	OnExecutionFinished func(executionID string)
	OnRateLimited       func(executionID string)
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

	c.Hooks.setDefaultHooks()
	return nil
}

func (l *EngineLimits) setDefaultLimits() {
	if l.MaxCapRegistryAccessRetries == 0 {
		l.MaxCapRegistryAccessRetries = defaultMaxCapRegistryAccessRetries
	}
	if l.CapRegistryAccessRetryIntervalMs == 0 {
		l.CapRegistryAccessRetryIntervalMs = defaultCapRegistryAccessRetryIntervalMs
	}
	if l.MaxTotalTriggerSubscriptions == 0 {
		l.MaxTotalTriggerSubscriptions = defaultMaxTotalTriggerSubscriptions
	}
	if l.MaxConcurrentCapabilityCalls == 0 {
		l.MaxConcurrentCapabilityCalls = defaultMaxConcurrentCapabilityCalls
	}
}

// set all to non-nil so the Engine doesn't have to check before each call
func (h *LifecycleHooks) setDefaultHooks() {
	if h.OnInitialized == nil {
		h.OnInitialized = func(err error) {}
	}
	if h.OnExecutionFinished == nil {
		h.OnExecutionFinished = func(executionID string) {}
	}
	if h.OnRateLimited == nil {
		h.OnRateLimited = func(executionID string) {}
	}
}
