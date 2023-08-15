package evm

import (
	"context"
	"math/big"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var _ ocr2keepers.ConditionalUpkeepProvider = &upkeepProvider{}

type upkeepProvider struct {
	activeUpkeeps ActiveUpkeepList
	reg           *EvmRegistry
	lp            logpoller.LogPoller
}

func NewUpkeepProvider(activeUpkeeps ActiveUpkeepList, reg *EvmRegistry, lp logpoller.LogPoller) *upkeepProvider {
	return &upkeepProvider{
		activeUpkeeps: activeUpkeeps,
		reg:           reg,
		lp:            lp,
	}
}

func (p *upkeepProvider) GetActiveUpkeeps(ctx context.Context) ([]ocr2keepers.UpkeepPayload, error) {
	latestBlock, err := p.lp.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}

	block := big.NewInt(latestBlock)
	blockHash, err := p.reg.getBlockHash(block)
	if err != nil {
		return nil, err
	}

	var payloads []ocr2keepers.UpkeepPayload
	for _, uid := range p.activeUpkeeps.View(ocr2keepers.ConditionTrigger) {
		payload, err := core.NewUpkeepPayload(
			uid,
			ocr2keepers.NewTrigger(ocr2keepers.BlockNumber(block.Int64()), blockHash),
			nil,
		)
		if err != nil {
			// skip invalid payloads
			continue
		}

		payloads = append(payloads, payload)
	}

	return payloads, nil
}
