package workflow

import (
	"os"

	"github.com/pkg/errors"

	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
)

func RegisterWithCRECLI(input cretypes.RegisterWorkflowWithCRECLIInput) error {
	if valErr := input.Validate(); valErr != nil {
		return errors.Wrap(valErr, "failed to validate RegisterWorkflowInput")
	}

	// This env var is required by the CRE CLI
	pkErr := os.Setenv("CRE_ETH_PRIVATE_KEY", input.CRECLIPrivateKey)
	if pkErr != nil {
		return errors.Wrap(pkErr, "failed to set CRE_ETH_PRIVATE_KEY")
	}

	var workflowURL string
	var workflowConfigURL *string
	var workflowSecretsURL *string

	// compile and upload the workflow, if we are not using an existing one
	if input.ShouldCompileNewWorkflow {
		compilationResult, compileErr := libcrecli.CompileWorkflow(input.CRECLIAbsPath, input.NewWorkflow.FolderLocation, input.NewWorkflow.WorkflowFileName, input.NewWorkflow.ConfigFilePath, input.CRESettingsFile)
		if compileErr != nil {
			return errors.Wrap(compileErr, "failed to compile workflow")
		}

		workflowURL = compilationResult.WorkflowURL
		workflowConfigURL = &compilationResult.ConfigURL

		if input.NewWorkflow.SecretsFilePath != nil && *input.NewWorkflow.SecretsFilePath != "" {
			secretsURL, secretsErr := libcrecli.EncryptSecrets(input.CRECLIAbsPath, *input.NewWorkflow.SecretsFilePath, input.NewWorkflow.Secrets, input.CRESettingsFile)
			if secretsErr != nil {
				return errors.Wrap(secretsErr, "failed to encrypt workflow secrets")
			}
			workflowSecretsURL = &secretsURL
		}
	} else {
		workflowURL = input.ExistingWorkflow.BinaryURL
		workflowConfigURL = input.ExistingWorkflow.ConfigURL
		workflowSecretsURL = input.ExistingWorkflow.SecretsURL
	}

	registerErr := libcrecli.DeployWorkflow(input.CRECLIAbsPath, input.WorkflowName, workflowURL, workflowConfigURL, workflowSecretsURL, input.CRESettingsFile)
	if registerErr != nil {
		return errors.Wrap(registerErr, "failed to register workflow")
	}

	return nil
}
