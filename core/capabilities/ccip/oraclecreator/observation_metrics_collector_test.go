package oraclecreator

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// mockPublisher captures published metrics for testing
type mockPublisher struct {
	mu      sync.Mutex
	metrics []metricRecord
}

type metricRecord struct {
	name   string
	value  float64
	labels map[string]string
}

func (m *mockPublisher) PublishMetric(_ context.Context, metricName string, value float64, labels map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics = append(m.metrics, metricRecord{
		name:   metricName,
		value:  value,
		labels: labels,
	})
}

func (m *mockPublisher) getMetrics() []metricRecord {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]metricRecord, len(m.metrics))
	copy(result, m.metrics)
	return result
}

func TestObservationMetricsCollector(t *testing.T) {
	lggr, err := logger.New()
	require.NoError(t, err)

	mockPub := &mockPublisher{}

	// Define constant labels that should be preserved (Prometheus labels)
	constantLabels := map[string]string{
		"name": "commit-1234",
		"env":  "test",
	}

	// Define Beholder labels with more details
	beholderLabels := map[string]string{
		"pluginType":    "commit",
		"chainId":       "1",
		"chainFamily":   "evm",
		"networkName":   "Ethereum",
		"chainSelector": "1234",
	}

	collector := NewObservationMetricsCollector(lggr, mockPub, constantLabels, beholderLabels)
	defer func() { _ = collector.Close() }()

	// Create a test registerer - we don't use WrapRegistererWith here anymore
	// as the collector handles the label wrapping internally
	registry := prometheus.NewRegistry()
	wrappedRegisterer := collector.CreateWrappedRegisterer(registry)

	// Simulate libocr registering the sent observations counter
	sentCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_sent_observations_total",
		Help: "Test counter for sent observations",
	})

	err = wrappedRegisterer.Register(sentCounter)
	require.NoError(t, err)

	// Simulate libocr registering the included observations counter
	includedCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_included_observations_total",
		Help: "Test counter for included observations",
	})

	err = wrappedRegisterer.Register(includedCounter)
	require.NoError(t, err)

	// Wait a moment for registration to complete
	time.Sleep(100 * time.Millisecond)

	// Increment the counters (simulating what libocr does)
	sentCounter.Inc()
	collector.sentObservationsCounter.readAndPublish()

	// Check that the metric was published with Beholder labels
	metrics := mockPub.getMetrics()
	assert.NotEmpty(t, metrics, "Expected at least one metric to be published")

	if len(metrics) > 0 {
		assert.Equal(t, "ocr3_sent_observations_total", metrics[0].name)
		assert.InEpsilon(t, 1.0, metrics[0].value, 0.01)
		// Verify Beholder labels are present
		assert.Equal(t, "commit", metrics[0].labels["pluginType"])
		assert.Equal(t, "1", metrics[0].labels["chainId"])
		assert.Equal(t, "evm", metrics[0].labels["chainFamily"])
		assert.Equal(t, "Ethereum", metrics[0].labels["networkName"])
		assert.Equal(t, "1234", metrics[0].labels["chainSelector"])
	}

	// Increment multiple times and trigger collections
	sentCounter.Inc()
	collector.sentObservationsCounter.readAndPublish()

	includedCounter.Inc()
	collector.includedObservationsCounter.readAndPublish()

	includedCounter.Inc()
	collector.includedObservationsCounter.readAndPublish()

	metrics = mockPub.getMetrics()
	assert.GreaterOrEqual(t, len(metrics), 3, "Expected at least 3 metrics to be published")

	// Verify we got the expected metrics with Beholder labels
	sentCount := 0
	includedCount := 0
	for _, m := range metrics {
		// Verify all metrics have the Beholder labels
		assert.Equal(t, "commit", m.labels["pluginType"], "Expected Beholder label 'pluginType' to be present")
		assert.Equal(t, "1", m.labels["chainId"], "Expected Beholder label 'chainId' to be present")
		assert.Equal(t, "evm", m.labels["chainFamily"], "Expected Beholder label 'chainFamily' to be present")
		assert.Equal(t, "Ethereum", m.labels["networkName"], "Expected Beholder label 'networkName' to be present")
		assert.Equal(t, "1234", m.labels["chainSelector"], "Expected Beholder label 'chainSelector' to be present")

		switch m.name {
		case "ocr3_sent_observations_total":
			sentCount++
		case "ocr3_included_observations_total":
			includedCount++
		}
	}

	assert.Equal(t, 2, sentCount, "Expected 2 sent observation metrics")
	assert.Equal(t, 2, includedCount, "Expected 2 included observation metrics")
}

