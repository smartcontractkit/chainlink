package client

import (
	"slices"
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

func newTestTransitionNode(t *testing.T, rpc *mockNodeClient[types.ID, Head]) testNode {
	return newTestNode(t, testNodeOpts{rpc: rpc})
}

func TestUnit_Node_StateTransitions(t *testing.T) {
	t.Parallel()

	t.Run("setState", func(t *testing.T) {
		n := newTestTransitionNode(t, nil)
		assert.Equal(t, nodeStateUndialed, n.State())
		n.setState(nodeStateAlive)
		assert.Equal(t, nodeStateAlive, n.State())
		n.setState(nodeStateUndialed)
		assert.Equal(t, nodeStateUndialed, n.State())
	})

	t.Run("transitionToAlive", func(t *testing.T) {
		const destinationState = nodeStateAlive
		allowedStates := []nodeState{nodeStateDialed, nodeStateInvalidChainID}
		rpc := newMockNodeClient[types.ID, Head](t)
		testTransition(t, rpc, testNode.transitionToAlive, destinationState, allowedStates...)
	})

	t.Run("transitionToInSync", func(t *testing.T) {
		const destinationState = nodeStateAlive
		allowedStates := []nodeState{nodeStateOutOfSync}
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
		allowedStates := []nodeState{nodeStateUndialed, nodeStateDialed, nodeStateAlive, nodeStateOutOfSync, nodeStateInvalidChainID}
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("DisconnectAll").Times(len(allowedStates))
		testTransition(t, rpc, testNode.transitionToUnreachable, destinationState, allowedStates...)
	})
	t.Run("transitionToInvalidChain", func(t *testing.T) {
		const destinationState = nodeStateInvalidChainID
		allowedStates := []nodeState{nodeStateDialed, nodeStateOutOfSync}
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("DisconnectAll").Times(len(allowedStates))
		testTransition(t, rpc, testNode.transitionToInvalidChainID, destinationState, allowedStates...)
	})
}

func testTransition(t *testing.T, rpc *mockNodeClient[types.ID, Head], transition func(node testNode, fn func()), destinationState nodeState, allowedStates ...nodeState) {
	node := newTestTransitionNode(t, rpc)
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
