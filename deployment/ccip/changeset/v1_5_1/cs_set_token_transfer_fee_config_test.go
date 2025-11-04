package v1_5_1_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/v1_5"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/helpers/pointer"
)

// Two weeks in seconds
const SetTokenTransferFeePriceRegStalenessThreshold = 60 * 60 * 24 * 14

// Helper to read back a token's config from the onramp
func ReadTokenTransferFeeConfig(t *testing.T, e testhelpers.DeployedEnv, srcSelector, dstSelector uint64, token common.Address) (evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfig, error) {
	t.Helper()

	s, err := stateview.LoadOnchainState(e.Env, stateview.WithLoadLegacyContracts(true))
	if err != nil {
		return evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfig{}, err
	}
	chainState, ok := s.EVMChainState(srcSelector)
	if !ok {
		return evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfig{}, fmt.Errorf("no EVM chain state for %d", srcSelector)
	}
	onramp, ok := chainState.EVM2EVMOnRamp[dstSelector]
	if !ok {
		return evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfig{}, fmt.Errorf("no onramp %d -> %d", srcSelector, dstSelector)
	}

	return onramp.GetTokenTransferFeeConfig(&bind.CallOpts{Context: e.Env.GetContext()}, token)
}

func TestSetTokenTransferFeeConfig_Validations(t *testing.T) {
	t.Parallel()

	// Spin up a memory env with minimal prerequisite deployment
	e, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithPrerequisiteDeploymentOnly(&changeset.V1_5DeploymentConfig{
			// NOTE: this property needs to be defined otherwise we will encounter an
			// error. For now it's set to a value that was found in another test case
			PriceRegStalenessThreshold: SetTokenTransferFeePriceRegStalenessThreshold,
		}),
	)

	// Ensure the chains exist
	allChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.Len(t, allChains, 2)
	src := allChains[0]
	dst := allChains[1]

	// Take a snapshot of the current state
	state, err := stateview.LoadOnchainState(e.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	// Wire up a bi-directional lane
	e.Env = v1_5.AddLanes(t, e.Env, state, []testhelpers.SourceDestPair{
		{SourceChainSelector: src, DestChainSelector: dst},
		{SourceChainSelector: dst, DestChainSelector: src},
	})

	// Define helper vars
	mcmCfg := &proposalutils.TimelockConfig{MinDelay: 0 * time.Second}
	tokenA := utils.RandomAddress()
	tokenB := utils.RandomAddress()

	// Define test cases
	tests := []struct {
		Config v1_5_1.SetTokenTransferFeeConfig
		Msg    string
		Err    string
	}{
		{
			Msg: "Nonexistent chain selector",
			Err: "failed to validate src chain",
			Config: v1_5_1.SetTokenTransferFeeConfig{
				MCMS: mcmCfg,
				InputsByChain: map[uint64]map[uint64]v1_5_1.SetTokenTransferFeeArgs{
					0: {},
				},
			},
		},
		{
			Msg: "Invalid EVM chain selector",
			Err: fmt.Sprintf("selector %d does not exist in environment", chain_selectors.SOLANA_DEVNET.Selector),
			Config: v1_5_1.SetTokenTransferFeeConfig{
				MCMS: mcmCfg,
				InputsByChain: map[uint64]map[uint64]v1_5_1.SetTokenTransferFeeArgs{
					chain_selectors.SOLANA_DEVNET.Selector: {},
				},
			},
		},
		{
			Msg: "Duplicate addresses in reset list",
			Err: "duplicate address in TokensToUseDefaultFeeConfigs",
			Config: v1_5_1.SetTokenTransferFeeConfig{
				MCMS: mcmCfg,
				InputsByChain: map[uint64]map[uint64]v1_5_1.SetTokenTransferFeeArgs{
					src: {
						dst: {
							TokensToUseDefaultFeeConfigs: []common.Address{tokenA, tokenA},
						},
					},
				},
			},
		},
		{
			Msg: "Zero address in reset list",
			Err: "zero address not allowed in TokensToUseDefaultFeeConfigs",
			Config: v1_5_1.SetTokenTransferFeeConfig{
				MCMS: mcmCfg,
				InputsByChain: map[uint64]map[uint64]v1_5_1.SetTokenTransferFeeArgs{
					src: {
						dst: {
							TokensToUseDefaultFeeConfigs: []common.Address{utils.ZeroAddress},
						},
					},
				},
			},
		},
		{
			Msg: "Zero token in update args",
			Err: "zero address not allowed in TokenTransferFeeConfigArgs",
			Config: v1_5_1.SetTokenTransferFeeConfig{
				MCMS: mcmCfg,
				InputsByChain: map[uint64]map[uint64]v1_5_1.SetTokenTransferFeeArgs{
					src: {
						dst: {
							TokenTransferFeeConfigArgs: map[common.Address]v1_5_1.TokenTransferFeeArgs{
								utils.ZeroAddress: {},
							},
						},
					},
				},
			},
		},
		{
			Msg: "Same address in updates and resets",
			Err: "the same address cannot be referenced in both TokensToUseDefaultFeeConfigs and TokenTransferFeeConfigArgs",
			Config: v1_5_1.SetTokenTransferFeeConfig{
				MCMS: mcmCfg,
				InputsByChain: map[uint64]map[uint64]v1_5_1.SetTokenTransferFeeArgs{
					src: {
						dst: {
							TokensToUseDefaultFeeConfigs: []common.Address{tokenB},
							TokenTransferFeeConfigArgs: map[common.Address]v1_5_1.TokenTransferFeeArgs{
								tokenB: {
									MinFeeUSDCents:            pointer.To(uint32(1)),
									MaxFeeUSDCents:            pointer.To(uint32(2)),
									DeciBps:                   pointer.To(uint16(10)),
									DestGasOverhead:           pointer.To(uint32(100)),
									DestBytesOverhead:         pointer.To(uint32(200)),
									AggregateRateLimitEnabled: pointer.To(true),
								},
							},
						},
					},
				},
			},
		},
	}

	// Run all tests
	for _, tt := range tests {
		t.Run(tt.Msg, func(t *testing.T) {
			_, err := commonchangeset.Apply(t, e.Env,
				commonchangeset.Configure(
					v1_5_1.SetTokenTransferFeeConfigChangeset,
					tt.Config,
				),
			)
			require.Error(t, err)
			require.ErrorContains(t, err, tt.Err)
		})
	}
}

