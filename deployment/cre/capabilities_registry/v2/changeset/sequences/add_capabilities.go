package sequences

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/modifier"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

type AddCapabilitiesDeps struct {
	Env           *cldf.Environment
	MCMSContracts *commonchangeset.MCMSWithTimelockState // Required if MCMSConfig is not nil
}

type AddCapabilitiesInput struct {
	// DonCapabilityConfigs maps DON name to the list of capability configs for that DON.
	DonCapabilityConfigs map[string][]contracts.CapabilityConfig

	// Force indicates whether to force the update even if we cannot validate that all forwarder contracts are ready to accept the new configure version.
	// This is very dangerous, and could break the whole platform if the forwarders are not ready. Be very careful with this option.
	Force bool

	RegistryRef datastore.AddressRefKey
	MCMSConfig  *crecontracts.MCMSConfig
}

func (i *AddCapabilitiesInput) Validate() error {
	if len(i.DonCapabilityConfigs) == 0 {
		return errors.New("donCapabilityConfigs must contain at least one DON entry")
	}
	for donName, configs := range i.DonCapabilityConfigs {
		if donName == "" {
			return errors.New("donCapabilityConfigs keys cannot be empty strings")
		}
		if len(configs) == 0 {
			return fmt.Errorf("donCapabilityConfigs[%q] must contain at least one capability config", donName)
		}
	}
	return nil
}

type AddCapabilitiesOutput struct {
	AddedCapabilities []contracts.RegisterableCapability
	DonInfos          []capabilities_registry_v2.CapabilitiesRegistryUpdateDONParams
	UpdatedNodes      []capabilities_registry_v2.CapabilitiesRegistryNodeParams
	Proposals         []mcmslib.TimelockProposal
}

