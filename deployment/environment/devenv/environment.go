package devenv

import (
	"context"
	"errors"
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
	offChain, err := NewJDClient(ctx(), config.JDConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create JD client: %w", err)
	}

	jd, ok := offChain.(*JobDistributor)
	if !ok {
		return nil, nil, errors.New("offchain client does not implement JobDistributor")
	}
	if jd == nil {
		return nil, nil, errors.New("offchain client is not set up")
	}
	var nodeIDs []string
	if jd.don != nil {
		// Gateway DON doesn't require any chain setup, and trying to create chains for it will fail,
		// because its nodes are missing chain-related configuration. Of course, we could add that configuration,
		// but its not how it is setup on production.
		if len(config.Chains) > 0 {
			err = jd.don.CreateSupportedChains(ctx(), config.Chains, *jd)
			if err != nil {
				return nil, nil, err
			}
		}
		nodeIDs = jd.don.NodeIds()
	}

	return deployment.NewEnvironment(
		DevEnv,
		lggr,
		deployment.NewMemoryAddressBook(),
		chains,
		nil, // sending nil for solana chains right now, we can build this when we need it
		nodeIDs,
		offChain,
		ctx,
		deployment.XXXGenerateTestOCRSecrets(),
	), jd.don, nil
}
