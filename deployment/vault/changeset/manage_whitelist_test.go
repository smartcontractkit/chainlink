package changeset

import (
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	whitelistTestAddr1 = "0x1111111111111111111111111111111111111111"
	whitelistTestAddr2 = "0x2222222222222222222222222222222222222222"
	whitelistTestAddr3 = "0x3333333333333333333333333333333333333333"
)

func TestSetWhitelistValidation(t *testing.T) {
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
		config    types.SetWhitelistConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "empty whitelist config",
			config: types.SetWhitelistConfig{
				WhitelistByChain: map[uint64][]types.WhitelistAddress{},
			},
			wantError: true,
			errorMsg:  "whitelist_by_chain must not be empty",
		},
		{
			name: "zero address in whitelist",
			config: types.SetWhitelistConfig{
				WhitelistByChain: map[uint64][]types.WhitelistAddress{
					testChainSel: {
						{
							Address:     "0x0000000000000000000000000000000000000000",
							Description: "Zero address",
							Labels:      []string{"test"},
						},
					},
				},
			},
			wantError: true,
			errorMsg:  "address cannot be zero address",
		},
		{
			name: "duplicate addresses in same chain",
			config: types.SetWhitelistConfig{
				WhitelistByChain: map[uint64][]types.WhitelistAddress{
					testChainSel: {
						{
							Address:     common.HexToAddress(whitelistTestAddr1).Hex(),
							Description: "First instance",
							Labels:      []string{"test"},
						},
						{
							Address:     common.HexToAddress(whitelistTestAddr1).Hex(),
							Description: "Duplicate instance",
							Labels:      []string{"test"},
						},
					},
				},
			},
			wantError: true,
			errorMsg:  "duplicate address",
		},
		{
			name: "invalid chain selector",
			config: types.SetWhitelistConfig{
				WhitelistByChain: map[uint64][]types.WhitelistAddress{
					999999: {
						{
							Address:     common.HexToAddress(whitelistTestAddr1).Hex(),
							Description: "Valid address",
							Labels:      []string{"test"},
						},
					},
				},
			},
			wantError: true,
			errorMsg:  "invalid chain selector",
		},
		{
			name: "valid whitelist config",
			config: types.SetWhitelistConfig{
				WhitelistByChain: map[uint64][]types.WhitelistAddress{
					testChainSel: {
						{
							Address:     common.HexToAddress(whitelistTestAddr1).Hex(),
							Description: "Test address 1",
							Labels:      []string{"team", "approved"},
						},
						{
							Address:     common.HexToAddress(whitelistTestAddr2).Hex(),
							Description: "Test address 2",
							Labels:      []string{"partner"},
						},
					},
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSetWhitelistConfig(env, tt.config)

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

func TestGetWhitelistedAddresses(t *testing.T) {
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

	t.Run("get whitelist from uninitialized datastore", func(t *testing.T) {
		envNoDS := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
			Chains: 2,
		})
		envNoDS.DataStore = nil

		_, err := GetWhitelistedAddresses(envNoDS, []uint64{chain1})
		require.Error(t, err)
		require.Contains(t, err.Error(), "datastore is nil")
	})

	t.Run("get whitelist from empty whitelist", func(t *testing.T) {
		env.DataStore = datastore.NewMemoryDataStore().Seal()

		whitelist, err := GetWhitelistedAddresses(env, []uint64{chain1})
		require.NoError(t, err)
		require.Empty(t, whitelist[chain1])
	})

	t.Run("get whitelist after setting addresses", func(t *testing.T) {
		config := types.SetWhitelistConfig{
			WhitelistByChain: map[uint64][]types.WhitelistAddress{
				chain1: {
					{
						Address:     common.HexToAddress(whitelistTestAddr1).Hex(),
						Description: "Test address 1",
						Labels:      []string{"team", "approved"},
					},
					{
						Address:     common.HexToAddress(whitelistTestAddr2).Hex(),
						Description: "Test address 2",
						Labels:      []string{"partner"},
					},
				},
				chain2: {
					{
						Address:     common.HexToAddress(whitelistTestAddr3).Hex(),
						Description: "Test address 3",
						Labels:      []string{"contractor"},
					},
				},
			},
		}

		output, err := SetWhitelistChangeset.Apply(env, config)
		require.NoError(t, err)
		env.DataStore = output.DataStore.Seal()

		whitelist, err := GetWhitelistedAddresses(env, []uint64{chain1, chain2})
		require.NoError(t, err)

		require.Len(t, whitelist[chain1], 2)
		require.Equal(t, whitelistTestAddr1, whitelist[chain1][0].Address)
		require.Equal(t, []string{"team", "approved"}, whitelist[chain1][0].Labels)
		require.Equal(t, whitelistTestAddr2, whitelist[chain1][1].Address)
		require.Equal(t, []string{"partner"}, whitelist[chain1][1].Labels)

		require.Len(t, whitelist[chain2], 1)
		require.Equal(t, whitelistTestAddr3, whitelist[chain2][0].Address)
		require.Equal(t, []string{"contractor"}, whitelist[chain2][0].Labels)
	})

	t.Run("get whitelist for specific chain only", func(t *testing.T) {
		whitelist, err := GetWhitelistedAddresses(env, []uint64{chain1})
		require.NoError(t, err)

		require.Len(t, whitelist, 1)
		require.Contains(t, whitelist, chain1)
		require.NotContains(t, whitelist, chain2)
		require.Len(t, whitelist[chain1], 2)
	})
}

func TestValidateWhitelist(t *testing.T) {
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})

	chainSelectors := make([]uint64, 0)
	for chainSel := range env.BlockChains.EVMChains() {
		chainSelectors = append(chainSelectors, chainSel)
	}
	require.Len(t, chainSelectors, 2)

	sort.Slice(chainSelectors, func(i, j int) bool {
		return chainSelectors[i] < chainSelectors[j]
	})

	chain1 := chainSelectors[0]
	chain2 := chainSelectors[1]

	whitelistConfig := types.SetWhitelistConfig{
		WhitelistByChain: map[uint64][]types.WhitelistAddress{
			chain1: {
				{
					Address:     common.HexToAddress(whitelistTestAddr1).Hex(),
					Description: "Approved address 1",
					Labels:      []string{"team"},
				},
				{
					Address:     common.HexToAddress(whitelistTestAddr2).Hex(),
					Description: "Approved address 2",
					Labels:      []string{"partner"},
				},
			},
			chain2: {
				{
					Address:     common.HexToAddress(whitelistTestAddr3).Hex(),
					Description: "Approved address 3",
					Labels:      []string{"contractor"},
				},
			},
		},
	}

	output, err := SetWhitelistChangeset.Apply(env, whitelistConfig)
	require.NoError(t, err)
	env.DataStore = output.DataStore.Seal()

	t.Run("validate transfers with all whitelisted addresses", func(t *testing.T) {
		config := types.BatchNativeTransferConfig{
			TransfersByChain: map[uint64][]types.NativeTransfer{
				chain1: {
					{To: whitelistTestAddr1, Amount: OneETH},
					{To: whitelistTestAddr2, Amount: TenETH},
				},
				chain2: {
					{To: whitelistTestAddr3, Amount: OneETH},
				},
			},
		}

		errors, err := ValidateWhitelist(env, config)
		require.NoError(t, err)
		require.Empty(t, errors)
	})

	t.Run("validate transfers with non-whitelisted addresses", func(t *testing.T) {
		nonWhitelistedAddr := "0x4444444444444444444444444444444444444444"
		config := types.BatchNativeTransferConfig{
			TransfersByChain: map[uint64][]types.NativeTransfer{
				chain1: {
					{To: whitelistTestAddr1, Amount: OneETH},
					{To: nonWhitelistedAddr, Amount: TenETH},
				},
				chain2: {
					{To: whitelistTestAddr3, Amount: OneETH},
					{To: nonWhitelistedAddr, Amount: OneETH},
				},
			},
		}

		validationErrors, err := ValidateWhitelist(env, config)
		require.NoError(t, err)
		require.Len(t, validationErrors, 2)

		errorsByChain := make(map[uint64]types.TransferValidationError)
		for _, err := range validationErrors {
			errorsByChain[err.ChainSelector] = err
		}

		require.Contains(t, errorsByChain, chain1)
		require.Equal(t, nonWhitelistedAddr, errorsByChain[chain1].Address)
		require.Contains(t, errorsByChain[chain1].Error, "address not in whitelist")

		require.Contains(t, errorsByChain, chain2)
		require.Equal(t, nonWhitelistedAddr, errorsByChain[chain2].Address)
		require.Contains(t, errorsByChain[chain2].Error, "address not in whitelist")
	})
}

