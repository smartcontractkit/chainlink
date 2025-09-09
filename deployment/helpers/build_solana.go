package helpers

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	ks_shared "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/shared"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// TODO add forwarder
var programToVanityKey = map[cldf.ContractType]string{
	shared.Router:    "Ccip",
	shared.FeeQuoter: "FeeQ",
	shared.OffRamp:   "off",
	shared.RMNRemote: "Rmn",
}

type LocalBuildConfig struct {
	BuildLocally         bool
	CleanDestinationDir  bool
	CreateDestinationDir bool
	// Forces re-clone of git directory. Useful for forcing regeneration of keys
	CleanGitDir bool
	// When building locally, this will be used to replace the keys in the Rust files
	GenerateVanityKeys bool
	UpgradeKeys        map[cldf.ContractType]string
}

type BuildSolanaConfig struct {
	GitCommitSha string
	// when running using CLD, this should be same as the secret (solana_program_path) or envvar (SOLANA_PROGRAM_PATH)
	DestinationDir string
	LocalBuild     LocalBuildConfig
}

type DomainParams struct {
	RepoURL, CloneDir, AnchorDir, DeployDir, BuildCmd, ReplaceKeysCmd string

	// Syncers are used in case if some programs need to get their keys synchronized before build
	Syncers []func() error

	// used to replace keys in solana programs for upgrades
	ProgramFilesView map[cldf.ContractType]string
}

// Run a command in a specific directory
func RunCommand(command string, args []string, workDir string) (string, error) {
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
func cloneRepo(e cldf.Environment, revision, cloneDir, repoURL string, forceClean bool) error {
	// Check if the repository already exists
	if forceClean {
		e.Logger.Debugw("Cleaning repository", "dir", cloneDir)
		if err := os.RemoveAll(cloneDir); err != nil {
			return fmt.Errorf("failed to clean repository: %w", err)
		}
	}
	if _, err := os.Stat(filepath.Join(cloneDir, ".git")); err == nil {
		e.Logger.Debugw("Repository already exists, discarding local changes and updating", "dir", cloneDir)

		// Discard any local changes
		_, err := RunCommand("git", []string{"reset", "--hard"}, cloneDir)
		if err != nil {
			return fmt.Errorf("failed to discard local changes: %w", err)
		}

		// Fetch the latest changes from the remote
		_, err = RunCommand("git", []string{"fetch", "origin"}, cloneDir)
		if err != nil {
			return fmt.Errorf("failed to fetch origin: %w", err)
		}
	} else {
		// Repository does not exist, clone it
		e.Logger.Debugw("Cloning repository", "url", repoURL, "revision", revision)
		_, err := RunCommand("git", []string{"clone", repoURL, cloneDir}, ".")
		if err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	}

	e.Logger.Debugw("Checking out revision", "revision", revision)
	_, err := RunCommand("git", []string{"checkout", revision}, cloneDir)
	if err != nil {
		return fmt.Errorf("failed to checkout revision %s: %w", revision, err)
	}

	return nil
}

// Replace keys in Rust files
func replaceKeys(e cldf.Environment, cloneDir, anchorDir, replaceKeyCmd string) error {
	solanaDir := filepath.Join(cloneDir, anchorDir, "..")
	e.Logger.Debugw("Replacing keys", "solanaDir", solanaDir)
	output, err := RunCommand("make", []string{replaceKeyCmd}, solanaDir)
	if err != nil {
		return fmt.Errorf("anchor key replacement failed: %s %w", output, err)
	}
	return nil
}

func replaceKeysForUpgrade(e cldf.Environment, cloneDir, anchorDir string, programToFileMap map[cldf.ContractType]string, keys map[cldf.ContractType]string) error {
	e.Logger.Debug("Replacing keys in Rust files...")
	for program, key := range keys {
		programStr := string(program)
		filePath, exists := programToFileMap[program]
		if !exists {
			return fmt.Errorf("no file path found for program %s", programStr)
		}

		fullPath := filepath.Join(cloneDir, anchorDir, filePath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", fullPath, err)
		}

		// Replace declare_id!("..."); with the new key
		updatedContent := regexp.MustCompile(`declare_id!\(".*?"\);`).ReplaceAllString(string(content), fmt.Sprintf(`declare_id!("%s");`, key))
		err = os.WriteFile(fullPath, []byte(updatedContent), 0600)
		if err != nil {
			return fmt.Errorf("failed to write updated keys to file %s: %w", fullPath, err)
		}
		e.Logger.Debugf("Updated key for program %s in file %s\n", programStr, filePath)
	}
	return nil
}

func generateVanityKeys(e cldf.Environment, cloneDir, deployDir string, keys map[cldf.ContractType]string) error {
	e.Logger.Debug("Generating vanity keys...")
	for program, prefix := range programToVanityKey {
		_, exists := keys[program]
		if exists {
			fmt.Printf("Vanity key for program %s already exists, skipping generation.", program)
			continue
		}

		// Construct command arguments
		args := []string{"grind", "--starts-with", prefix + ":1"}

		// Run command using helper function
		output, err := RunCommand("solana-keygen", args, "./")
		if err != nil {
			return fmt.Errorf("failed to generate vanity key for program %s: %w", program, err)
		}

		// Parse output for JSON filename
		scanner := bufio.NewScanner(strings.NewReader(output))
		jsonFilePattern := regexp.MustCompile(`Wrote keypair to (.*\.json)`) // Regex to match output
		var jsonFilePath string

		for scanner.Scan() {
			line := scanner.Text()
			matches := jsonFilePattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				jsonFilePath = matches[1]
				break
			}
		}

		if jsonFilePath == "" {
			return fmt.Errorf("failed to parse output for JSON file path when generating vanity key for program %s", program)
		}

		// Get absolute path
		absPath, err := filepath.Abs(jsonFilePath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for JSON file: %w", err)
		}

		// Extract file name
		fileName := filepath.Base(absPath)
		keys[program] = strings.TrimSuffix(fileName, ".json")

		destination := filepath.Join(cloneDir, deployDir, getTypeToProgramDeployName()[program]+"-keypair.json")
		if err := os.Rename(absPath, destination); err != nil {
			return fmt.Errorf("failed to move generated vanity key from %s to %s: %w", absPath, destination, err)
		}
		fmt.Println("File copied to:", destination)
	}
	return nil
}

