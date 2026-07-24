package capabilityrunner

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// countingBroadcaster counts how many upstream subscriptions are opened against
// it, so a test can assert the singleton only ever opens one.
type countingBroadcaster struct {
	*loop.AtomicSettings
	mu   sync.Mutex
	subs int
}

func (c *countingBroadcaster) Subscribe(ctx context.Context) (<-chan core.SettingsUpdate, error) {
	c.mu.Lock()
	c.subs++
	c.mu.Unlock()
	return c.AtomicSettings.Subscribe(ctx)
}

func (c *countingBroadcaster) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.subs
}

func TestFileBackedSettings_WritesOnceBeforeSubscribersSee(t *testing.T) {
	lggr := logger.TestLogger(t)
	path := filepath.Join(t.TempDir(), "settings.txt")
	up := &countingBroadcaster{AtomicSettings: &loop.AtomicSettings{}}
	f := NewFileBackedSettings(lggr, up, path)

	const nSubs = 8
	chans := make([]<-chan core.SettingsUpdate, nSubs)
	for i := range chans {
		ch, err := f.Subscribe(t.Context())
		require.NoError(t, err)
		chans[i] = ch
	}
	// One upstream subscription no matter how many subscribers.
	require.Equal(t, 1, up.count())

	var wg sync.WaitGroup
	seen := make([]atomic.Int64, nSubs)
	for i, ch := range chans {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range ch {
				// The file must already be on disk, and already current.
				got, err := os.ReadFile(path)
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, "A = 'one'\n", string(got))
				seen[i].Add(1)
			}
		}()
	}

	require.NoError(t, up.Store(core.SettingsUpdate{Settings: "A = 'one'\n", Hash: "h1"}))

	require.Eventually(t, func() bool {
		for i := range seen {
			if seen[i].Load() == 0 {
				return false
			}
		}
		return true
	}, 5*time.Second, 10*time.Millisecond, "every subscriber should see the update")

	for _, ch := range chans {
		f.Unsubscribe(ch)
	}
	wg.Wait()
}

func TestFileBackedSettings_UnsubscribeStopsDelivery(t *testing.T) {
	lggr := logger.TestLogger(t)
	path := filepath.Join(t.TempDir(), "settings.txt")
	as := &loop.AtomicSettings{}
	f := NewFileBackedSettings(lggr, as, path)

	ch, err := f.Subscribe(t.Context())
	require.NoError(t, err)
	f.Unsubscribe(ch)

	_, ok := <-ch
	require.False(t, ok, "unsubscribed channel should be closed")

	// Still writes even with no subscribers left.
	require.NoError(t, as.Store(core.SettingsUpdate{Settings: "B = 'two'\n", Hash: "h2"}))
	require.Eventually(t, func() bool {
		got, err := os.ReadFile(path)
		return err == nil && string(got) == "B = 'two'\n"
	}, 5*time.Second, 10*time.Millisecond)
}
