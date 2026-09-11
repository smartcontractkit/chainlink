package sharding

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// ShardExecutionStatusUpdateCapabilityID is the capability ID used for routing
// ExecutionStatusUpdate messages between shards via the d2d dispatcher.
const ShardExecutionStatusUpdateCapabilityID = "shard-execution-status-update"

// ExecutionStatusUpdateHandler is called when a quorum of
// ExecutionStatusUpdate messages has been collected for a workflow.
type ExecutionStatusUpdateHandler func(msg *ringpb.ExecutionStatusUpdate)

// workflowHandler bundles a per-workflow handler with the primary DON
// used for quorum validation on incoming messages.
type workflowHandler struct {
	handler    ExecutionStatusUpdateHandler
	primaryDon commoncap.DON
}

// ShardFailoverCommunicator handles d2d communication for shard failover
// across all workflows on a node. It registers a single receiver with the
// dispatcher and routes incoming messages to the appropriate workflow handler
// based on the workflow ID in the message payload. A single instance avoids
// dispatcher conflicts when multiple workflows are deployed.
//
// Primary/secondary DON info is registered per-workflow by each workflow's
// ShardFailoverManager. The communicator itself is role-agnostic — it
// always registers a receiver and can send to any target DON. The manager
// decides at runtime whether to send (primary) or receive (secondary) by
// checking ownership on each event.
type ShardFailoverCommunicator struct {
	services.StateMachine
	stopCh     services.StopChan
	dispatcher remotetypes.Dispatcher
	localDonID uint32
	lggr       logger.SugaredLogger

	mu       sync.RWMutex
	handlers map[string]*workflowHandler

	quorumMu         sync.Mutex
	seenByExecID     map[string]map[p2ptypes.PeerID]bool
	seenAt           map[string]time.Time
	deliveredExecIDs map[string]time.Time
	statusesByExecID map[string][]ringpb.ExecutionStatus
	expiryDuration   time.Duration
	wg               sync.WaitGroup
}

// NewShardFailoverCommunicator creates a communicator that handles
// ExecutionStatusUpdate messages for all workflows on this node.
// The local DON ID is used for receiver registration with the dispatcher.
func NewShardFailoverCommunicator(dispatcher remotetypes.Dispatcher, localDonID uint32, lggr logger.Logger) *ShardFailoverCommunicator {
	return &ShardFailoverCommunicator{
		stopCh:           make(services.StopChan),
		dispatcher:       dispatcher,
		localDonID:       localDonID,
		lggr:             logger.Sugared(logger.With(lggr, "component", "ShardFailoverCommunicator")),
		handlers:         make(map[string]*workflowHandler),
		seenByExecID:     make(map[string]map[p2ptypes.PeerID]bool),
		seenAt:           make(map[string]time.Time),
		deliveredExecIDs: make(map[string]time.Time),
		statusesByExecID: make(map[string][]ringpb.ExecutionStatus),
		expiryDuration:   10 * time.Minute,
	}
}

// RegisterHandler registers a handler for incoming ExecutionStatusUpdate
// messages for the given workflow ID, along with the primary DON used for
// quorum validation (i.e. the DON whose members are allowed to send).
// Called by each workflow's ShardFailoverManager. Multiple workflows can
// register handlers concurrently; the communicator routes by workflow ID
// and validates quorum against the per-workflow primary DON.
func (c *ShardFailoverCommunicator) RegisterHandler(workflowID string, primaryDon commoncap.DON, handler ExecutionStatusUpdateHandler) {
	c.mu.Lock()
	c.handlers[workflowID] = &workflowHandler{handler: handler, primaryDon: primaryDon}
	c.mu.Unlock()
}

// UnregisterHandler removes the handler for the given workflow ID.
func (c *ShardFailoverCommunicator) UnregisterHandler(workflowID string) {
	c.mu.Lock()
	delete(c.handlers, workflowID)
	c.mu.Unlock()
}

func (c *ShardFailoverCommunicator) Start(ctx context.Context) error {
	// Idempotent: if already started (or starting), return nil. The
	// communicator is shared across all workflows — the first manager
	// starts it, subsequent managers just register handlers.
	st := c.State()
	if st == "Started" || st == "Starting" {
		return nil
	}
	return c.StartOnce(c.Name(), func() error {
		if err := c.dispatcher.SetReceiverForMethod(
			ShardExecutionStatusUpdateCapabilityID, c.localDonID,
			remotetypes.MethodExecutionStatusUpdate, c,
		); err != nil {
			return fmt.Errorf("failed to register receiver: %w", err)
		}
		c.wg.Add(1)
		go c.pruneLoop()
		return nil
	})
}

