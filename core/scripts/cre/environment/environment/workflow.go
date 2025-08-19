package environment

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
)

const (
	DefaultArtifactsDir        = "/home/chainlink/workflows"
	DefaultWorkflowNodePattern = "workflow-node"

	// Might change if deployment sequence changes or if different config file than 'configs/workflow-don.toml' is used
	DefaultWorkflowRegistryAddress     = "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
	DefaultCapabilitiesRegistryAddress = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"

	DefaultWorkflowOwnerAddress = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
)

func workflowCmds() *cobra.Command {
	workflowCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Workflow management commands",
		Long:  `Commands to manage workflows`,
	}

	workflowCmd.AddCommand(deployAndVerifyExampleWorkflowCmd())
	workflowCmd.AddCommand(compileDeployWorkflowCmd())
	workflowCmd.AddCommand(deleteWorkflowCmd())
	workflowCmd.AddCommand(deleteAllWorkflowsCmd())
	workflowCmd.AddCommand(compileWorkflowCmd())
	workflowCmd.AddCommand(deployWorkflowCmd())

	return workflowCmd
}

func deleteAllWorkflows(ctx context.Context, rpcURL, workflowRegistryAddress string) error {
	if pkErr := creenv.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
		return pkErr
	}

	sethClient, scErr := seth.NewClientBuilder().
		WithRpcUrl(rpcURL).
		WithPrivateKeys([]string{os.Getenv("PRIVATE_KEY")}).
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

func compileWorkflowCmd() *cobra.Command {
	var (
		workflowFilePathFlag string
		workflowNameFlag     string
	)

	cmd := &cobra.Command{
		Use:   "compile",
		Short: "Compiles a workflow",
		Long:  `Compiles, compresses with Brotli and encodes with base64 a workflow`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, compileErr := compileWorkflow(workflowFilePathFlag, workflowNameFlag)
			if compileErr != nil {
				return errors.Wrap(compileErr, "❌ failed to compile workflow")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&workflowFilePathFlag, "workflow-file-path", "w", "", "Path to the workflow file")
	cmd.Flags().StringVarP(&workflowNameFlag, "workflow-name", "n", "exampleworkflow", "Workflow name")

	if err := cmd.MarkFlagRequired("workflow-file-path"); err != nil {
		panic(err)
	}

	return cmd
}

func deployWorkflowCmd() *cobra.Command {
	var (
		workflowFilePathFlag            string
		configFilePathFlag              string
		secretsFilePathFlag             string
		containerTargetDirFlag          string
		containerNamePatternFlag        string
		workflowNameFlag                string
		workflowOwnerAddressFlag        string
		workflowRegistryAddressFlag     string
		capabilitiesRegistryAddressFlag string
		deleteWorkflowFileFlag          bool
		donIDFlag                       uint32
		chainIDFlag                     uint64
		rpcURLFlag                      string
	)

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploys a workflow to the environment",
		Long:  `Deploys a workflow to the environment by copying it to workflow nodes and registering with the workflow registry`,
		RunE: func(cmd *cobra.Command, args []string) error {
			initDxTracker()
			var regErr error

			defer func() {
				metaData := map[string]any{}
				if regErr != nil {
					metaData["result"] = "failure"
					metaData["error"] = oneLineErrorMessage(regErr)
				} else {
					metaData["result"] = "success"
				}

				trackingErr := dxTracker.Track("cre.local.workflow.deploy", metaData)
				if trackingErr != nil {
					fmt.Fprintf(os.Stderr, "failed to track workflow deploy: %s\n", trackingErr)
				}
			}()

			regErr = deployWorkflow(cmd.Context(), workflowFilePathFlag, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddressFlag, capabilitiesRegistryAddressFlag, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, rpcURLFlag, donIDFlag, deleteWorkflowFileFlag)

			return regErr
		},
	}

	cmd.Flags().StringVarP(&workflowFilePathFlag, "wasm-file-path", "w", "", "Path to the workflow WASM file")
	cmd.Flags().StringVarP(&configFilePathFlag, "config-file-path", "c", "", "Path to the config file")
	cmd.Flags().StringVarP(&secretsFilePathFlag, "secrets-file-path", "s", "", "Path to the secrets file")
	cmd.Flags().StringVarP(&containerTargetDirFlag, "container-target-dir", "t", DefaultArtifactsDir, "Path to the target directory in the Docker container")
	cmd.Flags().StringVarP(&containerNamePatternFlag, "container-name-pattern", "o", DefaultWorkflowNodePattern, "Pattern to match the container name")
	cmd.Flags().Uint64VarP(&chainIDFlag, "chain-id", "i", 1337, "Chain ID")
	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "workflow-owner-address", "d", DefaultWorkflowOwnerAddress, "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", DefaultWorkflowRegistryAddress, "Workflow registry address")
	cmd.Flags().StringVarP(&capabilitiesRegistryAddressFlag, "capabilities-registry-address", "b", DefaultCapabilitiesRegistryAddress, "Capabilities registry address")
	cmd.Flags().Uint32VarP(&donIDFlag, "don-id", "e", 1, "DON ID")
	cmd.Flags().StringVarP(&workflowNameFlag, "workflow-name", "n", "exampleworkflow", "Workflow name")
	cmd.Flags().BoolVarP(&deleteWorkflowFileFlag, "delete-workflow-file", "l", false, "Delete the workflow file after deployment")

	if err := cmd.MarkFlagRequired("wasm-file-path"); err != nil {
		panic(err)
	}

	return cmd
}

