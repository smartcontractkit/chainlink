package client

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
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

func (fm *fnMock) AssertNumberOfCalls(t *testing.T, n int) {
	assert.Equal(t, n, fm.calls)
}

var _ ethereum.Subscription = (*subMock)(nil)

type subMock struct{ unsubbed bool }

func (s *subMock) Unsubscribe() {
	s.unsubbed = true
}
func (s *subMock) Err() <-chan error { return nil }

func TestUnit_Node_StateTransitions(t *testing.T) {
	s := testutils.NewWSServer(t, testutils.FixtureChainID, func(method string, params gjson.Result) (string, string) {
		return "", ""
	})
	defer s.Close()
	iN := NewNode(TestNodeConfig{}, logger.TestLogger(t), *s.WSURL(), nil, "test node", 42, nil)
	n := iN.(*node)

	assert.Equal(t, NodeStateUndialed, n.State())

	t.Run("setState", func(t *testing.T) {
		n.setState(NodeStateAlive)
		assert.Equal(t, NodeStateAlive, n.state)
		n.setState(NodeStateUndialed)
		assert.Equal(t, NodeStateUndialed, n.state)
	})

	// must dial to set rpc client for use in state transitions
	err := n.dial(context.Background())
	require.NoError(t, err)

	t.Run("transitionToAlive", func(t *testing.T) {
		m := new(fnMock)
		assert.Panics(t, func() {
			n.transitionToAlive(m.Fn)
		})
		m.AssertNotCalled(t)
		n.setState(NodeStateDialed)
		n.transitionToAlive(m.Fn)
		m.AssertNumberOfCalls(t, 1)
		n.setState(NodeStateInvalidChainID)
		n.transitionToAlive(m.Fn)
		m.AssertNumberOfCalls(t, 2)
	})

	t.Run("transitionToInSync", func(t *testing.T) {
		m := new(fnMock)
		n.setState(NodeStateAlive)
		assert.Panics(t, func() {
			n.transitionToInSync(m.Fn)
		})
		m.AssertNotCalled(t)
		n.setState(NodeStateOutOfSync)
		n.transitionToInSync(m.Fn)
		m.AssertCalled(t)
	})
	t.Run("transitionToOutOfSync", func(t *testing.T) {
		m := new(fnMock)
		n.setState(NodeStateOutOfSync)
		assert.Panics(t, func() {
			n.transitionToOutOfSync(m.Fn)
		})
		m.AssertNotCalled(t)
		n.setState(NodeStateAlive)
		n.transitionToOutOfSync(m.Fn)
		m.AssertCalled(t)
	})
	t.Run("transitionToOutOfSync unsubscribes everything", func(t *testing.T) {
		m := new(fnMock)
		n.setState(NodeStateAlive)
		sub := &subMock{}
		n.registerSub(sub)
		n.transitionToOutOfSync(m.Fn)
		m.AssertNumberOfCalls(t, 1)
		assert.True(t, sub.unsubbed)
	})
	t.Run("transitionToUnreachable", func(t *testing.T) {
		m := new(fnMock)
		n.setState(NodeStateUnreachable)
		assert.Panics(t, func() {
			n.transitionToUnreachable(m.Fn)
		})
		m.AssertNotCalled(t)
		n.setState(NodeStateDialed)
		n.transitionToUnreachable(m.Fn)
		m.AssertNumberOfCalls(t, 1)
		n.setState(NodeStateAlive)
		n.transitionToUnreachable(m.Fn)
		m.AssertNumberOfCalls(t, 2)
		n.setState(NodeStateOutOfSync)
		n.transitionToUnreachable(m.Fn)
		m.AssertNumberOfCalls(t, 3)
		n.setState(NodeStateUndialed)
		n.transitionToUnreachable(m.Fn)
		m.AssertNumberOfCalls(t, 4)
		n.setState(NodeStateInvalidChainID)
		n.transitionToUnreachable(m.Fn)
		m.AssertNumberOfCalls(t, 5)
	})
	t.Run("transitionToUnreachable unsubscribes everything", func(t *testing.T) {
		m := new(fnMock)
		n.setState(NodeStateDialed)
		sub := &subMock{}
		n.registerSub(sub)
		n.transitionToUnreachable(m.Fn)
		m.AssertNumberOfCalls(t, 1)
		assert.True(t, sub.unsubbed)
	})
	t.Run("transitionToInvalidChainID", func(t *testing.T) {
		m := new(fnMock)
		n.setState(NodeStateUnreachable)
		assert.Panics(t, func() {
			n.transitionToInvalidChainID(m.Fn)
		})
		m.AssertNotCalled(t)
		n.setState(NodeStateDialed)
		n.transitionToInvalidChainID(m.Fn)
		m.AssertCalled(t)
	})
	t.Run("transitionToInvalidChainID unsubscribes everything", func(t *testing.T) {
		m := new(fnMock)
		n.setState(NodeStateDialed)
		sub := &subMock{}
		n.registerSub(sub)
		n.transitionToInvalidChainID(m.Fn)
		m.AssertNumberOfCalls(t, 1)
		assert.True(t, sub.unsubbed)
	})
	t.Run("Close", func(t *testing.T) {
		// first attempt panics due to node being unstarted
		assert.Panics(t, n.Close)
		// must start to allow closing
		err := n.StartOnce("test node", func() error { return nil })
		assert.NoError(t, err)
		n.Close()
		assert.Equal(t, NodeStateClosed, n.State())
		// second attempt panics due to node being stopped twice
		assert.Panics(t, n.Close)
	})
}
