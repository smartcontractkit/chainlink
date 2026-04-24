package environment

import (
	"context"
	"fmt"
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

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
)

const (
	DefaultWorkflowOwnerAddress = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
)

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
	sethClient, err := newSethClient(rpcURL)
	if err != nil {
		return err
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
				if err := creworkflow.IsBase64File(workflowFilePathFlag); err != nil {
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

			rpcURL := resolveRPCURL(cmd, rpcURLFlag, resolver)

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

			nodeDBPort, nodeCount := resolveWorkflowDONNodeInfo(resolver)

			regErr = deployWorkflow(cmd.Context(), workflowFilePathFlag, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddress, capabilitiesRegistryAddress, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, secretsOutputFilePathFlag, rpcURL, gatewayURL, workflowRegistryVersion, capabilitiesRegistryVersion, donID, deleteWorkflowFileFlag, nodeDBPort, nodeCount)

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

			rpcURL := resolveRPCURL(cmd, rpcURLFlag, resolver)

			sethClient, err := newSethClient(rpcURL)
			if err != nil {
				return err
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

			rpcURL := resolveRPCURL(cmd, rpcURLFlag, resolver)

			sethClient, err := newSethClient(rpcURL)
			if err != nil {
				return err
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
	nodeDBPort, nodeCount int,
) error {
	copyErr := creworkflow.CopyArtifactsToDockerContainers(containerTargetDirFlag, containerNamePatternFlag, wasmWorkflowFilePathFlag)
	if copyErr != nil {
		return errors.Wrap(copyErr, "❌ failed to copy workflow to Docker container")
	}

	fmt.Printf("\n✅ Workflow copied to Docker containers\n")
	fmt.Printf("\n⚙️ Creating Seth client\n\n")

	sethClient, err := newSethClient(rpcURLFlag)
	if err != nil {
		return err
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

		vaultPublicKey, vpkErr := creworkflow.FetchVaultPublicKey(ctx, gatewayURL)
		if vpkErr != nil {
			return errors.Wrap(vpkErr, "❌ failed to fetch vault public key from gateway")
		}

		fmt.Printf("\n✅ Vault public key fetched\n")

		fmt.Printf("\n⚙️ Checking vault capability config in capabilities registry\n")

		vaultDON, existingVaultCfg, getErr := cre.GetVaultCapabilityDON(ctx, sethClient, capabilitiesRegistryAddress)
		if getErr != nil {
			return errors.Wrap(getErr, "❌ failed to get vault capability config from capabilities registry")
		}

		if !creworkflow.VaultConfigHasPublicKey(existingVaultCfg, vaultPublicKey) {
			fmt.Printf("\n⚙️ Updating vault capability config in capabilities registry\n")

			if updateErr := creworkflow.UpdateVaultCapabilityConfig(ctx, sethClient, capabilitiesRegistryAddress, vaultDON, vaultPublicKey, 1); updateErr != nil {
				return errors.Wrap(updateErr, "❌ failed to update vault capability config in capabilities registry")
			}

			fmt.Printf("\n✅ Vault capability config updated\n")
			fmt.Printf("\n⚙️ Waiting for registry syncer to propagate vault config change\n")
			if waitErr := creworkflow.WaitForVaultConfigPropagation(ctx, nodeDBPort, nodeCount); waitErr != nil {
				return errors.Wrap(waitErr, "❌ failed while waiting for vault config propagation")
			}
		} else {
			fmt.Printf("\n✅ Vault public key already configured, skipping update and propagation wait\n")
		}

		fmt.Printf("\n⚙️ Encrypting workflow secrets for vault\n")

		ownerAddr := common.HexToAddress(workflowOwnerAddressFlag)
		encryptedPath, prepErr := creworkflow.PrepareSecrets(secretsFilePathFlag, vaultPublicKey, ownerAddr, secretsOutputFilePathFlag)
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

	workflowID, registerErr := creworkflow.RegisterWithContract(ctx, sethClient, common.HexToAddress(workflowRegistryAddress), workflowRegistryVersion, uint64(donIDFlag), workflowNameFlag, "file://"+wasmWorkflowFilePathFlag, configPath, nil, nil, &containerTargetDirFlag)
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

		execErr := creworkflow.ExecuteSecrets(ctx, encryptedSecretsJSONPath, gatewayURL, sethClient, common.HexToAddress(workflowRegistryAddress))
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

	resolver, resolverErr := TryLoadLocalCREStateResolver()
	if resolverErr != nil {
		return errors.Wrap(resolverErr, "failed to load local CRE state")
	}
	nodeDBPort, nodeCount := resolveWorkflowDONNodeInfo(resolver)

	return deployWorkflow(ctx, compressedWorkflowWasmPath, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddress, capabilitiesRegistryAddress, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, secretsOutputFilePathFlag, rpcURLFlag, gatewayURL, workflowRegistryVersion, capabilitiesRegistryVersion, donIDFlag, true, nodeDBPort, nodeCount)
}

// newSethClient creates a Seth client for rpcURL, ensuring PRIVATE_KEY is set in the
// environment and falling back to blockchain.DefaultAnvilPrivateKey when it is not.
func newSethClient(rpcURL string) (*seth.Client, error) {
	if err := environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); err != nil {
		return nil, err
	}
	client, err := seth.NewClientBuilder().
		WithRpcUrl(rpcURL).
		WithPrivateKeys([]string{os.Getenv("PRIVATE_KEY")}).
		WithProtections(false, false, seth.MustMakeDuration(time.Minute)).
		Build()
	return client, errors.Wrap(err, "failed to create Seth client")
}

// resolveRPCURL returns the RPC URL from the state resolver when the rpc-url flag was not
// explicitly provided on the command line; otherwise it returns the flag value.
func resolveRPCURL(cmd *cobra.Command, flagValue string, resolver *LocalCREStateResolver) string {
	if !cmd.Flags().Changed("rpc-url") && resolver != nil {
		if stateRPC, err := resolver.RegistryRPC(); err == nil {
			return stateRPC
		}
	}
	return flagValue
}

// resolveWorkflowDONNodeInfo returns the workflow DON's shared DB port and worker count
// from the resolver. Returns (0, 0) if the resolver is nil or node info is unavailable
// (non-fatal: callers fall back to a static wait when these are zero).
func resolveWorkflowDONNodeInfo(resolver *LocalCREStateResolver) (dbPort, nodeCount int) {
	if resolver == nil {
		return 0, 0
	}
	dbPort, nodeCount, _ = resolver.WorkflowDONNodeInfo()
	return dbPort, nodeCount
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
