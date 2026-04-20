package environment

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"

	vault_helpers "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	chainlinkvalues "github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	crevault "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/vault"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	vaulttypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

const (
	DefaultWorkflowOwnerAddress = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
)

// vaultSecretsConfig defines the structure of the vault secrets YAML file.
type vaultSecretsConfig struct {
	Secrets []vaultSecretEntry `yaml:"secrets"`
}

// vaultSecretEntry represents a single secret to be stored in the vault.
type vaultSecretEntry struct {
	Key       string `yaml:"key"`       // Vault key identifier
	EnvVar    string `yaml:"envVar"`    // Name of the env var containing the secret value
	Namespace string `yaml:"namespace"` // Vault namespace (defaults to "main" if empty)
}

// secretsNamesConfig is the secrets YAML format used by CRE:
//
//	secretsNames:
//	  SECRET_KEY:
//	    - ENV_VAR_NAME
type secretsNamesConfig struct {
	SecretsNames map[string][]string `yaml:"secretsNames"`
}

func workflowCmds() *cobra.Command {
	workflowCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Workflow management commands",
		Long:  `Commands to manage workflows`,
	}

	workflowCmd.AddCommand(deployAndVerifyExampleWorkflowCmd())
	workflowCmd.AddCommand(deleteWorkflowCmd())
	workflowCmd.AddCommand(deleteAllWorkflowsCmd())
	workflowCmd.AddCommand(compileWorkflowCmd())
	workflowCmd.AddCommand(deployWorkflowCmd())

	return workflowCmd
}

