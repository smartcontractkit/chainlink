package ccipdeployment

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/gethwrappers"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
)

var (
	MockARM              deployment.ContractType = "MockRMN"
	LinkToken            deployment.ContractType = "LinkToken"
	ARMProxy             deployment.ContractType = "ARMProxy"
	WETH9                deployment.ContractType = "WETH9"
	Router               deployment.ContractType = "Router"
	TokenAdminRegistry   deployment.ContractType = "TokenAdminRegistry"
	NonceManager         deployment.ContractType = "NonceManager"
	PriceRegistry        deployment.ContractType = "PriceRegistry"
	ManyChainMultisig    deployment.ContractType = "ManyChainMultiSig"
	CCIPConfig           deployment.ContractType = "CCIPConfig"
	RBACTimelock         deployment.ContractType = "RBACTimelock"
	OnRamp               deployment.ContractType = "OnRamp"
	OffRamp              deployment.ContractType = "OffRamp"
	CCIPReceiver         deployment.ContractType = "CCIPReceiver"
	CapabilitiesRegistry deployment.ContractType = "CapabilitiesRegistry"
)

type Contracts interface {
	*capabilities_registry.CapabilitiesRegistry |
		*rmn_proxy_contract.RMNProxyContract |
		*ccip_config.CCIPConfig |
		*nonce_manager.NonceManager |
		*price_registry.PriceRegistry |
		*router.Router |
		*token_admin_registry.TokenAdminRegistry |
		*weth9.WETH9 |
		*mock_rmn_contract.MockRMNContract |
		*owner_helpers.ManyChainMultiSig |
		*owner_helpers.RBACTimelock |
		*offramp.OffRamp |
		*onramp.OnRamp |
		*burn_mint_erc677.BurnMintERC677 |
		*maybe_revert_message_receiver.MaybeRevertMessageReceiver
}

type ContractDeploy[C Contracts] struct {
	// We just keep all the deploy return values
	// since some will be empty if there's an error.
	Address  common.Address
	Contract C
	Tx       *types.Transaction
	Tv       deployment.TypeAndVersion
	Err      error
}

// TODO: pull up to general deployment pkg somehow
// without exposing all product specific contracts?
func deployContract[C Contracts](
	lggr logger.Logger,
	chain deployment.Chain,
	addressBook deployment.AddressBook,
	deploy func(chain deployment.Chain) ContractDeploy[C],
) (*ContractDeploy[C], error) {
	contractDeploy := deploy(chain)
	if contractDeploy.Err != nil {
		lggr.Errorw("Failed to deploy contract", "err", contractDeploy.Err)
		return nil, contractDeploy.Err
	}
	err := chain.Confirm(contractDeploy.Tx.Hash())
	if err != nil {
		lggr.Errorw("Failed to confirm deployment", "err", err)
		return nil, err
	}
	err = addressBook.Save(chain.Selector, contractDeploy.Address.String(), contractDeploy.Tv)
	if err != nil {
		lggr.Errorw("Failed to save contract address", "err", err)
		return nil, err
	}
	return &contractDeploy, nil
}

type DeployCCIPContractConfig struct {
	HomeChainSel uint64
	// Existing contracts which we want to skip deployment
	// Leave empty if we want to deploy everything
	// TODO: Add skips to deploy function.
	CCIPOnChainState
}

