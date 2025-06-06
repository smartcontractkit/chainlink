package v1_5_1_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_0_0/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployTokenPoolFactoryChangeset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Msg                     string
		ForgetPrerequisites     bool
		MultipleRegistryModules bool
		ExpectedErr             string
		ConfigFn                func(selectors []uint64, state stateview.CCIPOnChainState) v1_5_1.DeployTokenPoolFactoryConfig
	}{
		{
			Msg: "should deploy token pool factory on all chains",
			ConfigFn: func(selectors []uint64, state stateview.CCIPOnChainState) v1_5_1.DeployTokenPoolFactoryConfig {
				return v1_5_1.DeployTokenPoolFactoryConfig{
					Chains: selectors,
				}
			},
		},
		{
			Msg:                 "should fail to deploy due to missing prereqs",
			ForgetPrerequisites: true,
			ConfigFn: func(selectors []uint64, state stateview.CCIPOnChainState) v1_5_1.DeployTokenPoolFactoryConfig {
				return v1_5_1.DeployTokenPoolFactoryConfig{
					Chains: selectors,
				}
			},
			ExpectedErr: "token admin registry does not exist",
		},
		{
			Msg:                     "should fail to deploy due to multiple registry modules",
			MultipleRegistryModules: true,
			ConfigFn: func(selectors []uint64, state stateview.CCIPOnChainState) v1_5_1.DeployTokenPoolFactoryConfig {
				return v1_5_1.DeployTokenPoolFactoryConfig{
					Chains: selectors,
				}
			},
			ExpectedErr: "multiple registry modules with version 1.6.0 exist",
		},
		{
			Msg:                     "should fail when a registry module is specified incorrectly",
			MultipleRegistryModules: true,
			ConfigFn: func(selectors []uint64, state stateview.CCIPOnChainState) v1_5_1.DeployTokenPoolFactoryConfig {
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
			ConfigFn: func(selectors []uint64, state stateview.CCIPOnChainState) v1_5_1.DeployTokenPoolFactoryConfig {
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
			lggr := logger.TestLogger(t)
			e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
				Chains: 2,
			})
			selectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))

			if !test.ForgetPrerequisites {
				// NOTE: We don't use the DeployPrerequisites changeset because the TokenPoolFactory is a prerequisite in itself.
				for _, selector := range selectors {
					// Deploy token admin registry
					tokenAdminRegistry, err := cldf.DeployContract(e.Logger, e.BlockChains.EVMChains()[selector], e.ExistingAddresses,
						func(chain cldf_evm.Chain) cldf.ContractDeploy[*token_admin_registry.TokenAdminRegistry] {
							tokenAdminRegistryAddr, tx2, tokenAdminRegistry, err2 := token_admin_registry.DeployTokenAdminRegistry(
								chain.DeployerKey,
								chain.Client)
							return cldf.ContractDeploy[*token_admin_registry.TokenAdminRegistry]{
								Address: tokenAdminRegistryAddr, Contract: tokenAdminRegistry, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.TokenAdminRegistry, deployment.Version1_5_0), Err: err2,
							}
						})
					require.NoError(t, err, "failed to deploy token admin registry")
					// Deploy RMN proxy
					rmnProxy, err := cldf.DeployContract(lggr, e.BlockChains.EVMChains()[selector], e.ExistingAddresses,
						func(chain cldf_evm.Chain) cldf.ContractDeploy[*rmn_proxy_contract.RMNProxy] {
							rmnProxyAddr, tx2, rmnProxy2, err2 := rmn_proxy_contract.DeployRMNProxy(
								chain.DeployerKey,
								chain.Client,
								// We don't need the actual RMN address here,
								// just a random address to satisfy the constructor.
								utils.RandomAddress(),
							)
							return cldf.ContractDeploy[*rmn_proxy_contract.RMNProxy]{
								Address: rmnProxyAddr, Contract: rmnProxy2, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.ARMProxy, deployment.Version1_0_0), Err: err2,
							}
						})
					require.NoError(t, err, "failed to deploy RMN proxy")
					// Deploy router
					_, err = cldf.DeployContract(e.Logger, e.BlockChains.EVMChains()[selector], e.ExistingAddresses,
						func(chain cldf_evm.Chain) cldf.ContractDeploy[*router.Router] {
							routerAddr, tx2, routerC, err2 := router.DeployRouter(
								chain.DeployerKey,
								chain.Client,
								// We don't need the actual addresses here, just some random ones
								// to satisfy the constructor.
								utils.RandomAddress(),
								rmnProxy.Address,
							)
							return cldf.ContractDeploy[*router.Router]{
								Address: routerAddr, Contract: routerC, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.Router, deployment.Version1_2_0), Err: err2,
							}
						})
					require.NoError(t, err, "failed to deploy router")
					// Deploy registry module
					_, err = cldf.DeployContract(e.Logger, e.BlockChains.EVMChains()[selector], e.ExistingAddresses,
						func(chain cldf_evm.Chain) cldf.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom] {
							regModAddr, tx2, regMod, err2 := registry_module_owner_custom.DeployRegistryModuleOwnerCustom(
								chain.DeployerKey,
								chain.Client,
								tokenAdminRegistry.Address,
							)
							return cldf.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom]{
								Address: regModAddr, Contract: regMod, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.RegistryModule, deployment.Version1_6_0), Err: err2,
							}
						})
					require.NoError(t, err, "failed to deploy registry module")
				}
			}

			if test.MultipleRegistryModules {
				// Add a new registry module to each chain
				state, err := stateview.LoadOnchainState(e)
				require.NoError(t, err, "failed to load onchain state")
				for _, selector := range selectors {
					_, err := cldf.DeployContract(e.Logger, e.BlockChains.EVMChains()[selector], e.ExistingAddresses,
						func(chain cldf_evm.Chain) cldf.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom] {
							regModAddr, tx2, regMod, err2 := registry_module_owner_custom.DeployRegistryModuleOwnerCustom(
								chain.DeployerKey,
								chain.Client,
								state.Chains[selector].TokenAdminRegistry.Address())
							return cldf.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom]{
								Address: regModAddr, Contract: regMod, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.RegistryModule, deployment.Version1_6_0), Err: err2,
							}
						})
					require.NoError(t, err, "failed to deploy registry module")
				}
			}

			state, err := stateview.LoadOnchainState(e)
			require.NoError(t, err, "failed to load onchain state")

			e, err = commonchangeset.Apply(t, e, commonchangeset.Configure(
				v1_5_1.DeployTokenPoolFactoryChangeset,
				test.ConfigFn(selectors, state),
			))
			if test.ExpectedErr != "" {
				require.ErrorContains(t, err, test.ExpectedErr, "expected error not found")
				return
			}
			require.NoError(t, err, "failed to apply DeployTokenPoolFactoryChangeset")

			state, err = stateview.LoadOnchainState(e)
			require.NoError(t, err, "failed to load onchain state")

			for _, chainSel := range selectors {
				tpf := state.Chains[chainSel].TokenPoolFactory
				require.NotNil(t, tpf, "token pool factory should be deployed on chain %d", chainSel)
				typeAndVersion, err := tpf.TypeAndVersion(nil)
				require.NoError(t, err, "failed to get type and version of token pool factory on chain %d", chainSel)
				require.Equal(t, "TokenPoolFactory 1.5.1", typeAndVersion, "unexpected type and version of token pool factory on chain %d", chainSel)
			}

			// IDEMPOTENCY CHECK
			e, err = commonchangeset.Apply(t, e, commonchangeset.Configure(
				v1_5_1.DeployTokenPoolFactoryChangeset,
				v1_5_1.DeployTokenPoolFactoryConfig{
					Chains: selectors,
				},
			))
			require.ErrorContains(t, err, "token pool factory already deployed")
		})
	}
}
