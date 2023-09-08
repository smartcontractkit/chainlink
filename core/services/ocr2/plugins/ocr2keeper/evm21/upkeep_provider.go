package evm

import (
	"context"
	"fmt"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

var _ ocr2keepers.ConditionalUpkeepProvider = &upkeepProvider{}

type upkeepProvider struct {
	activeUpkeeps ActiveUpkeepList
	bs            *BlockSubscriber
	lp            logpoller.LogPoller
}

func NewUpkeepProvider(activeUpkeeps ActiveUpkeepList, bs *BlockSubscriber, lp logpoller.LogPoller) *upkeepProvider {
	return &upkeepProvider{
		activeUpkeeps: activeUpkeeps,
		bs:            bs,
		lp:            lp,
	}
}

func (p *upkeepProvider) GetActiveUpkeeps(_ context.Context) ([]ocr2keepers.UpkeepPayload, error) {
	latestBlock := p.bs.latestBlock.Load()
	if latestBlock == nil {
		return nil, fmt.Errorf("no latest block found when fetching active upkeeps")
	}
	var payloads []ocr2keepers.UpkeepPayload
	for _, uid := range p.activeUpkeeps.View(ocr2keepers.ConditionTrigger) {
		payload, err := core.NewUpkeepPayload(
			uid,
			ocr2keepers.NewTrigger(latestBlock.Number, latestBlock.Hash),
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
