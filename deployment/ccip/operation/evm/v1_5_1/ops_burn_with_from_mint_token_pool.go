package v1_5_1

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_with_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

var (
	DeployBurnWithFromMintTokenPoolOp = opsutil.NewEVMDeployOperation(
		"DeployBurnWithFromMintTokenPool",
		semver.MustParse("1.0.0"),
		"Deploys BurnWithFromMintTokenPool 1.5.1 contract on the specified evm chain",
		shared.BurnWithFromMintTokenPool,
		burn_with_from_mint_token_pool.BurnWithFromMintTokenPoolMetaData,
		&opsutil.ContractOpts{
			Version:          &deployment.Version1_5_1,
			EVMBytecode:      common.FromHex(burn_with_from_mint_token_pool.BurnWithFromMintTokenPoolBin),
			ZkSyncVMBytecode: burn_with_from_mint_token_pool.ZkBytecode,
		},
		func(input DeployBurnMintTokenPoolInput) []any { // input format is the same as BurnMintTokenPool
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