func deleteAllWorkflows(ctx context.Context, rpcURL, workflowRegistryAddress string, contractsVersion *semver.Version) error {
	if pkErr := environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
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

	deleteErr := creworkflow.DeleteAllWithContract(ctx, sethClient, common.HexToAddress(workflowRegistryAddress), contractsVersion)
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
			_, compileErr := compileWorkflow(cmd.Context(), workflowFilePathFlag, workflowNameFlag)
			if compileErr != nil {
				return errors.Wrap(compileErr, "❌ failed to compile workflow")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&workflowFilePathFlag, "workflow-file-path", "w", "", "Path to the workflow main Go file")
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
		gatewayURLFlag                  string
		deleteWorkflowFileFlag          bool
		donIDFlag                       uint32
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
			resolver, resolverErr := TryLoadLocalCREStateResolver()
			if resolverErr != nil {
				return errors.Wrap(resolverErr, "failed to load local CRE state")
			}

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
				compiledWorkflowPath, compileErr := compileWorkflow(cmd.Context(), workflowFilePathFlag, workflowNameFlag)
				if compileErr != nil {
					return errors.Wrap(compileErr, "❌ failed to compile workflow")
				}

				workflowFilePathFlag = compiledWorkflowPath
			}

			rpcURL := rpcURLFlag
			if !cmd.Flags().Changed("rpc-url") && resolver != nil {
				if stateRPC, err := resolver.RegistryRPC(); err == nil {
					rpcURL = stateRPC
				}
			}

			donID := donIDFlag
			if !cmd.Flags().Changed("don-id") && resolver != nil {
				if stateDONID, err := resolver.WorkflowDONID(); err == nil {
					donID = stateDONID
				}
			}

			gatewayURL := gatewayURLFlag
			if !cmd.Flags().Changed("gateway-url") && resolver != nil {
				if stateGatewayURL, err := resolver.GatewayURL(); err == nil {
					gatewayURL = stateGatewayURL
				}
			}

			workflowRegistryAddress, workflowRegistryVersion, resolveErr := resolveContractAddressAndVersion(cmd, resolver, keystone_changeset.WorkflowRegistry, workflowRegistryAddressFlag, contractsVersionFlag, "workflow-registry-address")
			if resolveErr != nil {
				return errors.Wrap(resolveErr, "❌ failed to resolve workflow registry")
			}

			capabilitiesRegistryAddress, capabilitiesRegistryVersion, resolveErr := resolveContractAddressAndVersion(cmd, resolver, keystone_changeset.CapabilitiesRegistry, capabilitiesRegistryAddressFlag, contractsVersionFlag, "capabilities-registry-address")
			if resolveErr != nil {
				return errors.Wrap(resolveErr, "❌ failed to resolve capabilities registry")
			}

			regErr = deployWorkflow(cmd.Context(), workflowFilePathFlag, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddress, capabilitiesRegistryAddress, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, secretsOutputFilePathFlag, rpcURL, gatewayURL, workflowRegistryVersion, capabilitiesRegistryVersion, donID, deleteWorkflowFileFlag)

			return regErr
		},
	}

	cmd.Flags().StringVarP(&workflowFilePathFlag, "workflow-file-path", "w", "", "Path to a base64-encoded workflow WASM file or to a Go file that contains the workflow (if --compile flag is used)")
	cmd.Flags().StringVarP(&configFilePathFlag, "config-file-path", "c", "", "Path to the workflow config file")
	cmd.Flags().StringVarP(&secretsFilePathFlag, "secrets-file-path", "s", "", "Path to the vault secrets YAML file (keys, env var names, namespaces)")
	cmd.Flags().StringVarP(&secretsOutputFilePathFlag, "secrets-output-file-path", "o", "", "Path to encrypted vault secrets output file (default \"./vault_secrets.json\")")
	cmd.Flags().StringVarP(&containerTargetDirFlag, "container-target-dir", "t", creworkflow.DefaultWorkflowTargetDir, "Path to the target directory in the Docker container")
	cmd.Flags().StringVarP(&containerNamePatternFlag, "container-name-pattern", "p", creworkflow.DefaultWorkflowNodePattern, "Pattern to match Docker containers workkflow DON containers (e.g. 'workflow-node')")
	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "workflow-owner-address", "d", DefaultWorkflowOwnerAddress, "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", "", "Workflow registry address (if not provided, address from the state file will be used)")
	cmd.Flags().StringVar(&capabilitiesRegistryAddressFlag, "capabilities-registry-address", "", "Capabilities registry address for vault config update (if not provided, address from the state file will be used)")
	cmd.Flags().StringVarP(&gatewayURLFlag, "gateway-url", "g", "", "Gateway URL for vault secrets (if not provided, URL from the state file will be used)")
	cmd.Flags().Uint32VarP(&donIDFlag, "don-id", "e", 1, "donID used in the workflow registry contract (integer starting with 1)")
	cmd.Flags().StringVarP(&workflowNameFlag, "name", "n", "", "Workflow name")
	cmd.Flags().BoolVarP(&deleteWorkflowFileFlag, "delete-workflow-file", "l", false, "Deletes the workflow file after deployment")
	cmd.Flags().BoolVarP(&compileWorkflowFlag, "compile", "x", false, "Compiles the workflow before deploying it")
	cmd.Flags().StringVar(&contractsVersionFlag, "with-contracts-version", "v2", "Version of workflow registry contract to use (v1 or v2)")

	if err := cmd.MarkFlagRequired("workflow-file-path"); err != nil {
		panic(err)
	}

	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	return cmd
}

