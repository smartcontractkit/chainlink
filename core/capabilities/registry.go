package capabilities

import (
	"context"
	"fmt"
	"sync"
)

type Registry struct {
	m  map[string]Capability
	mu sync.RWMutex
}

func (r *Registry) Get(_ context.Context, id string) (Capability, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.m[id]
	if !ok {
		return nil, fmt.Errorf("capability not found with id %s", id)
	}

	return c, nil
}

func (r *Registry) getCapabilityOfType(ctx context.Context, id string, ct CapabilityType) (Capability, error) {
	c, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if c.Info().CapabilityType != ct {
		return nil, fmt.Errorf("capability with id %s is not of type %s", id, ct)
	}

	return c, nil
}

func (r *Registry) GetAction(ctx context.Context, id string) (SynchronousCapability, error) {
	c, err := r.getCapabilityOfType(ctx, id, CapabilityTypeAction)
	if err != nil {
		return nil, err
	}
	return c.(SynchronousCapability), err
}

func (r *Registry) GetTarget(ctx context.Context, id string) (SynchronousCapability, error) {
	c, err := r.getCapabilityOfType(ctx, id, CapabilityTypeTarget)
	if err != nil {
		return nil, err
	}
	return c.(SynchronousCapability), err
}

func (r *Registry) GetTrigger(ctx context.Context, id string) (AsynchronousCapability, error) {
	c, err := r.getCapabilityOfType(ctx, id, CapabilityTypeTrigger)
	if err != nil {
		return nil, err
	}
	return c.(AsynchronousCapability), err
}

func (r *Registry) GetReport(ctx context.Context, id string) (AsynchronousCapability, error) {
	c, err := r.getCapabilityOfType(ctx, id, CapabilityTypeReport)
	if err != nil {
		return nil, err
	}
	return c.(AsynchronousCapability), err
}

func (r *Registry) List(_ context.Context) []Capability {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cl := []Capability{}
	for _, v := range r.m {
		cl = append(cl, v)
	}

	return cl
}

func (r *Registry) Add(_ context.Context, c Capability) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info := c.Info()
	id := info.Id
	_, ok := r.m[id]
	if ok {
		return fmt.Errorf("capability with id: %s already exists", id)
	}

	switch info.CapabilityType {
	case CapabilityTypeAction, CapabilityTypeTarget:
		_, ok := c.(SynchronousCapability)
		if !ok {
			return fmt.Errorf("capability with id %s, type %s does not satisfy SynchronousCapability interface", id, info.CapabilityType)
		}
	case CapabilityTypeReport, CapabilityTypeTrigger:
		_, ok := c.(AsynchronousCapability)
		if !ok {
			return fmt.Errorf("capability with id %s, type %s does not satisfy AsynchronousCapability interface", id, info.CapabilityType)
		}
	}

	r.m[id] = c
	return nil

}

func NewRegistry() *Registry {
	return &Registry{
		m: map[string]Capability{},
	}
}
