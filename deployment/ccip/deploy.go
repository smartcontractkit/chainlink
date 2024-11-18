package ccipdeployment

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
)

var (
	MockRMN                    deployment.ContractType = "MockRMN"
	RMNRemote                  deployment.ContractType = "RMNRemote"
	LinkToken                  deployment.ContractType = "LinkToken"
	ARMProxy                   deployment.ContractType = "ARMProxy"
	WETH9                      deployment.ContractType = "WETH9"
	Router                     deployment.ContractType = "Router"
	CommitStore                deployment.ContractType = "CommitStore"
	TokenAdminRegistry         deployment.ContractType = "TokenAdminRegistry"
	RegistryModule             deployment.ContractType = "RegistryModuleOwnerCustom"
	NonceManager               deployment.ContractType = "NonceManager"
	FeeQuoter                  deployment.ContractType = "FeeQuoter"
	AdminManyChainMultisig     deployment.ContractType = "AdminManyChainMultiSig"
	BypasserManyChainMultisig  deployment.ContractType = "BypasserManyChainMultiSig"
	CancellerManyChainMultisig deployment.ContractType = "CancellerManyChainMultiSig"
	ProposerManyChainMultisig  deployment.ContractType = "ProposerManyChainMultiSig"
	CCIPHome                   deployment.ContractType = "CCIPHome"
	CCIPConfig                 deployment.ContractType = "CCIPConfig"
	RMNHome                    deployment.ContractType = "RMNHome"
	RBACTimelock               deployment.ContractType = "RBACTimelock"
	OnRamp                     deployment.ContractType = "OnRamp"
	OffRamp                    deployment.ContractType = "OffRamp"
	CapabilitiesRegistry       deployment.ContractType = "CapabilitiesRegistry"
	PriceFeed                  deployment.ContractType = "PriceFeed"
	// Note test router maps to a regular router contract.
	TestRouter          deployment.ContractType = "TestRouter"
	CCIPReceiver        deployment.ContractType = "CCIPReceiver"
	BurnMintToken       deployment.ContractType = "BurnMintToken"
	BurnMintTokenPool   deployment.ContractType = "BurnMintTokenPool"
	USDCToken           deployment.ContractType = "USDCToken"
	USDCMockTransmitter deployment.ContractType = "USDCMockTransmitter"
	USDCTokenMessenger  deployment.ContractType = "USDCTokenMessenger"
	USDCTokenPool       deployment.ContractType = "USDCTokenPool"
)

func DeployPrerequisiteChainContracts(e deployment.Environment, ab deployment.AddressBook, selectors []uint64) error {
	state, err := LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err")
		return err
	}
	for _, sel := range selectors {
		chain := e.Chains[sel]
		err = DeployPrerequisiteContracts(e, ab, state, chain)
		if err != nil {
			return errors.Wrapf(err, "failed to deploy prerequisite contracts for chain %d", sel)
		}
	}
	return nil
}

