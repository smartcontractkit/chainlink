package grpcsourcemock

import (
	"errors"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/privateregistry"
)

// ErrWorkflowNotFound is returned when a workflow is not found in the store
var ErrWorkflowNotFound = errors.New("workflow not found")

// WorkflowStatus represents the status of a workflow.
// Values match the on-chain contract status (0=active, 1=paused).
type WorkflowStatus uint8

const (
	// WorkflowStatusActive indicates the workflow is active (matches contract status 0)
	WorkflowStatusActive WorkflowStatus = 0
	// WorkflowStatusPaused indicates the workflow is paused (matches contract status 1)
	WorkflowStatusPaused WorkflowStatus = 1
)

// StoredWorkflow represents a workflow stored in memory
type StoredWorkflow struct {
	Registration *privateregistry.WorkflowRegistration
	Status       WorkflowStatus
	// CreatedAt is the Unix timestamp in milliseconds when the workflow was first added
	CreatedAt int64
	// UpdatedAt is the Unix timestamp in milliseconds when the workflow was last modified
	UpdatedAt int64
}

// WorkflowStore is an in-memory store for workflows
type WorkflowStore struct {
	mu        sync.RWMutex
	workflows map[[32]byte]*StoredWorkflow
}

// NewWorkflowStore creates a new in-memory workflow store
func NewWorkflowStore() *WorkflowStore {
	return &WorkflowStore{
		workflows: make(map[[32]byte]*StoredWorkflow),
	}
}

// Add adds a workflow to the store. Concurrent safe.
// If the workflow already exists, it updates the existing workflow and bumps UpdatedAt.
func (s *WorkflowStore) Add(registration *privateregistry.WorkflowRegistration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	// Check if workflow already exists
	if existing, exists := s.workflows[registration.WorkflowID]; exists {
		// Update existing workflow and bump UpdatedAt
		existing.Registration = registration
		existing.UpdatedAt = now
		return nil
	}

	// Create new workflow with both timestamps set
	s.workflows[registration.WorkflowID] = &StoredWorkflow{
		Registration: registration,
		Status:       WorkflowStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return nil
}

// Update updates a workflow's status. Concurrent safe.
// It bumps the UpdatedAt timestamp whenever the workflow is modified.
func (s *WorkflowStore) Update(workflowID [32]byte, config *privateregistry.WorkflowStatusConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	wf, exists := s.workflows[workflowID]
	if !exists {
		return ErrWorkflowNotFound
	}

	if config.Paused {
		wf.Status = WorkflowStatusPaused
	} else {
		wf.Status = WorkflowStatusActive
	}

	// Bump UpdatedAt timestamp
	wf.UpdatedAt = time.Now().UnixMilli()
	return nil
}

// Delete removes a workflow from the store. Concurrent safe.
func (s *WorkflowStore) Delete(workflowID [32]byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.workflows[workflowID]; !exists {
		return ErrWorkflowNotFound
	}

	delete(s.workflows, workflowID)
	return nil
}

// List returns all workflows matching the given DON family filter
// If donFamilies is empty, all workflows are returned
func (s *WorkflowStore) List(donFamilies []string) []*StoredWorkflow {
	s.mu.RLock()
	defer s.mu.RUnlock()

	familySet := make(map[string]bool)
	for _, f := range donFamilies {
		familySet[f] = true
	}

	var result []*StoredWorkflow
	for _, wf := range s.workflows {
		// If no family filter, include all workflows
		if len(donFamilies) == 0 {
			result = append(result, wf)
			continue
		}
		// Otherwise, filter by family
		if familySet[wf.Registration.DonFamily] {
			result = append(result, wf)
		}
	}
	return result
}

// Get retrieves a workflow by ID. Concurrent safe.
func (s *WorkflowStore) Get(workflowID [32]byte) (*StoredWorkflow, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	wf, exists := s.workflows[workflowID]
	return wf, exists
}
