package solana_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	solBinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"
	burnmint "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/burnmint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/fee_quoter"
	lockrelease "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/lockrelease_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/rmn_remote"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
)

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
		contracts     ccipChangesetSolana.CCIPContractsToTransfer
		chainSelector uint64
		expectedError string
	}{
		{
			name:          "All required contracts present",
			state:         makeState(validPubkey, validPubkey),
			contracts:     ccipChangesetSolana.CCIPContractsToTransfer{Router: true},
			chainSelector: 12345,
		},
		{
			name:          "Missing Router contract",
			state:         makeState(zeroPubkey, validPubkey),
			contracts:     ccipChangesetSolana.CCIPContractsToTransfer{Router: true},
			chainSelector: 12345,
			expectedError: "missing required contract Router on chain 12345",
		},
		{
			name:          "Missing FeeQuoter contract",
			state:         makeState(validPubkey, zeroPubkey),
			contracts:     ccipChangesetSolana.CCIPContractsToTransfer{Router: true, FeeQuoter: true},
			chainSelector: 12345,
			expectedError: "missing required contract FeeQuoter on chain 12345",
		},
		{
			name:          "invalid pub key",
			state:         makeState(validPubkey, zeroPubkey),
			contracts:     ccipChangesetSolana.CCIPContractsToTransfer{Router: true, FeeQuoter: true},
			chainSelector: 12345,
			expectedError: "missing required contract FeeQuoter on chain 12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ccipChangesetSolana.ValidateContracts(tt.state, tt.chainSelector, tt.contracts)

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
	selector := chainselectors.TEST_22222222222222222222222222222222222222222222.Selector

	tests := []struct {
		name             string
		env              func(t *testing.T) cldf.Environment
		contractsByChain map[uint64]ccipChangesetSolana.CCIPContractsToTransfer
		expectedError    string
	}{
		{
			name: "No chains found in environment",
			env: func(t *testing.T) cldf.Environment {
				t.Helper()

				e, err := environment.New(t.Context())
				require.NoError(t, err)

				return *e
			},
			expectedError: "no chains found",
		},
		{
			name: "Chain selector not found in environment",
			env: func(t *testing.T) cldf.Environment {
				t.Helper()

				e, err := environment.New(t.Context())
				require.NoError(t, err)

				e.BlockChains = cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{
					cldf_solana.Chain{
						Selector: selector,
					},
				})

				return *e
			},
			contractsByChain: map[uint64]ccipChangesetSolana.CCIPContractsToTransfer{
				99999: {Router: true, FeeQuoter: true},
			},
			expectedError: "chain 99999 not found in environment",
		},
		{
			name: "Invalid chain family",
			env: func(t *testing.T) cldf.Environment {
				t.Helper()

				e, err := environment.New(t.Context())
				require.NoError(t, err)

				e.BlockChains = cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{
					cldf_solana.Chain{
						Selector: selector,
					},
				})

				return *e
			},
			contractsByChain: map[uint64]ccipChangesetSolana.CCIPContractsToTransfer{
				selector: {Router: true, FeeQuoter: true},
			},
			expectedError: "failed to load addresses for chain 12463857294658392847: chain selector 12463857294658392847: chain not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolanaConfig{
				ContractsByChain: tt.contractsByChain,
				MCMSCfg: proposalutils.TimelockConfig{
					MinDelay: 0 * time.Second,
				},
			}

			err := cfg.Validate(tt.env(t))

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

	homeChainSel := chainselectors.TEST_90000001.Selector
	solChainSel := chainselectors.TEST_22222222222222222222222222222222222222222222.Selector

	programsPath := t.TempDir()
	progIDs := soltestutils.LoadCCIPPrograms(t, programsPath)
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{homeChainSel}),
		environment.WithSolanaContainer(t, []uint64{solChainSel}, programsPath, progIDs),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	testhelpers.RegisterNodes(t, env, 4, homeChainSel)

	solChain := env.BlockChains.SolanaChains()[solChainSel]
	nodes, err := deployment.NodeInfo(env.NodeIDs, env.Offchain)
	require.NoError(t, err)

	// Fund account for fees
	testutils.FundAccounts(env.GetContext(), []solana.PrivateKey{*solChain.DeployerKey}, solChain.Client, t)
	err = testhelpers.SavePreloadedSolAddresses(*env, solChainSel)
	require.NoError(t, err)
	solLinkTokenPrivKey, _ := solana.NewRandomPrivateKey()

	e := *env

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
				ChainSelector: solChainSel,
				TokenPrivKey:  solLinkTokenPrivKey,
				TokenDecimals: 9,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			ccipChangesetSolana.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ChainSelector:     solChainSel,
				ContractParamsPerChain: ccipChangesetSolana.ChainContractParams{
					FeeQuoterParams: ccipChangesetSolana.FeeQuoterParams{
						DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil},
					},
					OffRampParams: ccipChangesetSolana.OffRampParams{
						EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
					},
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeploySolanaToken),
			ccipChangesetSolana.DeploySolanaTokenConfig{
				ChainSelector:    solChainSel,
				TokenProgramName: shared.SPL2022Tokens,
				TokenDecimals:    9,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeploySolanaToken),
			ccipChangesetSolana.DeploySolanaTokenConfig{
				ChainSelector:    solChainSel,
				TokenProgramName: shared.SPLTokens,
				TokenDecimals:    9,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			map[uint64]commontypes.MCMSWithTimelockConfigV2{
				solChainSel: {
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
	err = testhelpers.ValidateSolanaState(e, []uint64{solChainSel})
	require.NoError(t, err)
	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	tokenAddressLockRelease := state.SolChains[solChainSel].SPL2022Tokens[0]
	tokenAddressBurnMint := state.SolChains[solChainSel].SPLTokens[0]

	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.InitGlobalConfigTokenPoolProgram),
			ccipChangesetSolana.TokenPoolConfigWithMCM{
				ChainSelector: solChainSel,
				TokenPoolConfigs: []ccipChangesetSolana.TokenPoolConfig{
					{
						TokenPubKey: tokenAddressLockRelease,
						PoolType:    shared.LockReleaseTokenPool,
						Metadata:    shared.CLLMetadata,
					},
					{
						TokenPubKey: tokenAddressBurnMint,
						PoolType:    shared.BurnMintTokenPool,
						Metadata:    shared.CLLMetadata,
					},
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddTokenPoolAndLookupTable),
			ccipChangesetSolana.AddTokenPoolAndLookupTableConfig{
				ChainSelector: solChainSel,
				TokenPoolConfigs: []ccipChangesetSolana.TokenPoolConfig{
					{
						TokenPubKey: tokenAddressLockRelease,
						PoolType:    shared.LockReleaseTokenPool,
						Metadata:    shared.CLLMetadata,
					},
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddTokenPoolAndLookupTable),
			ccipChangesetSolana.AddTokenPoolAndLookupTableConfig{
				ChainSelector: solChainSel,
				TokenPoolConfigs: []ccipChangesetSolana.TokenPoolConfig{
					{
						TokenPubKey: tokenAddressBurnMint,
						PoolType:    shared.BurnMintTokenPool,
						Metadata:    shared.CLLMetadata,
					},
				},
			},
		),
	})
	require.NoError(t, err)
	return e, state
}