func TestSetTokenTransferFeeConfig_EmptyConfigIsGracefullyHandled(t *testing.T) {
	// Spin up a memory env with minimal prerequisite deployment
	e, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithPrerequisiteDeploymentOnly(&changeset.V1_5DeploymentConfig{
			// NOTE: this property needs to be defined otherwise we will encounter an
			// error. For now it's set to a value that was found in another test case
			PriceRegStalenessThreshold: SetTokenTransferFeePriceRegStalenessThreshold,
		}),
	)

	// Ensure the chains exist
	allChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.Len(t, allChains, 2)
	src := allChains[0]
	dst := allChains[1]

	// Take a snapshot of the current state
	state, err := stateview.LoadOnchainState(e.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	// Wire up a bi-directional lane
	e.Env = v1_5.AddLanes(t, e.Env, state, []testhelpers.SourceDestPair{
		{SourceChainSelector: src, DestChainSelector: dst},
		{SourceChainSelector: dst, DestChainSelector: src},
	})

	// Running the changeset with an empty config should exit early and do nothing
	_, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			v1_5_1.SetTokenTransferFeeConfigChangeset,
			v1_5_1.SetTokenTransferFeeConfig{},
		),
	)
	require.NoError(t, err)
}

func TestSetTokenTransferFeeConfig_Execution_WithoutMCMS(t *testing.T) {
	// Spin up a memory env with minimal prerequisite deployment
	e, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithPrerequisiteDeploymentOnly(&changeset.V1_5DeploymentConfig{
			// NOTE: this property needs to be defined otherwise we will encounter an
			// error. For now it's set to a value that was found in another test case
			PriceRegStalenessThreshold: SetTokenTransferFeePriceRegStalenessThreshold,
		}),
	)

	// Ensure the chains exist
	allChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.Len(t, allChains, 2)
	src := allChains[0]
	dst := allChains[1]

	// Take a snapshot of the current state
	state, err := stateview.LoadOnchainState(e.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	// Wire up a bi-directional lane
	e.Env = v1_5.AddLanes(t, e.Env, state, []testhelpers.SourceDestPair{
		{SourceChainSelector: src, DestChainSelector: dst},
		{SourceChainSelector: dst, DestChainSelector: src},
	})

	// Token under test (enable and verify)
	tokenA := utils.RandomAddress()

	// Run the changeset without MCMS
	e.Env, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			v1_5_1.SetTokenTransferFeeConfigChangeset,
			v1_5_1.SetTokenTransferFeeConfig{
				MCMS: nil, // direct execution
				InputsByChain: map[uint64]map[uint64]v1_5_1.SetTokenTransferFeeArgs{
					src: {
						dst: {
							TokensToUseDefaultFeeConfigs: []common.Address{},
							TokenTransferFeeConfigArgs: map[common.Address]v1_5_1.TokenTransferFeeArgs{
								tokenA: {
									MinFeeUSDCents:            pointer.To(uint32(100)),
									MaxFeeUSDCents:            pointer.To(uint32(5000)),
									DeciBps:                   pointer.To(uint16(25)),
									DestGasOverhead:           pointer.To(uint32(100_000)),
									DestBytesOverhead:         pointer.To(uint32(1200)),
									AggregateRateLimitEnabled: pointer.To(true),
								},
							},
						},
					},
				},
			},
		),
	)
	require.NoError(t, err, "direct execution should succeed")

	// Verify that the config was set
	cfgA, err := ReadTokenTransferFeeConfig(t, e, src, dst, tokenA)
	require.NoError(t, err)
	require.True(t, cfgA.IsEnabled)
	require.Equal(t, uint32(100), cfgA.MinFeeUSDCents)
	require.Equal(t, uint32(5000), cfgA.MaxFeeUSDCents)
	require.Equal(t, uint16(25), cfgA.DeciBps)
	require.Equal(t, uint32(100_000), cfgA.DestGasOverhead)
	require.Equal(t, uint32(1200), cfgA.DestBytesOverhead)
	require.True(t, cfgA.AggregateRateLimitEnabled)
}

