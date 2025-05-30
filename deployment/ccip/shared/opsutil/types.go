package opsutil

import (
	"github.com/smartcontractkit/mcms"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

// ConfigureDependencies is used on sequence and the operation which needs to configure contracts and make use of deployer group
type ConfigureDependencies struct {
	Env          cldf.Environment
	CurrentState stateview.CCIPOnChainState
	AddressBook  cldf.AddressBook
}

// DeployContractDependencies is used on operations which need to deploy contracts or operations which can be used without
// the deployer group
type DeployContractDependencies struct {
	Chain       cldf_evm.Chain
	AddressBook cldf.AddressBook
}

type OpOutput struct {
	Proposals                  []mcms.TimelockProposal
	DescribedTimelockProposals []string
}

func (o *OpOutput) Merge(other OpOutput) error {
	if len(other.Proposals) > 0 {
		o.Proposals = append(o.Proposals, other.Proposals...)
	}
	if len(other.DescribedTimelockProposals) > 0 {
		o.DescribedTimelockProposals = append(o.DescribedTimelockProposals, other.DescribedTimelockProposals...)
	}
	return nil
}

func (o *OpOutput) ToChangesetOutput(deps *DeployContractDependencies) cldf.ChangesetOutput {
	if deps == nil || deps.AddressBook == nil {
		return cldf.ChangesetOutput{
			MCMSTimelockProposals:      o.Proposals,
			DescribedTimelockProposals: o.DescribedTimelockProposals,
		}
	}
	return cldf.ChangesetOutput{
		MCMSTimelockProposals:      o.Proposals,
		DescribedTimelockProposals: o.DescribedTimelockProposals,
		AddressBook:                deps.AddressBook,
	}
}
