// Package capabilityregistry provides commands for managing capabilities on the Capability Registry contract.
// It enables creating, listing, and associating capabilities with DONs
package capabilityregistry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
)

var (
	configPath        string
	contractAddress   string
	capabilityType    string
	capabilityName    string
	capabilityVersion string
	donID             uint32
)

// CapabilityType constants define the different types of capabilities
const (
	CapabilityTypeTrigger   uint8 = 0
	CapabilityTypeAction    uint8 = 1
	CapabilityTypeConsensus uint8 = 2
	CapabilityTypeTarget    uint8 = 3
)

// RegistryCmd is the root command for capability registry operations.
// It provides subcommands for creating and listing capabilities.
var RegistryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage the capability registry",
	Long:  "Manage capabilities in the Chainlink Capability Registry contract, including creating new capabilities and listing existing ones.",
}

// createCmd handles the creation of new capabilities and associates them with a DON.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Adds a capability to the capability registry",
	Long:  "Creates a new capability in the registry and associates it with the specified DON. All nodes in the DON will be updated to support this capability.",
	Run: func(cmd *cobra.Command, args []string) {
		// Setup client
		client, seth, err := setupRegistryClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to setup capability registry client: %v\n", err)
			return
		}

		// Check if capability already exists
		caps, err := client.GetCapabilities(seth.NewCallOpts())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get capabilities: %v\n", err)
			return
		}

		capExists := false
		for _, cap := range caps {
			if cap.LabelledName == capabilityName && cap.Version == capabilityVersion {
				fmt.Printf("Capability '%s@%s' already exists\n", capabilityName, capabilityVersion)
				capExists = true
				break
			}
		}

		// Add capability
		if !capExists {
			txOpt, err2 := setupTransaction(seth)
			if err2 != nil {
				fmt.Fprintf(os.Stderr, "Failed to create transaction: %v\n", err2)
				return
			}
			output, err2 := client.AddCapabilities(txOpt, []capabilities_registry.CapabilitiesRegistryCapability{
				{
					LabelledName:   capabilityName,
					Version:        capabilityVersion,
					CapabilityType: stringToCapabilityType(capabilityType),
					ResponseType:   0,
				},
			})
			if err2 != nil {
				fmt.Fprintf(os.Stderr, "Failed to create capability: %v\n", err2)
				return
			}

			fmt.Printf("Capability '%s@%s' successfully created\n", capabilityName, capabilityVersion)
			fmt.Printf("Transaction hash: %s\n", output.Hash().Hex())
		}

		don, err := client.GetDON(seth.NewCallOpts(), donID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get nodes: %v\n", err)
			return
		}

		fmt.Printf("Don %d has %d nodes\n", donID, len(don.NodeP2PIds))

		// Fetch Node Data
		nodes, err := client.GetNodes(seth.NewCallOpts())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get nodes: %v\n", err)
			return
		}

		// Figure out which nodes we want to update
		nodesToUpdate := make([]capabilities_registry.CapabilitiesRegistryNodeParams, 0)

		for _, node := range nodes {
			// Check if the node's P2pId is contained in the don.NodeP2PIds slice
			isInDon := slices.Contains(don.NodeP2PIds, node.P2pId)
			if isInDon {
				nodesToUpdate = append(nodesToUpdate, capabilities_registry.CapabilitiesRegistryNodeParams{
					NodeOperatorId:      node.NodeOperatorId,
					Signer:              node.Signer,
					P2pId:               node.P2pId,
					EncryptionPublicKey: node.EncryptionPublicKey,
					HashedCapabilityIds: append(node.HashedCapabilityIds, getHashedCapabilityID(capabilityName, capabilityVersion)),
				})
			}
		}

		txOpt, err := setupTransaction(seth)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create transaction: %v\n", err)
			return
		}
		output, err := client.UpdateNodes(txOpt, nodesToUpdate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to update nodes: %v\n", err)
			return
		}
		fmt.Printf("Nodes successfully updated, tx hash: %s\n", output.Hash().Hex())

		txOpt, err = setupTransaction(seth)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create transaction: %v\n", err)
			return
		}
		// apend cap config
		don.CapabilityConfigurations = append(don.CapabilityConfigurations, capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			CapabilityId: getHashedCapabilityID(capabilityName, capabilityVersion),
			Config:       nil,
		})
		updateDON, err := client.UpdateDON(txOpt, donID, don.NodeP2PIds, don.CapabilityConfigurations, don.IsPublic, don.F)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to update don: %v\n", err)
			return
		}
		fmt.Printf("DON successfully updated, tx hash: %s\n", updateDON.Hash().Hex())

	},
}

