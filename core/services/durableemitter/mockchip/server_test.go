package mockchip_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	chipbatch "github.com/smartcontractkit/chainlink-common/pkg/chipingress/batch"
	"github.com/smartcontractkit/chainlink-common/pkg/durableemitter"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"

	"github.com/smartcontractkit/chainlink/v2/core/services/durableemitter/mockchip"
)

// memStore is a tiny in-memory DurableEventStore for these tests. The
// canonical in-memory store in chainlink-common lives under _test.go and is
// not importable, so we duplicate the minimum surface DurableEmitter calls.
type memStore struct {
	mu     sync.Mutex
	events map[int64]*durableemitter.DurableEvent
	nextID atomic.Int64
}

func newMemStore() *memStore {
	return &memStore{events: make(map[int64]*durableemitter.DurableEvent)}
}

func (m *memStore) Insert(_ context.Context, payload []byte) (int64, error) {
	id := m.nextID.Add(1)
	m.mu.Lock()
	m.events[id] = &durableemitter.DurableEvent{
		ID: id, Payload: append([]byte(nil), payload...), CreatedAt: time.Now(),
	}
	m.mu.Unlock()
	return id, nil
}

func (m *memStore) Delete(_ context.Context, id int64) error {
	m.mu.Lock()
	delete(m.events, id)
	m.mu.Unlock()
	return nil
}

func (m *memStore) MarkDelivered(ctx context.Context, id int64) error {
	return m.Delete(ctx, id)
}

func (m *memStore) MarkDeliveredBatch(_ context.Context, ids []int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var n int64
	for _, id := range ids {
		if _, ok := m.events[id]; ok {
			delete(m.events, id)
			n++
		}
	}
	return n, nil
}

func (m *memStore) PurgeDelivered(_ context.Context, _ int) (int64, error) { return 0, nil }

