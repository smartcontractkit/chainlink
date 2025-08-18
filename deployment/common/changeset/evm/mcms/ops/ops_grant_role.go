package ops

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/mcms/sdk/evm/bindings"

	"github.com/smartcontractkit/chainlink/deployment/common/opsutils"

	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

type OpEVMGrantRoleInput struct {
	Account common.Address `json:"account"`
	RoleID  [32]byte       `json:"roleID"`
}

var OpEVMGrantRole = opsutils.NewEVMCallOperation(
	"evm-timelock-grant-role",
	semver.MustParse("1.0.0"),
	"Grants specified role to the ManyChainMultiSig contract on the EVM Timelock contract",
	bindings.RBACTimelockABI,
	commontypes.RBACTimelock,
	bindings.NewRBACTimelock,
	func(timelock *bindings.RBACTimelock, opts *bind.TransactOpts, input OpEVMGrantRoleInput) (*types.Transaction, error) {
		return timelock.GrantRole(opts, input.RoleID, input.Account)
	},
)
