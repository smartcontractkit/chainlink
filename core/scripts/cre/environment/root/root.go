package root

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "go_run_main.go", // this thing displays only the first word of the command!!!
	Short: "CLI tool for the local CRE",
	Long:  `A CLI tool for the local CRE to create and manage environments`,
}
