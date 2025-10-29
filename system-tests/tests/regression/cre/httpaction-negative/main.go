//go:build wasip1

package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/httpaction-negative/config"
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
	}).Run(RunHTTPActionWorkflow)
}

func RunHTTPActionWorkflow(wfCfg config.Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	return cre.Workflow[config.Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onCronTrigger,
		),
	}, nil
}

func onCronTrigger(wfCfg config.Config, runtime cre.Runtime, payload *cron.Payload) (_ any, _ error) {
	logger := runtime.Logger()
	logger.Info("HTTP Action regression workflow triggered", "testCase", wfCfg.TestCase)

	return runCRUDFailureTest(wfCfg, runtime)
}

func runCRUDFailureTest(wfCfg config.Config, runtime cre.Runtime) (string, error) {
	logger := runtime.Logger()
	logger.Info("Running CRUD failure test using HTTP Action capability")

	// Use the configuration passed from the test to determine what kind of failure to test
	failurePromise := cre.RunInNodeMode(wfCfg, runtime,
		func(cfg config.Config, nodeRuntime cre.NodeRuntime) (string, error) {
			client := &http.Client{}

			// Build request based on configuration
			req := &http.Request{
				Url:     cfg.URL,
				Method:  cfg.Method,
				Headers: cfg.Headers,
				Timeout: &durationpb.Duration{
					Seconds: int64(cfg.TimeoutMs),
				},
			}

			// Set default method if not specified
			if req.Method == "" {
				req.Method = "GET"
			}

			// Set default timeout if not specified
			if req.Timeout.Seconds == 0 {
				req.Timeout.Seconds = 10 // Default to 10 seconds
			}

			// Add body if specified
			if cfg.Body != "" {
				req.Body = []byte(cfg.Body)
			}

			logger.Info("Testing HTTP Action with configuration",
				"url", req.Url,
				"method", req.Method,
				"timeout", req.Timeout,
				"hasBody", len(cfg.Body) > 0)

			resp, err := client.SendRequest(nodeRuntime, req).Await()
			if err != nil {
				// Expected failure - return specific error identifier based on the failure type
				failureType := determineFailureType(cfg)
				return fmt.Sprintf("HTTP Action failure test completed: %s", failureType), nil
			}
			// If we get here, the request unexpectedly succeeded
			return fmt.Sprintf("HTTP Action unexpected success with status: %d", resp.StatusCode), nil
		},
		cre.ConsensusIdenticalAggregation[string](),
	)

	result, err := failurePromise.Await()
	if err != nil {
		return "", fmt.Errorf("HTTP Action failure test error: %w", err)
	}

	logger.Info("HTTP Action failure test completed", "result", result)
	return result, nil
}

// determineFailureType maps the configuration to a specific failure type identifier
func determineFailureType(cfg config.Config) string {
	// Check URL patterns
	if cfg.URL == "not-a-valid-url" {
		return "invalid-url-format"
	}
	if strings.Contains(cfg.URL, "non-existing-domain") {
		return "non-existing-domain"
	}
	if strings.Contains(cfg.URL, "delay/10") && cfg.TimeoutMs == 1000 {
		return "request-timeout"
	}
	if len(cfg.URL) > 1000 { // Very long URL
		return "oversized-url"
	}

	// Check method
	if cfg.Method == "INVALID" {
		return "invalid-http-method"
	}

	// Check headers
	if cfg.Headers != nil {
		for header := range cfg.Headers {
			if strings.Contains(header, "\n") {
				return "invalid-headers"
			}
		}
	}

	// Check body
	if strings.Contains(cfg.Body, `{"invalid": json}`) {
		return "corrupt-json-body"
	}
	if len(cfg.Body) > 1024*1024 { // Large body
		return "oversized-request-body"
	}

	// Default fallback
	return "unknown-failure-type"
}
