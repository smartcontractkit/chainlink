package messagecache_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/messagecache"
)

const (
	eventID1 = "event1"
	eventID2 = "event2"
	peerID1  = "peer1"
	peerID2  = "peer2"
	payloadA = "payloadA"
)

func TestMessageCache_InsertReady(t *testing.T) {
	cache := messagecache.NewMessageCache[string, string]()

	// not ready with one message
	ts := cache.Insert(eventID1, peerID1, 100, []byte(payloadA))
	require.Equal(t, int64(100), ts)
	ready, _ := cache.Ready(eventID1, 2, 100, true)
	require.False(t, ready)

	// not ready with two messages but only one fresh enough
	ts = cache.Insert(eventID1, peerID2, 200, []byte(payloadA))
	require.Equal(t, int64(100), ts)
	ready, _ = cache.Ready(eventID1, 2, 150, true)
	require.False(t, ready)

	// ready with two messages (once only)
	ready, messages := cache.Ready(eventID1, 2, 100, true)
	require.True(t, ready)
	require.Equal(t, []byte(payloadA), messages[0])
	require.Equal(t, []byte(payloadA), messages[1])

	// not ready again for the same event ID
	ready, _ = cache.Ready(eventID1, 2, 100, true)
	require.False(t, ready)
}

func TestMessageCache_DeleteOlderThan(t *testing.T) {
	cache := messagecache.NewMessageCache[string, string]()

	ts := cache.Insert(eventID1, peerID1, 100, []byte(payloadA))
	require.Equal(t, int64(100), ts)
	ts = cache.Insert(eventID2, peerID2, 200, []byte(payloadA))
	require.Equal(t, int64(200), ts)

	deleted := cache.DeleteOlderThan(150)
	require.Equal(t, 1, deleted)

	deleted = cache.DeleteOlderThan(150)
	require.Equal(t, 0, deleted)

	deleted = cache.DeleteOlderThan(201)
	require.Equal(t, 1, deleted)
}
