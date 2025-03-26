package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	keystonecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
)

const (
	cronCapabilityAssetFile = "cron"
)

var (
	cronVersion           string
	creCliVersion         string
	outputDir             string
	ghReadTokenEnvVarName string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "download-cli",
		Short: "CLI tool for downloading binary dependencies",
		Long:  `A CLI tool that helps download binary dependencies for CRE testing`,
	}

	downloadCronCmd := &cobra.Command{
		Use:   "cron",
		Short: "Download CRE cron capability binary",
		Long:  `Download the cron capability binary from GitHub releases`,
		RunE: func(cmd *cobra.Command, args []string) error {
			githubToken := os.Getenv(ghReadTokenEnvVarName)
			if githubToken == "" {
				return fmt.Errorf("%s environment variable is not set", ghReadTokenEnvVarName)
			}

			fmt.Printf("Downloading cron capability binary version %s...\n", cronVersion)
			path, err := keystonecapabilities.DownloadCapabilityFromRelease(githubToken, cronVersion, cronCapabilityAssetFile)
			if err != nil {
				return fmt.Errorf("failed to download cron capability: %w", err)
			}

			fmt.Printf("Cron capability binary downloaded to: %s\n", path)

			// Move binary if output path is specified
			if outputDir != "" && outputDir != "." {
				if err := moveFile(path, outputDir); err != nil {
					return fmt.Errorf("failed to move binary to output path: %w", err)
				}
				fmt.Printf("Moved binary to: %s\n", filepath.Join(outputDir, filepath.Base(path)))
			}

			return nil
		},
	}

	downloadCreCliCmd := &cobra.Command{
		Use:   "cre-cli",
		Short: "Download CRE CLI binary",
		Long:  `Download the CRE CLI binary from GitHub releases`,
		RunE: func(cmd *cobra.Command, args []string) error {
			githubToken := os.Getenv(ghReadTokenEnvVarName)
			if githubToken == "" {
				return fmt.Errorf("%s environment variable is not set", ghReadTokenEnvVarName)
			}

			fmt.Printf("Downloading CRE CLI binary version %s...\n", creCliVersion)
			path, err := libcrecli.DownloadAndInstallChainlinkCLI(githubToken, creCliVersion)
			if err != nil {
				return fmt.Errorf("failed to download CRE CLI: %w", err)
			}

			fmt.Printf("CRE CLI binary downloaded to: %s\n", path)

			// Move binary if output path is specified
			if outputDir != "" && outputDir != "." {
				if err := moveFile(path, outputDir); err != nil {
					return fmt.Errorf("failed to move binary to output path: %w", err)
				}
				fmt.Printf("Moved binary to: %s\n", filepath.Join(outputDir, filepath.Base(path)))
			}

			return nil
		},
	}

	downloadAllCmd := &cobra.Command{
		Use:   "all",
		Short: "Download all binaries",
		Long:  `Download both the cron capability and CRE CLI binaries`,
		RunE: func(cmd *cobra.Command, args []string) error {
			githubToken := os.Getenv(ghReadTokenEnvVarName)
			if githubToken == "" {
				return fmt.Errorf("%s environment variable is not set", ghReadTokenEnvVarName)
			}

			fmt.Println("Downloading all binaries...")

			fmt.Printf("Downloading cron capability binary version %s...\n", cronVersion)
			cronPath, err := keystonecapabilities.DownloadCapabilityFromRelease(githubToken, cronVersion, cronCapabilityAssetFile)
			if err != nil {
				return fmt.Errorf("failed to download cron capability: %w", err)
			}
			fmt.Printf("Cron capability binary downloaded to: %s\n", cronPath)

			fmt.Printf("Downloading CRE CLI binary version %s...\n", creCliVersion)
			cliPath, err := libcrecli.DownloadAndInstallChainlinkCLI(githubToken, creCliVersion)
			if err != nil {
				return fmt.Errorf("failed to download CRE CLI: %w", err)
			}
			fmt.Printf("CRE CLI binary downloaded to: %s\n", cliPath)

			// Move binaries if output path is specified
			if outputDir != "" && outputDir != "." {
				if err := moveFile(cronPath, outputDir); err != nil {
					return fmt.Errorf("failed to move cron capability binary to output path: %w", err)
				}
				fmt.Printf("Moved cron capability binary to: %s\n", filepath.Join(outputDir, filepath.Base(cronPath)))

				if err := moveFile(cliPath, outputDir); err != nil {
					return fmt.Errorf("failed to move CRE CLI binary to output path: %w", err)
				}
				fmt.Printf("Moved CRE CLI binary to: %s\n", filepath.Join(outputDir, filepath.Base(cliPath)))
			}

			return nil
		},
	}

	// Adding flags for all commands
	rootCmd.PersistentFlags().StringVar(&outputDir, "output-dir", ".", "Directory to save the binaries (defaults to current directory)")
	rootCmd.PersistentFlags().StringVar(&ghReadTokenEnvVarName, "gh-token-env-var-name", "GITHUB_READ_TOKEN", "Name of the environment variable that contains the GitHub read token")

	downloadCronCmd.Flags().StringVar(&cronVersion, "version", "v1.0.2-alpha", "Version of the cron capability to download")
	downloadCreCliCmd.Flags().StringVar(&creCliVersion, "version", "v0.1.5", "Version of the CRE CLI to download")
	downloadAllCmd.Flags().StringVar(&cronVersion, "cron-version", "v1.0.2-alpha", "Version of the cron capability to download")
	downloadAllCmd.Flags().StringVar(&creCliVersion, "cre-cli-version", "v0.1.5", "Version of the CRE CLI to download")

	// Add commands to root
	rootCmd.AddCommand(downloadCronCmd)
	rootCmd.AddCommand(downloadCreCliCmd)
	rootCmd.AddCommand(downloadAllCmd)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func moveFile(src, dstDir string) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// Get the base filename
	filename := filepath.Base(src)
	dst := filepath.Join(dstDir, filename)

	// Create or truncate the destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Copy the contents
	if _, err = dstFile.ReadFrom(srcFile); err != nil {
		return err
	}

	// Make the destination file executable
	if err := os.Chmod(dst, 0755); err != nil {
		return err
	}

	// Uncomment if you want to delete the source after copying
	if err := os.Remove(src); err != nil {
		return err
	}

	return nil
}
