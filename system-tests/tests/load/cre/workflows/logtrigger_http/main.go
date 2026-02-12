//go:build wasip1

// Package main implements a load-test workflow that triggers on EVM log events
// and then makes an HTTP POST to a configurable URL with the decoded event data.
// This exercises both the log-event-trigger and http-action capability paths.
package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"google.golang.org/protobuf/types/known/durationpb"
	"gopkg.in/yaml.v3"

	wfconfig "github.com/smartcontractkit/chainlink/system-tests/tests/load/cre/workflows/logtrigger_http/config"
)

func main() {
	wasm.NewRunner(func(b []byte) (wfconfig.Config, error) {
		cfg := wfconfig.Config{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return wfconfig.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
		}
		return cfg, nil
	}).Run(RunLogTriggerHTTPWorkflow)
}

func RunLogTriggerHTTPWorkflow(input wfconfig.Config, logger *slog.Logger, _ sdk.SecretsProvider) (sdk.Workflow[wfconfig.Config], error) {
	logger.Info("RunLogTriggerHTTPWorkflow: setting up log trigger with HTTP action")

	cfg := &evm.FilterLogTriggerRequest{
		Addresses: toByteSlices(input.Addresses),
		Topics: []*evm.TopicValues{
			{
				Values: toByteSlices(input.Topics[0].Values),
			},
		},
		Confidence: 1, // LATEST
	}

	return cre.Workflow[wfconfig.Config]{
		cre.Handler(
			evm.LogTrigger(input.ChainSelector, cfg),
			onTrigger,
		),
	}, nil
}

func onTrigger(cfg wfconfig.Config, runtime sdk.Runtime, outputs *evm.Log) (string, error) {
	logger := runtime.Logger()
	txHash := hex.EncodeToString(outputs.TxHash)
	logger.Info(fmt.Sprintf("onTrigger: txHash=%s logIndex=%d", txHash, outputs.Index))

	// NOTE: time.Sleep crashes wasip1 WASM (nanotime1 not implemented).
	// The 77-second delay must be implemented via the backend echo server instead.

	// Decode the log event data
	decodedMsg, err := decodeEventData(cfg.Abi, cfg.Event, outputs.Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode log data: %w", err)
	}
	logger.Info(fmt.Sprintf("onTrigger: decoded message: %s", decodedMsg))

	// Make an HTTP POST to the configured URL with the decoded event data
	httpResult, err := cre.RunInNodeMode(cfg, runtime,
		func(c wfconfig.Config, nodeRuntime cre.NodeRuntime) (string, error) {
			client := &http.Client{}
			body := fmt.Sprintf(`{"event":"%s","txHash":"%s","logIndex":%d}`, decodedMsg, txHash, outputs.Index)
			req := &http.Request{
				Url:    c.URL,
				Method: "POST",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body:    []byte(body),
				Timeout: &durationpb.Duration{Seconds: 9},
			}

			resp, err := client.SendRequest(nodeRuntime, req).Await()
			if err != nil {
				return "", fmt.Errorf("HTTP action failed: %w", err)
			}
			return string(resp.Body), nil
		},
		cre.ConsensusIdenticalAggregation[string](),
	).Await()

	if err != nil {
		return "", fmt.Errorf("HTTP action consensus failed: %w", err)
	}

	logger.Info(fmt.Sprintf("onTrigger: HTTP action result: %s", httpResult))
	return httpResult, nil
}

func decodeEventData(eventABI string, eventName string, data []byte) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(eventABI))
	if err != nil {
		return "", err
	}
	event := parsedABI.Events[eventName]
	values := make(map[string]interface{})
	if err := event.Inputs.UnpackIntoMap(values, data); err != nil {
		return "", err
	}

	var sb strings.Builder
	first := true
	for k, v := range values {
		if !first {
			sb.WriteString("; ")
		}
		sb.WriteString(fmt.Sprintf("%s:%v", k, v))
		first = false
	}
	return sb.String(), nil
}

func toByteSlices(hexStrings []string) [][]byte {
	result := make([][]byte, len(hexStrings))
	for i, s := range hexStrings {
		b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
		if err != nil {
			panic(fmt.Sprintf("failed to decode hex string %s: %v", s, err))
		}
		result[i] = b
	}
	return result
}
