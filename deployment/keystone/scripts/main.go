package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/client"
)

// Config holds the TOML configuration with multiple nodes
type Config struct {
	Nodes []NodeConfig `toml:"nodes"`
}

// NodeConfig holds the configuration for a single node
type NodeConfig struct {
	Name        string      `toml:"name"`
	Bootstrap   bool        `toml:"bootstrap"`
	Credentials Credentials `toml:"credentials"`
	Connection  Connection  `toml:"connection"`
}

type Credentials struct {
	Email    string `toml:"email"`
	Password string `toml:"password"`
}
type Connection struct {
	BaseURL string `toml:"base_url"`
}

// JobDistributorConfig defines a job distributor configuration
type JobDistributorConfig struct {
	Name      string `toml:"name"`
	URI       string `toml:"uri"`
	PublicKey string `toml:"public_key"`
}

// NodeInfo represents the output JSON format
type NodeInfo struct {
	Name           string                 `json:"name"`
	URL            string                 `json:"url"`
	CSAPublicKey   string                 `json:"csa_public_key,omitempty"`
	P2PPeerID      string                 `json:"p2p_peer_id,omitempty"`
	OCR2KeyBundles []client.OCR2KeyBundle `json:"ocr2_key_bundles,omitempty"`
}

type BootstrapNodeInfo struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	CSAPublicKey string `json:"csa_public_key,omitempty"`
	P2PPeerID    string `json:"p2p_peer_id,omitempty"`
	OCRUrl       string `json:"ocr_url,omitempty"`
	Don2DonURL   string `json:"don2_don_url,omitempty"`
}

var (
	// Global flags
	configPath string

	// Info command flags
	outputPath          string
	bootstrapOutputPath string

	// JD command flags
	jdConfigPath string
)

func main() {
	// Root command
	rootCmd := &cobra.Command{
		Use:   "keystone",
		Short: "Chainlink Keystone node management tool",
		Long:  `A CLI tool for managing Chainlink nodes in Keystone environment`,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "nodes.toml", "Path to TOML config file")

	// Info command
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Gather information about nodes",
		Long:  `Connects to nodes and retrieves key information`,
		RunE:  runInfoCommand,
	}
	infoCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to output JSON file (stdout if not specified)")
	infoCmd.Flags().StringVarP(&bootstrapOutputPath, "bootstrap-output", "b", "", "Path to output JSON file for bootstrap nodes (stdout if not specified)")

	// JD command

	// In the main() function, update the JD command section:

	// JD command hierarchy
	jdCmd := &cobra.Command{
		Use:   "jd",
		Short: "Job Distributor operations",
		Long:  `Commands for managing Job Distributors and blockchain integrations`,
	}

	// Original JD create command (rename for clarity)
	jdCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create Job Distributors",
		Long:  `Creates Job Distributors for nodes based on configuration`,
		RunE:  runJDCommand,
	}
	jdCreateCmd.Flags().StringVarP(&jdConfigPath, "jd-config", "j", "jd.toml", "Path to Job Distributor TOML config file")

	// Aptos chain enablement command
	jdAptosCmd := newJDAptosCmd()

	// Build JD command hierarchy
	jdCmd.AddCommand(jdCreateCmd)
	jdCmd.AddCommand(jdAptosCmd)
	jdCmd.AddCommand(newJDAcceptCmd())

	jdCmd.Flags().StringVarP(&jdConfigPath, "jd-config", "j", "jd.toml", "Path to Job Distributor TOML config file")

	// Keys command hierarchy
	keysCmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage cryptographic keys",
		Long:  `Commands for creating and managing various types of cryptographic keys on Chainlink nodes`,
	}

	// Aptos keys subcommand
	aptosKeysCmd := &cobra.Command{
		Use:   "aptos",
		Short: "Manage Aptos keys",
		Long:  `Commands for managing Aptos cryptographic keys`,
	}

	// Aptos create command
	aptosCreateCmd := newCreateAptosKeysCmd()

	// Build command hierarchy
	aptosKeysCmd.AddCommand(aptosCreateCmd)
	aptosKeysCmd.AddCommand(newListAptosKeysCmd())
	keysCmd.AddCommand(aptosKeysCmd)

	// Add commands to root
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(jdCmd)
	rootCmd.AddCommand(keysCmd)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadConfig(path string) (Config, error) {
	var config Config

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("config file not found: %s", path)
	}

	if _, err := toml.DecodeFile(path, &config); err != nil {
		return Config{}, fmt.Errorf("failed to decode config file: %w", err)
	}

	// Check if we have any nodes configured
	if len(config.Nodes) == 0 {
		return Config{}, errors.New("no nodes configured in the config file")
	}

	return config, nil
}

