package v1_6

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

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
		func(address common.Address, backend bind.ContractBackend, opts *bind.TransactOpts, input RouterApplyRampUpdatesOpInput) (opsutil.EVMCallOutputWithError, error) {
			router, err := router.NewRouter(address, backend)
			if err != nil {
				return opsutil.EVMCallOutputWithError{}, fmt.Errorf("failed to create Router contract instance: %w", err)
			}
			tx, callErr := router.ApplyRampUpdates(opts, input.OnRampUpdates, input.OffRampRemoves, input.OffRampAdds)
			return opsutil.EVMCallOutputWithError{
				CallErr: callErr,
				EVMCallOutput: opsutil.EVMCallOutput{
					Tx:           tx,
					ContractType: shared.Router,
				},
			}, nil
		},
	)
)
