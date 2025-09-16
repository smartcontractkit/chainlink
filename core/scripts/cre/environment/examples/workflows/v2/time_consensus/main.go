//go:build wasip1

package main

import (
	"errors"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"log/slog"
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

	if dontime1.IsZero() {
		err := errors.New("DON Time is zero; plugin has not started yet")
		runtime.Logger().Error(err.Error())
		return "", err
	}

	if !dontime2.After(dontime1) {
		return "", errors.New("DON time not increasing")
	}
	promise := cre.RunInNodeMode(cfg, runtime,
		func(cfg None, nodeRuntime cre.NodeRuntime) (int64, error) {
			return dontime1.UnixMilli(), nil
		},
		cre.ConsensusIdenticalAggregation[int64](),
	)

	_, err := promise.Await()
	if err != nil {
		runtime.Logger().Error("Failed to get identical consensus on DON Time")
		return "", err
	}

	runtime.Logger().Info("Verified consensus on DON Time")
	return "Verified consensus on DON Time", nil
}
