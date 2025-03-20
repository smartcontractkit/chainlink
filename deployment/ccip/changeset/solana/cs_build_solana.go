package solana

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/v58/github"
	"github.com/smartcontractkit/chainlink/deployment"
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"golang.org/x/mod/modfile"
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

type LocalBuildConfig struct {
	BuildLocally         bool
	CleanDestinationDir  bool
	CreateDestinationDir bool
	// Forces re-clone of git directory. Useful for forcing regeneration of keys
	CleanGitDir bool
}

type BuildSolanaConfig struct {
	GitCommitSha string
	// when running using CLD, this should be same as the secret (solana_program_path) or envvar (SOLANA_PROGRAM_PATH)
	DestinationDir string
	LocalBuild     LocalBuildConfig
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

func buildLocally(e deployment.Environment, config BuildSolanaConfig) error {
	// Clone the repository
	if err := cloneRepo(e, config.GitCommitSha, config.LocalBuild.CleanGitDir); err != nil {
		return fmt.Errorf("error cloning repo: %w", err)
	}

	// Build the project with Anchor
	if err := buildProject(e); err != nil {
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

	return nil
}

func getSolanaCcipDependencyVersion(gomodPath string) (string, error) {
	const dependency = "github.com/smartcontractkit/chainlink-ccip/chains/solana"

	gomod, err := os.ReadFile(gomodPath)
	if err != nil {
		return "", err
	}

	modFile, err := modfile.ParseLax("go.mod", gomod, nil)
	if err != nil {
		return "", err
	}

	for _, dep := range modFile.Require {
		if dep.Mod.Path == dependency {
			return dep.Mod.Version, nil
		}
	}

	return "", fmt.Errorf("dependency %s not found", dependency)
}

func getVersion() (version string, err error) {
	_, currentFile, _, _ := runtime.Caller(0)
	// Get the root directory by walking up from current file until we find go.mod
	rootDir := filepath.Dir(currentFile)
	for {
		if _, err := os.Stat(filepath.Join(rootDir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(rootDir)
		if parent == rootDir {
			return "", fmt.Errorf("could not find project root directory containing go.mod")
		}
		rootDir = parent
	}

	go_mod_version, err := getSolanaCcipDependencyVersion(filepath.Join(rootDir, "go.mod"))
	if err != nil {
		return "", err
	}
	tokens := strings.Split(go_mod_version, "-")
	if len(tokens) == 3 {
		version := tokens[len(tokens)-1]
		return version, nil
	} else {
		return "", fmt.Errorf("invalid go.mod version: %s", go_mod_version)
	}
}

func DownloadRelease(cfg BuildSolanaConfig) error {
	owner := "smartcontractkit"
	repo := "chainlink-ccip"
	if cfg.GitCommitSha == "" {
		version, err := getVersion()
		if err != nil {
			return fmt.Errorf("error getting version: %w", err)
		}
		cfg.GitCommitSha = version
	}

	tag := fmt.Sprintf("solana-artifacts-localtest-%s", cfg.GitCommitSha)

	// ✅ Try getting GITHUB_TOKEN (only if needed)
	githubToken := os.Getenv("GITHUB_TOKEN")

	// ✅ Create a context
	ctx := context.Background()

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
		return fmt.Errorf("failed to find release %s: %w", tag, err)
	}

	fmt.Printf("Found release: %s (ID: %d)\n", *release.TagName, *release.ID)

	for _, asset := range release.Assets {
		fmt.Printf("Downloading asset: %s (ID: %d)\n", *asset.Name, *asset.ID)
		err := downloadFile(ctx, client, owner, repo, *asset.ID, cfg.DestinationDir, *asset.Name)
		if err != nil {
			fmt.Printf("failed to download asset %s: %v", *asset.Name, err)
			return err
		} else {
			fmt.Printf("Downloaded %s successfully\n", *asset.Name)
		}
	}
	return nil
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

func BuildSolana(e deployment.Environment, config BuildSolanaConfig) error {
	if config.LocalBuild.BuildLocally {
		DownloadRelease(config)
	} else {
		buildLocally(e, config)
	}

	return nil
}
