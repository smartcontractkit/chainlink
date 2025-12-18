package chainlink_test

import (
	"bytes"
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/metric/noop"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
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
	lggr := logger.TestLogger(t)

	// Use a thread-safe byte collector
	collector := &byteCollector{}
	client, err := beholder.NewWriterClient(collector)
	require.NoError(t, err)

	// Set the global beholder client
	beholder.SetClient(client)
	defer beholder.SetClient(beholder.NewNoopClient()) // Reset to prevent test side effects

	// Create custom metric counter to track gauge calls and set up heartbeat service
	var heartbeatCounter, heartbeatCountCounter int32
	mockMeter := newCountingMeter(t, &heartbeatCounter, &heartbeatCountCounter)
	c := chainlink.HeartbeatConfig{
		Beat:  50 * time.Millisecond,
		Lggr:  lggr,
		P2P:   "peer-id",
		AppID: "app-id",
	}
	heartbeat := chainlink.NewHeartbeat(c, chainlink.WithMeter(mockMeter))
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

	// Check the output buffer for heartbeat messages
	outputStr := collector.String()
	assert.Contains(t, outputStr, "heartbeat", "Output should contain heartbeat messages")
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

// byteCollector collects all bytes written to it in a thread-safe manner
type byteCollector struct {
	mu     sync.Mutex
	buffer bytes.Buffer
}

func (bc *byteCollector) Write(p []byte) (n int, err error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	return bc.buffer.Write(p)
}

func (bc *byteCollector) String() string {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	return bc.buffer.String()
}

func TestExtractDomainFromURLString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		rawURL   string
		expected string
	}{
		{
			name:     "empty string",
			rawURL:   "",
			expected: "",
		},
		{
			name:     "simple https URL",
			rawURL:   "https://eth-mainnet.example.com",
			expected: "eth-mainnet.example.com",
		},
		{
			name:     "https URL with path",
			rawURL:   "https://eth-mainnet.infura.io/v3/abc123secretkey",
			expected: "eth-mainnet.infura.io",
		},
		{
			name:     "https URL with port",
			rawURL:   "https://localhost:8545",
			expected: "localhost",
		},
		{
			name:     "wss URL",
			rawURL:   "wss://eth-mainnet.example.com/ws",
			expected: "eth-mainnet.example.com",
		},
		{
			name:     "wss URL with API key in path",
			rawURL:   "wss://mainnet.infura.io/ws/v3/abc123secretkey",
			expected: "mainnet.infura.io",
		},
		{
			name:     "URL with query params (should be stripped)",
			rawURL:   "https://api.example.com/rpc?apikey=secret",
			expected: "api.example.com",
		},
		{
			name:     "URL with user info (should be stripped)",
			rawURL:   "https://user:pass@api.example.com/rpc",
			expected: "api.example.com",
		},
		{
			name:     "invalid URL",
			rawURL:   "://invalid",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chainlink.ExtractDomainFromURLString(tt.rawURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewHeartbeat_WithEVMNodeURLs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		evmNodeURLs map[string]string
	}{
		{
			name:        "no EVM node URLs",
			evmNodeURLs: nil,
		},
		{
			name: "single chain with both URLs",
			evmNodeURLs: map[string]string{
				"WSURL_1":   "eth-mainnet.example.com",
				"HTTPURL_1": "eth-mainnet.infura.io",
			},
		},
		{
			name: "multiple chains",
			evmNodeURLs: map[string]string{
				"WSURL_1":     "eth-mainnet.example.com",
				"HTTPURL_1":   "eth-mainnet.infura.io",
				"WSURL_137":   "polygon-mainnet.example.com",
				"HTTPURL_137": "polygon-mainnet.infura.io",
			},
		},
		{
			name: "chain with multiple nodes (comma-separated domains)",
			evmNodeURLs: map[string]string{
				"WSURL_1":   "node1.example.com,node2.example.com",
				"HTTPURL_1": "rpc1.infura.io,rpc2.alchemy.com",
			},
		},
		{
			name: "chain with only HTTPURL",
			evmNodeURLs: map[string]string{
				"HTTPURL_42161": "arb-mainnet.g.alchemy.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lggr := logger.TestLogger(t)

			c := chainlink.HeartbeatConfig{
				Beat:        1 * time.Second,
				Lggr:        lggr,
				P2P:         "peer-id",
				AppID:       "app-id",
				EVMNodeURLs: tt.evmNodeURLs,
			}

			heartbeat := chainlink.NewHeartbeat(c)

			// Verify heartbeat was created successfully
			assert.Equal(t, 1*time.Second, heartbeat.GetBeat())
		})
	}
}
