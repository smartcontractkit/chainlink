package solana_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	solanachangesets "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	commonState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

// TODO: remove. These should be deployed as part of the test once deployment changesets are ready.
const TimelockProgramID = "LoCoNsJFuhTkSQjfdDfn3yuwqhSYoPujmviRHVCzsqn"
const MCMProgramID = "6UmMZr5MEqiKWD5jqTJd1WCR5kT8oZuFYBLJFi1o6GQX"

func TestValidateContracts(t *testing.T) {
	validPubkey := solana.NewWallet().PublicKey()

	zeroPubkey := solana.PublicKey{} // Zero public key

	makeState := func(router, feeQuoter solana.PublicKey) changeset.SolCCIPChainState {
		return changeset.SolCCIPChainState{
			Router:    router,
			FeeQuoter: feeQuoter,
		}
	}

	tests := []struct {
		name          string
		state         changeset.SolCCIPChainState
		contracts     solanachangesets.CCIPContractsToTransfer
		chainSelector uint64
		expectedError string
	}{
		{
			name:          "All required contracts present",
			state:         makeState(validPubkey, validPubkey),
			contracts:     solanachangesets.CCIPContractsToTransfer{Router: true},
			chainSelector: 12345,
		},
		{
			name:          "Missing Router contract",
			state:         makeState(zeroPubkey, validPubkey),
			contracts:     solanachangesets.CCIPContractsToTransfer{Router: true},
			chainSelector: 12345,
			expectedError: "missing required contract Router on chain 12345",
		},
		{
			name:          "Missing FeeQuoter contract",
			state:         makeState(validPubkey, zeroPubkey),
			contracts:     solanachangesets.CCIPContractsToTransfer{Router: true, FeeQuoter: true},
			chainSelector: 12345,
			expectedError: "missing required contract FeeQuoter on chain 12345",
		},
		{
			name:          "invalid pub key",
			state:         makeState(validPubkey, zeroPubkey),
			contracts:     solanachangesets.CCIPContractsToTransfer{Router: true, FeeQuoter: true},
			chainSelector: 12345,
			expectedError: "missing required contract FeeQuoter on chain 12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := solanachangesets.ValidateContracts(tt.state, tt.chainSelector, tt.contracts)

			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Equal(t, tt.expectedError, err.Error())
			}
		})
	}
}

func TestValidate(t *testing.T) {
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		SolChains:  1,
		Nodes:      4,
	})
	envWithInvalidSolChain := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		SolChains:  1,
		Nodes:      4,
	})
	envWithInvalidSolChain.SolChains[chainselectors.ETHEREUM_TESTNET_SEPOLIA_LENS_1.Selector] = deployment.SolChain{}
	timelockID := mcmsSolana.ContractAddress(solana.MustPublicKeyFromBase58(TimelockProgramID), [32]byte{'t', 'e', 's', 't'})
	mcmsID := mcmsSolana.ContractAddress(solana.MustPublicKeyFromBase58(MCMProgramID), [32]byte{'t', 'e', 's', 't'})
	err := env.ExistingAddresses.Save(env.AllChainSelectorsSolana()[0], timelockID, deployment.TypeAndVersion{Type: commontypes.RBACTimelock, Version: deployment.Version1_0_0})
	require.NoError(t, err)
	err = env.ExistingAddresses.Save(env.AllChainSelectorsSolana()[0], mcmsID, deployment.TypeAndVersion{Type: commontypes.ProposerManyChainMultisig, Version: deployment.Version1_0_0})
	require.NoError(t, err)

	tests := []struct {
		name             string
		env              deployment.Environment
		contractsByChain map[uint64]solanachangesets.CCIPContractsToTransfer
		expectedError    string
	}{
		{
			name: "No chains found in environment",
			env: memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
				Bootstraps: 1,
				Chains:     0,
				SolChains:  0,
				Nodes:      4,
			}),
			expectedError: "no chains found",
		},
		{
			name: "Chain selector not found in environment",
			env: memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
				Bootstraps: 1,
				Chains:     1,
				SolChains:  1,
				Nodes:      4,
			}),
			contractsByChain: map[uint64]solanachangesets.CCIPContractsToTransfer{
				99999: {Router: true, FeeQuoter: true},
			},
			expectedError: "chain 99999 not found in environment",
		},
		{
			name: "Invalid chain family",
			env:  envWithInvalidSolChain,
			contractsByChain: map[uint64]solanachangesets.CCIPContractsToTransfer{
				chainselectors.ETHEREUM_TESTNET_SEPOLIA_LENS_1.Selector: {Router: true, FeeQuoter: true},
			},
			expectedError: "failed to load addresses for chain 6827576821754315911: chain selector 6827576821754315911: chain not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := solanachangesets.TransferCCIPToMCMSWithTimelockSolanaConfig{
				ContractsByChain: tt.contractsByChain,
				MinDelay:         10 * time.Second,
			}

			err := cfg.Validate(tt.env)

			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

