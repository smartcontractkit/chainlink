package sequences

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	operations2 "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"

	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

const (
	DefaultBootstrapJobName = "OCR3 MultiChain Capability Bootstrap"
)

type ProposeVaultBootstrapJobsDeps struct {
	Env cldf.Environment
}

type ProposeVaultBootstrapJobsInput struct {
	DONName                 string
	Domain                  string
	ContractQualifierPrefix string
	EnvironmentLabel        string
	ChainSelectorEVM        uint64

	JobName string // Optional job name, if not provided, the default will be used.

	DONFilters  []offchain.TargetDONFilter
	ExtraLabels map[string]string
}

type ProposeVaultBootstrapJobsOutput struct {
	Specs map[string][]string
}

var ProposeVaultBootstrapJobs = operations.NewSequence[ProposeVaultBootstrapJobsInput, ProposeVaultBootstrapJobsOutput, ProposeVaultBootstrapJobsDeps](
	"propose-vault-boostrap-jobs-seq",
	semver.MustParse("1.0.0"),
	"Propose Vault Bootstrap Jobs",
	func(b operations.Bundle, deps ProposeVaultBootstrapJobsDeps, input ProposeVaultBootstrapJobsInput) (ProposeVaultBootstrapJobsOutput, error) {
		// 1. Deploy Vault Bootstrappers job for plugin contract
		pluginContractQualifier := fmt.Sprintf("%s_%s", input.ContractQualifierPrefix, "plugin")
		pluginRefKey := pkg.GetOCR3CapabilityAddressRefKey(input.ChainSelectorEVM, pluginContractQualifier)
		pluginAddrRef, err := deps.Env.DataStore.Addresses().Get(pluginRefKey)
		if err != nil {
			return ProposeVaultBootstrapJobsOutput{}, fmt.Errorf("failed to get Vault Plugin contract address for chain selector %d and qualifier %s: %w", input.ChainSelectorEVM, pluginContractQualifier, err)
		}

		r1, rErr := operations.ExecuteOperation(
			b,
			operations2.ProposeOCR3BootstrapJob,
			operations2.ProposeOCR3BootstrapJobDeps{Env: deps.Env},
			operations2.ProposeOCR3BootstrapJobInput{
				Domain:           input.Domain,
				DONName:          input.DONName,
				ContractID:       pluginAddrRef.Address,
				EnvironmentLabel: input.EnvironmentLabel,
				ChainSelectorEVM: input.ChainSelectorEVM,
				JobName:          input.JobName + " (Plugin)",
				DONFilters:       input.DONFilters,
				ExtraLabels:      input.ExtraLabels,
			},
		)
		if rErr != nil {
			return ProposeVaultBootstrapJobsOutput{}, fmt.Errorf("failed to propose Vault Plugin bootstrap job: %w", rErr)
		}

		output := ProposeVaultBootstrapJobsOutput{
			Specs: make(map[string][]string),
		}
		for k, v := range r1.Output.Specs {
			if v == nil {
				output.Specs[k] = []string{}
			}

			output.Specs[k] = append(output.Specs[k], v...)
		}

		// 2. Deploy Vault Bootstrappers job for DKG contract
		dkgContractQualifier := fmt.Sprintf("%s_%s", input.ContractQualifierPrefix, "dkg")
		dkgRefKey := pkg.GetOCR3CapabilityAddressRefKey(input.ChainSelectorEVM, dkgContractQualifier)
		dkgAddrRef, err := deps.Env.DataStore.Addresses().Get(dkgRefKey)
		if err != nil {
			return ProposeVaultBootstrapJobsOutput{}, fmt.Errorf("failed to get Vault DKG contract address for chain selector %d and qualifier %s: %w", input.ChainSelectorEVM, pluginContractQualifier, err)
		}

		r2, rErr := operations.ExecuteOperation(
			b,
			operations2.ProposeOCR3BootstrapJob,
			operations2.ProposeOCR3BootstrapJobDeps{Env: deps.Env},
			operations2.ProposeOCR3BootstrapJobInput{
				Domain:           input.Domain,
				DONName:          input.DONName,
				ContractID:       dkgAddrRef.Address,
				EnvironmentLabel: input.EnvironmentLabel,
				ChainSelectorEVM: input.ChainSelectorEVM,
				DONFilters:       input.DONFilters,
				JobName:          input.JobName + " (DKG)",
				ExtraLabels:      input.ExtraLabels,
			},
		)
		if rErr != nil {
			return ProposeVaultBootstrapJobsOutput{}, fmt.Errorf("failed to propose Vault Plugin bootstrap job: %w", rErr)
		}

		for k, v := range r2.Output.Specs {
			if v == nil {
				output.Specs[k] = []string{}
			}

			output.Specs[k] = append(output.Specs[k], v...)
		}

		return output, nil
	},
)
