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
	NodeIDs     []string
	Chains      []devenv.ChainConfig   // chain selector -> Chain Config
	AddressBook deployment.AddressBook // Addresses of all contracts
}

type DeployCCIPOutput struct {
	AddressBook deployment.AddressBookMap
	NodeIDs     []string
}

func NewDeployEnvironmentFromCribOutput(lggr logger.Logger, output DeployOutput) (*deployment.Environment, error) {
	chains, err := devenv.NewChains(lggr, output.Chains)
	if err != nil {
		return nil, err
	}
	return deployment.NewEnvironment(
		CRIB_ENV_NAME,
		lggr,
		output.AddressBook,
		chains,
		nil, // nil for solana chains, can use memory solana chain example when required
		output.NodeIDs,
		nil, // todo: populate the offchain client using output.DON
		func() context.Context { return context.Background() }, deployment.XXXGenerateTestOCRSecrets(),
	), nil
}