// prepareEnvironmentForOwnershipTransfer helper that deploys the necessary contracts as pre-requisite to
// the transfer ownership changeset.
func prepareEnvironmentForOwnershipTransfer(t *testing.T) (deployment.Environment, changeset.CCIPOnChainState) {
	t.Helper()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		SolChains:  1,
		Nodes:      4,
	})
	evmSelectors := e.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.AllChainSelectorsSolana()
	solChain1 := e.AllChainSelectorsSolana()[0]
	solChain := e.SolChains[solChain1]
	selectors := make([]uint64, 0, len(evmSelectors)+len(solChainSelectors))
	selectors = append(selectors, evmSelectors...)
	selectors = append(selectors, solChainSelectors...)
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	// Fund account for fees
	testutils.FundAccounts(e.GetContext(), []solana.PrivateKey{*solChain.DeployerKey}, solChain.Client, t)
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	contractParams := make(map[uint64]changeset.ChainContractParams)
	for _, chain := range solChainSelectors {
		contractParams[chain] = changeset.ChainContractParams{
			FeeQuoterParams: changeset.DefaultFeeQuoterParams(),
			OffRampParams:   changeset.DefaultOffRampParams(),
		}
	}
	prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	for _, chain := range e.AllChainSelectors() {
		prereqCfg = append(prereqCfg, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
	}
	testhelpers.SavePreloadedSolAddresses(t, e, solChainSelectors[0])
	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.DeployHomeChainChangeset),
			changeset.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: nodes.NonBootstraps().PeerIDs(),
				},
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			selectors,
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelock),
			cfg,
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
			changeset.DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(solanachangesets.DeployChainContractsChangesetSolana),
			changeset.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParams,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(solanachangesets.DeploySolanaToken),
			solanachangesets.DeploySolanaTokenConfig{
				ChainSelector:    solChain1,
				TokenProgramName: deployment.SPL2022Tokens,
				TokenDecimals:    9,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			map[uint64]commontypes.MCMSWithTimelockConfigV2{
				solChain1: {
					Canceller:        proposalutils.SingleGroupMCMSV2(t),
					Proposer:         proposalutils.SingleGroupMCMSV2(t),
					Bypasser:         proposalutils.SingleGroupMCMSV2(t),
					TimelockMinDelay: big.NewInt(0),
				},
			},
		),
	})
	require.NoError(t, err)

	// solana verification
	testhelpers.ValidateSolanaState(t, e, solChainSelectors)
	state, err := changeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	tokenAddress := state.SolChains[solChain1].SPL2022Tokens[0]

	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(solanachangesets.AddTokenPool),
			solanachangesets.TokenPoolConfig{
				ChainSelector:    solChain1,
				TokenPubKey:      tokenAddress.String(),
				TokenProgramName: deployment.SPL2022Tokens,
				PoolType:         test_token_pool.LockAndRelease_PoolType,
				Authority:        e.SolChains[solChain1].DeployerKey.PublicKey().String(),
			},
		),
	})
	require.NoError(t, err)
	return e, state
}
func TestTransferCCIPToMCMSWithTimelockSolana(t *testing.T) {
	t.Parallel()
	e, state := prepareEnvironmentForOwnershipTransfer(t)
	solChain1 := e.AllChainSelectorsSolana()[0]
	solChain := e.SolChains[solChain1]
	// tokenAddress := state.SolChains[solChain1].SPL2022Tokens[0]
	addresses, err := e.ExistingAddresses.AddressesForChain(solChain1)
	require.NoError(t, err)
	mcmState, err := commonState.MaybeLoadMCMSWithTimelockChainStateSolana(e.SolChains[solChain1], addresses)
	require.NoError(t, err)

	// Fund signer PDAs for timelock and mcm
	// If we don't fund, execute() calls will fail with "no funds" errors.
	timelockSignerPDA := commonState.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	mcmSignerPDA := commonState.GetMCMSignerPDA(mcmState.McmProgram, mcmState.ProposerMcmSeed)
	memory.FundSolanaAccounts(e.GetContext(), t, []solana.PublicKey{timelockSignerPDA, mcmSignerPDA},
		100, solChain.Client)
	t.Logf("funded timelock signer PDA: %s", timelockSignerPDA.String())
	t.Logf("funded mcm signer PDA: %s", mcmSignerPDA.String())
	// Apply transfer ownership changeset
	e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(solanachangesets.TransferCCIPToMCMSWithTimelockSolana),
			solanachangesets.TransferCCIPToMCMSWithTimelockSolanaConfig{
				MinDelay: 1 * time.Second,
				ContractsByChain: map[uint64]solanachangesets.CCIPContractsToTransfer{
					solChain1: {
						Router:    true,
						FeeQuoter: true,
						OffRamp:   true,
					},
				},
			},
		),
	})
	require.NoError(t, err)

	// 5. Now verify on-chain that each contract’s “config account” authority is the Timelock PDA.
	//    Typically, each contract has its own config account: RouterConfigPDA, FeeQuoterConfigPDA,
	//    Token Pool config PDAs, OffRamp config, etc.
	ctx := context.Background()

	// (A) Check Router ownership -  we need to add retries as the ownership transfer commitment is confirmed and not finalized.
	require.Eventually(t, func() bool {
		routerConfigPDA := state.SolChains[solChain1].RouterConfigPDA
		t.Logf("Checking Router Config PDA ownership data configPDA: %s", routerConfigPDA.String())
		programData := ccip_router.Config{}
		err = solChain.GetAccountDataBorshInto(ctx, routerConfigPDA, &programData)
		return timelockSignerPDA.String() == programData.Owner.String()
	}, 30*time.Second, 5*time.Second, "Router config PDA owner was not changed to timelock signer PDA")

	// (B) Check FeeQuoter ownership
	require.Eventually(t, func() bool {
		feeQuoterConfigPDA := state.SolChains[solChain1].FeeQuoterConfigPDA
		t.Logf("Checking Fee Quoter PDA ownership data configPDA: %s", feeQuoterConfigPDA.String())
		programData := fee_quoter.Config{}
		err = solChain.GetAccountDataBorshInto(ctx, feeQuoterConfigPDA, &programData)
		require.NoError(t, err)
		return timelockSignerPDA.String() == programData.Owner.String()
	}, 30*time.Second, 5*time.Second, "Fee Quoter config PDA owner was not changed to timelock signer PDA")

	// (C) Check OffRamp:
	require.Eventually(t, func() bool {
		offRampConfigPDA := state.SolChains[solChain1].OffRampConfigPDA
		programData := ccip_offramp.Config{}
		t.Logf("Checking Off Ramp PDA ownership data configPDA: %s", offRampConfigPDA.String())
		err = solChain.GetAccountDataBorshInto(ctx, offRampConfigPDA, &programData)
		require.NoError(t, err)
		return timelockSignerPDA.String() == programData.Owner.String()
	}, 30*time.Second, 5*time.Second, "OffRamp config PDA owner was not changed to timelock signer PDA")
}
