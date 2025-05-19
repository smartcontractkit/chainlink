package ocrcommon

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	ocr1types "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	serializablebig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median/config"
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

const defaultUpdateInterval = time.Minute * 5
const defaultStalenessAlertThreshold = time.Hour * 24
const dataSourceCacheKey = "dscache"

type DataSourceCacheService interface {
	Start(context.Context) error
	Close() error
	median.DataSource
}

func NewInMemoryDataSourceCache(ds median.DataSource, kvStore job.KVStore, cacheCfg *config.JuelsPerFeeCoinCache) (DataSourceCacheService, error) {
	inMemoryDS, ok := ds.(*inMemoryDataSource)
	if !ok {
		return nil, errors.Errorf("unsupported data source type: %T, only inMemoryDataSource supported", ds)
	}
	var updateInterval, stalenessAlertThreshold time.Duration
	if cacheCfg == nil {
		updateInterval = defaultUpdateInterval
		stalenessAlertThreshold = defaultStalenessAlertThreshold
	} else {
		updateInterval, stalenessAlertThreshold = cacheCfg.UpdateInterval.Duration(), cacheCfg.StalenessAlertThreshold.Duration()
		if updateInterval == 0 {
			updateInterval = defaultUpdateInterval
		}
		if stalenessAlertThreshold == 0 {
			stalenessAlertThreshold = defaultStalenessAlertThreshold
		}
	}

	dsCache := &inMemoryDataSourceCache{
		inMemoryDataSource:      inMemoryDS,
		kvStore:                 kvStore,
		updateInterval:          updateInterval,
		stalenessAlertThreshold: stalenessAlertThreshold,
		chStop:                  make(chan struct{}),
		chDone:                  make(chan struct{}),
	}
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
		ds.lggr.Warnf("unable to attach metadata for run, err: %v", err)
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

	run, trrs, err := ds.pipelineRunner.ExecuteRun(ctx, ds.spec, vars)
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

	finalResult := trrs.FinalResult()
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
	// updateInterval indicates duration between cache updates.
	// Even if update fail, previous values are returned.
	updateInterval time.Duration
	// stalenessAlertThreshold indicates duration before logs raise severity level because of stale cache.
	stalenessAlertThreshold time.Duration
	mu                      sync.RWMutex
	chStop                  services.StopChan
	chDone                  chan struct{}
	latestUpdateErr         error
	latestTrrs              pipeline.TaskRunResults
	latestResult            pipeline.FinalResult
	kvStore                 job.KVStore
}

func (ds *inMemoryDataSourceCache) Start(context.Context) error {
	go func() { ds.updater() }()
	return nil
}

func (ds *inMemoryDataSourceCache) Close() error {
	close(ds.chStop)
	<-ds.chDone
	return nil
}

// updater periodically updates data source cache.
func (ds *inMemoryDataSourceCache) updater() {
	ticker := time.NewTicker(ds.updateInterval)
	updateCache := func() {
		ctx, cancel := ds.chStop.CtxWithTimeout(time.Second * 10)
		defer cancel()
		if err := ds.updateCache(ctx); err != nil {
			ds.lggr.Warnf("failed to update cache, err: %v", err)
		}
	}

	updateCache()
	for {
		select {
		case <-ticker.C:
			updateCache()
		case <-ds.chStop:
			close(ds.chDone)
			return
		}
	}
}

type ResultTimePair struct {
	Result serializablebig.Big `json:"result"`
	Time   time.Time           `json:"time"`
}

func (ds *inMemoryDataSourceCache) updateCache(ctx context.Context) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	_, latestTrrs, err := ds.executeRun(ctx)
	if err != nil {
		previousUpdateErr := ds.latestUpdateErr
		ds.latestUpdateErr = err
		// warn log if previous cache update also errored
		if previousUpdateErr != nil {
			ds.lggr.Warnf("consecutive cache updates errored: previous err: %v new err: %v", previousUpdateErr, ds.latestUpdateErr)
		}

		return errors.Wrapf(ds.latestUpdateErr, "error updating in memory data source cache for spec ID %v", ds.spec.ID)
	}

	value, err := ds.inMemoryDataSource.parse(latestTrrs.FinalResult())
	if err != nil {
		ds.latestUpdateErr = errors.Wrapf(err, "invalid result")
		return ds.latestUpdateErr
	}

	// update cache values
	ds.latestTrrs = latestTrrs
	ds.latestResult = ds.latestTrrs.FinalResult()
	ds.latestUpdateErr = nil

	// backup in case data source fails continuously and node gets rebooted
	timePairBytes, err := json.Marshal(&ResultTimePair{Result: *serializablebig.New(value), Time: time.Now()})
	if err != nil {
		return fmt.Errorf("failed to marshal result time pair, err: %w", err)
	}

	if err = ds.kvStore.Store(ctx, dataSourceCacheKey, timePairBytes); err != nil {
		ds.lggr.Errorf("failed to persist latest task run value, err: %v", err)
	}

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
		ds.lggr.Warnf("failed to update cache, returning stale result now, err: %v", err)
	}

	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.latestResult, ds.latestTrrs
}

func (ds *inMemoryDataSourceCache) Observe(ctx context.Context, timestamp ocr2types.ReportTimestamp) (*big.Int, error) {
	var resTime ResultTimePair
	latestResult, latestTrrs := ds.get(ctx)
	if latestTrrs == nil {
		ds.lggr.Warnf("cache is empty, returning persisted value now")

		timePairBytes, err := ds.kvStore.Get(ctx, dataSourceCacheKey)
		if err != nil {
			return nil, fmt.Errorf("in memory data source cache is empty and failed to get backup persisted value, err: %w", err)
		}

		if err = json.Unmarshal(timePairBytes, &resTime); err != nil {
			return nil, fmt.Errorf("in memory data source cache is empty and failed to unmarshal backup persisted value, err: %w", err)
		}

		if time.Since(resTime.Time) >= ds.stalenessAlertThreshold {
			ds.lggr.Errorf("in memory data source cache is empty and the persisted value hasn't been updated for over %v, latestUpdateErr is: %v", ds.stalenessAlertThreshold, ds.latestUpdateErr)
		}
		return resTime.Result.ToInt(), nil
	}

	setEATelemetry(ds.inMemoryDataSource, latestResult, latestTrrs, ObservationTimestamp{
		Round:        timestamp.Round,
		Epoch:        timestamp.Epoch,
		ConfigDigest: timestamp.ConfigDigest.Hex(),
	})

	// if last update was unsuccessful, check how much time passed since a successful update
	if ds.latestUpdateErr != nil {
		if time.Since(ds.latestTrrs.GetTaskRunResultsFinishedAt()) >= ds.stalenessAlertThreshold {
			ds.lggr.Errorf("in memory cache is old and hasn't been updated for over %v, latestUpdateErr is: %v", ds.stalenessAlertThreshold, ds.latestUpdateErr)
		}
	}
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

	finalResult := trrs.FinalResult()
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
