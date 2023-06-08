package keeper

import (
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Keeper registry management",
	Long:  `This command provides an interface to manage keeper registry.`,
}

var deployRegistryCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy keeper registry",
	Long:  `This command deploys a new keeper registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)

		hdlr.DeployRegistry(cmd.Context())
	},
}

var updateRegistryCmd = &cobra.Command{
	Use:   "update",
	Short: "Update keeper registry",
	Long:  `This command updates existing keeper registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)

		hdlr.UpdateRegistry(cmd.Context())
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

func init() {
	registryCmd.AddCommand(deployRegistryCmd)
	registryCmd.AddCommand(updateRegistryCmd)
	registryCmd.AddCommand(withdrawFromRegistryCmd)
}