func TestWrappedCounter(t *testing.T) {
	lggr, err := logger.New()
	require.NoError(t, err)

	mockPub := &mockPublisher{}

	baseCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "Test counter",
	})

	// Use Beholder-style labels
	allLabels := map[string]string{
		"pluginType":    "exec",
		"chainId":       "42161",
		"chainFamily":   "evm",
		"networkName":   "Arbitrum",
		"chainSelector": "42161",
	}

	wrapped := &wrappedCounter{
		Collector:  baseCounter,
		metricName: "test_counter",
		labels:     allLabels,
		publisher:  mockPub,
		logger:     lggr,
	}

	// Test Inc() - increment the base counter and trigger collection
	baseCounter.Inc()
	wrapped.readAndPublish()

	metrics := mockPub.getMetrics()
	require.Len(t, metrics, 1)
	assert.Equal(t, "test_counter", metrics[0].name)
	assert.InEpsilon(t, 1.0, metrics[0].value, 0.01) // Delta of 1

	// Verify Beholder labels are present
	assert.Equal(t, "exec", metrics[0].labels["pluginType"])
	assert.Equal(t, "42161", metrics[0].labels["chainId"])
	assert.Equal(t, "evm", metrics[0].labels["chainFamily"])
	assert.Equal(t, "Arbitrum", metrics[0].labels["networkName"])

	// Test Add() - increment by 5 and trigger collection
	baseCounter.Add(5.0)
	wrapped.readAndPublish()

	metrics = mockPub.getMetrics()
	require.Len(t, metrics, 2)
	assert.InEpsilon(t, 5.0, metrics[1].value, 0.01) // Delta of 5, not cumulative 6

	// Verify labels are still present in the second metric
	assert.Equal(t, "exec", metrics[1].labels["pluginType"])
	assert.Equal(t, "42161", metrics[1].labels["chainId"])
	assert.Equal(t, "evm", metrics[1].labels["chainFamily"])
	assert.Equal(t, "Arbitrum", metrics[1].labels["networkName"])
}

// TestWrappedCounter_ConcurrentIncrements verifies that the total delta is correctly published
// when the underlying counter is incremented concurrently from multiple goroutines.
func TestWrappedCounter_ConcurrentIncrements(t *testing.T) {
	lggr, err := logger.New()
	require.NoError(t, err)

	mockPub := &mockPublisher{}

	baseCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "concurrent_test_counter",
		Help: "Test counter for concurrent operations",
	})

	wrapped := &wrappedCounter{
		Collector:  baseCounter,
		metricName: "concurrent_test_counter",
		labels:     map[string]string{"test": "concurrent"},
		publisher:  mockPub,
		logger:     lggr,
	}

	// Simulate concurrent increments
	const numGoroutines = 10
	const incrementsPerGoroutine = 100
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < incrementsPerGoroutine; j++ {
				baseCounter.Inc()
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	wrapped.readAndPublish()

	// Verify that we published the total delta
	metrics := mockPub.getMetrics()
	require.Len(t, metrics, 1, "Expected one metric published with total delta")

	expectedTotal := float64(numGoroutines * incrementsPerGoroutine)
	assert.InEpsilon(t, expectedTotal, metrics[0].value, 0.01, "Should publish total delta")
}

