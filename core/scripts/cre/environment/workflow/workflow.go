package workflow

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	workflow_registry_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
)

var (
	configPath      string
	contractAddress string
	workflowName    string
	workflowID      string
	workflowOwner   string
	binaryURL       string
	configURL       string
	secretsURL      string
	initialStatus   uint8
	donID           uint32
)

type RegisterWorkflowParameters struct {
	WorkflowName  string   // required: user specified human readable workflow name label
	WorkflowID    [32]byte // required: generated based on the workflow content and owner address
	BinaryURL     string   // required: URL location for the workflow binary WASM file
	ConfigURL     string   // optional: URL location for the workflow configuration file (default empty string)
	SecretsURL    string   // optional: URL location for the workflow secrets encrypted file (default empty string)
	InitialStatus uint8    // optional: 1 to pause workflow after registration, 0 to activate it (default is 0)
	DonID         uint32   // optional: DON where the workflow will run (default is specified in the DefaultDonId constant)
}

var WorkflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Manage workflows on the blockchain",
	Long:  `Commands to register, update, activate, pause, and query workflows on the Workflow Registry contract`,
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new workflow",
	Long:  `Register a new workflow on the Workflow Registry smart contract`,
	Run: func(cmd *cobra.Command, args []string) {
		// Setup client
		client, seth, err := setupWorkflowClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to setup workflow client: %v\n", err)
			return
		}

		// Parse workflowID to bytes32
		var workflowIDBytes [32]byte
		copy(workflowIDBytes[:], common.FromHex(workflowID))

		txOpt, err := bind.NewKeyedTransactorWithChainID(seth.MustGetRootPrivateKey(), big.NewInt(int64(chain_selectors.GETH_TESTNET.EvmChainID))) // TODO: @george-dorin get chainid from cre.yaml
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create transactor: %v\n", err)
		}
		// Register workflow
		output, err := client.RegisterWorkflow(txOpt, workflowName, workflowIDBytes, donID, initialStatus, binaryURL, configURL, secretsURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to register workflow: %v\n", err)
			return
		}

		fmt.Printf("Workflow '%s' successfully registered\n", workflowName)
		fmt.Printf("Transaction hash: %s\n", output.Hash().Hex())
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing workflow",
	Long:  `Update the binary, config, or secrets URL of an existing workflow`,
	Run: func(cmd *cobra.Command, args []string) {
		// Setup client
		client, _, err := setupWorkflowClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to setup workflow client: %v\n", err)
			return
		}

		// Parse workflowID and owner
		var workflowIDBytes [32]byte
		copy(workflowIDBytes[:], common.FromHex(workflowID))

		// Update workflow
		output, err := client.UpdateWorkflow(&bind.TransactOpts{}, ComputeHashKey(common.HexToAddress(workflowOwner), workflowName), workflowIDBytes, binaryURL, configURL, secretsURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to update workflow: %v\n", err)
			return
		}

		fmt.Printf("Workflow '%s' successfully updated\n", workflowName)
		fmt.Printf("Transaction hash: %s\n", output.Hash().Hex())
	},
}

