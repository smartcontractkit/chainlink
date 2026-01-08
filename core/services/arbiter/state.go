package arbiter

import (
	"sync"
)

// State holds the current scaling state.
type State struct {
	currentReplicas       map[string]ShardReplica
	lastScalingReason     string
	desiredReplicasCount  int
	approvedReplicasCount int
	mu                    sync.RWMutex
}

// NewState creates a new State with default values.
func NewState() *State {
	return &State{
		currentReplicas:       make(map[string]ShardReplica),
		desiredReplicasCount:  1,
		approvedReplicasCount: 1,
		lastScalingReason:     "Initial state",
	}
}

// Update updates the state with new scale intent data.
func (s *State) Update(currentReplicas map[string]ShardReplica, desiredCount int, reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentReplicas = currentReplicas
	s.desiredReplicasCount = desiredCount
	s.lastScalingReason = reason
}

// SetApprovedCount sets the approved replica count.
func (s *State) SetApprovedCount(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.approvedReplicasCount = count
}

// GetScalingSpec returns the current scaling specification.
func (s *State) GetScalingSpec() ScalingSpecResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return ScalingSpecResponse{
		CurrentReplicaCount:  len(s.currentReplicas),
		DesiredReplicaCount:  s.desiredReplicasCount,
		ApprovedReplicaCount: s.approvedReplicasCount,
		LastScalingReason:    s.lastScalingReason,
	}
}

// GetCurrentReplicaCount returns the current number of replicas.
func (s *State) GetCurrentReplicaCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.currentReplicas)
}

// GetDesiredReplicaCount returns the desired number of replicas.
func (s *State) GetDesiredReplicaCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.desiredReplicasCount
}

// GetApprovedReplicaCount returns the approved number of replicas.
func (s *State) GetApprovedReplicaCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.approvedReplicasCount
}
