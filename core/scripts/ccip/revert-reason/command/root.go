package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ccip-revert-reason",
	Short: "ChainLink CLI tool to resolve CCIP revert reasons",
	Long:  `ccip-revert-reason is a CLI for running the CCIP revert reason resolution commands.`,
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

	RootCmd.AddCommand(RevertReasonCmd)
	RevertReasonCmd.Flags().Bool("from-error", false, "Whether to decode an error string instead of transaction hash")
}
