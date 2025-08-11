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
}

func (c ConfigureCapabilitiesRegistrySeqInput) Validate() error {
	if c.ContractAddress == "" {
		return errors.New("ContractAddress is not set")
	}
	return nil
}

type ConfigureCapabilitiesRegistrySeqOutput struct {
	Nops []*capabilities_registry_v2.CapabilitiesRegistryNodeOperatorAdded
}

var ConfigureCapabilitiesRegistrySequence = operations.NewSequence(
	"configure-capabilities-registry",
	semver.MustParse("1.0.0"),
	"Configures the capabilities registry by registering node operators, nodes, dons and capabilities",
	func(b operations.Bundle, deps ConfigureCapabilitiesRegistrySeqDeps, input ConfigureCapabilitiesRegistrySeqInput) (ConfigureCapabilitiesRegistrySeqOutput, error) {
		// Execute operations in sequence
		deployReport, err := operations.ExecuteOperation(b, contracts.RegisterNopsOp, contracts.RegisterNopsOpDeps{Env: deps.Env}, contracts.RegisterNopsOpInput{
			ChainSelector: input.RegistryChainSel,
			Address:       input.ContractAddress,
			Nops:          input.Nops,
		})
		if err != nil {
			return ConfigureCapabilitiesRegistrySeqOutput{}, err
		}

		return ConfigureCapabilitiesRegistrySeqOutput{
			Nops: deployReport.Output.Nops,
		}, nil
	},
)