func TestTransferCCIPToMCMSWithTimelockSolana(t *testing.T) {
	t.Parallel()
	skipInCI(t)

	e, state := prepareEnvironmentForOwnershipTransfer(t)
	solSelector := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))[0]
	solChain := e.BlockChains.SolanaChains()[solSelector]
	solState := state.SolChains[solSelector]

	tokenAddressLockRelease := solState.SPL2022Tokens[0]
	tokenAddressBurnMint := solState.SPLTokens[0]

	burnMintPoolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenAddressBurnMint, solState.BurnMintTokenPools[shared.CLLMetadata])
	lockReleasePoolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenAddressLockRelease, solState.LockReleaseTokenPools[shared.CLLMetadata])
	timelockSignerPDA, _ := testhelpers.TransferOwnershipSolanaV0_1_1(
		t,
		&e,
		solSelector,
		false,
		ccipChangesetSolana.CCIPContractsToTransfer{
			Router:                true,
			FeeQuoter:             true,
			OffRamp:               true,
			RMNRemote:             true,
			BurnMintTokenPools:    map[string][]solana.PublicKey{shared.CLLMetadata: {tokenAddressBurnMint}},
			LockReleaseTokenPools: map[string][]solana.PublicKey{shared.CLLMetadata: {tokenAddressLockRelease}},
		})

	// 5. Now verify on-chain that each contract’s “config account” authority is the Timelock PDA.
	//    Typically, each contract has its own config account: RouterConfigPDA, FeeQuoterConfigPDA,
	//    Token Pool config PDAs, OffRamp config, etc.
	ctx := context.Background()

	// (A) Check Router ownership -  we need to add retries as the ownership transfer commitment is confirmed and not finalized.
	require.Eventually(t, func() bool {
		routerConfigPDA := solState.RouterConfigPDA
		t.Logf("Checking Router Config PDA ownership data configPDA: %s", routerConfigPDA.String())
		programData := ccip_router.Config{}
		err := solChain.GetAccountDataBorshInto(ctx, routerConfigPDA, &programData)
		require.NoError(t, err)
		return timelockSignerPDA.String() == programData.Owner.String()
	}, 30*time.Second, 5*time.Second, "Router config PDA owner was not changed to timelock signer PDA")

	// (B) Check FeeQuoter ownership
	require.Eventually(t, func() bool {
		feeQuoterConfigPDA := solState.FeeQuoterConfigPDA
		t.Logf("Checking Fee Quoter PDA ownership data configPDA: %s", feeQuoterConfigPDA.String())
		programData := fee_quoter.Config{}
		err := solChain.GetAccountDataBorshInto(ctx, feeQuoterConfigPDA, &programData)
		require.NoError(t, err)
		return timelockSignerPDA.String() == programData.Owner.String()
	}, 30*time.Second, 5*time.Second, "Fee Quoter config PDA owner was not changed to timelock signer PDA")

	// (C) Check OffRamp:
	require.Eventually(t, func() bool {
		offRampConfigPDA := solState.OffRampConfigPDA
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
		rmnRemoteConfigPDA := solState.RMNRemoteConfigPDA
		t.Logf("Checking RMNRemote PDA ownership data configPDA: %s", rmnRemoteConfigPDA.String())
		programData := rmn_remote.Config{}
		err := solChain.GetAccountDataBorshInto(ctx, rmnRemoteConfigPDA, &programData)
		require.NoError(t, err)
		return timelockSignerPDA.String() == programData.Owner.String()
	}, 30*time.Second, 5*time.Second, "RMNRemote config PDA owner was not changed to timelock signer PDA")
}

