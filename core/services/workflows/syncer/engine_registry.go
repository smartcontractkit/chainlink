package syncer

import (
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

var errNotFound = errors.New("engine not found")

type EngineRegistryKey struct {
	Owner []byte
	Name  string
}

// KeyFor generates a key that will be used to identify the engine in the engine registry.
// This is used instead of a Workflow ID, because the WID will change if the workflow code is modified.
func (k EngineRegistryKey) keyFor() string {
	return hex.EncodeToString(k.Owner) + "-" + k.Name
}

type ServiceWithMetadata struct {
	WorkflowID    types.WorkflowID
	WorkflowName  string
	WorkflowOwner []byte
	services.Service
}

type EngineRegistry struct {
	engines map[string]ServiceWithMetadata
	mu      sync.RWMutex
}

func NewEngineRegistry() *EngineRegistry {
	return &EngineRegistry{
		engines: make(map[string]ServiceWithMetadata),
	}
}

// Add adds an engine to the registry.
func (r *EngineRegistry) Add(key EngineRegistryKey, engine services.Service, workflowID types.WorkflowID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := key.keyFor()
	if _, found := r.engines[k]; found {
		return errors.New("attempting to register duplicate engine")
	}
	r.engines[k] = ServiceWithMetadata{
		WorkflowID:    workflowID,
		WorkflowName:  key.Name,
		WorkflowOwner: key.Owner,
		Service:       engine,
	}
	return nil
}

// Get retrieves an engine from the registry. The second argument indicates whether an engine was found or not.
func (r *EngineRegistry) Get(key EngineRegistryKey) (ServiceWithMetadata, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	engine, found := r.engines[key.keyFor()]
	if !found {
		return ServiceWithMetadata{}, false
	}
	return engine, true
}

// GetAll retrieves all engines from the engine registry.
func (r *EngineRegistry) GetAll() []ServiceWithMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	engines := []ServiceWithMetadata{}
	for _, enginWithMetadata := range r.engines {
		engines = append(engines, enginWithMetadata)
	}
	return engines
}

// Contains is true if the engine exists.
func (r *EngineRegistry) Contains(key EngineRegistryKey) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, found := r.engines[key.keyFor()]
	return found
}

// Pop removes an engine from the registry and returns the engine if found.
func (r *EngineRegistry) Pop(key EngineRegistryKey) (ServiceWithMetadata, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := key.keyFor()
	engine, ok := r.engines[k]
	if !ok {
		return ServiceWithMetadata{}, fmt.Errorf("pop failed: %w", errNotFound)
	}
	delete(r.engines, k)
	return engine, nil
}

// PopAll removes and returns all engines.
func (r *EngineRegistry) PopAll() []ServiceWithMetadata {
	r.mu.Lock()
	defer r.mu.Unlock()
	all := slices.Collect(maps.Values(r.engines))
	r.engines = make(map[string]ServiceWithMetadata)
	return all
}
