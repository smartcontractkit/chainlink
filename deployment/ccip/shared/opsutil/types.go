package opsutil

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

type OpDependencies struct {
	Env          cldf.Environment
	CurrentState stateview.CCIPOnChainState
	AddressBook  cldf.AddressBook
}

type OpOutput struct {
	Proposals                  []mcms.TimelockProposal
	DescribedTimelockProposals []string
}

func (o *OpOutput) Merge(other OpOutput) error {
	if len(other.Proposals) == 0 {
		o.Proposals = append(o.Proposals, other.Proposals...)
	}
	if len(other.DescribedTimelockProposals) == 0 {
		o.DescribedTimelockProposals = append(o.DescribedTimelockProposals, other.DescribedTimelockProposals...)
	}
	return nil
}

func (o *OpOutput) ToChangesetOutput(deps OpDependencies) cldf.ChangesetOutput {
	return cldf.ChangesetOutput{
		MCMSTimelockProposals:      o.Proposals,
		DescribedTimelockProposals: o.DescribedTimelockProposals,
		AddressBook:                deps.AddressBook,
	}
}
