package generic

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var _ core.PipelineRunnerService = (*PipelineRunnerAdapter)(nil)

type pipelineRunner interface {
	ExecuteAndInsertFinishedRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, saveSuccessfulTaskRuns bool) (runID int64, results pipeline.TaskRunResults, err error)
}

type PipelineRunnerAdapter struct {
	runner pipelineRunner
	job    job.Job
	logger logger.Logger
}

func (p *PipelineRunnerAdapter) ExecuteRun(ctx context.Context, spec string, vars core.Vars, options core.Options) (core.TaskResults, error) {
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
	_, trrs, err := p.runner.ExecuteAndInsertFinishedRun(ctx, s, finalVars, true)
	if err != nil {
		return nil, err
	}

	taskResults := make([]core.TaskResult, len(trrs))
	for i, trr := range trrs {
		taskResults[i] = core.TaskResult{
			ID:    trr.ID.String(),
			Type:  string(trr.Task.Type()),
			Index: int(trr.Task.OutputIndex()),
			TaskValue: core.TaskValue{
				Value:      trr.Result.OutputDB(),
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
