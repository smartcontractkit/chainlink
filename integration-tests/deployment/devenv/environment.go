package devenv

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

const (
	DevEnv = "devenv"
)

type EnvironmentConfig struct {
	Chains            []ChainConfig
	HomeChainSelector uint64
	FeedChainSelector uint64
	JDConfig          JDConfig
}

func NewEnvironment(ctx context.Context, lggr logger.Logger, config EnvironmentConfig) (*deployment.Environment, *DON, error) {
	chains, err := NewChains(lggr, config.Chains)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create chains: %w", err)
	}
	offChain, err := NewJDClient(ctx, config.JDConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create JD client: %w", err)
	}

	jd, ok := offChain.(*JobDistributor)
	if !ok {
		return nil, nil, fmt.Errorf("offchain client does not implement JobDistributor")
	}
	if jd == nil || jd.don == nil {
		return nil, nil, fmt.Errorf("offchain client does not have a DON")
	}

	err = jd.don.CreateSupportedChains(ctx, config.Chains, *jd)
	if err != nil {
		return nil, nil, err
	}
	nodeIDs := jd.don.NodeIds()

	return &deployment.Environment{
		Name:     DevEnv,
		Offchain: offChain,
		NodeIDs:  nodeIDs,
		Chains:   chains,
		Logger:   lggr,
	}, jd.don, nil
}
