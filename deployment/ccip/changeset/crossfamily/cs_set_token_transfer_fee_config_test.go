package crossfamily_test

import (
	"math"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	solfq "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/fee_quoter"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/helpers/pointer"

	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	solstateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/crossfamily"
	ccip_cs_sol_v0_1_1 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
)

const SetTokenTransferFeePriceRegStalenessThreshold = 60 * 60 * 24 * 14 // two weeks in seconds

var SetTokenTransferFeeMcmsConfig = proposalutils.TimelockConfig{MinDelay: 1 * time.Second}

func deploySolanaToken(t *testing.T, tEnv cldf.Environment, solSelector uint64, tokenSymbol string) (cldf.Environment, solana.PublicKey, error) {
	t.Helper()

	e, err := commonchangeset.Apply(t, tEnv,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccip_cs_sol_v0_1_1.DeploySolanaToken),
			ccip_cs_sol_v0_1_1.DeploySolanaTokenConfig{
				MintAmountToAddress: map[string]uint64{},
				TokenProgramName:    shared.SPLTokens,
				ChainSelector:       solSelector,
				TokenDecimals:       9,
				TokenSymbol:         tokenSymbol,
				ATAList:             []string{},
			},
		),
	)
	if err != nil {
		return cldf.Environment{}, solana.PublicKey{}, err
	}

	addresses, err := e.ExistingAddresses.AddressesForChain(solSelector)
	require.NoError(t, err)

	tokenAddress := solstateview.FindSolanaAddress(
		cldf.TypeAndVersion{
			Version: deployment.Version1_0_0,
			Labels:  cldf.NewLabelSet(tokenSymbol),
			Type:    shared.SPLTokens,
		},
		addresses,
	)

	return e, tokenAddress, nil
}

func getSolanaTokenTransferFeeConfig(t *testing.T, tEnv cldf.Environment, srcSelector, dstSelector uint64, fq, tokenPubKey solana.PublicKey) solfq.PerChainPerTokenConfig {
	t.Helper()

	pda, _, err := solstate.FindFqPerChainPerTokenConfigPDA(dstSelector, tokenPubKey, fq)
	require.NoError(t, err)

	var cfg solfq.PerChainPerTokenConfig
	err = tEnv.BlockChains.SolanaChains()[srcSelector].GetAccountDataBorshInto(tEnv.GetContext(), pda, &cfg)
	require.NoError(t, err)

	return cfg
}

