package v1_6

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/sync/errgroup"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"
)

var _ cldf.ChangeSet[ccipseq.DeployChainContractsConfig] = DeployChainContractsChangeset

// DeployChainContracts deploys all new CCIP v1.6 or later contracts for the given chains.
// It returns the new addresses for the contracts.
// DeployChainContractsChangeset is idempotent. If there is an error, it will return the successfully deployed addresses and the error so that the caller can call the
// changeset again with the same input to retry the failed deployment.
// Caller should update the environment's address book with the returned addresses.
// Points to note :
// In case of migrating from legacy ccip to 1.6, the previous RMN address should be set while deploying RMNRemote.
// if there is no existing RMN address found, RMNRemote will be deployed with 0x0 address for previous RMN address
// which will set RMN to 0x0 address immutably in RMNRemote.
func DeployChainContractsChangeset(env cldf.Environment, c ccipseq.DeployChainContractsConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid DeployChainContractsConfig: %w", err)
	}
	newAddresses := cldf.NewMemoryAddressBook()
	err := deployChainContractsForChains(env, newAddresses, c.HomeChainSelector, c.ContractParamsPerChain)
	if err != nil {
		env.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", newAddresses)
		return cldf.ChangesetOutput{AddressBook: newAddresses}, cldf.MaybeDataErr(err)
	}
	return cldf.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}

type FeeQuoterParamsOld struct {
	MaxFeeJuelsPerMsg              *big.Int
	TokenPriceStalenessThreshold   uint32
	LinkPremiumMultiplierWeiPerEth uint64
	WethPremiumMultiplierWeiPerEth uint64
}

func ValidateHomeChainState(e cldf.Environment, homeChainSel uint64, existingState stateview.CCIPOnChainState) error {
	existingState, err := stateview.LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return err
	}
	capReg := existingState.Chains[homeChainSel].CapabilityRegistry
	if capReg == nil {
		e.Logger.Errorw("Failed to get capability registry")
		return errors.New("capability registry not found")
	}
	cr, err := capReg.GetHashedCapabilityId(
		&bind.CallOpts{}, shared.CapabilityLabelledName, shared.CapabilityVersion)
	if err != nil {
		e.Logger.Errorw("Failed to get hashed capability id", "err", err)
		return err
	}
	if cr != shared.CCIPCapabilityID {
		return fmt.Errorf("unexpected mismatch between calculated ccip capability id (%s) and expected ccip capability id constant (%s)",
			hexutil.Encode(cr[:]),
			hexutil.Encode(shared.CCIPCapabilityID[:]))
	}
	capability, err := capReg.GetCapability(nil, shared.CCIPCapabilityID)
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
		return errors.New("ccip home address mismatch")
	}
	rmnHome := existingState.Chains[homeChainSel].RMNHome
	if rmnHome == nil {
		e.Logger.Errorw("Failed to get rmn home", "err", err)
		return errors.New("rmn home not found")
	}
	return nil
}

