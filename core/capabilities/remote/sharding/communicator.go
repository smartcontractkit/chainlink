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

	quorumMu             sync.Mutex
	seenByExecID         map[string]map[p2ptypes.PeerID]bool
	seenAt               map[string]time.Time
	deliveredExecIDs     map[string]time.Time
	statusCountsByExecID map[string]map[ringpb.ExecutionStatus]int
	expiryDuration       time.Duration
	wg                   sync.WaitGroup
}

// NewShardFailoverCommunicator creates a communicator that handles
// ExecutionStatusUpdate messages for all workflows on this node.
// The local DON ID is used for receiver registration with the dispatcher.
func NewShardFailoverCommunicator(dispatcher remotetypes.Dispatcher, localDonID uint32, lggr logger.Logger) *ShardFailoverCommunicator {
	return &ShardFailoverCommunicator{
		stopCh:               make(services.StopChan),
		dispatcher:           dispatcher,
		localDonID:           localDonID,
		lggr:                 logger.Sugared(logger.With(lggr, "component", "ShardFailoverCommunicator")),
		handlers:             make(map[string]*workflowHandler),
		seenByExecID:         make(map[string]map[p2ptypes.PeerID]bool),
		seenAt:               make(map[string]time.Time),
		deliveredExecIDs:     make(map[string]time.Time),
		statusCountsByExecID: make(map[string]map[ringpb.ExecutionStatus]int),
		expiryDuration:       10 * time.Minute,
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
// per-workflow primary DON, and collects F+1 matching reports by execution
// ID before routing to the registered handler.
//
// Quorum is per-status (like OCR): a status is delivered only when F+1
// peers agree on it. This ensures F faulty nodes cannot force a wrong
// outcome. If SYSTEM_ERROR reaches F+1, the handler triggers failover
// replay; if SUCCESS or USER_ERROR reaches F+1 first, the cached event
// is drained.
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

	dedupKey := fmt.Sprintf("%s:%s", execUpdate.WorkflowId, execUpdate.ExecutionId)
	quorum := 2*int(primary.F) + 1

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
		c.statusCountsByExecID[dedupKey] = make(map[ringpb.ExecutionStatus]int)
	}
	if peers[sender] {
		c.quorumMu.Unlock()
		return
	}
	peers[sender] = true
	c.statusCountsByExecID[dedupKey][execUpdate.Status]++

	// Check if any status has reached F+1 quorum.
	statusCounts := c.statusCountsByExecID[dedupKey]
	var quorumStatus ringpb.ExecutionStatus
	var reached bool
	for status, count := range statusCounts {
		if count >= quorum {
			quorumStatus = status
			reached = true
			break
		}
	}

	if reached {
		c.deliveredExecIDs[dedupKey] = time.Now()
		delete(c.seenByExecID, dedupKey)
		delete(c.seenAt, dedupKey)
		delete(c.statusCountsByExecID, dedupKey)
	}
	c.quorumMu.Unlock()

	if !reached {
		c.lggr.Debugw("ExecutionStatusUpdate quorum not yet reached",
			"dedupKey", dedupKey, "statusCounts", len(statusCounts), "peers", len(peers), "required", quorum)
		return
	}

	c.lggr.Debugw("ExecutionStatusUpdate quorum reached, routing to handler",
		"workflowID", execUpdate.WorkflowId,
		"triggerEventID", execUpdate.TriggerEventId,
		"executionID", execUpdate.ExecutionId,
		"status", quorumStatus,
		"peers", len(peers))

	execUpdate.Status = quorumStatus
	wh.handler(&execUpdate)
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
					delete(c.statusCountsByExecID, key)
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
