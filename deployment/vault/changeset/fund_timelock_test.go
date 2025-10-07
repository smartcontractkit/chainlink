package changeset

import (
	"math/big"
	"testing"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

var (
	TwoETH  = big.NewInt(0).Mul(OneETH, big.NewInt(2))
	FiveETH = big.NewInt(0).Mul(OneETH, big.NewInt(5))
)

func TestFundTimelockValidation(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
	)
	require.NoError(t, err)

	tests := []struct {
		name      string
		config    types.FundTimelockConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "empty funding config",
			config: types.FundTimelockConfig{
				FundingByChain: map[uint64]*big.Int{},
			},
			wantError: true,
			errorMsg:  "funding_by_chain must not be empty",
		},
		{
			name: "zero amount funding",
			config: types.FundTimelockConfig{
				FundingByChain: map[uint64]*big.Int{
					selector: big.NewInt(0),
				},
			},
			wantError: true,
			errorMsg:  "funding amount for chain",
		},
		{
			name: "negative amount funding",
			config: types.FundTimelockConfig{
				FundingByChain: map[uint64]*big.Int{
					selector: big.NewInt(-1),
				},
			},
			wantError: true,
			errorMsg:  "funding amount for chain",
		},
		{
			name: "invalid chain selector",
			config: types.FundTimelockConfig{
				FundingByChain: map[uint64]*big.Int{
					999999: OneETH,
				},
			},
			wantError: true,
			errorMsg:  "invalid chain selector",
		},
		{
			name: "valid funding config",
			config: types.FundTimelockConfig{
				FundingByChain: map[uint64]*big.Int{
					selector: OneETH,
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateFundTimelockConfig(t.Context(), *env, tt.config)

			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.ErrorContains(t, err, tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetTimelockBalances(t *testing.T) {
	t.Parallel()

	selectors := []uint64{chainselectors.TEST_90000001.Selector, chainselectors.TEST_90000002.Selector}
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, selectors),
	))
	require.NoError(t, err)

	setupMCMSInfrastructure(t, rt, selectors)

	t.Run("get balances for existing timelocks", func(t *testing.T) {
		balances, err := GetTimelockBalances(rt.Environment(), selectors)
		require.NoError(t, err)
		require.Len(t, balances, len(selectors))

		for _, chainSel := range selectors {
			balance, exists := balances[chainSel]
			require.True(t, exists)
			require.NotNil(t, balance.Balance)
			require.NotEmpty(t, balance.TimelockAddr)
		}
	})

	t.Run("get balances for non existent chain", func(t *testing.T) {
		_, err := GetTimelockBalances(rt.Environment(), []uint64{999999})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain 999999 not found")
	})

	t.Run("get balances with no timelock deployed", func(t *testing.T) {
		envNoTimelock, err := environment.New(t.Context(),
			environment.WithEVMSimulatedN(t, 1),
		)
		require.NoError(t, err)

		testChainSels := make([]uint64, 0)
		for chainSel := range envNoTimelock.BlockChains.EVMChains() {
			testChainSels = append(testChainSels, chainSel)
		}

		balances, err := GetTimelockBalances(*envNoTimelock, testChainSels)
		require.NoError(t, err)
		require.Empty(t, balances)
	})
}

func TestCalculateFundingRequirements(t *testing.T) {
	t.Parallel()

	selectors := []uint64{chainselectors.TEST_90000001.Selector, chainselectors.TEST_90000002.Selector}
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, selectors),
	))
	require.NoError(t, err)

	setupMCMSInfrastructure(t, rt, selectors)

	t.Run("calculate requirements for multiple chains", func(t *testing.T) {
		config := types.BatchNativeTransferConfig{
			TransfersByChain: map[uint64][]types.NativeTransfer{
				selectors[0]: {
					{To: testAddr1, Amount: OneETH},
					{To: testAddr2, Amount: TenETH},
				},
				selectors[1]: {
					{To: testRecipientMultiChain, Amount: FiveETH},
				},
			},
		}

		requirements, err := CalculateFundingRequirements(rt.Environment(), config)
		require.NoError(t, err)
		require.Len(t, requirements, 2)

		req1 := requirements[selectors[0]]
		require.NotNil(t, req1)
		require.Equal(t, selectors[0], req1.ChainSelector)
		require.NotNil(t, req1.CurrentBalance)
		expectedAmount1 := big.NewInt(0).Add(OneETH, TenETH)
		require.Equal(t, expectedAmount1, req1.RequiredAmount)
		require.Equal(t, 2, req1.TransferCount)

		req2 := requirements[selectors[1]]
		require.NotNil(t, req2)
		require.Equal(t, selectors[1], req2.ChainSelector)
		require.NotNil(t, req2.CurrentBalance)
		require.Equal(t, FiveETH, req2.RequiredAmount)
		require.Equal(t, 1, req2.TransferCount)
	})

	t.Run("calculate requirements with no transfers", func(t *testing.T) {
		config := types.BatchNativeTransferConfig{
			TransfersByChain: map[uint64][]types.NativeTransfer{},
		}

		requirements, err := CalculateFundingRequirements(rt.Environment(), config)
		require.NoError(t, err)
		require.Empty(t, requirements)
	})

	t.Run("calculate requirements with single transfer", func(t *testing.T) {
		config := types.BatchNativeTransferConfig{
			TransfersByChain: map[uint64][]types.NativeTransfer{
				selectors[0]: {
					{To: testAddr1, Amount: OneETH},
				},
			},
		}

		requirements, err := CalculateFundingRequirements(rt.Environment(), config)
		require.NoError(t, err)
		require.Len(t, requirements, 1)

		req := requirements[selectors[0]]
		require.NotNil(t, req)
		require.Equal(t, OneETH, req.RequiredAmount)
		require.Equal(t, 1, req.TransferCount)
	})
}

