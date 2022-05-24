package client

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	promEVMPoolRPCNodeHighestSeenBlock = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "evm_pool_rpc_node_highest_seen_block",
		Help: "The highest seen block for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeNumSeenBlocks = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_num_seen_blocks",
		Help: "The total number of new blocks seen by the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodePolls = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_polls_total",
		Help: "The total number of poll checks for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodePollsFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_polls_failed",
		Help: "The total number of failed poll checks for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodePollsSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_polls_success",
		Help: "The total number of successful poll checks for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
)

// zombieNodeCheckInterval controls how often to re-check to see if we need to
// state change in case we have to force a state transition due to no available
// nodes.
// NOTE: This only applies to out-of-sync nodes if they are the last available node
func zombieNodeCheckInterval(cfg NodeConfig) time.Duration {
	interval := cfg.NodeNoNewHeadsThreshold()
	if interval <= 0 || interval > queryTimeout {
		interval = queryTimeout
	}
	return utils.WithJitter(interval)
}

// Node is a FSM
// Each state has a loop that goes with it, which monitors the node and moves it into another state as necessary.
// Only one loop must run at a time.
// Each loop passes control onto the next loop as it exits, except when the node is Closed which terminates the loop permanently.

// This handles node lifecycle for the ALIVE state
// Should only be run ONCE per node, after a successful Dial
func (n *node) aliveLoop() {
	defer n.wg.Done()

	{
		// sanity check
		state := n.State()
		switch state {
		case NodeStateAlive:
		case NodeStateClosed:
			return
		default:
			panic(fmt.Sprintf("aliveLoop can only run for node in Alive state, got: %s", state))
		}
	}

	noNewHeadsTimeoutThreshold := n.cfg.NodeNoNewHeadsThreshold()
	pollFailureThreshold := n.cfg.NodePollFailureThreshold()
	pollInterval := n.cfg.NodePollInterval()

	lggr := n.lfcLog.Named("Alive").With("noNewHeadsTimeoutThreshold", noNewHeadsTimeoutThreshold, "pollInterval", pollInterval, "pollFailureThreshold", pollFailureThreshold)
	lggr.Tracew("Alive loop starting", "nodeState", n.State())

	var headsC <-chan *evmtypes.Head
	var outOfSyncT *time.Ticker
	var outOfSyncTC <-chan time.Time
	var sub ethereum.Subscription
	var subErrC <-chan error
	if noNewHeadsTimeoutThreshold > 0 {
		lggr.Debugw("Head liveness checking enabled", "nodeState", n.State())
		var err error
		writableCh := make(chan *evmtypes.Head)
		sub, err = n.EthSubscribe(context.Background(), writableCh, "newHeads")
		headsC = writableCh
		if err != nil {
			lggr.Errorw("Initial subscribe for liveness checking failed", "nodeState", n.State())
			n.declareUnreachable()
			return
		}
		defer sub.Unsubscribe()
		outOfSyncT = time.NewTicker(noNewHeadsTimeoutThreshold)
		defer outOfSyncT.Stop()
		outOfSyncTC = outOfSyncT.C
		subErrC = sub.Err()
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

	var latestReceivedBlockNumber int64 = -1
	var pollFailures uint32

	for {
		select {
		case <-n.chStop:
			return
		case <-pollCh:
			var version string
			promEVMPoolRPCNodePolls.WithLabelValues(n.chainID.String(), n.name).Inc()
			lggr.Tracew("Polling for version", "nodeState", n.State(), "pollFailures", pollFailures)
			ctx, cancel := context.WithTimeout(context.Background(), pollInterval)
			ctx, cancel2 := n.makeQueryCtx(ctx)
			err := n.CallContext(ctx, &version, "web3_clientVersion")
			cancel2()
			cancel()
			if err != nil {
				// prevent overflow
				if pollFailures < math.MaxUint32 {
					promEVMPoolRPCNodePollsFailed.WithLabelValues(n.chainID.String(), n.name).Inc()
					pollFailures++
				}
				lggr.Warnw(fmt.Sprintf("Poll failure, RPC endpoint %s failed to respond properly", n.String()), "err", err, "pollFailures", pollFailures, "nodeState", n.State())
			} else {
				lggr.Tracew("Version poll successful", "nodeState", n.State(), "clientVersion", version)
				promEVMPoolRPCNodePollsSuccess.WithLabelValues(n.chainID.String(), n.name).Inc()
				pollFailures = 0
			}
			if pollFailureThreshold > 0 && pollFailures >= pollFailureThreshold {
				lggr.Errorw(fmt.Sprintf("RPC endpoint failed to respond to %d consecutive polls", pollFailures), "pollFailures", pollFailures, "nodeState", n.State())
				if n.nLiveNodes != nil && n.nLiveNodes() < 2 {
					lggr.Critical("RPC endpoint failed to respond to polls; but cannot disable this connection because there are no other RPC endpoints, or all other RPC endpoints are dead. Chainlink is now operating in a degraded state and urgent action is required to resolve the issue")
					continue
				}
				n.declareUnreachable()
				return
			}
		case bh, open := <-headsC:
			if !open {
				lggr.Errorw("Subscription channel unexpectedly closed", "nodeState", n.State())
				n.declareUnreachable()
				return
			}
			promEVMPoolRPCNodeNumSeenBlocks.WithLabelValues(n.chainID.String(), n.name).Inc()
			lggr.Tracew("Got head", "head", bh)
			if bh.Number > latestReceivedBlockNumber {
				promEVMPoolRPCNodeHighestSeenBlock.WithLabelValues(n.chainID.String(), n.name).Set(float64(bh.Number))
				lggr.Tracew("Got higher block number, resetting timer", "latestReceivedBlockNumber", latestReceivedBlockNumber, "blockNumber", bh.Number, "nodeState", n.State())
				latestReceivedBlockNumber = bh.Number
			} else {
				lggr.Tracew("Ignoring previously seen block number", "latestReceivedBlockNumber", latestReceivedBlockNumber, "blockNumber", bh.Number, "nodeState", n.State())
			}
			outOfSyncT.Reset(noNewHeadsTimeoutThreshold)
		case err := <-subErrC:
			lggr.Errorw("Subscription was terminated", "err", err, "nodeState", n.State())
			n.declareUnreachable()
			return
		case <-outOfSyncTC:
			// We haven't received a head on the channel for at least the
			// threshold amount of time, mark it broken
			lggr.Errorw(fmt.Sprintf("RPC endpoint detected out of sync; no new heads received for %s (last head received was %v)", noNewHeadsTimeoutThreshold, latestReceivedBlockNumber), "nodeState", n.State(), "latestReceivedBlockNumber", latestReceivedBlockNumber, "noNewHeadsTimeoutThreshold", noNewHeadsTimeoutThreshold)
			if n.nLiveNodes != nil && n.nLiveNodes() < 2 {
				lggr.Critical("RPC endpoint detected out of sync; but cannot disable this connection because there are no other RPC endpoints, or all other RPC endpoints dead. Chainlink is now operating in a degraded state and urgent action is required to resolve the issue")
				// We don't necessarily want to wait the full timeout to check again, we should
				// check regularly and log noisily in this state
				outOfSyncT.Reset(zombieNodeCheckInterval(n.cfg))
				continue
			}
			n.declareOutOfSync(latestReceivedBlockNumber)
			return
		}
	}
}

// outOfSyncLoop takes an OutOfSync node and puts it back to live status if it
// receives a later head than one we have already seen
func (n *node) outOfSyncLoop(stuckAtBlockNumber int64) {
	defer n.wg.Done()

	{
		// sanity check
		state := n.State()
		switch state {
		case NodeStateOutOfSync:
		case NodeStateClosed:
			return
		default:
			panic(fmt.Sprintf("outOfSyncLoop can only run for node in OutOfSync state, got: %s", state))
		}
	}

	outOfSyncAt := time.Now()

	lggr := n.lfcLog.Named("OutOfSync")
	lggr.Debugw("Trying to revive out-of-sync RPC node", "nodeState", n.State())

	// Need to redial since out-of-sync nodes are automatically disconnected
	err := n.dial(context.Background())
	if err != nil {
		lggr.Errorw("Failed to dial out-of-sync RPC node", "nodeState", n.State())
		n.declareUnreachable()
		return
	}

	// Manually re-verify since out-of-sync nodes are automatically disconnected
	err = n.verify(context.Background())
	if err != nil {
		lggr.Errorw(fmt.Sprintf("Failed to verify out-of-sync RPC node: %v", err), "err", err)
		n.declareInvalidChainID()
		return
	}

	lggr.Tracew("Successfully subscribed to heads feed on out-of-sync RPC node", "stuckAtBlockNumber", stuckAtBlockNumber, "nodeState", n.State())

	ch := make(chan *evmtypes.Head)
	subCtx, cancel := n.makeQueryCtx(context.Background())
	// raw call here to bypass node state checking
	sub, err := n.ws.rpc.EthSubscribe(subCtx, ch, "newHeads")
	cancel()
	if err != nil {
		lggr.Errorw("Failed to subscribe heads on out-of-sync RPC node", "nodeState", n.State(), "err", err)
		n.declareUnreachable()
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-n.chStop:
			return
		case head, open := <-ch:
			if !open {
				lggr.Error("Subscription channel unexpectedly closed", "nodeState", n.State())
				n.declareUnreachable()
				return
			}
			if head.Number > stuckAtBlockNumber {
				// unstuck! flip back into alive loop
				lggr.Infow(fmt.Sprintf("Received new block for RPC node %s. Node was offline for %s", n.String(), time.Since(outOfSyncAt)), "latestReceivedBlockNumber", head.Number, "nodeState", n.State())
				n.declareInSync()
				return
			}
			lggr.Debugw("Received previously seen block for RPC node, waiting for new block before marking as live again", "stuckAtBlockNumber", stuckAtBlockNumber, "blockNumber", head.Number, "nodeState", n.State())
		case <-time.After(zombieNodeCheckInterval(n.cfg)):
			if n.nLiveNodes != nil && n.nLiveNodes() < 1 {
				lggr.Critical("RPC endpoint is still out of sync, but there are no other available nodes. This RPC node will be forcibly moved back into the live pool in a degraded state")
				n.declareInSync()
				return

			}
		case err := <-sub.Err():
			lggr.Errorw("Subscription was terminated", "nodeState", n.State(), "err", err)
			n.declareUnreachable()
			return
		}
	}
}

