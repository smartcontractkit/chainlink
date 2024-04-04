package remote_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
)

const (
	eventId1 = "event1"
	eventId2 = "event2"
	peerId1  = "peer1"
	peerId2  = "peer2"
	payloadA = "payloadA"
)

func TestMessageCache_InsertReady(t *testing.T) {
	cache := remote.NewMessageCache[string, string]()

	// not ready with one message
	ts := cache.Insert(eventId1, peerId1, 100, []byte(payloadA))
	require.Equal(t, int64(100), ts)
	ready, _ := cache.Ready(eventId1, 2, 100, true)
	require.False(t, ready)

	// not ready with two messages but only one fresh enough
	ts = cache.Insert(eventId1, peerId2, 200, []byte(payloadA))
	require.Equal(t, int64(100), ts)
	ready, _ = cache.Ready(eventId1, 2, 150, true)
	require.False(t, ready)

	// ready with two messages (once only)
	ready, messages := cache.Ready(eventId1, 2, 100, true)
	require.True(t, ready)
	require.Equal(t, []byte(payloadA), messages[0])
	require.Equal(t, []byte(payloadA), messages[1])

	// not ready again for the same event ID
	ready, _ = cache.Ready(eventId1, 2, 100, true)
	require.False(t, ready)
}

func TestMessageCache_DeleteOlderThan(t *testing.T) {
	cache := remote.NewMessageCache[string, string]()

	ts := cache.Insert(eventId1, peerId1, 100, []byte(payloadA))
	require.Equal(t, int64(100), ts)
	ts = cache.Insert(eventId2, peerId2, 200, []byte(payloadA))
	require.Equal(t, int64(200), ts)

	deleted := cache.DeleteOlderThan(150)
	require.Equal(t, 1, deleted)

	deleted = cache.DeleteOlderThan(150)
	require.Equal(t, 0, deleted)

	deleted = cache.DeleteOlderThan(201)
	require.Equal(t, 1, deleted)
}
