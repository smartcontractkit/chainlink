package ccipdeployment

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/configwrappers"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
)

var (
	RMNRemote            deployment.ContractType = "RMNRemote"
	LinkToken            deployment.ContractType = "LinkToken"
	ARMProxy             deployment.ContractType = "ARMProxy"
	WETH9                deployment.ContractType = "WETH9"
	Router               deployment.ContractType = "Router"
	TokenAdminRegistry   deployment.ContractType = "TokenAdminRegistry"
	NonceManager         deployment.ContractType = "NonceManager"
	FeeQuoter            deployment.ContractType = "FeeQuoter"
	ManyChainMultisig    deployment.ContractType = "ManyChainMultiSig"
	CCIPConfig           deployment.ContractType = "CCIPConfig"
	RBACTimelock         deployment.ContractType = "RBACTimelock"
	OnRamp               deployment.ContractType = "OnRamp"
	OffRamp              deployment.ContractType = "OffRamp"
	CapabilitiesRegistry deployment.ContractType = "CapabilitiesRegistry"
	PriceFeed            deployment.ContractType = "PriceFeed"
	// Note test router maps to a regular router contract.
	TestRouter   deployment.ContractType = "TestRouter"
	CCIPReceiver deployment.ContractType = "CCIPReceiver"

	TestXXXMCMSSigner *ecdsa.PrivateKey
)

func init() {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	TestXXXMCMSSigner = key
}

type Contracts interface {
	*capabilities_registry.CapabilitiesRegistry |
		*rmn_proxy_contract.RMNProxyContract |
		*ccip_config.CCIPConfig |
		*nonce_manager.NonceManager |
		*fee_quoter.FeeQuoter |
		*router.Router |
		*token_admin_registry.TokenAdminRegistry |
		*weth9.WETH9 |
		*rmn_remote.RMNRemote |
		*owner_helpers.ManyChainMultiSig |
		*owner_helpers.RBACTimelock |
		*offramp.OffRamp |
		*onramp.OnRamp |
		*burn_mint_erc677.BurnMintERC677 |
		*maybe_revert_message_receiver.MaybeRevertMessageReceiver |
		*aggregator_v3_interface.AggregatorV3Interface
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
	_, err := chain.Confirm(contractDeploy.Tx)
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
	HomeChainSel   uint64
	FeedChainSel   uint64
	ChainsToDeploy []uint64
	TokenConfig    TokenConfig
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
		return ab, fmt.Errorf("capability registry not found for home chain %d, needs to be deployed first", c.HomeChainSel)
	}
	cr, err := c.Chains[c.HomeChainSel].CapabilityRegistry.GetHashedCapabilityId(
		&bind.CallOpts{}, CapabilityLabelledName, CapabilityVersion)
	if err != nil {
		e.Logger.Errorw("Failed to get hashed capability id", "err", err)
		return ab, err
	}
	if cr != CCIPCapabilityID {
		return ab, fmt.Errorf("Capability registry does not support CCIP %s %s", hexutil.Encode(cr[:]), hexutil.Encode(CCIPCapabilityID[:]))
	}
	// Signal to CR that our nodes support CCIP capability.
	if err := AddNodes(
		c.Chains[c.HomeChainSel].CapabilityRegistry,
		e.Chains[c.HomeChainSel],
		nodes.NonBootstraps().PeerIDs(),
	); err != nil {
		return ab, err
	}

	for _, chainSel := range c.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return ab, fmt.Errorf("Chain %d not found", chainSel)
		}
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

		tokenInfo := c.TokenConfig.GetTokenInfo(e.Logger, chainState)
		// TODO: Do we want to extract this?
		// Add chain config for each chain.
		_, err = AddChainConfig(
			e.Logger,
			e.Chains[c.HomeChainSel],
			c.Chains[c.HomeChainSel].CCIPConfig,
			chain.Selector,
			nodes.NonBootstraps().PeerIDs())
		if err != nil {
			return ab, err
		}

		// For each chain, we create a DON on the home chain (2 OCR instances)
		if err := AddDON(
			e.Logger,
			c.Chains[c.HomeChainSel].CapabilityRegistry,
			c.Chains[c.HomeChainSel].CCIPConfig,
			chainState.OffRamp,
			c.FeedChainSel,
			tokenInfo,
			chain,
			e.Chains[c.HomeChainSel],
			nodes.NonBootstraps(),
		); err != nil {
			e.Logger.Errorw("Failed to add DON", "err", err)
			return ab, err
		}
	}

	return ab, nil
}