func TestSetTokenTransferFeeConfig_Execution_WithMCMS(t *testing.T) {
	// Spin up a memory env with minimal prerequisite deployment
	e, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithPrerequisiteDeploymentOnly(&changeset.V1_5DeploymentConfig{
			// NOTE: this property needs to be defined otherwise we will encounter an
			// error. For now it's set to a value that was found in another test case
			PriceRegStalenessThreshold: SetTokenTransferFeePriceRegStalenessThreshold,
		}),
	)

	// Ensure the chains exist
	allChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.Len(t, allChains, 2)
	src := allChains[0]
	dst := allChains[1]

	// Take a snapshot of the current state
	state, err := stateview.LoadOnchainState(e.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	// Wire up a bi-directional lane
	e.Env = v1_5.AddLanes(t, e.Env, state, []testhelpers.SourceDestPair{
		{SourceChainSelector: src, DestChainSelector: dst},
		{SourceChainSelector: dst, DestChainSelector: src},
	})

	// Define helper vars
	mcmCfg := proposalutils.TimelockConfig{MinDelay: 0 * time.Second}
	tokenA := utils.RandomAddress()
	tokenB := utils.RandomAddress() // will be reset via MCMS

	// Start by enabling tokenA directly so we have on-chain values to partially override via MCMS
	e.Env, err = commonchangeset.Apply(t, e.Env, commonchangeset.Configure(
		v1_5_1.SetTokenTransferFeeConfigChangeset,
		v1_5_1.SetTokenTransferFeeConfig{
			MCMS: nil,
			InputsByChain: map[uint64]map[uint64]v1_5_1.SetTokenTransferFeeArgs{
				src: {
					dst: {
						TokensToUseDefaultFeeConfigs: []common.Address{},
						TokenTransferFeeConfigArgs: map[common.Address]v1_5_1.TokenTransferFeeArgs{
							tokenA: {
								MinFeeUSDCents:            pointer.To(uint32(100)),
								MaxFeeUSDCents:            pointer.To(uint32(5000)),
								DeciBps:                   pointer.To(uint16(25)),
								DestGasOverhead:           pointer.To(uint32(100_000)),
								DestBytesOverhead:         pointer.To(uint32(1200)),
								AggregateRateLimitEnabled: pointer.To(true),
							},
						},
					},
				},
			},
		},
	))
	require.NoError(t, err)

	// Take a fresh snapshot of the state and transfer the OnRamps to timelock
	state, err = stateview.LoadOnchainState(e.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)
	e.Env, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
			commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: map[uint64][]common.Address{
					src: {state.MustGetEVMChainState(src).EVM2EVMOnRamp[dst].Address()},
					dst: {state.MustGetEVMChainState(dst).EVM2EVMOnRamp[src].Address()},
				},
				MCMSConfig: mcmCfg,
			},
		),
	)
	require.NoError(t, err)

	// Now mutate via MCMS (partial fields; others fallback to on-chain values)
	e.Env, err = commonchangeset.Apply(t, e.Env, commonchangeset.Configure(
		v1_5_1.SetTokenTransferFeeConfigChangeset,
		v1_5_1.SetTokenTransferFeeConfig{
			MCMS: &mcmCfg,
			InputsByChain: map[uint64]map[uint64]v1_5_1.SetTokenTransferFeeArgs{
				src: {
					dst: {
						TokensToUseDefaultFeeConfigs: []common.Address{tokenB},
						TokenTransferFeeConfigArgs: map[common.Address]v1_5_1.TokenTransferFeeArgs{
							tokenA: {
								MinFeeUSDCents:            nil,                    // keep current
								MaxFeeUSDCents:            nil,                    // keep current
								DeciBps:                   pointer.To(uint16(30)), // change
								DestGasOverhead:           nil,                    // keep current
								DestBytesOverhead:         nil,                    // keep current
								AggregateRateLimitEnabled: nil,                    // keep current
							},
						},
					},
				},
			},
		},
	))
	require.NoError(t, err, "MCMS execution should succeed")

	// Verify tokenA was updated
	cfgA, err := ReadTokenTransferFeeConfig(t, e, src, dst, tokenA)
	require.NoError(t, err)
	require.True(t, cfgA.IsEnabled)
	require.Equal(t, uint16(30), cfgA.DeciBps)

	// Verify tokenB was reset to default
	cfgB, err := ReadTokenTransferFeeConfig(t, e, src, dst, tokenB)
	require.NoError(t, err)
	require.False(t, cfgB.IsEnabled, "tokenB should be using default fee config after reset")
}

