//go:build wasip1

package main

import (
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/sdk"
	"github.com/smartcontractkit/cre-sdk-go/sdk/wasm"
)

func RunSimpleCronWorkflow(_ *sdk.Environment[struct{}]) (sdk.Workflow[struct{}], error) {
	cfg := &cron.Config{
		Schedule: "*/3 * * * * *", // every 3 seconds
	}

	return sdk.Workflow[struct{}]{
		sdk.Handler(
			cron.Trigger(cfg),
			onTrigger,
		),
	}, nil
}

func onTrigger(env *sdk.Environment[struct{}], runtime sdk.Runtime, outputs *cron.Payload) (string, error) {
	env.Logger.Info("inside onTrigger handler")
	return "success!", nil
}

func main() {
	wasm.NewRunner(func(_ []byte) (struct{}, error) { return struct{}{}, nil }).Run(RunSimpleCronWorkflow)
}
