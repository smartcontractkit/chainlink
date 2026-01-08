package sequences

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type AddCapabilitiesDeps struct {
	Env           *cldf.Environment
	MCMSContracts *commonchangeset.MCMSWithTimelockState // Required if MCMSConfig is not nil
}

type AddCapabilitiesInput struct {
	CapabilityConfigs []contracts.CapabilityConfig // if Config subfield is nil, a default config is used

	// DonName to update, this is required
	DonName string

	// Force indicates whether to force the update even if we cannot validate that all forwarder contracts are ready to accept the new configure version.
	// This is very dangerous, and could break the whole platform if the forwarders are not ready. Be very careful with this option.
	Force bool

	RegistryRef datastore.AddressRefKey
	MCMSConfig  *crecontracts.MCMSConfig
}

func (i *AddCapabilitiesInput) Validate() error {
	if i.DonName == "" {
		return errors.New("must specify DONName")
	}
	if len(i.CapabilityConfigs) == 0 {
		return errors.New("capabilityConfigs is required")
	}
	return nil
}

type AddCapabilitiesOutput struct {
	AddedCapabilities []contracts.RegisterableCapability
	DonInfo           capabilities_registry_v2.CapabilitiesRegistryUpdateDONParams
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

		don, nodes, err := GetDonNodes(input.DonName, capReg)
		if err != nil {
			return AddCapabilitiesOutput{}, fmt.Errorf("failed to get DON %s nodes: %w", input.DonName, err)
		}

		p2pIDs := make([]p2pkey.PeerID, 0)
		for _, node := range nodes {
			p2pIDs = append(p2pIDs, node.P2pId)
		}

		nodeUpdates := make(map[string]contracts.NodeConfig, len(p2pIDs))
		capabilities := make([]contracts.RegisterableCapability, len(input.CapabilityConfigs))
		for i, cfg := range input.CapabilityConfigs {
			metadataBytes, err := json.Marshal(cfg.Capability.Metadata)
			if err != nil {
				return AddCapabilitiesOutput{}, fmt.Errorf("failed to marshal capability metadata for capability %s: %w", cfg.Capability.CapabilityID, err)
			}
			capability := capabilities_registry_v2.CapabilitiesRegistryCapability{
				CapabilityId:          cfg.Capability.CapabilityID,
				ConfigurationContract: cfg.Capability.ConfigurationContract,
				Metadata:              metadataBytes,
			}
			capabilities[i] = contracts.RegisterableCapability{
				Metadata:              cfg.Capability.Metadata,
				CapabilityID:          cfg.Capability.CapabilityID,
				ConfigurationContract: cfg.Capability.ConfigurationContract,
			}
			for _, p2pID := range p2pIDs {
				p2pIDStr := p2pID.String()
				nodeUpdate, exists := nodeUpdates[p2pIDStr]
				if !exists {
					nodeUpdate = contracts.NodeConfig{
						Capabilities: make([]capabilities_registry_v2.CapabilitiesRegistryCapability, 0, len(input.CapabilityConfigs)),
					}
				}
				nodeUpdate.Capabilities = append(nodeUpdates[p2pIDStr].Capabilities, capability)
				nodeUpdates[p2pIDStr] = nodeUpdate
			}
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
			return AddCapabilitiesOutput{}, fmt.Errorf("failed to update nodes: %w", err)
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
				CapabilityConfigs:                 input.CapabilityConfigs,
				MergeCapabilityConfigsWithOnChain: true,
				DonName:                           input.DonName,
				F:                                 don.F,
				IsPrivate:                         !don.IsPublic,
				Force:                             input.Force,
				MCMSConfig:                        input.MCMSConfig,
			},
		)
		if err != nil {
			return AddCapabilitiesOutput{}, fmt.Errorf("failed to update don: %w", err)
		}

		var proposals []mcmslib.TimelockProposal

		if input.MCMSConfig != nil {
			ops := toOpsSlice(regCapsReport.Output.Operation, updateNodesReport.Output.Operation, updateDonReport.Output.Operation)
			if len(ops) > 0 {
				proposal, mcmsErr := strategy.BuildProposal(ops)
				if mcmsErr != nil {
					return AddCapabilitiesOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", mcmsErr)
				}

				proposals = append(proposals, *proposal)
			} else {
				deps.Env.Logger.Warnw("Add capability sequence has not produced any operations to execute")
			}
		}

		return AddCapabilitiesOutput{
			DonInfo:           updateDonReport.Output.DonInfo,
			UpdatedNodes:      updateNodesReport.Output.UpdatedNodes,
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
