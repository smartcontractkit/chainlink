package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/command/keeper"
)

var configFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "chaincli",
	Short: "ChainLink CLI tool to manage products such as keeper, vrf, etc.",
	Long:  `chaincli is a CLI for running the product management commands, e.g. keepers deployment.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is .env)")
	_ = viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))

	RootCmd.AddCommand(keeper.RootCmd)
	RootCmd.AddCommand(BootstrapNodeCmd)
	RootCmd.AddCommand(RevertReasonCmd)
}
