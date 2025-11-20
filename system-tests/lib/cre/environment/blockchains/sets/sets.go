package sets

import (
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/tron"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func NewDeployerSet(testLogger zerolog.Logger, provider *infra.Provider, cribConfigsDir string) map[blockchain.ChainFamily]blockchains.Deployer {
	return map[blockchain.ChainFamily]blockchains.Deployer{
		blockchain.FamilyEVM:    evm.NewDeployer(testLogger, provider, cribConfigsDir),
		blockchain.FamilySolana: solana.NewDeployer(testLogger, provider),
		blockchain.FamilyTron:   tron.NewDeployer(testLogger, provider),
	}
}
