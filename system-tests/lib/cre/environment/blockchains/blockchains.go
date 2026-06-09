package blockchains

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

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
	Start(ctx context.Context, input *blockchain.Input) (*blockchain.Output, error)
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

func StartChain(
	ctx context.Context,
	deployers map[blockchain.ChainFamily]Deployer,
	input *blockchain.Input,
) (*blockchain.Output, error) {
	if input == nil {
		return nil, pkgerrors.New("blockchain input is nil")
	}

	chainFamily, err := blockchain.TypeToFamily(input.Type)
	if err != nil {
		return nil, err
	}

	deployer, ok := deployers[chainFamily]
	if !ok {
		return nil, fmt.Errorf("no deployer found for blockchain type %s", input.Type)
	}
	deployed, err := deployer.Start(ctx, input)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain of type %s", input.Type)
	}
	return deployed, nil
}
