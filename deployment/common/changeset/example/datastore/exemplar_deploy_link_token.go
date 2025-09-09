package example

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/link_token"

	exemplarMd "github.com/smartcontractkit/chainlink/deployment/exemplar/metadata"
)

// ExemplarDeployLinkToken implements the deployment.ChangeSetV2 interface
var _ cldf.ChangeSetV2[uint64] = ExemplarDeployLinkToken{}

// ExemplarDeployLinkToken is an example changeset that deploys the LinkToken contract to an EVM chain and updates the
// environment datastore using the exemplar metadata types.
type ExemplarDeployLinkToken struct{}

// VerifyPreconditions checks if the chainSelector is a valid EVM chain selector.
func (cs ExemplarDeployLinkToken) VerifyPreconditions(_ cldf.Environment, chainSelector uint64) error {
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
func (cs ExemplarDeployLinkToken) Apply(e cldf.Environment, chainSelector uint64) (cldf.ChangesetOutput, error) {
	// Create an in-memory data store to store the address references, contract metadata and env metadata changes.
	dataStore := ds.NewMemoryDataStore()

	// Get the chain from the environment
	chain, ok := e.BlockChains.EVMChains()[chainSelector]
	if !ok {
		return cldf.ChangesetOutput{},
			fmt.Errorf("chain not found in environment: %d", chainSelector)
	}

	// Deploy the contract using geth bindings
	addr, tx, _, err := link_token.DeployLinkToken(chain.DeployerKey, chain.Client)
	if err != nil {
		return cldf.ChangesetOutput{},
			fmt.Errorf("failed to deploy link token contract: %w", err)
	}

	// Wait for the transaction to be confirmed and get the block number
	var blockNumber uint64
	if blockNumber, err = chain.Confirm(tx); err != nil {
		return cldf.ChangesetOutput{},
			fmt.Errorf("failed to confirm transaction: %w", err)
	}

	// Add a new AddressRef pointing to the deployed contract
	if err = dataStore.Addresses().Add(
		ds.AddressRef{
			ChainSelector: chainSelector,
			Address:       addr.String(),
			Type:          "LinkToken",
			Version:       semver.MustParse("1.0.0"),
			Qualifier:     "LinkTokenContractV1_" + addr.String(),
			Labels: ds.NewLabelSet(
				"LinkToken",
				"LinkTokenV1_0_0",
			),
		},
	); err != nil {
		return cldf.ChangesetOutput{},
			fmt.Errorf("failed to save address ref in datastore: %w", err)
	}

	// Add a new ContractMetadata entry for the deployed contract with information about the deployment.
	if err = dataStore.ContractMetadata().Add(
		ds.ContractMetadata{
			ChainSelector: chainSelector,
			Address:       addr.String(),
			Metadata: exemplarMd.SimpleContract{
				DeployedAt:  tx.Time(),
				TxHash:      tx.Hash(),
				BlockNumber: blockNumber,
			},
		},
	); err != nil {
		return cldf.ChangesetOutput{},
			fmt.Errorf("failed to save contract metadata in datastore: %w", err)
	}

	// Fetch the existing env metadata so we can update it with the new deployment count.
	envMetadata, err := e.DataStore.EnvMetadata().Get()
	if err != nil {
		// if err is different from ds.ErrEnvMetadataNotSet, return the error
		if !errors.Is(err, ds.ErrEnvMetadataNotSet) {
			return cldf.ChangesetOutput{},
				fmt.Errorf("failed to fetch existing env metadata: %w", err)
		}

		// Ensure the env metadata is initialized if it doesn't exist yet
		envMetadata = ds.EnvMetadata{
			Metadata: exemplarMd.SimpleEnv{
				DeployCounts: make(map[uint64]int64),
			},
		}
	}

	typedMeta, err := ds.As[exemplarMd.SimpleEnv](envMetadata.Metadata)
	if err != nil {
		return cldf.ChangesetOutput{},
			fmt.Errorf("failed to cast env metadata to SimpleEnv: %w", err)
	}

	// Increment the deployment count for the chain selector
	typedMeta.DeployCounts[chainSelector]++

	// Update the env metadata in the in-memory data store
	if err = dataStore.EnvMetadata().Set(
		ds.EnvMetadata{
			Metadata: typedMeta,
		},
	); err != nil {
		return cldf.ChangesetOutput{},
			fmt.Errorf("failed to save updated env metadata in datastore: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore: dataStore,
	}, nil
}
