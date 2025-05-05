package v1_5_1_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
		Input       v1_5_1.DeployTokenPoolContractsConfig
		ErrStr      string
	}{
		{
			Msg:    "Token symbol is missing",
			Input:  v1_5_1.DeployTokenPoolContractsConfig{},
			ErrStr: "token symbol must be defined",
		},
		{
			Msg: "Chain selector is not valid",
			Input: v1_5_1.DeployTokenPoolContractsConfig{
				TokenSymbol: "TEST",
				NewPools: map[uint64]v1_5_1.DeployTokenPoolInput{
					0: v1_5_1.DeployTokenPoolInput{},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Input: v1_5_1.DeployTokenPoolContractsConfig{
				TokenSymbol: "TEST",
				NewPools: map[uint64]v1_5_1.DeployTokenPoolInput{
					5009297550715157269: v1_5_1.DeployTokenPoolInput{},
				},
			},
			ErrStr: "does not exist in environment",
		},
		{
			Msg: "Router contract is missing from chain",
			Input: v1_5_1.DeployTokenPoolContractsConfig{
				TokenSymbol: "TEST",
				NewPools: map[uint64]v1_5_1.DeployTokenPoolInput{
					e.AllChainSelectors()[0]: v1_5_1.DeployTokenPoolInput{},
				},
			},
			ErrStr: "missing router",
		},
		{
			Msg: "Test router contract is missing from chain",
			Input: v1_5_1.DeployTokenPoolContractsConfig{
				TokenSymbol: "TEST",
				NewPools: map[uint64]v1_5_1.DeployTokenPoolInput{
					e.AllChainSelectors()[0]: v1_5_1.DeployTokenPoolInput{},
				},
				IsTestRouter: true,
			},
			ErrStr: "missing test router",
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

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               changeset.BurnMintTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, true)

	tests := []struct {
		Msg    string
		Symbol changeset.TokenSymbol
		Input  v1_5_1.DeployTokenPoolInput
		ErrStr string
	}{
		{
			Msg:    "Token address is missing",
			Input:  v1_5_1.DeployTokenPoolInput{},
			ErrStr: "token address must be defined",
		},
		{
			Msg: "Token pool type is missing",
			Input: v1_5_1.DeployTokenPoolInput{
				TokenAddress: invalidAddress,
			},
			ErrStr: "type must be defined",
		},
		{
			Msg: "Token pool type is invalid",
			Input: v1_5_1.DeployTokenPoolInput{
				TokenAddress: invalidAddress,
				Type:         deployment.ContractType("InvalidTokenPool"),
			},
			ErrStr: "requested token pool type InvalidTokenPool is unknown",
		},
		{
			Msg: "Token address is invalid",
			Input: v1_5_1.DeployTokenPoolInput{
				Type:         changeset.BurnMintTokenPool,
				TokenAddress: invalidAddress,
			},
			ErrStr: fmt.Sprintf("failed to fetch symbol from token with address %s", invalidAddress),
		},
		{
			Msg:    "Token symbol mismatch",
			Symbol: "WRONG",
			Input: v1_5_1.DeployTokenPoolInput{
				Type:         changeset.BurnMintTokenPool,
				TokenAddress: tokens[selectorA].Address,
			},
			ErrStr: fmt.Sprintf("symbol of token with address %s (%s) does not match expected symbol (WRONG)", tokens[selectorA].Address, testhelpers.TestTokenSymbol),
		},
		{
			Msg:    "Token decimal mismatch",
			Symbol: testhelpers.TestTokenSymbol,
			Input: v1_5_1.DeployTokenPoolInput{
				Type:               changeset.BurnMintTokenPool,
				TokenAddress:       tokens[selectorA].Address,
				LocalTokenDecimals: 17,
			},
			ErrStr: fmt.Sprintf("decimals of token with address %s (%d) does not match localTokenDecimals (17)", tokens[selectorA].Address, testhelpers.LocalTokenDecimals),
		},
		{
			Msg:    "Accept liquidity should be defined",
			Symbol: testhelpers.TestTokenSymbol,
			Input: v1_5_1.DeployTokenPoolInput{
				Type:               changeset.LockReleaseTokenPool,
				TokenAddress:       tokens[selectorA].Address,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			},
			ErrStr: "accept liquidity must be defined for lock release pools",
		},
		{
			Msg:    "Accept liquidity should be omitted",
			Symbol: testhelpers.TestTokenSymbol,
			Input: v1_5_1.DeployTokenPoolInput{
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
			Input: v1_5_1.DeployTokenPoolInput{
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

func TestDeployTokenPoolContracts(t *testing.T) {
	t.Parallel()

	acceptLiquidity := false

	type Ownable interface {
		Owner(opts *bind.CallOpts) (common.Address, error)
	}

	tests := []struct {
		Msg     string
		Input   v1_5_1.DeployTokenPoolInput
		GetPool func(changeset.CCIPChainState) Ownable
	}{
		{
			Msg: "BurnMint",
			Input: v1_5_1.DeployTokenPoolInput{
				Type:               changeset.BurnMintTokenPool,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				AllowList:          []common.Address{},
			},
			GetPool: func(cs changeset.CCIPChainState) Ownable {
				tokenPools, ok := cs.BurnMintTokenPools[testhelpers.TestTokenSymbol]
				require.True(t, ok)
				require.Len(t, tokenPools, 1)
				return tokenPools[deployment.Version1_5_1]
			},
		},
		{
			Msg: "BurnWithFromMint",
			Input: v1_5_1.DeployTokenPoolInput{
				Type:               changeset.BurnWithFromMintTokenPool,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				AllowList:          []common.Address{},
			},
			GetPool: func(cs changeset.CCIPChainState) Ownable {
				tokenPools, ok := cs.BurnWithFromMintTokenPools[testhelpers.TestTokenSymbol]
				require.True(t, ok)
				require.Len(t, tokenPools, 1)
				return tokenPools[deployment.Version1_5_1]
			},
		},
		{
			Msg: "BurnFromMint",
			Input: v1_5_1.DeployTokenPoolInput{
				Type:               changeset.BurnFromMintTokenPool,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				AllowList:          []common.Address{},
			},
			GetPool: func(cs changeset.CCIPChainState) Ownable {
				tokenPools, ok := cs.BurnFromMintTokenPools[testhelpers.TestTokenSymbol]
				require.True(t, ok)
				require.Len(t, tokenPools, 1)
				return tokenPools[deployment.Version1_5_1]
			},
		},
		{
			Msg: "LockRelease",
			Input: v1_5_1.DeployTokenPoolInput{
				Type:               changeset.LockReleaseTokenPool,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
				AllowList:          []common.Address{},
				AcceptLiquidity:    &acceptLiquidity,
			},
			GetPool: func(cs changeset.CCIPChainState) Ownable {
				tokenPools, ok := cs.LockReleaseTokenPools[testhelpers.TestTokenSymbol]
				require.True(t, ok)
				require.Len(t, tokenPools, 1)
				return tokenPools[deployment.Version1_5_1]
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			e, selectorA, _, tokens, timelockContracts := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

			test.Input.TokenAddress = tokens[selectorA].Address

			e, err := commonchangeset.Apply(t, e, timelockContracts,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_5_1.DeployTokenPoolContractsChangeset),
					v1_5_1.DeployTokenPoolContractsConfig{
						TokenSymbol: testhelpers.TestTokenSymbol,
						NewPools: map[uint64]v1_5_1.DeployTokenPoolInput{
							selectorA: test.Input,
						},
					},
				),
			)
			require.NoError(t, err)

			state, err := changeset.LoadOnchainState(e)
			require.NoError(t, err)

			pool := test.GetPool(state.Chains[selectorA])
			owner, err := pool.Owner(nil)
			require.NoError(t, err)
			require.Equal(t, e.Chains[selectorA].DeployerKey.From, owner)
		})
	}
}
