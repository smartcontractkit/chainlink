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

	return cre.Workflow[config.Config]{
		cre.Handler(
			solana.LogTrigger(chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector, filterLogTriggerRequest),
			onLogTrigger,
		),
		cre.Handler(
			solana.LogTrigger(chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector, filterLogTriggerRequestCPI),
			onLogTriggerCPI,
		),
	}, nil
}

func onLogTrigger(cfg config.Config, runtime cre.Runtime, payload *solana.Log) (string, error) {
	runtime.Logger().Info("TestEvent received!",
		"blockNumber", payload.BlockNumber,
		"txHash", fmt.Sprintf("%x", payload.TxHash),
	)
	return fmt.Sprintf("Log trigger received event at block %d", payload.BlockNumber), nil
}

func onLogTriggerCPI(cfg config.Config, runtime cre.Runtime, payload *solana.Log) (string, error) {
	runtime.Logger().Info("TestEvent CPI received!",
		"blockNumber", payload.BlockNumber,
		"txHash", fmt.Sprintf("%x", payload.TxHash),
	)
	return fmt.Sprintf("CPI log trigger received event at block %d", payload.BlockNumber), nil
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
