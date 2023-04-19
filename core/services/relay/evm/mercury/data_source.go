package mercury

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	sync "sync"

	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type datasource struct {
	pipelineRunner pipeline.Runner
	jb             job.Job
	spec           pipeline.Spec
	lggr           logger.Logger
	runResults     chan<- pipeline.Run

	mu sync.RWMutex
}

var _ relaymercury.DataSource = &datasource{}

func NewDataSource(pr pipeline.Runner, jb job.Job, spec pipeline.Spec, lggr logger.Logger, rr chan pipeline.Run) *datasource {
	return &datasource{pr, jb, spec, lggr, rr, sync.RWMutex{}}
}

func (ds *datasource) Observe(ctx context.Context) (relaymercury.Observation, error) {
	run, trrs, err := ds.executeRun(ctx)
	if err != nil {
		return relaymercury.Observation{}, fmt.Errorf("Observe failed while executing run: %w", err)
	}
	select {
	case ds.runResults <- run:
	default:
		ds.lggr.Warnf("unable to enqueue run save for job ID %d, buffer full", ds.spec.JobID)
	}

	return ds.parse(trrs)
}

func toBigInt(val interface{}) (*big.Int, error) {
	dec, err := utils.ToDecimal(val)
	if err != nil {
		return nil, err
	}
	return dec.BigInt(), nil
}

// parse expects the output of observe to be five values, in the following order:
// 1. benchmark price
// 2. bid
// 3. ask
// 4. current block number
// 5. current block hash
//
// returns error on parse errors: if something is the wrong type
func (ds *datasource) parse(trrs pipeline.TaskRunResults) (obs relaymercury.Observation, merr error) {
	// pipeline.TaskRunResults comes ordered asc by index, this is guaranteed
	// by the pipeline executor
	if len(trrs) != 5 {
		return obs, fmt.Errorf("invalid number of results, expected: 5, got: %d", len(trrs))
	}
	merr = errors.Join(
		setBenchmarkPrice(&obs, trrs[0].Result),
		setBid(&obs, trrs[1].Result),
		setAsk(&obs, trrs[2].Result),
		setCurrentBlockNum(&obs, trrs[3].Result),
		setCurrentBlockHash(&obs, trrs[4].Result),
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

func setCurrentBlockNum(obs *relaymercury.Observation, res pipeline.Result) error {
	if res.Error != nil {
		obs.CurrentBlockNum.Err = res.Error
	} else if val, is := res.Value.(int64); !is {
		return fmt.Errorf("failed to parse CurrentBlockNum: expected int64, got: %T (%v)", res.Value, res.Value)
	} else {
		obs.CurrentBlockNum.Val = val
	}
	return nil
}

func setCurrentBlockHash(obs *relaymercury.Observation, res pipeline.Result) error {
	if res.Error != nil {
		obs.CurrentBlockHash.Err = res.Error
	} else if val, is := res.Value.(common.Hash); !is {
		return fmt.Errorf("failed to parse CurrentBlockHash: expected hash, got: %T (%v)", res.Value, res.Value)
	} else {
		obs.CurrentBlockHash.Val = val.Bytes()
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

	// NOTE: trrs comes back as _all_ tasks, but we only want the terminal ones
	// They are guaranteed to be sorted by index asc so should be in the correct order
	run, trrs, err := ds.pipelineRunner.ExecuteRun(ctx, ds.spec, vars, ds.lggr)
	if err != nil {
		return pipeline.Run{}, nil, pkgerrors.Wrapf(err, "error executing run for spec ID %v", ds.spec.ID)
	}
	var finaltrrs []pipeline.TaskRunResult
	for _, trr := range trrs {
		// only return terminal trrs from executeRun
		if trr.IsTerminal() {
			finaltrrs = append(finaltrrs, trr)
		}
	}

	return run, finaltrrs, err
}
