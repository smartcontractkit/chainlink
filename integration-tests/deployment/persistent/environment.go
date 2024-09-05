package persistent

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func NewEnvironment(lggr logger.Logger, config persistent_types.EnvironmentConfig) (*deployment.Environment, error) {
	chains, rpcProviders, err := NewChains(lggr, config.ChainConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create Chains")
	}

	if config.EnvironmentHooks != nil {
		if hookErr := config.EnvironmentHooks.PostChainStartupHooks(chains, rpcProviders, &config); hookErr != nil {
			return nil, errors.Wrapf(hookErr, "failed to run post chain startup hooks")
		}
	}

	if config.DONConfig.NewDON != nil {
		clNodesConfigs, err := BuildEVMOnlyChainlinkConfigs(config.DONConfig.NewDON.ChainlinkDeployment, rpcProviders)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create chainlink configs")
		}
		config.DONConfig.NewDON.ChainlinkConfigs = clNodesConfigs
	}

	don, err := NewNodes(config.DONConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create nodes")
	}

	if config.EnvironmentHooks != nil {
		if hookErr := config.EnvironmentHooks.PostNodeStartupHooks(don, &config); hookErr != nil {
			return nil, errors.Wrapf(hookErr, "failed to run post node startup hooks")
		}
	}

	nodeIDs, keysErr := FetchNodeIds(don)
	if keysErr != nil {
		return nil, errors.Wrapf(keysErr, "failed to fetch nodeIds")
	}

	mocks, err := NewMocks(config.DONConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create mocks")
	}

	if config.EnvironmentHooks != nil {
		if hookErr := config.EnvironmentHooks.PostMocksStartupHooks(mocks, &config); hookErr != nil {
			return nil, errors.Wrapf(hookErr, "failed to run post mocks hooks")
		}
	}

	err = AppendMocksToDONConfig(don, mocks)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to append mocks to don config")
	}

	return &deployment.Environment{
		Name:     "persistent",
		Offchain: NewJobClient(don.ChainlinkClients),
		NodeIDs:  nodeIDs,
		Chains:   chains,
		Logger:   lggr,
	}, nil
}