func TestTransferCCIPFromMCMSWithTimelockSolana(t *testing.T) {
	t.Parallel()
	skipInCI(t)

	e, state := prepareEnvironmentForOwnershipTransfer(t)
	solSelector := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))[0]
	solChain := e.BlockChains.SolanaChains()[solSelector]
	solState := state.SolChains[solSelector]

	tokenAddressLockRelease := solState.SPL2022Tokens[0]
	tokenAddressBurnMint := solState.SPLTokens[0]

	burnMintPoolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenAddressBurnMint, solState.BurnMintTokenPools[shared.CLLMetadata])
	lockReleasePoolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenAddressLockRelease, solState.LockReleaseTokenPools[shared.CLLMetadata])
	timelockSignerPDA, _ := testhelpers.TransferOwnershipSolanaV0_1_1(
		t,
		&e,
		solSelector,
		false,
		ccipChangesetSolana.CCIPContractsToTransfer{
			Router:                true,
			FeeQuoter:             true,
			OffRamp:               true,
			RMNRemote:             true,
			BurnMintTokenPools:    map[string][]solana.PublicKey{shared.CLLMetadata: {tokenAddressBurnMint}},
			LockReleaseTokenPools: map[string][]solana.PublicKey{shared.CLLMetadata: {tokenAddressLockRelease}},
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
					solSelector: {
						Router:                true,
						FeeQuoter:             true,
						OffRamp:               true,
						RMNRemote:             true,
						BurnMintTokenPools:    map[string][]solana.PublicKey{shared.CLLMetadata: {tokenAddressBurnMint}},
						LockReleaseTokenPools: map[string][]solana.PublicKey{shared.CLLMetadata: {tokenAddressLockRelease}},
					},
				},
			},
		),
	})
	require.NoError(t, err)
	// we have to accept separate from the changeset because the proposal needs to execute
	// just spot check that the ownership transfer happened
	config := state.SolChains[solSelector].RouterConfigPDA
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
