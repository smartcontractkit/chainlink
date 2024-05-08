package main

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/test_workflow"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/wasm/sdk"
)

func main() {
	workflow, err := test_workflow.CreateWorkflow()
	if err != nil {
		panic(err)
	}

	runner := sdk.NewRunner()
	if err = runner.Run(workflow); err != nil {
		panic(err)
	}
}
