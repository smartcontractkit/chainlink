package fluxmonitorv2

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// PipelineRun wraps a the pipeline to execute a single pipeline run
type PipelineRun struct {
	runner pipeline.Runner
	spec   pipeline.Spec
	logger logger.Logger
}

// NewPipelineRun constructs a new PipelineRun
func NewPipelineRun(
	runner pipeline.Runner,
	spec pipeline.Spec,
	logger logger.Logger,
) PipelineRun {
	return PipelineRun{
		runner: runner,
		spec:   spec,
		logger: logger,
	}
}

// Execute executes a pipeline run, extracts the singular result and converts it
// to a decimal.
func (run *PipelineRun) Execute(meta map[string]interface{}) (int64, *decimal.Decimal, error) {
	ctx := context.Background()
	runID, results, err := run.runner.ExecuteAndInsertFinishedRun(ctx, run.spec, pipeline.JSONSerializable{Val: meta}, run.logger, false)
	if err != nil {
		return runID, nil, errors.Wrapf(err, "error executing new run for job ID %v name %v", run.spec.JobID, run.spec.JobName)
	}

	result, err := results.SingularResult()
	if err != nil {
		return runID, nil, errors.Wrapf(err, "error getting singular result for job ID %v name %v", run.spec.JobID, run.spec.JobName)
	}
	if result.Error != nil {
		return runID, nil, result.Error
	}

	dec, err := utils.ToDecimal(result.Value)
	if err != nil {
		return runID, nil, errors.Wrap(err, "cannot convert result to decimal")
	}

	return runID, &dec, nil
}