func (c *ShardFailoverCommunicator) Close() error {
	return c.StopOnce(c.Name(), func() error {
		c.dispatcher.RemoveReceiverForMethod(
			ShardExecutionStatusUpdateCapabilityID, c.localDonID,
			remotetypes.MethodExecutionStatusUpdate,
		)
		close(c.stopCh)
		c.wg.Wait()
		return nil
	})
}

// Send sends an ExecutionStatusUpdate to all members of the given target DON.
// The caller (ShardFailoverManager) is responsible for resolving which DON
// to send to — this is a per-workflow decision based on shard ownership.
//
// TODO: Currently the manager resolves the secondary DON per-event in
// forwardExecutionStatus. If shard topology grows beyond two shards, this
// needs to fan out to all secondary DONs for the workflow, not just one.
func (c *ShardFailoverCommunicator) Send(ctx context.Context, msg *ringpb.ExecutionStatusUpdate, targetDon commoncap.DON) {
	if len(targetDon.Members) == 0 {
		c.lggr.Warnw("ShardFailoverCommunicator: target DON has no members, cannot send",
			"workflowID", msg.WorkflowId)
		return
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		c.lggr.Errorw("failed to marshal ExecutionStatusUpdate", "err", err)
		return
	}

	messageID := fmt.Sprintf("%s:%s:%d", msg.WorkflowId, msg.TriggerEventId, msg.TriggerIndex)
	for _, peerID := range targetDon.Members {
		body := &remotetypes.MessageBody{
			CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
			Method:           remotetypes.MethodExecutionStatusUpdate,
			CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
			Payload:          payload,
			CallerDonId:      c.localDonID,
			CapabilityDonId:  targetDon.ID,
			MessageId:        []byte(messageID),
		}
		if err := c.dispatcher.Send(peerID, body); err != nil {
			c.lggr.Errorw("failed to send ExecutionStatusUpdate",
				"peerID", peerID, "err", err)
		}
	}
}

// Receive implements remotetypes.Receiver. It unmarshals the payload to
// extract the workflow ID and execution ID, validates the sender against the
// per-workflow primary DON, collects F+1 quorum by execution ID, and routes
// the message to the registered handler.
//
// When quorum is reached, the handler receives the worst-case status among
// all quorum messages: if any peer reports SYSTEM_ERROR, the aggregated
// result is SYSTEM_ERROR. This ensures the secondary shard replays the
// cached trigger event even if only a subset of primary nodes failed.
func (c *ShardFailoverCommunicator) Receive(ctx context.Context, msg *remotetypes.MessageBody) {
	if msg.Method != remotetypes.MethodExecutionStatusUpdate {
		return
	}

	var execUpdate ringpb.ExecutionStatusUpdate
	if err := proto.Unmarshal(msg.Payload, &execUpdate); err != nil {
		c.lggr.Errorw("failed to unmarshal ExecutionStatusUpdate", "err", err)
		return
	}

	c.mu.RLock()
	wh, ok := c.handlers[execUpdate.WorkflowId]
	c.mu.RUnlock()
	if !ok {
		c.lggr.Debugw("no handler registered for workflow, dropping",
			"workflowID", execUpdate.WorkflowId)
		return
	}

	primary := wh.primaryDon
	if len(primary.Members) == 0 {
		c.lggr.Warnw("ShardFailoverCommunicator: primary DON not set for workflow, dropping message",
			"workflowID", execUpdate.WorkflowId)
		return
	}

	// msg.Sender is the verified peer ID set by the dispatcher after
	// cryptographic validation of the incoming p2p message. It can be
	// trusted as the authentic origin of this ExecutionStatusUpdate.
	sender, err := remote.ToPeerID(msg.Sender)
	if err != nil {
		c.lggr.Errorw("failed to parse sender peer ID", "err", err)
		return
	}

	if !isPeerInDON(sender, primary.Members) {
		c.lggr.Warnw("ExecutionStatusUpdate from peer not in primary shard DON",
			"peerID", sender, "workflowID", execUpdate.WorkflowId)
		return
	}

	// Dedup by execution ID rather than payload hash. Different primary
	// nodes may produce slightly different payloads for the same execution,
	// but the execution ID is guaranteed to match.
	dedupKey := fmt.Sprintf("%s:%s", execUpdate.WorkflowId, execUpdate.ExecutionId)

	c.quorumMu.Lock()
	if _, delivered := c.deliveredExecIDs[dedupKey]; delivered {
		c.quorumMu.Unlock()
		return
	}

	peers, ok := c.seenByExecID[dedupKey]
	if !ok {
		peers = make(map[p2ptypes.PeerID]bool)
		c.seenByExecID[dedupKey] = peers
		c.seenAt[dedupKey] = time.Now()
	}
	if peers[sender] {
		c.quorumMu.Unlock()
		return
	}
	peers[sender] = true

	// Track the status from this peer so we can aggregate at quorum.
	c.recordStatus(dedupKey, sender, execUpdate.Status)

	quorum := int(primary.F) + 1
	reached := len(peers) >= quorum
	if reached {
		c.deliveredExecIDs[dedupKey] = time.Now()
		delete(c.seenByExecID, dedupKey)
		delete(c.seenAt, dedupKey)
	}
	c.quorumMu.Unlock()

	if !reached {
		c.lggr.Debugw("ExecutionStatusUpdate quorum not yet reached",
			"dedupKey", dedupKey, "received", len(peers), "required", quorum)
		return
	}

	// Aggregate: if any quorum peer reported SYSTEM_ERROR, propagate
	// SYSTEM_ERROR so the secondary replays the cached event. Otherwise
	// use the majority status among the received messages.
	aggregated := c.aggregateStatus(dedupKey)
	c.clearStatuses(dedupKey)

	c.lggr.Debugw("ExecutionStatusUpdate quorum reached, routing to handler",
		"workflowID", execUpdate.WorkflowId,
		"triggerEventID", execUpdate.TriggerEventId,
		"executionID", execUpdate.ExecutionId,
		"aggregatedStatus", aggregated,
		"peers", len(peers))

	execUpdate.Status = aggregated
	wh.handler(&execUpdate)
}

