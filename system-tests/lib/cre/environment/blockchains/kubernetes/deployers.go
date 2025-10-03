package kubernetes

import (
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

func NewDeployerSet(testLogger zerolog.Logger, namespace string, cribConfigsDir string) map[blockchain.ChainFamily]creblockchains.Deployer {
	return map[blockchain.ChainFamily]creblockchains.Deployer{
		blockchain.FamilyEVM: NewEVMDeployer(testLogger, namespace, cribConfigsDir),
	}
}
