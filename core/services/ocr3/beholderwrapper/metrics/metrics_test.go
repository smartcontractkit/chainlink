package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewPluginMetrics(t *testing.T) {
	metrics, err := NewPluginMetrics("platform_ocr3_reporting_plugin", "test-plugin", "abc123")
	require.NoError(t, err)
	require.NotNil(t, metrics)
}

func TestRecordDuration(t *testing.T) {
	metrics, err := NewPluginMetrics("platform_ocr3_reporting_plugin", "test-plugin", "abc123")
	require.NoError(t, err)

	// Should not panic and should complete without error
	metrics.RecordDuration(context.Background(), Query, 100*time.Millisecond, true)
	metrics.RecordDuration(context.Background(), Observation, 200*time.Millisecond, false)
	metrics.RecordDuration(context.Background(), ValidateObservation, 50*time.Millisecond, true)
	metrics.RecordDuration(context.Background(), Outcome, 150*time.Millisecond, true)
	metrics.RecordDuration(context.Background(), Reports, 75*time.Millisecond, true)
	metrics.RecordDuration(context.Background(), ShouldAccept, 10*time.Millisecond, true)
	metrics.RecordDuration(context.Background(), ShouldTransmit, 5*time.Millisecond, false)
}

func TestTrackReports(t *testing.T) {
	metrics, err := NewPluginMetrics("platform_ocr3_reporting_plugin", "test-plugin", "abc123")
	require.NoError(t, err)

	// Should not panic and should complete without error
	metrics.TrackReports(context.Background(), Reports, 5, true)
	metrics.TrackReports(context.Background(), ShouldAccept, 1, true)
	metrics.TrackReports(context.Background(), ShouldTransmit, 0, false)
}

func TestTrackSize(t *testing.T) {
	metrics, err := NewPluginMetrics("platform_ocr3_reporting_plugin", "test-plugin", "abc123")
	require.NoError(t, err)

	// Should not panic and should complete without error
	metrics.TrackSize(context.Background(), Observation, 1024)
	metrics.TrackSize(context.Background(), Outcome, 2048)
}

func TestUpdateStatus(t *testing.T) {
	metrics, err := NewPluginMetrics("platform_ocr3_reporting_plugin", "test-plugin", "abc123")
	require.NoError(t, err)

	// Should not panic and should complete without error
	metrics.UpdateStatus(context.Background(), true)
	metrics.UpdateStatus(context.Background(), false)
}

func TestMetricViews(t *testing.T) {
	views := MetricViews("platform_ocr3_reporting_plugin")
	require.Len(t, views, 2)
}

func TestFunctionTypeConstants(t *testing.T) {
	// Verify all expected function types exist
	require.Equal(t, Query, FunctionType("query"))
	require.Equal(t, Observation, FunctionType("observation"))
	require.Equal(t, ValidateObservation, FunctionType("validateObservation"))
	require.Equal(t, Outcome, FunctionType("outcome"))
	require.Equal(t, ObservationQuorum, FunctionType("observationQuorum"))
	require.Equal(t, StateTransition, FunctionType("stateTransition"))
	require.Equal(t, Committed, FunctionType("committed"))
	require.Equal(t, Reports, FunctionType("reports"))
	require.Equal(t, ShouldAccept, FunctionType("shouldAccept"))
	require.Equal(t, ShouldTransmit, FunctionType("shouldTransmit"))
}
