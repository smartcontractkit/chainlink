package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var envresolveCmd = &cobra.Command{
	Use:   "envresolve",
	Short: "Resolve environment variables in a string",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintln(cmd.OutOrStdout(), "Error: No input provided")
			return
		}
		input := args[0]

		fmt.Fprintln(cmd.OutOrStdout(), mustResolveEnvPlaceholder(input))
	},
}
