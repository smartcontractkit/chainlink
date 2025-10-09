package ops

import (
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

type OpEVMDeployTimelockInput struct {
	TimelockMinDelay *big.Int         `json:"timelockMinDelay"`
	Admin            common.Address   `json:"admin"`      // Admin of the timelock contract, usually the deployer key
	Proposers        []common.Address `json:"proposers"`  // Proposer of the timelock contract, usually the deployer key
	Executors        []common.Address `json:"executors"`  // Executor of the timelock contract, usually the call proxy
	Cancellers       []common.Address `json:"cancellers"` // Canceller of the timelock contract, usually the deployer key
	Bypassers        []common.Address `json:"bypassers"`  // Bypasser of the timelock contract, usually the deployer key
}

var OpEVMDeployTimelock = opsutils.NewEVMDeployOperation(
	"evm-timelock-deploy",
	semver.MustParse("1.0.0"),
	"Deploys Timelock contract on the specified EVM chains",
	commontypes.RBACTimelock,
	bindings.RBACTimelockMetaData,
	&opsutils.ContractOpts{
		Version:     &deployment.Version1_0_0,
		EVMBytecode: common.FromHex(bindings.RBACTimelockBin),
		// ZkSyncVMBytecode not supported
	},
	func(input OpEVMDeployTimelockInput) []any {
		return []any{
			input.TimelockMinDelay,
			input.Admin,
			input.Proposers,
			input.Executors,
			input.Cancellers,
			input.Bypassers,
		}
	},
)
