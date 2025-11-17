package v1_5_1_test

import (
	"math/big"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestFastTransferUpdateLaneConfigChangeset_ValidationErrors(t *testing.T) {
	e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
		selectorB: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, false)

	testCases := []struct {
		name             string
		update           v1_5_1.UpdateLaneConfig
		expectedErrorMsg string
	}{
		{
			name: "Invalid_FillerFeeBps_TooHigh",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 10001,
				FastTransferPoolFeeBps:   100,
				FillAmountMaxRequest:     big.NewInt(1000),
				FillerAllowlistEnabled:   false,
			},
			expectedErrorMsg: "fast transfer filler fee bps 10001 is greater than 10000",
		},
		{
			name: "Invalid_PoolFeeBps_TooHigh",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 100,
				FastTransferPoolFeeBps:   10001,
				FillAmountMaxRequest:     big.NewInt(1000),
				FillerAllowlistEnabled:   false,
			},
			expectedErrorMsg: "fast transfer pool fee bps 10001 is greater than 10000",
		},
		{
			name: "Invalid_FillAmount_Negative",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 100,
				FastTransferPoolFeeBps:   100,
				FillAmountMaxRequest:     big.NewInt(-1),
				FillerAllowlistEnabled:   false,
			},
			expectedErrorMsg: "fill amount max request must be a positive intege",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := v1_5_1.FastTransferUpdateLaneConfigConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.BurnMintFastTransferTokenPool,
				ContractVersion: shared.FastTransferTokenPoolVersion,
				Updates: map[uint64]map[uint64]v1_5_1.UpdateLaneConfig{
					selectorA: {selectorB: tc.update},
				},
			}

			_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferUpdateLaneConfigChangeset, config))
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErrorMsg)
		})
	}
}

func TestFastTransferFillerAllowlistChangeset_ValidationErrors(t *testing.T) {
	e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
		selectorB: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, false)

	t.Run("Empty_FillerUpdates", func(t *testing.T) {
		config := v1_5_1.FastTransferFillerAllowlistConfig{
			TokenSymbol:     testhelpers.TestTokenSymbol,
			ContractType:    shared.BurnMintFastTransferTokenPool,
			ContractVersion: shared.FastTransferTokenPoolVersion,
			Updates: map[uint64]v1_5_1.FillerAllowlistConfig{
				selectorA: {},
			},
		}

		_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferFillerAllowlistChangeset, config))
		require.Error(t, err)
		require.Contains(t, err.Error(), "at least one filler must be added or removed")
	})

	t.Run("Empty_FillerAddress_Add", func(t *testing.T) {
		config := v1_5_1.FastTransferFillerAllowlistConfig{
			TokenSymbol:     testhelpers.TestTokenSymbol,
			ContractType:    shared.BurnMintFastTransferTokenPool,
			ContractVersion: shared.FastTransferTokenPoolVersion,
			Updates: map[uint64]v1_5_1.FillerAllowlistConfig{
				selectorA: {AddFillers: []common.Address{{}}},
			},
		}

		_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferFillerAllowlistChangeset, config))
		require.Error(t, err)
		require.Contains(t, err.Error(), "filler address cannot be empty")
	})

	t.Run("Empty_FillerAddress_Remove", func(t *testing.T) {
		config := v1_5_1.FastTransferFillerAllowlistConfig{
			TokenSymbol:     testhelpers.TestTokenSymbol,
			ContractType:    shared.BurnMintFastTransferTokenPool,
			ContractVersion: shared.FastTransferTokenPoolVersion,
			Updates: map[uint64]v1_5_1.FillerAllowlistConfig{
				selectorA: {RemoveFillers: []common.Address{{}}},
			},
		}

		_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferFillerAllowlistChangeset, config))
		require.Error(t, err)
		require.Contains(t, err.Error(), "filler address cannot be empty")
	})

	t.Run("Empty_TokenSymbol", func(t *testing.T) {
		config := v1_5_1.FastTransferFillerAllowlistConfig{
			TokenSymbol:     "",
			ContractType:    shared.BurnMintFastTransferTokenPool,
			ContractVersion: shared.FastTransferTokenPoolVersion,
			Updates: map[uint64]v1_5_1.FillerAllowlistConfig{
				selectorA: {AddFillers: []common.Address{common.HexToAddress("0x1111111111111111111111111111111111111111")}},
			},
		}

		_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferFillerAllowlistChangeset, config))
		require.Error(t, err)
		require.Contains(t, err.Error(), "token symbol must be defined")
	})
}

