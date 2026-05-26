package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
)

// MCMSWithTimelockState holds the Go bindings
// for a MCMSWithTimelock contract deployment.
// It is public for use in product specific packages.
// Either all fields are nil or all fields are non-nil.
// Deprecated: use MCMSWithTimelockState from deployment/common/changeset/state/evm.go instead
// if you are changing this, please make the similar changes in deployment/common/changeset/state
type MCMSWithTimelockState struct {
	*cldfproposalutils.MCMSWithTimelockContracts
}

// Deprecated: use GenerateMCMSWithTimelockView from deployment/common/changeset/state/evm.go instead
// if you are changing this, please make the similar changes in deployment/common/changeset/state
func (state MCMSWithTimelockState) GenerateMCMSWithTimelockView() (v1_0.MCMSWithTimelockView, error) {
	if err := state.Validate(); err != nil {
		return v1_0.MCMSWithTimelockView{}, err
	}
	timelockView, err := v1_0.GenerateTimelockView(*state.Timelock)
	if err != nil {
		return v1_0.MCMSWithTimelockView{}, nil
	}
	callProxyView, err := v1_0.GenerateCallProxyView(*state.CallProxy)
	if err != nil {
		return v1_0.MCMSWithTimelockView{}, nil
	}
	bypasserView, err := v1_0.GenerateMCMSView(*state.BypasserMcm)
	if err != nil {
		return v1_0.MCMSWithTimelockView{}, nil
	}
	proposerView, err := v1_0.GenerateMCMSView(*state.ProposerMcm)
	if err != nil {
		return v1_0.MCMSWithTimelockView{}, nil
	}
	cancellerView, err := v1_0.GenerateMCMSView(*state.CancellerMcm)
	if err != nil {
		return v1_0.MCMSWithTimelockView{}, nil
	}
	return v1_0.MCMSWithTimelockView{
		Timelock:  timelockView,
		Bypasser:  bypasserView,
		Proposer:  proposerView,
		Canceller: cancellerView,
		CallProxy: callProxyView,
	}, nil
}

// MaybeLoadMCMSWithTimelockState loads the MCMSWithTimelockState state for each chain in the given environment.
// Deprecated: use MaybeLoadMCMSWithTimelockState from deployment/common/changeset/state/evm.go instead
// if you are changing this, please make the similar changes in deployment/common/changeset/state
func MaybeLoadMCMSWithTimelockState(env cldf.Environment, chainSelectors []uint64) (map[uint64]*MCMSWithTimelockState, error) {
	result := map[uint64]*MCMSWithTimelockState{}
	for _, chainSelector := range chainSelectors {
		chain, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return nil, fmt.Errorf("chain %d not found", chainSelector)
		}
		addressesChain, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return nil, err
		}
		state, err := MaybeLoadMCMSWithTimelockChainState(chain, addressesChain)
		if err != nil {
			return nil, err
		}
		result[chainSelector] = state
	}
	return result, nil
}

// MaybeLoadMCMSWithTimelockChainState looks for the addresses corresponding to
// contracts deployed with DeployMCMSWithTimelock and loads them into a
// MCMSWithTimelockState struct.  If none of the contracts are found, the state struct will be nil.
//
// An error indicates:
// - Found but was unable to load a contract
// - It only found part of the bundle of contracts
// - If found more than one instance of a contract (we expect one bundle in the given addresses)
// Deprecated: use MaybeLoadMCMSWithTimelockChainState from deployment/common/changeset/state/evm.go instead
// if you are changing this, please make the similar changes in deployment/common/changeset/state
func MaybeLoadMCMSWithTimelockChainState(
	chain cldf_evm.Chain,
	addresses map[string]cldf.TypeAndVersion,
) (*MCMSWithTimelockState, error) {
	var (
		state = MCMSWithTimelockState{
			MCMSWithTimelockContracts: &cldfproposalutils.MCMSWithTimelockContracts{},
		}

		// We expect one of each contract on the chain.
		timelock  = cldf.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0)
		callProxy = cldf.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0)
		proposer  = cldf.NewTypeAndVersion(types.ProposerManyChainMultisig, deployment.Version1_0_0)
		canceller = cldf.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0)
		bypasser  = cldf.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0)

		// the same contract can have different roles
		multichain    = cldf.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
		proposerMCMS  = cldf.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
		bypasserMCMS  = cldf.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
		cancellerMCMS = cldf.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
	)

	// Convert map keys to a slice
	proposerMCMS.Labels.Add(types.ProposerRole.String())
	bypasserMCMS.Labels.Add(types.BypasserRole.String())
	cancellerMCMS.Labels.Add(types.CancellerRole.String())
	wantTypes := []cldf.TypeAndVersion{timelock, proposer, canceller, bypasser, callProxy,
		proposerMCMS, bypasserMCMS, cancellerMCMS,
	}

	// Ensure we either have the bundle or not.
	_, err := cldf.EnsureDeduped(addresses, wantTypes)
	if err != nil {
		return nil, fmt.Errorf("unable to check MCMS contracts on chain %s error: %w", chain.Name(), err)
	}

	for address, tv := range addresses {
		switch {
		case tv.Type == timelock.Type && tv.Version.String() == timelock.Version.String():
			tl, err := owner_helpers.NewRBACTimelock(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.Timelock = tl
		case tv.Type == callProxy.Type && tv.Version.String() == callProxy.Version.String():
			cp, err := owner_helpers.NewCallProxy(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.CallProxy = cp
		case tv.Type == proposer.Type && tv.Version.String() == proposer.Version.String():
			mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.ProposerMcm = mcms
		case tv.Type == bypasser.Type && tv.Version.String() == bypasser.Version.String():
			mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.BypasserMcm = mcms
		case tv.Type == canceller.Type && tv.Version.String() == canceller.Version.String():
			mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.CancellerMcm = mcms
		case tv.Type == multichain.Type && tv.Version.String() == multichain.Version.String():
			// Contract of type ManyChainMultiSig must be labeled to assign to the proper state
			// field.  If a specifically typed contract already occupies the field, then this
			// contract will be ignored.
			mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			if tv.Labels.Contains(types.ProposerRole.String()) && state.ProposerMcm == nil {
				state.ProposerMcm = mcms
			}
			if tv.Labels.Contains(types.BypasserRole.String()) && state.BypasserMcm == nil {
				state.BypasserMcm = mcms
			}
			if tv.Labels.Contains(types.CancellerRole.String()) && state.CancellerMcm == nil {
				state.CancellerMcm = mcms
			}
		}
	}
	return &state, nil
}
