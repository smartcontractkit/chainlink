package changeset

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	TwoETH  = big.NewInt(0).Mul(OneETH, big.NewInt(2))
	FiveETH = big.NewInt(0).Mul(OneETH, big.NewInt(5))
)

func TestFundTimelockValidation(t *testing.T) {
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})

	chainSelectors := make([]uint64, 0)
	for chainSel := range env.BlockChains.EVMChains() {
		chainSelectors = append(chainSelectors, chainSel)
	}
	require.Len(t, chainSelectors, 1)
	testChainSel := chainSelectors[0]

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
					testChainSel: big.NewInt(0),
				},
			},
			wantError: true,
			errorMsg:  "funding amount for chain",
		},
		{
			name: "negative amount funding",
			config: types.FundTimelockConfig{
				FundingByChain: map[uint64]*big.Int{
					testChainSel: big.NewInt(-1),
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
					testChainSel: OneETH,
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFundTimelockConfig(env.GetContext(), env, tt.config)

			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetTimelockBalances(t *testing.T) {
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})

	chainSelectors := make([]uint64, 0)
	for chainSel := range env.BlockChains.EVMChains() {
		chainSelectors = append(chainSelectors, chainSel)
	}
	require.Len(t, chainSelectors, 2)
	env = setupMCMSInfrastructure(t, env, chainSelectors)

	t.Run("get balances for existing timelocks", func(t *testing.T) {
		balances, err := GetTimelockBalances(env, chainSelectors)
		require.NoError(t, err)
		require.Len(t, balances, len(chainSelectors))

		for _, chainSel := range chainSelectors {
			balance, exists := balances[chainSel]
			require.True(t, exists)
			require.NotNil(t, balance.Balance)
			require.NotEmpty(t, balance.TimelockAddr)
		}
	})

	t.Run("get balances for non existent chain", func(t *testing.T) {
		_, err := GetTimelockBalances(env, []uint64{999999})
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain 999999 not found")
	})

	t.Run("get balances with no timelock deployed", func(t *testing.T) {
		envNoTimelock := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
			Chains: 1,
		})
		envNoTimelock.DataStore = datastore.NewMemoryDataStore().Seal()

		testChainSels := make([]uint64, 0)
		for chainSel := range envNoTimelock.BlockChains.EVMChains() {
			testChainSels = append(testChainSels, chainSel)
		}

		balances, err := GetTimelockBalances(envNoTimelock, testChainSels)
		require.NoError(t, err)
		require.Empty(t, balances)
	})
}

func TestCalculateFundingRequirements(t *testing.T) {
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})

	chainSelectors := make([]uint64, 0)
	for chainSel := range env.BlockChains.EVMChains() {
		chainSelectors = append(chainSelectors, chainSel)
	}
	require.Len(t, chainSelectors, 2)

	chain1 := chainSelectors[0]
	chain2 := chainSelectors[1]

	env = setupMCMSInfrastructure(t, env, chainSelectors)

	t.Run("calculate requirements for multiple chains", func(t *testing.T) {
		config := types.BatchNativeTransferConfig{
			TransfersByChain: map[uint64][]types.NativeTransfer{
				chain1: {
					{To: testAddr1, Amount: OneETH},
					{To: testAddr2, Amount: TenETH},
				},
				chain2: {
					{To: testRecipientMultiChain, Amount: FiveETH},
				},
			},
		}

		requirements, err := CalculateFundingRequirements(env, config)
		require.NoError(t, err)
		require.Len(t, requirements, 2)

		req1 := requirements[chain1]
		require.NotNil(t, req1)
		require.Equal(t, chain1, req1.ChainSelector)
		require.NotNil(t, req1.CurrentBalance)
		expectedAmount1 := big.NewInt(0).Add(OneETH, TenETH)
		require.Equal(t, expectedAmount1, req1.RequiredAmount)
		require.Equal(t, 2, req1.TransferCount)

		req2 := requirements[chain2]
		require.NotNil(t, req2)
		require.Equal(t, chain2, req2.ChainSelector)
		require.NotNil(t, req2.CurrentBalance)
		require.Equal(t, FiveETH, req2.RequiredAmount)
		require.Equal(t, 1, req2.TransferCount)
	})

	t.Run("calculate requirements with no transfers", func(t *testing.T) {
		config := types.BatchNativeTransferConfig{
			TransfersByChain: map[uint64][]types.NativeTransfer{},
		}

		requirements, err := CalculateFundingRequirements(env, config)
		require.NoError(t, err)
		require.Empty(t, requirements)
	})

	t.Run("calculate requirements with single transfer", func(t *testing.T) {
		config := types.BatchNativeTransferConfig{
			TransfersByChain: map[uint64][]types.NativeTransfer{
				chain1: {
					{To: testAddr1, Amount: OneETH},
				},
			},
		}

		requirements, err := CalculateFundingRequirements(env, config)
		require.NoError(t, err)
		require.Len(t, requirements, 1)

		req := requirements[chain1]
		require.NotNil(t, req)
		require.Equal(t, OneETH, req.RequiredAmount)
		require.Equal(t, 1, req.TransferCount)
	})
}

