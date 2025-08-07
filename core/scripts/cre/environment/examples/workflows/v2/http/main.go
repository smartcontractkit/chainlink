//go:build wasip1

package main

import (
	"encoding/json"
	"fmt"

	sdk "github.com/smartcontractkit/cre-sdk-go/cre"

	http "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

type None struct{}

func main() {
	wasm.NewRunner(func(configBytes []byte) (None, error) {
		return None{}, nil
	}).Run(RunSimpleHttpWorkflow)
}

func RunSimpleHttpWorkflow(wcx *sdk.Environment[None]) (sdk.Workflow[None], error) {
	workflows := sdk.Workflow[None]{
		sdk.Handler(
			http.Trigger(&http.Config{
				AuthorizedKeys: []*http.AuthorizedKey{
					{
						Type:      http.KeyType_KEY_TYPE_ECDSA,
						PublicKey: "0xC3Ad031A27E1A6C692cBdBafD85359b0BE1B15DD", // ALICE
					},
					{
						Type:      http.KeyType_KEY_TYPE_ECDSA,
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
func orderPizza(env *sdk.NodeEnvironment[None], nodeRuntime sdk.NodeRuntime, inputs map[string]interface{}, customer string) (string, error) {
	client := &http.Client{}

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
			StoreInCache:  true,
			TtlMs:         10000,
		}
	}

	resp, err := client.SendRequest(nodeRuntime, req).Await()
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

func onTrigger(env *sdk.Environment[None], runtime sdk.Runtime, trigger *http.Payload) (string, error) {
	env.Logger.Info("Hello! Workflow triggered.")

	inputMap := trigger.Input.AsMap()
	env.Logger.Info("Processing pizza order with inputs", "inputs", inputMap)

	customer := "default"
	if trigger.Key != nil && trigger.Key.PublicKey == "0x4b8d44a7a1302011fbc119407f8ce3baee6ea2ff" {
		customer = "Bob"
	}

	pizzaPromise := sdk.RunInNodeMode(env, runtime,
		func(env *sdk.NodeEnvironment[None], nodeRuntime sdk.NodeRuntime) (string, error) {
			return orderPizza(env, nodeRuntime, inputMap, customer)
		},
		sdk.ConsensusIdenticalAggregation[string](),
	)

	// Await the final, aggregated result.
	result, err := pizzaPromise.Await()
	if err != nil {
		return "", err
	}

	env.Logger.Info("Successfully processed pizza order", "result", result)
	return "", nil
}
