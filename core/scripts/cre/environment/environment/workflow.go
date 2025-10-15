package environment

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	creconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
)

const (
	// Might change if deployment sequence changes or if different config file than 'configs/workflow-don.toml' is used
	DefaultWorkflowRegistryAddress     = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
	DefaultCapabilitiesRegistryAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3"

	DefaultWorkflowOwnerAddress = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
)

// getWorkflowRegistryTypeVersion returns the appropriate TypeAndVersion based on the contracts version flag
func getWorkflowRegistryTypeVersion(contractsVersion string) deployment.TypeAndVersion {
	switch strings.ToLower(contractsVersion) {
	case "v2":
		return deployment.TypeAndVersion{
			Version: *semver.MustParse(creconfig.WorkflowRegistryV2Semver),
		}
	default:
		// Default to v1 for backward compatibility
		return deployment.TypeAndVersion{
			Version: *semver.MustParse("1.0.0"),
		}
	}
}

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

func deleteAllWorkflows(ctx context.Context, rpcURL, workflowRegistryAddress, contractsVersion string) error {
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

	workflowRegistryTypeVersion := getWorkflowRegistryTypeVersion(contractsVersion)
	deleteErr := creworkflow.DeleteAllWithContract(ctx, sethClient, common.HexToAddress(workflowRegistryAddress), workflowRegistryTypeVersion)
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
		Use:              "compile",
		Short:            "Compiles a workflow",
		Long:             `Compiles, compresses with Brotli and encodes with base64 a workflow`,
		PersistentPreRun: globalPreRunFunc,
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
		secretsOutputFilePathFlag       string
		compileWorkflowFlag             bool
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
		contractsVersionFlag            string
	)

	cmd := &cobra.Command{
		Use:              "deploy",
		Short:            "Deploys a workflow to the environment",
		Long:             `Deploys a workflow to the environment by copying it to workflow nodes and registering with the workflow registry`,
		PersistentPreRun: globalPreRunFunc,
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

				trackingErr := dxTracker.Track(MetricWorkflowDeploy, metaData)
				if trackingErr != nil {
					fmt.Fprintf(os.Stderr, "failed to track workflow deploy: %s\n", trackingErr)
				}
			}()

			if !compileWorkflowFlag {
				if err := isBase64File(workflowFilePathFlag); err != nil {
					return errors.Wrap(err, "❌ invalid WASM workflow file. Please make sure you're passing a base64-encoded and compiled workflow WASM file. If you want to compile and deploy a workflow, add '--compile' flag to the command instead")
				}
			}

			if compileWorkflowFlag {
				compiledWorkflowPath, compileErr := compileWorkflow(workflowFilePathFlag, workflowNameFlag)
				if compileErr != nil {
					return errors.Wrap(compileErr, "❌ failed to compile workflow")
				}

				workflowFilePathFlag = compiledWorkflowPath
			}

			regErr = deployWorkflow(cmd.Context(), workflowFilePathFlag, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddressFlag, capabilitiesRegistryAddressFlag, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, secretsOutputFilePathFlag, rpcURLFlag, contractsVersionFlag, donIDFlag, deleteWorkflowFileFlag)

			return regErr
		},
	}

	cmd.Flags().StringVarP(&workflowFilePathFlag, "workflow-file-path", "w", "", "Path to a base64-encoded workflow WASM file or to a Go file that contains the workflow (if --compile flag is used)")
	cmd.Flags().StringVarP(&configFilePathFlag, "config-file-path", "c", "", "Path to the config file")
	cmd.Flags().StringVarP(&secretsFilePathFlag, "secrets-file-path", "s", "", "Path to the secrets file")
	cmd.Flags().StringVarP(&secretsOutputFilePathFlag, "secrets-output-file-path", "o", "", "Path to encrypted secrets output file (default \"./encrypted.secrets.json\")")
	cmd.Flags().StringVarP(&containerTargetDirFlag, "container-target-dir", "t", creworkflow.DefaultWorkflowTargetDir, "Path to the target directory in the Docker container")
	cmd.Flags().StringVarP(&containerNamePatternFlag, "container-name-pattern", "p", creworkflow.DefaultWorkflowNodePattern, "Pattern to match the container name")
	cmd.Flags().Uint64VarP(&chainIDFlag, "chain-id", "i", 1337, "Chain ID")
	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "workflow-owner-address", "d", DefaultWorkflowOwnerAddress, "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", DefaultWorkflowRegistryAddress, "Workflow registry address")
	cmd.Flags().StringVarP(&capabilitiesRegistryAddressFlag, "capabilities-registry-address", "b", DefaultCapabilitiesRegistryAddress, "Capabilities registry address")
	cmd.Flags().Uint32VarP(&donIDFlag, "don-id", "e", 1, "DON ID")
	cmd.Flags().StringVarP(&workflowNameFlag, "name", "n", "", "Workflow name")
	cmd.Flags().BoolVarP(&deleteWorkflowFileFlag, "delete-workflow-file", "l", false, "Delete the workflow file after deployment")
	cmd.Flags().BoolVarP(&compileWorkflowFlag, "compile", "x", false, "Compile the workflow before deploying it")
	cmd.Flags().StringVar(&contractsVersionFlag, "with-contracts-version", "v1", "Version of workflow and capabilities registry contracts to use (v1 or v2)")

	if err := cmd.MarkFlagRequired("workflow-file-path"); err != nil {
		panic(err)
	}

	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	return cmd
}

