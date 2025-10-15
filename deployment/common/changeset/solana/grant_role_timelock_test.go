package solana

import (
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	timelockbindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/timelock"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func TestGrantRoleTimelockSolana(t *testing.T) {
	t.Skip("fails with Program is not deployed (DoajfR5tK24xVw51fWcawUZWhAXD8yrBJVacc13neVQA) in CI")
	t.Parallel()

	// --- arrange ---
	rt, selector := setupTest(t)

	chain := rt.Environment().BlockChains.SolanaChains()[selector]
	executors1 := randomSolanaAccounts(t, 2)
	executors2 := randomSolanaAccounts(t, 2)
	addresses, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)
	mcmsState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	require.NoError(t, err)

	fundSignerPDAs(t, chain, mcmsState)

	// validate initial executors
	inspector := mcmssolanasdk.NewTimelockInspector(chain.Client)
	onChainExecutors, err := inspector.GetExecutors(t.Context(), timelockAddress(mcmsState))
	require.NoError(t, err)
	require.ElementsMatch(t, onChainExecutors, []string{chain.DeployerKey.PublicKey().String()})

	t.Run("without MCMS", func(t *testing.T) {
		err = rt.Exec(
			runtime.ChangesetTask(&GrantRoleTimelockSolana{}, GrantRoleTimelockSolanaConfig{
				Role:     timelockbindings.Executor_Role,
				Accounts: map[uint64][]solana.PublicKey{selector: executors1},
			}),
		)

		onChainExecutors, err = inspector.GetExecutors(t.Context(), timelockAddress(mcmsState))
		require.NoError(t, err)
		require.ElementsMatch(t, onChainExecutors, []string{
			chain.DeployerKey.PublicKey().String(), executors1[0].String(), executors1[1].String(),
		})
	})

	t.Run("with MCMS", func(t *testing.T) {
		err = rt.Exec(
			runtime.ChangesetTask(&TransferMCMSToTimelockSolana{}, TransferMCMSToTimelockSolanaConfig{
				Chains:  []uint64{selector},
				MCMSCfg: proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
			}),
			runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
		)
		require.NoError(t, err)

		err = rt.Exec(
			runtime.ChangesetTask(&GrantRoleTimelockSolana{}, GrantRoleTimelockSolanaConfig{
				Role:     timelockbindings.Executor_Role,
				Accounts: map[uint64][]solana.PublicKey{selector: executors2},
				MCMS: &proposalutils.TimelockConfig{
					MinDelay:   1 * time.Second,
					MCMSAction: mcmstypes.TimelockActionSchedule,
				},
			}),
			runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
		)
		require.NoError(t, err)

		// --- assert ---
		onChainExecutors, err = inspector.GetExecutors(t.Context(), timelockAddress(mcmsState))
		require.NoError(t, err)
		require.ElementsMatch(t, onChainExecutors, []string{
			chain.DeployerKey.PublicKey().String(),
			executors1[0].String(), executors1[1].String(),
			executors2[0].String(), executors2[1].String(),
		})
	})
}

func randomSolanaAccounts(t *testing.T, n int) []solana.PublicKey {
	t.Helper()
	accounts := make([]solana.PublicKey, n)
	for i := range n {
		privateKey, err := solana.NewRandomPrivateKey()
		require.NoError(t, err)
		accounts[i] = privateKey.PublicKey()
	}

	return accounts
}

func timelockAddress(chainState *state.MCMSWithTimelockStateSolana) string {
	return state.EncodeAddressWithSeed(chainState.TimelockProgram, chainState.TimelockSeed)
}
