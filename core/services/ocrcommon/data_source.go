package ocrcommon

import (
	"context"
	errjoin "errors"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	ocr1types "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

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

	chEnhancedTelemetry chan<- EnhancedTelemetryData
}

type Saver interface {
	Save(run *pipeline.Run)
}

type dataSourceBase struct {
	inMemoryDataSource
	saver Saver
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

func NewDataSourceV1(pr pipeline.Runner, jb job.Job, spec pipeline.Spec, lggr logger.Logger, s Saver, chEnhancedTelemetry chan EnhancedTelemetryData) ocr1types.DataSource {
	return &dataSource{
		dataSourceBase: dataSourceBase{
			inMemoryDataSource: inMemoryDataSource{
				pipelineRunner:      pr,
				jb:                  jb,
				spec:                spec,
				lggr:                lggr,
				chEnhancedTelemetry: chEnhancedTelemetry,
			},
			saver: s,
		},
	}
}

func NewDataSourceV2(pr pipeline.Runner, jb job.Job, spec pipeline.Spec, lggr logger.Logger, s Saver, enhancedTelemChan chan EnhancedTelemetryData) median.DataSource {
	return &dataSourceV2{
		dataSourceBase: dataSourceBase{
			inMemoryDataSource: inMemoryDataSource{
				pipelineRunner:      pr,
				jb:                  jb,
				spec:                spec,
				lggr:                lggr,
				chEnhancedTelemetry: enhancedTelemChan,
			},
			saver: s,
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

const defaultInMemoryCacheDuration = time.Minute * 5

func NewInMemoryDataSourceCache(ds median.DataSource, cacheExpiryDuration time.Duration) (median.DataSource, error) {
	inMemoryDS, ok := ds.(*inMemoryDataSource)
	if !ok {
		return nil, errors.Errorf("unsupported data source type: %T, only inMemoryDataSource supported", ds)
	}

	if cacheExpiryDuration == 0 {
		cacheExpiryDuration = defaultInMemoryCacheDuration
	}

	dsCache := &inMemoryDataSourceCache{
		cacheExpiration:    cacheExpiryDuration,
		inMemoryDataSource: inMemoryDS,
	}
	go func() { dsCache.updater() }()
	return dsCache, nil
}

var _ ocr1types.DataSource = (*dataSource)(nil)

func setEATelemetry(ds *inMemoryDataSource, finalResult pipeline.FinalResult, trrs pipeline.TaskRunResults, timestamp ObservationTimestamp) {
	promSetFinalResultMetrics(ds, &finalResult)
	promSetBridgeParseMetrics(ds, &trrs)
	if ShouldCollectEnhancedTelemetry(&ds.jb) {
		EnqueueEnhancedTelem(ds.chEnhancedTelemetry, EnhancedTelemetryData{
			TaskRunResults: trrs,
			FinalResults:   finalResult,
			RepTimestamp:   timestamp,
		})
	} else {
		ds.lggr.Infow("Enhanced telemetry is disabled for job", "job", ds.jb.Name)
	}
}

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
func (ds *inMemoryDataSource) executeRun(ctx context.Context) (*pipeline.Run, pipeline.TaskRunResults, error) {
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
		return nil, pipeline.TaskRunResults{}, errors.Wrapf(err, "error executing run for spec ID %v", ds.spec.ID)
	}

	return run, trrs, err
}

// parse uses the FinalResult into a big.Int and stores it in the bridge metadata
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
	_, trrs, err := ds.executeRun(ctx)
	if err != nil {
		return nil, err
	}

	finalResult := trrs.FinalResult(ds.lggr)
	setEATelemetry(ds, finalResult, trrs, ObservationTimestamp{
		Round:        timestamp.Round,
		Epoch:        timestamp.Epoch,
		ConfigDigest: timestamp.ConfigDigest.Hex(),
	})

	return ds.parse(finalResult)
}

// inMemoryDataSourceCache is a time based cache wrapper for inMemoryDataSource.
// If cache update is overdue Observe defaults to standard inMemoryDataSource behaviour.
type inMemoryDataSourceCache struct {
	*inMemoryDataSource
	cacheExpiration time.Duration
	mu              sync.RWMutex
	latestUpdateErr error
	latestTrrs      pipeline.TaskRunResults
	latestResult    pipeline.FinalResult
}

// updater periodically updates data source cache.
func (ds *inMemoryDataSourceCache) updater() {
	ticker := time.NewTicker(ds.cacheExpiration)
	for ; true; <-ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		if err := ds.updateCache(ctx); err != nil {
			ds.lggr.Infow("failed to update cache", "err", err)
		}
		cancel()
	}
}

func (ds *inMemoryDataSourceCache) updateCache(ctx context.Context) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	_, ds.latestTrrs, ds.latestUpdateErr = ds.executeRun(ctx)
	if ds.latestUpdateErr != nil {
		return errors.Wrapf(ds.latestUpdateErr, "error executing run for spec ID %v", ds.spec.ID)
	} else if ds.latestTrrs.FinalResult(ds.lggr).HasErrors() {
		ds.latestUpdateErr = errjoin.Join(ds.latestTrrs.FinalResult(ds.lggr).AllErrors...)
		return errors.Wrapf(ds.latestUpdateErr, "error executing run for spec ID %v", ds.spec.ID)
	}

	ds.latestResult = ds.latestTrrs.FinalResult(ds.lggr)
	return nil
}

func (ds *inMemoryDataSourceCache) get(ctx context.Context) (pipeline.FinalResult, pipeline.TaskRunResults) {
	ds.mu.RLock()
	// updater didn't error, so we know that the latestResult is fresh
	if ds.latestUpdateErr == nil {
		defer ds.mu.RUnlock()
		return ds.latestResult, ds.latestTrrs
	}
	ds.mu.RUnlock()

	if err := ds.updateCache(ctx); err != nil {
		ds.lggr.Errorw("failed to update cache, returning stale result now", "err", err)
	}

	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.latestResult, ds.latestTrrs
}

func (ds *inMemoryDataSourceCache) Observe(ctx context.Context, timestamp ocr2types.ReportTimestamp) (*big.Int, error) {
	latestResult, latestTrrs := ds.get(ctx)
	setEATelemetry(ds.inMemoryDataSource, latestResult, latestTrrs, ObservationTimestamp{
		Round:        timestamp.Round,
		Epoch:        timestamp.Epoch,
		ConfigDigest: timestamp.ConfigDigest.Hex(),
	})

	return ds.parse(latestResult)
}

func (ds *dataSourceBase) observe(ctx context.Context, timestamp ObservationTimestamp) (*big.Int, error) {
	run, trrs, err := ds.inMemoryDataSource.executeRun(ctx)
	if err != nil {
		return nil, err
	}

	// Save() does the database write in a non-blocking fashion
	// so we can return the observation results immediately.
	// This is helpful in the case of a blocking API call, where
	// we reach the passed in context deadline and we want to
	// immediately return any result we have and do not want to have
	// a db write block that.
	ds.saver.Save(run)

	finalResult := trrs.FinalResult(ds.lggr)
	setEATelemetry(&ds.inMemoryDataSource, finalResult, trrs, timestamp)

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
