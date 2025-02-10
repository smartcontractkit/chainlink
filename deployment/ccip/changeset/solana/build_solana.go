package changeset

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

type ProgramOutput struct {
	Name        deployment.ContractType
	Key         solana.PublicKey
	SoPath      string
	KeypairPath string
}

// Configuration
const (
	repoURL   = "https://github.com/smartcontractkit/chainlink-ccip.git"
	revision  = "2e4e27e1d64f8633b4742100a395936c13614fb8" // Dec 27 2023
	cloneDir  = "./temp-repo"
	anchorDir = "chains/solana/contracts" // Path to the Anchor project within the repo
	deployDir = "chains/solana/contracts/target/deploy"
)

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
func cloneRepo(revision string) error {
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
	_, err = runCommand("git", []string{"checkout"}, cloneDir)
	if err != nil {
		return fmt.Errorf("failed to checkout revision %s: %w", revision, err)
	}

	return nil
}

// Generate keys using Anchor
func generateKeys() error {
	log.Println("Generating keys with Anchor...")
	anchorPath := filepath.Join(cloneDir, anchorDir)
	_, err := runCommand("anchor", []string{"keys", "list"}, anchorPath)
	if err != nil {
		return fmt.Errorf("anchor key generation failed: %w", err)
	}

	return err
}

// Replace keys in Rust files
func replaceKeys() error {
	log.Println("Generating keys with Anchor...")
	anchorPath := filepath.Join(cloneDir, anchorDir)
	_, err := runCommand("anchor", []string{"keys", "sync"}, anchorPath)
	if err != nil {
		return fmt.Errorf("anchor key replacement failed: %w", err)
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

type BuildSolanaConfig struct {
	ChainSelector uint64
	GitCommitSha  string
}

func BuildSolana(e deployment.Environment, config BuildSolanaConfig) (deployment.ChangesetOutput, error) {
	_, ok := e.SolChains[config.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
	}
	family, err := chainsel.GetSelectorFamily(config.ChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	if family != chainsel.FamilySolana {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain is not solana chain %d", config.ChainSelector)
	}

	newAddresses := deployment.NewMemoryAddressBook()
	// Step 1: Clone the repository
	if err := cloneRepo(config.GitCommitSha); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error cloning repo: %v", err)
	}

	// Step 2: Generate keys using Anchor
	err = generateKeys()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error generating keys: %v", err)
	}

	// Step 3: Replace keys in Rust files
	if err := replaceKeys(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error replacing keys: %v", err)
	}

	// Step 4: Build the project with Anchor
	if err := buildProject(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error building project: %v", err)
	}

	deployDir := filepath.Join(cloneDir, deployDir)
	files, err := os.ReadDir(deployDir)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to read deploy directory: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			name := strings.TrimSuffix(file.Name(), ".so")
			var contractType deployment.ContractType
			switch name {
			case "ccip_router":
				contractType = changeset.Router
			case "token_pool":
				contractType = changeset.TokenPool
			case "fee_quoter":
				contractType = changeset.FeeQuoter
			case "ccip_offramp":
				contractType = changeset.OffRamp
			case "test_ccip_receiver":
				contractType = changeset.Receiver
			default:
				e.Logger.Warnf("Unknown contract type: %s", name)
				continue
			}
			keypairPath := filepath.Join(deployDir, name+"-keypair.json")
			key, err := os.ReadFile(keypairPath)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to read keypair file: %w", err)
			}
			publicKey := solana.PublicKeyFromBytes(key)
			err = newAddresses.Save(config.ChainSelector, publicKey.String(), deployment.NewTypeAndVersion(contractType, deployment.Version1_0_0))
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
		}
	}

	return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
}
