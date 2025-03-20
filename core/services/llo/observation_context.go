package llo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/services/llo/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// ObservationContext ensures that each pipeline is only executed once. It is
// intended to be instantiated and used then discarded as part of one
// Observation cycle. Subsequent calls to Observe will return the same cached
// values.

var _ ObservationContext = (*observationContext)(nil)

type ObservationContext interface {
	Observe(ctx context.Context, streamID streams.StreamID, opts llo.DSOpts) (val llo.StreamValue, err error)
}

type execution struct {
	done <-chan struct{}

	run  *pipeline.Run
	trrs pipeline.TaskRunResults
	err  error
}

type observationContext struct {
	l logger.Logger
	r Registry
	t Telemeter

	executionsMu sync.Mutex
	// only execute each pipeline once
	executions map[streams.Pipeline]*execution
}

func NewObservationContext(l logger.Logger, r Registry, t Telemeter) ObservationContext {
	return newObservationContext(l, r, t)
}

func newObservationContext(l logger.Logger, r Registry, t Telemeter) *observationContext {
	return &observationContext{l, r, t, sync.Mutex{}, make(map[streams.Pipeline]*execution)}
}

func (oc *observationContext) Observe(ctx context.Context, streamID streams.StreamID, opts llo.DSOpts) (val llo.StreamValue, err error) {
	run, trrs, err := oc.run(ctx, streamID)
	observationFinishedAt := time.Now()
	if err != nil {
		// FIXME: This is a hack specific for V3 telemetry, future schemas should
		// use a generic stream value telemetry instead
		// https://smartcontract-it.atlassian.net/browse/MERC-6290
		oc.t.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, val, err)
		return nil, err
	}
	// Extract stream value based on streamID attribute
	found := false
	for _, trr := range trrs {
		if trr.Task.TaskStreamID() != nil && *trr.Task.TaskStreamID() == streamID {
			val, err = resultToStreamValue(trr.Result.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to convert result to StreamValue for streamID %d: %w", streamID, err)
			}
			if trr.FinishedAt.Valid {
				observationFinishedAt = trr.FinishedAt.Time
			}
			found = true
			break
		}
	}
	if !found {
		// If no streamID attribute is found in the task results, then assume the
		// final output is the stream ID and return that. This is safe to do since
		// the registry will never return a spec that doesn't match either by tag
		// or by spec streamID.
		val, err = extractFinalResultAsStreamValue(trrs)
		if oc.t.CaptureEATelemetry() {
			// FIXME: This is a hack specific for V3 telemetry, future schemas should
			// use the generic stream value telemetry instead
			// https://smartcontract-it.atlassian.net/browse/MERC-6290
			oc.t.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, val, err)
		}
	}
	if ch := GetObservationTelemetryCh(ctx); ch != nil {
		cd := opts.ConfigDigest()
		ot := &telem.LLOObservationTelemetry{
			StreamId:              streamID,
			ObservationTimestamp:  opts.ObservationTimestamp().UnixNano(),
			ObservationFinishedAt: observationFinishedAt.UnixNano(),
			SeqNr:                 opts.SeqNr(),
			ConfigDigest:          cd[:],
		}
		if err != nil {
			ot.ObservationError = new(string)
			*ot.ObservationError = err.Error()
		}
		if val != nil {
			ot.StreamValueType = int32(val.Type())
			b, err := val.MarshalBinary()
			if err != nil {
				oc.l.Errorw("failed to MarshalBinary on stream value", "error", err)
			} else {
				ot.StreamValueBinary = b
			}
			s, err := val.MarshalText()
			if err != nil {
				oc.l.Errorw("failed to MarshalText on stream value", "error", err)
			} else {
				ot.StreamValueText = string(s)
			}
		}
		select {
		case ch <- ot:
		default:
			oc.l.Error("telemetry channel is full, dropping observation telemetry")
		}
	}
	return
}

func resultToStreamValue(val interface{}) (llo.StreamValue, error) {
	switch v := val.(type) {
	case decimal.Decimal:
		return llo.ToDecimal(v), nil
	case float64:
		return llo.ToDecimal(decimal.NewFromFloat(v)), nil
	case int64:
		return llo.ToDecimal(decimal.NewFromInt(v)), nil
	case pipeline.ObjectParam:
		switch v.Type {
		case pipeline.DecimalType:
			return llo.ToDecimal(decimal.Decimal(v.DecimalValue)), nil
		default:
			return nil, fmt.Errorf("don't know how to convert pipeline.ObjectParam with type %d to llo.StreamValue", v.Type)
		}
	default:
		return nil, fmt.Errorf("don't know how to convert pipeline output result of type %T to llo.StreamValue (got: %v)", val, val)
	}
}

// extractFinalResultAsStreamValue extracts a final StreamValue from a TaskRunResults
func extractFinalResultAsStreamValue(trrs pipeline.TaskRunResults) (llo.StreamValue, error) {
	// pipeline.TaskRunResults comes ordered asc by index, this is guaranteed
	// by the pipeline executor
	finaltrrs := trrs.Terminals()

	// HACK: Right now we rely on the number of outputs to determine whether
	// its a Decimal or a Quote.
	// This is a hack to support the legacy "Quote" case.
	// Future stream specs should use streamID tags instead.
	switch len(finaltrrs) {
	case 1:
		res := finaltrrs[0].Result
		if res.Error != nil {
			return nil, fmt.Errorf("terminal task error: %w; all task errors: %w", res.Error, trrs.AllErrors())
		}
		val, err := toDecimal(res.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse BenchmarkPrice: %w", err)
		}
		return llo.ToDecimal(val), nil
	case 3:
		// Expect ordering of Benchmark, Bid, Ask
		results := make([]decimal.Decimal, 3)
		for i, trr := range finaltrrs {
			res := trr.Result
			if res.Error != nil {
				return nil, fmt.Errorf("failed to parse stream output into Quote (task index: %d): %w", i, res.Error)
			}
			val, err := toDecimal(res.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse decimal: %w", err)
			}
			results[i] = val
		}
		return &llo.Quote{
			Benchmark: results[0],
			Bid:       results[1],
			Ask:       results[2],
		}, nil
	default:
		return nil, fmt.Errorf("invalid number of results, expected: 1 or 3, got: %d", len(finaltrrs))
	}
}

func toDecimal(val interface{}) (decimal.Decimal, error) {
	return utils.ToDecimal(val)
}

type MissingStreamError struct {
	StreamID streams.StreamID
}

func (e MissingStreamError) Error() string {
	return fmt.Sprintf("no pipeline for stream: %d", e.StreamID)
}

func (oc *observationContext) run(ctx context.Context, streamID streams.StreamID) (*pipeline.Run, pipeline.TaskRunResults, error) {
	p, exists := oc.r.Get(streamID)
	if !exists {
		return nil, nil, MissingStreamError{StreamID: streamID}
	}

	// In case of multiple streamIDs per pipeline then the
	// first call executes and the others wait for result
	oc.executionsMu.Lock()
	ex, isExecuting := oc.executions[p]
	if isExecuting {
		oc.executionsMu.Unlock()
		// wait for it to finish
		select {
		case <-ex.done:
			return ex.run, ex.trrs, ex.err
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}

	// execute here
	ch := make(chan struct{})
	ex = &execution{done: ch}
	oc.executions[p] = ex
	oc.executionsMu.Unlock()

	run, trrs, err := p.Run(ctx)
	ex.run = run
	ex.trrs = trrs
	ex.err = err
	close(ch)

	return run, trrs, err
}
