package operations

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type OpDependencies struct {
	Env          cldf.Environment
	CurrentState stateview.CCIPOnChainState
}

type OpOutput struct {
	Proposals   []mcms.TimelockProposal
	AddressBook cldf.AddressBook
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
	return nil
}

type postOpsInput struct {
	SolanaChainSelector uint64
	EVMChainSelector    uint64
	MCMSConfig          *proposalutils.TimelockConfig
	Proposals           []mcmslib.TimelockProposal
}

var (
	PostOps = operations.NewOperation(
		"postOpsToAggregateProposals",
		semver.MustParse("1.0.0"),
		"Post ops to aggregate proposals",
		func(b operations.Bundle, deps OpDependencies, input postOpsInput) ([]mcmslib.TimelockProposal, error) {
			allProposals := input.Proposals
			proposal, err := proposalutils.AggregateProposals(
				deps.Env, deps.EVMMCMSState, deps.SolanaMCMSState, allProposals,
				"Adding EVM and Solana lane", input.MCMSConfig)
			if err != nil {
				return nil, err
			}
			if proposal != nil {
				input.Proposals = []mcmslib.TimelockProposal{*proposal}
			}
			return input.Proposals, nil
		},
	)
)
