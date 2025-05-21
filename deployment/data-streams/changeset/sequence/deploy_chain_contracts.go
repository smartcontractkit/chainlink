package sequence

import (
	"errors"
	"fmt"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"

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
	Billing        types.BillingFeature
	Ownership      types.OwnershipFeature
}

func deployDataStreamsLogic(e cldf.Environment, cc DeployDataStreamsConfig) (cldf.ChangesetOutput, error) {
	deployedAddresses := ds.NewMemoryDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata]()

	// Prevents mutating environment state - injected environment is not expected to be updated during changeset Apply
	cloneEnv := e.Clone()

	var timelockProposals []mcms.TimelockProposal

	for chainSel, cfg := range cc.ChainsToDeploy {
		family, err := chainselectors.GetSelectorFamily(chainSel)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get family for chain %d: %w", chainSel, err)
		}
		switch family {
		case chainselectors.FamilyEVM:
			chainProposals, err := deployChainComponentsEVM(&cloneEnv, chainSel, cfg, deployedAddresses)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy components for chain %d: %w", chainSel, err)
			}
			timelockProposals = append(timelockProposals, chainProposals...)
		default:
			return cldf.ChangesetOutput{}, fmt.Errorf("unsupported chain family %s for chain %d", family, chainSel)
		}
	}

	if len(timelockProposals) > 0 {
		mergedTimelockProposal, err := mcmsutil.MergeSimilarTimelockProposals(timelockProposals)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge timelock proposals: %w", err)
		}
		timelockProposals = []mcms.TimelockProposal{mergedTimelockProposal}
	}

	sealedDS, err := ds.ToDefault(deployedAddresses.Seal())
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert data store to default format: %w", err)
	}

	ab, err := utils.DataStoreToAddressBook(sealedDS.Seal())
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert data store to address book: %w", err)
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab, // backwards compatibility. This will be removed in the future.
		DataStore:             sealedDS,
		MCMSTimelockProposals: timelockProposals,
	}, nil
}

func deployDataStreamsPrecondition(_ cldf.Environment, cc DeployDataStreamsConfig) error {
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
		if err := cldf.IsValidChainSelector(chain); err != nil {
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
