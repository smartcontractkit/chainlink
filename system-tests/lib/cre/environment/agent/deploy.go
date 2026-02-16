package agent

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

type OutputDeployer interface {
	DeployOutput(ctx context.Context, input *blockchain.Input) (*blockchain.Output, error)
}

func DeployBlockchainComponent(
	ctx context.Context,
	deployers map[blockchain.ChainFamily]blockchains.Deployer,
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

	if outputDeployer, ok := deployer.(OutputDeployer); ok {
		deployedOutput, err := outputDeployer.DeployOutput(ctx, input)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain output of type %s", input.Type)
		}
		return deployedOutput, nil
	}

	deployed, err := deployer.Deploy(ctx, input)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain of type %s", input.Type)
	}

	return deployed.CtfOutput(), nil
}