func TestFastTransferFillerAllowlistChangeset_RemoveFillers(t *testing.T) {
	e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
		selectorB: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, false)

	type testCase struct {
		name                    string
		initialFillers          []common.Address
		fillersToRemove         []common.Address
		expectedRemainingFiller common.Address
		expectRemaining         bool
		removedFiller           common.Address
		expectRemoved           bool
	}

	testCases := []testCase{
		{
			name: "RemoveOneFillerKeepOther",
			initialFillers: []common.Address{
				common.HexToAddress("0x3333333333333333333333333333333333333331"),
				common.HexToAddress("0x3333333333333333333333333333333333333332"),
			},
			fillersToRemove:         []common.Address{common.HexToAddress("0x3333333333333333333333333333333333333331")},
			expectedRemainingFiller: common.HexToAddress("0x3333333333333333333333333333333333333332"),
			expectRemaining:         true,
			removedFiller:           common.HexToAddress("0x3333333333333333333333333333333333333331"),
			expectRemoved:           true,
		},
		{
			name: "RemoveAllFillers",
			initialFillers: []common.Address{
				common.HexToAddress("0x3333333333333333333333333333333333333333"),
				common.HexToAddress("0x3333333333333333333333333333333333333334"),
			},
			fillersToRemove: []common.Address{
				common.HexToAddress("0x3333333333333333333333333333333333333333"),
				common.HexToAddress("0x3333333333333333333333333333333333333334"),
			},
			expectedRemainingFiller: common.HexToAddress("0x3333333333333333333333333333333333333333"),
			expectRemaining:         false,
			removedFiller:           common.HexToAddress("0x3333333333333333333333333333333333333334"),
			expectRemoved:           true,
		},
		{
			name: "RemoveSingleFiller",
			initialFillers: []common.Address{
				common.HexToAddress("0x3333333333333333333333333333333333333335"),
			},
			fillersToRemove:         []common.Address{common.HexToAddress("0x3333333333333333333333333333333333333335")},
			expectedRemainingFiller: common.Address{},
			expectRemaining:         false,
			removedFiller:           common.HexToAddress("0x3333333333333333333333333333333333333335"),
			expectRemoved:           true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			addConfig := v1_5_1.FastTransferFillerAllowlistConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.BurnMintFastTransferTokenPool,
				ContractVersion: shared.FastTransferTokenPoolVersion,
				Updates: map[uint64]v1_5_1.FillerAllowlistConfig{
					selectorA: {AddFillers: tc.initialFillers},
				},
			}

			_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferFillerAllowlistChangeset, addConfig))
			require.NoError(t, err)

			state, err := stateview.LoadOnchainState(e)
			require.NoError(t, err)
			pool := state.Chains[selectorA].BurnMintFastTransferTokenPools[testhelpers.TestTokenSymbol][shared.FastTransferTokenPoolVersion]

			for _, filler := range tc.initialFillers {
				isAllowlisted, err := pool.IsAllowedFiller(nil, filler)
				require.NoError(t, err)
				require.True(t, isAllowlisted, "Expected initial filler to be allowlisted")
			}

			removeConfig := v1_5_1.FastTransferFillerAllowlistConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.BurnMintFastTransferTokenPool,
				ContractVersion: shared.FastTransferTokenPoolVersion,
				Updates: map[uint64]v1_5_1.FillerAllowlistConfig{
					selectorA: {RemoveFillers: tc.fillersToRemove},
				},
			}

			_, err = commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferFillerAllowlistChangeset, removeConfig))
			require.NoError(t, err)

			state, err = stateview.LoadOnchainState(e)
			require.NoError(t, err)
			pool = state.Chains[selectorA].BurnMintFastTransferTokenPools[testhelpers.TestTokenSymbol][shared.FastTransferTokenPoolVersion]

			if tc.expectRemoved && tc.removedFiller != (common.Address{}) {
				isAllowlisted, err := pool.IsAllowedFiller(nil, tc.removedFiller)
				require.NoError(t, err)
				require.False(t, isAllowlisted, "Expected filler to be removed from allowlist")
			}

			if tc.expectRemaining && tc.expectedRemainingFiller != (common.Address{}) {
				isAllowlisted, err := pool.IsAllowedFiller(nil, tc.expectedRemainingFiller)
				require.NoError(t, err)
				require.True(t, isAllowlisted, "Expected filler to remain in allowlist")
			}
		})
	}
}

