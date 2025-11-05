//go:build wasip1

package main

import (
	"errors"
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

	defaultStringPromise := cre.RunInNodeMode(cfg, runtime, func(config None, nodeRuntime cre.NodeRuntime) ([]byte, error) {
		return nil, nil
	}, cre.ConsensusIdenticalAggregation[[]byte]().WithDefault([]byte("stuff")))
	resultBytes, err := defaultStringPromise.Await()
	if err != nil {
		runtime.Logger().Warn("Consensus error on default string", "error", err)
		return "", err
	}

	runtime.Logger().Info("Result bytes are here", "result", resultBytes)

	if string(resultBytes) == "stuff" {
		runtime.Logger().Info("Successfully reached identical consensus on default value", "result", resultBytes)
	} else {
		runtime.Logger().Warn("Failed to reach consensus on default value")
		return "failed", errors.New("Failed to reach consensus on default value")
	}

	mathPromise := cre.RunInNodeMode(cfg, runtime, fetchData, cre.ConsensusMedianAggregation[int]())
	offchainValue, err := mathPromise.Await()
	if err != nil {
		runtime.Logger().Warn("Consensus error", "error", err)
		return "", err
	}
	runtime.Logger().Info("Successfully fetched offchain value and reached consensus", "result", offchainValue)
	runtime.Logger().Info("Successfully passed all consensus tests")

	return "success", nil
}

func fetchData(cfg None, nodeRuntime cre.NodeRuntime) (int, error) {

	randomValue := rand.Intn(10000)
	nodeRuntime.Logger().Info("Generate random value", "randomValue", randomValue)

	// Generate a random int64
	return randomValue, nil
}