// listCmd displays all capabilities and DON information from the registry.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all capabilities inside the registry",
	Long:  "Displays detailed information about all capabilities and DONs registered in the Capability Registry contract.",
	Run: func(cmd *cobra.Command, args []string) {
		// Setup client
		client, seth, err := setupRegistryClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to setup capability registry client: %v\n", err)
			return
		}

		output, err := client.GetCapabilities(seth.NewCallOpts())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to list capabilities: %v\n", err)
			return
		}

		dons, err := client.GetDONs(seth.NewCallOpts())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get DONs: %v\n", err)
			return
		}

		for i, don := range dons {
			// Add separator between DONs
			if i > 0 {
				fmt.Println(strings.Repeat("-", 50))
			}

			// DON header with clear formatting
			fmt.Printf("DON ID: %d\n", don.Id)
			fmt.Println(strings.Repeat("=", 30))

			// Group basic DON information in a structured format
			fmt.Printf("  %-20s %d\n", "Config count:", don.ConfigCount)
			fmt.Printf("  %-20s %d\n", "F value:", don.F)
			fmt.Printf("  %-20s %t\n", "Is public:", don.IsPublic)
			fmt.Printf("  %-20s %t\n", "Accepts workflows:", don.AcceptsWorkflows)

			// Format node IDs in a cleaner way
			fmt.Println("\n  Node P2P IDs:")
			if len(don.NodeP2PIds) == 0 {
				fmt.Println("    None")
			} else {
				for j, nodeID := range don.NodeP2PIds {
					fmt.Printf("    %d. %x\n", j+1, nodeID)
				}
			}

			// Capabilities section with clear header
			fmt.Println("\n  Capabilities:")
			if len(don.CapabilityConfigurations) == 0 {
				fmt.Println("    None")
			} else {
				for j, cap := range don.CapabilityConfigurations {
					capFound := false
					for _, c := range output {
						if c.HashedId == cap.CapabilityId {
							capFound = true
							fmt.Printf("    Capability #%d:\n", j+1)
							fmt.Printf("      %-24s %s\n", "Name:", c.LabelledName)
							fmt.Printf("      %-24s %s\n", "Version:", c.Version)
							fmt.Printf("      %-24s %s (%d)\n", "Type:", capabilityTypeToString(c.CapabilityType), c.CapabilityType)
							fmt.Printf("      %-24s %x\n", "Hashed ID:", c.HashedId)
							fmt.Printf("      %-24s %s\n", "Configuration contract:", c.ConfigurationContract)
							fmt.Printf("      %-24s %t\n", "Is deprecated:", c.IsDeprecated)
							fmt.Printf("      %-24s %d\n", "Response type:", c.ResponseType)
							break
						}
					}
					if !capFound {
						fmt.Printf("    Capability #%d: (ID: %x) - Details not found\n", j+1, cap.CapabilityId)
					}
				}
			}
			fmt.Println()
		}
	},
}

