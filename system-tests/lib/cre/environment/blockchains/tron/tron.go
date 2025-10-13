package tron

import (
	"errors"
	"strconv"
	"strings"

	pkgerrors "github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type Deployer struct {
	provider infra.Provider
}

func NewDeployer(provider *infra.Provider) *Deployer {
	return &Deployer{
		provider: *provider,
	}
}

func (t *Deployer) Deploy(input *blockchain.Input) (*cre.Blockchain, error) {
	if t.provider.IsCRIB() {
		return nil, errors.New("CRIB deployment for Tron is not supported yet")
	}

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

	chainID, err := strconv.ParseUint(bcOut.ChainID, 10, 64)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to parse chain id %s", bcOut.ChainID)
	}
	selector, err := chainselectors.SelectorFromChainId(chainID)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "failed to get chain selector for chain id %s", bcOut.ChainID)
	}

	// if jsonrpc is not present, add it
	if !strings.HasSuffix(bcOut.Nodes[0].ExternalHTTPUrl, "/jsonrpc") {
		bcOut.Nodes[0].ExternalHTTPUrl += "/jsonrpc"
	}
	if !strings.HasSuffix(bcOut.Nodes[0].InternalHTTPUrl, "/jsonrpc") {
		bcOut.Nodes[0].InternalHTTPUrl += "/jsonrpc"
	}

	externalHTTPURL := bcOut.Nodes[0].ExternalHTTPUrl
	internalHTTPURL := bcOut.Nodes[0].InternalHTTPUrl

	return &cre.Blockchain{
		ChainSelector: selector,
		ChainID:       chainID,
		CtfOutput: &blockchain.Output{
			ChainID: bcOut.ChainID,
			Family:  blockchain.FamilyTron,
			Nodes: []*blockchain.Node{
				{
					InternalHTTPUrl: internalHTTPURL,
					ExternalHTTPUrl: externalHTTPURL,
				},
			},
		},
		SethClient:         nil,
		DeployerPrivateKey: blockchain.TRONAccounts.PrivateKeys[0],
	}, nil
}
