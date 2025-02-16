package changeset_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestValidateDeployUSDCTokenPoolContractsConfig(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})

	tests := []struct {
		Msg    string
		Input  changeset.DeployUSDCTokenPoolContractsConfig
		ErrStr string
	}{
		{
			Msg: "Chain selector is not valid",
			Input: changeset.DeployUSDCTokenPoolContractsConfig{
				USDCPools: map[uint64]changeset.DeployUSDCTokenPoolInput{
					0: changeset.DeployUSDCTokenPoolInput{},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Input: changeset.DeployUSDCTokenPoolContractsConfig{
				USDCPools: map[uint64]changeset.DeployUSDCTokenPoolInput{
					5009297550715157269: changeset.DeployUSDCTokenPoolInput{},
				},
			},
			ErrStr: "does not exist in environment",
		},
		{
			Msg: "Missing router",
			Input: changeset.DeployUSDCTokenPoolContractsConfig{
				USDCPools: map[uint64]changeset.DeployUSDCTokenPoolInput{
					e.AllChainSelectors()[0]: changeset.DeployUSDCTokenPoolInput{},
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

	usdcToken, tokenMessenger := testhelpers.DeployUSDCPrerequisites(t, lggr, chain, addressBook)

	nonUsdcToken, err := deployment.DeployContract(e.Logger, chain, addressBook,
		func(chain deployment.Chain) deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
				chain.DeployerKey,
				chain.Client,
				"NOTUSDC",
				"NOTUSDC",
				6,
				big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
			)
			return deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
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
		Input  changeset.DeployUSDCTokenPoolInput
		ErrStr string
	}{
		{
			Msg:    "Missing token address",
			Input:  changeset.DeployUSDCTokenPoolInput{},
			ErrStr: "token address must be defined",
		},
		{
			Msg: "Missing token messenger",
			Input: changeset.DeployUSDCTokenPoolInput{
				TokenAddress: utils.RandomAddress(),
			},
			ErrStr: "token messenger must be defined",
		},
		{
			Msg: "Can't reach token",
			Input: changeset.DeployUSDCTokenPoolInput{
				TokenAddress:   utils.RandomAddress(),
				TokenMessenger: utils.RandomAddress(),
			},
			ErrStr: "failed to fetch symbol from token",
		},
		{
			Msg: "Symbol is wrong",
			Input: changeset.DeployUSDCTokenPoolInput{
				TokenAddress:   nonUsdcToken.Address,
				TokenMessenger: utils.RandomAddress(),
			},
			ErrStr: "is not USDC",
		},
		{
			Msg: "Can't reach token messenger",
			Input: changeset.DeployUSDCTokenPoolInput{
				TokenAddress:   usdcToken.Address,
				TokenMessenger: utils.RandomAddress(),
			},
			ErrStr: "failed to fetch local message transmitter from address",
		},
		{
			Msg: "No error",
			Input: changeset.DeployUSDCTokenPoolInput{
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
					deployment.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
					changeset.DeployPrerequisiteConfig{
						Configs: prereqCfg,
					},
				),
			)
			require.NoError(t, err)

			newUSDCTokenPools := make(map[uint64]changeset.DeployUSDCTokenPoolInput, len(selectors))
			for _, selector := range selectors {
				usdcToken, tokenMessenger := testhelpers.DeployUSDCPrerequisites(t, lggr, e.Chains[selector], addressBook)

				newUSDCTokenPools[selector] = changeset.DeployUSDCTokenPoolInput{
					TokenAddress:   usdcToken.Address,
					TokenMessenger: tokenMessenger.Address,
				}
			}

			for i := range numRuns {
				e, err = commoncs.Apply(t, e, nil,
					commonchangeset.Configure(
						deployment.CreateLegacyChangeSet(changeset.DeployUSDCTokenPoolContractsChangeset),
						changeset.DeployUSDCTokenPoolContractsConfig{
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