// init initializes the registry command and its subcommands.
// It registers all command flags and their requirements.
func init() {
	// Add config flag first so it's loaded before other flags are processed
	RegistryCmd.PersistentFlags().StringVar(&configPath, "config", "cre.yaml", "Path to cre.yaml config file")

	createCmd.Flags().StringVar(&capabilityName, "name", "", "Capability name (required)")
	createCmd.Flags().StringVar(&capabilityVersion, "version", "", "Capability version (required)")
	createCmd.Flags().StringVar(&capabilityType, "type", "", "Capability type (required)")
	createCmd.Flags().Uint32Var(&donID, "don-id", 0, "DonID (required)")
	if err := createCmd.MarkFlagRequired("name"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'name' flag as required: %v\n", err)
		panic(err)
	}
	if err := createCmd.MarkFlagRequired("version"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'version' flag as required: %v\n", err)
		panic(err)
	}
	if err := createCmd.MarkFlagRequired("type"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'type' flag as required: %v\n", err)
		panic(err)
	}
	if err := createCmd.MarkFlagRequired("don-id"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'don-id' flag as required: %v\n", err)
		panic(err)
	}

	RegistryCmd.AddCommand(createCmd, listCmd)
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

// setupRegistryClient creates a capability registry client using RPC information from cre.yaml
// It loads the configuration file, extracts connection information, and initializes
// an Ethereum client connected to the specified RPC endpoint.
//
// Returns:
//   - Initialized CapabilitiesRegistry client
//   - Seth client for blockchain interactions
//   - Error if any step fails
func setupRegistryClient() (*capabilities_registry.CapabilitiesRegistry, *seth.Client, error) {
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
		if v.Name == "CapabilitiesRegistry" {
			contractAddress = v.Address
			fmt.Printf("Using contract address from config: %s\n", contractAddress)
			break
		}
	}

	// Validate we have at least one RPC URL
	if len(config.Test.Rpcs) == 0 {
		return nil, nil, errors.New("no RPC URLs found in config file")
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
		return nil, nil, fmt.Errorf("failed to create Ethereum client: %w", err)
	}

	if contractAddress == "" {
		return nil, nil, errors.New("capability registry contract address not provided in flags or config")
	}

	// Create workflow registry client
	registryClient, err := capabilities_registry.NewCapabilitiesRegistry(
		common.HexToAddress(contractAddress),
		ethClient.Client,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create workflow registry client: %w", err)
	}

	return registryClient, ethClient, nil
}

// stringToCapabilityType converts a string capability type to the corresponding integer constant
func stringToCapabilityType(typeStr string) uint8 {
	switch strings.ToUpper(typeStr) {
	case "TRIGGER":
		return CapabilityTypeTrigger
	case "ACTION":
		return CapabilityTypeAction
	case "CONSENSUS":
		return CapabilityTypeConsensus
	case "TARGET":
		return CapabilityTypeTarget
	default:
		fmt.Fprintf(os.Stderr, "Warning: Unknown capability type: %s. Using TRIGGER as default.\n", typeStr)
		return CapabilityTypeTrigger
	}
}

// capabilityTypeToString converts a uint8 capability type to its string representation
func capabilityTypeToString(typeInt uint8) string {
	switch typeInt {
	case CapabilityTypeTrigger:
		return "TRIGGER"
	case CapabilityTypeAction:
		return "ACTION"
	case CapabilityTypeConsensus:
		return "CONSENSUS"
	case CapabilityTypeTarget:
		return "TARGET"
	default:
		return "UNKNOWN"
	}
}

// getHashedCapabilityID creates a unique bytes32 key from capability name and version.
// The key is created by hashing the packed ABI encoding of the name and version strings.
//
// Parameters:
//
//	labelledName - The capability name
//	version - The capability version string
//
// Returns:
//
//	A bytes32 representation of the capability ID
func getHashedCapabilityID(labelledName string, version string) [32]byte {
	// Create proper ABI encoding for strings
	stringType, _ := abi.NewType("string", "", nil)
	arguments := abi.Arguments{
		{Type: stringType},
		{Type: stringType},
	}

	encodedData, err := arguments.Pack(labelledName, version)
	if err != nil {
		panic(err)
	}

	hash := crypto.Keccak256Hash(encodedData)
	return hash
}