func TestFundTimelockChangeset(t *testing.T) {
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})

	chainSelectors := make([]uint64, 0)
	for chainSel := range env.BlockChains.EVMChains() {
		chainSelectors = append(chainSelectors, chainSel)
	}
	require.Len(t, chainSelectors, 1)
	testChainSel := chainSelectors[0]

	env = setupMCMSInfrastructure(t, env, chainSelectors)
	fundDeployerAccounts(t, env, chainSelectors)

	t.Run("successful timelock funding", func(t *testing.T) {
		fundingAmount := OneETH
		config := types.FundTimelockConfig{
			FundingByChain: map[uint64]*big.Int{
				testChainSel: fundingAmount,
			},
		}

		balancesBefore, err := GetTimelockBalances(env, chainSelectors)
		require.NoError(t, err)
		balanceBefore := balancesBefore[testChainSel].Balance

		output, err := FundTimelockChangeset.Apply(env, config)
		require.NoError(t, err)
		require.NotNil(t, output.DataStore)

		balancesAfter, err := GetTimelockBalances(env, chainSelectors)
		require.NoError(t, err)
		balanceAfter := balancesAfter[testChainSel].Balance

		expectedBalance := big.NewInt(0).Add(balanceBefore, fundingAmount)
		require.Equal(t, expectedBalance, balanceAfter)
	})

	t.Run("funding with invalid config fails precondition", func(t *testing.T) {
		config := types.FundTimelockConfig{
			FundingByChain: map[uint64]*big.Int{},
		}

		err := FundTimelockChangeset.VerifyPreconditions(env, config)
		require.Error(t, err)
		require.Contains(t, err.Error(), "funding_by_chain must not be empty")
	})

	t.Run("funding multiple chains", func(t *testing.T) {
		env2 := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
			Chains: 2,
		})

		chainSels2 := make([]uint64, 0)
		for chainSel := range env2.BlockChains.EVMChains() {
			chainSels2 = append(chainSels2, chainSel)
		}
		require.Len(t, chainSels2, 2)

		env2 = setupMCMSInfrastructure(t, env2, chainSels2)
		fundDeployerAccounts(t, env2, chainSels2)

		fundingAmount1 := OneETH
		fundingAmount2 := TwoETH
		config := types.FundTimelockConfig{
			FundingByChain: map[uint64]*big.Int{
				chainSels2[0]: fundingAmount1,
				chainSels2[1]: fundingAmount2,
			},
		}

		output, err := FundTimelockChangeset.Apply(env2, config)
		require.NoError(t, err)
		require.NotNil(t, output.DataStore)
	})
}
