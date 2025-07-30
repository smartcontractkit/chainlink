package download

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	keystonecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
)

var (
	capabilitiesVersion   string
	capabilityNames       []string
	outputDir             string
	ghReadTokenEnvVarName string
)

var downloadCapabilitiesCmd = &cobra.Command{
	Use:   "capabilities",
	Short: "Download capability binaries",
	Long:  `Download capability binaries from GitHub releases`,
	RunE: func(cmd *cobra.Command, args []string) error {
		githubToken, err := ghToken()
		if err != nil {
			return err
		}

		return downloadCapabilities(githubToken, capabilitiesVersion, capabilityNames)
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

	downloadCapabilitiesCmd.Flags().StringSliceVar(&capabilityNames, "names", []string{}, "Names of the capabilities to download (requires GITHUB_READ_TOKEN)")
	downloadCapabilitiesCmd.Flags().StringVar(&capabilitiesVersion, "version", "", "Version of the capabilities to download (requires GITHUB_READ_TOKEN)")

	DownloadCmd.AddCommand(downloadCapabilitiesCmd)
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

func downloadCapabilities(githubToken, version string, names []string) error {
	if len(names) == 0 {
		return errors.New("names flag is required")
	}
	if version == "" {
		return errors.New("version flag is required")
	}

	for _, name := range names {
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
	}

	return nil
}
