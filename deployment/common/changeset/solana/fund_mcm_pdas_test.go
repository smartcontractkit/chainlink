package solana

import (
	"fmt"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/smartcontractkit/quarantine"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestFundMCMSignersChangeset_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	selector1 := chainselectors.TEST_22222222222222222222222222222222222222222222.Selector
	selector2 := chainselectors.TEST_33333333333333333333333333333333333333333333.Selector

	env, err := environment.New(t.Context(),
		environment.WithSolanaContainer(t, []uint64{selector1, selector2}, t.TempDir(), map[string]string{}),
	)
	require.NoError(t, err)

	// Setup selector1 to have a chain where programs are deployed to pass validation
	timelockID := mcmsSolana.ContractAddress(
		solana.NewWallet().PublicKey(),
		[32]byte{'t', 'e', 's', 't'},
	)
	mcmDummyProgram := solana.NewWallet().PublicKey()
	mcmsProposerID := mcmsSolana.ContractAddress(
		mcmDummyProgram,
		[32]byte{'t', 'e', 's', 't', '1'},
	)

	mcmsCancellerID := mcmsSolana.ContractAddress(
		mcmDummyProgram,
		[32]byte{'t', 'e', 's', 't', '2'},
	)

	mcmsBypasserID := mcmsSolana.ContractAddress(
		mcmDummyProgram,
		[32]byte{'t', 'e', 's', 't', '3'},
	)
	err = env.ExistingAddresses.Save(selector1, timelockID, cldf.TypeAndVersion{
		Type:    types.RBACTimelock,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)
	err = env.ExistingAddresses.Save(selector1, mcmsProposerID, cldf.TypeAndVersion{
		Type:    types.ProposerManyChainMultisig,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)
	err = env.ExistingAddresses.Save(selector1, mcmsCancellerID, cldf.TypeAndVersion{
		Type:    types.CancellerManyChainMultisig,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)
	err = env.ExistingAddresses.Save(selector1, mcmsBypasserID, cldf.TypeAndVersion{
		Type:    types.BypasserManyChainMultisig,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)

	// Setup selector2 to have a chain where the MCMS contracts have not been deployed,
	// e.g. missing the required addresses so that the state loader returns empty seeds.
	mcmsProposerIDEmpty := mcmsSolana.ContractAddress(
		mcmDummyProgram,
		[32]byte{},
	)

	err = env.ExistingAddresses.Save(selector2, mcmsProposerIDEmpty, cldf.TypeAndVersion{
		Type:    types.BypasserManyChainMultisig,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)

	tests := []struct {
		name          string
		env           func(t *testing.T) cldf.Environment
		config        FundMCMSignerConfig
		expectedError string
	}{
		{
			name: "All preconditions satisfied",
			env:  func(t *testing.T) cldf.Environment { t.Helper(); return *env },
			config: FundMCMSignerConfig{
				AmountsPerChain: map[uint64]AmountsToTransfer{selector1: {
					ProposeMCM:   100,
					CancellerMCM: 100,
					BypasserMCM:  100,
					Timelock:     100,
				}},
			},
			expectedError: "",
		},
		{
			name: "No Solana chains found in environment",
			env: func(t *testing.T) cldf.Environment {
				t.Helper()

				// Create an environment with no solana chains
				emptyEnv, err := environment.New(t.Context())
				require.NoError(t, err)

				return *emptyEnv
			},
			config: FundMCMSignerConfig{
				AmountsPerChain: map[uint64]AmountsToTransfer{selector1: {
					ProposeMCM:   100,
					CancellerMCM: 100,
					BypasserMCM:  100,
					Timelock:     100,
				}},
			},
			expectedError: fmt.Sprintf("solana chain not found for selector %d", selector1),
		},
		{
			name: "Chain selector not found in environment",
			env:  func(t *testing.T) cldf.Environment { t.Helper(); return *env },
			config: FundMCMSignerConfig{AmountsPerChain: map[uint64]AmountsToTransfer{99999: {
				ProposeMCM:   100,
				CancellerMCM: 100,
				BypasserMCM:  100,
				Timelock:     100,
			}}},
			expectedError: "solana chain not found for selector 99999",
		},
		{
			name: "MCMS contracts not deployed (empty seeds)",
			env:  func(t *testing.T) cldf.Environment { t.Helper(); return *env },
			config: FundMCMSignerConfig{
				AmountsPerChain: map[uint64]AmountsToTransfer{selector2: {
					ProposeMCM:   100,
					CancellerMCM: 100,
					BypasserMCM:  100,
					Timelock:     100,
				}},
			},
			expectedError: "mcm/timelock seeds are empty, please deploy MCMS contracts first",
		},
		{
			name: "Insufficient deployer balance",
			env:  func(t *testing.T) cldf.Environment { t.Helper(); return *env },
			config: FundMCMSignerConfig{
				AmountsPerChain: map[uint64]AmountsToTransfer{selector1: {
					ProposeMCM:   9999999999999999999,
					CancellerMCM: 9999999999999999999,
					BypasserMCM:  9999999999999999999,
					Timelock:     9999999999999999999,
				}},
			},
			expectedError: "deployer balance is insufficient",
		},
		{
			name: "Invalid Solana chain in environment",
			env: func(t *testing.T) cldf.Environment {
				t.Helper()

				invalidEnv, err := environment.New(t.Context())
				require.NoError(t, err)
				invalidEnv.BlockChains = cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
					selector1: cldf_solana.Chain{}, // Empty chain is invalid
				})

				return *invalidEnv
			},
			config: FundMCMSignerConfig{
				AmountsPerChain: map[uint64]AmountsToTransfer{selector1: {
					ProposeMCM:   100,
					CancellerMCM: 100,
					BypasserMCM:  100,
					Timelock:     100,
				}},
			},
			expectedError: "failed to get existing addresses: chain selector 12463857294658392847: chain not found",
		},
	}

	cs := FundMCMSignersChangeset{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cs.VerifyPreconditions(tt.env(t), tt.config)
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectedError)
			}
		})
	}
}

