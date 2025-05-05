package v1_6

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

/*
Future improvements:
- Enable connecting directly to the prod router through PromoteNewChainForConfigChangeset
- Align configuration structs with whatever is simplest for BIX team to use
- Add more validation coverage to the precondition functions
- Use within add_chain integration test
*/

var (
	// AddCandidatesForNewChainChangeset deploys a new chain and adds its exec and commit plugins as candidates on the home chain.
	// This changeset is not idempotent because the underlying AddDonAndSetCandidateChangeset is not idempotent.
	// Provide an MCMS config if the contracts on the existing chains are owned by MCMS (omit this config otherwise).
	AddCandidatesForNewChainChangeset = cldf.CreateChangeSet(addCandidatesForNewChainLogic, addCandidatesForNewChainPrecondition)
	// PromoteNewChainForConfigChangeset promotes exec and commit plugin candidates for the new chain on the home chain.
	// It also connects the new chain to various destination chains through the test router.
	// This changeset should be run after AddCandidatesForNewChainChangeset.
	// This changeset is not idempotent because the underlying PromoteCandidateChangeset is not idepotent.
	// Provide an MCMS config if the contracts on the existing chains are owned by MCMS (omit this config otherwise).
	PromoteNewChainForConfigChangeset = cldf.CreateChangeSet(promoteNewChainForConfigLogic, promoteNewChainForConfigPrecondition)
	// ConnectNewChainChangeset activates connects a new chain with other chains by updating onRamp, offRamp, and router contracts.
	// If connecting to production routers, you should have already run PromoteNewChainForConfigChangeset.
	// Rerunning this changeset with a given input will produce the same results each time (outside of ownership transfers, which only happen once).
	// Provide an MCMS config if the contracts on the existing chains are owned by MCMS (omit this config otherwise).
	ConnectNewChainChangeset = cldf.CreateChangeSet(connectNewChainLogic, connectNewChainPrecondition)
)

// /////////////////////////////////
// START AddCandidatesForNewChainChangeset
// /////////////////////////////////
// ChainDefinition defines how a chain should be configured on both remote chains and itself.
type ChainDefinition struct {
	// ConnectionConfig holds configuration for connection.
	ConnectionConfig `json:"connectionConfig"`
	// Selector is the chain selector of this chain.
	Selector uint64 `json:"selector"`
	// GasPrice defines the USD price (18 decimals) per unit gas for this chain as a destination.
	GasPrice *big.Int `json:"gasPrice"`
	// TokenPrices define the USD price (18 decimals) per 1e18 of the smallest token denomination for various tokens on this chain.
	TokenPrices map[common.Address]*big.Int `json:"tokenPrices"`
	// FeeQuoterDestChainConfig is the configuration on a fee quoter for this chain as a destination.
	FeeQuoterDestChainConfig fee_quoter.FeeQuoterDestChainConfig `json:"feeQuoterDestChainConfig"`
}

// NewChainDefinition defines how a NEW chain should be configured.
type NewChainDefinition struct {
	// ChainDefinition holds basic chain info.
	ChainDefinition `json:"chainDefinition"`
	// ChainContractParams defines contract parameters for the chain.
	ChainContractParams `json:"chainContractParams"`
	// ExistingContracts defines any contracts that are already deployed on this chain.
	ExistingContracts commoncs.ExistingContractsConfig `json:"existingContracts"`
	// ConfigOnHome defines how this chain should be configured on the CCIPHome contract.
	ConfigOnHome ChainConfig `json:"configOnHome"`
	// CommitOCRParams defines the OCR parameters for this chain's commit plugin.
	CommitOCRParams CCIPOCRParams `json:"commitOcrParams"`
	// ExecOCRParams defines the OCR parameters for this chain's exec plugin.
	ExecOCRParams CCIPOCRParams `json:"execOcrParams"`
	// RMNRemoteConfig is the config for the RMNRemote contract.
	RMNRemoteConfig *RMNRemoteConfig `json:"rmnRemoteConfig,omitempty"`
}

// AddCandidatesForNewChainConfig is a configuration struct for AddCandidatesForNewChainChangeset.
type AddCandidatesForNewChainConfig struct {
	// HomeChainSelector is the selector of the home chain.
	HomeChainSelector uint64 `json:"homeChainSelector"`
	// FeedChainSelector is the selector of the chain on which price feeds are deployed.
	FeedChainSelector uint64 `json:"feedChainSelector"`
	// NewChain defines the new chain to be deployed.
	NewChain NewChainDefinition `json:"newChain"`
	// RemoteChains defines the remote chains to be connected to the new chain.
	RemoteChains []ChainDefinition `json:"remoteChains"`
	// MCMSDeploymentConfig configures the MCMS deployment to the new chain.
	MCMSDeploymentConfig *commontypes.MCMSWithTimelockConfigV2 `json:"mcmsDeploymentConfig,omitempty"`
	// MCMSConfig defines the MCMS configuration for the changeset.
	MCMSConfig *proposalutils.TimelockConfig `json:"mcmsConfig,omitempty"`
}

