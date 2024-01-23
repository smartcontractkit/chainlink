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
	id := info.ID
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
