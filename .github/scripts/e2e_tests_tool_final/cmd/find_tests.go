package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// findTestsCmd represents the find-tests command
var findTestsCmd = &cobra.Command{
	Use:   "find-tests [path]",
	Short: "Find all Go test functions in a directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		tests, err := extractTests(path)
		if err != nil {
			fmt.Println("Error extracting tests:", err)
			os.Exit(1)
		}
		for _, t := range tests {
			fmt.Printf("%+v", t)
		}
	},
}
