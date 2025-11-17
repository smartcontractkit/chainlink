package capabilities

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

const (
	resultSuccess = "success"
	resultFailure = "failure"
	resultSkipped = "skipped"

	keyCapabilityID  = "capability_id"
	keyRemoteDONName = "remote_don_name"
)

type launcherMetrics struct {
	remoteAddedSuccess metric.Int64Gauge
	remoteAddedFailure metric.Int64Gauge
	remoteSkipped      metric.Int64Gauge

	localExposedSuccess metric.Int64Gauge
	localExposedFailure metric.Int64Gauge
	localSkipped        metric.Int64Gauge

	completedUpdates metric.Int64Counter
}

func resultToInt(result string, expectedResult string) int64 {
	if result == expectedResult {
		return 1
	}
	return 0
}

func (m *launcherMetrics) recordRemoteCapabilityAdded(ctx context.Context, capabilityID string, remoteDONName string, result string) {
	attrs := metric.WithAttributes(
		attribute.String(keyCapabilityID, capabilityID),
		attribute.String(keyRemoteDONName, remoteDONName),
	)
	m.remoteAddedSuccess.Record(ctx, resultToInt(result, resultSuccess), attrs)
	m.remoteAddedFailure.Record(ctx, resultToInt(result, resultFailure), attrs)
	m.remoteSkipped.Record(ctx, resultToInt(result, resultSkipped), attrs)
}

func (m *launcherMetrics) recordLocalCapabilityExposed(ctx context.Context, capabilityID string, result string) {
	attrs := metric.WithAttributes(
		attribute.String(keyCapabilityID, capabilityID),
	)
	m.localExposedSuccess.Record(ctx, resultToInt(result, resultSuccess), attrs)
	m.localExposedFailure.Record(ctx, resultToInt(result, resultFailure), attrs)
	m.localSkipped.Record(ctx, resultToInt(result, resultSkipped), attrs)
}

func (m *launcherMetrics) incrementCompletedUpdates(ctx context.Context) {
	attrs := metric.WithAttributes()
	m.completedUpdates.Add(ctx, 1, attrs)
}

func newLauncherMetrics() (*launcherMetrics, error) {
	remoteAddedSuccess, err := beholder.GetMeter().Int64Gauge("platform_launcher_remote_capability_added_success")
	if err != nil {
		return nil, err
	}

	remoteAddedFailure, err := beholder.GetMeter().Int64Gauge("platform_launcher_remote_capability_added_failure")
	if err != nil {
		return nil, err
	}

	remoteSkipped, err := beholder.GetMeter().Int64Gauge("platform_launcher_remote_capability_skipped")
	if err != nil {
		return nil, err
	}

	localExposedSuccess, err := beholder.GetMeter().Int64Gauge("platform_launcher_local_capability_exposed_success")
	if err != nil {
		return nil, err
	}

	localExposedFailure, err := beholder.GetMeter().Int64Gauge("platform_launcher_local_capability_exposed_failure")
	if err != nil {
		return nil, err
	}

	localSkipped, err := beholder.GetMeter().Int64Gauge("platform_launcher_local_capability_skipped")
	if err != nil {
		return nil, err
	}

	completedUpdates, err := beholder.GetMeter().Int64Counter("platform_launcher_completed_updates_total")
	if err != nil {
		return nil, err
	}

	return &launcherMetrics{
		remoteAddedSuccess:  remoteAddedSuccess,
		remoteAddedFailure:  remoteAddedFailure,
		remoteSkipped:       remoteSkipped,
		localExposedSuccess: localExposedSuccess,
		localExposedFailure: localExposedFailure,
		localSkipped:        localSkipped,
		completedUpdates:    completedUpdates,
	}, nil
}
