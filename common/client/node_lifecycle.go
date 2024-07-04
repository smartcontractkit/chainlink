package client

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/v2/common/types"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	bigmath "github.com/smartcontractkit/chainlink-common/pkg/utils/big_math"

	iutils "github.com/smartcontractkit/chainlink/v2/common/internal/utils"
)

var (
	promPoolRPCNodeHighestSeenBlock = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pool_rpc_node_highest_seen_block",
		Help: "The highest seen block for the given RPC node",
	}, []string{"chainID", "nodeName"})
	promPoolRPCNodeHighestFinalizedBlock = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pool_rpc_node_highest_finalized_block",
		Help: "The highest seen finalized block for the given RPC node",
	}, []string{"chainID", "nodeName"})
	promPoolRPCNodeNumSeenBlocks = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_seen_blocks",
		Help: "The total number of new blocks seen by the given RPC node",
	}, []string{"chainID", "nodeName"})
	promPoolRPCNodePolls = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_polls_total",
		Help: "The total number of poll checks for the given RPC node",
	}, []string{"chainID", "nodeName"})
	promPoolRPCNodePollsFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_polls_failed",
		Help: "The total number of failed poll checks for the given RPC node",
	}, []string{"chainID", "nodeName"})
	promPoolRPCNodePollsSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_polls_success",
		Help: "The total number of successful poll checks for the given RPC node",
	}, []string{"chainID", "nodeName"})
)

// zombieNodeCheckInterval controls how often to re-check to see if we need to
// state change in case we have to force a state transition due to no available
// nodes.
// NOTE: This only applies to out-of-sync nodes if they are the last available node
func zombieNodeCheckInterval(noNewHeadsThreshold time.Duration) time.Duration {
	interval := noNewHeadsThreshold
	if interval <= 0 || interval > QueryTimeout {
		interval = QueryTimeout
	}
	return utils.WithJitter(interval)
}

const (
	msgCannotDisable = "but cannot disable this connection because there are no other RPC endpoints, or all other RPC endpoints are dead."
	msgDegradedState = "Chainlink is now operating in a degraded state and urgent action is required to resolve the issue"
)

// Node is a FSM
// Each state has a loop that goes with it, which monitors the node and moves it into another state as necessary.
// Only one loop must run at a time.
// Each loop passes control onto the next loop as it exits, except when the node is Closed which terminates the loop permanently.

