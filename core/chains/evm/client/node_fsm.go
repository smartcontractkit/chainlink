package client

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	promEVMPoolRPCNodeTransitionsToAlive = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_num_transitions_to_alive",
		Help: fmt.Sprintf("Total number of times node has transitioned to %s", NodeStateAlive),
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeTransitionsToInSync = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_num_transitions_to_in_sync",
		Help: fmt.Sprintf("Total number of times node has transitioned from %s to %s", NodeStateOutOfSync, NodeStateAlive),
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeTransitionsToOutOfSync = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_num_transitions_to_out_of_sync",
		Help: fmt.Sprintf("Total number of times node has transitioned to %s", NodeStateOutOfSync),
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeTransitionsToUnreachable = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_num_transitions_to_unreachable",
		Help: fmt.Sprintf("Total number of times node has transitioned to %s", NodeStateUnreachable),
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeTransitionsToInvalidChainID = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_num_transitions_to_invalid_chain_id",
		Help: fmt.Sprintf("Total number of times node has transitioned to %s", NodeStateInvalidChainID),
	}, []string{"evmChainID", "nodeName"})
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

// State allows reading the current state of the node
func (n *node) State() NodeState {
	n.stateMu.RLock()
	defer n.stateMu.RUnlock()
	return n.state
}

// setState is only used by internal state management methods.
// This is low-level; care should be taken by the caller to ensure the new state is a valid transition.
// State changes should always be synchronous: only one goroutine at a time should change state.
// n.stateMu should not be locked for long periods of time because external clients expect a timely response from n.State()
func (n *node) setState(s NodeState) {
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	n.state = s
}

// declareXXX methods change the state and pass conrol off the new state
// management goroutine

func (n *node) declareAlive() {
	n.transitionToAlive(func() {
		n.lfcLog.Infow("RPC Node is online", "nodeState", n.state)
		n.wg.Add(1)
		go n.aliveLoop()
	})
}

func (n *node) transitionToAlive(fn func()) {
	promEVMPoolRPCNodeTransitionsToAlive.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateDialed, NodeStateInvalidChainID:
		n.state = NodeStateAlive
	default:
		panic(fmt.Sprintf("cannot transition from %#v to %#v", n.state, NodeStateAlive))
	}
	fn()
}

// declareInSync puts a node back into Alive state, allowing it to be used by
// pool consumers again
func (n *node) declareInSync() {
	n.transitionToInSync(func() {
		n.lfcLog.Infow("RPC Node is back in sync", "nodeState", n.state)
		n.wg.Add(1)
		go n.aliveLoop()
	})
}

func (n *node) transitionToInSync(fn func()) {
	promEVMPoolRPCNodeTransitionsToAlive.WithLabelValues(n.chainID.String(), n.name).Inc()
	promEVMPoolRPCNodeTransitionsToInSync.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateOutOfSync:
		n.state = NodeStateAlive
	default:
		panic(fmt.Sprintf("cannot transition from %#v to %#v", n.state, NodeStateAlive))
	}
	fn()
}

// declareOutOfSync puts a node into OutOfSync state, disconnecting all current
// clients and making it unavailable for use
func (n *node) declareOutOfSync(latestReceivedBlockNumber int64) {
	n.transitionToOutOfSync(func() {
		n.lfcLog.Errorw("RPC Node is out of sync", "nodeState", n.state)
		n.wg.Add(1)
		go n.outOfSyncLoop(latestReceivedBlockNumber)
	})
}

func (n *node) transitionToOutOfSync(fn func()) {
	promEVMPoolRPCNodeTransitionsToOutOfSync.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateAlive:
		n.disconnectAll()
		n.state = NodeStateOutOfSync
	default:
		panic(fmt.Sprintf("cannot transition from %#v to %#v", n.state, NodeStateOutOfSync))
	}
	fn()
}

func (n *node) declareUnreachable() {
	n.transitionToUnreachable(func() {
		n.lfcLog.Errorw("RPC Node is unreachable", "nodeState", n.state)
		n.wg.Add(1)
		go n.unreachableLoop()
	})
}

func (n *node) transitionToUnreachable(fn func()) {
	promEVMPoolRPCNodeTransitionsToUnreachable.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateUndialed, NodeStateDialed, NodeStateAlive, NodeStateOutOfSync, NodeStateInvalidChainID:
		n.disconnectAll()
		n.state = NodeStateUnreachable
	default:
		panic(fmt.Sprintf("cannot transition from %#v to %#v", n.state, NodeStateUnreachable))
	}
	fn()
}

func (n *node) declareInvalidChainID() {
	n.transitionToInvalidChainID(func() {
		n.lfcLog.Errorw("RPC Node has the wrong chain ID", "nodeState", n.state)
		n.wg.Add(1)
		go n.invalidChainIDLoop()
	})
}

func (n *node) transitionToInvalidChainID(fn func()) {
	promEVMPoolRPCNodeTransitionsToInvalidChainID.WithLabelValues(n.chainID.String(), n.name).Inc()
	n.stateMu.Lock()
	defer n.stateMu.Unlock()
	if n.state == NodeStateClosed {
		return
	}
	switch n.state {
	case NodeStateDialed:
		n.disconnectAll()
		n.state = NodeStateInvalidChainID
	default:
		panic(fmt.Sprintf("cannot transition from %#v to %#v", n.state, NodeStateInvalidChainID))
	}
	fn()
}
