package cliex

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func BindStringFlag(cmd *cobra.Command, key, flag, short, usage string) {
	cmd.Flags().StringP(flag, short, "", usage)
	if err := viper.BindPFlag(key, cmd.Flags().Lookup(flag)); err != nil {
		logrus.Fatal(err)
	}
}
