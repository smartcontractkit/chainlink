package v1_6_2_test

import (
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

func setupCCTPMsgTransmitterProxyEnvironmentForDeploy(t *testing.T, withPrereqs bool) (cldf.Environment, []uint64) {
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

func setupCCTPMsgTransmitterProxyContractsForDeploy(
	t *testing.T,
	logger logger.Logger,
	chain cldf_evm.Chain,
	addressBook cldf.AddressBook,
) *cldf.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger] {
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

	return messenger
}

func TestValidateDeployCCTPMessageTransmitterProxyInput(t *testing.T) {
	t.Parallel()

	env, selectors := setupCCTPMsgTransmitterProxyEnvironmentForDeploy(t, false)

	require.GreaterOrEqual(t, len(selectors), 1)
	chain := env.BlockChains.EVMChains()[selectors[0]]

	state, err := stateview.LoadOnchainState(env)
	require.NoError(t, err)

	tests := []struct {
		Msg    string
		Input  v1_6_2.DeployCCTPMessageTransmitterProxyInput
		ErrStr string
	}{
		{
			Msg:    "Empty token messenger address is not allowed",
			Input:  v1_6_2.DeployCCTPMessageTransmitterProxyInput{},
			ErrStr: "token messenger must be defined for chain",
		},
		{
			Msg: "Token messenger address cannot be the zero address",
			Input: v1_6_2.DeployCCTPMessageTransmitterProxyInput{
				TokenMessenger: utils.ZeroAddress,
			},
			ErrStr: "token messenger must be defined for chain",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.Input.Validate(env.GetContext(), chain, state.Chains[chain.Selector])
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestDeployCCTPMessageTransmitterProxy(t *testing.T) {
	t.Parallel()

	env, selectors := setupCCTPMsgTransmitterProxyEnvironmentForDeploy(t, true)

	newProxies := make(map[uint64]v1_6_2.DeployCCTPMessageTransmitterProxyInput, len(selectors))
	addressBook := cldf.NewMemoryAddressBook()
	for _, selector := range selectors {
		blockchain := env.BlockChains.EVMChains()[selector]
		tokenMsngr := setupCCTPMsgTransmitterProxyContractsForDeploy(t, env.Logger, blockchain, addressBook)
		newProxies[selector] = v1_6_2.DeployCCTPMessageTransmitterProxyInput{
			TokenMessenger: tokenMsngr.Address,
		}
	}

	env, err := commoncs.Apply(t, env,
		commoncs.Configure(
			v1_6_2.DeployCCTPMessageTransmitterProxyNew,
			v1_6_2.DeployCCTPMessageTransmitterProxyContractConfig{
				USDCProxies: newProxies,
			},
		),
	)
	require.NoError(t, err)

	state, err := stateview.LoadOnchainState(env)
	require.NoError(t, err)
	for _, selector := range selectors {
		proxies := state.Chains[selector].CCTPMessageTransmitterProxies
		require.Len(t, proxies, 1)

		owner, err := proxies[deployment.Version1_6_2].Owner(nil)
		require.NoError(t, err)

		deployer := env.BlockChains.EVMChains()[selector].DeployerKey.From
		require.Equal(t, deployer, owner)
	}
}
