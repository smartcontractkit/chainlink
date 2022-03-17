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
	Long:  `This command deploys keepers (keeper registry + upkeeps).`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)
		hdlr.DeployKeepers(cmd.Context())
	},
}
