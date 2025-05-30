package changeset

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/factory_burn_mint_erc20"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/token_pool_factory"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_0_0/rmn_proxy_contract"
	price_registry_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/price_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/usdc_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/multicall3"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/weth9"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/weth9_zksync"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var (
	_ cldf.ChangeSet[DeployPrerequisiteConfig] = DeployPrerequisitesChangeset
)

// DeployPrerequisitesChangeset deploys the pre-requisite contracts for CCIP
// pre-requisite contracts are the contracts which can be reused from previous versions of CCIP
// Or the contracts which are already deployed on the chain ( for example, tokens, feeds, etc)
// Caller should update the environment's address book with the returned addresses.
func DeployPrerequisitesChangeset(env cldf.Environment, cfg DeployPrerequisiteConfig) (cldf.ChangesetOutput, error) {
	err := cfg.Validate()
	if err != nil {
		return cldf.ChangesetOutput{}, errors.Wrapf(cldf.ErrInvalidConfig, "%v", err)
	}
	ab := cldf.NewMemoryAddressBook()
	err = deployPrerequisiteChainContracts(env, ab, cfg)
	if err != nil {
		env.Logger.Errorw("Failed to deploy prerequisite contracts", "err", err, "addressBook", ab)
		return cldf.ChangesetOutput{
			AddressBook: ab,
		}, fmt.Errorf("failed to deploy prerequisite contracts: %w", err)
	}
	return cldf.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

type DeployPrerequisiteContractsOpts struct {
	USDCEnabled             bool
	Multicall3Enabled       bool
	TokenPoolFactoryEnabled bool
	LegacyDeploymentCfg     *V1_5DeploymentConfig
}

type V1_5DeploymentConfig struct {
	RMNConfig                  *rmn_contract.RMNConfig
	PriceRegStalenessThreshold uint32
}

type DeployPrerequisiteConfig struct {
	Configs []DeployPrerequisiteConfigPerChain
}

type DeployPrerequisiteConfigPerChain struct {
	ChainSelector uint64
	Opts          []PrerequisiteOpt
}

func (c DeployPrerequisiteConfig) Validate() error {
	mapAllChainSelectors := make(map[uint64]struct{})
	for _, cfg := range c.Configs {
		cs := cfg.ChainSelector
		mapAllChainSelectors[cs] = struct{}{}
		if err := cldf.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
	}
	return nil
}

type PrerequisiteOpt func(o *DeployPrerequisiteContractsOpts)

func WithTokenPoolFactoryEnabled() PrerequisiteOpt {
	return func(o *DeployPrerequisiteContractsOpts) {
		o.TokenPoolFactoryEnabled = true
	}
}

func WithUSDCEnabled() PrerequisiteOpt {
	return func(o *DeployPrerequisiteContractsOpts) {
		o.USDCEnabled = true
	}
}

func WithMultiCall3Enabled() PrerequisiteOpt {
	return func(o *DeployPrerequisiteContractsOpts) {
		o.Multicall3Enabled = true
	}
}

func WithLegacyDeploymentEnabled(cfg V1_5DeploymentConfig) PrerequisiteOpt {
	return func(o *DeployPrerequisiteContractsOpts) {
		if cfg.PriceRegStalenessThreshold == 0 {
			panic("PriceRegStalenessThreshold must be set")
		}
		// TODO validate RMNConfig
		o.LegacyDeploymentCfg = &cfg
	}
}

func deployPrerequisiteChainContracts(e cldf.Environment, ab cldf.AddressBook, cfg DeployPrerequisiteConfig) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return err
	}
	deployGrp := errgroup.Group{}
	for _, c := range cfg.Configs {
		chain := e.BlockChains.EVMChains()[c.ChainSelector]
		deployGrp.Go(func() error {
			err := deployPrerequisiteContracts(e, ab, state, chain, c.Opts...)
			if err != nil {
				e.Logger.Errorw("Failed to deploy prerequisite contracts", "chain", chain.String(), "err", err)
				return err
			}
			return nil
		})
	}
	return deployGrp.Wait()
}

