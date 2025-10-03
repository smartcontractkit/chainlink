package docker

import (
	"maps"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

type TronDeployer struct {
	commonLogger logger.Logger
}

func NewTronDeployer(commonLogger logger.Logger) *TronDeployer {
	return &TronDeployer{
		commonLogger: commonLogger,
	}
}

func (t *TronDeployer) Deploy(input *blockchain.Input) (*creblockchains.DeployedBLockchain, error) {
	bcOut, err := blockchain.NewBlockchainNetwork(input)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain %s chainID: %s", input.Type, input.ChainID)
	}

	w, wrapErr := creblockchains.WrapTron(input, bcOut)
	if wrapErr != nil {
		return nil, pkgerrors.Wrap(wrapErr, "failed to wrap Tron")
	}

	chainsConfigs := make([]devenv.ChainConfig, 0)
	cfg, cfgErr := cre.ChainConfigFromWrapped(w)
	if cfgErr != nil {
		return nil, pkgerrors.Wrap(cfgErr, "failed to wrap blockchain output to chain config")
	}
	chainsConfigs = append(chainsConfigs, cfg)

	cldfBlockchain, err := devenv.NewChains(t.commonLogger, chainsConfigs)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create chains")
	}

	return &creblockchains.DeployedBLockchain{
		CldfBlockchain: maps.Collect(cldfBlockchain.All())[w.ChainSelector],
		Blockchain:     w,
	}, nil
}
