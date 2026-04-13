package vault

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestReplayGuard_FirstCallSucceeds(t *testing.T) {
	guard := NewRequestReplayGuard()
	futureExpiry := time.Now().UTC().Unix() + 100

	err := guard.CheckAndRecord("digest-1", futureExpiry)
	require.NoError(t, err)
}

func TestRequestReplayGuard_DuplicateRejected(t *testing.T) {
	guard := NewRequestReplayGuard()
	futureExpiry := time.Now().UTC().Unix() + 100

	err := guard.CheckAndRecord("digest-1", futureExpiry)
	require.NoError(t, err)

	err = guard.CheckAndRecord("digest-1", futureExpiry)
	require.ErrorIs(t, err, ErrRequestAlreadySeen)
}

func TestRequestReplayGuard_DifferentDigestsIndependent(t *testing.T) {
	guard := NewRequestReplayGuard()
	futureExpiry := time.Now().UTC().Unix() + 100

	require.NoError(t, guard.CheckAndRecord("digest-1", futureExpiry))
	require.NoError(t, guard.CheckAndRecord("digest-2", futureExpiry))
	require.NoError(t, guard.CheckAndRecord("digest-3", futureExpiry))

	require.ErrorIs(t, guard.CheckAndRecord("digest-1", futureExpiry), ErrRequestAlreadySeen)
	require.ErrorIs(t, guard.CheckAndRecord("digest-2", futureExpiry), ErrRequestAlreadySeen)
}

func TestRequestReplayGuard_ExpiredEntryCleaned(t *testing.T) {
	guard := NewRequestReplayGuard()
	now := time.Now()
	guard.nowFunc = func() time.Time { return now }

	pastExpiry := now.UTC().Unix() - 10
	err := guard.CheckAndRecord("digest-1", pastExpiry)
	require.NoError(t, err)

	// Advance time past the expiry — next call should clean up the entry
	guard.nowFunc = func() time.Time { return now.Add(20 * time.Second) }

	// Same digest should succeed because the expired entry was cleaned up
	err = guard.CheckAndRecord("digest-1", now.Add(20*time.Second).UTC().Unix()+100)
	require.NoError(t, err)
}

func TestRequestReplayGuard_NonExpiredEntryRetained(t *testing.T) {
	guard := NewRequestReplayGuard()
	now := time.Now()
	guard.nowFunc = func() time.Time { return now }

	futureExpiry := now.UTC().Unix() + 100
	require.NoError(t, guard.CheckAndRecord("digest-1", futureExpiry))

	// Advance time, but NOT past the expiry
	guard.nowFunc = func() time.Time { return now.Add(50 * time.Second) }

	err := guard.CheckAndRecord("digest-1", futureExpiry)
	require.ErrorIs(t, err, ErrRequestAlreadySeen)
}

func TestRequestReplayGuard_MixedExpiryCleanup(t *testing.T) {
	guard := NewRequestReplayGuard()
	now := time.Now()
	guard.nowFunc = func() time.Time { return now }

	shortExpiry := now.UTC().Unix() + 10
	longExpiry := now.UTC().Unix() + 200

	require.NoError(t, guard.CheckAndRecord("short-lived", shortExpiry))
	require.NoError(t, guard.CheckAndRecord("long-lived", longExpiry))

	// Advance past short expiry but before long expiry
	guard.nowFunc = func() time.Time { return now.Add(50 * time.Second) }

	// Short-lived should be re-recordable (cleaned up)
	require.NoError(t, guard.CheckAndRecord("short-lived", now.Add(50*time.Second).UTC().Unix()+100))

	// Long-lived should still be rejected
	require.ErrorIs(t, guard.CheckAndRecord("long-lived", longExpiry), ErrRequestAlreadySeen)
}

func TestRequestReplayGuard_ConcurrentAccess(t *testing.T) {
	guard := NewRequestReplayGuard()
	futureExpiry := time.Now().UTC().Unix() + 100

	const goroutines = 100
	results := make([]error, goroutines)
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := range goroutines {
		go func(idx int) {
			defer wg.Done()
			results[idx] = guard.CheckAndRecord("same-digest", futureExpiry)
		}(i)
	}
	wg.Wait()

	successCount := 0
	duplicateCount := 0
	for _, err := range results {
		if err == nil {
			successCount++
		} else {
			require.ErrorIs(t, err, ErrRequestAlreadySeen)
			duplicateCount++
		}
	}

	assert.Equal(t, 1, successCount, "exactly one goroutine should succeed")
	assert.Equal(t, goroutines-1, duplicateCount, "all others should be rejected as duplicates")
}

func TestRequestReplayGuard_ConcurrentDifferentDigests(t *testing.T) {
	guard := NewRequestReplayGuard()
	futureExpiry := time.Now().UTC().Unix() + 100

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	errors := make([]error, goroutines)
	for i := range goroutines {
		go func(idx int) {
			defer wg.Done()
			digest := "digest-" + string(rune('A'+idx))
			errors[idx] = guard.CheckAndRecord(digest, futureExpiry)
		}(i)
	}
	wg.Wait()

	for i, err := range errors {
		assert.NoError(t, err, "goroutine %d should succeed for unique digest", i)
	}
}

func TestRequestReplayGuard_ClearExpiredIndependently(t *testing.T) {
	guard := NewRequestReplayGuard()
	now := time.Now()
	guard.nowFunc = func() time.Time { return now }

	shortExpiry := now.UTC().Unix() + 5
	longExpiry := now.UTC().Unix() + 200

	require.NoError(t, guard.CheckAndRecord("ephemeral", shortExpiry))
	require.NoError(t, guard.CheckAndRecord("durable", longExpiry))

	// Advance past the short expiry
	guard.nowFunc = func() time.Time { return now.Add(30 * time.Second) }

	// ClearExpired should prune without needing a CheckAndRecord call
	guard.ClearExpired()

	guard.mu.Lock()
	_, ephemeralPresent := guard.seen["ephemeral"]
	_, durablePresent := guard.seen["durable"]
	guard.mu.Unlock()

	assert.False(t, ephemeralPresent, "expired entry should have been pruned")
	assert.True(t, durablePresent, "non-expired entry should remain")
}

func TestRequestReplayGuard_EmptyDigest(t *testing.T) {
	guard := NewRequestReplayGuard()
	futureExpiry := time.Now().UTC().Unix() + 100

	require.NoError(t, guard.CheckAndRecord("", futureExpiry))
	require.ErrorIs(t, guard.CheckAndRecord("", futureExpiry), ErrRequestAlreadySeen)
}
