package v1

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	pkgerrors "github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	v1types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	v1 "github.com/smartcontractkit/chainlink-data-streams/mercury/v1"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	insufficientBlocksCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_insufficient_blocks_count",
		Help: fmt.Sprintf("Count of times that there were not enough blocks in the chain during observation (need: %d)", nBlocksObservation),
	},
		[]string{"feedID"},
	)
	zeroBlocksCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_zero_blocks_count",
		Help: "Count of times that there were zero blocks in the chain during observation",
	},
		[]string{"feedID"},
	)
)

const nBlocksObservation int = v1.MaxAllowedBlocks

type Runner interface {
	ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger) (run *pipeline.Run, trrs pipeline.TaskRunResults, err error)
}

// Fetcher fetcher data from Mercury server
type Fetcher interface {
	// FetchInitialMaxFinalizedBlockNumber should fetch the initial max finalized block number
	FetchInitialMaxFinalizedBlockNumber(context.Context) (*int64, error)
}

type datasource struct {
	pipelineRunner Runner
	jb             job.Job
	spec           pipeline.Spec
	lggr           logger.Logger
	saver          ocrcommon.Saver
	orm            types.DataSourceORM
	codec          reportcodec.ReportCodec
	feedID         [32]byte

	mu sync.RWMutex

	chEnhancedTelem    chan<- ocrcommon.EnhancedTelemetryMercuryData
	mercuryChainReader mercury.ChainReader
	fetcher            Fetcher
	initialBlockNumber *int64

	insufficientBlocksCounter prometheus.Counter
	zeroBlocksCounter         prometheus.Counter
}

var _ v1.DataSource = &datasource{}

func NewDataSource(orm types.DataSourceORM, pr pipeline.Runner, jb job.Job, spec pipeline.Spec, lggr logger.Logger, s ocrcommon.Saver, enhancedTelemChan chan ocrcommon.EnhancedTelemetryMercuryData, mercuryChainReader mercury.ChainReader, fetcher Fetcher, initialBlockNumber *int64, feedID mercuryutils.FeedID) *datasource {
	return &datasource{pr, jb, spec, lggr, s, orm, reportcodec.ReportCodec{}, feedID, sync.RWMutex{}, enhancedTelemChan, mercuryChainReader, fetcher, initialBlockNumber, insufficientBlocksCount.WithLabelValues(feedID.String()), zeroBlocksCount.WithLabelValues(feedID.String())}
}

type ErrEmptyLatestReport struct {
	Err error
}

func (e ErrEmptyLatestReport) Unwrap() error { return e.Err }

func (e ErrEmptyLatestReport) Error() string {
	return fmt.Sprintf("FetchInitialMaxFinalizedBlockNumber returned empty LatestReport; this is a new feed. No initialBlockNumber was set, tried to use current block number to determine maxFinalizedBlockNumber but got error: %v", e.Err)
}

