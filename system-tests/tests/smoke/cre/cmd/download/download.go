package download

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	keystonecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
)

var (
	capabilityVersion     string
	capabilityName        string
	creCliVersion         string
	outputDir             string
	ghReadTokenEnvVarName string
)

var downloadCapabilityCmd = &cobra.Command{
	Use:   "capability",
	Short: "Download a capability binary",
	Long:  `Download a capability binary from GitHub releases`,
	RunE: func(cmd *cobra.Command, args []string) error {
		githubToken, err := ghToken()
		if err != nil {
			return err
		}

		return downloadCapability(githubToken, capabilityName, capabilityVersion)
	},
}

var downloadCreCliCmd = &cobra.Command{
	Use:   "cre-cli",
	Short: "Download CRE CLI binary",
	Long:  `Download the CRE CLI binary from GitHub releases`,
	RunE: func(cmd *cobra.Command, args []string) error {
		githubToken, err := ghToken()
		if err != nil {
			return err
		}

		return downloadCreCLI(githubToken, creCliVersion)
	},
}

var downloadAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Download all binaries",
	Long:  `Download both the cron capability and CRE CLI binaries`,
	RunE: func(cmd *cobra.Command, args []string) error {
		githubToken, err := ghToken()
		if err != nil {
			return err
		}

		fmt.Println("Downloading all binaries...")

		if err := downloadCapability(githubToken, capabilityName, capabilityVersion); err != nil {
			return err
		}

		return downloadCreCLI(githubToken, creCliVersion)
	},
}

var DownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download binaries",
	Long:  `Download binaries for capabilities and CRE CLI`,
}

func init() {
	DownloadCmd.PersistentFlags().StringVar(&outputDir, "output-dir", ".", "Directory to save the binaries (defaults to current directory)")
	DownloadCmd.PersistentFlags().StringVar(&ghReadTokenEnvVarName, "gh-token-env-var-name", "GITHUB_READ_TOKEN", "Name of the environment variable that contains the GitHub read token")

	downloadCapabilityCmd.Flags().StringVar(&capabilityName, "name", "", "Name of the capability to download (requires GITHUB_READ_TOKEN)")
	downloadCapabilityCmd.Flags().StringVar(&capabilityVersion, "version", "", "Version of the capability to download (requires GITHUB_READ_TOKEN)")
	downloadCreCliCmd.Flags().StringVar(&creCliVersion, "version", "", "Version of the CRE CLI to download (requires GITHUB_READ_TOKEN)")
	downloadAllCmd.Flags().StringVar(&capabilityName, "capability-name", "", "Name of the capability to download (requires GITHUB_READ_TOKEN)")
	downloadAllCmd.Flags().StringVar(&capabilityVersion, "capability-version", "", "Version of the capability to download (requires GITHUB_READ_TOKEN)")
	downloadAllCmd.Flags().StringVar(&creCliVersion, "cre-cli-version", "", "Version of the CRE CLI to download (requires GITHUB_READ_TOKEN)")

	DownloadCmd.AddCommand(downloadCapabilityCmd)
	DownloadCmd.AddCommand(downloadCreCliCmd)
	DownloadCmd.AddCommand(downloadAllCmd)
}

func moveFile(src, dstDir string) error {
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	filename := filepath.Base(src)
	dst := filepath.Join(dstDir, filename)

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if _, err = dstFile.ReadFrom(srcFile); err != nil {
		return err
	}

	if err := os.Chmod(dst, 0755); err != nil {
		return err
	}

	return os.Remove(src)
}

func ghToken() (string, error) {
	githubToken := os.Getenv(ghReadTokenEnvVarName)
	if githubToken == "" {
		return "", fmt.Errorf("%s environment variable is not set", ghReadTokenEnvVarName)
	}

	return githubToken, nil
}

func downloadCapability(githubToken, name, version string) error {
	if name == "" {
		return errors.New("name flag is required")
	}
	if version == "" {
		return errors.New("version flag is required")
	}

	fmt.Printf("Downloading %s capability binary version %s...\n", name, version)
	path, err := keystonecapabilities.DownloadCapabilityFromRelease(githubToken, version, name)
	if err != nil {
		return errors.Wrapf(err, "failed to download %s capability", name)
	}

	fmt.Printf("%s capability binary downloaded to: %s\n", name, path)

	if outputDir != "" && outputDir != "." {
		if err := moveFile(path, outputDir); err != nil {
			return fmt.Errorf("failed to move binary to output path: %w", err)
		}
		fmt.Printf("Moved binary to: %s\n", filepath.Join(outputDir, filepath.Base(path)))
	}

	return nil
}

func downloadCreCLI(githubToken, version string) error {
	if version == "" {
		return errors.New("version flag is required")
	}

	fmt.Printf("Downloading CRE CLI binary version %s...\n", version)
	path, err := libcrecli.DownloadAndInstallChainlinkCLI(githubToken, version)
	if err != nil {
		return fmt.Errorf("failed to download CRE CLI: %w", err)
	}

	fmt.Printf("CRE CLI binary downloaded to: %s\n", path)

	if outputDir != "" && outputDir != "." {
		if err := moveFile(path, outputDir); err != nil {
			return fmt.Errorf("failed to move binary to output path: %w", err)
		}
		fmt.Printf("Moved binary to: %s\n", filepath.Join(outputDir, filepath.Base(path)))
	}

	return nil
}