func TestSetTokenTransferFeeConfig_MultipleChains(t *testing.T) {
	// Spin up a memory env with minimal prerequisite deployment
	e, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithPrerequisiteDeploymentOnly(&changeset.V1_5DeploymentConfig{
			// NOTE: this property needs to be defined otherwise we will encounter an
			// error. For now it's set to a value that was found in another test case
			PriceRegStalenessThreshold: SetTokenTransferFeePriceRegStalenessThreshold,
		}),
	)

	// Ensure the chains exist
	allChains := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.Len(t, allChains, 2)
	src := allChains[0]
	dst := allChains[1]

	// Take a snapshot of the current state
	state, err := stateview.LoadOnchainState(e.Env, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)

	// Wire up a bi-directional lane
	e.Env = v1_5.AddLanes(t, e.Env, state, []testhelpers.SourceDestPair{
		{SourceChainSelector: src, DestChainSelector: dst},
		{SourceChainSelector: dst, DestChainSelector: src},
	})

	// For src->dst
	tokenA := utils.RandomAddress() // enable with full config
	tokenB := utils.RandomAddress() // include in resets (default/no-op)

	// For dst->src
	tokenC := utils.RandomAddress() // enable with full config
	tokenD := utils.RandomAddress() // include in resets (default/no-op)

	// Prepare a single config that touches both directions (multiple chains)
	cfg := v1_5_1.SetTokenTransferFeeConfig{
		MCMS: nil, // direct execution is fine; the point is fan-out to multiple chains
		InputsByChain: map[uint64]map[uint64]v1_5_1.SetTokenTransferFeeArgs{
			src: {
				dst: {
					TokensToUseDefaultFeeConfigs: []common.Address{tokenB},
					TokenTransferFeeConfigArgs: map[common.Address]v1_5_1.TokenTransferFeeArgs{
						tokenA: {
							MinFeeUSDCents:            pointer.To(uint32(101)),
							MaxFeeUSDCents:            pointer.To(uint32(5001)),
							DeciBps:                   pointer.To(uint16(26)),
							DestGasOverhead:           pointer.To(uint32(110_000)),
							DestBytesOverhead:         pointer.To(uint32(1300)),
							AggregateRateLimitEnabled: pointer.To(true),
						},
					},
				},
			},
			dst: {
				src: {
					TokensToUseDefaultFeeConfigs: []common.Address{tokenD},
					TokenTransferFeeConfigArgs: map[common.Address]v1_5_1.TokenTransferFeeArgs{
						tokenC: {
							MinFeeUSDCents:            pointer.To(uint32(202)),
							MaxFeeUSDCents:            pointer.To(uint32(6002)),
							DeciBps:                   pointer.To(uint16(31)),
							DestGasOverhead:           pointer.To(uint32(120_000)),
							DestBytesOverhead:         pointer.To(uint32(1400)),
							AggregateRateLimitEnabled: pointer.To(false),
						},
					},
				},
			},
		},
	}

	// Apply once; should fan out to both lanes/chains
	e.Env, err = commonchangeset.Apply(t, e.Env, commonchangeset.Configure(v1_5_1.SetTokenTransferFeeConfigChangeset, cfg))
	require.NoError(t, err, "multi-chain execution should succeed in a single apply")

	// ---- Verify src -> dst lane; token config A should be set
	cfgA1, err := ReadTokenTransferFeeConfig(t, e, src, dst, tokenA)
	require.NoError(t, err)
	require.True(t, cfgA1.IsEnabled)
	require.Equal(t, uint32(101), cfgA1.MinFeeUSDCents)
	require.Equal(t, uint32(5001), cfgA1.MaxFeeUSDCents)
	require.Equal(t, uint16(26), cfgA1.DeciBps)
	require.Equal(t, uint32(110_000), cfgA1.DestGasOverhead)
	require.Equal(t, uint32(1300), cfgA1.DestBytesOverhead)
	require.True(t, cfgA1.AggregateRateLimitEnabled)

	// ---- Verify src -> dst lane; token config B should still be disabled (no-op)
	cfgB, err := ReadTokenTransferFeeConfig(t, e, src, dst, tokenB)
	require.NoError(t, err)
	require.False(t, cfgB.IsEnabled, "tokenB should be using default fee config after reset")

	// ---- Verify dst -> src lane; token config C should be set
	cfgC1, err := ReadTokenTransferFeeConfig(t, e, dst, src, tokenC)
	require.NoError(t, err)
	require.True(t, cfgC1.IsEnabled)
	require.Equal(t, uint32(202), cfgC1.MinFeeUSDCents)
	require.Equal(t, uint32(6002), cfgC1.MaxFeeUSDCents)
	require.Equal(t, uint16(31), cfgC1.DeciBps)
	require.Equal(t, uint32(120_000), cfgC1.DestGasOverhead)
	require.Equal(t, uint32(1400), cfgC1.DestBytesOverhead)
	require.False(t, cfgC1.AggregateRateLimitEnabled)

	// ---- Verify dst -> src lane; token config D should still be disabled (no-op)
	cfgD, err := ReadTokenTransferFeeConfig(t, e, dst, src, tokenD)
	require.NoError(t, err)
	require.False(t, cfgD.IsEnabled, "tokenD should be using default fee config after reset")

	// Idempotency-ish check: re-applying the exact same config should be a no-op but still succeed
	_, err = commonchangeset.Apply(t, e.Env, commonchangeset.Configure(v1_5_1.SetTokenTransferFeeConfigChangeset, cfg))
	require.NoError(t, err, "re-applying same config should be a no-op and succeed")

	// Re-read to ensure nothing regressed
	cfgA2, err := ReadTokenTransferFeeConfig(t, e, src, dst, tokenA)
	require.NoError(t, err)
	require.Equal(t, cfgA1, cfgA2)
	cfgC2, err := ReadTokenTransferFeeConfig(t, e, dst, src, tokenC)
	require.NoError(t, err)
	require.Equal(t, cfgC1, cfgC2)
}
