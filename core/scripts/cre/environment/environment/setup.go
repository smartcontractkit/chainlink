package environment

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/docker/docker/client"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// TODO this can move to the toml configuration file
const (
	awsProfile    = "sdlc"
	creCLIVersion = "0.2.1"
)

// TODO these can move to the toml configuration file
var (
	ECR           = os.Getenv("AWS_ECR") // TODO this can be moved to an env file
	jdTag         = "0.12.7"
	JDBuildConfig = BuildConfig{
		RepoURL:    "https://github.com/smartcontractkit/job-distributor",
		Branch:     "v" + jdTag,
		Dockerfile: "e2e/Dockerfile.e2e",
		Dir:        ".",
		LocalImage: "job-distributor:" + jdTag,
	}
	JDPullConfig = PullConfig{
		LocalImage: "job-distributor:" + jdTag,
		EcrImage:   fmt.Sprintf("%s/job-distributor:%s", ECR, jdTag),
	}

	JDImageConfig = ImageConfig{
		BuildConfig: JDBuildConfig,
		PullConfig:  JDPullConfig,
	}

	chipRemoteTag = "qa-latest" // no released version yet. sha 1a9726faa5fe1d45138ca89143655e309ff65ae50cd3db5631f2b401c54d0c1f

	ChipBuildConfig = BuildConfig{
		RepoURL:    "https://github.com/smartcontractkit/atlas",
		Branch:     "master",
		Dockerfile: "chip-ingress/Dockerfile",
		Dir:        "chip-ingress",
		LocalImage: "chip-ingress:local-cre",
		PreRun:     chipVendor,
	}
	ChipPullConfig = PullConfig{
		LocalImage: "chip-ingress:local-cre",
		EcrImage:   fmt.Sprintf("%s/atlas-chip-ingress:%s", ECR, chipRemoteTag),
	}
	ChipImageConfig = ImageConfig{
		BuildConfig: ChipBuildConfig,
		PullConfig:  ChipPullConfig,
	}
)

// SetupConfig represents the configuration for the setup command
type SetupConfig struct {
	ConfigPath string
}

type BuildConfig struct {
	RepoURL    string
	LocalRepo  string
	Branch     string
	Dockerfile string
	Dir        string
	LocalImage string
	PreRun     func(ctx context.Context, c BuildConfig) error // Optional function to run before building
}

func (c BuildConfig) Build(ctx context.Context) (localImage string, err error) {
	var (
		repo = c.RepoURL
		tag  = c.Branch
	)
	logger := framework.L
	name := strings.ReplaceAll(strings.Split(localImage, ":")[0], "-", " ")
	name = cases.Title(language.English).String(name)
	logger.Info().Msgf("Building %s image...", name)

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
			return "", fmt.Errorf("failed to create temporary directory: %w", err2)
		}
		defer os.RemoveAll(tempDir)
		workingDir = tempDir

		// Clone the repository
		logger.Info().Msgf("Cloning repository from %s", repo)
		cmd := exec.CommandContext(ctx, "git", "clone", repo, tempDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err2 := cmd.Run(); err2 != nil {
			return "", fmt.Errorf("failed to clone repository: %w", err2)
		}
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
	if c.PreRun != nil {
		if err := c.PreRun(ctx, c); err != nil {
			return "", fmt.Errorf("pre-run step failed: %w", err)
		}
	}

	// Build Docker image
	cmd := exec.CommandContext(ctx, "docker", "build", "-t", c.LocalImage, "-f", c.Dockerfile, c.Dir) //nolint:gosec //G204: Subprocess launched with a potential tainted input or cmd arguments
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
	LocalImage string
	EcrImage   string
}

func (c PullConfig) Pull(ctx context.Context) (localImage string, err error) {
	if ECR == "" {
		return "", errors.New("AWS_ECR environment variable is not set. See README for more details and references to find the correct ECR URL or visit https://smartcontract-it.atlassian.net/wiki/spaces/INFRA/pages/1045495923/Configure+the+AWS+CLI")
	}
	return pullImage(ctx, c.LocalImage, c.EcrImage)
}

type ImageConfig struct {
	BuildConfig BuildConfig
	PullConfig  PullConfig
}