// recordStatus stores the status reported by a peer for a given execution.
// Called under c.quorumMu.
func (c *ShardFailoverCommunicator) recordStatus(dedupKey string, _ p2ptypes.PeerID, status ringpb.ExecutionStatus) {
	c.statusesByExecID[dedupKey] = append(c.statusesByExecID[dedupKey], status)
}

// aggregateStatus returns the worst-case status among all received messages
// for the given execution. If any peer reported SYSTEM_ERROR, the result is
// SYSTEM_ERROR — this ensures the secondary shard replays the cached trigger
// event even if only a subset of primary nodes experienced a system failure.
// Otherwise the majority status is returned.
// Called under c.quorumMu.
func (c *ShardFailoverCommunicator) aggregateStatus(dedupKey string) ringpb.ExecutionStatus {
	statuses := c.statusesByExecID[dedupKey]
	counts := make(map[ringpb.ExecutionStatus]int, len(statuses))
	for _, s := range statuses {
		if s == ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR {
			return ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR
		}
		counts[s]++
	}

	var best ringpb.ExecutionStatus
	var bestCount int
	for s, n := range counts {
		if n > bestCount {
			best = s
			bestCount = n
		}
	}
	return best
}

// clearStatuses removes the status tracking for a delivered execution.
// Called under c.quorumMu.
func (c *ShardFailoverCommunicator) clearStatuses(dedupKey string) {
	delete(c.statusesByExecID, dedupKey)
}

func (c *ShardFailoverCommunicator) pruneLoop() {
	defer c.wg.Done()
	ticker := time.NewTicker(c.expiryDuration / 2)
	defer ticker.Stop()

	ctx, cancel := c.stopCh.NewCtx()
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.quorumMu.Lock()
			now := time.Now()
			for key, deliveredAt := range c.deliveredExecIDs {
				if now.Sub(deliveredAt) >= c.expiryDuration {
					delete(c.deliveredExecIDs, key)
				}
			}
			for key, seenTime := range c.seenAt {
				if now.Sub(seenTime) >= c.expiryDuration {
					delete(c.seenByExecID, key)
					delete(c.seenAt, key)
					delete(c.statusesByExecID, key)
				}
			}
			c.quorumMu.Unlock()
		}
	}
}

func (c *ShardFailoverCommunicator) Name() string {
	return fmt.Sprintf("ShardFailoverCommunicator-don-%d", c.localDonID)
}

func (c *ShardFailoverCommunicator) HealthReport() map[string]error {
	return map[string]error{c.Name(): c.Healthy()}
}

func isPeerInDON(peer p2ptypes.PeerID, members []p2ptypes.PeerID) bool {
	return slices.Contains(members, peer)
}
