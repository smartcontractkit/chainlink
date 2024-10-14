package ccipdeployment

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"

	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
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
	RMNRemote                  deployment.ContractType = "RMNRemote"
	LinkToken                  deployment.ContractType = "LinkToken"
	ARMProxy                   deployment.ContractType = "ARMProxy"
	WETH9                      deployment.ContractType = "WETH9"
	Router                     deployment.ContractType = "Router"
	CommitStore                deployment.ContractType = "CommitStore"
	TokenAdminRegistry         deployment.ContractType = "TokenAdminRegistry"
	NonceManager               deployment.ContractType = "NonceManager"
	FeeQuoter                  deployment.ContractType = "FeeQuoter"
	AdminManyChainMultisig     deployment.ContractType = "AdminManyChainMultiSig"
	BypasserManyChainMultisig  deployment.ContractType = "BypasserManyChainMultiSig"
	CancellerManyChainMultisig deployment.ContractType = "CancellerManyChainMultiSig"
	ProposerManyChainMultisig  deployment.ContractType = "ProposerManyChainMultiSig"
	CCIPHome                   deployment.ContractType = "CCIPHome"
	RMNHome                    deployment.ContractType = "RMNHome"
	RBACTimelock               deployment.ContractType = "RBACTimelock"
	OnRamp                     deployment.ContractType = "OnRamp"
	OffRamp                    deployment.ContractType = "OffRamp"
	CapabilitiesRegistry       deployment.ContractType = "CapabilitiesRegistry"
	PriceFeed                  deployment.ContractType = "PriceFeed"
	// Note test router maps to a regular router contract.
	TestRouter   deployment.ContractType = "TestRouter"
	CCIPReceiver deployment.ContractType = "CCIPReceiver"
)

type Contracts interface {
	*capabilities_registry.CapabilitiesRegistry |
		*rmn_proxy_contract.RMNProxyContract |
		*ccip_home.CCIPHome |
		*rmn_home.RMNHome |
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
	HomeChainSel       uint64
	FeedChainSel       uint64
	ChainsToDeploy     []uint64
	TokenConfig        TokenConfig
	CapabilityRegistry common.Address
	FeeTokenContracts  map[uint64]FeeTokenContracts
	// I believe it makes sense to have the same signers across all chains
	// since that's the point MCMS.
	MCMSConfig MCMSConfig
}

