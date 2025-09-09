package v1_6_2_test

import (
	"maps"
	"math/big"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6_2"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_2/usdc_token_pool"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func setupCCTPMsgTransmitterProxyEnvironmentForConfigure(t *testing.T, withPrereqs bool) (cldf.Environment, []uint64) {
	env := memory.NewMemoryEnvironment(t,
		logger.Test(t),
		zapcore.InfoLevel,
		memory.MemoryEnvironmentConfig{Chains: 2},
	)

	selectors := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	if withPrereqs {
		var err error

		prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, len(selectors))
		for i, selector := range selectors {
			prereqCfg[i] = changeset.DeployPrerequisiteConfigPerChain{
				ChainSelector: selector,
			}
		}

		env, err = commoncs.Apply(t, env,
			commoncs.Configure(
				cldf.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
				changeset.DeployPrerequisiteConfig{
					Configs: prereqCfg,
				},
			),
		)
		require.NoError(t, err)
	}

	return env, selectors
}

func setupCCTPMsgTransmitterProxyContractsForConfigure(
	t *testing.T,
	logger logger.Logger,
	chain cldf_evm.Chain,
	addressBook cldf.AddressBook,
) (
	*cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677],
	*cldf.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger],
) {
	usdcToken, err := cldf.DeployContract(logger, chain, addressBook,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
				chain.DeployerKey,
				chain.Client,
				"USDC",
				"USDC",
				6,
				big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
			)
			return cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
				Address:  tokenAddress,
				Contract: token,
				Tv:       cldf.NewTypeAndVersion(shared.BurnMintToken, deployment.Version1_6_2),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	transmitter, err := cldf.DeployContract(logger, chain, addressBook,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter] {
			transmitterAddress, tx, transmitter, err := mock_usdc_token_transmitter.DeployMockE2EUSDCTransmitter(chain.DeployerKey, chain.Client, 0, 1, usdcToken.Address)
			return cldf.ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter]{
				Address:  transmitterAddress,
				Contract: transmitter,
				Tv:       cldf.NewTypeAndVersion(shared.USDCMockTransmitter, deployment.Version1_6_2),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	messenger, err := cldf.DeployContract(logger, chain, addressBook,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger] {
			messengerAddress, tx, messenger, err := mock_usdc_token_messenger.DeployMockE2EUSDCTokenMessenger(chain.DeployerKey, chain.Client, 0, transmitter.Address)
			return cldf.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger]{
				Address:  messengerAddress,
				Contract: messenger,
				Tv:       cldf.NewTypeAndVersion(shared.USDCTokenMessenger, deployment.Version1_6_2),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	return usdcToken, messenger
}

func TestValidateConfigureCCTPMessageTransmitterProxyInput(t *testing.T) {
	t.Parallel()

	env, selectors := setupCCTPMsgTransmitterProxyEnvironmentForConfigure(t, true)
	require.GreaterOrEqual(t, len(selectors), 1)
	chain := env.BlockChains.EVMChains()[selectors[0]]

	addressBook := cldf.NewMemoryAddressBook()
	_, tokenMsngr := setupCCTPMsgTransmitterProxyContractsForConfigure(t, env.Logger, chain, addressBook)

	env, err := commoncs.Apply(t, env,
		commoncs.Configure(
			v1_6_2.DeployCCTPMessageTransmitterProxyNew,
			v1_6_2.DeployCCTPMessageTransmitterProxyContractConfig{
				USDCProxies: map[uint64]v1_6_2.DeployCCTPMessageTransmitterProxyInput{
					chain.Selector: {
						TokenMessenger: tokenMsngr.Address,
					},
				},
			},
		),
	)
	require.NoError(t, err)

	state, err := stateview.LoadOnchainState(env)
	require.NoError(t, err)

	tests := []struct {
		Msg    string
		Input  v1_6_2.ConfigureCCTPMessageTransmitterProxyInput
		ErrStr string
	}{
		{
			Msg: "Allowed caller cannot be zero address",
			Input: v1_6_2.ConfigureCCTPMessageTransmitterProxyInput{
				AllowedCallerUpdates: []v1_6_2.AllowedCallerUpdate{
					{
						AllowedCaller: utils.ZeroAddress,
						Enabled:       true,
					},
				},
			},
			ErrStr: "allowed caller must be defined for chain",
		},
		{
			Msg: "Invalid allowed caller",
			Input: v1_6_2.ConfigureCCTPMessageTransmitterProxyInput{
				AllowedCallerUpdates: []v1_6_2.AllowedCallerUpdate{
					{
						AllowedCaller: utils.RandomAddress(),
						Enabled:       true,
					},
				},
			},
			ErrStr: "allowed caller does not match any existing 1.6 USDC token pools",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.Input.Validate(env.GetContext(), chain, state.Chains[chain.Selector])
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestConfigureCCTPMessageTransmitterProxy(t *testing.T) {
	t.Parallel()

	env, selectors := setupCCTPMsgTransmitterProxyEnvironmentForConfigure(t, true)
	addrBook := cldf.NewMemoryAddressBook()

	newUSDCMsgProxies := make(map[uint64]v1_6_2.DeployCCTPMessageTransmitterProxyInput, len(selectors))
	newUSDCTokenPools := make(map[uint64]v1_6_2.DeployUSDCTokenPoolInput, len(selectors))
	for _, selector := range selectors {
		usdcToken, tokenMessenger := setupCCTPMsgTransmitterProxyContractsForConfigure(t,
			env.Logger,
			env.BlockChains.EVMChains()[selector],
			addrBook,
		)

		newUSDCMsgProxies[selector] = v1_6_2.DeployCCTPMessageTransmitterProxyInput{
			TokenMessenger: tokenMessenger.Address,
		}

		newUSDCTokenPools[selector] = v1_6_2.DeployUSDCTokenPoolInput{
			PreviousPoolAddress: v1_6_2.USDCTokenPoolSentinelAddress,
			TokenMessenger:      tokenMessenger.Address,
			TokenAddress:        usdcToken.Address,
			PoolType:            shared.USDCTokenPool,
		}
	}

	env, err := commoncs.Apply(t, env,
		commoncs.Configure(
			v1_6_2.DeployCCTPMessageTransmitterProxyNew,
			v1_6_2.DeployCCTPMessageTransmitterProxyContractConfig{
				USDCProxies: newUSDCMsgProxies,
			},
		),
	)
	require.NoError(t, err)

	env, err = commoncs.Apply(t, env,
		commoncs.Configure(
			v1_6_2.DeployUSDCTokenPoolNew,
			v1_6_2.DeployUSDCTokenPoolContractsConfig{
				USDCPools: newUSDCTokenPools,
			},
		),
	)
	require.NoError(t, err)

	startState, err := stateview.LoadOnchainState(env)
	require.NoError(t, err)

	newUSDCProxyCnfgs := make(map[uint64]v1_6_2.ConfigureCCTPMessageTransmitterProxyInput, len(selectors))
	for _, selector := range selectors {
		pools := startState.Chains[selector].USDCTokenPoolsV1_6
		input := make([]v1_6_2.AllowedCallerUpdate, len(pools))

		for i, pool := range slices.AppendSeq([]*usdc_token_pool.USDCTokenPool{}, maps.Values(pools)) {
			input[i] = v1_6_2.AllowedCallerUpdate{
				AllowedCaller: pool.Address(),
				Enabled:       true,
			}
		}

		newUSDCProxyCnfgs[selector] = v1_6_2.ConfigureCCTPMessageTransmitterProxyInput{
			AllowedCallerUpdates: input,
		}
	}

	env, err = commoncs.Apply(t, env,
		commoncs.Configure(
			v1_6_2.ConfigureCCTPMessageTransmitterProxy,
			v1_6_2.ConfigureCCTPMessageTransmitterProxyContractConfig{
				CCTPProxies: newUSDCProxyCnfgs,
			},
		),
	)
	require.NoError(t, err)

	finalState, err := stateview.LoadOnchainState(env)
	require.NoError(t, err)
	for _, selector := range selectors {
		proxies := finalState.Chains[selector].CCTPMessageTransmitterProxies
		updates := newUSDCProxyCnfgs[selector].AllowedCallerUpdates
		require.Len(t, proxies, 1)

		expectedCallers := make([]common.Address, len(updates))
		for i, cfg := range updates {
			expectedCallers[i] = cfg.AllowedCaller
		}

		actualCallers, err := proxies[deployment.Version1_6_2].GetAllowedCallers(nil)
		require.NoError(t, err)

		require.ElementsMatch(t,
			expectedCallers,
			actualCallers,
		)
	}
}
