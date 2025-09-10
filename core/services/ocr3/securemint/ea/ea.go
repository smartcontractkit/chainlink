package ea

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core/securemint"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	sm_config "github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

// externalAdapter implements securemint.ExternalAdapter
var _ securemint.ExternalAdapter = &externalAdapter{}

type externalAdapter struct {
	config         *sm_config.SecureMintConfig
	chainSelectors []uint64 // use parsed chain selectors from config
	runner         pipeline.Runner
	job            job.Job
	spec           pipeline.Spec
	saver          ocrcommon.Saver
	lggr           logger.Logger
}

func NewExternalAdapter(config *sm_config.SecureMintConfig, runner pipeline.Runner, job job.Job, spec pipeline.Spec, saver ocrcommon.Saver, lggr logger.Logger) (*externalAdapter, error) {
	chainSelectors := make([]uint64, 0, len(config.ChainSelectors))
	for _, chainSelector := range config.ChainSelectors {
		chainSelectorUint64, err := strconv.ParseUint(chainSelector, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse chain selector: %s", chainSelector)
		}
		chainSelectors = append(chainSelectors, chainSelectorUint64)
	}

	return &externalAdapter{config: config, chainSelectors: chainSelectors, runner: runner, job: job, spec: spec, saver: saver, lggr: lggr}, nil
}

// GetPayload retrieves the payload for the given blocks by executing a pipeline run.
func (ea *externalAdapter) GetPayload(ctx context.Context, blocks securemint.Blocks) (securemint.ExternalAdapterPayload, error) {
	ea.lggr.Debugf("GetPayload called with blocks parameter: %v", blocks)

	// Create the request for the external adapter
	req := Request{
		Token:    ea.config.Token,
		Reserves: ea.config.Reserves,
	}

	// coalesce blocks with config.ChainSelectors
	coalescedBlocks := make(map[uint64]uint64)
	for _, chainSelector := range ea.chainSelectors {
		coalescedBlocks[chainSelector] = uint64(blocks[securemint.ChainSelector(chainSelector)])
	}

	// add coalesced blocks to request
	for chainSelector, blockNumber := range coalescedBlocks {
		req.SupplyChains = append(req.SupplyChains, strconv.FormatUint(chainSelector, 10))
		req.SupplyChainBlocks = append(req.SupplyChainBlocks, blockNumber)
	}

	// Serialize EA request to JSON
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return securemint.ExternalAdapterPayload{}, fmt.Errorf("failed to marshal ea request: %w (request: %#v)", err, req)
	}

	ea.lggr.Debugf("GetPayload serialized ea request to JSON: %v", string(reqJSON))

	// Execute the request
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
		return securemint.ExternalAdapterPayload{}, fmt.Errorf("failed to execute GetPayload: %w", err)
	}

	ea.saver.Save(run)

	// Parse and return results
	for _, trr := range trrs {
		if !trr.IsTerminal() {
			continue
		}

		resultMap, ok := trr.Result.Value.(map[string]any)
		if !ok {
			return securemint.ExternalAdapterPayload{}, fmt.Errorf("unexpected result type for GetPayload: %T", trr.Result.Value)
		}

		payload, err := ea.convertMapToPayload(resultMap)
		if err != nil {
			return securemint.ExternalAdapterPayload{}, fmt.Errorf("failed to convert EA response map to payload: %w, map: %#v", err, resultMap)
		}

		ea.lggr.Debugw("GetPayload result", "payload", payload)
		if len(blocks) == 0 {
			ea.lggr.Debugw("Plugin does not know about any chains or blocks yet, not returning any mintables")
			// set Mintables to empty map - plugin will error out if it's not empty when it hasn't requested any mintables yet
			// NB: this will be fixed in v0.5 of the plugin.
			payload.Mintables = make(securemint.Mintables)
		}
		ea.lggr.Debugw("GetPayload returning", "payload", payload)

		return payload, nil
	}

	return securemint.ExternalAdapterPayload{}, errors.New("no terminal result for GetPayload")
}

// convertMapToPayload converts a map[string]any response to securemint.ExternalAdapterPayload
func (ea *externalAdapter) convertMapToPayload(resultMap map[string]any) (securemint.ExternalAdapterPayload, error) {
	// Marshal and unmarshal to convert to Response struct
	b, err := json.Marshal(resultMap)
	if err != nil {
		return securemint.ExternalAdapterPayload{}, fmt.Errorf("failed to marshal EA payload map: %w", err)
	}

	ea.lggr.Debugf("EA response: %s", string(b))

	var eaResponse Response
	if err := json.Unmarshal(b, &eaResponse); err != nil {
		return securemint.ExternalAdapterPayload{}, fmt.Errorf("failed to unmarshal EA response: %w", err)
	}

	// Create the payload
	payload := securemint.ExternalAdapterPayload{
		Mintables:    make(securemint.Mintables),
		LatestBlocks: make(securemint.Blocks),
	}

	// Convert mintables
	for chainSelector, mintable := range eaResponse.Mintables {
		chainSelectorUint64, err := strconv.ParseUint(chainSelector, 10, 64)
		if err != nil {
			return securemint.ExternalAdapterPayload{}, fmt.Errorf("failed to parse chain selector: %s", chainSelector)
		}

		mintableAmount, ok := new(big.Int).SetString(mintable.Mintable, 10)
		if !ok {
			return securemint.ExternalAdapterPayload{}, fmt.Errorf("failed to parse mintable amount: %s", mintable.Mintable)
		}

		payload.Mintables[securemint.ChainSelector(chainSelectorUint64)] = securemint.BlockMintablePair{
			Block:    securemint.BlockNumber(mintable.Block),
			Mintable: mintableAmount,
		}
	}

	// Convert reserve info
	reserveAmount, ok := new(big.Int).SetString(eaResponse.ReserveInfo.ReserveAmount, 10)
	if !ok {
		return securemint.ExternalAdapterPayload{}, fmt.Errorf("failed to parse reserve amount: %s", eaResponse.ReserveInfo.ReserveAmount)
	}
	payload.ReserveInfo = securemint.ReserveInfo{
		ReserveAmount: reserveAmount,
		Timestamp:     time.UnixMilli(eaResponse.ReserveInfo.Timestamp),
	}

	// Convert latest blocks
	for chainSelector, block := range eaResponse.LatestBlocks {
		chainSelectorUint64, err := strconv.ParseUint(chainSelector, 10, 64)
		if err != nil {
			return securemint.ExternalAdapterPayload{}, fmt.Errorf("failed to parse chain selector: %s", chainSelector)
		}
		payload.LatestBlocks[securemint.ChainSelector(chainSelectorUint64)] = securemint.BlockNumber(block)
	}

	return payload, nil
}
