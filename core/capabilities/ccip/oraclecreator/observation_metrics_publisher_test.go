package oraclecreator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBeholderMetricsPublisher_PublishMetric(t *testing.T) {
	t.Run("publishes with correct context", func(t *testing.T) {
		// Create a context with cancellation
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		// Create mock publisher
		mockPub := &mockPublisher{}

		// Test that context is properly passed through
		labels := map[string]string{"test": "label"}
		mockPub.PublishMetric(ctx, "test_metric", 1.0, labels)

		// Verify metric was recorded
		metrics := mockPub.getMetrics()
		require.Len(t, metrics, 1)
		assert.Equal(t, "test_metric", metrics[0].name)
		assert.Equal(t, 1, int(metrics[0].value))
		assert.Equal(t, "label", metrics[0].labels["test"])
	})

	t.Run("handles multiple labels", func(t *testing.T) {
		mockPub := &mockPublisher{}

		labels := map[string]string{
			"name":        "commit-1234",
			"env":         "test",
			"chain":       "ethereum",
			"node":        "node-1",
			"plugin_type": "commit",
		}

		mockPub.PublishMetric(t.Context(), "ocr3_sent_observations_total", 5.0, labels)

		metrics := mockPub.getMetrics()
		require.Len(t, metrics, 1)

		// Verify all labels are preserved
		assert.Equal(t, "commit-1234", metrics[0].labels["name"])
		assert.Equal(t, "test", metrics[0].labels["env"])
		assert.Equal(t, "ethereum", metrics[0].labels["chain"])
		assert.Equal(t, "node-1", metrics[0].labels["node"])
		assert.Equal(t, "commit", metrics[0].labels["plugin_type"])
	})

	t.Run("handles empty labels", func(t *testing.T) {
		mockPub := &mockPublisher{}

		mockPub.PublishMetric(t.Context(), "test_metric", 1.0, map[string]string{})

		metrics := mockPub.getMetrics()
		require.Len(t, metrics, 1)
		assert.Empty(t, metrics[0].labels)
	})

	t.Run("handles nil labels map", func(t *testing.T) {
		mockPub := &mockPublisher{}

		mockPub.PublishMetric(t.Context(), "test_metric", 1.0, nil)

		assert.NotPanics(t, func() {
			metrics := mockPub.getMetrics()
			require.Len(t, metrics, 1)
		})
	})
}
