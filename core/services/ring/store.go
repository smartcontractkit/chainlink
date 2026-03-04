package ring

import (
	"context"
	"maps"
	"slices"
	"sync"
	"time"

	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
)

type MappingMeta struct {
	OldShardID   uint32
	NewShardID   uint32
	InTransition bool
	UpdatedAt    time.Time
}

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
	routingState     map[string]uint32       // workflow_id -> shard_id (cache of allocated workflows)
	routingStateMeta map[string]*MappingMeta // workflow_id -> mapping metadata
	shardHealth      map[uint32]bool         // shard_id -> is_healthy
	healthyShards    []uint32                // Sorted list of healthy shards
	currentState     *ringpb.RoutingState    // Current routing state (steady or transition)
	mappingVersion   uint64

	pendingAllocs map[string][]chan uint32 // workflow_id -> waiting channels
	allocRequests chan AllocationRequest   // Channel for new allocation requests

	mu sync.Mutex
}

const (
	AllocationRequestChannelCapacity = 1000
	getShardTransitionTimeout        = 30 * time.Second
	submitAllocRetries               = 5
	submitAllocRetryInterval         = 20 * time.Millisecond
)

func NewStore() *Store {
	return &Store{
		routingState:     make(map[string]uint32),
		routingStateMeta: make(map[string]*MappingMeta),
		shardHealth:      make(map[uint32]bool),
		healthyShards:    make([]uint32, 0),
		pendingAllocs:    make(map[string][]chan uint32),
		allocRequests:    make(chan AllocationRequest, AllocationRequestChannelCapacity),
		mu:               sync.Mutex{},
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

	slices.Sort(s.healthyShards)
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

	resultCh := make(chan uint32, 1)
	s.pendingAllocs[workflowID] = append(s.pendingAllocs[workflowID], resultCh)
	s.mu.Unlock()

	select {
	case s.allocRequests <- AllocationRequest{WorkflowID: workflowID, Result: resultCh}:
	case <-ctx.Done():
		return 0, ctx.Err()
	}

	runCtx := ctx
	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		runCtx, cancel = context.WithTimeout(ctx, getShardTransitionTimeout)
		defer cancel()
	}
	select {
	case shard := <-resultCh:
		return shard, nil
	case <-runCtx.Done():
		s.mu.Lock()
		ring := newShardRing(s.healthyShards)
		s.mu.Unlock()
		return locateShard(ring, workflowID)
	}
}

// SetShardForWorkflow is called by the RingOCR plugin whenever it finishes a round with allocations for a given workflow ID.
func (s *Store) SetShardForWorkflow(workflowID string, shardID uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldShardID := s.routingState[workflowID]
	s.routingState[workflowID] = shardID

	inTransition := !IsInSteadyState(s.currentState)
	s.routingStateMeta[workflowID] = &MappingMeta{
		OldShardID:   oldShardID,
		NewShardID:   shardID,
		InTransition: inTransition,
		UpdatedAt:    time.Now(),
	}
	s.mappingVersion++

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

func (s *Store) GetMappingVersion() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.mappingVersion
}

func (s *Store) GetWorkflowMappingsBatch(workflowIDs []string) (map[string]*MappingMeta, uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make(map[string]*MappingMeta, len(workflowIDs))
	for _, wfID := range workflowIDs {
		if shardID, ok := s.routingState[wfID]; ok {
			meta := s.routingStateMeta[wfID]
			if meta != nil {
				result[wfID] = &MappingMeta{
					OldShardID:   meta.OldShardID,
					NewShardID:   meta.NewShardID,
					InTransition: meta.InTransition,
					UpdatedAt:    meta.UpdatedAt,
				}
			} else {
				result[wfID] = &MappingMeta{
					NewShardID: shardID,
					UpdatedAt:  time.Now(),
				}
			}
		}
	}
	return result, s.mappingVersion
}

func (s *Store) RegisterWorkflowsFromShard(shardID uint32, workflowIDs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for _, wfID := range workflowIDs {
		if _, exists := s.routingState[wfID]; !exists {
			s.routingState[wfID] = shardID
			s.routingStateMeta[wfID] = &MappingMeta{
				OldShardID:   0,
				NewShardID:   shardID,
				InTransition: false,
				UpdatedAt:    now,
			}
		}
	}
	s.mappingVersion++
}

// SyncRoutes atomically replaces the routing map with the authoritative set
// from the OCR outcome, pruning any workflow IDs that are no longer present.
func (s *Store) SyncRoutes(routes map[string]uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	inTransition := !IsInSteadyState(s.currentState)

	for wfID, shardID := range routes {
		old := s.routingState[wfID]
		s.routingState[wfID] = shardID
		s.routingStateMeta[wfID] = &MappingMeta{
			OldShardID:   old,
			NewShardID:   shardID,
			InTransition: inTransition,
			UpdatedAt:    now,
		}
		if waiters, ok := s.pendingAllocs[wfID]; ok {
			for _, ch := range waiters {
				select {
				case ch <- shardID:
				default:
				}
			}
			delete(s.pendingAllocs, wfID)
		}
	}

	for wfID := range s.routingState {
		if _, keep := routes[wfID]; !keep {
			delete(s.routingState, wfID)
			delete(s.routingStateMeta, wfID)
		}
	}

	s.mappingVersion++
}

func (s *Store) SubmitWorkflowsForAllocation(workflowIDs []string) (dropped int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, wfID := range workflowIDs {
		if _, exists := s.routingState[wfID]; !exists {
			enqueued := false
			for attempt := 0; attempt < submitAllocRetries && !enqueued; attempt++ {
				select {
				case s.allocRequests <- AllocationRequest{WorkflowID: wfID, Result: nil}:
					enqueued = true
				default:
					if attempt < submitAllocRetries-1 {
						s.mu.Unlock()
						time.Sleep(submitAllocRetryInterval)
						s.mu.Lock()
					} else {
						dropped++
					}
				}
			}
		}
	}
	return dropped
}