func TestSetTokenTransferFeeConfig_Validations(t *testing.T) {
	// Build a mixed env with EVM + non-EVM so we can validate both families in one table-driven test
	env, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithCCIPSolanaContractVersion(ccip_cs_sol_v0_1_1.SolanaContractV0_1_1),
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSolChains(1),
	)

	// EVM selectors
	evmSelectors := env.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))
	require.GreaterOrEqual(t, len(evmSelectors), 2, "test requires at least 2 EVM chains in the memory env")
	evmSrc := evmSelectors[0]
	evmDst := evmSelectors[1]

	// Solana selectors
	solSelectors := env.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilySolana))
	require.GreaterOrEqual(t, len(solSelectors), 1, "test requires at least 1 Solana chain in the memory env")
	solSrc := solSelectors[0]

	// Helper vars
	evmAddr := utils.RandomAddress().String()
	solAddr := solana.SolMint.String()

	// Configure MCMS on Solana
	_, _ = testhelpers.TransferOwnershipSolanaV0_1_1(t, &env.Env, solSrc, true,
		ccip_cs_sol_v0_1_1.CCIPContractsToTransfer{
			FeeQuoter: true,
			Router:    true,
			OffRamp:   true,
		})
	_, err := ccip_cs_sol_v0_1_1.FetchTimelockSigner(env.Env, solSrc)
	require.NoError(t, err)

	// Define test cases
	tests := []struct {
		MsgStr string
		Config crossfamily.SetTokenTransferFeeConfigInput
		ErrStr string
	}{
		{
			ErrStr: "MCMS config is required",
			MsgStr: "MCMS required when inputs present",
			Config: crossfamily.SetTokenTransferFeeConfigInput{
				InputsByChain: map[uint64]map[uint64]crossfamily.TokenTransferFeeConfigArgs{
					evmSrc: {
						evmDst: {
							TokenAddressToFeeConfig: map[string]crossfamily.OptionalTokenTransferFeeConfig{
								evmAddr: {},
							},
						},
					},
				},
				MCMS: nil, // required -> expect error
			},
		},
		{
			ErrStr: "invalid hex EVM address detected",
			MsgStr: "Invalid EVM hex token address",
			Config: crossfamily.SetTokenTransferFeeConfigInput{
				InputsByChain: map[uint64]map[uint64]crossfamily.TokenTransferFeeConfigArgs{
					evmSrc: {
						evmDst: {
							TokenAddressToFeeConfig: map[string]crossfamily.OptionalTokenTransferFeeConfig{
								"not-a-hex-address": {},
							},
						},
					},
				},
				MCMS: &SetTokenTransferFeeMcmsConfig,
			},
		},
		{
			ErrStr: "failed to validate dst chain",
			MsgStr: "Missing destination chain selector (dst=0)",
			Config: crossfamily.SetTokenTransferFeeConfigInput{
				InputsByChain: map[uint64]map[uint64]crossfamily.TokenTransferFeeConfigArgs{
					evmSrc: {
						0: {
							TokenAddressToFeeConfig: map[string]crossfamily.OptionalTokenTransferFeeConfig{
								evmAddr: {},
							},
						},
					},
				},
				MCMS: &SetTokenTransferFeeMcmsConfig,
			},
		},
		{
			ErrStr: "invalid base58 solana address detected",
			MsgStr: "Solana: invalid base58 address",
			Config: crossfamily.SetTokenTransferFeeConfigInput{
				InputsByChain: map[uint64]map[uint64]crossfamily.TokenTransferFeeConfigArgs{
					solSrc: {
						evmDst: {
							TokenAddressToFeeConfig: map[string]crossfamily.OptionalTokenTransferFeeConfig{
								"this-is-not-base58!!": {},
							},
						},
					},
				},
				MCMS: &SetTokenTransferFeeMcmsConfig,
			},
		},
		{
			ErrStr: "failed to validate src chain",
			MsgStr: "Solana: selector not in env (src validation)",
			Config: crossfamily.SetTokenTransferFeeConfigInput{
				InputsByChain: map[uint64]map[uint64]crossfamily.TokenTransferFeeConfigArgs{
					chainsel.SOLANA_DEVNET.Selector: {
						evmDst: {
							TokenAddressToFeeConfig: map[string]crossfamily.OptionalTokenTransferFeeConfig{
								solAddr: {},
							},
						},
					},
				},
				MCMS: &SetTokenTransferFeeMcmsConfig,
			},
		},
		{
			ErrStr: "destination chain cannot be the same as source chain",
			MsgStr: "Reject same src/dst selector",
			Config: crossfamily.SetTokenTransferFeeConfigInput{
				InputsByChain: map[uint64]map[uint64]crossfamily.TokenTransferFeeConfigArgs{
					evmSrc: {
						evmSrc: { // same as src
							TokenAddressToFeeConfig: map[string]crossfamily.OptionalTokenTransferFeeConfig{
								evmAddr: {},
							},
						},
					},
				},
				MCMS: &SetTokenTransferFeeMcmsConfig,
			},
		},
		{
			ErrStr: "unsupported EVM version",
			MsgStr: "Unsupported EVM version hint fails fast",
			Config: crossfamily.SetTokenTransferFeeConfigInput{
				VersionHints: &crossfamily.OptionalVersions{
					Solana: pointer.To(ccip_cs_sol_v0_1_1.VersionSolanaV0_1_1), // valid solana (doesn't matter)
					Evm:    pointer.To("1.4.9"),                                // bogus/unsupported
				},
				InputsByChain: map[uint64]map[uint64]crossfamily.TokenTransferFeeConfigArgs{
					evmSrc: {
						evmDst: {
							TokenAddressToFeeConfig: map[string]crossfamily.OptionalTokenTransferFeeConfig{
								evmAddr: {},
							},
						},
					},
				},
				MCMS: &SetTokenTransferFeeMcmsConfig,
			},
		},
		{
			ErrStr: "unsupported Solana version",
			MsgStr: "Unsupported Solana version hint fails fast",
			Config: crossfamily.SetTokenTransferFeeConfigInput{
				VersionHints: &crossfamily.OptionalVersions{
					Evm:    pointer.To(deployment.Version1_6_0.String()),
					Solana: pointer.To("v0_0_9"), // bogus/unsupported
				},
				InputsByChain: map[uint64]map[uint64]crossfamily.TokenTransferFeeConfigArgs{
					solSrc: {
						evmDst: {
							TokenAddressToFeeConfig: map[string]crossfamily.OptionalTokenTransferFeeConfig{
								solAddr: {},
							},
						},
					},
				},
				MCMS: &SetTokenTransferFeeMcmsConfig,
			},
		},
	}

	// Run test cases
	for _, test := range tests {
		t.Run(test.MsgStr, func(t *testing.T) {
			_, err := commonchangeset.Apply(t, env.Env,
				commonchangeset.Configure(
					crossfamily.SetTokenTransferFeeConfig,
					test.Config,
				),
			)
			require.Error(t, err)
			require.ErrorContains(t, err, test.ErrStr)
		})
	}
}

