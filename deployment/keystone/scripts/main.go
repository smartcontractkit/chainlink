package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/deployment/environment/web/sdk/client"
)

// Config holds the TOML configuration with multiple nodes
type Config struct {
	Nodes []NodeConfig `toml:"nodes"`
}

// NodeConfig holds the configuration for a single node
type NodeConfig struct {
	Name        string `toml:"name"`
	Bootstrap   bool   `toml:"bootstrap"`
	Credentials struct {
		Email    string `toml:"email"`
		Password string `toml:"password"`
	} `toml:"credentials"`
	Connection struct {
		BaseURL string `toml:"base_url"`
	} `toml:"connection"`
}

// JobDistributorConfig defines a job distributor configuration
type JobDistributorConfig struct {
	Name      string `toml:"name"`
	Uri       string `toml:"uri"`
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
	jdCmd := &cobra.Command{
		Use:   "jd",
		Short: "Create Job Distributors",
		Long:  `Creates Job Distributors for nodes based on configuration`,
		RunE:  runJDCommand,
	}
	jdCmd.Flags().StringVarP(&jdConfigPath, "jd-config", "j", "jd.toml", "Path to Job Distributor TOML config file")

	// Add commands to root
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(jdCmd)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runInfoCommand(cmd *cobra.Command, args []string) error {
	// Load configuration
	var config Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("Config file not found: %s", configPath)
	}

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return fmt.Errorf("Failed to decode config file: %v", err)
	}

	// Check if we have any nodes configured
	if len(config.Nodes) == 0 {
		return fmt.Errorf("No nodes configured in the config file")
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

		// Initialize node info with known data
		nodeInfo := NodeInfo{
			Name: nodeConfig.Name,
			URL:  nodeConfig.Connection.BaseURL,
		}

		// Create timeout context
		timeout := time.Duration(10) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		// Create client credentials
		creds := client.Credentials{
			Email:    nodeConfig.Credentials.Email,
			Password: nodeConfig.Credentials.Password,
		}

		// Try to create client and fetch keys
		cl, err := client.NewWithContext(ctx, nodeConfig.Connection.BaseURL, creds)
		if err != nil {
			log.Printf("Failed to connect to node %s: %v", nodeConfig.Name, err)
			cancel()
			continue
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
					log.Printf("Failed to create OCR2 key bundle for Aptos on node %s: %v", nodeConfig.Name, err)
					panic(fmt.Sprintf("Failed to create OCR2 key bundle for Aptos on node %s: %v", nodeConfig.Name, err))
				} else {
					log.Printf("Created OCR2 key bundle for Aptos on node %s with ID: %s", nodeConfig.Name, aptosKeyBundleID)
					// Retry fetching OCR2 key bundles
					ocrKeyBundleIDs, err = cl.ListOCR2KeyBundles(ctx)
					if err == nil {
						nodeInfo.OCR2KeyBundles = ocrKeyBundleIDs
					} else {
						log.Printf("Failed to fetch OCR2 key bundles after creating new one on node %s: %v", nodeConfig.Name, err)
						panic(fmt.Sprintf("Failed to fetch OCR2 key bundles after creating new one on node %s: %v", nodeConfig.Name, err))
					}
				}
			}
			// Add node to the list
			nodes = append(nodes, nodeInfo)
		}
		if nodeConfig.Bootstrap {
			p := strings.TrimPrefix(nodeInfo.P2PPeerID, "p2p_")
			bootstrapNodeInfo := BootstrapNodeInfo{
				Name:         nodeConfig.Name,
				URL:          nodeConfig.Connection.BaseURL,
				CSAPublicKey: nodeInfo.CSAPublicKey,
				P2PPeerID:    nodeInfo.P2PPeerID,
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
		return fmt.Errorf("Failed to convert nodes to JSON: %v", err)
	}
	// Convert bootstrap nodes to JSON
	bootstrapJSONOutput, err := json.MarshalIndent(bootstrapNodes, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to convert bootstrap nodes to JSON: %v", err)
	}

	// Write output to file or stdout
	if outputPath == "" {
		fmt.Println(string(jsonOutput))
	} else {
		err = os.WriteFile(outputPath, jsonOutput, 0644)
		if err != nil {
			return fmt.Errorf("Failed to write output to file: %v", err)
		}
		log.Printf("Output written to %s", outputPath)
	}

	// Write bootstrap output to file or stdout
	if bootstrapOutputPath == "" {
		fmt.Println(string(bootstrapJSONOutput))
	} else {
		err = os.WriteFile(bootstrapOutputPath, bootstrapJSONOutput, 0644)
		if err != nil {
			return fmt.Errorf("Failed to write bootstrap output to file: %v", err)
		}
		log.Printf("Bootstrap output written to %s", bootstrapOutputPath)
	}

	return nil
}

func runJDCommand(cmd *cobra.Command, args []string) error {
	// Load node configuration
	var nodeConfig Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("Node config file not found: %s", configPath)
	}
	if _, err := toml.DecodeFile(configPath, &nodeConfig); err != nil {
		return fmt.Errorf("Failed to decode node config file: %v", err)
	}

	// Load JD configuration
	var jdConfig JobDistributorConfig
	if _, err := os.Stat(jdConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("Job Distributor config file not found: %s", jdConfigPath)
	}
	if _, err := toml.DecodeFile(jdConfigPath, &jdConfig); err != nil {
		return fmt.Errorf("Failed to decode Job Distributor config file: %v", err)
	}

	// Map node names to configurations for quick lookup
	nodeMap := make(map[string]NodeConfig)
	for _, node := range nodeConfig.Nodes {
		nodeMap[node.Name] = node
	}

	// for each node, check if it is already connected to the Job Distributor
	// if it is, skip creating a new Job Distributor
	for _, node := range nodeConfig.Nodes {

		log.Printf("Creating Job Distributor for node: %s", node.Name)

		{
			// Create timeout context
			timeout := time.Duration(20) * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

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
					Uri:       jdConfig.Uri,
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
