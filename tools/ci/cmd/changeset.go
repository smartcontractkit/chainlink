package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/changeset"
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/ghaction"
)

func newChangesetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changeset",
		Short: "Commands for changeset validation and formatting",
	}

	cmd.AddCommand(newChangesetCheckTagsCmd())

	return cmd
}

func newChangesetCheckTagsCmd() *cobra.Command {
	var (
		filePath   string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "check-tags [file]",
		Short: "Check if at least one release tag exists in a changeset file",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			act := ghaction.NewAction(cmd.OutOrStdout())
			if len(args) > 0 {
				filePath = args[0]
			}
			if filePath == "" {
				filePath = act.GetInput("file")
			}
			if filePath == "" {
				filePath = act.Getenv("CHANGESET_FILE_PATH")
			}
			paths := strings.Fields(filePath)
			if len(paths) == 0 {
				return errors.New("no changeset file path provided")
			}
			if len(paths) > 1 {
				return fmt.Errorf("multiple changeset file paths provided (%d); please run once per file or pass a single path: %q", len(paths), filePath)
			}
			filePath = paths[0]

			res, err := changeset.CheckTags(filePath)
			if err != nil {
				return err
			}

			if jsonOutput {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(res)
			}

			for _, tag := range res.FoundTags {
				fmt.Fprintf(cmd.OutOrStdout(), "Found tag: %s in %s\n", tag, filePath)
			}

			if !res.HasTags {
				fmt.Fprintf(cmd.OutOrStdout(), "Error: No tags found in %s\n", filePath)
			}

			if os.Getenv("GITHUB_OUTPUT") != "" {
				act := ghaction.New(cmd.OutOrStdout(), "", "")
				if err := act.SetOutput("has_tags", strconv.FormatBool(res.HasTags)); err != nil {
					return err
				}
				if err := act.SetOutput("found_tags", strings.Join(res.FoundTags, ",")); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "Path to the changeset file (env: CHANGESET_FILE_PATH)")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output result in JSON format")

	return cmd
}