func TestSetTokenTransferFeeConfig_EmptyConfigIsGracefullyHandled(t *testing.T) {
	env, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithCCIPSolanaContractVersion(ccip_cs_sol_v0_1_1.SolanaContractV0_1_1),
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSolChains(1),
	)

	_, err := commonchangeset.Apply(t, env.Env,
		commonchangeset.Configure(
			crossfamily.SetTokenTransferFeeConfig,
			crossfamily.SetTokenTransferFeeConfigInput{},
		),
	)
	require.NoError(t, err)
}

func TestSetTokenTransferFeeConfig_EVM_V1_6_0_Only(t *testing.T) {
	// Setup EVM environment
	env, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))

	// EVM selectors
	evm := env.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))
	require.GreaterOrEqual(t, len(evm), 2)
	src, dst := evm[0], evm[1]

	// Load state and transfer ownership to timelock so MCMS can execute
	state, err := stateview.LoadOnchainState(env.Env)
	require.NoError(t, err)
	testhelpers.TransferToTimelock(t, env, state, []uint64{src, dst}, true)

	// Helper vars
	link := state.MustGetEVMChainState(src).LinkToken.Address()
	opts := &bind.CallOpts{Context: env.Env.GetContext()}

	// Define the token transfer fee config
	cfg := crossfamily.SetTokenTransferFeeConfigInput{
		InputsByChain: map[uint64]map[uint64]crossfamily.TokenTransferFeeConfigArgs{
			src: {
				dst: {
					TokenAddressToFeeConfig: map[string]crossfamily.OptionalTokenTransferFeeConfig{
						link.Hex(): {
							// partial fields: defaults should be auto-filled by the EVM changeset when none exist
							MinFeeUsdCents:  pointer.To(uint32(800)),
							DestGasOverhead: pointer.To(uint32(100)),
						},
					},
				},
			},
		},
		VersionHints: &crossfamily.OptionalVersions{
			Evm: pointer.To(deployment.Version1_6_0.String()),
		},
		MCMS: &SetTokenTransferFeeMcmsConfig,
	}

	// Apply once
	_, err = commonchangeset.Apply(t, env.Env, commonchangeset.Configure(crossfamily.SetTokenTransferFeeConfig, cfg))
	require.NoError(t, err)

	// Verify on FeeQuoter(src) for (dst, link)
	cfg1, err := state.MustGetEVMChainState(src).FeeQuoter.GetTokenTransferFeeConfig(opts, dst, link)
	require.NoError(t, err)
	require.Equal(t, uint32(math.MaxUint32), cfg1.MaxFeeUSDCents)
	require.Equal(t, uint32(800), cfg1.MinFeeUSDCents)
	require.Equal(t, uint32(100), cfg1.DestGasOverhead)
	require.Equal(t, uint32(32), cfg1.DestBytesOverhead)
	require.Equal(t, uint16(0), cfg1.DeciBps)
	require.True(t, cfg1.IsEnabled)

	// Re-apply the exact same config: should be a no-op and succeed
	_, err = commonchangeset.Apply(t, env.Env, commonchangeset.Configure(crossfamily.SetTokenTransferFeeConfig, cfg))
	require.NoError(t, err)

	// Verify on FeeQuoter(src) for (dst, link)
	cfg2, err := state.MustGetEVMChainState(src).FeeQuoter.GetTokenTransferFeeConfig(opts, dst, link)
	require.NoError(t, err)
	require.Equal(t, cfg1, cfg2)
}

