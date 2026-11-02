package operations

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type ProposeCapabilitiesRegistryJobDeps struct {
	Env cldf.Environment
}

type ProposeCapabilitiesRegistryJobInput struct {
	Domain         string
	Environment    string
	DONName        string
	DONFilters     []offchain.TargetDONFilter
	ExtraLabels    map[string]string
	OffchainConfig string // proto-JSON
}

type ProposeCapabilitiesRegistryJobOutput struct {
	Specs map[string][]string
}

var ProposeCapabilitiesRegistryJob = operations.NewOperation[
	ProposeCapabilitiesRegistryJobInput,
	ProposeCapabilitiesRegistryJobOutput,
	ProposeCapabilitiesRegistryJobDeps,
](
	"propose-capabilities-registry-job-op",
	semver.MustParse("1.0.0"),
	"Propose CapabilitiesRegistry CRESettings Job",
	func(b operations.Bundle, deps ProposeCapabilitiesRegistryJobDeps, input ProposeCapabilitiesRegistryJobInput) (output ProposeCapabilitiesRegistryJobOutput, err error) {
		job := pkg.CapabilitiesRegistryJob{OffchainConfig: input.OffchainConfig}
		jobSpec, err := job.ResolveJob()
		if err != nil {
			return ProposeCapabilitiesRegistryJobOutput{}, fmt.Errorf("failed to resolve capabilities registry job spec: %w", err)
		}

		report, err := operations.ExecuteOperation(b, ProposeJobSpec, ProposeJobSpecDeps(deps), ProposeJobSpecInput{
			Domain:      input.Domain,
			Environment: input.Environment,
			DONName:     input.DONName,
			JobLabels:   input.ExtraLabels,
			DONFilters:  input.DONFilters,
			Spec:        jobSpec,
		})
		if err != nil {
			return ProposeCapabilitiesRegistryJobOutput{}, fmt.Errorf("failed to propose capabilities registry job: %w", err)
		}

		return ProposeCapabilitiesRegistryJobOutput{
			Specs: report.Output.Specs,
		}, nil
	},
)
