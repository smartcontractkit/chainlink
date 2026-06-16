package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var (
		pluginFile string
		pluginName string
	)

	rootCmd := &cobra.Command{
		Use:   "pluginref",
		Short: "Resolve a plugin gitRef from YAML to a git rev-parse ref",
		RunE: func(cmd *cobra.Command, args []string) error {
			if pluginFile == "" {
				return errors.New("--plugin-file is required")
			}
			if pluginName == "" {
				return errors.New("--plugin is required")
			}

			gitRef, err := gitRefForPlugin(pluginFile, pluginName)
			if err != nil {
				return err
			}

			commitRef, err := commitRefForGitRef(gitRef)
			if err != nil {
				return err
			}

			fmt.Println(commitRef)
			return nil
		},
	}

	rootCmd.Flags().StringVar(&pluginFile, "plugin-file", "", "Path to plugins YAML file")
	rootCmd.Flags().StringVar(&pluginName, "plugin", "", "Plugin name (YAML key under plugins:)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
