package v1_5_1_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestDeployTokenPoolFactoryChangeset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Msg                     string
		ForgetPrerequisites     bool
		MultipleRegistryModules bool
		ExpectedErr             string
		ConfigFn                func(selectors []uint64, state changeset.CCIPOnChainState) v1_5_1.DeployTokenPoolFactoryConfig
	}{
		{
			Msg: "should deploy token pool factory on all chains",
			ConfigFn: func(selectors []uint64, state changeset.CCIPOnChainState) v1_5_1.DeployTokenPoolFactoryConfig {
				return v1_5_1.DeployTokenPoolFactoryConfig{
					Chains: selectors,
				}
			},
		},
		{
			Msg:                 "should fail to deploy due to missing prereqs",
			ForgetPrerequisites: true,
			ConfigFn: func(selectors []uint64, state changeset.CCIPOnChainState) v1_5_1.DeployTokenPoolFactoryConfig {
				return v1_5_1.DeployTokenPoolFactoryConfig{
					Chains: selectors,
				}
			},
			ExpectedErr: "token admin registry does not exist",
		},
		{
			Msg:                     "should fail to deploy due to multiple registry modules",
			MultipleRegistryModules: true,
			ConfigFn: func(selectors []uint64, state changeset.CCIPOnChainState) v1_5_1.DeployTokenPoolFactoryConfig {
				return v1_5_1.DeployTokenPoolFactoryConfig{
					Chains: selectors,
				}
			},
			ExpectedErr: "multiple registry modules with version 1.6.0 exist",
		},
		{
			Msg:                     "should fail when a registry module is specified incorrectly",
			MultipleRegistryModules: true,
			ConfigFn: func(selectors []uint64, state changeset.CCIPOnChainState) v1_5_1.DeployTokenPoolFactoryConfig {
				addresses := make(map[uint64]common.Address, len(selectors))
				for _, selector := range selectors {
					addresses[selector] = utils.RandomAddress()
				}
				return v1_5_1.DeployTokenPoolFactoryConfig{
					Chains:                     selectors,
					RegistryModule1_6Addresses: addresses,
				}
			},
			ExpectedErr: "no registry module with version 1.6.0 and address",
		},
		{
			Msg:                     "should successfully deploy when a registry module is specified",
			MultipleRegistryModules: true,
			ConfigFn: func(selectors []uint64, state changeset.CCIPOnChainState) v1_5_1.DeployTokenPoolFactoryConfig {
				addresses := make(map[uint64]common.Address, len(selectors))
				for _, selector := range selectors {
					addresses[selector] = state.Chains[selector].RegistryModules1_6[0].Address()
				}
				return v1_5_1.DeployTokenPoolFactoryConfig{
					Chains:                     selectors,
					RegistryModule1_6Addresses: addresses,
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
				testCfg.Chains = 2
				testCfg.PrerequisiteDeploymentOnly = true
			})
			e := deployedEnvironment.Env
			selectors := e.AllChainSelectors()

			state, err := changeset.LoadOnchainState(e)
			require.NoError(t, err, "failed to load onchain state")

			if test.MultipleRegistryModules {
				// Add a new registry module to each chain
				for _, selector := range selectors {
					_, err := deployment.DeployContract(e.Logger, e.Chains[selector], e.ExistingAddresses,
						func(chain deployment.Chain) deployment.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom] {
							regModAddr, tx2, regMod, err2 := registry_module_owner_custom.DeployRegistryModuleOwnerCustom(
								chain.DeployerKey,
								chain.Client,
								state.Chains[selector].TokenAdminRegistry.Address())
							return deployment.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom]{
								Address: regModAddr, Contract: regMod, Tx: tx2, Tv: deployment.NewTypeAndVersion(changeset.RegistryModule, deployment.Version1_6_0), Err: err2,
							}
						})
					require.NoError(t, err, "failed to deploy registry module")
				}
			}
			if test.ForgetPrerequisites {
				// Clear the address book
				e.ExistingAddresses = deployment.NewMemoryAddressBook()
			}

			state, err = changeset.LoadOnchainState(e)
			require.NoError(t, err, "failed to load onchain state")

			e, err = commonchangeset.Apply(t, e, nil, commonchangeset.Configure(
				v1_5_1.DeployTokenPoolFactoryChangeset,
				test.ConfigFn(selectors, state),
			))
			if test.ExpectedErr != "" {
				require.ErrorContains(t, err, test.ExpectedErr, "expected error not found")
				return
			}
			require.NoError(t, err, "failed to apply DeployTokenPoolFactoryChangeset")

			state, err = changeset.LoadOnchainState(e)
			require.NoError(t, err, "failed to load onchain state")

			for _, chainSel := range selectors {
				tpf := state.Chains[chainSel].TokenPoolFactory
				require.NotNil(t, tpf, "token pool factory should be deployed on chain %d", chainSel)
				typeAndVersion, err := tpf.TypeAndVersion(nil)
				require.NoError(t, err, "failed to get type and version of token pool factory on chain %d", chainSel)
				require.Equal(t, "TokenPoolFactory 1.5.1", typeAndVersion, "unexpected type and version of token pool factory on chain %d", chainSel)
			}

			// IDEMPOTENCY CHECK
			e, err = commonchangeset.Apply(t, e, nil, commonchangeset.Configure(
				v1_5_1.DeployTokenPoolFactoryChangeset,
				v1_5_1.DeployTokenPoolFactoryConfig{
					Chains: selectors,
				},
			))
			require.ErrorContains(t, err, "token pool factory already deployed")
		})
	}
}
