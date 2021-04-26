package offchainreporting

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/store/models"

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
	pipelineRunner        pipeline.Runner
	spec                  pipeline.Spec
	ocrLogger             logger.Logger
	runResults            chan<- pipeline.RunWithResults
	currentBridgeMetadata models.BridgeMetaData
}

var _ ocrtypes.DataSource = (*dataSource)(nil)

// The context passed in here has a timeout of (ObservationTimeout + ObservationGracePeriod).
// Upon context cancellation, its expected that we return any usable values within ObservationGracePeriod.
func (ds *dataSource) Observe(ctx context.Context) (ocrtypes.Observation, error) {
	var observation ocrtypes.Observation
	start := time.Now()
	md, err := models.MarshalBridgeMetaData(ds.currentBridgeMetadata.LatestAnswer, ds.currentBridgeMetadata.UpdatedAt)
	if err != nil {
		logger.Warnw("unable to attach metadata for run", "err", err)
	}
	trrs, err := ds.pipelineRunner.ExecuteRun(ctx, ds.spec, pipeline.JSONSerializable{
		Val: md,
	}, ds.ocrLogger)
	if err != nil {
		return observation, errors.Wrapf(err, "error executing run for spec ID %v", ds.spec.ID)
	}
	end := time.Now()

	var run pipeline.Run
	run.PipelineSpecID = ds.spec.ID
	run.CreatedAt = start
	run.FinishedAt = &end

	finalResult := trrs.FinalResult()
	run.Outputs = finalResult.OutputsDB()
	run.Errors = finalResult.ErrorsDB()

	// Do the database write in a non-blocking fashion
	// so we can return the observation results immediately.
	// This is helpful in the case of a blocking API call, where
	// we reach the passed in context deadline and we want to
	// immediately return any result we have and do not want to have
	// a db write block that.
	select {
	case ds.runResults <- pipeline.RunWithResults{
		Run:            run,
		TaskRunResults: trrs,
	}:
	default:
		return nil, errors.Errorf("unable to enqueue run save for job ID %v, buffer full", ds.spec.JobID)
	}

	result, err := finalResult.SingularResult()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting singular result for job ID %v", ds.spec.JobID)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	asDecimal, err := utils.ToDecimal(result.Value)
	if err != nil {
		return nil, errors.Wrap(err, "cannot convert observation to decimal")
	}
	ds.currentBridgeMetadata = models.BridgeMetaData{
		LatestAnswer: asDecimal.BigInt(),
		UpdatedAt:    big.NewInt(time.Now().Unix()),
	}
	return asDecimal.BigInt(), nil
}