func (c AddCandidatesForNewChainConfig) prerequisiteConfigForNewChain() changeset.DeployPrerequisiteConfig {
	return changeset.DeployPrerequisiteConfig{
		Configs: []changeset.DeployPrerequisiteConfigPerChain{
			changeset.DeployPrerequisiteConfigPerChain{
				ChainSelector: c.NewChain.Selector,
			},
		},
	}
}

func (c AddCandidatesForNewChainConfig) deploymentConfigForNewChain() DeployChainContractsConfig {
	return DeployChainContractsConfig{
		HomeChainSelector: c.HomeChainSelector,
		ContractParamsPerChain: map[uint64]ChainContractParams{
			c.NewChain.Selector: c.NewChain.ChainContractParams,
		},
	}
}

func (c AddCandidatesForNewChainConfig) rmnRemoteConfigForNewChain() SetRMNRemoteConfig {
	if c.NewChain.RMNRemoteConfig == nil {
		return SetRMNRemoteConfig{}
	}
	return SetRMNRemoteConfig{
		HomeChainSelector: c.HomeChainSelector,
		RMNRemoteConfigs: map[uint64]RMNRemoteConfig{
			c.NewChain.Selector: *c.NewChain.RMNRemoteConfig,
		},
	}
}

func (c AddCandidatesForNewChainConfig) updateChainConfig() UpdateChainConfigConfig {
	return UpdateChainConfigConfig{
		HomeChainSelector: c.HomeChainSelector,
		RemoteChainAdds: map[uint64]ChainConfig{
			c.NewChain.Selector: c.NewChain.ConfigOnHome,
		},
		MCMS: c.MCMSConfig,
	}
}

func addCandidatesForNewChainPrecondition(e deployment.Environment, c AddCandidatesForNewChainConfig) error {
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	err = changeset.ValidateChain(e, state, c.HomeChainSelector, c.MCMSConfig)
	if err != nil {
		return fmt.Errorf("failed to validate home chain: %w", err)
	}
	homeChainState := state.Chains[c.HomeChainSelector]
	if homeChainState.CCIPHome == nil {
		return fmt.Errorf("home chain with selector %d does not have a CCIPHome", c.HomeChainSelector)
	}
	if homeChainState.CapabilityRegistry == nil {
		return fmt.Errorf("home chain with selector %d does not have a CapabilitiesRegistry", c.HomeChainSelector)
	}

	// We pre-validate any changesets that do not rely on contracts being deployed.
	// The following can't be easily pre-validated:
	// SetRMNRemoteOnRMNProxyChangeset, UpdateFeeQuoterDestsChangeset, UpdateFeeQuoterPricesChangeset
	if err := c.NewChain.ExistingContracts.Validate(); err != nil {
		return fmt.Errorf("failed to validate existing contracts on new chain: %w", err)
	}
	if err := c.prerequisiteConfigForNewChain().Validate(); err != nil {
		return fmt.Errorf("failed to validate prerequisite config for new chain: %w", err)
	}
	if err := c.deploymentConfigForNewChain().Validate(); err != nil {
		return fmt.Errorf("failed to validate deployment config for new chain: %w", err)
	}
	if c.NewChain.RMNRemoteConfig != nil {
		if err := c.rmnRemoteConfigForNewChain().Validate(); err != nil {
			return fmt.Errorf("failed to validate RMN remote config for new chain: %w", err)
		}
	}
	if err := c.updateChainConfig().Validate(e); err != nil {
		return fmt.Errorf("failed to validate update chain config: %w", err)
	}

	return nil
}

