package ccipdeployment

import (
	"encoding/hex"
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
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/arm_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
)

var (
	MockARM              deployment.ContractType = "MockARM"
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
	EVM2EVMMultiOnRamp   deployment.ContractType = "EVM2EVMMultiOnRamp"
	EVM2EVMMultiOffRamp  deployment.ContractType = "EVM2EVMMultiOffRamp"
	CCIPReceiver         deployment.ContractType = "CCIPReceiver"
	CapabilitiesRegistry deployment.ContractType = "CapabilitiesRegistry"
)

type Contracts interface {
	*capabilities_registry.CapabilitiesRegistry |
		*arm_proxy_contract.ARMProxyContract |
		*ccip_config.CCIPConfig |
		*nonce_manager.NonceManager |
		*price_registry.PriceRegistry |
		*router.Router |
		*token_admin_registry.TokenAdminRegistry |
		*weth9.WETH9 |
		*mock_arm_contract.MockARMContract |
		*owner_helpers.ManyChainMultiSig |
		*owner_helpers.RBACTimelock |
		*evm_2_evm_multi_offramp.EVM2EVMMultiOffRamp |
		*evm_2_evm_multi_onramp.EVM2EVMMultiOnRamp |
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
	ab := deployment.NewMemoryAddressBook()
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil || len(nodes) == 0 {
		e.Logger.Errorw("Failed to get node info", "err", err)
		return ab, err
	}
	if _, ok := c.CapabilityRegistry[c.HomeChainSel]; !ok {
		return ab, fmt.Errorf("Capability registry not found for home chain %d, needs to be deployed first", c.HomeChainSel)
	}
	cr, err := c.CapabilityRegistry[c.HomeChainSel].GetHashedCapabilityId(
		&bind.CallOpts{}, CapabilityLabelledName, CapabilityVersion)
	if err != nil {
		e.Logger.Errorw("Failed to get hashed capability id", "err", err)
		return ab, err
	}
	// Signal to CR that our nodes support CCIP capability.
	if err := AddNodes(
		c.CapabilityRegistry[c.HomeChainSel],
		e.Chains[c.HomeChainSel],
		nodes.PeerIDs(c.HomeChainSel), // Doesn't actually matter which sel here
		[][32]byte{cr},
	); err != nil {
		return ab, err
	}

	for sel, chain := range e.Chains {
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
			func(chain deployment.Chain) ContractDeploy[*mock_arm_contract.MockARMContract] {
				mockARMAddr, tx, mockARM, err2 := mock_arm_contract.DeployMockARMContract(
					chain.DeployerKey,
					chain.Client,
				)
				return ContractDeploy[*mock_arm_contract.MockARMContract]{
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

		armProxy, err := deployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) ContractDeploy[*arm_proxy_contract.ARMProxyContract] {
				armProxyAddr, tx, armProxy, err2 := arm_proxy_contract.DeployARMProxyContract(
					chain.DeployerKey,
					chain.Client,
					mockARM.Address,
				)
				return ContractDeploy[*arm_proxy_contract.ARMProxyContract]{
					armProxyAddr, armProxy, tx, deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_0_0), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy armProxy", "err", err)
			return ab, err
		}
		e.Logger.Infow("deployed armProxy", "addr", armProxy.Address)

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
					armProxy.Address,
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
			func(chain deployment.Chain) ContractDeploy[*evm_2_evm_multi_onramp.EVM2EVMMultiOnRamp] {
				onRampAddr, tx, onRamp, err2 := evm_2_evm_multi_onramp.DeployEVM2EVMMultiOnRamp(
					chain.DeployerKey,
					chain.Client,
					evm_2_evm_multi_onramp.EVM2EVMMultiOnRampStaticConfig{
						ChainSelector:      sel,
						RmnProxy:           armProxy.Address,
						NonceManager:       nonceManager.Address,
						TokenAdminRegistry: tokenAdminRegistry.Address,
					},
					evm_2_evm_multi_onramp.EVM2EVMMultiOnRampDynamicConfig{
						PriceRegistry: priceRegistry.Address,
						FeeAggregator: common.HexToAddress("0x1"), // TODO real fee aggregator
					},
					[]evm_2_evm_multi_onramp.EVM2EVMMultiOnRampDestChainConfigArgs{},
				)
				return ContractDeploy[*evm_2_evm_multi_onramp.EVM2EVMMultiOnRamp]{
					onRampAddr, onRamp, tx, deployment.NewTypeAndVersion(EVM2EVMMultiOnRamp, deployment.Version1_6_0_dev), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy onramp", "err", err)
			return ab, err
		}
		e.Logger.Infow("deployed onramp", "addr", onRamp.Address)

		offRamp, err := deployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) ContractDeploy[*evm_2_evm_multi_offramp.EVM2EVMMultiOffRamp] {
				offRampAddr, tx, offRamp, err2 := evm_2_evm_multi_offramp.DeployEVM2EVMMultiOffRamp(
					chain.DeployerKey,
					chain.Client,
					evm_2_evm_multi_offramp.EVM2EVMMultiOffRampStaticConfig{
						ChainSelector:      sel,
						RmnProxy:           armProxy.Address,
						NonceManager:       nonceManager.Address,
						TokenAdminRegistry: tokenAdminRegistry.Address,
					},
					evm_2_evm_multi_offramp.EVM2EVMMultiOffRampDynamicConfig{
						PriceRegistry:                           priceRegistry.Address,
						PermissionLessExecutionThresholdSeconds: uint32(86400),
						MaxTokenTransferGas:                     uint32(200_000),
						MaxPoolReleaseOrMintGas:                 uint32(200_000),
					},
					[]evm_2_evm_multi_offramp.EVM2EVMMultiOffRampSourceChainConfigArgs{},
				)
				return ContractDeploy[*evm_2_evm_multi_offramp.EVM2EVMMultiOffRamp]{
					offRampAddr, offRamp, tx, deployment.NewTypeAndVersion(EVM2EVMMultiOffRamp, deployment.Version1_6_0_dev), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy offramp", "err", err)
			return ab, err
		}
		e.Logger.Infow("deployed offramp", "addr", offRamp)

		// Enable ramps on price registry/nonce manager
		tx, err := priceRegistry.Contract.ApplyAuthorizedCallerUpdates(chain.DeployerKey, price_registry.AuthorizedCallersAuthorizedCallerArgs{
			// TODO: We enable the deployer initially to set prices
			AddedCallers: []common.Address{offRamp.Address, chain.DeployerKey.From},
		})
		if err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
			e.Logger.Errorw("Failed to confirm price registry authorized caller update", "err", err)
			return ab, err
		}

		tx, err = nonceManager.Contract.ApplyAuthorizedCallerUpdates(chain.DeployerKey, nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
			AddedCallers: []common.Address{offRamp.Address, onRamp.Address},
		})
		if err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
			e.Logger.Errorw("Failed to update nonce manager with ramps", "err", err)
			return ab, err
		}

		// Add chain config for each chain.
		_, err = AddChainConfig(e.Logger,
			e.Chains[c.HomeChainSel],
			c.CCIPOnChainState.CCIPConfig[c.HomeChainSel],
			chain.Selector,
			nodes.PeerIDs(chain.Selector),
			uint8(len(nodes)/3))
		if err != nil {
			return ab, err
		}

		// For each chain, we create a DON on the home chain.
		if err := AddDON(e.Logger,
			cr,
			c.CapabilityRegistry[c.HomeChainSel],
			c.CCIPConfig[c.HomeChainSel],
			offRamp.Contract,
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

func AddLane(e deployment.Environment, state CCIPOnChainState, from, to uint64) error {
	// TODO: Batch
	tx, err := state.Routers[from].ApplyRampUpdates(e.Chains[from].DeployerKey, []router.RouterOnRamp{
		{
			DestChainSelector: to,
			OnRamp:            state.EvmOnRampsV160[from].Address(),
		},
	}, []router.RouterOffRamp{}, []router.RouterOffRamp{})
	if err := deployment.ConfirmIfNoError(e.Chains[from], tx, err); err != nil {
		return err
	}
	tx, err = state.EvmOnRampsV160[from].ApplyDestChainConfigUpdates(e.Chains[from].DeployerKey,
		[]evm_2_evm_multi_onramp.EVM2EVMMultiOnRampDestChainConfigArgs{
			{
				DestChainSelector: to,
				Router:            state.Routers[from].Address(),
			},
		})
	if err := deployment.ConfirmIfNoError(e.Chains[from], tx, err); err != nil {
		return err
	}

	_, err = state.PriceRegistries[from].UpdatePrices(
		e.Chains[from].DeployerKey, price_registry.InternalPriceUpdates{
			TokenPriceUpdates: []price_registry.InternalTokenPriceUpdate{
				{
					SourceToken: state.LinkTokens[from].Address(),
					UsdPerToken: deployment.E18Mult(20),
				},
				{
					SourceToken: state.Weth9s[from].Address(),
					UsdPerToken: deployment.E18Mult(4000),
				},
			},
			GasPriceUpdates: []price_registry.InternalGasPriceUpdate{
				{
					DestChainSelector: to,
					UsdPerUnitGas:     big.NewInt(2e12),
				},
			}})
	if err := deployment.ConfirmIfNoError(e.Chains[from], tx, err); err != nil {
		return err
	}

	// Enable dest in price registry
	tx, err = state.PriceRegistries[from].ApplyDestChainConfigUpdates(e.Chains[from].DeployerKey,
		[]price_registry.PriceRegistryDestChainConfigArgs{
			{
				DestChainSelector: to,
				DestChainConfig:   defaultPriceRegistryDestChainConfig(),
			},
		})
	if err := deployment.ConfirmIfNoError(e.Chains[from], tx, err); err != nil {
		return err
	}

	tx, err = state.EvmOffRampsV160[to].ApplySourceChainConfigUpdates(e.Chains[to].DeployerKey,
		[]evm_2_evm_multi_offramp.EVM2EVMMultiOffRampSourceChainConfigArgs{
			{
				Router:              state.Routers[to].Address(),
				SourceChainSelector: from,
				IsEnabled:           true,
				OnRamp:              common.LeftPadBytes(state.EvmOnRampsV160[from].Address().Bytes(), 32),
			},
		})
	if err := deployment.ConfirmIfNoError(e.Chains[to], tx, err); err != nil {
		return err
	}
	tx, err = state.Routers[to].ApplyRampUpdates(e.Chains[to].DeployerKey, []router.RouterOnRamp{}, []router.RouterOffRamp{}, []router.RouterOffRamp{
		{
			SourceChainSelector: from,
			OffRamp:             state.EvmOffRampsV160[to].Address(),
		},
	})
	return deployment.ConfirmIfNoError(e.Chains[to], tx, err)
}

func defaultPriceRegistryDestChainConfig() price_registry.PriceRegistryDestChainConfig {
	// https://github.com/smartcontractkit/ccip/blob/c4856b64bd766f1ddbaf5d13b42d3c4b12efde3a/contracts/src/v0.8/ccip/libraries/Internal.sol#L337-L337
	/*
		```Solidity
			// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
			bytes4 public constant CHAIN_FAMILY_SELECTOR_EVM = 0x2812d52c;
		```
	*/
	evmFamilySelector, _ := hex.DecodeString("2812d52c")
	return price_registry.PriceRegistryDestChainConfig{
		IsEnabled:                         true,
		MaxNumberOfTokensPerMsg:           10,
		MaxDataBytes:                      256,
		MaxPerMsgGasLimit:                 3_000_000,
		DestGasOverhead:                   50_000,
		DefaultTokenFeeUSDCents:           1,
		DestGasPerPayloadByte:             10,
		DestDataAvailabilityOverheadGas:   0,
		DestGasPerDataAvailabilityByte:    100,
		DestDataAvailabilityMultiplierBps: 1,
		DefaultTokenDestGasOverhead:       125_000,
		DefaultTokenDestBytesOverhead:     32,
		DefaultTxGasLimit:                 200_000,
		GasMultiplierWeiPerEth:            1,
		NetworkFeeUSDCents:                1,
		ChainFamilySelector:               [4]byte(evmFamilySelector),
	}
}
