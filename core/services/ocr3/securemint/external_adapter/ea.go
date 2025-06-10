package external_adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

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

// Ensure externalAdapter implements por.ExternalAdapter
func (ea *externalAdapter) GetChains(ctx context.Context) ([]por.ChainSelector, error) {
	// TODO(gg): remove this when it's removed from the plugin's adapter

	ea.lggr.Warnf("GetChains not implemented yet, returning mock data")
	chains := []por.ChainSelector{
		8953668971247136127, // "bitcoin-testnet-rootstock"
		729797994450396300,  // "telos-evm-testnet"
	}
	return chains, nil
}

// GetPayload retrieves the payload for the given blocks by executing a pipeline run.
func (ea *externalAdapter) GetPayload(ctx context.Context, blocks por.Blocks) (por.ExternalAdapterPayload, error) {
	ea.lggr.Debugf("GetPayload called with blocks: %v", blocks)

	req := EARequest{
		Token:    "eth",
		Reserves: "platform",
	}

	for chainSelector, blockNumber := range blocks {
		req.SupplyChains = append(req.SupplyChains, fmt.Sprintf("%d", chainSelector))
		req.SupplyChainBlocks = append(req.SupplyChainBlocks, uint64(blockNumber))
	}

	// Serialize EA request as JSON string
	reqJSON, err := json.Marshal(req)
	if err != nil {
		ea.lggr.Errorw("Error marshaling ea request to JSON", "error", err, "request", req)
		return por.ExternalAdapterPayload{}, fmt.Errorf("failed to marshal ea request: %w", err)
	}

	ea.lggr.Debugf("GetPayload serialized blocks to JSON: %v", string(reqJSON))

	// execute
	vars := map[string]any{
		"jb": map[string]any{
			"databaseID":    ea.job.ID,
			"externalJobID": ea.job.ExternalJobID,
			"name":          ea.job.Name.ValueOrZero(),
		},
		"action":     "get_payload",
		"ea_request": req,
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

			// TODO(gg): clean up, depends also on EA and plugin types

			if m, ok := trr.Result.Value.(map[string]any); ok {
				ea.lggr.Debugw("GetPayload result as map", "result", m)
				b, err := json.Marshal(m)
				if err != nil {
					return por.ExternalAdapterPayload{}, fmt.Errorf("failed to marshal EA payload map: %w", err)
				}

				ea.lggr.Debugw("GetPayload result as map marshaled to JSON", "json", string(b))

				var eaResp EAResponse
				err = json.Unmarshal(b, &eaResp)
				if err != nil {
					return por.ExternalAdapterPayload{}, fmt.Errorf("failed to unmarshal EA response: %w", err)
				}

				// Convert eaResponse to por.ExternalAdapterPayload
				payload := por.ExternalAdapterPayload{
					Mintables:            make(por.Mintables),
					ReserveInfo:          por.ReserveInfo{},
					LatestRelevantBlocks: make(por.Blocks),
				}

				for chainSelector, mintable := range eaResp.Mintables {
					blockMintablePair := por.BlockMintablePair{
						Block:    por.BlockNumber(mintable.Block),
						Mintable: new(big.Int),
					}
					blockMintablePair.Mintable, ok = big.NewInt(0).SetString(mintable.Mintable, 10)
					if !ok {
						return por.ExternalAdapterPayload{}, fmt.Errorf("failed to parse mintable amount: %s", mintable.Mintable)
					}
					chainSelectorUint64, err := strconv.ParseUint(chainSelector, 10, 64)
					if err != nil {
						return por.ExternalAdapterPayload{}, fmt.Errorf("failed to parse chain selector: %s", chainSelector)
					}
					payload.Mintables[por.ChainSelector(chainSelectorUint64)] = blockMintablePair
				}
				payload.ReserveInfo = por.ReserveInfo{
					ReserveAmount: new(big.Int),
				}
				payload.ReserveInfo.ReserveAmount, ok = big.NewInt(0).SetString(eaResp.ReserveInfo.ReserveAmount, 10)
				if !ok {
					return por.ExternalAdapterPayload{}, fmt.Errorf("failed to parse reserve amount: %s", eaResp.ReserveInfo.ReserveAmount)
				}
				payload.ReserveInfo.Timestamp = time.UnixMilli(eaResp.ReserveInfo.Timestamp)
				for chainSelector, block := range eaResp.LatestRelevantBlocks {
					chainSelectorUint64, err := strconv.ParseUint(chainSelector, 10, 64)
					if err != nil {
						return por.ExternalAdapterPayload{}, fmt.Errorf("failed to parse chain selector: %s", chainSelector)
					}

					payload.LatestRelevantBlocks[por.ChainSelector(chainSelectorUint64)] = por.BlockNumber(block)
				}

				ea.lggr.Debugw("GetPayload result", "payload", payload)
				return payload, nil
			}
			return por.ExternalAdapterPayload{}, fmt.Errorf("unexpected result type for GetPayload: %T", trr.Result.Value)
		}
	}
	return por.ExternalAdapterPayload{}, errors.New("no terminal result for GetPayload")
}
