//go:build wasip1

package main

import (
	"fmt"
	"log/slog"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/httpaction/config"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"google.golang.org/protobuf/types/known/durationpb"
	"gopkg.in/yaml.v3"
)

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		wfCfg := config.Config{}
		if err := yaml.Unmarshal(b, &wfCfg); err != nil {
			return config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
		}
		return wfCfg, nil
	}).Run(RunHTTPActionSuccessWorkflow)
}

func RunHTTPActionSuccessWorkflow(wfCfg config.Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	return cre.Workflow[config.Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onCronTrigger,
		),
	}, nil
}

func onCronTrigger(wfCfg config.Config, runtime cre.Runtime, payload *cron.Payload) (_ any, _ error) {
	logger := runtime.Logger()
	logger.Info(
		"HTTP Action workflow triggered",
		"testCase",
		wfCfg.TestCase,
		"method",
		wfCfg.Method,
		"url",
		wfCfg.URL,
	)

	return runCRUDSuccessTest(wfCfg, runtime)
}

func runCRUDSuccessTest(wfCfg config.Config, runtime cre.Runtime) (string, error) {
	logger := runtime.Logger()
	logger.Info(
		"Running HTTP Action capability",
		"testCase",
		wfCfg.TestCase,
		"method",
		wfCfg.Method,
		"url",
		wfCfg.URL,
	)

	// Execute CRUD operations using HTTP Action capability
	crudPromise := cre.RunInNodeMode(wfCfg, runtime,
		func(cfg config.Config, nodeRuntime cre.NodeRuntime) (string, error) {
			client := &http.Client{}

			req := &http.Request{
				Url:     cfg.URL,
				Method:  cfg.Method,
				Headers: map[string]string{"Content-Type": "application/json"},
				Body:    []byte(cfg.Body),
				Timeout: &durationpb.Duration{Seconds: 10},
			}

			logger.Info("Testing HTTP Action capability with configuration",
				"url", req.Url,
				"method", req.Method,
				"hasBody", len(cfg.Body) > 0)

			resp, err := client.SendRequest(nodeRuntime, req).Await()
			if err != nil {
				logger.Error(
					"Failed to complete HTTP Action request",
					"error",
					err,
					"url",
					req.Url,
					"method",
					req.Method,
				)
				return "", fmt.Errorf("HTTP Action %s request failed: %w", req.Method, err)
			}

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				logger.Error(
					"Failed response to HTTP Action request",
					"status",
					resp.StatusCode,
					"url",
					req.Url,
					"method",
					req.Method,
				)
				return "", fmt.Errorf("HTTP Action %s request failed with status: %d", req.Method, resp.StatusCode)
			}

			logger.Info(
				"HTTP Action completed",
				"url",
				req.Url,
				"method",
				req.Method,
				"status",
				resp.StatusCode,
				"body",
				string(resp.Body),
			)

			return fmt.Sprintf("HTTP Action CRUD success test completed: %s", cfg.TestCase), nil
		},
		cre.ConsensusIdenticalAggregation[string](),
	)

	result, err := crudPromise.Await()
	if err != nil {
		logger.Error("Failed to complete HTTP Action capability", "error", err)
		return "", fmt.Errorf("HTTP Action test failed: %w", err)
	}

	logger.Info("HTTP Action test completed", "result", result)
	return result, nil
}
