package globalconfig

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

// Metric names and labels from the Offchain Capabilities Registry design (Observability).
const (
	MetricAppliedVersion   = "platform_cap_config_applied_version"
	MetricValidationErrors = "platform_cap_config_validation_errors_total"
	MetricApplyErrors      = "platform_cap_config_apply_errors_total"
	MetricAppliedHash      = "platform_cap_config_applied_hash"
	MetricAppliedTimestamp = "platform_cap_config_applied_timestamp_seconds"

	labelDomain = "domain"
	labelEnv    = "env"
)

// Metrics emits the node-side offchain capabilities registry metrics. A nil *Metrics is a
// no-op.
type Metrics struct {
	appliedVersion   metric.Int64Gauge
	validationErrors metric.Int64Counter
	applyErrors      metric.Int64Counter
	appliedHash      metric.Int64Gauge
	appliedTimestamp metric.Int64Gauge
}

// NewMetrics registers the metrics on meter.
func NewMetrics(meter metric.Meter) (*Metrics, error) {
	appliedVersion, err := meter.Int64Gauge(MetricAppliedVersion,
		metric.WithDescription("Currently applied OffchainCapabilitiesRegistry version (0 when none is applied)"))
	if err != nil {
		return nil, err
	}
	validationErrors, err := meter.Int64Counter(MetricValidationErrors,
		metric.WithDescription("Offchain capabilities registry payloads that failed validation"))
	if err != nil {
		return nil, err
	}
	applyErrors, err := meter.Int64Counter(MetricApplyErrors,
		metric.WithDescription("Valid offchain capabilities registry payloads that could not be applied (e.g. version not newer than the applied one)"))
	if err != nil {
		return nil, err
	}
	appliedHash, err := meter.Int64Gauge(MetricAppliedHash,
		metric.WithDescription("First 48 bits of the applied payload hash; equal values across nodes mean identical content"))
	if err != nil {
		return nil, err
	}
	appliedTimestamp, err := meter.Int64Gauge(MetricAppliedTimestamp,
		metric.WithDescription("Unix time at which the current offchain capabilities registry payload was applied"))
	if err != nil {
		return nil, err
	}
	return &Metrics{
		appliedVersion:   appliedVersion,
		validationErrors: validationErrors,
		applyErrors:      applyErrors,
		appliedHash:      appliedHash,
		appliedTimestamp: appliedTimestamp,
	}, nil
}

var defaultMetrics = sync.OnceValue(func() *Metrics {
	m, err := NewMetrics(beholder.GetMeter())
	if err != nil {
		return nil // metrics are best-effort; a nil *Metrics is a no-op
	}
	return m
})

// DefaultMetrics returns the process-wide metrics, registered on the Beholder meter.
func DefaultMetrics() *Metrics { return defaultMetrics() }

func labels(domain, env string) metric.MeasurementOption {
	return metric.WithAttributes(attribute.String(labelDomain, domain), attribute.String(labelEnv, env))
}

// RecordAppliedVersion records the applied version for (domain, env).
func (m *Metrics) RecordAppliedVersion(ctx context.Context, domain, env string, version uint64) {
	if m == nil {
		return
	}
	m.appliedVersion.Record(ctx, int64(min(version, math.MaxInt64)), labels(domain, env))
}

// RecordValidationError counts a payload rejected by validation (malformed, or not newer than
// the applied version).
func (m *Metrics) RecordValidationError(ctx context.Context, domain, env string) {
	if m == nil {
		return
	}
	m.validationErrors.Add(ctx, 1, labels(domain, env))
}

// RecordApplyError counts a payload that passed validation but could not be applied, e.g. a
// version that is not newer than the applied one.
func (m *Metrics) RecordApplyError(ctx context.Context, domain, env string) {
	if m == nil {
		return
	}
	m.applyErrors.Add(ctx, 1, labels(domain, env))
}

// RecordApplied records the content hash and the time a payload was applied for (domain, env).
// hash is the hex sha256 of the payload; only its first 48 bits are exported, so the value is
// exact in a float64 and nodes holding identical content report the same number.
func (m *Metrics) RecordApplied(ctx context.Context, domain, env, hash string, at time.Time) {
	if m == nil {
		return
	}
	m.appliedHash.Record(ctx, HashGaugeValue(hash), labels(domain, env))
	m.appliedTimestamp.Record(ctx, at.Unix(), labels(domain, env))
}

// HashGaugeValue maps a hex digest to the first 48 bits, or 0 if it is not valid hex.
func HashGaugeValue(hash string) int64 {
	if len(hash) < 12 {
		return 0
	}
	v, err := strconv.ParseInt(hash[:12], 16, 64)
	if err != nil {
		return 0
	}
	return v
}

// PayloadLabels extracts domain and env from a payload for metric labels. It is best-effort and
// works on payloads that fail validation; missing values are returned as "".
func PayloadLabels(raw string) (domain, env string) {
	var l struct {
		Domain string `json:"domain"`
		Env    string `json:"env"`
	}
	_ = json.Unmarshal([]byte(raw), &l)
	return l.Domain, l.Env
}
