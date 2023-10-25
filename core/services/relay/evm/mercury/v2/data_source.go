package mercury_v2

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	pkgerrors "github.com/pkg/errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	relaymercuryv2 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v2"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
	mercurytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v2/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Runner interface {
	ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger) (run *pipeline.Run, trrs pipeline.TaskRunResults, err error)
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
	runResults     chan<- *pipeline.Run
	orm            types.DataSourceORM
	codec          reportcodec.ReportCodec

	fetcher      LatestReportFetcher
	linkFeedID   mercuryutils.FeedID
	nativeFeedID mercuryutils.FeedID

	mu sync.RWMutex

	chEnhancedTelem chan<- ocrcommon.EnhancedTelemetryMercuryData
}

var _ relaymercuryv2.DataSource = &datasource{}

func NewDataSource(orm types.DataSourceORM, pr pipeline.Runner, jb job.Job, spec pipeline.Spec, feedID mercuryutils.FeedID, lggr logger.Logger, rr chan *pipeline.Run, enhancedTelemChan chan ocrcommon.EnhancedTelemetryMercuryData, fetcher LatestReportFetcher, linkFeedID, nativeFeedID mercuryutils.FeedID) *datasource {
	return &datasource{pr, jb, spec, feedID, lggr, rr, orm, reportcodec.ReportCodec{}, fetcher, linkFeedID, nativeFeedID, sync.RWMutex{}, enhancedTelemChan}
}

func (ds *datasource) Observe(ctx context.Context, repts ocrtypes.ReportTimestamp, fetchMaxFinalizedTimestamp bool) (obs relaymercuryv2.Observation, err error) {
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		var trrs pipeline.TaskRunResults
		var run *pipeline.Run
		run, trrs, err = ds.executeRun(ctx)
		if err != nil {
			cancel()
			err = fmt.Errorf("Observe failed while executing run: %w", err)
			return
		}
		select {
		case ds.runResults <- run:
		default:
			ds.lggr.Warnf("unable to enqueue run save for job ID %d, buffer full", ds.spec.JobID)
		}

		var parsed parseOutput
		parsed, err = ds.parse(trrs)
		if err != nil {
			cancel()
			// This is not expected under normal circumstances
			ds.lggr.Errorw("Observe failed while parsing run results", "err", err)
			err = fmt.Errorf("Observe failed while parsing run results: %w", err)
			return
		}
		obs.BenchmarkPrice = parsed.benchmarkPrice
	}()

	var isLink, isNative bool
	if ds.feedID == ds.linkFeedID {
		isLink = true
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			obs.LinkPrice.Val, obs.LinkPrice.Err = ds.fetcher.LatestPrice(ctx, ds.linkFeedID)
			if obs.LinkPrice.Val == nil && obs.LinkPrice.Err == nil {
				mercurytypes.PriceFeedMissingCount.WithLabelValues(ds.linkFeedID.String()).Inc()
				ds.lggr.Warnw(fmt.Sprintf("Mercury server was missing LINK feed, using sentinel value of %s", relaymercuryv2.MissingPrice), "linkFeedID", ds.linkFeedID)
				obs.LinkPrice.Val = relaymercuryv2.MissingPrice
			} else if obs.LinkPrice.Err != nil {
				mercurytypes.PriceFeedErrorCount.WithLabelValues(ds.linkFeedID.String()).Inc()
				ds.lggr.Errorw("Mercury server returned error querying LINK price feed", "err", obs.LinkPrice.Err, "linkFeedID", ds.linkFeedID)
			}
		}()
	}

	if ds.feedID == ds.nativeFeedID {
		isNative = true
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			obs.NativePrice.Val, obs.NativePrice.Err = ds.fetcher.LatestPrice(ctx, ds.nativeFeedID)
			if obs.NativePrice.Val == nil && obs.NativePrice.Err == nil {
				mercurytypes.PriceFeedMissingCount.WithLabelValues(ds.nativeFeedID.String()).Inc()
				ds.lggr.Warnw(fmt.Sprintf("Mercury server was missing native feed, using sentinel value of %s", relaymercuryv2.MissingPrice), "nativeFeedID", ds.nativeFeedID)
				obs.NativePrice.Val = relaymercuryv2.MissingPrice
			} else if obs.NativePrice.Err != nil {
				mercurytypes.PriceFeedErrorCount.WithLabelValues(ds.nativeFeedID.String()).Inc()
				ds.lggr.Errorw("Mercury server returned error querying native price feed", "err", obs.NativePrice.Err, "nativeFeedID", ds.nativeFeedID)
			}
		}()
	}

	wg.Wait()
	cancel()

	if isLink || isNative {
		// run has now completed so it is safe to use err or benchmark price
		if err != nil {
			return
		}
		if isLink {
			// This IS the LINK feed, use our observed price
			obs.LinkPrice.Val, obs.LinkPrice.Err = obs.BenchmarkPrice.Val, obs.BenchmarkPrice.Err
		}
		if isNative {
			// This IS the native feed, use our observed price
			obs.NativePrice.Val, obs.NativePrice.Err = obs.BenchmarkPrice.Val, obs.BenchmarkPrice.Err
		}
	}

	// todo: implement telemetry - https://smartcontract-it.atlassian.net/browse/MERC-1388
	// if ocrcommon.ShouldCollectEnhancedTelemetryMercury(&ds.jb) {
	// 	ocrcommon.EnqueueEnhancedTelem(ds.chEnhancedTelem, ocrcommon.EnhancedTelemetryMercuryData{
	// 		TaskRunResults: trrs,
	// 		Observation:    obs,
	// 		RepTimestamp:   repts,
	// 	})
	// }

	return obs, err
}

func toBigInt(val interface{}) (*big.Int, error) {
	dec, err := utils.ToDecimal(val)
	if err != nil {
		return nil, err
	}
	return dec.BigInt(), nil
}

type parseOutput struct {
	benchmarkPrice relaymercury.ObsResult[*big.Int]
}

func (ds *datasource) parse(trrs pipeline.TaskRunResults) (o parseOutput, merr error) {
	var finaltrrs []pipeline.TaskRunResult
	for _, trr := range trrs {
		// only return terminal trrs from executeRun
		if trr.IsTerminal() {
			finaltrrs = append(finaltrrs, trr)
		}
	}

	if len(finaltrrs) != 1 {
		return o, fmt.Errorf("invalid number of results, expected: 1, got: %d", len(finaltrrs))
	}

	return o, setBenchmarkPrice(&o, finaltrrs[0].Result)
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
