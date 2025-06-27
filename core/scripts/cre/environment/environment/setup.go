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
}

func (c BuildConfig) Build(ctx context.Context) (localImage string, err error) {
	return buildImage(ctx, c.RepoURL, c.Branch, c.Dockerfile, c.Dir, c.LocalImage)
}

type PullConfig struct {
	LocalImage string
	EcrImage   string
}

func (c PullConfig) Pull(ctx context.Context) (localImage string, err error) {
	if ECR == "" {
		return "", errors.New("AWS_ECR environment variable is not set. See README for more details and references to find the correct ECR URL")
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
		logger.Info().Msgf("Would you like to Pull (requires AWS SSO) or build the %s image? (P/b)", name)

		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
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
	_ = SetupCmd.MarkFlagRequired("config")

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
	dockerClient, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	_, err = dockerClient.Ping(ctx)
	if err != nil {
		return fmt.Errorf("docker is not running. Please start Docker and try again: %w", err)
	}
	logger.Info().Msg("✓ Docker is running")

	// Check Docker configuration
	if err2 := checkDockerConfiguration(ctx); err2 != nil {
		return err2
	}

	// Check if AWS CLI is installed
	if !isCommandAvailable("aws") {
		return errors.New("AWS CLI is not installed. Please install AWS CLI and try again")
	}
	logger.Info().Msg("✓ AWS CLI is installed")

	jdLocalImage, err := JDImageConfig.Ensure(ctx, dockerClient)
	if err != nil {
		return fmt.Errorf("failed to ensure Job Distributor image: %w", err)
	}
	chipLocalImage, err := ChipImageConfig.Ensure(ctx, dockerClient)
	if err != nil {
		return fmt.Errorf("failed to ensure Atlas Chip Ingress image: %w", err)
	}
	creCLI, err := checkCRECLI(ctx)
	if err != nil {
		return fmt.Errorf("failed to ensure CRE CLI: %w", err)
	}
	// Print summary
	logger.Info().Msg("\n✅ Setup Summary:")
	logger.Info().Msg("   ✓ Docker is installed and configured correctly")
	logger.Info().Msgf("   ✓ Job Distributor image %s is available", jdLocalImage)
	logger.Info().Msgf("   ✓ Atlas Chip Ingress image %s is available", chipLocalImage)
	if creCLI {
		logger.Info().Msg("   ✓ CRE CLI is installed")
	} else {
		logger.Warn().Msg("   ✗ CRE CLI is not installed")
	}

	logger.Info().Msg("\n🚀 Next Steps:")
	logger.Info().Msg("1. Navigate to the CRE environment directory: cd core/scripts/cre/environment")
	logger.Info().Msg("2. Start the environment: go run main.go env start")
	logger.Info().Msg("   Optional: Add --with-example to start with an example workflow")
	logger.Info().Msg("   Optional: Add --with-plugins-docker-image to use a pre-built image with capabilities")
	logger.Info().Msg("\nFor more information, see the documentation in core/scripts/cre/environment/docs.md")

	return nil
}

// isCommandAvailable checks if a command is available in the PATH
func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// checkDockerConfiguration checks if Docker is configured correctly
func checkDockerConfiguration(ctx context.Context) error {
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
		logger.Info().Msgf("  ✓ %s image (%s) is available from local build", name, localImage)
		return true, nil
	}

	// Check if ECR image exists
	_, err = dockerClient.ImageInspect(ctx, ecrImage)
	if err == nil {
		logger.Info().Msgf("  ✓ %s image (%s) is available", name, ecrImage)
		// Tag ECR image as local image
		if err := dockerClient.ImageTag(ctx, ecrImage, localImage); err != nil {
			return false, fmt.Errorf("failed to tag %s image: %w", name, err)
		}
		logger.Info().Msgf("  ✓ %s image tagged as %s", name, localImage)
		return true, nil
	}
	return false, nil
}

