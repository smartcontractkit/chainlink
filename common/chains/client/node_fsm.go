package client

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
	// nodeStateLen tracks the number of states
	nodeStateLen
)

// allNodeStates represents all possible states a node can be in
var allNodeStates []NodeState

func init() {
	for s := NodeState(0); s < nodeStateLen; s++ {
		allNodeStates = append(allNodeStates, s)
	}
}

// FSM methods

// State allows reading the current state of the node.
func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) State() NodeState {
	n.stateMu.RLock()
	defer n.stateMu.RUnlock()
	return n.state
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) StateAndLatest() (NodeState, int64, *utils.Big) {
	n.stateMu.RLock()
	defer n.stateMu.RUnlock()
	return n.state, n.stateLatestBlockNumber, n.stateLatestTotalDifficulty
}

// setState is only used by internal state management methods.
// This is low-level; care should be taken by the caller to ensure the new state is a valid transition.
// State changes should always be synchronous: only one goroutine at a time should change state.
// n.stateMu should not be locked for long periods of time because external clients expect a timely response from n.State()
func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) setState(s NodeState) {
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	n.SetState(s)
}

// declareXXX methods change the state and pass conrol off the new state
// management goroutine

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) declareAlive() {
	n.transitionToAlive(func() {
		n.lfcLog.Infow("RPC Node is online", "nodeState", n.state)
		n.wg.Add(1)
		go n.aliveLoop()
	})
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) transitionToAlive(fn func()) {
	promPoolRPCNodeTransitionsToAlive.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateDialed, NodeStateInvalidChainID:
		n.SetState(NodeStateAlive)
	default:
		panic(transitionFail(n.state, NodeStateAlive))
	}
	fn()
}

// declareInSync puts a node back into Alive state, allowing it to be used by
// pool consumers again
func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) declareInSync() {
	n.transitionToInSync(func() {
		n.lfcLog.Infow("RPC Node is back in sync", "nodeState", n.state)
		n.wg.Add(1)
		go n.aliveLoop()
	})
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) transitionToInSync(fn func()) {
	promPoolRPCNodeTransitionsToAlive.WithLabelValues(n.chainID.String(), n.name).Inc()
	promPoolRPCNodeTransitionsToInSync.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateOutOfSync:
		n.SetState(NodeStateAlive)
	default:
		panic(transitionFail(n.state, NodeStateAlive))
	}
	fn()
}

// declareOutOfSync puts a node into OutOfSync state, disconnecting all current
// clients and making it unavailable for use until back in-sync.
func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) declareOutOfSync(isOutOfSync func(num int64, td *utils.Big) bool) {
	n.transitionToOutOfSync(func() {
		n.lfcLog.Errorw("RPC Node is out of sync", "nodeState", n.state)
		n.wg.Add(1)
		go n.outOfSyncLoop(isOutOfSync)
	})
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) transitionToOutOfSync(fn func()) {
	promPoolRPCNodeTransitionsToOutOfSync.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateAlive:
		n.disconnectAll()
		n.SetState(NodeStateOutOfSync)
	default:
		panic(transitionFail(n.state, NodeStateOutOfSync))
	}
	fn()
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) declareUnreachable() {
	n.transitionToUnreachable(func() {
		n.lfcLog.Errorw("RPC Node is unreachable", "nodeState", n.state)
		n.wg.Add(1)
		go n.unreachableLoop()
	})
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) transitionToUnreachable(fn func()) {
	promPoolRPCNodeTransitionsToUnreachable.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateUndialed, NodeStateDialed, NodeStateAlive, NodeStateOutOfSync, NodeStateInvalidChainID:
		n.disconnectAll()
		n.SetState(NodeStateUnreachable)
	default:
		panic(transitionFail(n.state, NodeStateUnreachable))
	}
	fn()
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) declareInvalidChainID() {
	n.transitionToInvalidChainID(func() {
		n.lfcLog.Errorw("RPC Node has the wrong chain ID", "nodeState", n.state)
		n.wg.Add(1)
		go n.invalidChainIDLoop()
	})
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) transitionToInvalidChainID(fn func()) {
	promPoolRPCNodeTransitionsToInvalidChainID.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateDialed, NodeStateOutOfSync:
		n.disconnectAll()
		n.SetState(NodeStateInvalidChainID)
	default:
		panic(transitionFail(n.state, NodeStateInvalidChainID))
	}
	fn()
}

func transitionString(state NodeState) string {
	return fmt.Sprintf("Total number of times node has transitioned to %s", state)
}

func transitionFail(from NodeState, to NodeState) string {
	return fmt.Sprintf("cannot transition from %#v to %#v", from, to)
}
