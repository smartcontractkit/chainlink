package offchainreporting

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
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
	runID, err := ds.pipelineRunner.CreateRun(ctx, ds.jobID, nil)
	endCreate := time.Now()
	if ctx.Err() != nil {
		return nil, errors.Errorf("context cancelled due to timeout or shutdown, cancel create run. Runtime %v", endCreate.Sub(start))
	}
	if err != nil {
		logger.Errorw("Error creating new pipeline run", "jobID", ds.jobID, "error", err)
		return nil, err
	}

	err = ds.pipelineRunner.AwaitRun(ctx, runID)
	endAwait := time.Now()
	if ctx.Err() != nil {
		return nil, errors.Errorf("context cancelled due to timeout or shutdown, cancel await run. Runtime %v", endAwait.Sub(start))
	}
	if err != nil {
		return nil, err
	}

	results, err := ds.pipelineRunner.ResultsForRun(ctx, runID)
	endResults := time.Now()
	if ctx.Err() != nil {
		return nil, errors.Errorf("context cancelled due to timeout or shutdown, cancel get results for run. Runtime %v", endResults.Sub(start))
	} else if err != nil {
		return nil, errors.Wrapf(err, "pipeline error")
	} else if len(results) != 1 {
		return nil, errors.Errorf("offchain reporting pipeline should have a single output (job spec ID: %v, pipeline run ID: %v)", ds.jobID, runID)
	}

	if results[0].Error != nil {
		return nil, results[0].Error
	}

	asDecimal, err := utils.ToDecimal(results[0].Value)
	if err != nil {
		return nil, err
	}
	return ocrtypes.Observation(asDecimal.BigInt()), nil
}
