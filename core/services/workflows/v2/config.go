package v2

import (
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

type EngineConfig struct {
	Lggr            logger.Logger
	Module          host.ModuleV2
	CapRegistry     core.CapabilitiesRegistry
	ExecutionsStore store.Store

	WorkflowID string

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

	// TODO(CAPPL-733): add owner and name; validate format
	if c.WorkflowID == "" {
		return errors.New("workflowID not set")
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