// DeployPrerequisiteContracts deploys the contracts that can be ported from previous CCIP version to the new one.
// This is only required for staging and test environments where the contracts are not already deployed.
func DeployPrerequisiteContracts(e deployment.Environment, ab deployment.AddressBook, state CCIPOnChainState, chain deployment.Chain) error {
	lggr := e.Logger
	chainState, chainExists := state.Chains[chain.Selector]
	var weth9Contract *weth9.WETH9
	var linkTokenContract *burn_mint_erc677.BurnMintERC677
	var tokenAdminReg *token_admin_registry.TokenAdminRegistry
	var registryModule *registry_module_owner_custom.RegistryModuleOwnerCustom
	var rmnProxy *rmn_proxy_contract.RMNProxyContract
	var r *router.Router
	if chainExists {
		weth9Contract = chainState.Weth9
		linkTokenContract = chainState.LinkToken
		tokenAdminReg = chainState.TokenAdminRegistry
		registryModule = chainState.RegistryModule
		rmnProxy = chainState.RMNProxyExisting
		r = chainState.Router
	}
	if rmnProxy == nil {
		// we want to replicate the mainnet scenario where RMNProxy is already deployed with some existing RMN
		// This will need us to use two different RMNProxy contracts
		// 1. RMNProxyNew with RMNRemote - ( deployed later in chain contracts)
		// 2. RMNProxyExisting with mockRMN - ( deployed here, replicating the behavior of existing RMNProxy with already set RMN)
		rmn, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*mock_rmn_contract.MockRMNContract] {
				rmnAddr, tx2, rmn, err2 := mock_rmn_contract.DeployMockRMNContract(
					chain.DeployerKey,
					chain.Client,
				)
				return deployment.ContractDeploy[*mock_rmn_contract.MockRMNContract]{
					rmnAddr, rmn, tx2, deployment.NewTypeAndVersion(MockRMN, deployment.Version1_0_0), err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy mock RMN", "err", err)
			return err
		}
		lggr.Infow("deployed mock RMN", "addr", rmn.Address)
		rmnProxyContract, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*rmn_proxy_contract.RMNProxyContract] {
				rmnProxyAddr, tx2, rmnProxy, err2 := rmn_proxy_contract.DeployRMNProxyContract(
					chain.DeployerKey,
					chain.Client,
					rmn.Address,
				)
				return deployment.ContractDeploy[*rmn_proxy_contract.RMNProxyContract]{
					rmnProxyAddr, rmnProxy, tx2, deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_0_0), err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy RMNProxyNew", "err", err)
			return err
		}
		lggr.Infow("deployed RMNProxyNew", "addr", rmnProxyContract.Address)
		rmnProxy = rmnProxyContract.Contract
	}
	if tokenAdminReg == nil {
		tokenAdminRegistry, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*token_admin_registry.TokenAdminRegistry] {
				tokenAdminRegistryAddr, tx2, tokenAdminRegistry, err2 := token_admin_registry.DeployTokenAdminRegistry(
					chain.DeployerKey,
					chain.Client)
				return deployment.ContractDeploy[*token_admin_registry.TokenAdminRegistry]{
					tokenAdminRegistryAddr, tokenAdminRegistry, tx2, deployment.NewTypeAndVersion(TokenAdminRegistry, deployment.Version1_5_0), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy token admin registry", "err", err)
			return err
		}
		e.Logger.Infow("deployed tokenAdminRegistry", "addr", tokenAdminRegistry)
		tokenAdminReg = tokenAdminRegistry.Contract
	} else {
		e.Logger.Infow("tokenAdminRegistry already deployed", "addr", tokenAdminReg.Address)
	}
	if registryModule == nil {
		customRegistryModule, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom] {
				regModAddr, tx2, regMod, err2 := registry_module_owner_custom.DeployRegistryModuleOwnerCustom(
					chain.DeployerKey,
					chain.Client,
					tokenAdminReg.Address())
				return deployment.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom]{
					regModAddr, regMod, tx2, deployment.NewTypeAndVersion(RegistryModule, deployment.Version1_5_0), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy custom registry module", "err", err)
			return err
		}
		e.Logger.Infow("deployed custom registry module", "addr", customRegistryModule)
		registryModule = customRegistryModule.Contract
	} else {
		e.Logger.Infow("custom registry module already deployed", "addr", registryModule.Address)
	}
	isRegistryAdded, err := tokenAdminReg.IsRegistryModule(nil, registryModule.Address())
	if err != nil {
		e.Logger.Errorw("Failed to check if registry module is added on token admin registry", "err", err)
		return fmt.Errorf("failed to check if registry module is added on token admin registry: %w", err)
	}
	if !isRegistryAdded {
		tx, err := tokenAdminReg.AddRegistryModule(chain.DeployerKey, registryModule.Address())
		if err != nil {
			e.Logger.Errorw("Failed to assign registry module on token admin registry", "err", err)
			return fmt.Errorf("failed to assign registry module on token admin registry: %w", err)
		}

		_, err = chain.Confirm(tx)
		if err != nil {
			e.Logger.Errorw("Failed to confirm assign registry module on token admin registry", "err", err)
			return fmt.Errorf("failed to confirm assign registry module on token admin registry: %w", err)
		}
		e.Logger.Infow("assigned registry module on token admin registry")
	}
	if weth9Contract == nil {
		weth, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*weth9.WETH9] {
				weth9Addr, tx2, weth9c, err2 := weth9.DeployWETH9(
					chain.DeployerKey,
					chain.Client,
				)
				return deployment.ContractDeploy[*weth9.WETH9]{
					weth9Addr, weth9c, tx2, deployment.NewTypeAndVersion(WETH9, deployment.Version1_0_0), err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy weth9", "err", err)
			return err
		}
		lggr.Infow("deployed weth9", "addr", weth.Address)
		weth9Contract = weth.Contract
	} else {
		lggr.Infow("weth9 already deployed", "addr", weth9Contract.Address)
	}
	if linkTokenContract == nil {
		linkToken, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
				linkTokenAddr, tx2, linkToken, err2 := burn_mint_erc677.DeployBurnMintERC677(
					chain.DeployerKey,
					chain.Client,
					"Link Token",
					"LINK",
					uint8(18),
					big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
				)
				return deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
					linkTokenAddr, linkToken, tx2, deployment.NewTypeAndVersion(LinkToken, deployment.Version1_0_0), err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy linkToken", "err", err)
			return err
		}
		lggr.Infow("deployed linkToken", "addr", linkToken.Address)
		linkTokenContract = linkToken.Contract
	} else {
		lggr.Infow("linkToken already deployed", "addr", linkTokenContract.Address)
	}
	// if router is not already deployed, we deploy it
	if r == nil {
		routerContract, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*router.Router] {
				routerAddr, tx2, routerC, err2 := router.DeployRouter(
					chain.DeployerKey,
					chain.Client,
					weth9Contract.Address(),
					rmnProxy.Address(),
				)
				return deployment.ContractDeploy[*router.Router]{
					routerAddr, routerC, tx2, deployment.NewTypeAndVersion(Router, deployment.Version1_2_0), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy router", "err", err)
			return err
		}
		e.Logger.Infow("deployed router", "addr", routerContract.Address)
		r = routerContract.Contract
	} else {
		e.Logger.Infow("router already deployed", "addr", chainState.Router.Address)
	}
	return nil
}

