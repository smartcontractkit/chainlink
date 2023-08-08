package evm

import (
	"context"
	"math/big"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	keepersflows "github.com/smartcontractkit/ocr2keepers/pkg/v3/flows"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

var _ keepersflows.UpkeepProvider = &upkeepProvider{}

type upkeepProvider struct {
	reg *EvmRegistry
}

func NewUpkeepProvider(reg *EvmRegistry) *upkeepProvider {
	return &upkeepProvider{
		reg: reg,
	}
}

func (p *upkeepProvider) GetActiveUpkeeps(ctx context.Context, blockKey ocr2keepers.BlockKey) ([]ocr2keepers.UpkeepPayload, error) {
	ids, err := p.reg.GetActiveUpkeepIDsByType(ctx, uint8(ocr2keepers.ConditionTrigger))
	if err != nil {
		return nil, err
	}
	block, ok := big.NewInt(0).SetString(string(blockKey), 10)
	if !ok {
		return nil, ocr2keepers.ErrInvalidBlockKey
	}
	blockHash, err := p.reg.getBlockHash(block)
	if err != nil {
		return nil, err
	}

	var payloads []ocr2keepers.UpkeepPayload
	for _, uid := range ids {
		payload, err := core.NewUpkeepPayload(
			big.NewInt(0).SetBytes(uid),
			int(ocr2keepers.ConditionTrigger),
			ocr2keepers.NewTrigger(block.Int64(), blockHash.Hex(), struct{}{}),
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
