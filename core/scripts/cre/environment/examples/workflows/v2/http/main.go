//go:build wasip1

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/cre"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

type None struct{}

func main() {
	wasm.NewRunner(func(configBytes []byte) (None, error) {
		return None{}, nil
	}).Run(RunSimpleHttpWorkflow)
}

func RunSimpleHttpWorkflow(_ None, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[None], error) {
	workflows := cre.Workflow[None]{
		cre.Handler(
			http.Trigger(&http.Config{
				AuthorizedKeys: []*http.AuthorizedKey{
					{
						Type:      http.KeyType_KEY_TYPE_ECDSA_EVM,
						PublicKey: "0xC3Ad031A27E1A6C692cBdBafD85359b0BE1B15DD", // ALICE
					},
					{
						Type:      http.KeyType_KEY_TYPE_ECDSA_EVM,
						PublicKey: "0x4b8d44A7A1302011fbc119407F8Ce3baee6Ea2FF", // BOB
					},
				},
			}),
			onTrigger,
		),
	}
	return workflows, nil
}

// OrderResponse represents the response from the orders endpoint
type OrderResponse struct {
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// orderPizza posts a pizza order to the orders endpoint
func orderPizza(sendReqester *http.SendRequester, inputs map[string]interface{}, customer string) (string, error) {
	if customer == "Bob" {
		inputs["toppings"] = []string{"pineapples"}
	}

	// Send the entire inputs as JSON body
	requestBody, err := json.Marshal(inputs)
	if err != nil {
		return "", fmt.Errorf("failed to marshal order request: %w", err)
	}

	req := &http.Request{
		Url:    "http://host.docker.internal:2999/orders",
		Method: "POST",
		Body:   requestBody,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	dedupe := extractBoolFromInput(inputs, "dedupe")
	if dedupe {
		req.CacheSettings = &http.CacheSettings{
			ReadFromCache: true,
			MaxAgeMs:      10000,
		}
	}

	resp, err := sendReqester.SendRequest(req).Await()
	if err != nil {
		return "", fmt.Errorf("failed to post pizza order: %w", err)
	}

	// Parse the JSON response
	var orderResp OrderResponse
	if err := json.Unmarshal(resp.Body, &orderResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	if orderResp.Status == "success" {
		return fmt.Sprintf("Pizza order placed successfully! Order ID: %s", orderResp.OrderID), nil
	}

	return "", nil
}

func extractBoolFromInput(input map[string]interface{}, key string) bool {
	if val, ok := input[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

func onTrigger(config None, runtime cre.Runtime, trigger *http.Payload) (string, error) {
	logger := runtime.Logger()
	logger.Info("Hello! Workflow triggered.")

	inputMap := trigger.Input.AsMap()
	logger.Info("Processing pizza order with inputs", "inputs", inputMap)

	customer := "default"
	if trigger.Key != nil && trigger.Key.PublicKey == "0x4b8d44a7a1302011fbc119407f8ce3baee6ea2ff" {
		customer = "Bob"
	}

	client := &http.Client{}
	pizzaPromise := http.SendRequest(config, runtime, client, func(_ None, logger *slog.Logger, sendRequester *http.SendRequester) (string, error) {
		return orderPizza(sendRequester, inputMap, customer)
	}, cre.ConsensusIdenticalAggregation[string]())

	// Await the final, aggregated result.
	result, err := pizzaPromise.Await()
	if err != nil {
		return "", err
	}

	logger.Info("Successfully processed pizza order", "result", result)
	return "", nil
}
