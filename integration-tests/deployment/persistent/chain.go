package persistent

import (
	"fmt"
	"github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// NewChains returns chains and RPC providers for the given configuration. It will start new chains and return their RPC providers, and it will
// connect to existing chains and return their RPC providers.
func NewChains(lggr logger.Logger, config persistent_types.ChainConfig) (map[uint64]deployment.Chain, map[uint64]persistent_types.RpcProvider, error) {
	lggr.Info("Setting up persistent chains")
	existingChains, existingRPCs, err := BuildExistingChainsAndEndpoints(config.ExistingEVMChains)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create existing Chains")
	}
	newDockerChains, newDockerRPCs, err := StartNewDockerChainsAndPrepareEndpoints(config.NewEVMChains)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create new Chains")
	}

	return mergeExistingAndNewChainsAndEndpoints(newDockerChains, existingChains, newDockerRPCs, existingRPCs)
}

func StartNewDockerChainsAndPrepareEndpoints(configs []persistent_types.NewEVMChainProducer) (map[uint64]deployment.Chain, map[uint64]persistent_types.RpcProvider, error) {
	chains := make(map[uint64]deployment.Chain)
	rpcProviders := make(map[uint64]persistent_types.RpcProvider)
	for _, config := range configs {
		chain, rpcProvider, err := config.Chain()
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to create chain")
		}

		chainData, found := chainselectors.ChainBySelector(chain.Selector)
		if !found {
			return nil, nil, errors.Wrapf(err, "failed to get chain data for selector %d", chain.Selector)
		}

		chains[chainData.EvmChainID] = chain
		rpcProviders[chainData.EvmChainID] = rpcProvider
	}

	return chains, rpcProviders, nil
}

func BuildExistingChainsAndEndpoints(configs []persistent_types.ExistingEVMChainProducer) (map[uint64]deployment.Chain, map[uint64]persistent_types.RpcProvider, error) {
	chains := make(map[uint64]deployment.Chain)
	rpcProviders := make(map[uint64]persistent_types.RpcProvider)
	for _, config := range configs {
		chain, rpcProvider, err := config.Chain()
		if err != nil {
			return nil, nil, err
		}

		chainData, found := chainselectors.ChainBySelector(chain.Selector)
		if !found {
			return nil, nil, errors.Wrapf(err, "failed to get chain data for selector %d", chain.Selector)
		}

		chains[chainData.EvmChainID] = chain
		rpcProviders[chainData.EvmChainID] = rpcProvider
	}
	return chains, rpcProviders, nil
}

func mergeExistingAndNewChainsAndEndpoints(newDockerChains, existingChains map[uint64]deployment.Chain, newDockerRPCs, existingRPCs map[uint64]persistent_types.RpcProvider) (map[uint64]deployment.Chain, map[uint64]persistent_types.RpcProvider, error) {
	chains := make(map[uint64]deployment.Chain)
	rpcProviders := make(map[uint64]persistent_types.RpcProvider)
	for k, v := range existingChains {
		chains[k] = v
	}
	for k, v := range newDockerChains {
		if _, ok := chains[k]; ok {
			return nil, nil, fmt.Errorf("duplicate chain id %d used by new and existing Chains", k)
		}
		chains[k] = v
	}

	for k, v := range existingRPCs {
		rpcProviders[k] = v
	}
	for k, v := range newDockerRPCs {
		if _, ok := rpcProviders[k]; ok {
			return nil, nil, fmt.Errorf("duplicate chain id %d used by new and existing Chains", k)
		}
		rpcProviders[k] = v
	}

	return chains, rpcProviders, nil
}
