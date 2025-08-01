// Package workflow provides commands for managing workflows on the Workflow Registry blockchain contract.
// It enables registering, updating, deleting, and querying workflow metadata.
package workflow

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	workflow_registry_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
)

const (
	// Workflow status values
	WorkflowStatusActive = 0
	WorkflowStatusPaused = 1

	// Maximum number of workflows to fetch in a single query
	MaxWorkflowsToFetch = 9999999
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

// WorkflowCmd is the root command for workflow management operations.
// It provides subcommands for registering, updating, deleting, and querying workflows.
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
		err := validateStatus(initialStatus)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid status: %v\n", err)
			return
		}

		// Setup client
		client, seth, err := setupWorkflowClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to setup workflow client: %v\n", err)
			return
		}

		// Parse workflowID to bytes32
		var workflowIDBytes [32]byte
		copy(workflowIDBytes[:], common.FromHex(workflowID))

		txOpt, err := setupTransaction(seth)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create transactor: %v\n", err)
			return
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
		client, seth, err := setupWorkflowClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to setup workflow client: %v\n", err)
			return
		}

		// Parse workflowID and owner
		var workflowIDBytes [32]byte
		copy(workflowIDBytes[:], common.FromHex(workflowID))

		// Update workflow
		txOpt, err := setupTransaction(seth)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create transactor: %v\n", err)
			return
		}
		output, err := client.UpdateWorkflow(txOpt, ComputeWorkflowKey(common.HexToAddress(workflowOwner), workflowName), workflowIDBytes, binaryURL, configURL, secretsURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to update workflow: %v\n", err)
			return
		}

		fmt.Printf("Workflow '%s' successfully updated\n", workflowName)
		fmt.Printf("Transaction hash: %s\n", output.Hash().Hex())
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

		txOpt, err := setupTransaction(seth)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create transactor: %v\n", err)
			return
		}
		txHash, err := client.DeleteWorkflow(txOpt, ComputeWorkflowKey(common.HexToAddress(workflowOwner), workflowName))
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
		metadata, err := client.GetWorkflowMetadataListByOwner(seth.NewCallOpts(), owner, big.NewInt(0), big.NewInt(MaxWorkflowsToFetch))
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
			fmt.Printf("Status: %d (%s)\n", m.Status, formatStatus(m.Status))
			fmt.Printf("Binary URL: %s\n", m.BinaryURL)
			fmt.Printf("Config URL: %s\n", m.ConfigURL)
			fmt.Printf("Secrets URL: %s\n", m.SecretsURL)
		}

	},
}

// setupWorkflowClient creates a workflow registry client using RPC information from cre.yaml.
// It loads the configuration file, extracts necessary connection information,
// and initializes an Ethereum client connected to the specified RPC endpoint.
//
// The function tries to use contract address and workflow owner from config if not provided via flags.
//
// Returns:
//   - Initialized WorkflowRegistry client
//   - Seth client for blockchain interactions
//   - Error if any step fails
func setupWorkflowClient() (*workflow_registry_wrapper.WorkflowRegistry, *seth.Client, error) {
	// Load and parse the cre.yaml file
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading config file: %w", err)
	}
	config := crecli.Profiles{}

	if err = yaml.Unmarshal(configData, &config); err != nil {
		return nil, nil, fmt.Errorf("error parsing config file: %w", err)
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
		return nil, nil, errors.New("no RPC URLs found in config file")
	}

	// Use the first RPC URL by default
	rpcURL := config.Test.Rpcs[0].URL

	// Create Ethereum client
	ethClient, err := seth.NewClientBuilder().
		WithRpcUrl(rpcURL).
		WithPrivateKeys([]string{blockchain.DefaultAnvilPrivateKey}).
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Ethereum client: %w", err)
	}

	if contractAddress == "" {
		return nil, nil, errors.New("workflow registry contract address not provided in flags or config")
	}

	// Create workflow registry client
	workflowClient, err := workflow_registry_wrapper.NewWorkflowRegistry(
		common.HexToAddress(contractAddress),
		ethClient.Client,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create workflow registry client: %w", err)
	}

	return workflowClient, ethClient, nil
}