func compileDeployWorkflowCmd() *cobra.Command {
	var (
		workflowFilePathFlag            string
		configFilePathFlag              string
		secretsFilePathFlag             string
		secretsOutputFilePathFlag       string
		containerTargetDirFlag          string
		containerNamePatternFlag        string
		workflowNameFlag                string
		workflowOwnerAddressFlag        string
		workflowRegistryAddressFlag     string
		capabilitiesRegistryAddressFlag string
		donIDFlag                       uint32
		chainIDFlag                     uint64
		rpcURLFlag                      string
		contractsVersionFlag            string
	)

	cmd := &cobra.Command{
		Use:              "compile-deploy",
		Short:            "DEPRECATED: Use 'go run . env workflow deploy --compile' instead",
		Long:             `DEPRECATED: Use 'go run . env workflow deploy --compile' instead`,
		Hidden:           true,
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("\n⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️\n\n")
			fmt.Printf("\033[31m'go run . env workflow compile-deploy' is DEPRECATED. Use 'go run . env workflow deploy --compile' instead\033[0m\n")
			fmt.Printf("\n⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️ ⚠️\n\n")

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

				trackingErr := dxTracker.Track(MetricWorkflowDeploy, metaData)
				if trackingErr != nil {
					fmt.Fprintf(os.Stderr, "failed to track workflow deploy: %s\n", trackingErr)
				}
			}()

			regErr = compileCopyAndRegisterWorkflow(cmd.Context(), workflowFilePathFlag, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddressFlag, capabilitiesRegistryAddressFlag, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, secretsOutputFilePathFlag, rpcURLFlag, contractsVersionFlag, donIDFlag)

			return regErr
		},
	}

	cmd.Flags().StringVarP(&workflowFilePathFlag, "workflow-file-path", "w", "./examples/workflows/v2/cron/main.go", "Path to the Go file that contains the workflow")
	cmd.Flags().StringVarP(&configFilePathFlag, "config-file-path", "c", "", "Path to the config file")
	cmd.Flags().StringVarP(&secretsFilePathFlag, "secrets-file-path", "s", "", "Path to the secrets file")
	cmd.Flags().StringVarP(&secretsOutputFilePathFlag, "secrets-output-file-path", "o", "", "Path to encrypted secrets output file (default \"./encrypted.secrets.json\")")
	cmd.Flags().StringVarP(&containerTargetDirFlag, "container-target-dir", "t", creworkflow.DefaultWorkflowTargetDir, "Path to the target directory in the Docker container")
	cmd.Flags().StringVarP(&containerNamePatternFlag, "container-name-pattern", "p", creworkflow.DefaultWorkflowNodePattern, "Pattern to match the container name")
	cmd.Flags().Uint64VarP(&chainIDFlag, "chain-id", "i", 1337, "Chain ID")
	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "owner-address", "d", DefaultWorkflowOwnerAddress, "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", DefaultWorkflowRegistryAddress, "Workflow registry address")
	cmd.Flags().StringVarP(&capabilitiesRegistryAddressFlag, "capabilities-registry-address", "b", DefaultCapabilitiesRegistryAddress, "Capabilities registry address")
	cmd.Flags().Uint32VarP(&donIDFlag, "don-id", "e", 1, "DON ID")
	cmd.Flags().StringVarP(&workflowNameFlag, "name", "n", "", "Workflow name")
	cmd.Flags().StringVar(&contractsVersionFlag, "with-contracts-version", "v1", "Version of workflow and capabilities registry contracts to use (v1 or v2)")

	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	return cmd
}

