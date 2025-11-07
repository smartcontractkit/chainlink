package blockchains

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

type Blockchain interface {
	ChainSelector() uint64
	ChainID() uint64
	ChainFamily() string
	IsFamily(chainFamily string) bool

	Fund(ctx context.Context, address string, amount uint64) error

	CtfOutput() *blockchain.Output
	ToCldfChain() (cldf_chain.BlockChain, error)
}

type Deployer interface {
	Deploy(ctx context.Context, input *blockchain.Input) (Blockchain, error)
}

type DeployedBlockchains struct {
	Outputs         []Blockchain
	CldfBlockChains cldf_chain.BlockChains
}

func (s *DeployedBlockchains) RegistryChain() Blockchain {
	return s.Outputs[0]
}

func Start(
	ctx context.Context,
	testLogger zerolog.Logger,
	commonLogger logger.Logger,
	inputs []*blockchain.Input,
	deployers map[blockchain.ChainFamily]Deployer,
) (*DeployedBlockchains, error) {
	outputs := make([]Blockchain, 0, len(inputs))

	for _, input := range inputs {
		chainFamily, chErr := blockchain.TypeToFamily(input.Type)
		if chErr != nil {
			return nil, chErr
		}

		deployer, ok := deployers[chainFamily]
		if !ok {
			infra.PrintFailedContainerLogs(testLogger, 30)
			return nil, fmt.Errorf("no deployer found for blockchain type %s", input.Type)
		}

		deployedBlockchain, deployErr := deployer.Deploy(ctx, input)
		if deployErr != nil {
			return nil, pkgerrors.Wrapf(deployErr, "failed to deploy blockchain of type %s", input.Type)
		}

		outputs = append(outputs, deployedBlockchain)
	}

	cldfBlockchains := make([]cldf_chain.BlockChain, 0, len(outputs))
	for _, db := range outputs {
		chain, chainErr := db.ToCldfChain()
		if chainErr != nil {
			return nil, pkgerrors.Wrap(chainErr, "failed to create cldf chain from blockchain")
		}
		cldfBlockchains = append(cldfBlockchains, chain)
	}

	return &DeployedBlockchains{
		Outputs:         outputs,
		CldfBlockChains: cldf_chain.NewBlockChainsFromSlice(cldfBlockchains),
	}, nil
}
