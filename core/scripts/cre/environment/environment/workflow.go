package environment

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	secretsUtils "github.com/smartcontractkit/chainlink-common/pkg/workflows/secrets"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
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
	workflowCmd.AddCommand(deleteWorkflowCmd())
	workflowCmd.AddCommand(deleteAllWorkflowsCmd())

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
		donIDFlag                       uint32
		chainIDFlag                     uint64
		rpcURLFlag                      string
	)

	cmd := &cobra.Command{
		Use:   "deploy",
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
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "workflow-owner-address", "d", "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0", "Workflow registry address")
	cmd.Flags().StringVarP(&capabilitiesRegistryAddressFlag, "capabilities-registry-address", "b", "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512", "Capabilities registry address")
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
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "workflow-owner-address", "d", "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0", "Workflow registry address")
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
	cmd.Flags().StringVarP(&workflowOwnerAddressFlag, "workflow-owner-address", "d", "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "Workflow owner address")
	cmd.Flags().StringVarP(&workflowRegistryAddressFlag, "workflow-registry-address", "a", "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0", "Workflow registry address")

	return cmd
}

func compileCopyAndRegisterWorkflow(ctx context.Context, workflowFilePathFlag, workflowNameFlag, workflowOwnerAddressFlag, workflowRegistryAddressFlag, capabilitiesRegistryAddressFlag, containerNamePatternFlag, containerTargetDirFlag, configFilePathFlag, secretsFilePathFlag, rpcURLFlag string, donIDFlag uint32) error {
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
		fmt.Printf("\n⚙️ Loading workflow secrets\n")

		secretsConfig, err := newSecretsConfig(secretsFilePathFlag)
		if err != nil {
			return err
		}

		envSecrets, err := loadSecretsFromEnvironment(secretsConfig)
		if err != nil {
			return err
		}

		fmt.Printf("\n✅ Loaded workflow secrets\n\n")

		fmt.Printf("\n⚙️ Encrypting workflow secrets\n")

		encryptSecrets, err := encryptSecrets(sethClient, common.HexToAddress(capabilitiesRegistryAddressFlag), donIDFlag, workflowOwnerAddressFlag, envSecrets, secretsConfig)
		if err != nil {
			return err
		}

		fmt.Printf("\n✅ Encrypted workflow secrets\n\n")

		fmt.Printf("\n⚙️ Writing encrypted secrets file to disk\n")

		encryptedSecretsFilePath := "./encrypted.secrets.json"
		encryptedSecretsFile, err := os.Create(encryptedSecretsFilePath)
		if err != nil {
			return fmt.Errorf("failed to create secrets file: %w", err)
		}
		defer encryptedSecretsFile.Close()
		defer func() {
			_ = os.Remove(encryptedSecretsFilePath)
		}()

		encoder := json.NewEncoder(encryptedSecretsFile)
		if err := encoder.Encode(encryptSecrets); err != nil {
			return fmt.Errorf("failed to write to secrets file: %w", err)
		}

		fmt.Printf("\n✅ Wrote encrypted secrets file to disk\n\n")

		fmt.Printf("\n⚙️ Copying encrypted secrets file to Docker container\n")
		secretPathAbs, secretPathAbsErr := filepath.Abs(encryptedSecretsFilePath)
		if secretPathAbsErr != nil {
			return errors.Wrap(secretPathAbsErr, "failed to get absolute path of the encrypted secrets file")
		}

		secretsCopyErr := creworkflow.CopyWorkflowToDockerContainers(encryptedSecretsFilePath, containerNamePatternFlag, containerTargetDirFlag)
		if secretsCopyErr != nil {
			return errors.Wrap(secretsCopyErr, "❌ failed to copy encrypted secrets file to Docker container")
		}

		secretPathAbs = "file://" + secretPathAbs
		secretsPath = &secretPathAbs

		fmt.Printf("\n✅ Workflow encrypted secrets file copied to Docker container\n\n")
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

	registerErr := creworkflow.RegisterWithContract(ctx, sethClient, common.HexToAddress(workflowRegistryAddressFlag), uint64(donIDFlag), workflowNameFlag, "file://"+compressedWorkflowWasmPath, configPath, secretsPath, &containerTargetDirFlag)
	if registerErr != nil {
		return errors.Wrapf(registerErr, "❌ failed to register workflow %s", workflowNameFlag)
	}

	defer func() {
		_ = os.Remove(compressedWorkflowWasmPath)
	}()

	fmt.Printf("\n✅ Workflow registered successfully\n\n")

	return nil
}

