package state

import (
	"context"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	view "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
)

type PDASeed [32]byte

// MCMSWithTimelockProgramsSolana holds the solana publick keys
// and seeds for the MCM, AccessController and Timelock programs.
// It is public for use in product specific packages.
type MCMSWithTimelockProgramsSolana struct {
	McmProgram                       solana.PublicKey
	ProposerMcmSeed                  PDASeed
	CancellerMcmSeed                 PDASeed
	BypasserMcmSeed                  PDASeed
	TimelockProgram                  solana.PublicKey
	TimelockSeed                     PDASeed
	AccessControllerProgram          solana.PublicKey
	ProposerAccessControllerAccount  solana.PublicKey
	ExecutorAccessControllerAccount  solana.PublicKey
	CancellerAccessControllerAccount solana.PublicKey
	BypasserAccessControllerAccount  solana.PublicKey
}

// Validate checks that all fields are non-nil, ensuring it's ready
// for use generating views or interactions.
func (s *MCMSWithTimelockProgramsSolana) Validate() error {
	if s.McmProgram.IsZero() {
		return errors.New("mcm program not found")
	}
	if s.TimelockProgram.IsZero() {
		return errors.New("timelock program not found")
	}
	if s.AccessControllerProgram.IsZero() {
		return errors.New("access controller program not found")
	}
	if s.ProposerAccessControllerAccount.IsZero() {
		return errors.New("proposer access controller account not found")
	}
	if s.ExecutorAccessControllerAccount.IsZero() {
		return errors.New("executor access controller account not found")
	}
	if s.CancellerAccessControllerAccount.IsZero() {
		return errors.New("canceller access controller account not found")
	}
	if s.BypasserAccessControllerAccount.IsZero() {
		return errors.New("bypasser access controller account not found")
	}
	return nil
}

func (s *MCMSWithTimelockProgramsSolana) GenerateView(
	ctx context.Context, chain cldf_solana.Chain,
) (view.MCMSWithTimelockViewSolana, error) {
	if err := s.Validate(); err != nil {
		return view.MCMSWithTimelockViewSolana{}, fmt.Errorf("unable to validate state: %w", err)
	}

	inspector := mcmssolanasdk.NewInspector(chain.Client)
	timelockInspector := mcmssolanasdk.NewTimelockInspector(chain.Client)

	return view.GenerateMCMSWithTimelockViewSolana(ctx, inspector, timelockInspector, s.McmProgram,
		s.ProposerMcmSeed, s.CancellerMcmSeed, s.BypasserMcmSeed, s.TimelockProgram, s.TimelockSeed)
}

// MCMSWithTimelockStateSolana holds the Go bindings
// for a MCMSWithTimelock contract deployment.
// It is public for use in product specific packages.
// Either all fields are nil or all fields are non-nil.
type MCMSWithTimelockStateSolana struct {
	*MCMSWithTimelockProgramsSolana
}
