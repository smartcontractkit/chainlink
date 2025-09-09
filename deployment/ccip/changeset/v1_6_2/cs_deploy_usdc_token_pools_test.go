package v1_6_2_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6_2"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func setupUSDCTokenPoolsEnvironmentForDeploy(t *testing.T, withPrereqs bool) (cldf.Environment, []uint64) {
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

func setupUSDCTokenPoolsContractsForDeploy(
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

func TestValidateDeployUSDCTokenPoolContractsConfig(t *testing.T) {
	t.Parallel()

	env, selectors := setupUSDCTokenPoolsEnvironmentForDeploy(t, true)

	require.GreaterOrEqual(t, len(selectors), 1)
	selector := selectors[0]

	tests := []struct {
		Msg    string
		Input  v1_6_2.DeployUSDCTokenPoolContractsConfig
		ErrStr string
	}{
		{
			Msg: "Chain selector is not valid",
			Input: v1_6_2.DeployUSDCTokenPoolContractsConfig{
				USDCPools: map[uint64]v1_6_2.DeployUSDCTokenPoolInput{
					0: {},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Input: v1_6_2.DeployUSDCTokenPoolContractsConfig{
				USDCPools: map[uint64]v1_6_2.DeployUSDCTokenPoolInput{
					5009297550715157269: {},
				},
			},
			ErrStr: "does not exist in environment",
		},
		{
			Msg: "No proxy",
			Input: v1_6_2.DeployUSDCTokenPoolContractsConfig{
				USDCPools: map[uint64]v1_6_2.DeployUSDCTokenPoolInput{
					selector: {
						PreviousPoolAddress: v1_6_2.USDCTokenPoolSentinelAddress,
						TokenMessenger:      utils.RandomAddress(),
						TokenAddress:        utils.RandomAddress(),
					},
				},
			},
			ErrStr: fmt.Sprintf(
				"CCTP message transmitter proxy for version %s not found",
				deployment.Version1_6_2,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := v1_6_2.DeployUSDCTokenPoolNew.VerifyPreconditions(env, test.Input)
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestValidateDeployUSDCTokenPoolInput(t *testing.T) {
	t.Parallel()

	env, selectors := setupUSDCTokenPoolsEnvironmentForDeploy(t, true)
	blockchain := env.BlockChains.EVMChains()[selectors[0]]
	addrBook := cldf.NewMemoryAddressBook()

	usdcToken, tokenMessenger := setupUSDCTokenPoolsContractsForDeploy(t,
		env.Logger,
		blockchain,
		addrBook,
	)

	nonUsdcToken, err := cldf.DeployContract(env.Logger, blockchain, addrBook,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
				chain.DeployerKey,
				chain.Client,
				"NOTUSDC",
				"NOTUSDC",
				6,
				big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
			)
			return cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
				Address:  tokenAddress,
				Contract: token,
				Tv:       cldf.NewTypeAndVersion(shared.USDCTokenPool, deployment.Version1_6_2),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	env, err = commoncs.Apply(t, env,
		commoncs.Configure(
			v1_6_2.DeployCCTPMessageTransmitterProxyNew,
			v1_6_2.DeployCCTPMessageTransmitterProxyContractConfig{
				USDCProxies: map[uint64]v1_6_2.DeployCCTPMessageTransmitterProxyInput{
					blockchain.Selector: {
						TokenMessenger: tokenMessenger.Address,
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
		Input  v1_6_2.DeployUSDCTokenPoolInput
		ErrStr string
	}{
		{
			Msg: "Token address is not defined",
			Input: v1_6_2.DeployUSDCTokenPoolInput{
				TokenAddress: utils.ZeroAddress,
				PoolType:     shared.USDCTokenPool,
			},
			ErrStr: "token address must be defined",
		},
		{
			Msg: "Token messenger address is not defined",
			Input: v1_6_2.DeployUSDCTokenPoolInput{
				TokenMessenger: utils.ZeroAddress,
				TokenAddress:   utils.RandomAddress(),
				PoolType:       shared.USDCTokenPool,
			},
			ErrStr: "token messenger must be defined",
		},
		{
			Msg: "No previous pool",
			Input: v1_6_2.DeployUSDCTokenPoolInput{
				PreviousPoolAddress: utils.ZeroAddress,
				TokenMessenger:      utils.RandomAddress(),
				TokenAddress:        utils.RandomAddress(),
				PoolType:            shared.USDCTokenPool,
			},
			ErrStr: "unable to find a previous pool",
		},
		{
			Msg: "Can't reach token",
			Input: v1_6_2.DeployUSDCTokenPoolInput{
				PreviousPoolAddress: v1_6_2.USDCTokenPoolSentinelAddress,
				TokenAddress:        utils.RandomAddress(),
				TokenMessenger:      utils.RandomAddress(),
				PoolType:            shared.USDCTokenPool,
			},
			ErrStr: "failed to fetch symbol from token",
		},
		{
			Msg: "Symbol is wrong",
			Input: v1_6_2.DeployUSDCTokenPoolInput{
				PreviousPoolAddress: v1_6_2.USDCTokenPoolSentinelAddress,
				TokenAddress:        nonUsdcToken.Address,
				TokenMessenger:      utils.RandomAddress(),
				PoolType:            shared.USDCTokenPool,
			},
			ErrStr: "is not USDC",
		},
		{
			Msg: "Can't reach token messenger",
			Input: v1_6_2.DeployUSDCTokenPoolInput{
				PreviousPoolAddress: v1_6_2.USDCTokenPoolSentinelAddress,
				TokenAddress:        usdcToken.Address,
				TokenMessenger:      utils.RandomAddress(),
				PoolType:            shared.USDCTokenPool,
			},
			ErrStr: "failed to fetch local message transmitter from address",
		},
		{
			Msg: "Invalid pool type",
			Input: v1_6_2.DeployUSDCTokenPoolInput{
				PreviousPoolAddress: v1_6_2.USDCTokenPoolSentinelAddress,
				TokenAddress:        usdcToken.Address,
				TokenMessenger:      tokenMessenger.Address,
				PoolType:            "bad pool type",
			},
			ErrStr: "unsupported pool type",
		},
		{
			Msg: "No error",
			Input: v1_6_2.DeployUSDCTokenPoolInput{
				PreviousPoolAddress: v1_6_2.USDCTokenPoolSentinelAddress,
				TokenAddress:        usdcToken.Address,
				TokenMessenger:      tokenMessenger.Address,
				PoolType:            shared.USDCTokenPool,
			},
			ErrStr: "",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.Input.Validate(env.GetContext(), blockchain, state.Chains[blockchain.Selector])
			if test.ErrStr != "" {
				require.Contains(t, err.Error(), test.ErrStr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDeployUSDCTokenPool(t *testing.T) {
	t.Parallel()

	env, selectors := setupUSDCTokenPoolsEnvironmentForDeploy(t, true)
	addrBook := cldf.NewMemoryAddressBook()

	newUSDCMsgProxies := make(map[uint64]v1_6_2.DeployCCTPMessageTransmitterProxyInput, len(selectors))
	newUSDCTokenPools := make(map[uint64]v1_6_2.DeployUSDCTokenPoolInput, len(selectors))
	for _, selector := range selectors {
		blockchain := env.BlockChains.EVMChains()[selector]
		usdcToken, tokenMessenger := setupUSDCTokenPoolsContractsForDeploy(t, env.Logger, blockchain, addrBook)

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

	state, err := stateview.LoadOnchainState(env)
	require.NoError(t, err)
	for _, selector := range selectors {
		usdcTokenPools := state.Chains[selector].USDCTokenPoolsV1_6
		require.Len(t, usdcTokenPools, 1, selector)

		owner, err := usdcTokenPools[deployment.Version1_6_2].Owner(nil)
		require.NoError(t, err)

		deployer := env.BlockChains.EVMChains()[selector].DeployerKey.From
		require.Equal(t, deployer, owner)
	}
}

func TestDeployHybridLockReleaseUSDCTokenPool(t *testing.T) {
	t.Parallel()

	env, selectors := setupUSDCTokenPoolsEnvironmentForDeploy(t, true)
	addrBook := cldf.NewMemoryAddressBook()

	newUSDCMsgProxies := make(map[uint64]v1_6_2.DeployCCTPMessageTransmitterProxyInput, len(selectors))
	newUSDCTokenPools := make(map[uint64]v1_6_2.DeployUSDCTokenPoolInput, len(selectors))
	for _, selector := range selectors {
		blockchain := env.BlockChains.EVMChains()[selector]
		usdcToken, tokenMessenger := setupUSDCTokenPoolsContractsForDeploy(t, env.Logger, blockchain, addrBook)

		newUSDCMsgProxies[selector] = v1_6_2.DeployCCTPMessageTransmitterProxyInput{
			TokenMessenger: tokenMessenger.Address,
		}

		newUSDCTokenPools[selector] = v1_6_2.DeployUSDCTokenPoolInput{
			PreviousPoolAddress: v1_6_2.USDCTokenPoolSentinelAddress,
			TokenMessenger:      tokenMessenger.Address,
			TokenAddress:        usdcToken.Address,
			PoolType:            shared.HybridLockReleaseUSDCTokenPool,
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

	state, err := stateview.LoadOnchainState(env)
	require.NoError(t, err)
	for _, selector := range selectors {
		usdcTokenPools := state.Chains[selector].USDCTokenPoolsV1_6
		require.Len(t, usdcTokenPools, 1, selector)

		owner, err := usdcTokenPools[deployment.Version1_6_2].Owner(nil)
		require.NoError(t, err)

		deployer := env.BlockChains.EVMChains()[selector].DeployerKey.From
		require.Equal(t, deployer, owner)
	}
}
