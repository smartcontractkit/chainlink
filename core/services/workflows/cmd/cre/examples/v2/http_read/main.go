//go:build wasip1

package main

import (
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
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
	logger := runtime.Logger()
	logger.Info("onTrigger called")

	httpClient := &http.Client{}
	ret, err := http.SendRequest(config, runtime, httpClient, func(_ struct{}, _ *slog.Logger, sendRequester *http.SendRequester) (string, error) {
		resp, err := sendRequester.SendRequest(&http.Request{
			Method:  "GET",
			Url:     "https://dummyjson.com/test",
			Headers: map[string]string{"Content-Type": "application/json"},
		}).Await()
		return string(resp.Body), err
	}, cre.ConsensusIdenticalAggregation[string]()).Await()

	if err != nil {
		logger.Error("Error in RunInNodeMode", "err", err)
	} else {
		logger.Info("Successfully aggregated HTTP responses", "aggregatedResponse", ret)
	}
	return ret, err
}

func main() {
	wasm.NewRunner(func(_ []byte) (struct{}, error) { return struct{}{}, nil }).Run(RunSimpleCronWorkflow)
}
