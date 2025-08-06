package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	workflow_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
)

// SetMetadataConfigRequest configures metadata validation settings
type SetMetadataConfigRequest struct {
	ChainSelector uint64
	NameLen       uint8  // Maximum length for workflow names
	TagLen        uint8  // Maximum length for workflow tags
	URLLen        uint8  // Maximum length for URLs
	AttrLen       uint16 // Maximum length for attributes
}

// SetMetadataConfig changeset to configure metadata validation settings
var _ cldf.ChangeSet[SetMetadataConfigRequest] = SetMetadataConfig

func SetMetadataConfig(env cldf.Environment, req SetMetadataConfigRequest) (cldf.ChangesetOutput, error) {
	chain, ok := env.BlockChains.EVMChains()[req.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", req.ChainSelector)
	}

	registry, err := getWorkflowRegistryV2(env, req.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
	}

	tx, err := registry.SetMetadataConfig(chain.DeployerKey, req.NameLen, req.TagLen, req.URLLen, req.AttrLen)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to set metadata config: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm set metadata config transaction: %w", err)
	}

	env.Logger.Infof("Successfully set metadata config on chain %d, tx: %s", req.ChainSelector, tx.Hash().String())
	return cldf.ChangesetOutput{}, nil
}

// SetWorkflowOwnerConfigRequest configures workflow owner-specific settings
type SetWorkflowOwnerConfigRequest struct {
	ChainSelector uint64
	Owner         common.Address // Workflow owner address
	Config        []byte         // Owner-specific configuration data
}

// SetWorkflowOwnerConfig changeset to configure workflow owner settings
var _ cldf.ChangeSet[SetWorkflowOwnerConfigRequest] = SetWorkflowOwnerConfig

func SetWorkflowOwnerConfig(env cldf.Environment, req SetWorkflowOwnerConfigRequest) (cldf.ChangesetOutput, error) {
	chain, ok := env.BlockChains.EVMChains()[req.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", req.ChainSelector)
	}

	registry, err := getWorkflowRegistryV2(env, req.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
	}

	tx, err := registry.SetWorkflowOwnerConfig(chain.DeployerKey, req.Owner, req.Config)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to set workflow owner config: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm set workflow owner config transaction: %w", err)
	}

	env.Logger.Infof("Successfully set workflow owner config for %s on chain %d, tx: %s", req.Owner.String(), req.ChainSelector, tx.Hash().String())
	return cldf.ChangesetOutput{}, nil
}

// SetDONLimitRequest configures DON workflow limits
type SetDONLimitRequest struct {
	ChainSelector uint64
	DONFamily     string // DON family identifier
	Limit         uint32 // Maximum number of workflows per owner
	Enabled       bool   // Whether the limit is enabled
}

// SetDONLimit changeset to configure DON workflow limits
var _ cldf.ChangeSet[SetDONLimitRequest] = SetDONLimit

func SetDONLimit(env cldf.Environment, req SetDONLimitRequest) (cldf.ChangesetOutput, error) {
	chain, ok := env.BlockChains.EVMChains()[req.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", req.ChainSelector)
	}

	registry, err := getWorkflowRegistryV2(env, req.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
	}

	tx, err := registry.SetDONLimit(chain.DeployerKey, req.DONFamily, req.Limit, req.Enabled)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to set DON limit: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm set DON limit transaction: %w", err)
	}

	env.Logger.Infof("Successfully set DON limit for family %s on chain %d, tx: %s", req.DONFamily, req.ChainSelector, tx.Hash().String())
	return cldf.ChangesetOutput{}, nil
}

// SetUserDONOverrideRequest configures user-specific DON overrides
type SetUserDONOverrideRequest struct {
	ChainSelector uint64
	User          common.Address // User address
	DONFamily     string         // DON family identifier
	Limit         uint32         // User-specific limit
	Enabled       bool           // Whether the override is enabled
}

