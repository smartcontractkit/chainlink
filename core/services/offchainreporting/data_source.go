package offchainreporting

import (
	"context"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

// dataSource is an abstraction over the process of initiating a pipeline run
// and capturing the result. Additionally, it converts the result to an
// ocrtypes.Observation (*big.Int), as expected by the offchain reporting library.
type dataSource struct {
	pipelineRunner pipeline.Runner
	jobID          int32
	spec           pipeline.Spec
	ocrLogger      logger.Logger
}

var _ ocrtypes.DataSource = (*dataSource)(nil)

// The context passed in here has a timeout of observationTimeout.
// Gorm/pgx doesn't return a helpful error upon cancellation, so we manually check for cancellation and return a
// appropriate error (FIXME: How does this work after the current refactoring?)
func (ds dataSource) Observe(ctx context.Context) (ocrtypes.Observation, error) {
	_, results, err := ds.pipelineRunner.ExecuteAndInsertNewRun(ctx, ds.spec, ds.ocrLogger)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing new run for job ID %v", ds.jobID)
	}

	result, err := results.SingularResult()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting singular result for job ID %v", ds.jobID)
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