// This handles node lifecycle for the ALIVE state
// Should only be run ONCE per node, after a successful Dial
func (n *node[CHAIN_ID, HEAD, RPC]) aliveLoop() {
	defer n.wg.Done()
	ctx, cancel := n.newCtx()
	defer cancel()

	{
		// sanity check
		state := n.getCachedState()
		switch state {
		case NodeStateAlive:
		case NodeStateClosed:
			return
		default:
			panic(fmt.Sprintf("aliveLoop can only run for node in Alive state, got: %s", state))
		}
	}

	noNewHeadsTimeoutThreshold := n.chainCfg.NodeNoNewHeadsThreshold()
	pollFailureThreshold := n.nodePoolCfg.PollFailureThreshold()
	pollInterval := n.nodePoolCfg.PollInterval()

	lggr := logger.Sugared(n.lfcLog).Named("Alive").With("noNewHeadsTimeoutThreshold", noNewHeadsTimeoutThreshold, "pollInterval", pollInterval, "pollFailureThreshold", pollFailureThreshold)
	lggr.Tracew("Alive loop starting", "nodeState", n.getCachedState())

	headsC, sub, err := n.rpc.SubscribeToHeads(ctx)
	if err != nil {
		lggr.Errorw("Initial subscribe for heads failed", "nodeState", n.getCachedState())
		n.declareUnreachable()
		return
	}

	n.stateMu.Lock()
	n.aliveLoopSub = sub
	n.stateMu.Unlock()
	defer sub.Unsubscribe()

	var outOfSyncT *time.Ticker
	var outOfSyncTC <-chan time.Time
	if noNewHeadsTimeoutThreshold > 0 {
		lggr.Debugw("Head liveness checking enabled", "nodeState", n.getCachedState())
		outOfSyncT = time.NewTicker(noNewHeadsTimeoutThreshold)
		defer outOfSyncT.Stop()
		outOfSyncTC = outOfSyncT.C
	} else {
		lggr.Debug("Head liveness checking disabled")
	}

	var pollCh <-chan time.Time
	if pollInterval > 0 {
		lggr.Debug("Polling enabled")
		pollT := time.NewTicker(pollInterval)
		defer pollT.Stop()
		pollCh = pollT.C
		if pollFailureThreshold > 0 {
			// polling can be enabled with no threshold to enable polling but
			// the node will not be marked offline regardless of the number of
			// poll failures
			lggr.Debug("Polling liveness checking enabled")
		}
	} else {
		lggr.Debug("Polling disabled")
	}

	var finalizedHeadCh <-chan HEAD
	var finalizedHeadSub types.Subscription
	if n.chainCfg.FinalityTagEnabled() {
		lggr.Debugw("Finalized block polling enabled")
		finalizedHeadCh, finalizedHeadSub, err = n.rpc.SubscribeToFinalizedHeads(ctx)
		if err != nil {
			lggr.Errorw("Failed to subscribe to finalized heads", "err", err)
			n.declareUnreachable()
			return
		}
		defer finalizedHeadSub.Unsubscribe()
	}
	n.finalizedBlockSub = finalizedHeadSub

	localHighestChainInfo, _ := n.rpc.GetInterceptedChainInfo()
	var pollFailures uint32

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollCh:
			promPoolRPCNodePolls.WithLabelValues(n.chainID.String(), n.name).Inc()
			lggr.Tracew("Pinging RPC", "nodeState", n.State(), "pollFailures", pollFailures)
			pollCtx, cancel := context.WithTimeout(ctx, pollInterval)
			err := n.RPC().Ping(pollCtx)
			cancel()
			if err != nil {
				// prevent overflow
				if pollFailures < math.MaxUint32 {
					promPoolRPCNodePollsFailed.WithLabelValues(n.chainID.String(), n.name).Inc()
					pollFailures++
				}
				lggr.Warnw(fmt.Sprintf("Poll failure, RPC endpoint %s failed to respond properly", n.String()), "err", err, "pollFailures", pollFailures, "nodeState", n.getCachedState())
			} else {
				lggr.Debugw("Ping successful", "nodeState", n.State())
				promPoolRPCNodePollsSuccess.WithLabelValues(n.chainID.String(), n.name).Inc()
				pollFailures = 0
			}
			if pollFailureThreshold > 0 && pollFailures >= pollFailureThreshold {
				lggr.Errorw(fmt.Sprintf("RPC endpoint failed to respond to %d consecutive polls", pollFailures), "pollFailures", pollFailures, "nodeState", n.getCachedState())
				if n.poolInfoProvider != nil {
					if l, _ := n.poolInfoProvider.LatestChainInfo(); l < 2 {
						lggr.Criticalf("RPC endpoint failed to respond to polls; %s %s", msgCannotDisable, msgDegradedState)
						continue
					}
				}
				n.declareUnreachable()
				return
			}
			_, ci := n.StateAndLatest()
			if outOfSync, liveNodes := n.syncStatus(ci.BlockNumber, ci.TotalDifficulty); outOfSync {
				// note: there must be another live node for us to be out of sync
				lggr.Errorw("RPC endpoint has fallen behind", "blockNumber", ci.BlockNumber, "totalDifficulty", ci.TotalDifficulty, "nodeState", n.getCachedState())
				if liveNodes < 2 {
					lggr.Criticalf("RPC endpoint has fallen behind; %s %s", msgCannotDisable, msgDegradedState)
					continue
				}
				n.declareOutOfSync(n.isOutOfSync)
				return
			}
		case bh, open := <-headsC:
			if !open {
				lggr.Errorw("Subscription channel unexpectedly closed", "nodeState", n.getCachedState())
				n.declareUnreachable()
				return
			}
			promPoolRPCNodeNumSeenBlocks.WithLabelValues(n.chainID.String(), n.name).Inc()
			lggr.Tracew("Got head", "head", bh)
			if bh.BlockNumber() > localHighestChainInfo.BlockNumber {
				promPoolRPCNodeHighestSeenBlock.WithLabelValues(n.chainID.String(), n.name).Set(float64(bh.BlockNumber()))
				lggr.Tracew("Got higher block number, resetting timer", "latestReceivedBlockNumber", localHighestChainInfo.BlockNumber, "blockNumber", bh.BlockNumber(), "nodeState", n.getCachedState())
				localHighestChainInfo.BlockNumber = bh.BlockNumber()
			} else {
				lggr.Tracew("Ignoring previously seen block number", "latestReceivedBlockNumber", localHighestChainInfo.BlockNumber, "blockNumber", bh.BlockNumber(), "nodeState", n.getCachedState())
			}
			if outOfSyncT != nil {
				outOfSyncT.Reset(noNewHeadsTimeoutThreshold)
			}
			if !n.chainCfg.FinalityTagEnabled() {
				latestFinalizedBN := max(bh.BlockNumber()-int64(n.chainCfg.FinalityDepth()), 0)
				if latestFinalizedBN > localHighestChainInfo.FinalizedBlockNumber {
					promPoolRPCNodeHighestFinalizedBlock.WithLabelValues(n.chainID.String(), n.name).Set(float64(latestFinalizedBN))
					localHighestChainInfo.FinalizedBlockNumber = latestFinalizedBN
				}
			}
		case err := <-sub.Err():
			lggr.Errorw("Subscription was terminated", "err", err, "nodeState", n.getCachedState())
			n.declareUnreachable()
			return
		case <-outOfSyncTC:
			// We haven't received a head on the channel for at least the
			// threshold amount of time, mark it broken
			lggr.Errorw(fmt.Sprintf("RPC endpoint detected out of sync; no new heads received for %s (last head received was %v)", noNewHeadsTimeoutThreshold, localHighestChainInfo.BlockNumber), "nodeState", n.getCachedState(), "latestReceivedBlockNumber", localHighestChainInfo.BlockNumber, "noNewHeadsTimeoutThreshold", noNewHeadsTimeoutThreshold)
			if n.poolInfoProvider != nil {
				if l, _ := n.poolInfoProvider.LatestChainInfo(); l < 2 {
					lggr.Criticalf("RPC endpoint detected out of sync; %s %s", msgCannotDisable, msgDegradedState)
					// We don't necessarily want to wait the full timeout to check again, we should
					// check regularly and log noisily in this state
					outOfSyncT.Reset(zombieNodeCheckInterval(noNewHeadsTimeoutThreshold))
					continue
				}
			}
			n.declareOutOfSync(func(num int64, td *big.Int) bool { return num < localHighestChainInfo.BlockNumber })
			return
		case latestFinalized, open := <-finalizedHeadCh:
			if !open {
				lggr.Errorw("Subscription channel unexpectedly closed", "nodeState", n.State())
				n.declareUnreachable()
				return
			}
			if !latestFinalized.IsValid() {
				lggr.Warn("Latest finalized block is not valid")
				continue
			}

			n.stateMu.Lock()
			latestFinalizedBN := latestFinalized.BlockNumber()
			if latestFinalizedBN > localHighestChainInfo.FinalizedBlockNumber {
				promPoolRPCNodeHighestFinalizedBlock.WithLabelValues(n.chainID.String(), n.name).Set(float64(latestFinalizedBN))
				localHighestChainInfo.FinalizedBlockNumber = latestFinalizedBN
			}
			n.stateMu.Unlock()
		}
	}
}

