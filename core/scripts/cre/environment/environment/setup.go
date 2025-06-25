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

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// TODO this can move to the toml configuration file
const (
	JD_ECR             = "804282218731.dkr.ecr.us-west-2.amazonaws.com"
	AWS_PROFILE        = "sdlc"
	CRE_CLI_VERSION    = "0.2.1"
	JD_VERSION_DEFAULT = "0.12.7"
)

// SetupConfig represents the configuration for the setup command
type SetupConfig struct {
	ConfigPath string
	JDVersion  string
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
	SetupCmd.Flags().StringVarP(&config.JDVersion, "jd-version", "j", JD_VERSION_DEFAULT, "Version of Job Distributor to use")
	_ = SetupCmd.MarkFlagRequired("config")

	EnvironmentCmd.AddCommand(SetupCmd)
}

// RunSetup performs the setup for the CRE environment
func RunSetup(ctx context.Context, config SetupConfig) error {
	logger := framework.L
	logger.Info().Msg("🔍 Checking prerequisites for CRE environment...")

	// Check if Docker is installed
	if !isCommandAvailable("docker") {
		return fmt.Errorf("Docker is not installed. Please install Docker and try again")
	}
	logger.Info().Msg("✓ Docker is installed")

	// Check if Docker is running
	dockerClient, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	_, err = dockerClient.Ping(ctx)
	if err != nil {
		return fmt.Errorf("Docker is not running. Please start Docker and try again: %w", err)
	}
	logger.Info().Msg("✓ Docker is running")

	// Check Docker configuration
	if err := checkDockerConfiguration(ctx); err != nil {
		return err
	}

	// Check if AWS CLI is installed
	if !isCommandAvailable("aws") {
		return fmt.Errorf("AWS CLI is not installed. Please install AWS CLI and try again")
	}
	logger.Info().Msg("✓ AWS CLI is installed")

	/*
		// Load TOML configuration
		in, err := framework.Load[Config](nil)
		if err != nil {
			return fmt.Errorf("failed to load test configuration: %w", err)
		}
	*/
	// Set JD version if provided
	jdVersion := config.JDVersion
	jdLocalImage := fmt.Sprintf("job-distributor:%s", jdVersion)
	jdEcrImage := fmt.Sprintf("%s/%s", JD_ECR, jdLocalImage)

	// Check Job Distributor image
	logger.Info().Msg("\n🔍 Checking for Job Distributor image...")
	if err := checkJobDistributorImage(ctx, dockerClient, jdLocalImage, jdEcrImage); err != nil {
		return err
	}

	// Check CRE CLI
	logger.Info().Msg("\n🔍 Checking for CRE CLI...")
	if err := checkCRECLI(ctx); err != nil {
		return err
	}
	/*
		// Check Capability Binaries based on config
		if err := checkCapabilityBinaries(ctx, in.ExtraCapabilities); err != nil {
			return err
		}
	*/
	// Print summary
	logger.Info().Msg("\n✅ Setup Summary:")
	logger.Info().Msg("   ✓ Docker is installed and configured correctly")
	logger.Info().Msgf("   ✓ Job Distributor image %s is available", jdLocalImage)
	logger.Info().Msg("   ✓ CRE CLI is installed")

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
			return fmt.Errorf("Docker settings file not found at expected macOS locations")
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
		return fmt.Errorf("Docker is not configured correctly. Please fix the issues and try again")
	}

	return nil
}

// checkJobDistributorImage checks if the Job Distributor image is available
func checkJobDistributorImage(ctx context.Context, dockerClient *client.Client, jdLocalImage, jdEcrImage string) error {
	logger := framework.L
	// Check if local image exists
	_, err := dockerClient.ImageInspect(ctx, jdLocalImage)
	if err == nil {
		logger.Info().Msgf("  ✓ Job Distributor image (%s) is available from local build", jdLocalImage)
		return nil
	}

	// Check if ECR image exists
	_, err = dockerClient.ImageInspect(ctx, jdEcrImage)
	if err == nil {
		logger.Info().Msgf("  ✓ Job Distributor image (%s) is available", jdEcrImage)
		// Tag ECR image as local image
		if err := dockerClient.ImageTag(ctx, jdEcrImage, jdLocalImage); err != nil {
			return fmt.Errorf("failed to tag Job Distributor image: %w", err)
		}
		logger.Info().Msgf("  ✓ Job Distributor image tagged as %s", jdLocalImage)
		return nil
	}

	// Ask user what to do
	logger.Info().Msgf("  ✗ Job Distributor image (%s) is not available", jdLocalImage)
	logger.Info().Msg("Would you like to Pull (requires AWS SSO) or build the Job Distributor image? (P/b)")

	var input string
	fmt.Scanln(&input)

	if strings.ToLower(input) == "b" {
		return buildJobDistributorImage(ctx, jdLocalImage, strings.TrimPrefix(filepath.Base(jdLocalImage), "job-distributor:"))
	}

	return pullJobDistributorImage(ctx, jdLocalImage, jdEcrImage)
}

