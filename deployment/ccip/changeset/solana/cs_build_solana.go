package solana

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
)

var _ deployment.ChangeSet[BuildSolanaConfig] = BuildSolanaChangeset

// Configuration
const (
	repoURL   = "https://github.com/smartcontractkit/chainlink-ccip.git"
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
func cloneRepo(e deployment.Environment, revision string) error {
	// Remove the clone directory if it already exists
	if _, err := os.Stat(cloneDir); !os.IsNotExist(err) {
		os.RemoveAll(cloneDir)
	}

	e.Logger.Debugw("Cloning repository", "url", repoURL, "revision", revision)
	_, err := runCommand("git", []string{"clone", repoURL, cloneDir}, ".")
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	e.Logger.Debugw("Checking out revision", "revision", revision)
	_, err = runCommand("git", []string{"checkout", revision}, cloneDir)
	if err != nil {
		return fmt.Errorf("failed to checkout revision %s: %w", revision, err)
	}

	return nil
}

// Replace keys in Rust files
func replaceKeys(e deployment.Environment) error {
	solanaDir := filepath.Join(cloneDir, anchorDir, "..")
	e.Logger.Debugw("Replacing keys", "solanaDir", solanaDir)
	output, err := runCommand("make", []string{"docker-update-contracts"}, solanaDir)
	if err != nil {
		fmt.Println(output)
		return fmt.Errorf("anchor key replacement failed: %s %w", output, err)
	}
	return nil
}

func copyFile(srcFile string, destDir string) error {
	output, err := runCommand("cp", []string{srcFile, destDir}, ".")
	if err != nil {
		return fmt.Errorf("failed to copy file: %s %w", output, err)
	}
	return nil
}

// Build the project with Anchor
func buildProject(e deployment.Environment) error {
	solanaDir := filepath.Join(cloneDir, anchorDir, "..")
	e.Logger.Debugw("Building project", "solanaDir", solanaDir)
	output, err := runCommand("make", []string{"docker-build-contracts"}, solanaDir)
	if err != nil {
		return fmt.Errorf("anchor build failed: %s %w", output, err)
	}
	return nil
}

type BuildSolanaConfig struct {
	ChainSelector        uint64
	GitCommitSha         string
	DestinationDir       string
	IsUpgrade            bool
	CleanDestinationDir  bool
	CreateDestinationDir bool
}

func BuildSolanaChangeset(e deployment.Environment, config BuildSolanaConfig) (deployment.ChangesetOutput, error) {
	_, ok := e.SolChains[config.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", config.ChainSelector)
	}
	family, err := chainsel.GetSelectorFamily(config.ChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	if family != chainsel.FamilySolana {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain is not solana chain %d", config.ChainSelector)
	}

	// Clone the repository
	if err := cloneRepo(e, config.GitCommitSha); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error cloning repo: %w", err)
	}

	// Upgrades don't need to generate keys, we upgrade the program in place
	if !config.IsUpgrade {
		if err := replaceKeys(e); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error replacing keys: %w", err)
		}
	}

	// Build the project with Anchor
	if err := buildProject(e); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error building project: %w", err)
	}

	if config.CleanDestinationDir {
		e.Logger.Debugw("Cleaning destination dir", "destinationDir", config.DestinationDir)
		if err := os.RemoveAll(config.DestinationDir); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error cleaning build folder: %w", err)
		}
		e.Logger.Debugw("Creating destination dir", "destinationDir", config.DestinationDir)
		err = os.MkdirAll(config.DestinationDir, os.ModePerm)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create build directory: %w", err)
		}
	} else if config.CreateDestinationDir {
		e.Logger.Debugw("Creating destination dir", "destinationDir", config.DestinationDir)
		err := os.MkdirAll(config.DestinationDir, os.ModePerm)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create build directory: %w", err)
		}
	}

	deployFilePath := filepath.Join(cloneDir, deployDir)
	e.Logger.Debugw("Reading deploy directory", "deployFilePath", deployFilePath)
	files, err := os.ReadDir(deployFilePath)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to read deploy directory: %w", err)
	}

	for _, file := range files {
		filePath := filepath.Join(deployFilePath, file.Name())
		e.Logger.Debugw("Copying file", "filePath", filePath, "destinationDir", config.DestinationDir)
		err := copyFile(filePath, config.DestinationDir)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to copy file: %w", err)
		}
	}
	return deployment.ChangesetOutput{}, nil
}