// DeployCCIPContracts assumes that the capability registry and ccip home contracts
// are already deployed (needed as a first step because the chainlink nodes point to them).
// It then deploys
func DeployCCIPContracts(e deployment.Environment, ab deployment.AddressBook, c DeployCCIPContractConfig) error {
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil || len(nodes) == 0 {
		e.Logger.Errorw("Failed to get node info", "err", err)
		return err
	}
	capReg, err := capabilities_registry.NewCapabilitiesRegistry(c.CapabilityRegistry, e.Chains[c.HomeChainSel].Client)
	if err != nil {
		e.Logger.Errorw("Failed to get capability registry", "err", err)
		return err
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
	ccipHome, err := ccip_home.NewCCIPHome(capability.ConfigurationContract, e.Chains[c.HomeChainSel].Client)
	if err != nil {
		e.Logger.Errorw("Failed to get ccip config", "err", err)
		return err
	}

	// Signal to CR that our nodes support CCIP capability.
	if err := AddNodes(
		e.Logger,
		capReg,
		e.Chains[c.HomeChainSel],
		nodes.NonBootstraps().PeerIDs(),
	); err != nil {
		return err
	}

	for _, chainSel := range c.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain %d not found", chainSel)
		}
		chainConfig, ok := c.FeeTokenContracts[chainSel]
		if !ok {
			return fmt.Errorf("chain %d config not found", chainSel)
		}
		err = DeployChainContracts(e, chain, ab, chainConfig, c.MCMSConfig)
		if err != nil {
			return err
		}
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

		tokenInfo := c.TokenConfig.GetTokenInfo(e.Logger, c.FeeTokenContracts[chainSel].LinkToken)
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

		// For each chain, we create a DON on the home chain (2 OCR instances)
		if err := AddDON(
			e.Logger,
			capReg,
			ccipHome,
			chainState.OffRamp,
			c.FeedChainSel,
			tokenInfo,
			chain,
			e.Chains[c.HomeChainSel],
			nodes.NonBootstraps(),
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
) (*ContractDeploy[*owner_helpers.ManyChainMultiSig], error) {
	groupQuorums, groupParents, signerAddresses, signerGroups := mcmConfig.ExtractSetConfigInputs()
	mcm, err := deployContract(lggr, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*owner_helpers.ManyChainMultiSig] {
			mcmAddr, tx, mcm, err2 := owner_helpers.DeployManyChainMultiSig(
				chain.DeployerKey,
				chain.Client,
			)
			return ContractDeploy[*owner_helpers.ManyChainMultiSig]{
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
	Admin     *ContractDeploy[*owner_helpers.ManyChainMultiSig]
	Canceller *ContractDeploy[*owner_helpers.ManyChainMultiSig]
	Bypasser  *ContractDeploy[*owner_helpers.ManyChainMultiSig]
	Proposer  *ContractDeploy[*owner_helpers.ManyChainMultiSig]
	Timelock  *ContractDeploy[*owner_helpers.RBACTimelock]
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

	timelock, err := deployContract(lggr, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*owner_helpers.RBACTimelock] {
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
			return ContractDeploy[*owner_helpers.RBACTimelock]{
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

func DeployFeeTokensToChains(lggr logger.Logger, ab deployment.AddressBook, chains map[uint64]deployment.Chain) (map[uint64]FeeTokenContracts, error) {
	feeTokenContractsByChain := make(map[uint64]FeeTokenContracts)
	for _, chain := range chains {
		feeTokenContracts, err := DeployFeeTokens(lggr, chain, ab)
		if err != nil {
			return nil, err
		}
		feeTokenContractsByChain[chain.Selector] = feeTokenContracts
	}
	return feeTokenContractsByChain, nil
}

// DeployFeeTokens deploys link and weth9. This is _usually_ for test environments only,
// real environments they tend to already exist, but sometimes we still have to deploy them to real chains.
func DeployFeeTokens(lggr logger.Logger, chain deployment.Chain, ab deployment.AddressBook) (FeeTokenContracts, error) {
	weth9, err := deployContract(lggr, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*weth9.WETH9] {
			weth9Addr, tx2, weth9c, err2 := weth9.DeployWETH9(
				chain.DeployerKey,
				chain.Client,
			)
			return ContractDeploy[*weth9.WETH9]{
				weth9Addr, weth9c, tx2, deployment.NewTypeAndVersion(WETH9, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy weth9", "err", err)
		return FeeTokenContracts{}, err
	}
	lggr.Infow("deployed weth9", "addr", weth9.Address)

	linkToken, err := deployContract(lggr, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			linkTokenAddr, tx2, linkToken, err2 := burn_mint_erc677.DeployBurnMintERC677(
				chain.DeployerKey,
				chain.Client,
				"Link Token",
				"LINK",
				uint8(18),
				big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
			)
			return ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
				linkTokenAddr, linkToken, tx2, deployment.NewTypeAndVersion(LinkToken, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy linkToken", "err", err)
		return FeeTokenContracts{}, err
	}
	lggr.Infow("deployed linkToken", "addr", linkToken.Address)
	return FeeTokenContracts{
		LinkToken: linkToken.Contract,
		Weth9:     weth9.Contract,
	}, nil
}

type FeeTokenContracts struct {
	LinkToken *burn_mint_erc677.BurnMintERC677
	Weth9     *weth9.WETH9
}

func DeployChainContracts(
	e deployment.Environment,
	chain deployment.Chain,
	ab deployment.AddressBook,
	contractConfig FeeTokenContracts,
	mcmsConfig MCMSConfig,
) error {
	mcmsContracts, err := DeployMCMSContracts(e.Logger, chain, ab, mcmsConfig)
	if err != nil {
		return err
	}
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
		return err
	}
	e.Logger.Infow("deployed receiver", "addr", ccipReceiver.Address)

	// TODO: Correctly configure RMN remote.
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
		return err
	}
	e.Logger.Infow("deployed RMNRemote", "addr", rmnRemote.Address)

	// TODO: Correctly configure RMN remote with config digest from RMN home.
	tx, err := rmnRemote.Contract.SetConfig(chain.DeployerKey, rmn_remote.RMNRemoteConfig{
		RmnHomeContractConfigDigest: [32]byte{1},
		Signers:                     []rmn_remote.RMNRemoteSigner{},
		MinSigners:                  0, // TODO: update when we have signers
	})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		e.Logger.Errorw("Failed to confirm RMNRemote config", "err", err)
		return err
	}

	rmnProxy, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*rmn_proxy_contract.RMNProxyContract] {
			rmnProxyAddr, tx2, rmnProxy, err2 := rmn_proxy_contract.DeployRMNProxyContract(
				chain.DeployerKey,
				chain.Client,
				rmnRemote.Address,
			)
			return ContractDeploy[*rmn_proxy_contract.RMNProxyContract]{
				rmnProxyAddr, rmnProxy, tx2, deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy rmnProxy", "err", err)
		return err
	}
	e.Logger.Infow("deployed rmnProxy", "addr", rmnProxy.Address)

	routerContract, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*router.Router] {
			routerAddr, tx2, routerC, err2 := router.DeployRouter(
				chain.DeployerKey,
				chain.Client,
				contractConfig.Weth9.Address(),
				rmnProxy.Address,
			)
			return ContractDeploy[*router.Router]{
				routerAddr, routerC, tx2, deployment.NewTypeAndVersion(Router, deployment.Version1_2_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy router", "err", err)
		return err
	}
	e.Logger.Infow("deployed router", "addr", routerContract.Address)

	testRouterContract, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*router.Router] {
			routerAddr, tx2, routerC, err2 := router.DeployRouter(
				chain.DeployerKey,
				chain.Client,
				contractConfig.Weth9.Address(),
				rmnProxy.Address,
			)
			return ContractDeploy[*router.Router]{
				routerAddr, routerC, tx2, deployment.NewTypeAndVersion(TestRouter, deployment.Version1_2_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy test router", "err", err)
		return err
	}
	e.Logger.Infow("deployed test router", "addr", testRouterContract.Address)

	tokenAdminRegistry, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*token_admin_registry.TokenAdminRegistry] {
			tokenAdminRegistryAddr, tx2, tokenAdminRegistry, err2 := token_admin_registry.DeployTokenAdminRegistry(
				chain.DeployerKey,
				chain.Client)
			return ContractDeploy[*token_admin_registry.TokenAdminRegistry]{
				tokenAdminRegistryAddr, tokenAdminRegistry, tx2, deployment.NewTypeAndVersion(TokenAdminRegistry, deployment.Version1_5_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy token admin registry", "err", err)
		return err
	}
	e.Logger.Infow("deployed tokenAdminRegistry", "addr", tokenAdminRegistry)

	nonceManager, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*nonce_manager.NonceManager] {
			nonceManagerAddr, tx2, nonceManager, err2 := nonce_manager.DeployNonceManager(
				chain.DeployerKey,
				chain.Client,
				[]common.Address{}, // Need to add onRamp after
			)
			return ContractDeploy[*nonce_manager.NonceManager]{
				nonceManagerAddr, nonceManager, tx2, deployment.NewTypeAndVersion(NonceManager, deployment.Version1_6_0_dev), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy nonce manager", "err", err)
		return err
	}
	e.Logger.Infow("Deployed nonce manager", "addr", nonceManager.Address)

	feeQuoter, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*fee_quoter.FeeQuoter] {
			prAddr, tx2, pr, err2 := fee_quoter.DeployFeeQuoter(
				chain.DeployerKey,
				chain.Client,
				fee_quoter.FeeQuoterStaticConfig{
					MaxFeeJuelsPerMsg:  big.NewInt(0).Mul(big.NewInt(2e2), big.NewInt(1e18)),
					LinkToken:          contractConfig.LinkToken.Address(),
					StalenessThreshold: uint32(24 * 60 * 60),
				},
				[]common.Address{mcmsContracts.Timelock.Address},                                     // timelock should be able to update, ramps added after
				[]common.Address{contractConfig.Weth9.Address(), contractConfig.LinkToken.Address()}, // fee tokens
				[]fee_quoter.FeeQuoterTokenPriceFeedUpdate{},
				[]fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{}, // TODO: tokens
				[]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs{
					{
						PremiumMultiplierWeiPerEth: 9e17, // 0.9 ETH
						Token:                      contractConfig.LinkToken.Address(),
					},
					{
						PremiumMultiplierWeiPerEth: 1e18,
						Token:                      contractConfig.Weth9.Address(),
					},
				},
				[]fee_quoter.FeeQuoterDestChainConfigArgs{},
			)
			return ContractDeploy[*fee_quoter.FeeQuoter]{
				prAddr, pr, tx2, deployment.NewTypeAndVersion(FeeQuoter, deployment.Version1_6_0_dev), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy fee quoter", "err", err)
		return err
	}
	e.Logger.Infow("Deployed fee quoter", "addr", feeQuoter.Address)

	onRamp, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*onramp.OnRamp] {
			onRampAddr, tx2, onRamp, err2 := onramp.DeployOnRamp(
				chain.DeployerKey,
				chain.Client,
				onramp.OnRampStaticConfig{
					ChainSelector:      chain.Selector,
					RmnRemote:          rmnProxy.Address,
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
				onRampAddr, onRamp, tx2, deployment.NewTypeAndVersion(OnRamp, deployment.Version1_6_0_dev), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy onramp", "err", err)
		return err
	}
	e.Logger.Infow("Deployed onramp", "addr", onRamp.Address)

	offRamp, err := deployContract(e.Logger, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*offramp.OffRamp] {
			offRampAddr, tx2, offRamp, err2 := offramp.DeployOffRamp(
				chain.DeployerKey,
				chain.Client,
				offramp.OffRampStaticConfig{
					ChainSelector:      chain.Selector,
					RmnRemote:          rmnProxy.Address,
					NonceManager:       nonceManager.Address,
					TokenAdminRegistry: tokenAdminRegistry.Address,
				},
				offramp.OffRampDynamicConfig{
					FeeQuoter:                               feeQuoter.Address,
					PermissionLessExecutionThresholdSeconds: uint32(86400),
				},
				[]offramp.OffRampSourceChainConfigArgs{},
			)
			return ContractDeploy[*offramp.OffRamp]{
				offRampAddr, offRamp, tx2, deployment.NewTypeAndVersion(OffRamp, deployment.Version1_6_0_dev), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy offramp", "err", err)
		return err
	}
	e.Logger.Infow("Deployed offramp", "addr", offRamp.Address)

	// Basic wiring is always needed.
	tx, err = feeQuoter.Contract.ApplyAuthorizedCallerUpdates(chain.DeployerKey, fee_quoter.AuthorizedCallersAuthorizedCallerArgs{
		// TODO: We enable the deployer initially to set prices
		// Should be removed after.
		AddedCallers: []common.Address{offRamp.Contract.Address(), chain.DeployerKey.From},
	})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		e.Logger.Errorw("Failed to confirm fee quoter authorized caller update", "err", err)
		return err
	}

	tx, err = nonceManager.Contract.ApplyAuthorizedCallerUpdates(chain.DeployerKey, nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
		AddedCallers: []common.Address{offRamp.Contract.Address(), onRamp.Contract.Address()},
	})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		e.Logger.Errorw("Failed to update nonce manager with ramps", "err", err)
		return err
	}
	return nil
}
