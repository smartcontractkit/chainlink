package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/changelog"
)

func newChangelogCmd(stdout io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changelog",
		Short: "Changelog preview and release formatting tools",
	}

	cmd.AddCommand(newChangelogFormatCmd(stdout))

	return cmd
}

func newChangelogFormatCmd(stdout io.Writer) *cobra.Command {
	var (
		changelogPath   string
		packageJSONPath string
		githubOutput    bool
	)

	cmd := &cobra.Command{
		Use:     "format",
		Short:   "Group changelog entries by tag for release preview PR description",
		Example: `  go run ./tools/ci changelog format --changelog=CHANGELOG.md --package-json=package.json --github-output=true`,
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := changelog.Format(changelogPath, packageJSONPath, githubOutput)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(stdout, "Formatted changelog for version %s\n", res.Version)
			return err
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&changelogPath, "changelog", "CHANGELOG.md", "Path to CHANGELOG.md")
	flags.StringVar(&packageJSONPath, "package-json", "package.json", "Path to package.json to read release version")
	flags.BoolVar(&githubOutput, "github-output", true, "Write version and pr_body to $GITHUB_OUTPUT")

	return cmd
}
