package operations

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type ProposeOCR3JobDeps struct {
	Env cldf.Environment
}

type ProposeOCR3JobInput struct {
	Domain  string
	EnvName string

	DONName      string
	JobName      string
	TemplateName string

	ContractAddress      string
	ChainSelectorEVM     uint64
	ChainSelectorAptos   uint64
	BootstrapperOCR3Urls []string

	DONFilters  []offchain.TargetDONFilter
	ExtraLabels map[string]string
}

type ProposeOCR3JobOutput struct {
	Specs map[string][]string
}

var ProposeOCR3Job = operations.NewSequence[ProposeOCR3JobInput, ProposeOCR3JobOutput, ProposeOCR3JobDeps](
	"propose-ocr3-job-seq",
	semver.MustParse("1.0.0"),
	"Propose OCR3 Job",
	func(b operations.Bundle, deps ProposeOCR3JobDeps, input ProposeOCR3JobInput) (ProposeOCR3JobOutput, error) {
		nodes, err := pkg.FetchNodesFromJD(b.GetContext(), deps.Env, pkg.FetchNodesRequest{
			Domain:  input.Domain,
			Filters: input.DONFilters,
		})
		if err != nil {
			return ProposeOCR3JobOutput{}, fmt.Errorf("failed to fetch nodes from JD: %w", err)
		}

		specs, err := pkg.BuildOCR3JobConfigSpecs(
			deps.Env.Offchain, deps.Env.Logger, input.ContractAddress, input.ChainSelectorEVM,
			input.ChainSelectorAptos, nodes, input.BootstrapperOCR3Urls, input.DONName, input.JobName, input.TemplateName,
		)
		if err != nil {
			return ProposeOCR3JobOutput{}, fmt.Errorf("failed to build OCR3 job config specs: %w", err)
		}

		finalSpecs := make(map[string][]string)

		var mergedErrs error
		for _, spec := range specs {
			opReport, opErr := operations.ExecuteOperation(b, ProposeJobSpec, ProposeJobSpecDeps(deps), ProposeJobSpecInput{
				Domain:     input.Domain,
				DONName:    input.DONName,
				Spec:       spec.Spec,
				DONFilters: input.DONFilters,
				JobLabels:  input.ExtraLabels,
			})
			if opErr != nil {
				// Do not fail the sequence if a single proposal fails, make it through all proposals.
				mergedErrs = fmt.Errorf("error proposing job to node %s spec %s: %w", spec.NodeID, spec.Spec, opErr)
				continue
			}

			for nodeID, s := range opReport.Output.Specs {
				finalSpecs[nodeID] = append(finalSpecs[nodeID], s...)
			}
		}

		return ProposeOCR3JobOutput{
			Specs: finalSpecs,
		}, mergedErrs
	},
)
