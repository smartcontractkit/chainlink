package client

import (
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type fnMock struct{ calls int }

func (fm *fnMock) Fn() {
	fm.calls++
}

func (fm *fnMock) AssertNotCalled(t *testing.T) {
	assert.Equal(t, 0, fm.calls)
}

func (fm *fnMock) AssertCalled(t *testing.T) {
	assert.Greater(t, fm.calls, 0)
}

func TestUnit_Node_StateTransitions(t *testing.T) {
	t.Parallel()

	t.Run("setState", func(t *testing.T) {
		n := newTestNode(t, testNodeOpts{rpc: nil, config: testNodeConfig{nodeIsSyncingEnabled: true}})
		assert.Equal(t, NodeStateUndialed, n.State())
		n.setState(NodeStateAlive)
		assert.Equal(t, NodeStateAlive, n.State())
		n.setState(NodeStateUndialed)
		assert.Equal(t, NodeStateUndialed, n.State())
	})

	t.Run("transitionToAlive", func(t *testing.T) {
		const destinationState = NodeStateAlive
		allowedStates := []NodeState{NodeStateDialed, NodeStateInvalidChainID, NodeStateSyncing}
		rpc := newMockRPCClient[types.ID, Head](t)
		testTransition(t, rpc, testNode.transitionToAlive, destinationState, allowedStates...)
	})

	t.Run("transitionToInSync", func(t *testing.T) {
		const destinationState = NodeStateAlive
		allowedStates := []NodeState{NodeStateOutOfSync, NodeStateSyncing}
		rpc := newMockRPCClient[types.ID, Head](t)
		testTransition(t, rpc, testNode.transitionToInSync, destinationState, allowedStates...)
	})
	t.Run("transitionToOutOfSync", func(t *testing.T) {
		const destinationState = NodeStateOutOfSync
		allowedStates := []NodeState{NodeStateAlive}
		rpc := newMockRPCClient[types.ID, Head](t)
		rpc.On("UnsubscribeAllExcept", mock.Anything)
		testTransition(t, rpc, testNode.transitionToOutOfSync, destinationState, allowedStates...)
	})
	t.Run("transitionToUnreachable", func(t *testing.T) {
		const destinationState = NodeStateUnreachable
		allowedStates := []NodeState{NodeStateUndialed, NodeStateDialed, NodeStateAlive, NodeStateOutOfSync, NodeStateInvalidChainID, NodeStateSyncing}
		rpc := newMockRPCClient[types.ID, Head](t)
		rpc.On("UnsubscribeAllExcept", mock.Anything)
		testTransition(t, rpc, testNode.transitionToUnreachable, destinationState, allowedStates...)
	})
	t.Run("transitionToInvalidChain", func(t *testing.T) {
		const destinationState = NodeStateInvalidChainID
		allowedStates := []NodeState{NodeStateDialed, NodeStateOutOfSync, NodeStateSyncing}
		rpc := newMockRPCClient[types.ID, Head](t)
		rpc.On("UnsubscribeAllExcept", mock.Anything)
		testTransition(t, rpc, testNode.transitionToInvalidChainID, destinationState, allowedStates...)
	})
	t.Run("transitionToSyncing", func(t *testing.T) {
		const destinationState = NodeStateSyncing
		allowedStates := []NodeState{NodeStateDialed, NodeStateOutOfSync, NodeStateInvalidChainID}
		rpc := newMockRPCClient[types.ID, Head](t)
		rpc.On("UnsubscribeAllExcept", mock.Anything)
		testTransition(t, rpc, testNode.transitionToSyncing, destinationState, allowedStates...)
	})
	t.Run("transitionToSyncing panics if nodeIsSyncing is disabled", func(t *testing.T) {
		rpc := newMockRPCClient[types.ID, Head](t)
		rpc.On("UnsubscribeAllExcept", mock.Anything)
		node := newTestNode(t, testNodeOpts{rpc: rpc})
		node.setState(NodeStateDialed)
		fn := new(fnMock)
		defer fn.AssertNotCalled(t)
		assert.PanicsWithValue(t, "unexpected transition to NodeStateSyncing, while it's disabled", func() {
			node.transitionToSyncing(fn.Fn)
		})
	})
}

func testTransition(t *testing.T, rpc *mockRPCClient[types.ID, Head], transition func(node testNode, fn func()), destinationState NodeState, allowedStates ...NodeState) {
	node := newTestNode(t, testNodeOpts{rpc: rpc, config: testNodeConfig{nodeIsSyncingEnabled: true}})
	for _, allowedState := range allowedStates {
		m := new(fnMock)
		node.setState(allowedState)
		transition(node, m.Fn)
		assert.Equal(t, destinationState, node.State(), "Expected node to successfully transition from %s to %s state", allowedState, destinationState)
		m.AssertCalled(t)
	}
	// noop on attempt to transition from Closed state
	m := new(fnMock)
	node.setState(NodeStateClosed)
	transition(node, m.Fn)
	m.AssertNotCalled(t)
	assert.Equal(t, NodeStateClosed, node.State(), "Expected node to remain in closed state on transition attempt")

	for _, nodeState := range allNodeStates {
		if slices.Contains(allowedStates, nodeState) || nodeState == NodeStateClosed {
			continue
		}

		m := new(fnMock)
		node.setState(nodeState)
		assert.Panics(t, func() {
			transition(node, m.Fn)
		}, "Expected transition from `%s` to `%s` to panic", nodeState, destinationState)
		m.AssertNotCalled(t)
		assert.Equal(t, nodeState, node.State(), "Expected node to remain in initial state on invalid transition")
	}
}

func TestNodeState_String(t *testing.T) {
	t.Run("Ensure all states are meaningful when converted to string", func(t *testing.T) {
		for _, ns := range allNodeStates {
			// ensure that string representation is not NodeState(%d)
			assert.NotContains(t, ns.String(), strconv.FormatInt(int64(ns), 10), "Expected node state to have readable name")
		}
	})
}
