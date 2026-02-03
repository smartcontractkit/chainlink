package ring

import (
	"context"
	"maps"
	"slices"
	"sync"

	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
)

// AllocationRequest represents a pending workflow allocation request during transition
type AllocationRequest struct {
	WorkflowID string
	Result     chan uint32
}

// Store manages shard routing state and workflow mappings.
// It serves as a shared data layer across three components:
//   - RingOCR plugin: produces consensus-driven routing updates
//   - Arbiter: provides shard health and scaling decisions
//   - ShardOrchestrator: consumes routing state to direct workflow execution
type Store struct {
	routingState  map[string]uint32    // workflow_id -> shard_id (cache of allocated workflows)
	shardHealth   map[uint32]bool      // shard_id -> is_healthy
	healthyShards []uint32             // Sorted list of healthy shards
	currentState  *ringpb.RoutingState // Current routing state (steady or transition)

	pendingAllocs map[string][]chan uint32 // workflow_id -> waiting channels
	allocRequests chan AllocationRequest   // Channel for new allocation requests

	mu sync.Mutex
}

const AllocationRequestChannelCapacity = 1000

func NewStore() *Store {
	return &Store{
		routingState:  make(map[string]uint32),
		shardHealth:   make(map[uint32]bool),
		healthyShards: make([]uint32, 0),
		pendingAllocs: make(map[string][]chan uint32),
		allocRequests: make(chan AllocationRequest, AllocationRequestChannelCapacity),
		mu:            sync.Mutex{},
	}
}

func (s *Store) updateHealthyShards() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.healthyShards = make([]uint32, 0)

	for shardID, healthy := range s.shardHealth {
		if healthy {
			s.healthyShards = append(s.healthyShards, shardID)
		}
	}

	// Sort for determinism
	slices.Sort(s.healthyShards)

	// If no healthy shards, add shard 0 as fallback
	if len(s.healthyShards) == 0 {
		s.healthyShards = []uint32{0}
	}
}

// GetShardForWorkflow called by Workflow Registry Syncers of all shards via ShardOrchestratorService.
func (s *Store) GetShardForWorkflow(ctx context.Context, workflowID string) (uint32, error) {
	s.mu.Lock()

	// Only trust the cache in steady state; during transition OCR may have invalidated it
	if IsInSteadyState(s.currentState) {
		// Check if already allocated in cache
		if shard, ok := s.routingState[workflowID]; ok {
			s.mu.Unlock()
			return shard, nil
		}
		ring := newShardRing(s.healthyShards)
		s.mu.Unlock()
		return locateShard(ring, workflowID)
	}

	// During transition, defer to OCR consensus for consistent shard assignment across nodes
	resultCh := make(chan uint32, 1)
	s.pendingAllocs[workflowID] = append(s.pendingAllocs[workflowID], resultCh)
	s.mu.Unlock()

	select {
	case s.allocRequests <- AllocationRequest{WorkflowID: workflowID, Result: resultCh}:
	case <-ctx.Done():
		return 0, ctx.Err()
	}

	select {
	case shard := <-resultCh:
		return shard, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

// SetShardForWorkflow is called by the RingOCR plugin whenever it finishes a round with allocations for a given workflow ID.
func (s *Store) SetShardForWorkflow(workflowID string, shardID uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.routingState[workflowID] = shardID

	// Signal any waiting allocation requests
	if waiters, ok := s.pendingAllocs[workflowID]; ok {
		for _, ch := range waiters {
			select {
			case ch <- shardID:
			default:
			}
		}
		delete(s.pendingAllocs, workflowID)
	}
}

// SetRoutingState is called by the RingOCR plugin whenever a state transition happens.
func (s *Store) SetRoutingState(state *ringpb.RoutingState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentState = state
}

func (s *Store) GetRoutingState() *ringpb.RoutingState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.currentState
}

// GetPendingAllocations called by the RingOCR plugin in the observation phase
// to collect all allocation requests (only applicable to the TRANSITION phase).
func (s *Store) GetPendingAllocations() []string {
	var pending []string
	for {
		select {
		case req := <-s.allocRequests:
			pending = append(pending, req.WorkflowID)
		default:
			return pending
		}
	}
}

func (s *Store) IsInTransition() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return !IsInSteadyState(s.currentState)
}

func (s *Store) GetShardHealth() map[uint32]bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return maps.Clone(s.shardHealth)
}

func (s *Store) SetShardHealth(shardID uint32, healthy bool) {
	s.mu.Lock()
	s.shardHealth[shardID] = healthy
	s.mu.Unlock()
	s.updateHealthyShards()
}

func (s *Store) SetAllShardHealth(health map[uint32]bool) {
	s.mu.Lock()
	s.shardHealth = make(map[uint32]bool)
	for k, v := range health {
		s.shardHealth[k] = v
	}

	// Uninitialized store must wait for OCR consensus before serving requests
	if s.currentState == nil {
		numHealthy := uint32(0)
		for _, healthy := range health {
			if healthy {
				numHealthy++
			}
		}
		s.currentState = &ringpb.RoutingState{
			State: &ringpb.RoutingState_Transition{
				Transition: &ringpb.Transition{
					WantShards: numHealthy,
				},
			},
		}
	}
	s.mu.Unlock()

	s.updateHealthyShards()
}

func (s *Store) GetAllRoutingState() map[string]uint32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return maps.Clone(s.routingState)
}

func (s *Store) DeleteWorkflow(workflowID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.routingState, workflowID)
}

func (s *Store) GetHealthyShardCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.healthyShards)
}

func (s *Store) GetHealthyShards() []uint32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.healthyShards)
}
