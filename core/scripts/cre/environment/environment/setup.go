package environment

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/Masterminds/semver/v3"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/tracking"
)

var SetupCmd *cobra.Command

func init() {
	var (
		config      SetupConfig
		noPrompt    bool
		purge       bool
		withBilling bool
	)

	SetupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Setup the CRE environment prerequisites",
		Long:  `Checks and sets up prerequisites for the CRE environment including Docker, AWS, Job Distributor, and CRE CLI`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunSetup(cmd.Context(), config, noPrompt, purge, withBilling, relativePathToRepoRoot)
		},
	}

	SetupCmd.Flags().StringVarP(&config.ConfigPath, "config", "c", DefaultSetupConfigPath, "Path to the TOML configuration file")
	SetupCmd.Flags().BoolVarP(&noPrompt, "no-prompt", "y", false, "Automatically accept defaults and do not prompt for user input")
	SetupCmd.Flags().BoolVarP(&purge, "purge", "p", false, "Purge all existing images and re-download/re-build them")
	SetupCmd.Flags().BoolVar(&withBilling, "with-billing", false, "Include billing service in the setup")

	EnvironmentCmd.AddCommand(SetupCmd)

	BuildCapabilitiesCmd := &cobra.Command{
		Use:   "build-caps",
		Short: "Build capabilities binaries",
		Long:  `Builds the capabilities binaries for the CRE environment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return BuildCapabilities(cmd.Context(), config, noPrompt)
		},
	}

	BuildCapabilitiesCmd.Flags().StringVarP(&config.ConfigPath, "config", "c", DefaultSetupConfigPath, "Path to the TOML configuration file")
	BuildCapabilitiesCmd.Flags().BoolVarP(&noPrompt, "no-prompt", "y", false, "Automatically accept defaults and do not prompt for user input")
	EnvironmentCmd.AddCommand(BuildCapabilitiesCmd)
}

type config struct {
	General        generalConfig         `toml:"general"`
	JobDistributor jobDistributorConfig  `toml:"job_distributor"`
	ChipIngress    *chipIngressConfig    `toml:"chip_ingress"`
	ChipConfig     *chipConfigConfig     `toml:"chip_config"`
	BillingService *billingServiceConfig `toml:"billing_platform_service"`
	Capabilities   capabilitiesConfig    `toml:"capabilities"`
	Observability  observabilityConfig   `toml:"observability"`
}

type generalConfig struct {
	AWSProfile      string `toml:"aws_profile"`
	MinGHCLIVersion string `toml:"min_gh_cli_version"`
}

type jobDistributorConfig struct {
	BuildConfig BuildConfig `toml:"build_config"`
	PullConfig  PullConfig  `toml:"pull_config"`
}

type chipIngressConfig struct {
	BuildConfig BuildConfig `toml:"build_config"`
	PullConfig  PullConfig  `toml:"pull_config"`
}

type chipConfigConfig struct {
	BuildConfig BuildConfig `toml:"build_config"`
	PullConfig  PullConfig  `toml:"pull_config"`
}

type billingServiceConfig struct {
	BuildConfig BuildConfig `toml:"build_config"`
	PullConfig  PullConfig  `toml:"pull_config"`
}

type capabilitiesConfig struct {
	TargetPath   string   `toml:"target_path"`
	MakeCommands []string `toml:"make_commands"`
}

type observabilityConfig struct {
	RepoURL    string `toml:"repository"`
	Branch     string `toml:"branch"`
	TargetPath string `toml:"target_path"`
}

var (
	ECR = os.Getenv("AWS_ECR") // TODO this can be moved to an env file
)

const DefaultSetupConfigPath = "configs/setup.toml"
const DefaultCapabilityBinariesPath = ".binaries"

type EnsureOption = string

const (
	PullOption  EnsureOption = "p"
	BuildOption EnsureOption = "b"
)

// SetupConfig represents the configuration for the setup command
type SetupConfig struct {
	ConfigPath string
}

type BuildConfig struct {
	RepoURL            string `toml:"repository"`
	LocalRepo          string `toml:"local_repo"`
	Branch             string `toml:"branch"`
	Commit             string `toml:"commit"`
	RequireGithubToken bool   `toml:"require_github_token"`
	Dockerfile         string `toml:"dockerfile"`
	DockerCtx          string `toml:"docker_ctx"`
	LocalImage         string `toml:"local_image"`
	PreRun             string `toml:"pre_run"` // Optional function to run before building
}

// setupRepo clones the repository if it's a remote URL or uses the local path if it's a directory.
// It returns the working directory path, a boolean indicating if it's a local repo, and an error if any.
// It will checkout the specified reference branch/tag and commit if provided.
func setupRepo(ctx context.Context, logger zerolog.Logger, repo, reference, commit, workingDir string) (string, bool, error) {
	if repo == "" {
		return "", false, errors.New("repository URL or path is empty")
	}

	// Expand ~ to home directory in workingDir if present
	if workingDir != "" && strings.HasPrefix(workingDir, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", false, fmt.Errorf("failed to get user home directory: %w", err)
		}
		workingDir = filepath.Join(homeDir, workingDir[2:])
	}

	// Check if repo is a local directory
	isLocalRepo := false
	if _, err2 := os.Stat(repo); err2 == nil {
		fileInfo, err3 := os.Stat(repo)
		if err3 == nil && fileInfo.IsDir() {
			isLocalRepo = true
			logger.Info().Msgf("Using local repository at %s", repo)
		}
	}

	if isLocalRepo {
		// Use the local repo path directly
		workingDir = repo
	} else {
		if reference == "" {
			return "", false, errors.New("branch or tag reference is required for remote repositories")
		}

		if workingDir == "" {
			// Create a temporary directory for cloning the remote repo
			tempDir, err2 := os.MkdirTemp("", filepath.Base(repo)+"-*")
			if err2 != nil {
				return "", false, fmt.Errorf("failed to create temporary directory: %w", err2)
			}
			workingDir = tempDir
		} else {
			// Clear or create the working directory
			if _, err := os.Stat(workingDir); err == nil {
				if err = os.RemoveAll(workingDir); err != nil {
					return "", false, fmt.Errorf("failed to clear existing working directory: %w", err)
				}
			} else {
				if err = os.MkdirAll(workingDir, 0o755); err != nil {
					return "", false, fmt.Errorf("failed to create working directory: %w", err)
				}
			}
		}

		// Clone the repository
		logger.Info().Msgf("Cloning repository from %s", repo)
		cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", "--branch", reference, "--single-branch", repo, workingDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err2 := cmd.Run(); err2 != nil {
			return "", false, fmt.Errorf("failed to clone repository: %w", err2)
		}
		if commit != "" {
			// Checkout the specific commit if provided
			logger.Info().Msgf("Checking out commit %s", commit)
			cmd := exec.CommandContext(ctx, "git", "fetch", "--depth", "1", "origin", commit)
			cmd.Dir = workingDir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err2 := cmd.Run(); err2 != nil {
				return "", false, fmt.Errorf("failed to checkout commit %s: %w", commit, err2)
			}
			cmd = exec.CommandContext(ctx, "git", "checkout", commit)
			cmd.Dir = workingDir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err2 := cmd.Run(); err2 != nil {
				return "", false, fmt.Errorf("failed to checkout commit %s: %w", commit, err2)
			}
		}
	}

	return workingDir, isLocalRepo, nil
}

func (c BuildConfig) Build(ctx context.Context) (localImage string, err error) {
	var (
		repo   = c.RepoURL
		tag    = c.Branch
		commit = c.Commit
	)
	logger := framework.L
	name := strings.ReplaceAll(strings.Split(c.LocalImage, ":")[0], "-", " ")
	name = cases.Title(language.English).String(name)
	logger.Info().Msgf("Building %s image...", name)

	if c.RequireGithubToken {
		if os.Getenv("GITHUB_TOKEN") == "" {
			return "", errors.New("GITHUB_TOKEN environment variable is required to build the billing service from source")
		}
	}

	workingDir, isLocalRepo, err := setupRepo(ctx, logger, repo, tag, commit, "")
	if err != nil {
		return "", fmt.Errorf("failed to setup repository: %w", err)
	}

	if !isLocalRepo {
		defer func() {
			_ = os.RemoveAll(workingDir)
		}()
	}

	// Save current directory and change to working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(workingDir); err != nil {
		return "", fmt.Errorf("failed to change to working directory: %w", err)
	}
	defer func() {
		_ = os.Chdir(currentDir)
	}()

	// If pre-run function is specified, run it
	if c.PreRun != "" {
		logger.Info().Msgf("Running pre-run step: %s", c.PreRun)
		if err := exec.CommandContext(ctx, "bash", "-c", c.PreRun).Run(); err != nil { //nolint:gosec //G204: Subprocess launched with a potential tainted input or cmd arguments
			return "", fmt.Errorf("pre-run step failed: %w", err)
		}
	}

	// Build Docker image
	args := []string{"build", "-t", c.LocalImage, "-f", c.Dockerfile, c.DockerCtx}
	if c.RequireGithubToken {
		args = append(args, "--build-arg", "GITHUB_TOKEN="+os.Getenv("GITHUB_TOKEN"))
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Info("Running command:", "cmd", cmd.String(), "dir", workingDir)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to build Docker image: %w", err)
	}

	logger.Info().Msgf("  ✓ %s image built successfully", name)
	return c.LocalImage, nil
}

type PullConfig struct {
	LocalImage string `toml:"local_image"`
	EcrImage   string `toml:"ecr_image"`
}

func (c PullConfig) Pull(ctx context.Context, awsProfile string) (localImage string, err error) {
	if ECR == "" {
		return "", errors.New("AWS_ECR environment variable is not set. See README for more details and references to find the correct ECR URL or visit https://smartcontract-it.atlassian.net/wiki/spaces/INFRA/pages/1045495923/Configure+the+AWS+CLI")
	}

	tmpl, tmplErr := template.New("ecr-image").Parse(c.EcrImage)
	if tmplErr != nil {
		return "", errors.Wrapf(tmplErr, "failed to parse ECR image template")
	}

	templateData := map[string]string{
		"ECR": ECR,
	}

	var configBuffer bytes.Buffer
	if err := tmpl.Execute(&configBuffer, templateData); err != nil {
		return "", errors.Wrapf(err, "failed to execute ECR image template")
	}
	ecrImage := configBuffer.String()

	return pullImage(ctx, awsProfile, c.LocalImage, ecrImage)
}

type ImageConfig struct {
	BuildConfig BuildConfig
	PullConfig  PullConfig
}

func (c ImageConfig) Ensure(ctx context.Context, dockerClient *client.Client, awsProfile string, noPrompt bool, defaultOption EnsureOption, purge bool) (localImage string, err error) {
	// If purge flag is set, remove existing images first
	if purge {
		logger := framework.L
		name := strings.ReplaceAll(strings.Split(c.BuildConfig.LocalImage, ":")[0], "-", " ")
		name = cases.Title(language.English).String(name)
		logger.Info().Msgf("🗑️  Purging existing %s images...", name)

		// Remove local image if it exists
		_, err = dockerClient.ImageRemove(ctx, c.BuildConfig.LocalImage, image.RemoveOptions{Force: true})
		if err != nil {
			logger.Warn().Msgf("Failed to remove local image %s: %v", c.BuildConfig.LocalImage, err)
		}

		// Remove ECR image if it exists
		_, err = dockerClient.ImageRemove(ctx, c.PullConfig.EcrImage, image.RemoveOptions{Force: true})
		if err != nil {
			logger.Warn().Msgf("Failed to remove ECR image %s: %v", c.PullConfig.EcrImage, err)
		}

		logger.Info().Msgf("  ✓ %s images purged", name)
	}

	exist, err := localImageExists(ctx, dockerClient, c.BuildConfig.LocalImage, c.PullConfig.EcrImage)
	if err != nil {
		return "", fmt.Errorf("failed to check if image exists: %w", err)
	}
	if !exist {
		// If not exist, ask to pull or build
		logger := framework.L
		name := strings.ReplaceAll(strings.Split(c.BuildConfig.LocalImage, ":")[0], "-", " ")
		name = cases.Title(language.English).String(name)
		logger.Info().Msgf("🔍 %s image not found.", name)
		logger.Info().Msgf("Would you like to Pull (requires AWS SSO) or build the %s image? (P/b) [B]", name)

		var input = PullOption // Default to Pull
		if !noPrompt {
			_, err := fmt.Scanln(&input)
			if err != nil {
				// If error is due to empty input (just pressing Enter), use default
				if err.Error() != "unexpected newline" {
					return "", errors.Wrap(err, "failed to read input")
				}
			}
		}
		// check that input is valid
		input = strings.TrimSpace(strings.ToLower(input))
		if input != PullOption && input != BuildOption {
			logger.Warn().Msg("Invalid input. Please enter 'p' or 'b'.")
			return "", fmt.Errorf("invalid input: %s", input)
		}

		if strings.ToLower(input) == BuildOption {
			return c.BuildConfig.Build(ctx)
		}

		return c.PullConfig.Pull(ctx, awsProfile)
	}
	return c.BuildConfig.LocalImage, nil
}

// RunSetup performs the setup for the CRE environment
func RunSetup(ctx context.Context, config SetupConfig, noPrompt, purge, withBilling bool, relativePathToRepoRoot string) (setupErr error) {
	logger := framework.L
	var localDXTracker tracking.Tracker
	localDXTracker = &tracking.NoOpTracker{}

	defer func() {
		var trackingErr error
		if setupErr != nil {
			trackingErr = localDXTracker.Track(MetricSetupResult, map[string]any{"result": "failure", "no_prompt": noPrompt, "error": oneLineErrorMessage(setupErr)})
		} else {
			trackingErr = localDXTracker.Track(MetricSetupResult, map[string]any{"result": "success", "no_prompt": noPrompt})
		}
		if trackingErr != nil {
			fmt.Fprintf(os.Stderr, "failed to track setup: %s\n", trackingErr)
		}
	}()

	logger.Info().Msg("🔍 Checking prerequisites for CRE environment...")

	// Check if Docker is installed
	if !isCommandAvailable("docker") {
		setupErr = errors.New("docker is not installed. Please install Docker and try again")
		return
	}
	logger.Info().Msg("✓ Docker is installed")

	// Check if Docker is running
	dockerClient, dockerClientErr := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if dockerClientErr != nil {
		setupErr = errors.Wrap(dockerClientErr, "failed to create Docker client")
		return
	}

	_, pingErr := dockerClient.Ping(ctx)
	if pingErr != nil {
		setupErr = errors.Wrap(pingErr, "docker is not running. Please start Docker and try again")
		return
	}
	logger.Info().Msg("✓ Docker is running")

	// Check Docker configuration
	if dockerConfigErr := checkDockerConfiguration(); dockerConfigErr != nil {
		setupErr = errors.Wrap(dockerConfigErr, "failed to check Docker configuration")
		return
	}

	// Check if AWS CLI is installed
	if !noPrompt {
		if !isCommandAvailable("aws") {
			setupErr = errors.New("AWS CLI is not installed. Please install AWS CLI and try again")
			return
		}
		logger.Info().Msg("✓ AWS CLI is installed")
	}

	cfg, cfgErr := readConfig(config.ConfigPath)
	if cfgErr != nil {
		setupErr = errors.Wrap(cfgErr, "failed to read config")
		return
	}

	ghCli, ghCliErr := checkGHCli(ctx, cfg.General.MinGHCLIVersion, noPrompt)
	if ghCliErr != nil {
		setupErr = errors.Wrap(ghCliErr, "failed to ensure GitHub CLI")
		return
	}

	// once we have GH CLI setup we can try to create the DX tracker
	if ghCli {
		var trackerErr error
		localDXTracker, trackerErr = tracking.NewDxTracker(GetDXGitHubVariableName, GetDXProductName)
		if trackerErr != nil {
			fmt.Fprintf(os.Stderr, "failed to create DX tracker: %s\n", trackerErr)
		}
	}

	bun, bunErr := checkBun(ctx, noPrompt)
	if bunErr != nil {
		setupErr = errors.Wrap(bunErr, "failed to ensure Bun")
		return
	}

	if bun {
		err := ensurePackageJSON(".")
		if err != nil {
			setupErr = errors.Wrap(err, "failed to ensure package.json")
			return
		}

		if err := installBunPackages(ctx); err != nil {
			setupErr = errors.Wrap(err, "failed to install Bun packages")
			return
		}
	}

	jdConfig := ImageConfig{
		BuildConfig: cfg.JobDistributor.BuildConfig,
		PullConfig:  cfg.JobDistributor.PullConfig,
	}

	jdLocalImage, jdErr := jdConfig.Ensure(ctx, dockerClient, cfg.General.AWSProfile, noPrompt, PullOption, purge)
	if jdErr != nil {
		setupErr = errors.Wrap(jdErr, "failed to ensure Job Distributor image")
		return
	}

	var chipIngressLocalImage string
	if cfg.ChipIngress != nil {
		chipConfig := ImageConfig{
			BuildConfig: cfg.ChipIngress.BuildConfig,
			PullConfig:  cfg.ChipIngress.PullConfig,
		}

		var err error
		chipIngressLocalImage, err = chipConfig.Ensure(ctx, dockerClient, cfg.General.AWSProfile, noPrompt, PullOption, purge)
		if err != nil {
			setupErr = errors.Wrap(err, "failed to ensure Atlas Chip Ingress image")
			return
		}
	} else {
		logger.Warn().Str("config file", config.ConfigPath).Msgf("Skipping Atlas Chip Ingress setup, because configuration is not provided in the config file")
	}

	var chipConfigLocalImage string
	if cfg.ChipConfig != nil {
		chipConfig := ImageConfig{
			BuildConfig: cfg.ChipConfig.BuildConfig,
			PullConfig:  cfg.ChipConfig.PullConfig,
		}

		var err error
		chipConfigLocalImage, err = chipConfig.Ensure(ctx, dockerClient, cfg.General.AWSProfile, noPrompt, PullOption, purge)
		if err != nil {
			setupErr = errors.Wrap(err, "failed to ensure Atlas Chip Config image")
			return
		}
	} else {
		logger.Warn().Str("config file", config.ConfigPath).Msgf("Skipping Atlas Chip Config setup, because configuration is not provided in the config file")
	}

	var billingLocalImage string
	if withBilling {
		if cfg.BillingService == nil {
			setupErr = errors.New("billing service configuration is required when using --with-billing flag")
			return
		}

		billingConfig := ImageConfig{
			BuildConfig: cfg.BillingService.BuildConfig,
			PullConfig:  cfg.BillingService.PullConfig,
		}

		var billingErr error
		// Try to build Billing service since almost noone has access the ECR that stores the image
		billingLocalImage, billingErr = billingConfig.Ensure(ctx, dockerClient, cfg.General.AWSProfile, noPrompt, BuildOption, purge)
		if billingErr != nil {
			setupErr = errors.Wrap(billingErr, "failed to ensure Billing Platform Service image")
			return
		}
	} else {
		logger.Warn().Msgf("Skipping Billing Platform Service setup, because the --with-billing flag was not provided")
	}

	if err := runGHSetupGit(ctx); err != nil {
		return errors.Wrap(err, "failed to run 'gh auth setup-git'")
	}

	observabilityRepoPath, _, err := setupRepo(ctx, logger, cfg.Observability.RepoURL, cfg.Observability.Branch,
		"", cfg.Observability.TargetPath)
	if err != nil {
		setupErr = errors.Wrap(err, "failed to clone observability repo")
		return
	}

	installedCapabilities, capErr := makeCapabilities(ctx, cfg.Capabilities, relativePathToRepoRoot)
	if capErr != nil {
		return errors.Wrap(capErr, "failed to install capabilities")
	}

	// Print summary
	fmt.Println()
	logger.Info().Msg("✅ Setup Summary:")
	logger.Info().Msg("   ✓ Docker is installed and configured correctly")
	logger.Info().Msgf("   ✓ Job Distributor image %s is available", jdLocalImage)
	if chipIngressLocalImage != "" {
		logger.Info().Msgf("   ✓ Atlas Chip Ingress image %s is available", chipIngressLocalImage)
	}
	if chipConfigLocalImage != "" {
		logger.Info().Msgf("   ✓ Atlas Chip Config image %s is available", chipConfigLocalImage)
	}
	logger.Info().Msgf("   ✓ Observability repo cloned to %s", observabilityRepoPath)
	if billingLocalImage != "" {
		logger.Info().Msgf("   ✓ Billing Platform Service image %s is available", billingLocalImage)
	}
	if ghCli {
		logger.Info().Msg("   ✓ GitHub CLI is installed")
	} else {
		logger.Warn().Msg("   ✗ GitHub CLI is not installed")
	}
	if bun {
		logger.Info().Msg("   ✓ Bun is installed")
	} else {
		logger.Warn().Msg("   ✗ Bun is not installed")
	}
	if len(cfg.Capabilities.MakeCommands) > 0 {
		logger.Info().Msg("   ✓ Capabilities binaries installed")
		logger.Info().Msgf("     - capabilities: %s", strings.Join(installedCapabilities, ", "))
	}

	fmt.Println()
	logger.Info().Msg("🚀 Next Steps:")
	logger.Info().Msg("1. Navigate to the CRE environment directory: cd core/scripts/cre/environment")
	logger.Info().Msg("2. Start the environment: go run . env start")
	logger.Info().Msg("   Optional: Add --with-example to start with an example workflow")
	logger.Info().Msg("   Optional: Add --with-plugins-docker-image to use a pre-built image with capabilities")
	logger.Info().Msg("   Optional: Add --with-beholder to start the Beholder")
	logger.Info().Msg("\nFor more information, see the documentation in core/scripts/cre/environment/README.md")

	return nil
}

func BuildCapabilities(ctx context.Context, config SetupConfig, noPrompt bool) error {
	cfg, cfgErr := readConfig(config.ConfigPath)
	if cfgErr != nil {
		return errors.Wrap(cfgErr, "failed to read config")
	}

	_, ghCliErr := checkGHCli(ctx, cfg.General.MinGHCLIVersion, noPrompt)
	if ghCliErr != nil {
		return errors.Wrap(ghCliErr, "failed to ensure GitHub CLI")
	}

	if err := runGHSetupGit(ctx); err != nil {
		return errors.Wrap(err, "failed to run 'gh auth setup-git'")
	}

	installedCapabilities, capErr := makeCapabilities(ctx, cfg.Capabilities, relativePathToRepoRoot)
	if capErr != nil {
		return errors.Wrap(capErr, "failed to install capabilities")
	}

	fmt.Println()
	logger := framework.L
	logger.Info().Msg("✅ Build Capabilities Summary:")
	for _, capability := range installedCapabilities {
		logger.Info().Msgf("   ✓ %s", capability)
	}
	logger.Info().Msgf("   ✓ %d capabilities installed", len(installedCapabilities))

	return nil
}

func runGHSetupGit(ctx context.Context) error {
	logger := framework.L
	logger.Info().Msg("🔍 Checking if GitHub CLI authentication is set up for Git...")
	cmd := exec.CommandContext(ctx, "bash", "-c", `printf "protocol=https\nhost=github.com\n\n" | git credential fill | sed -n 's/^password=//p' | head -n1`)
	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to run git credential fill")
	}

	if out.String() == "" {
		logger.Info().Msg("  GitHub CLI authentication is not set up for Git. Running 'gh auth setup-git'...")
		setupCmd := exec.CommandContext(ctx, "gh", "auth", "setup-git")
		setupCmd.Stdout = os.Stdout
		setupCmd.Stderr = os.Stderr
		if err := setupCmd.Run(); err != nil {
			return errors.Wrap(err, "failed to run 'gh auth setup-git'")
		}
		logger.Info().Msg("  ✓ GitHub CLI authentication is now set up for Git.")
	} else {
		logger.Info().Msg("  ✓ GitHub CLI authentication is already set up for Git.")
	}

	return nil
}

func makeCapabilities(ctx context.Context, capabilitiesConfig capabilitiesConfig, repoRootRelativePath string) ([]string, error) {
	if len(capabilitiesConfig.MakeCommands) == 0 {
		framework.L.Info().Msg("No make commands specified for capabilities. Skipping capabilities build.")
		return nil, nil
	}

	logger := framework.L
	logger.Info().Msg("🔍 Installing capabilities binaries...")

	tempDir, tempErr := os.MkdirTemp(".", ".tmp-capability-binaries")
	if tempErr != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", tempErr)
	}

	tempDirAbsPath, tAbsErr := filepath.Abs(tempDir)
	if tAbsErr != nil {
		return nil, fmt.Errorf("failed to get absolute path of temporary directory: %w", tAbsErr)
	}

	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	for _, makeCommand := range capabilitiesConfig.MakeCommands {
		cmd := exec.CommandContext(ctx, "make", makeCommand)
		cmd.Dir = repoRootRelativePath
		// Set GOBIN to the absolute path of the target path, so that binaries are placed there
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GOBIN="+tempDirAbsPath)
		// cross-compile for linux/amd64 with CGO disabled, because our Chainlink Docker images use linux/amd64
		cmd.Env = append(cmd.Env, "CL_PLUGIN_ENVVARS=GOOS=linux GOARCH=amd64 CGO_ENABLED=0")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to run make command '%s': %w", makeCommand, err)
		}
	}

	if capabilitiesConfig.TargetPath == "" {
		capabilitiesConfig.TargetPath = DefaultCapabilityBinariesPath
	}

	absPath, absErr := filepath.Abs(capabilitiesConfig.TargetPath)
	if absErr != nil {
		return nil, fmt.Errorf("failed to get absolute path of target path: %w", absErr)
	}

	if err := os.MkdirAll(absPath, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create target path: %w", err)
	}

	cmd := exec.CommandContext(ctx, "cp", "-R", tempDirAbsPath+string(os.PathSeparator)+".", absPath) //nolint:gosec //G204: Subprocess launched with a potential tainted input or cmd arguments
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to copy binaries to target path: %w", err)
	}

	files, err := os.ReadDir(tempDirAbsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read temporary directory: %w", err)
	}

	fmt.Println("Dir: ", tempDirAbsPath)

	installedCapabilities := []string{}
	for _, f := range files {
		if f.Type().IsRegular() {
			installedCapabilities = append(installedCapabilities, f.Name())
		}
	}

	logger.Info().Msgf("  ✓ %d capabilities binaries installed in %s", len(installedCapabilities), absPath)

	return installedCapabilities, nil
}

func readConfig(configPath string) (*config, error) {
	cfg := &config{}

	cfgBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	if err := toml.Unmarshal(cfgBytes, cfg); err != nil {
		return nil, errors.Wrap(err, "failed to decode config")
	}

	return cfg, nil
}

// isCommandAvailable checks if a command is available in the PATH
func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// checkDockerConfiguration checks if Docker is configured correctly
func checkDockerConfiguration() error {
	logger := framework.L
	logger.Info().Msg("🔍 Checking Docker settings...")

	dockerSettingsOK := true
	osType := runtime.GOOS

	// Check for settings based on OS
	switch osType {
	case "darwin":
		logger.Info().Msg("  Detected macOS system")
		configPaths := []string{
			filepath.Join(os.Getenv("HOME"), "Library/Group Containers/group.com.docker/settings-store.json"),
			filepath.Join(os.Getenv("HOME"), "Library/Group Containers/group.com.docker/settings.json"),
		}

		configFile := ""
		for _, path := range configPaths {
			if _, err := os.Stat(path); err == nil {
				configFile = path
				break
			}
		}

		if configFile == "" {
			logger.Warn().Msgf(" ! Could not find Docker settings files in %s. Your Docker installation may be misconfigured.", strings.Join(configPaths, ", "))
		}

		logger.Info().Msgf("  Found Docker settings file at %s", configFile)

		// Check settings
		settings, err := os.ReadFile(configFile)
		if err != nil {
			if strings.Contains(err.Error(), "operation not permitted") {
				logger.Warn().Msgf("  ! Could not check Docker settings due to restrictive TCC policies (can't read file). You need to manually verify the settings in the Docker Desktop UI.")
				return nil
			}
			return fmt.Errorf("failed to read Docker settings: %w", err)
		}

		// Check required settings using gjson
		settingsChecks := map[string]string{
			"UseVirtualizationFramework":         "true",
			"UseVirtualizationFrameworkVirtioFS": "true",
			"EnableDefaultDockerSocket":          "true",
		}

		for setting, expected := range settingsChecks {
			value := gjson.GetBytes(settings, setting).String()
			switch {
			case value == expected:
				logger.Info().Msgf("  ✓ %s is correctly set to %s", setting, expected)
			case strings.TrimSpace(value) == "":
				// some users may not have this setting at all; warn instead of error
				logger.Warn().Msgf("  ! Could not find setting for %s (should be %s). Manually check Docker settings in the UI", setting, expected)
			default:
				logger.Error().Msgf("  ✗ %s is set to %s (should be %s)", setting, value, expected)
				dockerSettingsOK = false
			}
		}

		// Check CPU requirements (minimum 4 cores)
		cpuValue := gjson.GetBytes(settings, "Cpus").Int()
		switch {
		case cpuValue >= 4:
			logger.Info().Msgf("  ✓ CPU allocation is sufficient (%d cores)", cpuValue)
		case cpuValue == 0:
			logger.Warn().Msg("  ! Could not find CPU setting. Manually check Docker settings in the UI (should be at least 4 cores)")
		default:
			logger.Error().Msgf("  ✗ CPU allocation is insufficient (%d cores, should be at least 4)", cpuValue)
		}

		// Check memory requirements (minimum 10 GB = 10240 MiB)
		memoryValue := gjson.GetBytes(settings, "MemoryMiB").Int()
		switch {
		case memoryValue >= 10240:
			logger.Info().Msgf("  ✓ Memory allocation is sufficient (%d MiB / %.1f GB)", memoryValue, float64(memoryValue)/1024)
		case memoryValue == 0:
			logger.Warn().Msg("  ! Could not find memory setting. Manually check Docker settings in the UI (should be at least 10 GB)")
		default:
			logger.Error().Msgf("  ✗ Memory allocation is insufficient (%d MiB / %.1f GB, should be at least 10 GB)", memoryValue, float64(memoryValue)/1024)
		}

	case "linux":
		logger.Info().Msg("  Detected Linux system")
		logger.Info().Msg("  Docker daemon configuration typically doesn't need macOS-specific checks")

	default:
		logger.Warn().Msgf("  Unknown operating system: %s", osType)
		logger.Warn().Msg("  Cannot check Docker settings automatically")
		logger.Warn().Msg("  Please ensure Docker is properly configured for your system")
	}

	if !dockerSettingsOK {
		return errors.New("docker is not configured correctly. Please fix the issues and try again")
	}

	return nil
}

// localImageExists checks if the local image or ECR image exists
// if ECR image exists, it tags it as the local image
func localImageExists(ctx context.Context, dockerClient *client.Client, localImage, ecrImage string) (bool, error) {
	logger := framework.L
	name := strings.ReplaceAll(strings.Split(localImage, ":")[0], "-", " ")
	name = cases.Title(language.English).String(name)
	// Check if local image exists
	_, err := dockerClient.ImageInspect(ctx, localImage)
	if err == nil {
		logger.Info().Msgf("✓ %s image (%s) is available from local build", name, localImage)
		return true, nil
	}

	// Check if ECR image exists
	_, err = dockerClient.ImageInspect(ctx, ecrImage)
	if err == nil {
		logger.Info().Msgf("✓ %s image (%s) is available", name, ecrImage)
		// Tag ECR image as local image
		if err := dockerClient.ImageTag(ctx, ecrImage, localImage); err != nil {
			return false, fmt.Errorf("failed to tag %s image: %w", name, err)
		}
		logger.Info().Msgf("  ✓ %s image tagged as %s", name, localImage)
		return true, nil
	}
	return false, nil
}

// pullImage pulls the Job Distributor image from ECR
func pullImage(ctx context.Context, awsProfile string, localImage, ecrImage string) (string, error) {
	logger := framework.L
	name := strings.ReplaceAll(strings.Split(localImage, ":")[0], "-", " ")
	name = cases.Title(language.English).String(name)

	// Try pulling the image we need and login only if it doesn't succeed
	logger.Info().Msgf("Trying to pull Docker image %s...", ecrImage)
	pullCmd := exec.CommandContext(ctx, "docker", "pull", ecrImage)
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr
	if err := pullCmd.Run(); err != nil {
		// Check if AWS profile exists
		configureCmd := exec.CommandContext(ctx, "aws", "configure", "list-profiles")
		output, configureCmdErr := configureCmd.Output()
		if configureCmdErr != nil {
			return "", errors.Wrap(configureCmdErr, "failed to list AWS profiles")
		}

		if !strings.Contains(string(output), awsProfile) {
			return "", fmt.Errorf("AWS profile '%s' not found. Please ensure you have the correct AWS profile configured. Please see https://smartcontract-it.atlassian.net/wiki/spaces/INFRA/pages/1045495923/Configure+the+AWS+CLI", awsProfile)
		}

		// Get ECR login password
		// Check if we already have a valid AWS SSO session
		logger.Info().Msgf("Checking for valid AWS SSO session for profile %s...", awsProfile)
		checkCmd := exec.CommandContext(ctx, "aws", "sts", "get-caller-identity", "--profile", awsProfile)
		if err := checkCmd.Run(); err == nil {
			logger.Info().Msgf("  ✓ Valid AWS SSO session exists for profile %s", awsProfile)
		} else {
			// No valid session, need to log in
			logger.Info().Msgf("AWS SSO Login required for profile %s...", awsProfile)
			loginCmd := exec.CommandContext(ctx, "aws", "sso", "login", "--profile", awsProfile)
			loginCmd.Stdout = os.Stdout
			loginCmd.Stderr = os.Stderr

			if err := loginCmd.Run(); err != nil {
				return "", errors.Wrap(err, "failed to complete AWS SSO login")
			}
			logger.Info().Msgf("  ✓ AWS SSO login successful for profile %s", awsProfile)
		}

		// Get ECR login password after successful SSO login
		ecrHostname := strings.Split(ecrImage, "/")[0]
		ecrLoginCmd := exec.CommandContext(ctx, "aws", "ecr", "get-login-password", "--region", "us-west-2", "--profile", awsProfile)
		password, passErr := ecrLoginCmd.Output()
		if passErr != nil {
			return "", errors.Wrap(passErr, "failed to get ECR login password")
		}

		// Login to ECR
		dockerLoginCmd := exec.CommandContext(ctx, "docker", "login", "--username", "AWS", "--password-stdin", ecrHostname)
		dockerLoginCmd.Stdin = bytes.NewBuffer(password)
		dockerLoginCmd.Stdout = os.Stdout
		dockerLoginCmd.Stderr = os.Stderr
		if err := dockerLoginCmd.Run(); err != nil {
			return "", errors.Wrap(err, "docker login to ECR failed")
		}
		logger.Info().Msg("  ✓ Docker login to ECR successful")
		// Pull image
		logger.Info().Msgf("🔍 Pulling %s image from ECR...", name)

		pullCmd = exec.CommandContext(ctx, "docker", "pull", ecrImage)
		pullCmd.Stdout = os.Stdout
		pullCmd.Stderr = os.Stderr
		if err := pullCmd.Run(); err != nil {
			return "", errors.Wrapf(err, "failed to pull %s image", name)
		}
	}

	// Tag image
	tagCmd := exec.CommandContext(ctx, "docker", "tag", ecrImage, localImage)
	tagCmd.Stdout = os.Stdout
	tagCmd.Stderr = os.Stderr
	if err := tagCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to tag %s image: %w", name, err)
	}

	logger.Info().Msgf("  ✓ %s image pulled successfully", name)
	return localImage, nil
}

func checkIfGHLIIsInstalled(ctx context.Context, minGHCLIVersion string, noPrompt bool) (installed bool, err error) {
	logger := framework.L

	if isCommandAvailable("gh") {
		logger.Info().Msg("✓ GitHub CLI is already installed")

		ghVersionCmd := exec.CommandContext(ctx, "gh", "--version")
		output, outputErr := ghVersionCmd.Output()
		if outputErr != nil {
			logger.Warn().Msgf("failed to get GH CLI version: %s", outputErr.Error())
			return false, nil
		}

		re := regexp.MustCompile(`gh version (\d+\.\d+\.\d+)`)
		matches := re.FindStringSubmatch(string(output))
		if len(matches) < 2 {
			logger.Warn().Msgf("failed to parse GH CLI version: %s", string(output))
			return false, nil
		}

		version, versionErr := semver.NewVersion(matches[1])
		if versionErr != nil {
			logger.Warn().Msgf("failed to parse GH CLI version: %s", versionErr.Error())
			return false, nil
		}

		isEnoughVersion := version.Compare(semver.MustParse(minGHCLIVersion)) >= 0
		if isEnoughVersion {
			logger.Info().Msgf("  ✓ GitHub CLI is up to date (v%s)", version)
			return true, nil
		}

		logger.Info().Msg("  ✗ GitHub CLI is outdated, upgrading to latest via Homebrew")
		brewInfoCmd := exec.CommandContext(ctx, "brew", "info", "gh")
		brewInfoOutput, brewInfoErr := brewInfoCmd.Output()
		if brewInfoErr != nil {
			fmt.Fprint(os.Stderr, string(brewInfoOutput))
			logger.Warn().Msgf("GH CLI wasn't installed via brew, please update it manually to at least %s", minGHCLIVersion)
			return false, nil
		}

		brewUpgradeCmd := exec.CommandContext(ctx, "brew", "upgrade", "gh")
		brewUpdateOutput, brewUpdateErr := brewUpgradeCmd.Output()
		if brewUpdateErr != nil {
			fmt.Fprint(os.Stderr, string(brewUpdateOutput))
			logger.Warn().Msgf("failed to upgrade GitHub CLI via Homebrew, please update it manually to at least %s", minGHCLIVersion)
			return false, nil
		}
		logger.Info().Msg("  ✓ GitHub CLI upgraded to latest via Homebrew")

		return true, nil
	}

	logger.Info().Msg("Would you like to download and install the GitHub CLI now? (y/n) [y]")

	var input = "y" // Default to yes
	if !noPrompt {
		_, err = fmt.Scanln(&input)
		if err != nil {
			// If error is due to empty input (just pressing Enter), treat as 'y' (yes)
			if err.Error() != "unexpected newline" {
				return false, errors.Wrap(err, "failed to read input")
			}
		}
	}
	// check that input is valid
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "n" {
		logger.Warn().Msg("Invalid input. Please enter 'y' or 'n'.")
		return false, fmt.Errorf("invalid input: %s", input)
	}

	if strings.ToLower(input) != "y" {
		logger.Warn().Msg("  ! You will need to install GitHub CLI manually")
		return false, nil
	}

	logger.Info().Msg("Installing GitHub CLI...")
	installCmd := exec.CommandContext(ctx, "brew", "install", "gh")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return false, errors.Wrap(err, "failed to install GitHub CLI")
	}

	return true, nil
}

func checkGHCli(ctx context.Context, minGHCLIVersion string, noPrompt bool) (installed bool, err error) {
	installed, installErr := checkIfGHLIIsInstalled(ctx, minGHCLIVersion, noPrompt)
	if installErr != nil {
		return false, errors.Wrap(installErr, "failed to check if GitHub CLI is installed")
	}

	if installed {
		loginErr := logInToGithubWithGHCLI(ctx)
		if loginErr != nil {
			return false, errors.Wrap(loginErr, "failed to login to GitHub CLI")
		}
	}

	return installed, nil
}

func logInToGithubWithGHCLI(ctx context.Context) error {
	logger := framework.L
	var outputBuffer bytes.Buffer

	logger.Info().Msg("  Checking GitHub CLI authentication status...")

	ghAuthStatus := exec.CommandContext(ctx, "gh", "auth", "status")
	ghAuthStatus.Stdout = &outputBuffer
	ghAuthStatus.Stderr = &outputBuffer
	statusErr := ghAuthStatus.Run()
	if statusErr == nil {
		logger.Info().Msg("  ✓ GitHub CLI is already authenticated")
		return nil
	}

	// Get the exit code
	var exitError *exec.ExitError
	if !errors.As(statusErr, &exitError) {
		return errors.Wrap(statusErr, "failed to check GitHub CLI authentication status")
	}

	exitCode := exitError.ExitCode()
	logger.Info().Msgf("GitHub CLI authentication status check failed with exit code: %d", exitCode)

	// Exit code 1  means not authenticated
	if exitCode != 1 {
		fmt.Fprintf(os.Stderr, "failed to check GitHub CLI authentication status (exit code: %d): %s\n", exitCode, outputBuffer.String())
		return errors.Wrapf(statusErr, "failed to check GitHub CLI authentication status (exit code: %d)", exitCode)
	}
	logger.Info().Msg("GitHub CLI is not authenticated. Starting login process...")

	logger.Info().Msg("Logging in to GitHub CLI...")

	loginCmd := exec.CommandContext(ctx, "gh", "auth", "login")
	loginCmd.Stdout = os.Stdout
	loginCmd.Stderr = os.Stderr
	if err := loginCmd.Run(); err != nil {
		return errors.Wrap(err, "failed to login to GitHub CLI")
	}

	logger.Info().Msg("  ✓ GitHub CLI logged in successfully")
	return nil
}

func checkBun(ctx context.Context, noPrompt bool) (installed bool, err error) {
	installed, installErr := checkIfBunIsInstalled(ctx, noPrompt)
	if installErr != nil {
		return false, errors.Wrap(installErr, "failed to check if Bun is installed")
	}

	return installed, nil
}

func checkIfBunIsInstalled(ctx context.Context, noPrompt bool) (installed bool, err error) {
	logger := framework.L

	if isCommandAvailable("bun") {
		logger.Info().Msg("✓ Bun is already installed")

		return true, nil
	}

	logger.Info().Msg("Would you like to install Bun now? (y/n) [y]")

	var input = "y" // Default to yes
	if !noPrompt {
		_, err = fmt.Scanln(&input)
		if err != nil {
			// If error is due to empty input (just pressing Enter), treat as 'y' (yes)
			if err.Error() != "unexpected newline" {
				return false, errors.Wrap(err, "failed to read input")
			}
		}
	}
	// check that input is valid
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "n" {
		logger.Warn().Msg("Invalid input. Please enter 'y' or 'n'.")
		return false, fmt.Errorf("invalid input: %s", input)
	}

	if strings.ToLower(input) != "y" {
		logger.Warn().Msg("  ! You will need to install Bun manually")
		return false, nil
	}

	logger.Info().Msg("Installing Bun...")
	tapCmd := exec.CommandContext(ctx, "brew", "tap", "oven-sh/bun")
	tapCmd.Stdout = os.Stdout
	tapCmd.Stderr = os.Stderr
	if err := tapCmd.Run(); err != nil {
		return false, errors.Wrap(err, "failed to tap Bun repository")
	}

	installCmd := exec.CommandContext(ctx, "brew", "install", "bun")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return false, errors.Wrap(err, "failed to install Bun")
	}

	return true, nil
}

func installBunPackages(ctx context.Context) error {
	logger := framework.L
	logger.Info().Msg("Installing Bun packages...")

	installCmd := exec.CommandContext(ctx, "bun", "install")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return errors.Wrap(err, "failed to install Bun packages")
	}

	logger.Info().Msg("  ✓ Bun packages installed successfully")
	return nil
}

func ensurePackageJSON(dir string) error {
	packageJSONPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		return nil
	}

	content := `{
  "name": "typescript-cre-workflow",
  "version": "1.0.0",
  "main": "dist/main.js",
  "private": true,
  "scripts": {
    "postinstall": "bunx cre-setup"
  },
  "license": "UNLICENSED",
  "dependencies": {
    "@chainlink/cre-sdk": "^1.0.0",
    "viem": "2.34.0",
    "zod": "3.25.76"
  },
  "devDependencies": {
    "@types/bun": "1.2.21"
  }
}`

	if err := os.WriteFile(packageJSONPath, []byte(content), 0644); err != nil { //nolint:gosec //G306: Expect WriteFile permissions to be 0600 or less. We want broad read access here.
		return errors.Wrap(err, "failed to create package.json")
	}

	return nil
}
