//go:build wasip1

package main

import (
	"time"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

type None struct{}

func main() {
	wasm.NewRunner(func(configBytes []byte) (None, error) {
		return None{}, nil
	}).Run(RunSimpleCronWorkflow)
}

func RunSimpleCronWorkflow(wcx *sdk.Environment[None]) (sdk.Workflow[None], error) {
	workflows := sdk.Workflow[None]{
		sdk.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onTrigger,
		),
	}
	return workflows, nil
}

func onTrigger(wcx *sdk.Environment[None], runtime sdk.Runtime, trigger *cron.Payload) (string, error) {
	donTime := time.Now()
	wcx.Logger.Info("Requested DON Time", "donTime", donTime)
	return "Requested DON Time", nil
}
