package ops

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	zkbindings "github.com/smartcontractkit/mcms/sdk/zksync/bindings"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

type OpEVMDeployCallProxyInput struct {
	Timelock common.Address `json:"timelock"`
}

var OpEVMDeployCallProxy = opsutils.NewEVMDeployOperation(
	"evm-call-proxy-deploy",
	semver.MustParse("1.0.0"),
	"Deploys CallProxy contract on the specified EVM chains",
	commontypes.CallProxy,
	bindings.CallProxyMetaData,
	&opsutils.ContractOpts{
		Version:          &deployment.Version1_0_0,
		EVMBytecode:      common.FromHex(bindings.CallProxyBin),
		ZkSyncVMBytecode: zkbindings.CallProxyZkBytecode,
	},
	func(input OpEVMDeployCallProxyInput) []any {
		return []any{
			input.Timelock,
		}
	},
)
