package crib

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"

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
	blockChains, err := devenv.NewChains(lggr, output.Chains)
	if err != nil {
		return nil, err
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
		focr.XXXGenerateTestOCRSecrets(),
		blockChains,
	), nil
}
