package preprocessors

import (
	"context"

	ocr2keepersv3 "github.com/smartcontractkit/chainlink-automation/pkg/v3"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

func NewProposalFilterer(metadata types.MetadataStore, upkeepType types.UpkeepType) ocr2keepersv3.PreProcessor[ocr2keepers.UpkeepPayload] {
	return &proposalFilterer{
		upkeepType: upkeepType,
		metadata:   metadata,
	}
}

type proposalFilterer struct {
	metadata   types.MetadataStore
	upkeepType types.UpkeepType
}

var _ ocr2keepersv3.PreProcessor[ocr2keepers.UpkeepPayload] = (*proposalFilterer)(nil)

func (p *proposalFilterer) PreProcess(ctx context.Context, payloads []ocr2keepers.UpkeepPayload) ([]ocr2keepers.UpkeepPayload, error) {
	all := p.metadata.ViewProposals(p.upkeepType)
	flatten := map[string]bool{}
	for _, proposal := range all {
		flatten[proposal.WorkID] = true
	}
	filtered := make([]ocr2keepers.UpkeepPayload, 0)
	for _, payload := range payloads {
		if _, ok := flatten[payload.WorkID]; !ok {
			filtered = append(filtered, payload)
		}
	}

	return filtered, nil
}
