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

type OrderRequest struct {
	Customer string   `json:"customer"`
	Toppings []string `json:"toppings"`
	Dedup    bool     `json:"dedupe"`
}

// orderPizza posts a pizza order to the orders endpoint
func orderPizza(sendReqester *http.SendRequester, inputs []byte, customer string) (string, error) {
	var orderRequest OrderRequest
	if err := json.Unmarshal(inputs, &orderRequest); err != nil {
		return "", fmt.Errorf("failed to unmarshal order request: %w", err)
	}
	// this demonstrates that workflows can have custom logic based on the identity that invoked HTTP trigger
	// see `onTrigger()` function for how customer can be set based on the authorized key
	if customer == "Bob" {
		orderRequest.Toppings = []string{"pineapples"}
	}

	req := &http.Request{
		Url:    "http://host.docker.internal:2999/orders",
		Method: "POST",
		Body:   inputs,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	if orderRequest.Dedup {
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

func onTrigger(config None, runtime cre.Runtime, trigger *http.Payload) (string, error) {
	logger := runtime.Logger()
	logger.Info("Hello! Workflow triggered.")

	logger.Info("Processing pizza order with inputs", "inputs", string(trigger.Input))

	customer := "default"
    // this demonstrates that workflows can have custom logic based on the identity that invoked HTTP trigger
	if trigger.Key != nil && trigger.Key.PublicKey == "0x4b8d44a7a1302011fbc119407f8ce3baee6ea2ff" {
		customer = "Bob"
	}

	client := &http.Client{}
	pizzaPromise := http.SendRequest(config, runtime, client, func(_ None, logger *slog.Logger, sendRequester *http.SendRequester) (string, error) {
		return orderPizza(sendRequester, trigger.Input, customer)
	}, cre.ConsensusIdenticalAggregation[string]())

	// Await the final, aggregated result.
	result, err := pizzaPromise.Await()
	if err != nil {
		return "", err
	}

	logger.Info("Successfully processed pizza order", "result", result)
	return "", nil
}
