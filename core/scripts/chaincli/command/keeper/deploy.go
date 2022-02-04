package keeper

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

// deployCmd represents the command to run the service
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy keepers",
	Long:  `This command deploys keepers (keeper registry + upkeeps). Accepts a filename argument else will look for .env`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New(args)
		hdlr := handler.NewKeeper(cfg)

		hdlr.DeployKeepers(cmd.Context())
	},
}

var updateRegistryCmd = &cobra.Command{
	Use:   "update",
	Short: "Update keeper registry",
	Long:  `This command updates existing keeper registry. Accepts a filename argument else will look for .env`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New(args)
		hdlr := handler.NewKeeper(cfg)

		hdlr.GetRegistry(cmd.Context())
	},
}
