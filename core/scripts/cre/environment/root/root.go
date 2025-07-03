package root

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "go run .",
	Short: "CLI tool for the local CRE",
	Long:  `A CLI tool for the local CRE to create and manage environments`,
}
