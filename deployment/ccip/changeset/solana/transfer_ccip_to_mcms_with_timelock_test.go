package solana_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	solBinary "github.com/gagliardetto/binary"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	burnmint "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/burnmint_token_pool"
	lockrelease "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/lockrelease_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/rmn_remote"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	solanachangesets "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TODO: remove. These should be deployed as part of the test once deployment changesets are ready.
const TimelockProgramID = "LoCoNsJFuhTkSQjfdDfn3yuwqhSYoPujmviRHVCzsqn"
const MCMProgramID = "6UmMZr5MEqiKWD5jqTJd1WCR5kT8oZuFYBLJFi1o6GQX"

func TestValidateContracts(t *testing.T) {
	validPubkey := solana.NewWallet().PublicKey()

	zeroPubkey := solana.PublicKey{} // Zero public key

	makeState := func(router, feeQuoter solana.PublicKey) solanastateview.CCIPChainState {
		return solanastateview.CCIPChainState{
			Router:    router,
			FeeQuoter: feeQuoter,
		}
	}

	tests := []struct {
		name          string
		state         solanastateview.CCIPChainState
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
		SolChains: 1,
	})
	envWithInvalidSolChain := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{})
	envWithInvalidSolChain.BlockChains = cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		chainselectors.ETHEREUM_TESTNET_SEPOLIA_LENS_1.Selector: cldf_solana.Chain{},
	})
	timelockID := mcmsSolana.ContractAddress(solana.MustPublicKeyFromBase58(TimelockProgramID), [32]byte{'t', 'e', 's', 't'})
	mcmsID := mcmsSolana.ContractAddress(solana.MustPublicKeyFromBase58(MCMProgramID), [32]byte{'t', 'e', 's', 't'})
	solChain := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))[0]
	err := env.ExistingAddresses.Save(solChain, timelockID, cldf.TypeAndVersion{Type: commontypes.RBACTimelock, Version: deployment.Version1_0_0})
	require.NoError(t, err)
	err = env.ExistingAddresses.Save(solChain, mcmsID, cldf.TypeAndVersion{Type: commontypes.ProposerManyChainMultisig, Version: deployment.Version1_0_0})
	require.NoError(t, err)

	tests := []struct {
		name             string
		env              cldf.Environment
		contractsByChain map[uint64]solanachangesets.CCIPContractsToTransfer
		expectedError    string
	}{
		{
			name:          "No chains found in environment",
			env:           memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{}),
			expectedError: "no chains found",
		},
		{
			name: "Chain selector not found in environment",
			env: memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
				SolChains: 1,
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
				MCMSCfg: proposalutils.TimelockConfig{
					MinDelay: 0 * time.Second,
				},
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
func prepareEnvironmentForOwnershipTransfer(t *testing.T) (cldf.Environment, stateview.CCIPOnChainState) {
	t.Helper()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		SolChains:  1,
		Nodes:      4,
	})
	evmSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))
	solChain1 := solChainSelectors[0]
	solChain := e.BlockChains.SolanaChains()[solChain1]
	selectors := make([]uint64, 0, len(evmSelectors)+len(solChainSelectors))
	selectors = append(selectors, evmSelectors...)
	selectors = append(selectors, solChainSelectors...)
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	// Fund account for fees
	testutils.FundAccounts(e.GetContext(), []solana.PrivateKey{*solChain.DeployerKey}, solChain.Client, t)
	err = testhelpers.SavePreloadedSolAddresses(e, solChainSelectors[0])
	require.NoError(t, err)
	solLinkTokenPrivKey, _ := solana.NewRandomPrivateKey()
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.BlockChains.EVMChains()[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: nodes.NonBootstraps().PeerIDs(),
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeploySolanaLinkToken),
			commonchangeset.DeploySolanaLinkTokenConfig{
				ChainSelector: solChain1,
				TokenPrivKey:  solLinkTokenPrivKey,
				TokenDecimals: 9,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(solanachangesets.DeployChainContractsChangeset),
			solanachangesets.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ChainSelector:     solChain1,
				ContractParamsPerChain: solanachangesets.ChainContractParams{
					FeeQuoterParams: solanachangesets.FeeQuoterParams{
						DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil},
					},
					OffRampParams: solanachangesets.OffRampParams{
						EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
					},
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(solanachangesets.DeploySolanaToken),
			solanachangesets.DeploySolanaTokenConfig{
				ChainSelector:    solChain1,
				TokenProgramName: shared.SPL2022Tokens,
				TokenDecimals:    9,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(solanachangesets.DeploySolanaToken),
			solanachangesets.DeploySolanaTokenConfig{
				ChainSelector:    solChain1,
				TokenProgramName: shared.SPLTokens,
				TokenDecimals:    9,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
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
	err = testhelpers.ValidateSolanaState(e, solChainSelectors)
	require.NoError(t, err)
	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	tokenAddressLockRelease := state.SolChains[solChain1].SPL2022Tokens[0]
	tokenAddressBurnMint := state.SolChains[solChain1].SPLTokens[0]

	lnr := test_token_pool.LockAndRelease_PoolType
	bnm := test_token_pool.BurnAndMint_PoolType
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(solanachangesets.AddTokenPoolAndLookupTable),
			solanachangesets.TokenPoolConfig{
				ChainSelector: solChain1,
				TokenPubKey:   tokenAddressLockRelease,
				PoolType:      &lnr,
				Metadata:      shared.CLLMetadata,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(solanachangesets.AddTokenPoolAndLookupTable),
			solanachangesets.TokenPoolConfig{
				ChainSelector: solChain1,
				TokenPubKey:   tokenAddressBurnMint,
				PoolType:      &bnm,
				Metadata:      shared.CLLMetadata,
			},
		),
	})
	require.NoError(t, err)
	return e, state
}
func TestTransferCCIPToMCMSWithTimelockSolana(t *testing.T) {
	t.Parallel()
	e, state := prepareEnvironmentForOwnershipTransfer(t)
	solChain1 := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))[0]
	solChain := e.BlockChains.SolanaChains()[solChain1]

	tokenAddressLockRelease := state.SolChains[solChain1].SPL2022Tokens[0]
	tokenAddressBurnMint := state.SolChains[solChain1].SPLTokens[0]

	burnMintPoolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenAddressBurnMint, state.SolChains[solChain1].BurnMintTokenPools[shared.CLLMetadata])
	lockReleasePoolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenAddressLockRelease, state.SolChains[solChain1].LockReleaseTokenPools[shared.CLLMetadata])
	timelockSignerPDA, _ := testhelpers.TransferOwnershipSolana(
		t,
		&e,
		solChain1,
		false,
		solanachangesets.CCIPContractsToTransfer{
			Router:                true,
			FeeQuoter:             true,
			OffRamp:               true,
			RMNRemote:             true,
			BurnMintTokenPools:    map[string]map[solana.PublicKey]solana.PublicKey{shared.CLLMetadata: {burnMintPoolConfigPDA: tokenAddressBurnMint}},
			LockReleaseTokenPools: map[string]map[solana.PublicKey]solana.PublicKey{shared.CLLMetadata: {lockReleasePoolConfigPDA: tokenAddressLockRelease}},
		})

	// 5. Now verify on-chain that each contract’s “config account” authority is the Timelock PDA.
	//    Typically, each contract has its own config account: RouterConfigPDA, FeeQuoterConfigPDA,
	//    Token Pool config PDAs, OffRamp config, etc.
	ctx := context.Background()

	// (A) Check Router ownership -  we need to add retries as the ownership transfer commitment is confirmed and not finalized.
	require.Eventually(t, func() bool {
		routerConfigPDA := state.SolChains[solChain1].RouterConfigPDA
		t.Logf("Checking Router Config PDA ownership data configPDA: %s", routerConfigPDA.String())
		programData := ccip_router.Config{}
		err := solChain.GetAccountDataBorshInto(ctx, routerConfigPDA, &programData)
		require.NoError(t, err)
		return timelockSignerPDA.String() == programData.Owner.String()
	}, 30*time.Second, 5*time.Second, "Router config PDA owner was not changed to timelock signer PDA")

	// (B) Check FeeQuoter ownership
	require.Eventually(t, func() bool {
		feeQuoterConfigPDA := state.SolChains[solChain1].FeeQuoterConfigPDA
		t.Logf("Checking Fee Quoter PDA ownership data configPDA: %s", feeQuoterConfigPDA.String())
		programData := fee_quoter.Config{}
		err := solChain.GetAccountDataBorshInto(ctx, feeQuoterConfigPDA, &programData)
		require.NoError(t, err)
		return timelockSignerPDA.String() == programData.Owner.String()
	}, 30*time.Second, 5*time.Second, "Fee Quoter config PDA owner was not changed to timelock signer PDA")

	// (C) Check OffRamp:
	require.Eventually(t, func() bool {
		offRampConfigPDA := state.SolChains[solChain1].OffRampConfigPDA
		programData := ccip_offramp.Config{}
		t.Logf("Checking Off Ramp PDA ownership data configPDA: %s", offRampConfigPDA.String())
		err := solChain.GetAccountDataBorshInto(ctx, offRampConfigPDA, &programData)
		require.NoError(t, err)
		return timelockSignerPDA.String() == programData.Owner.String()
	}, 30*time.Second, 5*time.Second, "OffRamp config PDA owner was not changed to timelock signer PDA")

	// (D) Check BurnMintTokenPools ownership:
	require.Eventually(t, func() bool {
		programData := burnmint.State{}
		t.Logf("Checking BurnMintTokenPools ownership data. configPDA: %s", burnMintPoolConfigPDA.String())
		err := solChain.GetAccountDataBorshInto(ctx, burnMintPoolConfigPDA, &programData)
		require.NoError(t, err)
		return timelockSignerPDA.String() == programData.Config.Owner.String()
	}, 30*time.Second, 5*time.Second, "BurnMintTokenPool owner was not changed to timelock signer PDA")

	// (E) Check LockReleaseTokenPools ownership:
	require.Eventually(t, func() bool {
		programData := lockrelease.State{}
		t.Logf("Checking LockReleaseTokenPools ownership data. configPDA: %s", lockReleasePoolConfigPDA.String())
		err := solChain.GetAccountDataBorshInto(ctx, lockReleasePoolConfigPDA, &programData)
		require.NoError(t, err)
		return timelockSignerPDA.String() == programData.Config.Owner.String()
	}, 30*time.Second, 5*time.Second, "LockReleaseTokenPool owner was not changed to timelock signer PDA")

	// (F) Check RMNRemote ownership
	require.Eventually(t, func() bool {
		rmnRemoteConfigPDA := state.SolChains[solChain1].RMNRemoteConfigPDA
		t.Logf("Checking RMNRemote PDA ownership data configPDA: %s", rmnRemoteConfigPDA.String())
		programData := rmn_remote.Config{}
		err := solChain.GetAccountDataBorshInto(ctx, rmnRemoteConfigPDA, &programData)
		require.NoError(t, err)
		return timelockSignerPDA.String() == programData.Owner.String()
	}, 30*time.Second, 5*time.Second, "RMNRemote config PDA owner was not changed to timelock signer PDA")
}

