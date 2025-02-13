package changeset

import (
	"fmt"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type ConfiguredChangeSet interface {
	Apply(e deployment.Environment) (deployment.ChangesetOutput, error)
}

func Configure[C any](
	changeset deployment.ChangeSetV2[C],
	config C,
) ConfiguredChangeSet {
	return configuredChangeSetImpl[C]{
		changeset: changeset,
		config:    config,
	}
}

type configuredChangeSetImpl[C any] struct {
	changeset deployment.ChangeSetV2[C]
	config    C
}

func (ca configuredChangeSetImpl[C]) Apply(e deployment.Environment) (deployment.ChangesetOutput, error) {
	err := ca.changeset.VerifyPreconditions(e, ca.config)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return ca.changeset.Apply(e, ca.config)
}

// Apply applies the changeset applications to the environment and returns the updated environment. This is the
// variadic function equivalent of ApplyChangesets, but allowing you to simply pass in one or more changesets as
// parameters at the end of the function. e.g. `changeset.Apply(t, e, nil, configuredCS1, configuredCS2)` etc.
func Apply(t *testing.T, e deployment.Environment, timelockContractsPerChain map[uint64]*proposalutils.TimelockExecutionContracts, first ConfiguredChangeSet, rest ...ConfiguredChangeSet) (deployment.Environment, error) {
	return ApplyChangesets(t, e, timelockContractsPerChain, append([]ConfiguredChangeSet{first}, rest...))
}

// ApplyChangesets applies the changeset applications to the environment and returns the updated environment.
func ApplyChangesets(t *testing.T, e deployment.Environment, timelockContractsPerChain map[uint64]*proposalutils.TimelockExecutionContracts, changesetApplications []ConfiguredChangeSet) (deployment.Environment, error) {
	currentEnv := e
	for i, csa := range changesetApplications {

		out, err := csa.Apply(currentEnv)
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
				mcmProp := proposalutils.SignMCMSTimelockProposal(t, e, &prop)
				proposalutils.ExecuteMCMSProposalV2(t, e, mcmProp)
				proposalutils.ExecuteMCMSTimelockProposalV2(t, e, &prop)
			}
		}
		if out.MCMSProposals != nil {
			for _, prop := range out.MCMSProposals {
				p := proposalutils.SignMCMSProposal(t, e, &prop)
				proposalutils.ExecuteMCMSProposalV2(t, e, p)
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
