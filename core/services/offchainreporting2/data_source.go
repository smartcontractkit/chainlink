package offchainreporting2

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/bridges"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
)

// dataSource is an abstraction over the process of initiating a pipeline run
// and capturing the result. Additionally, it converts the result to an
// ocrtypes.Observation (*big.Int), as expected by the offchain reporting library.
type dataSource struct {
	runResults         chan<- pipeline.Run
	inMemoryDataSource inMemoryDataSource
}

var _ median.DataSource = (*dataSource)(nil)

func NewDataSource(pipeline pipeline.Runner, jobSpec job.Job, spec pipeline.Spec, log logger.Logger, runResults chan<- pipeline.Run) *dataSource {
	ds := NewInMemoryDataSource(pipeline, jobSpec, spec, log)

	return &dataSource{
		inMemoryDataSource: *ds,
		runResults: runResults,
	}
}

// The context passed in here has a timeout of (ObservationTimeout + ObservationGracePeriod).
// Upon context cancellation, its expected that we return any usable values within ObservationGracePeriod.
func (ds *dataSource) Observe(ctx context.Context) (*big.Int, error) {
	run, finalResult, err := ds.inMemoryDataSource.executeRun(ctx)
	if err != nil {
		return nil, err
	}

	// Do the database write in a non-blocking fashion
	// so we can return the observation results immediately.
	// This is helpful in the case of a blocking API call, where
	// we reach the passed in context deadline and we want to
	// immediately return any result we have and do not want to have
	// a db write block that.
	select {
	case ds.runResults <- run:
	default:
		return nil, errors.Errorf("unable to enqueue run save for job ID %v, buffer full", ds.inMemoryDataSource.spec.JobID)
	}

	return ds.inMemoryDataSource.parse(finalResult)
}

// same implementation as dataSource but without sending to channel for saving to DB
type inMemoryDataSource struct {
	pipelineRunner        pipeline.Runner
	jobSpec               job.Job
	spec                  pipeline.Spec
	ocrLogger             logger.Logger
	currentBridgeMetadata bridges.BridgeMetaData
}

func NewInMemoryDataSource(pipeline pipeline.Runner, jobSpec job.Job, spec pipeline.Spec, log logger.Logger) *inMemoryDataSource {
	return &inMemoryDataSource{
		pipelineRunner: pipeline,
		jobSpec:        jobSpec,
		spec:           spec,
		ocrLogger:      log,
	}
}

func (ds *inMemoryDataSource) executeRun(ctx context.Context) (pipeline.Run, pipeline.FinalResult, error) {
	md, err := bridges.MarshalBridgeMetaData(ds.currentBridgeMetadata.LatestAnswer, ds.currentBridgeMetadata.UpdatedAt)
	if err != nil {
		logger.Warnw("unable to attach metadata for run", "err", err)
	}

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    ds.jobSpec.ID,
			"externalJobID": ds.jobSpec.ExternalJobID,
			"name":          ds.jobSpec.Name.ValueOrZero(),
		},
		"jobRun": map[string]interface{}{
			"meta": md,
		},
	})

	run, trrs, err := ds.pipelineRunner.ExecuteRun(ctx, ds.spec, vars, ds.ocrLogger)
	if err != nil {
		return pipeline.Run{}, pipeline.FinalResult{}, errors.Wrapf(err, "error executing run for spec ID %v", ds.spec.ID)
	}
	finalResult := trrs.FinalResult()
	return run, finalResult, err
}

func (ds *inMemoryDataSource) parse(finalResult pipeline.FinalResult) (*big.Int, error) {
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
	ds.currentBridgeMetadata = bridges.BridgeMetaData{
		LatestAnswer: asDecimal.BigInt(),
		UpdatedAt:    big.NewInt(time.Now().Unix()),
	}
	return asDecimal.BigInt(), nil
}

func (ds *inMemoryDataSource) Observe(ctx context.Context) (*big.Int, error) {
	_, finalResult, err := ds.executeRun(ctx)
	if err != nil {
		return nil, err
	}
	return ds.parse(finalResult)
}