// TestWrappedCounter_DeltaPublishing verifies that deltas (not cumulative values) are published
func TestWrappedCounter_DeltaPublishing(t *testing.T) {
	lggr, err := logger.New()
	require.NoError(t, err)

	mockPub := &mockPublisher{}

	baseCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "delta_test_counter",
		Help: "Test counter for delta publishing",
	})

	wrapped := &wrappedCounter{
		Collector:  baseCounter,
		metricName: "delta_test_counter",
		labels:     map[string]string{"test": "delta"},
		publisher:  mockPub,
		logger:     lggr,
	}

	collectMetrics := func() { wrapped.readAndPublish() }

	// Test sequence: Inc(), Inc(), Add(5), Inc(), Add(10)
	baseCounter.Inc() // Should publish 1
	collectMetrics()

	baseCounter.Inc() // Should publish 1 (not 2)
	collectMetrics()

	baseCounter.Add(5.0) // Should publish 5 (not 7)
	collectMetrics()

	baseCounter.Inc() // Should publish 1 (not 8)
	collectMetrics()

	baseCounter.Add(10.0) // Should publish 10 (not 18)
	collectMetrics()

	metrics := mockPub.getMetrics()
	require.Len(t, metrics, 5)

	// Verify each published value is the delta, not cumulative
	expectedDeltas := []float64{1, 1, 5, 1, 10}
	for i, expected := range expectedDeltas {
		assert.InEpsilon(t, expected, metrics[i].value, 0.01,
			"Metric %d should publish delta %f, not cumulative value", i, expected)
	}

	// If we were to sum the deltas, we should get the cumulative value
	var sum int
	for _, m := range metrics {
		sum += int(m.value)
	}
	assert.Equal(t, 18, sum, "Sum of deltas should equal cumulative value")
}

// TestWrappedCounter_AddWithFractionalValues tests Add with non-integer values
func TestWrappedCounter_AddWithFractionalValues(t *testing.T) {
	lggr, err := logger.New()
	require.NoError(t, err)

	mockPub := &mockPublisher{}

	baseCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "fractional_test_counter",
		Help: "Test counter for fractional values",
	})

	wrapped := &wrappedCounter{
		Collector:  baseCounter,
		metricName: "fractional_test_counter",
		labels:     map[string]string{"test": "fractional"},
		publisher:  mockPub,
		logger:     lggr,
	}

	collectMetrics := func() { wrapped.readAndPublish() }

	// Test with fractional values
	baseCounter.Add(2.7)
	collectMetrics()

	baseCounter.Add(1.3)
	collectMetrics()

	metrics := mockPub.getMetrics()
	require.Len(t, metrics, 2)

	// Verify the published deltas are correct (fractional values)
	assert.InEpsilon(t, 2.7, metrics[0].value, 0.01)
	assert.InEpsilon(t, 1.3, metrics[1].value, 0.01)
}

// TestObservationMetricsCollector_NonTargetMetrics verifies non-observation metrics pass through unchanged
func TestObservationMetricsCollector_NonTargetMetrics(t *testing.T) {
	lggr, err := logger.New()
	require.NoError(t, err)

	mockPub := &mockPublisher{}
	prometheusLabels := map[string]string{"name": "test"}
	beholderLabels := map[string]string{
		"pluginType":  "commit",
		"chainId":     "1",
		"chainFamily": "evm",
		"networkName": "Ethereum",
	}

	collector := NewObservationMetricsCollector(lggr, mockPub, prometheusLabels, beholderLabels)
	defer func() { _ = collector.Close() }()

	registry := prometheus.NewRegistry()
	wrappedRegisterer := collector.CreateWrappedRegisterer(registry)

	// Register a counter that should NOT be wrapped
	otherCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "other_metric_total",
		Help: "Some other metric",
	})

	err = wrappedRegisterer.Register(otherCounter)
	require.NoError(t, err)

	// Increment the non-target counter
	otherCounter.Inc()
	time.Sleep(50 * time.Millisecond)

	// Should not be published to our mock publisher
	metrics := mockPub.getMetrics()
	assert.Empty(t, metrics, "Non-target metrics should not be published")

	// Verify the collector doesn't have it
	assert.Nil(t, collector.sentObservationsCounter)
	assert.Nil(t, collector.includedObservationsCounter)
}