type USDCConfig struct {
	Enabled bool
	USDCAttestationConfig
}

type USDCAttestationConfig struct {
	API         string
	APITimeout  *commonconfig.Duration
	APIInterval *commonconfig.Duration
}

type DeployCCIPContractConfig struct {
	HomeChainSel   uint64
	FeedChainSel   uint64
	ChainsToDeploy []uint64
	TokenConfig    TokenConfig
	// I believe it makes sense to have the same signers across all chains
	// since that's the point MCMS.
	MCMSConfig MCMSConfig
	USDCConfig USDCConfig
	// For setting OCR configuration
	OCRSecrets deployment.OCRSecrets
}

// DeployCCIPContracts assumes the following contracts are deployed:
// - Capability registry
// - CCIP home
// - RMN home
// - Fee tokens on all chains.
// and present in ExistingAddressBook.
// It then deploys the rest of the CCIP chain contracts to the selected chains
// registers the nodes with the capability registry and creates a DON for
// each new chain. TODO: Might be better to break this down a bit?
func DeployCCIPContracts(e deployment.Environment, ab deployment.AddressBook, c DeployCCIPContractConfig) error {
	if c.OCRSecrets.IsEmpty() {
		return fmt.Errorf("OCR secrets are empty")
	}
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil || len(nodes) == 0 {
		e.Logger.Errorw("Failed to get node info", "err", err)
		return err
	}
	existingState, err := LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err")
		return err
	}
	capReg := existingState.Chains[c.HomeChainSel].CapabilityRegistry
	if capReg == nil {
		e.Logger.Errorw("Failed to get capability registry")
		return fmt.Errorf("capability registry not found")
	}
	ccipHome := existingState.Chains[c.HomeChainSel].CCIPHome
	if ccipHome == nil {
		e.Logger.Errorw("Failed to get ccip home", "err", err)
		return fmt.Errorf("ccip home not found")
	}
	rmnHome := existingState.Chains[c.HomeChainSel].RMNHome
	if rmnHome == nil {
		e.Logger.Errorw("Failed to get rmn home", "err", err)
		return fmt.Errorf("rmn home not found")
	}

	usdcConfiguration := make(map[cciptypes.ChainSelector]pluginconfig.USDCCCTPTokenConfig)
	for _, chainSel := range c.ChainsToDeploy {
		chain, exists := e.Chains[chainSel]
		if !exists {
			return fmt.Errorf("chain %d not found", chainSel)
		}
		if c.USDCConfig.Enabled {
			token, pool, messenger, transmitter, err1 := DeployUSDC(e.Logger, chain, ab, existingState.Chains[chainSel])
			if err1 != nil {
				return err1
			}
			e.Logger.Infow("Deployed USDC contracts",
				"chainSelector", chainSel,
				"token", token.Address(),
				"pool", pool.Address(),
				"transmitter", transmitter.Address(),
				"messenger", messenger.Address(),
			)

			usdcConfiguration[cciptypes.ChainSelector(chainSel)] = pluginconfig.USDCCCTPTokenConfig{
				SourcePoolAddress:            pool.Address().Hex(),
				SourceMessageTransmitterAddr: transmitter.Address().Hex(),
			}
		}
	}
	err = DeployChainContractsForChains(e, ab, c.HomeChainSel, c.ChainsToDeploy, c.MCMSConfig)
	if err != nil {
		e.Logger.Errorw("Failed to deploy chain contracts", "err", err)
		return err
	}
	for _, chainSel := range c.ChainsToDeploy {
		chain, _ := e.Chains[chainSel]
		chainAddresses, err := ab.AddressesForChain(chain.Selector)
		if err != nil {
			e.Logger.Errorw("Failed to get chain addresses", "err", err)
			return err
		}
		chainState, err := LoadChainState(chain, chainAddresses)
		if err != nil {
			e.Logger.Errorw("Failed to load chain state", "err", err)
			return err
		}

		tokenInfo := c.TokenConfig.GetTokenInfo(e.Logger, existingState.Chains[chainSel].LinkToken, existingState.Chains[chainSel].Weth9)
		// TODO: Do we want to extract this?
		// Add chain config for each chain.
		_, err = AddChainConfig(
			e.Logger,
			e.Chains[c.HomeChainSel],
			ccipHome,
			chain.Selector,
			nodes.NonBootstraps().PeerIDs())
		if err != nil {
			return err
		}
		var tokenDataObserversConf []pluginconfig.TokenDataObserverConfig
		if c.USDCConfig.Enabled {
			tokenDataObserversConf = []pluginconfig.TokenDataObserverConfig{{
				Type:    pluginconfig.USDCCCTPHandlerType,
				Version: "1.0",
				USDCCCTPObserverConfig: &pluginconfig.USDCCCTPObserverConfig{
					Tokens:                 usdcConfiguration,
					AttestationAPI:         c.USDCConfig.API,
					AttestationAPITimeout:  c.USDCConfig.APITimeout,
					AttestationAPIInterval: c.USDCConfig.APIInterval,
				},
			}}
		}
		// For each chain, we create a DON on the home chain (2 OCR instances)
		if err := AddDON(
			e.Logger,
			c.OCRSecrets,
			capReg,
			ccipHome,
			rmnHome.Address(),
			chainState.OffRamp,
			c.FeedChainSel,
			tokenInfo,
			chain,
			e.Chains[c.HomeChainSel],
			nodes.NonBootstraps(),
			tokenDataObserversConf,
		); err != nil {
			e.Logger.Errorw("Failed to add DON", "err", err)
			return err
		}
	}

	return nil
}

