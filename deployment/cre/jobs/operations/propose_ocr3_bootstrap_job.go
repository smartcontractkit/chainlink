package operations

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

const (
	DefaultBootstrapJobName = "OCR3 MultiChain Capability Bootstrap"
)

type ProposeOCR3BootstrapJobDeps struct {
	Env cldf.Environment
}

type ProposeOCR3BootstrapJobInput struct {
	DONName          string
	Domain           string
	ContractID       string
	EnvironmentLabel string
	ChainSelectorEVM uint64

	JobName string // Optional job name, if not provided, the default will be used.

	DONFilters  []offchain.TargetDONFilter
	ExtraLabels map[string]string
}

type ProposeOCR3BootstrapJobOutput struct {
	Specs map[string][]string
}

var ProposeOCR3BootstrapJob = operations.NewOperation[ProposeOCR3BootstrapJobInput, ProposeOCR3BootstrapJobOutput, ProposeOCR3BootstrapJobDeps](
	"propose-ocr3-bootstrap-job-op",
	semver.MustParse("1.0.0"),
	"Propose OCR3 Bootstrap Job",
	func(b operations.Bundle, deps ProposeOCR3BootstrapJobDeps, input ProposeOCR3BootstrapJobInput) (ProposeOCR3BootstrapJobOutput, error) {
		extJobID, err := pkg.BootstrapExternalJobID(input.DONName, input.ContractID, input.ChainSelectorEVM)
		if err != nil {
			return ProposeOCR3BootstrapJobOutput{}, fmt.Errorf("failed to generate external job ID: %w", err)
		}

		chainID, err := chainsel.GetChainIDFromSelector(input.ChainSelectorEVM)
		if err != nil {
			return ProposeOCR3BootstrapJobOutput{}, fmt.Errorf("failed to get chain ID from selector: %w", err)
		}

		jobName := DefaultBootstrapJobName
		if input.JobName != "" {
			jobName = input.JobName
		}

		cfg := pkg.BootstrapCfg{
			JobName:       jobName,
			ExternalJobID: extJobID,
			ContractID:    input.ContractID,
			ChainID:       chainID,
		}
		if err = cfg.Validate(); err != nil {
			return ProposeOCR3BootstrapJobOutput{}, fmt.Errorf("invalid bootstrap config: %w", err)
		}

		spec, err := cfg.ResolveSpec()
		if err != nil {
			return ProposeOCR3BootstrapJobOutput{}, fmt.Errorf("failed to resolve bootstrap job spec: %w", err)
		}

		report, err := operations.ExecuteOperation(b, ProposeJobSpec, ProposeJobSpecDeps(deps), ProposeJobSpecInput{
			Domain:      input.Domain,
			DONName:     input.DONName,
			Spec:        spec,
			JobLabels:   input.ExtraLabels,
			DONFilters:  input.DONFilters,
			IsBootstrap: true,
		})
		if err != nil {
			return ProposeOCR3BootstrapJobOutput{}, fmt.Errorf("failed to propose bootstrap job: %w", err)
		}

		return ProposeOCR3BootstrapJobOutput{
			Specs: report.Output.Specs,
		}, nil
	},
)
