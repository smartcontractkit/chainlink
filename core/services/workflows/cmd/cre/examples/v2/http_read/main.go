//go:build wasip1

package main

import (
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
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
	env.Logger.Info("onTrigger called")
	ret, err := sdk.RunInNodeMode(env, runtime, func(env *sdk.NodeEnvironment[struct{}], nrt sdk.NodeRuntime) (string, error) {
		httpClient := http.Client{}
		resp, err := httpClient.SendRequest(nrt, &http.Request{
			Method:  "GET",
			Url:     "https://dummyjson.com/test",
			Headers: map[string]string{"Content-Type": "application/json"},
		}).Await()
		return string(resp.Body), err
	}, sdk.ConsensusIdenticalAggregation[string]()).Await()

	if err != nil {
		env.Logger.Error("Error in RunInNodeMode", "err", err)
	} else {
		env.Logger.Info("Successfully aggregated HTTP responses", "aggregatedResponse", ret)
	}
	return ret, err
}

func main() {
	wasm.NewRunner(func(_ []byte) (struct{}, error) { return struct{}{}, nil }).Run(RunSimpleCronWorkflow)
}
