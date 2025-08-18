//go:build wasip1

package main

import (
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

func RunSimpleCronWorkflow(_ *cre.Environment[struct{}]) (cre.Workflow[struct{}], error) {
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

func onTrigger(env *cre.Environment[struct{}], runtime cre.Runtime, outputs *cron.Payload) (string, error) {
	env.Logger.Info("onTrigger called")
	ret, err := cre.RunInNodeMode(env, runtime, func(env *cre.NodeEnvironment[struct{}], nrt cre.NodeRuntime) (string, error) {
		httpClient := http.Client{}
		resp, err := httpClient.SendRequest(nrt, &http.Request{
			Method:  "GET",
			Url:     "https://dummyjson.com/test",
			Headers: map[string]string{"Content-Type": "application/json"},
		}).Await()
		return string(resp.Body), err
	}, cre.ConsensusIdenticalAggregation[string]()).Await()

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
