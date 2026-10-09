package v2

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink/v2/core/platform"
)

const enclaveLastFailureTimestamp = "enclave_execution_last_failure_timestamp_seconds"

// These metrics measure the enclave round-trip on the node, not in the enclave.
type confidentialModuleMetrics struct {
	executionDuration metric.Int64Histogram
	executionFailures metric.Int64Counter
	lastFailure       metric.Int64ObservableGauge
	meter             metric.Meter
	attrs             metric.MeasurementOption

	mu           sync.Mutex
	registration metric.Registration
	closed       bool
	failures     map[string]int64
}

func newConfidentialModuleMetrics(meter metric.Meter, workflowID, workflowOwner, workflowName string) (*confidentialModuleMetrics, error) {
	executionDuration, err := meter.Int64Histogram("enclave_execution_time_ms")
	if err != nil {
		return nil, err
	}
	executionFailures, err := meter.Int64Counter("enclave_execution_failures")
	if err != nil {
		return nil, err
	}
	lastFailure, err := meter.Int64ObservableGauge(enclaveLastFailureTimestamp,
		metric.WithUnit("s"),
		metric.WithDescription("Unix time of the last failed enclave round-trip; zero until a failure occurs"))
	if err != nil {
		return nil, err
	}
	return &confidentialModuleMetrics{
		executionDuration: executionDuration,
		executionFailures: executionFailures,
		lastFailure:       lastFailure,
		meter:             meter,
		attrs: metric.WithAttributes(
			attribute.String(platform.KeyWorkflowID, workflowID),
			attribute.String(platform.KeyWorkflowOwner, workflowOwner),
			attribute.String(platform.KeyWorkflowName, workflowName),
		),
		failures: map[string]int64{"system": 0, "user": 0},
	}, nil
}

func (m *confidentialModuleMetrics) start() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed || m.registration != nil {
		return nil
	}
	// Register only on Start: constructing an engine can fail before it owns the module.
	registration, err := m.meter.RegisterCallback(func(_ context.Context, observer metric.Observer) error {
		m.mu.Lock()
		defer m.mu.Unlock()
		for errorType, timestamp := range m.failures {
			observer.ObserveInt64(m.lastFailure, timestamp, m.attrs,
				metric.WithAttributes(attribute.String(errorTypeAttribute, errorType)))
		}
		return nil
	}, m.lastFailure)
	if err != nil {
		return err
	}
	m.registration = registration
	return nil
}

func (m *confidentialModuleMetrics) recordFailure(ctx context.Context, err error, at time.Time) {
	errorType := errorTypeFor(err)
	m.executionFailures.Add(ctx, 1, m.attrs,
		metric.WithAttributes(attribute.String(errorTypeAttribute, errorType)))
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.closed {
		// Concurrent completions must not replace a newer event with an older one.
		m.failures[errorType] = max(m.failures[errorType], at.Unix())
	}
}

func (m *confidentialModuleMetrics) close() error {
	m.mu.Lock()
	m.closed = true
	registration := m.registration
	m.registration = nil
	m.failures = nil
	m.mu.Unlock()
	// Unregister may wait for a collection using mu.
	if registration != nil {
		return registration.Unregister()
	}
	return nil
}
