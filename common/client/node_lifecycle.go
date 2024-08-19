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
		case nodeStateAlive:
		case nodeStateClosed:
			return
		default:
			panic(fmt.Sprintf("aliveLoop can only run for node in Alive state, got: %s", state))
		}
	}

	noNewHeadsTimeoutThreshold := n.chainCfg.NodeNoNewHeadsThreshold()
	noNewFinalizedBlocksTimeoutThreshold := n.chainCfg.NoNewFinalizedHeadsThreshold()
	pollFailureThreshold := n.nodePoolCfg.PollFailureThreshold()
	pollInterval := n.nodePoolCfg.PollInterval()

	lggr := logger.Sugared(n.lfcLog).Named("Alive").With("noNewHeadsTimeoutThreshold", noNewHeadsTimeoutThreshold, "pollInterval", pollInterval, "pollFailureThreshold", pollFailureThreshold)
	lggr.Tracew("Alive loop starting", "nodeState", n.getCachedState())

	headsSub, err := n.registerNewSubscription(ctx, lggr.With("subscriptionType", "heads"),
		n.chainCfg.NodeNoNewHeadsThreshold(), n.rpc.SubscribeToHeads)
	if err != nil {
		lggr.Errorw("Initial subscribe for heads failed", "nodeState", n.getCachedState(), "err", err)
		n.declareUnreachable()
		return
	}

	// TODO: will be removed as part of merging effort with BCI-2875
	n.rpc.SetAliveLoopSub(headsSub.sub)

	defer headsSub.Unsubscribe()

	var finalizedHeadsSub headSubscription[HEAD]
	if n.chainCfg.FinalityTagEnabled() {
		finalizedHeadsSub, err = n.registerNewSubscription(ctx, lggr.With("subscriptionType", "finalizedHeads"),
			n.chainCfg.NoNewFinalizedHeadsThreshold(), n.rpc.SubscribeToFinalizedHeads)
		if err != nil {
			lggr.Errorw("Failed to subscribe to finalized heads", "err", err)
			n.declareUnreachable()
			return
		}

		defer finalizedHeadsSub.Unsubscribe()
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

	localHighestChainInfo, _ := n.rpc.GetInterceptedChainInfo()
	var pollFailures uint32

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollCh:
			promPoolRPCNodePolls.WithLabelValues(n.chainID.String(), n.name).Inc()
			lggr.Tracew("Polling for version", "nodeState", n.getCachedState(), "pollFailures", pollFailures)
			var version string
			version, err = func(ctx context.Context) (string, error) {
				ctx, cancel := context.WithTimeout(ctx, pollInterval)
				defer cancel()
				return n.RPC().ClientVersion(ctx)
			}(ctx)
			if err != nil {
				// prevent overflow
				if pollFailures < math.MaxUint32 {
					promPoolRPCNodePollsFailed.WithLabelValues(n.chainID.String(), n.name).Inc()
					pollFailures++
				}
				lggr.Warnw(fmt.Sprintf("Poll failure, RPC endpoint %s failed to respond properly", n.String()), "err", err, "pollFailures", pollFailures, "nodeState", n.getCachedState())
			} else {
				lggr.Debugw("Version poll successful", "nodeState", n.getCachedState(), "clientVersion", version)
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
			_, latestChainInfo := n.StateAndLatest()
			if outOfSync, liveNodes := n.isOutOfSyncWithPool(latestChainInfo); outOfSync {
				// note: there must be another live node for us to be out of sync
				lggr.Errorw("RPC endpoint has fallen behind", "blockNumber", latestChainInfo.BlockNumber, "totalDifficulty", latestChainInfo.TotalDifficulty, "nodeState", n.getCachedState())
				if liveNodes < 2 {
					lggr.Criticalf("RPC endpoint has fallen behind; %s %s", msgCannotDisable, msgDegradedState)
					continue
				}
				n.declareOutOfSync(syncStatusNotInSyncWithPool)
				return
			}
		case bh, open := <-headsSub.Heads:
			if !open {
				lggr.Errorw("Subscription channel unexpectedly closed", "nodeState", n.getCachedState())
				n.declareUnreachable()
				return
			}

			receivedNewHead := n.onNewHead(lggr, &localHighestChainInfo, bh)
			if receivedNewHead && noNewHeadsTimeoutThreshold > 0 {
				headsSub.ResetTimer(noNewHeadsTimeoutThreshold)
			}
		case err = <-headsSub.Errors:
			lggr.Errorw("Subscription was terminated", "err", err, "nodeState", n.getCachedState())
			n.declareUnreachable()
			return
		case <-headsSub.NoNewHeads:
			// We haven't received a head on the channel for at least the
			// threshold amount of time, mark it broken
			lggr.Errorw(fmt.Sprintf("RPC endpoint detected out of sync; no new heads received for %s (last head received was %v)", noNewHeadsTimeoutThreshold, localHighestChainInfo.BlockNumber), "nodeState", n.getCachedState(), "latestReceivedBlockNumber", localHighestChainInfo.BlockNumber, "noNewHeadsTimeoutThreshold", noNewHeadsTimeoutThreshold)
			if n.poolInfoProvider != nil {
				if l, _ := n.poolInfoProvider.LatestChainInfo(); l < 2 {
					lggr.Criticalf("RPC endpoint detected out of sync; %s %s", msgCannotDisable, msgDegradedState)
					// We don't necessarily want to wait the full timeout to check again, we should
					// check regularly and log noisily in this state
					headsSub.ResetTimer(zombieNodeCheckInterval(noNewHeadsTimeoutThreshold))
					continue
				}
			}
			n.declareOutOfSync(syncStatusNoNewHead)
			return
		case latestFinalized, open := <-finalizedHeadsSub.Heads:
			if !open {
				lggr.Errorw("Finalized heads subscription channel unexpectedly closed", "nodeState", n.getCachedState())
				n.declareUnreachable()
				return
			}

			receivedNewHead := n.onNewFinalizedHead(lggr, &localHighestChainInfo, latestFinalized)
			if receivedNewHead && noNewFinalizedBlocksTimeoutThreshold > 0 {
				finalizedHeadsSub.ResetTimer(noNewFinalizedBlocksTimeoutThreshold)
			}
		case <-finalizedHeadsSub.NoNewHeads:
			// We haven't received a finalized head on the channel for at least the
			// threshold amount of time, mark it broken
			lggr.Errorw(fmt.Sprintf("RPC's finalized state is out of sync; no new finalized heads received for %s (last finalized head received was %v)", noNewFinalizedBlocksTimeoutThreshold, localHighestChainInfo.FinalizedBlockNumber), "latestReceivedBlockNumber", localHighestChainInfo.BlockNumber)
			if n.poolInfoProvider != nil {
				if l, _ := n.poolInfoProvider.LatestChainInfo(); l < 2 {
					lggr.Criticalf("RPC's finalized state is out of sync; %s %s", msgCannotDisable, msgDegradedState)
					// We don't necessarily want to wait the full timeout to check again, we should
					// check regularly and log noisily in this state
					finalizedHeadsSub.ResetTimer(zombieNodeCheckInterval(noNewFinalizedBlocksTimeoutThreshold))
					continue
				}
			}
			n.declareOutOfSync(syncStatusNoNewFinalizedHead)
			return
		case <-finalizedHeadsSub.Errors:
			lggr.Errorw("Finalized heads subscription was terminated", "err", err)
			n.declareUnreachable()
			return
		}
	}
}

