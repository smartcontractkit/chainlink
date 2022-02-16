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

var withdrawFromRegistryCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "cancel upkeeps and withdraw funds from registry",
	Long:  `This command will cancel all registered upkeeps and withdraw the funds left. args = Registry address`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)

		hdlr.Withdraw(cmd.Context(), args[0])
	},
}
