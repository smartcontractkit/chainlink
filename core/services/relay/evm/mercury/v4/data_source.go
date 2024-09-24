package v4

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	pkgerrors "github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	v4types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v4"
	v4 "github.com/smartcontractkit/chainlink-data-streams/mercury/v4"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
	mercurytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v4/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Runner interface {
	ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars) (run *pipeline.Run, trrs pipeline.TaskRunResults, err error)
}

type LatestReportFetcher interface {
	LatestPrice(ctx context.Context, feedID [32]byte) (*big.Int, error)
	LatestTimestamp(context.Context) (int64, error)
}

type datasource struct {
	pipelineRunner Runner
	jb             job.Job
	spec           pipeline.Spec
	feedID         mercuryutils.FeedID
	lggr           logger.Logger
	saver          ocrcommon.Saver
	orm            types.DataSourceORM
	codec          reportcodec.ReportCodec

	fetcher      LatestReportFetcher
	linkFeedID   mercuryutils.FeedID
	nativeFeedID mercuryutils.FeedID

	mu sync.RWMutex

	chEnhancedTelem chan<- ocrcommon.EnhancedTelemetryMercuryData
}

var _ v4.DataSource = &datasource{}

func NewDataSource(orm types.DataSourceORM, pr pipeline.Runner, jb job.Job, spec pipeline.Spec, feedID mercuryutils.FeedID, lggr logger.Logger, s ocrcommon.Saver, enhancedTelemChan chan ocrcommon.EnhancedTelemetryMercuryData, fetcher LatestReportFetcher, linkFeedID, nativeFeedID mercuryutils.FeedID) *datasource {
	return &datasource{pr, jb, spec, feedID, lggr, s, orm, reportcodec.ReportCodec{}, fetcher, linkFeedID, nativeFeedID, sync.RWMutex{}, enhancedTelemChan}
}

