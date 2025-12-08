package contracts

import (
	"errors"
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

type UpdateDONDeps struct {
	Env                  *cldf.Environment
	Strategy             strategies.TransactionStrategy
	CapabilitiesRegistry *capabilities_registry_v2.CapabilitiesRegistry
}

type UpdateDONInput struct {
	ChainSelector uint64

	// P2PIDs are the peer ids that compose the don. Optional, only provided if the DON composition is changing.
	P2PIDs                            []p2pkey.PeerID
	CapabilityConfigs                 []CapabilityConfig
	MergeCapabilityConfigsWithOnChain bool

	// DonName to update, this is required
	DonName string

	// NewDonName is optional
	NewDonName string

	// F is the fault tolerance level
	// if omitted, the existing value fetched from the registry is used
	F uint8

	// IsPrivate indicates whether the DON is public or private
	// If omitted, the existing value fetched from the registry is used
	IsPrivate bool

	// Force indicates whether to force the update even if we cannot validate that all forwarder contracts are ready to accept the new configure version.
	// This is very dangerous, and could break the whole platform if the forwarders are not ready. Be very careful with this option.
	Force bool

	MCMSConfig *contracts.MCMSConfig
}

func (r *UpdateDONInput) Validate() error {
	if r.DonName == "" {
		return errors.New("must specify DONName")
	}

	return nil
}

type UpdateDONOutput struct {
	DonInfo   capabilities_registry_v2.CapabilitiesRegistryUpdateDONParams
	Operation *mcmstypes.BatchOperation
}

// CapabilityConfig is a struct that holds a capability and its configuration
type CapabilityConfig struct {
	Capability Capability
	// Config is the capability configuration. It will be marshalled to proto config.
	// It is untyped here because is has to be deserialized from JSON/YAML for any possible capability
	// If nil, a default config based on the capability type is used
	Config map[string]any
}

type Capability struct {
	CapabilityID          string         `json:"capabilityID" yaml:"capabilityID"`
	ConfigurationContract common.Address `json:"configurationContract" yaml:"configurationContract"`
	// Metadata is the capability metadata. It will be marshalled to json config.
	Metadata map[string]any `json:"metadata" yaml:"metadata"`
}

var UpdateDON = operations.NewOperation[UpdateDONInput, UpdateDONOutput, UpdateDONDeps](
	"update-don-op",
	semver.MustParse("1.0.0"),
	"Update DON in Capabilities Registry",
	func(b operations.Bundle, deps UpdateDONDeps, input UpdateDONInput) (UpdateDONOutput, error) {
		if err := input.Validate(); err != nil {
			return UpdateDONOutput{}, err
		}

		registry := deps.CapabilitiesRegistry

		// DonName is required
		don, err := registry.GetDONByName(&bind.CallOpts{}, input.DonName)
		if err != nil {
			err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
			return UpdateDONOutput{}, fmt.Errorf("failed to call GetDONByName: %w", err)
		}

		if don.AcceptsWorkflows && !input.Force {
			// TODO: CRE-277 ensure forwarders are support the next DON version
			// https://github.com/smartcontractkit/chainlink/blob/4fc61bb156fe57bfd939b836c02c413ad1209ebb/contracts/src/v0.8/keystone/CapabilitiesRegistry.sol#L812
			// and
			// https://github.com/smartcontractkit/chainlink/blob/4fc61bb156fe57bfd939b836c02c413ad1209ebb/contracts/src/v0.8/keystone/KeystoneForwarder.sol#L274
			return UpdateDONOutput{}, fmt.Errorf("refusing to update workflow don %d at config version %d because we cannot validate that all forwarder contracts are ready to accept the new configure version", don.Id, don.ConfigCount)
		}

		cfgs, err := computeConfigs(input.CapabilityConfigs, don.CapabilityConfigurations, input.MergeCapabilityConfigsWithOnChain)
		if err != nil {
			return UpdateDONOutput{}, fmt.Errorf("failed to compute configs: %w", err)
		}

		f := input.F
		if f == 0 {
			f = don.F
		}
		// this is implement as such to maintain backwards compatibility; the default (omitted) value of a bool is false
		var isPublic bool
		if input.IsPrivate {
			isPublic = false
		} else {
			isPublic = don.IsPublic
		}

		name := don.Name
		if input.NewDonName != "" {
			name = input.NewDonName
		}

		donUpdate := capabilities_registry_v2.CapabilitiesRegistryUpdateDONParams{
			Name:                     name,
			Nodes:                    pkg.PeerIDsToBytes(input.P2PIDs),
			CapabilityConfigurations: cfgs,
			IsPublic:                 isPublic,
			F:                        f,
			Config:                   don.Config,
		}

		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return registry.UpdateDONByName(opts, input.DonName, donUpdate)
		})
		if err != nil {
			err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
			return UpdateDONOutput{}, fmt.Errorf("failed to execute UpdateDON: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for UpdateDON '%s' on chain %d", input.DonName, input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully updated DON '%s' on chain %d", input.DonName, input.ChainSelector)
		}

		return UpdateDONOutput{
			DonInfo:   donUpdate,
			Operation: operation,
		}, nil
	},
)

func computeConfigs(
	capCfgs []CapabilityConfig,
	existingCapConfigs []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration,
	mergeCapabilities bool) ([]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, error) {
	capSet := make(map[string]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration)
	for _, capCfg := range capCfgs {
		onChainCap, err := capabilityConfigToOnChain(capCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to convert capability config to on-chain format: %w", err)
		}

		_, ok := capSet[onChainCap.CapabilityId]
		if ok {
			return nil, fmt.Errorf("duplicate capability configuration for id: %s", onChainCap.CapabilityId)
		}

		capSet[onChainCap.CapabilityId] = *onChainCap
	}

	if mergeCapabilities {
		for _, existingCapConfig := range existingCapConfigs {
			_, ok := capSet[existingCapConfig.CapabilityId]
			if !ok {
				capSet[existingCapConfig.CapabilityId] = existingCapConfig
			}
		}
	}
	var out []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration
	for _, capCfg := range capSet {
		out = append(out, capCfg)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CapabilityId < out[j].CapabilityId
	})
	return out, nil
}

func capabilityConfigToOnChain(capCfg CapabilityConfig) (*capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, error) {
	cfg := capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration{}
	cfg.CapabilityId = capCfg.Capability.CapabilityID
	var err error
	x := pkg.CapabilityConfig(capCfg.Config)
	cfg.Config, err = x.MarshalProto()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capability configuration config: %w", err)
	}
	if cfg.Config == nil {
		return nil, fmt.Errorf("config is required for capability %s", capCfg.Capability.CapabilityID)
	}

	return &cfg, nil
}