func copyFile(srcFile string, destDir string) error {
	output, err := RunCommand("cp", []string{srcFile, destDir}, ".")
	if err != nil {
		return fmt.Errorf("failed to copy file: %s %w", output, err)
	}
	return nil
}

// Build the project with Anchor
func buildProject(e cldf.Environment, cloneDir, anchorDir, buildCmd string) error {
	solanaDir := filepath.Join(cloneDir, anchorDir, "..")
	e.Logger.Debugw("Building project", "solanaDir", solanaDir)
	args := []string{buildCmd}
	output, err := RunCommand("make", args, solanaDir)
	if err != nil {
		return fmt.Errorf("anchor build failed: %s %w", output, err)
	}
	return nil
}

func buildLocally(e cldf.Environment, config BuildSolanaConfig, params DomainParams) error {
	e.Logger.Debugw("Starting local build process", "destinationDir", config.DestinationDir)
	// Clone the repository
	if err := cloneRepo(e, config.GitCommitSha, params.CloneDir, params.RepoURL, config.LocalBuild.CleanGitDir); err != nil {
		return fmt.Errorf("error cloning repo: %w", err)
	}

	// Replace keys in Rust files using anchor keys sync
	if err := replaceKeys(e, params.CloneDir, params.AnchorDir, params.ReplaceKeysCmd); err != nil {
		return fmt.Errorf("error replacing keys: %w", err)
	}

	if config.LocalBuild.GenerateVanityKeys {
		if len(config.LocalBuild.UpgradeKeys) == 0 {
			config.LocalBuild.UpgradeKeys = make(map[cldf.ContractType]string)
		}
		if err := generateVanityKeys(e, params.CloneDir, params.DeployDir, config.LocalBuild.UpgradeKeys); err != nil {
			return fmt.Errorf("error generating vanity keys: %w", err)
		}
	}

	// Replace keys in Rust files for upgrade by replacing the declare_id!() macro explicitly
	// We need to do this so the keys will match the existing deployed program
	if err := replaceKeysForUpgrade(e, params.CloneDir, params.AnchorDir, params.ProgramFilesView, config.LocalBuild.UpgradeKeys); err != nil {
		return fmt.Errorf("error replacing keys for upgrade: %w", err)
	}

	// run sync to replace keys in programs that need to be synchonized
	for _, sync := range params.Syncers {
		if err := sync(); err != nil {
			return fmt.Errorf("error syncing program files: %w", err)
		}
	}

	// Build the project with Anchor
	if err := buildProject(e, params.CloneDir, params.AnchorDir, params.BuildCmd); err != nil {
		return fmt.Errorf("error building project: %w", err)
	}

	if config.LocalBuild.CleanDestinationDir {
		e.Logger.Debugw("Cleaning destination dir", "destinationDir", config.DestinationDir)
		if err := os.RemoveAll(config.DestinationDir); err != nil {
			return fmt.Errorf("error cleaning build folder: %w", err)
		}
		e.Logger.Debugw("Creating destination dir", "destinationDir", config.DestinationDir)
		err := os.MkdirAll(config.DestinationDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create build directory: %w", err)
		}
	} else if config.LocalBuild.CreateDestinationDir {
		e.Logger.Debugw("Creating destination dir", "destinationDir", config.DestinationDir)
		err := os.MkdirAll(config.DestinationDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create build directory: %w", err)
		}
	}

	deployFilePath := filepath.Join(params.CloneDir, params.DeployDir)
	e.Logger.Debugw("Reading deploy directory", "deployFilePath", deployFilePath)
	files, err := os.ReadDir(deployFilePath)
	if err != nil {
		return fmt.Errorf("failed to read deploy directory: %w", err)
	}

	for _, file := range files {
		filePath := filepath.Join(deployFilePath, file.Name())
		e.Logger.Debugw("Copying file", "filePath", filePath, "destinationDir", config.DestinationDir)
		err := copyFile(filePath, config.DestinationDir)
		if err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
	}
	return nil
}

