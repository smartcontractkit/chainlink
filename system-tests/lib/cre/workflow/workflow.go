package workflow

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"

	workflow_registry_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"

	libnet "github.com/smartcontractkit/chainlink/system-tests/lib/net"

	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
)

func RegisterWithCRECLI(input cretypes.ManageWorkflowWithCRECLIInput) error {
	if registerValErr := validateRegisterWorkflowInput(input); registerValErr != nil {
		return errors.Wrap(registerValErr, "failed to validate RegisterWorkflowInput")
	}

	creCLIWorkflowSettingsFile := input.CRESettingsFile
	// TODO: make sure it's ok to remove the code below
	//creCLIWorkflowSettingsFile, err := prepareCRECLI(input)
	//if err != nil {
	//	return err
	//}

	var workflowURL string
	var workflowConfigURL *string
	var workflowSecretsURL *string

	// compile and upload the workflow, if we are not using an existing one
	if input.ShouldCompileNewWorkflow {
		wasmFilePath := filepath.Join(input.NewWorkflow.FolderLocation, "mytestworkflow.wasm.br")

		if input.S3ProviderOutput != nil {
			creCLI := libcrecli.NewCreCli(input.CRECLIAbsPath)

			err := creCLI.Compile(
				filepath.Join(input.NewWorkflow.FolderLocation, input.NewWorkflow.WorkflowFileName),
				libcrecli.Deref(input.NewWorkflow.ConfigFilePath),
				creCLIWorkflowSettingsFile.Name(),
				wasmFilePath,
			)
			if err != nil {
				return err
			}

			// TODO: once CRE CLI is fixed remove the rename hack
			//       hack to address CRE CLI bugs:
			//       - BUG #1:
			//         current: flag `-o fileName.wasm.br` produces `fileName.wasm.br.b64`
			//         expected: flag `-o fileName.wasm.br` should produce `fileName.wasm.br`
			//       - BUG #2:
			//         when using `upload` command for `fileName.wasm.br.b64` it returns:
			//         `Error: failed to create or update object: ... supported extensions: .wasm.br, .json, .yaml, .yml
			err = os.Rename(wasmFilePath+".b64", wasmFilePath)
			if err != nil {
				fmt.Printf("failed to rename workflow binary file (.b64 to .br): %v\n", err)
			}

			uploadOutput, err := creCLI.Upload(
				libcrecli.MINIO,
				wasmFilePath,
				libcrecli.Deref(input.NewWorkflow.ConfigFilePath),
				creCLIWorkflowSettingsFile.Name(),
			)
			if err != nil {
				return err
			}

			// Register to a Node with S3 storage, so need to replace the endpoint
			// with the base endpoint subjective to the target host.
			workflowURL = strings.Replace(uploadOutput.BinaryURL, input.S3ProviderOutput.Endpoint, input.S3ProviderOutput.BaseEndpoint, -1)
			tmpConfigURL := strings.Replace(uploadOutput.ConfigURL, input.S3ProviderOutput.Endpoint, input.S3ProviderOutput.BaseEndpoint, -1)
			workflowConfigURL = &tmpConfigURL
			// TODO: handle secrets file upload: probably could be done via Upload() adding argument -f <secrets file>
		} else {
			compilationResult, compileErr := libcrecli.CompileWorkflow(input.CRECLIAbsPath, input.NewWorkflow.FolderLocation, input.NewWorkflow.WorkflowFileName, input.NewWorkflow.ConfigFilePath, creCLIWorkflowSettingsFile, input.CRESettingsFile)
			if compileErr != nil {
				return errors.Wrap(compileErr, "failed to compile workflow")
			}

			workflowURL = compilationResult.WorkflowURL
			workflowConfigURL = &compilationResult.ConfigURL

			if input.NewWorkflow.SecretsFilePath != nil && *input.NewWorkflow.SecretsFilePath != "" {
				secretsURL, secretsErr := libcrecli.EncryptSecrets(input.CRECLIAbsPath, *input.NewWorkflow.SecretsFilePath, input.NewWorkflow.Secrets, creCLIWorkflowSettingsFile)
				if secretsErr != nil {
					return errors.Wrap(secretsErr, "failed to encrypt workflow secrets")
				}
				workflowSecretsURL = &secretsURL
			}
		}
	} else {
		workflowURL = input.ExistingWorkflow.BinaryURL
		workflowConfigURL = input.ExistingWorkflow.ConfigURL
		workflowSecretsURL = input.ExistingWorkflow.SecretsURL
	}

	// TODO: another fix required to CRE CLI:
	//       - BUG #3: Downloads remote files even if they are already present in the local filesystem.
	//                 The files are required to compute Workflow ID which is used to register the workflow.
	//                 This is not a problem for the local minio, because Nodes are running in docker
	// 			       and the URLs need to be relative to them, i.e. `minio:9000`
	//			       (which is different than the host machine `127.0.0.1:9000`).
	//                 CRE CLI executes in the local host context, so it can't reach docker internal hostname.
	//                 Downloading files from remote is excessive in this case and could be avoided.
	registerErr := libcrecli.DeployWorkflow(
		input.CRECLIAbsPath,
		workflowURL,
		workflowConfigURL,
		workflowSecretsURL,
		creCLIWorkflowSettingsFile,
		nil, // without workflow path
	)
	if registerErr != nil {
		return errors.Wrap(registerErr, "failed to register workflow")
	}

	return nil
}

