package sequences

import (
	"errors"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

type ConfigureCapabilitiesRegistryDeps struct {
	Env *cldf.Environment
}

type ConfigureCapabilitiesRegistryInput struct {
	RegistryChainSel uint64

	UseMCMS         bool
	ContractAddress string
	Nops            []capabilities_registry_v2.CapabilitiesRegistryNodeOperator
	Nodes           []capabilities_registry_v2.CapabilitiesRegistryNodeParams
	Capabilities    []capabilities_registry_v2.CapabilitiesRegistryCapability
	DONs            []capabilities_registry_v2.CapabilitiesRegistryNewDONParams
}

func (c ConfigureCapabilitiesRegistryInput) Validate() error {
	if c.ContractAddress == "" {
		return errors.New("ContractAddress is not set")
	}
	return nil
}

type ConfigureCapabilitiesRegistryOutput struct {
	Nops         []*capabilities_registry_v2.CapabilitiesRegistryNodeOperatorAdded
	Nodes        []*capabilities_registry_v2.CapabilitiesRegistryNodeAdded
	Capabilities []*capabilities_registry_v2.CapabilitiesRegistryCapabilityConfigured
	DONs         []capabilities_registry_v2.CapabilitiesRegistryDONInfo
}

var ConfigureCapabilitiesRegistry = operations.NewSequence(
	"configure-capabilities-registry",
	semver.MustParse("1.0.0"),
	"Configures the capabilities registry by registering node operators, nodes, dons and capabilities",
	func(b operations.Bundle, deps ConfigureCapabilitiesRegistryDeps, input ConfigureCapabilitiesRegistryInput) (ConfigureCapabilitiesRegistryOutput, error) {
		// Register Node Operators
		registerNopsReport, err := operations.ExecuteOperation(b, contracts.RegisterNops, contracts.RegisterNopsDeps{Env: deps.Env}, contracts.RegisterNopsInput{
			ChainSelector: input.RegistryChainSel,
			Address:       input.ContractAddress,
			Nops:          input.Nops,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistryOutput{}, err
		}

		// Register capabilities
		registerCapabilitiesReport, err := operations.ExecuteOperation(b, contracts.RegisterCapabilities, contracts.RegisterCapabilitiesDeps{Env: deps.Env}, contracts.RegisterCapabilitiesInput{
			ChainSelector: input.RegistryChainSel,
			Address:       input.ContractAddress,
			Capabilities:  input.Capabilities,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistryOutput{}, err
		}

		// Register Nodes
		registerNodesReport, err := operations.ExecuteOperation(b, contracts.RegisterNodes, contracts.RegisterNodesDeps{Env: deps.Env}, contracts.RegisterNodesInput{
			ChainSelector: input.RegistryChainSel,
			Address:       input.ContractAddress,
			Nodes:         input.Nodes,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistryOutput{}, err
		}

		// Register DONs
		registerDONsReport, err := operations.ExecuteOperation(b, contracts.RegisterDons, contracts.RegisterDonsDeps{Env: deps.Env}, contracts.RegisterDonsInput{
			ChainSelector: input.RegistryChainSel,
			Address:       input.ContractAddress,
			DONs:          input.DONs,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistryOutput{}, err
		}

		return ConfigureCapabilitiesRegistryOutput{
			Nops:         registerNopsReport.Output.Nops,
			Nodes:        registerNodesReport.Output.Nodes,
			Capabilities: registerCapabilitiesReport.Output.Capabilities,
			DONs:         registerDONsReport.Output.DONs,
		}, nil
	},
)
