package v1_5_1

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

type DeployLockReleaseTokenPoolInput struct {
	TokenAddress    common.Address   `json:"token"`
	Decimals        uint8            `json:"localTokenDecimals"`
	AllowList       []common.Address `json:"allowlist"`
	RMNProxyAddress common.Address   `json:"rmnProxy"`
	AcceptLiquidity bool             `json:"acceptliquidity"`
	RouterAddress   common.Address   `json:"router"`
}

var (
	DeployLockReleaseTokenPoolOp = opsutil.NewEVMDeployOperation(
		"DeployLockReleaseTokenPool",
		semver.MustParse("1.0.0"),
		"Deploys LockReleaseTokenPool 1.5.1 contract on the specified evm chain",
		shared.LockReleaseTokenPool,
		lock_release_token_pool.LockReleaseTokenPoolMetaData,
		&opsutil.ContractOpts{
			Version:          &deployment.Version1_5_1,
			EVMBytecode:      common.FromHex(lock_release_token_pool.LockReleaseTokenPoolBin),
			ZkSyncVMBytecode: lock_release_token_pool.ZkBytecode,
		},
		func(input DeployLockReleaseTokenPoolInput) []any {
			return []any{
				input.TokenAddress,
				input.Decimals,
				input.AllowList,
				input.RMNProxyAddress,
				input.AcceptLiquidity,
				input.RouterAddress,
			}
		},
	)
)
