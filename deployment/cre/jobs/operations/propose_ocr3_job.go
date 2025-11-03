package operations

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

const defaultVaultRequestExpiryDuration = "10s"

type ProposeOCR3JobDeps struct {
	Env cldf.Environment
}

type ProposeOCR3JobInput struct {
	Domain  string
	EnvName string

	DONName string
	JobName string

	TemplateName         string
	ContractAddress      string
	ChainSelectorEVM     uint64
	ChainSelectorAptos   uint64
	BootstrapperOCR3Urls []string

	// Optionals: specific to the worker vault OCR3 Job spec
	DKGContractAddress         string
	VaultRequestExpiryDuration string

	JobNodesSpecifier offchain.NodesSpecifier
	ExtraLabels       map[string]string
}

type ProposeOCR3JobOutput struct {
	Specs map[string][]string
}

var ProposeOCR3Job = operations.NewSequence[ProposeOCR3JobInput, ProposeOCR3JobOutput, ProposeOCR3JobDeps](
	"propose-ocr3-job-seq",
	semver.MustParse("1.0.0"),
	"Propose OCR3 Job",
	func(b operations.Bundle, deps ProposeOCR3JobDeps, input ProposeOCR3JobInput) (ProposeOCR3JobOutput, error) {
		err := input.JobNodesSpecifier.Validate()
		if err != nil {
			return ProposeOCR3JobOutput{}, fmt.Errorf("invalid job nodes specifier: %w", err)
		}
		if input.JobNodesSpecifier.LabelFilters != nil {
			// Backward compatibility: append to existing filters.
			// We only want to target plugin nodes for OCR3 jobs.
			input.JobNodesSpecifier.LabelFilters = append(input.JobNodesSpecifier.LabelFilters, offchain.NodeLabelFilter{
				Key:   "type",
				Value: "plugin",
			})
		}

		nodes, err := offchain.FetchNodesFromJD(b.GetContext(), deps.Env.Offchain, input.JobNodesSpecifier.Filter(input.DONName, input.EnvName, input.Domain))
		if err != nil {
			return ProposeOCR3JobOutput{}, fmt.Errorf("failed to fetch nodes from JD: %w", err)
		}
		if len(nodes) == 0 {
			return ProposeOCR3JobOutput{}, fmt.Errorf("no nodes found for DON `%s` with provided specifiers %+v, filter %+v", input.DONName, input.JobNodesSpecifier, input.JobNodesSpecifier.Filter(input.DONName, input.EnvName, input.Domain))
		}
		nodesIDs := make([]string, 0, len(nodes))
		for _, n := range nodes {
			nodesIDs = append(nodesIDs, n.Id)
		}
		b.Logger.Debugw("Proposing OCR3 job", "DON", input.DONName, "domain", input.Domain, "environment", input.EnvName, "nodes_count", len(nodes), "node_ids", nodesIDs)

		vaultReqExpiry := input.VaultRequestExpiryDuration
		if vaultReqExpiry == "" {
			vaultReqExpiry = defaultVaultRequestExpiryDuration
		}

		specs, err := pkg.BuildOCR3JobConfigSpecs(
			deps.Env.Offchain, deps.Env.Logger, input.ContractAddress, input.ChainSelectorEVM,
			input.ChainSelectorAptos, nodesIDs, input.BootstrapperOCR3Urls, input.DONName, input.JobName, input.TemplateName, input.DKGContractAddress, vaultReqExpiry,
		)
		if err != nil {
			return ProposeOCR3JobOutput{}, fmt.Errorf("failed to build OCR3 job config specs: %w", err)
		}

		finalSpecs := make(map[string][]string)

		var mergedErrs error
		// Propose each spec to its target node.
		for _, spec := range specs {
			// limit to the specific node
			specifier := offchain.NodesSpecifier{
				NodeIDs: []string{spec.NodeID},
			}
			opReport, opErr := operations.ExecuteOperation(b, ProposeJobSpec, ProposeJobSpecDeps(deps), ProposeJobSpecInput{
				Domain:            input.Domain,
				DONName:           input.DONName,
				Spec:              spec.Spec,
				JobNodesSpecifier: specifier,
				JobLabels:         input.ExtraLabels,
			})
			if opErr != nil {
				// Do not fail the sequence if a single proposal fails, make it through all proposals.
				mergedErrs = fmt.Errorf("error proposing OCR3 job to node %s spec %s: %w", spec.NodeID, spec.Spec, opErr)
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
