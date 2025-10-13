package blockchains

import (
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

type Deployer interface {
	Deploy(input *blockchain.Input) (*cre.Blockchain, error)
}

type DeployedBlockchains struct {
	Outputs         []*cre.Blockchain
	CldfBlockChains chain.BlockChains
}

func (s *DeployedBlockchains) RegistryChain() *cre.Blockchain {
	return s.Outputs[0]
}

func Start(
	commonLogger logger.Logger,
	inputs []*blockchain.Input,
	deployers map[blockchain.ChainFamily]Deployer,
) (*DeployedBlockchains, error) {
	outputs := make([]*cre.Blockchain, 0, len(inputs))

	for _, input := range inputs {
		chainFamily, chErr := blockchain.TypeToFamily(input.Type)
		if chErr != nil {
			return nil, chErr
		}

		deployer, ok := deployers[chainFamily]
		if !ok {
			return nil, fmt.Errorf("no deployer found for blockchain type %s", input.Type)
		}

		deployedBlockchain, deployErr := deployer.Deploy(input)
		if deployErr != nil {
			return nil, pkgerrors.Wrapf(deployErr, "failed to deploy blockchain of type %s", input.Type)
		}

		outputs = append(outputs, deployedBlockchain)
	}

	chainsConfigs := make([]cre.ChainConfig, 0, len(outputs))
	for _, db := range outputs {
		cfg, cfgErr := cre.ChainConfigFromWrapped(db)
		if cfgErr != nil {
			return nil, pkgerrors.Wrap(cfgErr, "failed to wrap blockchain output to chain config")
		}
		chainsConfigs = append(chainsConfigs, cfg)
	}

	cldfBlockchains, err := cre.NewChains(commonLogger, chainsConfigs)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create chains")
	}

	return &DeployedBlockchains{
		Outputs:         outputs,
		CldfBlockChains: cldfBlockchains,
	}, nil
}