func TestGetChainWhitelist(t *testing.T) {
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

	t.Run("get whitelist from empty datastore", func(t *testing.T) {
		ds := datastore.NewMemoryDataStore().Seal()
		metadata, err := GetChainWhitelist(ds, testChainSel)
		require.NoError(t, err)
		require.NotNil(t, metadata)
		require.Empty(t, metadata.Addresses)
	})

	t.Run("get whitelist after setting addresses", func(t *testing.T) {
		ds := datastore.NewMemoryDataStore()

		whitelistMetadata := types.WhitelistMetadata{
			Addresses: []types.WhitelistAddress{
				{
					Address:     common.HexToAddress(whitelistTestAddr1).Hex(),
					Description: "Test address 1",
					Labels:      []string{"team"},
				},
				{
					Address:     common.HexToAddress(whitelistTestAddr2).Hex(),
					Description: "Test address 2",
					Labels:      []string{"partner"},
				},
			},
		}

		err := ds.ChainMetadata().Upsert(datastore.ChainMetadata{
			ChainSelector: testChainSel,
			Metadata:      whitelistMetadata,
		})
		require.NoError(t, err)

		sealedDS := ds.Seal()
		metadata, err := GetChainWhitelist(sealedDS, testChainSel)
		require.NoError(t, err)
		require.NotNil(t, metadata)
		require.Len(t, metadata.Addresses, 2)

		require.Equal(t, common.HexToAddress(whitelistTestAddr1).Hex(), metadata.Addresses[0].Address)
		require.Equal(t, "Test address 1", metadata.Addresses[0].Description)
		require.Equal(t, []string{"team"}, metadata.Addresses[0].Labels)

		require.Equal(t, common.HexToAddress(whitelistTestAddr2).Hex(), metadata.Addresses[1].Address)
		require.Equal(t, "Test address 2", metadata.Addresses[1].Description)
		require.Equal(t, []string{"partner"}, metadata.Addresses[1].Labels)
	})
}

