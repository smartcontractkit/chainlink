package state

import (
	"errors"
	"fmt"

	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

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
