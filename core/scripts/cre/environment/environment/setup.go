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
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/tracking"
)

var SetupCmd *cobra.Command

func init() {
	var (
		config   SetupConfig
		noPrompt bool
		purge    bool
	)
	SetupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Setup the CRE environment prerequisites",
		Long:  `Checks and sets up prerequisites for the CRE environment including Docker, AWS, Job Distributor, and CRE CLI`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunSetup(cmd.Context(), config, noPrompt, purge)
		},
	}

	SetupCmd.Flags().StringVarP(&config.ConfigPath, "config", "c", DefaultSetupConfigPath, "Path to the TOML configuration file")
	SetupCmd.Flags().BoolVarP(&noPrompt, "no-prompt", "y", false, "Automatically accept defaults and do not prompt for user input")
	SetupCmd.Flags().BoolVarP(&purge, "purge", "p", false, "Purge all existing images and re-download/re-build them")

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
	General        generalConfig        `toml:"general"`
	JobDistributor jobDistributorConfig `toml:"job_distributor"`
	ChipIngress    chipIngressConfig    `toml:"chip_ingress"`
	Capabilities   capabilitiesConfig   `toml:"capabilities"`
}

type generalConfig struct {
	AWSProfile      string `toml:"aws_profile"`
	MinGHCLIVersion string `toml:"min_gh_cli_version"`
	CTFVersion      string `toml:"ctf_version"`
}

type jobDistributorConfig struct {
	BuildConfig BuildConfig `toml:"build_config"`
	PullConfig  PullConfig  `toml:"pull_config"`
}

type chipIngressConfig struct {
	BuildConfig BuildConfig `toml:"build_config"`
	PullConfig  PullConfig  `toml:"pull_config"`
}

type capabilitiesConfig struct {
	TargetPath   string                 `toml:"target_path"`
	Repositories []capabilityRepository `toml:"repositories"`
}

type capabilityRepository struct {
	RepoURL      string `toml:"repository"`
	Branch       string `toml:"branch"`
	BuildCommand string `toml:"build_command"`
	ArtifactsDir string `toml:"artifacts_dir"`
}

var (
	ECR = os.Getenv("AWS_ECR") // TODO this can be moved to an env file
)

const DefaultSetupConfigPath = "configs/setup.toml"

// SetupConfig represents the configuration for the setup command
type SetupConfig struct {
	ConfigPath string
}

type BuildConfig struct {
	RepoURL    string `toml:"repository"`
	LocalRepo  string `toml:"local_repo"`
	Branch     string `toml:"branch"`
	Dockerfile string `toml:"dockerfile"`
	DockerCtx  string `toml:"docker_ctx"`
	LocalImage string `toml:"local_image"`
	PreRun     string `toml:"pre_run"` // Optional function to run before building
}

func checkoutAndBuildRepo(ctx context.Context, logger zerolog.Logger, repo, reference string) (string, bool, error) {
	// Check if repo is a local directory
	isLocalRepo := false
	if _, err2 := os.Stat(repo); err2 == nil {
		fileInfo, err3 := os.Stat(repo)
		if err3 == nil && fileInfo.IsDir() {
			isLocalRepo = true
			logger.Info().Msgf("Using local repository at %s", repo)
		}
	}

	var workingDir string

	if isLocalRepo {
		// Use the local repo path directly
		workingDir = repo
	} else {
		// Create a temporary directory for cloning the remote repo
		tempDir, err2 := os.MkdirTemp("", filepath.Base(repo)+"-*")
		if err2 != nil {
			return "", false, fmt.Errorf("failed to create temporary directory: %w", err2)
		}
		workingDir = tempDir

		// Clone the repository
		logger.Info().Msgf("Cloning repository from %s", repo)
		cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", "--branch", reference, "--single-branch", repo, tempDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err2 := cmd.Run(); err2 != nil {
			return "", false, fmt.Errorf("failed to clone repository: %w", err2)
		}
	}

	return workingDir, isLocalRepo, nil
}

