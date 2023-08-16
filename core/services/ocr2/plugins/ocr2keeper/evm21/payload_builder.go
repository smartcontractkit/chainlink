package evm

import (
	"context"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

type payloadBuilder struct {
	lggr logger.Logger
}

var _ ocr2keepers.PayloadBuilder = &payloadBuilder{}

func NewPayloadBuilder(lggr logger.Logger) *payloadBuilder {
	return &payloadBuilder{
		lggr: lggr,
	}
}

func (b *payloadBuilder) BuildPayloads(ctx context.Context, proposals ...ocr2keepers.CoordinatedBlockProposal) ([]ocr2keepers.UpkeepPayload, error) {
	payloads := make([]ocr2keepers.UpkeepPayload, len(proposals))
	for i, p := range proposals {
		b.lggr.Debugf("building payload for coordinated block proposal %+v", p)
		var checkData []byte
		switch core.GetUpkeepType(p.UpkeepID) {
		case ocr2keepers.LogTrigger:
			checkData = []byte{} // TODO: call recoverer
		case ocr2keepers.ConditionTrigger:
			// Trigger.BlockNumber and Trigger.BlockHash are already coordinated
			checkData = []byte{} // CheckData derived on chain for conditionals
			// TODO: check for upkeepID being active upkeep here
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
