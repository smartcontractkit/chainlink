package solana_test

import (
	"fmt"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonSolana "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// setupFundingTestEnv deploys all required contracts for the funding test
func setupFundingTestEnv(t *testing.T) cldf.Environment {
	lggr := logger.TestLogger(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:     1,
		SolChains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	chainSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))[0]

	config := proposalutils.SingleGroupTimelockConfigV2(t)
	err := testhelpers.SavePreloadedSolAddresses(env, chainSelector)
	require.NoError(t, err)
	// Initialize the address book with a dummy address to avoid deploy precondition errors.
	err = env.ExistingAddresses.Save(chainSelector, "dummyAddress", cldf.TypeAndVersion{Type: "dummy", Version: deployment.Version1_0_0})
	require.NoError(t, err)

	// Deploy MCMS and Timelock
	env, err = changeset.Apply(t, env, nil,
		changeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.DeployMCMSWithTimelockV2),
			map[uint64]types.MCMSWithTimelockConfigV2{
				chainSelector: config,
			},
		),
	)
	require.NoError(t, err)

	return env
}

func TestFundMCMSignersChangeset_VerifyPreconditions(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	validEnv := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{SolChains: 1})
	validSolChainSelector := validEnv.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))[0]

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
	err := validEnv.ExistingAddresses.Save(validSolChainSelector, timelockID, cldf.TypeAndVersion{
		Type:    types.RBACTimelock,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)
	err = validEnv.ExistingAddresses.Save(validSolChainSelector, mcmsProposerID, cldf.TypeAndVersion{
		Type:    types.ProposerManyChainMultisig,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)
	err = validEnv.ExistingAddresses.Save(validSolChainSelector, mcmsCancellerID, cldf.TypeAndVersion{
		Type:    types.CancellerManyChainMultisig,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)
	err = validEnv.ExistingAddresses.Save(validSolChainSelector, mcmsBypasserID, cldf.TypeAndVersion{
		Type:    types.BypasserManyChainMultisig,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)
	mcmsProposerIDEmpty := mcmsSolana.ContractAddress(
		mcmDummyProgram,
		[32]byte{},
	)

	// Create an environment that simulates a chain where the MCMS contracts have not been deployed,
	// e.g. missing the required addresses so that the state loader returns empty seeds.
	noMCMSEnv := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{})
	noMCMSEnv.BlockChains = cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		chainselectors.SOLANA_DEVNET.Selector: cldf_solana.Chain{},
	})
	err = noMCMSEnv.ExistingAddresses.Save(chainselectors.SOLANA_DEVNET.Selector, mcmsProposerIDEmpty, cldf.TypeAndVersion{
		Type:    types.BypasserManyChainMultisig,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)

	// Create an environment with a Solana chain that has an invalid (zero) underlying chain.
	invalidSolChainEnv := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{})
	invalidSolChainEnv.BlockChains = cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		validSolChainSelector: cldf_solana.Chain{},
	})

	tests := []struct {
		name          string
		env           cldf.Environment
		config        commonSolana.FundMCMSignerConfig
		expectedError string
	}{
		{
			name: "All preconditions satisfied",
			env:  validEnv,
			config: commonSolana.FundMCMSignerConfig{
				AmountsPerChain: map[uint64]commonSolana.AmountsToTransfer{validSolChainSelector: {
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
			env: memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
				Bootstraps: 1,
				Chains:     1,
				SolChains:  0,
				Nodes:      1,
			}),
			config: commonSolana.FundMCMSignerConfig{
				AmountsPerChain: map[uint64]commonSolana.AmountsToTransfer{validSolChainSelector: {
					ProposeMCM:   100,
					CancellerMCM: 100,
					BypasserMCM:  100,
					Timelock:     100,
				}},
			},
			expectedError: fmt.Sprintf("solana chain not found for selector %d", validSolChainSelector),
		},
		{
			name: "Chain selector not found in environment",
			env:  validEnv,
			config: commonSolana.FundMCMSignerConfig{AmountsPerChain: map[uint64]commonSolana.AmountsToTransfer{99999: {
				ProposeMCM:   100,
				CancellerMCM: 100,
				BypasserMCM:  100,
				Timelock:     100,
			}}},
			expectedError: "solana chain not found for selector 99999",
		},
		{
			name: "MCMS contracts not deployed (empty seeds)",
			env:  noMCMSEnv,
			config: commonSolana.FundMCMSignerConfig{
				AmountsPerChain: map[uint64]commonSolana.AmountsToTransfer{chainselectors.SOLANA_DEVNET.Selector: {
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
			env:  validEnv,
			config: commonSolana.FundMCMSignerConfig{
				AmountsPerChain: map[uint64]commonSolana.AmountsToTransfer{validSolChainSelector: {
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
			env:  invalidSolChainEnv,
			config: commonSolana.FundMCMSignerConfig{
				AmountsPerChain: map[uint64]commonSolana.AmountsToTransfer{validSolChainSelector: {
					ProposeMCM:   100,
					CancellerMCM: 100,
					BypasserMCM:  100,
					Timelock:     100,
				}},
			},
			expectedError: "failed to get existing addresses: chain selector 12463857294658392847: chain not found",
		},
	}

	cs := commonSolana.FundMCMSignersChangeset{}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			err := cs.VerifyPreconditions(tt.env, tt.config)
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestFundMCMSignersChangeset_Apply(t *testing.T) {
	t.Parallel()
	env := setupFundingTestEnv(t)
	cfgAmounts := commonSolana.AmountsToTransfer{
		ProposeMCM:   100 * solana.LAMPORTS_PER_SOL,
		CancellerMCM: 350 * solana.LAMPORTS_PER_SOL,
		BypasserMCM:  75 * solana.LAMPORTS_PER_SOL,
		Timelock:     83 * solana.LAMPORTS_PER_SOL,
	}
	solChains := env.BlockChains.SolanaChains()
	amountsPerChain := make(map[uint64]commonSolana.AmountsToTransfer)
	for chainSelector := range solChains {
		amountsPerChain[chainSelector] = cfgAmounts
	}
	config := commonSolana.FundMCMSignerConfig{
		AmountsPerChain: amountsPerChain,
	}

	changesetInstance := commonSolana.FundMCMSignersChangeset{}

	env, _, err := changeset.ApplyChangesetsV2(t, env, []changeset.ConfiguredChangeSet{
		changeset.Configure(changesetInstance, config),
	})
	require.NoError(t, err)

	chainSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))[0]
	solChain := solChains[chainSelector]
	addresses, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)

	// Check balances of MCM Signer PDAS
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addresses)
	require.NoError(t, err)

	accounts := []solana.PublicKey{
		state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed),
		state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.ProposerMcmSeed),
		state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.CancellerMcmSeed),
		state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.BypasserMcmSeed),
	}
	var balances []uint64
	for _, account := range accounts {
		balance, err := solChain.Client.GetBalance(env.GetContext(), account, rpc.CommitmentConfirmed)
		require.NoError(t, err)
		t.Logf("Account: %s, Balance: %d", account, balance.Value)
		balances = append(balances, balance.Value)
	}

	require.Equal(t, cfgAmounts.Timelock, balances[0])
	require.Equal(t, cfgAmounts.ProposeMCM, balances[1])
	require.Equal(t, cfgAmounts.CancellerMCM, balances[2])
	require.Equal(t, cfgAmounts.BypasserMCM, balances[3])
}
