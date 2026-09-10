package sharding

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

// ShardHeartbeatCapabilityID is the capability ID used for routing
// ShardHeartbeat messages between shards via the d2d dispatcher.
const ShardHeartbeatCapabilityID = "shard-heartbeat"

// ShardFailoverCommunicator handles d2d communication for shard failover
// across all workflows on a node. It registers a single receiver with the
// dispatcher and routes incoming messages to the appropriate workflow handler
// based on the workflow ID in the message payload. A single instance avoids
// dispatcher conflicts when multiple workflows are deployed.
//
// The communicator is role-agnostic: it always registers a receiver and can
// always send to the secondary shard. The ShardFailoverManager decides at
// runtime whether to send (primary) or receive (secondary) by checking
// ownership on each event — no static role assignment or re-wiring needed.
type ShardFailoverCommunicator struct {
	services.StateMachine
	stopCh     services.StopChan
	dispatcher remotetypes.Dispatcher
	localDonID uint32
	lggr       logger.SugaredLogger

	mu           sync.RWMutex
	primaryDon   commoncap.DON
	secondaryDon commoncap.DON
	handlers     map[string]ExecutionStatusUpdateHandler

	quorumMu        sync.Mutex
	seenByHash      map[string]map[p2ptypes.PeerID]bool
	seenAt          map[string]time.Time
	deliveredHashes map[string]time.Time
	expiryDuration  time.Duration
	wg              sync.WaitGroup
}

// NewShardFailoverCommunicator creates a communicator that handles
// ExecutionStatusUpdate messages for all workflows on this node.
// The local DON ID is used for receiver registration with the dispatcher.
// Call SetShardDons before Start so quorum checking and sending work correctly.
func NewShardFailoverCommunicator(dispatcher remotetypes.Dispatcher, localDonID uint32, lggr logger.Logger) *ShardFailoverCommunicator {
	return &ShardFailoverCommunicator{
		stopCh:          make(services.StopChan),
		dispatcher:      dispatcher,
		localDonID:      localDonID,
		lggr:            logger.Sugared(logger.With(lggr, "component", "ShardFailoverCommunicator")),
		handlers:        make(map[string]ExecutionStatusUpdateHandler),
		seenByHash:      make(map[string]map[p2ptypes.PeerID]bool),
		seenAt:          make(map[string]time.Time),
		deliveredHashes: make(map[string]time.Time),
		expiryDuration:  10 * time.Minute,
	}
}

// SetShardDons provides the primary and secondary shard DON info.
// primaryDon is used for quorum validation on Receive (who to accept from).
// secondaryDon is used for Send (where to send to).
// Must be called before Start. Can be called again at runtime to update
// when shard config changes.
func (c *ShardFailoverCommunicator) SetShardDons(primaryDon, secondaryDon commoncap.DON) {
	c.mu.Lock()
	c.primaryDon = primaryDon
	c.secondaryDon = secondaryDon
	c.mu.Unlock()
}

// RegisterHandler registers a handler for incoming ExecutionStatusUpdate
// messages for the given workflow ID. Called by the secondary shard's
// ShardFailoverManager. Multiple workflows can register handlers
// concurrently; the communicator routes by workflow ID.
func (c *ShardFailoverCommunicator) RegisterHandler(workflowID string, handler ExecutionStatusUpdateHandler) {
	c.mu.Lock()
	c.handlers[workflowID] = handler
	c.mu.Unlock()
}

// UnregisterHandler removes the handler for the given workflow ID.
func (c *ShardFailoverCommunicator) UnregisterHandler(workflowID string) {
	c.mu.Lock()
	delete(c.handlers, workflowID)
	c.mu.Unlock()
}

