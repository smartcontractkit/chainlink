package solana_test

import (
	"fmt"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

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
func setupFundingTestEnv(t *testing.T) deployment.Environment {
	lggr := logger.TestLogger(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:     1,
		SolChains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	chainSelector := env.AllChainSelectorsSolana()[0]

	config := proposalutils.SingleGroupTimelockConfigV2(t)
	testhelpers.SavePreloadedSolAddresses(t, env, chainSelector)
	// Initialize the address book with a dummy address to avoid deploy precondition errors.
	err := env.ExistingAddresses.Save(chainSelector, "dummyAddress", deployment.TypeAndVersion{Type: "dummy", Version: deployment.Version1_0_0})
	require.NoError(t, err)

	// Deploy MCMS and Timelock
	env, err = changeset.Apply(t, env, nil,
		changeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.DeployMCMSWithTimelockV2),
			map[uint64]types.MCMSWithTimelockConfigV2{
				chainSelector: config,
			},
		),
	)
	require.NoError(t, err)

	return env
}

func TestFundMCMSignersChangeset_VerifyPreconditions(t *testing.T) {
	// Create a logger instance.
	lggr := logger.TestLogger(t)
	// Create a valid in–memory environment with one Solana chain.
	validEnv := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{SolChains: 1})
	validEnv.SolChains[chainselectors.SOLANA_DEVNET.Selector] = deployment.SolChain{}
	// Get the sole Solana chain selector.
	validSolChainSelector := validEnv.AllChainSelectorsSolana()[0]

	// Save MCMS contract addresses into the environment to simulate that MCMS contracts
	// have been deployed. The addresses here are generated using dummy seeds.
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
	err := validEnv.ExistingAddresses.Save(validSolChainSelector, timelockID, deployment.TypeAndVersion{
		Type:    types.RBACTimelock,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)
	err = validEnv.ExistingAddresses.Save(validSolChainSelector, mcmsProposerID, deployment.TypeAndVersion{
		Type:    types.ProposerManyChainMultisig,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)
	err = validEnv.ExistingAddresses.Save(validSolChainSelector, mcmsCancellerID, deployment.TypeAndVersion{
		Type:    types.CancellerManyChainMultisig,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)
	err = validEnv.ExistingAddresses.Save(validSolChainSelector, mcmsBypasserID, deployment.TypeAndVersion{
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
	noMCMSEnv := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains:    0,
		SolChains: 1,
		Nodes:     1,
	})
	noMCMSEnv.SolChains[chainselectors.SOLANA_DEVNET.Selector] = deployment.SolChain{}
	err = noMCMSEnv.ExistingAddresses.Save(chainselectors.SOLANA_DEVNET.Selector, mcmsProposerIDEmpty, deployment.TypeAndVersion{
		Type:    types.BypasserManyChainMultisig,
		Version: deployment.Version1_0_0,
	})
	require.NoError(t, err)

	// Note: We do not call ExistingAddresses.Save on noMCMSEnv.

	// Create an environment with a Solana chain that has an invalid (zero) underlying chain.
	invalidSolChainEnv := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains:    0,
		SolChains: 0,
		Nodes:     1,
	})
	// Overwrite the solana chain with an empty one.
	invalidSolChainEnv.SolChains[validSolChainSelector] = deployment.SolChain{}

	tests := []struct {
		name          string
		env           deployment.Environment
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
			// Create an environment with zero Solana chains.
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
			// Use a chain selector that is not present in validEnv.
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
	// Set up the test environment
	env := setupFundingTestEnv(t)

	// Build a funding configuration.
	// Here, we assume that we want to fund each chain with an amount equal to 1000 SOL per MCMS signer.
	// There are 4 signers (bypasser, canceller, proposer, timelock).
	amountPerSigner := 589 * solana.LAMPORTS_PER_SOL
	amountsPerChain := make(map[uint64]commonSolana.AmountsToTransfer)
	for chainSelector := range env.SolChains {
		amountsPerChain[chainSelector] = commonSolana.AmountsToTransfer{
			ProposeMCM:   amountPerSigner,
			CancellerMCM: amountPerSigner,
			BypasserMCM:  amountPerSigner,
			Timelock:     amountPerSigner,
		}
	}
	config := commonSolana.FundMCMSignerConfig{
		AmountsPerChain: amountsPerChain,
	}

	// Create the changeset instance.
	changesetInstance := commonSolana.FundMCMSignersChangeset{}

	env, err := changeset.ApplyChangesetsV2(t, env, []changeset.ConfiguredChangeSet{
		changeset.Configure(changesetInstance, config),
	})
	require.NoError(t, err)

	chainSelector := env.AllChainSelectorsSolana()[0]
	solChain := env.SolChains[chainSelector]
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
	// Check if the accounts are funded
	for _, account := range accounts {
		balance, err := solChain.Client.GetBalance(env.GetContext(), account, rpc.CommitmentConfirmed)
		t.Logf("Account: %s, Balance: %d", account, balance.Value)
		require.NoError(t, err)
		require.GreaterOrEqual(t, amountPerSigner, balance.Value)
	}
}
