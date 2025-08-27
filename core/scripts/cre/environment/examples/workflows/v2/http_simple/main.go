//go:build wasip1

package main

import (
	"encoding/json"
	"fmt"

	http "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
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

func RunSimpleHttpWorkflow(wcx *sdk.Environment[Config]) (sdk.Workflow[Config], error) {
	config := wcx.Config

	workflows := sdk.Workflow[Config]{
		sdk.Handler(
			http.Trigger(&http.Config{
				AuthorizedKeys: []*http.AuthorizedKey{
					{
						Type:      http.KeyType_KEY_TYPE_ECDSA,
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

func onTrigger(env *sdk.Environment[Config], runtime sdk.Runtime, trigger *http.Payload) (string, error) {
	env.Logger.Info("Simple HTTP workflow triggered.")

	inputMap := trigger.Input.AsMap()
	env.Logger.Info("Processing order with inputs", "inputs", inputMap)

	orderPromise := sdk.RunInNodeMode(env, runtime,
		func(env *sdk.NodeEnvironment[Config], nodeRuntime sdk.NodeRuntime) (string, error) {
			client := &http.Client{}

			requestBody, err := json.Marshal(inputMap)
			if err != nil {
				return "", fmt.Errorf("failed to marshal order request: %w", err)
			}

			req := &http.Request{
				Url:    env.Config.URL,
				Method: "POST",
				Body:   requestBody,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
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

	env.Logger.Info("Successfully processed order", "result", result)
	return result, nil
}
