package docker

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

func NewDeployerSet(commonLogger logger.Logger) map[creblockchains.ChainFamily]creblockchains.Deployer {
	return map[creblockchains.ChainFamily]creblockchains.Deployer{
		blockchain.FamilyEVM:    NewEVMDeployer(commonLogger),
		blockchain.FamilySolana: NewSolanaDeployer(commonLogger),
		blockchain.FamilyTron:   NewTronDeployer(commonLogger),
	}
}