func deleteWorkflowCmd() *cobra.Command {
	var (
		workflowNameFlag            string
		workflowRegistryAddressFlag string
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
			resolver, resolverErr := TryLoadLocalCREStateResolver()
			if resolverErr != nil {
				return errors.Wrap(resolverErr, "failed to load local CRE state")
			}

			rpcURL := rpcURLFlag
			if !cmd.Flags().Changed("rpc-url") && resolver != nil {
				if stateRPC, err := resolver.RegistryRPC(); err == nil {
					rpcURL = stateRPC
				}
			}

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

			workflowRegistryAddress, contractsVersion, err := resolveContractAddressAndVersion(cmd, resolver, keystone_changeset.WorkflowRegistry, workflowRegistryAddressFlag, contractsVersionFlag, "workflow-registry-address")
			if err != nil {
				return errors.Wrap(err, "❌ failed to resolve workflow registry")
			}

			workflowNames, workflowNamesErr := creworkflow.GetWorkflowNames(cmd.Context(), sethClient, common.HexToAddress(workflowRegistryAddress), contractsVersion)
			if workflowNamesErr != nil {
				return errors.Wrap(workflowNamesErr, "failed to get workflows from the registry")
			}

			if !slices.Contains(workflowNames, workflowNameFlag) {
				fmt.Printf("\n✅ Workflow '%s' not found in the registry %s. Skipping...\n\n", workflowNameFlag, workflowRegistryAddress)

				return nil
			}

			deleteErr := creworkflow.DeleteWithContract(cmd.Context(), sethClient, common.HexToAddress(workflowRegistryAddress), contractsVersion, workflowNameFlag)
			if deleteErr != nil {
				return errors.Wrapf(deleteErr, "❌ failed to delete workflow '%s' from the registry %s", workflowNameFlag, workflowRegistryAddress)
			}

			fmt.Printf("\n✅ Workflow deleted from the workflow registry\n\n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", "", "Workflow registry address (if not provided, address from the state file will be used)")
	cmd.Flags().StringVarP(&workflowNameFlag, "name", "n", "", "Workflow name")
	cmd.Flags().StringVar(&contractsVersionFlag, "with-contracts-version", "v2", "Version of workflow registry contract to use (v1 or v2)")

	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	return cmd
}

func deleteAllWorkflowsCmd() *cobra.Command {
	var (
		workflowRegistryAddressFlag string
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
			resolver, resolverErr := TryLoadLocalCREStateResolver()
			if resolverErr != nil {
				return errors.Wrap(resolverErr, "failed to load local CRE state")
			}

			rpcURL := rpcURLFlag
			if !cmd.Flags().Changed("rpc-url") && resolver != nil {
				if stateRPC, err := resolver.RegistryRPC(); err == nil {
					rpcURL = stateRPC
				}
			}

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

			workflowRegistryAddress, contractsVersion, err := resolveContractAddressAndVersion(cmd, resolver, keystone_changeset.WorkflowRegistry, workflowRegistryAddressFlag, contractsVersionFlag, "workflow-registry-address")
			if err != nil {
				return errors.Wrap(err, "❌ failed to resolve workflow registry")
			}

			deleteErr := creworkflow.DeleteAllWithContract(cmd.Context(), sethClient, common.HexToAddress(workflowRegistryAddress), contractsVersion)
			if deleteErr != nil {
				return errors.Wrapf(deleteErr, "❌ failed to delete all workflows from the registry %s", workflowRegistryAddress)
			}

			fmt.Printf("\n✅ All workflows deleted from the workflow registry\n\n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&rpcURLFlag, "rpc-url", "r", "http://localhost:8545", "RPC URL")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", "", "Workflow registry address (if not provided, address from the state file will be used)")
	cmd.Flags().StringVar(&contractsVersionFlag, "with-contracts-version", "v2", "Version of workflow registry contract to use (v1 or v2)")

	return cmd
}

func compileWorkflow(ctx context.Context, workflowFilePathFlag, workflowNameFlag string) (string, error) {
	fmt.Printf("\n⚙️ Compiling workflow from %s\n", workflowFilePathFlag)

	compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(ctx, workflowFilePathFlag, workflowNameFlag)
	if compileErr != nil {
		return "", errors.Wrap(compileErr, "❌ failed to compile workflow")
	}

	fmt.Printf("\n✅ Workflow saved to %s\n\n", compressedWorkflowWasmPath)

	return compressedWorkflowWasmPath, nil
}

func deployWorkflow(
	ctx context.Context,
	wasmWorkflowFilePathFlag, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddress, capabilitiesRegistryAddress, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, secretsOutputFilePathFlag, rpcURLFlag, gatewayURL string,
	workflowRegistryVersion, capabilitiesRegistryVersion *semver.Version,
	donIDFlag uint32,
	deleteWorkflowFile bool,
) error {
	copyErr := creworkflow.CopyArtifactsToDockerContainers(containerTargetDirFlag, containerNamePatternFlag, wasmWorkflowFilePathFlag)
	if copyErr != nil {
		return errors.Wrap(copyErr, "❌ failed to copy workflow to Docker container")
	}

	fmt.Printf("\n✅ Workflow copied to Docker containers\n")
	fmt.Printf("\n⚙️ Creating Seth client\n\n")

	if pkErr := environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
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

	var encryptedSecretsJSONPath string
	if secretsFilePathFlag != "" {
		if workflowRegistryVersion == nil || workflowRegistryVersion.Major() != 2 {
			return fmt.Errorf("❌ vault secrets flow requires v2 workflow registry contract, got %v", workflowRegistryVersion)
		}
		if capabilitiesRegistryVersion == nil || capabilitiesRegistryVersion.Major() != 2 {
			return fmt.Errorf("❌ vault secrets flow requires v2 capabilities registry contract, got %v", capabilitiesRegistryVersion)
		}

		if gatewayURL == "" {
			return errors.New("❌ --gateway-url (or a local CRE state file with gateway configuration) is required when --secrets-file-path is provided")
		}

		fmt.Printf("\n⚙️ Fetching vault public key from gateway\n")

		vaultPublicKey, vpkErr := fetchVaultPublicKey(ctx, gatewayURL)
		if vpkErr != nil {
			return errors.Wrap(vpkErr, "❌ failed to fetch vault public key from gateway")
		}

		fmt.Printf("\n✅ Vault public key fetched\n")

		if capabilitiesRegistryAddress != "" {
			fmt.Printf("\n⚙️ Updating vault capability config in capabilities registry\n")

			if updateErr := updateVaultCapabilityConfig(ctx, sethClient, capabilitiesRegistryAddress, vaultPublicKey); updateErr != nil {
				return errors.Wrap(updateErr, "❌ failed to update vault capability config in capabilities registry")
			}

			fmt.Printf("\n✅ Vault capability config updated\n")
			fmt.Printf("\n⚙️ Waiting for registry syncer to propagate vault config change\n")
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(15 * time.Second):
			}
		}

		fmt.Printf("\n⚙️ Encrypting workflow secrets for vault\n")

		ownerAddr := common.HexToAddress(workflowOwnerAddressFlag)
		encryptedPath, prepErr := prepareVaultSecrets(secretsFilePathFlag, vaultPublicKey, ownerAddr, secretsOutputFilePathFlag)
		if prepErr != nil {
			return errors.Wrap(prepErr, "❌ failed to prepare vault secrets")
		}

		encryptedSecretsJSONPath = encryptedPath

		fmt.Printf("\n✅ Vault secrets prepared at: %s\n\n", encryptedSecretsJSONPath)
	}

	fmt.Printf("\n⚙️ Deleting workflow '%s' from the workflow registry\n\n", workflowNameFlag)

	workflowNames, workflowNamesErr := creworkflow.GetWorkflowNames(ctx, sethClient, common.HexToAddress(workflowRegistryAddress), workflowRegistryVersion)
	if workflowNamesErr != nil {
		return errors.Wrap(workflowNamesErr, "failed to get workflows from the registry")
	}

	if !slices.Contains(workflowNames, workflowNameFlag) {
		fmt.Printf("\n✅ Workflow '%s' not found in the registry %s. Skipping...\n\n", workflowNameFlag, workflowRegistryAddress)
	} else {
		deleteErr := creworkflow.DeleteWithContract(ctx, sethClient, common.HexToAddress(workflowRegistryAddress), workflowRegistryVersion, workflowNameFlag)
		if deleteErr != nil {
			return errors.Wrapf(deleteErr, "❌ failed to delete workflow '%s' from the registry %s", workflowNameFlag, workflowRegistryAddress)
		}

		fmt.Printf("\n✅ Workflow '%s' deleted from the workflow registry\n\n", workflowNameFlag)
	}

	fmt.Printf("\n⚙️ Registering workflow '%s' with the workflow registry\n\n", workflowNameFlag)

	workflowID, registerErr := creworkflow.RegisterWithContract(ctx, sethClient, common.HexToAddress(workflowRegistryAddress), workflowRegistryVersion, uint64(donIDFlag), workflowNameFlag, "file://"+wasmWorkflowFilePathFlag, configPath, nil, &containerTargetDirFlag)
	if registerErr != nil {
		return errors.Wrapf(registerErr, "❌ failed to register workflow %s", workflowNameFlag)
	}

	if deleteWorkflowFile {
		defer func() {
			_ = os.Remove(wasmWorkflowFilePathFlag)
		}()
	}

	fmt.Printf("\n✅ Workflow registered successfully: workflowID='%s'\n\n", workflowID)

	if encryptedSecretsJSONPath != "" {
		fmt.Printf("\n⚙️ Sending encrypted secrets to vault via gateway\n\n")

		defer func() {
			_ = os.Remove(encryptedSecretsJSONPath)
		}()

		execErr := executeVaultSecrets(ctx, encryptedSecretsJSONPath, gatewayURL, sethClient, common.HexToAddress(workflowRegistryAddress))
		if execErr != nil {
			return errors.Wrap(execErr, "❌ failed to send secrets to vault gateway")
		}

		fmt.Printf("\n✅ Secrets sent to vault successfully\n\n")
	}

	return nil
}

func compileCopyAndRegisterWorkflow(ctx context.Context, workflowFilePathFlag, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddress, capabilitiesRegistryAddress, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, secretsOutputFilePathFlag, rpcURLFlag, gatewayURL string, workflowRegistryVersion, capabilitiesRegistryVersion *semver.Version, donIDFlag uint32) error {
	compressedWorkflowWasmPath, compileErr := compileWorkflow(ctx, workflowFilePathFlag, workflowNameFlag)
	if compileErr != nil {
		return errors.Wrap(compileErr, "❌ failed to compile workflow")
	}

	return deployWorkflow(ctx, compressedWorkflowWasmPath, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddress, capabilitiesRegistryAddress, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, secretsOutputFilePathFlag, rpcURLFlag, gatewayURL, workflowRegistryVersion, capabilitiesRegistryVersion, donIDFlag, true)
}

// updateVaultCapabilityConfig finds the DON that has the vault capability registered and injects
// the vault public key and a threshold of 1 into its DefaultConfig in the capabilities registry.
// This is required so that workflow nodes can unwrap the capability config when calling runtime.GetSecret().
func updateVaultCapabilityConfig(ctx context.Context, sethClient *seth.Client, capabilitiesRegistryAddr, vaultPublicKey string) error {
	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		common.HexToAddress(capabilitiesRegistryAddr), sethClient.Client,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create capabilities registry wrapper")
	}

	const pageSize int64 = 100

	var targetDON *capabilities_registry_v2.CapabilitiesRegistryDONInfo
	for start := int64(0); targetDON == nil; start += pageSize {
		donsPage, getErr := capReg.GetDONs(&bind.CallOpts{Context: ctx}, big.NewInt(start), big.NewInt(pageSize))
		if getErr != nil {
			return errors.Wrap(getErr, "failed to get DONs from capabilities registry")
		}

		for i := range donsPage {
			for _, cc := range donsPage[i].CapabilityConfigurations {
				if cc.CapabilityId == vault_helpers.CapabilityID {
					don := donsPage[i]
					targetDON = &don
					break
				}
			}
			if targetDON != nil {
				break
			}
		}

		if len(donsPage) < int(pageSize) {
			break
		}
	}
	if targetDON == nil {
		return fmt.Errorf("no DON with %s capability found in capabilities registry", vault_helpers.CapabilityID)
	}

	newConfigs := make([]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, 0, len(targetDON.CapabilityConfigurations))
	for _, cc := range targetDON.CapabilityConfigurations {
		if cc.CapabilityId == vault_helpers.CapabilityID {
			existingCfg := &capabilitiespb.CapabilityConfig{}
			if len(cc.Config) > 0 {
				if unmarshalErr := proto.Unmarshal(cc.Config, existingCfg); unmarshalErr != nil {
					return errors.Wrap(unmarshalErr, "failed to unmarshal existing vault capability config")
				}
			}

			valueMap, wrapErr := chainlinkvalues.WrapMap(map[string]interface{}{
				"VaultPublicKey": vaultPublicKey,
				"Threshold":      1,
			})
			if wrapErr != nil {
				return errors.Wrap(wrapErr, "failed to wrap vault capability config values")
			}

			existingCfg.DefaultConfig = chainlinkvalues.ProtoMap(valueMap)

			configBytes, marshalErr := proto.Marshal(existingCfg)
			if marshalErr != nil {
				return errors.Wrap(marshalErr, "failed to marshal updated vault capability config")
			}

			cc.Config = configBytes
		}
		newConfigs = append(newConfigs, cc)
	}

	updateParams := capabilities_registry_v2.CapabilitiesRegistryUpdateDONParams{
		Name:                     targetDON.Name,
		Config:                   targetDON.Config,
		CapabilityConfigurations: newConfigs,
		Nodes:                    targetDON.NodeP2PIds,
		F:                        targetDON.F,
		IsPublic:                 targetDON.IsPublic,
	}

	_, updateErr := sethClient.Decode(capReg.UpdateDONByName(sethClient.NewTXOpts(), targetDON.Name, updateParams))
	return errors.Wrap(updateErr, "UpdateDONByName tx failed")
}

// sendToVaultGateway sends an HTTP POST request to the vault gateway and returns the status code and body.
func sendToVaultGateway(ctx context.Context, gatewayURL string, body []byte) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, gatewayURL, bytes.NewReader(body))
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to build vault gateway request")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, errors.Wrap(err, "vault gateway HTTP request failed")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, errors.Wrap(err, "failed to read vault gateway response body")
	}

	return resp.StatusCode, respBody, nil
}

