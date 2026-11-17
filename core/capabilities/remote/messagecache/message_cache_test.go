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
	t.Parallel()

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

	// ready with two messages — Ready does not commit delivery state;
	// it returns true every time the threshold is met until MarkDelivered.
	ready, messages := cache.Ready(eventID1, 2, 100, true)
	require.True(t, ready)
	require.Equal(t, []byte(payloadA), messages[0])
	require.Equal(t, []byte(payloadA), messages[1])

	// still ready before MarkDelivered (failed aggregation can retry)
	ready, _ = cache.Ready(eventID1, 2, 100, true)
	require.True(t, ready)
	require.False(t, cache.WasReady(eventID1))

	// after MarkDelivered, once=true returns false and WasReady reports true
	cache.MarkDelivered(eventID1)
	ready, _ = cache.Ready(eventID1, 2, 100, true)
	require.False(t, ready)
	require.True(t, cache.WasReady(eventID1))
}

func TestMessageCache_Peers(t *testing.T) {
	t.Parallel()

	cache := messagecache.NewMessageCache[string, string]()
	require.Nil(t, cache.Peers(eventID1))

	cache.Insert(eventID1, peerID1, 100, []byte(payloadA))
	cache.Insert(eventID1, peerID1, 110, []byte(payloadA)) // re-insert from the same peer
	cache.Insert(eventID1, peerID2, 120, []byte(payloadA))
	cache.Insert(eventID2, peerID1, 130, []byte(payloadA))
	require.Equal(t, map[string]bool{peerID1: true, peerID2: true}, cache.Peers(eventID1))

	// the result is a snapshot
	peers := cache.Peers(eventID1)
	delete(peers, peerID1)
	require.Len(t, cache.Peers(eventID1), 2)

	cache.MarkDelivered(eventID1)
	cache.Delete(eventID1)
	require.Nil(t, cache.Peers(eventID1))
	require.False(t, cache.WasReady(eventID1))
	require.Len(t, cache.Peers(eventID2), 1)
}

func TestMessageCache_DeleteOlderThan(t *testing.T) {
	t.Parallel()

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
