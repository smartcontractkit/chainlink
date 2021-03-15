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
	jobID  int32
	logger logger.Logger
}

// NewPipelineRun constructs a new PipelineRun
func NewPipelineRun(
	runner pipeline.Runner,
	spec pipeline.Spec,
	jobID int32,
	logger logger.Logger,
) PipelineRun {
	return PipelineRun{
		runner: runner,
		spec:   spec,
		jobID:  jobID,
		logger: logger,
	}
}

// Execute executes a pipeline run, extracts the singular result and converts it
// to a decimal.
func (run *PipelineRun) Execute() (int64, *decimal.Decimal, error) {
	ctx := context.Background()
	runID, results, err := run.runner.ExecuteAndInsertNewRun(ctx, run.spec, run.logger)
	if err != nil {
		return runID, nil, errors.Wrapf(err, "error executing new run for job ID %v", run.jobID)
	}

	result, err := results.SingularResult()
	if err != nil {
		return runID, nil, errors.Wrapf(err, "error getting singular result for job ID %v", run.jobID)
	}
	if result.Error != nil {
		return runID, nil, result.Error
	}

	dec, err := utils.ToDecimal(result.Value)
	if err != nil {
		return runID, nil, err
	}

	return runID, &dec, nil
}
