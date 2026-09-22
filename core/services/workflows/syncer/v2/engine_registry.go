package v2

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/hashutil"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

var ErrNotFound = errors.New("engine not found")
var ErrAlreadyExists = errors.New("attempting to register duplicate engine")

type ServiceWithMetadata struct {
	WorkflowID types.WorkflowID
	Source     string // Which source this workflow came from (e.g., "ContractWorkflowSource", "GRPCWorkflowSource")
	// ReconcileKey fingerprints the on-chain record (owner/name) the engine was started for.
	// Empty when the engine was registered without identity metadata (e.g. via Add).
	ReconcileKey string
	// Coordinated is true if the TriggerCoordinator owns this workflow's trigger
	// registration, handles, and acknowledgement. This is true when the engine
	// behind this entry is a v2.ExecutionEngine, not the legacy v2.Engine. Set
	// once at Add and never changed: the flag decision that produced this entry
	// is fixed for its lifetime, even if the flag itself later flips. This is
	// the single source of truth cleanup uses to decide whether to call the
	// coordinator and whether to free the syncer-owned workflow-count limit.
	Coordinated bool
	// Owner is the workflow owner, needed at teardown to Free the
	// syncer-owned workflow-count limit under the same contexts.CRE{Owner: ...}
	// it was acquired under — the per-owner resource pool keys on that context
	// value, so a Free with the wrong (or missing) owner leaks the slot rather
	// than releasing it. Empty when the engine was registered without identity
	// metadata (e.g. via Add).
	Owner string
	services.Service
}

// engineEntry holds the engine and its associated source for internal storage.
type engineEntry struct {
	engine       services.Service
	source       string
	reconcileKey string
	coordinated  bool
	owner        string
}

// ReconcileKey fingerprints the workflow record identity that a WorkflowID is expected to map to.
func ReconcileKey(owner []byte, name string) (string, error) {
	h, err := hashutil.BytesOfBytesKeccak([][]byte{owner, []byte(name)})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h[:]), nil
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
	return r.AddWithReconcileKey(workflowID, source, "", engine)
}

// AddWithReconcileKey adds an engine to the registry with its source and identity fingerprint.
func (r *EngineRegistry) AddWithReconcileKey(workflowID types.WorkflowID, source, reconcileKey string, engine services.Service) error {
	return r.AddCoordinated(workflowID, source, reconcileKey, false, "", engine)
}

// AddCoordinated adds an engine to the registry, recording whether the
// TriggerCoordinator owns its trigger registration/handles/ACK (CRE-6176) and
// the workflow owner (needed at teardown to Free the workflow-count limit
// under the right per-owner context). coordinated and owner are decided once,
// by the caller, at construction time — stored verbatim, never re-derived.
func (r *EngineRegistry) AddCoordinated(workflowID types.WorkflowID, source, reconcileKey string, coordinated bool, owner string, engine services.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, found := r.engines[workflowID]; found {
		return ErrAlreadyExists
	}
	r.engines[workflowID] = engineEntry{
		engine:       engine,
		source:       source,
		reconcileKey: reconcileKey,
		coordinated:  coordinated,
		owner:        owner,
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
		WorkflowID:   workflowID,
		Source:       entry.source,
		ReconcileKey: entry.reconcileKey,
		Coordinated:  entry.coordinated,
		Owner:        entry.owner,
		Service:      entry.engine,
	}, true
}

// GetAll retrieves all engines from the engine registry.
func (r *EngineRegistry) GetAll() []ServiceWithMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	engines := make([]ServiceWithMetadata, 0, len(r.engines))
	for workflowID, entry := range r.engines {
		engines = append(engines, ServiceWithMetadata{
			WorkflowID:  workflowID,
			Source:      entry.source,
			Coordinated: entry.coordinated,
			Owner:       entry.owner,
			Service:     entry.engine,
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
				WorkflowID:  workflowID,
				Source:      entry.source,
				Coordinated: entry.coordinated,
				Owner:       entry.owner,
				Service:     entry.engine,
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
		WorkflowID:  workflowID,
		Source:      entry.source,
		Coordinated: entry.coordinated,
		Owner:       entry.owner,
		Service:     entry.engine,
	}, nil
}

// PopAll removes and returns all engines.
func (r *EngineRegistry) PopAll() []ServiceWithMetadata {
	r.mu.Lock()
	defer r.mu.Unlock()
	engines := make([]ServiceWithMetadata, 0, len(r.engines))
	for workflowID, entry := range r.engines {
		engines = append(engines, ServiceWithMetadata{
			WorkflowID:  workflowID,
			Source:      entry.source,
			Coordinated: entry.coordinated,
			Owner:       entry.owner,
			Service:     entry.engine,
		})
	}
	r.engines = make(map[[32]byte]engineEntry)
	return engines
}
