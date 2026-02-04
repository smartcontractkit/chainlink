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
	Source     string // Which source this workflow came from (e.g., "ContractWorkflowSource", "GRPCWorkflowSource")
	services.Service
}

// engineEntry holds the engine and its associated source for internal storage
type engineEntry struct {
	engine services.Service
	source string
}

type EngineRegistry struct {
	engines map[[32]byte]engineEntry
	mu      sync.RWMutex
}

func NewEngineRegistry() *EngineRegistry {
	return &EngineRegistry{
		engines: make(map[[32]byte]engineEntry),
	}
}

// Add adds an engine to the registry with its source.
func (r *EngineRegistry) Add(workflowID types.WorkflowID, source string, engine services.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, found := r.engines[workflowID]; found {
		return ErrAlreadyExists
	}
	r.engines[workflowID] = engineEntry{
		engine: engine,
		source: source,
	}
	return nil
}

// Get retrieves an engine from the registry. The second return value indicates whether an engine was found or not.
func (r *EngineRegistry) Get(workflowID types.WorkflowID) (ServiceWithMetadata, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, found := r.engines[workflowID]
	if !found {
		return ServiceWithMetadata{}, false
	}
	return ServiceWithMetadata{
		WorkflowID: workflowID,
		Source:     entry.source,
		Service:    entry.engine,
	}, true
}

// GetAll retrieves all engines from the engine registry.
func (r *EngineRegistry) GetAll() []ServiceWithMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	engines := []ServiceWithMetadata{}
	for workflowID, entry := range r.engines {
		engines = append(engines, ServiceWithMetadata{
			WorkflowID: workflowID,
			Source:     entry.source,
			Service:    entry.engine,
		})
	}
	return engines
}

// GetBySource retrieves all engines from a specific source.
func (r *EngineRegistry) GetBySource(source string) []ServiceWithMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []ServiceWithMetadata
	for workflowID, entry := range r.engines {
		if entry.source == source {
			result = append(result, ServiceWithMetadata{
				WorkflowID: workflowID,
				Source:     entry.source,
				Service:    entry.engine,
			})
		}
	}
	return result
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
	entry, ok := r.engines[workflowID]
	if !ok {
		return ServiceWithMetadata{}, fmt.Errorf("pop failed: %w", ErrNotFound)
	}
	delete(r.engines, workflowID)
	return ServiceWithMetadata{
		WorkflowID: workflowID,
		Source:     entry.source,
		Service:    entry.engine,
	}, nil
}

// PopAll removes and returns all engines.
func (r *EngineRegistry) PopAll() []ServiceWithMetadata {
	r.mu.Lock()
	defer r.mu.Unlock()
	engines := []ServiceWithMetadata{}
	for workflowID, entry := range r.engines {
		engines = append(engines, ServiceWithMetadata{
			WorkflowID: workflowID,
			Source:     entry.source,
			Service:    entry.engine,
		})
	}
	r.engines = make(map[[32]byte]engineEntry)
	return engines
}
