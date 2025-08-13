package changeset

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

type CapabilitiesRegistryNodeOperator struct {
	Admin common.Address `json:"admin" yaml:"admin"`
	Name  string         `json:"name" yaml:"name"`
}

func (nop CapabilitiesRegistryNodeOperator) ToWrapper() capabilities_registry_v2.CapabilitiesRegistryNodeOperator {
	return capabilities_registry_v2.CapabilitiesRegistryNodeOperator{
		Admin: nop.Admin,
		Name:  nop.Name,
	}
}

type CapabilityMetadata struct {
	CapabilityType uint8 `json:"capabilityType" yaml:"capabilityType"`
	ResponseType   uint8 `json:"responseType" yaml:"responseType"`
}

type CapabilitiesRegistryCapability struct {
	CapabilityID          string             `json:"capabilityID" yaml:"capabilityID"`
	ConfigurationContract common.Address     `json:"configurationContract" yaml:"configurationContract"`
	Metadata              CapabilityMetadata `json:"metadata" yaml:"metadata"`
}

func (cap CapabilitiesRegistryCapability) ToWrapper() (capabilities_registry_v2.CapabilitiesRegistryCapability, error) {
	metadataBytes, err := json.Marshal(cap.Metadata)
	if err != nil {
		return capabilities_registry_v2.CapabilitiesRegistryCapability{}, fmt.Errorf("failed to marshal metadata: %w", err)
	}
	return capabilities_registry_v2.CapabilitiesRegistryCapability{
		CapabilityId:          cap.CapabilityID,
		ConfigurationContract: cap.ConfigurationContract,
		Metadata:              metadataBytes,
	}, nil
}

type CapabilitiesRegistryNodeParams struct {
	NodeOperatorID      uint32   `json:"nodeOperatorID" yaml:"nodeOperatorID"`
	Signer              [32]byte `json:"signer" yaml:"signer"`
	P2pID               [32]byte `json:"p2pID" yaml:"p2pID"`
	EncryptionPublicKey [32]byte `json:"encryptionPublicKey" yaml:"encryptionPublicKey"`
	CsaKey              [32]byte `json:"csaKey" yaml:"csaKey"`
	CapabilityIDs       []string `json:"capabilityIDs" yaml:"capabilityIDs"`
}

func (node CapabilitiesRegistryNodeParams) ToWrapper() capabilities_registry_v2.CapabilitiesRegistryNodeParams {
	return capabilities_registry_v2.CapabilitiesRegistryNodeParams{
		NodeOperatorId:      node.NodeOperatorID,
		Signer:              node.Signer,
		P2pId:               node.P2pID,
		EncryptionPublicKey: node.EncryptionPublicKey,
		CsaKey:              node.CsaKey,
		CapabilityIds:       node.CapabilityIDs,
	}
}

type CapabilitiesRegistryCapabilityConfiguration struct {
	CapabilityID string `json:"capabilityID" yaml:"capabilityID"`
	Config       []byte `json:"config" yaml:"config"`
}

type CapabilitiesRegistryNewDONParams struct {
	Name                     string                                        `json:"name" yaml:"name"`
	DonFamilies              []string                                      `json:"donFamilies" yaml:"donFamilies"`
	Config                   []byte                                        `json:"config" yaml:"config"`
	CapabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration `json:"capabilityConfigurations" yaml:"capabilityConfigurations"`
	Nodes                    [][32]byte                                    `json:"nodes" yaml:"nodes"`
	F                        uint8                                         `json:"f" yaml:"f"`
	IsPublic                 bool                                          `json:"isPublic" yaml:"isPublic"`
	AcceptsWorkflows         bool                                          `json:"acceptsWorkflows" yaml:"acceptsWorkflows"`
}

func (don CapabilitiesRegistryNewDONParams) ToWrapper() capabilities_registry_v2.CapabilitiesRegistryNewDONParams {
	capabilityConfigurations := make([]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, len(don.CapabilityConfigurations))
	for j, capConfig := range don.CapabilityConfigurations {
		capabilityConfigurations[j] = capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration{
			CapabilityId: capConfig.CapabilityID,
			Config:       capConfig.Config,
		}
	}
	return capabilities_registry_v2.CapabilitiesRegistryNewDONParams{
		Name:                     don.Name,
		DonFamilies:              don.DonFamilies,
		Config:                   don.Config,
		CapabilityConfigurations: capabilityConfigurations,
		Nodes:                    don.Nodes,
		F:                        don.F,
		IsPublic:                 don.IsPublic,
		AcceptsWorkflows:         don.AcceptsWorkflows,
	}
}
