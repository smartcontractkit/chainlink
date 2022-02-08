package feed

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

// deployCmd represents the command to run the service.
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy price feed",
	Long:  `This command deploys price feeds.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewFeed(cfg)

		hdlr.DeployDerivedPriceFeed(cmd.Context())
	},
}
