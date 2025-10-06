package devenv

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
)

const (
	DevEnv = "devenv"
)

type EnvironmentConfig struct {
	Chains   []ChainConfig
	JDConfig JDConfig
}

func NewEnvironment(ctx func() context.Context, lggr logger.Logger, config EnvironmentConfig) (*cldf.Environment, *DON, error) {
	blockChains, err := NewChains(lggr, config.Chains)
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
		err = jd.don.CreateSupportedChains(ctx(), config.Chains, *jd)
		if err != nil {
			return nil, nil, err
		}
		nodeIDs = jd.don.NodeIds()
	}

	return cldf.NewEnvironment(
		DevEnv,
		lggr,
		cldf.NewMemoryAddressBook(),
		datastore.NewMemoryDataStore().Seal(),
		nodeIDs,
		offChain,
		ctx,
		focr.XXXGenerateTestOCRSecrets(),
		blockChains,
	), jd.don, nil
}
