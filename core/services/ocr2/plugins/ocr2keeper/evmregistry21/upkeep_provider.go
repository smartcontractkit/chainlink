package evm

import (
	"context"
	"math/big"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	keepersflows "github.com/smartcontractkit/ocr2keepers/pkg/v3/flows"
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
	ids, err := p.reg.GetActiveUpkeepIDsByType(ctx, uint8(conditionTrigger))
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
		payloads = append(payloads, ocr2keepers.NewUpkeepPayload(
			big.NewInt(0).SetBytes(uid),
			int(conditionTrigger),
			blockKey,
			ocr2keepers.NewTrigger(block.Int64(), blockHash.Hex(), struct{}{}),
			nil,
		))
	}

	return payloads, nil
}