func TestFastTransferFillerAllowlistChangeset_AddAndRemoveSimultaneously(t *testing.T) {
	e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
		selectorB: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, false)

	type testCase struct {
		name                    string
		initialFillers          []common.Address
		fillersToAdd            []common.Address
		fillersToRemove         []common.Address
		expectedAddedFiller     common.Address
		expectAdded             bool
		expectedRemovedFiller   common.Address
		expectRemoved           bool
		expectedRemainingFiller common.Address
		expectRemaining         bool
	}

	testCases := []testCase{
		{
			name: "SwapSingleFiller",
			initialFillers: []common.Address{
				common.HexToAddress("0x4444444444444444444444444444444444444441"),
			},
			fillersToAdd: []common.Address{
				common.HexToAddress("0x4444444444444444444444444444444444444442"),
			},
			fillersToRemove: []common.Address{
				common.HexToAddress("0x4444444444444444444444444444444444444441"),
			},
			expectedAddedFiller:   common.HexToAddress("0x4444444444444444444444444444444444444442"),
			expectAdded:           true,
			expectedRemovedFiller: common.HexToAddress("0x4444444444444444444444444444444444444441"),
			expectRemoved:         true,
			expectRemaining:       false,
		},
		{
			name: "AddNewKeepExisting",
			initialFillers: []common.Address{
				common.HexToAddress("0x4444444444444444444444444444444444444443"),
			},
			fillersToAdd: []common.Address{
				common.HexToAddress("0x4444444444444444444444444444444444444444"),
				common.HexToAddress("0x4444444444444444444444444444444444444445"),
			},
			fillersToRemove:         []common.Address{},
			expectedAddedFiller:     common.HexToAddress("0x4444444444444444444444444444444444444444"),
			expectAdded:             true,
			expectedRemovedFiller:   common.Address{},
			expectRemoved:           false,
			expectedRemainingFiller: common.HexToAddress("0x4444444444444444444444444444444444444443"),
			expectRemaining:         true,
		},
		{
			name: "ComplexSwapMultipleFillers",
			initialFillers: []common.Address{
				common.HexToAddress("0x4444444444444444444444444444444444444446"),
				common.HexToAddress("0x4444444444444444444444444444444444444447"),
				common.HexToAddress("0x4444444444444444444444444444444444444448"),
			},
			fillersToAdd: []common.Address{
				common.HexToAddress("0x4444444444444444444444444444444444444449"),
				common.HexToAddress("0x444444444444444444444444444444444444444a"),
			},
			fillersToRemove: []common.Address{
				common.HexToAddress("0x4444444444444444444444444444444444444446"),
				common.HexToAddress("0x4444444444444444444444444444444444444447"),
			},
			expectedAddedFiller:     common.HexToAddress("0x4444444444444444444444444444444444444449"),
			expectAdded:             true,
			expectedRemovedFiller:   common.HexToAddress("0x4444444444444444444444444444444444444446"),
			expectRemoved:           true,
			expectedRemainingFiller: common.HexToAddress("0x4444444444444444444444444444444444444448"),
			expectRemaining:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.initialFillers) > 0 {
				addConfig := v1_5_1.FastTransferFillerAllowlistConfig{
					TokenSymbol:     testhelpers.TestTokenSymbol,
					ContractType:    shared.BurnMintFastTransferTokenPool,
					ContractVersion: shared.FastTransferTokenPoolVersion,
					Updates: map[uint64]v1_5_1.FillerAllowlistConfig{
						selectorA: {AddFillers: tc.initialFillers},
					},
				}

				_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferFillerAllowlistChangeset, addConfig))
				require.NoError(t, err)

				state, err := stateview.LoadOnchainState(e)
				require.NoError(t, err)
				pool := state.Chains[selectorA].BurnMintFastTransferTokenPools[testhelpers.TestTokenSymbol][shared.FastTransferTokenPoolVersion]

				for _, filler := range tc.initialFillers {
					isAllowlisted, err := pool.IsAllowedFiller(nil, filler)
					require.NoError(t, err)
					require.True(t, isAllowlisted, "Expected initial filler to be allowlisted")
				}
			}

			updateConfig := v1_5_1.FastTransferFillerAllowlistConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.BurnMintFastTransferTokenPool,
				ContractVersion: shared.FastTransferTokenPoolVersion,
				Updates: map[uint64]v1_5_1.FillerAllowlistConfig{
					selectorA: {
						AddFillers:    tc.fillersToAdd,
						RemoveFillers: tc.fillersToRemove,
					},
				},
			}

			_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferFillerAllowlistChangeset, updateConfig))
			require.NoError(t, err)

			state, err := stateview.LoadOnchainState(e)
			require.NoError(t, err)
			pool := state.Chains[selectorA].BurnMintFastTransferTokenPools[testhelpers.TestTokenSymbol][shared.FastTransferTokenPoolVersion]

			if tc.expectAdded && tc.expectedAddedFiller != (common.Address{}) {
				isAllowlisted, err := pool.IsAllowedFiller(nil, tc.expectedAddedFiller)
				require.NoError(t, err)
				require.True(t, isAllowlisted, "Expected new filler to be added")
			}

			if tc.expectRemoved && tc.expectedRemovedFiller != (common.Address{}) {
				isAllowlisted, err := pool.IsAllowedFiller(nil, tc.expectedRemovedFiller)
				require.NoError(t, err)
				require.False(t, isAllowlisted, "Expected existing filler to be removed")
			}

			if tc.expectRemaining && tc.expectedRemainingFiller != (common.Address{}) {
				isAllowlisted, err := pool.IsAllowedFiller(nil, tc.expectedRemainingFiller)
				require.NoError(t, err)
				require.True(t, isAllowlisted, "Expected filler to remain in allowlist")
			}
		})
	}
}

type testCase struct {
	name            string
	mcmsEnabled     bool
	contractType    cldf.ContractType
	contractVersion semver.Version
}

