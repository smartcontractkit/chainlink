package v1_5_1_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func deployUSDCPrerequisites(
	t *testing.T,
	logger logger.Logger,
	chain deployment.Chain,
	addressBook deployment.AddressBook,
) (*cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677], *cldf.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger]) {
	usdcToken, err := cldf.DeployContract(logger, chain, addressBook,
		func(chain deployment.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
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
				Tv:       deployment.NewTypeAndVersion(changeset.USDCTokenPool, deployment.Version1_5_1),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	transmitter, err := cldf.DeployContract(logger, chain, addressBook,
		func(chain deployment.Chain) cldf.ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter] {
			transmitterAddress, tx, transmitter, err := mock_usdc_token_transmitter.DeployMockE2EUSDCTransmitter(chain.DeployerKey, chain.Client, 0, 1, usdcToken.Address)
			return cldf.ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter]{
				Address:  transmitterAddress,
				Contract: transmitter,
				Tv:       deployment.NewTypeAndVersion(changeset.USDCMockTransmitter, deployment.Version1_0_0),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	messenger, err := cldf.DeployContract(logger, chain, addressBook,
		func(chain deployment.Chain) cldf.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger] {
			messengerAddress, tx, messenger, err := mock_usdc_token_messenger.DeployMockE2EUSDCTokenMessenger(chain.DeployerKey, chain.Client, 0, transmitter.Address)
			return cldf.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger]{
				Address:  messengerAddress,
				Contract: messenger,
				Tv:       deployment.NewTypeAndVersion(changeset.USDCTokenMessenger, deployment.Version1_0_0),
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

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})

	tests := []struct {
		Msg    string
		Input  v1_5_1.DeployUSDCTokenPoolContractsConfig
		ErrStr string
	}{
		{
			Msg: "Chain selector is not valid",
			Input: v1_5_1.DeployUSDCTokenPoolContractsConfig{
				USDCPools: map[uint64]v1_5_1.DeployUSDCTokenPoolInput{
					0: v1_5_1.DeployUSDCTokenPoolInput{},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Input: v1_5_1.DeployUSDCTokenPoolContractsConfig{
				USDCPools: map[uint64]v1_5_1.DeployUSDCTokenPoolInput{
					5009297550715157269: v1_5_1.DeployUSDCTokenPoolInput{},
				},
			},
			ErrStr: "does not exist in environment",
		},
		{
			Msg: "Missing router",
			Input: v1_5_1.DeployUSDCTokenPoolContractsConfig{
				USDCPools: map[uint64]v1_5_1.DeployUSDCTokenPoolInput{
					e.AllChainSelectors()[0]: v1_5_1.DeployUSDCTokenPoolInput{},
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

func TestValidateDeployUSDCTokenPoolInput(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})
	selector := e.AllChainSelectors()[0]
	chain := e.Chains[selector]
	addressBook := deployment.NewMemoryAddressBook()

	usdcToken, tokenMessenger := deployUSDCPrerequisites(t, lggr, chain, addressBook)

	nonUsdcToken, err := cldf.DeployContract(e.Logger, chain, addressBook,
		func(chain deployment.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
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
				Tv:       deployment.NewTypeAndVersion(changeset.USDCTokenPool, deployment.Version1_5_1),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	tests := []struct {
		Msg    string
		Input  v1_5_1.DeployUSDCTokenPoolInput
		ErrStr string
	}{
		{
			Msg:    "Missing token address",
			Input:  v1_5_1.DeployUSDCTokenPoolInput{},
			ErrStr: "token address must be defined",
		},
		{
			Msg: "Missing token messenger",
			Input: v1_5_1.DeployUSDCTokenPoolInput{
				TokenAddress: utils.RandomAddress(),
			},
			ErrStr: "token messenger must be defined",
		},
		{
			Msg: "Can't reach token",
			Input: v1_5_1.DeployUSDCTokenPoolInput{
				TokenAddress:   utils.RandomAddress(),
				TokenMessenger: utils.RandomAddress(),
			},
			ErrStr: "failed to fetch symbol from token",
		},
		{
			Msg: "Symbol is wrong",
			Input: v1_5_1.DeployUSDCTokenPoolInput{
				TokenAddress:   nonUsdcToken.Address,
				TokenMessenger: utils.RandomAddress(),
			},
			ErrStr: "is not USDC",
		},
		{
			Msg: "Can't reach token messenger",
			Input: v1_5_1.DeployUSDCTokenPoolInput{
				TokenAddress:   usdcToken.Address,
				TokenMessenger: utils.RandomAddress(),
			},
			ErrStr: "failed to fetch local message transmitter from address",
		},
		{
			Msg: "No error",
			Input: v1_5_1.DeployUSDCTokenPoolInput{
				TokenAddress:   usdcToken.Address,
				TokenMessenger: tokenMessenger.Address,
			},
			ErrStr: "",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.Input.Validate(e.GetContext(), chain, state.Chains[selector])
			if test.ErrStr != "" {
				require.Contains(t, err.Error(), test.ErrStr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDeployUSDCTokenPoolContracts(t *testing.T) {
	t.Parallel()

	for _, numRuns := range []int{1, 2} {
		t.Run(fmt.Sprintf("Run deployment %d time(s)", numRuns), func(t *testing.T) {
			lggr := logger.TestLogger(t)
			e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
				Chains: 2,
			})
			selectors := e.AllChainSelectors()

			addressBook := deployment.NewMemoryAddressBook()
			prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, len(selectors))
			for i, selector := range selectors {
				prereqCfg[i] = changeset.DeployPrerequisiteConfigPerChain{
					ChainSelector: selector,
				}
			}

			e, err := commoncs.Apply(t, e, nil,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
					changeset.DeployPrerequisiteConfig{
						Configs: prereqCfg,
					},
				),
			)
			require.NoError(t, err)

			newUSDCTokenPools := make(map[uint64]v1_5_1.DeployUSDCTokenPoolInput, len(selectors))
			for _, selector := range selectors {
				usdcToken, tokenMessenger := deployUSDCPrerequisites(t, lggr, e.Chains[selector], addressBook)

				newUSDCTokenPools[selector] = v1_5_1.DeployUSDCTokenPoolInput{
					TokenAddress:   usdcToken.Address,
					TokenMessenger: tokenMessenger.Address,
				}
			}

			for i := range numRuns {
				e, err = commoncs.Apply(t, e, nil,
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(v1_5_1.DeployUSDCTokenPoolContractsChangeset),
						v1_5_1.DeployUSDCTokenPoolContractsConfig{
							USDCPools: newUSDCTokenPools,
						},
					),
				)

				if i > 0 {
					require.ErrorContains(t, err, "already exists")
				} else {
					require.NoError(t, err)

					state, err := changeset.LoadOnchainState(e)
					require.NoError(t, err)

					for _, selector := range selectors {
						usdcTokenPools := state.Chains[selector].USDCTokenPools
						require.Len(t, usdcTokenPools, 1)
						owner, err := usdcTokenPools[deployment.Version1_5_1].Owner(nil)
						require.NoError(t, err)
						require.Equal(t, e.Chains[selector].DeployerKey.From, owner)
					}
				}
			}
		})
	}
}
