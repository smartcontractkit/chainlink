package example_test

import (
	"crypto/ecdsa"
	"fmt"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/smartcontractkit/quarantine"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/example"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
	"github.com/smartcontractkit/chainlink/deployment/utils/solutils"
)

func TestTransferFromTimelockConfig_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	receiverKey := solana.NewWallet().PublicKey()
	selector := chainselectors.TEST_22222222222222222222222222222222222222222222.Selector

	// Save the timelock contract address to the address book
	ab := cldf.NewMemoryAddressBook()
	timelockID := mcmsSolana.ContractAddress(
		solana.NewWallet().PublicKey(),
		[32]byte{'t', 'e', 's', 't'},
	)
	require.NoError(t, ab.Save(selector, timelockID, cldf.TypeAndVersion{
		Type:    types.RBACTimelock,
		Version: deployment.Version1_0_0,
	}))

	env, err := environment.New(t.Context(),
		environment.WithSolanaContainer(t, []uint64{selector}, t.TempDir(), map[string]string{}),
		environment.WithAddressBook(ab),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	// Create an environment with a Solana chain that has an invalid (zero value) underlying chain.
	invalidEnv, err := environment.New(t.Context())
	require.NoError(t, err)
	invalidEnv.BlockChains = cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		selector: cldf_solana.Chain{},
	})

	tests := []struct {
		name          string
		env           func(t *testing.T) cldf.Environment
		config        example.TransferFromTimelockConfig
		expectedError string
	}{
		{
			name: "All preconditions satisfied",
			env:  func(t *testing.T) cldf.Environment { t.Helper(); return *env },
			config: example.TransferFromTimelockConfig{
				AmountsPerChain: map[uint64]example.TransferData{selector: {
					Amount: 100,
					To:     receiverKey,
				}},
			},
			expectedError: "",
		},
		{
			name: "No Solana chains found in environment",
			env: func(t *testing.T) cldf.Environment {
				t.Helper()

				emptyEnv, gerr := environment.New(t.Context())
				require.NoError(t, gerr)

				return *emptyEnv
			},
			config: example.TransferFromTimelockConfig{
				AmountsPerChain: map[uint64]example.TransferData{selector: {
					Amount: 100,
					To:     receiverKey,
				}},
			},
			expectedError: fmt.Sprintf("solana chain not found for selector %d", selector),
		},
		{
			name: "Chain selector not found in environment",
			env:  func(t *testing.T) cldf.Environment { t.Helper(); return *env },
			config: example.TransferFromTimelockConfig{
				AmountsPerChain: map[uint64]example.TransferData{99999: {
					Amount: 100,
					To:     receiverKey,
				}}},
			expectedError: "solana chain not found for selector 99999",
		},
		{
			name: "timelock contracts not deployed (empty seeds)",
			env: func(t *testing.T) cldf.Environment {
				t.Helper()

				// Create an environment that simulates a chain where the MCMS contracts have not been deployed,
				// e.g. missing the required addresses so that the state loader returns empty seeds.
				emptyEnv, gerr := environment.New(t.Context())

				require.NoError(t, gerr)

				emptyEnv.BlockChains = cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{
					cldf_solana.Chain{Selector: selector},
				})
				gerr = emptyEnv.ExistingAddresses.Save(selector, "dummy", cldf.TypeAndVersion{
					Type:    "Sometype",
					Version: deployment.Version1_0_0,
				})
				require.NoError(t, gerr)

				return *emptyEnv
			},
			config: example.TransferFromTimelockConfig{
				AmountsPerChain: map[uint64]example.TransferData{selector: {
					Amount: 100,
					To:     receiverKey,
				}},
			},
			expectedError: "timelock seeds are empty, please deploy MCMS contracts first",
		},
		{
			name: "Insufficient deployer balance",
			env:  func(t *testing.T) cldf.Environment { t.Helper(); return *env },
			config: example.TransferFromTimelockConfig{
				AmountsPerChain: map[uint64]example.TransferData{
					selector: {
						Amount: 999999999999999999,
						To:     receiverKey,
					},
				},
			},
			expectedError: "deployer balance is insufficient",
		},
		{
			name: "Invalid Solana chain in environment",
			env:  func(t *testing.T) cldf.Environment { t.Helper(); return *invalidEnv },
			config: example.TransferFromTimelockConfig{
				AmountsPerChain: map[uint64]example.TransferData{selector: {
					Amount: 100,
					To:     receiverKey,
				}},
			},
			expectedError: "failed to get existing addresses: chain selector 12463857294658392847: chain not found",
		},
		{
			name: "empty from field",
			env:  func(t *testing.T) cldf.Environment { t.Helper(); return *invalidEnv },
			config: example.TransferFromTimelockConfig{
				AmountsPerChain: map[uint64]example.TransferData{selector: {
					Amount: 100,
					To:     solana.PublicKey{},
				}},
			},
			expectedError: "destination address is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := example.TransferFromTimelock{}.VerifyPreconditions(tt.env(t), tt.config)
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestTransferFromTimelockConfig_Apply(t *testing.T) {
	quarantine.Flaky(t, "DX-1754")
	t.Parallel()

	selector := chainselectors.TEST_22222222222222222222222222222222222222222222.Selector
	programsPath, programIDs, ab := soltestutils.PreloadMCMS(t, selector)

	// Initialize the address book with a dummy address to avoid deploy precondition errors.
	err := ab.Save(selector, "dummyAddress", cldf.TypeAndVersion{Type: "dummy", Version: deployment.Version1_0_0})
	require.NoError(t, err)

	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithSolanaContainer(t, []uint64{selector}, programsPath, programIDs),
		environment.WithAddressBook(ab),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.SolanaChains()[selector]

	// Deploy MCMS and Timelock
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployMCMSWithTimelockV2), map[uint64]types.MCMSWithTimelockConfigV2{
			selector: proposalutils.SingleGroupTimelockConfigV2(t),
		}),
	)
	require.NoError(t, err)

	// Fund the signer PDAs for the MCMS contracts
	mcmState := soltestutils.GetMCMSStateFromAddressBook(t, rt.State().AddressBook, chain)
	timelockSigner := state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	mcmSigner := state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.ProposerMcmSeed)

	err = solutils.FundAccounts(t.Context(), chain.Client, []solana.PublicKey{timelockSigner, mcmSigner, chain.DeployerKey.PublicKey()}, 150)
	require.NoError(t, err)

	// Execute the transfer from timelock changeset
	cfgAmounts := example.TransferData{
		Amount: 100 * solana.LAMPORTS_PER_SOL,
		To:     solana.NewWallet().PublicKey(),
	}

	err = rt.Exec(
		runtime.ChangesetTask(example.TransferFromTimelock{}, example.TransferFromTimelockConfig{
			TimelockCfg: proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
			AmountsPerChain: map[uint64]example.TransferData{
				selector: cfgAmounts,
			},
		}),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
	)
	require.NoError(t, err)

	balance, err := chain.Client.GetBalance(t.Context(), cfgAmounts.To, rpc.CommitmentConfirmed)
	require.NoError(t, err)
	t.Logf("Account: %s, Balance: %d", cfgAmounts.To, balance.Value)

	require.Equal(t, cfgAmounts.Amount, balance.Value)
}
