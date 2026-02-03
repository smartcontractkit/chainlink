package shardorchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// TransitionState represents the state of a workflow's shard assignment
type TransitionState uint8

const (
	StateSteady TransitionState = iota
	StateTransitioning
)

// String returns the string representation of the TransitionState
func (s TransitionState) String() string {
	switch s {
	case StateSteady:
		return "steady"
	case StateTransitioning:
		return "transitioning"
	default:
		return "unknown"
	}
}

// InTransition returns true if the state is transitioning
func (s TransitionState) InTransition() bool {
	return s == StateTransitioning
}

// WorkflowMappingState represents the state of a workflow assignment
type WorkflowMappingState struct {
	WorkflowID      string
	OldShardID      uint32
	NewShardID      uint32
	TransitionState TransitionState
	UpdatedAt       time.Time
}

// Store manages workflow-to-shard mappings that will be exposed via gRPC
// RingOCR plugin updates this store, and the gRPC service reads from it
type Store struct {
	// workflowMappings tracks the current shard assignment for each workflow
	workflowMappings map[string]*WorkflowMappingState // workflow_id -> mapping state

	// shardRegistrations tracks what workflows each shard has registered
	// This is populated by ReportWorkflowTriggerRegistration calls from shards
	shardRegistrations map[uint32]map[string]bool // shard_id -> set of workflow_ids

	// mappingVersion increments on any change to workflowMappings
	// Used by clients for cache invalidation
	mappingVersion uint64

	// lastUpdateTime tracks when mappings were last modified
	lastUpdateTime time.Time

	mu     sync.RWMutex
	logger logger.Logger
}

func NewStore(lggr logger.Logger) *Store {
	return &Store{
		workflowMappings:   make(map[string]*WorkflowMappingState),
		shardRegistrations: make(map[uint32]map[string]bool),
		mappingVersion:     0,
		lastUpdateTime:     time.Now(),
		logger:             logger.Named(lggr, "ShardOrchestratorStore"),
	}
}

// UpdateWorkflowMapping is called by RingOCR to update workflow assignments
// This is the primary data source for shard orchestration
func (s *Store) UpdateWorkflowMapping(ctx context.Context, workflowID string, oldShardID, newShardID uint32, state TransitionState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.workflowMappings[workflowID] = &WorkflowMappingState{
		WorkflowID:      workflowID,
		OldShardID:      oldShardID,
		NewShardID:      newShardID,
		TransitionState: state,
		UpdatedAt:       now,
	}

	s.mappingVersion++
	s.lastUpdateTime = now

	s.logger.Debugw("Updated workflow mapping",
		"workflowID", workflowID,
		"oldShardID", oldShardID,
		"newShardID", newShardID,
		"state", state.String(),
		"version", s.mappingVersion,
	)

	return nil
}

// BatchUpdateWorkflowMappings allows RingOCR to update multiple mappings atomically
func (s *Store) BatchUpdateWorkflowMappings(ctx context.Context, mappings []*WorkflowMappingState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for _, mapping := range mappings {
		s.workflowMappings[mapping.WorkflowID] = &WorkflowMappingState{
			WorkflowID:      mapping.WorkflowID,
			OldShardID:      mapping.OldShardID,
			NewShardID:      mapping.NewShardID,
			TransitionState: mapping.TransitionState,
			UpdatedAt:       now,
		}
	}

	s.mappingVersion++
	s.lastUpdateTime = now

	s.logger.Debugw("Batch updated workflow mappings", "count", len(mappings), "version", s.mappingVersion)
	return nil
}

// GetWorkflowMapping retrieves the shard assignment for a specific workflow
// This is called by the gRPC service to respond to GetWorkflowShardMapping requests
func (s *Store) GetWorkflowMapping(ctx context.Context, workflowID string) (*WorkflowMappingState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mapping, ok := s.workflowMappings[workflowID]
	if !ok {
		return nil, fmt.Errorf("workflow %s not found in shard mappings", workflowID)
	}

	// Return a copy to avoid external mutations
	return &WorkflowMappingState{
		WorkflowID:      mapping.WorkflowID,
		OldShardID:      mapping.OldShardID,
		NewShardID:      mapping.NewShardID,
		TransitionState: mapping.TransitionState,
		UpdatedAt:       mapping.UpdatedAt,
	}, nil
}

// GetAllWorkflowMappings returns all current workflow-to-shard assignments
func (s *Store) GetAllWorkflowMappings(ctx context.Context) ([]*WorkflowMappingState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mappings := make([]*WorkflowMappingState, 0, len(s.workflowMappings))
	for _, mapping := range s.workflowMappings {
		mappings = append(mappings, &WorkflowMappingState{
			WorkflowID:      mapping.WorkflowID,
			OldShardID:      mapping.OldShardID,
			NewShardID:      mapping.NewShardID,
			TransitionState: mapping.TransitionState,
			UpdatedAt:       mapping.UpdatedAt,
		})
	}

	return mappings, nil
}

// ReportShardRegistration is called when a shard reports its registered workflows
// This helps track which workflows each shard has successfully loaded
// It also updates workflowMappings so GetWorkflowShardMapping returns correct data
func (s *Store) ReportShardRegistration(ctx context.Context, shardID uint32, workflowIDs []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	// Clear and update shard registrations
	s.shardRegistrations[shardID] = make(map[string]bool)
	for _, wfID := range workflowIDs {
		s.shardRegistrations[shardID][wfID] = true
	}

	// Also update workflowMappings - when a shard reports it has a workflow,
	// that's authoritative information about where the workflow is running
	for _, wfID := range workflowIDs {
		existing, ok := s.workflowMappings[wfID]
		if !ok || existing.NewShardID != shardID {
			s.workflowMappings[wfID] = &WorkflowMappingState{
				WorkflowID:      wfID,
				OldShardID:      0,
				NewShardID:      shardID,
				TransitionState: StateSteady,
				UpdatedAt:       now,
			}
			if ok {
				s.workflowMappings[wfID].OldShardID = existing.NewShardID
			}
		}
	}

	s.mappingVersion++
	s.lastUpdateTime = now

	s.logger.Debugw("Updated shard registrations",
		"shardID", shardID,
		"workflowCount", len(workflowIDs),
		"version", s.mappingVersion,
	)

	return nil
}

// GetShardRegistrations returns the workflows registered on a specific shard
func (s *Store) GetShardRegistrations(ctx context.Context, shardID uint32) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	workflows, ok := s.shardRegistrations[shardID]
	if !ok {
		return []string{}, nil
	}

	result := make([]string, 0, len(workflows))
	for wfID := range workflows {
		result = append(result, wfID)
	}

	return result, nil
}

// GetWorkflowMappingsBatch retrieves mappings for multiple workflows
func (s *Store) GetWorkflowMappingsBatch(ctx context.Context, workflowIDs []string) (map[string]*WorkflowMappingState, uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*WorkflowMappingState, len(workflowIDs))

	for _, workflowID := range workflowIDs {
		if mapping, ok := s.workflowMappings[workflowID]; ok {
			// Return a copy to avoid external mutations
			result[workflowID] = &WorkflowMappingState{
				WorkflowID:      mapping.WorkflowID,
				OldShardID:      mapping.OldShardID,
				NewShardID:      mapping.NewShardID,
				TransitionState: mapping.TransitionState,
				UpdatedAt:       mapping.UpdatedAt,
			}
		}
	}

	return result, s.mappingVersion, nil
}

// GetMappingVersion returns the current version of the mapping set
func (s *Store) GetMappingVersion() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.mappingVersion
}