func deleteWorkflowCmd() *cobra.Command {
	var (
		workflowNameFlag            string
		workflowOwnerAddressFlag    string
		workflowRegistryAddressFlag string
		chainIDFlag                 uint64
		rpcURLFlag                  string
		contractsVersionFlag        string
	)

	cmd := &cobra.Command{
		Use:              "delete",
		Short:            "Deletes a workflow from the workflow registry contract",
		Long:             `Deletes a workflow from the workflow registry contract (but doesn't remove it from the Docker containers)`,
		PersistentPreRun: globalPreRunFunc,
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

			workflowRegistryTypeVersion := getWorkflowRegistryTypeVersion(contractsVersionFlag)
			workflowNames, workflowNamesErr := creworkflow.GetWorkflowNames(cmd.Context(), sethClient, common.HexToAddress(workflowRegistryAddressFlag), workflowRegistryTypeVersion)
			if workflowNamesErr != nil {
				return errors.Wrap(workflowNamesErr, "failed to get workflows from the registry")
			}

			if !slices.Contains(workflowNames, workflowNameFlag) {
				fmt.Printf("\n✅ Workflow '%s' not found in the registry %s. Skipping...\n\n", workflowNameFlag, workflowRegistryAddressFlag)

				return nil
			}

			deleteErr := creworkflow.DeleteWithContract(cmd.Context(), sethClient, common.HexToAddress(workflowRegistryAddressFlag), workflowRegistryTypeVersion, workflowNameFlag)
			if deleteErr != nil {
				return errors.Wrapf(deleteErr, "❌ failed to delete workflow '%s' from the registry %s", workflowNameFlag, workflowRegistryAddressFlag)
			}

			fmt.Printf("\n✅ Workflow deleted from the workflow registry\n\n")

			return nil
		},
	}

	cmd.Flags().Uint64VarP(&chainIDFlag, "chain-id", "i", 1337, "Chain ID")
	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "owner-address", "d", DefaultWorkflowOwnerAddress, "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", DefaultWorkflowRegistryAddress, "Workflow registry address")
	cmd.Flags().StringVarP(&workflowNameFlag, "name", "n", "", "Workflow name")
	cmd.Flags().StringVar(&contractsVersionFlag, "with-contracts-version", "v1", "Version of workflow and capabilities registry contracts to use (v1 or v2)")

	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	return cmd
}

func deleteAllWorkflowsCmd() *cobra.Command {
	var (
		workflowOwnerAddressFlag    string
		workflowRegistryAddressFlag string
		chainIDFlag                 uint64
		rpcURLFlag                  string
		contractsVersionFlag        string
	)

	cmd := &cobra.Command{
		Use:              "delete-all",
		Short:            "Deletes all workflows from the workflow registry contract",
		Long:             `Deletes all workflows from the workflow registry contract (but doesn't remove them from the Docker containers)`,
		PersistentPreRun: globalPreRunFunc,
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

			workflowRegistryTypeVersion := getWorkflowRegistryTypeVersion(contractsVersionFlag)
			deleteErr := creworkflow.DeleteAllWithContract(cmd.Context(), sethClient, common.HexToAddress(workflowRegistryAddressFlag), workflowRegistryTypeVersion)
			if deleteErr != nil {
				return errors.Wrapf(deleteErr, "❌ failed to delete all workflows from the registry %s", workflowRegistryAddressFlag)
			}

			fmt.Printf("\n✅ All workflows deleted from the workflow registry\n\n")

			return nil
		},
	}

	cmd.Flags().Uint64VarP(&chainIDFlag, "chain-id", "i", 1337, "Chain ID")
	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "owner-address", "d", DefaultWorkflowOwnerAddress, "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", DefaultWorkflowRegistryAddress, "Workflow registry address")
	cmd.Flags().StringVar(&contractsVersionFlag, "with-contracts-version", "v1", "Version of workflow and capabilities registry contracts to use (v1 or v2)")

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

