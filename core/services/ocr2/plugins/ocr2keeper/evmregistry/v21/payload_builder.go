package evm

import (
	"context"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/logprovider"
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
		if !b.upkeepList.IsActive(proposal.UpkeepID.BigInt()) {
			b.lggr.Warnw("upkeep is not active, skipping", "upkeepID", proposal.UpkeepID)
			continue
		}
		b.lggr.Debugf("building payload for coordinated block proposal %+v", proposal)
		var checkData []byte
		var err error
		switch core.GetUpkeepType(proposal.UpkeepID) {
		case types.LogTrigger:
			checkData, err = b.recoverer.GetProposalData(ctx, proposal)
			if err != nil {
				b.lggr.Warnw("failed to get log proposal data", "err", err, "upkeepID", proposal.UpkeepID, "trigger", proposal.Trigger)
				continue
			}
		case types.ConditionTrigger:
			// Empty checkData for conditionals
		}
		payload, err = core.NewUpkeepPayload(proposal.UpkeepID.BigInt(), proposal.Trigger, checkData)
		if err != nil {
			b.lggr.Warnw("error building upkeep payload", "err", err, "upkeepID", proposal.UpkeepID)
			continue
		}

		payloads[i] = payload
	}

	return payloads, nil
}
