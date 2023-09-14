package streams

import (
	"context"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-data-streams/streams"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Runner interface {
	ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger) (run *pipeline.Run, trrs pipeline.TaskRunResults, err error)
}

// TODO: Generalize to beyond simply an int
type DataPoint *big.Int

type Stream interface {
	Observe(ctx context.Context) (DataPoint, error)
}

type stream struct {
	id     streams.StreamID
	lggr   logger.Logger
	spec   pipeline.Spec
	runner Runner
}

func NewStream(lggr logger.Logger, id streams.StreamID, spec pipeline.Spec, runner Runner) Stream {
	return newStream(lggr, id, spec, runner)
}

func newStream(lggr logger.Logger, id streams.StreamID, spec pipeline.Spec, runner Runner) *stream {
	return &stream{id, lggr, spec, runner}
}

func (s *stream) Observe(ctx context.Context) (DataPoint, error) {
	var run *pipeline.Run
	run, trrs, err := s.executeRun(ctx)
	if err != nil {
		return nil, fmt.Errorf("Observe failed while executing run: %w", err)
	}
	s.lggr.Tracew("Observe executed run", "run", run)
	// FIXME: runResults??
	// select {
	// case s.runResults <- run:
	// default:
	//     s.lggr.Warnf("unable to enqueue run save for job ID %d, buffer full", s.spec.JobID)
	// }

	// NOTE: trrs comes back as _all_ tasks, but we only want the terminal ones
	// They are guaranteed to be sorted by index asc so should be in the correct order
	var finaltrrs []pipeline.TaskRunResult
	for _, trr := range trrs {
		if trr.IsTerminal() {
			finaltrrs = append(finaltrrs, trr)
		}
	}

	// FIXME: How to handle arbitrary-shaped inputs?
	// For now just assume everything is one *big.Int
	parsed, err := s.parse(finaltrrs)
	if err != nil {
		return nil, fmt.Errorf("Observe failed while parsing run results: %w", err)
	}
	return parsed, nil

}

// The context passed in here has a timeout of (ObservationTimeout + ObservationGracePeriod).
// Upon context cancellation, its expected that we return any usable values within ObservationGracePeriod.
func (s *stream) executeRun(ctx context.Context) (*pipeline.Run, pipeline.TaskRunResults, error) {
	// TODO: does it need some kind of debugging stuff here?
	vars := pipeline.NewVarsFrom(map[string]interface{}{})

	run, trrs, err := s.runner.ExecuteRun(ctx, s.spec, vars, s.lggr)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing run for spec ID %v: %w", s.spec.ID, err)
	}

	return run, trrs, err
}

// returns error on parse errors: if something is the wrong type
func (s *stream) parse(trrs pipeline.TaskRunResults) (*big.Int, error) {
	var finaltrrs []pipeline.TaskRunResult
	for _, trr := range trrs {
		// only return terminal trrs from executeRun
		if trr.IsTerminal() {
			finaltrrs = append(finaltrrs, trr)
		}
	}

	// pipeline.TaskRunResults comes ordered asc by index, this is guaranteed
	// by the pipeline executor
	if len(finaltrrs) != 1 {
		return nil, fmt.Errorf("invalid number of results, expected: 1, got: %d", len(finaltrrs))
	}
	res := finaltrrs[0].Result
	if res.Error != nil {
		return nil, res.Error
	} else if val, err := toBigInt(res.Value); err != nil {
		return nil, fmt.Errorf("failed to parse BenchmarkPrice: %w", err)
	} else {
		return val, nil
	}
}

func toBigInt(val interface{}) (*big.Int, error) {
	dec, err := utils.ToDecimal(val)
	if err != nil {
		return nil, err
	}
	return dec.BigInt(), nil
}
