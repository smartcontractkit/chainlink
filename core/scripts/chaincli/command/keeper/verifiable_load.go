package keeper

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

// verifiableLoad represents the command to get verifiable load testing details
var verifiableLoad = &cobra.Command{
	Use:   "verifiable-load",
	Short: "Print verifiable load testing details to console",
	Long:  `Print verifiable load testing details to console, including details of every active upkeep and total result`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)
		hdlr.GetVerifiableLoadStats(cmd.Context())
	},
}