func TestFundTimelockChangeset(t *testing.T) {
	t.Parallel()

	selector1 := chainselectors.TEST_90000001.Selector
	selector2 := chainselectors.TEST_90000002.Selector

	t.Run("single chain", func(t *testing.T) {
		t.Parallel()

		selectors := []uint64{selector1}

		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, selectors),
		))
		require.NoError(t, err)

		setupMCMSInfrastructure(t, rt, selectors)
		fundDeployerAccounts(t, rt.Environment(), selectors)

		t.Run("successful timelock funding", func(t *testing.T) {
			fundingAmount := OneETH
			config := types.FundTimelockConfig{
				FundingByChain: map[uint64]*big.Int{
					selector1: fundingAmount,
				},
			}

			balancesBefore, err := GetTimelockBalances(rt.Environment(), []uint64{selector1})
			require.NoError(t, err)
			balanceBefore := balancesBefore[selector1].Balance

			output, err := FundTimelockChangeset.Apply(rt.Environment(), config)
			require.NoError(t, err)
			require.NotNil(t, output.DataStore)

			balancesAfter, err := GetTimelockBalances(rt.Environment(), selectors)
			require.NoError(t, err)
			balanceAfter := balancesAfter[selector1].Balance

			expectedBalance := big.NewInt(0).Add(balanceBefore, fundingAmount)
			require.Equal(t, expectedBalance, balanceAfter)
		})

		t.Run("funding with invalid config fails precondition", func(t *testing.T) {
			config := types.FundTimelockConfig{
				FundingByChain: map[uint64]*big.Int{},
			}

			err := FundTimelockChangeset.VerifyPreconditions(rt.Environment(), config)
			require.Error(t, err)
			require.Contains(t, err.Error(), "funding_by_chain must not be empty")
		})
	})

	t.Run("multiple chains", func(t *testing.T) {
		t.Parallel()

		selectors := []uint64{selector1, selector2}

		rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
			environment.WithEVMSimulated(t, selectors),
		))
		require.NoError(t, err)

		setupMCMSInfrastructure(t, rt, selectors)
		fundDeployerAccounts(t, rt.Environment(), selectors)

		fundingAmount1 := OneETH
		fundingAmount2 := TwoETH
		config := types.FundTimelockConfig{
			FundingByChain: map[uint64]*big.Int{
				selector1: fundingAmount1,
				selector2: fundingAmount2,
			},
		}

		output, err := FundTimelockChangeset.Apply(rt.Environment(), config)
		require.NoError(t, err)
		require.NotNil(t, output.DataStore)
	})
}