// init initializes the workflow command and its subcommands.
// It registers all command flags and their requirements, and adds
// the subcommands to the main workflow command.
func init() {
	// Add config flag first so it's loaded before other flags are processed
	WorkflowCmd.PersistentFlags().StringVar(&configPath, "config", "cre.yaml", "Path to cre.yaml config file")

	// Add common flags - now optional since they can come from config
	WorkflowCmd.PersistentFlags().StringVar(&contractAddress, "contract", "", "Workflow Registry contract address (optional if in config)")

	// Register command flags
	registerCmd.Flags().StringVar(&workflowID, "id", "", "Workflow id")
	registerCmd.Flags().StringVar(&workflowName, "name", "", "Workflow name (required)")
	registerCmd.Flags().StringVar(&binaryURL, "binary-url", "", "URL to workflow binary WASM file (required)")
	registerCmd.Flags().StringVar(&configURL, "config-url", "", "URL to workflow configuration file")
	registerCmd.Flags().StringVar(&secretsURL, "secrets-url", "", "URL to encrypted workflow secrets file")
	registerCmd.Flags().Uint8Var(&initialStatus, "status", 0, "Initial status (0=active, 1=paused)")
	registerCmd.Flags().Uint32Var(&donID, "don-id", 1, "DON ID where the workflow will run")
	// Required
	if err := MarkFlagsRequired(registerCmd,
		"id",
		"name",
		"binary-url",
		"config-url",
		"secrets-url"); err != nil {
		panic(err)
	}

	// Update command flags
	updateCmd.Flags().StringVar(&workflowName, "name", "", "Workflow name (required)")
	updateCmd.Flags().StringVar(&workflowOwner, "owner", "", "Workflow owner address (optional if in config)")
	updateCmd.Flags().StringVar(&binaryURL, "new-binary-url", "", "URL to workflow binary WASM file (required)")
	updateCmd.Flags().StringVar(&configURL, "new-config-url", "", "URL to workflow configuration file")
	updateCmd.Flags().StringVar(&secretsURL, "new-secrets-url", "", "URL to encrypted workflow secrets file")
	updateCmd.Flags().StringVar(&workflowID, "new-id", "", "Workflow ID in hex (required)")
	// Required
	if err := MarkFlagsRequired(updateCmd,
		"name",
		"new-id",
		"new-binary-url",
		"new-config-url",
		"new-secrets-url"); err != nil {
		panic(err)
	}

	deleteCmd.Flags().StringVar(&workflowName, "name", "", "Workflow name (required)")
	deleteCmd.Flags().StringVar(&workflowOwner, "owner", "", "Workflow owner address (optional if in config)")
	// Required
	err := deleteCmd.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}
	// Get command flags
	getCmd.Flags().StringVar(&workflowOwner, "owner", "", "Workflow owner address (optional if in config)")

	// Add subcommands to workflow command
	WorkflowCmd.AddCommand(registerCmd, updateCmd, deleteCmd, getCmd)
}

// ComputeWorkflowKey computes a unique bytes32 key from owner address and workflow name.
// The key is created by hashing the packed concatenation of the owner address and name.
// This key uniquely identifies a workflow in the registry contract.
//
// Parameters:
//
//	owner - The Ethereum address of the workflow owner
//	name - The workflow name as a string
//
// Returns:
//
//	A bytes32 representation of the workflow key
func ComputeWorkflowKey(owner common.Address, name string) [32]byte {
	// Pack the values together (equivalent to abi.encodePacked)
	packed := append(owner.Bytes(), []byte(name)...)

	// Compute keccak256 hash
	hash := crypto.Keccak256(packed)

	// Convert to bytes32
	var result [32]byte
	copy(result[:], hash)

	return result
}

// setupTransaction creates a keyed transactor with the current chain ID.
// This is used to sign and send transactions to the blockchain.
//
// Parameters:
//
//	seth - The Seth client used to get chain ID and private key
//
// Returns:
//   - Configured transaction options for blockchain transactions
//   - Error if fetching chain ID fails
func setupTransaction(seth *seth.Client) (*bind.TransactOpts, error) {
	chainID, err := seth.Client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	return bind.NewKeyedTransactorWithChainID(seth.MustGetRootPrivateKey(), chainID)
}

// formatStatus converts numeric status to a human-readable string.
// It handles the standard workflow statuses and provides a fallback for unknown values.
//
// Parameters:
//
//	status - Numeric workflow status (0=active, 1=paused)
//
// Returns:
//
//	A descriptive string representation of the status
func formatStatus(status uint8) string {
	switch status {
	case WorkflowStatusActive:
		return "Active (0)"
	case WorkflowStatusPaused:
		return "Paused (1)"
	default:
		return fmt.Sprintf("Unknown (%d)", status)
	}
}

// validateStatus checks if the provided status is valid.
// Valid statuses are WorkflowStatusActive (0) and WorkflowStatusPaused (1).
//
// Returns:
//
//	nil if status is valid, error with description otherwise
func validateStatus(status uint8) error {
	if status != WorkflowStatusActive && status != WorkflowStatusPaused {
		return fmt.Errorf("invalid status: %d (must be 0=active or 1=paused)", status)
	}
	return nil
}

// MarkFlagsRequired marks multiple flags as required for the given command.
// It takes a command and a variable number of flag names, marking each one as required.
// If any flag cannot be marked as required, it returns an error with details.
//
// Example usage:
//
//	if err := MarkFlagsRequired(updateCmd, "name", "new-id", "new-binary-url"); err != nil {
//	    fmt.Fprintf(os.Stderr, "Error marking required flags: %v\n", err)
//	}
func MarkFlagsRequired(cmd *cobra.Command, flagNames ...string) error {
	for _, flagName := range flagNames {
		if err := cmd.MarkFlagRequired(flagName); err != nil {
			return fmt.Errorf("failed to mark flag '%s' as required: %w", flagName, err)
		}
	}
	return nil
}
