//go:build wasip1

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	http "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"google.golang.org/protobuf/types/known/durationpb"

	// "github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

type WorkflowConfig struct {
	Schedule string `yaml:"schedule,omitempty"`
}

// func main() {
// 	wasm.NewRunner(func(configBytes []byte) (WorkflowConfig, error) {
// 		return WorkflowConfig{
// 			Schedule: "*/30 * * * * *",
// 		}, nil
// 	}).Run(RunSimpleCronWorkflow)
// }

func main() {
	wasm.NewRunner(func(configBytes []byte) (None, error) {
		return None{}, nil
	}).Run(RunSimpleHttpWorkflow)
}

type None struct{}

// func RunSimpleCronWorkflow(config WorkflowConfig, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[WorkflowConfig], error) {
// 	workflows := cre.Workflow[WorkflowConfig]{
// 		cre.Handler(
// 			cron.Trigger(&cron.Config{Schedule: config.Schedule}),
// 			onTrigger,
// 		),
// 	}
// 	return workflows, nil
// }

func RunSimpleHttpWorkflow(_ None, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[None], error) {
	workflows := cre.Workflow[None]{
		cre.Handler(
			http.Trigger(&http.Config{
				AuthorizedKeys: []*http.AuthorizedKey{
					{
						Type:      http.KeyType_KEY_TYPE_ECDSA_EVM,
						PublicKey: "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", // dev key
					},
				},
			}),
			onTrigger,
		),
	}
	return workflows, nil
}

type MietekResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Alive bool   `json:"alive"`
}

func onTrigger(cfg None, runtime cre.Runtime, trigger *http.Payload) (string, error) {
	logger := runtime.Logger()
	logger.Info("Simple HTTP workflow triggered.")

	logger.Info("Processing " + trigger.Key.PublicKey)

	var input []byte
	var err error

	b := map[string]string{
		"key": trigger.Key.PublicKey,
	}

	input, err = json.Marshal(b)
	if err != nil {
		return "", fmt.Errorf("failed to marshal input: %w", err)
	}

	mietekPromise := sdk.RunInNodeMode(cfg, runtime,
		func(cfg None, nodeRuntime sdk.NodeRuntime) (string, error) {
			client := &http.Client{}

			req := &http.Request{
				Url:    "https://mieteks-power.free.beeceptor.com/get-mietek",
				Method: "POST",
				Body:   input,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Timeout: &durationpb.Duration{
					Seconds: 10,
				},
			}

			resp, err := client.SendRequest(nodeRuntime, req).Await()
			if err != nil {
				return "", fmt.Errorf("failed to call mietek: %w", err)
			}

			var mietekResp MietekResponse
			if err := json.Unmarshal(resp.Body, &mietekResp); err != nil {
				return "", fmt.Errorf("failed to unmarshal Mietek response: %w", err)
			}

			if mietekResp.Alive {
				return fmt.Sprintf("Mietek is alive! ID: %s, Name: %s", mietekResp.ID, mietekResp.Name), nil
			}

			return "Mietek is not alive", nil
		},
		sdk.ConsensusIdenticalAggregation[string](),
	)

	result, err := mietekPromise.Await()
	if err != nil {
		return "", err
	}

	logger.Info("Successfully processed Mietek", "result", result)
	return result, nil
}

// func onTrigger(cfg WorkflowConfig, runtime sdk.Runtime, _ *cron.Payload) (string, error) {
// 	logger := runtime.Logger()
// 	logger.Info("Simple HTTP workflow triggered.")

// 	logger.Info("Processing Mietek")

// 	mietekPromise := sdk.RunInNodeMode(cfg, runtime,
// 		func(cfg WorkflowConfig, nodeRuntime sdk.NodeRuntime) (string, error) {
// 			client := &http.Client{}

// 			req := &http.Request{
// 				Url:    "https://mieteks-power.free.beeceptor.com/get-mietek",
// 				Method: "POST",
// 				// Body:   trigger.Input,
// 				Headers: map[string]string{
// 					"Content-Type": "application/json",
// 				},
// 				Timeout: &durationpb.Duration{
// 					Seconds: 10,
// 				},
// 			}

// 			resp, err := client.SendRequest(nodeRuntime, req).Await()
// 			if err != nil {
// 				return "", fmt.Errorf("failed to call mietek: %w", err)
// 			}

// 			var mietekResp MietekResponse
// 			if err := json.Unmarshal(resp.Body, &mietekResp); err != nil {
// 				return "", fmt.Errorf("failed to unmarshal Mietek response: %w", err)
// 			}

// 			if mietekResp.Alive {
// 				return fmt.Sprintf("Mietek is alive! ID: %s, Name: %s", mietekResp.ID, mietekResp.Name), nil
// 			}

// 			return "Mietek is not alive", nil
// 		},
// 		sdk.ConsensusIdenticalAggregation[string](),
// 	)

// 	result, err := mietekPromise.Await()
// 	if err != nil {
// 		return "", err
// 	}

// 	logger.Info("Successfully processed Mietek", "result", result)
// 	return result, nil
// }
