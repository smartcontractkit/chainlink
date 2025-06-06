package external_adapter

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// externalAdapter implements por.externalAdapter
type externalAdapter struct {
	runner pipeline.Runner
	job    job.Job
	spec   pipeline.Spec
	saver  ocrcommon.Saver
	lggr   logger.Logger
	mu     sync.RWMutex
}

func NewExternalAdapter(runner pipeline.Runner, job job.Job, spec pipeline.Spec, saver ocrcommon.Saver, lggr logger.Logger) *externalAdapter {
	return &externalAdapter{runner: runner, job: job, spec: spec, saver: saver, lggr: lggr}
}

func (ea *externalAdapter) GetChains(ctx context.Context) ([]por.ChainSelector, error) {
	// TODO(gg): pass this through to the EA

	ea.lggr.Warnf("GetChains not implemented yet, returning mock data")
	chains := []por.ChainSelector{
		8953668971247136127, // "bitcoin-testnet-rootstock"
		729797994450396300,  // "telos-evm-testnet"
	}
	return chains, nil
}

func (ea *externalAdapter) GetMintables(ctx context.Context, blocks por.Blocks) (por.Mintables, error) {
	ea.lggr.Debugf("GetMintables called with blocks: %v", blocks)
	// execute
	vars := map[string]any{
		"jb": map[string]any{
			"databaseID":    ea.job.ID,
			"externalJobID": ea.job.ExternalJobID,
			"name":          ea.job.Name.ValueOrZero(),
		},
		"action": "get_mintables",
		"blocks": blocks,
	}

	run, trrs, err := ea.runner.ExecuteRun(ctx, ea.spec, pipeline.NewVarsFrom(vars))
	if err != nil {
		ea.lggr.Errorw("Error executing GetMintables", "error", err)
		return por.Mintables{}, err
	}

	// save run
	ea.saver.Save(run)

	// parse and return results
	for _, trr := range trrs {
		if trr.IsTerminal() {
			if m, ok := trr.Result.Value.(por.Mintables); ok {
				return m, nil
			}
			return por.Mintables{}, fmt.Errorf("unexpected result type for GetMintables: %T", trr.Result.Value)
		}
	}
	return por.Mintables{}, fmt.Errorf("no terminal result for GetMintables")
}

func (ea *externalAdapter) GetLatestBlocks(ctx context.Context) (por.Blocks, error) {
	// TODO(gg): pass this through to the EA

	ea.lggr.Warnf("GetLatestBlocks not implemented yet, returning mock data")
	blocks := make(por.Blocks)

	for _, chain := range []por.ChainSelector{
		8953668971247136127, // "bitcoin-testnet-rootstock"
		729797994450396300,  // "telos-evm-testnet"
	} {
		blocks[chain] = 1
	}

	return blocks, nil
}

func (ea *externalAdapter) GetReserveInfo(ctx context.Context) (*big.Int, time.Time, error) {
	// TODO(gg): pass this through to the EA

	ea.lggr.Warnf("GetReserveInfo not implemented yet, returning mock data")
	return big.NewInt(1000), time.Now(), nil

}

func (ea *externalAdapter) executeRun(ctx context.Context, extraVars map[string]any) (*pipeline.Run, pipeline.TaskRunResults, error) {
	vars := map[string]any{
		"jb": map[string]any{
			"databaseID":    ea.job.ID,
			"externalJobID": ea.job.ExternalJobID,
			"name":          ea.job.Name.ValueOrZero(),
		},
	}
	for k, v := range extraVars {
		vars[k] = v
	}
	return ea.runner.ExecuteRun(ctx, ea.spec, pipeline.NewVarsFrom(vars))
}