var testCases = []testCase{
	{
		name:            "BurnMintFastTransferTokenPool",
		mcmsEnabled:     false,
		contractType:    shared.BurnMintFastTransferTokenPool,
		contractVersion: shared.FastTransferTokenPoolVersion,
	},
	{
		name:            "BurnMintFastTransferTokenPool with MCMS",
		mcmsEnabled:     true,
		contractType:    shared.BurnMintFastTransferTokenPool,
		contractVersion: shared.FastTransferTokenPoolVersion,
	},
	{
		name:            "BurnMintWithExternalMintFastTransferTokenPool",
		mcmsEnabled:     false,
		contractType:    shared.BurnMintWithExternalMinterFastTransferTokenPool,
		contractVersion: shared.BurnMintWithExternalMinterFastTransferTokenPoolVersion,
	},
	{
		name:            "BurnMintWithExternalMintFastTransferTokenPool with MCMS",
		mcmsEnabled:     true,
		contractType:    shared.BurnMintWithExternalMinterFastTransferTokenPool,
		contractVersion: shared.BurnMintWithExternalMinterFastTransferTokenPoolVersion,
	},
	{
		name:            "HybridWithExternalMintFastTransferTokenPool",
		mcmsEnabled:     false,
		contractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
		contractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
	},
	{
		name:            "HybridWithExternalMintFastTransferTokenPool with MCMS",
		mcmsEnabled:     true,
		contractType:    shared.HybridWithExternalMinterFastTransferTokenPool,
		contractVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
	},
}

func TestFastTransferUpdateLaneConfigChangeset_WithMCMS(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), tc.mcmsEnabled)

			externalMinterA := common.Address{}
			externalMinterB := common.Address{}
			if tc.contractType == shared.BurnMintWithExternalMinterFastTransferTokenPool || tc.contractType == shared.HybridWithExternalMinterFastTransferTokenPool {
				externalMinterA, _ = testhelpers.DeployTokenGovernor(t, e, selectorA, tokens[selectorA].Address)
				externalMinterB, _ = testhelpers.DeployTokenGovernor(t, e, selectorB, tokens[selectorB].Address)
			}

			e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
				selectorA: {
					Type:               tc.contractType,
					TokenAddress:       tokens[selectorA].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
					ExternalMinter:     externalMinterA,
				},
				selectorB: {
					Type:               tc.contractType,
					TokenAddress:       tokens[selectorB].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
					ExternalMinter:     externalMinterB,
				},
			}, tc.mcmsEnabled)

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}

			update := v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 250,
				FastTransferPoolFeeBps:   150,
				FillAmountMaxRequest:     big.NewInt(2000),
				FillerAllowlistEnabled:   true,
				SkipAllowlistValidation:  true,
			}
			config := v1_5_1.FastTransferUpdateLaneConfigConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    tc.contractType,
				ContractVersion: tc.contractVersion,
				MCMS:            mcmsConfig,
				Updates: map[uint64]map[uint64]v1_5_1.UpdateLaneConfig{
					selectorA: {selectorB: update},
					selectorB: {selectorA: update},
				},
			}

			_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferUpdateLaneConfigChangeset, config))
			require.NoError(t, err)

			pool, err := bindings.GetFastTransferTokenPoolContract(e, testhelpers.TestTokenSymbol, tc.contractType, tc.contractVersion, selectorA)
			require.NoError(t, err)

			result, _, err := pool.GetDestChainConfig(nil, selectorB)
			require.NoError(t, err)
			require.Equal(t, update.FastTransferFillerFeeBps, result.FastTransferFillerFeeBps)
			require.Equal(t, update.FastTransferPoolFeeBps, result.FastTransferPoolFeeBps)
			require.Equal(t, update.FillAmountMaxRequest, result.MaxFillAmountPerRequest)
			require.True(t, result.FillerAllowlistEnabled, "Expected filler allowlist to be enabled")
		})
	}
}

func TestFastTransferFillerAllowlistChangeset_WithMCMS(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), tc.mcmsEnabled)

			externalMinterA := common.Address{}
			externalMinterB := common.Address{}
			if tc.contractType == shared.BurnMintWithExternalMinterFastTransferTokenPool || tc.contractType == shared.HybridWithExternalMinterFastTransferTokenPool {
				externalMinterA, _ = testhelpers.DeployTokenGovernor(t, e, selectorA, tokens[selectorA].Address)
				externalMinterB, _ = testhelpers.DeployTokenGovernor(t, e, selectorB, tokens[selectorB].Address)
			}

			e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
				selectorA: {
					Type:               tc.contractType,
					TokenAddress:       tokens[selectorA].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
					ExternalMinter:     externalMinterA,
				},
				selectorB: {
					Type:               tc.contractType,
					TokenAddress:       tokens[selectorB].Address,
					LocalTokenDecimals: testhelpers.LocalTokenDecimals,
					ExternalMinter:     externalMinterB,
				},
			}, tc.mcmsEnabled)

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}

			config := v1_5_1.FastTransferFillerAllowlistConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    tc.contractType,
				ContractVersion: tc.contractVersion,
				MCMS:            mcmsConfig,
				Updates: map[uint64]v1_5_1.FillerAllowlistConfig{
					selectorA: {AddFillers: []common.Address{common.HexToAddress("0x5555555555555555555555555555555555555551")}},
					selectorB: {AddFillers: []common.Address{common.HexToAddress("0x5555555555555555555555555555555555555552")}},
				},
			}
			_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferFillerAllowlistChangeset, config))
			require.NoError(t, err)

			pool, err := bindings.GetFastTransferTokenPoolContract(e, testhelpers.TestTokenSymbol, tc.contractType, tc.contractVersion, selectorA)
			require.NoError(t, err)
			destinationPool, err := bindings.GetFastTransferTokenPoolContract(e, testhelpers.TestTokenSymbol, tc.contractType, tc.contractVersion, selectorB)
			require.NoError(t, err)

			isFillerAllowlisted, err := pool.IsAllowedFiller(nil, common.HexToAddress("0x5555555555555555555555555555555555555551"))
			require.NoError(t, err)
			require.True(t, isFillerAllowlisted, "Expected filler to be allowlisted")

			isFillerAllowlisted, err = destinationPool.IsAllowedFiller(nil, common.HexToAddress("0x5555555555555555555555555555555555555552"))
			require.NoError(t, err)
			require.True(t, isFillerAllowlisted, "Expected filler to be allowlisted in destination pool")
		})
	}
}

