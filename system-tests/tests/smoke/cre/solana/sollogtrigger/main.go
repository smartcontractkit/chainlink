//go:build wasip1

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/sollogtrigger/config"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/solana"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"gopkg.in/yaml.v3"
)

// getTestEventDiscriminator returns the 8-byte Anchor event discriminator
func getTestEventDiscriminator() [8]byte {
	hash := sha256.Sum256([]byte("event:TestEvent"))
	var discriminator [8]byte
	copy(discriminator[:], hash[:8])
	return discriminator
}

func getTestEventIdlJson() ([]byte, error) {
	idl := map[string]interface{}{
		"event": map[string]interface{}{
			"name": "TestEvent",
			"fields": []map[string]interface{}{
				{"name": "str_val", "type": "string", "index": false},
				{"name": "u64_value", "type": "u64", "index": false},
			},
		},
		"types": []interface{}{},
	}
	return json.Marshal(idl)
}

func RunSolLogTriggerWorkflow(cfg config.Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	logger.Info("RunSolLogTriggerWorkflow called")

	discriminator := getTestEventDiscriminator()
	eventIdlJson, err := getTestEventIdlJson()
	if err != nil {
		return nil, fmt.Errorf("failed to generate event IDL JSON: %w", err)
	}

	filterLogTriggerRequest := &solana.FilterLogTriggerRequest{
		Name:         "test-event-filter",
		Address:      cfg.LogReadTestProgramID[:],
		EventName:    "TestEvent",
		EventSig:     discriminator[:],
		EventIdlJson: eventIdlJson,
	}

	return cre.Workflow[config.Config]{
		cre.Handler(
			solana.LogTrigger(chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector, filterLogTriggerRequest),
			onLogTrigger,
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

func main() {
	wasm.NewRunner(func(configBytes []byte) (config.Config, error) {
		cfg := config.Config{}
		if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
			return config.Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		return cfg, nil
	}).Run(RunSolLogTriggerWorkflow)
}
