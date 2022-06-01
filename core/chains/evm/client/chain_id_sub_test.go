package client

import (
	"errors"
	"math/big"
	"testing"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
)

type mockSubscription struct {
	unsubscribed bool
	Errors       chan error
}

func newMockSubscription() *mockSubscription {
	return &mockSubscription{Errors: make(chan error)}
}

func (mes *mockSubscription) Err() <-chan error { return mes.Errors }

func (mes *mockSubscription) Unsubscribe() {
	mes.unsubscribed = true
	close(mes.Errors)
}

func TestChainIDSubForwarder(t *testing.T) {
	t.Parallel()

	chainID := big.NewInt(123)

	t.Run("unsubscribe forwarder", func(t *testing.T) {
		t.Parallel()

		ch := make(chan *evmtypes.Head)
		forwarder := newChainIDSubForwarder(chainID, ch)
		sub := newMockSubscription()
		err := forwarder.start(sub, nil)
		assert.NoError(t, err)
		forwarder.Unsubscribe()

		assert.True(t, sub.unsubscribed)
		_, ok := <-sub.Err()
		assert.False(t, ok)
		_, ok = <-forwarder.Err()
		assert.False(t, ok)
	})

	t.Run("unsubscribe forwarder with error", func(t *testing.T) {
		t.Parallel()

		ch := make(chan *evmtypes.Head)
		forwarder := newChainIDSubForwarder(chainID, ch)
		sub := newMockSubscription()
		err := forwarder.start(sub, nil)
		assert.NoError(t, err)
		sub.Errors <- errors.New("boo")
		forwarder.Unsubscribe()

		assert.True(t, sub.unsubscribed)
		_, ok := <-sub.Err()
		assert.False(t, ok)
		_, ok = <-forwarder.Err()
		assert.False(t, ok)
	})

	t.Run("unsubscribe forwarder with message", func(t *testing.T) {
		t.Parallel()

		ch := make(chan *evmtypes.Head)
		forwarder := newChainIDSubForwarder(chainID, ch)
		sub := newMockSubscription()
		err := forwarder.start(sub, nil)
		assert.NoError(t, err)
		forwarder.srcCh <- &evmtypes.Head{}
		forwarder.Unsubscribe()

		assert.True(t, sub.unsubscribed)
		_, ok := <-sub.Err()
		assert.False(t, ok)
		_, ok = <-forwarder.Err()
		assert.False(t, ok)
	})

	t.Run("non nil error parameter", func(t *testing.T) {
		t.Parallel()

		ch := make(chan *evmtypes.Head)
		forwarder := newChainIDSubForwarder(chainID, ch)
		sub := newMockSubscription()
		errIn := errors.New("foo")
		errOut := forwarder.start(sub, errIn)
		assert.Equal(t, errIn, errOut)
	})

	t.Run("forwarding", func(t *testing.T) {
		t.Parallel()

		ch := make(chan *evmtypes.Head)
		forwarder := newChainIDSubForwarder(chainID, ch)
		sub := newMockSubscription()
		err := forwarder.start(sub, nil)
		assert.NoError(t, err)

		head := &evmtypes.Head{
			ID: 1,
		}
		forwarder.srcCh <- head
		receivedHead := <-ch
		assert.Equal(t, head, receivedHead)
		assert.Equal(t, utils.NewBig(chainID), receivedHead.EVMChainID)

		expectedErr := errors.New("error")
		sub.Errors <- expectedErr
		receivedErr := <-forwarder.Err()
		assert.Equal(t, expectedErr, receivedErr)
	})
}
