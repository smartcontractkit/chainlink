package bridgeconn

import (
	"context"
	"encoding/hex"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/bridgeconn/streamspb"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func testBridge(t *testing.T, name string) bridges.BridgeType {
	t.Helper()
	u, err := url.Parse("http://" + name + ".example.invalid:8080")
	require.NoError(t, err)
	return bridges.BridgeType{
		Name:                 bridges.BridgeName(name),
		URL:                  models.WebURL(*u),
		UseConnectionManager: true,
	}
}

func newTestManager() *bridgeConnManager {
	return &bridgeConnManager{
		cache: make(map[[32]byte]cacheEntry),
		conns: make(map[string]*eaConn),
		lggr:  logger.Nop(),
		dial: func(_ context.Context, _ string, _ bool) (eaStreamClient, error) {
			return nil, errStreamDialingDisabledForTest
		},
		clock: clockwork.NewRealClock(),
	}
}

func TestBridgeConnManager_GetOrCreateConn_OnePerBridge(t *testing.T) {
	t.Parallel()
	m := newTestManager()
	bridgeA := testBridge(t, "bridgea")
	bridgeB := testBridge(t, "bridgeb")

	var wg sync.WaitGroup
	conns := make([]*eaConn, 20)
	for idx := range 10 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			conns[idx] = m.getOrCreateConn(bridgeA.Name.String(), bridgeA.URL)
		}()
		go func() {
			defer wg.Done()
			conns[idx+10] = m.getOrCreateConn(bridgeB.Name.String(), bridgeB.URL)
		}()
	}
	wg.Wait()

	for i := 1; i < 10; i++ {
		assert.Same(t, conns[0], conns[i], "all concurrent lookups for bridgeA must return the same EAConn")
	}
	for i := 11; i < 20; i++ {
		assert.Same(t, conns[10], conns[i], "all concurrent lookups for bridgeB must return the same EAConn")
	}
	assert.NotSame(t, conns[0], conns[10], "separate bridges must get separate EAConns")

	m.connsMu.Lock()
	assert.Len(t, m.conns, 2)
	m.connsMu.Unlock()
}

func TestEAConn_RegisterAsset_RefreshAndIdlePrune(t *testing.T) {
	t.Parallel()
	m := newTestManager()
	bridge := testBridge(t, "idlebridge")
	clock := clockwork.NewFakeClock()
	m.clock = clock
	conn := newEAConn(bridge.Name.String(), bridge.URL, m)

	payload, err := structpb.NewStruct(map[string]any{"endpoint": "crypto"})
	require.NoError(t, err)

	var keyA, keyB [32]byte
	keyA[0] = 1
	keyB[0] = 2
	conn.registerAsset(keyA, payload)
	conn.registerAsset(keyB, payload)

	// Neither asset is idle yet: both survive a snapshot.
	active, req := conn.pruneAndSnapshot()
	assert.Len(t, active, 2)
	assert.Len(t, req.Subscriptions, 2)

	// Refresh keyA only, then advance past the idle timeout.
	clock.Advance(assetIdleTimeout / 2)
	conn.registerAsset(keyA, payload)
	clock.Advance(assetIdleTimeout/2 + time.Millisecond)

	active, req = conn.pruneAndSnapshot()
	assert.Len(t, active, 1)
	assert.Len(t, req.Subscriptions, 1)
	_, aStillActive := active[keyA]
	assert.True(t, aStillActive, "refreshed asset must survive the idle prune")
	_, bStillActive := active[keyB]
	assert.False(t, bStillActive, "unrefreshed asset must be pruned (indirect unsubscribe)")
}

