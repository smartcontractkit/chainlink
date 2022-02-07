package keeper

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

var launchAndTestCmd = &cobra.Command{
	Use:   "launch-and-test",
	Short: "Launches keepers and starts performing",
	Long:  `This command launches chainlink nodes, keeper setup and starts performing upkeeps.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)
		hdlr.LaunchAndTest(cmd.Context())
	},
}
