package forwarder

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	forwarderwrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

// AcceptOwnershipInput identifies which KeystoneForwarder contract should accept pending ownership.
type AcceptOwnershipInput struct {
	// ChainSelector of the chain where the forwarder is deployed.
	ChainSelector uint64 `json:"chainSelector" yaml:"chainSelector"`
	// Qualifier optionally disambiguates the forwarder in the datastore.
	// Leave empty if there is only one forwarder on the chain.
	Qualifier string `json:"qualifier,omitempty" yaml:"qualifier,omitempty"`
}

// AcceptOwnershipForwarder directly calls AcceptOwnership on a KeystoneForwarder contract
// using the environment's deployer key. Use this after the previous owner has called
// transferOwnership to the deployer EOA.
type AcceptOwnershipForwarder struct{}

var _ cldf.ChangeSetV2[AcceptOwnershipInput] = AcceptOwnershipForwarder{}

func (AcceptOwnershipForwarder) VerifyPreconditions(e cldf.Environment, input AcceptOwnershipInput) error {
	if _, ok := e.BlockChains.EVMChains()[input.ChainSelector]; !ok {
		return fmt.Errorf("chain selector %d not found in environment", input.ChainSelector)
	}
	filters := []datastore.FilterFunc[datastore.AddressRefKey, datastore.AddressRef]{
		datastore.AddressRefByChainSelector(input.ChainSelector),
		datastore.AddressRefByType(datastore.ContractType(contracts.KeystoneForwarder)),
	}
	if input.Qualifier != "" {
		filters = append(filters, datastore.AddressRefByQualifier(input.Qualifier))
	}
	refs := e.DataStore.Addresses().Filter(filters...)
	if len(refs) == 0 {
		return fmt.Errorf("no KeystoneForwarder found for chain %d (qualifier %q)", input.ChainSelector, input.Qualifier)
	}
	return nil
}

func (AcceptOwnershipForwarder) Apply(e cldf.Environment, input AcceptOwnershipInput) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.EVMChains()[input.ChainSelector]

	filters := []datastore.FilterFunc[datastore.AddressRefKey, datastore.AddressRef]{
		datastore.AddressRefByChainSelector(input.ChainSelector),
		datastore.AddressRefByType(datastore.ContractType(contracts.KeystoneForwarder)),
	}
	if input.Qualifier != "" {
		filters = append(filters, datastore.AddressRefByQualifier(input.Qualifier))
	}

	refs := e.DataStore.Addresses().Filter(filters...)
	for _, ref := range refs {
		contract, err := forwarderwrapper.NewKeystoneForwarder(common.HexToAddress(ref.Address), chain.Client)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to instantiate forwarder %s on chain %d: %w", ref.Address, input.ChainSelector, err)
		}

		tx, err := contract.AcceptOwnership(chain.DeployerKey)
		if _, confErr := cldf.ConfirmIfNoError(chain, tx, err); confErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("AcceptOwnership failed for %s on chain %d: %w", ref.Address, input.ChainSelector, confErr)
		}

		e.Logger.Infow("Accepted ownership of KeystoneForwarder", "address", ref.Address, "chainSelector", input.ChainSelector)
	}

	return cldf.ChangesetOutput{}, nil
}
