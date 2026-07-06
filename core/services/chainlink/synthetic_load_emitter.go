package chainlink

// THROWAWAY: synthetic workflow-event load generator.
//
// This service emits BaseMessages through the node's global beholder emitter at
// a fixed, hardcoded rate so we can load-test the downstream workflow service
// from the reliability staging DON. Because it goes through the standard
// beholder emitter, every message automatically carries the node's telemetry
// resource attributes (e.g. platformEnv=staging-reliability) as Kafka headers,
// which is exactly what the downstream IGNORE_DONS filter matches on.
//
// This whole file (and its wiring in application.go) is meant to be deleted
// before this branch is ever merged.

import (
	"context"
	"time"

	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/timeutil"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Tuning knobs for the throwaway load generator. Edit and redeploy to change
// the load; there is intentionally no TOML/env wiring.
const (
	// syntheticLoadEnabled turns the generator on. Keep false on any branch
	// that is not the dedicated reliability load-test deploy.
	syntheticLoadEnabled = true
	// syntheticLoadEventsPerSecond is the number of BaseMessages emitted every
	// second.
	syntheticLoadEventsPerSecond = 1000
)

type syntheticLoadEmitter struct {
	commonservices.Service
	eng *commonservices.Engine

	emitter         custmsg.MessageEmitter
	eventsPerSecond int
}

func newSyntheticLoadEmitter(lggr logger.Logger, eventsPerSecond int) *syntheticLoadEmitter {
	s := &syntheticLoadEmitter{
		emitter:         custmsg.NewLabeler().With("source", "synthetic-load"),
		eventsPerSecond: eventsPerSecond,
	}

	s.Service, s.eng = commonservices.Config{
		Name:  "SyntheticLoadEmitter",
		Start: s.start,
	}.NewServiceEngine(lggr)
	return s
}

func (s *syntheticLoadEmitter) start(_ context.Context) error {
	s.eng.Info("SyntheticLoad: starting")
	if s.eventsPerSecond <= 0 {
		s.eng.Warnw("synthetic load emitter started with non-positive rate; nothing will be emitted", "eventsPerSecond", s.eventsPerSecond)
		return nil
	}

	// Emit the full per-second budget once per tick. A one-second tick keeps the
	// rate predictable regardless of ticker precision; the burst is fine for a
	// load test.
	tick := func(ctx context.Context) {
		s.eng.Info("SyntheticLoad: emitting load on tick")
		for i := 0; i < s.eventsPerSecond; i++ {
			if err := s.emitter.Emit(ctx, "synthetic load test event"); err != nil {
				s.eng.Errorw("synthetic load emit failed", "err", err, "emitted", i)
				return
			}
		}
	}

	s.eng.GoTick(timeutil.NewTicker(func() time.Duration { return time.Second }), tick)
	return nil
}