func TestFastTransferUpdateLaneConfigChangeset_EdgeCases(t *testing.T) {
	e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
		selectorB: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, false)

	type testCase struct {
		name                   string
		update                 v1_5_1.UpdateLaneConfig
		bidirectional          bool
		expectError            bool
		expectedErrorMsg       string
		validateSpecificFields bool
	}

	testCases := []testCase{
		{
			name: "Valid_BoundaryValues_MaxFees_Filler",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 9999,
				FastTransferPoolFeeBps:   0,
				FillAmountMaxRequest:     big.NewInt(1),
				FillerAllowlistEnabled:   false,
			},
			bidirectional:          false,
			expectError:            false,
			validateSpecificFields: true,
		},
		{
			name: "Valid_BoundaryValues_MaxFees_Pool",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 0,
				FastTransferPoolFeeBps:   9999,
				FillAmountMaxRequest:     big.NewInt(1),
				FillerAllowlistEnabled:   false,
			},
			bidirectional:          false,
			expectError:            false,
			validateSpecificFields: true,
		},
		{
			name: "Valid_BoundaryValues_MaxFees",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 9999,
				FastTransferPoolFeeBps:   9999,
				FillAmountMaxRequest:     big.NewInt(1),
				FillerAllowlistEnabled:   false,
			},
			bidirectional:          false,
			expectError:            true,
			validateSpecificFields: true,
		},
		{
			name: "Valid_BoundaryValues_MinFees",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 0,
				FastTransferPoolFeeBps:   0,
				FillAmountMaxRequest:     big.NewInt(1),
				FillerAllowlistEnabled:   false,
			},
			bidirectional:          true,
			expectError:            false,
			validateSpecificFields: true,
		},
		{
			name: "Valid_LargeAmount",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 5000,
				FastTransferPoolFeeBps:   3000,
				FillAmountMaxRequest:     new(big.Int).Mul(big.NewInt(1000000), big.NewInt(1e18)),
				FillerAllowlistEnabled:   false,
			},
			bidirectional:          true,
			expectError:            false,
			validateSpecificFields: true,
		},
		{
			name: "Single_ChainUpdate_OnlyAtoB",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 500,
				FastTransferPoolFeeBps:   300,
				FillAmountMaxRequest:     big.NewInt(3000),
				FillerAllowlistEnabled:   false,
			},
			bidirectional:          false,
			expectError:            false,
			validateSpecificFields: true,
		},
		{
			name: "AllowlistEnabled_NoValidation",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 200,
				FastTransferPoolFeeBps:   150,
				FillAmountMaxRequest:     big.NewInt(5000),
				FillerAllowlistEnabled:   true,
				SkipAllowlistValidation:  true,
			},
			bidirectional:          true,
			expectError:            false,
			validateSpecificFields: true,
		},
		{
			name: "ZeroAmount_Valid",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 100,
				FastTransferPoolFeeBps:   100,
				FillAmountMaxRequest:     big.NewInt(0),
				FillerAllowlistEnabled:   false,
			},
			bidirectional:          true,
			expectError:            true,
			validateSpecificFields: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updates := make(map[uint64]map[uint64]v1_5_1.UpdateLaneConfig)
			updates[selectorA] = map[uint64]v1_5_1.UpdateLaneConfig{selectorB: tc.update}
			if tc.bidirectional {
				updates[selectorB] = map[uint64]v1_5_1.UpdateLaneConfig{selectorA: tc.update}
			}

			config := v1_5_1.FastTransferUpdateLaneConfigConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.BurnMintFastTransferTokenPool,
				ContractVersion: shared.FastTransferTokenPoolVersion,
				Updates:         updates,
			}

			_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferUpdateLaneConfigChangeset, config))

			if tc.expectError {
				require.Error(t, err)
				if tc.expectedErrorMsg != "" {
					require.Contains(t, err.Error(), tc.expectedErrorMsg)
				}
				return
			}

			require.NoError(t, err)

			if tc.validateSpecificFields {
				state, err := stateview.LoadOnchainState(e)
				require.NoError(t, err)

				poolA := state.Chains[selectorA].BurnMintFastTransferTokenPools[testhelpers.TestTokenSymbol][shared.FastTransferTokenPoolVersion]

				resultAB, _, err := poolA.GetDestChainConfig(nil, selectorB)
				require.NoError(t, err)
				require.Equal(t, tc.update.FastTransferFillerFeeBps, resultAB.FastTransferFillerFeeBps)
				require.Equal(t, tc.update.FastTransferPoolFeeBps, resultAB.FastTransferPoolFeeBps)
				require.Equal(t, tc.update.FillAmountMaxRequest, resultAB.MaxFillAmountPerRequest)
				require.Equal(t, tc.update.FillerAllowlistEnabled, resultAB.FillerAllowlistEnabled)

				if tc.bidirectional {
					poolB := state.Chains[selectorB].BurnMintFastTransferTokenPools[testhelpers.TestTokenSymbol][shared.FastTransferTokenPoolVersion]
					resultBA, _, err := poolB.GetDestChainConfig(nil, selectorA)
					require.NoError(t, err)
					require.Equal(t, tc.update.FastTransferFillerFeeBps, resultBA.FastTransferFillerFeeBps)
					require.Equal(t, tc.update.FastTransferPoolFeeBps, resultBA.FastTransferPoolFeeBps)
					require.Equal(t, tc.update.FillAmountMaxRequest, resultBA.MaxFillAmountPerRequest)
					require.Equal(t, tc.update.FillerAllowlistEnabled, resultBA.FillerAllowlistEnabled)
				}
			}
		})
	}
}

