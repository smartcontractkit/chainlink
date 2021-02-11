package offchainreporting

import (
	"context"

	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/services/job"

	"github.com/pkg/errors"
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
}

var _ ocrtypes.DataSource = (*dataSource)(nil)

func newDatasource(db *gorm.DB, jobID int32, pipelineRunner pipeline.Runner) (*dataSource, error) {
	var j job.SpecDB
	err := db.Preload("PipelineSpec.PipelineTaskSpecs").Find(&j, "id = ?", jobID).Error
	if err != nil {
		return nil, errors.Wrapf(err, "could not load pipeline_spec for job ID %v", jobID)
	}
	return &dataSource{jobID: jobID, spec: *j.PipelineSpec, pipelineRunner: pipelineRunner}, nil
}

// The context passed in here has a timeout of observationTimeout.
// Gorm/pgx doesn't return a helpful error upon cancellation, so we manually check for cancellation and return a
// appropriate error (FIXME: How does this work after the current refactoring?)
func (ds dataSource) Observe(ctx context.Context) (ocrtypes.Observation, error) {
	results, err := ds.pipelineRunner.ExecuteAndInsertNewRun(ctx, ds.spec)
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
