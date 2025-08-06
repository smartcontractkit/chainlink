package testhelpers

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-ton/ops/ccip/config"
	toncs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/ton"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/xssnick/tonutils-go/address"
)

const CHAINSEL_EVM_TEST_90000001 = 909606746561742123

func DeployChainContractsToTonCS(t *testing.T, e DeployedEnv, chainSelector uint64) commoncs.ConfiguredChangeSet {
	// TODO update the contract config
	ccipConfig := toncs.DeployCCIPContractsCfg{
		TonChainSelector: chainSelector,
		ChainContractParams: config.ChainContractParams{
			FeeQuoterParams: config.FeeQuoterParams{
				MaxFeeJuelsPerMsg:                    big.NewInt(1),
				TokenPriceStalenessThreshold:         0,
				FeeTokens:                            []*address.Address{},
				PremiumMultiplierWeiPerEthByFeeToken: map[shared.TokenSymbol]uint64{},
			},
			OffRampParams: config.OffRampParams{
				// ...
			},
			OnRampParams: config.OnRampParams{
				ChainSelector: CHAINSEL_EVM_TEST_90000001,
				// TODO:
				// AllowlistAdmin: &address.Address{},
				FeeAggregator: e.Env.BlockChains.TonChains()[0].WalletAddress,
			},
		},
	}
	return commoncs.Configure(toncs.DeployCCIPContracts{}, ccipConfig)
}

func AddLaneTONChangesets(e *DeployedEnv, from, to uint64, fromFamily, toFamily string) commoncs.ConfiguredChangeSet {
	laneConfig := toncs.AddLaneCfg{
		FromChainSelector: from,
		ToChainSelector:   to,
		FromFamily:        fromFamily,
		ToFamily:          toFamily,
	}
	return commoncs.Configure(toncs.AddLane{}, laneConfig)
}
