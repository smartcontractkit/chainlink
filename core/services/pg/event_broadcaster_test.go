package pg_test

import (
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/pg"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
)

func TestEventBroadcaster(t *testing.T) {
	config, _ := heavyweight.FullTestDBNoFixtures(t, "event_broadcaster")

	eventBroadcaster := cltest.NewEventBroadcaster(t, config.DatabaseURL())
	require.NoError(t, eventBroadcaster.Start(testutils.Context(t)))
	t.Cleanup(func() { require.NoError(t, eventBroadcaster.Close()) })

	t.Run("doesn't broadcast unrelated events (no payload filter)", func(t *testing.T) {
		sub, err := eventBroadcaster.Subscribe("foo", "")
		require.NoError(t, err)
		defer sub.Close()

		go func() {
			err := eventBroadcaster.Notify("bar", "123")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("fooo", "123")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("fo", "123")
			require.NoError(t, err)
		}()

		ch := sub.Events()
		gomega.NewWithT(t).Consistently(ch).ShouldNot(gomega.Receive())
	})

	t.Run("doesn't broadcast unrelated events (with payload filter)", func(t *testing.T) {
		sub, err := eventBroadcaster.Subscribe("foo", "123")
		require.NoError(t, err)
		defer sub.Close()

		go func() {
			err := eventBroadcaster.Notify("foo", "asdf")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("bar", "123")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("fooo", "123")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("fo", "123")
			require.NoError(t, err)
		}()

		ch := sub.Events()
		gomega.NewWithT(t).Consistently(ch).ShouldNot(gomega.Receive())
	})

	t.Run("does broadcast related events (no payload filter)", func(t *testing.T) {
		sub, err := eventBroadcaster.Subscribe("foo", "")
		require.NoError(t, err)
		defer sub.Close()

		go func() {
			err := eventBroadcaster.Notify("foo", "123")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("foo", "aslkdjslkdfj")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("foo", "true")
			require.NoError(t, err)
		}()

		ch := sub.Events()
		gomega.NewWithT(t).Eventually(ch).Should(gomega.Receive())
		gomega.NewWithT(t).Eventually(ch).Should(gomega.Receive())
		gomega.NewWithT(t).Eventually(ch).Should(gomega.Receive())
	})

	t.Run("does broadcast related events (with payload filter)", func(t *testing.T) {
		sub, err := eventBroadcaster.Subscribe("foo", "123")
		require.NoError(t, err)
		defer sub.Close()

		go func() {
			err := eventBroadcaster.Notify("foo", "asdf")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("foo", "123")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("foo", "123")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("foo", "true")
			require.NoError(t, err)
		}()

		ch := sub.Events()
		gomega.NewWithT(t).Eventually(ch).Should(gomega.Receive())
		gomega.NewWithT(t).Eventually(ch).Should(gomega.Receive())
		gomega.NewWithT(t).Consistently(ch).ShouldNot(gomega.Receive())
	})

	t.Run("broadcasts to the correct subscribers", func(t *testing.T) {
		sub1, err := eventBroadcaster.Subscribe("foo", "")
		require.NoError(t, err)
		defer sub1.Close()

		sub2, err := eventBroadcaster.Subscribe("foo", "123")
		require.NoError(t, err)
		defer sub2.Close()

		sub3, err := eventBroadcaster.Subscribe("bar", "")
		require.NoError(t, err)
		defer sub3.Close()

		sub4, err := eventBroadcaster.Subscribe("bar", "asdf")
		require.NoError(t, err)
		defer sub4.Close()

		var wg sync.WaitGroup
		wg.Add(5)

		recv := func(ch <-chan pg.Event) pg.Event {
			select {
			case e := <-ch:
				return e
			case <-time.After(5 * time.Second):
				t.Fatal("did not receive")
			}
			return pg.Event{}
		}

		go func() {
			defer wg.Done()
			err := eventBroadcaster.Notify("foo", "asdf")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("foo", "123")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("foo", "123")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("foo", "true")
			require.NoError(t, err)

			err = eventBroadcaster.Notify("bar", "asdf")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("bar", "123")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("bar", "123")
			require.NoError(t, err)
			err = eventBroadcaster.Notify("bar", "true")
			require.NoError(t, err)
		}()

		go func() {
			defer wg.Done()
			e := recv(sub1.Events())
			require.Equal(t, "foo", e.Channel)
			require.Equal(t, "asdf", e.Payload)

			e = recv(sub1.Events())
			require.Equal(t, "foo", e.Channel)
			require.Equal(t, "123", e.Payload)

			e = recv(sub1.Events())
			require.Equal(t, "foo", e.Channel)
			require.Equal(t, "123", e.Payload)

			e = recv(sub1.Events())
			require.Equal(t, "foo", e.Channel)
			require.Equal(t, "true", e.Payload)

			gomega.NewWithT(t).Consistently(sub1.Events()).ShouldNot(gomega.Receive())
		}()

		go func() {
			defer wg.Done()
			e := recv(sub2.Events())
			require.Equal(t, "foo", e.Channel)
			require.Equal(t, "123", e.Payload)

			e = recv(sub2.Events())
			require.Equal(t, "foo", e.Channel)
			require.Equal(t, "123", e.Payload)

			gomega.NewWithT(t).Consistently(sub2.Events()).ShouldNot(gomega.Receive())
		}()

		go func() {
			defer wg.Done()
			e := recv(sub3.Events())
			require.Equal(t, "bar", e.Channel)
			require.Equal(t, "asdf", e.Payload)

			e = recv(sub3.Events())
			require.Equal(t, "bar", e.Channel)
			require.Equal(t, "123", e.Payload)

			e = recv(sub3.Events())
			require.Equal(t, "bar", e.Channel)
			require.Equal(t, "123", e.Payload)

			e = recv(sub3.Events())
			require.Equal(t, "bar", e.Channel)
			require.Equal(t, "true", e.Payload)

			gomega.NewWithT(t).Consistently(sub3.Events()).ShouldNot(gomega.Receive())
		}()

		go func() {
			defer wg.Done()
			e := recv(sub4.Events())
			require.Equal(t, "bar", e.Channel)
			require.Equal(t, "asdf", e.Payload)

			gomega.NewWithT(t).Consistently(sub4.Events()).ShouldNot(gomega.Receive())
		}()

		wg.Wait()
	})
}