// TODO: Likely we'll want to further parameterize the deployment
// For example a list of contracts to skip deploying if they already exist.
// Or mock vs real RMN.
// Deployment produces an address book of everything it deployed.
func DeployCCIPContracts(e deployment.Environment, c DeployCCIPContractConfig) (deployment.AddressBook, error) {
	var ab deployment.AddressBook = deployment.NewMemoryAddressBook()
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil || len(nodes) == 0 {
		e.Logger.Errorw("Failed to get node info", "err", err)
		return ab, err
	}
	if c.Chains[c.HomeChainSel].CapabilityRegistry == nil {
		return ab, fmt.Errorf("Capability registry not found for home chain %d, needs to be deployed first", c.HomeChainSel)
	}
	cr, err := c.Chains[c.HomeChainSel].CapabilityRegistry.GetHashedCapabilityId(
		&bind.CallOpts{}, CapabilityLabelledName, CapabilityVersion)
	if err != nil {
		e.Logger.Errorw("Failed to get hashed capability id", "err", err)
		return ab, err
	}
	// Signal to CR that our nodes support CCIP capability.
	if err := AddNodes(
		c.Chains[c.HomeChainSel].CapabilityRegistry,
		e.Chains[c.HomeChainSel],
		nodes.PeerIDs(c.HomeChainSel), // Doesn't actually matter which sel here
		[][32]byte{cr},
	); err != nil {
		return ab, err
	}

	for _, chain := range e.Chains {
		ab, err = DeployChainContracts(e, chain, ab)
		if err != nil {
			return ab, err
		}
		chainAddresses, err := ab.AddressesForChain(chain.Selector)
		if err != nil {
			e.Logger.Errorw("Failed to get chain addresses", "err", err)
			return ab, err
		}
		chainState, err := LoadChainState(chain, chainAddresses)
		if err != nil {
			e.Logger.Errorw("Failed to load chain state", "err", err)
			return ab, err
		}
		// Enable ramps on price registry/nonce manager
		tx, err := chainState.PriceRegistry.ApplyAuthorizedCallerUpdates(chain.DeployerKey, price_registry.AuthorizedCallersAuthorizedCallerArgs{
			// TODO: We enable the deployer initially to set prices
			AddedCallers: []common.Address{chainState.EvmOffRampV160.Address(), chain.DeployerKey.From},
		})
		if err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
			e.Logger.Errorw("Failed to confirm price registry authorized caller update", "err", err)
			return ab, err
		}

		tx, err = chainState.NonceManager.ApplyAuthorizedCallerUpdates(chain.DeployerKey, nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
			AddedCallers: []common.Address{chainState.EvmOffRampV160.Address(), chainState.EvmOnRampV160.Address()},
		})
		if err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
			e.Logger.Errorw("Failed to update nonce manager with ramps", "err", err)
			return ab, err
		}

		// Add chain config for each chain.
		_, err = AddChainConfig(e.Logger,
			e.Chains[c.HomeChainSel],
			c.Chains[c.HomeChainSel].CCIPConfig,
			chain.Selector,
			nodes.PeerIDs(chain.Selector),
			uint8(len(nodes)/3))
		if err != nil {
			return ab, err
		}

		// For each chain, we create a DON on the home chain.
		if err := AddDON(e.Logger,
			cr,
			c.Chains[c.HomeChainSel].CapabilityRegistry,
			c.Chains[c.HomeChainSel].CCIPConfig,
			chainState.EvmOffRampV160,
			chain,
			e.Chains[c.HomeChainSel],
			uint8(len(nodes)/3),
			nodes.BootstrapPeerIDs(chain.Selector)[0],
			nodes.PeerIDs(chain.Selector),
			nodes,
		); err != nil {
			e.Logger.Errorw("Failed to add DON", "err", err)
			return ab, err
		}
	}

	return ab, nil
}