type MCMSConfig struct {
	Admin     config.Config
	Canceller config.Config
	Bypasser  config.Config
	Proposer  config.Config
	Executors []common.Address
}

func DeployMCMSWithConfig(
	contractType deployment.ContractType,
	lggr logger.Logger,
	chain deployment.Chain,
	ab deployment.AddressBook,
	mcmConfig config.Config,
) (*deployment.ContractDeploy[*owner_helpers.ManyChainMultiSig], error) {
	groupQuorums, groupParents, signerAddresses, signerGroups := mcmConfig.ExtractSetConfigInputs()
	mcm, err := deployment.DeployContract(lggr, chain, ab,
		func(chain deployment.Chain) deployment.ContractDeploy[*owner_helpers.ManyChainMultiSig] {
			mcmAddr, tx, mcm, err2 := owner_helpers.DeployManyChainMultiSig(
				chain.DeployerKey,
				chain.Client,
			)
			return deployment.ContractDeploy[*owner_helpers.ManyChainMultiSig]{
				mcmAddr, mcm, tx, deployment.NewTypeAndVersion(contractType, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy mcm", "err", err)
		return mcm, err
	}
	mcmsTx, err := mcm.Contract.SetConfig(chain.DeployerKey,
		signerAddresses,
		signerGroups, // Signer 1 is int group 0 (root group) with quorum 1.
		groupQuorums,
		groupParents,
		false,
	)
	if _, err := deployment.ConfirmIfNoError(chain, mcmsTx, err); err != nil {
		lggr.Errorw("Failed to confirm mcm config", "err", err)
		return mcm, err
	}
	return mcm, nil
}

type MCMSContracts struct {
	Admin     *deployment.ContractDeploy[*owner_helpers.ManyChainMultiSig]
	Canceller *deployment.ContractDeploy[*owner_helpers.ManyChainMultiSig]
	Bypasser  *deployment.ContractDeploy[*owner_helpers.ManyChainMultiSig]
	Proposer  *deployment.ContractDeploy[*owner_helpers.ManyChainMultiSig]
	Timelock  *deployment.ContractDeploy[*owner_helpers.RBACTimelock]
}

// DeployMCMSContracts deploys the MCMS contracts for the given configuration
// as well as the timelock.
func DeployMCMSContracts(
	lggr logger.Logger,
	chain deployment.Chain,
	ab deployment.AddressBook,
	mcmConfig MCMSConfig,
) (*MCMSContracts, error) {
	adminMCM, err := DeployMCMSWithConfig(AdminManyChainMultisig, lggr, chain, ab, mcmConfig.Admin)
	if err != nil {
		return nil, err
	}
	bypasser, err := DeployMCMSWithConfig(BypasserManyChainMultisig, lggr, chain, ab, mcmConfig.Bypasser)
	if err != nil {
		return nil, err
	}
	canceller, err := DeployMCMSWithConfig(CancellerManyChainMultisig, lggr, chain, ab, mcmConfig.Canceller)
	if err != nil {
		return nil, err
	}
	proposer, err := DeployMCMSWithConfig(ProposerManyChainMultisig, lggr, chain, ab, mcmConfig.Proposer)
	if err != nil {
		return nil, err
	}

	timelock, err := deployment.DeployContract(lggr, chain, ab,
		func(chain deployment.Chain) deployment.ContractDeploy[*owner_helpers.RBACTimelock] {
			timelock, tx2, cc, err2 := owner_helpers.DeployRBACTimelock(
				chain.DeployerKey,
				chain.Client,
				big.NewInt(0), // minDelay
				adminMCM.Address,
				[]common.Address{proposer.Address},  // proposers
				mcmConfig.Executors,                 //executors
				[]common.Address{canceller.Address}, // cancellers
				[]common.Address{bypasser.Address},  // bypassers
			)
			return deployment.ContractDeploy[*owner_helpers.RBACTimelock]{
				timelock, cc, tx2, deployment.NewTypeAndVersion(RBACTimelock, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy timelock", "err", err)
		return nil, err
	}
	lggr.Infow("deployed timelock", "addr", timelock.Address)
	return &MCMSContracts{
		Admin:     adminMCM,
		Canceller: canceller,
		Bypasser:  bypasser,
		Proposer:  proposer,
		Timelock:  timelock,
	}, nil
}

func DeployChainContractsForChains(e deployment.Environment, ab deployment.AddressBook, homeChainSel uint64, chainsToDeploy []uint64, mcmsConfig MCMSConfig) error {
	existingState, err := LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err")
		return err
	}

	capReg := existingState.Chains[homeChainSel].CapabilityRegistry
	if capReg == nil {
		e.Logger.Errorw("Failed to get capability registry")
		return fmt.Errorf("capability registry not found")
	}
	cr, err := capReg.GetHashedCapabilityId(
		&bind.CallOpts{}, CapabilityLabelledName, CapabilityVersion)
	if err != nil {
		e.Logger.Errorw("Failed to get hashed capability id", "err", err)
		return err
	}
	if cr != CCIPCapabilityID {
		return fmt.Errorf("capability registry does not support CCIP %s %s", hexutil.Encode(cr[:]), hexutil.Encode(CCIPCapabilityID[:]))
	}
	capability, err := capReg.GetCapability(nil, CCIPCapabilityID)
	if err != nil {
		e.Logger.Errorw("Failed to get capability", "err", err)
		return err
	}
	ccipHome, err := ccip_home.NewCCIPHome(capability.ConfigurationContract, e.Chains[homeChainSel].Client)
	if err != nil {
		e.Logger.Errorw("Failed to get ccip config", "err", err)
		return err
	}
	if ccipHome.Address() != existingState.Chains[homeChainSel].CCIPHome.Address() {
		return fmt.Errorf("ccip home address mismatch")
	}
	rmnHome := existingState.Chains[homeChainSel].RMNHome
	if rmnHome == nil {
		e.Logger.Errorw("Failed to get rmn home", "err", err)
		return fmt.Errorf("rmn home not found")
	}
	for _, chainSel := range chainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain %d not found", chainSel)
		}
		if existingState.Chains[chainSel].LinkToken == nil || existingState.Chains[chainSel].Weth9 == nil {
			return fmt.Errorf("fee tokens not found for chain %d", chainSel)
		}
		err := DeployChainContracts(e, chain, ab, mcmsConfig, rmnHome)
		if err != nil {
			e.Logger.Errorw("Failed to deploy chain contracts", "chain", chainSel, "err", err)
			return fmt.Errorf("failed to deploy chain contracts for chain %d: %w", chainSel, err)
		}
	}
	return nil
}

func DeployChainContracts(
	e deployment.Environment,
	chain deployment.Chain,
	ab deployment.AddressBook,
	mcmsConfig MCMSConfig,
	rmnHome *rmn_home.RMNHome,
) error {
	mcmsContracts, err := DeployMCMSContracts(e.Logger, chain, ab, mcmsConfig)
	if err != nil {
		return err
	}
	// check for existing contracts
	state, err := LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err")
		return err
	}
	chainState, chainExists := state.Chains[chain.Selector]
	if !chainExists {
		return fmt.Errorf("chain %d not found in existing state, deploy the prerequisites first", chain.Selector)
	}
	if chainState.Weth9 == nil {
		return fmt.Errorf("weth9 not found for chain %d, deploy the prerequisites first", chain.Selector)
	}
	weth9Contract := chainState.Weth9
	if chainState.LinkToken == nil {
		return fmt.Errorf("link token not found for chain %d, deploy the prerequisites first", chain.Selector)
	}
	linkTokenContract := chainState.LinkToken
	if chainState.TokenAdminRegistry == nil {
		return fmt.Errorf("token admin registry not found for chain %d, deploy the prerequisites first", chain.Selector)
	}
	tokenAdminReg := chainState.TokenAdminRegistry
	if chainState.RegistryModule == nil {
		return fmt.Errorf("registry module not found for chain %d, deploy the prerequisites first", chain.Selector)
	}
	if chainState.Router == nil {
		return fmt.Errorf("router not found for chain %d, deploy the prerequisites first", chain.Selector)
	}
	if chainState.Receiver == nil {
		ccipReceiver, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*maybe_revert_message_receiver.MaybeRevertMessageReceiver] {
				receiverAddr, tx, receiver, err2 := maybe_revert_message_receiver.DeployMaybeRevertMessageReceiver(
					chain.DeployerKey,
					chain.Client,
					false,
				)
				return deployment.ContractDeploy[*maybe_revert_message_receiver.MaybeRevertMessageReceiver]{
					receiverAddr, receiver, tx, deployment.NewTypeAndVersion(CCIPReceiver, deployment.Version1_0_0), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy receiver", "err", err)
			return err
		}
		e.Logger.Infow("deployed receiver", "addr", ccipReceiver.Address)
	} else {
		e.Logger.Infow("receiver already deployed", "addr", chainState.Receiver.Address)
	}
	rmnRemoteContract := chainState.RMNRemote
	if chainState.RMNRemote == nil {
		// TODO: Correctly configure RMN remote.
		rmnRemote, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*rmn_remote.RMNRemote] {
				rmnRemoteAddr, tx, rmnRemote, err2 := rmn_remote.DeployRMNRemote(
					chain.DeployerKey,
					chain.Client,
					chain.Selector,
				)
				return deployment.ContractDeploy[*rmn_remote.RMNRemote]{
					rmnRemoteAddr, rmnRemote, tx, deployment.NewTypeAndVersion(RMNRemote, deployment.Version1_6_0_dev), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy RMNRemote", "err", err)
			return err
		}
		e.Logger.Infow("deployed RMNRemote", "addr", rmnRemote.Address)
		rmnRemoteContract = rmnRemote.Contract
	} else {
		e.Logger.Infow("rmn remote already deployed", "addr", chainState.RMNRemote.Address)
	}
	activeDigest, err := rmnHome.GetActiveDigest(&bind.CallOpts{})
	if err != nil {
		e.Logger.Errorw("Failed to get active digest", "err", err)
		return err
	}
	e.Logger.Infow("setting active home digest to rmn remote", "digest", activeDigest)

	tx, err := rmnRemoteContract.SetConfig(chain.DeployerKey, rmn_remote.RMNRemoteConfig{
		RmnHomeContractConfigDigest: activeDigest,
		Signers: []rmn_remote.RMNRemoteSigner{
			{NodeIndex: 0, OnchainPublicKey: common.Address{1}},
		},
		F: 0, // TODO: update when we have signers
	})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		e.Logger.Errorw("Failed to confirm RMNRemote config", "err", err)
		return err
	}

	// we deploy a new RMNProxy so that RMNRemote can be tested first before pointing it to the main Existing RMNProxy
	// To differentiate between the two RMNProxies, we will deploy new one with Version1_6_0_dev
	rmnProxyContract := chainState.RMNProxyNew
	if chainState.RMNProxyNew == nil {
		// we deploy a new rmnproxy contract to test RMNRemote
		rmnProxy, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*rmn_proxy_contract.RMNProxyContract] {
				rmnProxyAddr, tx, rmnProxy, err2 := rmn_proxy_contract.DeployRMNProxyContract(
					chain.DeployerKey,
					chain.Client,
					rmnRemoteContract.Address(),
				)
				return deployment.ContractDeploy[*rmn_proxy_contract.RMNProxyContract]{
					rmnProxyAddr, rmnProxy, tx, deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_6_0_dev), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy RMNProxyNew", "err", err)
			return err
		}
		e.Logger.Infow("deployed new RMNProxyNew", "addr", rmnProxy.Address)
		rmnProxyContract = rmnProxy.Contract
	} else {
		e.Logger.Infow("rmn proxy already deployed", "addr", chainState.RMNProxyNew.Address)
	}
	if chainState.TestRouter == nil {
		testRouterContract, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*router.Router] {
				routerAddr, tx2, routerC, err2 := router.DeployRouter(
					chain.DeployerKey,
					chain.Client,
					weth9Contract.Address(),
					rmnProxyContract.Address(),
				)
				return deployment.ContractDeploy[*router.Router]{
					routerAddr, routerC, tx2, deployment.NewTypeAndVersion(TestRouter, deployment.Version1_2_0), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy test router", "err", err)
			return err
		}
		e.Logger.Infow("deployed test router", "addr", testRouterContract.Address)
	} else {
		e.Logger.Infow("test router already deployed", "addr", chainState.TestRouter.Address)
	}

	nmContract := chainState.NonceManager
	if chainState.NonceManager == nil {
		nonceManager, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*nonce_manager.NonceManager] {
				nonceManagerAddr, tx2, nonceManager, err2 := nonce_manager.DeployNonceManager(
					chain.DeployerKey,
					chain.Client,
					[]common.Address{}, // Need to add onRamp after
				)
				return deployment.ContractDeploy[*nonce_manager.NonceManager]{
					nonceManagerAddr, nonceManager, tx2, deployment.NewTypeAndVersion(NonceManager, deployment.Version1_6_0_dev), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy nonce manager", "err", err)
			return err
		}
		e.Logger.Infow("Deployed nonce manager", "addr", nonceManager.Address)
		nmContract = nonceManager.Contract
	} else {
		e.Logger.Infow("nonce manager already deployed", "addr", chainState.NonceManager.Address)
	}
	feeQuoterContract := chainState.FeeQuoter
	if chainState.FeeQuoter == nil {
		feeQuoter, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*fee_quoter.FeeQuoter] {
				prAddr, tx2, pr, err2 := fee_quoter.DeployFeeQuoter(
					chain.DeployerKey,
					chain.Client,
					fee_quoter.FeeQuoterStaticConfig{
						MaxFeeJuelsPerMsg:            big.NewInt(0).Mul(big.NewInt(2e2), big.NewInt(1e18)),
						LinkToken:                    linkTokenContract.Address(),
						TokenPriceStalenessThreshold: uint32(24 * 60 * 60),
					},
					[]common.Address{mcmsContracts.Timelock.Address},                       // timelock should be able to update, ramps added after
					[]common.Address{weth9Contract.Address(), linkTokenContract.Address()}, // fee tokens
					[]fee_quoter.FeeQuoterTokenPriceFeedUpdate{},
					[]fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{}, // TODO: tokens
					[]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs{
						{
							PremiumMultiplierWeiPerEth: 9e17, // 0.9 ETH
							Token:                      linkTokenContract.Address(),
						},
						{
							PremiumMultiplierWeiPerEth: 1e18,
							Token:                      weth9Contract.Address(),
						},
					},
					[]fee_quoter.FeeQuoterDestChainConfigArgs{},
				)
				return deployment.ContractDeploy[*fee_quoter.FeeQuoter]{
					prAddr, pr, tx2, deployment.NewTypeAndVersion(FeeQuoter, deployment.Version1_6_0_dev), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy fee quoter", "err", err)
			return err
		}
		e.Logger.Infow("Deployed fee quoter", "addr", feeQuoter.Address)
		feeQuoterContract = feeQuoter.Contract
	} else {
		e.Logger.Infow("fee quoter already deployed", "addr", chainState.FeeQuoter.Address)
	}
	onRampContract := chainState.OnRamp
	if onRampContract == nil {
		onRamp, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*onramp.OnRamp] {
				onRampAddr, tx2, onRamp, err2 := onramp.DeployOnRamp(
					chain.DeployerKey,
					chain.Client,
					onramp.OnRampStaticConfig{
						ChainSelector:      chain.Selector,
						RmnRemote:          rmnProxyContract.Address(),
						NonceManager:       nmContract.Address(),
						TokenAdminRegistry: tokenAdminReg.Address(),
					},
					onramp.OnRampDynamicConfig{
						FeeQuoter:     feeQuoterContract.Address(),
						FeeAggregator: common.HexToAddress("0x1"), // TODO real fee aggregator
					},
					[]onramp.OnRampDestChainConfigArgs{},
				)
				return deployment.ContractDeploy[*onramp.OnRamp]{
					onRampAddr, onRamp, tx2, deployment.NewTypeAndVersion(OnRamp, deployment.Version1_6_0_dev), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy onramp", "err", err)
			return err
		}
		e.Logger.Infow("Deployed onramp", "addr", onRamp.Address)
		onRampContract = onRamp.Contract
	} else {
		e.Logger.Infow("onramp already deployed", "addr", chainState.OnRamp.Address)
	}
	offRampContract := chainState.OffRamp
	if offRampContract == nil {
		offRamp, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*offramp.OffRamp] {
				offRampAddr, tx2, offRamp, err2 := offramp.DeployOffRamp(
					chain.DeployerKey,
					chain.Client,
					offramp.OffRampStaticConfig{
						ChainSelector:      chain.Selector,
						RmnRemote:          rmnProxyContract.Address(),
						NonceManager:       nmContract.Address(),
						TokenAdminRegistry: tokenAdminReg.Address(),
					},
					offramp.OffRampDynamicConfig{
						FeeQuoter:                               feeQuoterContract.Address(),
						PermissionLessExecutionThresholdSeconds: uint32(86400),
						IsRMNVerificationDisabled:               true,
					},
					[]offramp.OffRampSourceChainConfigArgs{},
				)
				return deployment.ContractDeploy[*offramp.OffRamp]{
					Address: offRampAddr, Contract: offRamp, Tx: tx2, Tv: deployment.NewTypeAndVersion(OffRamp, deployment.Version1_6_0_dev), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy offramp", "err", err)
			return err
		}
		e.Logger.Infow("Deployed offramp", "addr", offRamp.Address)
		offRampContract = offRamp.Contract
	} else {
		e.Logger.Infow("offramp already deployed", "addr", chainState.OffRamp.Address)
	}
	// Basic wiring is always needed.
	tx, err = feeQuoterContract.ApplyAuthorizedCallerUpdates(chain.DeployerKey, fee_quoter.AuthorizedCallersAuthorizedCallerArgs{
		// TODO: We enable the deployer initially to set prices
		// Should be removed after.
		AddedCallers: []common.Address{offRampContract.Address(), chain.DeployerKey.From},
	})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		e.Logger.Errorw("Failed to confirm fee quoter authorized caller update", "err", err)
		return err
	}

	tx, err = nmContract.ApplyAuthorizedCallerUpdates(chain.DeployerKey, nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
		AddedCallers: []common.Address{offRampContract.Address(), onRampContract.Address()},
	})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		e.Logger.Errorw("Failed to update nonce manager with ramps", "err", err)
		return err
	}
	return nil
}
