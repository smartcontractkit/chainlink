//go:build wasip1

package main

import (
	croncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/v2"
)

func RunSimpleCronWorkflow(_ *sdk.WorkflowContext[struct{}]) (sdk.Workflow[struct{}], error) {
	cron := &croncap.Cron{}
	cfg := &croncap.Config{
		Schedule: "*/3 * * * * *", // every 3 seconds
	}

	return sdk.Workflow[struct{}]{
		sdk.On(
			cron.Trigger(cfg),
			onTrigger,
		),
	}, nil
}

func onTrigger(_ *sdk.WorkflowContext[struct{}], runtime sdk.Runtime, outputs *croncap.Payload) (string, error) {
	return "ping", nil
}

func main() {
	wasm.NewRunner(func(_ []byte) (struct{}, error) { return struct{}{}, nil }).Run(RunSimpleCronWorkflow)
}
