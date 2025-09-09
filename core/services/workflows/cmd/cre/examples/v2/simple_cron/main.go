//go:build wasip1

package main

import (
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

func RunSimpleCronWorkflow(_ struct{}, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[struct{}], error) {
	cfg := &cron.Config{
		Schedule: "*/3 * * * * *", // every 3 seconds
	}

	return cre.Workflow[struct{}]{
		cre.Handler(
			cron.Trigger(cfg),
			onTrigger,
		),
	}, nil
}

func onTrigger(config struct{}, runtime cre.Runtime, outputs *cron.Payload) (string, error) {
	runtime.Logger().Info("inside onTrigger handler")
	return "success!", nil
}

func main() {
	wasm.NewRunner(func(_ []byte) (struct{}, error) { return struct{}{}, nil }).Run(RunSimpleCronWorkflow)
}