// buildJobDistributorImage builds the Job Distributor image
func buildJobDistributorImage(ctx context.Context, jdLocalImage, jdVersion string) error {
	logger := framework.L
	logger.Info().Msg("Building Job Distributor image...")

	// Create a temporary directory for cloning
	tempDir, err := os.MkdirTemp("", "job-distributor-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Clone the repository
	cmd := exec.CommandContext(ctx, "git", "clone", "https://github.com/smartcontractkit/job-distributor", tempDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Save current directory and change to temp dir
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		return fmt.Errorf("failed to change to temporary directory: %w", err)
	}
	defer os.Chdir(currentDir)

	// Checkout specific version
	cmd = exec.CommandContext(ctx, "git", "checkout", fmt.Sprintf("v%s", jdVersion))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout version v%s: %w", jdVersion, err)
	}

	// Build Docker image
	cmd = exec.CommandContext(ctx, "docker", "build", "-t", jdLocalImage, "-f", "e2e/Dockerfile.e2e", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build Docker image: %w", err)
	}

	logger.Info().Msg("  ✓ Job Distributor image built successfully")
	return nil
}

// pullJobDistributorImage pulls the Job Distributor image from ECR
func pullJobDistributorImage(ctx context.Context, jdLocalImage, jdEcrImage string) error {
	logger := framework.L

	// Check if AWS profile exists
	cmd := exec.Command("aws", "configure", "list-profiles")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list AWS profiles: %w", err)
	}

	if !strings.Contains(string(output), AWS_PROFILE) {
		return fmt.Errorf("AWS profile '%s' not found. Please ensure you have the correct AWS profile configured", AWS_PROFILE)
	}

	// Login to AWS SSO
	logger.Info().Msgf("AWS SSO Login for profile %s...", AWS_PROFILE)
	cmd = exec.CommandContext(ctx, "aws", "sso", "login", "--profile", AWS_PROFILE)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("AWS SSO login failed: %w", err)
	}
	logger.Info().Msgf("  ✓ AWS SSO login successful for profile %s", AWS_PROFILE)

	// Get ECR login password
	logger.Info().Msg("🔍 Pulling Job Distributor image from ECR...")
	cmd = exec.CommandContext(ctx, "aws", "ecr", "get-login-password", "--region", "us-west-2", "--profile", AWS_PROFILE)
	password, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get ECR login password: %w", err)
	}

	// Login to ECR
	ecrHostname := strings.Split(jdEcrImage, "/")[0]
	dockerLoginCmd := exec.CommandContext(ctx, "docker", "login", "--username", "AWS", "--password-stdin", ecrHostname)
	dockerLoginCmd.Stdin = bytes.NewBuffer(password)
	dockerLoginCmd.Stdout = os.Stdout
	dockerLoginCmd.Stderr = os.Stderr
	if err := dockerLoginCmd.Run(); err != nil {
		return fmt.Errorf("Docker login to ECR failed: %w", err)
	}
	logger.Info().Msg("  ✓ Docker login to ECR successful")

	// Pull image
	logger.Info().Msg("Pulling Job Distributor image from ECR...")
	cmd = exec.CommandContext(ctx, "docker", "pull", jdEcrImage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull Job Distributor image: %w", err)
	}

	// Tag image
	cmd = exec.CommandContext(ctx, "docker", "tag", jdEcrImage, jdLocalImage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to tag Job Distributor image: %w", err)
	}

	logger.Info().Msg("  ✓ Job Distributor image pulled successfully")
	return nil
}

