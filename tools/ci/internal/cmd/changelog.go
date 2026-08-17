package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/changelog"
)

func newChangelogCmd(stdout io.Writer) *cobra.Command {
	changelogCmd := &cobra.Command{
		Use:   "changelog",
		Short: "Changelog and release notes automation",
		Long:  "Commands to parse, format, group tags, and update CHANGELOG.md and PR release previews.",
	}

	changelogCmd.AddCommand(newFormatCmd(stdout))

	return changelogCmd
}

func newFormatCmd(stdout io.Writer) *cobra.Command {
	var (
		changelogPath string
		packageJSON   string
		githubOutput  bool
	)

	cmd := &cobra.Command{
		Use:   "format",
		Short: "Parse and group changeset preview notes in CHANGELOG.md",
		Example: `  go run ./tools/ci changelog format
  go run ./tools/ci changelog format --changelog=CHANGELOG.md --package-json=package.json --github-output`,
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := changelog.Format(changelogPath, packageJSON, githubOutput)
			if err != nil {
				return err
			}

			if _, err := fmt.Fprintf(stdout, "Formatted changelog for version %s\n", res.Version); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&changelogPath, "changelog", "CHANGELOG.md", "Path to CHANGELOG.md")
	flags.StringVar(&packageJSON, "package-json", "package.json", "Path to package.json")
	flags.BoolVar(&githubOutput, "github-output", true, "Write version and pr_body to $GITHUB_OUTPUT")

	return cmd
}
