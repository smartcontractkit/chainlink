package monitoring

import (
	"context"
	"strconv"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

const (
	SteadySignalStateUnobserved int64 = iota
	SteadySignalStateTransition
	SteadySignalStateSteady
)

type SteadySignalMetrics struct {
	stateGauge        metric.Int64Gauge
	observeCounter    metric.Int64Counter
	invalidateCounter metric.Int64Counter
}

var (
	globalSteadySignalMetrics     *SteadySignalMetrics
	errGlobalSteadySignalMetrics  error
	globalSteadySignalMetricsOnce sync.Once
)

func GlobalSteadySignalMetrics() (*SteadySignalMetrics, error) {
	globalSteadySignalMetricsOnce.Do(func() {
		globalSteadySignalMetrics, errGlobalSteadySignalMetrics = newSteadySignalMetrics()
	})
	return globalSteadySignalMetrics, errGlobalSteadySignalMetrics
}

func newSteadySignalMetrics() (*SteadySignalMetrics, error) {
	m := &SteadySignalMetrics{}
	var err error
	m.stateGauge, err = beholder.GetMeter().Int64Gauge("platform_shard_routing_steady_signal_state")
	if err != nil {
		return nil, err
	}
	m.observeCounter, err = beholder.GetMeter().Int64Counter("platform_shard_routing_steady_signal_observe_total")
	if err != nil {
		return nil, err
	}
	m.invalidateCounter, err = beholder.GetMeter().Int64Counter("platform_shard_routing_steady_signal_invalidate_total")
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *SteadySignalMetrics) RecordInitial(ctx context.Context) {
	if m == nil {
		return
	}
	m.stateGauge.Record(ctx, SteadySignalStateUnobserved)
}

func (m *SteadySignalMetrics) RecordObserved(ctx context.Context, routingSteady bool) {
	if m == nil {
		return
	}
	var state int64
	if routingSteady {
		state = SteadySignalStateSteady
	} else {
		state = SteadySignalStateTransition
	}
	m.stateGauge.Record(ctx, state)
	steadyStr := strconv.FormatBool(routingSteady)
	m.observeCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("routing_steady", steadyStr)))
}

func (m *SteadySignalMetrics) RecordInvalidate(ctx context.Context) {
	if m == nil {
		return
	}
	m.invalidateCounter.Add(ctx, 1)
	m.stateGauge.Record(ctx, SteadySignalStateTransition)
}
