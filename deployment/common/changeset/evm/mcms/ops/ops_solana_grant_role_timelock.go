package ops

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"

	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	accesscontrollerbindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/access_controller"
	timelockbindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/timelock"
	cldfsolana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

type OpSolanaGrantRoleTimelockDeps struct {
	Chain cldfsolana.Chain
}

type OpSolanaGrantRoleTimelockInput struct {
	ChainState         *state.MCMSWithTimelockStateSolana `json:"chainState"`
	Role               timelockbindings.Role              `json:"role"`
	Account            solana.PublicKey                   `json:"account"`
	IsDeployerKeyAdmin bool                               `json:"isDeployerKeyAdmin"`
}

type OpSolanaGrantRoleTimelockOutput struct {
	MCMSTransaction mcmstypes.Transaction `json:"mcmsTransaction"`
}

var OpSolanaGrantRoleTimelock = operations.NewOperation(
	"solana-grant-role-timelock",
	semver.MustParse("1.0.0"),
	"Grant a role to an account in a Solana Timelock instance",
	func(b operations.Bundle, deps OpSolanaGrantRoleTimelockDeps, in OpSolanaGrantRoleTimelockInput) (OpSolanaGrantRoleTimelockOutput, error) {
		accessController, err := selectAccessController(in)
		if err != nil {
			return OpSolanaGrantRoleTimelockOutput{}, fmt.Errorf("failed to select access controller: %w", err)
		}

		timelockbindings.SetProgramID(in.ChainState.TimelockProgram)
		accesscontrollerbindings.SetProgramID(in.ChainState.AccessControllerProgram)
		var signer solana.PublicKey
		if in.IsDeployerKeyAdmin {
			signer = deps.Chain.DeployerKey.PublicKey()
		} else {
			signer = state.GetTimelockSignerPDA(in.ChainState.TimelockProgram, in.ChainState.TimelockSeed)
		}

		ix, err := accesscontrollerbindings.NewAddAccessInstruction(accessController, signer, in.Account).ValidateAndBuild()
		if err != nil {
			return OpSolanaGrantRoleTimelockOutput{}, fmt.Errorf("failed to create update delay instruction: %w", err)
		}

		if in.IsDeployerKeyAdmin {
			cerr := deps.Chain.SendAndConfirm(b.GetContext(), []solana.Instruction{ix})
			if cerr != nil {
				return OpSolanaGrantRoleTimelockOutput{}, fmt.Errorf("failed to confirm instructions: %w", cerr)
			}

			return OpSolanaGrantRoleTimelockOutput{}, nil
		}

		transaction, err := mcmssolanasdk.NewTransactionFromInstruction(ix, "AccessController", []string{})
		if err != nil {
			return OpSolanaGrantRoleTimelockOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}

		return OpSolanaGrantRoleTimelockOutput{MCMSTransaction: transaction}, nil
	},
)

func selectAccessController(in OpSolanaGrantRoleTimelockInput) (solana.PublicKey, error) {
	switch in.Role {
	case timelockbindings.Admin_Role:
		return solana.PublicKey{}, errors.New("admin role not supported")
	case timelockbindings.Proposer_Role:
		return in.ChainState.ProposerAccessControllerAccount, nil
	case timelockbindings.Executor_Role:
		return in.ChainState.ExecutorAccessControllerAccount, nil
	case timelockbindings.Canceller_Role:
		return in.ChainState.CancellerAccessControllerAccount, nil
	case timelockbindings.Bypasser_Role:
		return in.ChainState.BypasserAccessControllerAccount, nil
	default:
		return solana.PublicKey{}, fmt.Errorf("unknown role %s", in.Role)
	}
}
