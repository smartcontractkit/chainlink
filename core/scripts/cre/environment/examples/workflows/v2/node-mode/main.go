//go:build wasip1

package main

import (
	"fmt"

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

func onTrigger(env *sdk.Environment[None], runtime sdk.Runtime, trigger *cron.Payload) (string, error) {
	mathPromise := sdk.RunInNodeMode(env, runtime, fetchData, sdk.ConsensusIdenticalAggregation[float64]())
	offchainValue, err := mathPromise.Await()
	if err != nil {
		return "", err
	}
	env.Logger.Info("Successfully fetched offchain value", "result", offchainValue)
	return fmt.Sprintf("value: %f", offchainValue), nil
}

func fetchData(env *sdk.NodeEnvironment[None], nodeRuntime sdk.NodeRuntime) (float64, error) {
	// pretend we're fetching some node-mode data
	return 420.69, nil
}
