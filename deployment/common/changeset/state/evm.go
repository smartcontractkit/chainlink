package state

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	view "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
)

// MCMSWithTimelockState holds the Go bindings
// for a MCMSWithTimelock contract deployment.
// It is public for use in product specific packages.
// Either all fields are nil or all fields are non-nil.
type MCMSWithTimelockState struct {
	CancellerMcm *bindings.ManyChainMultiSig
	BypasserMcm  *bindings.ManyChainMultiSig
	ProposerMcm  *bindings.ManyChainMultiSig
	Timelock     *bindings.RBACTimelock
	CallProxy    *bindings.CallProxy
}

// Validate checks that all fields are non-nil, ensuring it's ready
// for use generating views or interactions.
func (state MCMSWithTimelockState) Validate() error {
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

func (state MCMSWithTimelockState) GenerateMCMSWithTimelockView() (view.MCMSWithTimelockView, error) {
	if err := state.Validate(); err != nil {
		return view.MCMSWithTimelockView{}, fmt.Errorf("unable to validate McmsWithTimelock state: %w", err)
	}

	return view.GenerateMCMSWithTimelockView(*state.BypasserMcm, *state.CancellerMcm, *state.ProposerMcm,
		*state.Timelock, *state.CallProxy)
}

// MaybeLoadMCMSWithTimelockState loads the MCMSWithTimelockState state for each chain in the given environment.
func MaybeLoadMCMSWithTimelockState(env deployment.Environment, chainSelectors []uint64) (map[uint64]*MCMSWithTimelockState, error) {
	result := map[uint64]*MCMSWithTimelockState{}
	for _, chainSelector := range chainSelectors {
		chain, ok := env.Chains[chainSelector]
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
// MCMSWithTimelockState struct. If none of the contracts are found, the state struct will be nil.
// An error indicates:
// - Found but was unable to load a contract
// - It only found part of the bundle of contracts
// - If found more than one instance of a contract (we expect one bundle in the given addresses)
func MaybeLoadMCMSWithTimelockChainState(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (*MCMSWithTimelockState, error) {
	state := MCMSWithTimelockState{}
	var (
		// We expect one of each contract on the chain.
		timelock  = deployment.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0)
		callProxy = deployment.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0)
		proposer  = deployment.NewTypeAndVersion(types.ProposerManyChainMultisig, deployment.Version1_0_0)
		canceller = deployment.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0)
		bypasser  = deployment.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0)

		// the same contract can have different roles
		multichain    = deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
		proposerMCMS  = deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
		bypasserMCMS  = deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
		cancellerMCMS = deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
	)

	// Convert map keys to a slice
	proposerMCMS.Labels.Add(types.ProposerRole.String())
	bypasserMCMS.Labels.Add(types.BypasserRole.String())
	cancellerMCMS.Labels.Add(types.CancellerRole.String())
	wantTypes := []deployment.TypeAndVersion{timelock, proposer, canceller, bypasser, callProxy,
		proposerMCMS, bypasserMCMS, cancellerMCMS,
	}

	// Ensure we either have the bundle or not.
	_, err := deployment.EnsureDeduped(addresses, wantTypes)
	if err != nil {
		return nil, fmt.Errorf("unable to check MCMS contracts on chain %s error: %w", chain.Name(), err)
	}

	for address, tv := range addresses {
		switch {
		case tv.Type == timelock.Type && tv.Version.String() == timelock.Version.String():
			tl, err := bindings.NewRBACTimelock(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.Timelock = tl
		case tv.Type == callProxy.Type && tv.Version.String() == callProxy.Version.String():
			cp, err := bindings.NewCallProxy(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.CallProxy = cp
		case tv.Type == proposer.Type && tv.Version.String() == proposer.Version.String():
			mcms, err := bindings.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.ProposerMcm = mcms
		case tv.Type == bypasser.Type && tv.Version.String() == bypasser.Version.String():
			mcms, err := bindings.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.BypasserMcm = mcms
		case tv.Type == canceller.Type && tv.Version.String() == canceller.Version.String():
			mcms, err := bindings.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.CancellerMcm = mcms
		case tv.Type == multichain.Type && tv.Version.String() == multichain.Version.String():
			// Contract of type ManyChainMultiSig must be labeled to assign to the proper state
			// field.  If a specifically typed contract already occupies the field, then this
			// contract will be ignored.
			mcms, err := bindings.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
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

type LinkTokenState struct {
	LinkToken *link_token.LinkToken
}

func (s LinkTokenState) GenerateLinkView() (view.LinkTokenView, error) {
	if s.LinkToken == nil {
		return view.LinkTokenView{}, errors.New("link token not found")
	}
	return view.GenerateLinkTokenView(s.LinkToken)
}

// MaybeLoadLinkTokenState loads the LinkTokenState state for each chain in the given environment.
func MaybeLoadLinkTokenState(env deployment.Environment, chainSelectors []uint64) (map[uint64]*LinkTokenState, error) {
	result := map[uint64]*LinkTokenState{}
	for _, chainSelector := range chainSelectors {
		chain, ok := env.Chains[chainSelector]
		if !ok {
			return nil, fmt.Errorf("chain %d not found", chainSelector)
		}
		addressesChain, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return nil, err
		}
		state, err := MaybeLoadLinkTokenChainState(chain, addressesChain)
		if err != nil {
			return nil, err
		}
		result[chainSelector] = state
	}
	return result, nil
}

func MaybeLoadLinkTokenChainState(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (*LinkTokenState, error) {
	state := LinkTokenState{}
	linkToken := deployment.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0)

	// Convert map keys to a slice
	wantTypes := []deployment.TypeAndVersion{linkToken}

	// Ensure we either have the bundle or not.
	_, err := deployment.EnsureDeduped(addresses, wantTypes)
	if err != nil {
		return nil, fmt.Errorf("unable to check link token on chain %s error: %w", chain.Name(), err)
	}

	for address, tvStr := range addresses {
		if tvStr.Type == linkToken.Type && tvStr.Version.String() == linkToken.Version.String() {
			lt, err := link_token.NewLinkToken(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.LinkToken = lt
		}
	}
	return &state, nil
}

type StaticLinkTokenState struct {
	StaticLinkToken *link_token_interface.LinkToken
}

func (s StaticLinkTokenState) GenerateStaticLinkView() (view.StaticLinkTokenView, error) {
	if s.StaticLinkToken == nil {
		return view.StaticLinkTokenView{}, errors.New("static link token not found")
	}
	return view.GenerateStaticLinkTokenView(s.StaticLinkToken)
}

func MaybeLoadStaticLinkTokenState(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (*StaticLinkTokenState, error) {
	state := StaticLinkTokenState{}
	staticLinkToken := deployment.NewTypeAndVersion(types.StaticLinkToken, deployment.Version1_0_0)

	// Convert map keys to a slice
	wantTypes := []deployment.TypeAndVersion{staticLinkToken}

	// Ensure we either have the bundle or not.
	_, err := deployment.EnsureDeduped(addresses, wantTypes)
	if err != nil {
		return nil, fmt.Errorf("unable to check static link token on chain %s error: %w", chain.Name(), err)
	}

	for address, tvStr := range addresses {
		if tvStr.Type == staticLinkToken.Type && tvStr.Version.String() == staticLinkToken.Version.String() {
			lt, err := link_token_interface.NewLinkToken(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.StaticLinkToken = lt
		}
	}
	return &state, nil
}