var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Activate a workflow",
	Long:  `Activate a paused workflow on the Workflow Registry contract`,
	Run: func(cmd *cobra.Command, args []string) {
		// // Setup client
		// client, _, err := setupWorkflowClient()
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to setup workflow client: %v\n", err)
		// 	return
		// }
		//
		// // Activate workflow
		// _, err = client.ActivateWorkflow(&bind.TransactOpts{}, workflowID)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to activate workflow: %v\n", err)
		// 	return
		// }
		//
		// fmt.Printf("Workflow '%s' successfully activated\n", workflowName)
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause a workflow",
	Long:  `Pause an active workflow on the Workflow Registry contract`,
	Run: func(cmd *cobra.Command, args []string) {
		// // Setup client
		// client, err := setupWorkflowClient()
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to setup workflow client: %v\n", err)
		// 	return
		// }
		//
		// // Parse owner
		// owner := common.HexToAddress(workflowOwner)
		//
		// // Pause workflow
		// err = client.PauseWorkflow(workflowName, owner)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to pause workflow: %v\n", err)
		// 	return
		// }
		//
		// fmt.Printf("Workflow '%s' successfully paused\n", workflowName)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a workflow",
	Long:  `Delete a workflow from the Workflow Registry contract`,
	Run: func(cmd *cobra.Command, args []string) {
		// Setup client
		client, seth, err := setupWorkflowClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to setup workflow client: %v\n", err)
			return
		}

		// Delete workflow
		txOpt, err := bind.NewKeyedTransactorWithChainID(seth.MustGetRootPrivateKey(), big.NewInt(int64(chain_selectors.GETH_TESTNET.EvmChainID))) // TODO: @george-dorin get chainid from cre.yaml
		txHash, err := client.DeleteWorkflow(txOpt, ComputeHashKey(common.HexToAddress(workflowOwner), workflowName))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to delete workflow: %v\n", err)
			return
		}

		fmt.Printf("Workflow '%s' successfully deleted\n", workflowName)
		fmt.Printf("Transaction hash: %s\n", txHash.Hash().Hex())
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get workflow details",
	Long:  `Get details of a specific workflow from the Workflow Registry contract`,
	Run: func(cmd *cobra.Command, args []string) {
		// Setup client
		client, seth, err := setupWorkflowClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to setup workflow client: %v\n", err)
			return
		}
		// Parse owner
		owner := common.HexToAddress(workflowOwner)

		// Get workflow details
		metadata, err := client.GetWorkflowMetadataListByOwner(seth.NewCallOpts(), owner, big.NewInt(0), big.NewInt(9999999))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get workflow: %v\n", err)
			return
		}

		// Display workflow information
		for _, m := range metadata {
			fmt.Print("### Workflow ###\n")
			fmt.Printf("Name: %s\n", m.WorkflowName)
			fmt.Printf("Owner: %s\n", m.Owner.Hex())
			fmt.Printf("ID: 0x%x\n", m.WorkflowID)
			fmt.Printf("DON ID: %d\n", m.DonID)
			fmt.Printf("Status: %d\n", m.Status)
			fmt.Printf("Binary URL: %s\n", m.BinaryURL)
			fmt.Printf("Config URL: %s\n", m.ConfigURL)
			fmt.Printf("Secrets URL: %s\n", m.SecretsURL)
		}

	},
}

// setupWorkflowClient creates a workflow registry client using RPC information from cre.yaml
func setupWorkflowClient() (*workflow_registry_wrapper.WorkflowRegistry, *seth.Client, error) {
	// Load and parse the cre.yaml file
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading config file: %v", err)
	}
	config := crecli.Profiles{}

	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, nil, fmt.Errorf("error parsing config file: %v", err)
	}

	// Use contract address from config if not provided via flag
	for _, v := range config.Test.Contracts.ContractRegistry {
		if v.Name == "WorkflowRegistry" {
			contractAddress = v.Address
			fmt.Printf("Using contract address from config: %s\n", contractAddress)
			break
		}
	}

	// Use owner from config if not provided via flag
	if workflowOwner == "" && config.Test.UserWorkflow.WorkflowOwnerAddress != "" {
		workflowOwner = config.Test.UserWorkflow.WorkflowOwnerAddress
		fmt.Printf("Using owner address from config: %s\n", workflowOwner)
	}

	// Validate we have at least one RPC URL
	if len(config.Test.Rpcs) == 0 {
		return nil, nil, fmt.Errorf("no RPC URLs found in config file")
	}

	// Use the first RPC URL by default
	rpcURL := config.Test.Rpcs[0].URL
	privateKey := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" // Default key

	// Create Ethereum client
	ethClient, err := seth.NewClientBuilder().
		WithRpcUrl(rpcURL).
		WithPrivateKeys([]string{privateKey}).
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Ethereum client: %v", err)
	}

	if contractAddress == "" {
		return nil, nil, fmt.Errorf("workflow registry contract address not provided in flags or config")
	}

	// Create workflow registry client
	workflowClient, err := workflow_registry_wrapper.NewWorkflowRegistry(
		common.HexToAddress(contractAddress),
		ethClient.Client,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create workflow registry client: %v", err)
	}

	return workflowClient, ethClient, nil
}

