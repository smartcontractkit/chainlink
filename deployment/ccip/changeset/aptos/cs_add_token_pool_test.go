package aptos_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/smartcontractkit/quarantine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/managed_token_pool"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var testTokenTransferFeeConfig = fee_quoter.TokenTransferFeeConfig{
	MinFeeUsdCents:    50,
	MaxFeeUsdCents:    1600,
	DeciBps:           0,
	DestGasOverhead:   300_000,
	DestBytesOverhead: 100,
	IsEnabled:         true,
}

func TestAddTokenPool_Apply(t *testing.T) {
	quarantine.Flaky(t, "DX-2088")
	t.Parallel()
	// Setup environment and config with 1 Aptos chain
	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithAptosChains(1),
	)
	env := deployedEnvironment.Env

	// Get chain selectors for Aptos
	aptosChainSelectors := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyAptos))
	require.Len(t, aptosChainSelectors, 1, "Expected exactly 1 Aptos chain")
	aptosSelector := aptosChainSelectors[0]

	mockEVMPool := "0xbd10ffa3815c010d5cf7d38815a0eaabc959eb84"
	// Get EVM chain selectors
	emvSelector1 := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	emvSelector2 := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[1]

	// Configure token pool settings
	cfg := config.AddTokenPoolConfig{
		MCMSConfig: &proposalutils.TimelockConfig{
			MinDelay:     time.Duration(1) * time.Second,
			MCMSAction:   mcmstypes.TimelockActionSchedule,
			OverrideRoot: false,
		},
		ChainSelector: aptosSelector,
		PoolType:      shared.AptosManagedTokenPoolType,
		EVMRemoteConfigs: map[uint64]config.EVMRemoteConfig{
			emvSelector1: {
				TokenAddress:     common.HexToAddress("0xa"),
				TokenPoolAddress: common.HexToAddress(mockEVMPool),
				RateLimiterConfig: config.RateLimiterConfig{
					RemoteChainSelector: emvSelector1,
					OutboundIsEnabled:   false,
					OutboundCapacity:    0,
					OutboundRate:        0,
					InboundIsEnabled:    true,
					InboundCapacity:     110,
					InboundRate:         20,
				},
			},
			emvSelector2: {
				TokenAddress:     common.HexToAddress("0xa"),
				TokenPoolAddress: common.HexToAddress(mockEVMPool),
				RateLimiterConfig: config.RateLimiterConfig{
					RemoteChainSelector: emvSelector2,
					OutboundIsEnabled:   false,
					OutboundCapacity:    0,
					OutboundRate:        0,
					InboundIsEnabled:    true,
					InboundCapacity:     110,
					InboundRate:         20,
				},
			},
		},
		TokenTransferFeeByRemoteChainConfig: map[uint64]fee_quoter.TokenTransferFeeConfig{
			emvSelector2: testTokenTransferFeeConfig,
		},
		TokenParams: config.TokenParams{
			MaxSupply: big.NewInt(1000000),
			Name:      "BnMTest",
			Symbol:    "BnM",
			Decimals:  8,
			Icon:      "",
			Project:   "",
		},
	}

	// Apply the AddTokenPool changeset
	env, output, err := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(aptoscs.AddTokenPool{}, cfg),
	})
	require.NoError(t, err)

	// Load onchain state for assertions
	state, err := stateview.LoadOnchainState(env)
	require.NoError(t, err, "must load onchain state")
	require.NotNil(t, state.AptosChains[aptosSelector].AptosManagedTokenPools)

	client := env.BlockChains.AptosChains()[aptosSelector].Client
	aptosCCIPAddr := state.AptosChains[aptosSelector].CCIPAddress
	ccipContract := ccip.Bind(aptosCCIPAddr, client)
	mcmsAddress := state.AptosChains[aptosSelector].MCMSAddress
	mcmsContract := mcmsbind.Bind(mcmsAddress, client)

	expectedAdm, err := mcmsContract.MCMSRegistry().GetRegisteredOwnerAddress(nil, aptosCCIPAddr)
	require.NoError(t, err)

	for tokenAddress, pool := range state.AptosChains[aptosSelector].AptosManagedTokenPools {
		poolBind := managed_token_pool.Bind(pool, env.BlockChains.AptosChains()[aptosSelector].Client)

		remotePools, err := poolBind.ManagedTokenPool().GetRemotePools(nil, emvSelector1)
		require.NoError(t, err)
		require.NotEmpty(t, remotePools)
		hexString := fmt.Sprintf("0x%x", remotePools[0])
		assert.Equal(t, hexString, mockEVMPool)

		remotePools2, err := poolBind.ManagedTokenPool().GetRemotePools(nil, emvSelector2)
		require.NoError(t, err)
		require.NotEmpty(t, remotePools2)
		hexString = fmt.Sprintf("0x%x", remotePools2[0])
		assert.Equal(t, hexString, mockEVMPool)

		poolAdd, admin, pendingAdm, err := ccipContract.TokenAdminRegistry().GetTokenConfig(nil, tokenAddress)
		require.NoError(t, err)
		require.Equal(t, pool, poolAdd, "Expected the registered pool to match the deployed pool")
		require.Equal(t, expectedAdm, admin, "Admin should match ccipOwnerAddress")
		require.Equal(t, aptos.AccountAddress{}, pendingAdm, "Pending admin should be empty")

		ttfcfgs1, err := ccipContract.FeeQuoter().GetTokenTransferFeeConfig(nil, emvSelector1, tokenAddress)
		require.NoError(t, err)
		require.Equal(t, fee_quoter.TokenTransferFeeConfig{}, ttfcfgs1)
		ttfcfgs2, err := ccipContract.FeeQuoter().GetTokenTransferFeeConfig(nil, emvSelector2, tokenAddress)
		require.NoError(t, err)
		require.Equal(t, testTokenTransferFeeConfig, ttfcfgs2)

	}

	// The output should include MCMS proposals
	require.Len(t, output[0].MCMSTimelockProposals, 1, "Expected exactly 1 MCMS proposal")
	require.Len(t, output[0].MCMSTimelockProposals[0].Operations, 7, "Expected exactly 7 MCMS proposal operations, received:", len(output[0].MCMSTimelockProposals[0].Operations))
}