func TestSetWhitelistChangeset(t *testing.T) {
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

	t.Run("set whitelist for multiple chains", func(t *testing.T) {
		config := types.SetWhitelistConfig{
			WhitelistByChain: map[uint64][]types.WhitelistAddress{
				chain1: {
					{
						Address:     common.HexToAddress(whitelistTestAddr1).Hex(),
						Description: "Team A wallet",
						Labels:      []string{"team", "payments"},
					},
					{
						Address:     common.HexToAddress(whitelistTestAddr2).Hex(),
						Description: "Team B wallet",
						Labels:      []string{"team", "payments"},
					},
				},
				chain2: {
					{
						Address:     common.HexToAddress(whitelistTestAddr3).Hex(),
						Description: "Partner wallet",
						Labels:      []string{"partner", "contractor"},
					},
				},
			},
		}

		output, err := SetWhitelistChangeset.Apply(env, config)
		require.NoError(t, err)
		require.NotNil(t, output.DataStore)

		env.DataStore = output.DataStore.Seal()

		whitelist, err := GetWhitelistedAddresses(env, []uint64{chain1, chain2})
		require.NoError(t, err)

		require.Len(t, whitelist[chain1], 2)
		require.Len(t, whitelist[chain2], 1)

		require.Equal(t, whitelistTestAddr1, whitelist[chain1][0].Address)
		require.Equal(t, []string{"team", "payments"}, whitelist[chain1][0].Labels)

		require.Equal(t, whitelistTestAddr3, whitelist[chain2][0].Address)
		require.Equal(t, []string{"partner", "contractor"}, whitelist[chain2][0].Labels)
	})

	t.Run("update existing whitelist", func(t *testing.T) {
		updatedConfig := types.SetWhitelistConfig{
			WhitelistByChain: map[uint64][]types.WhitelistAddress{
				chain1: {
					{
						Address:     common.HexToAddress(whitelistTestAddr2).Hex(),
						Description: "Team B wallet updated",
						Labels:      []string{"team", "payments", "updated"},
					},
				},
				chain2: {
					{
						Address:     common.HexToAddress(whitelistTestAddr3).Hex(),
						Description: "Partner wallet",
						Labels:      []string{"partner", "contractor"},
					},
					{
						Address:     common.HexToAddress("0x5555555555555555555555555555555555555555").Hex(),
						Description: "New partner wallet",
						Labels:      []string{"partner", "new"},
					},
				},
			},
		}

		output, err := SetWhitelistChangeset.Apply(env, updatedConfig)
		require.NoError(t, err)
		require.NotNil(t, output.DataStore)

		env.DataStore = output.DataStore.Seal()

		whitelist, err := GetWhitelistedAddresses(env, []uint64{chain1, chain2})
		require.NoError(t, err)

		require.Len(t, whitelist[chain1], 1)
		require.Len(t, whitelist[chain2], 2)

		require.Equal(t, whitelistTestAddr2, whitelist[chain1][0].Address)
		require.Equal(t, []string{"team", "payments", "updated"}, whitelist[chain1][0].Labels)
	})

	t.Run("clear whitelist for a chain", func(t *testing.T) {
		clearConfig := types.SetWhitelistConfig{
			WhitelistByChain: map[uint64][]types.WhitelistAddress{
				chain1: {},
			},
		}

		output, err := SetWhitelistChangeset.Apply(env, clearConfig)
		require.NoError(t, err)
		require.NotNil(t, output.DataStore)

		env.DataStore = output.DataStore.Seal()

		whitelist, err := GetWhitelistedAddresses(env, []uint64{chain1})
		require.NoError(t, err)
		require.Empty(t, whitelist[chain1])
	})
}
