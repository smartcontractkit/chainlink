package capabilities

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Registry is a struct for the registry of capabilities.
// Registry is safe for concurrent use.
type Registry struct {
	m    map[string]capabilities.BaseCapability
	mu   sync.RWMutex
	lggr logger.Logger
}

// Get gets a capability from the registry.
func (r *Registry) Get(_ context.Context, id string) (capabilities.BaseCapability, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	r.lggr.Debugw("get capability", "id", id)
	c, ok := r.m[id]
	if !ok {
		return nil, fmt.Errorf("capability not found with id %s", id)
	}

	return c, nil
}

// GetTrigger gets a capability from the registry and tries to coerce it to the TriggerCapability interface.
func (r *Registry) GetTrigger(ctx context.Context, id string) (capabilities.TriggerCapability, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	tc, ok := c.(capabilities.TriggerCapability)
	if !ok {
		return nil, fmt.Errorf("capability with id: %s does not satisfy the capability interface", id)
	}

	return tc, nil
}

// GetAction gets a capability from the registry and tries to coerce it to the ActionCapability interface.
func (r *Registry) GetAction(ctx context.Context, id string) (capabilities.ActionCapability, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	ac, ok := c.(capabilities.ActionCapability)
	if !ok {
		return nil, fmt.Errorf("capability with id: %s does not satisfy the capability interface", id)
	}

	return ac, nil
}

// GetConsensus gets a capability from the registry and tries to coerce it to the ConsensusCapability interface.
func (r *Registry) GetConsensus(ctx context.Context, id string) (capabilities.ConsensusCapability, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	cc, ok := c.(capabilities.ConsensusCapability)
	if !ok {
		return nil, fmt.Errorf("capability with id: %s does not satisfy the capability interface", id)
	}

	return cc, nil
}

// GetTarget gets a capability from the registry and tries to coerce it to the TargetCapability interface.
func (r *Registry) GetTarget(ctx context.Context, id string) (capabilities.TargetCapability, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	tc, ok := c.(capabilities.TargetCapability)
	if !ok {
		return nil, fmt.Errorf("capability with id: %s does not satisfy the capability interface", id)
	}

	return tc, nil
}

// List lists all the capabilities in the registry.
func (r *Registry) List(_ context.Context) ([]capabilities.BaseCapability, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cl := []capabilities.BaseCapability{}
	for _, v := range r.m {
		cl = append(cl, v)
	}

	return cl, nil
}

// Add adds a capability to the registry.
func (r *Registry) Add(ctx context.Context, c capabilities.BaseCapability) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info, err := c.Info(ctx)
	if err != nil {
		return err
	}

	switch info.CapabilityType {
	case capabilities.CapabilityTypeTrigger:
		_, ok := c.(capabilities.TriggerCapability)
		if !ok {
			return fmt.Errorf("trigger capability does not satisfy TriggerCapability interface")
		}
	case capabilities.CapabilityTypeAction:
		_, ok := c.(capabilities.ActionCapability)
		if !ok {
			return fmt.Errorf("action does not satisfy ActionCapability interface")
		}
	case capabilities.CapabilityTypeConsensus:
		_, ok := c.(capabilities.ConsensusCapability)
		if !ok {
			return fmt.Errorf("consensus capability does not satisfy ConsensusCapability interface")
		}
	case capabilities.CapabilityTypeTarget:
		_, ok := c.(capabilities.TargetCapability)
		if !ok {
			return fmt.Errorf("target capability does not satisfy TargetCapability interface")
		}
	default:
		return fmt.Errorf("unknown capability type: %s", info.CapabilityType)
	}

	id := info.ID
	_, ok := r.m[id]
	if ok {
		return fmt.Errorf("capability with id: %s already exists", id)
	}

	r.m[id] = c
	r.lggr.Infow("capability added", "id", id, "type", info.CapabilityType, "description", info.Description, "version", info.Version)
	return nil
}

// NewRegistry returns a new Registry.
func NewRegistry(lggr logger.Logger) *Registry {
	return &Registry{
		m:    map[string]capabilities.BaseCapability{},
		lggr: lggr.Named("CapabilityRegistry"),
	}
}