func TestSetTokenTransferFeeConfig_EVM_V1_5_1_Only(t *testing.T) {
	// Setup EVM environment
	env, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithPrerequisiteDeploymentOnly(&changeset.V1_5DeploymentConfig{
			// NOTE: this property needs to be defined otherwise we will encounter an
			// error. For now it's set to a value that was found in another test case
			PriceRegStalenessThreshold: SetTokenTransferFeePriceRegStalenessThreshold,
		}),
	)

	// EVM selectors
	evm := env.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))
	require.GreaterOrEqual(t, len(evm), 2)
	src, dst := evm[0], evm[1]

	// Load state and transfer ownership to timelock so MCMS can execute
	state, err := stateview.LoadOnchainState(env.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	// Wire up a bi-directional lane
	env.Env = v1_5.AddLanes(t, env.Env, state, []testhelpers.SourceDestPair{
		{SourceChainSelector: src, DestChainSelector: dst},
		{SourceChainSelector: dst, DestChainSelector: src},
	})

	// Take a fresh snapshot of the state and transfer the OnRamps to timelock
	state, err = stateview.LoadOnchainState(env.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)
	env.Env, err = commonchangeset.Apply(t, env.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
			commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: map[uint64][]common.Address{
					src: {state.MustGetEVMChainState(src).EVM2EVMOnRamp[dst].Address()},
					dst: {state.MustGetEVMChainState(dst).EVM2EVMOnRamp[src].Address()},
				},
				MCMSConfig: SetTokenTransferFeeMcmsConfig,
			},
		),
	)
	require.NoError(t, err)

	// Helper vars
	link := state.MustGetEVMChainState(src).LinkToken.Address()
	opts := &bind.CallOpts{Context: env.Env.GetContext()}

	// Define the token transfer fee config
	cfg := crossfamily.SetTokenTransferFeeConfigInput{
		InputsByChain: map[uint64]map[uint64]crossfamily.TokenTransferFeeConfigArgs{
			src: {
				dst: {
					TokenAddressToFeeConfig: map[string]crossfamily.OptionalTokenTransferFeeConfig{
						link.Hex(): {
							DestBytesOverhead: pointer.To(uint32(32)),
							DestGasOverhead:   pointer.To(uint32(10)),
							MaxFeeUsdCents:    pointer.To(uint32(math.MaxUint32)),
							MinFeeUsdCents:    pointer.To(uint32(800)),
							DeciBps:           pointer.To(uint16(0)),
							IsEnabled:         pointer.To(true),
						},
					},
				},
			},
		},
		VersionHints: &crossfamily.OptionalVersions{
			Evm: pointer.To(deployment.Version1_5_1.String()),
		},
		MCMS: &SetTokenTransferFeeMcmsConfig,
	}

	// Apply once
	_, err = commonchangeset.Apply(t, env.Env, commonchangeset.Configure(crossfamily.SetTokenTransferFeeConfig, cfg))
	require.NoError(t, err)

	// Verify on OnRamp[src][dst] for link
	cfg1, err := state.MustGetEVMChainState(src).EVM2EVMOnRamp[dst].GetTokenTransferFeeConfig(opts, link)
	require.NoError(t, err)
	require.Equal(t, uint32(math.MaxUint32), cfg1.MaxFeeUSDCents, cfg1)
	require.Equal(t, uint32(800), cfg1.MinFeeUSDCents)
	require.Equal(t, uint32(32), cfg1.DestBytesOverhead)
	require.Equal(t, uint32(10), cfg1.DestGasOverhead)
	require.Equal(t, uint16(0), cfg1.DeciBps)
	require.True(t, cfg1.IsEnabled)

	// Re-apply the exact same config: should be a no-op and succeed
	_, err = commonchangeset.Apply(t, env.Env, commonchangeset.Configure(crossfamily.SetTokenTransferFeeConfig, cfg))
	require.NoError(t, err)

	// Verify on OnRamp[src][dst] for link
	cfg2, err := state.MustGetEVMChainState(src).EVM2EVMOnRamp[dst].GetTokenTransferFeeConfig(opts, link)
	require.NoError(t, err)
	require.Equal(t, cfg1, cfg2)
}