func loadJDConfig(path string) (JobDistributorConfig, error) {
	var jdConfig JobDistributorConfig
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return JobDistributorConfig{}, fmt.Errorf("job Distributor config file not found: %s", path)
	}
	if _, err := toml.DecodeFile(path, &jdConfig); err != nil {
		return JobDistributorConfig{}, fmt.Errorf("failed to decode job Distributor config file: %w", err)
	}
	return jdConfig, nil
}

func runInfoCommand(cmd *cobra.Command, args []string) error {
	config, err := loadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	// Create slice to hold all node information
	var nodes []NodeInfo
	var bootstrapNodes []BootstrapNodeInfo

	// Process each node
	for _, nodeConfig := range config.Nodes {
		// Validate node configuration
		if nodeConfig.Name == "" {
			log.Println("Skipping node with no name")
			continue
		}

		if nodeConfig.Credentials.Email == "" || nodeConfig.Credentials.Password == "" {
			log.Printf("Skipping node %s: missing credentials", nodeConfig.Name)
			continue
		}

		if nodeConfig.Connection.BaseURL == "" {
			log.Printf("Skipping node %s: missing base URL", nodeConfig.Name)
			continue
		}

		// Create timeout context
		timeout := time.Duration(10) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		info, err := nodeInfo(ctx, nodeConfig)
		if err != nil {
			log.Printf("Failed to get info for node %s: %v", nodeConfig.Name, err)
			cancel()
			continue
		}
		if !nodeConfig.Bootstrap {
			nodes = append(nodes, info)
		}

		if nodeConfig.Bootstrap {
			p := strings.TrimPrefix(info.P2PPeerID, "p2p_")
			bootstrapNodeInfo := BootstrapNodeInfo{
				Name:         nodeConfig.Name,
				URL:          nodeConfig.Connection.BaseURL,
				CSAPublicKey: info.CSAPublicKey,
				P2PPeerID:    info.P2PPeerID,
				OCRUrl:       fmt.Sprintf("%s@%s:5001", p, nodeConfig.Name),
				Don2DonURL:   fmt.Sprintf("%s@%s:6690", p, nodeConfig.Name),
			}
			bootstrapNodes = append(bootstrapNodes, bootstrapNodeInfo)
		}
		cancel()
	}

	// Convert nodes to JSON
	jsonOutput, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to convert nodes to JSON: %w", err)
	}
	// Convert bootstrap nodes to JSON
	bootstrapJSONOutput, err := json.MarshalIndent(bootstrapNodes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to convert bootstrap nodes to JSON: %w", err)
	}

	// Write output to file or stdout
	if outputPath == "" {
		fmt.Println(string(jsonOutput))
	} else {
		err = os.WriteFile(outputPath, jsonOutput, 0600)
		if err != nil {
			return fmt.Errorf("failed to write output to file: %w", err)
		}
		log.Printf("Output written to %s", outputPath)
	}

	// Write bootstrap output to file or stdout
	if bootstrapOutputPath == "" {
		fmt.Println(string(bootstrapJSONOutput))
	} else {
		err = os.WriteFile(bootstrapOutputPath, bootstrapJSONOutput, 0600)
		if err != nil {
			return fmt.Errorf("failed to write bootstrap output to file: %w", err)
		}
		log.Printf("Bootstrap output written to %s", bootstrapOutputPath)
	}

	return nil
}

