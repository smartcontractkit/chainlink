//go:build wasip1

package main

import (
	"log/slog"
	"math/rand"

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
	runtime.Logger().Info("Triggered fetch of value")

	mathPromise := cre.RunInNodeMode(cfg, runtime, fetchData, cre.ConsensusMedianAggregation[int]())
	offchainValue, err := mathPromise.Await()
	if err != nil {
		runtime.Logger().Warn("Consensus error", "error", err)
		return "", err
	}
	runtime.Logger().Info("Successfully fetched offchain value and reached consensus", "result", offchainValue)

	return "success", nil
}

func fetchData(cfg None, nodeRuntime cre.NodeRuntime) (int, error) {

	randomValue := rand.Intn(10000)
	nodeRuntime.Logger().Info("Generate random value", "randomValue", randomValue)

	// Generate a random int64
	return randomValue, nil
}
