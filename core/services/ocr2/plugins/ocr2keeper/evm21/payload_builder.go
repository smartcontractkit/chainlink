package evm

import (
	"context"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"
)

type payloadBuilder struct {
	upkeepList ActiveUpkeepList
	lggr       logger.Logger
	recoverer  logprovider.LogRecoverer
}

var _ ocr2keepers.PayloadBuilder = &payloadBuilder{}

func NewPayloadBuilder(activeUpkeepList ActiveUpkeepList, recoverer logprovider.LogRecoverer, lggr logger.Logger) *payloadBuilder {
	return &payloadBuilder{
		upkeepList: activeUpkeepList,
		lggr:       lggr,
		recoverer:  recoverer,
	}
}

func (b *payloadBuilder) BuildPayloads(ctx context.Context, proposals ...ocr2keepers.CoordinatedBlockProposal) ([]ocr2keepers.UpkeepPayload, error) {
	payloads := make([]ocr2keepers.UpkeepPayload, 0, len(proposals))
	for _, proposal := range proposals {
		var payload ocr2keepers.UpkeepPayload
		var err error
		if b.upkeepList.IsActive(proposal.UpkeepID.BigInt()) {
			b.lggr.Debugf("building payload for coordinated block proposal %+v", proposal)
			switch core.GetUpkeepType(proposal.UpkeepID) {
			case ocr2keepers.LogTrigger:
				payload, err = b.BuildPayload(ctx, proposal)
			case ocr2keepers.ConditionTrigger:
				// Trigger.BlockNumber and Trigger.BlockHash are already coordinated
				// TODO: check for upkeepID being active upkeep here using b.active
			}
			if err != nil {
				b.lggr.Warnw("failed to build payload", "err", err, "upkeepID", proposal.UpkeepID)
			}
		} else {
			b.lggr.Warnw("upkeep is not active, skipping", "upkeepID", proposal.UpkeepID)
		}

		if !payload.IsEmpty() {
			payloads = append(payloads, payload)
		}
	}

	return payloads, nil
}

func (b *payloadBuilder) BuildPayload(ctx context.Context, proposal ocr2keepers.CoordinatedBlockProposal) (ocr2keepers.UpkeepPayload, error) {
	return b.recoverer.BuildPayload(ctx, proposal)
}