func (n *node[CHAIN_ID, HEAD, RPC]) isOutOfSync(num int64, td *big.Int) (outOfSync bool) {
	outOfSync, _ = n.syncStatus(num, td)
	return
}

// syncStatus returns outOfSync true if num or td is more than SyncThresold behind the best node.
// Always returns outOfSync false for SyncThreshold 0.
// liveNodes is only included when outOfSync is true.
func (n *node[CHAIN_ID, HEAD, RPC]) syncStatus(num int64, td *big.Int) (outOfSync bool, liveNodes int) {
	if n.poolInfoProvider == nil {
		return // skip for tests
	}
	threshold := n.nodePoolCfg.SyncThreshold()
	if threshold == 0 {
		return // disabled
	}
	// Check against best node
	ln, ci := n.poolInfoProvider.LatestChainInfo()
	mode := n.nodePoolCfg.SelectionMode()
	switch mode {
	case NodeSelectionModeHighestHead, NodeSelectionModeRoundRobin, NodeSelectionModePriorityLevel:
		return num < ci.BlockNumber-int64(threshold), ln
	case NodeSelectionModeTotalDifficulty:
		bigThreshold := big.NewInt(int64(threshold))
		return td.Cmp(bigmath.Sub(ci.TotalDifficulty, bigThreshold)) < 0, ln
	default:
		panic("unrecognized NodeSelectionMode: " + mode)
	}
}

