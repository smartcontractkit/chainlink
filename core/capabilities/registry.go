package capabilities

import (
	"fmt"
	"sync"
)

type Registry struct {
	m  map[string]Capability
	mu sync.RWMutex
}

func (r *Registry) Get(id fmt.Stringer) (Capability, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.m[id.String()]
	if !ok {
		return nil, fmt.Errorf("capability not found with id %s", id)
	}

	return c, nil
}

func (r *Registry) getCapabilityOfType(id fmt.Stringer, ct CapabilityType) (Capability, error) {
	c, err := r.Get(id)
	if err != nil {
		return nil, err
	}

	if c.Info().CapabilityType != ct {
		return nil, fmt.Errorf("capability with id %s is not of type %s", id, ct)
	}

	return c, nil
}

func (r *Registry) GetAction(id fmt.Stringer) (SynchronousCapability, error) {
	c, err := r.getCapabilityOfType(id, CapabilityTypeAction)
	if err != nil {
		return nil, err
	}
	return c.(SynchronousCapability), err
}

func (r *Registry) GetTarget(id fmt.Stringer) (SynchronousCapability, error) {
	c, err := r.getCapabilityOfType(id, CapabilityTypeTarget)
	if err != nil {
		return nil, err
	}
	return c.(SynchronousCapability), err
}

func (r *Registry) GetTrigger(id fmt.Stringer) (AsynchronousCapability, error) {
	c, err := r.getCapabilityOfType(id, CapabilityTypeTrigger)
	if err != nil {
		return nil, err
	}
	return c.(AsynchronousCapability), err
}

func (r *Registry) GetReport(id fmt.Stringer) (AsynchronousCapability, error) {
	c, err := r.getCapabilityOfType(id, CapabilityTypeReport)
	if err != nil {
		return nil, err
	}
	return c.(AsynchronousCapability), err
}

func (r *Registry) List() []Capability {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cl := []Capability{}
	for _, v := range r.m {
		cl = append(cl, v)
	}

	return cl
}

func (r *Registry) Add(c Capability) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info := c.Info()
	id := info.Id.String()
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
