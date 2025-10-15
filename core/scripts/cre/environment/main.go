package main

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/environment"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/minio"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/mock"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/root"
)

func init() {
	root.RootCmd.AddCommand(environment.EnvironmentCmd)
	root.RootCmd.AddCommand(examples.ExamplesCmd)
	root.RootCmd.AddCommand(minio.MinioCommand)
	root.RootCmd.AddCommand(mock.MockCommand)
}

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "shell" || os.Args[1] == "sh") {
		_ = os.Setenv("CTF_CONFIGS", "configs/workflow-don.toml") // Set default config for shell

		StartShell()
		return
	}
	if err := root.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