const (
	msgReceivedBlock = "Received block for RPC node, waiting until back in-sync to mark as live again"
	msgInSync        = "RPC node back in sync"
)

// outOfSyncLoop takes an OutOfSync node and waits until isOutOfSync returns false to go back to live status
func (n *node[CHAIN_ID, HEAD, RPC]) outOfSyncLoop(isOutOfSync func(num int64, td *big.Int) bool) {
	defer n.wg.Done()
	ctx, cancel := n.newCtx()
	defer cancel()

	{
		// sanity check
		state := n.getCachedState()
		switch state {
		case NodeStateOutOfSync:
		case NodeStateClosed:
			return
		default:
			panic(fmt.Sprintf("outOfSyncLoop can only run for node in OutOfSync state, got: %s", state))
		}
	}

	outOfSyncAt := time.Now()

	lggr := logger.Sugared(logger.Named(n.lfcLog, "OutOfSync"))
	lggr.Debugw("Trying to revive out-of-sync RPC node", "nodeState", n.getCachedState())

	// Need to redial since out-of-sync nodes are automatically disconnected
	state := n.createVerifiedConn(ctx, lggr)
	if state != NodeStateAlive {
		n.declareState(state)
		return
	}

	lggr.Tracew("Successfully subscribed to heads feed on out-of-sync RPC node", "nodeState", n.getCachedState())

	ch, sub, err := n.rpc.SubscribeToHeads(ctx)
	if err != nil {
		lggr.Errorw("Failed to subscribe heads on out-of-sync RPC node", "nodeState", n.getCachedState(), "err", err)
		n.declareUnreachable()
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return
		case head, open := <-ch:
			if !open {
				lggr.Error("Subscription channel unexpectedly closed", "nodeState", n.getCachedState())
				n.declareUnreachable()
				return
			}
			if !isOutOfSync(head.BlockNumber(), head.BlockDifficulty()) {
				// back in-sync! flip back into alive loop
				lggr.Infow(fmt.Sprintf("%s: %s. Node was out-of-sync for %s", msgInSync, n.String(), time.Since(outOfSyncAt)), "blockNumber", head.BlockNumber(), "blockDifficulty", head.BlockDifficulty(), "nodeState", n.getCachedState())
				n.declareInSync()
				return
			}
			lggr.Debugw(msgReceivedBlock, "blockNumber", head.BlockNumber(), "blockDifficulty", head.BlockDifficulty(), "nodeState", n.getCachedState())
		case <-time.After(zombieNodeCheckInterval(n.chainCfg.NodeNoNewHeadsThreshold())):
			if n.poolInfoProvider != nil {
				if l, _ := n.poolInfoProvider.LatestChainInfo(); l < 1 {
					lggr.Critical("RPC endpoint is still out of sync, but there are no other available nodes. This RPC node will be forcibly moved back into the live pool in a degraded state")
					n.declareInSync()
					return
				}
			}
		case err := <-sub.Err():
			lggr.Errorw("Subscription was terminated", "nodeState", n.getCachedState(), "err", err)
			n.declareUnreachable()
			return
		}
	}
}

