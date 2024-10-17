package devenv

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/pelletier/go-toml/v2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

const (
	DevEnv = "devenv"
)

type EnvironmentConfig struct {
	Chains               []ChainConfig `toml:",omitempty"`
	HomeChainSelectorStr string        `toml:"HomeChainSelector,omitempty"`
	FeedChainSelectorStr string        `toml:"FeedChainSelector,omitempty"`
	JDConfig             JDConfig      `toml:",omitempty"`
}

func (c EnvironmentConfig) HomeChainSelector() (uint64, error) {
	return strconv.ParseUint(c.HomeChainSelectorStr, 10, 64)
}

func (c EnvironmentConfig) FeedChainSelector() (uint64, error) {
	return strconv.ParseUint(c.FeedChainSelectorStr, 10, 64)
}

func LoadEnvironmentConfig(path string, privateKeys []string) (EnvironmentConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return EnvironmentConfig{}, fmt.Errorf("file does not exist: %w", err)
	}
	cBytes, err := os.ReadFile(path)
	if err != nil {
		return EnvironmentConfig{}, fmt.Errorf("error reading environment config: %w", err)
	}
	var config EnvironmentConfig
	err = toml.Unmarshal(cBytes, &config)
	if err != nil {
		return config, fmt.Errorf("failed to decode environment config: %w", err)
	}
	// if no private keys are provided, set deployer key using KMS
	if privateKeys == nil {
		for i := range config.Chains {
			err := config.Chains[i].SetDeployerKey(nil)
			if err != nil {
				return EnvironmentConfig{}, fmt.Errorf("failed to set deployer key for chain id %d, have you set KMS env vars or provided private key? : %w", config.Chains[i].ChainID, err)
			}
		}
		return config, nil
	}
	if len(config.Chains) != len(privateKeys) {
		return EnvironmentConfig{}, fmt.Errorf("number of private keys %d does not match number of chains %d", len(privateKeys), len(config.Chains))
	}
	for i, key := range privateKeys {
		err := config.Chains[i].SetDeployerKey(&key)
		if err != nil {
			return EnvironmentConfig{}, fmt.Errorf("failed to set deployer key for chain id %d: %w", config.Chains[i].ChainID, err)
		}
	}
	return config, nil
}

func NewEnvironment(ctx context.Context, lggr logger.Logger, config EnvironmentConfig) (*deployment.Environment, *DON, error) {
	chains, err := NewChains(lggr, config.Chains)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create chains: %w", err)
	}
	var nodeIDs []string
	var offChain deployment.OffchainClient
	var don *DON
	if !config.JDConfig.IsEmpty() {
		offChain, err = NewJDClient(ctx, config.JDConfig)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create JD client: %w", err)
		}

		jd, ok := offChain.(*JobDistributor)
		if !ok {
			return nil, nil, fmt.Errorf("offchain client does not implement JobDistributor")
		}
		if jd == nil {
			return nil, nil, fmt.Errorf("offchain client is nil")
		}
		if jd.don != nil {
			err = jd.don.CreateSupportedChains(ctx, config.Chains, *jd)
			if err != nil {
				return nil, nil, err
			}
			nodeIDs = jd.don.NodeIds()
			don = jd.don
		}
	}

	return &deployment.Environment{
		Name:     DevEnv,
		Offchain: offChain,
		NodeIDs:  nodeIDs,
		Chains:   chains,
		Logger:   lggr,
	}, don, nil
}
