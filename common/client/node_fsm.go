package client

import (
	"fmt"
	"math/big"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	promPoolRPCNodeTransitionsToAlive = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_alive",
		Help: transitionString(NodeStateAlive),
	}, []string{"chainID", "nodeName"})
	promPoolRPCNodeTransitionsToInSync = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_in_sync",
		Help: fmt.Sprintf("%s to %s", transitionString(NodeStateOutOfSync), NodeStateAlive),
	}, []string{"chainID", "nodeName"})
	promPoolRPCNodeTransitionsToOutOfSync = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_out_of_sync",
		Help: transitionString(NodeStateOutOfSync),
	}, []string{"chainID", "nodeName"})
	promPoolRPCNodeTransitionsToUnreachable = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_unreachable",
		Help: transitionString(NodeStateUnreachable),
	}, []string{"chainID", "nodeName"})
	promPoolRPCNodeTransitionsToInvalidChainID = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_invalid_chain_id",
		Help: transitionString(NodeStateInvalidChainID),
	}, []string{"chainID", "nodeName"})
	promPoolRPCNodeTransitionsToUnusable = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_unusable",
		Help: transitionString(NodeStateUnusable),
	}, []string{"chainID", "nodeName"})
	promPoolRPCNodeTransitionsToSyncing = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_num_transitions_to_syncing",
		Help: transitionString(NodeStateSyncing),
	}, []string{"chainID", "nodeName"})
)

// NodeState represents the current state of the node
// Node is a FSM (finite state machine)
type NodeState int

func (n NodeState) String() string {
	switch n {
	case NodeStateUndialed:
		return "Undialed"
	case NodeStateDialed:
		return "Dialed"
	case NodeStateInvalidChainID:
		return "InvalidChainID"
	case NodeStateAlive:
		return "Alive"
	case NodeStateUnreachable:
		return "Unreachable"
	case NodeStateUnusable:
		return "Unusable"
	case NodeStateOutOfSync:
		return "OutOfSync"
	case NodeStateClosed:
		return "Closed"
	case NodeStateSyncing:
		return "Syncing"
	case NodeStateFinalizedBlockOutOfSync:
		return "FinalizedBlockOutOfSync"
	default:
		return fmt.Sprintf("NodeState(%d)", n)
	}
}

// GoString prints a prettier state
func (n NodeState) GoString() string {
	return fmt.Sprintf("NodeState%s(%d)", n.String(), n)
}

const (
	// NodeStateUndialed is the first state of a virgin node
	NodeStateUndialed = NodeState(iota)
	// NodeStateDialed is after a node has successfully dialed but before it has verified the correct chain ID
	NodeStateDialed
	// NodeStateInvalidChainID is after chain ID verification failed
	NodeStateInvalidChainID
	// NodeStateAlive is a healthy node after chain ID verification succeeded
	NodeStateAlive
	// NodeStateUnreachable is a node that cannot be dialed or has disconnected
	NodeStateUnreachable
	// NodeStateOutOfSync is a node that is accepting connections but exceeded
	// the failure threshold without sending any new heads. It will be
	// disconnected, then put into a revive loop and re-awakened after redial
	// if a new head arrives
	NodeStateOutOfSync
	// NodeStateUnusable is a sendonly node that has an invalid URL that can never be reached
	NodeStateUnusable
	// NodeStateClosed is after the connection has been closed and the node is at the end of its lifecycle
	NodeStateClosed
	// NodeStateSyncing is a node that is actively back-filling blockchain. Usually, it's a newly set up node that is
	// still syncing the chain. The main difference from `NodeStateOutOfSync` is that it represents state relative
	// to other primary nodes configured in the MultiNode. In contrast, `NodeStateSyncing` represents the internal state of
	// the node (RPC).
	NodeStateSyncing
	// nodeStateFinalizedBlockOutOfSync - node is lagging behind on latest finalized block
	NodeStateFinalizedBlockOutOfSync
	// nodeStateLen tracks the number of states
	NodeStateLen
)

// allNodeStates represents all possible states a node can be in
var allNodeStates []NodeState

func init() {
	for s := NodeState(0); s < NodeStateLen; s++ {
		allNodeStates = append(allNodeStates, s)
	}
}

// FSM methods

