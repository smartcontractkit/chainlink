package crib

import (
	"context"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

const (
	CRIB_ENV_NAME = "Crib Environment"
)

type DeployOutput struct {
	NodeIDs          []string
	Chains           []devenv.ChainConfig        // chain selector -> Chain Config
	AddressesByChain deployment.AddressesByChain // Addresses of all contracts
}

type DeployCCIPOutput struct {
	AddressesByChain deployment.AddressesByChain
	NodeIDs          []string
}

func NewDeployEnvironmentFromCribOutput(lggr logger.Logger, output DeployOutput) (*deployment.Environment, error) {
	chains, err := devenv.NewChains(lggr, output.Chains)
	if err != nil {
		return nil, err
	}
	addressBook := deployment.NewMemoryAddressBookFromMap(output.AddressesByChain)

	return deployment.NewEnvironment(
		CRIB_ENV_NAME,
		lggr,
		addressBook,
		chains,
		output.NodeIDs,
		nil, // todo: populate the offchain client using output.DON
		func() context.Context { return context.Background() },
	), nil
}
