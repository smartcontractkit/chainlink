package sequences

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

type UpdateDONDeps struct {
	Env *cldf.Environment
}

type UpdateDONInput struct {
	RegistryChainSel uint64
	// P2PIDs are the peer ids that compose the don
	P2PIDs            []p2pkey.PeerID
	CapabilityConfigs []contracts.CapabilityConfig // if Config subfield is nil, a default config is used

	// DonName to update, this is required
	DonName string

	// F is the fault tolerance level
	// if omitted, the existing value fetched from the registry is used
	F uint8

	// IsPrivate indicates whether the DON is public or private
	// If omitted, the existing value fetched from the registry is used
	IsPrivate bool

	// Force indicates whether to force the update even if we cannot validate that all forwarder contracts are ready to accept the new configure version.
	// This is very dangerous, and could break the whole platform if the forwarders are not ready. Be very careful with this option.
	Force bool

	RegistryRef datastore.AddressRefKey
}

func (i *UpdateDONInput) Validate() error {
	if i.DonName == "" {
		return errors.New("must specify DONName")
	}
	if len(i.CapabilityConfigs) == 0 {
		return errors.New("capabilityConfigs is required")
	}
	return nil
}

type UpdateDONOutput struct {
	DonInfo           capabilities_registry_v2.CapabilitiesRegistryDONInfo
	UpdatedNodes      []*capabilities_registry_v2.CapabilitiesRegistryNodeUpdated
	AddedCapabilities []*capabilities_registry_v2.CapabilitiesRegistryCapabilityConfigured
}

var UpdateDON = operations.NewSequence[UpdateDONInput, UpdateDONOutput, UpdateDONDeps](
	"update-don",
	semver.MustParse("1.0.0"),
	"Updates a DON in the capabilities registry",
	func(b operations.Bundle, deps UpdateDONDeps, input UpdateDONInput) (UpdateDONOutput, error) {
		if err := input.Validate(); err != nil {
			return UpdateDONOutput{}, fmt.Errorf("invalid input: %w", err)
		}

		chain, ok := deps.Env.BlockChains.EVMChains()[input.RegistryChainSel]
		if !ok {
			return UpdateDONOutput{}, fmt.Errorf("chain not found for selector %d", input.RegistryChainSel)
		}

		registryAddressRef, err := deps.Env.DataStore.Addresses().Get(input.RegistryRef)
		if err != nil {
			return UpdateDONOutput{}, fmt.Errorf("failed to get registry address: %w", err)
		}

		capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
			common.HexToAddress(registryAddressRef.Address), chain.Client,
		)
		if err != nil {
			return UpdateDONOutput{}, fmt.Errorf("failed to create CapabilitiesRegistry: %w", err)
		}

		nodes, err := getDonNodes(input.DonName, capReg)
		if err != nil {
			return UpdateDONOutput{}, fmt.Errorf("failed to get DON %s nodes: %w", input.DonName, err)
		}

		p2pIDs := input.P2PIDs
		if len(p2pIDs) == 0 {
			// If no P2P IDs are provided to change the DON composition, we use the existing DON's P2P IDs
			p2pIDs = make([]p2pkey.PeerID, 0)

			for _, node := range nodes {
				p2pIDs = append(p2pIDs, node.P2pId)
			}
		}

		nodeUpdates := make(map[p2pkey.PeerID]contracts.NodeConfig, len(p2pIDs))
		capabilities := make([]capabilities_registry_v2.CapabilitiesRegistryCapability, len(input.CapabilityConfigs))
		for i, cfg := range input.CapabilityConfigs {
			capabilities[i] = cfg.Capability
			for _, p2pID := range p2pIDs {
				nodeUpdate, exists := nodeUpdates[p2pID]
				if !exists {
					nodeUpdate = contracts.NodeConfig{
						Capabilities: make([]capabilities_registry_v2.CapabilitiesRegistryCapability, 0, len(input.CapabilityConfigs)),
					}
				}
				nodeUpdate.Capabilities = append(nodeUpdates[p2pID].Capabilities, cfg.Capability)
				nodeUpdates[p2pID] = nodeUpdate
			}
		}

		regCapsReport, err := operations.ExecuteOperation(
			b,
			contracts.RegisterCapabilities,
			contracts.RegisterCapabilitiesDeps(deps),
			contracts.RegisterCapabilitiesInput{
				Address:       registryAddressRef.Address,
				ChainSelector: input.RegistryChainSel,
				Capabilities:  capabilities,
			},
		)
		if err != nil {
			return UpdateDONOutput{}, fmt.Errorf("failed to register capabilities: %w", err)
		}

		updateNodesReport, err := operations.ExecuteOperation(
			b,
			contracts.UpdateNodes,
			contracts.UpdateNodesDeps{
				Env:                  deps.Env,
				CapabilitiesRegistry: capReg,
			},
			contracts.UpdateNodesInput{
				ChainSelector: input.RegistryChainSel,
				NodesUpdates:  nodeUpdates,
			},
		)
		if err != nil {
			return UpdateDONOutput{}, fmt.Errorf("failed to update nodes: %w", err)
		}

		updateDonReport, err := operations.ExecuteOperation(
			b,
			contracts.UpdateDON,
			contracts.UpdateDONDeps{
				Env:                  deps.Env,
				CapabilitiesRegistry: capReg,
			},
			contracts.UpdateDONInput{
				ChainSelector:     input.RegistryChainSel,
				P2PIDs:            p2pIDs,
				CapabilityConfigs: input.CapabilityConfigs,
				DonName:           input.DonName,
				F:                 input.F,
				IsPrivate:         input.IsPrivate,
				Force:             input.Force,
			},
		)
		if err != nil {
			return UpdateDONOutput{}, fmt.Errorf("failed to update don: %w", err)
		}

		return UpdateDONOutput{
			DonInfo:           updateDonReport.Output.DonInfo,
			UpdatedNodes:      updateNodesReport.Output.UpdatedNodes,
			AddedCapabilities: regCapsReport.Output.Capabilities,
		}, nil
	},
)

func getDonNodes(donName string, capReg *capabilities_registry_v2.CapabilitiesRegistry) (
	[]capabilities_registry_v2.INodeInfoProviderNodeInfo,
	error,
) {
	don, err := capReg.GetDONByName(&bind.CallOpts{}, donName)
	if err != nil {
		err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to get DON by name %s: %w", donName, err)
	}

	nodes, err := capReg.GetNodesByP2PIds(&bind.CallOpts{}, don.NodeP2PIds)
	if err != nil {
		err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to get nodes by P2P IDs for DON %s: %w", donName, err)
	}

	return nodes, nil
}