func compileDeployWorkflowCmd() *cobra.Command {
	var (
		workflowFilePathFlag            string
		configFilePathFlag              string
		secretsFilePathFlag             string
		containerTargetDirFlag          string
		containerNamePatternFlag        string
		workflowNameFlag                string
		workflowOwnerAddressFlag        string
		workflowRegistryAddressFlag     string
		capabilitiesRegistryAddressFlag string
		donIDFlag                       uint32
		chainIDFlag                     uint64
		rpcURLFlag                      string
	)

	cmd := &cobra.Command{
		Use:   "compile-deploy",
		Short: "Compiles and uploads a workflow to the environment",
		Long:  `Compiles and uploads a workflow to the environment by copying it to workflow nodes and registering with the workflow registry`,
		RunE: func(cmd *cobra.Command, args []string) error {
			initDxTracker()
			var regErr error

			defer func() {
				metaData := map[string]any{}
				if regErr != nil {
					metaData["result"] = "failure"
					metaData["error"] = oneLineErrorMessage(regErr)
				} else {
					metaData["result"] = "success"
				}

				trackingErr := dxTracker.Track("cre.local.workflow.deploy", metaData)
				if trackingErr != nil {
					fmt.Fprintf(os.Stderr, "failed to track workflow deploy: %s\n", trackingErr)
				}
			}()

			regErr = compileCopyAndRegisterWorkflow(cmd.Context(), workflowFilePathFlag, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddressFlag, capabilitiesRegistryAddressFlag, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, rpcURLFlag, donIDFlag)

			return regErr
		},
	}

	cmd.Flags().StringVarP(&workflowFilePathFlag, "workflow-file-path", "w", "./examples/workflows/v2/cron/main.go", "Path to the workflow file")
	cmd.Flags().StringVarP(&configFilePathFlag, "config-file-path", "c", "", "Path to the config file")
	cmd.Flags().StringVarP(&secretsFilePathFlag, "secrets-file-path", "s", "", "Path to the secrets file")
	cmd.Flags().StringVarP(&containerTargetDirFlag, "container-target-dir", "t", DefaultArtifactsDir, "Path to the target directory in the Docker container")
	cmd.Flags().StringVarP(&containerNamePatternFlag, "container-name-pattern", "o", DefaultWorkflowNodePattern, "Pattern to match the container name")
	cmd.Flags().Uint64VarP(&chainIDFlag, "chain-id", "i", 1337, "Chain ID")
	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "workflow-owner-address", "d", DefaultWorkflowOwnerAddress, "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", DefaultWorkflowRegistryAddress, "Workflow registry address")
	cmd.Flags().StringVarP(&capabilitiesRegistryAddressFlag, "capabilities-registry-address", "b", DefaultCapabilitiesRegistryAddress, "Capabilities registry address")
	cmd.Flags().Uint32VarP(&donIDFlag, "don-id", "e", 1, "DON ID")
	cmd.Flags().StringVarP(&workflowNameFlag, "workflow-name", "n", "exampleworkflow", "Workflow name")

	return cmd
}

