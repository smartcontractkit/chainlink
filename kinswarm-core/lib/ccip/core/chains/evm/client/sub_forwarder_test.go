package client

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

func TestChainIDSubForwarder(t *testing.T) {
	t.Parallel()

	newChainIDSubForwarder := func(chainID *big.Int, ch chan<- *evmtypes.Head) *subForwarder[*evmtypes.Head] {
		return newSubForwarder(ch, func(head *evmtypes.Head) (*evmtypes.Head, error) {
			head.EVMChainID = ubig.New(chainID)
			return head, nil
		}, nil)
	}

	chainID := big.NewInt(123)

	t.Run("unsubscribe forwarder", func(t *testing.T) {
		t.Parallel()

		ch := make(chan *evmtypes.Head)
		forwarder := newChainIDSubForwarder(chainID, ch)
		sub := NewMockSubscription()
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
		sub := NewMockSubscription()
		err := forwarder.start(sub, nil)
		assert.NoError(t, err)
		expectedError := errors.New("boo")
		sub.Errors <- expectedError
		forwarder.Unsubscribe()

		assert.True(t, sub.unsubscribed)
		err, ok := <-forwarder.Err()
		assert.True(t, ok)
		require.ErrorIs(t, err, expectedError)
		_, ok = <-forwarder.Err()
		assert.False(t, ok)
	})

	t.Run("unsubscribe forwarder with message", func(t *testing.T) {
		t.Parallel()

		ch := make(chan *evmtypes.Head)
		forwarder := newChainIDSubForwarder(chainID, ch)
		sub := NewMockSubscription()
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
		sub := NewMockSubscription()
		errIn := errors.New("foo")
		errOut := forwarder.start(sub, errIn)
		assert.Equal(t, errIn, errOut)
	})

	t.Run("forwarding", func(t *testing.T) {
		t.Parallel()

		ch := make(chan *evmtypes.Head)
		forwarder := newChainIDSubForwarder(chainID, ch)
		sub := NewMockSubscription()
		err := forwarder.start(sub, nil)
		assert.NoError(t, err)

		head := &evmtypes.Head{
			ID: 1,
		}
		forwarder.srcCh <- head
		receivedHead := <-ch
		assert.Equal(t, head, receivedHead)
		assert.Equal(t, ubig.New(chainID), receivedHead.EVMChainID)

		expectedErr := errors.New("error")
		sub.Errors <- expectedErr
		receivedErr := <-forwarder.Err()
		assert.Equal(t, expectedErr, receivedErr)
	})
}

func TestSubscriptionForwarder(t *testing.T) {
	t.Run("Error returned by interceptResult is forwarded to err channel", func(t *testing.T) {
		t.Parallel()

		ch := make(chan *evmtypes.Head)
		expectedErr := errors.New("something went wrong during result interception")
		forwarder := newSubForwarder(ch, func(head *evmtypes.Head) (*evmtypes.Head, error) {
			return nil, expectedErr
		}, nil)
		mockedSub := NewMockSubscription()
		require.NoError(t, forwarder.start(mockedSub, nil))

		head := &evmtypes.Head{
			ID: 1,
		}
		forwarder.srcCh <- head
		err := <-forwarder.Err()
		require.ErrorIs(t, err, expectedErr)
		// ensure forwarder is closed
		_, ok := <-forwarder.Err()
		assert.False(t, ok)
		assert.True(t, mockedSub.unsubscribed)
	})
}

func TestSubscriptionErrorWrapper(t *testing.T) {
	t.Parallel()
	newSubscriptionErrorWrapper := func(t *testing.T, sub commontypes.Subscription, errorPrefix string) ethereum.Subscription {
		ch := make(chan *evmtypes.Head)
		result := newSubForwarder(ch, nil, func(err error) error {
			return fmt.Errorf("%s: %w", errorPrefix, err)
		})
		require.NoError(t, result.start(sub, nil))
		return result
	}
	t.Run("Unsubscribe wrapper releases resources", func(t *testing.T) {
		t.Parallel()

		mockedSub := NewMockSubscription()
		const prefix = "RPC returned error"
		wrapper := newSubscriptionErrorWrapper(t, mockedSub, prefix)
		wrapper.Unsubscribe()

		// mock's resources were released
		assert.True(t, mockedSub.unsubscribed)
		_, ok := <-mockedSub.Err()
		assert.False(t, ok)
		// wrapper's channels are closed
		_, ok = <-wrapper.Err()
		assert.False(t, ok)
		//  subsequence unsubscribe does not causes panic
		wrapper.Unsubscribe()
	})
	t.Run("Successfully wraps error", func(t *testing.T) {
		t.Parallel()
		sub := NewMockSubscription()
		const prefix = "RPC returned error"
		wrapper := newSubscriptionErrorWrapper(t, sub, prefix)
		sub.Errors <- fmt.Errorf("root error")

		err, ok := <-wrapper.Err()
		assert.True(t, ok)
		assert.Equal(t, "RPC returned error: root error", err.Error())

		wrapper.Unsubscribe()
		_, ok = <-wrapper.Err()
		assert.False(t, ok)
	})
	t.Run("Unsubscribe on root does not cause panic", func(t *testing.T) {
		t.Parallel()
		mockedSub := NewMockSubscription()
		wrapper := newSubscriptionErrorWrapper(t, mockedSub, "")

		mockedSub.Unsubscribe()
		// mock's resources were released
		assert.True(t, mockedSub.unsubscribed)
		_, ok := <-mockedSub.Err()
		assert.False(t, ok)
		// wrapper's channels are eventually closed
		tests.AssertEventually(t, func() bool {
			_, ok = <-wrapper.Err()
			return !ok
		})
	})
}
