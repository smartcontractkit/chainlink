package opsutil

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

type DeployContractInput struct {
	Chain uint64
	AB    cldf.AddressBook
}

type OpDependencies struct {
	Env          cldf.Environment
	CurrentState stateview.CCIPOnChainState
}

type OpOutput struct {
	Proposals                  []mcms.TimelockProposal
	DescribedTimelockProposals []string
	AddressBook                cldf.AddressBook
}

func (o *OpOutput) Merge(other OpOutput, env cldf.Environment) error {
	if o.AddressBook == nil {
		o.AddressBook = other.AddressBook
	} else if other.AddressBook != nil {
		if err := o.AddressBook.Merge(other.AddressBook); err != nil {
			return fmt.Errorf("failed to merge address book: %w", err)
		}
		if err := env.ExistingAddresses.Merge(other.AddressBook); err != nil {
			return fmt.Errorf("failed to merge existing addresses to environment: %w", err)
		}
	}
	if len(other.Proposals) == 0 {
		o.Proposals = append(o.Proposals, other.Proposals...)
	}
	if len(other.DescribedTimelockProposals) == 0 {
		o.DescribedTimelockProposals = append(o.DescribedTimelockProposals, other.DescribedTimelockProposals...)
	}
	return nil
}

func (o *OpOutput) ToChangesetOutput() cldf.ChangesetOutput {
	return cldf.ChangesetOutput{
		MCMSTimelockProposals:      o.Proposals,
		DescribedTimelockProposals: o.DescribedTimelockProposals,
		AddressBook:                o.AddressBook,
	}
}