func TestSetTokenTransferFeeConfig_Solana_V0_1_0_Only(t *testing.T) {
	// Setup Solana environment
	env, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithCCIPSolanaContractVersion(ccip_cs_sol_v0_1_1.SolanaContractV0_1_1),
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSolChains(1),
	)

	// EVM selectors
	evm := env.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))
	require.GreaterOrEqual(t, len(evm), 1)
	dst := evm[0]

	// SOL selectors
	sol := env.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilySolana))
	require.GreaterOrEqual(t, len(sol), 1)
	src := sol[0]

	// Deploy a Solana test token X
	e, tokenAddressX, err := deploySolanaToken(t, env.Env, src, "TEST_TOKEN_X")
	require.NoError(t, err)
	env.Env = e

	// Deploy a Solana test token Y
	e, tokenAddressY, err := deploySolanaToken(t, env.Env, src, "TEST_TOKEN_Y")
	require.NoError(t, err)
	env.Env = e

	// Configure MCMS on Solana
	_, _ = testhelpers.TransferOwnershipSolanaV0_1_1(t, &env.Env, src, true,
		ccip_cs_sol_v0_1_1.CCIPContractsToTransfer{
			FeeQuoter: true,
			Router:    true,
			OffRamp:   true,
		})
	_, err = ccip_cs_sol_v0_1_1.FetchTimelockSigner(env.Env, src)
	require.NoError(t, err)

	// Define the token transfer fee config
	cfg := crossfamily.SetTokenTransferFeeConfigInput{
		InputsByChain: map[uint64]map[uint64]crossfamily.TokenTransferFeeConfigArgs{
			src: {
				dst: {
					TokenAddressToFeeConfig: map[string]crossfamily.OptionalTokenTransferFeeConfig{
						tokenAddressX.String(): {
							// partial fields: defaults should be auto-filled by the SOL changeset when none exist
							DestGasOverhead: pointer.To(uint32(10)),
						},
						tokenAddressY.String(): {
							// partial fields: defaults should be auto-filled by the SOL changeset when none exist
							DestBytesOverhead: pointer.To(uint32(64)),
							DestGasOverhead:   pointer.To(uint32(10)),
						},
					},
				},
			},
		},
		VersionHints: &crossfamily.OptionalVersions{
			Solana: pointer.To(ccip_cs_sol_v0_1_1.VersionSolanaV0_1_1),
		},
		MCMS: &SetTokenTransferFeeMcmsConfig,
	}

	// Set the token transfer fee config
	e, err = commonchangeset.Apply(t, env.Env, commonchangeset.Configure(crossfamily.SetTokenTransferFeeConfig, cfg))
	require.NoError(t, err)
	env.Env = e

	// Refresh state
	state, err := stateview.LoadOnchainState(env.Env)
	require.NoError(t, err)

	// Verify on FeeQuoter(src) for (dst, tokenX)
	solCfgX := getSolanaTokenTransferFeeConfig(t, env.Env, src, dst, state.SolChains[src].FeeQuoter, tokenAddressX)
	require.Equal(t, tokenAddressX, solCfgX.Mint)
	require.Equal(t, uint32(math.MaxUint32), solCfgX.TokenTransferConfig.MaxFeeUsdcents)
	require.Equal(t, uint32(175_000), solCfgX.TokenTransferConfig.MinFeeUsdcents)
	require.Equal(t, uint32(32), solCfgX.TokenTransferConfig.DestBytesOverhead)
	require.Equal(t, uint32(10), solCfgX.TokenTransferConfig.DestGasOverhead)
	require.Equal(t, uint16(0), solCfgX.TokenTransferConfig.DeciBps)
	require.True(t, solCfgX.TokenTransferConfig.IsEnabled)

	// Verify on FeeQuoter(src) for (dst, tokenY)
	solCfgY := getSolanaTokenTransferFeeConfig(t, env.Env, src, dst, state.SolChains[src].FeeQuoter, tokenAddressY)
	require.Equal(t, tokenAddressY, solCfgY.Mint)
	require.Equal(t, uint32(math.MaxUint32), solCfgY.TokenTransferConfig.MaxFeeUsdcents)
	require.Equal(t, uint32(175_000), solCfgY.TokenTransferConfig.MinFeeUsdcents)
	require.Equal(t, uint32(64), solCfgY.TokenTransferConfig.DestBytesOverhead)
	require.Equal(t, uint32(10), solCfgY.TokenTransferConfig.DestGasOverhead)
	require.Equal(t, uint16(0), solCfgY.TokenTransferConfig.DeciBps)
	require.True(t, solCfgY.TokenTransferConfig.IsEnabled)
}

