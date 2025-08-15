package changeset

import (
	"errors"
	"fmt"

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
		c, err := cap.ToWrapper()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert capability %d: %w", i, err)
		}
		capabilities[i] = c
	}

	nodes := make([]capabilities_registry_v2.CapabilitiesRegistryNodeParams, len(config.Nodes))
	for i, node := range config.Nodes {
		n, err := node.ToWrapper()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert node %d: %w", i, err)
		}
		nodes[i] = n
	}

	dons := make([]capabilities_registry_v2.CapabilitiesRegistryNewDONParams, len(config.DONs))
	for i, don := range config.DONs {
		d, err := don.ToWrapper()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert DON %d: %w", i, err)
		}
		dons[i] = d
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
