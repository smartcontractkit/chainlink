package v2

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

type EngineConfig struct {
	Lggr            logger.Logger
	Module          host.ModuleV2
	CapRegistry     core.CapabilitiesRegistry
	ExecutionsStore store.Store

	WorkflowID    string // hex-encoded [32]byte, no "0x" prefix
	WorkflowOwner string // hex-encoded [20]byte, no "0x" prefix
	WorkflowName  types.WorkflowName

	Limits EngineLimits
	Hooks  LifecycleHooks
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

	c.setDefaultLimits()
	c.setDefaultHooks()
	return nil
}

func (c *EngineConfig) setDefaultLimits() {
	if c.Limits.MaxCapRegistryAccessRetries == 0 {
		c.Limits.MaxCapRegistryAccessRetries = defaultMaxCapRegistryAccessRetries
	}
	if c.Limits.CapRegistryAccessRetryIntervalMs == 0 {
		c.Limits.CapRegistryAccessRetryIntervalMs = defaultCapRegistryAccessRetryIntervalMs
	}
	if c.Limits.MaxTotalTriggerSubscriptions == 0 {
		c.Limits.MaxTotalTriggerSubscriptions = defaultMaxTotalTriggerSubscriptions
	}
	if c.Limits.MaxConcurrentCapabilityCalls == 0 {
		c.Limits.MaxConcurrentCapabilityCalls = defaultMaxConcurrentCapabilityCalls
	}
}

// set all to non-nil so the Engine doesn't have to check before each call
func (c *EngineConfig) setDefaultHooks() {
	if c.Hooks.OnInitialized == nil {
		c.Hooks.OnInitialized = func(err error) {}
	}
	if c.Hooks.OnExecutionFinished == nil {
		c.Hooks.OnExecutionFinished = func(executionID string) {}
	}
	if c.Hooks.OnRateLimited == nil {
		c.Hooks.OnRateLimited = func(executionID string) {}
	}
}