func TestFastTransferFillerAllowlistChangeset_DuplicateFillerValidation(t *testing.T) {
	e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
		selectorB: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, false)

	fillerAddr := common.HexToAddress("0x6666666666666666666666666666666666666661")

	addConfig := v1_5_1.FastTransferFillerAllowlistConfig{
		TokenSymbol:     testhelpers.TestTokenSymbol,
		ContractType:    shared.BurnMintFastTransferTokenPool,
		ContractVersion: shared.FastTransferTokenPoolVersion,
		Updates: map[uint64]v1_5_1.FillerAllowlistConfig{
			selectorA: {AddFillers: []common.Address{fillerAddr}},
		},
	}

	_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferFillerAllowlistChangeset, addConfig))
	require.NoError(t, err)

	duplicateConfig := v1_5_1.FastTransferFillerAllowlistConfig{
		TokenSymbol:     testhelpers.TestTokenSymbol,
		ContractType:    shared.BurnMintFastTransferTokenPool,
		ContractVersion: shared.FastTransferTokenPoolVersion,
		Updates: map[uint64]v1_5_1.FillerAllowlistConfig{
			selectorA: {AddFillers: []common.Address{fillerAddr}},
		},
	}

	_, err = commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferFillerAllowlistChangeset, duplicateConfig))
	require.Error(t, err)
	require.Contains(t, err.Error(), "is already in the allowlist")

	nonExistentFiller := common.HexToAddress("0x7777777777777777777777777777777777777771")
	removeConfig := v1_5_1.FastTransferFillerAllowlistConfig{
		TokenSymbol:     testhelpers.TestTokenSymbol,
		ContractType:    shared.BurnMintFastTransferTokenPool,
		ContractVersion: shared.FastTransferTokenPoolVersion,
		Updates: map[uint64]v1_5_1.FillerAllowlistConfig{
			selectorA: {RemoveFillers: []common.Address{nonExistentFiller}},
		},
	}

	_, err = commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferFillerAllowlistChangeset, removeConfig))
	require.Error(t, err)
	require.Contains(t, err.Error(), "is not in the allowlist")
}

