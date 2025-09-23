package sequences

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ocr3ops "github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
)

type DeployVaultDeps struct {
	Env *cldf.Environment
}

type DeployVaultInput struct {
	ChainSelector uint64
	Qualifier     string
}

type DeployVaultOutput struct {
	PluginAddress   string
	PluginQualifier string
	DKGAddress      string
	DKGQualifier    string
	ChainSelector   uint64
	Type            string
	Version         string
	Labels          []string
	Datastore       datastore.DataStore
}

// DeployOCR3 is an operation that deploys the OCR3 contract.
// This atomic operation performs the single side effect of deploying and registering the contract.
var DeployVault = operations.NewSequence(
	"deploy-vault-seq",
	semver.MustParse("1.0.0"),
	"Deploy Vault Contracts (OCR3 and DKG)",
	func(b operations.Bundle, deps DeployVaultDeps, input DeployVaultInput) (DeployVaultOutput, error) {
		// 1. Deploy Vault Plugin OCR3 contract
		output1, err := operations.ExecuteOperation(
			b,
			ocr3ops.DeployOCR3,
			ocr3ops.DeployOCR3Deps{
				Env: deps.Env,
			},
			ocr3ops.DeployOCR3Input{
				ChainSelector: input.ChainSelector,
				Qualifier:     fmt.Sprintf("%s_%s", input.Qualifier, "plugin"),
			},
		)
		if err != nil {
			return DeployVaultOutput{}, fmt.Errorf("failed to deploy Vault Plugin contract: %w", err)
		}

		vaultOutput := DeployVaultOutput{
			PluginAddress:   output1.Output.Address,
			PluginQualifier: output1.Output.Qualifier,
			ChainSelector:   output1.Output.ChainSelector,
			Type:            output1.Output.Type,
			Version:         output1.Output.Version,
			Labels:          output1.Output.Labels,
		}

		ds := datastore.NewMemoryDataStore()
		err = ds.Merge(output1.Output.Datastore.Seal())
		if err != nil {
			return DeployVaultOutput{}, fmt.Errorf("failed to merge datastore from plugin contract: %w", err)
		}

		// 2. Deploy Vault DKG OCR3 contract
		output2, err := operations.ExecuteOperation(
			b,
			ocr3ops.DeployOCR3,
			ocr3ops.DeployOCR3Deps{
				Env: deps.Env,
			},
			ocr3ops.DeployOCR3Input{
				ChainSelector: input.ChainSelector,
				Qualifier:     fmt.Sprintf("%s_%s", input.Qualifier, "dkg"),
			},
		)
		if err != nil {
			return DeployVaultOutput{}, fmt.Errorf("failed to deploy Vault Plugin contract: %w", err)
		}

		vaultOutput.DKGAddress = output2.Output.Address
		vaultOutput.DKGQualifier = output2.Output.Qualifier

		err = ds.Merge(output2.Output.Datastore.Seal())
		if err != nil {
			return DeployVaultOutput{}, fmt.Errorf("failed to merge datastore from dkg contract: %w", err)
		}

		vaultOutput.Datastore = ds.Seal()

		return vaultOutput, nil
	},
)
