//go:build wasip1

package main

import (
	"fmt"
	"log/slog"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/sollogtrigger/config"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/sollogtrigger/contracts/solana/src/generated/log_read_test"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/solana"
	solanabindings "github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/solana/bindings"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"
)

func formatTestEvent(evt log_read_test.TestEvent) string {
	return fmt.Sprintf("str_val=%s u64_value=%d", evt.StrVal, evt.U64Value)
}

func validateTestEvent(cfg config.Config, evt log_read_test.TestEvent) error {
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

	chainSelector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	client := &solana.Client{ChainSelector: chainSelector}

	logReadTest, err := log_read_test.NewLogReadTest(client)
	if err != nil {
		return nil, fmt.Errorf("create log_read_test bindings: %w", err)
	}

	filters := []log_read_test.TestEventFilters{
		{U64Value: &cfg.ExpectedU64Value},
	}

	return createLogTrigger(logReadTest, chainSelector, filters, cfg.CPILogTrigger)
}

func createLogTrigger(logReadTest *log_read_test.LogReadTest, chainSelector uint64, filters []log_read_test.TestEventFilters, cpi bool) (cre.Workflow[config.Config], error) {
	opts := &solanabindings.LogTriggerOptions{CPI: cpi}
	name := "test-event-filter"
	function := onLogTrigger
	if cpi {
		name = "test-cpi-event-filter"
		function = onLogTriggerCPI
	}
	trigger, err := logReadTest.LogTriggerTestEventLog(
		chainSelector,
		name,
		filters,
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register log trigger: %w", err)
	}
	return cre.Workflow[config.Config]{
		cre.Handler(trigger, function),
	}, nil
}

func onLogTrigger(cfg config.Config, runtime cre.Runtime, payload *solanabindings.DecodedLog[log_read_test.TestEvent]) (string, error) {
	return handleTestEventLog(cfg, runtime, payload,
		"TestEvent received: ",
		"Log trigger received event at block %d",
	)
}

func onLogTriggerCPI(cfg config.Config, runtime cre.Runtime, payload *solanabindings.DecodedLog[log_read_test.TestEvent]) (string, error) {
	return handleTestEventLog(cfg, runtime, payload,
		"TestEvent CPI received: ",
		"CPI log trigger received event at block %d",
	)
}

func handleTestEventLog(
	cfg config.Config,
	runtime cre.Runtime,
	payload *solanabindings.DecodedLog[log_read_test.TestEvent],
	logPrefix string,
	resultFmt string,
) (string, error) {
	if err := validateTestEvent(cfg, payload.Data); err != nil {
		return "", err
	}

	decoded := formatTestEvent(payload.Data)
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
