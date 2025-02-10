package changeset

import (
	"fmt"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type ChangesetApplication struct {
	Changeset deployment.ChangeSet[any]
	Config    any
}

func WrapChangeSet[C any](fn deployment.ChangeSet[C]) func(e deployment.Environment, config any) (deployment.ChangesetOutput, error) {
	return func(e deployment.Environment, config any) (deployment.ChangesetOutput, error) {
		var zeroC C
		if config != nil {
			c, ok := config.(C)
			if !ok {
				return deployment.ChangesetOutput{}, fmt.Errorf("invalid config type, expected %T", c)
			}
			return fn(e, config.(C))
		}

		return fn(e, zeroC)
	}
}

// ApplyChangesets applies the changeset applications to the environment and returns the updated environment.
func ApplyChangesets(t *testing.T, e deployment.Environment, timelockContractsPerChain map[uint64]*proposalutils.TimelockExecutionContracts, changesetApplications []ChangesetApplication) (deployment.Environment, error) {
	currentEnv := e
	for i, csa := range changesetApplications {
		out, err := csa.Changeset(currentEnv, csa.Config)
		if err != nil {
			return e, fmt.Errorf("failed to apply changeset at index %d: %w", i, err)
		}
		var addresses deployment.AddressBook
		if out.AddressBook != nil {
			addresses = out.AddressBook
			err := addresses.Merge(currentEnv.ExistingAddresses)
			if err != nil {
				return e, fmt.Errorf("failed to merge address book: %w", err)
			}
		} else {
			addresses = currentEnv.ExistingAddresses
		}
		if out.Jobs != nil {
			// do nothing, as these jobs auto-accept.
		}
		if out.Proposals != nil {
			for _, prop := range out.Proposals {
				chains := mapset.NewSet[uint64]()
				for _, op := range prop.Transactions {
					chains.Add(uint64(op.ChainIdentifier))
				}

				signed := proposalutils.SignProposal(t, e, &prop)
				for _, sel := range chains.ToSlice() {
					timelockContracts, ok := timelockContractsPerChain[sel]
					if !ok || timelockContracts == nil {
						return deployment.Environment{}, fmt.Errorf("timelock contracts not found for chain %d", sel)
					}

					err := proposalutils.ExecuteProposal(t, e, signed, timelockContracts, sel) //nolint:staticcheck //SA1019 ignoring deprecated function for compatibility; we don't have tools to generate the new field
					if err != nil {
						return e, fmt.Errorf("failed to execute proposal: %w", err)
					}
				}
			}
		}
		if out.MCMSTimelockProposals != nil {
			for _, prop := range out.MCMSTimelockProposals {
				chains := mapset.NewSet[uint64]()
				for _, op := range prop.Operations {
					chains.Add(uint64(op.ChainSelector))
				}

				p := proposalutils.SignMCMSTimelockProposal(t, e, &prop)
				for _, sel := range chains.ToSlice() {
					timelockContracts, ok := timelockContractsPerChain[sel]
					if !ok || timelockContracts == nil {
						return deployment.Environment{}, fmt.Errorf("timelock contracts not found for chain %d", sel)
					}

					proposalutils.ExecuteMCMSProposalV2(t, e, p, sel)
					proposalutils.ExecuteMCMSTimelockProposalV2(t, e, &prop, sel)
				}
			}
		}
		if out.MCMSProposals != nil {
			for _, prop := range out.MCMSProposals {
				chains := mapset.NewSet[uint64]()
				for _, op := range prop.Operations {
					chains.Add(uint64(op.ChainSelector))
				}

				p := proposalutils.SignMCMSProposal(t, e, &prop)
				for _, sel := range chains.ToSlice() {
					proposalutils.ExecuteMCMSProposalV2(t, e, p, sel)
				}
			}
		}
		currentEnv = deployment.Environment{
			Name:              e.Name,
			Logger:            e.Logger,
			ExistingAddresses: addresses,
			Chains:            e.Chains,
			SolChains:         e.SolChains,
			NodeIDs:           e.NodeIDs,
			Offchain:          e.Offchain,
			OCRSecrets:        e.OCRSecrets,
			GetContext:        e.GetContext,
		}
	}
	return currentEnv, nil
}
