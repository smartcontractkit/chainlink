package globalconfig

import (
	"context"
	"encoding/json"
	"math"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Metric names and labels from the Offchain Capabilities Registry design (Observability).
const (
	MetricAppliedVersion   = "platform_cap_config_applied_version"
	MetricValidationErrors = "platform_cap_config_validation_errors_total"

	labelDomain = "domain"
	labelEnv    = "env"
)

// Metrics emits the node-side offchain capabilities registry metrics. A nil *Metrics is a
// no-op.
type Metrics struct {
	appliedVersion   metric.Int64Gauge
	validationErrors metric.Int64Counter
}

// NewMetrics registers the metrics on meter. A nil *Metrics is a no-op.
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
	return &Metrics{appliedVersion: appliedVersion, validationErrors: validationErrors}, nil
}

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