var AddCapabilities = operations.NewSequence[AddCapabilitiesInput, AddCapabilitiesOutput, AddCapabilitiesDeps](
	"add-capabilities-seq",
	semver.MustParse("1.0.0"),
	"Add Capabilities to the capabilities registry",
	func(b operations.Bundle, deps AddCapabilitiesDeps, input AddCapabilitiesInput) (AddCapabilitiesOutput, error) {
		if err := input.Validate(); err != nil {
			return AddCapabilitiesOutput{}, fmt.Errorf("invalid input: %w", err)
		}

		chainSel := input.RegistryRef.ChainSelector()
		chain, ok := deps.Env.BlockChains.EVMChains()[chainSel]
		if !ok {
			return AddCapabilitiesOutput{}, fmt.Errorf("chain not found for selector %d", chainSel)
		}

		registryAddressRef, err := deps.Env.DataStore.Addresses().Get(input.RegistryRef)
		if err != nil {
			return AddCapabilitiesOutput{}, fmt.Errorf("failed to get registry address: %w", err)
		}

		capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
			common.HexToAddress(registryAddressRef.Address), chain.Client,
		)
		if err != nil {
			return AddCapabilitiesOutput{}, fmt.Errorf("failed to create CapabilitiesRegistry: %w", err)
		}

		// Build capabilities list once (registry-level; union across all DONs).
		capabilities, err := buildCapabilitiesFromAllDONConfigs(input.DonCapabilityConfigs)
		if err != nil {
			return AddCapabilitiesOutput{}, err
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			common.HexToAddress(registryAddressRef.Address),
			contracts.AddCapabilitiesDescription,
		)
		if err != nil {
			return AddCapabilitiesOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		// Register capabilities once for the registry.
		regCapsReport, err := operations.ExecuteOperation(
			b,
			contracts.RegisterCapabilities,
			contracts.RegisterCapabilitiesDeps{
				Env:      deps.Env,
				Strategy: strategy,
			},
			contracts.RegisterCapabilitiesInput{
				Address:       registryAddressRef.Address,
				ChainSelector: chainSel,
				Capabilities:  capabilities,
				MCMSConfig:    input.MCMSConfig,
			},
		)
		if err != nil {
			return AddCapabilitiesOutput{}, fmt.Errorf("failed to register capabilities: %w", err)
		}

		var allOps []types.BatchOperation
		if regCapsReport.Output.Operation != nil {
			allOps = append(allOps, *regCapsReport.Output.Operation)
		}

		var donInfos []capabilities_registry_v2.CapabilitiesRegistryUpdateDONParams
		var allUpdatedNodes []capabilities_registry_v2.CapabilitiesRegistryNodeParams

		// Update each DON: get nodes, update node configs, update DON.
		for donName, donCapConfigs := range input.DonCapabilityConfigs {
			don, nodes, err := GetDonNodes(donName, capReg)
			if err != nil {
				return AddCapabilitiesOutput{}, fmt.Errorf("failed to get DON %s nodes: %w", donName, err)
			}

			p2pIDs := make([]p2pkey.PeerID, 0, len(nodes))
			for _, node := range nodes {
				p2pIDs = append(p2pIDs, node.P2pId)
			}

			nodeUpdates, err := buildNodeUpdatesForDON(p2pIDs, donCapConfigs)
			if err != nil {
				return AddCapabilitiesOutput{}, fmt.Errorf("failed to build node updates for DON %s: %w", donName, err)
			}

			// apply modifiers to capability configs
			// currently we add p2pToTransmitterMap to the specConfig for Aptos capabilities
			// more modifiers can be added here as needed
			modifierParams := modifier.CapabilityConfigModifierParams{
				Env:     deps.Env,
				DonName: donName,
				P2PIDs:  p2pIDs,
				Configs: donCapConfigs, // modified in place
			}
			for _, mod := range modifier.DefaultCapabilityConfigModifiers() {
				if err := mod.Modify(modifierParams); err != nil {
					return AddCapabilitiesOutput{}, fmt.Errorf("modify capability configs for DON %s: %w", donName, err)
				}
			}

			updateNodesReport, err := operations.ExecuteOperation(
				b,
				contracts.UpdateNodes,
				contracts.UpdateNodesDeps{
					Env:                  deps.Env,
					CapabilitiesRegistry: capReg,
					Strategy:             strategy,
				},
				contracts.UpdateNodesInput{
					ChainSelector: chainSel,
					NodesUpdates:  nodeUpdates,
					MCMSConfig:    input.MCMSConfig,
				},
			)
			if err != nil {
				return AddCapabilitiesOutput{}, fmt.Errorf("failed to update nodes for DON %s: %w", donName, err)
			}

			updateDonReport, err := operations.ExecuteOperation(
				b,
				contracts.UpdateDON,
				contracts.UpdateDONDeps{
					Env:                  deps.Env,
					CapabilitiesRegistry: capReg,
					Strategy:             strategy,
				},
				contracts.UpdateDONInput{
					ChainSelector:                     chainSel,
					P2PIDs:                            p2pIDs,
					CapabilityConfigs:                 donCapConfigs,
					MergeCapabilityConfigsWithOnChain: true,
					DonName:                           donName,
					F:                                 don.F,
					IsPrivate:                         !don.IsPublic,
					Force:                             input.Force,
					MCMSConfig:                        input.MCMSConfig,
				},
			)
			if err != nil {
				return AddCapabilitiesOutput{}, fmt.Errorf("failed to update DON %s: %w", donName, err)
			}

			allOps = append(allOps, toOpsSlice(updateNodesReport.Output.Operation, updateDonReport.Output.Operation)...)
			donInfos = append(donInfos, updateDonReport.Output.DonInfo)
			allUpdatedNodes = append(allUpdatedNodes, updateNodesReport.Output.UpdatedNodes...)
		}

		var proposals []mcmslib.TimelockProposal
		if input.MCMSConfig != nil {
			if len(allOps) > 0 {
				proposal, mcmsErr := strategy.BuildProposal(allOps)
				if mcmsErr != nil {
					return AddCapabilitiesOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", mcmsErr)
				}
				proposals = append(proposals, *proposal)
			} else {
				deps.Env.Logger.Warnw("Add capability sequence has not produced any operations to execute")
			}
		}

		return AddCapabilitiesOutput{
			DonInfos:          donInfos,
			UpdatedNodes:      allUpdatedNodes,
			AddedCapabilities: regCapsReport.Output.Capabilities,
			Proposals:         proposals,
		}, nil
	},
)

