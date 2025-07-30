package capabilityregistry

import (
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
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

// setupWorkflowClient creates a workflow registry client using RPC information from cre.yaml
func setupRegistryClient() (*capabilities_registry.CapabilitiesRegistry, *seth.Client, error) {
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
		if v.Name == "CapabilitiesRegistry" {
			contractAddress = v.Address
			fmt.Printf("Using contract address from config: %s\n", contractAddress)
			break
		}
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
		return nil, nil, fmt.Errorf("capability registry contract address not provided in flags or config")
	}

	// Create workflow registry client
	registryClient, err := capabilities_registry.NewCapabilitiesRegistry(
		common.HexToAddress(contractAddress),
		ethClient.Client,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create workflow registry client: %v", err)
	}

	return registryClient, ethClient, nil
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Adds a capability to the capability registry",
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
			txOpt, err := bind.NewKeyedTransactorWithChainID(seth.MustGetRootPrivateKey(), big.NewInt(int64(chain_selectors.GETH_TESTNET.EvmChainID)))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create transaction: %v\n", err)
				return
			}
			output, err := client.AddCapabilities(txOpt, []capabilities_registry.CapabilitiesRegistryCapability{
				{
					LabelledName:   capabilityName,
					Version:        capabilityVersion,
					CapabilityType: stringToCapabilityType(capabilityType),
					ResponseType:   0,
				},
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create capability: %v\n", err)
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
			isInDon := false
			for _, donNodeP2pId := range don.NodeP2PIds {
				if node.P2pId == donNodeP2pId {
					isInDon = true
					break
				}
			}
			if isInDon {
				nodesToUpdate = append(nodesToUpdate, capabilities_registry.CapabilitiesRegistryNodeParams{
					NodeOperatorId:      node.NodeOperatorId,
					Signer:              node.Signer,
					P2pId:               node.P2pId,
					EncryptionPublicKey: node.EncryptionPublicKey,
					HashedCapabilityIds: append(node.HashedCapabilityIds, getHashedCapabilityId(capabilityName, capabilityVersion)),
				})
			}
		}

		txOpt, err := bind.NewKeyedTransactorWithChainID(seth.MustGetRootPrivateKey(), big.NewInt(int64(chain_selectors.GETH_TESTNET.EvmChainID)))
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

		txOpt, err = bind.NewKeyedTransactorWithChainID(seth.MustGetRootPrivateKey(), big.NewInt(int64(chain_selectors.GETH_TESTNET.EvmChainID)))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create transaction: %v\n", err)
			return
		}
		// apend cap config
		don.CapabilityConfigurations = append(don.CapabilityConfigurations, capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			CapabilityId: getHashedCapabilityId(capabilityName, capabilityVersion),
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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all capability inside the registries",
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

		for _, don := range dons {
			// Print all info in the don
			fmt.Printf("DON ID: %d\n", don.Id)
			fmt.Printf("Config count: %d\n", don.ConfigCount)
			fmt.Printf("F: %d\n", don.F)
			fmt.Printf("isPublic: %t\n", don.IsPublic)
			fmt.Printf("Accepts Workflows: %t\n", don.AcceptsWorkflows)
			fmt.Printf("NodeP2PIds: %x\n", don.NodeP2PIds)

			fmt.Println("Capabilities:")
			for _, cap := range don.CapabilityConfigurations {
				for _, c := range output {
					if c.HashedId == cap.CapabilityId {
						fmt.Printf("\tLabelledName: %s\n", c.LabelledName)
						fmt.Printf("\tVersion: %s\n", c.Version)
						fmt.Printf("\tCapability type: %d (%s)\n", c.CapabilityType, capabilityTypeToString(c.CapabilityType))
						fmt.Printf("\tHashedID: %x\n", c.HashedId)
						fmt.Printf("\tConfiguration contract: %s\n", c.ConfigurationContract)
						fmt.Printf("\tisDepecrated: %t\n", c.IsDeprecated)
						fmt.Printf("\tResponse Type: %d\n\n", c.ResponseType)
					}
				}

			}

		}

	},
}

// StringToCapabilityType converts a string capability type to the corresponding integer constant
func stringToCapabilityType(typeStr string) uint8 {
	switch strings.ToUpper(typeStr) {
	case "TRIGGER":
		return 0
	case "ACTION":
		return 1
	case "CONSENSUS":
		return 2
	case "TARGET":
		return 3
	default:
		panic(fmt.Errorf("unknown capability type: %s", typeStr))
	}
}

// CapabilityTypeToString converts a uint8 capability type to its string representation
func capabilityTypeToString(typeInt uint8) string {
	switch typeInt {
	case 0:
		return "TRIGGER"
	case 1:
		return "ACTION"
	case 2:
		return "CONSENSUS"
	case 3:
		return "TARGET"
	default:
		return "UNKNOWN"
	}
}

// GetHashedCapabilityId mimics the Solidity function that computes the hashed capability ID
func getHashedCapabilityId(labelledName string, version string) [32]byte {
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

var RegistryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage the capability registry",
}

func init() {
	// Add config path flag first so it's loaded before other flags are processed
	RegistryCmd.PersistentFlags().StringVar(&configPath, "config", "cre.yaml", "Path to cre.yaml config file")

	createCmd.Flags().StringVar(&capabilityName, "name", "", "Capability name (required)")
	createCmd.Flags().StringVar(&capabilityVersion, "version", "", "Capability version (required)")
	createCmd.Flags().StringVar(&capabilityType, "type", "", "Capability type (required)")
	createCmd.Flags().Uint32Var(&donID, "don-id", 0, "DonID (required)")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("version")
	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagRequired("don-id")

	RegistryCmd.AddCommand(createCmd, listCmd)
}
