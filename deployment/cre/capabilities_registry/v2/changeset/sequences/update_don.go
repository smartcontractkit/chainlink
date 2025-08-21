package sequences

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
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

	// DonName to update
	// If omitted, the don will be inferred from the P2P keys
	// If the update request intended to change the nodes in the don, the DonName must be specified
	DonName string

	// F is the fault tolerance level
	// if omitted, the existing value fetched from the registry is used
	F uint8

	// IsPrivate indicates whether the DON is public or private
	// If omitted, the existing value fetched from the registry is used
	IsPrivate bool

	RegistryRef datastore.AddressRefKey
}

func (i *UpdateDONInput) Validate() error {
	if len(i.P2PIDs) == 0 {
		return errors.New("p2pIDs is required")
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

		nodeUpdates := make(map[p2pkey.PeerID]contracts.UpdateNodesNodeUpdate, len(input.P2PIDs))
		capabilities := make([]capabilities_registry_v2.CapabilitiesRegistryCapability, len(input.CapabilityConfigs))
		for i, cfg := range input.CapabilityConfigs {
			capabilities[i] = cfg.Capability
			for _, p2pID := range input.P2PIDs {
				nodeUpdate, exists := nodeUpdates[p2pID]
				if !exists {
					nodeUpdate = contracts.UpdateNodesNodeUpdate{
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

		capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
			common.HexToAddress(registryAddressRef.Address), chain.Client,
		)

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
				P2PIDs:            input.P2PIDs,
				CapabilityConfigs: input.CapabilityConfigs,
				DonName:           input.DonName,
				F:                 input.F,
				IsPrivate:         input.IsPrivate,
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
