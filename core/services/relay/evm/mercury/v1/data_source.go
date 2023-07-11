package mercury_v1

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	pkgerrors "github.com/pkg/errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	relaymercuryv1 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v1"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Runner interface {
	ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger) (run pipeline.Run, trrs pipeline.TaskRunResults, err error)
}

type datasource struct {
	pipelineRunner Runner
	jb             job.Job
	spec           pipeline.Spec
	lggr           logger.Logger
	runResults     chan<- pipeline.Run

	mu sync.RWMutex

	chEnhancedTelem chan<- ocrcommon.EnhancedTelemetryMercuryData
}

var _ relaymercuryv1.DataSource = &datasource{}

func NewDataSource(pr pipeline.Runner, jb job.Job, spec pipeline.Spec, lggr logger.Logger, rr chan pipeline.Run, enhancedTelemChan chan ocrcommon.EnhancedTelemetryMercuryData) *datasource {
	return &datasource{pr, jb, spec, lggr, rr, sync.RWMutex{}, enhancedTelemChan}
}

func (ds *datasource) Observe(ctx context.Context, repts ocrtypes.ReportTimestamp) (obs relaymercuryv1.Observation, err error) {
	var trrs pipeline.TaskRunResults
	var run pipeline.Run
	run, trrs, err = ds.executeRun(ctx)
	if err != nil {
		err = fmt.Errorf("Observe failed while executing run: %w", err)
		return
	}
	select {
	case ds.runResults <- run:
	default:
		ds.lggr.Warnf("unable to enqueue run save for job ID %d, buffer full", ds.spec.JobID)
	}

	var finaltrrs []pipeline.TaskRunResult
	for _, trr := range trrs {
		if trr.IsTerminal() {
			finaltrrs = append(finaltrrs, trr)
		}
	}

	var parsed parseOutput
	parsed, err = ds.parse(finaltrrs)
	if err != nil {
		err = fmt.Errorf("Observe failed while parsing run results: %w", err)
		return
	}
	obs.BenchmarkPrice = parsed.benchmarkPrice
	obs.Bid = parsed.bid
	obs.Ask = parsed.ask

	// todo: implement telemetry
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
