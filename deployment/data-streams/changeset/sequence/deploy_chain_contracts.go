package sequence

import (
	"errors"
	"fmt"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verification"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
)

// DeployDataStreamsChainContractsChangeset deploys the entire data streams destination chain contracts. It should be kept up to date
// with the latest contract versions and deployment logic.
var DeployDataStreamsChainContractsChangeset = cldf.CreateChangeSet(deployDataStreamsLogic, deployDataStreamsPrecondition)

type DeployDataStreamsConfig struct {
	ChainsToDeploy map[uint64]DeployDataStreams
}

type DeployDataStreams struct {
	VerifierConfig verification.SetConfig

	Billing   types.BillingFeature
	Ownership types.OwnershipFeature
}

func deployDataStreamsLogic(e deployment.Environment, cc DeployDataStreamsConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook() // changeset output expects only new addresses

	// Clone env to avoid mutation
	cloneEnv, err := cloneEnvironment(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	var timelockProposals []mcms.TimelockProposal

	for chainSel, cfg := range cc.ChainsToDeploy {
		family, err := chain_selectors.GetSelectorFamily(chainSel)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get family for chain %d: %w", chainSel, err)
		}
		switch family {
		case chain_selectors.FamilyEVM:
			// Deploy each component of the system for this chain
			chainProposals, err := deployChainComponentsEVM(cloneEnv, chainSel, cfg, newAddresses)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy components for chain %d: %w", chainSel, err)
			}
			timelockProposals = append(timelockProposals, chainProposals...)
		default:
			return deployment.ChangesetOutput{}, fmt.Errorf("unsupported chain family %s for chain %d", family, chainSel)
		}
	}

	if len(timelockProposals) > 0 {
		mergedTimelockProposal, err := mcmsutil.MergeSimilarTimelockProposals(timelockProposals)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge timelock proposals: %w", err)
		}
		timelockProposals = []mcms.TimelockProposal{mergedTimelockProposal}
	}

	return deployment.ChangesetOutput{
		AddressBook:           newAddresses,
		MCMSTimelockProposals: timelockProposals,
	}, nil
}

func deployDataStreamsPrecondition(_ deployment.Environment, cc DeployDataStreamsConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployDataStreams config: %w", err)
	}
	return nil
}

func (cc DeployDataStreamsConfig) Validate() error {
	if len(cc.ChainsToDeploy) == 0 {
		return errors.New("ChainsToDeploy is empty")
	}

	if len(cc.ChainsToDeploy) > 1 {
		// MergeSimilarTimelockProposals only supports a single chain.
		// Add this support when chain deployment frequency increases.
		return errors.New("changeset currently does not support multiple chains")
	}

	for chain, cfg := range cc.ChainsToDeploy {
		if err := deployment.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}

		if err := cfg.Ownership.Validate(); err != nil {
			return fmt.Errorf("invalid ownership settings for chain %d: %w", chain, err)
		}

		if err := cfg.Billing.Validate(); err != nil {
			return fmt.Errorf("invalid billing settings for chain %d: %w", chain, err)
		}
	}
	return nil
}

// cloneEnvironment creates a copy of the environment to prevent mutations
func cloneEnvironment(e deployment.Environment) (deployment.Environment, error) {
	existingAddresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return deployment.Environment{}, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	abClone := deployment.NewMemoryAddressBookFromMap(existingAddresses)

	return deployment.Environment{
		Name:              e.Name,
		Logger:            e.Logger,
		ExistingAddresses: abClone,
		Chains:            e.Chains,
		SolChains:         e.SolChains,
		NodeIDs:           e.NodeIDs,
		Offchain:          e.Offchain,
		OCRSecrets:        e.OCRSecrets,
		GetContext:        e.GetContext,
		OperationsBundle:  e.OperationsBundle,
	}, nil
}
