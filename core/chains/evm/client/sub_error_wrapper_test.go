package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestSubscriptionErrorWrapper(t *testing.T) {
	t.Parallel()
	t.Run("Unsubscribe wrapper releases resources", func(t *testing.T) {
		t.Parallel()

		mockedSub := NewMockSubscription()
		const prefix = "RPC returned error"
		wrapper := newSubscriptionErrorWrapper(mockedSub, prefix)
		wrapper.Unsubscribe()

		// mock's resources were relased
		assert.True(t, mockedSub.unsubscribed)
		_, ok := <-mockedSub.Err()
		assert.False(t, ok)
		// wrapper's channels are closed
		_, ok = <-wrapper.Err()
		assert.False(t, ok)
		//  subsequence unsubscribe does not causes panic
		wrapper.Unsubscribe()
	})
	t.Run("Unsubscribe interrupts error delivery", func(t *testing.T) {
		t.Parallel()
		sub := NewMockSubscription()
		const prefix = "RPC returned error"
		wrapper := newSubscriptionErrorWrapper(sub, prefix)
		sub.Errors <- fmt.Errorf("error")

		wrapper.Unsubscribe()
		_, ok := <-wrapper.Err()
		assert.False(t, ok)
	})
	t.Run("Successfully wraps error", func(t *testing.T) {
		t.Parallel()
		sub := NewMockSubscription()
		const prefix = "RPC returned error"
		wrapper := newSubscriptionErrorWrapper(sub, prefix)
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
		wrapper := newSubscriptionErrorWrapper(mockedSub, "")

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
