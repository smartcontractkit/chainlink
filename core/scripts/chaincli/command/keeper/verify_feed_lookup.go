package keeper

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

// jobCmd represents the command to run the service
var verifyMercuryCmd = &cobra.Command{
	Use:   "verify mercury",
	Short: "temporary method to verify mercury support",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)
		hdlr.VerifyFeedLookup(cmd.Context())
	},
}
