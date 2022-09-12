package command

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

// BootstrapNodeCmd launches a chainlink node with a bootstrap job
var BootstrapNodeCmd = &cobra.Command{
	Use:   "bootstrap [address]",
	Short: "Setup a bootstrap node.",
	Long:  `This commands launches a chainlink node inside the docker container and sets up the bootstrap job`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		baseHandler := handler.NewBaseHandler(cfg)
		baseHandler.StartBootstrapNode(cmd.Context(), args[0])
	},
}
