package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the current version of the ci tool.
var Version = "0.1.0"

func newVersionCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of the ci CLI tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOutput {
				payload := map[string]string{
					"version": Version,
				}
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(payload)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "ci version %s\n", Version)
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output version information in JSON format")

	return cmd
}