func TestFastTransferUpdateLaneConfigChangeset_SettlementOverheadGasAndCustomExtraArgs(t *testing.T) {
	e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), false)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
		selectorB: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, false)

	testCases := []struct {
		name                        string
		update                      v1_5_1.UpdateLaneConfig
		expectedSettlementGas       uint32
		expectedCustomExtraArgsSize int
	}{
		{
			name: "Default_Values_NilFields",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 100,
				FastTransferPoolFeeBps:   50,
				FillAmountMaxRequest:     big.NewInt(1000),
				FillerAllowlistEnabled:   false,
				SettlementOverheadGas:    nil,
				CustomExtraArgs:          nil,
			},
			expectedSettlementGas:       0,
			expectedCustomExtraArgsSize: 0,
		},
		{
			name: "Custom_SettlementOverheadGas",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 100,
				FastTransferPoolFeeBps:   50,
				FillAmountMaxRequest:     big.NewInt(1000),
				FillerAllowlistEnabled:   false,
				SettlementOverheadGas:    &[]uint32{50000}[0],
				CustomExtraArgs:          []byte{},
			},
			expectedSettlementGas:       50000,
			expectedCustomExtraArgsSize: 0,
		},
		{
			name: "Custom_ExtraArgs",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 100,
				FastTransferPoolFeeBps:   50,
				FillAmountMaxRequest:     big.NewInt(1000),
				FillerAllowlistEnabled:   false,
				SettlementOverheadGas:    nil,
				CustomExtraArgs:          []byte{0x01, 0x02, 0x03},
			},
			expectedSettlementGas:       0,
			expectedCustomExtraArgsSize: 3,
		},
		{
			name: "Both_Custom_Fields",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 200,
				FastTransferPoolFeeBps:   100,
				FillAmountMaxRequest:     big.NewInt(2000),
				FillerAllowlistEnabled:   false,
				SettlementOverheadGas:    &[]uint32{100000}[0],
				CustomExtraArgs:          []byte{0x04, 0x05, 0x06, 0x07, 0x08},
			},
			expectedSettlementGas:       100000,
			expectedCustomExtraArgsSize: 5,
		},
		{
			name: "Large_SettlementGas_And_ExtraArgs",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 150,
				FastTransferPoolFeeBps:   75,
				FillAmountMaxRequest:     big.NewInt(1500),
				FillerAllowlistEnabled:   false,
				SettlementOverheadGas:    &[]uint32{2000000}[0],
				CustomExtraArgs:          make([]byte, 2048),
			},
			expectedSettlementGas:       2000000,
			expectedCustomExtraArgsSize: 2048,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := v1_5_1.FastTransferUpdateLaneConfigConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.BurnMintFastTransferTokenPool,
				ContractVersion: shared.FastTransferTokenPoolVersion,
				Updates: map[uint64]map[uint64]v1_5_1.UpdateLaneConfig{
					selectorA: {selectorB: tc.update},
				},
			}

			_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferUpdateLaneConfigChangeset, config))
			require.NoError(t, err)

			pool, err := bindings.GetFastTransferTokenPoolContract(e, testhelpers.TestTokenSymbol, shared.BurnMintFastTransferTokenPool, shared.FastTransferTokenPoolVersion, selectorA)
			require.NoError(t, err)

			result, _, err := pool.GetDestChainConfig(nil, selectorB)
			require.NoError(t, err)
			require.Equal(t, tc.update.FastTransferFillerFeeBps, result.FastTransferFillerFeeBps)
			require.Equal(t, tc.update.FastTransferPoolFeeBps, result.FastTransferPoolFeeBps)
			require.Equal(t, tc.update.FillAmountMaxRequest, result.MaxFillAmountPerRequest)
			require.Equal(t, tc.update.FillerAllowlistEnabled, result.FillerAllowlistEnabled)
			require.Equal(t, tc.expectedSettlementGas, result.SettlementOverheadGas)
			require.Len(t, result.CustomExtraArgs, tc.expectedCustomExtraArgsSize)
			if tc.expectedCustomExtraArgsSize > 0 {
				require.Equal(t, tc.update.CustomExtraArgs, result.CustomExtraArgs)
			}
		})
	}
}

