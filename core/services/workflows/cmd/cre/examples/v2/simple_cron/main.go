//go:build wasip1

package main

import (
	croncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/v2"
)

func RunSimpleCronWorkflow(runner sdk.DonRunner) {
	cron := &croncap.Cron{}
	cfg := &croncap.Config{
		Schedule: "*/3 * * * * *", // every three seconds
	}

	runner.Run(&sdk.WorkflowArgs[sdk.DonRuntime]{
		Handlers: []sdk.Handler[sdk.DonRuntime]{
			sdk.NewDonHandler(
				cron.Trigger(cfg),
				onTrigger,
			),
		},
	})
}

func onTrigger(runtime sdk.DonRuntime, outputs *croncap.Payload) (string, error) {
	return "ping", nil
}

func main() {
	RunSimpleCronWorkflow(wasm.NewDonRunner())
}
