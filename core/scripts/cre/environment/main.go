package main

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/environment"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/root"
)

func init() {
	root.RootCmd.AddCommand(environment.EnvironmentCmd)
	root.RootCmd.AddCommand(examples.ExamplesCmd)
}

func main() {
	if err := root.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
