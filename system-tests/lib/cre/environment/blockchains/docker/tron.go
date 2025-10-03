package docker

import (
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

type TronDeployer struct{}

func (t *TronDeployer) Deploy(input *blockchain.Input) (*cre.Blockchain, error) {
	var bcOut *blockchain.Output
	var err error

	if input.Out != nil {
		bcOut = input.Out
	} else {
		bcOut, err = blockchain.NewBlockchainNetwork(input)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, "failed to deploy blockchain %s chainID: %s", input.Type, input.ChainID)
		}
	}

	return creblockchains.WrapTron(input, bcOut)
}
