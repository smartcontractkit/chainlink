package example

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"
	"github.com/smartcontractkit/chainlink/deployment"
	exemplarMd "github.com/smartcontractkit/chainlink/deployment/exemplar/metadata"
)

// ExemplarDeployLinkToken implements the deployment.ChangeSetV2 interface
var _ deployment.ChangeSetV2[uint64] = ExemplarDeployLinkToken{}

// ExemplarDeployLinkToken is an example changeset that deploys the LinkToken contract to an EVM chain and updates the
// environment datastore using the exemplar metadata types.
type ExemplarDeployLinkToken struct{}

// VerifyPreconditions checks if the chainSelector is a valid EVM chain selector.
func (cs ExemplarDeployLinkToken) VerifyPreconditions(_ deployment.Environment, chainSelector uint64) error {
	fam, err := chain_selectors.GetSelectorFamily(chainSelector)
	if err != nil {
		return fmt.Errorf("failed to get chain selector family: %w", err)
	}

	if fam != chain_selectors.FamilyEVM {
		return fmt.Errorf("invalid chain selector for EVM: %d", chainSelector)
	}

	return nil
}

// Apply deploys the LinkToken contract to the specified EVM chain and updates the environment datastore.
func (cs ExemplarDeployLinkToken) Apply(e deployment.Environment, chainSelector uint64) (deployment.ChangesetOutput, error) {
	// Read contents of the environment datastore. Since it is loaded using the DefaultMetadata types, it needs to be converted
	// to the exemplar metadata types before use using the FromDefault utility function.
	envDatastore, err := ds.FromDefault[
		exemplarMd.SimpleContract,
		exemplarMd.SimpleEnv,
	](e.DataStore)
	if err != nil {
		return deployment.ChangesetOutput{},
			fmt.Errorf("failed to convert environment datastore: %w", err)
	}

	// Create an in-memory data store to store the  address references, contract metadata and env metadata changes.
	dataStore := ds.NewMemoryDataStore[
		exemplarMd.SimpleContract,
		exemplarMd.SimpleEnv,
	]()

	// Get the chain from the environment
	chain, ok := e.Chains[chainSelector]
	if !ok {
		return deployment.ChangesetOutput{},
			fmt.Errorf("chain not found in environment: %d", chainSelector)
	}

	// Deploy the contract using geth bindings
	addr, tx, _, err := link_token.DeployLinkToken(chain.DeployerKey, chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{},
			fmt.Errorf("failed to deploy link token contract: %w", err)
	}

	// Wait for the transaction to be confirmed and get the block number
	var blockNumber uint64
	if blockNumber, err = chain.Confirm(tx); err != nil {
		return deployment.ChangesetOutput{},
			fmt.Errorf("failed to confirm transaction: %w", err)
	}

	// Add a new AddressRef pointing to the deployed contract
	if err = dataStore.Addresses().Add(
		ds.AddressRef{
			ChainSelector: chainSelector,
			Address:       addr.String(),
			Type:          "LinkToken",
			Version:       semver.MustParse("1.0.0"),
			Qualifier:     fmt.Sprintf("LinkTokenContractV1_%s", addr.String()),
			Labels: ds.NewLabelSet(
				"LinkToken",
				"LinkTokenV1_0_0",
			),
		},
	); err != nil {
		return deployment.ChangesetOutput{},
			fmt.Errorf("failed to save address ref in datastore: %w", err)
	}

	// Add a new ContractMetadata entry for the deployed contract with information about the deployment.
	if err = dataStore.ContractMetadata().Add(
		ds.ContractMetadata[exemplarMd.SimpleContract]{
			ChainSelector: chainSelector,
			Address:       addr.String(),
			Metadata: exemplarMd.SimpleContract{
				DeployedAt:  tx.Time(),
				TxHash:      tx.Hash(),
				BlockNumber: blockNumber,
			},
		},
	); err != nil {
		return deployment.ChangesetOutput{},
			fmt.Errorf("failed to save contract metadata in datastore: %w", err)
	}

	// Fetch the existing env metadata so we can update it with the new deployment count.
	envMetadata, err := envDatastore.EnvMetadata().Get()
	if err != nil {
		if err != ds.ErrEnvMetadataNotSet {
			return deployment.ChangesetOutput{},
				fmt.Errorf("failed to fetch existing env metadata: %w", err)
		}

		// Ensure the env metadata is initialized if it doesn't exist yet
		envMetadata = ds.EnvMetadata[exemplarMd.SimpleEnv]{
			Domain:      "exemplar",
			Environment: e.Name,
			Metadata: exemplarMd.SimpleEnv{
				DeployCounts: make(map[uint64]int64),
			},
		}
	}
	// Increment the deployment count for the chain selector
	envMetadata.Metadata.DeployCounts[chainSelector]++

	// Update the env metadata in the in-memory data store
	if err = dataStore.EnvMetadata().Set(
		ds.EnvMetadata[exemplarMd.SimpleEnv]{
			Domain:      "exemplar",
			Environment: e.Name,
			Metadata:    envMetadata.Metadata,
		},
	); err != nil {
		return deployment.ChangesetOutput{},
			fmt.Errorf("failed to save updated env metadata in datastore: %w", err)
	}

	// ChangesetOutput accepts a DataStore that uses the DefaultMetadata types, so we need to convert the in-memory data store,
	// this conversion can be performed using the ToDefault utility function.
	ds, err := ds.ToDefault(dataStore.Seal())
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to convert data store to default format: %w", err)
	}

	return deployment.ChangesetOutput{
		DataStore: ds,
	}, nil
}
