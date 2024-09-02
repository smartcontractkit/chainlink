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
func NewChains(lggr logger.Logger, config persistent_types.ChainConfig) (map[uint64]deployment.Chain, error) {
	lggr.Info("Creating persistent Chains")
	existingChains, err := newExistingChains(config.ExistingEVMChains)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create existing Chains")
	}
	createdChains, err := newChains(config.NewEVMChains)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create new Chains")
	}
	chains := make(map[uint64]deployment.Chain)
	for k, v := range existingChains {
		if _, ok := chains[k]; ok {
			return nil, errors.Wrapf(err, "duplicate chain id %d used by new and existing Chains", k)
		}
		chains[k] = v
	}
	for k, v := range createdChains {
		chains[k] = v
	}
	return chains, nil
}

func newChains(configs []persistent_types.NewEVMChainConfig) (map[uint64]deployment.Chain, error) {
	chains := make(map[uint64]deployment.Chain)
	for _, config := range configs {
		chain, err := config.Chain()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create chain")
		}

		chainData, found := chainselectors.ChainBySelector(chain.Selector)
		if !found {
			return nil, errors.Wrapf(err, "failed to get chain data for selector %d", chain.Selector)
		}

		chains[chainData.EvmChainID] = chain
	}

	return chains, nil
}

func newExistingChains(configs []persistent_types.ExistingEVMChainConfig) (map[uint64]deployment.Chain, error) {
	chains := make(map[uint64]deployment.Chain)
	for _, config := range configs {
		chain, err := config.Chain()
		if err != nil {
			return chains, err
		}

		chainData, found := chainselectors.ChainBySelector(chain.Selector)
		if !found {
			return nil, errors.Wrapf(err, "failed to get chain data for selector %d", chain.Selector)
		}

		chains[chainData.EvmChainID] = chain
	}
	return chains, nil
}