// State allows reading the current state of the node.
func (n *node[CHAIN_ID, HEAD, RPC]) State() NodeState {
	n.stateMu.RLock()
	defer n.stateMu.RUnlock()
	return n.recalculateState()
}

func (n *node[CHAIN_ID, HEAD, RPC]) getCachedState() NodeState {
	n.stateMu.RLock()
	defer n.stateMu.RUnlock()
	return n.state
}

func (n *node[CHAIN_ID, HEAD, RPC]) recalculateState() NodeState {
	if n.state != NodeStateAlive {
		return n.state
	}

	// double check that node is not lagging on finalized block
	if n.nodePoolCfg.EnforceRepeatableRead() && n.isFinalizedBlockOutOfSync() {
		return NodeStateFinalizedBlockOutOfSync
	}

	return NodeStateAlive
}

func (n *node[CHAIN_ID, HEAD, RPC]) isFinalizedBlockOutOfSync() bool {
	if n.poolInfoProvider == nil {
		return false
	}

	highestObservedByCaller := n.poolInfoProvider.HighestUserObservations()
	latest, _ := n.rpc.GetInterceptedChainInfo()
	if n.chainCfg.FinalityTagEnabled() {
		return latest.FinalizedBlockNumber < highestObservedByCaller.FinalizedBlockNumber-int64(n.chainCfg.FinalizedBlockOffset())
	}

	return latest.BlockNumber < highestObservedByCaller.BlockNumber-int64(n.chainCfg.FinalizedBlockOffset())
}

// StateAndLatest returns nodeState with the latest ChainInfo observed by Node during current lifecycle.
func (n *node[CHAIN_ID, HEAD, RPC]) StateAndLatest() (NodeState, ChainInfo) {
	n.stateMu.RLock()
	defer n.stateMu.RUnlock()
	latest, _ := n.rpc.GetInterceptedChainInfo()
	return n.recalculateState(), latest
}

// HighestUserObservations - returns highest ChainInfo ever observed by external user of the Node
func (n *node[CHAIN_ID, HEAD, RPC]) HighestUserObservations() ChainInfo {
	_, highestUserObservations := n.rpc.GetInterceptedChainInfo()
	return highestUserObservations
}
func (n *node[CHAIN_ID, HEAD, RPC]) SetPoolChainInfoProvider(poolInfoProvider PoolChainInfoProvider) {
	n.poolInfoProvider = poolInfoProvider
}

// setState is only used by internal state management methods.
// This is low-level; care should be taken by the caller to ensure the new state is a valid transition.
// State changes should always be synchronous: only one goroutine at a time should change state.
// n.stateMu should not be locked for long periods of time because external clients expect a timely response from n.State()
func (n *node[CHAIN_ID, HEAD, RPC]) setState(s NodeState) {
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	n.state = s
}

// declareXXX methods change the state and pass conrol off the new state
// management goroutine

func (n *node[CHAIN_ID, HEAD, RPC]) declareAlive() {
	n.transitionToAlive(func() {
		n.lfcLog.Infow("RPC Node is online", "nodeState", n.state)
		n.wg.Add(1)
		go n.aliveLoop()
	})
}

func (n *node[CHAIN_ID, HEAD, RPC]) transitionToAlive(fn func()) {
	promPoolRPCNodeTransitionsToAlive.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateDialed, NodeStateInvalidChainID, NodeStateSyncing:
		n.state = NodeStateAlive
	default:
		panic(transitionFail(n.state, NodeStateAlive))
	}
	fn()
}

// declareInSync puts a node back into Alive state, allowing it to be used by
// pool consumers again
func (n *node[CHAIN_ID, HEAD, RPC]) declareInSync() {
	n.transitionToInSync(func() {
		n.lfcLog.Infow("RPC Node is back in sync", "nodeState", n.state)
		n.wg.Add(1)
		go n.aliveLoop()
	})
}

func (n *node[CHAIN_ID, HEAD, RPC]) transitionToInSync(fn func()) {
	promPoolRPCNodeTransitionsToAlive.WithLabelValues(n.chainID.String(), n.name).Inc()
	promPoolRPCNodeTransitionsToInSync.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateOutOfSync, NodeStateSyncing:
		n.state = NodeStateAlive
	default:
		panic(transitionFail(n.state, NodeStateAlive))
	}
	fn()
}

