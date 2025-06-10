package ops

import (
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/mcms/sdk/evm/bindings"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type OpEVMGrantRoleDeps struct {
	Chain cldf_evm.Chain
}

type OpEVMGrantRoleInput struct {
	TimelockAddress    common.Address `json:"timelockAddress"` // Address of the EVM Timelock contract
	Address            common.Address `json:"address"`         // Address to grant the role to
	IsDeployerKeyAdmin bool           `json:"isDeployerKeyAdmin"`
	RoleID             [32]byte       `json:"roleID"`
}

type OpEVMGrantRoleOutput struct {
	MCMSTx mcmstypes.Transaction `json:"mcmsTx"`
}

var OpEVMGrantRole = operations.NewOperation(
	"evm-timelock-grant-role",
	semver.MustParse("1.0.0"),
	"Grants specified role to the ManyChainMultiSig contract on the EVM Timelock contract",
	func(b operations.Bundle, deps OpEVMGrantRoleDeps, input OpEVMGrantRoleInput) (OpEVMGrantRoleOutput, error) {
		txOpts := cldf.SimTransactOpts()
		if input.IsDeployerKeyAdmin {
			txOpts = deps.Chain.DeployerKey
		}

		timelock, err := bindings.NewRBACTimelock(input.TimelockAddress, deps.Chain.Client)
		if err != nil {
			b.Logger.Errorw("Failed to create Timelock instance",
				"chainSelector", deps.Chain.ChainSelector(),
				"chainName", deps.Chain.Name(),
				"contractAddr", input.Address.String(),
				"err", err,
			)
			return OpEVMGrantRoleOutput{}, err
		}

		grantRoleTx, err := timelock.GrantRole(
			txOpts, input.RoleID, input.Address,
		)
		if input.IsDeployerKeyAdmin {
			if _, err2 := cldf.ConfirmIfNoErrorWithABI(deps.Chain, grantRoleTx, bindings.RBACTimelockABI, err); err2 != nil {
				b.Logger.Errorw("Failed to grant timelock role",
					"chain", deps.Chain.Name(),
					"timelock", timelock.Address().Hex(),
					"Address to grant role", input.Address.Hex(),
					"TxHash", grantRoleTx.Hash().Hex(),
					"err", err2)
				return OpEVMGrantRoleOutput{}, err2
			}
			return OpEVMGrantRoleOutput{}, err
		}
		if err != nil {
			b.Logger.Errorw("Failed to grant timelock role",
				"chain", deps.Chain.Name(),
				"timelock", timelock.Address().Hex(),
				"Address to grant role", input.Address.Hex(),
				"err", err)
			return OpEVMGrantRoleOutput{}, err
		}

		// Create MCMS Tx
		mcmsTx, err := proposalutils.TransactionForChain(deps.Chain.Selector, timelock.Address().Hex(), grantRoleTx.Data(),
			big.NewInt(0), commontypes.RBACTimelock.String(), []string{})
		if err != nil {
			b.Logger.Errorw("Failed to create transaction for chain",
				"chain", deps.Chain.Name(),
				"timelock", timelock.Address().Hex(),
				"Address to grant role", input.Address.Hex(),
				"err", err)
			return OpEVMGrantRoleOutput{}, err
		}

		return OpEVMGrantRoleOutput{
			MCMSTx: mcmsTx,
		}, nil
	})