func deployWorkflow(ctx context.Context, wasmWorkflowFilePathFlag, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddressFlag, capabilitiesRegistryAddressFlag, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, secretsOutputFilePathFlag, rpcURLFlag, contractsVersionFlag string, donIDFlag uint32, deleteWorkflowFile bool) error {
	copyErr := creworkflow.CopyArtifactsToDockerContainers(containerTargetDirFlag, containerNamePatternFlag, wasmWorkflowFilePathFlag)
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

		configCopyErr := creworkflow.CopyArtifactsToDockerContainers(containerTargetDirFlag, containerNamePatternFlag, configFilePathFlag)
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

		secretPathAbs, secretsErr := creworkflow.PrepareSecrets(sethClient, donIDFlag, common.HexToAddress(capabilitiesRegistryAddressFlag), common.HexToAddress(workflowOwnerAddressFlag), secretsFilePathFlag, secretsOutputFilePathFlag)
		if secretsErr != nil {
			return errors.Wrap(secretsErr, "failed to prepare secrets")
		}

		defer func() {
			_ = os.Remove(secretPathAbs)
		}()

		fmt.Printf("\n✅ Encrypted workflow secrets file created at: %s\n\n", secretPathAbs)

		fmt.Printf("\n⚙️ Copying encrypted secrets file to Docker container\n")
		secretsCopyErr := creworkflow.CopyArtifactsToDockerContainers(containerTargetDirFlag, containerNamePatternFlag, secretPathAbs)
		if secretsCopyErr != nil {
			return errors.Wrap(secretsCopyErr, "❌ failed to copy encrypted secrets file to Docker container")
		}

		secretPathAbs = "file://" + secretPathAbs
		secretsPath = &secretPathAbs

		fmt.Printf("\n✅ Encrypted workflow secrets file copied to Docker container\n\n")
	}

	fmt.Printf("\n⚙️ Deleting workflow '%s' from the workflow registry\n\n", workflowNameFlag)

	workflowRegistryTypeVersion := getWorkflowRegistryTypeVersion(contractsVersionFlag)
	workflowNames, workflowNamesErr := creworkflow.GetWorkflowNames(ctx, sethClient, common.HexToAddress(workflowRegistryAddressFlag), workflowRegistryTypeVersion)
	if workflowNamesErr != nil {
		return errors.Wrap(workflowNamesErr, "failed to get workflows from the registry")
	}

	if !slices.Contains(workflowNames, workflowNameFlag) {
		fmt.Printf("\n✅ Workflow '%s' not found in the registry %s. Skipping...\n\n", workflowNameFlag, workflowRegistryAddressFlag)
	} else {
		deleteErr := creworkflow.DeleteWithContract(ctx, sethClient, common.HexToAddress(workflowRegistryAddressFlag), workflowRegistryTypeVersion, workflowNameFlag)
		if deleteErr != nil {
			return errors.Wrapf(deleteErr, "❌ failed to delete workflow '%s' from the registry %s", workflowNameFlag, workflowRegistryAddressFlag)
		}

		fmt.Printf("\n✅ Workflow '%s' deleted from the workflow registry\n\n", workflowNameFlag)
	}

	fmt.Printf("\n⚙️ Registering workflow '%s' with the workflow registry\n\n", workflowNameFlag)

	workflowID, registerErr := creworkflow.RegisterWithContract(ctx, sethClient, common.HexToAddress(workflowRegistryAddressFlag), workflowRegistryTypeVersion, uint64(donIDFlag), workflowNameFlag, "file://"+wasmWorkflowFilePathFlag, configPath, secretsPath, &containerTargetDirFlag)
	if registerErr != nil {
		return errors.Wrapf(registerErr, "❌ failed to register workflow %s", workflowNameFlag)
	}

	if deleteWorkflowFile {
		defer func() {
			_ = os.Remove(wasmWorkflowFilePathFlag)
		}()
	}

	fmt.Printf("\n✅ Workflow registered successfully: workflowID='%s'\n\n", workflowID)

	return nil
}

func compileCopyAndRegisterWorkflow(ctx context.Context, workflowFilePathFlag, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddressFlag, capabilitiesRegistryAddressFlag, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, secretsOutputFilePathFlag, rpcURLFlag, contractsVersionFlag string, donIDFlag uint32) error {
	compressedWorkflowWasmPath, compileErr := compileWorkflow(workflowFilePathFlag, workflowNameFlag)
	if compileErr != nil {
		return errors.Wrap(compileErr, "❌ failed to compile workflow")
	}

	return deployWorkflow(ctx, compressedWorkflowWasmPath, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddressFlag, capabilitiesRegistryAddressFlag, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, secretsOutputFilePathFlag, rpcURLFlag, contractsVersionFlag, donIDFlag, true)
}

func isBase64File(filename string) error {
	fileInfo, fErr := os.Stat(filename)
	if fErr != nil {
		return errors.Wrap(fErr, "failed to get file info")
	}

	readSize := min(fileInfo.Size(), 4*1024*1024) // 4MB

	file, oErr := os.Open(filename)
	if oErr != nil {
		return errors.Wrap(oErr, "failed to open file")
	}
	defer file.Close()

	buffer := make([]byte, readSize)
	n, rErr := file.Read(buffer)
	if rErr != nil && rErr != io.EOF {
		return errors.Wrap(rErr, "failed to read file")
	}

	if !isBase64Content(string(buffer[:n])) {
		return fmt.Errorf("❌ file %s is not a base64-encoded file", filename)
	}

	return nil
}

func isBase64Content(content string) bool {
	// Remove whitespace and newlines, just to be safe
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, " ", "")
	content = strings.ReplaceAll(content, "\t", "")

	if len(content) == 0 {
		return false
	}

	_, err := base64.StdEncoding.DecodeString(content)
	return err == nil
}
