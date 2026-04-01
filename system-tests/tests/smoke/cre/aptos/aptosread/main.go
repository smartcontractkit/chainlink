//go:build wasip1

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/aptos"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	sdk "github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/aptos/aptosread/config"
)

var aptosCoinTypeTag = &aptos.TypeTag{
	Kind: aptos.TypeTagKind_TYPE_TAG_KIND_STRUCT,
	Value: &aptos.TypeTag_Struct{
		Struct: &aptos.StructTag{
			Address: []byte{0x1},
			Module:  "aptos_coin",
			Name:    "AptosCoin",
		},
	},
}

func main() {
	wasm.NewRunner(func(b []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("unmarshal config: %w", err)
		}
		return cfg, nil
	}).Run(RunReadWorkflow)
}

func RunReadWorkflow(cfg config.Config, logger *slog.Logger, secretsProvider sdk.SecretsProvider) (sdk.Workflow[config.Config], error) {
	return sdk.Workflow[config.Config]{
		sdk.Handler(
			cron.Trigger(&cron.Config{Schedule: "*/30 * * * * *"}),
			onAptosReadTrigger,
		),
	}, nil
}

func onAptosReadTrigger(cfg config.Config, runtime sdk.Runtime, payload *cron.Payload) (_ any, err error) {
	runtime.Logger().Info("onAptosReadTrigger called", "workflow", cfg.WorkflowName)
	defer func() {
		if r := recover(); r != nil {
			runtime.Logger().Info("Aptos read failed: panic in onAptosReadTrigger", "workflow", cfg.WorkflowName, "panic", fmt.Sprintf("%v", r))
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	client := aptos.Client{ChainSelector: cfg.ChainSelector}
	reply, err := client.View(runtime, &aptos.ViewRequest{
		Payload: &aptos.ViewPayload{
			Module: &aptos.ModuleID{
				Address: []byte{0x1},
				Name:    "coin",
			},
			Function: "name",
			ArgTypes: []*aptos.TypeTag{aptosCoinTypeTag},
		},
	}).Await()
	if err != nil {
		msg := fmt.Sprintf("Aptos read failed: View error: %v", err)
		runtime.Logger().Info(msg, "workflow", cfg.WorkflowName, "chainSelector", cfg.ChainSelector)
		return nil, fmt.Errorf("Aptos View(0x1::coin::name): %w", err)
	}
	if reply == nil {
		runtime.Logger().Info("Aptos read failed: View reply is nil", "workflow", cfg.WorkflowName)
		return nil, errors.New("View reply is nil")
	}
	if len(reply.Data) == 0 {
		runtime.Logger().Info("Aptos read failed: View reply data is empty", "workflow", cfg.WorkflowName)
		return nil, errors.New("View reply data is empty")
	}

	onchainValue, parseErr := parseSingleStringViewReply(reply.Data)
	if parseErr != nil {
		msg := fmt.Sprintf("Aptos read failed: cannot parse view reply data %q: %v", string(reply.Data), parseErr)
		runtime.Logger().Info(msg, "workflow", cfg.WorkflowName)
		return nil, fmt.Errorf("invalid Aptos view reply payload: %w", parseErr)
	}

	if onchainValue != cfg.ExpectedCoinName {
		msg := fmt.Sprintf("Aptos read failed: onchain value %q does not match expected %q", onchainValue, cfg.ExpectedCoinName)
		runtime.Logger().Info(msg, "workflow", cfg.WorkflowName)
		return nil, fmt.Errorf("onchain value %q does not match expected %q", onchainValue, cfg.ExpectedCoinName)
	}

	msg := "Aptos read consensus succeeded"
	runtime.Logger().Info(msg, "onchain_value", strings.TrimSpace(onchainValue), "workflow", cfg.WorkflowName)
	return nil, nil
}

func parseSingleStringViewReply(data []byte) (string, error) {
	var values []string
	if err := json.Unmarshal(data, &values); err != nil {
		return "", fmt.Errorf("decode json string array: %w", err)
	}
	if len(values) == 0 {
		return "", errors.New("empty json array")
	}
	return values[0], nil
}
