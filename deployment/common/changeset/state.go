package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

// MCMSWithTimelockState holds the Go bindings
// for a MCMSWithTimelock contract deployment.
// It is public for use in product specific packages.
// Either all fields are nil or all fields are non-nil.
type MCMSWithTimelockState struct {
	*proposalutils.MCMSWithTimelockContracts
}

func MaybeLoadMCMSWithTimelockState(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (*MCMSWithTimelockState, error) {
	contracts, err := proposalutils.MaybeLoadMCMSWithTimelockContracts(chain, addresses)
	if err != nil {
		return nil, err
	}

	return &MCMSWithTimelockState{
		MCMSWithTimelockContracts: contracts,
	}, nil
}

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

type LinkTokenState struct {
	LinkToken *link_token.LinkToken
}

func (s LinkTokenState) GenerateLinkView() (v1_0.LinkTokenView, error) {
	if s.LinkToken == nil {
		return v1_0.LinkTokenView{}, errors.New("link token not found")
	}
	return v1_0.GenerateLinkTokenView(s.LinkToken)
}

func MaybeLoadLinkTokenState(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (*LinkTokenState, error) {
	state := LinkTokenState{}
	linkToken := deployment.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0)
	// Perhaps revisit if we have a use case for multiple.
	_, err := deployment.AddressesContainBundle(addresses, map[deployment.TypeAndVersion]struct{}{linkToken: {}})
	if err != nil {
		return nil, fmt.Errorf("unable to check link token on chain %s error: %w", chain.Name(), err)
	}
	for address, tvStr := range addresses {
		switch tvStr {
		case linkToken:
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

func (s StaticLinkTokenState) GenerateStaticLinkView() (v1_0.StaticLinkTokenView, error) {
	if s.StaticLinkToken == nil {
		return v1_0.StaticLinkTokenView{}, errors.New("static link token not found")
	}
	return v1_0.GenerateStaticLinkTokenView(s.StaticLinkToken)
}

func MaybeLoadStaticLinkTokenState(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (*StaticLinkTokenState, error) {
	state := StaticLinkTokenState{}
	staticLinkToken := deployment.NewTypeAndVersion(types.StaticLinkToken, deployment.Version1_0_0)
	// Perhaps revisit if we have a use case for multiple.
	_, err := deployment.AddressesContainBundle(addresses, map[deployment.TypeAndVersion]struct{}{staticLinkToken: {}})
	if err != nil {
		return nil, fmt.Errorf("unable to check static link token on chain %s error: %w", chain.Name(), err)
	}
	for address, tvStr := range addresses {
		switch tvStr {
		case staticLinkToken:
			lt, err := link_token_interface.NewLinkToken(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.StaticLinkToken = lt
		}
	}
	return &state, nil
}