func nodeInfo(ctx context.Context, nodeConfig NodeConfig) (NodeInfo, error) {
	nodeInfo := NodeInfo{
		Name: nodeConfig.Name,
		URL:  nodeConfig.Connection.BaseURL,
	}
	// Create client credentials
	creds := client.Credentials{
		Email:    nodeConfig.Credentials.Email,
		Password: nodeConfig.Credentials.Password,
	}

	// Try to create client and fetch keys
	cl, err := client.NewWithContext(ctx, nodeConfig.Connection.BaseURL, creds)
	if err != nil {
		return NodeInfo{}, fmt.Errorf("failed to connect to node %s: %w", nodeConfig.Name, err)
	}

	// Fetch CSA Public Key
	csaKey, err := cl.FetchCSAPublicKey(ctx)
	if err == nil && csaKey != nil {
		nodeInfo.CSAPublicKey = *csaKey
	}

	// Fetch P2P Peer ID
	peerID, err := cl.FetchP2PPeerID(ctx)
	if err == nil && peerID != nil {
		nodeInfo.P2PPeerID = *peerID
	}

	if !nodeConfig.Bootstrap {
		ocrKeyBundleIDs, err := cl.ListOCR2KeyBundles(ctx)
		if err == nil {
			nodeInfo.OCR2KeyBundles = ocrKeyBundleIDs
		}
		// if no aptos ocr2 key bundles are found, create one and retry fetching
		foundAptosKeyBundle := false
		for _, keyBundle := range nodeInfo.OCR2KeyBundles {
			if keyBundle.ChainType == client.OCR2ChainTypeAptos {
				foundAptosKeyBundle = true
				break
			}
		}
		if !foundAptosKeyBundle {
			// Create a new OCR2 key bundle for Aptos
			aptosKeyBundleID, err := cl.CreateOCR2KeyBundle(ctx, client.OCR2ChainTypeAptos)
			if err != nil {
				return nodeInfo, fmt.Errorf("failed to create OCR2 key bundle for Aptos on node %s: %w", nodeConfig.Name, err)
			}
			log.Printf("Created OCR2 key bundle for Aptos on node %s with ID: %s", nodeConfig.Name, aptosKeyBundleID)
			// Retry fetching OCR2 key bundles
			ocrKeyBundleIDs, err = cl.ListOCR2KeyBundles(ctx)
			if err != nil {
				return nodeInfo, fmt.Errorf("failed to fetch OCR2 key bundles after creating new one on node %s: %w", nodeConfig.Name, err)
			}
			nodeInfo.OCR2KeyBundles = ocrKeyBundleIDs
		}
	}
	return nodeInfo, nil
}

func runJDCommand(cmd *cobra.Command, args []string) error {
	nodeConfig, err := loadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load node configuration: %w", err)
	}

	// Load JD configuration
	jdConfig, err := loadJDConfig(jdConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load Job Distributor configuration: %w", err)
	}

	// Map node names to configurations for quick lookup
	nodeMap := make(map[string]NodeConfig)
	for _, node := range nodeConfig.Nodes {
		nodeMap[node.Name] = node
	}

	// for each node, check if it is already connected to the Job Distributor
	// if it is, skip creating a new Job Distributor
	// Create timeout context
	timeout := time.Duration(40) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, node := range nodeConfig.Nodes {
		log.Printf("Creating Job Distributor for node: %s", node.Name)
		{
			// Create client credentials
			creds := client.Credentials{
				Email:    node.Credentials.Email,
				Password: node.Credentials.Password,
			}

			// Connect to the node
			cl, err := client.NewWithContext(ctx, node.Connection.BaseURL, creds)
			if err != nil {
				log.Printf("Failed to connect to node %s: %v", node.Name, err)
				continue
			}

			// Create the Job Distributor if it doesn't already exist
			existingDistributors, err := cl.ListJobDistributors(ctx)
			if err != nil {
				log.Printf("Failed to list Job Distributors on node %s: %v", node.Name, err)
				continue
			}

			jdFound := false
			for _, jd := range existingDistributors.FeedsManagers.GetResults() {
				if jd.PublicKey == jdConfig.PublicKey {
					log.Printf("Job Distributor %s already exists on node %s", jdConfig.Name, node.Name)
					jdFound = true
					break
				}
			}

			// Note: Assuming CreateJobDistributor exists in the client package
			// You'll need to adjust parameters based on the actual implementation
			if !jdFound {
				jobID, err := cl.CreateJobDistributor(ctx, client.JobDistributorInput{
					Name:      jdConfig.Name,
					Uri:       jdConfig.URI,
					PublicKey: jdConfig.PublicKey,
				})

				if err != nil {
					log.Printf("Failed to create Job Distributor on node %s: %v", node.Name, err)
				} else {
					log.Printf("Created Job Distributor on node %s with job ID: %s", node.Name, jobID)
				}
			}
		}
	}

	return nil
}