// declareOutOfSync puts a node into OutOfSync state, disconnecting all current
// clients and making it unavailable for use until back in-sync.
func (n *node[CHAIN_ID, HEAD, RPC]) declareOutOfSync(isOutOfSync func(num int64, td *big.Int) bool) {
	n.transitionToOutOfSync(func() {
		n.lfcLog.Errorw("RPC Node is out of sync", "nodeState", n.state)
		n.wg.Add(1)
		go n.outOfSyncLoop(isOutOfSync)
	})
}

func (n *node[CHAIN_ID, HEAD, RPC]) transitionToOutOfSync(fn func()) {
	promPoolRPCNodeTransitionsToOutOfSync.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateAlive:
		n.UnsubscribeAllExceptAliveLoop()
		n.state = NodeStateOutOfSync
	default:
		panic(transitionFail(n.state, NodeStateOutOfSync))
	}
	fn()
}

func (n *node[CHAIN_ID, HEAD, RPC]) declareUnreachable() {
	n.transitionToUnreachable(func() {
		n.lfcLog.Errorw("RPC Node is unreachable", "nodeState", n.state)
		n.wg.Add(1)
		go n.unreachableLoop()
	})
}

func (n *node[CHAIN_ID, HEAD, RPC]) transitionToUnreachable(fn func()) {
	promPoolRPCNodeTransitionsToUnreachable.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateUndialed, NodeStateDialed, NodeStateAlive, NodeStateOutOfSync, NodeStateInvalidChainID, NodeStateSyncing:
		n.UnsubscribeAllExceptAliveLoop()
		n.state = NodeStateUnreachable
	default:
		panic(transitionFail(n.state, NodeStateUnreachable))
	}
	fn()
}

func (n *node[CHAIN_ID, HEAD, RPC]) declareState(state NodeState) {
	if n.getCachedState() == NodeStateClosed {
		return
	}
	switch state {
	case NodeStateInvalidChainID:
		n.declareInvalidChainID()
	case NodeStateUnreachable:
		n.declareUnreachable()
	case NodeStateSyncing:
		n.declareSyncing()
	case NodeStateAlive:
		n.declareAlive()
	default:
		panic(fmt.Sprintf("%#v state declaration is not implemented", state))
	}
}

func (n *node[CHAIN_ID, HEAD, RPC]) declareInvalidChainID() {
	n.transitionToInvalidChainID(func() {
		n.lfcLog.Errorw("RPC Node has the wrong chain ID", "nodeState", n.state)
		n.wg.Add(1)
		go n.invalidChainIDLoop()
	})
}

func (n *node[CHAIN_ID, HEAD, RPC]) transitionToInvalidChainID(fn func()) {
	promPoolRPCNodeTransitionsToInvalidChainID.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateDialed, NodeStateOutOfSync, NodeStateSyncing:
		n.UnsubscribeAllExceptAliveLoop()
		n.state = NodeStateInvalidChainID
	default:
		panic(transitionFail(n.state, NodeStateInvalidChainID))
	}
	fn()
}

func (n *node[CHAIN_ID, HEAD, RPC]) declareSyncing() {
	n.transitionToSyncing(func() {
		n.lfcLog.Errorw("RPC Node is syncing", "nodeState", n.state)
		n.wg.Add(1)
		go n.syncingLoop()
	})
}

func (n *node[CHAIN_ID, HEAD, RPC]) transitionToSyncing(fn func()) {
	promPoolRPCNodeTransitionsToSyncing.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateDialed, NodeStateOutOfSync, NodeStateInvalidChainID:
		n.UnsubscribeAllExceptAliveLoop()
		n.state = NodeStateSyncing
	default:
		panic(transitionFail(n.state, NodeStateSyncing))
	}

	if !n.nodePoolCfg.NodeIsSyncingEnabled() {
		panic("unexpected transition to NodeStateSyncing, while it's disabled")
	}
	fn()
}

func transitionString(state NodeState) string {
	return fmt.Sprintf("Total number of times node has transitioned to %s", state)
}

func transitionFail(from NodeState, to NodeState) string {
	return fmt.Sprintf("cannot transition from %#v to %#v", from, to)
}