func TestFundMCMSignersChangeset_Apply(t *testing.T) {
	quarantine.Flaky(t, "DX-1776")
	t.Parallel()

	rt, selector := setupTest(t)

	chain := rt.Environment().BlockChains.SolanaChains()[selector]
	cfgAmounts := AmountsToTransfer{
		ProposeMCM:   100 * solana.LAMPORTS_PER_SOL,
		CancellerMCM: 350 * solana.LAMPORTS_PER_SOL,
		BypasserMCM:  75 * solana.LAMPORTS_PER_SOL,
		Timelock:     83 * solana.LAMPORTS_PER_SOL,
	}

	err := rt.Exec(
		runtime.ChangesetTask(FundMCMSignersChangeset{}, FundMCMSignerConfig{
			AmountsPerChain: map[uint64]AmountsToTransfer{selector: cfgAmounts},
		}),
	)
	require.NoError(t, err)

	addresses, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)

	// Check balances of MCM Signer PDAS
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	require.NoError(t, err)

	accounts := []solana.PublicKey{
		state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed),
		state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.ProposerMcmSeed),
		state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.CancellerMcmSeed),
		state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.BypasserMcmSeed),
	}
	var balances []uint64
	for _, account := range accounts {
		balance, err := chain.Client.GetBalance(t.Context(), account, rpc.CommitmentConfirmed)
		require.NoError(t, err)
		t.Logf("Account: %s, Balance: %d", account, balance.Value)
		balances = append(balances, balance.Value)
	}

	require.Equal(t, cfgAmounts.Timelock, balances[0])
	require.Equal(t, cfgAmounts.ProposeMCM, balances[1])
	require.Equal(t, cfgAmounts.CancellerMCM, balances[2])
	require.Equal(t, cfgAmounts.BypasserMCM, balances[3])
}