func init() {
	// Add config path flag first so it's loaded before other flags are processed
	WorkflowCmd.PersistentFlags().StringVar(&configPath, "config", "cre.yaml", "Path to cre.yaml config file")

	// Add common flags - now optional since they can come from config
	WorkflowCmd.PersistentFlags().StringVar(&contractAddress, "contract", "", "Workflow Registry contract address (optional if in config)")

	// Register command flags
	registerCmd.Flags().StringVar(&workflowName, "name", "", "Workflow name (required)")
	registerCmd.Flags().StringVar(&workflowID, "id", "", "Workflow ID in hex (required)")
	registerCmd.Flags().StringVar(&binaryURL, "binary-url", "", "URL to workflow binary WASM file (required)")
	registerCmd.Flags().StringVar(&configURL, "config-url", "", "URL to workflow configuration file")
	registerCmd.Flags().StringVar(&secretsURL, "secrets-url", "", "URL to encrypted workflow secrets file")
	registerCmd.Flags().Uint8Var(&initialStatus, "status", 0, "Initial status (0=active, 1=paused)")
	registerCmd.Flags().Uint32Var(&donID, "don-id", 1, "DON ID where the workflow will run")
	registerCmd.MarkFlagRequired("name")
	registerCmd.MarkFlagRequired("id")
	registerCmd.MarkFlagRequired("binary-url")

	// Update command flags
	updateCmd.Flags().StringVar(&workflowName, "name", "", "Workflow name (required)")
	updateCmd.Flags().StringVar(&workflowOwner, "owner", "", "Workflow owner address (optional if in config)")
	updateCmd.Flags().StringVar(&binaryURL, "new-binary-url", "", "URL to workflow binary WASM file (required)")
	updateCmd.Flags().StringVar(&configURL, "new-config-url", "", "URL to workflow configuration file")
	updateCmd.Flags().StringVar(&secretsURL, "new-secrets-url", "", "URL to encrypted workflow secrets file")
	updateCmd.Flags().StringVar(&workflowID, "new-id", "", "Workflow ID in hex (required)")
	updateCmd.MarkFlagRequired("name")
	updateCmd.MarkFlagRequired("new-id")
	updateCmd.MarkFlagRequired("new-binary-url")
	updateCmd.MarkFlagRequired("new-config-url")
	updateCmd.MarkFlagRequired("new-secrets-url")

	deleteCmd.Flags().StringVar(&workflowName, "name", "", "Workflow name (required)")
	deleteCmd.Flags().StringVar(&workflowOwner, "owner", "", "Workflow owner address (optional if in config)")
	deleteCmd.MarkFlagRequired("name")

	getCmd.Flags().StringVar(&workflowOwner, "owner", "", "Workflow owner address (optional if in config)")

	// Add subcommands to main workflow command
	WorkflowCmd.AddCommand(registerCmd, updateCmd, activateCmd, pauseCmd, deleteCmd, getCmd)
}

// ComputeWorkflowKey creates a workflowKey fron owner and workflow name
func ComputeHashKey(owner common.Address, name string) [32]byte {
	// Pack the values together (equivalent to abi.encodePacked)
	packed := append(owner.Bytes(), []byte(name)...)

	// Compute keccak256 hash
	hash := crypto.Keccak256(packed)

	// Convert to bytes32
	var result [32]byte
	copy(result[:], hash)

	return result
}