func toOpsSlice(opPtrs ...*types.BatchOperation) []types.BatchOperation {
	var result []types.BatchOperation
	for _, opPtr := range opPtrs {
		if opPtr != nil {
			result = append(result, *opPtr)
		}
	}

	return result
}

// buildCapabilitiesFromAllDONConfigs collects the unique capabilities across all DONs' configs
// for registry-level registration.
func buildCapabilitiesFromAllDONConfigs(donConfigs map[string][]contracts.CapabilityConfig) ([]contracts.RegisterableCapability, error) {
	uniqueCaps := make(map[string]contracts.RegisterableCapability)
	for _, configs := range donConfigs {
		for _, cfg := range configs {
			if _, ok := uniqueCaps[cfg.Capability.CapabilityID]; ok {
				continue
			}
			uniqueCaps[cfg.Capability.CapabilityID] = contracts.RegisterableCapability{
				Metadata:              cfg.Capability.Metadata,
				CapabilityID:          cfg.Capability.CapabilityID,
				ConfigurationContract: cfg.Capability.ConfigurationContract,
			}
		}
	}
	return slices.Collect(maps.Values(uniqueCaps)), nil
}

// buildNodeUpdatesForDON builds node config updates for a DON's nodes (adds the new capabilities to each node).
func buildNodeUpdatesForDON(p2pIDs []p2pkey.PeerID, configs []contracts.CapabilityConfig) (map[string]contracts.NodeConfig, error) {
	nodeUpdates := make(map[string]contracts.NodeConfig, len(p2pIDs))
	for _, cfg := range configs {
		metadataBytes, err := json.Marshal(cfg.Capability.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal capability metadata for capability %s: %w", cfg.Capability.CapabilityID, err)
		}
		capability := capabilities_registry_v2.CapabilitiesRegistryCapability{
			CapabilityId:          cfg.Capability.CapabilityID,
			ConfigurationContract: cfg.Capability.ConfigurationContract,
			Metadata:              metadataBytes,
		}
		for _, p2pID := range p2pIDs {
			p2pIDStr := p2pID.String()
			nodeUpdate := nodeUpdates[p2pIDStr]
			if nodeUpdate.Capabilities == nil {
				nodeUpdate.Capabilities = make([]capabilities_registry_v2.CapabilitiesRegistryCapability, 0, len(configs))
			}
			nodeUpdate.Capabilities = append(nodeUpdate.Capabilities, capability)
			nodeUpdates[p2pIDStr] = nodeUpdate
		}
	}
	return nodeUpdates, nil
}

func GetDonNodes(donName string, capReg *capabilities_registry_v2.CapabilitiesRegistry) (
	*capabilities_registry_v2.CapabilitiesRegistryDONInfo,
	[]capabilities_registry_v2.INodeInfoProviderNodeInfo,
	error,
) {
	don, err := capReg.GetDONByName(&bind.CallOpts{}, donName)
	if err != nil {
		err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
		return nil, nil, fmt.Errorf("failed to get DON by name %s: %w", donName, err)
	}

	nodes, err := capReg.GetNodesByP2PIds(&bind.CallOpts{}, don.NodeP2PIds)
	if err != nil {
		err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
		return nil, nil, fmt.Errorf("failed to get nodes by P2P IDs for DON %s: %w", donName, err)
	}

	return &don, nodes, nil
}
