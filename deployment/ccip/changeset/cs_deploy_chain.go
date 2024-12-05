package changeset

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

var _ deployment.ChangeSet[DeployChainContractsConfig] = DeployChainContracts

// DeployChainContracts deploys all new CCIP v1.6 or later contracts for the given chains.
// It returns the new addresses for the contracts.
// DeployChainContracts is idempotent. If there is an error, it will return the successfully deployed addresses and the error so that the caller can call the
// changeset again with the same input to retry the failed deployment.
// Caller should update the environment's address book with the returned addresses.
func DeployChainContracts(env deployment.Environment, c DeployChainContractsConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid DeployChainContractsConfig: %w", err)
	}
	newAddresses := deployment.NewMemoryAddressBook()
	err := deployChainContractsForChains(env, newAddresses, c.HomeChainSelector, c.ChainSelectors)
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

type DeployChainContractsConfig struct {
	ChainSelectors    []uint64
	HomeChainSelector uint64
}

func (c DeployChainContractsConfig) Validate() error {
	for _, cs := range c.ChainSelectors {
		if err := deployment.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
	}
	if err := deployment.IsValidChainSelector(c.HomeChainSelector); err != nil {
		return fmt.Errorf("invalid home chain selector: %d - %w", c.HomeChainSelector, err)
	}
	return nil
}

// deployCCIPContracts assumes the following contracts are deployed:
// - Capability registry
// - CCIP home
// - RMN home
// - Fee tokens on all chains.
// and present in ExistingAddressBook.
// It then deploys the rest of the CCIP chain contracts to the selected chains
// registers the nodes with the capability registry and creates a DON for
// each new chain.
func deployCCIPContracts(
	e deployment.Environment,
	ab deployment.AddressBook,
	c NewChainsConfig) error {
	err := deployChainContractsForChains(e, ab, c.HomeChainSel, c.Chains())
	if err != nil {
		e.Logger.Errorw("Failed to deploy chain contracts", "err", err)
		return err
	}
	err = e.ExistingAddresses.Merge(ab)
	if err != nil {
		e.Logger.Errorw("Failed to merge address book", "err", err)
		return err
	}
	err = configureChain(e, c)
	if err != nil {
		e.Logger.Errorw("Failed to add chain", "err", err)
		return err
	}

	return nil
}

func deployChainContractsForChains(
	e deployment.Environment,
	ab deployment.AddressBook,
	homeChainSel uint64,
	chainsToDeploy []uint64) error {
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
		&bind.CallOpts{}, internal.CapabilityLabelledName, internal.CapabilityVersion)
	if err != nil {
		e.Logger.Errorw("Failed to get hashed capability id", "err", err)
		return err
	}
	if cr != internal.CCIPCapabilityID {
		return fmt.Errorf("unexpected mismatch between calculated ccip capability id (%s) and expected ccip capability id constant (%s)",
			hexutil.Encode(cr[:]),
			hexutil.Encode(internal.CCIPCapabilityID[:]))
	}
	capability, err := capReg.GetCapability(nil, internal.CCIPCapabilityID)
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
	deployGrp := errgroup.Group{}
	for _, chainSel := range chainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain %d not found", chainSel)
		}
		if existingState.Chains[chainSel].LinkToken == nil || existingState.Chains[chainSel].Weth9 == nil {
			return fmt.Errorf("fee tokens not found for chain %d", chainSel)
		}
		deployGrp.Go(
			func() error {
				err := deployChainContracts(e, chain, ab, rmnHome)
				if err != nil {
					e.Logger.Errorw("Failed to deploy chain contracts", "chain", chainSel, "err", err)
					return fmt.Errorf("failed to deploy chain contracts for chain %d: %w", chainSel, err)
				}
				return nil
			})
	}
	if err := deployGrp.Wait(); err != nil {
		e.Logger.Errorw("Failed to deploy chain contracts", "err", err)
		return err
	}
	return nil
}

func deployChainContracts(
	e deployment.Environment,
	chain deployment.Chain,
	ab deployment.AddressBook,
	rmnHome *rmn_home.RMNHome,
) error {
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
	if chainState.Timelock == nil {
		return fmt.Errorf("timelock not found for chain %d, deploy the mcms contracts first", chain.Selector)
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
					// Indicates no legacy RMN contract
					common.HexToAddress("0x0"),
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
					[]common.Address{state.Chains[chain.Selector].Timelock.Address()},      // timelock should be able to update, ramps added after
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
