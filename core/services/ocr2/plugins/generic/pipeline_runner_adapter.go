package generic

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var _ types.PipelineRunnerService = (*PipelineRunnerAdapter)(nil)

type pipelineRunner interface {
	ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger) (run *pipeline.Run, trrs pipeline.TaskRunResults, err error)
}

type PipelineRunnerAdapter struct {
	runner pipelineRunner
	job    job.Job
	logger logger.Logger
}

func (p *PipelineRunnerAdapter) ExecuteRun(ctx context.Context, spec string, vars types.Vars, options types.Options) (types.TaskResults, error) {
	s := pipeline.Spec{
		DotDagSource:    spec,
		CreatedAt:       time.Now(),
		MaxTaskDuration: models.Interval(options.MaxTaskDuration),
		JobID:           p.job.ID,
		JobName:         p.job.Name.ValueOrZero(),
		JobType:         string(p.job.Type),
	}

	defaultVars := map[string]interface{}{
		"jb": map[string]interface{}{
			"databaseID":    p.job.ID,
			"externalJobID": p.job.ExternalJobID,
			"name":          p.job.Name.ValueOrZero(),
		},
	}

	merge(defaultVars, vars.Vars)

	finalVars := pipeline.NewVarsFrom(defaultVars)
	_, trrs, err := p.runner.ExecuteRun(ctx, s, finalVars, p.logger)
	if err != nil {
		return nil, err
	}

	taskResults := make([]types.TaskResult, len(trrs))
	for i, trr := range trrs {
		taskResults[i] = types.TaskResult{
			ID:    trr.ID.String(),
			Type:  string(trr.Task.Type()),
			Index: int(trr.Task.OutputIndex()),

			TaskValue: types.TaskValue{
				Value:      trr.Result.Value,
				Error:      trr.Result.Error,
				IsTerminal: len(trr.Task.Outputs()) == 0,
			},
		}
	}
	return taskResults, nil
}

func NewPipelineRunnerAdapter(logger logger.Logger, job job.Job, runner pipelineRunner) *PipelineRunnerAdapter {
	return &PipelineRunnerAdapter{
		logger: logger,
		job:    job,
		runner: runner,
	}
}

// merge merges mapTwo into mapOne, modifying mapOne in the process.
func merge(mapOne, mapTwo map[string]interface{}) {
	for k, v := range mapTwo {
		// if `mapOne` doesn't have `k`, then nothing to do, just assign v to `mapOne`.
		if _, ok := mapOne[k]; !ok {
			mapOne[k] = v
		} else {
			vAsMap, vOK := v.(map[string]interface{})
			mapOneVAsMap, moOK := mapOne[k].(map[string]interface{})
			if vOK && moOK {
				merge(mapOneVAsMap, vAsMap)
			} else {
				mapOne[k] = v
			}
		}
	}
}
