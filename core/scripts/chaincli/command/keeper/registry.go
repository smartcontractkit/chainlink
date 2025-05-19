package keeper

import (
	"log"

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

		verify, err := cmd.Flags().GetBool("verify")
		if err != nil {
			log.Fatal("failed to get verify flag: ", err)
		}

		hdlr.DeployRegistry(cmd.Context(), verify)
	},
}

var verifyRegistryCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify keeper registry",
	Long:  `This command verifys a keeper registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)

		hdlr.VerifyContract(args...)
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
	deployRegistryCmd.Flags().BoolP("verify", "v", false, "Specify if contracts should be verified on Etherscan")
	registryCmd.AddCommand(deployRegistryCmd)
	registryCmd.AddCommand(verifyRegistryCmd)
	registryCmd.AddCommand(updateRegistryCmd)
	registryCmd.AddCommand(withdrawFromRegistryCmd)
}
