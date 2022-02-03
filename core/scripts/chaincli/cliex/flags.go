package cliex

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// BindStringFlag binds the given string flag.
func BindStringFlag(cmd *cobra.Command, key, flag, short, usage string) {
	cmd.Flags().StringP(flag, short, "", usage)
	if err := viper.BindPFlag(key, cmd.Flags().Lookup(flag)); err != nil {
		log.Fatal(err)
	}
}
