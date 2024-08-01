package postprocessors

import (
	"context"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

type addProposalToMetadataStore struct {
	metadataStore types.MetadataStore
}

func NewAddProposalToMetadataStorePostprocessor(store types.MetadataStore) *addProposalToMetadataStore {
	return &addProposalToMetadataStore{metadataStore: store}
}

func (a *addProposalToMetadataStore) PostProcess(_ context.Context, results []ocr2keepers.CheckResult, _ []ocr2keepers.UpkeepPayload) error {
	// should only add values and not remove them
	for _, r := range results {
		if r.PipelineExecutionState == 0 && r.Eligible {
			proposal := ocr2keepers.CoordinatedBlockProposal{
				UpkeepID: r.UpkeepID,
				Trigger:  r.Trigger,
				WorkID:   r.WorkID,
			}
			a.metadataStore.AddProposals(proposal)
		}
	}

	return nil
}