// fetchVaultPublicKey polls the vault gateway until it returns a public key.
func fetchVaultPublicKey(ctx context.Context, gatewayURL string) (string, error) {
	getPublicKeyRequest := jsonrpc.Request[vault_helpers.GetPublicKeyRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uuid.New().String(),
		Method:  vaulttypes.MethodPublicKeyGet,
		Params:  &vault_helpers.GetPublicKeyRequest{},
	}

	reqBody, err := json.Marshal(getPublicKeyRequest)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal public key request")
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	timeout := time.After(2 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-timeout:
			return "", errors.New("timed out waiting for vault public key from gateway")
		case <-ticker.C:
			statusCode, respBody, sendErr := sendToVaultGateway(ctx, gatewayURL, reqBody)
			if sendErr != nil || statusCode != http.StatusOK || respBody == nil {
				continue
			}

			var jsonResponse jsonrpc.Response[vault_helpers.GetPublicKeyResponse]
			if jsonErr := json.Unmarshal(respBody, &jsonResponse); jsonErr != nil {
				continue
			}

			if jsonResponse.Result.PublicKey != "" {
				return jsonResponse.Result.PublicKey, nil
			}
		}
	}
}

// prepareVaultSecrets reads the vault secrets YAML file, encrypts each secret using the vault
// public key, and writes the encrypted secrets list to a JSON file. The JSON file path is returned.
//
// Two YAML formats are accepted:
//
// Format 1 (explicit):
//
//	secrets:
//	  - key: "my-secret"
//	    envVar: "MY_SECRET_ENV_VAR"
//	    namespace: "main"   # optional, defaults to "main"
//
// Format 2 (secretsNames, shared with other CRE tools):
//
//	secretNames:
//	  SECRET_KEY:
//	    - ENV_VAR_NAME
func prepareVaultSecrets(secretsFilePath, vaultPublicKey string, ownerAddress common.Address, outputFilePath string) (string, error) {
	data, err := os.ReadFile(secretsFilePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to read secrets file")
	}

	var cfg vaultSecretsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "", errors.Wrap(err, "failed to parse secrets YAML file")
	}

	if len(cfg.Secrets) == 0 {
		// Try the alternative secretsNames format.
		var altCfg secretsNamesConfig
		if altErr := yaml.Unmarshal(data, &altCfg); altErr == nil {
			for key, envVars := range altCfg.SecretsNames {
				if len(envVars) == 0 {
					continue
				}
				cfg.Secrets = append(cfg.Secrets, vaultSecretEntry{
					Key:    key,
					EnvVar: envVars[0],
				})
			}
		}
	}

	if len(cfg.Secrets) == 0 {
		return "", errors.New("no secrets found in secrets file")
	}

	encryptedSecrets := make([]*vault_helpers.EncryptedSecret, 0, len(cfg.Secrets))
	for _, entry := range cfg.Secrets {
		value := os.Getenv(entry.EnvVar)
		if value == "" {
			return "", fmt.Errorf("environment variable %q is not set for secret key %q", entry.EnvVar, entry.Key)
		}

		namespace := entry.Namespace
		if namespace == "" {
			namespace = "main"
		}

		encryptedValue, encErr := crevault.EncryptSecret(value, vaultPublicKey, ownerAddress)
		if encErr != nil {
			return "", errors.Wrapf(encErr, "failed to encrypt secret %q", entry.Key)
		}

		encryptedSecrets = append(encryptedSecrets, &vault_helpers.EncryptedSecret{
			Id: &vault_helpers.SecretIdentifier{
				Key:       entry.Key,
				Owner:     ownerAddress.Hex(),
				Namespace: namespace,
			},
			EncryptedValue: encryptedValue,
		})
	}

	if outputFilePath == "" {
		outputFilePath = "./vault_secrets.json"
	}

	absPath, absErr := filepath.Abs(outputFilePath)
	if absErr != nil {
		return "", errors.Wrap(absErr, "failed to resolve absolute path for secrets output file")
	}

	jsonData, marshalErr := json.Marshal(encryptedSecrets)
	if marshalErr != nil {
		return "", errors.Wrap(marshalErr, "failed to marshal encrypted secrets to JSON")
	}

	if writeErr := os.WriteFile(absPath, jsonData, 0600); writeErr != nil {
		return "", errors.Wrap(writeErr, "failed to write encrypted secrets file")
	}

	return absPath, nil
}