func (c BuildConfig) Build(ctx context.Context) (localImage string, err error) {
	var (
		repo = c.RepoURL
		tag  = c.Branch
	)
	logger := framework.L
	name := strings.ReplaceAll(strings.Split(c.LocalImage, ":")[0], "-", " ")
	name = cases.Title(language.English).String(name)
	logger.Info().Msgf("Building %s image...", name)

	workingDir, isLocalRepo, err := checkoutAndBuildRepo(ctx, logger, repo, tag)
	if err != nil {
		return "", fmt.Errorf("failed to checkout and build repository: %w", err)
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

	// Only checkout specific version if using a git repo and version is specified
	if !isLocalRepo && tag != "" {
		logger.Info().Msgf("Checking out version %s", tag)
		cmd := exec.CommandContext(ctx, "git", "checkout", tag)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to checkout version %s: %w", tag, err)
		}
	}

	// If pre-run function is specified, run it
	if c.PreRun != "" {
		logger.Info().Msgf("Running pre-run step: %s", c.PreRun)
		if err := exec.CommandContext(ctx, "bash", "-c", c.PreRun).Run(); err != nil { //nolint:gosec //G204: Subprocess launched with a potential tainted input or cmd arguments
			return "", fmt.Errorf("pre-run step failed: %w", err)
		}
	}

	// Build Docker image
	cmd := exec.CommandContext(ctx, "docker", "build", "-t", c.LocalImage, "-f", c.Dockerfile, c.DockerCtx) //nolint:gosec //G204: Subprocess launched with a potential tainted input or cmd arguments
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

func (c ImageConfig) Ensure(ctx context.Context, dockerClient *client.Client, awsProfile string, noPrompt bool, purge bool) (localImage string, err error) {
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
		logger.Info().Msgf("Would you like to Pull (requires AWS SSO) or build the %s image? (P/b) [P]", name)

		var input = "b" // Default to Build; TODO default to Pull when AWS access is sorted
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
		if input != "p" && input != "b" {
			logger.Warn().Msg("Invalid input. Please enter 'p' or 'b'.")
			return "", fmt.Errorf("invalid input: %s", input)
		}

		if strings.ToLower(input) == "b" {
			return c.BuildConfig.Build(ctx)
		}

		return c.PullConfig.Pull(ctx, awsProfile)
	}
	return c.BuildConfig.LocalImage, nil
}

