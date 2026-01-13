package arbiter

import (
	"sync"
)

// State holds the current scaling state for the Arbiter.
type State struct {
	currentReplicas     map[string]ShardReplica
	consensusWantShards int // Number of shards the Ring consensus wants
	mu                  sync.RWMutex
}

// NewState creates a new State with default values.
func NewState() *State {
	return &State{
		currentReplicas: make(map[string]ShardReplica),
	}
}

// SetCurrentReplicas updates the current replicas map.
func (s *State) SetCurrentReplicas(replicas map[string]ShardReplica) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentReplicas = replicas
}

// GetCurrentReplicaCount returns the current number of replicas.
func (s *State) GetCurrentReplicaCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.currentReplicas)
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
		isHealthy := replica.Status == StatusReady
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