func (n *node[CHAIN_ID, HEAD, RPC]) unreachableLoop() {
	defer n.wg.Done()
	ctx, cancel := n.newCtx()
	defer cancel()

	{
		// sanity check
		state := n.getCachedState()
		switch state {
		case NodeStateUnreachable:
		case NodeStateClosed:
			return
		default:
			panic(fmt.Sprintf("unreachableLoop can only run for node in Unreachable state, got: %s", state))
		}
	}

	unreachableAt := time.Now()

	lggr := logger.Sugared(logger.Named(n.lfcLog, "Unreachable"))
	lggr.Debugw("Trying to revive unreachable RPC node", "nodeState", n.getCachedState())

	dialRetryBackoff := iutils.NewRedialBackoff()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(dialRetryBackoff.Duration()):
			lggr.Tracew("Trying to re-dial RPC node", "nodeState", n.getCachedState())

			err := n.rpc.Dial(ctx)
			if err != nil {
				lggr.Errorw(fmt.Sprintf("Failed to redial RPC node; still unreachable: %v", err), "err", err, "nodeState", n.getCachedState())
				continue
			}

			n.setState(NodeStateDialed)

			state := n.verifyConn(ctx, lggr)
			switch state {
			case NodeStateUnreachable:
				n.setState(NodeStateUnreachable)
				continue
			case NodeStateAlive:
				lggr.Infow(fmt.Sprintf("Successfully redialled and verified RPC node %s. Node was offline for %s", n.String(), time.Since(unreachableAt)), "nodeState", n.getCachedState())
				fallthrough
			default:
				n.declareState(state)
				return
			}
		}
	}
}

func (n *node[CHAIN_ID, HEAD, RPC]) invalidChainIDLoop() {
	defer n.wg.Done()
	ctx, cancel := n.newCtx()
	defer cancel()

	{
		// sanity check
		state := n.getCachedState()
		switch state {
		case NodeStateInvalidChainID:
		case NodeStateClosed:
			return
		default:
			panic(fmt.Sprintf("invalidChainIDLoop can only run for node in InvalidChainID state, got: %s", state))
		}
	}

	fmt.Println("invalidChainIDLoop")

	invalidAt := time.Now()

	lggr := logger.Named(n.lfcLog, "InvalidChainID")

	// Need to redial since invalid chain ID nodes are automatically disconnected
	state := n.createVerifiedConn(ctx, lggr)
	if state != NodeStateInvalidChainID {
		n.declareState(state)
		return
	}

	lggr.Debugw(fmt.Sprintf("Periodically re-checking RPC node %s with invalid chain ID", n.String()), "nodeState", n.getCachedState())

	chainIDRecheckBackoff := iutils.NewRedialBackoff()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(chainIDRecheckBackoff.Duration()):
			state := n.verifyConn(ctx, lggr)
			switch state {
			case NodeStateInvalidChainID:
				continue
			case NodeStateAlive:
				lggr.Infow(fmt.Sprintf("Successfully verified RPC node. Node was offline for %s", time.Since(invalidAt)), "nodeState", n.getCachedState())
				fallthrough
			default:
				n.declareState(state)
				return
			}
		}
	}
}

func (n *node[CHAIN_ID, HEAD, RPC]) syncingLoop() {
	defer n.wg.Done()
	ctx, cancel := n.newCtx()
	defer cancel()

	{
		// sanity check
		state := n.getCachedState()
		switch state {
		case NodeStateSyncing:
		case NodeStateClosed:
			return
		default:
			panic(fmt.Sprintf("syncingLoop can only run for node in NodeStateSyncing state, got: %s", state))
		}
	}

	syncingAt := time.Now()

	lggr := logger.Sugared(logger.Named(n.lfcLog, "Syncing"))
	lggr.Debugw(fmt.Sprintf("Periodically re-checking RPC node %s with syncing status", n.String()), "nodeState", n.getCachedState())
	// Need to redial since syncing nodes are automatically disconnected
	state := n.createVerifiedConn(ctx, lggr)
	if state != NodeStateSyncing {
		n.declareState(state)
		return
	}

	recheckBackoff := iutils.NewRedialBackoff()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(recheckBackoff.Duration()):
			lggr.Tracew("Trying to recheck if the node is still syncing", "nodeState", n.getCachedState())
			isSyncing, err := n.rpc.IsSyncing(ctx)
			if err != nil {
				lggr.Errorw("Unexpected error while verifying RPC node synchronization status", "err", err, "nodeState", n.getCachedState())
				n.declareUnreachable()
				return
			}

			if isSyncing {
				lggr.Errorw("Verification failed: Node is syncing", "nodeState", n.getCachedState())
				continue
			}

			lggr.Infow(fmt.Sprintf("Successfully verified RPC node. Node was syncing for %s", time.Since(syncingAt)), "nodeState", n.getCachedState())
			n.declareAlive()
			return
		}
	}
}
