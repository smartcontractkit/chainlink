package blockchains

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
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

// ValidateKubernetesBlockchainOutput validates that the blockchain output is configured for Kubernetes
// Returns an error if output is nil or missing nodes, nil otherwise
func ValidateKubernetesBlockchainOutput(input *blockchain.Input) error {
	if input.Out == nil || len(input.Out.Nodes) == 0 {
		return fmt.Errorf("kubernetes provider requires blockchain URLs to be configured in config file for blockchain type %s chainID: %s", input.Type, input.ChainID)
	}
	return nil
}

func Start(
	ctx context.Context,
	testLogger zerolog.Logger,
	commonLogger logger.Logger,
	inputs []*blockchain.Input,
	deployers map[blockchain.ChainFamily]Deployer,
) (*DeployedBlockchains, error) {
	outputs := make([]Blockchain, len(inputs))

	eg, egCtx := errgroup.WithContext(ctx)
	for i := range inputs {
		eg.Go(func() error {
			in := inputs[i]

			chainFamily, chErr := blockchain.TypeToFamily(in.Type)
			if chErr != nil {
				return chErr
			}

			deployer, ok := deployers[chainFamily]
			if !ok {
				if err := framework.PrintFailedContainerLogs(30); err != nil {
					testLogger.Error().Err(err).Msg("failed to print failed Docker container logs")
				}
				return fmt.Errorf("no deployer found for blockchain type %s", in.Type)
			}

			deployedBlockchain, deployErr := deployer.Deploy(egCtx, in)
			if deployErr != nil {
				return pkgerrors.Wrapf(deployErr, "failed to deploy blockchain of type %s", in.Type)
			}

			outputs[i] = deployedBlockchain
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
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