func (m *memStore) ListPending(_ context.Context, createdBefore time.Time, limit int) ([]durableemitter.DurableEvent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []durableemitter.DurableEvent
	for _, e := range m.events {
		if e.CreatedAt.Before(createdBefore) {
			out = append(out, *e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	if len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (m *memStore) DeleteExpired(_ context.Context, ttl time.Duration) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cutoff := time.Now().Add(-ttl)
	var n int64
	for id, e := range m.events {
		if e.CreatedAt.Before(cutoff) {
			delete(m.events, id)
			n++
		}
	}
	return n, nil
}

// startMockServer returns a fully started mock + cleanup closure.
func startMockServer(t *testing.T) *mockchip.Server {
	t.Helper()
	srv := mockchip.NewServer()
	addr, err := srv.Start("127.0.0.1:0")
	require.NoError(t, err)
	require.NotEmpty(t, addr)
	t.Cleanup(srv.Stop)
	return srv
}

// newDurableEmitterAgainstMock builds a DurableEmitter wired to the mock
// server's gRPC endpoint with retransmit enabled and a short retransmit
// cadence so the drain-after-restore assertion does not have to wait long.
func newDurableEmitterAgainstMock(t *testing.T, srv *mockchip.Server) *durableemitter.DurableEmitter {
	t.Helper()

	chipClient, err := chipingress.NewClient(srv.Addr(), chipingress.WithInsecureConnection())
	require.NoError(t, err)

	batchClient, err := chipbatch.NewBatchClient(chipClient,
		chipbatch.WithBatchSize(10),
		chipbatch.WithBatchInterval(20*time.Millisecond),
		chipbatch.WithMaxConcurrentSends(2),
		chipbatch.WithMaxPublishTimeout(2*time.Second),
		chipbatch.WithShutdownTimeout(2*time.Second),
	)
	require.NoError(t, err)

	cfg := durableemitter.DefaultConfig()
	cfg.RetransmitInterval = 100 * time.Millisecond
	cfg.RetransmitAfter = 50 * time.Millisecond
	cfg.PurgeInterval = 100 * time.Millisecond
	cfg.PublishTimeout = 2 * time.Second
	cfg.DisablePruning = false

	em, err := durableemitter.NewDurableEmitter(
		newMemStore(),
		batchClient,
		nil, // no single-event fallback — retransmit is the path we want to exercise
		true,
		cfg,
		logger.Test(t),
		nil,
	)
	require.NoError(t, err)
	return em
}

func TestMockServer_CapturesPublishedEvents(t *testing.T) {
	srv := startMockServer(t)

	em := newDurableEmitterAgainstMock(t, srv)
	servicetest.Run(t, em)

	ctx := t.Context()
	require.NoError(t, em.Emit(ctx, []byte(`{"hello":"world"}`),
		"source", "mockchip-test", "type", "test.event", "subject", "sub-1"))

	require.NoError(t, srv.WaitFor(testCtx(t, 5*time.Second), 1))

	got := srv.Captured()
	require.Len(t, got, 1)
	assert.Equal(t, "mockchip-test", got[0].Event.GetSource())
	assert.Equal(t, "test.event", got[0].Event.GetType())
	assert.NotEmpty(t, got[0].Event.GetId())
	assert.False(t, got[0].ReceivedAt.IsZero())

	stats := srv.Stats()
	assert.Equal(t, 1, stats.Captured)
	assert.GreaterOrEqual(t, stats.BatchCalls+stats.PublishCalls, int64(1))
	assert.Equal(t, int64(0), stats.FailedCalls)
	assert.False(t, stats.OutageActive)
}

func TestMockServer_OutageRejectsRPCs(t *testing.T) {
	srv := startMockServer(t)
	srv.SetOutage(true)
	require.True(t, srv.OutageActive())

	client, err := chipingress.NewClient(srv.Addr(), chipingress.WithInsecureConnection())
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	event, err := chipingress.NewEvent("mockchip-test", "test.event", []byte("payload"), nil)
	require.NoError(t, err)
	eventPb, err := chipingress.EventToProto(event)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	_, err = client.Publish(ctx, eventPb)
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok, "expected gRPC status error, got %T: %v", err, err)
	assert.Equal(t, codes.Unavailable, st.Code())

	stats := srv.Stats()
	assert.Equal(t, 0, stats.Captured)
	assert.GreaterOrEqual(t, stats.FailedCalls, int64(1))
	assert.Equal(t, int64(1), stats.OutageDurations)
}

func TestMockServer_DrainsAfterOutageRestored(t *testing.T) {
	srv := startMockServer(t)

	em := newDurableEmitterAgainstMock(t, srv)
	servicetest.Run(t, em)

	srv.SetOutage(true)

	const sent = 25
	ctx := t.Context()
	for i := range sent {
		require.NoError(t, em.Emit(ctx, fmt.Appendf(nil, "payload-%d", i),
			"source", "mockchip-test", "type", "test.event"))
	}

	// While outage is active, the mock must reject every publish attempt and
	// DurableEmitter must hold the events in its pending queue.
	require.Eventually(t, func() bool {
		return srv.Stats().FailedCalls > 0
	}, 5*time.Second, 25*time.Millisecond, "expected outage to produce failed RPCs")
	assert.Zero(t, srv.CapturedCount(), "no events should be captured during outage")
	assert.Equal(t, int64(sent), em.PendingDepth(), "events must remain pending in DurableEmitter")

	srv.SetOutage(false)

	// Retransmit loop should drain the queue into the mock once the outage
	// is lifted. Allow up to 10s — generous for a 100ms retransmit tick.
	require.NoError(t, srv.WaitFor(testCtx(t, 10*time.Second), sent))

	// And the durable store should drain to zero pending soon after.
	require.Eventually(t, func() bool {
		return em.PendingDepth() == 0
	}, 10*time.Second, 50*time.Millisecond, "DurableEmitter pending depth should fall to 0 after drain")

	// All emitted payloads should be present, no duplicates required (the
	// retransmit loop may re-deliver, so we just assert >= sent).
	assert.GreaterOrEqual(t, srv.CapturedCount(), sent)
}

func TestHTTPController_OutageToggleAndStats(t *testing.T) {
	srv := startMockServer(t)

	ctrl := mockchip.NewHTTPController(srv)
	httpAddr, err := ctrl.Start("127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = ctrl.Stop(shutdownCtx)
	})

	base := "http://" + httpAddr
	ctx := t.Context()

	// Healthz
	resp := doRequest(t, ctx, http.MethodGet, base+"/healthz")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	// Outage on -> stats reflects it
	resp = doRequest(t, ctx, http.MethodPost, base+"/outage/on")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	resp = doRequest(t, ctx, http.MethodGet, base+"/stats")
	var stats mockchip.Stats
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&stats))
	_ = resp.Body.Close()
	assert.True(t, stats.OutageActive)

	// Outage off -> stats reflects it
	resp = doRequest(t, ctx, http.MethodPost, base+"/outage/off")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	resp = doRequest(t, ctx, http.MethodGet, base+"/stats")
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&stats))
	_ = resp.Body.Close()
	assert.False(t, stats.OutageActive)

	// Capture a real event then list it via /events
	em := newDurableEmitterAgainstMock(t, srv)
	servicetest.Run(t, em)
	require.NoError(t, em.Emit(t.Context(), []byte("via-http"),
		"source", "mockchip-test", "type", "test.event"))
	require.NoError(t, srv.WaitFor(testCtx(t, 5*time.Second), 1))

	resp = doRequest(t, ctx, http.MethodGet, base+"/events")
	var events []mockchip.EventSummary
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&events))
	_ = resp.Body.Close()
	require.NotEmpty(t, events)
	assert.Equal(t, "mockchip-test", events[0].Source)
	assert.Equal(t, "test.event", events[0].Type)

	// Reset clears captured events
	resp = doRequest(t, ctx, http.MethodPost, base+"/reset")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()
	assert.Equal(t, 0, srv.CapturedCount())
}

func doRequest(t *testing.T, ctx context.Context, method, url string) *http.Response {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

// testCtx returns a context that's cancelled at test cleanup or after timeout,
// whichever comes first.
func testCtx(t *testing.T, timeout time.Duration) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	t.Cleanup(cancel)
	return ctx
}
