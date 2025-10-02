//go:build wasip1

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/http/config"
)

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		wfCfg := config.Config{}
		if err := yaml.Unmarshal(b, &wfCfg); err != nil {
			return config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
		}
		return wfCfg, nil
	}).Run(RunHTTPRegressionWorkflow)
}

func RunHTTPRegressionWorkflow(wfCfg config.Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	// Create HTTP trigger with potentially invalid configuration based on test case
	var triggerConfig *http.Config

	switch wfCfg.TestCase {
	case "invalid-key-type":
		// Use an invalid key type (non-existent enum value)
		triggerConfig = &http.Config{
			AuthorizedKeys: []*http.AuthorizedKey{
				{
					Type:      999, // Invalid key type
					PublicKey: wfCfg.AuthorizedKey,
				},
			},
		}
	case "invalid-public-key":
		// Use an invalid public key format
		triggerConfig = &http.Config{
			AuthorizedKeys: []*http.AuthorizedKey{
				{
					Type:      http.KeyType_KEY_TYPE_ECDSA_EVM,
					PublicKey: "invalid-public-key-format",
				},
			},
		}
	case "non-existing-public-key":
		// Use a non-existing but properly formatted public key
		triggerConfig = &http.Config{
			AuthorizedKeys: []*http.AuthorizedKey{
				{
					Type:      http.KeyType_KEY_TYPE_ECDSA_EVM,
					PublicKey: "0x0000000000000000000000000000000000000000",
				},
			},
		}
	default:
		// Default case with valid configuration (should not be used in regression tests)
		triggerConfig = &http.Config{
			AuthorizedKeys: []*http.AuthorizedKey{
				{
					Type:      http.KeyType_KEY_TYPE_ECDSA_EVM,
					PublicKey: wfCfg.AuthorizedKey,
				},
			},
		}
	}

	return cre.Workflow[config.Config]{
		cre.Handler(
			http.Trigger(triggerConfig),
			onHTTPTrigger,
		),
	}, nil
}

func onHTTPTrigger(wfCfg config.Config, runtime cre.Runtime, trigger *http.Payload) (string, error) {
	logger := runtime.Logger()
	logger.Info("HTTP regression workflow triggered", "testCase", wfCfg.TestCase)

	// This should not be reached if the trigger validation fails properly
	logger.Error("HTTP trigger should have failed but succeeded", "testCase", wfCfg.TestCase)

	// Try to make HTTP request to validate the workflow execution
	inputMap := trigger.Input.AsMap()
	logger.Info("Processing request with inputs", "inputs", inputMap)

	// Use http.SendRequest to make HTTP requests
	client := &http.Client{}
	requestPromise := http.SendRequest(wfCfg, runtime, client, func(_ config.Config, logger *slog.Logger, sendRequester *http.SendRequester) (string, error) {
		requestBody, err := json.Marshal(inputMap)
		if err != nil {
			return "", fmt.Errorf("failed to marshal request: %w", err)
		}

		req := &http.Request{
			Url:    wfCfg.URL,
			Method: "POST",
			Body:   requestBody,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			TimeoutMs: 5000,
		}

		resp, err := sendRequester.SendRequest(req).Await()
		if err != nil {
			return "", fmt.Errorf("failed to send request: %w", err)
		}

		return fmt.Sprintf("Request completed with status: %d", resp.StatusCode), nil
	}, cre.ConsensusIdenticalAggregation[string]())

	result, err := requestPromise.Await()
	if err != nil {
		return "", err
	}
	logger.Info("Successfully processed request", "result", result)
	return result, nil
}
