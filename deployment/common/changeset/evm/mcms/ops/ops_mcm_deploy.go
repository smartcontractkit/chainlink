package ops

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"

	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

type OpEVMDeployMCMOutput struct {
	Address common.Address `json:"address"`
}

var OpEVMDeployProposerMCM = opsutils.NewEVMDeployOperation(
	"evm-proposer-mcm-deploy",
	semver.MustParse("1.0.0"),
	"Deploys Proposer MCM contract",
	commontypes.ProposerManyChainMultisig,
	bindings.ManyChainMultiSigMetaData,
	&opsutils.ContractOpts{
		Version:     &deployment.Version1_0_0,
		EVMBytecode: common.FromHex(bindings.ManyChainMultiSigBin),
		// ZkSyncVMBytecode not supported
	},
	func(input any) []any {
		return []any{}
	},
)

var OpEVMDeployBypasserMCM = opsutils.NewEVMDeployOperation(
	"evm-bypasser-mcm-deploy",
	semver.MustParse("1.0.0"),
	"Deploys Bypasser MCM contract",
	commontypes.BypasserManyChainMultisig,
	bindings.ManyChainMultiSigMetaData,
	&opsutils.ContractOpts{
		Version:     &deployment.Version1_0_0,
		EVMBytecode: common.FromHex(bindings.ManyChainMultiSigBin),
		// ZkSyncVMBytecode not supported
	},
	func(input any) []any {
		return []any{}
	},
)

var OpEVMDeployCancellerMCM = opsutils.NewEVMDeployOperation(
	"evm-canceller-mcm-deploy",
	semver.MustParse("1.0.0"),
	"Deploys Canceller MCM contract",
	commontypes.CancellerManyChainMultisig,
	bindings.ManyChainMultiSigMetaData,
	&opsutils.ContractOpts{
		Version:     &deployment.Version1_0_0,
		EVMBytecode: common.FromHex(bindings.ManyChainMultiSigBin),
		// ZkSyncVMBytecode not supported
	},
	func(input any) []any {
		return []any{}
	},
)
