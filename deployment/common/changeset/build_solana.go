package changeset

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/deployment"
)

type ProgramOutput struct {
	Name        deployment.ContractType
	Key         string
	SoPath      string
	KeypairPath string
}

// Configuration
const (
	repoURL   = "https://github.com/smartcontractkit/chainlink-ccip.git"
	revision  = "2e4e27e1d64f8633b4742100a395936c13614fb8" // Dec 27 2023
	cloneDir  = "./temp-repo"
	anchorDir = "chains/solana/contracts" // Path to the Anchor project within the repo
)

// Map program names to their Rust file paths (relative to the Anchor project root)
var programToFileMap = map[string]string{
	"access_controller":         "programs/access-controller/src/lib.rs",
	"timelock":                  "programs/timelock/src/lib.rs",
	"ccip_router":               "programs/ccip-router/src/lib.rs",
	"token_pool":                "programs/token-pool/src/lib.rs",
	"external_program_cpi_stub": "programs/external-program-cpi-stub/src/lib.rs",
	"mcm":                       "programs/mcm/src/lib.rs",
	"ccip_receiver":             "programs/ccip-receiver/src/lib.rs",
	"ccip_invalid_receiver":     "programs/ccip-invalid-receiver/src/lib.rs",
}

// Run a command in a specific directory
func runCommand(command string, args []string, workDir string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = workDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}
	return stdout.String(), nil
}

// Clone and checkout the specific revision of the repo
func cloneRepo() error {
	// Remove the clone directory if it already exists
	if _, err := os.Stat(cloneDir); !os.IsNotExist(err) {
		log.Println("Temporary directory exists, removing...")
		os.RemoveAll(cloneDir)
	}

	log.Println("Cloning repository...")
	_, err := runCommand("git", []string{"clone", repoURL, cloneDir}, ".")
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Check out the specific revision (branch or commit SHA)
	log.Println("Checking out specific revision...")
	_, err = runCommand("git", []string{"checkout", revision}, cloneDir)
	if err != nil {
		return fmt.Errorf("failed to checkout revision %s: %w", revision, err)
	}

	return nil
}

// Generate keys using Anchor
func generateKeys() (map[string]string, error) {
	log.Println("Generating keys with Anchor...")
	anchorPath := filepath.Join(cloneDir, anchorDir)
	output, err := runCommand("anchor", []string{"keys", "list"}, anchorPath)
	if err != nil {
		return nil, fmt.Errorf("anchor key generation failed: %w", err)
	}

	keys := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			program := strings.TrimSpace(parts[0])
			key := strings.TrimSpace(parts[1])
			keys[program] = key
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error parsing keys: %w", err)
	}
	return keys, nil
}

// Replace keys in Rust files
func replaceKeys(keys map[string]string) error {
	log.Println("Replacing keys in Rust files...")
	for program, filePath := range programToFileMap {
		key, exists := keys[program]
		if !exists {
			log.Printf("Warning: No key found for program %s. Skipping.\n", program)
			continue
		}

		fullPath := filepath.Join(cloneDir, anchorDir, filePath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", fullPath, err)
		}

		// Replace declare_id!("..."); with the new key
		updatedContent := regexp.MustCompile(`declare_id!\(".*?"\);`).ReplaceAllString(string(content), fmt.Sprintf(`declare_id!("%s");`, key))
		err = os.WriteFile(fullPath, []byte(updatedContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to write updated keys to file %s: %w", fullPath, err)
		}
		log.Printf("Updated key for program %s in file %s\n", program, filePath)
	}
	return nil
}

// Build the project with Anchor
func buildProject() error {
	log.Println("Building project with Anchor...")
	anchorPath := filepath.Join(cloneDir, anchorDir)
	_, err := runCommand("anchor", []string{"build"}, anchorPath)
	if err != nil {
		return fmt.Errorf("anchor build failed: %w", err)
	}
	return nil
}

func collectProgramOutputs(deployDir string, keys map[string]string) ([]ProgramOutput, error) {
	outputs := []ProgramOutput{}
	files, err := os.ReadDir(deployDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read deploy directory: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			name := strings.TrimSuffix(file.Name(), ".so")
			key := keys[name]
			outputs = append(outputs, ProgramOutput{
				Name:        deployment.ContractType(name),
				Key:         key,
				SoPath:      filepath.Join(deployDir, file.Name()),
				KeypairPath: filepath.Join(deployDir, name+"-keypair.json"),
			})
		}
	}

	return outputs, nil
}

func BuildSolana(e deployment.Environment, chains []uint64) (deployment.ChangesetOutput, error) {
	for _, chain := range chains {
		_, ok := e.SolChains[chain]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
		}
	}
	newAddresses := deployment.NewMemoryAddressBook()
	// Step 1: Clone the repository
	if err := cloneRepo(); err != nil {
		e.Logger.Fatalf("Error: %v", err)
	}

	// Step 2: Generate keys using Anchor
	keys, err := generateKeys()
	if err != nil {
		e.Logger.Fatalf("Error: %v", err)
	}

	// Step 3: Replace keys in Rust files
	if err := replaceKeys(keys); err != nil {
		e.Logger.Fatalf("Error: %v", err)
	}

	// Step 4: Build the project with Anchor
	if err := buildProject(); err != nil {
		e.Logger.Fatalf("Error: %v", err)
	}

	deployDir := filepath.Join(cloneDir, "chains/solana/contracts/target/deploy")
	outputs, err := collectProgramOutputs(deployDir, keys)
	if err != nil {
		e.Logger.Fatalf("Error collecting program outputs: %v", err)
	}

	for _, output := range outputs {
		for _, chain := range chains {
			family, err := chainsel.GetSelectorFamily(chain)
			if err != nil {
				return deployment.ChangesetOutput{AddressBook: newAddresses}, err
			}
			if family == chainsel.FamilySolana {
				newAddresses.Save(chain, output.Key, deployment.NewTypeAndVersion(output.Name, deployment.Version1_6_0_dev))
			} else {
				e.Logger.Errorf("Unknown chain family %s", family)
			}
			if err != nil {
				return deployment.ChangesetOutput{AddressBook: newAddresses}, err
			}
		}
	}

	return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
}
