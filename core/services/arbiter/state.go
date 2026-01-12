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
	consensusWantShards   int // Number of shards the Ring consensus wants
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

// SetConsensusWantShards sets the number of shards the Ring consensus wants.
func (s *State) SetConsensusWantShards(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.consensusWantShards = count
}

// GetConsensusWantShards returns the number of shards the Ring consensus wants.
func (s *State) GetConsensusWantShards() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.consensusWantShards
}

// GetRoutableShards returns the count and status of shards ready for routing.
// This is used by Ring OCR to determine which shards can receive traffic.
// Only shards with Status == READY are counted as routable.
func (s *State) GetRoutableShards() RoutableShardsInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	readyCount := 0
	shardInfo := make(map[uint32]ShardHealth)

	// Iterate through current replicas and count READY ones
	shardID := uint32(0)
	for _, replica := range s.currentReplicas {
		isHealthy := replica.Status == StatusReady.String()
		shardInfo[shardID] = ShardHealth{
			IsHealthy: isHealthy,
		}
		if isHealthy {
			readyCount++
		}
		shardID++
	}

	return RoutableShardsInfo{
		ReadyCount: readyCount,
		ShardInfo:  shardInfo,
	}
}
