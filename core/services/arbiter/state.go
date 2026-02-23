package arbiter

import (
	"slices"
	"strconv"
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
// Replica keys are iterated in sorted order so that shard IDs are deterministic
// and match the client's shard indices (e.g. from GetDesiredReplicas).
func (s *State) GetRoutableShards() RoutableShardsInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	readyCount := 0
	shardInfo := make(map[uint32]ShardHealth)

	// Sort keys so shard ID assignment is deterministic and matches caller's shard indices.
	// Callers (e.g. test env) send status with keys 0, 1, 2...; we must preserve that mapping.
	keys := make([]string, 0, len(s.currentReplicas))
	for k := range s.currentReplicas {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for i, key := range keys {
		replica := s.currentReplicas[key]
		isHealthy := replica.Status == StatusReady
		shardID := uint32(i) //nolint:gosec // G115: range index i is non-negative; shard count is small
		// If key is numeric, use it as shard ID so Ring's shard 0/1 matches topology's ShardIndex 0/1.
		if id, err := strconv.ParseUint(key, 10, 32); err == nil {
			shardID = uint32(id)
		}
		shardInfo[shardID] = ShardHealth{
			IsHealthy: isHealthy,
		}
		if isHealthy {
			readyCount++
		}
	}

	return RoutableShardsInfo{
		ReadyCount: readyCount,
		ShardInfo:  shardInfo,
	}
}
