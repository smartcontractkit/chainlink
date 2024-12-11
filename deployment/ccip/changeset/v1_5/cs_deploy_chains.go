package v1_5

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
)

var _ deployment.ChangeSet[DeployChainContractsConfig] = DeployChainContracts

type DeployChainContractsConfig struct {
	configs []DeployChainContractsConfigPerChain
}

type DeployChainContractsConfigPerChain struct {
	ChainSelector              uint64
	RMNConfig                  *rmn_contract.RMNConfig
	PriceRegStalenessThreshold uint32
}

func DeployChainContracts(env deployment.Environment, c DeployChainContractsConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid DeployChainContractsConfig: %w", err)
	}
	newAddresses := deployment.NewMemoryAddressBook()
	err := deployChainContractsForChains(env, newAddresses, c)
	if err != nil {
		env.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", newAddresses)
		return deployment.ChangesetOutput{AddressBook: newAddresses}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: newAddresses,
		JobSpecs:    nil,
	}, nil
}

func (c DeployChainContractsConfig) Validate() error {
	for _, cs := range c.configs {
		if err := deployment.IsValidChainSelector(cs.ChainSelector); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
	}
	// TODO: Add rest of the config validation
	return nil
}

func deployChainContractsForChains(env deployment.Environment, ab deployment.AddressBook, cfg DeployChainContractsConfig) error {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for _, c := range cfg.configs {
		if err := deployChainContracts(env, ab, state, c); err != nil {
			return fmt.Errorf("failed to deploy CCIP contracts for chain %d: %w", c.ChainSelector, err)
		}
	}
	return nil
}