func (ds *datasource) Observe(ctx context.Context, repts ocrtypes.ReportTimestamp, fetchMaxFinalizedBlockNum bool) (obs v1types.Observation, pipelineExecutionErr error) {
	// setLatestBlocks must come chronologically before observations, along
	// with observationTimestamp, to avoid front-running

	// Errors are not expected when reading from the underlying ChainReader
	if err := ds.setLatestBlocks(ctx, &obs); err != nil {
		return obs, err
	}

	var wg sync.WaitGroup
	if fetchMaxFinalizedBlockNum {
		wg.Add(1)
		go func() {
			defer wg.Done()
			latest, dbErr := ds.orm.LatestReport(ctx, ds.feedID)
			if dbErr != nil {
				obs.MaxFinalizedBlockNumber.Err = dbErr
				return
			}
			if latest != nil {
				obs.MaxFinalizedBlockNumber.Val, obs.MaxFinalizedBlockNumber.Err = ds.codec.CurrentBlockNumFromReport(latest)
				return
			}
			val, fetchErr := ds.fetcher.FetchInitialMaxFinalizedBlockNumber(ctx)
			if fetchErr != nil {
				obs.MaxFinalizedBlockNumber.Err = fetchErr
				return
			}
			if val != nil {
				obs.MaxFinalizedBlockNumber.Val = *val
				return
			}
			if ds.initialBlockNumber == nil {
				if obs.CurrentBlockNum.Err != nil {
					obs.MaxFinalizedBlockNumber.Err = ErrEmptyLatestReport{Err: obs.CurrentBlockNum.Err}
				} else {
					// Subract 1 here because we will later add 1 to the
					// maxFinalizedBlockNumber to get the first validFromBlockNum, which
					// ought to be the same as current block num.
					obs.MaxFinalizedBlockNumber.Val = obs.CurrentBlockNum.Val - 1
					ds.lggr.Infof("FetchInitialMaxFinalizedBlockNumber returned empty LatestReport; this is a new feed so maxFinalizedBlockNumber=%d (initialBlockNumber unset, using currentBlockNum=%d-1)", obs.MaxFinalizedBlockNumber.Val, obs.CurrentBlockNum.Val)
				}
			} else {
				// NOTE: It's important to subtract 1 if the server is missing any past
				// report (brand new feed) since we will add 1 to the
				// maxFinalizedBlockNumber to get the first validFromBlockNum, which
				// ought to be zero.
				//
				// If "initialBlockNumber" is set to zero, this will give a starting block of zero.
				obs.MaxFinalizedBlockNumber.Val = *ds.initialBlockNumber - 1
				ds.lggr.Infof("FetchInitialMaxFinalizedBlockNumber returned empty LatestReport; this is a new feed so maxFinalizedBlockNumber=%d (initialBlockNumber=%d)", obs.MaxFinalizedBlockNumber.Val, *ds.initialBlockNumber)
			}
		}()
	} else {
		obs.MaxFinalizedBlockNumber.Err = errors.New("fetchMaxFinalizedBlockNum=false")
	}
	var trrs pipeline.TaskRunResults
	wg.Add(1)
	go func() {
		defer wg.Done()
		var run *pipeline.Run
		run, trrs, pipelineExecutionErr = ds.executeRun(ctx)
		if pipelineExecutionErr != nil {
			pipelineExecutionErr = fmt.Errorf("Observe failed while executing run: %w", pipelineExecutionErr)
			return
		}

		ds.saver.Save(run)

		// NOTE: trrs comes back as _all_ tasks, but we only want the terminal ones
		// They are guaranteed to be sorted by index asc so should be in the correct order
		var finaltrrs []pipeline.TaskRunResult
		for _, trr := range trrs {
			if trr.IsTerminal() {
				finaltrrs = append(finaltrrs, trr)
			}
		}

		var parsed parseOutput
		parsed, pipelineExecutionErr = ds.parse(finaltrrs)
		if pipelineExecutionErr != nil {
			pipelineExecutionErr = fmt.Errorf("Observe failed while parsing run results: %w", pipelineExecutionErr)
			return
		}
		obs.BenchmarkPrice = parsed.benchmarkPrice
		obs.Bid = parsed.bid
		obs.Ask = parsed.ask
	}()

	wg.Wait()

	if pipelineExecutionErr != nil {
		return
	}

	ocrcommon.MaybeEnqueueEnhancedTelem(ds.jb, ds.chEnhancedTelem, ocrcommon.EnhancedTelemetryMercuryData{
		V1Observation:  &obs,
		TaskRunResults: trrs,
		RepTimestamp:   repts,
		FeedVersion:    mercuryutils.REPORT_V1,
	})

	return obs, nil
}

func toBigInt(val interface{}) (*big.Int, error) {
	dec, err := utils.ToDecimal(val)
	if err != nil {
		return nil, err
	}
	return dec.BigInt(), nil
}

type parseOutput struct {
	benchmarkPrice mercury.ObsResult[*big.Int]
	bid            mercury.ObsResult[*big.Int]
	ask            mercury.ObsResult[*big.Int]
}