func TestTransferCCIPFromMCMSWithTimelockSolana(t *testing.T) {
	t.Parallel()
	e, state := prepareEnvironmentForOwnershipTransfer(t)
	solChain1 := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))[0]
	solChain := e.BlockChains.SolanaChains()[solChain1]

	tokenAddressLockRelease := state.SolChains[solChain1].SPL2022Tokens[0]
	tokenAddressBurnMint := state.SolChains[solChain1].SPLTokens[0]

	burnMintPoolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenAddressBurnMint, state.SolChains[solChain1].BurnMintTokenPools[shared.CLLMetadata])
	lockReleasePoolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenAddressLockRelease, state.SolChains[solChain1].LockReleaseTokenPools[shared.CLLMetadata])
	timelockSignerPDA, _ := testhelpers.TransferOwnershipSolana(
		t,
		&e,
		solChain1,
		false,
		solanachangesets.CCIPContractsToTransfer{
			Router:                true,
			FeeQuoter:             true,
			OffRamp:               true,
			RMNRemote:             true,
			BurnMintTokenPools:    map[string]map[solana.PublicKey]solana.PublicKey{shared.CLLMetadata: {burnMintPoolConfigPDA: tokenAddressBurnMint}},
			LockReleaseTokenPools: map[string]map[solana.PublicKey]solana.PublicKey{shared.CLLMetadata: {lockReleasePoolConfigPDA: tokenAddressLockRelease}},
		})
	// Transfer ownership back to the deployer
	e, _, err := commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolana),
			ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolanaConfig{
				MCMSCfg:       proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
				CurrentOwner:  timelockSignerPDA,
				ProposedOwner: solChain.DeployerKey.PublicKey(),
				ContractsByChain: map[uint64]ccipChangesetSolana.CCIPContractsToTransfer{
					solChain1: ccipChangesetSolana.CCIPContractsToTransfer{
						Router:    true,
						FeeQuoter: true,
						OffRamp:   true,
						RMNRemote: true,
						BurnMintTokenPools: map[string]map[solana.PublicKey]solana.PublicKey{
							shared.CLLMetadata: {
								burnMintPoolConfigPDA: tokenAddressBurnMint,
							},
						},
						LockReleaseTokenPools: map[string]map[solana.PublicKey]solana.PublicKey{
							shared.CLLMetadata: {
								lockReleasePoolConfigPDA: tokenAddressLockRelease,
							},
						},
					},
				},
			},
		),
	})
	require.NoError(t, err)
	// we have to accept separate from the changeset because the proposal needs to execute
	// just spot check that the ownership transfer happened
	config := state.SolChains[solChain1].RouterConfigPDA
	ix, err := ccip_router.NewAcceptOwnershipInstruction(
		config, solChain.DeployerKey.PublicKey(),
	).ValidateAndBuild()
	require.NoError(t, err)
	err = solChain.Confirm([]solana.Instruction{ix})
	require.NoError(t, err)

	// lnr
	lnrIx, err := lockrelease.NewAcceptOwnershipInstruction(
		lockReleasePoolConfigPDA, tokenAddressLockRelease, solChain.DeployerKey.PublicKey(),
	).ValidateAndBuild()
	require.NoError(t, err)
	err = solChain.Confirm([]solana.Instruction{lnrIx})
	require.NoError(t, err)

	// bnm
	bnmIx, err := burnmint.NewAcceptOwnershipInstruction(
		burnMintPoolConfigPDA, tokenAddressBurnMint, solChain.DeployerKey.PublicKey(),
	).ValidateAndBuild()
	require.NoError(t, err)
	err = solChain.Confirm([]solana.Instruction{bnmIx})
	require.NoError(t, err)
}
