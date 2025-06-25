package v1_6

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Masterminds/semver/v3"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

	ccipopsv1_2 "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_2"
	ccipopsv1_6 "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

type ChainContractParams struct {
	FeeQuoterParams ccipopsv1_6.FeeQuoterParams
	OffRampParams   ccipopsv1_6.OffRampParams
}

func (c ChainContractParams) Validate(selector uint64) error {
	if err := c.FeeQuoterParams.Validate(); err != nil {
		return fmt.Errorf("invalid FeeQuoterParams: %w", err)
	}
	if err := c.OffRampParams.Validate(false); err != nil {
		return fmt.Errorf("invalid OffRampParams: %w", err)
	}

	return nil
}

type DeployChainContractsConfig struct {
	HomeChainSelector      uint64
	ContractParamsPerChain map[uint64]ChainContractParams
	GasBoostConfigPerChain map[uint64]commontypes.GasBoostConfig
}

func (c DeployChainContractsConfig) Validate() error {
	if err := deployment.IsValidChainSelector(c.HomeChainSelector); err != nil {
		return fmt.Errorf("invalid home chain selector: %d - %w", c.HomeChainSelector, err)
	}
	for cs, args := range c.ContractParamsPerChain {
		if err := deployment.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
		if err := args.Validate(cs); err != nil {
			return fmt.Errorf("invalid contract args for chain %d: %w", cs, err)
		}
	}
	return nil
}

type CCIPAddresses struct {
	LegacyRMNAddress          common.Address
	RMNProxyAddress           common.Address
	WrappedNativeAddress      common.Address
	TimelockAddress           common.Address
	LinkAddress               common.Address
	FeeAggregatorAddress      common.Address
	TokenAdminRegistryAddress common.Address

	// New addresses deployed in this sequence
	OnRampAddress       common.Address
	TestRouterAddress   common.Address
	OffRampAddress      common.Address
	NonceManagerAddress common.Address
	FeeQuoterAddress    common.Address
	RMNRemoteAddress    common.Address
}

func (c CCIPAddresses) Validate(selector uint64) error {
	if c.RMNProxyAddress == (common.Address{}) {
		return fmt.Errorf("rmn proxy address is not defined for chain with selector %d, deploy the prerequisites first", selector)
	}
	if c.WrappedNativeAddress == (common.Address{}) {
		return fmt.Errorf("wrapped native address is not defined for chain with selector %d, deploy the prerequisites first", selector)
	}
	if c.TimelockAddress == (common.Address{}) {
		return fmt.Errorf("timelock address is not defined for chain with selector %d, deploy the mcms contracts first", selector)
	}
	if c.LinkAddress == (common.Address{}) {
		return fmt.Errorf("link address is not defined for chain with selector %d", selector)
	}
	if c.TokenAdminRegistryAddress == (common.Address{}) {
		return fmt.Errorf("token admin registry address is not defined for chain with selector %d, deploy the prerequisites first", selector)
	}

	return nil
}

type DeployChainContractsSeqConfig struct {
	DeployChainContractsConfig
	RMNHomeAddress         common.Address
	AddressesPerChain      map[uint64]CCIPAddresses
	GasBoostConfigPerChain map[uint64]commontypes.GasBoostConfig
}

func (c DeployChainContractsSeqConfig) Validate() error {
	if err := c.DeployChainContractsConfig.Validate(); err != nil {
		return fmt.Errorf("invalid DeployChainContractsConfig: %w", err)
	}
	if c.RMNHomeAddress == (common.Address{}) {
		return errors.New("rmn home is not defined, deploy the home chain first")
	}
	for chainSelector, addresses := range c.AddressesPerChain {
		if _, ok := c.ContractParamsPerChain[chainSelector]; !ok {
			return fmt.Errorf("no contract params defined for chain %d, but addresses are provided", chainSelector)
		}
		if err := addresses.Validate(chainSelector); err != nil {
			return fmt.Errorf("invalid addresses for chain %d: %w", chainSelector, err)
		}
	}
	for chainSelector := range c.ContractParamsPerChain {
		if _, ok := c.AddressesPerChain[chainSelector]; !ok {
			return fmt.Errorf("no addresses defined for chain %d, but contract params are provided", chainSelector)
		}
	}

	return nil
}

