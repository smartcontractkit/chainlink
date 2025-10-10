package internal

import (
	"fmt"
	"maps"

	"github.com/ethereum/go-ethereum/common"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type configureForwarderContractsRequest struct {
	Dons []RegisteredDon

	Chains  map[uint64]struct{} // list of chains for which request will be executed. If empty, request is applied to all chains
	UseMCMS bool
}
type configureForwarderContractsResponse struct {
	// ForwarderAddresses is a map of chain selector to forwarder contract address that has been configured (non-MCMS),
	// or will be configured (MCMS).
	ForwarderAddresses map[uint64]common.Address
	OpsPerChain        map[uint64]mcmstypes.BatchOperation
	Config             map[uint64]ForwarderConfig
}

// Depreciated: use [changeset.configureForwardContracts] instead
// configureForwardContracts configures the forwarder contracts on all chains for the given DONS
// the address book is required to contain the an address of the deployed forwarder contract for every chain in the environment
func configureForwardContracts(env *cldf.Environment, req configureForwarderContractsRequest) (*configureForwarderContractsResponse, error) {
	evmChains := env.BlockChains.EVMChains()
	contractSetsResp, err := GetContractSets(env.Logger, &GetContractSetsRequest{
		Chains:      evmChains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}

	opPerChain := make(map[uint64]mcmstypes.BatchOperation)
	forwarderAddresses := make(map[uint64]common.Address)
	configs := make(map[uint64]ForwarderConfig)
	// configure forwarders on all chains
	for _, chain := range evmChains {
		if _, shouldInclude := req.Chains[chain.Selector]; len(req.Chains) > 0 && !shouldInclude {
			continue
		}
		// get the forwarder contract for the chain
		contracts, ok := contractSetsResp.ContractSets[chain.Selector]
		if !ok {
			return nil, fmt.Errorf("failed to get contract set for chain %d", chain.Selector)
		}
		r, err := configureForwarder(env.Logger, chain, contracts.Forwarder, req.Dons, req.UseMCMS)
		if err != nil {
			return nil, fmt.Errorf("failed to configure forwarder for chain selector %d: %w", chain.Selector, err)
		}
		configs[chain.Selector] = r.Config
		maps.Copy(opPerChain, r.Ops)
		forwarderAddresses[chain.Selector] = contracts.Forwarder.Address()
	}
	return &configureForwarderContractsResponse{
		ForwarderAddresses: forwarderAddresses,
		OpsPerChain:        opPerChain,
		Config:             configs,
	}, nil
}
