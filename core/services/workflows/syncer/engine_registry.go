package syncer

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

var errNotFound = errors.New("engine not found")

type EngineRegistry struct {
	engines map[string]services.Service
	mu      sync.RWMutex
}

func NewEngineRegistry() *EngineRegistry {
	return &EngineRegistry{
		engines: make(map[string]services.Service),
	}
}

// Add adds an engine to the registry.
func (r *EngineRegistry) Add(id string, engine services.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, found := r.engines[id]; found {
		return errors.New("attempting to register duplicate engine")
	}
	r.engines[id] = engine
	return nil
}

// Get retrieves an engine from the registry.
func (r *EngineRegistry) Get(id string) (services.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	engine, found := r.engines[id]
	if !found {
		return nil, errNotFound
	}
	return engine, nil
}

// Contains is true if the engine exists.
func (r *EngineRegistry) Contains(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, found := r.engines[id]
	return found
}

// Pop removes an engine from the registry and returns the engine if found.
func (r *EngineRegistry) Pop(id string) (services.Service, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	engine, ok := r.engines[id]
	if !ok {
		return nil, fmt.Errorf("pop failed: %w", errNotFound)
	}
	delete(r.engines, id)
	return engine, nil
}

// PopAll removes and returns all engines.
func (r *EngineRegistry) PopAll() []services.Service {
	r.mu.Lock()
	defer r.mu.Unlock()
	all := slices.Collect(maps.Values(r.engines))
	r.engines = make(map[string]services.Service)
	return all
}