func deployChainContracts(
	env deployment.Environment,
	ab deployment.AddressBook,
	state changeset.CCIPOnChainState,
	cfg DeployChainContractsConfigPerChain) error {
	chain, ok := env.Chains[cfg.ChainSelector]
	if !ok {
		return fmt.Errorf("chain %d not found", cfg.ChainSelector)
	}
	chainState, ok := state.Chains[cfg.ChainSelector]
	if !ok {
		return fmt.Errorf("chain %d state not found", cfg.ChainSelector)
	}
	lggr := env.Logger
	// ================================================================
	// │                         Deploy RMN                           │
	// ================================================================
	var rmnAddr common.Address
	if cfg.RMNConfig == nil {
		if chainState.MockRMN == nil {
			rmn, err := deployment.DeployContract(lggr, chain, ab,
				func(chain deployment.Chain) deployment.ContractDeploy[*mock_rmn_contract.MockRMNContract] {
					rmnAddr, tx2, rmn, err2 := mock_rmn_contract.DeployMockRMNContract(
						chain.DeployerKey,
						chain.Client,
					)
					return deployment.ContractDeploy[*mock_rmn_contract.MockRMNContract]{
						Address: rmnAddr, Contract: rmn, Tx: tx2, Tv: deployment.NewTypeAndVersion(changeset.MockRMN, deployment.Version1_0_0), Err: err2,
					}
				})
			if err != nil {
				lggr.Errorw("Failed to deploy mock RMN", "chain", chain.String(), "err", err)
				return err
			}
			rmnAddr = rmn.Address
		} else {
			lggr.Infow("MockRMN already deployed", "chain", chain.String(), "address", chainState.MockRMN.Address())
			rmnAddr = chainState.MockRMN.Address()
		}
	} else {
		if chainState.RMN == nil {
			rmn, err := deployment.DeployContract(lggr, chain, ab,
				func(chain deployment.Chain) deployment.ContractDeploy[*rmn_contract.RMNContract] {
					rmnAddr, tx2, rmn, err2 := rmn_contract.DeployRMNContract(
						chain.DeployerKey,
						chain.Client,
						*cfg.RMNConfig,
					)
					return deployment.ContractDeploy[*rmn_contract.RMNContract]{
						Address: rmnAddr, Contract: rmn, Tx: tx2, Tv: deployment.NewTypeAndVersion(changeset.RMN, deployment.Version1_5_0), Err: err2,
					}
				})
			if err != nil {
				lggr.Errorw("Failed to deploy RMN", "chain", chain.String(), "err", err)
				return err
			}
			rmnAddr = rmn.Address
		} else {
			lggr.Infow("RMN already deployed", "chain", chain.String(), "address", chainState.RMN.Address)
			rmnAddr = chainState.RMN.Address()
		}
	}
	var rmnProxyAddress common.Address
	if chainState.RMNProxy == nil {
		rmnProxy, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*rmn_proxy_contract.RMNProxyContract] {
				rmnProxyAddr, tx2, rmnProxy, err2 := rmn_proxy_contract.DeployRMNProxyContract(
					chain.DeployerKey,
					chain.Client,
					rmnAddr,
				)
				return deployment.ContractDeploy[*rmn_proxy_contract.RMNProxyContract]{
					Address: rmnProxyAddr, Contract: rmnProxy, Tx: tx2, Tv: deployment.NewTypeAndVersion(changeset.ARMProxy, deployment.Version1_0_0), Err: err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy RMNProxy", "chain", chain.String(), "err", err)
			return err
		}
		rmnProxyAddress = rmnProxy.Address
	} else {
		lggr.Infow("RMNProxy already deployed", "chain", chain.String(), "address", chainState.RMNProxy.Address)
		// check if the RMNProxy is pointing to the correct RMN
		rmnProxy := chainState.RMNProxy
		setRMN, err := rmnProxy.GetARM(nil)
		if err != nil {
			return err
		}
		if setRMN != rmnAddr {
			lggr.Infow("RMNProxy pointing to wrong RMN, changing the RMN", "chain", chain.String(), "rmnProxy", rmnProxy.Address, "rmn", rmnAddr)
			tx, err := rmnProxy.SetARM(chain.DeployerKey, rmnAddr)
			if err != nil {
				return fmt.Errorf("failed to set RMN on RMNProxy for chain %s RMN %s RMNProxy %s: %w",
					chain.String(), rmnAddr.String(), rmnProxy.Address().String(), err)
			}

			_, err = chain.Confirm(tx)
			if err != nil {
				lggr.Errorw("Failed to confirm RMNProxy SetARM",
					"chain", chain.String(), "rmn", rmnAddr.String(), "rmnProxy", rmnProxy.Address().String(), "err", err)
				return fmt.Errorf("failed to confirm RMNProxy SetARM for chain %s RMN %s RMNProxy %s: %w",
					chain.String(), rmnAddr.String(), rmnProxy.Address().String(), err)
			}
			lggr.Infow("RMNProxy SetARM confirmed",
				"chain", chain.String(), "rmn", rmnAddr.String(), "rmnProxy", rmnProxy.Address().String())
		}
		rmnProxyAddress = rmnProxy.Address()
	}
	// ================================================================
	// │                 Deploy TokenAdminRegistry                    │
	// ================================================================
	var tokenAdminReg *token_admin_registry.TokenAdminRegistry
	var registryModule *registry_module_owner_custom.RegistryModuleOwnerCustom
	if chainState.TokenAdminRegistry == nil {
		tokenAdminRegistry, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*token_admin_registry.TokenAdminRegistry] {
				tokenAdminRegistryAddr, tx2, tokenAdminRegistry, err2 := token_admin_registry.DeployTokenAdminRegistry(
					chain.DeployerKey,
					chain.Client)
				return deployment.ContractDeploy[*token_admin_registry.TokenAdminRegistry]{
					Address: tokenAdminRegistryAddr, Contract: tokenAdminRegistry, Tx: tx2,
					Tv: deployment.NewTypeAndVersion(changeset.TokenAdminRegistry, deployment.Version1_5_0), Err: err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy token admin registry", "chain", chain.String(), "err", err)
			return err
		}
		tokenAdminReg = tokenAdminRegistry.Contract
	} else {
		lggr.Infow("TokenAdminRegistry already deployed", "chain", chain.String(), "address", chainState.TokenAdminRegistry.Address)
		tokenAdminReg = chainState.TokenAdminRegistry
	}
	if chainState.RegistryModule == nil {
		customRegistryModule, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom] {
				regModAddr, tx2, regMod, err2 := registry_module_owner_custom.DeployRegistryModuleOwnerCustom(
					chain.DeployerKey,
					chain.Client,
					tokenAdminReg.Address())
				return deployment.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom]{
					Address: regModAddr, Contract: regMod, Tx: tx2,
					Tv: deployment.NewTypeAndVersion(changeset.RegistryModule, deployment.Version1_5_0), Err: err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy custom registry module", "chain", chain.String(), "err", err)
			return err
		}
		registryModule = customRegistryModule.Contract
	} else {
		lggr.Infow("custom registry module already deployed", "chain", chain.String(), "addr", registryModule.Address)
	}
	isRegistryAdded, err := tokenAdminReg.IsRegistryModule(nil, registryModule.Address())
	if err != nil {
		lggr.Errorw("Failed to check if registry module is added on token admin registry", "chain", chain.String(), "err", err)
		return fmt.Errorf("failed to check if registry module is added on token admin registry: %w", err)
	}
	if !isRegistryAdded {
		tx, err := tokenAdminReg.AddRegistryModule(chain.DeployerKey, registryModule.Address())
		if err != nil {
			lggr.Errorw("Failed to assign registry module on token admin registry", "chain", chain.String(), "err", err)
			return fmt.Errorf("failed to assign registry module on token admin registry: %w", err)
		}

		_, err = chain.Confirm(tx)
		if err != nil {
			lggr.Errorw("Failed to confirm assign registry module on token admin registry", "chain", chain.String(), "err", err)
			return fmt.Errorf("failed to confirm assign registry module on token admin registry: %w", err)
		}
		lggr.Infow("assigned registry module on token admin registry")
	}
	// ================================================================
	// │                       Deploy Tokens                          │
	// ================================================================
	var weth9Address common.Address
	var linkAddress common.Address
	if chainState.Weth9 == nil {
		weth, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*weth9.WETH9] {
				weth9Addr, tx2, weth9c, err2 := weth9.DeployWETH9(
					chain.DeployerKey,
					chain.Client,
				)
				return deployment.ContractDeploy[*weth9.WETH9]{
					Address: weth9Addr, Contract: weth9c, Tx: tx2,
					Tv: deployment.NewTypeAndVersion(changeset.WETH9, deployment.Version1_0_0), Err: err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy weth9", "chain", chain.String(), "err", err)
			return err
		}
		weth9Address = weth.Address
	} else {
		lggr.Infow("weth9 already deployed", "chain", chain.String(), "addr", chainState.Weth9.Address())
	}
	// check if link token is already deployed
	if chainState.LinkToken == nil {
		lggr.Errorw("Link token not deployed", "chain", chain.String())
		return fmt.Errorf("link token not deployed for chain %s, deploy link first with DeployLinkToken changeset", chain.String())
	} else {
		linkAddress = chainState.LinkToken.Address()
	}

	// ================================================================
	// │                       Deploy Routers                         │
	// ================================================================
	if chainState.Router == nil {
		_, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*router.Router] {
				routerAddr, tx2, routerC, err2 := router.DeployRouter(
					chain.DeployerKey,
					chain.Client,
					weth9Address,
					rmnProxyAddress,
				)
				return deployment.ContractDeploy[*router.Router]{
					Address: routerAddr, Contract: routerC, Tx: tx2,
					Tv: deployment.NewTypeAndVersion(changeset.Router, deployment.Version1_2_0), Err: err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy router", "chain", chain.String(), "err", err)
			return err
		}
	} else {
		lggr.Infow("router already deployed", "chain", chain.String(), "addr", chainState.Router.Address)
	}
	// ================================================================
	// │                    Deploy Price Registry                     │
	// ================================================================
	if chainState.PriceRegistry == nil {
		_, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*price_registry_1_2_0.PriceRegistry] {
				priceRegAddr, tx2, priceRegAddrC, err2 := price_registry_1_2_0.DeployPriceRegistry(
					chain.DeployerKey,
					chain.Client,
					nil,
					[]common.Address{weth9Address, linkAddress},
					cfg.PriceRegStalenessThreshold,
				)
				return deployment.ContractDeploy[*price_registry_1_2_0.PriceRegistry]{
					Address: priceRegAddr, Contract: priceRegAddrC, Tx: tx2,
					Tv: deployment.NewTypeAndVersion(changeset.PriceRegistry, deployment.Version1_2_0), Err: err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy PriceRegistry", "chain", chain.String(), "err", err)
			return err
		}
	} else {
		lggr.Infow("PriceRegistry already deployed", "chain", chain.String(), "addr", chainState.PriceRegistry.Address)
	}
	// ================================================================
	// │                    Deploy Receiver                     	  │
	// ================================================================
	if chainState.Receiver == nil {
		_, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*maybe_revert_message_receiver.MaybeRevertMessageReceiver] {
				receiverAddr, tx, receiver, err2 := maybe_revert_message_receiver.DeployMaybeRevertMessageReceiver(
					chain.DeployerKey,
					chain.Client,
					false,
				)
				return deployment.ContractDeploy[*maybe_revert_message_receiver.MaybeRevertMessageReceiver]{
					Address: receiverAddr, Contract: receiver, Tx: tx,
					Tv: deployment.NewTypeAndVersion(changeset.CCIPReceiver, deployment.Version1_0_0), Err: err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy receiver", "chain", chain.String(), "err", err)
			return err
		}
	} else {
		lggr.Infow("Receiver already deployed", "chain", chain.String(), "addr", chainState.Receiver.Address)
	}
	return nil
}
