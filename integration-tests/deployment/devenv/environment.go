package devenv

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	DevEnv = "devenv"
)

type EnvironmentConfig struct {
	Chains   []ChainConfig
	nodeInfo []NodeInfo
	JDConfig JDConfig
}

func NewEnvironment(ctx context.Context, lggr logger.Logger, config EnvironmentConfig) (*deployment.Environment, error) {
	chains, err := NewChains(lggr, config.Chains)
	if err != nil {
		return nil, fmt.Errorf("failed to create chains: %w", err)
	}
	offChain, err := NewJDClient(config.JDConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create JD client: %w", err)
	}

	jd, ok := offChain.(JobDistributor)
	if !ok {
		return nil, fmt.Errorf("offchain client does not implement JobDistributor")
	}
	don, err := NewRegisteredDON(ctx, config.nodeInfo, jd)
	if err != nil {
		return nil, fmt.Errorf("failed to create registered DON: %w", err)
	}
	nodeIDs := don.NodeIds()

	err = don.CreateSupportedChains(ctx, config.Chains)
	if err != nil {
		return nil, err
	}

	return &deployment.Environment{
		Name:     DevEnv,
		Offchain: offChain,
		NodeIDs:  nodeIDs,
		Chains:   chains,
		Logger:   lggr,
	}, nil
}
