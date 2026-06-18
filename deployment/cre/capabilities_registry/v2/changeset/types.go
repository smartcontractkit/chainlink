package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/modifier"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
)

type CapabilitiesRegistryNodeOperator struct {
	Admin common.Address `json:"admin" yaml:"admin"`
	Name  string         `json:"name" yaml:"name"`
}

func (nop CapabilitiesRegistryNodeOperator) ToWrapper() capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams {
	return capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams{
		Admin: nop.Admin,
		Name:  nop.Name,
	}
}

type CapabilitiesRegistryCapability struct {
	CapabilityID          string         `json:"capabilityID" yaml:"capabilityID"`
	ConfigurationContract common.Address `json:"configurationContract" yaml:"configurationContract"`
	Metadata              map[string]any `json:"metadata" yaml:"metadata"`
}

func (cap CapabilitiesRegistryCapability) ToWrapper() (contracts.RegisterableCapability, error) {
	return contracts.RegisterableCapability{
		CapabilityID:          cap.CapabilityID,
		ConfigurationContract: cap.ConfigurationContract,
		Metadata:              cap.Metadata,
	}, nil
}

type CapabilitiesRegistryNodeParams struct {
	NOP                 string   `json:"nop" yaml:"nop"`
	Signer              string   `json:"signer" yaml:"signer"`
	P2pID               string   `json:"p2pID" yaml:"p2pID"`
	EncryptionPublicKey string   `json:"encryptionPublicKey" yaml:"encryptionPublicKey"`
	CsaKey              string   `json:"csaKey" yaml:"csaKey"`
	CapabilityIDs       []string `json:"capabilityIDs" yaml:"capabilityIDs"`
}

func (node CapabilitiesRegistryNodeParams) ToWrapper() (contracts.NodesInput, error) {
	csaKeyBytes, err := pkg.HexStringTo32Bytes(node.CsaKey)
	if err != nil {
		return contracts.NodesInput{}, fmt.Errorf("failed to convert CSA key: %w", err)
	}

	signerBytes, err := pkg.HexStringTo32Bytes(node.Signer)
	if err != nil {
		return contracts.NodesInput{}, fmt.Errorf("failed to convert signer: %w", err)
	}

	// P2PID is not a hex value
	p2pIDBytes, err := p2pkey.MakePeerID(node.P2pID)
	if err != nil {
		return contracts.NodesInput{}, fmt.Errorf("failed to convert P2P ID: %w", err)
	}

	encryptionPublicKeyBytes, err := pkg.HexStringTo32Bytes(node.EncryptionPublicKey)
	if err != nil {
		return contracts.NodesInput{}, fmt.Errorf("failed to convert encryption public key: %w", err)
	}

	if node.NOP == "" {
		return contracts.NodesInput{}, errors.New("NOP name cannot be empty")
	}

	return contracts.NodesInput{
		NOP:                 node.NOP,
		Signer:              signerBytes,
		P2pID:               p2pIDBytes,
		EncryptionPublicKey: encryptionPublicKeyBytes,
		CsaKey:              csaKeyBytes,
		CapabilityIDs:       node.CapabilityIDs,
	}, nil
}

type CapabilitiesRegistryCapabilityConfiguration struct {
	CapabilityID string         `json:"capabilityID" yaml:"capabilityID"`
	Config       map[string]any `json:"config" yaml:"config"`
}

type CapabilitiesRegistryNewDONParams struct {
	Name                     string                                        `json:"name" yaml:"name"`
	DonFamilies              []string                                      `json:"donFamilies" yaml:"donFamilies"`
	Config                   map[string]any                                `json:"config" yaml:"config"`
	CapabilityConfigurations []CapabilitiesRegistryCapabilityConfiguration `json:"capabilityConfigurations" yaml:"capabilityConfigurations"`
	Nodes                    []string                                      `json:"nodes" yaml:"nodes"`
	F                        uint8                                         `json:"f" yaml:"f"`
	IsPublic                 bool                                          `json:"isPublic" yaml:"isPublic"`
	AcceptsWorkflows         bool                                          `json:"acceptsWorkflows" yaml:"acceptsWorkflows"`
}

func (don CapabilitiesRegistryNewDONParams) ToWrapper(e cldf.Environment) (capabilities_registry_v2.CapabilitiesRegistryNewDONParams, error) {
	p2pIDs := make([]p2pkey.PeerID, 0)
	nodes := make([][32]byte, len(don.Nodes))
	// These are P2P IDs, they are not hex values
	for i, node := range don.Nodes {
		n, err := p2pkey.MakePeerID(node)
		if err != nil {
			return capabilities_registry_v2.CapabilitiesRegistryNewDONParams{}, fmt.Errorf("failed to convert node ID: %w", err)
		}
		nodes[i] = n
		p2pIDs = append(p2pIDs, n)
	}

	donCapCfg := pkg.CapabilityConfig(don.Config)
	donConfigBytes, err := donCapCfg.MarshalProto()
	if err != nil {
		return capabilities_registry_v2.CapabilitiesRegistryNewDONParams{}, fmt.Errorf("failed to marshal DON config: %w", err)
	}

	capabilityConfigurations, err := don.applyModifiersToCapabilityConfigs(e, p2pIDs)
	if err != nil {
		return capabilities_registry_v2.CapabilitiesRegistryNewDONParams{}, fmt.Errorf("failed to apply modifiers to capability configs: %w", err)
	}

	return capabilities_registry_v2.CapabilitiesRegistryNewDONParams{
		Name:                     don.Name,
		DonFamilies:              don.DonFamilies,
		Config:                   donConfigBytes,
		CapabilityConfigurations: capabilityConfigurations,
		Nodes:                    nodes,
		F:                        don.F,
		IsPublic:                 don.IsPublic,
		AcceptsWorkflows:         don.AcceptsWorkflows,
	}, nil
}

func (don CapabilitiesRegistryNewDONParams) applyModifiersToCapabilityConfigs(e cldf.Environment, p2pIDs []p2pkey.PeerID) ([]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, error) {
	capabilityConfigurations := make([]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, len(don.CapabilityConfigurations))

	capConfigsClone := make([]contracts.CapabilityConfig, len(don.CapabilityConfigurations))
	for i, capConfig := range don.CapabilityConfigurations {
		capConfigsClone[i] = contracts.CapabilityConfig{
			Capability: contracts.Capability{
				CapabilityID: capConfig.CapabilityID,
			},
			Config: capConfig.Config,
		}
	}

	// apply modifiers to capability configs
	// currently we add p2pToTransmitterMap to the specConfig for Aptos capabilities
	// more modifiers can be added here as needed
	modifierParams := modifier.CapabilityConfigModifierParams{
		Env:     &e,
		DonName: don.Name,
		P2PIDs:  p2pIDs,
		Configs: capConfigsClone, // modified in place
	}
	for _, mod := range modifier.DefaultCapabilityConfigModifiers() {
		if err := mod.Modify(modifierParams); err != nil {
			return nil, fmt.Errorf("modify capability configs for DON %s: %w", don.Name, err)
		}
	}

	for j, capConfig := range capConfigsClone {
		x := pkg.CapabilityConfig(capConfig.Config)
		capCfgBytes, err := x.MarshalProto()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal capability configuration config: %w", err)
		}
		capabilityConfigurations[j] = capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration{
			CapabilityId: capConfig.Capability.CapabilityID,
			Config:       capCfgBytes,
		}
	}

	return capabilityConfigurations, nil
}