type headSubscription[HEAD any] struct {
	Heads      <-chan HEAD
	Errors     <-chan error
	NoNewHeads <-chan time.Time

	noNewHeadsTicker *time.Ticker
	sub              types.Subscription
	cleanUpTasks     []func()
}

func (sub *headSubscription[HEAD]) ResetTimer(duration time.Duration) {
	sub.noNewHeadsTicker.Reset(duration)
}

func (sub *headSubscription[HEAD]) Unsubscribe() {
	for _, doCleanUp := range sub.cleanUpTasks {
		doCleanUp()
	}
}

func (n *node[CHAIN_ID, HEAD, PRC]) registerNewSubscription(ctx context.Context, lggr logger.SugaredLogger,
	noNewDataThreshold time.Duration, newSub func(ctx context.Context) (<-chan HEAD, types.Subscription, error)) (headSubscription[HEAD], error) {
	result := headSubscription[HEAD]{}
	var err error
	var sub types.Subscription
	result.Heads, sub, err = newSub(ctx)
	if err != nil {
		return result, err
	}

	result.Errors = sub.Err()
	lggr.Debug("Successfully subscribed")

	// TODO: will be removed as part of merging effort with BCI-2875
	result.sub = sub
	//n.stateMu.Lock()
	//n.healthCheckSubs = append(n.healthCheckSubs, sub)
	//n.stateMu.Unlock()

	result.cleanUpTasks = append(result.cleanUpTasks, sub.Unsubscribe)

	if noNewDataThreshold > 0 {
		lggr.Debugw("Subscription liveness checking enabled")
		result.noNewHeadsTicker = time.NewTicker(noNewDataThreshold)
		result.NoNewHeads = result.noNewHeadsTicker.C
		result.cleanUpTasks = append(result.cleanUpTasks, result.noNewHeadsTicker.Stop)
	} else {
		lggr.Debug("Subscription liveness checking disabled")
	}

	return result, nil
}

