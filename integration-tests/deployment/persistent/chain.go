package persistent

import (
	"github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// NewChains creates Chains based on the provided configuration. It returns a map of chain id to chain.
// You can mix existing and new Chains in the configuration, meaning that you can have Chains that are already running and Chains that will be started by the test environment.
func NewChains(lggr logger.Logger, config persistent_types.ChainConfig) (map[uint64]deployment.Chain, map[uint64]persistent_types.RpcProvider, error) {
	lggr.Info("Setting up persistent chains")
	existingChains, existingRpcs, err := newExistingChains(config.ExistingEVMChains)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create existing Chains")
	}
	createdChains, createdRpcs, err := newChains(config.NewEVMChains)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create new Chains")
	}
	chains := make(map[uint64]deployment.Chain)
	rpcProviders := make(map[uint64]persistent_types.RpcProvider)
	for k, v := range existingChains {
		chains[k] = v
	}
	for k, v := range createdChains {
		if _, ok := chains[k]; ok {
			return nil, nil, errors.Wrapf(err, "duplicate chain id %d used by new and existing Chains", k)
		}
		chains[k] = v
	}

	for k, v := range existingRpcs {
		rpcProviders[k] = v
	}
	for k, v := range createdRpcs {
		if _, ok := rpcProviders[k]; ok {
			return nil, nil, errors.Wrapf(err, "duplicate chain id %d used by new and existing Chains", k)
		}
		rpcProviders[k] = v
	}

	return chains, rpcProviders, nil
}

func newChains(configs []persistent_types.NewEVMChainProducer) (map[uint64]deployment.Chain, map[uint64]persistent_types.RpcProvider, error) {
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

func newExistingChains(configs []persistent_types.ExistingEVMChainProducer) (map[uint64]deployment.Chain, map[uint64]persistent_types.RpcProvider, error) {
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
