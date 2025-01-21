package streams

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

type Runner interface {
	ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars) (run *pipeline.Run, trrs pipeline.TaskRunResults, err error)
	InitializePipeline(spec pipeline.Spec) (*pipeline.Pipeline, error)
}

type RunResultSaver interface {
	Save(run *pipeline.Run)
}

type Pipeline interface {
	Run(ctx context.Context) (*pipeline.Run, pipeline.TaskRunResults, error)
	StreamIDs() []StreamID
}

type multiStreamPipeline struct {
	lggr      logger.Logger
	spec      pipeline.Spec
	runner    Runner
	rrs       RunResultSaver
	streamIDs []StreamID
	newVars   func() pipeline.Vars
}

func NewMultiStreamPipeline(lggr logger.Logger, jb job.Job, runner Runner, rrs RunResultSaver) (Pipeline, error) {
	return newMultiStreamPipeline(lggr, jb, runner, rrs)
}

func newMultiStreamPipeline(lggr logger.Logger, jb job.Job, runner Runner, rrs RunResultSaver) (*multiStreamPipeline, error) {
	if jb.PipelineSpec == nil {
		// should never happen
		return nil, errors.New("job has no pipeline spec")
	}
	spec := *jb.PipelineSpec
	spec.JobID = jb.ID
	spec.JobName = jb.Name.ValueOrZero()
	spec.JobType = string(jb.Type)
	if spec.Pipeline == nil {
		pipeline, err := spec.ParsePipeline()
		if err != nil {
			return nil, fmt.Errorf("unparseable pipeline: %w", err)
		}

		spec.Pipeline = pipeline
		// initialize it for the given runner
		if _, err := runner.InitializePipeline(spec); err != nil {
			return nil, fmt.Errorf("error while initializing pipeline: %w", err)
		}
	}
	var streamIDs []StreamID
	for _, t := range spec.Pipeline.Tasks {
		if t.TaskStreamID() != nil {
			streamIDs = append(streamIDs, *t.TaskStreamID())
		}
	}
	if jb.StreamID != nil {
		streamIDs = append(streamIDs, *jb.StreamID)
	}
	if err := validateStreamIDs(streamIDs); err != nil {
		return nil, fmt.Errorf("invalid stream IDs: %w", err)
	}
	vars := func() pipeline.Vars {
		return pipeline.NewVarsFrom(map[string]interface{}{
			"pipelineSpec": map[string]interface{}{
				"id": jb.PipelineSpecID,
			},
			"jb": map[string]interface{}{
				"databaseID":    jb.ID,
				"externalJobID": jb.ExternalJobID,
				"name":          jb.Name.ValueOrZero(),
			},
		})
	}

	return &multiStreamPipeline{
		lggr.Named("MultiStreamPipeline").With("spec.ID", spec.ID, "jobID", spec.JobID, "jobName", spec.JobName, "jobType", spec.JobType),
		spec,
		runner,
		rrs,
		streamIDs,
		vars}, nil
}

func validateStreamIDs(streamIDs []StreamID) error {
	seen := make(map[StreamID]struct{})
	for _, id := range streamIDs {
		if _, ok := seen[id]; ok {
			return fmt.Errorf("duplicate stream ID: %v", id)
		}
		seen[id] = struct{}{}
	}
	return nil
}

func (s *multiStreamPipeline) Run(ctx context.Context) (run *pipeline.Run, trrs pipeline.TaskRunResults, err error) {
	run, trrs, err = s.executeRun(ctx)

	if err != nil {
		return nil, nil, fmt.Errorf("Run failed: %w", err)
	}
	if s.rrs != nil {
		s.rrs.Save(run)
	}

	return
}

func (s *multiStreamPipeline) StreamIDs() []StreamID {
	return s.streamIDs
}

// The context passed in here has a timeout of (ObservationTimeout + ObservationGracePeriod).
// Upon context cancellation, its expected that we return any usable values within ObservationGracePeriod.
func (s *multiStreamPipeline) executeRun(ctx context.Context) (*pipeline.Run, pipeline.TaskRunResults, error) {
	run, trrs, err := s.runner.ExecuteRun(ctx, s.spec, s.newVars())
	if err != nil {
		return nil, nil, fmt.Errorf("error executing run for spec ID %v: %w", s.spec.ID, err)
	}

	return run, trrs, err
}
