package keeper

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

// verifyCmd represents the command to run the service
var verifyCmd = &cobra.Command{
	Use:   "verify-config",
	Short: "Verify configs",
	Long:  `This command verifies configs for automation contracts.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewBaseHandler(cfg)
		hdlr.Verify()
	},
}