func addCandidatesForNewChainLogic(e deployment.Environment, c AddCandidatesForNewChainConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()
	var allProposals []mcmslib.TimelockProposal

	// Save existing contracts
	err := runAndSaveAddresses(func() (deployment.ChangesetOutput, error) {
		return commoncs.SaveExistingContractsChangeset(e, c.NewChain.ExistingContracts)
	}, newAddresses, e.ExistingAddresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run SaveExistingContractsChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// Deploy the prerequisite contracts to the new chain
	err = runAndSaveAddresses(func() (deployment.ChangesetOutput, error) {
		return changeset.DeployPrerequisitesChangeset(e, c.prerequisiteConfigForNewChain())
	}, newAddresses, e.ExistingAddresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run DeployPrerequisitesChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// Deploy MCMS contracts
	if c.MCMSDeploymentConfig != nil {
		err = runAndSaveAddresses(func() (deployment.ChangesetOutput, error) {
			return commoncs.DeployMCMSWithTimelockV2(e, map[uint64]commontypes.MCMSWithTimelockConfigV2{
				c.NewChain.Selector: *c.MCMSDeploymentConfig,
			})
		}, newAddresses, e.ExistingAddresses)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to run DeployMCMSWithTimelockV2 on chain with selector %d: %w", c.NewChain.Selector, err)
		}
	}

	// Deploy chain contracts to the new chain
	err = runAndSaveAddresses(func() (deployment.ChangesetOutput, error) {
		return DeployChainContractsChangeset(e, c.deploymentConfigForNewChain())
	}, newAddresses, e.ExistingAddresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run DeployChainContractsChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// Set RMN remote config & set RMN on proxy on the new chain (if config provided)
	if c.NewChain.RMNRemoteConfig != nil {
		_, err = SetRMNRemoteConfigChangeset(e, c.rmnRemoteConfigForNewChain())
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to run SetRMNRemoteConfigChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
		}
	}
	// Set the RMN remote on the RMN proxy, using MCMS if RMN proxy is owned by Timelock
	// RMN proxy will already exist on chains that supported CCIPv1.5.0, in which case RMN proxy will be owned by Timelock
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	owner, err := state.Chains[c.NewChain.Selector].RMNProxy.Owner(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get owner of RMN proxy on chain with selector %d: %w", c.NewChain.Selector, err)
	}
	var mcmsConfig *proposalutils.TimelockConfig
	if owner == state.Chains[c.NewChain.Selector].Timelock.Address() {
		mcmsConfig = c.MCMSConfig
	}
	out, err := SetRMNRemoteOnRMNProxyChangeset(e, SetRMNRemoteOnRMNProxyConfig{
		ChainSelectors: []uint64{c.NewChain.Selector},
		MCMSConfig:     mcmsConfig,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run SetRMNRemoteOnRMNProxyChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// Update the fee quoter destinations on the new chain
	destChainConfigs := make(map[uint64]fee_quoter.FeeQuoterDestChainConfig, len(c.RemoteChains))
	for _, remoteChain := range c.RemoteChains {
		destChainConfigs[remoteChain.Selector] = remoteChain.FeeQuoterDestChainConfig
	}
	_, err = UpdateFeeQuoterDestsChangeset(e, UpdateFeeQuoterDestsConfig{
		UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
			c.NewChain.Selector: destChainConfigs,
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateFeeQuoterDestsChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// Update the fee quoter prices on the new chain
	gasPrices := make(map[uint64]*big.Int, len(c.RemoteChains))
	for _, remoteChain := range c.RemoteChains {
		gasPrices[remoteChain.Selector] = remoteChain.GasPrice
	}
	_, err = UpdateFeeQuoterPricesChangeset(e, UpdateFeeQuoterPricesConfig{
		PricesByChain: map[uint64]FeeQuoterPriceUpdatePerSource{
			c.NewChain.Selector: FeeQuoterPriceUpdatePerSource{
				TokenPrices: c.NewChain.TokenPrices,
				GasPrices:   gasPrices,
			},
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateFeeQuoterPricesChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// Fetch the next DON ID from the capabilities registry
	state, err = changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	donID, err := state.Chains[c.HomeChainSelector].CapabilityRegistry.GetNextDONId(&bind.CallOpts{
		Context: e.GetContext(),
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get next DON ID: %w", err)
	}

	// Add new chain config to the home chain
	out, err = UpdateChainConfigChangeset(e, c.updateChainConfig())
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateChainConfigChangeset on home chain: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// Add the DON to the registry and set candidate for the commit plugin
	out, err = AddDonAndSetCandidateChangeset(e, AddDonAndSetCandidateChangesetConfig{
		SetCandidateConfigBase: SetCandidateConfigBase{
			HomeChainSelector: c.HomeChainSelector,
			FeedChainSelector: c.FeedChainSelector,
			MCMS:              c.MCMSConfig,
		},
		PluginInfo: SetCandidatePluginInfo{
			PluginType: types.PluginTypeCCIPCommit,
			OCRConfigPerRemoteChainSelector: map[uint64]CCIPOCRParams{
				c.NewChain.Selector: c.NewChain.CommitOCRParams,
			},
			SkipChainConfigValidation: true,
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run AddDonAndSetCandidateChangeset on home chain: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// Set the candidate for the exec plugin
	out, err = SetCandidateChangeset(e, SetCandidateChangesetConfig{
		SetCandidateConfigBase: SetCandidateConfigBase{
			HomeChainSelector: c.HomeChainSelector,
			FeedChainSelector: c.FeedChainSelector,
			MCMS:              c.MCMSConfig,
		},
		PluginInfo: []SetCandidatePluginInfo{
			{
				PluginType: types.PluginTypeCCIPExec,
				OCRConfigPerRemoteChainSelector: map[uint64]CCIPOCRParams{
					c.NewChain.Selector: c.NewChain.ExecOCRParams,
				},
				SkipChainConfigValidation: true,
			},
		},
		DonIDOverrides: map[uint64]uint32{c.NewChain.Selector: donID},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run SetCandidateChangeset on home chain: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// Reset existing addresses
	// This is for compatibility with in-memory tests, where we merge the new addresses into the environment immediately after running the changeset
	// If we don't reset the existing addresses mapping, merging will fail because the addresses will already exist there
	err = e.ExistingAddresses.Remove(newAddresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to reset existing addresses: %w", err)
	}

	proposal, err := proposalutils.AggregateProposals(
		e,
		state.EVMMCMSStateByChain(),
		allProposals,
		nil,
		fmt.Sprintf("Deploy and set candidates for chain with selector %d", c.NewChain.Selector),
		c.MCMSConfig,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	if proposal == nil {
		return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
	}
	return deployment.ChangesetOutput{AddressBook: newAddresses, MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

///////////////////////////////////
// END AddCandidatesForNewChainChangeset
///////////////////////////////////

///////////////////////////////////
// START PromoteNewChainForConfigChangeset
///////////////////////////////////

// PromoteNewChainForConfig is a configuration struct for PromoteNewChainForConfigChangeset.
type PromoteNewChainForConfig struct {
	// HomeChainSelector is the selector of the home chain.
	HomeChainSelector uint64 `json:"homeChainSelector"`
	// NewChain defines the new chain to be deployed.
	NewChain NewChainDefinition `json:"newChain"`
	// RemoteChains defines the remote chains to be connected to the new chain.
	RemoteChains []ChainDefinition `json:"remoteChains"`
	// TestRouter is true if we want to connect via test routers.
	TestRouter *bool `json:"testRouter,omitempty"`
	// MCMSConfig defines the MCMS configuration for the changeset.
	MCMSConfig *proposalutils.TimelockConfig `json:"mcmsConfig,omitempty"`
}

func (c PromoteNewChainForConfig) promoteCandidateConfig() PromoteCandidateChangesetConfig {
	return PromoteCandidateChangesetConfig{
		HomeChainSelector: c.HomeChainSelector,
		MCMS:              c.MCMSConfig,
		PluginInfo: []PromoteCandidatePluginInfo{
			{
				PluginType:           types.PluginTypeCCIPCommit,
				RemoteChainSelectors: []uint64{c.NewChain.Selector},
			},
			{
				PluginType:           types.PluginTypeCCIPExec,
				RemoteChainSelectors: []uint64{c.NewChain.Selector},
			},
		},
	}
}

func (c PromoteNewChainForConfig) setOCR3OffRampConfig() SetOCR3OffRampConfig {
	candidate := globals.ConfigTypeActive
	if c.MCMSConfig != nil {
		candidate = globals.ConfigTypeCandidate // If going through MCMS, the config will be candidate during changeset validation
	}
	return SetOCR3OffRampConfig{
		HomeChainSel:       c.HomeChainSelector,
		RemoteChainSels:    []uint64{c.NewChain.Selector},
		CCIPHomeConfigType: candidate,
	}
}

func (c PromoteNewChainForConfig) updateFeeQuoterDestsConfig(remoteChain ChainDefinition) UpdateFeeQuoterDestsConfig {
	return UpdateFeeQuoterDestsConfig{
		UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
			remoteChain.Selector: map[uint64]fee_quoter.FeeQuoterDestChainConfig{
				c.NewChain.Selector: c.NewChain.FeeQuoterDestChainConfig,
			},
		},
		MCMS: c.MCMSConfig,
	}
}

func (c PromoteNewChainForConfig) updateFeeQuoterPricesConfig(remoteChain ChainDefinition) UpdateFeeQuoterPricesConfig {
	return UpdateFeeQuoterPricesConfig{
		PricesByChain: map[uint64]FeeQuoterPriceUpdatePerSource{
			remoteChain.Selector: FeeQuoterPriceUpdatePerSource{
				TokenPrices: remoteChain.TokenPrices,
				GasPrices:   map[uint64]*big.Int{c.NewChain.Selector: c.NewChain.GasPrice},
			},
		},
		MCMS: c.MCMSConfig,
	}
}

func (c PromoteNewChainForConfig) connectNewChainConfig(testRouter bool) ConnectNewChainConfig {
	connections := make(map[uint64]ConnectionConfig, len(c.RemoteChains))
	for _, remoteChain := range c.RemoteChains {
		connections[remoteChain.Selector] = remoteChain.ConnectionConfig
	}
	return ConnectNewChainConfig{
		RemoteChains:             connections,
		NewChainSelector:         c.NewChain.Selector,
		NewChainConnectionConfig: c.NewChain.ConnectionConfig,
		TestRouter:               &testRouter,
		MCMSConfig:               c.MCMSConfig,
	}
}

func promoteNewChainForConfigPrecondition(e deployment.Environment, c PromoteNewChainForConfig) error {
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	if _, err := c.promoteCandidateConfig().Validate(e); err != nil {
		return fmt.Errorf("failed to validate promote candidate config: %w", err)
	}

	if err := c.setOCR3OffRampConfig().Validate(e, state); err != nil {
		return fmt.Errorf("failed to validate set OCR3 off ramp config: %w", err)
	}

	for _, remoteChain := range c.RemoteChains {
		if err := c.updateFeeQuoterDestsConfig(remoteChain).Validate(e); err != nil {
			return fmt.Errorf("failed to validate update fee quoter dests config for remote chain with selector %d: %w", remoteChain.Selector, err)
		}
		if err := c.updateFeeQuoterPricesConfig(remoteChain).Validate(e); err != nil {
			return fmt.Errorf("failed to validate update fee quoter prices config for remote chain with selector %d: %w", remoteChain.Selector, err)
		}
	}

	err = ConnectNewChainChangeset.VerifyPreconditions(e, c.connectNewChainConfig(*c.TestRouter))
	if err != nil {
		return fmt.Errorf("failed to validate ConnectNewChainChangeset: %w", err)
	}

	return nil
}

func promoteNewChainForConfigLogic(e deployment.Environment, c PromoteNewChainForConfig) (deployment.ChangesetOutput, error) {
	var allProposals []mcmslib.TimelockProposal
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Promote the candidates for the commit and exec plugins
	out, err := PromoteCandidateChangeset(e, c.promoteCandidateConfig())
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run PromoteCandidateChangeset on home chain: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	// Set the OCR3 config on the off ramp on the new chain
	out, err = SetOCR3OffRampChangeset(e, c.setOCR3OffRampConfig())
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run SetOCR3OffRampChangeset on chain with selector %d: %w", c.NewChain.Selector, err)
	}

	// Update the fee quoter prices and destinations on the remote chains
	for _, remoteChain := range c.RemoteChains {
		out, err = UpdateFeeQuoterDestsChangeset(e, c.updateFeeQuoterDestsConfig(remoteChain))
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateFeeQuoterDestsChangeset on chain with selector %d: %w", remoteChain.Selector, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)

		out, err = UpdateFeeQuoterPricesChangeset(e, c.updateFeeQuoterPricesConfig(remoteChain))
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to run UpdateFeeQuoterPricesChangeset on chain with selector %d: %w", remoteChain.Selector, err)
		}
		allProposals = append(allProposals, out.MCMSTimelockProposals...)
	}

	// Connect the new chain to the existing chains (use the test router)
	out, err = ConnectNewChainChangeset.Apply(e, c.connectNewChainConfig(*c.TestRouter))
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run ConnectNewChainChangeset: %w", err)
	}
	allProposals = append(allProposals, out.MCMSTimelockProposals...)

	proposal, err := proposalutils.AggregateProposals(
		e,
		state.EVMMCMSStateByChain(),
		allProposals,
		nil,
		fmt.Sprintf("Promote chain with selector %d for testing", c.NewChain.Selector),
		c.MCMSConfig,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	if proposal == nil {
		return deployment.ChangesetOutput{}, nil
	}
	return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

///////////////////////////////////
// END PromoteNewChainForConfigChangeset
///////////////////////////////////

///////////////////////////////////
// START ConnectNewChainChangeset
///////////////////////////////////

// ConnectionConfig defines how a chain should connect with other chains.
type ConnectionConfig struct {
	// RMNVerificationDisabled is true if we do not want the RMN to bless messages FROM this chain.
	RMNVerificationDisabled bool `json:"rmnVerificationDisabled"`
	// AllowListEnabled is true if we want an allowlist to dictate who can send messages TO this chain.
	AllowListEnabled bool `json:"allowListEnabled"`
}

// ConnectNewChainConfig is a configuration struct for ConnectNewChainChangeset.
type ConnectNewChainConfig struct {
	// NewChainSelector is the selector of the new chain to connect.
	NewChainSelector uint64 `json:"newChainSelector"`
	// NewChainConnectionConfig defines how the new chain should connect with other chains.
	NewChainConnectionConfig ConnectionConfig `json:"newChainConnectionConfig"`
	// RemoteChains are the chains to connect the new chain to.
	RemoteChains map[uint64]ConnectionConfig `json:"remoteChains"`
	// TestRouter is true if we want to connect via test routers.
	TestRouter *bool `json:"testRouter,omitempty"`
	// MCMSConfig is the MCMS configuration, omit to use deployer key only.
	MCMSConfig *proposalutils.TimelockConfig `json:"mcmsConfig,omitempty"`
}

func (c ConnectNewChainConfig) validateNewChain(env deployment.Environment, state changeset.CCIPOnChainState) error {
	// When running this changeset, there is no case in which the new chain contract should be owned by MCMS,
	// which is why we do not use MCMSConfig to determine the ownedByMCMS variable.
	err := c.validateChain(env, state, c.NewChainSelector, false)
	if err != nil {
		return fmt.Errorf("failed to validate chain with selector %d: %w", c.NewChainSelector, err)
	}

	return nil
}

func (c ConnectNewChainConfig) validateRemoteChains(env deployment.Environment, state changeset.CCIPOnChainState) error {
	for remoteChainSelector := range c.RemoteChains {
		// The remote chain may or may not be owned by MCMS, as MCMS is not really used in staging.
		// Therefore, we use the presence of MCMSConfig to determine the ownedByMCMS variable.
		err := c.validateChain(env, state, remoteChainSelector, c.MCMSConfig != nil)
		if err != nil {
			return fmt.Errorf("failed to validate chain with selector %d: %w", remoteChainSelector, err)
		}
	}

	return nil
}

func (c ConnectNewChainConfig) validateChain(e deployment.Environment, state changeset.CCIPOnChainState, chainSelector uint64, ownedByMCMS bool) error {
	err := changeset.ValidateChain(e, state, chainSelector, c.MCMSConfig)
	if err != nil {
		return fmt.Errorf("failed to validate chain with selector %d: %w", chainSelector, err)
	}
	chainState := state.Chains[chainSelector]
	deployerKey := e.Chains[chainSelector].DeployerKey.From

	if chainState.OnRamp == nil {
		return errors.New("onRamp contract not found")
	}
	if chainState.OffRamp == nil {
		return errors.New("offRamp contract not found")
	}
	if chainState.Router == nil {
		return errors.New("router contract not found")
	}
	if chainState.TestRouter == nil {
		return errors.New("test router contract not found")
	}

	err = commoncs.ValidateOwnership(e.GetContext(), ownedByMCMS, deployerKey, chainState.Timelock.Address(), chainState.OnRamp)
	if err != nil {
		return fmt.Errorf("failed to validate ownership of onRamp: %w", err)
	}
	err = commoncs.ValidateOwnership(e.GetContext(), ownedByMCMS, deployerKey, chainState.Timelock.Address(), chainState.OffRamp)
	if err != nil {
		return fmt.Errorf("failed to validate ownership of offRamp: %w", err)
	}
	err = commoncs.ValidateOwnership(e.GetContext(), ownedByMCMS, deployerKey, chainState.Timelock.Address(), chainState.Router)
	if err != nil {
		return fmt.Errorf("failed to validate ownership of router: %w", err)
	}

	// Test router should always be owned by deployer key
	err = commoncs.ValidateOwnership(e.GetContext(), false, deployerKey, chainState.Timelock.Address(), chainState.TestRouter)
	if err != nil {
		return fmt.Errorf("failed to validate ownership of test router: %w", err)
	}

	return nil
}

func connectNewChainPrecondition(env deployment.Environment, c ConnectNewChainConfig) error {
	if c.TestRouter == nil {
		return errors.New("must define whether to use the test router")
	}

	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	err = c.validateNewChain(env, state)
	if err != nil {
		return fmt.Errorf("failed to validate new chain: %w", err)
	}

	err = c.validateRemoteChains(env, state)
	if err != nil {
		return fmt.Errorf("failed to validate remote chains: %w", err)
	}

	return nil
}

func connectNewChainLogic(env deployment.Environment, c ConnectNewChainConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	readOpts := &bind.CallOpts{Context: env.GetContext()}

	var ownershipTransferProposals []timelock.MCMSWithTimelockProposal
	if !*c.TestRouter && c.MCMSConfig != nil {
		// If using the production router, transfer ownership of all contracts on the new chain to MCMS.
		allContracts := []commoncs.Ownable{
			state.Chains[c.NewChainSelector].OnRamp,
			state.Chains[c.NewChainSelector].OffRamp,
			state.Chains[c.NewChainSelector].FeeQuoter,
			state.Chains[c.NewChainSelector].RMNProxy,
			state.Chains[c.NewChainSelector].NonceManager,
			state.Chains[c.NewChainSelector].TokenAdminRegistry,
			state.Chains[c.NewChainSelector].Router,
			state.Chains[c.NewChainSelector].RMNRemote,
		}
		addressesToTransfer := make([]common.Address, 0, len(allContracts))
		for _, contract := range allContracts {
			if contract == nil {
				continue
			}
			owner, err := contract.Owner(readOpts)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to get owner of contract %s: %w", contract.Address().Hex(), err)
			}
			if owner == env.Chains[c.NewChainSelector].DeployerKey.From {
				addressesToTransfer = append(addressesToTransfer, contract.Address())
			}
		}
		out, err := commoncs.TransferToMCMSWithTimelock(env, commoncs.TransferToMCMSWithTimelockConfig{
			ContractsByChain: map[uint64][]common.Address{
				c.NewChainSelector: addressesToTransfer,
			},
			MCMSConfig: *c.MCMSConfig,
		})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to run TransferToMCMSWithTimelock on chain with selector %d: %w", c.NewChainSelector, err)
		}
		ownershipTransferProposals = out.Proposals //nolint:staticcheck //SA1019 ignoring deprecated function for compatibility

		// Also, renounce the admin role on the Timelock (if not already done).
		adminRole, err := state.Chains[c.NewChainSelector].Timelock.ADMINROLE(readOpts)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get admin role of timelock on chain with selector %d: %w", c.NewChainSelector, err)
		}
		hasRole, err := state.Chains[c.NewChainSelector].Timelock.HasRole(readOpts, adminRole, env.Chains[c.NewChainSelector].DeployerKey.From)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to check if deployer key has admin role on timelock on chain with selector %d: %w", c.NewChainSelector, err)
		}
		if hasRole {
			out, err = commoncs.RenounceTimelockDeployer(env, commoncs.RenounceTimelockDeployerConfig{
				ChainSel: c.NewChainSelector,
			})
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to run RenounceTimelockDeployer on chain with selector %d: %w", c.NewChainSelector, err)
			}
		}
	}

	// Enable the production router on [new chain -> each remote chain] and [each remote chain -> new chain].
	var allEnablementProposals []mcmslib.TimelockProposal
	var mcmsConfig *proposalutils.TimelockConfig
	if !*c.TestRouter {
		mcmsConfig = c.MCMSConfig
	}
	allEnablementProposals, err = connectRampsAndRouters(env, c.NewChainSelector, c.RemoteChains, mcmsConfig, *c.TestRouter, allEnablementProposals)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to enable production router on chain with selector %d: %w", c.NewChainSelector, err)
	}
	for remoteChainSelector := range c.RemoteChains {
		allEnablementProposals, err = connectRampsAndRouters(env, remoteChainSelector, map[uint64]ConnectionConfig{c.NewChainSelector: c.NewChainConnectionConfig}, c.MCMSConfig, *c.TestRouter, allEnablementProposals)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to enable production router on chain with selector %d: %w", remoteChainSelector, err)
		}
	}

	proposal, err := proposalutils.AggregateProposals(
		env,
		state.EVMMCMSStateByChain(),
		allEnablementProposals,
		ownershipTransferProposals,
		fmt.Sprintf("Connect chain with selector %d to other chains", c.NewChainSelector),
		c.MCMSConfig,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	if proposal == nil {
		return deployment.ChangesetOutput{}, nil
	}
	return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

// connectRampsAndRouters updates the onRamp and offRamp to point at the router for the given remote chains.
// It also sets the onRamp and offRamp on the router for the given remote chains.
// This function will add the proposals required to make these changes to the proposalAggregate slice.
func connectRampsAndRouters(
	e deployment.Environment,
	chainSelector uint64,
	remoteChains map[uint64]ConnectionConfig,
	mcmsConfig *proposalutils.TimelockConfig,
	testRouter bool,
	proposalAggregate []mcmslib.TimelockProposal,
) ([]mcmslib.TimelockProposal, error) {
	// Update offRamp sources on the new chain.
	offRampUpdatesOnNew := make(map[uint64]OffRampSourceUpdate, len(remoteChains))
	for remoteChainSelector, remoteChain := range remoteChains {
		offRampUpdatesOnNew[remoteChainSelector] = OffRampSourceUpdate{
			TestRouter:                testRouter,
			IsRMNVerificationDisabled: remoteChain.RMNVerificationDisabled,
			IsEnabled:                 true,
		}
	}
	out, err := UpdateOffRampSourcesChangeset(e, UpdateOffRampSourcesConfig{
		UpdatesByChain: map[uint64]map[uint64]OffRampSourceUpdate{
			chainSelector: offRampUpdatesOnNew,
		},
		MCMS:               mcmsConfig,
		SkipOwnershipCheck: true,
	})
	if err != nil {
		return []mcmslib.TimelockProposal{}, fmt.Errorf("failed to run UpdateOffRampSourcesChangeset on chain with selector %d: %w", chainSelector, err)
	}
	proposalAggregate = append(proposalAggregate, out.MCMSTimelockProposals...)

	// Update onRamp destinations on the new chain.
	onRampUpdatesOnNew := make(map[uint64]OnRampDestinationUpdate, len(remoteChains))
	for remoteChainSelector, remoteChain := range remoteChains {
		onRampUpdatesOnNew[remoteChainSelector] = OnRampDestinationUpdate{
			TestRouter:       testRouter,
			AllowListEnabled: remoteChain.AllowListEnabled,
			IsEnabled:        true,
		}
	}
	out, err = UpdateOnRampsDestsChangeset(e, UpdateOnRampDestsConfig{
		UpdatesByChain: map[uint64]map[uint64]OnRampDestinationUpdate{
			chainSelector: onRampUpdatesOnNew,
		},
		MCMS:               mcmsConfig,
		SkipOwnershipCheck: true,
	})
	if err != nil {
		return []mcmslib.TimelockProposal{}, fmt.Errorf("failed to run UpdateOnRampsDestsChangeset on chain with selector %d: %w", chainSelector, err)
	}
	proposalAggregate = append(proposalAggregate, out.MCMSTimelockProposals...)

	// Update router ramps on the new chain.
	offRampUpdates := make(map[uint64]bool, len(remoteChains))
	onRampUpdates := make(map[uint64]bool, len(remoteChains))
	for remoteChainSelector := range remoteChains {
		offRampUpdates[remoteChainSelector] = true
		onRampUpdates[remoteChainSelector] = true
	}
	cfg := mcmsConfig
	if testRouter { // Again, test router does not use MCMS. We are making this assumption.
		cfg = nil
	}
	out, err = UpdateRouterRampsChangeset(e, UpdateRouterRampsConfig{
		TestRouter: testRouter,
		UpdatesByChain: map[uint64]RouterUpdates{
			chainSelector: RouterUpdates{
				OnRampUpdates:  onRampUpdates,
				OffRampUpdates: offRampUpdates,
			},
		},
		MCMS:               cfg,
		SkipOwnershipCheck: true,
	})
	if err != nil {
		return []mcmslib.TimelockProposal{}, fmt.Errorf("failed to run UpdateRouterRampsChangeset on chain with selector %d: %w", chainSelector, err)
	}
	proposalAggregate = append(proposalAggregate, out.MCMSTimelockProposals...)

	return proposalAggregate, nil
}

///////////////////////////////////
// END ConnectNewChainChangeset
///////////////////////////////////

func runAndSaveAddresses(fn func() (deployment.ChangesetOutput, error), newAddresses deployment.AddressBook, existingAddresses deployment.AddressBook) error {
	output, err := fn()
	if err != nil {
		return fmt.Errorf("failed to run changeset: %w", err)
	}
	err = newAddresses.Merge(output.AddressBook)
	if err != nil {
		return fmt.Errorf("failed to update new address book: %w", err)
	}
	err = existingAddresses.Merge(output.AddressBook)
	if err != nil {
		return fmt.Errorf("failed to update existing address book: %w", err)
	}

	return nil
}
