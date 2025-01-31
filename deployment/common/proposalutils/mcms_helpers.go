package proposalutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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

type RunTimelockExecutorConfig struct {
	Executor          *mcms.Executor
	TimelockContracts *TimelockExecutionContracts
	ChainSelector     uint64
	// BlockStart is optional. It filter the timelock scheduled events.
	// If not provided, the executor assumes that the operations have not been executed yet
	// executes all the operations for the given chain.
	BlockStart *uint64
	BlockEnd   *uint64
}

func (cfg RunTimelockExecutorConfig) Validate() error {
	if cfg.Executor == nil {
		return errors.New("executor is nil")
	}
	if cfg.TimelockContracts == nil {
		return errors.New("timelock contracts is nil")
	}
	if cfg.ChainSelector == 0 {
		return errors.New("chain selector is 0")
	}
	if cfg.BlockStart != nil && cfg.BlockEnd == nil {
		if *cfg.BlockStart > *cfg.BlockEnd {
			return errors.New("block start is greater than block end")
		}
	}
	if cfg.BlockStart == nil && cfg.BlockEnd != nil {
		return errors.New("block start must not be nil when block end is not nil")
	}

	if len(cfg.Executor.Operations[mcms.ChainIdentifier(cfg.ChainSelector)]) == 0 {
		return fmt.Errorf("no operations for chain %d", cfg.ChainSelector)
	}
	return nil
}

// RunTimelockExecutor runs the scheduled operations for the given chain.
// If the block start is not provided, it assumes that the operations have not been scheduled yet
// and executes all the operations for the given chain.
// It is an error if there are no operations for the given chain.
func RunTimelockExecutor(env deployment.Environment, cfg RunTimelockExecutorConfig) error {
	// TODO: This sort of helper probably should move to the MCMS lib.
	// Execute all the transactions in the proposal which are for this chain.
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("error validating config: %w", err)
	}
	for _, chainOp := range cfg.Executor.Operations[mcms.ChainIdentifier(cfg.ChainSelector)] {
		for idx, op := range cfg.Executor.ChainAgnosticOps {
			start := cfg.BlockStart
			end := cfg.BlockEnd
			if bytes.Equal(op.Data, chainOp.Data) && op.To == chainOp.To {
				if start == nil {
					opTx, err2 := cfg.Executor.ExecuteOnChain(env.Chains[cfg.ChainSelector].Client, env.Chains[cfg.ChainSelector].DeployerKey, idx)
					if err2 != nil {
						return fmt.Errorf("error executing on chain: %w", err2)
					}
					block, err2 := env.Chains[cfg.ChainSelector].Confirm(opTx)
					if err2 != nil {
						return fmt.Errorf("error confirming on chain: %w", err2)
					}
					start = &block
					end = &block
				}

				it, err2 := cfg.TimelockContracts.Timelock.FilterCallScheduled(&bind.FilterOpts{
					Start:   *start,
					End:     end,
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
					verboseDebug(env.Logger, it.Event)
					env.Logger.Info("scheduled", "event", it.Event)
					calls = append(calls, owner_helpers.RBACTimelockCall{
						Target: it.Event.Target,
						Data:   it.Event.Data,
						Value:  it.Event.Value,
					})
				}
				if len(calls) == 0 {
					return fmt.Errorf("no calls found for chain %d in blocks [%d, %d]", cfg.ChainSelector, *start, *end)
				}
				timelockExecutorProxy, err := owner_helpers.NewRBACTimelock(cfg.TimelockContracts.CallProxy.Address(), env.Chains[cfg.ChainSelector].Client)
				if err != nil {
					return fmt.Errorf("error creating timelock executor proxy: %w", err)
				}
				tx, err := timelockExecutorProxy.ExecuteBatch(
					env.Chains[cfg.ChainSelector].DeployerKey, calls, pred, salt)
				if err != nil {
					return fmt.Errorf("error executing batch: %w", err)
				}
				_, err = env.Chains[cfg.ChainSelector].Confirm(tx)
				if err != nil {
					return fmt.Errorf("error confirming batch: %w", err)
				}
			}
		}
	}
	return nil
}

func verboseDebug(lggr logger.Logger, event *owner_helpers.RBACTimelockCallScheduled) {
	b, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	lggr.Debug("scheduled", "event", string(b))
}

// MCMSWithTimelockContracts holds the Go bindings
// for a MCMSWithTimelock contract deployment.
// It is public for use in product specific packages.
// Either all fields are nil or all fields are non-nil.
type MCMSWithTimelockContracts struct {
	CancellerMcm *owner_helpers.ManyChainMultiSig
	BypasserMcm  *owner_helpers.ManyChainMultiSig
	ProposerMcm  *owner_helpers.ManyChainMultiSig
	Timelock     *owner_helpers.RBACTimelock
	CallProxy    *owner_helpers.CallProxy
}

// Validate checks that all fields are non-nil, ensuring it's ready
// for use generating views or interactions.
func (state MCMSWithTimelockContracts) Validate() error {
	if state.Timelock == nil {
		return errors.New("timelock not found")
	}
	if state.CancellerMcm == nil {
		return errors.New("canceller not found")
	}
	if state.ProposerMcm == nil {
		return errors.New("proposer not found")
	}
	if state.BypasserMcm == nil {
		return errors.New("bypasser not found")
	}
	if state.CallProxy == nil {
		return errors.New("call proxy not found")
	}
	return nil
}

// MaybeLoadMCMSWithTimelockContracts looks for the addresses corresponding to
// contracts deployed with DeployMCMSWithTimelock and loads them into a
// MCMSWithTimelockState struct. If none of the contracts are found, the state struct will be nil.
// An error indicates:
// - Found but was unable to load a contract
// - It only found part of the bundle of contracts
// - If found more than one instance of a contract (we expect one bundle in the given addresses)
func MaybeLoadMCMSWithTimelockContracts(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (*MCMSWithTimelockContracts, error) {
	state := MCMSWithTimelockContracts{}
	// We expect one of each contract on the chain.
	timelock := deployment.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0)
	callProxy := deployment.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0)
	proposer := deployment.NewTypeAndVersion(types.ProposerManyChainMultisig, deployment.Version1_0_0)
	canceller := deployment.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0)
	bypasser := deployment.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0)

	// Convert map keys to a slice
	wantTypes := []deployment.TypeAndVersion{timelock, proposer, canceller, bypasser, callProxy}
	_, err := deployment.AddressesContainBundle(addresses, wantTypes)
	if err != nil {
		return nil, fmt.Errorf("unable to check MCMS contracts on chain %s error: %w", chain.Name(), err)
	}

	for address, tvStr := range addresses {
		switch {
		case tvStr.Type == timelock.Type && tvStr.Version.String() == timelock.Version.String():
			tl, err := owner_helpers.NewRBACTimelock(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.Timelock = tl
		case tvStr.Type == callProxy.Type && tvStr.Version.String() == callProxy.Version.String():
			cp, err := owner_helpers.NewCallProxy(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.CallProxy = cp
		case tvStr.Type == proposer.Type && tvStr.Version.String() == proposer.Version.String():
			mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.ProposerMcm = mcms
		case tvStr.Type == bypasser.Type && tvStr.Version.String() == bypasser.Version.String():
			mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.BypasserMcm = mcms
		case tvStr.Type == canceller.Type && tvStr.Version.String() == canceller.Version.String():
			mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.CancellerMcm = mcms
		}
	}
	return &state, nil
}
