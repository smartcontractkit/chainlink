package mercury

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	pkgerrors "github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

//go:generate mockery --quiet --name ChainHeadTracker --output ./mocks/ --case=underscore
type ChainHeadTracker interface {
	Client() evmclient.Client
	HeadTracker() httypes.HeadTracker
}

type datasource struct {
	pipelineRunner pipeline.Runner
	jb             job.Job
	spec           pipeline.Spec
	lggr           logger.Logger
	runResults     chan<- pipeline.Run

	mu sync.RWMutex

	chEnhancedTelem  chan<- ocrcommon.EnhancedTelemetryMercuryData
	chainHeadTracker ChainHeadTracker
}

var _ relaymercury.DataSource = &datasource{}

func NewDataSource(pr pipeline.Runner, jb job.Job, spec pipeline.Spec, lggr logger.Logger, rr chan pipeline.Run, enhancedTelemChan chan ocrcommon.EnhancedTelemetryMercuryData, chainHeadTracker ChainHeadTracker) *datasource {
	return &datasource{pr, jb, spec, lggr, rr, sync.RWMutex{}, enhancedTelemChan, chainHeadTracker}
}

func (ds *datasource) Observe(ctx context.Context, repts ocrtypes.ReportTimestamp) (relaymercury.Observation, error) {
	run, trrs, err := ds.executeRun(ctx)
	if err != nil {
		return relaymercury.Observation{}, fmt.Errorf("Observe failed while executing run: %w", err)
	}
	select {
	case ds.runResults <- run:
	default:
		ds.lggr.Warnf("unable to enqueue run save for job ID %d, buffer full", ds.spec.JobID)
	}

	// NOTE: trrs comes back as _all_ tasks, but we only want the terminal ones
	// They are guaranteed to be sorted by index asc so should be in the correct order
	var finaltrrs []pipeline.TaskRunResult
	for _, trr := range trrs {
		if trr.IsTerminal() {
			finaltrrs = append(finaltrrs, trr)
		}
	}

	parsed, err := ds.parse(finaltrrs)
	if err != nil {
		return relaymercury.Observation{}, fmt.Errorf("Observe failed while parsing run results: %w", err)
	}
	ds.setCurrentBlock(ctx, &parsed)

	if ocrcommon.ShouldCollectEnhancedTelemetryMercury(&ds.jb) {
		ocrcommon.EnqueueEnhancedTelem(ds.chEnhancedTelem, ocrcommon.EnhancedTelemetryMercuryData{
			TaskRunResults: trrs,
			Observation:    parsed,
			RepTimestamp:   repts,
		})

	}

	return parsed, nil
}

func toBigInt(val interface{}) (*big.Int, error) {
	dec, err := utils.ToDecimal(val)
	if err != nil {
		return nil, err
	}
	return dec.BigInt(), nil
}

// parse expects the output of observe to be three values, in the following order:
// 1. benchmark price
// 2. bid
// 3. ask
//
// returns error on parse errors: if something is the wrong type
func (ds *datasource) parse(trrs pipeline.TaskRunResults) (obs relaymercury.Observation, merr error) {
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
		return obs, fmt.Errorf("invalid number of results, expected: 3, got: %d", len(finaltrrs))
	}
	merr = errors.Join(
		setBenchmarkPrice(&obs, finaltrrs[0].Result),
		setBid(&obs, finaltrrs[1].Result),
		setAsk(&obs, finaltrrs[2].Result),
	)

	return obs, merr
}

func setBenchmarkPrice(obs *relaymercury.Observation, res pipeline.Result) error {
	if res.Error != nil {
		obs.BenchmarkPrice.Err = res.Error
	} else if val, err := toBigInt(res.Value); err != nil {
		return fmt.Errorf("failed to parse BenchmarkPrice: %w", err)
	} else {
		obs.BenchmarkPrice.Val = val
	}
	return nil
}

func setBid(obs *relaymercury.Observation, res pipeline.Result) error {
	if res.Error != nil {
		obs.Bid.Err = res.Error
	} else if val, err := toBigInt(res.Value); err != nil {
		return fmt.Errorf("failed to parse Bid: %w", err)
	} else {
		obs.Bid.Val = val
	}
	return nil
}

func setAsk(obs *relaymercury.Observation, res pipeline.Result) error {
	if res.Error != nil {
		obs.Ask.Err = res.Error
	} else if val, err := toBigInt(res.Value); err != nil {
		return fmt.Errorf("failed to parse Ask: %w", err)
	} else {
		obs.Ask.Val = val
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

func (ds *datasource) setCurrentBlock(ctx context.Context, obs *relaymercury.Observation) {
	latestHead, err := ds.getCurrentBlock(ctx)
	if err != nil {
		obs.CurrentBlockNum.Err = err
		obs.CurrentBlockHash.Err = err
		obs.CurrentBlockTimestamp.Err = err
		return
	}
	obs.CurrentBlockNum.Val = latestHead.Number
	obs.CurrentBlockHash.Val = latestHead.Hash.Bytes()

	if latestHead.Timestamp.IsZero() {
		obs.CurrentBlockTimestamp.Val = 0
	} else {
		obs.CurrentBlockTimestamp.Val = uint64(latestHead.Timestamp.Unix())
	}
}

func (ds *datasource) getCurrentBlock(ctx context.Context) (*evmtypes.Head, error) {
	// Use the headtracker's view of the latest block, this is very fast since
	// it doesn't make any external network requests, and it is the
	// headtracker's job to ensure it has an up-to-date view of the chain based
	// on responses from all available RPC nodes
	latestHead := ds.chainHeadTracker.HeadTracker().LatestChain()
	if latestHead == nil {
		logger.Sugared(ds.lggr).AssumptionViolation("HeadTracker unexpectedly returned nil head, falling back to RPC call")
		var err error
		latestHead, err = ds.chainHeadTracker.Client().HeadByNumber(ctx, nil)
		if err != nil {
			return nil, err
		}
	}
	return latestHead, nil
}