func (n *node) unreachableLoop() {
	defer n.wg.Done()

	{
		// sanity check
		state := n.State()
		switch state {
		case NodeStateUnreachable:
		case NodeStateClosed:
			return
		default:
			panic(fmt.Sprintf("unreachableLoop can only run for node in Unreachable state, got: %s", state))
		}
	}

	unreachableAt := time.Now()

	lggr := n.lfcLog.Named("Unreachable")
	lggr.Debugw("Trying to revive unreachable RPC node", "nodeState", n.State())

	dialRetryBackoff := utils.NewRedialBackoff()

	for {
		select {
		case <-n.chStop:
			return
		case <-time.After(dialRetryBackoff.Duration()):
			lggr.Tracew("Trying to re-dial RPC node", "nodeState", n.State())

			err := n.dial(context.Background())
			if err != nil {
				lggr.Errorw(fmt.Sprintf("Failed to redial RPC node; still unreachable: %v", err), "err", err, "nodeState", n.State())
				continue
			}

			n.setState(NodeStateDialed)

			err = n.verify(context.Background())
			if errors.Is(err, errInvalidChainID) {
				lggr.Errorw("Failed to redial RPC node; remote endpoint returned the wrong chain ID", "err", err)
				n.declareInvalidChainID()
				return
			} else if err != nil {
				lggr.Errorw(fmt.Sprintf("Failed to redial RPC node; verify failed: %v", err), "err", err)
				n.declareUnreachable()
				return
			}

			lggr.Infow(fmt.Sprintf("Successfully redialled and verified RPC node %s. Node was offline for %s", n.String(), time.Since(unreachableAt)), "nodeState", n.State())
			n.declareAlive()
			return
		}
	}
}

