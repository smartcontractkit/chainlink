package aptos_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/managed_token_pool"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func TestAddTokenPool_Apply(t *testing.T) {
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
	emvSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]

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
			emvSelector: {
				TokenAddress:     common.HexToAddress("0xa"),
				TokenPoolAddress: common.HexToAddress(mockEVMPool),
				RateLimiterConfig: config.RateLimiterConfig{
					RemoteChainSelector: emvSelector,
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
			emvSelector: {
				MinFeeUsdCents:    800,
				MaxFeeUsdCents:    1600,
				DeciBps:           0,
				DestGasOverhead:   300_000,
				DestBytesOverhead: 100,
				IsEnabled:         true,
			},
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

	for _, pool := range state.AptosChains[aptosSelector].AptosManagedTokenPools {
		poolBind := managed_token_pool.Bind(pool, env.BlockChains.AptosChains()[aptosSelector].Client)
		remotePools, err := poolBind.ManagedTokenPool().GetRemotePools(nil, emvSelector)
		require.NoError(t, err)
		require.NotEmpty(t, remotePools)
		hexString := fmt.Sprintf("0x%x", remotePools[0])
		assert.Equal(t, hexString, mockEVMPool)
	}

	// The output should include MCMS proposals
	require.Len(t, output[0].MCMSTimelockProposals, 1, "Expected exactly 1 MCMS proposal")
	require.Len(t, output[0].MCMSTimelockProposals[0].Operations, 8, "Expected exactly 8 MCMS proposal operations, received:", len(output[0].MCMSTimelockProposals[0].Operations))
}
