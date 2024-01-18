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

	id := c.Info().Id.String()
	_, ok := r.m[id]
	if ok {
		return fmt.Errorf("capability with id: %s already exists", id)
	}

	r.m[id] = c
	return nil

}

func NewRegistry() *Registry {
	return &Registry{
		m: map[string]Capability{},
	}
}
