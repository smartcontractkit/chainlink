// Command generate_file_source creates a workflow metadata JSON file for the file-based workflow source.
// This tool generates the correct workflowID based on the binary, config, owner, and name.
//
// The binary file should be in .br.b64 format (base64-encoded brotli-compressed WASM).
// This is the format used by the workflow deploy command.
//
// Usage:
//
//	go run ./cmd/generate_file_source \
//	  --binary /path/to/workflow.br.b64 \
//	  --config /path/to/config.yaml \
//	  --name my-workflow \
//	  --owner f39fd6e51aad88f6f4ce6ab8827279cfffb92266 \
//	  --output /tmp/workflows_metadata.json \
//	  --don-family workflow
package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/andybalholm/brotli"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
)

type FileWorkflowMetadata struct {
	WorkflowID   string `json:"workflow_id"`
	Owner        string `json:"owner"`
	CreatedAt    uint64 `json:"created_at"`
	Status       uint8  `json:"status"`
	WorkflowName string `json:"workflow_name"`
	BinaryURL    string `json:"binary_url"`
	ConfigURL    string `json:"config_url"`
	Tag          string `json:"tag"`
	DonFamily    string `json:"don_family"`
}

type FileWorkflowSourceData struct {
	Workflows []FileWorkflowMetadata `json:"workflows"`
}

func main() {
	var (
		binaryPath      string
		configPath      string
		workflowName    string
		owner           string
		outputPath      string
		donFamily       string
		tag             string
		binaryURLPrefix string
		configURLPrefix string
		status          int
	)

	flag.StringVar(&binaryPath, "binary", "", "Path to the compiled workflow binary (required)")
	flag.StringVar(&configPath, "config", "", "Path to the workflow config file (optional)")
	flag.StringVar(&workflowName, "name", "file-source-workflow", "Workflow name")
	flag.StringVar(&owner, "owner", "f39fd6e51aad88f6f4ce6ab8827279cfffb92266", "Workflow owner address (hex without 0x)")
	flag.StringVar(&outputPath, "output", "/tmp/workflows_metadata.json", "Output path for the JSON file")
	flag.StringVar(&donFamily, "don-family", "workflow", "DON family name")
	flag.StringVar(&tag, "tag", "v1.0.0", "Workflow tag")
	flag.StringVar(&binaryURLPrefix, "binary-url-prefix", "file:///home/chainlink/workflows/", "URL prefix for binary (will append filename)")
	flag.StringVar(&configURLPrefix, "config-url-prefix", "file:///home/chainlink/workflows/", "URL prefix for config (will append filename)")
	flag.IntVar(&status, "status", 0, "Workflow status (0=active, 1=paused)")
	flag.Parse()

	if binaryPath == "" {
		fmt.Println("Error: --binary is required")
		flag.Usage()
		os.Exit(1)
	}

	// Read binary file
	binaryRaw, err := os.ReadFile(binaryPath)
	if err != nil {
		fmt.Printf("Error reading binary file: %v\n", err)
		os.Exit(1)
	}

	// Decompress binary if it's in .br.b64 format
	var binary []byte
	if strings.HasSuffix(binaryPath, ".br.b64") {
		// Base64 decode
		decoded, decodeErr := base64.StdEncoding.DecodeString(string(binaryRaw))
		if decodeErr != nil {
			fmt.Printf("Error base64 decoding binary: %v\n", decodeErr)
			os.Exit(1)
		}
		// Brotli decompress
		reader := brotli.NewReader(strings.NewReader(string(decoded)))
		var decompressErr error
		binary, decompressErr = io.ReadAll(reader)
		if decompressErr != nil {
			fmt.Printf("Error brotli decompressing binary: %v\n", decompressErr)
			os.Exit(1)
		}
		fmt.Printf("Decompressed binary from %d bytes (compressed) to %d bytes (WASM)\n", len(binaryRaw), len(binary))
	} else {
		binary = binaryRaw
	}

	// Read config file (optional)
	var config []byte
	if configPath != "" {
		config, err = os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Error reading config file: %v\n", err)
			os.Exit(1)
		}
	}

	// Decode owner
	ownerBytes, err := hex.DecodeString(owner)
	if err != nil {
		fmt.Printf("Error decoding owner hex: %v\n", err)
		os.Exit(1)
	}

	// Generate workflow ID
	workflowID, err := pkgworkflows.GenerateWorkflowID(ownerBytes, workflowName, binary, config, "")
	if err != nil {
		fmt.Printf("Error generating workflow ID: %v\n", err)
		os.Exit(1)
	}

	// Get binary and config filenames - use .br.b64 for compressed binary
	binaryFilename := "file_source_workflow.br.b64"
	configFilename := "file_source_config.json"

	// Build the metadata
	now := time.Now().Unix()
	var createdAt uint64
	if now >= 0 {
		createdAt = uint64(now) // #nosec G115 -- time is always positive
	}
	var statusUint8 uint8
	if status >= 0 && status <= 255 {
		statusUint8 = uint8(status) // #nosec G115 -- status is validated in range
	}
	metadata := FileWorkflowSourceData{
		Workflows: []FileWorkflowMetadata{
			{
				WorkflowID:   hex.EncodeToString(workflowID[:]),
				Owner:        owner,
				CreatedAt:    createdAt,
				Status:       statusUint8,
				WorkflowName: workflowName,
				BinaryURL:    binaryURLPrefix + binaryFilename,
				ConfigURL:    configURLPrefix + configFilename,
				Tag:          tag,
				DonFamily:    donFamily,
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// Write to output file
	if err := os.WriteFile(outputPath, jsonData, 0600); err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated workflow metadata file: %s\n", outputPath)
	fmt.Printf("Workflow ID: %s\n", hex.EncodeToString(workflowID[:]))
	fmt.Printf("Workflow Name: %s\n", workflowName)
	fmt.Printf("Owner: %s\n", owner)
	fmt.Printf("DON Family: %s\n", donFamily)
	fmt.Printf("\nTo use this workflow:\n")
	fmt.Printf("1. Copy the binary to Docker containers: docker cp %s workflow-node1:/home/chainlink/workflows/%s\n", binaryPath, binaryFilename)
	if configPath != "" {
		fmt.Printf("2. Copy the config to Docker containers: docker cp %s workflow-node1:/home/chainlink/workflows/%s\n", configPath, configFilename)
	}
	fmt.Printf("3. Copy the metadata JSON to Docker containers: docker cp %s workflow-node1:/tmp/workflows_metadata.json\n", outputPath)
	fmt.Printf("4. Repeat steps 1-3 for all workflow nodes\n")
	fmt.Printf("5. Wait for syncer to pick up the workflow (default 12 second interval)\n")
}