// executeVaultSecrets reads the encrypted secrets JSON file produced by prepareVaultSecrets,
// allowlists the vault request in the workflow registry, and sends the secrets to the vault gateway.
func executeVaultSecrets(ctx context.Context, encryptedSecretsJSONPath, gatewayURL string, sethClient *seth.Client, workflowRegistryAddress common.Address) error {
	data, err := os.ReadFile(encryptedSecretsJSONPath)
	if err != nil {
		return errors.Wrap(err, "failed to read encrypted secrets file")
	}

	var encryptedSecrets []*vault_helpers.EncryptedSecret
	if err = json.Unmarshal(data, &encryptedSecrets); err != nil {
		return errors.Wrap(err, "failed to unmarshal encrypted secrets")
	}

	uniqueRequestID := uuid.New().String()
	createSecretsRequest := vault_helpers.CreateSecretsRequest{
		RequestId:        uniqueRequestID,
		EncryptedSecrets: encryptedSecrets,
	}

	requestBody, err := json.Marshal(&createSecretsRequest)
	if err != nil {
		return errors.Wrap(err, "failed to marshal create secrets request")
	}
	requestBodyJSON := json.RawMessage(requestBody)

	jsonRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsCreate,
		Params:  &requestBodyJSON,
	}

	requestDigest, err := jsonRequest.Digest()
	if err != nil {
		return errors.Wrap(err, "failed to compute request digest")
	}

	requestDigestBytes, err := hex.DecodeString(requestDigest)
	if err != nil {
		return errors.Wrap(err, "failed to decode request digest hex")
	}

	reqDigestBytes := [32]byte(requestDigestBytes)

	wfReg, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(workflowRegistryAddress, sethClient.Client)
	if err != nil {
		return errors.Wrap(err, "failed to instantiate workflow registry v2 wrapper")
	}

	expiry := uint32(time.Now().Add(time.Hour).Unix()) //nolint:gosec // G115: timestamp fits uint32 until year 2106
	_, decErr := sethClient.Decode(wfReg.AllowlistRequest(sethClient.NewTXOpts(), reqDigestBytes, expiry))
	if decErr != nil {
		return errors.Wrap(decErr, "failed to allowlist vault request in workflow registry")
	}

	fmt.Printf("\n✅ Vault request allowlisted in workflow registry\n")

	reqBody, err := json.Marshal(jsonRequest)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON-RPC request")
	}

	statusCode, respBody, sendErr := sendToVaultGateway(ctx, gatewayURL, reqBody)
	if sendErr != nil {
		return errors.Wrap(sendErr, "failed to send request to vault gateway")
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("vault gateway responded with status %d: %s", statusCode, string(respBody))
	}

	var jsonResponse jsonrpc.Response[json.RawMessage]
	if err := json.Unmarshal(respBody, &jsonResponse); err != nil {
		return errors.Wrap(err, "failed to unmarshal vault gateway response")
	}

	if jsonResponse.Error != nil && jsonResponse.Error.Error() != "" {
		return fmt.Errorf("vault gateway returned error: %s", jsonResponse.Error.Error())
	}

	return nil
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

