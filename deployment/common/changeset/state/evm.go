package state

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/link_token"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

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
func MaybeLoadMCMSWithTimelockState(env cldf.Environment, chainSelectors []uint64) (map[uint64]*MCMSWithTimelockState, error) {
	return MaybeLoadMCMSWithTimelockStateWithQualifier(env, chainSelectors, "")
}

// MaybeLoadMCMSWithTimelockStateWithQualifier loads the MCMSWithTimelockState state for each chain in the given environment,
// supporting qualifiers for filtering addresses. This uses the merged approach searching both AddressBook and DataStore.
func MaybeLoadMCMSWithTimelockStateWithQualifier(env cldf.Environment, chainSelectors []uint64, qualifier string) (map[uint64]*MCMSWithTimelockState, error) {
	result := map[uint64]*MCMSWithTimelockState{}
	for _, chainSelector := range chainSelectors {
		chain, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return nil, fmt.Errorf("chain %d not found", chainSelector)
		}

		// Use merged addresses from both AddressBook and DataStore for backward compatibility
		addressesChain, err := LoadAddressesFromDataStore(env.DataStore, chainSelector, qualifier)
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

// MaybeLoadMCMSWithTimelockStateDataStore loads the MCMSWithTimelockState state for each chain in the given environment from the DataStore.
func MaybeLoadMCMSWithTimelockStateDataStore(env cldf.Environment, chainSelectors []uint64) (map[uint64]*MCMSWithTimelockState, error) {
	return MaybeLoadMCMSWithTimelockStateDataStoreWithQualifier(env, chainSelectors, "")
}

func MaybeLoadMCMSWithTimelockStateDataStoreWithQualifier(env cldf.Environment, chainSelectors []uint64, qualifier string) (map[uint64]*MCMSWithTimelockState, error) {
	result := map[uint64]*MCMSWithTimelockState{}
	for _, chainSelector := range chainSelectors {
		chain, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return nil, fmt.Errorf("chain %d not found", chainSelector)
		}

		addressesChain, err := LoadAddressesFromDataStore(env.DataStore, chainSelector, qualifier)
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

// LoadAddressesFromDataStore loads addresses from DataStore with optional qualifier.
// This is a public utility function that can be used by other packages to avoid duplication.
func LoadAddressesFromDataStore(ds datastore.DataStore, chainSelector uint64, qualifier string) ([]datastore.AddressRef, error) {
	// Build filter list starting with chain selector
	filters := []datastore.FilterFunc[datastore.AddressRefKey, datastore.AddressRef]{datastore.AddressRefByChainSelector(chainSelector)}

	// Add qualifier filter if provided
	if qualifier != "" {
		filters = append(filters, datastore.AddressRefByQualifier(qualifier))
	}

	addresses := ds.Addresses().Filter(filters...)
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no addresses found for chain %d", chainSelector)
	}

	return addresses, nil
}

// MaybeLoadMCMSWithTimelockChainState looks for the addresses corresponding to
// contracts deployed with DeployMCMSWithTimelock and loads them into a
// MCMSWithTimelockState struct. If none of the contracts are found, the state struct will be nil.
// An error indicates:
// - Found but was unable to load a contract
// - It only found part of the bundle of contracts
// - If found more than one instance of a contract (we expect one bundle in the given addresses)
func MaybeLoadMCMSWithTimelockChainState(chain cldf_evm.Chain, addresses []datastore.AddressRef) (*MCMSWithTimelockState, error) {
	state := MCMSWithTimelockState{}
	var (
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

	for _, ref := range addresses {
		tv := cldf.NewTypeAndVersion(cldf.ContractType(ref.Type), *ref.Version)
		address := ref.Address
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

func MaybeLoadLinkTokenChainState(chain cldf_evm.Chain, addresses []datastore.AddressRef) (*LinkTokenState, error) {
	state := LinkTokenState{}
	linkToken := cldf.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0)

	found := false
	for _, ref := range addresses {
		if ref.Type == datastore.ContractType(linkToken.Type) &&
			ref.Version.String() == linkToken.Version.String() {
			if found {
				return nil, fmt.Errorf("multiple link tokens found on chain %s", chain.Name())
			}
			found = true
			lt, err := link_token.NewLinkToken(common.HexToAddress(ref.Address), chain.Client)
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

func MaybeLoadStaticLinkTokenState(chain cldf_evm.Chain, addresses []datastore.AddressRef) (*StaticLinkTokenState, error) {
	state := StaticLinkTokenState{}
	staticLinkToken := cldf.NewTypeAndVersion(types.StaticLinkToken, deployment.Version1_0_0)

	found := false
	for _, ref := range addresses {
		if ref.Type == datastore.ContractType(staticLinkToken.Type) &&
			ref.Version.String() == staticLinkToken.Version.String() {
			if found {
				return nil, fmt.Errorf("multiple static link tokens found on chain %s", chain.Name())
			}
			found = true
			lt, err := link_token_interface.NewLinkToken(common.HexToAddress(ref.Address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.StaticLinkToken = lt
		}
	}
	return &state, nil
}