func BuildSolana(e cldf.Environment, config BuildSolanaConfig, params DomainParams) error {
	if !config.LocalBuild.BuildLocally {
		e.Logger.Debug("Downloading Solana CCIP program artifacts...")
		err := memory.DownloadSolanaCCIPProgramArtifacts(e.GetContext(), config.DestinationDir, e.Logger, config.GitCommitSha)
		if err != nil {
			return fmt.Errorf("error downloading solana ccip program artifacts: %w", err)
		}
	} else {
		e.Logger.Debug("Building Solana program artifacts locally...")
		err := buildLocally(e, config, params)
		if err != nil {
			return fmt.Errorf("error building solana program artifacts: %w", err)
		}
	}

	return nil
}

func getTypeToProgramDeployName() map[cldf.ContractType]string {
	return map[cldf.ContractType]string{
		shared.Router:                  deployment.RouterProgramName,
		shared.OffRamp:                 deployment.OffRampProgramName,
		shared.FeeQuoter:               deployment.FeeQuoterProgramName,
		shared.BurnMintTokenPool:       deployment.BurnMintTokenPoolProgramName,
		shared.LockReleaseTokenPool:    deployment.LockReleaseTokenPoolProgramName,
		shared.RMNRemote:               deployment.RMNRemoteProgramName,
		types.AccessControllerProgram:  deployment.AccessControllerProgramName,
		types.ManyChainMultisigProgram: deployment.McmProgramName,
		types.RBACTimelockProgram:      deployment.TimelockProgramName,
		shared.Receiver:                deployment.ReceiverProgramName,
		ks_shared.Forwarder:            deployment.KeystoneForwarderProgramName,
		ks_shared.DataFeedsCache:       deployment.DataFeedsCacheProgramName,
	}
}
