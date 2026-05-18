package client

import (
	"slices"
	"strconv"
	"testing"

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
		assert.Equal(t, nodeStateUndialed, n.State())
		n.setState(nodeStateAlive)
		assert.Equal(t, nodeStateAlive, n.State())
		n.setState(nodeStateUndialed)
		assert.Equal(t, nodeStateUndialed, n.State())
	})

	t.Run("transitionToAlive", func(t *testing.T) {
		const destinationState = nodeStateAlive
		allowedStates := []nodeState{nodeStateDialed, nodeStateInvalidChainID, nodeStateSyncing}
		rpc := newMockNodeClient[types.ID, Head](t)
		testTransition(t, rpc, testNode.transitionToAlive, destinationState, allowedStates...)
	})

	t.Run("transitionToInSync", func(t *testing.T) {
		const destinationState = nodeStateAlive
		allowedStates := []nodeState{nodeStateOutOfSync, nodeStateSyncing}
		rpc := newMockNodeClient[types.ID, Head](t)
		testTransition(t, rpc, testNode.transitionToInSync, destinationState, allowedStates...)
	})
	t.Run("transitionToOutOfSync", func(t *testing.T) {
		const destinationState = nodeStateOutOfSync
		allowedStates := []nodeState{nodeStateAlive}
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("DisconnectAll").Once()
		testTransition(t, rpc, testNode.transitionToOutOfSync, destinationState, allowedStates...)
	})
	t.Run("transitionToUnreachable", func(t *testing.T) {
		const destinationState = nodeStateUnreachable
		allowedStates := []nodeState{nodeStateUndialed, nodeStateDialed, nodeStateAlive, nodeStateOutOfSync, nodeStateInvalidChainID, nodeStateSyncing}
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("DisconnectAll").Times(len(allowedStates))
		testTransition(t, rpc, testNode.transitionToUnreachable, destinationState, allowedStates...)
	})
	t.Run("transitionToInvalidChain", func(t *testing.T) {
		const destinationState = nodeStateInvalidChainID
		allowedStates := []nodeState{nodeStateDialed, nodeStateOutOfSync, nodeStateSyncing}
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("DisconnectAll").Times(len(allowedStates))
		testTransition(t, rpc, testNode.transitionToInvalidChainID, destinationState, allowedStates...)
	})
	t.Run("transitionToSyncing", func(t *testing.T) {
		const destinationState = nodeStateSyncing
		allowedStates := []nodeState{nodeStateDialed, nodeStateOutOfSync, nodeStateInvalidChainID}
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("DisconnectAll").Times(len(allowedStates))
		testTransition(t, rpc, testNode.transitionToSyncing, destinationState, allowedStates...)
	})
	t.Run("transitionToSyncing panics if nodeIsSyncing is disabled", func(t *testing.T) {
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("DisconnectAll").Once()
		node := newTestNode(t, testNodeOpts{rpc: rpc})
		node.setState(nodeStateDialed)
		fn := new(fnMock)
		defer fn.AssertNotCalled(t)
		assert.PanicsWithValue(t, "unexpected transition to nodeStateSyncing, while it's disabled", func() {
			node.transitionToSyncing(fn.Fn)
		})
	})
}

func testTransition(t *testing.T, rpc *mockNodeClient[types.ID, Head], transition func(node testNode, fn func()), destinationState nodeState, allowedStates ...nodeState) {
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
	node.setState(nodeStateClosed)
	transition(node, m.Fn)
	m.AssertNotCalled(t)
	assert.Equal(t, nodeStateClosed, node.State(), "Expected node to remain in closed state on transition attempt")

	for _, nodeState := range allNodeStates {
		if slices.Contains(allowedStates, nodeState) || nodeState == nodeStateClosed {
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
			// ensure that string representation is not nodeState(%d)
			assert.NotContains(t, ns.String(), strconv.FormatInt(int64(ns), 10), "Expected node state to have readable name")
		}
	})
}
