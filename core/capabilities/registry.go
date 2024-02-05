package capabilities

import (
	"context"
	"fmt"
	"sync"
)

// Registry is a struct for the registry of capabilities.
// Registry is safe for concurrent use.
type Registry struct {
	m  map[string]BaseCapability
	mu sync.RWMutex
}

// Get gets a capability from the registry.
func (r *Registry) Get(_ context.Context, id string) (BaseCapability, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.m[id]
	if !ok {
		return nil, fmt.Errorf("capability not found with id %s", id)
	}

	return c, nil
}

// GetTrigger gets a capability from the registry and tries to coerce it to the TriggerCapability interface.
func (r *Registry) GetTrigger(ctx context.Context, id string) (TriggerCapability, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	tc, ok := c.(TriggerCapability)
	if !ok {
		return nil, fmt.Errorf("capability with id: %s does not satisfy the capability interface", id)
	}

	return tc, nil
}

// GetAction gets a capability from the registry and tries to coerce it to the ActionCapability interface.
func (r *Registry) GetAction(ctx context.Context, id string) (ActionCapability, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	ac, ok := c.(ActionCapability)
	if !ok {
		return nil, fmt.Errorf("capability with id: %s does not satisfy the capability interface", id)
	}

	return ac, nil
}

// GetConsensus gets a capability from the registry and tries to coerce it to the ActionCapability interface.
func (r *Registry) GetConsensus(ctx context.Context, id string) (ConsensusCapability, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	cc, ok := c.(ConsensusCapability)
	if !ok {
		return nil, fmt.Errorf("capability with id: %s does not satisfy the capability interface", id)
	}

	return cc, nil
}

// GetTarget gets a capability from the registry and tries to coerce it to the ActionCapability interface.
func (r *Registry) GetTarget(ctx context.Context, id string) (TargetCapability, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	tc, ok := c.(TargetCapability)
	if !ok {
		return nil, fmt.Errorf("capability with id: %s does not satisfy the capability interface", id)
	}

	return tc, nil
}

// List lists all the capabilities in the registry.
func (r *Registry) List(_ context.Context) []BaseCapability {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cl := []BaseCapability{}
	for _, v := range r.m {
		cl = append(cl, v)
	}

	return cl
}

// Add adds a capability to the registry.
func (r *Registry) Add(_ context.Context, c BaseCapability) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info := c.Info()

	switch info.CapabilityType {
	case CapabilityTypeTrigger:
		_, ok := c.(TriggerCapability)
		if !ok {
			return fmt.Errorf("trigger capability does not satisfy TriggerCapability interface")
		}
	case CapabilityTypeAction:
		_, ok := c.(ActionCapability)
		if !ok {
			return fmt.Errorf("action does not satisfy ActionCapability interface")
		}
	case CapabilityTypeConsensus:
		_, ok := c.(ConsensusCapability)
		if !ok {
			return fmt.Errorf("consensus capability does not satisfy ConsensusCapability interface")
		}
	case CapabilityTypeTarget:
		_, ok := c.(TargetCapability)
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
	return nil

}

// NewRegistry returns a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		m: map[string]BaseCapability{},
	}
}
