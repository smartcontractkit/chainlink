package v1_5_1

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

type DeployBurnMintTokenPoolInput struct {
	TokenAddress    common.Address   `json:"token"`
	Decimals        uint8            `json:"localTokenDecimals"`
	AllowList       []common.Address `json:"allowlist"`
	RMNProxyAddress common.Address   `json:"rmnProxy"`
	RouterAddress   common.Address   `json:"router"`
}

var (
	DeployBurnMintTokenPoolOp = opsutil.NewEVMDeployOperation(
		"DeployBurnMintTokenPool",
		semver.MustParse("1.0.0"),
		"Deploys BurnMintTokenPool 1.5.1 contract on the specified evm chain",
		shared.BurnMintTokenPool,
		burn_mint_token_pool.BurnMintTokenPoolMetaData,
		&opsutil.ContractOpts{
			Version:          &deployment.Version1_5_1,
			EVMBytecode:      common.FromHex(burn_mint_token_pool.BurnMintTokenPoolBin),
			ZkSyncVMBytecode: burn_mint_token_pool.ZkBytecode,
		},
		func(input DeployBurnMintTokenPoolInput) []any {
			return []any{
				input.TokenAddress,
				input.Decimals,
				input.AllowList,
				input.RMNProxyAddress,
				input.RouterAddress,
			}
		},
	)
)
