package ops

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"

	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	accesscontrollerbindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/access_controller"
	timelockbindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/timelock"
	cldfsolana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

type OpSolanaGrantRoleDeps struct {
	Chain cldfsolana.Chain
}

type OpSolanaGrantRolesInput struct {
	ChainState         *state.MCMSWithTimelockStateSolana `json:"chainState"`
	Role               timelockbindings.Role              `json:"role"`
	Accounts           []solana.PublicKey                 `json:"accounts"`
	IsDeployerKeyAdmin bool                               `json:"isDeployerKeyAdmin"`
}

type OpSolanaGrantRolesOutput struct {
	MCMSBatchOperation mcmstypes.BatchOperation `json:"mcmsBatchOperation"`
}

var OpSolanaGrantRoles = operations.NewOperation(
	"solana-grant-role",
	semver.MustParse("1.0.0"),
	"Grant Role in a Timelock instance",
	func(b operations.Bundle, deps OpSolanaGrantRoleDeps, in OpSolanaGrantRolesInput) (OpSolanaGrantRolesOutput, error) {
		accessController, err := selectAccessController(in)
		if err != nil {
			return OpSolanaGrantRolesOutput{}, fmt.Errorf("failed to select access controller: %w", err)
		}

		timelockbindings.SetProgramID(in.ChainState.TimelockProgram)
		accesscontrollerbindings.SetProgramID(in.ChainState.AccessControllerProgram)
		var signer solana.PublicKey
		if in.IsDeployerKeyAdmin {
			signer = deps.Chain.DeployerKey.PublicKey()
		} else {
			signer = state.GetTimelockSignerPDA(in.ChainState.TimelockProgram, in.ChainState.TimelockSeed)
		}

		transactions := make([]mcmstypes.Transaction, len(in.Accounts))
		for i, account := range in.Accounts {
			ix, err := accesscontrollerbindings.NewAddAccessInstruction(accessController, signer, account).ValidateAndBuild()
			if err != nil {
				return OpSolanaGrantRolesOutput{}, fmt.Errorf("failed to create update delay instruction: %w", err)
			}
			if in.IsDeployerKeyAdmin {
				cerr := deps.Chain.SendAndConfirm(b.GetContext(), []solana.Instruction{ix})
				if cerr != nil {
					return OpSolanaGrantRolesOutput{}, fmt.Errorf("failed to confirm instructions: %w", cerr)
				}
				continue
			}

			transactions[i], err = mcmssolanasdk.NewTransactionFromInstruction(ix, "AccessController", []string{})
			if err != nil {
				return OpSolanaGrantRolesOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
		}

		return OpSolanaGrantRolesOutput{
			MCMSBatchOperation: mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(deps.Chain.ChainSelector()),
				Transactions:  transactions,
			},
		}, nil
	},
)

func selectAccessController(in OpSolanaGrantRolesInput) (solana.PublicKey, error) {
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