// buildImage builds the Job Distributor image
func buildImage(ctx context.Context, repo, tag, dockerFile, dir, localImage string) (string, error) {
	logger := framework.L
	name := strings.ReplaceAll(strings.Split(localImage, ":")[0], "-", " ")
	name = cases.Title(language.English).String(name)
	logger.Info().Msgf("Building %s image...", name)

	// Check if repo is a local directory
	isLocalRepo := false
	if _, err := os.Stat(repo); err == nil {
		fileInfo, err := os.Stat(repo)
		if err == nil && fileInfo.IsDir() {
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
		tempDir, err := os.MkdirTemp("", filepath.Base(repo)+"-*")
		if err != nil {
			return "", fmt.Errorf("failed to create temporary directory: %w", err)
		}
		defer os.RemoveAll(tempDir)
		workingDir = tempDir

		// Clone the repository
		logger.Info().Msgf("Cloning repository from %s", repo)
		cmd := exec.CommandContext(ctx, "git", "clone", repo, tempDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to clone repository: %w", err)
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

	// Build Docker image
	cmd := exec.CommandContext(ctx, "docker", "build", "-t", localImage, "-f", dockerFile, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to build Docker image: %w", err)
	}

	logger.Info().Msgf("  ✓ %s image built successfully", name)
	return localImage, nil
}

// pullImage pulls the Job Distributor image from ECR
func pullImage(ctx context.Context, localImage, ecrImage string) (string, error) {
	logger := framework.L
	name := strings.ReplaceAll(strings.Split(localImage, ":")[0], "-", " ")
	name = cases.Title(language.English).String(name)

	// Check if AWS profile exists
	cmd := exec.Command("aws", "configure", "list-profiles")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to list AWS profiles: %w", err)
	}

	if !strings.Contains(string(output), awsProfile) {
		return "", fmt.Errorf("AWS profile '%s' not found. Please ensure you have the correct AWS profile configured", awsProfile)
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
		cmd = exec.CommandContext(ctx, "aws", "sso", "login", "--profile", awsProfile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		password, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("AWS SSO login failed: %w", err)
		}
		logger.Info().Msgf("  ✓ AWS SSO login successful for profile %s", awsProfile)

		// Login to ECR
		ecrHostname := strings.Split(ecrImage, "/")[0]
		dockerLoginCmd := exec.CommandContext(ctx, "docker", "login", "--username", "AWS", "--password-stdin", ecrHostname)
		dockerLoginCmd.Stdin = bytes.NewBuffer(password)
		dockerLoginCmd.Stdout = os.Stdout
		dockerLoginCmd.Stderr = os.Stderr
		if err := dockerLoginCmd.Run(); err != nil {
			return "", fmt.Errorf("docker login to ECR failed: %w", err)
		}
		logger.Info().Msg("  ✓ Docker login to ECR successful")
	}
	// Pull image
	logger.Info().Msgf("🔍 Pulling %s image from ECR...", name)

	cmd = exec.CommandContext(ctx, "docker", "pull", ecrImage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to pull %s image: %w", name, err)
	}

	// Tag image
	cmd = exec.CommandContext(ctx, "docker", "tag", ecrImage, localImage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to tag %s image: %w", name, err)
	}

	logger.Info().Msgf("  ✓ %s image pulled successfully", name)
	return localImage, nil
}

// checkCRECLI checks if the CRE CLI is installed
func checkCRECLI(ctx context.Context) (installed bool, err error) {
	logger := framework.L

	// Check for CRE CLI
	osType := runtime.GOOS
	archType := runtime.GOARCH

	creBinaryName := fmt.Sprintf("cre_v%s_%s_%s", creCLIVersion, osType, archType)
	if isCommandAvailable(creBinaryName) || isCommandAvailable("cre") {
		logger.Info().Msg("  ✓ CRE CLI is already installed")
		return true, nil
	}

	// CRE CLI not found
	logger.Info().Msg("  ✗ CRE CLI is not installed")
	logger.Info().Msg("Would you like to download and install the CRE CLI now? (y/n)")

	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	if strings.ToLower(input) != "y" {
		logger.Warn().Msg("  ! You will need to install CRE CLI manually")
		return false, nil
	}

	// Check for GitHub CLI
	if !isCommandAvailable("gh") {
		return false, fmt.Errorf("GitHub CLI is not installed. Please install GitHub CLI or download CRE CLI manually from https://github.com/smartcontractkit/dev-platform/releases/tag/v%s", creCLIVersion)
	}

	// Download CRE CLI
	logger.Info().Msgf("Downloading CRE CLI v%s for %s_%s...", creCLIVersion, osType, archType)
	archivePattern := fmt.Sprintf("*%s_%s.tar.gz", osType, archType)
	cmd := exec.CommandContext(ctx, "gh", "release", "download", "v"+creCLIVersion, "--repo", "smartcontractkit/dev-platform", "--pattern", archivePattern)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err2 := cmd.Run(); err2 != nil {
		return false, fmt.Errorf("failed to download CRE CLI: %w", err2)
	}

	// Extract archive
	archiveName := fmt.Sprintf("cre_v%s_%s_%s.tar.gz", creCLIVersion, osType, archType)
	logger.Info().Msg("Extracting CRE CLI...")
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
