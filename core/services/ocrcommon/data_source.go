package ocrcommon

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	ocr1types "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// inMemoryDataSource is an abstraction over the process of initiating a pipeline run
// and returning the result
type inMemoryDataSource struct {
	pipelineRunner pipeline.Runner
	jb             job.Job
	spec           pipeline.Spec
	lggr           logger.Logger

	current bridges.BridgeMetaData
	mu      sync.RWMutex

	monitoringEndpoint commontypes.MonitoringEndpoint
}

type dataSourceBase struct {
	inMemoryDataSource
	runResults chan<- pipeline.Run
}

// dataSource implements dataSourceBase with the proper Observe return type for ocr1
type dataSource struct {
	dataSourceBase
}

// dataSourceV2 implements dataSourceBase with the proper Observe return type for ocr2
type dataSourceV2 struct {
	dataSourceBase
}

// ObservationTimestamp abstracts ocr2types.ReportTimestamp and ocr1types.ReportTimestamp
type ObservationTimestamp struct {
	Round        uint8
	Epoch        uint32
	ConfigDigest string
}

func NewDataSourceV1(pr pipeline.Runner, jb job.Job, spec pipeline.Spec, lggr logger.Logger, runResults chan<- pipeline.Run, endpoint commontypes.MonitoringEndpoint) ocr1types.DataSource {
	return &dataSource{
		dataSourceBase: dataSourceBase{
			inMemoryDataSource: inMemoryDataSource{
				pipelineRunner:     pr,
				jb:                 jb,
				spec:               spec,
				lggr:               lggr,
				monitoringEndpoint: endpoint,
			},
			runResults: runResults,
		},
	}
}

func NewDataSourceV2(pr pipeline.Runner, jb job.Job, spec pipeline.Spec, lggr logger.Logger, runResults chan<- pipeline.Run, endpoint commontypes.MonitoringEndpoint) median.DataSource {
	return &dataSourceV2{
		dataSourceBase: dataSourceBase{
			inMemoryDataSource: inMemoryDataSource{
				pipelineRunner:     pr,
				jb:                 jb,
				spec:               spec,
				lggr:               lggr,
				monitoringEndpoint: endpoint,
			},
			runResults: runResults,
		},
	}
}

func NewInMemoryDataSource(pr pipeline.Runner, jb job.Job, spec pipeline.Spec, lggr logger.Logger) median.DataSource {
	return &inMemoryDataSource{
		pipelineRunner: pr,
		jb:             jb,
		spec:           spec,
		lggr:           lggr,
	}
}

var _ ocr1types.DataSource = (*dataSource)(nil)

func (ds *inMemoryDataSource) updateAnswer(a *big.Int) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.current = bridges.BridgeMetaData{
		LatestAnswer: a,
		UpdatedAt:    big.NewInt(time.Now().Unix()),
	}
}

func (ds *inMemoryDataSource) currentAnswer() (*big.Int, *big.Int) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.current.LatestAnswer, ds.current.UpdatedAt
}

// The context passed in here has a timeout of (ObservationTimeout + ObservationGracePeriod).
// Upon context cancellation, its expected that we return any usable values within ObservationGracePeriod.
func (ds *inMemoryDataSource) executeRun(ctx context.Context, timestamp ObservationTimestamp) (pipeline.Run, pipeline.FinalResult, error) {
	md, err := bridges.MarshalBridgeMetaData(ds.currentAnswer())
	if err != nil {
		ds.lggr.Warnw("unable to attach metadata for run", "err", err)
	}

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jb": map[string]interface{}{
			"databaseID":    ds.jb.ID,
			"externalJobID": ds.jb.ExternalJobID,
			"name":          ds.jb.Name.ValueOrZero(),
		},
		"jobRun": map[string]interface{}{
			"meta": md,
		},
	})

	run, trrs, err := ds.pipelineRunner.ExecuteRun(ctx, ds.spec, vars, ds.lggr)
	if err != nil {
		return pipeline.Run{}, pipeline.FinalResult{}, errors.Wrapf(err, "error executing run for spec ID %v", ds.spec.ID)
	}
	finalResult := trrs.FinalResult(ds.lggr)
	promSetBridgeParseMetrics(ds, &trrs)
	promSetFinalResultMetrics(ds, &finalResult)
	collectEATelemetry(ds, &trrs, &finalResult, timestamp)

	return run, finalResult, err
}

// parse uses the finalResult into a big.Int and stores it in the bridge metadata
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
	ds.updateAnswer(asDecimal.BigInt())
	return asDecimal.BigInt(), nil
}

// Observe without saving to DB
func (ds *inMemoryDataSource) Observe(ctx context.Context, timestamp ocr2types.ReportTimestamp) (*big.Int, error) {
	_, finalResult, err := ds.executeRun(ctx, ObservationTimestamp{
		Round:        timestamp.Round,
		Epoch:        timestamp.Epoch,
		ConfigDigest: timestamp.ConfigDigest.Hex(),
	})
	if err != nil {
		return nil, err
	}
	return ds.parse(finalResult)
}

func (ds *dataSourceBase) observe(ctx context.Context, timestamp ObservationTimestamp) (*big.Int, error) {
	run, finalResult, err := ds.inMemoryDataSource.executeRun(ctx, timestamp)
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
		// If we're unable to enqueue a write, still return the value we have but warn.
		ds.lggr.Warnf("unable to enqueue run save for job ID %d, buffer full", ds.inMemoryDataSource.spec.JobID)
		return ds.inMemoryDataSource.parse(finalResult)
	}

	return ds.inMemoryDataSource.parse(finalResult)
}

// Observe with saving to DB, satisfies ocr1 interface
func (ds *dataSource) Observe(ctx context.Context, timestamp ocr1types.ReportTimestamp) (ocr1types.Observation, error) {
	return ds.observe(ctx, ObservationTimestamp{
		Round:        timestamp.Round,
		Epoch:        timestamp.Epoch,
		ConfigDigest: timestamp.ConfigDigest.Hex(),
	})
}

// Observe with saving to DB, satisfies ocr2 interface
func (ds *dataSourceV2) Observe(ctx context.Context, timestamp ocr2types.ReportTimestamp) (*big.Int, error) {
	return ds.observe(ctx, ObservationTimestamp{
		Round:        timestamp.Round,
		Epoch:        timestamp.Epoch,
		ConfigDigest: timestamp.ConfigDigest.Hex(),
	})
}
