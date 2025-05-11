package v1_5_1_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestAddTokenE2E(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		externalAdmin bool
		withMCMS      bool
		withNewToken  bool
	}{
		{
			name:          "e2e token configuration with external admin",
			externalAdmin: true,
			withMCMS:      true,
		},
		{
			name:          "e2e token configuration with external admin without mcms",
			externalAdmin: true,
			withMCMS:      false,
		},
		{
			name:          "e2e token configuration with admin as token admin registry with MCMS",
			externalAdmin: false,
			withMCMS:      true,
		},
		{
			name:          "e2e token configuration with external token admin registry without MCMS",
			externalAdmin: false,
			withMCMS:      false,
		},
		{
			name:          "e2e token configuration with admin as token admin registry with MCMS with new token",
			externalAdmin: false,
			withMCMS:      true,
			withNewToken:  true,
		},
		{
			name:          "e2e token configuration with admin as token admin registry without MCMS with new token",
			externalAdmin: false,
			withMCMS:      false,
			withNewToken:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			var (
				e                    deployment.Environment
				selectorA, selectorB uint64
				mcmsConfig           *proposalutils.TimelockConfig
				err                  error
			)

			tokens := make(map[uint64]*cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677])
			timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts)
			if test.withMCMS {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay:   0,
					MCMSAction: mcmsTypes.TimelockActionSchedule,
				}
			}
			// we deploy the token separately as part of env set up
			if !test.withNewToken {
				e, selectorA, selectorB, tokens, timelockContracts = testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), test.withMCMS)
			} else {
				// we deploy the token as part of AddTokenE2E changeset
				tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithPrerequisiteDeploymentOnly(nil))
				e = tenv.Env
				state, err := changeset.LoadOnchainState(e)
				require.NoError(t, err)
				selectors := e.AllChainSelectors()
				selectorA = selectors[0]
				selectorB = selectors[1]
				// We only need the token admin registry to be owned by the timelock in these tests
				timelockOwnedContractsByChain := make(map[uint64][]common.Address)
				for _, selector := range selectors {
					timelockOwnedContractsByChain[selector] = []common.Address{state.Chains[selector].TokenAdminRegistry.Address()}
				}

				// Assemble map of addresses required for Timelock scheduling & execution
				for _, selector := range selectors {
					timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
						Timelock:  state.Chains[selector].Timelock,
						CallProxy: state.Chains[selector].CallProxy,
					}
				}
				if test.withMCMS {
					e, err = commonchangeset.Apply(t, e, timelockContracts,
						commonchangeset.Configure(
							cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
							commonchangeset.TransferToMCMSWithTimelockConfig{
								ContractsByChain: timelockOwnedContractsByChain,
								MCMSConfig:       *mcmsConfig,
							},
						),
					)
					require.NoError(t, err)
				}
			}

			externalAdmin := utils.ZeroAddress
			if test.externalAdmin {
				externalAdmin = utils.RandomAddress()
			}

			SelectorA2B := createSymmetricRateLimits(100, 1000)
			SelectorB2A := createSymmetricRateLimits(100, 1000)
			addTokenE2EConfig := v1_5_1.AddTokensE2EConfig{
				MCMS: mcmsConfig,
			}
			recipientAddress := utils.RandomAddress()
			topupAmount := big.NewInt(1000)
			// form the changeset input config
			for _, chain := range e.AllChainSelectors() {
				if addTokenE2EConfig.Tokens == nil {
					addTokenE2EConfig.Tokens = make(map[changeset.TokenSymbol]v1_5_1.AddTokenE2EConfig)
				}
				if _, ok := addTokenE2EConfig.Tokens[testhelpers.TestTokenSymbol]; !ok {
					addTokenE2EConfig.Tokens[testhelpers.TestTokenSymbol] = v1_5_1.AddTokenE2EConfig{
						PoolConfig: make(map[uint64]v1_5_1.E2ETokenAndPoolConfig),
					}
				}
				rateLimiterPerChain := make(map[uint64]v1_5_1.RateLimiterConfig)
				for range []uint64{selectorA, selectorB} {
					switch chain {
					case selectorA:
						rateLimiterPerChain[selectorB] = SelectorA2B
					case selectorB:
						rateLimiterPerChain[selectorA] = SelectorB2A
					}
				}
				poolConfig := addTokenE2EConfig.Tokens[testhelpers.TestTokenSymbol].PoolConfig
				var deployPoolConfig *v1_5_1.DeployTokenPoolInput
				var deployTokenConfig *v1_5_1.DeployTokenConfig
				if test.withNewToken {
					deployTokenConfig = &v1_5_1.DeployTokenConfig{
						TokenName:     string(testhelpers.TestTokenSymbol),
						TokenSymbol:   testhelpers.TestTokenSymbol,
						TokenDecimals: testhelpers.LocalTokenDecimals,
						MaxSupply:     big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
						Type:          changeset.BurnMintToken,
						PoolType:      changeset.BurnMintTokenPool,
						MintTokenForRecipients: map[common.Address]*big.Int{
							recipientAddress: topupAmount,
						},
					}
				} else {
					token := tokens[chain]
					deployPoolConfig = &v1_5_1.DeployTokenPoolInput{
						Type:               changeset.BurnMintTokenPool,
						TokenAddress:       token.Address,
						LocalTokenDecimals: testhelpers.LocalTokenDecimals,
					}
				}
				poolConfig[chain] = v1_5_1.E2ETokenAndPoolConfig{
					TokenDeploymentConfig: deployTokenConfig,
					DeployPoolConfig:      deployPoolConfig,
					PoolVersion:           deployment.Version1_5_1,
					ExternalAdmin:         externalAdmin,
					RateLimiterConfig:     rateLimiterPerChain,
				}
			}
			// apply the changeset
			e, err = commonchangeset.Apply(t, e, timelockContracts,
				commonchangeset.Configure(v1_5_1.AddTokensE2E, addTokenE2EConfig))
			require.NoError(t, err)

			state, err := changeset.LoadOnchainState(e)
			require.NoError(t, err)

			// populate token details in case of token deployment as part of changeset
			if test.withNewToken {
				// ensure the token is deployed
				for chain, chainState := range state.Chains {
					token, ok := chainState.BurnMintTokens677[testhelpers.TestTokenSymbol]
					require.True(t, ok)
					tokens[chain] = &cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
						Address:  token.Address(),
						Contract: token,
					}
					// check token balance
					balance, err := token.BalanceOf(&bind.CallOpts{Context: ctx}, recipientAddress)
					require.NoError(t, err)
					require.Equal(t, balance, topupAmount)
					// check minter role
					minterCheck, err := token.IsMinter(&bind.CallOpts{Context: ctx}, recipientAddress)
					require.NoError(t, err)
					require.True(t, minterCheck)
				}
			}
			registryOnA := state.Chains[selectorA].TokenAdminRegistry
			registryOnB := state.Chains[selectorB].TokenAdminRegistry

			poolOnSelectorB := state.Chains[selectorB].BurnMintTokenPools[testhelpers.TestTokenSymbol][deployment.Version1_5_1]
			poolOnSelectorA := state.Chains[selectorA].BurnMintTokenPools[testhelpers.TestTokenSymbol][deployment.Version1_5_1]
			// validate end results
			for chain, token := range tokens {
				// check token pool is deployed
				tokenpools, ok := state.Chains[chain].BurnMintTokenPools[testhelpers.TestTokenSymbol]
				require.True(t, ok)
				require.Len(t, tokenpools, 1)
				tokenPoolC, err := token_pool.NewTokenPool(tokenpools[deployment.Version1_5_1].Address(), e.Chains[chain].Client)
				require.NoError(t, err)
				var rateLimiterConfig v1_5_1.RateLimiterConfig
				var remotePoolAddr common.Address
				var registry *token_admin_registry.TokenAdminRegistry
				switch chain {
				case selectorA:
					rateLimiterConfig = SelectorA2B
					remotePoolAddr = poolOnSelectorB.Address()
					registry = registryOnA
				case selectorB:
					rateLimiterConfig = SelectorB2A
					remotePoolAddr = poolOnSelectorA.Address()
					registry = registryOnB
				}
				// check token pool is configured
				validateMemberOfTokenPoolPair(
					t,
					state,
					tokenPoolC,
					[]common.Address{remotePoolAddr},
					tokens,
					testhelpers.TestTokenSymbol,
					chain,
					rateLimiterConfig.Inbound.Rate, // inbound & outbound are the same in this test
					rateLimiterConfig.Inbound.Capacity,
					e.Chains[chain].DeployerKey.From, // the pools are still owned by the deployer
				)
				if test.withNewToken {
					// check token pool is added as minter
					minterCheck, err := token.Contract.IsMinter(&bind.CallOpts{Context: ctx}, tokenPoolC.Address())
					require.NoError(t, err)
					require.True(t, minterCheck)

					// check token pool is added as burner
					burnerCheck, err := token.Contract.IsBurner(&bind.CallOpts{Context: ctx}, tokenPoolC.Address())
					require.NoError(t, err)
					require.True(t, burnerCheck)
				}
				// check if admin and set pool is set correctly
				regConfig, err := registry.GetTokenConfig(&bind.CallOpts{Context: ctx}, token.Address)
				require.NoError(t, err)

				if !test.externalAdmin {
					// if not external admin then admin should be token admin registry
					// and pool should be set for token
					require.Equal(t, tokenPoolC.Address(), regConfig.TokenPool)
					if test.withMCMS {
						require.Equal(t, state.Chains[chain].Timelock.Address(), regConfig.Administrator)
					} else {
						require.Equal(t, e.Chains[chain].DeployerKey.From, regConfig.Administrator)
					}
				} else {
					// if external admin then PendingAdministrator should be external admin
					// as external admin has not accepted the admin role yet
					require.Equal(t, externalAdmin, regConfig.PendingAdministrator)
					require.Empty(t, regConfig.TokenPool)
				}
			}
		})
	}
}