var (
	DeployChainContractsSeq = operations.NewSequence(
		"DeployChainContractsSeq",
		semver.MustParse("1.0.0"),
		"Deploys all 1.6 chain contracts for the specified evm chain(s)",
		func(b operations.Bundle, deps map[uint64]cldf_evm.Chain, input DeployChainContractsSeqConfig) (map[uint64]map[string]string, error) {
			err := input.Validate()
			if err != nil {
				return nil, fmt.Errorf("invalid DeployChainContractsSeqConfig: %w", err)
			}
			gasBoostConfigs := opsutil.GasBoostConfigsForChainMap(input.ContractParamsPerChain, input.GasBoostConfigPerChain)
			out := make(map[uint64]map[string]string)
			grp := errgroup.Group{}
			homeChain, ok := deps[input.HomeChainSelector]
			if !ok {
				return nil, fmt.Errorf("home chain with selector %d not defined in dependencies", input.HomeChainSelector)
			}
			mu := sync.Mutex{}
			for chainSelector, contractParams := range input.ContractParamsPerChain {
				newAddresses := make(map[string]string)
				chainSelector := chainSelector
				contractParams := contractParams
				chainAddresses := input.AddressesPerChain[chainSelector]
				chain, ok := deps[chainSelector]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined in dependencies", chainSelector)
				}
				var rmnRemoteAddress common.Address
				var feeQuoterAddress common.Address
				var nonceManagerAddress common.Address
				var offRampAddress common.Address
				var onRampAddress common.Address
				var testRouterAddress common.Address
				grp.Go(func() error {
					if chainAddresses.RMNRemoteAddress == (common.Address{}) {
						if chainAddresses.LegacyRMNAddress == (common.Address{}) {
							b.Logger.Warnf("No legacy RMN contract found for %s, will not setRMN in RMNRemote", chain.String())
						}
						report, err := operations.ExecuteOperation(b, ccipopsv1_6.DeployRMNRemoteOp, chain, opsutil.EVMDeployInput[ccipopsv1_6.DeployRMNRemoteInput]{
							ChainSelector: chainSelector,
							DeployInput: ccipopsv1_6.DeployRMNRemoteInput{
								ChainSelector: chainSelector,
								RMNLegacyAddr: chainAddresses.LegacyRMNAddress,
							},
						}, opsutil.RetryDeploymentWithGasBoost[ccipopsv1_6.DeployRMNRemoteInput](gasBoostConfigs[chainSelector]))
						if err != nil {
							return fmt.Errorf("failed to deploy RMNRemote for %s: %w", chain, err)
						}
						rmnRemoteAddress = report.Output.Address
						newAddresses[rmnRemoteAddress.Hex()] = report.Output.TypeAndVersion
					} else {
						rmnRemoteAddress = chainAddresses.RMNRemoteAddress
					}
					// Set RMNRemote config if not already set
					// If no config is set, we need to set it with active digest and initial empty signers
					digest, set, err := isRMNRemoteInitialSetUpCompleted(b.GetContext(), input.RMNHomeAddress, homeChain, rmnRemoteAddress, chain)
					if err != nil {
						return fmt.Errorf("failed to check if RMNRemote config is set for %s: %w", chain, err)
					}
					if !set {
						_, err = operations.ExecuteOperation(b, ccipopsv1_6.SetRMNRemoteConfigOp, chain, opsutil.EVMCallInput[rmn_remote.RMNRemoteConfig]{
							Address:       rmnRemoteAddress,
							NoSend:        false,
							ChainSelector: chainSelector,
							CallInput: rmn_remote.RMNRemoteConfig{
								RmnHomeContractConfigDigest: digest,
								Signers: []rmn_remote.RMNRemoteSigner{
									{NodeIndex: 0, OnchainPublicKey: common.Address{1}},
								},
								FSign: 0,
							},
						}, opsutil.RetryCallWithGasBoost[rmn_remote.RMNRemoteConfig](gasBoostConfigs[chainSelector]))
						if err != nil {
							return fmt.Errorf("failed to set RMNRemote config for chain %d: %w", chainSelector, err)
						}
					} else {
						b.Logger.Infow("RMNRemote config already set", "chain", chain.String())
					}
					// Deploy Test Router if not already deployed
					if chainAddresses.TestRouterAddress == (common.Address{}) {
						report, err := operations.ExecuteOperation(b, ccipopsv1_2.DeployTestRouter, chain, opsutil.EVMDeployInput[ccipopsv1_2.DeployRouterInput]{
							ChainSelector: chainSelector,
							DeployInput: ccipopsv1_2.DeployRouterInput{
								ChainSelector: chainSelector,
								RMNProxy:      chainAddresses.RMNProxyAddress,
								WethAddress:   chainAddresses.WrappedNativeAddress,
							},
						}, opsutil.RetryDeploymentWithGasBoost[ccipopsv1_2.DeployRouterInput](gasBoostConfigs[chainSelector]))
						if err != nil {
							return fmt.Errorf("failed to deploy test router for %s: %w", chain, err)
						}
						testRouterAddress = report.Output.Address
						newAddresses[testRouterAddress.Hex()] = report.Output.TypeAndVersion
					}
					// Deploy NonceManager if not already deployed
					if chainAddresses.NonceManagerAddress == (common.Address{}) {
						report, err := operations.ExecuteOperation(b, ccipopsv1_6.DeployNonceManagerOp, chain, opsutil.EVMDeployInput[[]common.Address]{
							ChainSelector: chainSelector,
							DeployInput:   []common.Address{},
						}, opsutil.RetryDeploymentWithGasBoost[[]common.Address](gasBoostConfigs[chainSelector]))
						if err != nil {
							return fmt.Errorf("failed to deploy nonce manager for %s: %w", chain, err)
						}
						nonceManagerAddress = report.Output.Address
						newAddresses[nonceManagerAddress.Hex()] = report.Output.TypeAndVersion
					} else {
						nonceManagerAddress = chainAddresses.NonceManagerAddress
					}
					// Deploy FeeQuoter if not already deployed
					if chainAddresses.FeeQuoterAddress == (common.Address{}) {
						report, err := operations.ExecuteOperation(b, ccipopsv1_6.DeployFeeQuoterOp, chain, opsutil.EVMDeployInput[ccipopsv1_6.DeployFeeQInput]{
							ChainSelector: chainSelector,
							DeployInput: ccipopsv1_6.DeployFeeQInput{
								Chain:    chainSelector,
								Params:   contractParams.FeeQuoterParams,
								LinkAddr: chainAddresses.LinkAddress,
								WethAddr: chainAddresses.WrappedNativeAddress,
								// Allow timelock and deployer key to set prices.
								// Deployer key should be removed sometime after initial deployment
								PriceUpdaters: []common.Address{chainAddresses.TimelockAddress, chain.DeployerKey.From},
							},
						}, opsutil.RetryDeploymentWithGasBoost[ccipopsv1_6.DeployFeeQInput](gasBoostConfigs[chainSelector]))
						if err != nil {
							return fmt.Errorf("failed to deploy fee quoter for %s: %w", chain, err)
						}
						feeQuoterAddress = report.Output.Address
						newAddresses[feeQuoterAddress.Hex()] = report.Output.TypeAndVersion
					} else {
						feeQuoterAddress = chainAddresses.FeeAggregatorAddress
					}
					// Deploy OnRamp if not already deployed
					if chainAddresses.OnRampAddress == (common.Address{}) {
						// if the fee aggregator is not set, use the deployer key address
						// this is to ensure that feeAggregator is not set to zero address, otherwise there is chance of
						// fund loss when WithdrawFeeToken is called on OnRamp
						feeAggregator := chainAddresses.FeeAggregatorAddress
						if feeAggregator == (common.Address{}) {
							feeAggregator = chain.DeployerKey.From
						}
						report, err := operations.ExecuteOperation(b, ccipopsv1_6.DeployOnRampOp, chain, opsutil.EVMDeployInput[ccipopsv1_6.DeployOnRampInput]{
							ChainSelector: chainSelector,
							DeployInput: ccipopsv1_6.DeployOnRampInput{
								ChainSelector:      chainSelector,
								TokenAdminRegistry: chainAddresses.TokenAdminRegistryAddress,
								NonceManager:       nonceManagerAddress,
								RmnRemote:          chainAddresses.RMNProxyAddress,
								FeeQuoter:          feeQuoterAddress,
								FeeAggregator:      feeAggregator,
							},
						}, opsutil.RetryDeploymentWithGasBoost[ccipopsv1_6.DeployOnRampInput](gasBoostConfigs[chainSelector]))
						if err != nil {
							return fmt.Errorf("failed to deploy on ramp for %s: %w", chain, err)
						}
						onRampAddress = report.Output.Address
						newAddresses[onRampAddress.Hex()] = report.Output.TypeAndVersion
					} else {
						onRampAddress = chainAddresses.OnRampAddress
					}
					// Deploy OffRamp if not already deployed
					if chainAddresses.OffRampAddress == (common.Address{}) {
						report, err := operations.ExecuteOperation(b, ccipopsv1_6.DeployOffRampOp, chain, opsutil.EVMDeployInput[ccipopsv1_6.DeployOffRampInput]{
							ChainSelector: chainSelector,
							DeployInput: ccipopsv1_6.DeployOffRampInput{
								Chain:              chainSelector,
								Params:             contractParams.OffRampParams,
								RmnRemote:          chainAddresses.RMNProxyAddress,
								NonceManager:       nonceManagerAddress,
								TokenAdminRegistry: chainAddresses.TokenAdminRegistryAddress,
								FeeQuoter:          feeQuoterAddress,
							},
						}, opsutil.RetryDeploymentWithGasBoost[ccipopsv1_6.DeployOffRampInput](gasBoostConfigs[chainSelector]))
						if err != nil {
							return fmt.Errorf("failed to deploy off ramp for %s: %w", chain, err)
						}
						offRampAddress = report.Output.Address
						newAddresses[offRampAddress.Hex()] = report.Output.TypeAndVersion
					} else {
						offRampAddress = chainAddresses.OffRampAddress
					}

					// Add offRamp as an authorized caller on the FeeQuoter
					// ApplyAuthorizedCaller updates are idempotent, so we can safely call them even if the offRamp and onRamp are already authorized callers
					_, err = operations.ExecuteOperation(b, ccipopsv1_6.FeeQApplyAuthorizedCallerOp, chain, opsutil.EVMCallInput[fee_quoter.AuthorizedCallersAuthorizedCallerArgs]{
						ChainSelector: chainSelector,
						NoSend:        false,
						Address:       feeQuoterAddress,
						CallInput: fee_quoter.AuthorizedCallersAuthorizedCallerArgs{
							AddedCallers: []common.Address{offRampAddress},
						},
					}, opsutil.RetryCallWithGasBoost[fee_quoter.AuthorizedCallersAuthorizedCallerArgs](gasBoostConfigs[chainSelector]))
					if err != nil {
						return fmt.Errorf("failed to set off ramp as authorized caller of FeeQuoter on chain %s: %w", chain, err)
					}
					// Add offRamp and onRamp as authorized callers on the NonceManager
					_, err = operations.ExecuteOperation(b, ccipopsv1_6.NonceManagerUpdateAuthorizedCallerOp, chain,
						opsutil.EVMCallInput[nonce_manager.AuthorizedCallersAuthorizedCallerArgs]{
							ChainSelector: chainSelector,
							NoSend:        false,
							Address:       nonceManagerAddress,
							CallInput: nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
								AddedCallers: []common.Address{
									offRampAddress,
									onRampAddress,
								},
							},
						}, opsutil.RetryCallWithGasBoost[nonce_manager.AuthorizedCallersAuthorizedCallerArgs](gasBoostConfigs[chainSelector]))
					if err != nil {
						return fmt.Errorf("failed to set off ramp and on ramp as authorized callers of NonceManager on chain %s: %w", chain, err)
					}

					mu.Lock()
					out[chainSelector] = newAddresses
					mu.Unlock()

					return nil
				})
			}
			if err := grp.Wait(); err != nil {
				return nil, fmt.Errorf("failed to deploy chain contracts: %w", err)
			}
			return out, nil
		})
)