func (n *node[CHAIN_ID, HEAD, RPC]) onNewFinalizedHead(lggr logger.SugaredLogger, chainInfo *ChainInfo, latestFinalized HEAD) bool {
	if !latestFinalized.IsValid() {
		lggr.Warn("Latest finalized block is not valid")
		return false
	}

	latestFinalizedBN := latestFinalized.BlockNumber()
	lggr.Tracew("Got latest finalized head", "latestFinalized", latestFinalized)
	if latestFinalizedBN <= chainInfo.FinalizedBlockNumber {
		lggr.Tracew("Ignoring previously seen finalized block number")
		return false
	}

	promPoolRPCNodeHighestFinalizedBlock.WithLabelValues(n.chainID.String(), n.name).Set(float64(latestFinalizedBN))
	chainInfo.FinalizedBlockNumber = latestFinalizedBN
	return true
}

func (n *node[CHAIN_ID, HEAD, RPC]) onNewHead(lggr logger.SugaredLogger, chainInfo *ChainInfo, head HEAD) bool {
	if !head.IsValid() {
		lggr.Warn("Latest head is not valid")
		return false
	}

	promPoolRPCNodeNumSeenBlocks.WithLabelValues(n.chainID.String(), n.name).Inc()
	lggr.Tracew("Got head", "head", head)
	lggr = lggr.With("latestReceivedBlockNumber", chainInfo.BlockNumber, "blockNumber", head.BlockNumber(), "nodeState", n.getCachedState())
	if head.BlockNumber() <= chainInfo.BlockNumber {
		lggr.Tracew("Ignoring previously seen block number")
		return false
	}

	promPoolRPCNodeHighestSeenBlock.WithLabelValues(n.chainID.String(), n.name).Set(float64(head.BlockNumber()))
	chainInfo.BlockNumber = head.BlockNumber()

	if !n.chainCfg.FinalityTagEnabled() {
		latestFinalizedBN := max(head.BlockNumber()-int64(n.chainCfg.FinalityDepth()), 0)
		if latestFinalizedBN > chainInfo.FinalizedBlockNumber {
			promPoolRPCNodeHighestFinalizedBlock.WithLabelValues(n.chainID.String(), n.name).Set(float64(latestFinalizedBN))
			chainInfo.FinalizedBlockNumber = latestFinalizedBN
		}
	}

	return true
}

// isOutOfSyncWithPool returns outOfSync true if num or td is more than SyncThresold behind the best node.
// Always returns outOfSync false for SyncThreshold 0.
// liveNodes is only included when outOfSync is true.
func (n *node[CHAIN_ID, HEAD, RPC]) isOutOfSyncWithPool(localState ChainInfo) (outOfSync bool, liveNodes int) {
	if n.poolInfoProvider == nil {
		n.lfcLog.Warn("skipping sync state against the pool - should only occur in tests")
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
		return localState.BlockNumber < ci.BlockNumber-int64(threshold), ln
	case NodeSelectionModeTotalDifficulty:
		bigThreshold := big.NewInt(int64(threshold))
		return localState.TotalDifficulty.Cmp(bigmath.Sub(ci.TotalDifficulty, bigThreshold)) < 0, ln
	default:
		panic("unrecognized NodeSelectionMode: " + mode)
	}
}

