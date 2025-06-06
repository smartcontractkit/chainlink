package v1_6

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	ccipopsv1_2 "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_2"
	ccipopsv1_6 "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

type DeployChainContractsConfig struct {
	HomeChainSelector      uint64
	ContractParamsPerChain map[uint64]ChainContractParams
}

func (c DeployChainContractsConfig) Validate() error {
	if err := deployment.IsValidChainSelector(c.HomeChainSelector); err != nil {
		return fmt.Errorf("invalid home chain selector: %d - %w", c.HomeChainSelector, err)
	}
	for cs, args := range c.ContractParamsPerChain {
		if err := deployment.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
		if err := args.Validate(); err != nil {
			return fmt.Errorf("invalid contract args for chain %d: %w", cs, err)
		}
	}
	return nil
}

type ChainContractParams struct {
	FeeQuoterParams ccipopsv1_6.FeeQuoterParams
	OffRampParams   ccipopsv1_6.OffRampParams
}

func (c ChainContractParams) Validate() error {
	if err := c.FeeQuoterParams.Validate(); err != nil {
		return fmt.Errorf("invalid FeeQuoterParams: %w", err)
	}
	if err := c.OffRampParams.Validate(false); err != nil {
		return fmt.Errorf("invalid OffRampParams: %w", err)
	}
	return nil
}