// deployPrerequisiteContracts deploys the contracts that can be ported from previous CCIP version to the new one.
// This is only required for staging and test environments where the contracts are not already deployed.
func deployPrerequisiteContracts(e cldf.Environment, ab cldf.AddressBook, state stateview.CCIPOnChainState, chain cldf_evm.Chain, opts ...PrerequisiteOpt) error {
	deployOpts := &DeployPrerequisiteContractsOpts{}
	for _, opt := range opts {
		if opt != nil {
			opt(deployOpts)
		}
	}
	lggr := e.Logger
	chainState, chainExists := state.EVMChainState(chain.Selector)
	var weth9Contract *weth9.WETH9
	var tokenAdminReg *token_admin_registry.TokenAdminRegistry
	var tokenPoolFactory *token_pool_factory.TokenPoolFactory
	var factoryBurnMintERC20 *factory_burn_mint_erc20.FactoryBurnMintERC20
	var rmnProxy *rmn_proxy_contract.RMNProxy
	var r *router.Router
	var mc3 *multicall3.Multicall3
	if chainExists {
		weth9Contract = chainState.Weth9
		tokenAdminReg = chainState.TokenAdminRegistry
		tokenPoolFactory = chainState.TokenPoolFactory
		factoryBurnMintERC20 = chainState.FactoryBurnMintERC20Token
		rmnProxy = chainState.RMNProxy
		r = chainState.Router
		mc3 = chainState.Multicall3
	}
	var rmnAddr common.Address
	// if we are setting up 1.5 version, deploy RMN contract based on the config provided
	// else deploy the mock RMN contract
	switch {
	// if RMN is found in state use that
	case chainState.RMN != nil && chainState.RMN.Address() != (common.Address{}):
		lggr.Infow("RMN already deployed", "chain", chain.String(), "address", chainState.RMN.Address)
		rmnAddr = chainState.RMN.Address()
	// if RMN is not found in state and LegacyDeploymentCfg is provided, deploy RMN contract based on the config
	case deployOpts.LegacyDeploymentCfg != nil && deployOpts.LegacyDeploymentCfg.RMNConfig != nil:
		rmn, err := cldf.DeployContract(lggr, chain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*rmn_contract.RMNContract] {
				var (
					rmnAddress common.Address
					tx2        *types.Transaction
					rmnC       *rmn_contract.RMNContract
					err2       error
				)

				if chain.IsZkSyncVM {
					rmnAddress, _, rmnC, err2 = rmn_contract.DeployRMNContractZk(
						nil,
						chain.ClientZkSyncVM,
						chain.DeployerKeyZkSyncVM,
						chain.Client,
						*deployOpts.LegacyDeploymentCfg.RMNConfig,
					)
				} else {
					rmnAddress, tx2, rmnC, err2 = rmn_contract.DeployRMNContract(
						chain.DeployerKey,
						chain.Client,
						*deployOpts.LegacyDeploymentCfg.RMNConfig,
					)
				}
				return cldf.ContractDeploy[*rmn_contract.RMNContract]{
					Address: rmnAddress, Contract: rmnC, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.RMN, deployment.Version1_5_0), Err: err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy RMN", "chain", chain.String(), "err", cldf.MaybeDataErr(err))
			return err
		}
		rmnAddr = rmn.Address
	default:
		// otherwise deploy the mock RMN contract
		if chainState.MockRMN == nil {
			rmn, err := cldf.DeployContract(lggr, chain, ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*mock_rmn_contract.MockRMNContract] {
					var (
						rmnAddress common.Address
						tx2        *types.Transaction
						rmnC       *mock_rmn_contract.MockRMNContract
						err2       error
					)
					if chain.IsZkSyncVM {
						rmnAddress, _, rmnC, err2 = mock_rmn_contract.DeployMockRMNContractZk(
							nil,
							chain.ClientZkSyncVM,
							chain.DeployerKeyZkSyncVM,
							chain.Client,
						)
					} else {
						rmnAddress, tx2, rmnC, err2 = mock_rmn_contract.DeployMockRMNContract(
							chain.DeployerKey,
							chain.Client,
						)
					}
					return cldf.ContractDeploy[*mock_rmn_contract.MockRMNContract]{
						Address: rmnAddress, Contract: rmnC, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.MockRMN, deployment.Version1_0_0), Err: err2,
					}
				})
			if err != nil {
				lggr.Errorw("Failed to deploy mock RMN", "chain", chain.String(), "err", err)
				return err
			}
			rmnAddr = rmn.Address
		} else {
			lggr.Infow("Mock RMN already deployed", "chain", chain.String(), "addr", chainState.MockRMN.Address)
			rmnAddr = chainState.MockRMN.Address()
		}
	}
	if rmnProxy == nil {
		RMNProxy, err := cldf.DeployContract(lggr, chain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*rmn_proxy_contract.RMNProxy] {
				var (
					rmnProxyAddr common.Address
					tx2          *types.Transaction
					rmnProxy2    *rmn_proxy_contract.RMNProxy
					err2         error
				)
				if chain.IsZkSyncVM {
					rmnProxyAddr, _, rmnProxy2, err2 = rmn_proxy_contract.DeployRMNProxyZk(
						nil,
						chain.ClientZkSyncVM,
						chain.DeployerKeyZkSyncVM,
						chain.Client,
						rmnAddr,
					)
				} else {
					rmnProxyAddr, tx2, rmnProxy2, err2 = rmn_proxy_contract.DeployRMNProxy(
						chain.DeployerKey,
						chain.Client,
						rmnAddr,
					)
				}
				return cldf.ContractDeploy[*rmn_proxy_contract.RMNProxy]{
					Address: rmnProxyAddr, Contract: rmnProxy2, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.ARMProxy, deployment.Version1_0_0), Err: err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy RMNProxy", "chain", chain.String(), "err", err)
			return err
		}
		rmnProxy = RMNProxy.Contract
	} else {
		lggr.Infow("RMNProxy already deployed", "chain", chain.String(), "addr", rmnProxy.Address)
		// check if the RMNProxy is pointing to the correct RMN contract
		currentRMNAddr, err := rmnProxy.GetARM(nil)
		if err != nil {
			lggr.Errorw("Failed to get RMN from RMNProxy", "chain", chain.String(), "err", err)
			return err
		}
		if currentRMNAddr != rmnAddr {
			lggr.Infow("RMNProxy is not pointing to the correct RMN contract, updating RMN", "chain", chain.String(), "currentRMN", currentRMNAddr, "expectedRMN", rmnAddr)
			rmnOwner, err := rmnProxy.Owner(nil)
			if err != nil {
				lggr.Errorw("Failed to get owner of RMNProxy", "chain", chain.String(), "err", err)
				return err
			}
			if rmnOwner != chain.DeployerKey.From {
				lggr.Warnw(
					"RMNProxy is not owned by the deployer and RMNProxy is not pointing to the correct RMN contract, "+
						"run SetRMNRemoteOnRMNProxyChangeset to update RMN with a proposal",
					"chain", chain.String(), "owner", rmnOwner, "currentRMN", currentRMNAddr, "expectedRMN", rmnAddr)
			} else {
				tx, err := rmnProxy.SetARM(chain.DeployerKey, rmnAddr)
				if err != nil {
					lggr.Errorw("Failed to set RMN on RMNProxy", "chain", chain.String(), "err", err)
					return err
				}
				_, err = chain.Confirm(tx)
				if err != nil {
					lggr.Errorw("Failed to confirm setRMN on RMNProxy", "chain", chain.String(), "err", err)
					return err
				}
			}
		}
	}
	if tokenAdminReg == nil {
		tokenAdminRegistry, err := cldf.DeployContract(e.Logger, chain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*token_admin_registry.TokenAdminRegistry] {
				var (
					tokenAdminRegistryAddr common.Address
					tx2                    *types.Transaction
					tokenAdminRegistry     *token_admin_registry.TokenAdminRegistry
					err2                   error
				)
				if chain.IsZkSyncVM {
					tokenAdminRegistryAddr, _, tokenAdminRegistry, err2 = token_admin_registry.DeployTokenAdminRegistryZk(
						nil,
						chain.ClientZkSyncVM,
						chain.DeployerKeyZkSyncVM,
						chain.Client,
					)
				} else {
					tokenAdminRegistryAddr, tx2, tokenAdminRegistry, err2 = token_admin_registry.DeployTokenAdminRegistry(
						chain.DeployerKey,
						chain.Client)
				}
				return cldf.ContractDeploy[*token_admin_registry.TokenAdminRegistry]{
					Address: tokenAdminRegistryAddr, Contract: tokenAdminRegistry, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.TokenAdminRegistry, deployment.Version1_5_0), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy token admin registry", "chain", chain.String(), "err", err)
			return err
		}
		tokenAdminReg = tokenAdminRegistry.Contract
	} else {
		e.Logger.Infow("tokenAdminRegistry already deployed", "chain", chain.String(), "addr", tokenAdminReg.Address)
	}
	// fetch addresses of both version of the registry module
	var regAddresses []common.Address
	for _, reg := range chainState.RegistryModules1_6 {
		regAddresses = append(regAddresses, reg.Address())
	}
	for _, reg := range chainState.RegistryModules1_5 {
		regAddresses = append(regAddresses, reg.Address())
	}
	if len(regAddresses) == 0 {
		customRegistryModule, err := cldf.DeployContract(e.Logger, chain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom] {
				var (
					regModAddr common.Address
					tx2        *types.Transaction
					regMod     *registry_module_owner_custom.RegistryModuleOwnerCustom
					err2       error
				)
				if chain.IsZkSyncVM {
					regModAddr, _, regMod, err2 = registry_module_owner_custom.DeployRegistryModuleOwnerCustomZk(
						nil,
						chain.ClientZkSyncVM,
						chain.DeployerKeyZkSyncVM,
						chain.Client,
						tokenAdminReg.Address(),
					)
				} else {
					regModAddr, tx2, regMod, err2 = registry_module_owner_custom.DeployRegistryModuleOwnerCustom(
						chain.DeployerKey,
						chain.Client,
						tokenAdminReg.Address())
				}
				return cldf.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom]{
					Address: regModAddr, Contract: regMod, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.RegistryModule, deployment.Version1_6_0), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy custom registry module", "chain", chain.String(), "err", err)
			return err
		}
		regAddresses = append(regAddresses, customRegistryModule.Address)
		e.Logger.Infow("deployed custom registry module", "chain", chain.String(), "addr", customRegistryModule.Address)
	} else {
		e.Logger.Infow("custom registry module already deployed", "chain", chain.String(), "addr", regAddresses)
	}
	for _, reg := range regAddresses {
		isRegistryAdded, err := tokenAdminReg.IsRegistryModule(nil, reg)
		if err != nil {
			e.Logger.Errorw("Failed to check if registry module is added on token admin registry", "chain", chain.String(), "err", err)
			return fmt.Errorf("failed to check if registry module is added on token admin registry: %w", err)
		}
		if !isRegistryAdded {
			owner, err := tokenAdminReg.Owner(nil)
			if err != nil {
				e.Logger.Errorw("Failed to get owner of token admin registry", "chain", chain.String(), "err", err)
				return fmt.Errorf("failed to get owner of token admin registry: %w", err)
			}
			if owner != chain.DeployerKey.From {
				e.Logger.Errorw("Owner is not deployer key, cannot add registry module", "chain", chain.String(), "owner", owner)
				return fmt.Errorf("owner %s is not deployer key, cannot add registry module", owner)
			}
			tx, err := tokenAdminReg.AddRegistryModule(chain.DeployerKey, reg)
			if err != nil {
				e.Logger.Errorw("Failed to assign registry module on token admin registry", "chain", chain.String(), "err", err)
				return fmt.Errorf("failed to assign registry module on token admin registry: %w", err)
			}

			_, err = chain.Confirm(tx)
			if err != nil {
				e.Logger.Errorw("Failed to confirm assign registry module on token admin registry", "chain", chain.String(), "err", err)
				return fmt.Errorf("failed to confirm assign registry module on token admin registry: %w", err)
			}
			e.Logger.Infow("assigned registry module on token admin registry")
		}
	}

	if weth9Contract == nil {
		deployWeth9ZkAndPort := func(chain cldf_evm.Chain) (*weth9.WETH9, common.Address, error) {
			weth9AddrZk, _, weth9zk, err := weth9_zksync.DeployWETH9ZKSyncZk(
				nil,
				chain.ClientZkSyncVM,
				chain.DeployerKeyZkSyncVM,
				chain.Client,
			)
			if err != nil {
				return nil, common.Address{}, err
			}
			weth9ZkPorted, err := weth9.NewWETH9(weth9zk.Address(), chain.Client)
			if err != nil {
				return nil, common.Address{}, err
			}

			return weth9ZkPorted, weth9AddrZk, nil
		}

		weth, err := cldf.DeployContract(lggr, chain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*weth9.WETH9] {
				var (
					weth9Addr common.Address
					tx2       *types.Transaction
					weth9c    *weth9.WETH9
					err2      error
				)
				if chain.IsZkSyncVM {
					weth9c, weth9Addr, err2 = deployWeth9ZkAndPort(chain)
				} else {
					weth9Addr, tx2, weth9c, err2 = weth9.DeployWETH9(
						chain.DeployerKey,
						chain.Client,
					)
				}
				return cldf.ContractDeploy[*weth9.WETH9]{
					Address: weth9Addr, Contract: weth9c, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.WETH9, deployment.Version1_0_0), Err: err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy weth9", "chain", chain.String(), "err", err)
			return err
		}
		weth9Contract = weth.Contract
	} else {
		lggr.Infow("weth9 already deployed", "chain", chain.String(), "addr", weth9Contract.Address)
		weth9Contract = chainState.Weth9
	}

	// if router is not already deployed, we deploy it
	if r == nil {
		routerContract, err := cldf.DeployContract(e.Logger, chain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*router.Router] {
				var (
					routerAddr common.Address
					tx2        *types.Transaction
					routerC    *router.Router
					err2       error
				)
				if chain.IsZkSyncVM {
					routerAddr, _, routerC, err2 = router.DeployRouterZk(
						nil,
						chain.ClientZkSyncVM,
						chain.DeployerKeyZkSyncVM,
						chain.Client,
						weth9Contract.Address(),
						rmnProxy.Address(),
					)
				} else {
					routerAddr, tx2, routerC, err2 = router.DeployRouter(
						chain.DeployerKey,
						chain.Client,
						weth9Contract.Address(),
						rmnProxy.Address(),
					)
				}
				return cldf.ContractDeploy[*router.Router]{
					Address: routerAddr, Contract: routerC, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.Router, deployment.Version1_2_0), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy router", "chain", chain.String(), "err", err)
			return err
		}

		r = routerContract.Contract
	} else {
		e.Logger.Infow("router already deployed", "chain", chain.String(), "addr", chainState.Router.Address)
	}
	if deployOpts.TokenPoolFactoryEnabled {
		if tokenPoolFactory == nil {
			_, err := cldf.DeployContract(e.Logger, chain, ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*token_pool_factory.TokenPoolFactory] {
					var (
						tpfAddr  common.Address
						tx2      *types.Transaction
						contract *token_pool_factory.TokenPoolFactory
						err2     error
					)
					if chain.IsZkSyncVM {
						tpfAddr, _, contract, err2 = token_pool_factory.DeployTokenPoolFactoryZk(
							nil,
							chain.ClientZkSyncVM,
							chain.DeployerKeyZkSyncVM,
							chain.Client,
							tokenAdminReg.Address(),
							regAddresses[0],
							rmnProxy.Address(),
							r.Address(),
						)
					} else {
						tpfAddr, tx2, contract, err2 = token_pool_factory.DeployTokenPoolFactory(
							chain.DeployerKey,
							chain.Client,
							tokenAdminReg.Address(),
							// There will always be at least one registry module deployed at this point.
							// We just use the first one here. If a different RegistryModule is desired,
							// users can run DeployTokenPoolFactoryChangeset with the desired address.
							regAddresses[0],
							rmnProxy.Address(),
							r.Address(),
						)
					}
					return cldf.ContractDeploy[*token_pool_factory.TokenPoolFactory]{
						Address: tpfAddr, Contract: contract, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.TokenPoolFactory, deployment.Version1_5_1), Err: err2,
					}
				},
			)
			if err != nil {
				e.Logger.Errorw("Failed to deploy token pool factory", "chain", chain.String(), "err", err)
				return err
			}
		} else {
			e.Logger.Infow("Token pool factory already deployed", "chain", chain.String(), "addr", tokenPoolFactory.Address)
		}
		// FactoryBurnMintERC20 is a contract that gets deployed by the TokenPoolFactory.
		// We deploy it here so that we can verify it. All subsequent user deployments would then be verified.
		if factoryBurnMintERC20 == nil {
			_, err := cldf.DeployContract(e.Logger, chain, ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*factory_burn_mint_erc20.FactoryBurnMintERC20] {
					var (
						factoryBurnMintERC20Addr common.Address
						tx2                      *types.Transaction
						contract                 *factory_burn_mint_erc20.FactoryBurnMintERC20
						err2                     error
					)
					if chain.IsZkSyncVM {
						factoryBurnMintERC20Addr, _, contract, err2 = factory_burn_mint_erc20.DeployFactoryBurnMintERC20Zk(
							nil,
							chain.ClientZkSyncVM,
							chain.DeployerKeyZkSyncVM,
							chain.Client,
							string(shared.FactoryBurnMintERC20Symbol),
							string(shared.FactoryBurnMintERC20Symbol),
							uint8(18),
							big.NewInt(0),
							big.NewInt(0),
							chain.DeployerKey.From,
						)
					} else {
						factoryBurnMintERC20Addr, tx2, contract, err2 = factory_burn_mint_erc20.DeployFactoryBurnMintERC20(
							chain.DeployerKey,
							chain.Client,
							string(shared.FactoryBurnMintERC20Symbol),
							string(shared.FactoryBurnMintERC20Symbol),
							18,
							big.NewInt(0),
							big.NewInt(0),
							chain.DeployerKey.From,
						)
					}
					return cldf.ContractDeploy[*factory_burn_mint_erc20.FactoryBurnMintERC20]{
						Address: factoryBurnMintERC20Addr, Contract: contract, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.FactoryBurnMintERC20Token, deployment.Version1_0_0), Err: err2,
					}
				},
			)
			if err != nil {
				e.Logger.Errorw("Failed to deploy factory burn mint erc20", "chain", chain.String(), "err", err)
				return err
			}
		} else {
			e.Logger.Infow("factory burn mint erc20 already deployed", "chain", chain.String(), "addr", factoryBurnMintERC20.Address)
		}
	}
	if deployOpts.Multicall3Enabled && mc3 == nil {
		_, err := cldf.DeployContract(e.Logger, chain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*multicall3.Multicall3] {
				var (
					multicall3Addr    common.Address
					tx2               *types.Transaction
					multicall3Wrapper *multicall3.Multicall3
					err2              error
				)
				if chain.IsZkSyncVM {
					multicall3Addr, _, multicall3Wrapper, err2 = multicall3.DeployMulticall3Zk(
						nil,
						chain.ClientZkSyncVM,
						chain.DeployerKeyZkSyncVM,
						chain.Client,
					)
				} else {
					multicall3Addr, tx2, multicall3Wrapper, err2 = multicall3.DeployMulticall3(
						chain.DeployerKey,
						chain.Client,
					)
				}
				return cldf.ContractDeploy[*multicall3.Multicall3]{
					Address: multicall3Addr, Contract: multicall3Wrapper, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.Multicall3, deployment.Version1_0_0), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy ccip multicall", "chain", chain.String(), "err", err)
			return err
		}
	} else if mc3 != nil {
		e.Logger.Info("ccip multicall already deployed", "chain", chain.String(), "addr", mc3.Address)

	}
	if deployOpts.USDCEnabled {
		token, pool, messenger, transmitter, err1 := deployUSDC(e.Logger, chain, ab, rmnProxy.Address(), r.Address())
		if err1 != nil {
			return err1
		}
		e.Logger.Infow("Deployed USDC contracts",
			"chain", chain.String(),
			"token", token.Address(),
			"pool", pool.Address(),
			"transmitter", transmitter.Address(),
			"messenger", messenger.Address(),
		)
	}
	if chainState.Receiver == nil {
		_, err := cldf.DeployContract(e.Logger, chain, ab,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*maybe_revert_message_receiver.MaybeRevertMessageReceiver] {
				var (
					receiverAddr common.Address
					tx           *types.Transaction
					receiver     *maybe_revert_message_receiver.MaybeRevertMessageReceiver
					err2         error
				)
				if chain.IsZkSyncVM {
					receiverAddr, _, receiver, err2 = maybe_revert_message_receiver.DeployMaybeRevertMessageReceiverZk(
						nil,
						chain.ClientZkSyncVM,
						chain.DeployerKeyZkSyncVM,
						chain.Client,
						false,
					)
				} else {
					receiverAddr, tx, receiver, err2 = maybe_revert_message_receiver.DeployMaybeRevertMessageReceiver(
						chain.DeployerKey,
						chain.Client,
						false,
					)
				}
				return cldf.ContractDeploy[*maybe_revert_message_receiver.MaybeRevertMessageReceiver]{
					Address: receiverAddr, Contract: receiver, Tx: tx, Tv: cldf.NewTypeAndVersion(shared.CCIPReceiver, deployment.Version1_0_0), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy receiver", "chain", chain.String(), "err", err)
			return err
		}
	} else {
		e.Logger.Infow("receiver already deployed", "addr", chainState.Receiver.Address, "chain", chain.String())
	}
	// Only applicable if setting up for 1.5 version, remove this once we have fully migrated to 1.6
	if deployOpts.LegacyDeploymentCfg != nil {
		if chainState.PriceRegistry == nil {
			linkAddr, err1 := chainState.LinkTokenAddress()
			if err1 != nil {
				return fmt.Errorf("failed to get link token address for chain %s: %w", chain.String(), err1)
			}
			_, err := cldf.DeployContract(lggr, chain, ab,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*price_registry_1_2_0.PriceRegistry] {
					var (
						priceRegAddr  common.Address
						tx2           *types.Transaction
						priceRegAddrC *price_registry_1_2_0.PriceRegistry
						err2          error
					)
					if chain.IsZkSyncVM {
						priceRegAddr, _, priceRegAddrC, err2 = price_registry_1_2_0.DeployPriceRegistryZk(
							nil,
							chain.ClientZkSyncVM,
							chain.DeployerKeyZkSyncVM,
							chain.Client,
							nil,
							[]common.Address{weth9Contract.Address(), linkAddr},
							deployOpts.LegacyDeploymentCfg.PriceRegStalenessThreshold,
						)
					} else {
						priceRegAddr, tx2, priceRegAddrC, err2 = price_registry_1_2_0.DeployPriceRegistry(
							chain.DeployerKey,
							chain.Client,
							nil,
							[]common.Address{weth9Contract.Address(), linkAddr},
							deployOpts.LegacyDeploymentCfg.PriceRegStalenessThreshold,
						)
					}
					return cldf.ContractDeploy[*price_registry_1_2_0.PriceRegistry]{
						Address: priceRegAddr, Contract: priceRegAddrC, Tx: tx2,
						Tv: cldf.NewTypeAndVersion(shared.PriceRegistry, deployment.Version1_2_0), Err: err2,
					}
				})
			if err != nil {
				lggr.Errorw("Failed to deploy PriceRegistry", "chain", chain.String(), "err", err)
				return err
			}
		} else {
			lggr.Infow("PriceRegistry already deployed", "chain", chain.String(), "addr", chainState.PriceRegistry.Address)
		}
	}
	return nil
}