func deleteWorkflowCmd() *cobra.Command {
	var (
		workflowNameFlag            string
		workflowOwnerAddressFlag    string
		workflowRegistryAddressFlag string
		chainIDFlag                 uint64
		rpcURLFlag                  string
	)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes a workflow from the workflow registry contract",
		Long:  `Deletes a workflow from the workflow registry contract (but doesn't remove it from the Docker containers)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("\n⚙️ Deleting workflow '%s' from the workflow registry\n\n", workflowNameFlag)

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

			workflowNames, workflowNamesErr := creworkflow.GetWorkflowNames(cmd.Context(), sethClient, common.HexToAddress(workflowRegistryAddressFlag))
			if workflowNamesErr != nil {
				return errors.Wrap(workflowNamesErr, "failed to get workflows from the registry")
			}

			if !slices.Contains(workflowNames, workflowNameFlag) {
				fmt.Printf("\n✅ Workflow '%s' not found in the registry %s. Skipping...\n\n", workflowNameFlag, workflowRegistryAddressFlag)

				return nil
			}

			deleteErr := creworkflow.DeleteWithContract(cmd.Context(), sethClient, common.HexToAddress(workflowRegistryAddressFlag), workflowNameFlag)
			if deleteErr != nil {
				return errors.Wrapf(deleteErr, "❌ failed to delete workflow '%s' from the registry %s", workflowNameFlag, workflowRegistryAddressFlag)
			}

			fmt.Printf("\n✅ Workflow deleted from the workflow registry\n\n")

			return nil
		},
	}

	cmd.Flags().Uint64VarP(&chainIDFlag, "chain-id", "i", 1337, "Chain ID")
	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "workflow-owner-address", "d", DefaultWorkflowOwnerAddress, "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", DefaultWorkflowRegistryAddress, "Workflow registry address")
	cmd.Flags().StringVarP(&workflowNameFlag, "name", "n", "exampleworkflow", "Workflow name")

	return cmd
}

func deleteAllWorkflowsCmd() *cobra.Command {
	var (
		workflowOwnerAddressFlag    string
		workflowRegistryAddressFlag string
		chainIDFlag                 uint64
		rpcURLFlag                  string
	)

	cmd := &cobra.Command{
		Use:   "delete-all",
		Short: "Deletes all workflows from the workflow registry contract",
		Long:  `Deletes all workflows from the workflow registry contract (but doesn't remove them from the Docker containers)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("\n⚙️ Deleting all workflows from the workflow registry\n\n")

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

			deleteErr := creworkflow.DeleteAllWithContract(cmd.Context(), sethClient, common.HexToAddress(workflowRegistryAddressFlag))
			if deleteErr != nil {
				return errors.Wrapf(deleteErr, "❌ failed to delete all workflows from the registry %s", workflowRegistryAddressFlag)
			}

			fmt.Printf("\n✅ All workflows deleted from the workflow registry\n\n")

			return nil
		},
	}

	cmd.Flags().Uint64VarP(&chainIDFlag, "chain-id", "i", 1337, "Chain ID")
	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "workflow-owner-address", "d", DefaultWorkflowOwnerAddress, "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", DefaultWorkflowRegistryAddress, "Workflow registry address")

	return cmd
}

func compileWorkflow(workflowFilePathFlag, workflowNameFlag string) (string, error) {
	fmt.Printf("\n⚙️ Compiling workflow from %s\n", workflowFilePathFlag)

	compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(workflowFilePathFlag, workflowNameFlag)
	if compileErr != nil {
		return "", errors.Wrap(compileErr, "❌ failed to compile workflow")
	}

	fmt.Printf("\n✅ Workflow saved to %s\n\n", compressedWorkflowWasmPath)

	return compressedWorkflowWasmPath, nil
}

