package changeset

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

// TimelockExecutionContracts is a helper struct for executing timelock proposals. it contains
// the timelock and call proxy contracts.
type TimelockExecutionContracts struct {
	Timelock  *owner_helpers.RBACTimelock
	CallProxy *owner_helpers.CallProxy
}

// NewTimelockExecutionContracts creates a new TimelockExecutionContracts struct.
// If there are multiple timelocks or call proxy on the chain, an error is returned.
// If there is a missing timelocks or call proxy on the chain, an error is returned.
func NewTimelockExecutionContracts(env deployment.Environment, chainSelector uint64) (*TimelockExecutionContracts, error) {
	addrTypeVer, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	if err != nil {
		return nil, fmt.Errorf("error getting addresses for chain: %w", err)
	}
	var timelock *owner_helpers.RBACTimelock
	var callProxy *owner_helpers.CallProxy
	for addr, tv := range addrTypeVer {
		if tv.Type == types.RBACTimelock {
			if timelock != nil {
				return nil, fmt.Errorf("multiple timelocks found on chain %d", chainSelector)
			}
			var err error
			timelock, err = owner_helpers.NewRBACTimelock(common.HexToAddress(addr), env.Chains[chainSelector].Client)
			if err != nil {
				return nil, fmt.Errorf("error creating timelock: %w", err)
			}
		}
		if tv.Type == types.CallProxy {
			if callProxy != nil {
				return nil, fmt.Errorf("multiple call proxies found on chain %d", chainSelector)
			}
			var err error
			callProxy, err = owner_helpers.NewCallProxy(common.HexToAddress(addr), env.Chains[chainSelector].Client)
			if err != nil {
				return nil, fmt.Errorf("error creating call proxy: %w", err)
			}
		}
	}
	if timelock == nil || callProxy == nil {
		return nil, fmt.Errorf("missing timelock (%T) or call proxy(%T) on chain %d", timelock == nil, callProxy == nil, chainSelector)
	}
	return &TimelockExecutionContracts{
		Timelock:  timelock,
		CallProxy: callProxy,
	}, nil
}

// RunTimelockExecutor executes all the operation in the given executor on the given chain.
// It is an error if there are no operations for the given chain.
func RunTimelockExecutor(env deployment.Environment, executor *mcms.Executor, timelockContracts *TimelockExecutionContracts, sel uint64) error {
	// TODO: This sort of helper probably should move to the MCMS lib.
	// Execute all the transactions in the proposal which are for this chain.
	if len(executor.Operations[mcms.ChainIdentifier(sel)]) == 0 {
		return fmt.Errorf("no operations for chain %d", sel)
	}
	for _, chainOp := range executor.Operations[mcms.ChainIdentifier(sel)] {
		for idx, op := range executor.ChainAgnosticOps {
			if bytes.Equal(op.Data, chainOp.Data) && op.To == chainOp.To {
				opTx, err2 := executor.ExecuteOnChain(env.Chains[sel].Client, env.Chains[sel].DeployerKey, idx)
				if err2 != nil {
					return fmt.Errorf("error executing on chain: %w", err2)
				}
				block, err2 := env.Chains[sel].Confirm(opTx)
				if err2 != nil {
					return fmt.Errorf("error confirming on chain: %w", err2)
				}
				it, err2 := timelockContracts.Timelock.FilterCallScheduled(&bind.FilterOpts{
					Start:   block,
					End:     &block,
					Context: env.GetContext(),
				}, nil, nil)
				if err2 != nil {
					return fmt.Errorf("error filtering call scheduled: %w", err2)
				}
				var calls []owner_helpers.RBACTimelockCall
				var pred, salt [32]byte
				for it.Next() {
					// Note these are the same for the whole batch, can overwrite
					pred = it.Event.Predecessor
					salt = it.Event.Salt
					env.Logger.Info("scheduled", "event", it.Event)
					calls = append(calls, owner_helpers.RBACTimelockCall{
						Target: it.Event.Target,
						Data:   it.Event.Data,
						Value:  it.Event.Value,
					})
				}

				timelockExecutorProxy, err := owner_helpers.NewRBACTimelock(timelockContracts.CallProxy.Address(), env.Chains[sel].Client)
				if err != nil {
					return fmt.Errorf("error creating timelock executor proxy: %w", err)
				}
				tx, err := timelockExecutorProxy.ExecuteBatch(
					env.Chains[sel].DeployerKey, calls, pred, salt)
				if err != nil {
					return fmt.Errorf("error executing batch: %w", err)
				}
				_, err = env.Chains[sel].Confirm(tx)
				if err != nil {
					return fmt.Errorf("error confirming batch: %w", err)
				}
			}
		}
	}
	return nil
}
