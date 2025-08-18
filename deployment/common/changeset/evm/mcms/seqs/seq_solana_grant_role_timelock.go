package seqs

import (
	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	timelockbindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/timelock"
	cldfsolana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/evm/mcms/ops"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

type SeqSolanaGrantRoleTimelockDeps struct {
	Chain cldfsolana.Chain
}

type SeqSolanaGrantRoleTimelockInput struct {
	ChainState         *state.MCMSWithTimelockStateSolana `json:"chainState"`
	Role               timelockbindings.Role              `json:"role"`
	Accounts           []solana.PublicKey                 `json:"accounts"`
	IsDeployerKeyAdmin bool                               `json:"isDeployerKeyAdmin"`
}

type SeqSolanaGrantRoleTimelockOutput struct {
	McmsTransactions []mcmsTypes.Transaction `json:"mcmsTxs"`
}

var SeqSolanaGrantRoleTimelock = operations.NewSequence(
	"seq-solana-grant-role-timelock",
	semver.MustParse("1.0.0"),
	"Grant a role to multiple accounts in a Solana Timelock instance",
	func(b operations.Bundle, deps SeqSolanaGrantRoleTimelockDeps, in SeqSolanaGrantRoleTimelockInput) (SeqSolanaGrantRoleTimelockOutput, error) {
		mcmsTxs := make([]mcmsTypes.Transaction, 0, len(in.Accounts))

		for _, account := range in.Accounts {
			opReport, err := operations.ExecuteOperation(b, ops.OpSolanaGrantRoleTimelock,
				ops.OpSolanaGrantRoleTimelockDeps{
					Chain: deps.Chain,
				},
				ops.OpSolanaGrantRoleTimelockInput{
					ChainState:         in.ChainState,
					Role:               in.Role,
					Account:            account,
					IsDeployerKeyAdmin: in.IsDeployerKeyAdmin,
				},
			)
			if err != nil {
				b.Logger.Errorw("Failed to grant role", "chainSelector", deps.Chain.ChainSelector(), "chainName", deps.Chain.Name(),
					"timelock", state.EncodeAddressWithSeed(in.ChainState.TimelockProgram, in.ChainState.TimelockSeed),
					"role", in.Role, "account", account)
				return SeqSolanaGrantRoleTimelockOutput{}, err
			}

			mcmsTxs = append(mcmsTxs, opReport.Output.MCMSTransaction)
		}

		return SeqSolanaGrantRoleTimelockOutput{McmsTransactions: mcmsTxs}, nil
	},
)
