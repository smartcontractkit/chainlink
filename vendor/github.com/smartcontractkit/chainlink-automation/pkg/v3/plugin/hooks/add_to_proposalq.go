package hooks

import (
	"fmt"
	"log"

	ocr2keepersv3 "github.com/smartcontractkit/chainlink-automation/pkg/v3"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/telemetry"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
)

func NewAddToProposalQHook(proposalQ types.ProposalQueue, logger *log.Logger) AddToProposalQHook {
	return AddToProposalQHook{
		proposalQ: proposalQ,
		logger:    log.New(logger.Writer(), fmt.Sprintf("[%s | pre-build hook:add-to-proposalq]", telemetry.ServiceName), telemetry.LogPkgStdFlags),
	}
}

type AddToProposalQHook struct {
	proposalQ types.ProposalQueue
	logger    *log.Logger
}

func (hook *AddToProposalQHook) RunHook(outcome ocr2keepersv3.AutomationOutcome) {
	addedProposals := 0
	for _, roundProposals := range outcome.SurfacedProposals {
		err := hook.proposalQ.Enqueue(roundProposals...)
		if err != nil {
			// Do not return error, just log and skip this round's proposals
			hook.logger.Printf("Error adding proposals to queue: %v", err)
			continue
		}
		addedProposals += len(roundProposals)
	}
	hook.logger.Printf("Added %d proposals from outcome", addedProposals)

}