// checkCRECLI checks if the CRE CLI is installed
func checkCRECLI(ctx context.Context) error {
	logger := framework.L

	// Check for CRE CLI
	osType := runtime.GOOS
	archType := runtime.GOARCH

	creBinaryName := fmt.Sprintf("cre_v%s_%s_%s", CRE_CLI_VERSION, osType, archType)
	if isCommandAvailable(creBinaryName) || isCommandAvailable("cre") {
		logger.Info().Msg("  ✓ CRE CLI is already installed")
		return nil
	}

	// CRE CLI not found
	logger.Info().Msg("  ✗ CRE CLI is not installed")
	logger.Info().Msg("Would you like to download and install the CRE CLI now? (y/n)")

	var input string
	fmt.Scanln(&input)

	if strings.ToLower(input) != "y" {
		logger.Warn().Msg("  ! You will need to install CRE CLI manually")
		return nil
	}

	// Check for GitHub CLI
	if !isCommandAvailable("gh") {
		return fmt.Errorf("GitHub CLI is not installed. Please install GitHub CLI or download CRE CLI manually from https://github.com/smartcontractkit/dev-platform/releases/tag/v%s", CRE_CLI_VERSION)
	}

	// Download CRE CLI
	logger.Info().Msgf("Downloading CRE CLI v%s for %s_%s...", CRE_CLI_VERSION, osType, archType)
	archivePattern := fmt.Sprintf("*%s_%s.tar.gz", osType, archType)
	cmd := exec.CommandContext(ctx, "gh", "release", "download", fmt.Sprintf("v%s", CRE_CLI_VERSION), "--repo", "smartcontractkit/dev-platform", "--pattern", archivePattern)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to download CRE CLI: %w", err)
	}

	// Extract archive
	archiveName := fmt.Sprintf("cre_v%s_%s_%s.tar.gz", CRE_CLI_VERSION, osType, archType)
	logger.Info().Msg("Extracting CRE CLI...")
	cmd = exec.CommandContext(ctx, "tar", "-xf", archiveName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract CRE CLI: %w", err)
	}

	// Remove archive
	if err := os.Remove(archiveName); err != nil {
		logger.Warn().Msgf("Failed to remove %s. Please remove it manually.", archiveName)
	}

	// Remove quarantine attribute on macOS
	if osType == "darwin" {
		cmd = exec.CommandContext(ctx, "xattr", "-d", "com.apple.quarantine", creBinaryName)
		_ = cmd.Run() // Ignore errors
	}

	// Make executable
	if err := os.Chmod(creBinaryName, 0755); err != nil {
		return fmt.Errorf("failed to make CRE CLI executable: %w", err)
	}

	// Create symlink
	if err := os.Symlink(creBinaryName, "cre"); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	logger.Info().Msgf("  ✓ CRE CLI installed to %s/cre", currentDir)
	logger.Warn().Msgf("  ! Add this directory to your PATH or move the CRE binary to a directory in your PATH")
	logger.Info().Msgf("   You can run: export PATH=\"%s:$PATH\"", currentDir)

	return nil
}

/*
// checkCapabilityBinaries checks if capability binaries are available
func checkCapabilityBinaries(ctx context.Context, extraCapabilities ExtraCapabilitiesConfig) error {
	logger := framework.L

	logger.Info().Msg("\n🔍 Checking for capability binaries...")

	// Check if capability binaries exist
	if extraCapabilities.CronCapabilityBinaryPath == "" {
		logger.Warn().Msg("  ! Cron capability binary path not set in configuration")
	} else if _, err := os.Stat(extraCapabilities.CronCapabilityBinaryPath); os.IsNotExist(err) {
		logger.Warn().Msgf("  ! Cron capability binary not found at %s", extraCapabilities.CronCapabilityBinaryPath)
	} else {
		logger.Info().Msgf("  ✓ Cron capability binary found at %s", extraCapabilities.CronCapabilityBinaryPath)
	}

	if extraCapabilities.LogEventTriggerBinaryPath == "" {
		logger.Warn().Msg("  ! Log event trigger capability binary path not set in configuration")
	} else if _, err := os.Stat(extraCapabilities.LogEventTriggerBinaryPath); os.IsNotExist(err) {
		logger.Warn().Msgf("  ! Log event trigger capability binary not found at %s", extraCapabilities.LogEventTriggerBinaryPath)
	} else {
		logger.Info().Msgf("  ✓ Log event trigger capability binary found at %s", extraCapabilities.LogEventTriggerBinaryPath)
	}

	if extraCapabilities.ReadContractBinaryPath == "" {
		logger.Warn().Msg("  ! Read contract capability binary path not set in configuration")
	} else if _, err := os.Stat(extraCapabilities.ReadContractBinaryPath); os.IsNotExist(err) {
		logger.Warn().Msgf("  ! Read contract capability binary not found at %s", extraCapabilities.ReadContractBinaryPath)
	} else {
		logger.Info().Msgf("  ✓ Read contract capability binary found at %s", extraCapabilities.ReadContractBinaryPath)
	}

	logger.Warn().Msg("  ! Some capabilities like cron, log-event-trigger, or read-contract might not be embedded in your Chainlink image")
	logger.Warn().Msg("  ! If needed, download from https://github.com/smartcontractkit/capabilities/releases/tag/v1.0.2-alpha")
	logger.Warn().Msg("  ! Or use: gh release download v1.0.2-alpha --repo smartcontractkit/capabilities --pattern 'amd64_cron'")

	return nil
}
*/
