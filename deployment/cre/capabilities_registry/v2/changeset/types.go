package changeset

import (
	"encoding/hex"
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

type CapabilitiesRegistryCapability struct {
	CapabilityID          string                 `json:"capabilityID" yaml:"capabilityID"`
	ConfigurationContract common.Address         `json:"configurationContract" yaml:"configurationContract"`
	Metadata              map[string]interface{} `json:"metadata" yaml:"metadata"`
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
	Signer              string   `json:"signer" yaml:"signer"`
	P2pID               string   `json:"p2pID" yaml:"p2pID"`
	EncryptionPublicKey string   `json:"encryptionPublicKey" yaml:"encryptionPublicKey"`
	CsaKey              string   `json:"csaKey" yaml:"csaKey"`
	CapabilityIDs       []string `json:"capabilityIDs" yaml:"capabilityIDs"`
}

func (node CapabilitiesRegistryNodeParams) ToWrapper() (capabilities_registry_v2.CapabilitiesRegistryNodeParams, error) {
	csaKeyBytes, err := hexStringTo32Bytes(node.CsaKey)
	if err != nil {
		return capabilities_registry_v2.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert CSA key: %w", err)
	}

	signerBytes, err := hexStringTo32Bytes(node.Signer)
	if err != nil {
		return capabilities_registry_v2.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert signer: %w", err)
	}

	p2pIDBytes, err := hexStringTo32Bytes(node.P2pID)
	if err != nil {
		return capabilities_registry_v2.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert P2P ID: %w", err)
	}
	encryptionPublicKeyBytes, err := hexStringTo32Bytes(node.EncryptionPublicKey)
	if err != nil {
		return capabilities_registry_v2.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert encryption public key: %w", err)
	}

	return capabilities_registry_v2.CapabilitiesRegistryNodeParams{
		NodeOperatorId:      node.NodeOperatorID,
		Signer:              signerBytes,
		P2pId:               p2pIDBytes,
		EncryptionPublicKey: encryptionPublicKeyBytes,
		CsaKey:              csaKeyBytes,
		CapabilityIds:       node.CapabilityIDs,
	}, nil
}

type CapabilitiesRegistryCapabilityConfiguration struct {
	CapabilityID string                 `json:"capabilityID" yaml:"capabilityID"`
	Config       map[string]interface{} `json:"config" yaml:"config"`
}

type CapabilitiesRegistryNewDONParams struct {
	Name                     string                                        `json:"name" yaml:"name"`
	DonFamilies              []string                                      `json:"donFamilies" yaml:"donFamilies"`
	Config                   map[string]interface{}                        `json:"config" yaml:"config"`
	CapabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration `json:"capabilityConfigurations" yaml:"capabilityConfigurations"`
	Nodes                    []string                                      `json:"nodes" yaml:"nodes"`
	F                        uint8                                         `json:"f" yaml:"f"`
	IsPublic                 bool                                          `json:"isPublic" yaml:"isPublic"`
	AcceptsWorkflows         bool                                          `json:"acceptsWorkflows" yaml:"acceptsWorkflows"`
}

func (don CapabilitiesRegistryNewDONParams) ToWrapper() (capabilities_registry_v2.CapabilitiesRegistryNewDONParams, error) {
	capabilityConfigurations := make([]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, len(don.CapabilityConfigurations))
	for j, capConfig := range don.CapabilityConfigurations {
		configBytes, err := json.Marshal(capConfig.Config)
		if err != nil {
			return capabilities_registry_v2.CapabilitiesRegistryNewDONParams{}, fmt.Errorf("failed to marshal capability configuration config: %w", err)
		}
		capabilityConfigurations[j] = capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration{
			CapabilityId: capConfig.CapabilityID,
			Config:       configBytes,
		}
	}

	nodes := make([][32]byte, len(don.Nodes))
	for i, node := range don.Nodes {
		n, err := hexStringTo32Bytes(node)
		if err != nil {
			return capabilities_registry_v2.CapabilitiesRegistryNewDONParams{}, fmt.Errorf("failed to convert node ID: %w", err)
		}
		nodes[i] = n
	}

	configBytes, err := json.Marshal(don.Config)
	if err != nil {
		return capabilities_registry_v2.CapabilitiesRegistryNewDONParams{}, fmt.Errorf("failed to marshal DON config: %w", err)
	}

	return capabilities_registry_v2.CapabilitiesRegistryNewDONParams{
		Name:                     don.Name,
		DonFamilies:              don.DonFamilies,
		Config:                   configBytes,
		CapabilityConfigurations: capabilityConfigurations,
		Nodes:                    nodes,
		F:                        don.F,
		IsPublic:                 don.IsPublic,
		AcceptsWorkflows:         don.AcceptsWorkflows,
	}, nil
}

// hexStringTo32Bytes converts a hex string (with or without 0x prefix) to [32]byte
func hexStringTo32Bytes(hexStr string) ([32]byte, error) {
	var result [32]byte

	// Remove 0x prefix if present
	if len(hexStr) >= 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}

	// Validate length
	if len(hexStr) != 64 {
		return result, fmt.Errorf("invalid hex string length: expected 64 hex characters, got %d", len(hexStr))
	}

	// Decode hex string
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return result, fmt.Errorf("invalid hex string: %w", err)
	}

	// Copy to fixed-size array
	copy(result[:], bytes)
	return result, nil
}