func deployChainContractsForChains(
	e cldf.Environment,
	ab cldf.AddressBook,
	homeChainSel uint64,
	contractParamsPerChain map[uint64]ccipseq.ChainContractParams) error {
	existingState, err := stateview.LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return err
	}

	err = ValidateHomeChainState(e, homeChainSel, existingState)
	if err != nil {
		return err
	}

	rmnHome := existingState.Chains[homeChainSel].RMNHome

	deployGrp := errgroup.Group{}

	for chainSel, contractParams := range contractParamsPerChain {
		if _, exists := existingState.SupportedChains()[chainSel]; !exists {
			return fmt.Errorf("chain %d not supported", chainSel)
		}
		// already validated family
		family, _ := chainsel.GetSelectorFamily(chainSel)
		var deployFn func() error
		switch family {
		case chainsel.FamilyEVM:
			staticLinkExists := existingState.Chains[chainSel].StaticLinkToken != nil
			linkExists := existingState.Chains[chainSel].LinkToken != nil
			weth9Exists := existingState.Chains[chainSel].Weth9 != nil
			feeTokensAreValid := weth9Exists && (linkExists != staticLinkExists)
			if !feeTokensAreValid {
				return fmt.Errorf("fee tokens not valid for chain %d, staticLinkExists: %t, linkExists: %t, weth9Exists: %t", chainSel, staticLinkExists, linkExists, weth9Exists)
			}
			chain := e.Chains[chainSel]
			deployFn = func() error { return deployChainContractsEVM(e, chain, ab, rmnHome, contractParams) }
		default:
			return fmt.Errorf("unsupported chain family for chain %d", chainSel)
		}
		deployGrp.Go(func() error {
			err := deployFn()
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

func deployChainContractsEVM(e cldf.Environment, chain cldf.Chain, ab cldf.AddressBook, rmnHome *rmn_home.RMNHome, contractParams ccipseq.ChainContractParams) error {
	// check for existing contracts
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return err
	}
	chainState, chainExists := state.Chains[chain.Selector]
	if !chainExists {
		return fmt.Errorf("chain %s not found in existing state, deploy the prerequisites first", chain.String())
	}

	activeDigest, err := rmnHome.GetActiveDigest(&bind.CallOpts{})
	if err != nil {
		e.Logger.Errorw("Failed to get active digest", "chain", chain.String(), "err", err)
		return err
	}

	// get existing config
	existingConfig, err := rmnRemoteContract.GetVersionedConfig(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		e.Logger.Errorw("Failed to get existing config from rmnRemote",
			"chain", chain.String(),
			"rmnRemote", rmnRemoteContract.Address(),
			"err", err)
		return err
	}
	// is the config already set?
	// if the config is already set, the version should be more than 0, and we can check if it matches the active digest on RMNHome
	// In this case we don't need to set it again on existing RMNRemote
	if existingConfig.Version > 0 && existingConfig.Config.RmnHomeContractConfigDigest == activeDigest {
		e.Logger.Infow("rmn remote config already set", "chain", chain.String(),
			"RmnHomeContractConfigDigest", existingConfig.Config.RmnHomeContractConfigDigest,
			"Signers", existingConfig.Config.Signers,
			"FSign", existingConfig.Config.FSign,
		)
	} else {
		e.Logger.Infow("setting active home digest to rmn remote", "chain", chain.String(), "digest", activeDigest)
		tx, err := rmnRemoteContract.SetConfig(chain.DeployerKey, rmn_remote.RMNRemoteConfig{
			RmnHomeContractConfigDigest: activeDigest,
			Signers: []rmn_remote.RMNRemoteSigner{
				{NodeIndex: 0, OnchainPublicKey: common.Address{1}},
			},
			FSign: 0,
		})
		if _, err := cldf.ConfirmIfNoErrorWithABI(chain, tx, rmn_remote.RMNRemoteABI, err); err != nil {
			e.Logger.Errorw("Failed to confirm RMNRemote config", "chain", chain.String(), "err", err)
			return err
		}
		e.Logger.Infow("rmn remote config set", "chain", chain.String())
	}

	nmContract := chainState.NonceManager
	if chainState.NonceManager == nil {
		nonceManager, err := cldf.DeployContract(e.Logger, chain, ab,
			func(chain cldf.Chain) cldf.ContractDeploy[*nonce_manager.NonceManager] {
				nonceManagerAddr, tx2, nonceManager, err2 := nonce_manager.DeployNonceManager(
					chain.DeployerKey,
					chain.Client,
					[]common.Address{}, // Need to add onRamp after
				)
				return cldf.ContractDeploy[*nonce_manager.NonceManager]{
					Address: nonceManagerAddr, Contract: nonceManager, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.NonceManager, deployment.Version1_6_0), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy nonce manager", "chain", chain.String(), "err", err)
			return err
		}
		nmContract = nonceManager.Contract
	} else {
		e.Logger.Infow("nonce manager already deployed", "chain", chain.String(), "addr", chainState.NonceManager.Address)
	}
	feeQuoterContract := chainState.FeeQuoter
	if chainState.FeeQuoter == nil {
		feeQuoter, err := cldf.DeployContract(e.Logger, chain, ab,
			func(chain cldf.Chain) cldf.ContractDeploy[*fee_quoter.FeeQuoter] {
				prAddr, tx2, pr, err2 := fee_quoter.DeployFeeQuoter(
					chain.DeployerKey,
					chain.Client,
					fee_quoter.FeeQuoterStaticConfig{
						MaxFeeJuelsPerMsg:            contractParams.FeeQuoterParams.MaxFeeJuelsPerMsg,
						LinkToken:                    linkTokenContractAddr,
						TokenPriceStalenessThreshold: contractParams.FeeQuoterParams.TokenPriceStalenessThreshold,
					},
					[]common.Address{state.Chains[chain.Selector].Timelock.Address()}, // timelock should be able to update, ramps added after
					[]common.Address{weth9Contract.Address(), linkTokenContractAddr},  // fee tokens
					contractParams.FeeQuoterParams.TokenPriceFeedUpdates,
					contractParams.FeeQuoterParams.TokenTransferFeeConfigArgs,
					append([]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs{
						{
							PremiumMultiplierWeiPerEth: contractParams.FeeQuoterParams.LinkPremiumMultiplierWeiPerEth,
							Token:                      linkTokenContractAddr,
						},
						{
							PremiumMultiplierWeiPerEth: contractParams.FeeQuoterParams.WethPremiumMultiplierWeiPerEth,
							Token:                      weth9Contract.Address(),
						},
					}, contractParams.FeeQuoterParams.MorePremiumMultiplierWeiPerEth...),
					contractParams.FeeQuoterParams.DestChainConfigArgs,
				)
				return cldf.ContractDeploy[*fee_quoter.FeeQuoter]{
					Address: prAddr, Contract: pr, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.FeeQuoter, deployment.Version1_6_0), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy fee quoter", "chain", chain.String(), "err", err)
			return err
		}
		feeQuoterContract = feeQuoter.Contract
	} else {
		e.Logger.Infow("fee quoter already deployed", "chain", chain.String(), "addr", chainState.FeeQuoter.Address)
	}
	onRampContract := chainState.OnRamp
	if onRampContract == nil {
		onRamp, err := cldf.DeployContract(e.Logger, chain, ab,
			func(chain cldf.Chain) cldf.ContractDeploy[*onramp.OnRamp] {
				onRampAddr, tx2, onRamp, err2 := onramp.DeployOnRamp(
					chain.DeployerKey,
					chain.Client,
					onramp.OnRampStaticConfig{
						ChainSelector:      chain.Selector,
						RmnRemote:          RMNProxy.Address(),
						NonceManager:       nmContract.Address(),
						TokenAdminRegistry: tokenAdminReg.Address(),
					},
					onramp.OnRampDynamicConfig{
						FeeQuoter:     feeQuoterContract.Address(),
						FeeAggregator: chain.DeployerKey.From, // TODO real fee aggregator, using deployer key for now
					},
					[]onramp.OnRampDestChainConfigArgs{},
				)
				return cldf.ContractDeploy[*onramp.OnRamp]{
					Address: onRampAddr, Contract: onRamp, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.OnRamp, deployment.Version1_6_0), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy onramp", "chain", chain.String(), "err", err)
			return err
		}
		onRampContract = onRamp.Contract
	} else {
		e.Logger.Infow("onramp already deployed", "chain", chain.String(), "addr", chainState.OnRamp.Address)
	}
	offRampContract := chainState.OffRamp
	if offRampContract == nil {
		offRamp, err := cldf.DeployContract(e.Logger, chain, ab,
			func(chain cldf.Chain) cldf.ContractDeploy[*offramp.OffRamp] {
				offRampAddr, tx2, offRamp, err2 := offramp.DeployOffRamp(
					chain.DeployerKey,
					chain.Client,
					offramp.OffRampStaticConfig{
						ChainSelector:        chain.Selector,
						GasForCallExactCheck: contractParams.OffRampParams.GasForCallExactCheck,
						RmnRemote:            RMNProxy.Address(),
						NonceManager:         nmContract.Address(),
						TokenAdminRegistry:   tokenAdminReg.Address(),
					},
					offramp.OffRampDynamicConfig{
						FeeQuoter:                               feeQuoterContract.Address(),
						PermissionLessExecutionThresholdSeconds: contractParams.OffRampParams.PermissionLessExecutionThresholdSeconds,
						MessageInterceptor:                      contractParams.OffRampParams.MessageInterceptor,
					},
					[]offramp.OffRampSourceChainConfigArgs{},
				)
				return cldf.ContractDeploy[*offramp.OffRamp]{
					Address: offRampAddr, Contract: offRamp, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.OffRamp, deployment.Version1_6_0), Err: err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy offramp", "chain", chain.String(), "err", err)
			return err
		}
		offRampContract = offRamp.Contract
	} else {
		e.Logger.Infow("offramp already deployed", "chain", chain.String(), "addr", chainState.OffRamp.Address)
	}
	// Basic wiring is always needed.
	// check if there is already a wiring
	callers, err := feeQuoterContract.GetAllAuthorizedCallers(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		e.Logger.Errorw("Failed to get fee quoter authorized callers",
			"chain", chain.String(),
			"feeQuoter", feeQuoterContract.Address(),
			"err", err)
		return err
	}
	// should only update callers if there are none, otherwise we might overwrite some existing callers for existing fee quoter
	if len(callers) == 1 && (callers[0] == chain.DeployerKey.From || callers[0] == state.Chains[chain.Selector].Timelock.Address()) {
		tx, err := feeQuoterContract.ApplyAuthorizedCallerUpdates(chain.DeployerKey, fee_quoter.AuthorizedCallersAuthorizedCallerArgs{
			// TODO: We enable the deployer initially to set prices
			// Should be removed after.
			AddedCallers: []common.Address{offRampContract.Address(), chain.DeployerKey.From},
		})
		if _, err := cldf.ConfirmIfNoErrorWithABI(chain, tx, fee_quoter.FeeQuoterABI, err); err != nil {
			e.Logger.Errorw("Failed to confirm fee quoter authorized caller update", "chain", chain.String(), "err", err)
			return err
		}
		e.Logger.Infow("Added fee quoter authorized callers", "chain", chain.String(), "callers", []common.Address{offRampContract.Address(), chain.DeployerKey.From})
	}
	// get all authorized callers
	// check if there is already a wiring
	nmCallers, err := nmContract.GetAllAuthorizedCallers(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		e.Logger.Errorw("Failed to get nonce manager authorized callers",
			"chain", chain.String(),
			"nonceManager", nmContract.Address(),
			"err", err)
		return err
	}
	// should only update callers if there are none, otherwise we might overwrite some existing callers for existing nonce manager
	if len(nmCallers) == 0 {
		tx, err := nmContract.ApplyAuthorizedCallerUpdates(chain.DeployerKey, nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
			AddedCallers: []common.Address{offRampContract.Address(), onRampContract.Address()},
		})
		if _, err := cldf.ConfirmIfNoErrorWithABI(chain, tx, nonce_manager.NonceManagerABI, err); err != nil {
			e.Logger.Errorw("Failed to update nonce manager with ramps", "chain", chain.String(), "err", err)
			return err
		}
		e.Logger.Infow("Added nonce manager authorized callers", "chain", chain.String(), "callers", []common.Address{offRampContract.Address(), onRampContract.Address()})
	}
	return nil
}
