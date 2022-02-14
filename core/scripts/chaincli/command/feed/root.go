package feed

import (
	"github.com/spf13/cobra"
)

// RootCmd represents the root price feed sub-command to manage feeds.
var RootCmd = &cobra.Command{
	Use:   "feed",
	Short: "Manage price feeds",
	Long:  `This command represents a CLI interface to manage Chainlink Price Feeds.`,
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