func newCreateAptosKeysCmd() *cobra.Command {
	var (
		nodeURL  string
		email    string
		password string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Aptos keys on Chainlink nodes",
		Long:  "Creates Aptos keys on one or more Chainlink nodes using their REST API",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			nodeConfig, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load node configuration: %w", err)
			}

			// If nodeURL is specified, create key for single node
			if nodeURL != "" {
				return createAptosKeyForNode(ctx, nodeURL, email, password)
			}

			// Otherwise, create keys for all nodes
			for _, node := range nodeConfig.Nodes {
				log.Printf("Creating Aptos key for node: %s", node.Connection.BaseURL)
				if err := createAptosKeyForNode(ctx, node.Connection.BaseURL, node.Credentials.Email, node.Credentials.Password); err != nil {
					log.Printf("Failed to create Aptos key for node %s: %v", node.Connection.BaseURL, err)
					continue
				}
				log.Printf("Successfully created Aptos key for node: %s", node.Connection.BaseURL)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&nodeURL, "node-url", "", "URL of a specific Chainlink node (optional)")
	cmd.Flags().StringVar(&email, "email", "", "Email for node authentication (required if using --node-url)")
	cmd.Flags().StringVar(&password, "password", "", "Password for node authentication (required if using --node-url)")

	return cmd
}

// Helper function to create Aptos key for a single node
func createAptosKeyForNode(ctx context.Context, nodeURL, email, password string) error {
	c, err := nodeclient.NewChainlinkClient(&nodeclient.ChainlinkConfig{
		URL:      nodeURL,
		Email:    email,
		Password: password,
	}, zerolog.Logger{})
	if err != nil {
		return fmt.Errorf("failed to create Chainlink client for %s: %w", nodeURL, err)
	}
	k, _, err := c.CreateAptosKey()
	if err != nil {
		return fmt.Errorf("failed to create Aptos key: %w", err)
	}
	fmt.Printf("Created Aptos key with ID: '%s' and account '%s'\n", k.Data.ID, k.Data.Attributes.Account)
	return nil
}