// parse expects the output of observe to be three values, in the following order:
// 1. benchmark price
// 2. bid
// 3. ask
//
// returns error on parse errors: if something is the wrong type
func (ds *datasource) parse(trrs pipeline.TaskRunResults) (o parseOutput, merr error) {
	var finaltrrs []pipeline.TaskRunResult
	for _, trr := range trrs {
		// only return terminal trrs from executeRun
		if trr.IsTerminal() {
			finaltrrs = append(finaltrrs, trr)
		}
	}

	// pipeline.TaskRunResults comes ordered asc by index, this is guaranteed
	// by the pipeline executor
	if len(finaltrrs) != 3 {
		return o, fmt.Errorf("invalid number of results, expected: 3, got: %d", len(finaltrrs))
	}
	merr = errors.Join(
		setBenchmarkPrice(&o, finaltrrs[0].Result),
		setBid(&o, finaltrrs[1].Result),
		setAsk(&o, finaltrrs[2].Result),
	)

	return o, merr
}

func setBenchmarkPrice(o *parseOutput, res pipeline.Result) error {
	if res.Error != nil {
		o.benchmarkPrice.Err = res.Error
	} else if val, err := toBigInt(res.Value); err != nil {
		return fmt.Errorf("failed to parse BenchmarkPrice: %w", err)
	} else {
		o.benchmarkPrice.Val = val
	}
	return nil
}

func setBid(o *parseOutput, res pipeline.Result) error {
	if res.Error != nil {
		o.bid.Err = res.Error
	} else if val, err := toBigInt(res.Value); err != nil {
		return fmt.Errorf("failed to parse Bid: %w", err)
	} else {
		o.bid.Val = val
	}
	return nil
}

func setAsk(o *parseOutput, res pipeline.Result) error {
	if res.Error != nil {
		o.ask.Err = res.Error
	} else if val, err := toBigInt(res.Value); err != nil {
		return fmt.Errorf("failed to parse Ask: %w", err)
	} else {
		o.ask.Val = val
	}
	return nil
}

// The context passed in here has a timeout of (ObservationTimeout + ObservationGracePeriod).
// Upon context cancellation, its expected that we return any usable values within ObservationGracePeriod.
func (ds *datasource) executeRun(ctx context.Context) (*pipeline.Run, pipeline.TaskRunResults, error) {
	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jb": map[string]interface{}{
			"databaseID":    ds.jb.ID,
			"externalJobID": ds.jb.ExternalJobID,
			"name":          ds.jb.Name.ValueOrZero(),
		},
	})

	run, trrs, err := ds.pipelineRunner.ExecuteRun(ctx, ds.spec, vars, ds.lggr)
	if err != nil {
		return nil, nil, pkgerrors.Wrapf(err, "error executing run for spec ID %v", ds.spec.ID)
	}

	return run, trrs, err
}

func (ds *datasource) setLatestBlocks(ctx context.Context, obs *v1types.Observation) error {
	latestBlocks, err := ds.mercuryChainReader.LatestHeads(ctx, nBlocksObservation)

	if err != nil {
		ds.lggr.Errorw("failed to read latest blocks", "err", err)
		return err
	}

	if len(latestBlocks) < nBlocksObservation {
		ds.insufficientBlocksCounter.Inc()
		ds.lggr.Warnw("Insufficient blocks", "latestBlocks", latestBlocks, "lenLatestBlocks", len(latestBlocks), "nBlocksObservation", nBlocksObservation)
	}

	// TODO: remove with https://smartcontract-it.atlassian.net/browse/BCF-2209
	if len(latestBlocks) == 0 {
		obsErr := fmt.Errorf("no blocks available")
		ds.zeroBlocksCounter.Inc()
		obs.CurrentBlockNum.Err = obsErr
		obs.CurrentBlockHash.Err = obsErr
		obs.CurrentBlockTimestamp.Err = obsErr
	} else {
		obs.CurrentBlockNum.Val = int64(latestBlocks[0].Number)
		obs.CurrentBlockHash.Val = latestBlocks[0].Hash
		obs.CurrentBlockTimestamp.Val = latestBlocks[0].Timestamp
	}

	for _, block := range latestBlocks {
		obs.LatestBlocks = append(
			obs.LatestBlocks,
			v1types.NewBlock(int64(block.Number), block.Hash, block.Timestamp))
	}

	return nil
}