func validateRegisterWorkflowInput(input cretypes.ManageWorkflowWithCRECLIInput) error {
	if input.ShouldCompileNewWorkflow && input.NewWorkflow == nil {
		return errors.New("NewWorkflow is required when ShouldCompileNewWorkflow is true")
	}
	if !input.ShouldCompileNewWorkflow && input.ExistingWorkflow == nil {
		return errors.New("ExistingWorkflow is required when ShouldCompileNewWorkflow is false")
	}
	if input.NewWorkflow != nil && input.NewWorkflow.FolderLocation == "" {
		return errors.New("WorkflowFolderLocation is required when ShouldCompileNewWorkflow is true")
	}
	if input.NewWorkflow != nil && input.NewWorkflow.WorkflowFileName == "" {
		return errors.New("WorkflowFileName is required when ShouldCompileNewWorkflow is true")
	}
	if input.ExistingWorkflow != nil && input.ExistingWorkflow.BinaryURL == "" {
		return errors.New("BinaryURL is required when ShouldCompileNewWorkflow is false")
	}

	return nil
}

func prepareCRECLI(input cretypes.ManageWorkflowWithCRECLIInput) (*os.File, error) {
	if valErr := input.Validate(); valErr != nil {
		return nil, errors.Wrap(valErr, "failed to validate WorkflowInput")
	}

	if pkErr := os.Setenv("CRE_ETH_PRIVATE_KEY", input.CRECLIPrivateKey); pkErr != nil {
		return nil, errors.Wrap(pkErr, "failed to set CRE_ETH_PRIVATE_KEY")
	}

	if profileErr := os.Setenv("CRE_PROFILE", input.CRECLIProfile); profileErr != nil {
		return nil, errors.Wrap(profileErr, "failed to set CRE_PROFILE")
	}

	return libcrecli.PrepareCRECLIWorkflowSettingsFile(input.CRECLIProfile, input.WorkflowOwnerAddress, input.WorkflowName)
}

func PauseWithCRECLI(input cretypes.ManageWorkflowWithCRECLIInput) error {
	creCLIWorkflowSettingsFile, err := prepareCRECLI(input)
	if err != nil {
		return err
	}

	pauseErr := libcrecli.PauseWorkflow(input.CRECLIAbsPath, creCLIWorkflowSettingsFile)
	if pauseErr != nil {
		return errors.Wrap(pauseErr, "failed to pause workflow")
	}

	return nil
}

func ActivateWithCRECLI(input cretypes.ManageWorkflowWithCRECLIInput) error {
	creCLIWorkflowSettingsFile, err := prepareCRECLI(input)
	if err != nil {
		return err
	}

	activateErr := libcrecli.ActivateWorkflow(input.CRECLIAbsPath, creCLIWorkflowSettingsFile)
	if activateErr != nil {
		return errors.Wrap(activateErr, "failed to activate workflow")
	}

	return nil
}

func RegisterWithContract(ctx context.Context, sc *seth.Client, workflowRegistryAddr common.Address, donID uint32, workflowName, binaryURL string, configURL, secretsURL *string) error {
	workFlowData, err := libnet.DownloadAndDecodeBase64(ctx, binaryURL)
	if err != nil {
		return errors.Wrap(err, "failed to download and decode workflow binary")
	}

	var configData []byte
	configURLToUse := ""
	if configURL != nil {
		configData, err = libnet.Download(ctx, *configURL)
		if err != nil {
			return errors.Wrap(err, "failed to download workflow config")
		}
		configURLToUse = *configURL
	}

	secretsURLToUse := ""
	if secretsURL != nil {
		secretsURLToUse = *secretsURL
	}

	// use non-encoded workflow name
	workflowID, idErr := generateWorkflowIDFromStrings(sc.MustGetRootKeyAddress().Hex(), workflowName, workFlowData, configData, secretsURLToUse)
	if idErr != nil {
		return errors.Wrap(idErr, "failed to generate workflow ID")
	}

	workflowRegistryInstance, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	if err != nil {
		return errors.Wrap(err, "failed to create workflow registry instance")
	}

	// use non-encoded workflow name
	_, decodeErr := sc.Decode(workflowRegistryInstance.RegisterWorkflow(sc.NewTXOpts(), workflowName, [32]byte(common.Hex2Bytes(workflowID)), donID, uint8(0), binaryURL, configURLToUse, secretsURLToUse))
	if decodeErr != nil {
		return errors.Wrap(decodeErr, "failed to register workflow")
	}

	return nil
}

func generateWorkflowIDFromStrings(owner string, name string, workflow []byte, config []byte, secretsURL string) (string, error) {
	ownerWithoutPrefix := owner
	if strings.HasPrefix(owner, "0x") {
		ownerWithoutPrefix = owner[2:]
	}

	ownerb, err := hex.DecodeString(ownerWithoutPrefix)
	if err != nil {
		return "", err
	}

	wid, err := pkgworkflows.GenerateWorkflowID(ownerb, name, workflow, config, secretsURL)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(wid[:]), nil
}
