package root

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "cli",
	Short: "CLI tool for system tests",
	Long:  `A CLI tool for system tests that helps download dependencies and start envioronments`,
}
