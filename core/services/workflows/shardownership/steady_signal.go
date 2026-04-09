package shardownership

import (
	"context"
	"sync/atomic"

	wfmon "github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

type SteadySignal struct {
	steady atomic.Bool
	seen   atomic.Bool

	metrics *wfmon.SteadySignalMetrics
}

type SteadySignalOption func(*SteadySignal)

func WithSteadySignalMetrics(m *wfmon.SteadySignalMetrics) SteadySignalOption {
	return func(s *SteadySignal) {
		s.metrics = m
	}
}

func NewSteadySignal(opts ...SteadySignalOption) *SteadySignal {
	s := &SteadySignal{}
	for _, o := range opts {
		o(s)
	}
	if s.metrics != nil {
		s.metrics.RecordInitial(context.Background())
	}
	return s
}

func (s *SteadySignal) ObserveRoutingSteady(steady bool) {
	if s == nil {
		return
	}
	s.seen.Store(true)
	s.steady.Store(steady)
	if s.metrics != nil {
		s.metrics.RecordObserved(context.Background(), steady)
	}
}

func (s *SteadySignal) Invalidate() {
	if s == nil {
		return
	}
	s.steady.Store(false)
	if s.metrics != nil {
		s.metrics.RecordInvalidate(context.Background())
	}
}

func (s *SteadySignal) SkipCommittedOwnerCheck() bool {
	if s == nil {
		return false
	}
	return s.seen.Load() && s.steady.Load()
}
