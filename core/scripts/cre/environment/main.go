package main

import (
	"fmt"
	"os"

	capabilityregistry "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/capability-registry"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/environment"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/minio"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/mock"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/root"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/workflow"
)

func init() {
	root.RootCmd.AddCommand(environment.EnvironmentCmd)
	root.RootCmd.AddCommand(examples.ExamplesCmd)
	root.RootCmd.AddCommand(minio.MinioCommand)
	root.RootCmd.AddCommand(workflow.WorkflowCmd)
	root.RootCmd.AddCommand(mock.MockCommand)
	root.RootCmd.AddCommand(capabilityregistry.RegistryCmd)
}

func main() {
	if err := root.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
