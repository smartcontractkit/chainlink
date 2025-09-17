//go:build wasip1

package main

import (
	"errors"
	"log/slog"
	"time"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

type None struct{}

func main() {
	wasm.NewRunner(func(configBytes []byte) (None, error) {
		return None{}, nil
	}).Run(RunSimpleCronWorkflow)
}

func RunSimpleCronWorkflow(_ None, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[None], error) {
	workflows := cre.Workflow[None]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onTrigger,
		),
	}
	return workflows, nil
}

func onTrigger(cfg None, runtime cre.Runtime, _ *cron.Payload) (string, error) {
	dontime1 := runtime.Now()
	dontime2 := runtime.Now()

	if !dontime2.After(dontime1) {
		return "", errors.New("DON time not increasing")
	}
	promise := cre.RunInNodeMode(cfg, runtime,
		func(cfg None, nodeRuntime cre.NodeRuntime) (time.Time, error) {
			return dontime1, nil
		},
		cre.ConsensusIdenticalAggregation[time.Time](),
	)

	_, err := promise.Await()
	if err != nil {
		runtime.Logger().Error("Failed to get identical consensus on DON Time")
		return "", err
	}

	runtime.Logger().Info("Verified consensus on DON Time")
	return "Verified consensus on DON Time", nil
}
