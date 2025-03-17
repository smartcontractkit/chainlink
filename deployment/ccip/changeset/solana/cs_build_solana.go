package solana

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/go-github/v58/github"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	csState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"golang.org/x/oauth2"
)

// Configuration
const (
	repoURL        = "https://github.com/smartcontractkit/chainlink-ccip.git"
	repoOrgAndName = "smartcontractkit/chainlink-ccip"
	cloneDir       = "./temp-repo"
	anchorDir      = "chains/solana/contracts" // Path to the Anchor project within the repo
	deployDir      = "chains/solana/contracts/target/deploy"
)

// Map program names to their Rust file paths (relative to the Anchor project root)
// Needed for upgrades in place
var programToFileMap = map[deployment.ContractType]string{
	cs.Router:                      "programs/ccip-router/src/lib.rs",
	cs.FeeQuoter:                   "programs/fee-quoter/src/lib.rs",
	cs.OffRamp:                     "programs/ccip-offramp/src/lib.rs",
	cs.BurnMintTokenPool:           "programs/example-burnmint-token-pool/src/lib.rs",
	cs.LockReleaseTokenPool:        "programs/example-lockrelease-token-pool/src/lib.rs",
	cs.RMNRemote:                   "programs/rmn-remote/src/lib.rs",
	types.AccessControllerProgram:  "programs/access-controller/src/lib.rs",
	types.ManyChainMultisigProgram: "programs/mcm/src/lib.rs",
	types.RBACTimelockProgram:      "programs/timelock/src/lib.rs",
}

// Run a command in a specific directory
func runCommand(command string, args []string, workDir string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = workDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	fmt.Println("Running command", cmd.String())
	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}
	return stdout.String(), nil
}

// Clone and checkout the specific revision of the repo
func cloneRepo(e deployment.Environment, revision string, forceClean bool) error {
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
		_, err := runCommand("git", []string{"reset", "--hard"}, cloneDir)
		if err != nil {
			return fmt.Errorf("failed to discard local changes: %w", err)
		}

		// Fetch the latest changes from the remote
		_, err = runCommand("git", []string{"fetch", "origin"}, cloneDir)
		if err != nil {
			return fmt.Errorf("failed to fetch origin: %w", err)
		}
	} else {
		// Repository does not exist, clone it
		e.Logger.Debugw("Cloning repository", "url", repoURL, "revision", revision)
		_, err := runCommand("git", []string{"clone", repoURL, cloneDir}, ".")
		if err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	}

	e.Logger.Debugw("Checking out revision", "revision", revision)
	_, err := runCommand("git", []string{"checkout", revision}, cloneDir)
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
		return fmt.Errorf("anchor key replacement failed: %s %w", output, err)
	}
	return nil
}

func replaceKeysForUpgrade(e deployment.Environment, keys map[deployment.ContractType]string) error {
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
	args := []string{"docker-build-contracts"}
	output, err := runCommand("make", args, solanaDir)
	if err != nil {
		return fmt.Errorf("anchor build failed: %s %w", output, err)
	}
	return nil
}

type BuildSolanaConfig struct {
	GitCommitSha string
	// when running using CLD, this should be same as the secret (solana_program_path) or envvar (SOLANA_PROGRAM_PATH)
	DestinationDir       string
	CleanDestinationDir  bool
	CreateDestinationDir bool
	// Forces re-clone of git directory. Useful for forcing regeneration of keys
	CleanGitDir   bool
	ReplaceKeys   bool
	UpgradeKeys   map[deployment.ContractType]string
	VerifiedBuild bool
}

// https://solana.com/developers/guides/advanced/verified-builds
type VerifyBuildConfig struct {
	GitCommitSha               string
	ChainSelector              uint64
	VerifyFeeQuoter            bool
	VerifyRouter               bool
	VerifyOffRamp              bool
	VerifyRMNRemote            bool
	VerifyBurnMintTokenPool    bool
	VerifyLockReleaseTokenPool bool
	VerifyAccessController     bool
	VerifyMCM                  bool
	VerifyTimelock             bool
	RemoteVerification         bool
	MCMSSolana                 *MCMSConfigSolana
}