var (
	DeployChainContractsSeq = operations.NewSequence(
		"DeployChainContractsSeq",
		semver.MustParse("1.0.0"),
		"Deploys all 1.6 chain contracts for the specified evm chain(s)",
		func(b operations.Bundle, deps opsutil.ConfigureDependencies, input DeployChainContractsConfig) (any, error) {
			state := deps.CurrentState
			e := deps.Env
			grp := errgroup.Group{}
			rmnHome := state.MustGetEVMChainState(input.HomeChainSelector).RMNHome
			if rmnHome == nil {
				return nil, errors.New("rmn home not found, deploy home chain contracts first")
			}
			for chainSelector, contractParams := range input.ContractParamsPerChain {
				chainSelector := chainSelector
				contractParams := contractParams
				grp.Go(func() error {
					chain, ok := e.BlockChains.EVMChains()[chainSelector]
					if !ok {
						return fmt.Errorf("chain %d not found in env", chainSelector)
					}
					chainState, chainExists := state.EVMChainState(chainSelector)
					if !chainExists {
						return fmt.Errorf("chain %s not found in existing state, deploy the prerequisites first", chain.String())
					}
					deployDeps := opsutil.DeployContractDependencies{
						Chain:       chain,
						AddressBook: deps.AddressBook,
					}
					staticLinkExists := state.MustGetEVMChainState(chainSelector).StaticLinkToken != nil
					linkExists := state.MustGetEVMChainState(chainSelector).LinkToken != nil
					weth9Exists := state.MustGetEVMChainState(chainSelector).Weth9 != nil
					feeTokensAreValid := weth9Exists && (linkExists != staticLinkExists)
					if !feeTokensAreValid {
						return fmt.Errorf("fee tokens not valid for chain %d, staticLinkExists: %t, linkExists: %t, weth9Exists: %t", chainSelector, staticLinkExists, linkExists, weth9Exists)
					}

					if chainState.Weth9 == nil {
						return fmt.Errorf("weth9 not found for chain %s, deploy the prerequisites first", chain.String())
					}
					if chainState.Timelock == nil {
						return fmt.Errorf("timelock not found for chain %s, deploy the mcms contracts first", chain.String())
					}
					linkAddr, err := chainState.LinkTokenAddress()
					if err != nil {
						return fmt.Errorf("failed to get link token address for chain %s: %w", chain.String(), err)
					}
					if chainState.TokenAdminRegistry == nil {
						return fmt.Errorf("token admin registry not found for chain %s, deploy the prerequisites first", chain.String())
					}
					if len(chainState.RegistryModules1_6) == 0 && len(chainState.RegistryModules1_5) == 0 {
						return fmt.Errorf("registry module not found for chain %s, deploy the prerequisites first", chain.String())
					}
					if chainState.Router == nil {
						return fmt.Errorf("router not found for chain %s, deploy the prerequisites first", chain.String())
					}
					if chainState.RMNProxy == nil {
						b.Logger.Errorw("RMNProxy not found", "chain", chain.String())
						return fmt.Errorf("rmn proxy not found for chain %s, deploy the prerequisites first", chain.String())
					}

					// Deploy RMNRemote if not already deployed
					var rmnLegacyAddr common.Address
					if chainState.MockRMN != nil {
						rmnLegacyAddr = chainState.MockRMN.Address()
					}
					// If RMN is deployed, set rmnLegacyAddr to the RMN address
					if chainState.RMN != nil {
						rmnLegacyAddr = chainState.RMN.Address()
					}
					if rmnLegacyAddr == (common.Address{}) {
						b.Logger.Warnf("No legacy RMN contract found for chain %s, will not setRMN in RMNRemote", chain.String())
					}

					rmnRemote := chainState.RMNRemote
					if rmnRemote == nil {
						report, err := operations.ExecuteOperation(b, ccipopsv1_6.DeployRMNRemoteOp, deployDeps, ccipopsv1_6.DeployRMNRemoteInput{
							ChainSelector: chainSelector,
							RMNLegacyAddr: rmnLegacyAddr,
						})
						if err != nil {
							return fmt.Errorf("failed to deploy RMNRemote for chain %d: %w", chainSelector, err)
						}
						if report.Output == (common.Address{}) {
							return fmt.Errorf("failed to deploy RMNRemote for chain %d: empty address", chainSelector)
						}
						rmnRemote, err = rmn_remote.NewRMNRemote(
							report.Output,
							chain.Client,
						)
						if err != nil {
							return fmt.Errorf("failed to create RMNRemote contract for chain after deployment %d: %w", chainSelector, err)
						}
						chainState.RMNRemote = rmnRemote
					} else {
						b.Logger.Infow("rmn remote already deployed", "chain", chain.String(), "addr", chainState.RMNRemote.Address())
					}
					// set RMNRemote config if not already set
					set, err := isRMNRemoteInitialSetUpCompleted(e, rmnHome, rmnRemote, chain)
					if err != nil {
						return fmt.Errorf("failed to check if RMNRemote config is set for chain %d: %w", chainSelector, err)
					}
					// if no config is set, we need to set it with active digest and initial empty signers
					if !set {
						// RMNRemote needs to be set in state
						deps.CurrentState.WriteEVMChainState(chainSelector, chainState)
						// set RMNRemote config
						_, err = operations.ExecuteOperation(b, ccipopsv1_6.SetRMNRemoteConfigOp, deps, ccipopsv1_6.SetRMNRemoteConfig{
							RMNRemoteConfig: ccipopsv1_6.RMNRemoteConfig{
								Signers: []rmn_remote.RMNRemoteSigner{
									{NodeIndex: 0, OnchainPublicKey: common.Address{1}},
								},
								F: 0,
							},
							ChainSelector: chainSelector,
						})
						if err != nil {
							return fmt.Errorf("failed to set RMNRemote config for chain %d: %w", chainSelector, err)
						}
					}
					// Deploy Test Router if not already deployed
					if chainState.TestRouter != nil {
						b.Logger.Infow("test router already deployed", "chain", chain.String(), "addr", chainState.TestRouter.Address)
					} else {
						_, err = operations.ExecuteOperation(b, ccipopsv1_2.DeployRouter, deployDeps, ccipopsv1_2.DeployRouterInput{
							IsTestRouter: true,
							RMNProxy:     chainState.RMNProxy.Address(),
							WethAddress:  chainState.Weth9.Address(),
						})
						if err != nil {
							return fmt.Errorf("failed to deploy test router for chain %d: %w", chainSelector, err)
						}
					}
					// Deploy NonceManager if not already deployed
					if chainState.NonceManager == nil {
						nmReport, err := operations.ExecuteOperation(b, ccipopsv1_6.DeployNonceManagerOp, deployDeps, chainSelector)
						if err != nil {
							return fmt.Errorf("failed to deploy nonce manager for chain %d: %w", chainSelector, err)
						}
						if err != nil {
							return fmt.Errorf("failed to deploy nonce manager for chain %d: %w", chainSelector, err)
						}
						if nmReport.Output == (common.Address{}) {
							return fmt.Errorf("failed to deploy nonce manager for chain %d: empty address", chainSelector)
						}
						chainState.NonceManager, err = nonce_manager.NewNonceManager(
							nmReport.Output,
							chain.Client,
						)
						if err != nil {
							return fmt.Errorf("failed to create nonce manager contract for chain after deployment %d: %w", chainSelector, err)
						}
					} else {
						b.Logger.Infow("nonce manager already deployed", "chain", chain.String(), "addr", chainState.NonceManager.Address())
					}
					// Deploy FeeQuoter if not already deployed
					if chainState.FeeQuoter == nil {
						feeQReport, err := operations.ExecuteOperation(b, ccipopsv1_6.DeployFeeQuoterOp, deployDeps, ccipopsv1_6.DeployFeeQInput{
							Chain:         chainSelector,
							Params:        contractParams.FeeQuoterParams,
							LinkAddr:      linkAddr,
							WethAddr:      chainState.Weth9.Address(),
							PriceUpdaters: []common.Address{chainState.Timelock.Address()}, // timelock should be able to update, ramps added after
						})
						if err != nil {
							return fmt.Errorf("failed to deploy fee quoter for chain %d: %w", chainSelector, err)
						}
						if feeQReport.Output == (common.Address{}) {
							return fmt.Errorf("failed to deploy fee quoter for chain %d: empty address", chainSelector)
						}
						chainState.FeeQuoter, err = fee_quoter.NewFeeQuoter(
							feeQReport.Output,
							chain.Client,
						)
						if err != nil {
							return fmt.Errorf("failed to create fee quoter contract for chain after deployment %d: %w", chainSelector, err)
						}
					} else {
						b.Logger.Infow("fee quoter already deployed", "chain", chain.String(), "addr", chainState.FeeQuoter.Address())
					}
					if chainState.FeeQuoter == nil {
						b.Logger.Errorw("FeeQuoter not found", "chain", chain.String())
						return fmt.Errorf("fee quoter not found for chain %s, needed for onRamp/OffRamp deployment", chain.String())
					}
					if chainState.NonceManager == nil {
						b.Logger.Errorw("NonceManager not found", "chain", chain.String())
						return fmt.Errorf("nonce manager not found for chain %s, needed for onRamp/OffRamp deployment", chain.String())
					}
					// Deploy OnRamp if not already deployed
					if chainState.OnRamp == nil {
						// if the fee aggregator is not set, use the deployer key address
						// this is to ensure that feeAggregator is not set to zero address, otherwise there is chance of
						// fund loss when WithdrawFeeToken is called on OnRamp
						feeAggregator := chainState.FeeAggregator
						if feeAggregator == (common.Address{}) {
							feeAggregator = chain.DeployerKey.From
						}
						onRampReport, err := operations.ExecuteOperation(b, ccipopsv1_6.DeployOnRampOp, deployDeps, ccipopsv1_6.DeployOnRampInput{
							ChainSelector:      chainSelector,
							TokenAdminRegistry: chainState.TokenAdminRegistry.Address(),
							NonceManager:       chainState.NonceManager.Address(),
							RmnRemote:          chainState.RMNProxy.Address(),
							FeeQuoter:          chainState.FeeQuoter.Address(),
							FeeAggregator:      feeAggregator,
						})
						if err != nil {
							return fmt.Errorf("failed to deploy onRamp for chain %d: %w", chainSelector, err)
						}
						if onRampReport.Output == (common.Address{}) {
							return fmt.Errorf("failed to deploy onRamp for chain %d: empty address", chainSelector)
						}
						chainState.OnRamp, err = onramp.NewOnRamp(
							onRampReport.Output,
							chain.Client,
						)
						if err != nil {
							return fmt.Errorf("failed to create onRamp contract for chain after deployment %d: %w", chainSelector, err)
						}
					} else {
						b.Logger.Infow("onRamp already deployed", "chain", chain.String(), "addr", chainState.OnRamp.Address())
					}
					// Deploy OffRamp if not already deployed
					if chainState.OffRamp == nil {
						offRampReport, err := operations.ExecuteOperation(b, ccipopsv1_6.DeployOffRampOp, deployDeps, ccipopsv1_6.DeployOffRampInput{
							Chain:              chainSelector,
							Params:             contractParams.OffRampParams,
							RmnRemote:          chainState.RMNProxy.Address(),
							NonceManager:       chainState.NonceManager.Address(),
							TokenAdminRegistry: chainState.TokenAdminRegistry.Address(),
							FeeQuoter:          chainState.FeeQuoter.Address(),
						})
						if err != nil {
							return fmt.Errorf("failed to deploy offramp for chain %d: %w", chainSelector, err)
						}

						if offRampReport.Output == (common.Address{}) {
							return fmt.Errorf("failed to deploy offramp for chain %d: empty address", chainSelector)
						}
						chainState.OffRamp, err = offramp.NewOffRamp(
							offRampReport.Output,
							chain.Client,
						)
						if err != nil {
							return fmt.Errorf("failed to create offramp contract for chain after deployment %d: %w", chainSelector, err)
						}
					} else {
						b.Logger.Infow("offramp already deployed", "chain", chain.String(), "addr", chainState.OffRamp.Address())
					}
					deps.CurrentState.WriteEVMChainState(chainSelector, chainState)
					callers, err := chainState.FeeQuoter.GetAllAuthorizedCallers(&bind.CallOpts{
						Context: e.GetContext(),
					})
					if err != nil {
						b.Logger.Errorw("Failed to get fee quoter authorized callers",
							"chain", chain.String(),
							"feeQuoter", chainState.FeeQuoter.Address(),
							"err", err)
						return err
					}
					// should only update callers if there are none, otherwise we might overwrite some existing callers for existing fee quoter
					if len(callers) == 1 && (callers[0] == chain.DeployerKey.From || callers[0] == state.MustGetEVMChainState(chain.Selector).Timelock.Address()) {
						_, err = operations.ExecuteOperation(b, ccipopsv1_6.FeeQApplyAuthorizedCallerOp, deps, ccipopsv1_6.FeeQApplyAuthorizedCallerOpInput{
							ChainSelector: chainSelector,
							Callers: fee_quoter.AuthorizedCallersAuthorizedCallerArgs{
								// TODO: We enable the deployer initially to set prices
								// Should be removed after.
								AddedCallers: []common.Address{chainState.OffRamp.Address(), chain.DeployerKey.From},
							},
						})
						if err != nil {
							b.Logger.Errorw("Failed to apply fee quoter authorized callers",
								"chain", chain.String(),
								"feeQuoter", chainState.FeeQuoter.Address(),
								"err", err)
							return err
						}
						b.Logger.Infow("Added fee quoter authorized callers", "chain", chain.String(), "callers",
							[]common.Address{chainState.OffRamp.Address(), chain.DeployerKey.From})
					}
					nmCallers, err := chainState.NonceManager.GetAllAuthorizedCallers(&bind.CallOpts{
						Context: e.GetContext(),
					})
					if err != nil {
						b.Logger.Errorw("Failed to get nonce manager authorized callers",
							"chain", chain.String(),
							"nonceManager", chainState.NonceManager.Address(),
							"err", err)
						return err
					}
					// should only update callers if there are none, otherwise we might overwrite some existing callers for existing nonce manager
					if len(nmCallers) == 0 {
						_, err = operations.ExecuteOperation(b, ccipopsv1_6.NonceManagerUpdateAuthorizedCallerOp, deps,
							ccipopsv1_6.NonceManagerUpdateAuthorizedCallerInput{
								ChainSelector: chainSelector,
								Callers: nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
									AddedCallers: []common.Address{
										chainState.OffRamp.Address(),
										chainState.OnRamp.Address(),
									},
								},
							})
						if err != nil {
							b.Logger.Errorw("Failed to apply nonce manager authorized callers",
								"chain", chain.String(),
								"nonceManager", chainState.NonceManager.Address(),
								"err", err)
							return err
						}
						b.Logger.Infow("Added nonce manager authorized callers",
							"chain", chain.String(), "callers", []common.Address{
								chainState.OffRamp.Address(),
								chainState.OnRamp.Address(),
							})
					}
					return nil
				})
			}
			if err := grp.Wait(); err != nil {
				return nil, fmt.Errorf("failed to deploy chain contracts: %w", err)
			}
			return nil, nil
		})
)

func isRMNRemoteInitialSetUpCompleted(e cldf.Environment, rmnHome *rmn_home.RMNHome, rmnRemoteContract *rmn_remote.RMNRemote, chain cldf_evm.Chain) (bool, error) {
	activeDigest, err := rmnHome.GetActiveDigest(&bind.CallOpts{})
	if err != nil {
		e.Logger.Errorw("Failed to get active digest", "chain", chain.String(), "err", err)
		return false, err
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
		return false, err
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
		return true, nil
	}
	return false, nil
}
