package mercury_v3

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	pkgerrors "github.com/pkg/errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	relaymercuryv3 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v3"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Runner interface {
	ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger) (run pipeline.Run, trrs pipeline.TaskRunResults, err error)
}

type LatestReportFetcher interface {
	LatestPrice(ctx context.Context, feedID [32]byte) (*big.Int, error)
	LatestTimestamp(context.Context) (uint32, error)
}

type datasource struct {
	pipelineRunner Runner
	jb             job.Job
	spec           pipeline.Spec
	feedID         types.FeedID
	lggr           logger.Logger
	runResults     chan<- pipeline.Run

	fetcher      LatestReportFetcher
	linkFeedID   types.FeedID
	nativeFeedID types.FeedID

	mu sync.RWMutex

	chEnhancedTelem chan<- ocrcommon.EnhancedTelemetryMercuryData
}

var _ relaymercuryv3.DataSource = &datasource{}

var maxInt192 *big.Int = new(big.Int).Exp(big.NewInt(2), big.NewInt(191), nil)

func NewDataSource(pr pipeline.Runner, jb job.Job, spec pipeline.Spec, feedID types.FeedID, lggr logger.Logger, rr chan pipeline.Run, enhancedTelemChan chan ocrcommon.EnhancedTelemetryMercuryData, fetcher LatestReportFetcher, linkFeedID, nativeFeedID types.FeedID) *datasource {
	return &datasource{pr, jb, spec, feedID, lggr, rr, fetcher, linkFeedID, nativeFeedID, sync.RWMutex{}, enhancedTelemChan}
}

func (ds *datasource) Observe(ctx context.Context, repts ocrtypes.ReportTimestamp, fetchMaxFinalizedTimestamp bool) (obs relaymercuryv3.Observation, merr error) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		run, trrs, err := ds.executeRun(ctx)
		if err != nil {
			merr = fmt.Errorf("Observe failed while executing run: %w", err)
			obs.BenchmarkPrice.Err = err
			obs.Bid.Err = err
			obs.Ask.Err = err
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
			// This is not expected under normal circumstances
			ds.lggr.Errorw("Observe failed while parsing run results", "err", err)
			merr = fmt.Errorf("Observe failed while parsing run results: %w", err)
			return
		}
		obs.BenchmarkPrice = parsed.benchmarkPrice
		obs.Bid = parsed.bid
		obs.Ask = parsed.ask
	}()

	// wait for runner to finish
	wg.Wait()

	if fetchMaxFinalizedTimestamp {
		wg.Add(1)
		go func() {
			defer wg.Done()
			obs.MaxFinalizedTimestamp.Val, obs.MaxFinalizedTimestamp.Err = ds.fetcher.LatestTimestamp(ctx)
		}()
	}

	if ds.feedID == ds.linkFeedID {
		// This IS the LINK feed, use our observed price
		obs.LinkPrice.Val, obs.LinkPrice.Err = obs.BenchmarkPrice.Val, obs.BenchmarkPrice.Err
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			obs.LinkPrice.Val, obs.LinkPrice.Err = ds.fetcher.LatestPrice(ctx, ds.linkFeedID)
			if obs.LinkPrice.Val == nil && obs.LinkPrice.Err == nil {
				ds.lggr.Warnw("Mercury server was missing LINK feed, falling back to max int192", "linkFeedID", ds.linkFeedID)
				obs.LinkPrice.Val = maxInt192
			}
		}()
	}

	if ds.feedID == ds.nativeFeedID {
		// This IS the native feed, use our observed price
		obs.NativePrice.Val, obs.NativePrice.Err = obs.BenchmarkPrice.Val, obs.BenchmarkPrice.Err
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			obs.NativePrice.Val, obs.NativePrice.Err = ds.fetcher.LatestPrice(ctx, ds.nativeFeedID)
			if obs.NativePrice.Val == nil && obs.NativePrice.Err == nil {
				ds.lggr.Warnw("Mercury server was missing native feed, falling back to max int192", "nativeFeedID", ds.nativeFeedID)
				obs.NativePrice.Val = maxInt192
			}
		}()
	}

	wg.Wait()

	// todo: implement telemetry
	// if ocrcommon.ShouldCollectEnhancedTelemetryMercury(&ds.jb) {
	// 	ocrcommon.EnqueueEnhancedTelem(ds.chEnhancedTelem, ocrcommon.EnhancedTelemetryMercuryData{
	// 		TaskRunResults: trrs,
	// 		Observation:    obs,
	// 		RepTimestamp:   repts,
	// 	})
	// }

	return obs, merr
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
	bid            relaymercury.ObsResult[*big.Int]
	ask            relaymercury.ObsResult[*big.Int]
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
		return res.Error
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
		return res.Error
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
		return res.Error
	} else if val, err := toBigInt(res.Value); err != nil {
		return fmt.Errorf("failed to parse Ask: %w", err)
	} else {
		o.ask.Val = val
	}
	return nil
}

// The context passed in here has a timeout of (ObservationTimeout + ObservationGracePeriod).
// Upon context cancellation, its expected that we return any usable values within ObservationGracePeriod.
func (ds *datasource) executeRun(ctx context.Context) (pipeline.Run, pipeline.TaskRunResults, error) {
	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jb": map[string]interface{}{
			"databaseID":    ds.jb.ID,
			"externalJobID": ds.jb.ExternalJobID,
			"name":          ds.jb.Name.ValueOrZero(),
		},
	})

	run, trrs, err := ds.pipelineRunner.ExecuteRun(ctx, ds.spec, vars, ds.lggr)
	if err != nil {
		return pipeline.Run{}, nil, pkgerrors.Wrapf(err, "error executing run for spec ID %v", ds.spec.ID)
	}

	return run, trrs, err
}