func BuildSolana(e deployment.Environment, config BuildSolanaConfig) error {
	// to use verified builds and actually have them work you need to:
	// 1. have the gh cli installed (brew install gh). This is already installed on GH runners so this will work in CI.
	// 2. have the private keypair files sourced from somewhere and already located in the DestinationDir. This is orthoganal to the verified build process
	if config.VerifiedBuild {
		if config.CreateDestinationDir || config.CleanDestinationDir {
			return errors.New("create or clean destination dir cannot be true when using verified builds. This could delete the private keypair files")
		}
		downloadRelease(e, config.DestinationDir)
	} else {
		// Clone the repository
		if err := cloneRepo(e, config.GitCommitSha, config.CleanGitDir); err != nil {
			return fmt.Errorf("error cloning repo: %w", err)
		}

		if config.ReplaceKeys {
			// Replace keys in Rust files using anchor keys sync
			if err := replaceKeys(e); err != nil {
				return fmt.Errorf("error replacing keys: %w", err)
			}

			// Replace keys in Rust files for upgrade by replacing the declare_id!() macro explicitly
			// We need to do this so the keys will match the existing deployed program
			if err := replaceKeysForUpgrade(e, config.UpgradeKeys); err != nil {
				return fmt.Errorf("error replacing keys for upgrade: %w", err)
			}
		}

		// Build the project with Anchor
		if err := buildProject(e); err != nil {
			return fmt.Errorf("error building project: %w", err)
		}

		if config.CleanDestinationDir {
			e.Logger.Debugw("Cleaning destination dir", "destinationDir", config.DestinationDir)
			if err := os.RemoveAll(config.DestinationDir); err != nil {
				return fmt.Errorf("error cleaning build folder: %w", err)
			}
			e.Logger.Debugw("Creating destination dir", "destinationDir", config.DestinationDir)
			err := os.MkdirAll(config.DestinationDir, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create build directory: %w", err)
			}
		} else if config.CreateDestinationDir {
			e.Logger.Debugw("Creating destination dir", "destinationDir", config.DestinationDir)
			err := os.MkdirAll(config.DestinationDir, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create build directory: %w", err)
			}
		}

		deployFilePath := filepath.Join(cloneDir, deployDir)
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
	}

	return nil
}

func runSolanaVerify(networkURL, programID, libraryName, commitHash, mountPath string, remote bool) error {
	cmdArgs := []string{
		"verify-from-repo",
		"-u", networkURL,
		"--program-id", programID,
		"--library-name", libraryName,
		strings.TrimSuffix(repoURL, ".git"),
		"--commit-hash", commitHash,
		"--mount-path", mountPath,
	}

	// Add --remote flag if remote verification is enabled
	if remote {
		cmdArgs = append(cmdArgs, "--remote")
	}

	output, err := runCommand("solana-verify", cmdArgs, ".")
	fmt.Println(output)
	if err != nil {
		return fmt.Errorf("solana program verification failed: %s %w", output, err)
	}
	return nil
}

