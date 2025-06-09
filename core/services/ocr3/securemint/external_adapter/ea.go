package external_adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// externalAdapter implements por.ExternalAdapter
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

func (ea *externalAdapter) GetPayload(ctx context.Context, blocks por.Blocks) (por.ExternalAdapterPayload, error) {
	ea.lggr.Debugf("GetPayload called with blocks: %v", blocks)

	// 		ds1 [type=bridge name="%s" timeout=0 requestData=<{"data": {"address": "0x1234"}}>]
	// ds1 [type=bridge name=\"bridge-api0\" requestData="{\\\"data\\": {\\\"from\\\":\\\"LINK\\\",\\\"to\\\":\\\"ETH\\\"}}"];
	//     submit [type=bridge name="substrate-adapter1" requestData=<{ "value": $(parse) }>]

	// {
	//     "data": {
	//         "token": "eth",
	//         "reserves": "platform",
	//         "supplyChains": [
	//             "5009297550715157269"
	//         ],
	//         "supplyChainBlocks": [
	//             0
	//         ]
	//     }
	// }

	// execute
	vars := map[string]any{
		"jb": map[string]any{
			"databaseID":    ea.job.ID,
			"externalJobID": ea.job.ExternalJobID,
			"name":          ea.job.Name.ValueOrZero(),
		},
		"action": "get_payload",
		"blocks": blocks,
	}

	run, trrs, err := ea.runner.ExecuteRun(ctx, ea.spec, pipeline.NewVarsFrom(vars))
	if err != nil {
		ea.lggr.Errorw("Error executing GetPayload", "error", err)
		return por.ExternalAdapterPayload{}, err
	}

	// save run
	ea.saver.Save(run)

	// parse and return results
	for _, trr := range trrs {
		if trr.IsTerminal() {
			if m, ok := trr.Result.Value.(por.ExternalAdapterPayload); ok {
				return m, nil
			}
			// Try to parse from map[string]interface{} using JSON marshal/unmarshal (TODO(gg): clean up if needed)
			if m, ok := trr.Result.Value.(map[string]any); ok {
				b, err := json.Marshal(m)
				if err != nil {
					return por.ExternalAdapterPayload{}, fmt.Errorf("failed to marshal EA payload map: %w", err)
				}
				var payload por.ExternalAdapterPayload
				err = json.Unmarshal(b, &payload)
				if err != nil {
					return por.ExternalAdapterPayload{}, fmt.Errorf("failed to unmarshal EA payload: %w", err)
				}
				ea.lggr.Debugw("GetPayload result", "payload", payload)
				return payload, nil
			}
			return por.ExternalAdapterPayload{}, fmt.Errorf("unexpected result type for GetPayload: %T", trr.Result.Value)
		}
	}
	return por.ExternalAdapterPayload{}, errors.New("no terminal result for GetPayload")
}