func TestEAConn_HandleObservation(t *testing.T) {
	t.Parallel()
	m := newTestManager()
	bridge := testBridge(t, "obsbridge")
	conn := newEAConn(bridge.Name.String(), bridge.URL, m)

	payload, err := structpb.NewStruct(map[string]any{"endpoint": "crypto"})
	require.NoError(t, err)
	var key [32]byte
	key[0] = 42
	conn.registerAsset(key, payload)

	assetKey := hex.EncodeToString(key[:])
	observationsMetric := func() float64 {
		return testutil.ToFloat64(promEAConnObservationsTotal.WithLabelValues(conn.bridgeName, assetKey))
	}

	// Registered key: cached, and the per-asset counter is incremented.
	conn.handleObservation(&streamspb.SubscribeResponse{
		PayloadHash:     key[:],
		ObservationJson: []byte(`{"result":"1"}`),
	})
	m.mu.RLock()
	cached, ok := m.cache[key]
	m.mu.RUnlock()
	require.True(t, ok)
	assert.JSONEq(t, `{"result":"1"}`, string(cached.payload))
	assert.InEpsilon(t, float64(1), observationsMetric(), 0)

	// Unregistered key: discarded, counter unchanged.
	var unknownKey [32]byte
	unknownKey[0] = 99
	conn.handleObservation(&streamspb.SubscribeResponse{
		PayloadHash:     unknownKey[:],
		ObservationJson: []byte(`{"result":"2"}`),
	})
	m.mu.RLock()
	_, unknownCached := m.cache[unknownKey]
	m.mu.RUnlock()
	assert.False(t, unknownCached, "observation for an unregistered key must not be cached")
	assert.InEpsilon(t, float64(1), observationsMetric(), 0, "counter must not increment for an unregistered key")

	// Malformed (wrong-length) payload_hash: discarded, no panic, counter unchanged.
	conn.handleObservation(&streamspb.SubscribeResponse{
		PayloadHash:     []byte{1, 2, 3},
		ObservationJson: []byte(`{"result":"3"}`),
	})
	assert.InEpsilon(t, float64(1), observationsMetric(), 0, "counter must not increment for a malformed payload_hash")
}

func TestNextBackoff(t *testing.T) {
	t.Parallel()
	d := reconnectBackoffInitial
	seen := make([]time.Duration, 0, 11)
	seen = append(seen, d)
	for range 10 {
		d = nextBackoff(d)
		seen = append(seen, d)
	}
	for i := 1; i < len(seen); i++ {
		assert.LessOrEqual(t, seen[i-1], seen[i], "backoff must never decrease")
		assert.LessOrEqual(t, seen[i], reconnectBackoffMax, "backoff must never exceed the configured max")
	}
	assert.Equal(t, reconnectBackoffMax, seen[len(seen)-1], "backoff must saturate at the max")
}

func TestBridgeConnManager_GetObservation_CacheHitAndMiss(t *testing.T) {
	t.Parallel()
	m := newTestManager()
	bridge := testBridge(t, "cachebridge")
	requestData := map[string]any{"data": map[string]any{"endpoint": "crypto"}}

	_, err := m.GetObservation(bridge, requestData)
	require.ErrorIs(t, err, ErrBridgeObservationNotFound)

	key, err := bridgeObservationCacheKey(bridge.Name.String(), requestData["data"].(map[string]any))
	require.NoError(t, err)

	m.PutObservation(key, []byte(`{"result":"9700"}`))
	got, err := m.GetObservation(bridge, requestData)
	require.NoError(t, err)
	assert.JSONEq(t, `{"result":"9700"}`, string(got))

	// GetObservation must still have registered the asset with the bridge's EAConn.
	m.connsMu.Lock()
	conn, ok := m.conns[bridge.Name.String()]
	m.connsMu.Unlock()
	require.True(t, ok)
	conn.mu.Lock()
	_, registered := conn.assets[key]
	conn.mu.Unlock()
	assert.True(t, registered)
}

func TestBridgeConnManager_GetObservation_ExpiresAfterTTL(t *testing.T) {
	t.Parallel()
	m := newTestManager()
	clock := clockwork.NewFakeClock()
	m.clock = clock
	bridge := testBridge(t, "ttlbridge")
	requestData := map[string]any{"data": map[string]any{"endpoint": "crypto"}}

	key, err := bridgeObservationCacheKey(bridge.Name.String(), requestData["data"].(map[string]any))
	require.NoError(t, err)

	m.PutObservation(key, []byte(`{"result":"9700"}`))

	clock.Advance(observationTTL)
	got, err := m.GetObservation(bridge, requestData)
	require.NoError(t, err, "an entry at exactly the TTL boundary must not be treated as expired")
	assert.JSONEq(t, `{"result":"9700"}`, string(got))

	clock.Advance(time.Nanosecond)
	_, err = m.GetObservation(bridge, requestData)
	require.ErrorIs(t, err, ErrBridgeObservationExpired)
}

func TestBridgeConnManager_GetObservation_MissingDataField(t *testing.T) {
	t.Parallel()
	m := newTestManager()
	bridge := testBridge(t, "nodatabridge")

	_, err := m.GetObservation(bridge, map[string]any{"foo": "bar"})
	require.Error(t, err, "requestData without a \"data\" field must be rejected, not silently subscribed")

	_, err = m.GetObservation(bridge, map[string]any{"data": map[string]any{}})
	assert.Error(t, err, "an empty \"data\" field must be rejected")
}
