package v1_6

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

type RouterApplyRampUpdatesOpInput struct {
	OnRampUpdates  []router.RouterOnRamp
	OffRampRemoves []router.RouterOffRamp
	OffRampAdds    []router.RouterOffRamp
}

var (
	RouterApplyRampUpdatesOp = opsutil.NewEVMCallOperation(
		"RouterApplyRampUpdatesOp",
		semver.MustParse("1.0.0"),
		"Updates OnRamps and OffRamps on the Router contract",
		router.RouterABI,
		shared.Router,
		router.NewRouter,
		func(router *router.Router, opts *bind.TransactOpts, input RouterApplyRampUpdatesOpInput) (*types.Transaction, error) {
			return router.ApplyRampUpdates(opts, input.OnRampUpdates, input.OffRampRemoves, input.OffRampAdds)
		},
	)
)