func resolveContractAddressAndVersion(cmd *cobra.Command, resolver *LocalCREStateResolver, contractType deployment.ContractType, explicitAddress, versionFlag, addressFlagName string) (string, *semver.Version, error) {
	if cmd.Flags().Changed(addressFlagName) {
		if strings.TrimSpace(explicitAddress) == "" {
			return "", nil, fmt.Errorf("❌ %s is required when %s is provided", addressFlagName, addressFlagName)
		}

		if strings.TrimSpace(versionFlag) == "" {
			return "", nil, fmt.Errorf("❌ %s is required when %s is provided", versionFlag, addressFlagName)
		}

		version, err := semverFromFlag(versionFlag)
		if err != nil {
			return "", nil, err
		}

		return explicitAddress, version, nil
	}

	if resolver != nil {
		addrRef, err := resolver.AddressRef(contractType)
		if err != nil {
			return "", nil, err
		}

		return addrRef.Address, addrRef.Version, nil
	}

	if strings.TrimSpace(versionFlag) == "" {
		return "", nil, fmt.Errorf("❌ %s is required when no %s is provided", versionFlag, addressFlagName)
	}

	version, err := semverFromFlag(versionFlag)
	if err != nil {
		return "", nil, err
	}

	if strings.TrimSpace(explicitAddress) != "" {
		return explicitAddress, version, nil
	}

	return "", nil, fmt.Errorf("no %s available from flags or local CRE state", contractType)
}
