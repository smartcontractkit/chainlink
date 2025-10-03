package kubernetes

import (
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

func NewDeployerSet(commonLogger logger.Logger, testLogger zerolog.Logger, namespace string, cribConfigsDir string) map[creblockchains.ChainFamily]creblockchains.Deployer {
	return map[creblockchains.ChainFamily]creblockchains.Deployer{
		blockchain.FamilyEVM: NewEVMDeployer(testLogger, commonLogger, namespace, cribConfigsDir),
	}
}