func TestFastTransferUpdateLaneConfigChangeset_DestinationPoolTypeAndVersion(t *testing.T) {
	e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), false)

	// Deploy different types of pools on different chains
	externalMinterA, _ := testhelpers.DeployTokenGovernor(t, e, selectorA, tokens[selectorA].Address)
	externalMinterB, _ := testhelpers.DeployTokenGovernor(t, e, selectorB, tokens[selectorB].Address)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.HybridWithExternalMinterFastTransferTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			ExternalMinter:     externalMinterA,
		},
		selectorB: {
			Type:               shared.BurnMintWithExternalMinterFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			ExternalMinter:     externalMinterB,
		},
	}, false)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorB: {
			Type:               shared.HybridWithExternalMinterFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			ExternalMinter:     externalMinterB,
		},
	}, false)

	testCases := []struct {
		name                       string
		sourceChain                uint64
		destChain                  uint64
		sourceType                 cldf.ContractType
		sourceVersion              semver.Version
		destinationContractType    *cldf.ContractType
		destinationContractVersion *semver.Version
		expectError                bool
		expectedErrorMsg           string
	}{
		{
			name:          "Valid_SameTypeAndVersion",
			sourceChain:   selectorA,
			destChain:     selectorB,
			sourceType:    shared.HybridWithExternalMinterFastTransferTokenPool,
			sourceVersion: shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
			// No destination type/version specified - should use default
			expectError: false,
		},
		{
			name:                       "Valid_DifferentType",
			sourceChain:                selectorA,
			destChain:                  selectorB,
			sourceType:                 shared.HybridWithExternalMinterFastTransferTokenPool,
			sourceVersion:              shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
			destinationContractType:    &shared.BurnMintWithExternalMinterFastTransferTokenPool,
			destinationContractVersion: &shared.BurnMintWithExternalMinterFastTransferTokenPoolVersion,
			expectError:                false,
		},
		{
			name:                       "Invalid_NonExistentDestinationType",
			sourceChain:                selectorB,
			destChain:                  selectorA,
			sourceType:                 shared.HybridWithExternalMinterFastTransferTokenPool,
			sourceVersion:              shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
			destinationContractType:    &shared.BurnMintWithExternalMinterFastTransferTokenPool, // This type doesn't exist on selectorA
			destinationContractVersion: &shared.BurnMintWithExternalMinterFastTransferTokenPoolVersion,
			expectError:                true,
			expectedErrorMsg:           "destination pool validation failed",
		},
		{
			name:                    "Invalid_NonExistentDestinationChain",
			sourceChain:             selectorA,
			destChain:               99999, // Non-existent chain
			sourceType:              shared.HybridWithExternalMinterFastTransferTokenPool,
			sourceVersion:           shared.HybridWithExternalMinterFastTransferTokenPoolVersion,
			destinationContractType: &shared.BurnMintWithExternalMinterFastTransferTokenPool,
			expectError:             true,
			expectedErrorMsg:        "unknown chain selector 99999",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			update := v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps:   100,
				FastTransferPoolFeeBps:     50,
				FillAmountMaxRequest:       big.NewInt(1000),
				FillerAllowlistEnabled:     false,
				DestinationContractType:    tc.destinationContractType,
				DestinationContractVersion: tc.destinationContractVersion,
			}

			config := v1_5_1.FastTransferUpdateLaneConfigConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    tc.sourceType,
				ContractVersion: tc.sourceVersion,
				Updates: map[uint64]map[uint64]v1_5_1.UpdateLaneConfig{
					tc.sourceChain: {tc.destChain: update},
				},
			}

			_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferUpdateLaneConfigChangeset, config))
			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrorMsg)
			} else {
				require.NoError(t, err)

				// Verify the configuration was applied correctly
				pool, err := bindings.GetFastTransferTokenPoolContract(e, testhelpers.TestTokenSymbol, tc.sourceType, tc.sourceVersion, tc.sourceChain)
				require.NoError(t, err)

				result, _, err := pool.GetDestChainConfig(nil, tc.destChain)
				require.NoError(t, err)
				require.Equal(t, update.FastTransferFillerFeeBps, result.FastTransferFillerFeeBps)
				require.Equal(t, update.FastTransferPoolFeeBps, result.FastTransferPoolFeeBps)
				require.Equal(t, update.FillAmountMaxRequest, result.MaxFillAmountPerRequest)

				// Verify the destination pool address is correctly set
				expectedDestType := tc.sourceType
				expectedDestVersion := tc.sourceVersion
				if tc.destinationContractType != nil {
					expectedDestType = *tc.destinationContractType
				}
				if tc.destinationContractVersion != nil {
					expectedDestVersion = *tc.destinationContractVersion
				}

				expectedDestPool, err := bindings.GetFastTransferTokenPoolContract(e, testhelpers.TestTokenSymbol, expectedDestType, expectedDestVersion, tc.destChain)
				require.NoError(t, err)

				expectedDestPoolPadded := common.LeftPadBytes(expectedDestPool.Address().Bytes(), 32)
				require.Equal(t, expectedDestPoolPadded, result.DestinationPool)
			}
		})
	}
}

func TestFastTransferUpdateLaneConfigChangeset_ValidationErrors_DestinationFields(t *testing.T) {
	e, selectorA, selectorB, tokens := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), false)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
		selectorB: {
			Type:               shared.BurnMintFastTransferTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, false)

	testCases := []struct {
		name             string
		update           v1_5_1.UpdateLaneConfig
		expectedErrorMsg string
	}{
		{
			name: "Invalid_DestinationPoolType_NonExistent",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps:   100,
				FastTransferPoolFeeBps:     50,
				FillAmountMaxRequest:       big.NewInt(1000),
				FillerAllowlistEnabled:     false,
				DestinationContractType:    &shared.BurnMintWithExternalMinterFastTransferTokenPool, // This type doesn't exist on selectorB
				DestinationContractVersion: &shared.BurnMintWithExternalMinterFastTransferTokenPoolVersion,
			},
			expectedErrorMsg: "destination pool validation failed",
		},
		{
			name: "Invalid_DestinationPoolVersion_NonExistent",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps:   100,
				FastTransferPoolFeeBps:     50,
				FillAmountMaxRequest:       big.NewInt(1000),
				FillerAllowlistEnabled:     false,
				DestinationContractVersion: func() *semver.Version { v, _ := semver.NewVersion("9.9.9"); return v }(), // Non-existent version
			},
			expectedErrorMsg: "destination pool validation failed",
		},
		{
			name: "Valid_OnlyDestinationContractType",
			update: v1_5_1.UpdateLaneConfig{
				FastTransferFillerFeeBps: 100,
				FastTransferPoolFeeBps:   50,
				FillAmountMaxRequest:     big.NewInt(1000),
				FillerAllowlistEnabled:   false,
				DestinationContractType:  &shared.BurnMintFastTransferTokenPool, // Same type as deployed on selectorB
				// Version will default to root config version
			},
			expectedErrorMsg: "", // Should not error
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := v1_5_1.FastTransferUpdateLaneConfigConfig{
				TokenSymbol:     testhelpers.TestTokenSymbol,
				ContractType:    shared.BurnMintFastTransferTokenPool,
				ContractVersion: shared.FastTransferTokenPoolVersion,
				Updates: map[uint64]map[uint64]v1_5_1.UpdateLaneConfig{
					selectorA: {selectorB: tc.update},
				},
			}

			_, err := commonchangeset.Apply(t, e, commonchangeset.Configure(v1_5_1.FastTransferUpdateLaneConfigChangeset, config))
			if tc.expectedErrorMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
