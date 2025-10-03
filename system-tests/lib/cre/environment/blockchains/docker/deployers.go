package docker

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

func NewDeployerSet() map[blockchain.ChainFamily]creblockchains.Deployer {
	return map[blockchain.ChainFamily]creblockchains.Deployer{
		blockchain.FamilyEVM:    &EVMDeployer{},
		blockchain.FamilySolana: &SolanaDeployer{},
		blockchain.FamilyTron:   &TronDeployer{},
	}
}