const (
	msgReceivedBlock          = "Received block for RPC node, waiting until back in-sync to mark as live again"
	msgReceivedFinalizedBlock = "Received new finalized block for RPC node, waiting until back in-sync to mark as live again"
	msgInSync                 = "RPC node back in sync"
)

// outOfSyncLoop takes an OutOfSync node and waits until isOutOfSync returns false to go back to live status
func (n *node[CHAIN_ID, HEAD, RPC]) outOfSyncLoop(syncIssues syncStatus) {
	defer n.wg.Done()
	ctx, cancel := n.newCtx()
	defer cancel()

	{
		// sanity check
		state := n.getCachedState()
		switch state {
		case nodeStateOutOfSync:
		case nodeStateClosed:
			return
		default:
			panic(fmt.Sprintf("outOfSyncLoop can only run for node in OutOfSync state, got: %s", state))
		}
	}

	outOfSyncAt := time.Now()

	// set logger name to OutOfSync or FinalizedBlockOutOfSync
	lggr := logger.Sugared(logger.Named(n.lfcLog, n.getCachedState().String())).With("nodeState", n.getCachedState())
	lggr.Debugw("Trying to revive out-of-sync RPC node")

	// Need to redial since out-of-sync nodes are automatically disconnected
	state := n.createVerifiedConn(ctx, lggr)
	if state != nodeStateAlive {
		n.declareState(state)
		return
	}

	noNewHeadsTimeoutThreshold := n.chainCfg.NodeNoNewHeadsThreshold()
	headsSub, err := n.registerNewSubscription(ctx, lggr.With("subscriptionType", "heads"),
		noNewHeadsTimeoutThreshold, n.rpc.SubscribeToHeads)
	if err != nil {
		lggr.Errorw("Failed to subscribe heads on out-of-sync RPC node", "err", err)
		n.declareUnreachable()
		return
	}

	lggr.Tracew("Successfully subscribed to heads feed on out-of-sync RPC node")
	defer headsSub.Unsubscribe()

	noNewFinalizedBlocksTimeoutThreshold := n.chainCfg.NoNewFinalizedHeadsThreshold()
	var finalizedHeadsSub headSubscription[HEAD]
	if n.chainCfg.FinalityTagEnabled() {
		finalizedHeadsSub, err = n.registerNewSubscription(ctx, lggr.With("subscriptionType", "finalizedHeads"),
			noNewFinalizedBlocksTimeoutThreshold, n.rpc.SubscribeToFinalizedHeads)
		if err != nil {
			lggr.Errorw("Subscribe to finalized heads failed on out-of-sync RPC node", "err", err)
			n.declareUnreachable()
			return
		}

		lggr.Tracew("Successfully subscribed to finalized heads feed on out-of-sync RPC node")
		defer finalizedHeadsSub.Unsubscribe()
	}

	_, localHighestChainInfo := n.rpc.GetInterceptedChainInfo()
	for {
		if syncIssues == syncStatusSynced {
			// back in-sync! flip back into alive loop
			lggr.Infow(fmt.Sprintf("%s: %s. Node was out-of-sync for %s", msgInSync, n.String(), time.Since(outOfSyncAt)))
			n.declareInSync()
			return
		}

		select {
		case <-ctx.Done():
			return
		case head, open := <-headsSub.Heads:
			if !open {
				lggr.Errorw("Subscription channel unexpectedly closed", "nodeState", n.getCachedState())
				n.declareUnreachable()
				return
			}

			if !n.onNewHead(lggr, &localHighestChainInfo, head) {
				continue
			}

			// received a new head - clear NoNewHead flag
			syncIssues &= ^syncStatusNoNewHead
			if outOfSync, _ := n.isOutOfSyncWithPool(localHighestChainInfo); !outOfSync {
				// we caught up with the pool - clear NotInSyncWithPool flag
				syncIssues &= ^syncStatusNotInSyncWithPool
			} else {
				// we've received new head, but lagging behind the pool, add NotInSyncWithPool flag to prevent false transition to alive
				syncIssues |= syncStatusNotInSyncWithPool
			}

			if noNewHeadsTimeoutThreshold > 0 {
				headsSub.ResetTimer(noNewHeadsTimeoutThreshold)
			}

			lggr.Debugw(msgReceivedBlock, "blockNumber", head.BlockNumber(), "blockDifficulty", head.BlockDifficulty(), "syncIssues", syncIssues)
		case <-time.After(zombieNodeCheckInterval(noNewHeadsTimeoutThreshold)):
			if n.poolInfoProvider != nil {
				if l, _ := n.poolInfoProvider.LatestChainInfo(); l < 1 {
					lggr.Criticalw("RPC endpoint is still out of sync, but there are no other available nodes. This RPC node will be forcibly moved back into the live pool in a degraded state", "syncIssues", syncIssues)
					n.declareInSync()
					return
				}
			}
		case err := <-headsSub.Errors:
			lggr.Errorw("Subscription was terminated", "err", err)
			n.declareUnreachable()
			return
		case <-headsSub.NoNewHeads:
			// we are not resetting the timer, as there is no need to add syncStatusNoNewHead until it's removed on new head.
			syncIssues |= syncStatusNoNewHead
			lggr.Debugw(fmt.Sprintf("No new heads received for %s. Node stays out-of-sync due to sync issues: %s", noNewHeadsTimeoutThreshold, syncIssues))
		case latestFinalized, open := <-finalizedHeadsSub.Heads:
			if !open {
				lggr.Errorw("Finalized heads subscription channel unexpectedly closed")
				n.declareUnreachable()
				return
			}
			if !latestFinalized.IsValid() {
				lggr.Warn("Latest finalized block is not valid")
				continue
			}

			receivedNewHead := n.onNewFinalizedHead(lggr, &localHighestChainInfo, latestFinalized)
			if !receivedNewHead {
				continue
			}

			// on new finalized head remove NoNewFinalizedHead flag from the mask
			syncIssues &= ^syncStatusNoNewFinalizedHead
			if noNewFinalizedBlocksTimeoutThreshold > 0 {
				finalizedHeadsSub.ResetTimer(noNewFinalizedBlocksTimeoutThreshold)
			}

			lggr.Debugw(msgReceivedFinalizedBlock, "blockNumber", latestFinalized.BlockNumber(), "syncIssues", syncIssues)
		case err := <-finalizedHeadsSub.Errors:
			lggr.Errorw("Finalized head subscription was terminated", "err", err)
			n.declareUnreachable()
			return
		case <-finalizedHeadsSub.NoNewHeads:
			// we are not resetting the timer, as there is no need to add syncStatusNoNewFinalizedHead until it's removed on new finalized head.
			syncIssues |= syncStatusNoNewFinalizedHead
			lggr.Debugw(fmt.Sprintf("No new finalized heads received for %s. Node stays out-of-sync due to sync issues: %s", noNewFinalizedBlocksTimeoutThreshold, syncIssues))
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
		case nodeStateUnreachable:
		case nodeStateClosed:
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

			n.setState(nodeStateDialed)

			state := n.verifyConn(ctx, lggr)
			switch state {
			case nodeStateUnreachable:
				n.setState(nodeStateUnreachable)
				continue
			case nodeStateAlive:
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
		case nodeStateInvalidChainID:
		case nodeStateClosed:
			return
		default:
			panic(fmt.Sprintf("invalidChainIDLoop can only run for node in InvalidChainID state, got: %s", state))
		}
	}

	invalidAt := time.Now()

	lggr := logger.Named(n.lfcLog, "InvalidChainID")

	// Need to redial since invalid chain ID nodes are automatically disconnected
	state := n.createVerifiedConn(ctx, lggr)
	if state != nodeStateInvalidChainID {
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
			case nodeStateInvalidChainID:
				continue
			case nodeStateAlive:
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
		case nodeStateSyncing:
		case nodeStateClosed:
			return
		default:
			panic(fmt.Sprintf("syncingLoop can only run for node in nodeStateSyncing state, got: %s", state))
		}
	}

	syncingAt := time.Now()

	lggr := logger.Sugared(logger.Named(n.lfcLog, "Syncing"))
	lggr.Debugw(fmt.Sprintf("Periodically re-checking RPC node %s with syncing status", n.String()), "nodeState", n.getCachedState())
	// Need to redial since syncing nodes are automatically disconnected
	state := n.createVerifiedConn(ctx, lggr)
	if state != nodeStateSyncing {
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
