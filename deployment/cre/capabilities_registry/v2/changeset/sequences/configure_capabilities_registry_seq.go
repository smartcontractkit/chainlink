package sequences

import (
	"errors"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

type ConfigureCapabilitiesRegistrySeqDeps struct {
	Env *cldf.Environment
}

type ConfigureCapabilitiesRegistrySeqInput struct {
	RegistryChainSel uint64

	UseMCMS         bool
	ContractAddress string
	Nops            []capabilities_registry_v2.CapabilitiesRegistryNodeOperator
	Nodes           []capabilities_registry_v2.CapabilitiesRegistryNodeParams
	Capabilities    []capabilities_registry_v2.CapabilitiesRegistryCapability
	DONs            []capabilities_registry_v2.CapabilitiesRegistryNewDONParams
}

func (c ConfigureCapabilitiesRegistrySeqInput) Validate() error {
	if c.ContractAddress == "" {
		return errors.New("ContractAddress is not set")
	}
	return nil
}

type ConfigureCapabilitiesRegistrySeqOutput struct {
	Nops         []*capabilities_registry_v2.CapabilitiesRegistryNodeOperatorAdded
	Nodes        []*capabilities_registry_v2.CapabilitiesRegistryNodeAdded
	Capabilities []*capabilities_registry_v2.CapabilitiesRegistryCapabilityConfigured
	DONs         []capabilities_registry_v2.CapabilitiesRegistryDONInfo
}

var ConfigureCapabilitiesRegistrySequence = operations.NewSequence(
	"configure-capabilities-registry",
	semver.MustParse("1.0.0"),
	"Configures the capabilities registry by registering node operators, nodes, dons and capabilities",
	func(b operations.Bundle, deps ConfigureCapabilitiesRegistrySeqDeps, input ConfigureCapabilitiesRegistrySeqInput) (ConfigureCapabilitiesRegistrySeqOutput, error) {
		// Register Node Operators
		registerNopsReport, err := operations.ExecuteOperation(b, contracts.RegisterNopsOp, contracts.RegisterNopsOpDeps{Env: deps.Env}, contracts.RegisterNopsOpInput{
			ChainSelector: input.RegistryChainSel,
			Address:       input.ContractAddress,
			Nops:          input.Nops,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, err
		}

		// Register capabilities
		registerCapabilitiesReport, err := operations.ExecuteOperation(b, contracts.RegisterCapabilitiesOp, contracts.RegisterCapabilitiesOpDeps{Env: deps.Env}, contracts.RegisterCapabilitiesOpInput{
			ChainSelector: input.RegistryChainSel,
			Address:       input.ContractAddress,
			Capabilities:  input.Capabilities,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, err
		}

		// Register Nodes
		registerNodesReport, err := operations.ExecuteOperation(b, contracts.RegisterNodesOp, contracts.RegisterNodesOpDeps{Env: deps.Env}, contracts.RegisterNodesOpInput{
			ChainSelector: input.RegistryChainSel,
			Address:       input.ContractAddress,
			Nodes:         input.Nodes,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, err
		}

		// Register DONs
		registerDONsReport, err := operations.ExecuteOperation(b, contracts.RegisterDonsOp, contracts.RegisterDonsOpDeps{Env: deps.Env}, contracts.RegisterDonsOpInput{
			ChainSelector: input.RegistryChainSel,
			Address:       input.ContractAddress,
			DONs:          input.DONs,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, err
		}

		return ConfigureCapabilitiesRegistrySeqOutput{
			Nops:         registerNopsReport.Output.Nops,
			Nodes:        registerNodesReport.Output.Nodes,
			Capabilities: registerCapabilitiesReport.Output.Capabilities,
			DONs:         registerDONsReport.Output.DONs,
		}, nil
	},
)
