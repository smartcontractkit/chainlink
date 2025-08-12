package changeset

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
)

var _ cldf.ChangeSetV2[ConfigureCapabilitiesRegistryInput] = ConfigureCapabilitiesRegistry{}

// ConfigureCapabilitiesRegistryInput must be JSON Serializable with no private fields
type ConfigureCapabilitiesRegistryInput struct {
	ChainSelector               uint64                                                      `json:"chainSelector"`
	CapabilitiesRegistryAddress string                                                      `json:"capabilitiesRegistryAddress"`
	UseMCMS                     bool                                                        `json:"useMCMS"`
	Nops                        []capabilities_registry_v2.CapabilitiesRegistryNodeOperator `json:"nops,omitempty"`
	Capabilities                []capabilities_registry_v2.CapabilitiesRegistryCapability   `json:"capabilities,omitempty"`
	Nodes                       []capabilities_registry_v2.CapabilitiesRegistryNodeParams   `json:"nodes,omitempty"`
	DONs                        []capabilities_registry_v2.CapabilitiesRegistryNewDONParams `json:"dons,omitempty"`
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
	capabilitiesRegistryConfigurationReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		sequences.ConfigureCapabilitiesRegistry,
		sequences.ConfigureCapabilitiesRegistryDeps{Env: &e},
		sequences.ConfigureCapabilitiesRegistryInput{
			RegistryChainSel: config.ChainSelector,
			UseMCMS:          config.UseMCMS,
			ContractAddress:  config.CapabilitiesRegistryAddress,
			Nops:             config.Nops,
			Capabilities:     config.Capabilities,
			Nodes:            config.Nodes,
			DONs:             config.DONs,
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