func DeployChainContracts(e deployment.Environment, chain deployment.Chain, ab deployment.AddressBook) (deployment.AddressBook, error) {
	ccipReceiver, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*maybe_revert_message_receiver.MaybeRevertMessageReceiver] {
			receiverAddr, tx, receiver, err2 := maybe_revert_message_receiver.DeployMaybeRevertMessageReceiver(
				chain.DeployerKey,
				chain.Client,
				false,
			)
			return ContractDeploy[*maybe_revert_message_receiver.MaybeRevertMessageReceiver]{
				receiverAddr, receiver, tx, deployment.NewTypeAndVersion(CCIPReceiver, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy receiver", "err", err)
		return ab, err
	}
	e.Logger.Infow("deployed receiver", "addr", ccipReceiver.Address)

	// TODO: Still waiting for RMNRemote/RMNHome contracts etc.
	mockARM, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*mock_rmn_contract.MockRMNContract] {
			mockARMAddr, tx, mockARM, err2 := mock_rmn_contract.DeployMockRMNContract(
				chain.DeployerKey,
				chain.Client,
			)
			return ContractDeploy[*mock_rmn_contract.MockRMNContract]{
				mockARMAddr, mockARM, tx, deployment.NewTypeAndVersion(MockARM, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy mockARM", "err", err)
		return ab, err
	}
	e.Logger.Infow("deployed mockARM", "addr", mockARM)

	mcm, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*owner_helpers.ManyChainMultiSig] {
			mcmAddr, tx, mcm, err2 := owner_helpers.DeployManyChainMultiSig(
				chain.DeployerKey,
				chain.Client,
			)
			return ContractDeploy[*owner_helpers.ManyChainMultiSig]{
				mcmAddr, mcm, tx, deployment.NewTypeAndVersion(ManyChainMultisig, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy mcm", "err", err)
		return ab, err
	}
	// TODO: Address soon
	e.Logger.Infow("deployed mcm", "addr", mcm.Address)

	_, err = deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*owner_helpers.RBACTimelock] {
			timelock, tx, cc, err2 := owner_helpers.DeployRBACTimelock(
				chain.DeployerKey,
				chain.Client,
				big.NewInt(0), // minDelay
				mcm.Address,
				[]common.Address{mcm.Address},            // proposers
				[]common.Address{chain.DeployerKey.From}, //executors
				[]common.Address{mcm.Address},            // cancellers
				[]common.Address{mcm.Address},            // bypassers
			)
			return ContractDeploy[*owner_helpers.RBACTimelock]{
				timelock, cc, tx, deployment.NewTypeAndVersion(RBACTimelock, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy timelock", "err", err)
		return ab, err
	}
	e.Logger.Infow("deployed timelock", "addr", mcm.Address)

	rmnProxy, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*rmn_proxy_contract.RMNProxyContract] {
			rmnProxyAddr, tx, rmnProxy, err2 := rmn_proxy_contract.DeployRMNProxyContract(
				chain.DeployerKey,
				chain.Client,
				mockARM.Address,
			)
			return ContractDeploy[*rmn_proxy_contract.RMNProxyContract]{
				rmnProxyAddr, rmnProxy, tx, deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy rmnProxy", "err", err)
		return ab, err
	}
	e.Logger.Infow("deployed rmnProxy", "addr", rmnProxy.Address)

	weth9, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*weth9.WETH9] {
			weth9Addr, tx, weth9c, err2 := weth9.DeployWETH9(
				chain.DeployerKey,
				chain.Client,
			)
			return ContractDeploy[*weth9.WETH9]{
				weth9Addr, weth9c, tx, deployment.NewTypeAndVersion(WETH9, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy weth9", "err", err)
		return ab, err
	}

	linkToken, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			linkTokenAddr, tx, linkToken, err2 := burn_mint_erc677.DeployBurnMintERC677(
				chain.DeployerKey,
				chain.Client,
				"Link Token",
				"LINK",
				uint8(18),
				big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
			)
			return ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
				linkTokenAddr, linkToken, tx, deployment.NewTypeAndVersion(LinkToken, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy linkToken", "err", err)
		return ab, err
	}

	routerContract, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*router.Router] {
			routerAddr, tx, routerC, err2 := router.DeployRouter(
				chain.DeployerKey,
				chain.Client,
				weth9.Address,
				rmnProxy.Address,
			)
			return ContractDeploy[*router.Router]{
				routerAddr, routerC, tx, deployment.NewTypeAndVersion(Router, deployment.Version1_2_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy router", "err", err)
		return ab, err
	}
	e.Logger.Infow("deployed router", "addr", routerContract)

	tokenAdminRegistry, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*token_admin_registry.TokenAdminRegistry] {
			tokenAdminRegistryAddr, tx, tokenAdminRegistry, err2 := token_admin_registry.DeployTokenAdminRegistry(
				chain.DeployerKey,
				chain.Client)
			return ContractDeploy[*token_admin_registry.TokenAdminRegistry]{
				tokenAdminRegistryAddr, tokenAdminRegistry, tx, deployment.NewTypeAndVersion(TokenAdminRegistry, deployment.Version1_5_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy token admin registry", "err", err)
		return ab, err
	}
	e.Logger.Infow("deployed tokenAdminRegistry", "addr", tokenAdminRegistry)

	nonceManager, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*nonce_manager.NonceManager] {
			nonceManagerAddr, tx, nonceManager, err2 := nonce_manager.DeployNonceManager(
				chain.DeployerKey,
				chain.Client,
				[]common.Address{}, // Need to add onRamp after
			)
			return ContractDeploy[*nonce_manager.NonceManager]{
				nonceManagerAddr, nonceManager, tx, deployment.NewTypeAndVersion(NonceManager, deployment.Version1_6_0_dev), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy router", "err", err)
		return ab, err
	}

	priceRegistry, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*price_registry.PriceRegistry] {
			prAddr, tx, pr, err2 := price_registry.DeployPriceRegistry(
				chain.DeployerKey,
				chain.Client,
				price_registry.PriceRegistryStaticConfig{
					MaxFeeJuelsPerMsg:  big.NewInt(0).Mul(big.NewInt(2e2), big.NewInt(1e18)),
					LinkToken:          linkToken.Address,
					StalenessThreshold: uint32(24 * 60 * 60),
				},
				[]common.Address{}, // ramps added after
				[]common.Address{weth9.Address, linkToken.Address}, // fee tokens
				[]price_registry.PriceRegistryTokenPriceFeedUpdate{},
				[]price_registry.PriceRegistryTokenTransferFeeConfigArgs{}, // TODO: tokens
				[]price_registry.PriceRegistryPremiumMultiplierWeiPerEthArgs{
					{
						PremiumMultiplierWeiPerEth: 9e17, // 0.9 ETH
						Token:                      linkToken.Address,
					},
					{
						PremiumMultiplierWeiPerEth: 1e18,
						Token:                      weth9.Address,
					},
				},
				[]price_registry.PriceRegistryDestChainConfigArgs{},
			)
			return ContractDeploy[*price_registry.PriceRegistry]{
				prAddr, pr, tx, deployment.NewTypeAndVersion(PriceRegistry, deployment.Version1_6_0_dev), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy price registry", "err", err)
		return ab, err
	}

	onRamp, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*onramp.OnRamp] {
			onRampAddr, tx, onRamp, err2 := onramp.DeployOnRamp(
				chain.DeployerKey,
				chain.Client,
				onramp.OnRampStaticConfig{
					ChainSelector:      chain.Selector,
					RmnProxy:           rmnProxy.Address,
					NonceManager:       nonceManager.Address,
					TokenAdminRegistry: tokenAdminRegistry.Address,
				},
				onramp.OnRampDynamicConfig{
					PriceRegistry: priceRegistry.Address,
					FeeAggregator: common.HexToAddress("0x1"), // TODO real fee aggregator
				},
				[]onramp.OnRampDestChainConfigArgs{},
			)
			return ContractDeploy[*onramp.OnRamp]{
				onRampAddr, onRamp, tx, deployment.NewTypeAndVersion(OnRamp, deployment.Version1_6_0_dev), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy onramp", "err", err)
		return ab, err
	}
	e.Logger.Infow("deployed onramp", "addr", onRamp.Address)

	offRamp, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*offramp.OffRamp] {
			offRampAddr, tx, offRamp, err2 := offramp.DeployOffRamp(
				chain.DeployerKey,
				chain.Client,
				offramp.OffRampStaticConfig{
					ChainSelector:      chain.Selector,
					RmnProxy:           rmnProxy.Address,
					NonceManager:       nonceManager.Address,
					TokenAdminRegistry: tokenAdminRegistry.Address,
				},
				offramp.OffRampDynamicConfig{
					PriceRegistry:                           priceRegistry.Address,
					PermissionLessExecutionThresholdSeconds: uint32(86400),
					MaxTokenTransferGas:                     uint32(200_000),
					MaxPoolReleaseOrMintGas:                 uint32(200_000),
				},
				[]offramp.OffRampSourceChainConfigArgs{},
			)
			return ContractDeploy[*offramp.OffRamp]{
				offRampAddr, offRamp, tx, deployment.NewTypeAndVersion(OffRamp, deployment.Version1_6_0_dev), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy offramp", "err", err)
		return ab, err
	}
	e.Logger.Infow("deployed offramp", "addr", offRamp)
	return ab, nil
}