// RunSetup performs the setup for the CRE environment
func RunSetup(ctx context.Context, config SetupConfig, noPrompt bool, purge bool) (setupErr error) {
	logger := framework.L
	var localDXTracker tracking.Tracker
	localDXTracker = &tracking.NoOpTracker{}

	defer func() {
		var trackingErr error
		if setupErr != nil {
			trackingErr = localDXTracker.Track("cre.local.setup.result", map[string]any{"result": "failure", "no_prompt": noPrompt, "error": oneLineErrorMessage(setupErr)})
		} else {
			trackingErr = localDXTracker.Track("cre.local.setup.result", map[string]any{"result": "success", "no_prompt": noPrompt})
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
		localDXTracker, trackerErr = tracking.NewDxTracker()
		if trackerErr != nil {
			fmt.Fprintf(os.Stderr, "failed to create DX tracker: %s\n", trackerErr)
		}
	}

	jdConfig := ImageConfig{
		BuildConfig: cfg.JobDistributor.BuildConfig,
		PullConfig:  cfg.JobDistributor.PullConfig,
	}

	jdLocalImage, jdErr := jdConfig.Ensure(ctx, dockerClient, cfg.General.AWSProfile, noPrompt, purge)
	if jdErr != nil {
		setupErr = errors.Wrap(jdErr, "failed to ensure Job Distributor image")
		return
	}

	chipConfig := ImageConfig{
		BuildConfig: cfg.ChipIngress.BuildConfig,
		PullConfig:  cfg.ChipIngress.PullConfig,
	}

	chipLocalImage, chipErr := chipConfig.Ensure(ctx, dockerClient, cfg.General.AWSProfile, noPrompt, purge)
	if chipErr != nil {
		setupErr = errors.Wrap(chipErr, "failed to ensure Atlas Chip Ingress image")
		return
	}

	ctfInstalled, ctfErr := checkCTF(ctx, cfg.General.CTFVersion, noPrompt, purge)
	if ctfErr != nil {
		setupErr = errors.Wrap(ctfErr, "failed to ensure CTF CLI")
		return
	}

	buildErr := buildCapabilityBinaries(ctx, cfg.Capabilities)
	if buildErr != nil {
		setupErr = errors.Wrap(buildErr, "failed to build capabilities")
		return
	}

	// Print summary
	fmt.Println()
	logger.Info().Msg("✅ Setup Summary:")
	logger.Info().Msg("   ✓ Docker is installed and configured correctly")
	logger.Info().Msgf("   ✓ Job Distributor image %s is available", jdLocalImage)
	logger.Info().Msgf("   ✓ Atlas Chip Ingress image %s is available", chipLocalImage)
	if ghCli {
		logger.Info().Msg("   ✓ GitHub CLI is installed")
	} else {
		logger.Warn().Msg("   ✗ GitHub CLI is not installed")
	}
	if ctfInstalled {
		logger.Info().Msg("   ✓ CTF CLI is installed")
	} else {
		logger.Warn().Msg("   ✗ CTF CLI is not installed")
	}
	if len(cfg.Capabilities.Repositories) > 0 {
		logger.Info().Msg("   ✓ Capabilities binaries built")
	} else {
		logger.Warn().Msg("   ✗ Capabilities binaries not built")
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

	buildErr := buildCapabilityBinaries(ctx, cfg.Capabilities)
	if buildErr != nil {
		return errors.Wrap(buildErr, "failed to build capabilities")
	}

	fmt.Println()
	logger := framework.L
	logger.Info().Msg("✅ Build Capabilities Summary:")
	for _, repo := range cfg.Capabilities.Repositories {
		logger.Info().Msgf("   ✓ %s", repo.RepoURL)
	}

	return nil
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

func buildCapabilityBinaries(ctx context.Context, capabilitiesConfig capabilitiesConfig) error {
	logger := framework.L
	logger.Info().Msg("🔍 Building capabilities binaries...")

	// Save current directory and change to working directory
	currentDir, cErr := os.Getwd()
	if cErr != nil {
		return fmt.Errorf("failed to get current directory: %w", cErr)
	}

	dirsToDelete := []string{}

	for _, repo := range capabilitiesConfig.Repositories {
		logger.Info().Msgf("🔍 Building %s...", repo.RepoURL)

		workingDir, isLocalRepo, err := checkoutAndBuildRepo(ctx, logger, repo.RepoURL, repo.Branch)
		if err != nil {
			return fmt.Errorf("failed to checkout and build repository: %w", err)
		}

		if !isLocalRepo {
			dirsToDelete = append(dirsToDelete, workingDir)
		}

		if err := os.Chdir(workingDir); err != nil {
			return fmt.Errorf("failed to change to working directory: %w", err)
		}

		// Only checkout specific version if using a git repo and version is specified
		if !isLocalRepo && repo.Branch != "" {
			logger.Info().Msgf("Checking out version %s", repo.Branch)
			cmd := exec.CommandContext(ctx, "git", "checkout", repo.Branch) //nolint:gosec //G204: Subprocess launched with a potential tainted input or cmd arguments
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to checkout version %s: %w", repo.Branch, err)
			}
		}

		// Run build command
		logger.Info().Msgf("Running build command: %s", repo.BuildCommand)
		cmd := exec.CommandContext(ctx, "bash", "-c", repo.BuildCommand) //nolint:gosec //G204: Subprocess launched with a potential tainted input or cmd arguments
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run build command: %w", err)
		}

		// Copy artifacts to target path
		targetPath := filepath.Join(currentDir, capabilitiesConfig.TargetPath)
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}

		logger.Info().Msgf("Copying build artifacts from %s to %s", repo.ArtifactsDir, targetPath)
		artifactsDir := filepath.Join(workingDir, repo.ArtifactsDir)
		copyCmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("cp -r %s/* %s/", artifactsDir, targetPath)) //nolint:gosec //G204: Subprocess launched with a potential tainted input or cmd arguments
		if err := copyCmd.Run(); err != nil {
			return fmt.Errorf("failed to copy directory: %w", err)
		}

		logger.Info().Msgf("✓ Build artifacts copied to %s", targetPath)
	}

	defer func() {
		_ = os.Chdir(currentDir)
		for _, dir := range dirsToDelete {
			_ = os.RemoveAll(dir)
		}
	}()

	return nil
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

	// Check if AWS profile exists
	configureCmd := exec.Command("aws", "configure", "list-profiles")
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

	pullCmd := exec.CommandContext(ctx, "docker", "pull", ecrImage)
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr
	if err := pullCmd.Run(); err != nil {
		return "", errors.Wrapf(err, "failed to pull %s image", name)
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

		ghVersionCmd := exec.Command("gh", "--version")
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
		brewInfoCmd := exec.Command("brew", "info", "gh")
		brewInfoOutput, brewInfoErr := brewInfoCmd.Output()
		if brewInfoErr != nil {
			fmt.Fprint(os.Stderr, string(brewInfoOutput))
			logger.Warn().Msgf("GH CLI wasn't installed via brew, please update it manually to at least %s", minGHCLIVersion)
			return false, nil
		}

		brewUpgradeCmd := exec.Command("brew", "upgrade", "gh")
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

// checkCTF checks if the CTF CLI is installed prompts to install if not
func checkCTF(ctx context.Context, requiredVersion string, noPrompt bool, purge bool) (installed bool, err error) {
	logger := framework.L

	if purge {
		_ = os.Remove(filepath.Join(binDir, "ctf"))
	}
	// Check for CTF CLI is in binDir
	if _, statErr := os.Stat(filepath.Join(binDir, "ctf")); statErr == nil {
		logger.Info().Msg("✓ CTF CLI is already installed")
		return true, nil
	}

	logger.Info().Msg("✗ CTF CLI is not installed")
	logger.Info().Msg("  Would you like to download and install the CTF CLI now? (y/n) [y]")

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
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "n" {
		logger.Warn().Msg("Invalid input. Please enter 'y' or 'n'.")
		return false, fmt.Errorf("invalid input: %s", input)
	}

	if strings.ToLower(input) != "y" {
		logger.Warn().Msg("  ! You will need to install CTF CLI manually")
		return false, nil
	}

	logger.Info().Msgf("Installing CTF CLI v%s...", requiredVersion)
	// change to temp directory and download and extract; change back to original directory on exit
	tempDir, err := os.MkdirTemp("", "ctf-download")
	if err != nil {
		return false, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)
	wd, _ := os.Getwd()
	if err := os.Chdir(tempDir); err != nil {
		return false, fmt.Errorf("failed to change directory: %w", err)
	}
	defer func() { _ = os.Chdir(wd) }()
	// install ctf by pulling from GitHub releases https://github.com/smartcontractkit/chainlink-testing-framework/releases/tag/framework%2Fv0.10.3
	// for MacOS users, it will be framework-vX.X.X-darwin-arm64.tar.gz
	// for Linux users, it will be framework-vX.X.X-linux-amd64.tar.gz or framework-vX.X.X-linux-arm64.tar.gz
	osType := runtime.GOOS
	archType := runtime.GOARCH
	archiveName := fmt.Sprintf("framework-v%s-%s-%s.tar.gz", requiredVersion, osType, archType)
	// gh release download framework/v0.10.3 --repo smartcontractkit/chainlink-testing-framework --pattern framework-v0.10.3-darwin-arm64.tar.gz
	cmd := exec.CommandContext(ctx, "gh", "release", "download", "framework/v"+requiredVersion, "--repo", "smartcontractkit/chainlink-testing-framework", "--pattern", archiveName) //nolint:gosec //G204: Subprocess launched with a potential tainted input or cmd arguments
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err2 := cmd.Run(); err2 != nil {
		return false, fmt.Errorf("failed to download CTF CLI: %w", err2)
	}

	logger.Info().Msg("Extracting CTF CLI...")
	cmd = exec.CommandContext(ctx, "tar", "-C", binDir, "-xf", archiveName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err2 := cmd.Run(); err2 != nil {
		return false, fmt.Errorf("failed to extract CTF CLI: %w", err2)
	}
	// Remove archive
	if err2 := os.Remove(archiveName); err2 != nil {
		logger.Warn().Msgf("Failed to remove %s. Please remove it manually.", archiveName)
	}

	logger.Info().Msgf("  ✓ CTF CLI installed to %s/ctf", binDir)
	logger.Warn().Msg("")
	logger.Warn().Msgf("   * -------------------------- I M P O R T A N T -------------------------------------- *")
	logger.Warn().Msgf("   *                                                                                     *")
	logger.Warn().Msgf("   * Add this directory to your PATH or move the CTF binary to a directory in your PATH  *")
	logger.Warn().Msgf("   *                                                                                     *")
	logger.Warn().Msgf("   * ----------------------------------------------------------------------------------- *")
	logger.Warn().Msg("")
	logger.Warn().Msgf("   You can run: export PATH=\"%s:$PATH\"", binDir)
	logger.Warn().Msg("")
	return true, nil
}
