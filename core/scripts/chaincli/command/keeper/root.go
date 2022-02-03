package keeper

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/command"
)

// rootCmd represents the root keeper sub-command to manage keepers
var rootCmd = &cobra.Command{
	Use:   "keeper",
	Short: "Manage keepers",
	Long:  `This command represents a CLI interface to manage keepers.`,
}

func init() {
	command.RootCmd.AddCommand(rootCmd)
}