// SetUserDONOverride changeset to configure user-specific DON overrides
var _ cldf.ChangeSet[SetUserDONOverrideRequest] = SetUserDONOverride

func SetUserDONOverride(env cldf.Environment, req SetUserDONOverrideRequest) (cldf.ChangesetOutput, error) {
	chain, ok := env.BlockChains.EVMChains()[req.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", req.ChainSelector)
	}

	registry, err := getWorkflowRegistryV2(env, req.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
	}

	tx, err := registry.SetUserDONOverride(chain.DeployerKey, req.User, req.DONFamily, req.Limit, req.Enabled)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to set user DON override: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm set user DON override transaction: %w", err)
	}

	env.Logger.Infof("Successfully set user DON override for %s on chain %d, tx: %s", req.User.String(), req.ChainSelector, tx.Hash().String())
	return cldf.ChangesetOutput{}, nil
}

// SetDONRegistryRequest configures the DON registry address
type SetDONRegistryRequest struct {
	ChainSelector    uint64
	Registry         common.Address // DON registry contract address
	ChainSelectorDON uint64         // Chain selector where the DON registry exists
}

// SetDONRegistry changeset to configure the DON registry address
var _ cldf.ChangeSet[SetDONRegistryRequest] = SetDONRegistry

func SetDONRegistry(env cldf.Environment, req SetDONRegistryRequest) (cldf.ChangesetOutput, error) {
	chain, ok := env.BlockChains.EVMChains()[req.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", req.ChainSelector)
	}

	registry, err := getWorkflowRegistryV2(env, req.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
	}

	tx, err := registry.SetDONRegistry(chain.DeployerKey, req.Registry, req.ChainSelectorDON)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to set DON registry: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm set DON registry transaction: %w", err)
	}

	env.Logger.Infof("Successfully set DON registry %s on chain %d, tx: %s", req.Registry.String(), req.ChainSelector, tx.Hash().String())
	return cldf.ChangesetOutput{}, nil
}

// UpdateAllowedSignersRequest updates the list of allowed signers
type UpdateAllowedSignersRequest struct {
	ChainSelector uint64
	Signers       []common.Address // List of signer addresses
	Allowed       bool             // Whether to allow or disallow these signers
}

// UpdateAllowedSigners changeset to update allowed signers
var _ cldf.ChangeSet[UpdateAllowedSignersRequest] = UpdateAllowedSigners

func UpdateAllowedSigners(env cldf.Environment, req UpdateAllowedSignersRequest) (cldf.ChangesetOutput, error) {
	chain, ok := env.BlockChains.EVMChains()[req.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", req.ChainSelector)
	}

	registry, err := getWorkflowRegistryV2(env, req.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
	}

	tx, err := registry.UpdateAllowedSigners(chain.DeployerKey, req.Signers, req.Allowed)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to update allowed signers: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm update allowed signers transaction: %w", err)
	}

	env.Logger.Infof("Successfully updated allowed signers on chain %d, tx: %s", req.ChainSelector, tx.Hash().String())
	return cldf.ChangesetOutput{}, nil
}

// getWorkflowRegistryV2 is a helper function to get the workflow registry v2 contract instance
func getWorkflowRegistryV2(env cldf.Environment, chainSelector uint64) (*workflow_registry_v2.WorkflowRegistry, error) {
	// Look up the workflow registry address from the datastore
	addresses := env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSelector))
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no addresses found for chain selector %d", chainSelector)
	}

	var registryAddr common.Address
	found := false
	for _, addr := range addresses {
		if addr.Type == "WorkflowRegistry" {
			registryAddr = common.HexToAddress(addr.Address)
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("workflow registry address not found for chain selector %d", chainSelector)
	}

	chain, ok := env.BlockChains.EVMChains()[chainSelector]
	if !ok {
		return nil, fmt.Errorf("chain with selector %d not found", chainSelector)
	}

	registry, err := workflow_registry_v2.NewWorkflowRegistry(registryAddr, chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow registry instance: %w", err)
	}

	return registry, nil
}