func deployWorkflow(ctx context.Context, wasmWorkflowFilePathFlag, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddressFlag, capabilitiesRegistryAddressFlag, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, rpcURLFlag string, donIDFlag uint32, deleteWorkflowFile bool) error {
	copyErr := creworkflow.CopyWorkflowToDockerContainers(wasmWorkflowFilePathFlag, containerNamePatternFlag, containerTargetDirFlag)
	if copyErr != nil {
		return errors.Wrap(copyErr, "❌ failed to copy workflow to Docker container")
	}

	fmt.Printf("\n✅ Workflow copied to Docker containers\n")
	fmt.Printf("\n⚙️ Creating Seth client\n\n")

	if pkErr := creenv.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
		return pkErr
	}

	sethClient, scErr := seth.NewClientBuilder().
		WithRpcUrl(rpcURLFlag).
		WithPrivateKeys([]string{os.Getenv("PRIVATE_KEY")}).
		WithProtections(false, false, seth.MustMakeDuration(time.Minute)).
		Build()
	if scErr != nil {
		return errors.Wrap(scErr, "failed to create Seth client")
	}

	var configPath *string
	if configFilePathFlag != "" {
		fmt.Printf("\n⚙️ Copying workflow config file to Docker container\n")
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

	var secretsPath *string
	if secretsFilePathFlag != "" {
		fmt.Printf("\n⚙️ Loading and encrypting workflow secrets\n")

		secretPathAbs, secretsErr := creworkflow.PrepareSecrets(sethClient, donIDFlag, common.HexToAddress(capabilitiesRegistryAddressFlag), common.HexToAddress(workflowOwnerAddressFlag), secretsFilePathFlag)
		if secretsErr != nil {
			return errors.Wrap(secretsErr, "failed to prepare secrets")
		}

		fmt.Printf("\n✅ Encrypted workflow secrets file prepared\n\n")

		fmt.Printf("\n⚙️ Copying encrypted secrets file to Docker container\n")
		secretsCopyErr := creworkflow.CopyWorkflowToDockerContainers(secretPathAbs, containerNamePatternFlag, containerTargetDirFlag)
		if secretsCopyErr != nil {
			return errors.Wrap(secretsCopyErr, "❌ failed to copy encrypted secrets file to Docker container")
		}

		secretPathAbs = "file://" + secretPathAbs
		secretsPath = &secretPathAbs

		fmt.Printf("\n✅ Encrypted workflow secrets file copied to Docker container\n\n")
	}

	fmt.Printf("\n⚙️ Deleting workflow '%s' from the workflow registry\n\n", workflowNameFlag)

	workflowNames, workflowNamesErr := creworkflow.GetWorkflowNames(ctx, sethClient, common.HexToAddress(workflowRegistryAddressFlag))
	if workflowNamesErr != nil {
		return errors.Wrap(workflowNamesErr, "failed to get workflows from the registry")
	}

	if !slices.Contains(workflowNames, workflowNameFlag) {
		fmt.Printf("\n✅ Workflow '%s' not found in the registry %s. Skipping...\n\n", workflowNameFlag, workflowRegistryAddressFlag)
	} else {
		deleteErr := creworkflow.DeleteWithContract(ctx, sethClient, common.HexToAddress(workflowRegistryAddressFlag), workflowNameFlag)
		if deleteErr != nil {
			return errors.Wrapf(deleteErr, "❌ failed to delete workflow '%s' from the registry %s", workflowNameFlag, workflowRegistryAddressFlag)
		}

		fmt.Printf("\n✅ Workflow '%s' deleted from the workflow registry\n\n", workflowNameFlag)
	}

	fmt.Printf("\n⚙️ Registering workflow '%s' with the workflow registry\n\n", workflowNameFlag)

	registerErr := creworkflow.RegisterWithContract(ctx, sethClient, common.HexToAddress(workflowRegistryAddressFlag), uint64(donIDFlag), workflowNameFlag, "file://"+wasmWorkflowFilePathFlag, configPath, secretsPath, &containerTargetDirFlag)
	if registerErr != nil {
		return errors.Wrapf(registerErr, "❌ failed to register workflow %s", workflowNameFlag)
	}

	if deleteWorkflowFile {
		defer func() {
			_ = os.Remove(wasmWorkflowFilePathFlag)
		}()
	}

	fmt.Printf("\n✅ Workflow registered successfully\n\n")

	return nil
}

func compileCopyAndRegisterWorkflow(ctx context.Context, workflowFilePathFlag, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddressFlag, capabilitiesRegistryAddressFlag, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, rpcURLFlag string, donIDFlag uint32) error {
	compressedWorkflowWasmPath, compileErr := compileWorkflow(workflowFilePathFlag, workflowNameFlag)
	if compileErr != nil {
		return errors.Wrap(compileErr, "❌ failed to compile workflow")
	}

	return deployWorkflow(ctx, compressedWorkflowWasmPath, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddressFlag, capabilitiesRegistryAddressFlag, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, rpcURLFlag, donIDFlag, true)
}