func isRMNRemoteInitialSetUpCompleted(ctx context.Context, rmnHomeAddress common.Address, homeChain cldf_evm.Chain, rmnRemoteAddress common.Address, chain cldf_evm.Chain) ([32]byte, bool, error) {
	rmnHome, err := rmn_home.NewRMNHome(rmnHomeAddress, homeChain.Client)
	if err != nil {
		return [32]byte{}, false, fmt.Errorf("failed to bind to RMNHome contract with address %s: %w", rmnHomeAddress, err)
	}
	activeDigest, err := rmnHome.GetActiveDigest(&bind.CallOpts{Context: ctx})
	if err != nil {
		return [32]byte{}, false, fmt.Errorf("failed to get active digest from RMNHome contract with address %s: %w", rmnHomeAddress, err)
	}

	rmnRemote, err := rmn_remote.NewRMNRemote(rmnRemoteAddress, chain.Client)
	if err != nil {
		return [32]byte{}, false, fmt.Errorf("failed to bind to RMNRemote contract on chain %s with address %s: %w", chain, rmnRemoteAddress, err)
	}
	// Get the existing config from RMNRemote
	existingConfig, err := rmnRemote.GetVersionedConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return [32]byte{}, false, fmt.Errorf("failed to get versioned config from RMNRemote on chain %s with address %s: %w", chain, rmnRemoteAddress, err)
	}
	// Is the config already set?
	// If the config is already set, the version should be more than 0, and we can check if it matches the active digest on RMNHome
	// In this case, we don't need to set it again on existing RMNRemote
	if existingConfig.Version > 0 && existingConfig.Config.RmnHomeContractConfigDigest == activeDigest {
		return activeDigest, true, nil
	}
	return activeDigest, false, nil
}
