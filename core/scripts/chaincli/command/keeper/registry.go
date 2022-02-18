package keeper

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

var updateRegistryCmd = &cobra.Command{
	Use:   "update",
	Short: "Update keeper registry",
	Long:  `This command updates existing keeper registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)

		hdlr.GetRegistry(cmd.Context())
	},
}
