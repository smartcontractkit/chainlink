package offchainreporting

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

// dataSource is an abstraction over the process of initiating a pipeline run
// and capturing the result.  Additionally, it converts the result to an
// ocrtypes.Observation (*big.Int), as expected by the offchain reporting library.
type dataSource struct {
	pipelineRunner pipeline.Runner
	jobID          int32
}

var _ ocrtypes.DataSource = (*dataSource)(nil)

// The context passed in here has a timeout of observationTimeout.
// Gorm/pgx doesn't return a helpful error upon cancellation, so we manually check for cancellation and return a
// appropriate error.
func (ds dataSource) Observe(ctx context.Context) (ocrtypes.Observation, error) {
	start := time.Now()

	// FIXME: Pull out the spec load from NewRun and make it pure
	run, err := ds.pipelineRunner.NewRun(ctx, ds.jobID, start)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating new run for job ID %v", ds.jobID)
	}

	trrs, err := ds.pipelineRunner.ExecuteRun(ctx, run)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing run for job ID %v", ds.jobID)
	}

	run.CreatedAt = start
	end := time.Now()
	run.FinishedAt = &end

	// TODO: Might wanna add some logging with runID
	if _, err := ds.pipelineRunner.InsertFinishedRunWithResults(ctx, run, trrs); err != nil {
		return nil, errors.Wrapf(err, "error inserting finished results for job ID %v", ds.jobID)
	}

	// TODO: Can we pull this into a function on []TaskRunResult?
	var result pipeline.Result
	for _, trr := range trrs {
		if trr.IsFinal {
			// FIXME: This assumes there is only one final result and will
			// have to change when the magical "__result__" type is removed
			// https://www.pivotaltracker.com/story/show/176557536
			result = trr.Result
		}
	}

	if result.Error != nil {
		return nil, result.Error
	}

	asDecimal, err := utils.ToDecimal(result.Value)
	if err != nil {
		return nil, err
	}
	return ocrtypes.Observation(asDecimal.BigInt()), nil
}
