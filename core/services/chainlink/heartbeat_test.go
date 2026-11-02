package chainlink_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/metric/noop"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func TestNewHeartbeat_ConfiguresHeartbeatInterval(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		interval time.Duration
	}{
		{"default interval", 1 * time.Second},
		{"custom interval", 5 * time.Second},
		{"long interval", 1 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test logger
			lggr := logger.TestLogger(t)

			c := chainlink.HeartbeatConfig{
				Beat:  tt.interval,
				Lggr:  lggr,
				P2P:   "peer-id",
				AppID: "app-id",
			}
			// Create a new heartbeat service
			heartbeat := chainlink.NewHeartbeat(c)

			// Verify the heartbeat interval was set correctly
			assert.Equal(t, tt.interval, heartbeat.GetBeat())
		})
	}
}

func TestHeartbeat_MeterEvents(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}
	t.Parallel()

	// Inject a per-test recorder instead of swapping the global beholder client
	rec := beholdertest.NewRecorder()
	labels := map[string]string{
		"system":  "Application",
		"peer_id": "peer-id",
		"appID":   "app-id",
	}

	// Create custom metric counter to track gauge calls and set up heartbeat service
	var heartbeatCounter, heartbeatCountCounter int32
	mockMeter := newCountingMeter(t, &heartbeatCounter, &heartbeatCountCounter)
	c := chainlink.HeartbeatConfig{
		Beat:  50 * time.Millisecond,
		Lggr:  logger.TestLogger(t),
		P2P:   "peer-id",
		AppID: "app-id",
	}
	heartbeat := chainlink.NewHeartbeat(c,
		chainlink.WithMeter(mockMeter),
		chainlink.WithEmitter(custmsg.NewLabelerWithEmitter(rec).WithMapLabels(labels)),
	)
	require.NoError(t, heartbeat.Start(t.Context()))

	// Wait for ~10 heartbeats
	expectedCalls := 10
	time.Sleep(time.Duration(expectedCalls) * c.Beat)
	require.NoError(t, heartbeat.Close())

	// Assert both counts are in the expected range (8-12)
	hb := atomic.LoadInt32(&heartbeatCounter)
	hbCount := atomic.LoadInt32(&heartbeatCountCounter)
	assert.InDelta(t, expectedCalls, hb, 2, "Expected ~%d heartbeat gauge calls", expectedCalls)
	assert.InDelta(t, expectedCalls, hbCount, 2, "Expected ~%d heartbeat count gauge calls", expectedCalls)

	// Each tick emits one heartbeat custom message through the injected
	// emitter, so the message count tracks the gauge counts
	msgs, err := rec.BaseMessagesForLabels(t, labels)
	require.NoError(t, err)
	heartbeatMsgs := 0
	for _, m := range msgs {
		if m.Msg == "heartbeat" {
			heartbeatMsgs++
		}
	}
	assert.InDelta(t, expectedCalls, heartbeatMsgs, 2, "Expected ~%d heartbeat messages", expectedCalls)
}

// mockMeter is a custom implementation of metric.Meter that counts gauge creation specifically for heartbeat metrics.
type mockMeter struct {
	noop.Meter
	hb      *int32
	hbCount *int32
}

var _ metric.Meter = (*mockMeter)(nil)

func newCountingMeter(t *testing.T, heartbeatCount, heartbeatCountCount *int32) metric.Meter {
	return &mockMeter{
		hb:      heartbeatCount,
		hbCount: heartbeatCountCount,
	}
}

func (m *mockMeter) Int64Gauge(name string, options ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	// Return a counting gauge based on the name
	switch name {
	case "heartbeat":
		return &countingGauge{counter: m.hb}, nil
	case "heartbeat_count":
		return &countingGauge{counter: m.hbCount}, nil
	}
	// Default gauge
	return &countingGauge{}, nil
}

// countingGauge implements metric.Int64Gauge and counts how many times Record is called
type countingGauge struct {
	embedded.Int64Gauge
	counter *int32
}

func (g *countingGauge) Record(ctx context.Context, value int64, options ...metric.RecordOption) {
	if g.counter != nil {
		atomic.AddInt32(g.counter, 1)
	}
}

func (g *countingGauge) Enabled(ctx context.Context) bool {
	return true
}