func newListAptosKeysCmd() *cobra.Command {
	var (
		nodeURL  string
		email    string
		password string
		output   string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Aptos keys on Chainlink nodes",
		Long:  "Lists all Aptos keys on one or more Chainlink nodes using their REST API",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			nodeConfig, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load node configuration: %w", err)
			}

			var allKeys []NodeAptosKeys

			// If nodeURL is specified, list keys for single node
			if nodeURL != "" {
				keys, err := listAptosKeysForNode(ctx, nodeURL, email, password)
				if err != nil {
					return err
				}
				allKeys = append(allKeys, NodeAptosKeys{
					NodeName: "single-node",
					NodeURL:  nodeURL,
					Keys:     keys,
				})
			} else {
				// Otherwise, list keys for all nodes
				for _, node := range nodeConfig.Nodes {
					log.Printf("Listing Aptos keys for node: %s", node.Connection.BaseURL)
					keys, err := listAptosKeysForNode(ctx, node.Connection.BaseURL, node.Credentials.Email, node.Credentials.Password)
					if err != nil {
						log.Printf("Failed to list Aptos keys for node %s: %v", node.Connection.BaseURL, err)
						continue
					}

					allKeys = append(allKeys, NodeAptosKeys{
						NodeName: node.Name,
						NodeURL:  node.Connection.BaseURL,
						Keys:     keys,
					})
				}
			}

			// Output results
			if output == "json" {
				jsonOutput, err := json.MarshalIndent(allKeys, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal keys to JSON: %w", err)
				}
				fmt.Println(string(jsonOutput))
			} else {
				// Human-readable output
				for _, nodeKeys := range allKeys {
					fmt.Printf("\nNode: %s (%s)\n", nodeKeys.NodeName, nodeKeys.NodeURL)
					if len(nodeKeys.Keys) == 0 {
						fmt.Println("  No Aptos keys found")
					} else {
						for _, key := range nodeKeys.Keys {
							fmt.Printf("  ID: %s\n", key.ID)
							fmt.Printf("  Account: %s\n", key.Attributes.Account)
							fmt.Printf("  Public Key: %s\n", key.Attributes.PublicKey)
							fmt.Printf("  Created: %s\n", key.Attributes.CreatedAt)
							fmt.Printf("  Updated: %s\n", key.Attributes.UpdatedAt)
							fmt.Println("  ---")
						}
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&nodeURL, "node-url", "", "URL of a specific Chainlink node (optional)")
	cmd.Flags().StringVar(&email, "email", "", "Email for node authentication (required if using --node-url)")
	cmd.Flags().StringVar(&password, "password", "", "Password for node authentication (required if using --node-url)")
	cmd.Flags().StringVarP(&output, "output", "o", "table", "Output format: table or json")

	return cmd
}

// Add these type definitions near the top with your other types

// NodeAptosKeys represents Aptos keys for a specific node
type NodeAptosKeys struct {
	NodeName string                    `json:"node_name"`
	NodeURL  string                    `json:"node_url"`
	Keys     []nodeclient.AptosKeyData `json:"keys"`
}

// Helper function to list Aptos keys for a single node
func listAptosKeysForNode(ctx context.Context, nodeURL, email, password string) ([]nodeclient.AptosKeyData, error) {
	c, err := nodeclient.NewChainlinkClient(&nodeclient.ChainlinkConfig{
		URL:      nodeURL,
		Email:    email,
		Password: password,
	}, zerolog.Logger{})
	if err != nil {
		return nil, fmt.Errorf("failed to create Chainlink client for %s: %w", nodeURL, err)
	}

	keys, _, err := c.ReadAptosKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to read Aptos keys: %w", err)
	}

	if keys == nil {
		return []nodeclient.AptosKeyData{}, nil
	}

	return keys.Data, nil
}

func newJDAptosCmd() *cobra.Command {
	var (
		nodeURL   string
		email     string
		password  string
		chainID   string
		chainName string
		adminAddr string
	)

	cmd := &cobra.Command{
		Use:   "aptos",
		Short: "Enable Aptos chain on Chainlink nodes",
		Long:  "Lists Aptos keys and enables Aptos chains using the accounts from those keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			nodeConfig, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load node configuration: %w", err)
			}
			jdConfig, err := loadJDConfig(jdConfigPath)
			if err != nil {
				return fmt.Errorf("failed to load Job Distributor configuration: %w", err)
			}

			// If nodeURL is specified, enable Aptos chain for single node
			if nodeURL != "" {
				return enableAptosChainForNode(ctx, nodeURL, email, password, chainID, chainName, adminAddr, jdConfig.PublicKey)
			}

			// Otherwise, enable Aptos chains for all nodes
			for _, node := range nodeConfig.Nodes {
				if node.Bootstrap {
					log.Printf("Skipping bootstrap node: %s", node.Name)
					continue
				}
				adminAddr = "0x0000000000000000000000000000000000000000" // Default admin address, can be overridden by command line flag
				log.Printf("Enabling Aptos chain for node: %s", node.Connection.BaseURL)
				if err := enableAptosChainForNode(ctx, node.Connection.BaseURL, node.Credentials.Email, node.Credentials.Password, chainID, chainName, adminAddr, jdConfig.PublicKey); err != nil {
					log.Printf("Failed to enable Aptos chain for node %s: %v", node.Connection.BaseURL, err)
					continue
				}
				log.Printf("Successfully enabled Aptos chain for node: %s", node.Connection.BaseURL)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&nodeURL, "node-url", "", "URL of a specific Chainlink node (optional)")
	cmd.Flags().StringVar(&email, "email", "", "Email for node authentication (required if using --node-url)")
	cmd.Flags().StringVar(&password, "password", "", "Password for node authentication (required if using --node-url)")
	cmd.Flags().StringVar(&chainID, "chain-id", "2", "Aptos chain ID (default: 2 for testnet)")
	cmd.Flags().StringVar(&chainName, "chain-name", "aptos-testnet", "Aptos chain name")
	cmd.Flags().StringVar(&adminAddr, "admin-addr", "0x0000000000000000000000000000000000000000", "Admin address for Aptos chain")
	return cmd
}