func (n *node) invalidChainIDLoop() {
	defer n.wg.Done()

	{
		// sanity check
		state := n.State()
		switch state {
		case NodeStateInvalidChainID:
		case NodeStateClosed:
			return
		default:
			panic(fmt.Sprintf("invalidChainIDLoop can only run for node in InvalidChainID state, got: %s", state))
		}
	}

	invalidAt := time.Now()

	lggr := n.lfcLog.Named("InvalidChainID")
	lggr.Debugw(fmt.Sprintf("Periodically re-checking RPC node %s with invalid chain ID", n.String()), "nodeState", n.State())

	chainIDRecheckBackoff := utils.NewRedialBackoff()

	for {
		select {
		case <-n.chStop:
			return
		case <-time.After(chainIDRecheckBackoff.Duration()):
			err := n.verify(context.Background())
			if errors.Is(err, errInvalidChainID) {
				lggr.Errorw("Failed to verify RPC node; remote endpoint returned the wrong chain ID", "err", err)
				continue
			} else if err != nil {
				lggr.Errorw(fmt.Sprintf("Unexpected error while verifying RPC node chain ID; %v", err), "err", err)
				n.declareUnreachable()
				return
			}
			lggr.Infow(fmt.Sprintf("Successfully verified RPC node. Node was offline for %s", time.Since(invalidAt)), "nodeState", n.State())
			n.declareAlive()
			return
		}
	}
}
