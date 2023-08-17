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
	payloads := make([]ocr2keepers.UpkeepPayload, len(proposals))

	for i, proposal := range proposals {
		var payload ocr2keepers.UpkeepPayload
		if b.upkeepList.IsActive(proposal.UpkeepID.BigInt()) {
			b.lggr.Debugf("building payload for coordinated block proposal %+v", proposal)
			switch core.GetUpkeepType(proposal.UpkeepID) {
			case ocr2keepers.LogTrigger:

				checkData, err := b.recoverer.GetProposalData(ctx, proposal)
				if err != nil {
					b.lggr.Warnw("failed to get proposal data", "err", err, "upkeepID", proposal.UpkeepID)
					break
				}

				payload, err = core.NewUpkeepPayload(
					proposal.UpkeepID.BigInt(),
					proposal.Trigger,
					checkData,
				)
				if err != nil {
					b.lggr.Warnw("error building upkeep payload", "err", err, "upkeepID", proposal.UpkeepID)
					break
				}

			case ocr2keepers.ConditionTrigger:
				// Trigger.BlockNumber and Trigger.BlockHash are already coordinated
				// TODO: check for upkeepID being active upkeep here using b.active
			}

		} else {
			b.lggr.Warnw("upkeep is not active, skipping", "upkeepID", proposal.UpkeepID)
		}

		payloads[i] = payload
	}

	return payloads, nil
}
