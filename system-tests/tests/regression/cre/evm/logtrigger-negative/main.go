//go:build wasip1

package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"

	logtrigger_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/evm/logtrigger-negative/config"
)

func main() {
	wasm.NewRunner(func(b []byte) (logtrigger_negative_config.Config, error) {
		cfg := logtrigger_negative_config.Config{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return logtrigger_negative_config.Config{}, fmt.Errorf("error unmarshalling config: %w", err)
		}
		return cfg, nil
	}).Run(RunLogTriggerNegativeWorkflow)
}

func RunLogTriggerNegativeWorkflow(input logtrigger_negative_config.Config, logger *slog.Logger, secretsProvider sdk.SecretsProvider) (sdk.Workflow[logtrigger_negative_config.Config], error) {
	logger.Info("Trigger RunLogTriggerNegativeWorkflow called")

	// create a LogTrigger with an EOA address (not a contract), this should fail during registration in log poller first and bubble up in log trigger
	invalidAddress := input.InvalidAddress

	// dummy event signature (topic 0) - any valid topic will work (error is in registration, not execution time)
	dummyEventSig := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" // Transfer event signature

	cfg := &evm.FilterLogTriggerRequest{
		Addresses: [][]byte{toBytes(invalidAddress)},
		Topics: []*evm.TopicValues{
			{
				Values: [][]byte{toBytes(dummyEventSig)},
			},
		},
		Confidence: 1, // LATEST
	}

	logger.Info(fmt.Sprintf(
		"FilterLogTriggerRequest content: Addresses=%v Topics=%+v Confidence=%d",
		formatAddresses(cfg.Addresses),
		formatTopics(cfg.Topics),
		cfg.Confidence,
	))

	// should fail during workflow registration because the address is not a contract
	return cre.Workflow[logtrigger_negative_config.Config]{
		cre.Handler(
			evm.LogTrigger(input.ChainSelector, cfg),
			onTrigger,
		),
	}, nil
}

func formatAddresses(addresses [][]byte) []string {
	result := make([]string, len(addresses))
	for i, addr := range addresses {
		result[i] = "0x" + hex.EncodeToString(addr)
	}
	return result
}

func formatTopics(topics []*evm.TopicValues) []string {
	result := make([]string, len(topics))
	for i, topic := range topics {
		vals := make([]string, len(topic.Values))
		for j, v := range topic.Values {
			vals[j] = "0x" + hex.EncodeToString(v)
		}
		result[i] = fmt.Sprintf("values: [%s]", strings.Join(vals, ", "))
	}
	return result
}

func toBytes(hexStr string) []byte {
	// remove 0x prefix if present
	hexStr = strings.TrimPrefix(hexStr, "0x")
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(fmt.Sprintf("failed to decode hex string %s: %v", hexStr, err))
	}
	return b
}

func onTrigger(cfg logtrigger_negative_config.Config, runtime sdk.Runtime, outputs *evm.Log) (string, error) {
	// should never be called because the trigger registration should fail
	runtime.Logger().Error("onTrigger should not be called - trigger registration should have failed")
	return "", fmt.Errorf("onTrigger should not be called - trigger registration should have failed for EOA address")
}