func (ds *datasource) Observe(ctx context.Context, repts ocrtypes.ReportTimestamp, fetchMaxFinalizedTimestamp bool) (obs v4types.Observation, pipelineExecutionErr error) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)

	if fetchMaxFinalizedTimestamp {
		wg.Add(1)
		go func() {
			defer wg.Done()
			latest, dbErr := ds.orm.LatestReport(ctx, ds.feedID)
			if dbErr != nil {
				obs.MaxFinalizedTimestamp.Err = dbErr
				return
			}
			if latest != nil {
				maxFinalizedBlockNumber, decodeErr := ds.codec.ObservationTimestampFromReport(latest)
				obs.MaxFinalizedTimestamp.Val, obs.MaxFinalizedTimestamp.Err = int64(maxFinalizedBlockNumber), decodeErr
				return
			}
			obs.MaxFinalizedTimestamp.Val, obs.MaxFinalizedTimestamp.Err = ds.fetcher.LatestTimestamp(ctx)
		}()
	}

	var trrs pipeline.TaskRunResults
	wg.Add(1)
	go func() {
		defer wg.Done()
		var run *pipeline.Run
		run, trrs, pipelineExecutionErr = ds.executeRun(ctx)
		if pipelineExecutionErr != nil {
			cancel()
			pipelineExecutionErr = fmt.Errorf("Observe failed while executing run: %w", pipelineExecutionErr)
			return
		}

		ds.saver.Save(run)

		var parsed parseOutput
		parsed, pipelineExecutionErr = ds.parse(trrs)
		if pipelineExecutionErr != nil {
			cancel()
			// This is not expected under normal circumstances
			ds.lggr.Errorw("Observe failed while parsing run results", "err", pipelineExecutionErr)
			pipelineExecutionErr = fmt.Errorf("Observe failed while parsing run results: %w", pipelineExecutionErr)
			return
		}
		obs.BenchmarkPrice = parsed.benchmarkPrice
		obs.MarketStatus = parsed.marketStatus
	}()

	var isLink, isNative bool
	if len(ds.jb.OCR2OracleSpec.PluginConfig) == 0 {
		obs.LinkPrice.Val = v4.MissingPrice
	} else if ds.feedID == ds.linkFeedID {
		isLink = true
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			obs.LinkPrice.Val, obs.LinkPrice.Err = ds.fetcher.LatestPrice(ctx, ds.linkFeedID)
			if obs.LinkPrice.Val == nil && obs.LinkPrice.Err == nil {
				mercurytypes.PriceFeedMissingCount.WithLabelValues(ds.linkFeedID.String()).Inc()
				ds.lggr.Warnw(fmt.Sprintf("Mercury server was missing LINK feed, using sentinel value of %s", v4.MissingPrice), "linkFeedID", ds.linkFeedID)
				obs.LinkPrice.Val = v4.MissingPrice
			} else if obs.LinkPrice.Err != nil {
				mercurytypes.PriceFeedErrorCount.WithLabelValues(ds.linkFeedID.String()).Inc()
				ds.lggr.Errorw("Mercury server returned error querying LINK price feed", "err", obs.LinkPrice.Err, "linkFeedID", ds.linkFeedID)
			}
		}()
	}

	if len(ds.jb.OCR2OracleSpec.PluginConfig) == 0 {
		obs.NativePrice.Val = v4.MissingPrice
	} else if ds.feedID == ds.nativeFeedID {
		isNative = true
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			obs.NativePrice.Val, obs.NativePrice.Err = ds.fetcher.LatestPrice(ctx, ds.nativeFeedID)
			if obs.NativePrice.Val == nil && obs.NativePrice.Err == nil {
				mercurytypes.PriceFeedMissingCount.WithLabelValues(ds.nativeFeedID.String()).Inc()
				ds.lggr.Warnw(fmt.Sprintf("Mercury server was missing native feed, using sentinel value of %s", v4.MissingPrice), "nativeFeedID", ds.nativeFeedID)
				obs.NativePrice.Val = v4.MissingPrice
			} else if obs.NativePrice.Err != nil {
				mercurytypes.PriceFeedErrorCount.WithLabelValues(ds.nativeFeedID.String()).Inc()
				ds.lggr.Errorw("Mercury server returned error querying native price feed", "err", obs.NativePrice.Err, "nativeFeedID", ds.nativeFeedID)
			}
		}()
	}

	wg.Wait()
	cancel()

	if pipelineExecutionErr != nil {
		return
	}

	if isLink || isNative {
		// run has now completed so it is safe to use benchmark price
		if isLink {
			// This IS the LINK feed, use our observed price
			obs.LinkPrice.Val, obs.LinkPrice.Err = obs.BenchmarkPrice.Val, obs.BenchmarkPrice.Err
		}
		if isNative {
			// This IS the native feed, use our observed price
			obs.NativePrice.Val, obs.NativePrice.Err = obs.BenchmarkPrice.Val, obs.BenchmarkPrice.Err
		}
	}

	ocrcommon.MaybeEnqueueEnhancedTelem(ds.jb, ds.chEnhancedTelem, ocrcommon.EnhancedTelemetryMercuryData{
		V4Observation:              &obs,
		TaskRunResults:             trrs,
		RepTimestamp:               repts,
		FeedVersion:                mercuryutils.REPORT_V4,
		FetchMaxFinalizedTimestamp: fetchMaxFinalizedTimestamp,
		IsLinkFeed:                 isLink,
		IsNativeFeed:               isNative,
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
	marketStatus   mercury.ObsResult[uint32]
}

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
	if len(finaltrrs) != 2 {
		return o, fmt.Errorf("invalid number of results, expected: 2, got: %d", len(finaltrrs))
	}

	merr = errors.Join(
		setBenchmarkPrice(&o, finaltrrs[0].Result),
		setMarketStatus(&o, finaltrrs[1].Result),
	)

	return o, merr
}

func setBenchmarkPrice(o *parseOutput, res pipeline.Result) error {
	if res.Error != nil {
		o.benchmarkPrice.Err = res.Error
		return res.Error
	}
	val, err := toBigInt(res.Value)
	if err != nil {
		return fmt.Errorf("failed to parse BenchmarkPrice: %w", err)
	}
	o.benchmarkPrice.Val = val
	return nil
}

func setMarketStatus(o *parseOutput, res pipeline.Result) error {
	if res.Error != nil {
		o.marketStatus.Err = res.Error
		return res.Error
	}
	val, err := toBigInt(res.Value)
	if err != nil {
		return fmt.Errorf("failed to parse MarketStatus: %w", err)
	}
	o.marketStatus.Val = uint32(val.Int64())
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

	run, trrs, err := ds.pipelineRunner.ExecuteRun(ctx, ds.spec, vars)
	if err != nil {
		return nil, nil, pkgerrors.Wrapf(err, "error executing run for spec ID %v", ds.spec.ID)
	}

	return run, trrs, err
}
