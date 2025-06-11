package crib

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

const (
	CRIB_ENV_NAME = "Crib Environment"
)

type DeployOutput struct {
	NodeIDs     []string
	Chains      []devenv.ChainConfig // chain selector -> Chain Config
	AddressBook cldf.AddressBook     // Addresses of all contracts
}

type DeployCCIPOutput struct {
	AddressBook cldf.AddressBookMap
	NodeIDs     []string
}

func NewDeployEnvironmentFromCribOutput(lggr logger.Logger, output DeployOutput) (*cldf.Environment, error) {
	chains, solChains, err := devenv.NewChains(lggr, output.Chains)
	if err != nil {
		return nil, err
	}

	blockChains := map[uint64]chain.BlockChain{}
	for _, c := range chains {
		blockChains[c.Selector] = c
	}
	for _, c := range solChains {
		blockChains[c.Selector] = c
	}

	return cldf.NewEnvironment(
		CRIB_ENV_NAME,
		lggr,
		output.AddressBook,
		datastore.NewMemoryDataStore().Seal(),
		output.NodeIDs,
		nil, // todo: populate the offchain client using output.DON
		//nolint:gocritic // intentionally use a lambda to allow dynamic context replacement in Environment Commit 90ee880
		func() context.Context { return context.Background() },
		cldf.XXXGenerateTestOCRSecrets(),
		chain.NewBlockChains(blockChains),
	), nil
}