func (c ImageConfig) Ensure(ctx context.Context, dockerClient *client.Client) (localImage string, err error) {
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

		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			// If error is due to empty input (just pressing Enter), use default
			if err.Error() != "unexpected newline" {
				return "", errors.Wrap(err, "failed to read input")
			}

			input = "p" // Default to Pull
		}

		if strings.ToLower(input) == "b" {
			return c.BuildConfig.Build(ctx)
		}

		return c.PullConfig.Pull(ctx)
	}
	return c.BuildConfig.LocalImage, nil
}

var SetupCmd *cobra.Command

func init() {
	var config SetupConfig
	SetupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Setup the CRE environment prerequisites",
		Long:  `Checks and sets up prerequisites for the CRE environment including Docker, AWS, Job Distributor, and CRE CLI`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunSetup(cmd.Context(), config)
		},
	}

	SetupCmd.Flags().StringVarP(&config.ConfigPath, "config", "c", "", "Path to the TOML configuration file")
	// _ = SetupCmd.MarkFlagRequired("config")

	EnvironmentCmd.AddCommand(SetupCmd)
}

// RunSetup performs the setup for the CRE environment
func RunSetup(ctx context.Context, config SetupConfig) error {
	logger := framework.L
	logger.Info().Msg("🔍 Checking prerequisites for CRE environment...")

	// Check if Docker is installed
	if !isCommandAvailable("docker") {
		return errors.New("docker is not installed. Please install Docker and try again")
	}
	logger.Info().Msg("✓ Docker is installed")

	// Check if Docker is running
	dockerClient, clientErr := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if clientErr != nil {
		return errors.Wrap(clientErr, "failed to create Docker client")
	}

	_, pingErr := dockerClient.Ping(ctx)
	if pingErr != nil {
		return errors.Wrap(pingErr, "docker is not running. Please start Docker and try again")
	}
	logger.Info().Msg("✓ Docker is running")

	// Check Docker configuration
	if dockerConfigErr := checkDockerConfiguration(); dockerConfigErr != nil {
		return errors.Wrap(dockerConfigErr, "failed to check Docker configuration")
	}

	// Check if AWS CLI is installed
	if !isCommandAvailable("aws") {
		return errors.New("AWS CLI is not installed. Please install AWS CLI and try again")
	}
	logger.Info().Msg("✓ AWS CLI is installed")

	ghCli, ghCliErr := checkGHCli(ctx)
	if ghCliErr != nil {
		return errors.Wrap(ghCliErr, "failed to ensure GitHub CLI")
	}

	jdLocalImage, jdErr := JDImageConfig.Ensure(ctx, dockerClient)
	if jdErr != nil {
		return errors.Wrap(jdErr, "failed to ensure Job Distributor image")
	}
	chipLocalImage, chipErr := ChipImageConfig.Ensure(ctx, dockerClient)
	if chipErr != nil {
		return errors.Wrap(chipErr, "failed to ensure Atlas Chip Ingress image")
	}

	creCLI, creCliErr := checkCRECLI(ctx)
	if creCliErr != nil {
		return errors.Wrap(creCliErr, "failed to ensure CRE CLI")
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
	if creCLI {
		logger.Info().Msg("   ✓ CRE CLI is installed")
	} else {
		logger.Warn().Msg("   ✗ CRE CLI is not installed")
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
			return errors.New("docker settings file not found at expected macOS locations")
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
			if value == expected {
				logger.Info().Msgf("  ✓ %s is correctly set to %s", setting, expected)
			} else {
				logger.Error().Msgf("  ✗ %s is set to %s (should be %s)", setting, value, expected)
				dockerSettingsOK = false
			}
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
func pullImage(ctx context.Context, localImage, ecrImage string) (string, error) {
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

func checkIfGHLIIsInstalled(ctx context.Context) (installed bool, err error) {
	logger := framework.L

	if isCommandAvailable("gh") {
		logger.Info().Msg("✓ GitHub CLI is already installed")
		return true, nil
	}

	logger.Info().Msg("Would you like to download and install the GitHub CLI now? (y/n) [y]")

	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		// If error is due to empty input (just pressing Enter), treat as 'y' (yes)
		if err.Error() != "unexpected newline" {
			return false, errors.Wrap(err, "failed to read input")
		}
		input = "y"
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

func checkGHCli(ctx context.Context) (installed bool, err error) {
	installed, installErr := checkIfGHLIIsInstalled(ctx)
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

// checkCRECLI checks if the CRE CLI is installed
func checkCRECLI(ctx context.Context) (installed bool, err error) {
	logger := framework.L

	// Check for CRE CLI
	osType := runtime.GOOS
	archType := runtime.GOARCH

	creBinaryName := fmt.Sprintf("cre_v%s_%s_%s", creCLIVersion, osType, archType)
	if isCommandAvailable(creBinaryName) || isCommandAvailable("cre") {
		logger.Info().Msg("✓ CRE CLI is already installed")
		return true, nil
	}

	// CRE CLI not found
	logger.Info().Msg("✗ CRE CLI is not installed")
	logger.Info().Msg("  Would you like to download and install the CRE CLI now? (y/n) [y]")

	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		// If error is due to empty input (just pressing Enter), treat as 'n' (no)
		if err.Error() != "unexpected newline" {
			return false, errors.Wrap(err, "failed to read input")
		}
		input = "y"
	}

	if strings.ToLower(input) != "y" {
		logger.Warn().Msg("  ! You will need to install CRE CLI manually")
		return false, nil
	}

	// Download CRE CLI
	logger.Info().Msgf("  Downloading CRE CLI v%s for %s_%s...", creCLIVersion, osType, archType)
	archivePattern := fmt.Sprintf("*%s_%s.tar.gz", osType, archType)
	cmd := exec.CommandContext(ctx, "gh", "release", "download", "v"+creCLIVersion, "--repo", "smartcontractkit/dev-platform", "--pattern", archivePattern)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err2 := cmd.Run(); err2 != nil {
		return false, fmt.Errorf("failed to download CRE CLI: %w", err2)
	}

	// Extract archive
	archiveName := fmt.Sprintf("cre_v%s_%s_%s.tar.gz", creCLIVersion, osType, archType)
	logger.Info().Msg("  Extracting CRE CLI...")
	cmd = exec.CommandContext(ctx, "tar", "-xf", archiveName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err2 := cmd.Run(); err2 != nil {
		return false, fmt.Errorf("failed to extract CRE CLI: %w", err2)
	}

	// Remove archive
	if err2 := os.Remove(archiveName); err2 != nil {
		logger.Warn().Msgf("Failed to remove %s. Please remove it manually.", archiveName)
	}

	// Remove quarantine attribute on macOS
	if osType == "darwin" {
		cmd = exec.CommandContext(ctx, "xattr", "-d", "com.apple.quarantine", creBinaryName)
		_ = cmd.Run() // Ignore errors
	}

	// Make executable
	if err2 := os.Chmod(creBinaryName, 0755); err2 != nil {
		return false, fmt.Errorf("failed to make CRE CLI executable: %w", err2)
	}

	// Create symlink
	if err2 := os.Symlink(creBinaryName, "cre"); err2 != nil && !os.IsExist(err2) {
		return false, fmt.Errorf("failed to create symlink: %w", err2)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("failed to get current directory: %w", err)
	}

	logger.Info().Msgf("  ✓ CRE CLI installed to %s/cre", currentDir)
	logger.Warn().Msgf("  ! Add this directory to your PATH or move the CRE binary to a directory in your PATH")
	logger.Info().Msgf("   You can run: export PATH=\"%s:$PATH\"", currentDir)

	return true, nil
}

// chipVendor changes to the directory specified in the config
// and executes go mod vendor command
func chipVendor(ctx context.Context, config BuildConfig) error {
	logger := framework.L

	// Save current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Change to the target directory
	logger.Info().Msgf("Changing directory to %s", config.Dir)
	if err := os.Chdir(config.Dir); err != nil {
		return fmt.Errorf("failed to change to directory %s: %w", config.Dir, err)
	}

	// Restore original directory when function completes
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			logger.Error().Err(err).Msg("Failed to restore original directory")
		}
	}()

	// Execute go mod vendor
	logger.Info().Msg("Running go mod vendor...")
	cmd := exec.CommandContext(ctx, "go", "mod", "vendor")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod vendor failed: %w", err)
	}

	logger.Info().Msg("Vendor directory successfully created")
	return nil
}
