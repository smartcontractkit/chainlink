//go:build wasip1

package main

import (
	"encoding/binary"
	"fmt"
	"log/slog"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/sollogtrigger/config"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/solana"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"
)

const testEventDiscriminatorLen = 8

type testEvent struct {
	StrVal   string
	U64Value uint64
}

func decodeTestEvent(data []byte) (testEvent, error) {
	if len(data) < testEventDiscriminatorLen+4+8 {
		return testEvent{}, fmt.Errorf("event data too short: %d bytes", len(data))
	}

	payload := data[testEventDiscriminatorLen:]
	strLen := binary.LittleEndian.Uint32(payload[:4])
	strEnd := 4 + int(strLen)
	if strLen == 0 || strEnd+8 > len(payload) {
		return testEvent{}, fmt.Errorf("invalid TestEvent payload length")
	}

	return testEvent{
		StrVal:   string(payload[4:strEnd]),
		U64Value: binary.LittleEndian.Uint64(payload[strEnd : strEnd+8]),
	}, nil
}

func formatTestEvent(evt testEvent) string {
	return fmt.Sprintf("str_val=%s u64_value=%d", evt.StrVal, evt.U64Value)
}

func validateTestEvent(cfg config.Config, evt testEvent) error {
	if evt.StrVal != cfg.ExpectedStrVal {
		return fmt.Errorf("unexpected str_val: got %q want %q", evt.StrVal, cfg.ExpectedStrVal)
	}
	if evt.U64Value != cfg.ExpectedU64Value {
		return fmt.Errorf("unexpected u64_value: got %d want %d", evt.U64Value, cfg.ExpectedU64Value)
	}
	return nil
}

func RunSolLogTriggerWorkflow(cfg config.Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	logger.Info("RunSolLogTriggerWorkflow called")

	if cfg.ContractIdlJSON == "" {
		return nil, fmt.Errorf("contract_idl_json is required in workflow config")
	}
	eventIdlJson := []byte(cfg.ContractIdlJSON)

	expectedValueBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(expectedValueBytes, cfg.ExpectedU64Value)

	filterLogTriggerRequest := &solana.FilterLogTriggerRequest{
		Name:            "test-event-filter",
		Address:         cfg.LogReadTestProgramID[:],
		EventName:       "TestEvent",
		ContractIdlJson: eventIdlJson,
		Subkeys: []*solana.SubkeyConfig{
			{Path: []string{"U64Value"}, Comparers: []*solana.ValueComparator{
				{Value: expectedValueBytes, Operator: solana.ComparisonOperator_COMPARISON_OPERATOR_EQ},
			}},
		},
	}

	filterLogTriggerRequestCPI := &solana.FilterLogTriggerRequest{
		Name:            "test-cpi-event-filter",
		Address:         cfg.LogReadTestProgramID[:],
		EventName:       "TestEvent",
		ContractIdlJson: eventIdlJson,
		CpiFilterConfig: &solana.CPIFilterConfig{
			DestAddress: cfg.LogReadTestProgramID[:],
			MethodName:  []byte("anchor:event"),
		},
	}

	chainSelector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	if cfg.CPILogTrigger {
		return cre.Workflow[config.Config]{
			cre.Handler(
				solana.LogTrigger(chainSelector, filterLogTriggerRequestCPI),
				onLogTriggerCPI,
			),
		}, nil
	}

	return cre.Workflow[config.Config]{
		cre.Handler(
			solana.LogTrigger(chainSelector, filterLogTriggerRequest),
			onLogTrigger,
		),
	}, nil
}

func onLogTrigger(cfg config.Config, runtime cre.Runtime, payload *solana.Log) (string, error) {
	return handleTestEventLog(cfg, runtime, payload,
		"TestEvent received: ",
		"Log trigger received event at block %d",
	)
}

func onLogTriggerCPI(cfg config.Config, runtime cre.Runtime, payload *solana.Log) (string, error) {
	return handleTestEventLog(cfg, runtime, payload,
		"TestEvent CPI received: ",
		"CPI log trigger received event at block %d",
	)
}

func handleTestEventLog(
	cfg config.Config,
	runtime cre.Runtime,
	payload *solana.Log,
	logPrefix string,
	resultFmt string,
) (string, error) {
	evt, err := decodeTestEvent(payload.Data)
	if err != nil {
		return "", fmt.Errorf("decode TestEvent: %w", err)
	}
	if err := validateTestEvent(cfg, evt); err != nil {
		return "", err
	}

	decoded := formatTestEvent(evt)
	runtime.Logger().Info(logPrefix+decoded,
		"blockNumber", payload.BlockNumber,
		"txHash", fmt.Sprintf("%x", payload.TxHash),
	)
	return fmt.Sprintf(resultFmt, payload.BlockNumber), nil
}

func main() {
	wasm.NewRunner(func(configBytes []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return cfg, nil
	}).Run(RunSolLogTriggerWorkflow)
}
