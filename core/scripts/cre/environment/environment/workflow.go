package environment

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
)

const (
	DefaultArtifactsDir        = "/home/chainlink/workflows"
	DefaultWorkflowNodePattern = "workflow-node"
)

func workflowCmds() *cobra.Command {
	workflowCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Workflow management commands",
		Long:  `Commands to manage workflows`,
	}

	workflowCmd.AddCommand(deployAndVerifyExampleWorkflowCmd())
	workflowCmd.AddCommand(deployWorkflowCmd())

	return workflowCmd
}

func deleteAllWorkflows(ctx context.Context, rpcURL, workflowRegistryAddress string) error {
	var privateKey string
	if os.Getenv("PRIVATE_KEY") != "" {
		privateKey = os.Getenv("PRIVATE_KEY")
	} else {
		privateKey = blockchain.DefaultAnvilPrivateKey
	}

	sethClient, scErr := seth.NewClientBuilder().
		WithRpcUrl(rpcURL).
		WithPrivateKeys([]string{privateKey}).
		WithProtections(false, false, seth.MustMakeDuration(time.Minute)).
		Build()
	if scErr != nil {
		return errors.Wrap(scErr, "failed to create Seth client")
	}

	fmt.Printf("\n⚙️ Deleting all workflows from the workflow registry\n\n")

	deleteErr := creworkflow.DeleteAllWithContract(ctx, sethClient, common.HexToAddress(workflowRegistryAddress))
	if deleteErr != nil {
		return errors.Wrapf(deleteErr, "❌ failed to delete all workflows from the registry %s", workflowRegistryAddress)
	}

	fmt.Printf("\n✅ All workflows deleted from the workflow registry\n\n")

	return nil
}

func deployWorkflowCmd() *cobra.Command {
	var (
		workflowFilePathFlag        string
		configFilePathFlag          string
		containerTargetDirFlag      string
		containerNamePatternFlag    string
		workflowNameFlag            string
		workflowOwnerAddressFlag    string
		workflowRegistryAddressFlag string
		chainIDFlag                 uint64
		rpcURLFlag                  string
	)

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Compiles and uploads a workflow to the environment",
		Long:  `Compiles and uploads a workflow to the environment by copying it to workflow nodes and registering with the workflow registry`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return compileCopyAndRegisterWorkflow(cmd.Context(), workflowFilePathFlag, workflowNameFlag, workflowRegistryAddressFlag, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, rpcURLFlag)
		},
	}

	cmd.Flags().StringVarP(&workflowFilePathFlag, "workflow-file-path", "w", "./examples/workflows/v2/cron/main.go", "Path to the workflow file")
	cmd.Flags().StringVarP(&configFilePathFlag, "config-file-path", "c", "", "Path to the config file")
	cmd.Flags().StringVarP(&containerTargetDirFlag, "container-target-dir", "t", DefaultArtifactsDir, "Path to the target directory in the Docker container")
	cmd.Flags().StringVarP(&containerNamePatternFlag, "container-name-pattern", "n", DefaultWorkflowNodePattern, "Pattern to match the container name")
	cmd.Flags().Uint64VarP(&chainIDFlag, "chain-id", "i", 1337, "Chain ID")
	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "workflow-owner-address", "o", "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0", "Workflow registry address")
	cmd.Flags().StringVarP(&workflowNameFlag, "workflow-name", "m", "exampleworkflow", "Workflow name")

	return cmd
}

func compileCopyAndRegisterWorkflow(ctx context.Context, workflowFilePathFlag, workflowNameFlag, workflowRegistryAddressFlag, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, rpcURLFlag string) error {
	fmt.Printf("\n⚙️ Compiling workflow from %s\n", workflowFilePathFlag)

	compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(workflowFilePathFlag, workflowNameFlag)
	if compileErr != nil {
		return errors.Wrap(compileErr, "❌ failed to compile workflow")
	}

	fmt.Printf("\n✅ Workflow compiled and compressed successfully\n\n")

	copyErr := creworkflow.CopyWorkflowToDockerContainers(compressedWorkflowWasmPath, containerNamePatternFlag, containerTargetDirFlag)
	if copyErr != nil {
		return errors.Wrap(copyErr, "❌ failed to copy workflow to Docker container")
	}

	fmt.Printf("\n✅ Workflow copied to Docker containers\n")
	fmt.Printf("\n⚙️ Creating Seth client\n\n")

	var privateKey string
	if os.Getenv("PRIVATE_KEY") != "" {
		privateKey = os.Getenv("PRIVATE_KEY")
	} else {
		privateKey = blockchain.DefaultAnvilPrivateKey
	}

	sethClient, scErr := seth.NewClientBuilder().
		WithRpcUrl(rpcURLFlag).
		WithPrivateKeys([]string{privateKey}).
		WithProtections(false, false, seth.MustMakeDuration(time.Minute)).
		Build()
	if scErr != nil {
		return errors.Wrap(scErr, "failed to create Seth client")
	}

	var configPath *string
	if configFilePathFlag != "" {
		fmt.Printf("\n⚙️ Copying workflow config file to Docker container\n\n")
		configPathAbs, configPathAbsErr := filepath.Abs(configFilePathFlag)
		if configPathAbsErr != nil {
			return errors.Wrap(configPathAbsErr, "failed to get absolute path of the config file")
		}

		configCopyErr := creworkflow.CopyWorkflowToDockerContainers(configFilePathFlag, containerNamePatternFlag, containerTargetDirFlag)
		if configCopyErr != nil {
			return errors.Wrap(configCopyErr, "❌ failed to copy config file to Docker container")
		}

		configPathAbs = "file://" + configPathAbs
		configPath = &configPathAbs

		fmt.Printf("\n✅ Workflow config file copied to Docker container\n\n")
	}

	fmt.Printf("\n⚙️ Deleting all workflows from the workflow registry\n\n")

	deleteErr := creworkflow.DeleteAllWithContract(ctx, sethClient, common.HexToAddress(workflowRegistryAddressFlag))
	if deleteErr != nil {
		return errors.Wrapf(deleteErr, "❌ failed to delete all workflows from the registry %s", workflowRegistryAddressFlag)
	}

	fmt.Printf("\n✅ All workflows deleted from the workflow registry\n\n")
	fmt.Printf("\n⚙️ Registering workflow %s with the workflow registry\n\n", workflowNameFlag)

	registerErr := creworkflow.RegisterWithContract(ctx, sethClient, common.HexToAddress(workflowRegistryAddressFlag), 1, workflowNameFlag, "file://"+compressedWorkflowWasmPath, configPath, nil, &containerTargetDirFlag)
	if registerErr != nil {
		return errors.Wrapf(registerErr, "❌ failed to register workflow %s", workflowNameFlag)
	}

	defer func() {
		_ = os.Remove(compressedWorkflowWasmPath)
	}()

	fmt.Printf("\n✅ Workflow registered successfully\n\n")

	return nil
}
