package sharding

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	shardingv1 "github.com/smartcontractkit/chainlink-protos/cre/go/sharding/v1"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type ExecutionCompletedHandler func(msg *shardingv1.ExecutionCompleted)

type ExecutionCompletedReceiver struct {
	services.StateMachine
	stopCh  services.StopChan
	primary commoncap.DON
	handler ExecutionCompletedHandler
	lggr    logger.SugaredLogger
	wg      sync.WaitGroup

	mu              sync.Mutex
	seenByHash      map[string]map[p2ptypes.PeerID]bool
	deliveredHashes map[string]bool
	expiryDuration  time.Duration
}

func NewExecutionCompletedReceiver(primary commoncap.DON, handler ExecutionCompletedHandler, lggr logger.Logger) *ExecutionCompletedReceiver {
	return &ExecutionCompletedReceiver{
		stopCh:          make(services.StopChan),
		primary:         primary,
		handler:         handler,
		lggr:            logger.Sugared(logger.With(lggr, "component", "ExecutionCompletedReceiver")),
		seenByHash:      make(map[string]map[p2ptypes.PeerID]bool),
		deliveredHashes: make(map[string]bool),
		expiryDuration:  10 * time.Minute,
	}
}

func (r *ExecutionCompletedReceiver) Start(ctx context.Context) error {
	return r.StartOnce(r.Name(), func() error {
		r.wg.Add(1)
		go r.pruneLoop()
		return nil
	})
}

func (r *ExecutionCompletedReceiver) Close() error {
	return r.StopOnce(r.Name(), func() error {
		close(r.stopCh)
		r.wg.Wait()
		return nil
	})
}

func (r *ExecutionCompletedReceiver) pruneLoop() {
	defer r.wg.Done()
	ticker := time.NewTicker(r.expiryDuration / 2)
	defer ticker.Stop()

	ctx, cancel := r.stopCh.NewCtx()
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.mu.Lock()
			r.seenByHash = make(map[string]map[p2ptypes.PeerID]bool)
			r.deliveredHashes = make(map[string]bool)
			r.mu.Unlock()
		}
	}
}

func (r *ExecutionCompletedReceiver) Receive(ctx context.Context, msg *remotetypes.MessageBody) {
	if msg.Method != remotetypes.MethodExecutionCompleted {
		return
	}

	sender, err := remote.ToPeerID(msg.Sender)
	if err != nil {
		r.lggr.Errorw("failed to parse sender peer ID", "err", err)
		return
	}

	if !isPeerInDON(sender, r.primary.Members) {
		r.lggr.Warnw("ExecutionCompleted from peer not in primary shard DON", "peerID", sender)
		return
	}

	hash := sha256.Sum256(msg.Payload)
	hashKey := hex.EncodeToString(hash[:])

	r.mu.Lock()
	if r.deliveredHashes[hashKey] {
		r.mu.Unlock()
		return
	}

	peers, ok := r.seenByHash[hashKey]
	if !ok {
		peers = make(map[p2ptypes.PeerID]bool)
		r.seenByHash[hashKey] = peers
	}
	if peers[sender] {
		r.mu.Unlock()
		return
	}
	peers[sender] = true

	quorum := int(r.primary.F) + 1
	reached := len(peers) >= quorum
	if reached {
		r.deliveredHashes[hashKey] = true
	}
	r.mu.Unlock()

	if !reached {
		r.lggr.Debugw("ExecutionCompleted quorum not yet reached",
			"hash", hashKey, "received", len(peers), "required", quorum)
		return
	}

	var execCompleted shardingv1.ExecutionCompleted
	if err := proto.Unmarshal(msg.Payload, &execCompleted); err != nil {
		r.lggr.Errorw("failed to unmarshal ExecutionCompleted", "err", err)
		return
	}

	r.lggr.Infow("ExecutionCompleted quorum reached, invoking handler",
		"workflowID", execCompleted.WorkflowId,
		"triggerEventID", execCompleted.TriggerEventId,
		"status", execCompleted.Status,
		"peers", len(peers))

	r.handler(&execCompleted)
}

func (r *ExecutionCompletedReceiver) Name() string {
	return fmt.Sprintf("ExecutionCompletedReceiver-primary-%d", r.primary.ID)
}

func (r *ExecutionCompletedReceiver) HealthReport() map[string]error {
	return map[string]error{r.Name(): r.Healthy()}
}

type ShardHeartbeatHandler func(msg *shardingv1.ShardHeartbeat)

type ShardHeartbeatReceiver struct {
	services.StateMachine
	stopCh  services.StopChan
	primary commoncap.DON
	handler ShardHeartbeatHandler
	lggr    logger.SugaredLogger

	mu         sync.Mutex
	seenByHash map[string]map[p2ptypes.PeerID]bool
	lastSeen   int64
}

func NewShardHeartbeatReceiver(primary commoncap.DON, handler ShardHeartbeatHandler, lggr logger.Logger) *ShardHeartbeatReceiver {
	return &ShardHeartbeatReceiver{
		stopCh:     make(services.StopChan),
		primary:    primary,
		handler:    handler,
		lggr:       logger.Sugared(logger.With(lggr, "component", "ShardHeartbeatReceiver")),
		seenByHash: make(map[string]map[p2ptypes.PeerID]bool),
	}
}

func (r *ShardHeartbeatReceiver) Start(ctx context.Context) error {
	return r.StartOnce(r.Name(), func() error { return nil })
}

func (r *ShardHeartbeatReceiver) Close() error {
	return r.StopOnce(r.Name(), func() error {
		close(r.stopCh)
		return nil
	})
}

func (r *ShardHeartbeatReceiver) Receive(ctx context.Context, msg *remotetypes.MessageBody) {
	if msg.Method != remotetypes.MethodShardHeartbeat {
		return
	}

	sender, err := remote.ToPeerID(msg.Sender)
	if err != nil {
		return
	}

	if !isPeerInDON(sender, r.primary.Members) {
		return
	}

	hash := sha256.Sum256(msg.Payload)
	hashKey := hex.EncodeToString(hash[:])

	r.mu.Lock()
	peers, ok := r.seenByHash[hashKey]
	if !ok {
		peers = make(map[p2ptypes.PeerID]bool)
		r.seenByHash[hashKey] = peers
	}
	if peers[sender] {
		r.mu.Unlock()
		return
	}
	peers[sender] = true

	quorum := int(r.primary.F) + 1
	reached := len(peers) >= quorum
	r.mu.Unlock()

	if !reached {
		return
	}

	var hb shardingv1.ShardHeartbeat
	if err := proto.Unmarshal(msg.Payload, &hb); err != nil {
		r.lggr.Errorw("failed to unmarshal ShardHeartbeat", "err", err)
		return
	}

	r.mu.Lock()
	r.lastSeen = time.Now().Unix()
	r.seenByHash = make(map[string]map[p2ptypes.PeerID]bool)
	r.mu.Unlock()

	r.handler(&hb)
}

func (r *ShardHeartbeatReceiver) LastSeen() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lastSeen
}

func (r *ShardHeartbeatReceiver) Name() string {
	return fmt.Sprintf("ShardHeartbeatReceiver-primary-%d", r.primary.ID)
}

func (r *ShardHeartbeatReceiver) HealthReport() map[string]error {
	return map[string]error{r.Name(): r.Healthy()}
}

func isPeerInDON(peer p2ptypes.PeerID, members []p2ptypes.PeerID) bool {
	return slices.Contains(members, peer)
}