func TestSetTokenTransferFeeConfig_MixedFamilies_SingleApply(t *testing.T) {
	// Build a mixed env with EVM + non-EVM
	env, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithCCIPSolanaContractVersion(ccip_cs_sol_v0_1_1.SolanaContractV0_1_1),
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSolChains(1),
	)

	// EVM selectors
	evm := env.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))
	require.GreaterOrEqual(t, len(evm), 2)
	evmSrc, evmDst := evm[0], evm[1]

	// SOL selectors
	sol := env.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilySolana))
	require.GreaterOrEqual(t, len(sol), 1)
	solSrc := sol[0]

	// EVM: transfer to timelock so we can MCMS-apply
	state, err := stateview.LoadOnchainState(env.Env)
	require.NoError(t, err)
	testhelpers.TransferToTimelock(t, env, state, []uint64{evmSrc, evmDst}, true)

	// Solana: transfer to timelock so we can MCMS-apply
	_, _ = testhelpers.TransferOwnershipSolanaV0_1_1(t, &env.Env, solSrc, true, ccip_cs_sol_v0_1_1.CCIPContractsToTransfer{
		FeeQuoter: true,
		Router:    true,
		OffRamp:   true,
	})
	_, err = ccip_cs_sol_v0_1_1.FetchTimelockSigner(env.Env, solSrc)
	require.NoError(t, err)

	// EVM: token to configure = LINK on source chain
	link := state.MustGetEVMChainState(evmSrc).LinkToken.Address()
	opts := &bind.CallOpts{Context: env.Env.GetContext()}

	// Solana: token to configure = deploy a test token
	var mint solana.PublicKey
	env.Env, mint, err = deploySolanaToken(t, env.Env, solSrc, "TEST_TOKEN_MIXED")
	require.NoError(t, err)

	// Define the token transfer fee config
	cfg := crossfamily.SetTokenTransferFeeConfigInput{
		InputsByChain: map[uint64]map[uint64]crossfamily.TokenTransferFeeConfigArgs{
			// EVM src -> EVM dst
			evmSrc: {
				evmDst: {
					TokenAddressToFeeConfig: map[string]crossfamily.OptionalTokenTransferFeeConfig{
						link.Hex(): {
							MinFeeUsdCents:    pointer.To(uint32(12)),
							DeciBps:           pointer.To(uint16(7)),
							DestGasOverhead:   pointer.To(uint32(222)),
							DestBytesOverhead: pointer.To(uint32(512)),
							IsEnabled:         pointer.To(true),
						},
					},
				},
			},
			// Solana src -> EVM dst
			solSrc: {
				evmDst: {
					TokenAddressToFeeConfig: map[string]crossfamily.OptionalTokenTransferFeeConfig{
						mint.String(): {
							// partial: defaults will fill any omitted fields
							MinFeeUsdCents:    pointer.To(uint32(900)),
							DestGasOverhead:   pointer.To(uint32(100)),
							DestBytesOverhead: pointer.To(uint32(128)),
							// DeciBps/IsEnabled left nil -> defaults (0 / true)
						},
					},
				},
			},
		},
		VersionHints: &crossfamily.OptionalVersions{
			Solana: pointer.To(ccip_cs_sol_v0_1_1.VersionSolanaV0_1_1),
			Evm:    pointer.To(deployment.Version1_6_0.String()),
		},
		MCMS: &SetTokenTransferFeeMcmsConfig,
	}

	// Apply
	_, err = commonchangeset.Apply(t, env.Env, commonchangeset.Configure(crossfamily.SetTokenTransferFeeConfig, cfg))
	require.NoError(t, err)

	// Refresh state
	state, err = stateview.LoadOnchainState(env.Env)
	require.NoError(t, err)

	// Get EVM token transfer fee config
	evmCfg, err := state.MustGetEVMChainState(evmSrc).FeeQuoter.GetTokenTransferFeeConfig(opts, evmDst, link)
	require.NoError(t, err)

	// Verify EVM config
	require.Equal(t, uint32(math.MaxUint32), evmCfg.MaxFeeUSDCents)
	require.Equal(t, uint32(12), evmCfg.MinFeeUSDCents)
	require.Equal(t, uint32(512), evmCfg.DestBytesOverhead)
	require.Equal(t, uint32(222), evmCfg.DestGasOverhead)
	require.Equal(t, uint16(7), evmCfg.DeciBps)
	require.True(t, evmCfg.IsEnabled)

	// Verify Solana config
	solCfg := getSolanaTokenTransferFeeConfig(t, env.Env, solSrc, evmDst, state.SolChains[solSrc].FeeQuoter, mint)
	require.Equal(t, mint, solCfg.Mint)
	require.Equal(t, uint32(math.MaxUint32), solCfg.TokenTransferConfig.MaxFeeUsdcents) // default (not provided)
	require.Equal(t, uint32(900), solCfg.TokenTransferConfig.MinFeeUsdcents)            // set explicitly
	require.Equal(t, uint32(128), solCfg.TokenTransferConfig.DestBytesOverhead)         // set explicitly
	require.Equal(t, uint32(100), solCfg.TokenTransferConfig.DestGasOverhead)           // set explicitly
	require.Equal(t, uint16(0), solCfg.TokenTransferConfig.DeciBps)                     // default (not provided)
	require.True(t, solCfg.TokenTransferConfig.IsEnabled)                               // defaulted to true
}
