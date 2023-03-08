package mercury

import (
	"context"
	"math/big"
	sync "sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
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
	run, finalResult, err := ds.executeRun(ctx)
	if err != nil {
		return relaymercury.Observation{}, err
	}
	select {
	case ds.runResults <- run:
	default:
		ds.lggr.Warnf("unable to enqueue run save for job ID %d, buffer full", ds.spec.JobID)
	}

	return ds.parse(finalResult)
}

func toBigInt(val interface{}) (*big.Int, error) {
	dec, err := utils.ToDecimal(val)
	if err != nil {
		return nil, err
	}
	return dec.BigInt(), nil
}

// parse converts the FinalResult into a Observation and stores it in the bridge metadata
// expects the output of observe to be five values, in the following order:
// - benchmark price
// - bid
// - ask
// - current block number
// - current block hash
func (ds *datasource) parse(result pipeline.FinalResult) (obs relaymercury.Observation, merr error) {
	if result.HasErrors() {
		return obs, result.CombinedError()
	}
	vals := result.Values
	if len(vals) != 5 {
		return obs, errors.Errorf("invalid number of results, got: %s", vals)
	}
	for i := 0; i < len(vals); i++ {
		var err error
		switch i {
		case 0:
			obs.BenchmarkPrice, err = toBigInt(vals[i])
		case 1:
			obs.Bid, err = toBigInt(vals[i])
		case 2:
			obs.Ask, err = toBigInt(vals[i])
		case 3:
			if currentblocknum, is := vals[i].(int64); is {
				obs.CurrentBlockNum = currentblocknum
			} else {
				err = errors.Errorf("expected int64, got: %v", vals[i])
			}
		case 4:
			if currentblockhash, is := vals[i].(common.Hash); is {
				obs.CurrentBlockHash = currentblockhash.Bytes()
			} else {
				err = errors.Errorf("expected hash, got: %v", vals[i])
			}
		}
		merr = multierr.Combine(merr, err)
	}

	return obs, merr
}

// The context passed in here has a timeout of (ObservationTimeout + ObservationGracePeriod).
// Upon context cancellation, its expected that we return any usable values within ObservationGracePeriod.
func (ds *datasource) executeRun(ctx context.Context) (pipeline.Run, pipeline.FinalResult, error) {
	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jb": map[string]interface{}{
			"databaseID":    ds.jb.ID,
			"externalJobID": ds.jb.ExternalJobID,
			"name":          ds.jb.Name.ValueOrZero(),
		},
	})

	run, trrs, err := ds.pipelineRunner.ExecuteRun(ctx, ds.spec, vars, ds.lggr)
	if err != nil {
		return pipeline.Run{}, pipeline.FinalResult{}, errors.Wrapf(err, "error executing run for spec ID %v", ds.spec.ID)
	}
	finalResult := trrs.FinalResult(ds.lggr)

	return run, finalResult, err
}
