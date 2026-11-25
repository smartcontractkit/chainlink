package contracts

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

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

	// RemoveCapabilityIDs names the on-chain capabilities of the DON that this update is allowed to remove.
	// The registry replaces the whole capability set on every update, so any on-chain capability that is
	// missing from the new set and not named here makes the operation fail instead of silently dropping it.
	RemoveCapabilityIDs []string

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
	Capability Capability `json:"capability" yaml:"capability"`
	// Config is the capability configuration. It will be marshalled to proto config.
	// It is untyped here because is has to be deserialized from JSON/YAML for any possible capability
	// If nil, a default config based on the capability type is used
	Config map[string]any `json:"config" yaml:"config"`
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

		registered, err := ReadRegisteredCapabilities(registry)
		if err != nil {
			return UpdateDONOutput{}, err
		}

		cfgs, err := computeConfigs(input.CapabilityConfigs, don.CapabilityConfigurations, input.MergeCapabilityConfigsWithOnChain, input.RemoveCapabilityIDs, registered)
		if err != nil {
			return UpdateDONOutput{}, fmt.Errorf("failed to compute configs: %w", err)
		}
		// A named removal that is not on the DON is accepted so the same input can be run again. Say so,
		// because it is also what a removal aimed at the wrong capability id looks like.
		for _, id := range input.RemoveCapabilityIDs {
			if !slices.ContainsFunc(don.CapabilityConfigurations, func(c capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration) bool {
				return c.CapabilityId == id
			}) {
				deps.Env.Logger.Warnf("removeCapabilityIDs names capability %q which is not on DON '%s': nothing to remove", id, input.DonName)
			}
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

// RegisteredCapabilities maps the id of every capability registered in the registry to whether it is
// deprecated.
type RegisteredCapabilities map[string]bool

// ReadRegisteredCapabilities reads every capability registered in the registry.
func ReadRegisteredCapabilities(registry *capabilities_registry_v2.CapabilitiesRegistry) (RegisteredCapabilities, error) {
	caps, err := pkg.GetCapabilities(&bind.CallOpts{}, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to list the capabilities registered in the registry: %w", err)
	}
	out := make(RegisteredCapabilities, len(caps))
	for _, c := range caps {
		out[c.CapabilityId] = c.IsDeprecated
	}
	return out, nil
}

// CheckCapabilityUpdate validates a DON capability update against the DON's on-chain capabilities
// (existingCapConfigs) and the capabilities registered in the registry. It only looks at capability ids, and
// it is the check the operation itself runs, so preconditions and apply cannot disagree. It fails when:
//   - an id is listed twice, or is both in capCfgs and in removeIDs
//   - an id in removeIDs is neither on the DON nor registered (a typo). An id that is registered but no
//     longer on the DON is accepted: the removal already happened, so the same input can be run again
//   - the DON holds a deprecated capability that is not named in removeIDs (the contract rejects an update
//     that keeps one, and dropping one must be named like any other removal)
//   - with merge off, an on-chain capability is neither in capCfgs nor named in removeIDs
//
// It returns the on-chain capabilities that the update removes.
func CheckCapabilityUpdate(
	capCfgs []CapabilityConfig,
	existingCapConfigs []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration,
	mergeCapabilities bool,
	removeIDs []string,
	registered RegisteredCapabilities) (map[string]struct{}, error) {
	onChainIDs := make(map[string]struct{}, len(existingCapConfigs))
	for _, existingCapConfig := range existingCapConfigs {
		onChainIDs[existingCapConfig.CapabilityId] = struct{}{}
	}

	var deprecated []string
	newIDs := make(map[string]struct{}, len(capCfgs))
	for _, capCfg := range capCfgs {
		id := capCfg.Capability.CapabilityID
		if _, ok := newIDs[id]; ok {
			return nil, fmt.Errorf("duplicate capability configuration for id: %s", id)
		}
		// An id the registry does not know yet is not an error here: with MCMS the capability can be
		// registered by an earlier operation of the same proposal, which has not executed at this point.
		if registered[id] {
			deprecated = append(deprecated, id)
		}
		newIDs[id] = struct{}{}
	}

	removals := make(map[string]struct{}, len(removeIDs))
	for _, id := range removeIDs {
		if _, ok := newIDs[id]; ok {
			return nil, fmt.Errorf("capability %q is both in the capability configs and in removeCapabilityIDs", id)
		}
		if _, ok := onChainIDs[id]; ok {
			removals[id] = struct{}{}
		} else if _, known := registered[id]; !known {
			return nil, fmt.Errorf("removeCapabilityIDs names capability %q which is not on the DON and is not registered in the registry; check the id", id)
		}
	}

	var unexpected []string
	for id := range onChainIDs {
		_, inInput := newIDs[id]
		_, removed := removals[id]
		switch {
		case inInput || removed:
		case registered[id]:
			deprecated = append(deprecated, id)
		case !mergeCapabilities:
			unexpected = append(unexpected, id)
		}
	}
	if len(deprecated) > 0 {
		sort.Strings(deprecated)
		return nil, fmt.Errorf(
			"capabilities %s are deprecated in the registry and can only be removed from the DON; "+
				"name them in removeCapabilityIDs of an update_don input",
			strings.Join(deprecated, ", "))
	}
	if len(unexpected) > 0 {
		sort.Strings(unexpected)
		return nil, fmt.Errorf(
			"update would remove capabilities from the DON that are not listed for removal: %s. "+
				"To keep them, set mergeCapabilityConfigsWithOnChain to true or include them in the capability configs; "+
				"to remove them, name them in removeCapabilityIDs",
			strings.Join(unexpected, ", "))
	}

	return removals, nil
}

// computeConfigs builds the new capability set of a DON. The registry replaces the whole set on every
// update, so the set is only built once CheckCapabilityUpdate has accepted it.
func computeConfigs(
	capCfgs []CapabilityConfig,
	existingCapConfigs []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration,
	mergeCapabilities bool,
	removeIDs []string,
	registered RegisteredCapabilities) ([]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, error) {
	removals, err := CheckCapabilityUpdate(capCfgs, existingCapConfigs, mergeCapabilities, removeIDs, registered)
	if err != nil {
		return nil, err
	}

	capSet := make(map[string]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration)
	for _, capCfg := range capCfgs {
		onChainCap, err := capabilityConfigToOnChain(capCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to convert capability config to on-chain format: %w", err)
		}
		capSet[onChainCap.CapabilityId] = *onChainCap
	}

	if mergeCapabilities {
		for _, existingCapConfig := range existingCapConfigs {
			if _, removed := removals[existingCapConfig.CapabilityId]; removed {
				continue
			}
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
