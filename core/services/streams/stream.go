package streams

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

type Runner interface {
	ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars) (run *pipeline.Run, trrs pipeline.TaskRunResults, err error)
	InitializePipeline(spec pipeline.Spec) (*pipeline.Pipeline, error)
}

type RunResultSaver interface {
	Save(run *pipeline.Run)
}

type Stream interface {
	Run(ctx context.Context) (*pipeline.Run, pipeline.TaskRunResults, error)
}

type stream struct {
	sync.RWMutex
	id     StreamID
	lggr   logger.Logger
	spec   *pipeline.Spec
	runner Runner
	rrs    RunResultSaver
}

func NewStream(lggr logger.Logger, id StreamID, spec pipeline.Spec, runner Runner, rrs RunResultSaver) Stream {
	return newStream(lggr, id, spec, runner, rrs)
}

func newStream(lggr logger.Logger, id StreamID, spec pipeline.Spec, runner Runner, rrs RunResultSaver) *stream {
	return &stream{sync.RWMutex{}, id, lggr.Named("Stream").With("streamID", id), &spec, runner, rrs}
}

func (s *stream) Run(ctx context.Context) (run *pipeline.Run, trrs pipeline.TaskRunResults, err error) {
	run, trrs, err = s.executeRun(ctx)

	if err != nil {
		return nil, nil, fmt.Errorf("Run failed: %w", err)
	}
	if s.rrs != nil {
		s.rrs.Save(run)
	}

	return
}

// The context passed in here has a timeout of (ObservationTimeout + ObservationGracePeriod).
// Upon context cancellation, its expected that we return any usable values within ObservationGracePeriod.
func (s *stream) executeRun(ctx context.Context) (*pipeline.Run, pipeline.TaskRunResults, error) {
	// the hot path here is to avoid parsing and use the pre-parsed, cached, pipeline
	s.RLock()
	initialize := s.spec.Pipeline == nil
	s.RUnlock()
	if initialize {
		pipeline, err := s.spec.ParsePipeline()
		if err != nil {
			return nil, nil, fmt.Errorf("Run failed due to unparseable pipeline: %w", err)
		}

		s.Lock()
		if s.spec.Pipeline == nil {
			s.spec.Pipeline = pipeline
			// initialize it for the given runner
			if _, err := s.runner.InitializePipeline(*s.spec); err != nil {
				return nil, nil, fmt.Errorf("Run failed due to error while initializing pipeline: %w", err)
			}
		}
		s.Unlock()
	}

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"pipelineSpec": map[string]interface{}{
			"id": s.spec.ID,
		},
		"stream": map[string]interface{}{
			"id": s.id,
		},
	})

	run, trrs, err := s.runner.ExecuteRun(ctx, *s.spec, vars)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing run for spec ID %v: %w", s.spec.ID, err)
	}

	return run, trrs, err
}
