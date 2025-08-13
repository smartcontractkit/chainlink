package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

var _ cldf.ChangeSetV2[ConfigureCapabilitiesRegistryInput] = ConfigureCapabilitiesRegistry{}

// ConfigureCapabilitiesRegistryInput must be JSON and YAML Serializable with no private fields
type ConfigureCapabilitiesRegistryInput struct {
	ChainSelector               uint64                             `json:"chainSelector" yaml:"chainSelector"`
	CapabilitiesRegistryAddress string                             `json:"capabilitiesRegistryAddress" yaml:"capabilitiesRegistryAddress"`
	UseMCMS                     bool                               `json:"useMCMS" yaml:"useMCMS"`
	Nops                        []CapabilitiesRegistryNodeOperator `json:"nops,omitempty" yaml:"nops,omitempty"`
	Capabilities                []CapabilitiesRegistryCapability   `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
	Nodes                       []CapabilitiesRegistryNodeParams   `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	DONs                        []CapabilitiesRegistryNewDONParams `json:"dons,omitempty" yaml:"dons,omitempty"`
}

type ConfigureCapabilitiesRegistryDeps struct {
	Env *cldf.Environment
}

type ConfigureCapabilitiesRegistry struct{}

func (l ConfigureCapabilitiesRegistry) VerifyPreconditions(e cldf.Environment, config ConfigureCapabilitiesRegistryInput) error {
	if config.CapabilitiesRegistryAddress == "" {
		return errors.New("capabilitiesRegistryAddress is not set")
	}
	if _, ok := e.BlockChains.EVMChains()[config.ChainSelector]; !ok {
		return fmt.Errorf("chain %d not found in environment", config.ChainSelector)
	}

	return nil
}

func (l ConfigureCapabilitiesRegistry) Apply(e cldf.Environment, config ConfigureCapabilitiesRegistryInput) (cldf.ChangesetOutput, error) {
	nops := make([]capabilities_registry_v2.CapabilitiesRegistryNodeOperator, len(config.Nops))
	for i, nop := range config.Nops {
		nops[i] = nop.ToWrapper()
	}

	capabilities := make([]capabilities_registry_v2.CapabilitiesRegistryCapability, len(config.Capabilities))
	for i, cap := range config.Capabilities {
		capabilities[i] = cap.ToWrapper()
	}

	nodes := make([]capabilities_registry_v2.CapabilitiesRegistryNodeParams, len(config.Nodes))
	for i, node := range config.Nodes {
		nodes[i] = node.ToWrapper()
	}

	dons := make([]capabilities_registry_v2.CapabilitiesRegistryNewDONParams, len(config.DONs))
	for i, don := range config.DONs {
		dons[i] = don.ToWrapper()
	}

	capabilitiesRegistryConfigurationReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		sequences.ConfigureCapabilitiesRegistry,
		sequences.ConfigureCapabilitiesRegistryDeps{Env: &e},
		sequences.ConfigureCapabilitiesRegistryInput{
			RegistryChainSel: config.ChainSelector,
			UseMCMS:          config.UseMCMS,
			ContractAddress:  config.CapabilitiesRegistryAddress,
			Nops:             nops,
			Capabilities:     capabilities,
			Nodes:            nodes,
			DONs:             dons,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure capabilities registry: %w", err)
	}

	reports := make([]operations.Report[any, any], 0)
	reports = append(reports, capabilitiesRegistryConfigurationReport.ToGenericReport())

	return cldf.ChangesetOutput{
		Reports: reports,
	}, nil
}

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
	CapabilityID          string         `json:"capabilityID" yaml:"capabilityID"`
	ConfigurationContract common.Address `json:"configurationContract" yaml:"configurationContract"`
	Metadata              []byte         `json:"metadata" yaml:"metadata"`
}

func (cap CapabilitiesRegistryCapability) ToWrapper() capabilities_registry_v2.CapabilitiesRegistryCapability {
	return capabilities_registry_v2.CapabilitiesRegistryCapability{
		CapabilityId:          cap.CapabilityID,
		ConfigurationContract: cap.ConfigurationContract,
		Metadata:              cap.Metadata,
	}
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