func newSecretsConfig(configPath string) (*secretsUtils.SecretsConfig, error) {
	secretsConfigFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening secrets config file: %w", err)
	}
	defer secretsConfigFile.Close()

	var config secretsUtils.SecretsConfig
	err = yaml.NewDecoder(secretsConfigFile).Decode(&config)
	if err != nil && errors.Is(err, io.EOF) {
		return &secretsUtils.SecretsConfig{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error decoding secrets config file: %w", err)
	}

	return &config, nil
}

func loadSecretsFromEnvironment(config *secretsUtils.SecretsConfig) (map[string][]string, error) {
	secrets := make(map[string][]string)
	for secretName, envVars := range config.SecretsNames {
		for _, envVar := range envVars {
			secretValue := os.Getenv(envVar)
			if secretValue == "" {
				return nil, fmt.Errorf("missing environment variable: %s", envVar)
			}
			secrets[secretName] = append(secrets[secretName], secretValue)
		}
	}
	return secrets, nil
}

func encryptSecrets(c *seth.Client, capabilitiesRegistry common.Address, donID uint32, workflowOwner string, secrets map[string][]string, config *secretsUtils.SecretsConfig) (secretsUtils.EncryptedSecretsResult, error) {
	cr, err := capabilities_registry.NewCapabilitiesRegistry(capabilitiesRegistry, c.Client)
	if err != nil {
		return secretsUtils.EncryptedSecretsResult{}, fmt.Errorf("failed to attach to the Capabilities Registry contract: %w", err)
	}

	nodeInfos, err := cr.GetNodes(c.NewCallOpts())
	if err != nil {
		return secretsUtils.EncryptedSecretsResult{}, fmt.Errorf("failed to get node information from the Capabilities Registry contract: %w", err)
	}

	donInfo, err := cr.GetDON(c.NewCallOpts(), donID)
	if err != nil {
		return secretsUtils.EncryptedSecretsResult{}, fmt.Errorf("failed to get DON information from the Capabilities Registry contract: %w", err)
	}

	encryptionPublicKeys := make(map[string][32]byte)
	for _, nodeInfo := range nodeInfos {
		// Filter only the nodes that are part of the DON
		if secretsUtils.ContainsP2pId(nodeInfo.P2pId, donInfo.NodeP2PIds) {
			encryptionPublicKeys[hex.EncodeToString(nodeInfo.P2pId[:])] = nodeInfo.EncryptionPublicKey
		}
	}

	if len(encryptionPublicKeys) == 0 {
		return secretsUtils.EncryptedSecretsResult{}, errors.New("no nodes found for the don")
	}

	// Encrypt secrets for each node
	encryptedSecrets, secretsEnvVarsByNode, err := secretsUtils.EncryptSecretsForNodes(
		workflowOwner,
		secrets,
		encryptionPublicKeys,
		secretsUtils.SecretsConfig{SecretsNames: config.SecretsNames},
	)
	if err != nil {
		return secretsUtils.EncryptedSecretsResult{}, fmt.Errorf("node public keys not found: %w", err)
	}

	// Convert encryptionPublicKey to hex strings for including in the metadata
	nodePublicEncryptionKeys := make(map[string]string)
	for p2pID, encryptionPublicKey := range encryptionPublicKeys {
		nodePublicEncryptionKeys[p2pID] = hex.EncodeToString(encryptionPublicKey[:])
	}

	result := secretsUtils.EncryptedSecretsResult{
		EncryptedSecrets: encryptedSecrets,
		Metadata: secretsUtils.Metadata{
			WorkflowOwner:            workflowOwner,
			CapabilitiesRegistry:     capabilitiesRegistry.String(),
			DonId:                    strconv.FormatUint(uint64(donID), 10),
			DateEncrypted:            time.Now().Format(time.RFC3339),
			NodePublicEncryptionKeys: nodePublicEncryptionKeys,
			EnvVarsAssignedToNodes:   secretsEnvVarsByNode,
		},
	}
	return result, nil
}