func DeployChainContracts(
	e deployment.Environment,
	chain deployment.Chain,
	ab deployment.AddressBook,
) (deployment.AddressBook, error) {
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
	rmnRemote, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*rmn_remote.RMNRemote] {
			rmnRemoteAddr, tx, rmnRemote, err2 := rmn_remote.DeployRMNRemote(
				chain.DeployerKey,
				chain.Client,
				chain.Selector,
			)
			return ContractDeploy[*rmn_remote.RMNRemote]{
				rmnRemoteAddr, rmnRemote, tx, deployment.NewTypeAndVersion(RMNRemote, deployment.Version1_6_0_dev), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy RMNRemote", "err", err)
		return ab, err
	}
	e.Logger.Infow("deployed RMNRemote", "addr", rmnRemote.Address)

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
	// TODO: Parameterize this.
	e.Logger.Infow("deployed mcm", "addr", mcm.Address)
	publicKey := TestXXXMCMSSigner.Public().(*ecdsa.PublicKey)
	// Convert the public key to an Ethereum address
	address := crypto.PubkeyToAddress(*publicKey)
	c, err := configwrappers.NewConfig(1, []common.Address{address}, []configwrappers.Config{})
	if err != nil {
		e.Logger.Errorw("Failed to create config", "err", err)
		return ab, err
	}
	groupQuorums, groupParents, signerAddresses, signerGroups := c.ExtractSetConfigInputs()
	mcmsTx, err := mcm.Contract.SetConfig(chain.DeployerKey,
		signerAddresses,
		signerGroups, // Signer 1 is int group 0 (root group) with quorum 1.
		groupQuorums,
		groupParents,
		false,
	)
	if _, err := deployment.ConfirmIfNoError(chain, mcmsTx, err); err != nil {
		e.Logger.Errorw("Failed to confirm mcm config", "err", err)
		return ab, err
	}

	timelock, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*owner_helpers.RBACTimelock] {
			timelock, tx, cc, err2 := owner_helpers.DeployRBACTimelock(
				chain.DeployerKey,
				chain.Client,
				big.NewInt(0), // minDelay
				mcm.Address,
				// TODO: Actual MCM groups need to be parameterized.
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
				rmnRemote.Address,
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

	// TODO: Need general configuration for using pre-existing weth9
	// link tokens.
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

	testRouterContract, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*router.Router] {
			routerAddr, tx, routerC, err2 := router.DeployRouter(
				chain.DeployerKey,
				chain.Client,
				weth9.Address,
				rmnProxy.Address,
			)
			return ContractDeploy[*router.Router]{
				routerAddr, routerC, tx, deployment.NewTypeAndVersion(TestRouter, deployment.Version1_2_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy test router", "err", err)
		return ab, err
	}
	e.Logger.Infow("deployed test router", "addr", testRouterContract.Address)

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

	feeQuoter, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*fee_quoter.FeeQuoter] {
			prAddr, tx, pr, err2 := fee_quoter.DeployFeeQuoter(
				chain.DeployerKey,
				chain.Client,
				fee_quoter.FeeQuoterStaticConfig{
					MaxFeeJuelsPerMsg:  big.NewInt(0).Mul(big.NewInt(2e2), big.NewInt(1e18)),
					LinkToken:          linkToken.Address,
					StalenessThreshold: uint32(24 * 60 * 60),
				},
				[]common.Address{timelock.Address},                 // timelock should be able to update, ramps added after
				[]common.Address{weth9.Address, linkToken.Address}, // fee tokens
				[]fee_quoter.FeeQuoterTokenPriceFeedUpdate{},
				[]fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{}, // TODO: tokens
				[]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs{
					{
						PremiumMultiplierWeiPerEth: 9e17, // 0.9 ETH
						Token:                      linkToken.Address,
					},
					{
						PremiumMultiplierWeiPerEth: 1e18,
						Token:                      weth9.Address,
					},
				},
				[]fee_quoter.FeeQuoterDestChainConfigArgs{},
			)
			return ContractDeploy[*fee_quoter.FeeQuoter]{
				prAddr, pr, tx, deployment.NewTypeAndVersion(FeeQuoter, deployment.Version1_6_0_dev), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy fee quoter", "err", err)
		return ab, err
	}

	onRamp, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*onramp.OnRamp] {
			onRampAddr, tx, onRamp, err2 := onramp.DeployOnRamp(
				chain.DeployerKey,
				chain.Client,
				onramp.OnRampStaticConfig{
					ChainSelector:      chain.Selector,
					Rmn:                rmnProxy.Address,
					NonceManager:       nonceManager.Address,
					TokenAdminRegistry: tokenAdminRegistry.Address,
				},
				onramp.OnRampDynamicConfig{
					FeeQuoter:     feeQuoter.Address,
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
					Rmn:                rmnProxy.Address,
					NonceManager:       nonceManager.Address,
					TokenAdminRegistry: tokenAdminRegistry.Address,
				},
				offramp.OffRampDynamicConfig{
					FeeQuoter:                               feeQuoter.Address,
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

	// Basic wiring is always needed.
	tx, err := feeQuoter.Contract.ApplyAuthorizedCallerUpdates(chain.DeployerKey, fee_quoter.AuthorizedCallersAuthorizedCallerArgs{
		// TODO: We enable the deployer initially to set prices
		// Should be removed after.
		AddedCallers: []common.Address{offRamp.Contract.Address(), chain.DeployerKey.From},
	})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		e.Logger.Errorw("Failed to confirm fee quoter authorized caller update", "err", err)
		return ab, err
	}

	tx, err = nonceManager.Contract.ApplyAuthorizedCallerUpdates(chain.DeployerKey, nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
		AddedCallers: []common.Address{offRamp.Contract.Address(), onRamp.Contract.Address()},
	})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		e.Logger.Errorw("Failed to update nonce manager with ramps", "err", err)
		return ab, err
	}
	return ab, nil
}