func (c *ShardFailoverCommunicator) Start(ctx context.Context) error {
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

// Send sends an ExecutionStatusUpdate to all members of the secondary shard DON.
// Called by the primary shard's ShardFailoverManager.
func (c *ShardFailoverCommunicator) Send(ctx context.Context, msg *ringpb.ExecutionStatusUpdate) {
	c.mu.RLock()
	secondary := c.secondaryDon
	c.mu.RUnlock()

	if len(secondary.Members) == 0 {
		c.lggr.Warnw("ShardFailoverCommunicator: secondary DON not set, cannot send",
			"workflowID", msg.WorkflowId)
		return
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		c.lggr.Errorw("failed to marshal ExecutionStatusUpdate", "err", err)
		return
	}

	messageID := fmt.Sprintf("%s:%s:%d", msg.WorkflowId, msg.TriggerEventId, msg.TriggerIndex)
	for _, peerID := range secondary.Members {
		body := &remotetypes.MessageBody{
			CapabilityId:     ShardExecutionStatusUpdateCapabilityID,
			Method:           remotetypes.MethodExecutionStatusUpdate,
			CapabilityMethod: remotetypes.MethodExecutionStatusUpdate,
			Payload:          payload,
			CallerDonId:      c.localDonID,
			CapabilityDonId:  secondary.ID,
			MessageId:        []byte(messageID),
		}
		if err := c.dispatcher.Send(peerID, body); err != nil {
			c.lggr.Errorw("failed to send ExecutionStatusUpdate",
				"peerID", peerID, "err", err)
		}
	}
}

// Receive implements remotetypes.Receiver. It validates the sender against the
// primary DON, collects quorum, and routes the unmarshalled message to the
// handler registered for the workflow ID in the message.
func (c *ShardFailoverCommunicator) Receive(ctx context.Context, msg *remotetypes.MessageBody) {
	if msg.Method != remotetypes.MethodExecutionStatusUpdate {
		return
	}

	c.mu.RLock()
	primary := c.primaryDon
	c.mu.RUnlock()

	if len(primary.Members) == 0 {
		c.lggr.Warnw("ShardFailoverCommunicator: primary DON not set, dropping message")
		return
	}

	sender, err := remote.ToPeerID(msg.Sender)
	if err != nil {
		c.lggr.Errorw("failed to parse sender peer ID", "err", err)
		return
	}

	if !isPeerInDON(sender, primary.Members) {
		c.lggr.Warnw("ExecutionStatusUpdate from peer not in primary shard DON", "peerID", sender)
		return
	}

	hash := sha256.Sum256(msg.Payload)
	hashKey := hex.EncodeToString(hash[:])

	c.quorumMu.Lock()
	if _, ok := c.deliveredHashes[hashKey]; ok {
		c.quorumMu.Unlock()
		return
	}

	peers, ok := c.seenByHash[hashKey]
	if !ok {
		peers = make(map[p2ptypes.PeerID]bool)
		c.seenByHash[hashKey] = peers
		c.seenAt[hashKey] = time.Now()
	}
	if peers[sender] {
		c.quorumMu.Unlock()
		return
	}
	peers[sender] = true

	quorum := int(primary.F) + 1
	reached := len(peers) >= quorum
	if reached {
		c.deliveredHashes[hashKey] = time.Now()
		delete(c.seenByHash, hashKey)
		delete(c.seenAt, hashKey)
	}
	c.quorumMu.Unlock()

	if !reached {
		c.lggr.Debugw("ExecutionStatusUpdate quorum not yet reached",
			"hash", hashKey, "received", len(peers), "required", quorum)
		return
	}

	var execUpdate ringpb.ExecutionStatusUpdate
	if err := proto.Unmarshal(msg.Payload, &execUpdate); err != nil {
		c.lggr.Errorw("failed to unmarshal ExecutionStatusUpdate", "err", err)
		return
	}

	c.lggr.Infow("ExecutionStatusUpdate quorum reached, routing to handler",
		"workflowID", execUpdate.WorkflowId,
		"triggerEventID", execUpdate.TriggerEventId,
		"status", execUpdate.Status,
		"peers", len(peers))

	c.mu.RLock()
	handler, ok := c.handlers[execUpdate.WorkflowId]
	c.mu.RUnlock()
	if !ok {
		c.lggr.Debugw("no handler registered for workflow, dropping",
			"workflowID", execUpdate.WorkflowId)
		return
	}
	handler(&execUpdate)
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
			for hash, deliveredAt := range c.deliveredHashes {
				if now.Sub(deliveredAt) >= c.expiryDuration {
					delete(c.deliveredHashes, hash)
				}
			}
			for hash, seenTime := range c.seenAt {
				if now.Sub(seenTime) >= c.expiryDuration {
					delete(c.seenByHash, hash)
					delete(c.seenAt, hash)
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