func deployUSDC(
	lggr logger.Logger,
	chain cldf_evm.Chain,
	addresses cldf.AddressBook,
	rmnProxy common.Address,
	router common.Address,
) (
	*burn_mint_erc677.BurnMintERC677,
	*usdc_token_pool.USDCTokenPool,
	*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger,
	*mock_usdc_token_transmitter.MockE2EUSDCTransmitter,
	error,
) {
	token, err := cldf.DeployContract(lggr, chain, addresses,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			var (
				tokenAddress  common.Address
				tx            *types.Transaction
				tokenContract *burn_mint_erc677.BurnMintERC677
				err2          error
			)
			if chain.IsZkSyncVM {
				tokenAddress, _, tokenContract, err2 = burn_mint_erc677.DeployBurnMintERC677Zk(
					nil,
					chain.ClientZkSyncVM,
					chain.DeployerKeyZkSyncVM,
					chain.Client,
					shared.USDCName,
					string(shared.USDCSymbol),
					shared.UsdcDecimals,
					big.NewInt(0),
				)
			} else {
				tokenAddress, tx, tokenContract, err2 = burn_mint_erc677.DeployBurnMintERC677(
					chain.DeployerKey,
					chain.Client,
					shared.USDCName,
					string(shared.USDCSymbol),
					shared.UsdcDecimals,
					big.NewInt(0),
				)
			}
			return cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
				Address:  tokenAddress,
				Contract: tokenContract,
				Tx:       tx,
				Tv:       cldf.NewTypeAndVersion(shared.USDCToken, deployment.Version1_0_0),
				Err:      err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy USDC token", "chain", chain.String(), "err", err)
		return nil, nil, nil, nil, err
	}

	tx, err := token.Contract.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	if err != nil {
		lggr.Errorw("Failed to grant mint role", "chain", chain.String(), "token", token.Contract.Address(), "err", err)
		return nil, nil, nil, nil, err
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	transmitter, err := cldf.DeployContract(lggr, chain, addresses,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter] {
			var (
				transmitterAddress  common.Address
				tx                  *types.Transaction
				transmitterContract *mock_usdc_token_transmitter.MockE2EUSDCTransmitter
				err2                error
			)
			if chain.IsZkSyncVM {
				transmitterAddress, _, transmitterContract, err2 = mock_usdc_token_transmitter.DeployMockE2EUSDCTransmitterZk(
					nil,
					chain.ClientZkSyncVM,
					chain.DeployerKeyZkSyncVM,
					chain.Client,
					0,
					reader.AllAvailableDomains()[chain.Selector],
					token.Address,
				)
			} else {
				transmitterAddress, tx, transmitterContract, err2 = mock_usdc_token_transmitter.DeployMockE2EUSDCTransmitter(
					chain.DeployerKey,
					chain.Client,
					0,
					reader.AllAvailableDomains()[chain.Selector],
					token.Address,
				)
			}
			return cldf.ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter]{
				Address:  transmitterAddress,
				Contract: transmitterContract,
				Tx:       tx,
				Tv:       cldf.NewTypeAndVersion(shared.USDCMockTransmitter, deployment.Version1_0_0),
				Err:      err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy mock USDC transmitter", "chain", chain.String(), "err", err)
		return nil, nil, nil, nil, err
	}

	messenger, err := cldf.DeployContract(lggr, chain, addresses,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger] {
			var (
				messengerAddress  common.Address
				tx                *types.Transaction
				messengerContract *mock_usdc_token_messenger.MockE2EUSDCTokenMessenger
				err2              error
			)
			if chain.IsZkSyncVM {
				messengerAddress, _, messengerContract, err2 = mock_usdc_token_messenger.DeployMockE2EUSDCTokenMessengerZk(
					nil,
					chain.ClientZkSyncVM,
					chain.DeployerKeyZkSyncVM,
					chain.Client,
					0,
					transmitter.Address,
				)
			} else {
				messengerAddress, tx, messengerContract, err2 = mock_usdc_token_messenger.DeployMockE2EUSDCTokenMessenger(
					chain.DeployerKey,
					chain.Client,
					0,
					transmitter.Address,
				)
			}
			return cldf.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger]{
				Address:  messengerAddress,
				Contract: messengerContract,
				Tx:       tx,
				Tv:       cldf.NewTypeAndVersion(shared.USDCTokenMessenger, deployment.Version1_0_0),
				Err:      err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy USDC token messenger", "chain", chain.String(), "err", err)
		return nil, nil, nil, nil, err
	}

	tokenPool, err := cldf.DeployContract(lggr, chain, addresses,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*usdc_token_pool.USDCTokenPool] {
			var (
				tokenPoolAddress  common.Address
				tx                *types.Transaction
				tokenPoolContract *usdc_token_pool.USDCTokenPool
				err2              error
			)
			if chain.IsZkSyncVM {
				tokenPoolAddress, _, tokenPoolContract, err2 = usdc_token_pool.DeployUSDCTokenPoolZk(
					nil,
					chain.ClientZkSyncVM,
					chain.DeployerKeyZkSyncVM,
					chain.Client,
					messenger.Address,
					token.Address,
					[]common.Address{},
					rmnProxy,
					router,
				)
			} else {
				tokenPoolAddress, tx, tokenPoolContract, err2 = usdc_token_pool.DeployUSDCTokenPool(
					chain.DeployerKey,
					chain.Client,
					messenger.Address,
					token.Address,
					[]common.Address{},
					rmnProxy,
					router,
				)
			}
			return cldf.ContractDeploy[*usdc_token_pool.USDCTokenPool]{
				Address:  tokenPoolAddress,
				Contract: tokenPoolContract,
				Tx:       tx,
				Tv:       cldf.NewTypeAndVersion(shared.USDCTokenPool, deployment.Version1_5_1),
				Err:      err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy USDC token pool", "chain", chain.String(), "err", err)
		return nil, nil, nil, nil, err
	}

	return token.Contract, tokenPool.Contract, messenger.Contract, transmitter.Contract, nil
}
