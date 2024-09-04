package devenv

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/ccip/integration-tests/deployment"
	csav1 "github.com/smartcontractkit/ccip/integration-tests/deployment/jd/csa/v1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	DevEnv = "devenv"
)

type EnvironmentConfig struct {
	Chains   []ChainConfig
	DON      DON
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

	keypairs, err := offChain.ListKeypairs(ctx, &csav1.ListKeypairsRequest{})
	if err != nil {
		return nil, err
	}

	nodes := NewNodes(t, chains, config.Nodes, config.Bootstraps, config.RegistryConfig)
	var nodeIDs []string
	for id := range nodes {
		nodeIDs = append(nodeIDs, id)
	}

	return &deployment.Environment{
		Name:     DevEnv,
		Offchain: offChain,
		NodeIDs:  nodeIDs,
		Chains:   chains,
		Logger:   lggr,
	}, nil
}
