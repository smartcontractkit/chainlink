package evm

import (
	"context"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"
)

type payloadBuilder struct {
	active    ActiveUpkeepList
	lggr      logger.Logger
	recoverer logprovider.LogRecoverer
}

var _ ocr2keepers.PayloadBuilder = &payloadBuilder{}

func NewPayloadBuilder(al ActiveUpkeepList, recoverer logprovider.LogRecoverer, lggr logger.Logger) *payloadBuilder {
	return &payloadBuilder{
		active:    al,
		lggr:      lggr,
		recoverer: recoverer,
	}
}

func (b *payloadBuilder) BuildPayloads(ctx context.Context, proposals ...ocr2keepers.CoordinatedBlockProposal) ([]ocr2keepers.UpkeepPayload, error) {
	payloads := make([]ocr2keepers.UpkeepPayload, len(proposals))
	for i, p := range proposals {
		var payload ocr2keepers.UpkeepPayload
		var err error
		if b.active.IsActive(p.UpkeepID.BigInt()) {
			b.lggr.Debugf("building payload for coordinated block proposal %+v", p)
			switch core.GetUpkeepType(p.UpkeepID) {
			case ocr2keepers.LogTrigger:
				payload, err = b.BuildPayload(ctx, p)
			case ocr2keepers.ConditionTrigger:
				// Trigger.BlockNumber and Trigger.BlockHash are already coordinated
				// TODO: check for upkeepID being active upkeep here using b.active
			default:
			}
			if err != nil {
				b.lggr.Warnw("failed to build payload", "err", err, "upkeepID", p.UpkeepID)
			}
		} else {
			b.lggr.Warnw("upkeep is not active, skipping", "upkeepID", p.UpkeepID)
		}
		payloads[i] = payload
	}

	return payloads, nil
}

func (b *payloadBuilder) BuildPayload(ctx context.Context, proposal ocr2keepers.CoordinatedBlockProposal) (ocr2keepers.UpkeepPayload, error) {
	return b.recoverer.BuildPayload(ctx, proposal)
}
