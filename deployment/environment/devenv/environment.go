package devenv

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
)

const (
	DevEnv = "devenv"
)

type EnvironmentConfig struct {
	Chains   []ChainConfig
	JDConfig JDConfig
}

func NewEnvironment(ctx func() context.Context, lggr logger.Logger, config EnvironmentConfig) (*deployment.Environment, *DON, error) {
	chains, err := NewChains(lggr, config.Chains)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create chains: %w", err)
	}

	var nodeIDs []string
	var offChain deployment.OffchainClient
	var don *DON
	if !config.JDConfig.IsEmpty() {
		offChain, err := NewJDClient(ctx(), config.JDConfig)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create JD client: %w", err)
		}

		jd, ok := offChain.(*JobDistributor)
		if !ok {
			return nil, nil, fmt.Errorf("offchain client does not implement JobDistributor")
		}
		if jd == nil {
			return nil, nil, fmt.Errorf("offchain client is not set up")
		}
		if jd.don != nil {
			err = jd.don.CreateSupportedChains(ctx(), config.Chains, *jd)
			if err != nil {
				return nil, nil, err
			}
			nodeIDs = jd.don.NodeIds()
			don = jd.don
		}
	}

	return deployment.NewEnvironment(
		DevEnv,
		lggr,
		deployment.NewMemoryAddressBook(),
		chains,
		nodeIDs,
		offChain,
		ctx,
	), don, nil
}
