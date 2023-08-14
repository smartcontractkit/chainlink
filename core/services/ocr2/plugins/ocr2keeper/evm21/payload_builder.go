package evm

import (
	"context"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

type payloadBuilder struct {
	lggr logger.Logger
	//cl   client.Client
}

var _ ocr2keepers.PayloadBuilder = &payloadBuilder{}

func NewPayloadBuilder(lggr logger.Logger) *payloadBuilder {
	return &payloadBuilder{
		lggr: lggr,
		//cl:   cl,
	}
}

func (b *payloadBuilder) BuildPayloads(ctx context.Context, proposals ...ocr2keepers.CoordinatedBlockProposal) ([]ocr2keepers.UpkeepPayload, error) {
	payloads := make([]ocr2keepers.UpkeepPayload, len(proposals))
	//block, err := b.latestBlock(ctx)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get latest block: %w", err)
	//}

	for i, p := range proposals {
		var checkData []byte
		switch core.GetUpkeepType(p.UpkeepID) {
		case ocr2keepers.LogTrigger:
			checkData = []byte{} // TODO: call recoverer
		case ocr2keepers.ConditionTrigger:
			checkData = []byte{} // CheckData derived on chain for conditionals
			// Trigger.BlockNumber and Trigger.BlockHash are already coordinated
			// TODO: check for upkeepID being active upkeep here
			//p.Trigger.BlockNumber = block.Number
			//p.Trigger.BlockHash = block.Hash
		default:
		}
		payload, err := core.NewUpkeepPayload(p.UpkeepID.BigInt(), p.Trigger, checkData)
		if err != nil {
			b.lggr.Warnw("failed to build payload", "err", err, "upkeepID", p.UpkeepID)
			payloads[i] = ocr2keepers.UpkeepPayload{}
			continue
		}
		payloads[i] = payload
	}

	return payloads, nil
}

/*
func (b *payloadBuilder) latestBlock(ctx context.Context) (ocr2keepers.BlockKey, error) {
	latest, err := b.cl.LatestBlockHeight(ctx)
	if err != nil {
		return ocr2keepers.BlockKey{}, err
	}
	block, err := b.cl.BlockByNumber(ctx, latest)
	if err != nil {
		return ocr2keepers.BlockKey{}, err
	}
	return ocr2keepers.BlockKey{
		Number: ocr2keepers.BlockNumber(latest.Uint64()),
		Hash:   block.Hash(),
	}, nil
}*/
