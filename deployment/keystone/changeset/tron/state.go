package tron

import (
	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type ContractsState struct {
	Forwarder address.Address
}

type KeystoneTronChainState struct {
	Chains map[uint64]ContractsState
}

func LoadTronOnchainState(e cldf.Environment) (KeystoneTronChainState, error) {
	state := KeystoneTronChainState{
		Chains: make(map[uint64]ContractsState),
	}

	for chainSelector := range e.BlockChains.TronChains() {
		records := e.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSelector))
		e.Logger.Infof("Records: %+v, Datastore: %+v", records, e.DataStore.Addresses())
		contractsState, err := LoadTronContractsState(e.Logger, records)
		if err != nil {
			return state, err
		}
		state.Chains[chainSelector] = *contractsState
	}
	return state, nil
}

// LoadTronContractsState Loads all contracts for tron chain into state
func LoadTronContractsState(logger logger.Logger, addresses []datastore.AddressRef) (*ContractsState, error) {
	var state ContractsState

	for _, addr := range addresses {
		if addr.Type == ForwarderContract {
			forwarderAddress, err := address.Base58ToAddress(addr.Address)
			state.Forwarder = forwarderAddress

			return &state, err
		}
	}
	return &state, nil
}
