//go:build wasip1

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	http "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

type Config struct {
	AuthorizedKey string `json:"authorizedKey"`
	URL           string `json:"url"`
}

func main() {
	wasm.NewRunner(func(configBytes []byte) (Config, error) {
		var config Config
		if err := json.Unmarshal(configBytes, &config); err != nil {
			return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
		}
		return config, nil
	}).Run(RunSimpleHttpWorkflow)
}

func RunSimpleHttpWorkflow(config Config, _ *slog.Logger, _ cre.SecretsProvider) (sdk.Workflow[Config], error) {
	workflows := sdk.Workflow[Config]{
		sdk.Handler(
			http.Trigger(&http.Config{
				AuthorizedKeys: []*http.AuthorizedKey{
					{
						Type:      http.KeyType_KEY_TYPE_ECDSA_EVM,
						PublicKey: config.AuthorizedKey,
					},
				},
			}),
			onTrigger,
		),
	}
	return workflows, nil
}

type OrderResponse struct {
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func onTrigger(cfg Config, runtime sdk.Runtime, trigger *http.Payload) (string, error) {
	logger := runtime.Logger()
	logger.Info("Simple HTTP workflow triggered.")

	logger.Info("Processing order with inputs", "inputs", string(trigger.Input))

	orderPromise := sdk.RunInNodeMode(cfg, runtime,
		func(cfg Config, nodeRuntime sdk.NodeRuntime) (string, error) {
			client := &http.Client{}

			req := &http.Request{
				Url:    cfg.URL,
				Method: "POST",
				Body:   trigger.Input,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				TimeoutMs: 5000,
			}

			resp, err := client.SendRequest(nodeRuntime, req).Await()
			if err != nil {
				return "", fmt.Errorf("failed to post order: %w", err)
			}

			var orderResp OrderResponse
			if err := json.Unmarshal(resp.Body, &orderResp); err != nil {
				return "", fmt.Errorf("failed to unmarshal order response: %w", err)
			}

			if orderResp.Status == "success" {
				return fmt.Sprintf("Order placed successfully! Order ID: %s", orderResp.OrderID), nil
			}

			return "Order completed", nil
		},
		sdk.ConsensusIdenticalAggregation[string](),
	)

	result, err := orderPromise.Await()
	if err != nil {
		return "", err
	}

	logger.Info("Successfully processed order", "result", result)
	return result, nil
}
