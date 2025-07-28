package v2

import (
	"errors"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

var ErrNotFound = errors.New("engine not found")
var ErrAlreadyExists = errors.New("attempting to register duplicate engine")

type ServiceWithMetadata struct {
	WorkflowID types.WorkflowID
	services.Service
}

type EngineRegistry struct {
	engines map[[32]byte]services.Service
	mu      sync.RWMutex
}

func NewEngineRegistry() *EngineRegistry {
	return &EngineRegistry{
		engines: make(map[[32]byte]services.Service),
	}
}

// Add adds an engine to the registry.
func (r *EngineRegistry) Add(workflowID types.WorkflowID, engine services.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, found := r.engines[workflowID]; found {
		return ErrAlreadyExists
	}
	r.engines[workflowID] = engine
	return nil
}

// Get retrieves an engine from the registry. The second return value indicates whether an engine was found or not.
func (r *EngineRegistry) Get(workflowID types.WorkflowID) (ServiceWithMetadata, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	engine, found := r.engines[workflowID]
	if !found {
		return ServiceWithMetadata{}, false
	}
	return ServiceWithMetadata{
		WorkflowID: workflowID,
		Service:    engine,
	}, true
}

// GetAll retrieves all engines from the engine registry.
func (r *EngineRegistry) GetAll() []ServiceWithMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	engines := []ServiceWithMetadata{}
	for workflowID, engine := range r.engines {
		engines = append(engines, ServiceWithMetadata{
			WorkflowID: workflowID,
			Service:    engine,
		})
	}
	return engines
}

// Contains is true if the engine exists.
func (r *EngineRegistry) Contains(workflowID types.WorkflowID) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, found := r.engines[workflowID]
	return found
}

// Pop removes an engine from the registry and returns the engine if found.
func (r *EngineRegistry) Pop(workflowID types.WorkflowID) (ServiceWithMetadata, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	engine, ok := r.engines[workflowID]
	if !ok {
		return ServiceWithMetadata{}, fmt.Errorf("pop failed: %w", ErrNotFound)
	}
	delete(r.engines, workflowID)
	return ServiceWithMetadata{
		WorkflowID: workflowID,
		Service:    engine,
	}, nil
}

// PopAll removes and returns all engines.
func (r *EngineRegistry) PopAll() []ServiceWithMetadata {
	r.mu.Lock()
	defer r.mu.Unlock()
	engines := []ServiceWithMetadata{}
	for workflowID, engine := range r.engines {
		engines = append(engines, ServiceWithMetadata{
			WorkflowID: workflowID,
			Service:    engine,
		})
	}
	r.engines = make(map[[32]byte]services.Service)
	return engines
}