// TestObservationMetricsCollector_BackgroundPolling verifies that Start publishes metrics
// on a timer without any explicit Collect call from the outside (i.e. without a Prometheus scrape).
func TestObservationMetricsCollector_BackgroundPolling(t *testing.T) {
	lggr, err := logger.New()
	require.NoError(t, err)

	mockPub := &mockPublisher{}
	collector := NewObservationMetricsCollector(lggr, mockPub,
		map[string]string{"name": "commit-1234"},
		map[string]string{"pluginType": "commit"},
	)
	defer func() { _ = collector.Close() }()

	registry := prometheus.NewRegistry()
	wrappedRegisterer := collector.CreateWrappedRegisterer(registry)

	sentCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_sent_observations_total",
		Help: "Test counter",
	})
	require.NoError(t, wrappedRegisterer.Register(sentCounter))

	sentCounter.Add(3)

	// Start with a short interval — no external Collect call is made.
	const pollInterval = 50 * time.Millisecond
	collector.Start(pollInterval)

	// Wait for at least one poll to fire.
	require.Eventually(t, func() bool {
		return len(mockPub.getMetrics()) > 0
	}, 500*time.Millisecond, 10*time.Millisecond)

	metrics := mockPub.getMetrics()
	require.Len(t, metrics, 1)
	assert.Equal(t, "ocr3_sent_observations_total", metrics[0].name)
	assert.InEpsilon(t, 3.0, metrics[0].value, 0.01)

	// Subsequent polls with no new increments should not publish.
	time.Sleep(3 * pollInterval)
	assert.Len(t, mockPub.getMetrics(), 1, "no new publishes expected when counter has not changed")

	// A new increment should be picked up on the next poll.
	sentCounter.Inc()
	require.Eventually(t, func() bool {
		return len(mockPub.getMetrics()) >= 2
	}, 500*time.Millisecond, 10*time.Millisecond)

	assert.InEpsilon(t, 1.0, mockPub.getMetrics()[1].value, 0.01)
}

// TestObservationMetricsCollector_CloseStopsPolling verifies that Close stops the background
// goroutine and no further publishes occur even if the counter keeps incrementing.
func TestObservationMetricsCollector_CloseStopsPolling(t *testing.T) {
	lggr, err := logger.New()
	require.NoError(t, err)

	mockPub := &mockPublisher{}
	collector := NewObservationMetricsCollector(lggr, mockPub,
		map[string]string{"name": "commit-1234"},
		map[string]string{"pluginType": "commit"},
	)

	registry := prometheus.NewRegistry()
	wrappedRegisterer := collector.CreateWrappedRegisterer(registry)

	sentCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ocr3_sent_observations_total",
		Help: "Test counter",
	})
	require.NoError(t, wrappedRegisterer.Register(sentCounter))

	sentCounter.Inc()

	const pollInterval = 50 * time.Millisecond
	collector.Start(pollInterval)

	// Wait for the first publish.
	require.Eventually(t, func() bool {
		return len(mockPub.getMetrics()) > 0
	}, 500*time.Millisecond, 10*time.Millisecond)

	require.NoError(t, collector.Close())

	// Increment again after Close — the poller should no longer be running.
	sentCounter.Inc()
	time.Sleep(3 * pollInterval)

	assert.Len(t, mockPub.getMetrics(), 1, "no new publishes expected after Close")
}

// TestObservationMetricsCollector_Close verifies proper cleanup
func TestObservationMetricsCollector_Close(t *testing.T) {
	lggr, err := logger.New()
	require.NoError(t, err)

	mockPub := &mockPublisher{}
	prometheusLabels := map[string]string{"name": "test"}
	beholderLabels := map[string]string{
		"pluginType":  "commit",
		"chainId":     "1",
		"chainFamily": "evm",
		"networkName": "Ethereum",
	}

	collector := NewObservationMetricsCollector(lggr, mockPub, prometheusLabels, beholderLabels)

	// Close the collector - should not error
	err = collector.Close()
	require.NoError(t, err)

	// Verify we can call Close multiple times without panic
	err = collector.Close()
	require.NoError(t, err)
}

// TestWrappedCounter_NilPublisher verifies behavior when publisher is nil
func TestWrappedCounter_NilPublisher(t *testing.T) {
	lggr, err := logger.New()
	require.NoError(t, err)

	baseCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "nil_publisher_test",
		Help: "Test counter with nil publisher",
	})

	wrapped := &wrappedCounter{
		Collector:  baseCounter,
		metricName: "nil_publisher_test",
		labels:     map[string]string{"test": "nil"},
		publisher:  nil, // Explicitly nil
		logger:     lggr,
	}

	// Should not panic
	require.NotPanics(t, func() {
		baseCounter.Inc()
		baseCounter.Add(5.0)
		wrapped.readAndPublish()
	})
}