// Helper function to enable Aptos chain for a single node
func enableAptosChainForNode(ctx context.Context, nodeURL, email, password, chainID, chainName, adminAddr, jdCSAKey string) error {
	cc, err := nodeclient.NewChainlinkClient(&nodeclient.ChainlinkConfig{
		URL:      nodeURL,
		Email:    email,
		Password: password,
	}, zerolog.Logger{})
	if err != nil {
		return fmt.Errorf("failed to create Chainlink client for %s: %w", nodeURL, err)
	}

	creds := client.Credentials{
		Email:    email,
		Password: password,
	}
	// Connect to the node
	gc, err := client.NewWithContext(ctx, nodeURL, creds)
	if err != nil {
		log.Printf("Failed to connect to node %s: %v", nodeURL, err)
		return fmt.Errorf("failed to connect to node %s: %w", nodeURL, err)
	}

	// First, list existing Aptos keys to get the account
	keys, _, err := cc.ReadAptosKeys()
	if err != nil {
		return fmt.Errorf("failed to read Aptos keys: %w", err)
	}

	if keys == nil || len(keys.Data) == 0 {
		return fmt.Errorf("no Aptos keys found on node %s", nodeURL)
	}

	// Use the first Aptos key's account
	aptosAccount := keys.Data[0].Attributes.Account
	log.Printf("Using Aptos account: %s", aptosAccount)

	jds, err := gc.ListJobDistributors(ctx)
	if err != nil {
		return fmt.Errorf("failed to list job distributors: %w", err)
	}

	var jdID string
	for _, jd := range jds.FeedsManagers.GetResults() {
		if jd.PublicKey == "" {
			log.Printf("Warning: Job Distributor %s has no public key", jd.Name)
			continue
		}

		if jd.GetPublicKey() == jdCSAKey {
			jdID = jd.Id
			break
		}
	}

	// now create the chain configuration. we need to list the existing keys
	info, err := nodeInfo(ctx, NodeConfig{
		Name:        "Aptos Chain Config",
		Connection:  Connection{BaseURL: nodeURL},
		Credentials: Credentials{Email: email, Password: password},
	})
	if err != nil {
		return fmt.Errorf("failed to get node info: %w", err)
	}

	if len(info.OCR2KeyBundles) == 0 {
		return fmt.Errorf("no OCR2 key bundles found for Aptos chain on node %s", nodeURL)
	}
	// get the aptos key bundle ID
	var aptosKeyBundleID string
	for _, bundle := range info.OCR2KeyBundles {
		if bundle.ChainType == client.OCR2ChainTypeAptos {
			aptosKeyBundleID = bundle.ID
			break
		}
	}
	if aptosKeyBundleID == "" {
		return fmt.Errorf("no OCR2 key bundle found for Aptos chain on node %s", nodeURL)
	}

	log.Printf("Using Aptos OCR2 Key Bundle ID: %s", aptosKeyBundleID)

	// Check if Aptos chain already exists
	r, err := gc.CreateJobDistributorChainConfig(ctx, client.JobDistributorChainConfigInput{
		JobDistributorID: jdID,
		ChainID:          chainID,
		ChainType:        client.OCR2ChainTypeAptos,
		AccountAddr:      aptosAccount,
		AdminAddr:        adminAddr,
		Ocr2Enabled:      true,
		Ocr2P2PPeerID:    info.P2PPeerID,
		Ocr2KeyBundleID:  aptosKeyBundleID,
		Ocr2Plugins:      `{"commit":false,"execute":false,"median":false,"mercury":false}`,
	})
	if err != nil {
		return fmt.Errorf("failed to create Aptos chain configuration: %w", err)
	}

	log.Printf("Created Aptos node with name: %s", r)

	fmt.Printf("Successfully enabled Aptos chain '%s' (ID: %s) for node %s using account %s\n",
		chainName, chainID, nodeURL, aptosAccount)

	return nil
}

