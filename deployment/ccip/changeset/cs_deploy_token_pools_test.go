package changeset_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/evm/utils"
)

func TestValidateDeployTokenPoolContractsConfig(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})

	tests := []struct {
		Msg         string
		TokenSymbol changeset.TokenSymbol
		Input       changeset.DeployTokenPoolContractsConfig
		ErrStr      string
	}{
		{
			Msg:    "Token symbol is missing",
			Input:  changeset.DeployTokenPoolContractsConfig{},
			ErrStr: "token symbol must be defined",
		},
		{
			Msg: "Chain selector is not valid",
			Input: changeset.DeployTokenPoolContractsConfig{
				TokenSymbol: "TEST",
				NewPools: map[uint64]changeset.DeployTokenPoolInput{
					0: changeset.DeployTokenPoolInput{},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Input: changeset.DeployTokenPoolContractsConfig{
				TokenSymbol: "TEST",
				NewPools: map[uint64]changeset.DeployTokenPoolInput{
					5009297550715157269: changeset.DeployTokenPoolInput{},
				},
			},
			ErrStr: "does not exist in environment",
		},
		{
			Msg: "Router contract is missing from chain",
			Input: changeset.DeployTokenPoolContractsConfig{
				TokenSymbol: "TEST",
				NewPools: map[uint64]changeset.DeployTokenPoolInput{
					e.AllChainSelectors()[0]: changeset.DeployTokenPoolInput{},
				},
			},
			ErrStr: "missing router",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.Input.Validate(e)
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestValidateDeployTokenPoolInput(t *testing.T) {
	t.Parallel()

	e, selectorA, _, tokens, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)
	acceptLiquidity := false
	invalidAddress := utils.RandomAddress()

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]changeset.DeployTokenPoolInput{
		selectorA: {
			Type:               changeset.BurnMintTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, true)

	tests := []struct {
		Msg    string
		Symbol changeset.TokenSymbol
		Input  changeset.DeployTokenPoolInput
		ErrStr string
	}{
		{
			Msg:    "Token address is missing",
			Input:  changeset.DeployTokenPoolInput{},
			ErrStr: "token address must be defined",
		},
		{
			Msg: "Token pool type is missing",
			Input: changeset.DeployTokenPoolInput{
				TokenAddress: invalidAddress,
			},
			ErrStr: "type must be defined",
		},
		{
			Msg: "Token pool type is invalid",
			Input: changeset.DeployTokenPoolInput{
				TokenAddress: invalidAddress,
				Type:         deployment.ContractType("InvalidTokenPool"),
			},
			ErrStr: "requested token pool type InvalidTokenPool is unknown",
		},
		{
			Msg: "Token address is invalid",
			Input: changeset.DeployTokenPoolInput{
				Type:         changeset.BurnMintTokenPool,
				TokenAddress: invalidAddress,
			},
			ErrStr: fmt.Sprintf("failed to fetch symbol from token with address %s", invalidAddress),
		},
		{
			Msg:    "Token symbol mismatch",
			Symbol: "WRONG",
			Input: changeset.DeployTokenPoolInput{
				Type:         changeset.BurnMintTokenPool,
				TokenAddress: tokens[selectorA].Address,
			},
			ErrStr: fmt.Sprintf("symbol of token with address %s (%s) does not match expected symbol (WRONG)", tokens[selectorA].Address, testhelpers.TestTokenSymbol),
		},
		{
			Msg:    "Token decimal mismatch",
			Symbol: testhelpers.TestTokenSymbol,
			Input: changeset.DeployTokenPoolInput{
				Type:               changeset.BurnMintTokenPool,
				TokenAddress:       tokens[selectorA].Address,
				LocalTokenDecimals: 17,
			},
			ErrStr: fmt.Sprintf("decimals of token with address %s (%d) does not match localTokenDecimals (17)", tokens[selectorA].Address, testhelpers.LocalTokenDecimals),
		},
		{
			Msg:    "Accept liquidity should be defined",
			Symbol: testhelpers.TestTokenSymbol,
			Input: changeset.DeployTokenPoolInput{
				Type:               changeset.LockReleaseTokenPool,
				TokenAddress:       tokens[selectorA].Address,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			},
			ErrStr: "accept liquidity must be defined for lock release pools",
		},
		{
			Msg:    "Accept liquidity should be omitted",
			Symbol: testhelpers.TestTokenSymbol,
			Input: changeset.DeployTokenPoolInput{
				Type:               changeset.BurnMintTokenPool,
				TokenAddress:       tokens[selectorA].Address,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				AcceptLiquidity:    &acceptLiquidity,
			},
			ErrStr: "accept liquidity must be nil for burn mint pools",
		},
		{
			Msg:    "Token pool already exists",
			Symbol: testhelpers.TestTokenSymbol,
			Input: changeset.DeployTokenPoolInput{
				Type:               changeset.BurnMintTokenPool,
				TokenAddress:       tokens[selectorA].Address,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			},
			ErrStr: fmt.Sprintf("token pool with type BurnMintTokenPool and version %s already exists", deployment.Version1_5_1),
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			state, err := changeset.LoadOnchainState(e)
			require.NoError(t, err)

			err = test.Input.Validate(context.Background(), e.Chains[selectorA], state.Chains[selectorA], test.Symbol)
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestDeployTokenPool(t *testing.T) {
	t.Parallel()

	e, selectorA, _, tokens, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)
	acceptLiquidity := false

	tests := []struct {
		Msg   string
		Input changeset.DeployTokenPoolInput
	}{
		{
			Msg: "BurnMint",
			Input: changeset.DeployTokenPoolInput{
				TokenAddress:       tokens[selectorA].Address,
				Type:               changeset.BurnMintTokenPool,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				AllowList:          []common.Address{},
			},
		},
		{
			Msg: "BurnWithFromMint",
			Input: changeset.DeployTokenPoolInput{
				TokenAddress:       tokens[selectorA].Address,
				Type:               changeset.BurnWithFromMintTokenPool,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				AllowList:          []common.Address{},
			},
		},
		{
			Msg: "BurnFromMint",
			Input: changeset.DeployTokenPoolInput{
				TokenAddress:       tokens[selectorA].Address,
				Type:               changeset.BurnFromMintTokenPool,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				AllowList:          []common.Address{},
			},
		},
		{
			Msg: "LockRelease",
			Input: changeset.DeployTokenPoolInput{
				TokenAddress:       tokens[selectorA].Address,
				Type:               changeset.LockReleaseTokenPool,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				AllowList:          []common.Address{},
				AcceptLiquidity:    &acceptLiquidity,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			addressBook := deployment.NewMemoryAddressBook()
			state, err := changeset.LoadOnchainState(e)
			require.NoError(t, err)

			_, err = changeset.DeployTokenPool(
				e.Logger,
				e.Chains[selectorA],
				state.Chains[selectorA],
				addressBook,
				test.Input,
			)
			require.NoError(t, err)

			err = e.ExistingAddresses.Merge(addressBook)
			require.NoError(t, err)

			state, err = changeset.LoadOnchainState(e)
			require.NoError(t, err)

			switch test.Input.Type {
			case changeset.BurnMintTokenPool:
				_, ok := state.Chains[selectorA].BurnMintTokenPools[testhelpers.TestTokenSymbol]
				require.True(t, ok)
			case changeset.LockReleaseTokenPool:
				_, ok := state.Chains[selectorA].LockReleaseTokenPools[testhelpers.TestTokenSymbol]
				require.True(t, ok)
			case changeset.BurnWithFromMintTokenPool:
				_, ok := state.Chains[selectorA].BurnWithFromMintTokenPools[testhelpers.TestTokenSymbol]
				require.True(t, ok)
			case changeset.BurnFromMintTokenPool:
				_, ok := state.Chains[selectorA].BurnFromMintTokenPools[testhelpers.TestTokenSymbol]
				require.True(t, ok)
			}
		})
	}
}

func TestDeployTokenPoolContracts(t *testing.T) {
	t.Parallel()

	e, selectorA, _, tokens, timelockContracts := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	e, err := commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
		commonchangeset.ChangesetApplication{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployTokenPoolContractsChangeset),
			Config: changeset.DeployTokenPoolContractsConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				NewPools: map[uint64]changeset.DeployTokenPoolInput{
					selectorA: {
						TokenAddress:       tokens[selectorA].Address,
						Type:               changeset.BurnMintTokenPool,
						LocalTokenDecimals: testhelpers.LocalTokenDecimals,
						AllowList:          []common.Address{},
					},
				},
			},
		},
	})
	require.NoError(t, err)

	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	burnMintTokenPools, ok := state.Chains[selectorA].BurnMintTokenPools[testhelpers.TestTokenSymbol]
	require.True(t, ok)
	require.Len(t, burnMintTokenPools, 1)
	owner, err := burnMintTokenPools[deployment.Version1_5_1].Owner(nil)
	require.NoError(t, err)
	require.Equal(t, e.Chains[selectorA].DeployerKey.From, owner)
}