func VerifyBuild(e deployment.Environment, cfg VerifyBuildConfig) (deployment.ChangesetOutput, error) {
	chain := e.SolChains[cfg.ChainSelector]
	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]

	addresses, err := e.ExistingAddresses.AddressesForChain(cfg.ChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := csState.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	verifications := []struct {
		name       string
		programID  string
		programLib string
		enabled    bool
	}{
		{"Fee Quoter", chainState.FeeQuoter.String(), deployment.FeeQuoterProgramName, cfg.VerifyFeeQuoter},
		{"Router", chainState.Router.String(), deployment.RouterProgramName, cfg.VerifyRouter},
		{"OffRamp", chainState.OffRamp.String(), deployment.OffRampProgramName, cfg.VerifyOffRamp},
		{"RMN Remote", chainState.RMNRemote.String(), deployment.RMNRemoteProgramName, cfg.VerifyRMNRemote},
		{"Burn Mint Token Pool", chainState.BurnMintTokenPool.String(), deployment.BurnMintTokenPoolProgramName, cfg.VerifyBurnMintTokenPool},
		{"Lock Release Token Pool", chainState.LockReleaseTokenPool.String(), deployment.LockReleaseTokenPoolProgramName, cfg.VerifyLockReleaseTokenPool},
		{"Access Controller", mcmState.AccessControllerProgram.String(), deployment.AccessControllerProgramName, cfg.VerifyAccessController},
		{"MCM", mcmState.McmProgram.String(), deployment.McmProgramName, cfg.VerifyMCM},
		{"Timelock", mcmState.TimelockProgram.String(), deployment.TimelockProgramName, cfg.VerifyTimelock},
	}

	for _, v := range verifications {
		if !v.enabled {
			continue
		}

		e.Logger.Debugw("Verifying program", "name", v.name, "programID", v.programID, "programLib", v.programLib)
		err := runSolanaVerify(
			chain.URL,
			v.programID,
			v.programLib,
			cfg.GitCommitSha,
			anchorDir,
			cfg.RemoteVerification,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("error verifying %s: %w", v.name, err)
		}
	}

	return deployment.ChangesetOutput{}, nil
}

func downloadRelease(e deployment.Environment, destination string) {
	owner := "smartcontractkit"
	repo := "chainlink-ccip"
	tag := "solana-artifacts-localtest-b0785c190f700710bcf5bbb0ce6b9b03e86e67ec"
	// destination := "../../../domains/ccip/testnet/inputs/"

	// ✅ Try getting GITHUB_TOKEN (only if needed)
	githubToken := os.Getenv("GITHUB_TOKEN")

	// ✅ Create a context
	ctx := e.GetContext()

	var client *github.Client

	if githubToken != "" {
		// ✅ Use authentication if a token is provided
		fmt.Println("🔑 Using authenticated GitHub client")
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		// ✅ Use an **unauthenticated** GitHub client (public repo access)
		fmt.Println("🌍 No GITHUB_TOKEN found, using public GitHub client")
		client = github.NewClient(&http.Client{})
	}

	release, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)

	if err != nil {
		fmt.Printf("failed to find release %s: %v", tag, err)
		// return fmt.Errorf("failed to find release %s: %w", tag, err)
	}

	fmt.Printf("Found release: %s (ID: %d)\n", *release.TagName, *release.ID)

	for _, asset := range release.Assets {
		fmt.Printf("Downloading asset: %s (ID: %d)\n", *asset.Name, *asset.ID)
		err := downloadFile(ctx, client, owner, repo, *asset.ID, destination, *asset.Name)
		if err != nil {
			fmt.Printf("failed to download asset %s: %v", *asset.Name, err)
		} else {
			fmt.Printf("Downloaded %s successfully\n", *asset.Name)
		}
	}
	return
}

func downloadFile(ctx context.Context, client *github.Client, owner, repo string, assetID int64, destDir, filename string) error {
	// 🔹 Step 1: Try to get the asset download URL
	fmt.Printf("🔍 Getting asset download URL for asset ID: %d\n", assetID)
	asset, _, err := client.Repositories.GetReleaseAsset(ctx, owner, repo, assetID)
	if err != nil {
		return fmt.Errorf("failed to get asset metadata: %w", err)
	}

	downloadURL := asset.GetBrowserDownloadURL()
	if downloadURL == "" {
		return fmt.Errorf("no browser download URL available for asset %d", assetID)
	}

	fmt.Println("🔗 Downloading from:", downloadURL)

	// 🔹 Step 2: Download using HTTP (handles both auth & public repos)
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}
	defer resp.Body.Close()

	// 🔹 Step 3: Create directory and save file
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	filePath := filepath.Join(destDir, filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	// 🔹 Step 4: Write to file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save asset: %w", err)
	}

	fmt.Printf("✅ Asset saved to %s\n", filePath)
	return nil
}