func newJDAcceptCmd() *cobra.Command {
	var (
		nodeURL    string
		email      string
		password   string
		proposalID string
		force      bool
		all        bool
	)

	cmd := &cobra.Command{
		Use:   "accept",
		Short: "Accept job proposals on Chainlink nodes",
		Long:  "Accepts job proposals from Job Distributors using the GraphQL API",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			nodeConfig, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load node configuration: %w", err)
			}

			// If nodeURL is specified, accept proposals for single node
			if nodeURL != "" {
				return acceptJobProposalsForNode(ctx, nodeURL, email, password, proposalID, force, all)
			}

			// Otherwise, accept proposals for all nodes (excluding bootstrap nodes)
			for _, node := range nodeConfig.Nodes {
				if node.Bootstrap {
					log.Printf("Skipping bootstrap node: %s", node.Name)
					continue
				}

				log.Printf("Accepting job proposals for node: %s", node.Connection.BaseURL)
				if err := acceptJobProposalsForNode(ctx, node.Connection.BaseURL, node.Credentials.Email, node.Credentials.Password, proposalID, force, all); err != nil {
					log.Printf("Failed to accept job proposals for node %s: %v", node.Connection.BaseURL, err)
					continue
				}
				log.Printf("Successfully processed job proposals for node: %s", node.Connection.BaseURL)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&nodeURL, "node-url", "", "URL of a specific Chainlink node (optional)")
	cmd.Flags().StringVar(&email, "email", "", "Email for node authentication (required if using --node-url)")
	cmd.Flags().StringVar(&password, "password", "", "Password for node authentication (required if using --node-url)")
	cmd.Flags().StringVar(&proposalID, "proposal-id", "", "Specific proposal ID to accept (optional)")
	cmd.Flags().BoolVar(&force, "force", false, "Force acceptance even if proposal is already approved")
	cmd.Flags().BoolVar(&all, "all", false, "Accept all pending proposals")

	return cmd
}

// Helper function to accept job proposals for a single node
func acceptJobProposalsForNode(ctx context.Context, nodeURL, email, password, proposalID string, force, acceptAll bool) error {
	creds := client.Credentials{
		Email:    email,
		Password: password,
	}
	gc, err := client.NewWithContext(ctx, nodeURL, creds)
	if err != nil {
		return fmt.Errorf("failed to connect to node %s: %w", nodeURL, err)
	}

	return acceptSingleProposal(ctx, gc, proposalID, force)
}

// Accept a single proposal by ID
func acceptSingleProposal(ctx context.Context, gc client.Client, proposalID string, force bool) error {
	log.Printf("Accepting job proposal with ID: %s", proposalID)

	// Get the proposal details first
	proposal, err := gc.GetJobProposal(ctx, proposalID)
	if err != nil {
		return fmt.Errorf("failed to get job proposal %s: %w", proposalID, err)
	}
	if proposal == nil {
		return fmt.Errorf("job proposal %s not found", proposalID)
	}

	log.Printf("Proposal: %s - Status: %s", proposal.LatestSpec.Definition, proposal.LatestSpec.Status)

	// Check if already approved (unless force is true)
	if !force && proposal.LatestSpec.Status == "APPROVED" {
		log.Printf("Proposal %s is already approved, skipping", proposalID)
		return nil
	}

	// Accept the proposal
	result, err := gc.ApproveJobProposalSpec(ctx, proposalID, force)
	if err != nil {
		return fmt.Errorf("failed to approve job proposal %s: %w", proposalID, err)
	}
	if result != nil {
		log.Printf("Successfully approved job proposal %s", proposalID)
		fmt.Printf("Approved job proposal: %s\n", proposalID)
	}

	return nil
}
