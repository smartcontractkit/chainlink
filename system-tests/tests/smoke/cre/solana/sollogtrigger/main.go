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

func getTestEventIdlJson() []byte {
	return []byte(`{
  "address": "J1zQwrBNBngz26jRPNWsUSZMHJwBwpkoDitXRV95LdK4",
  "metadata": {
    "name": "log_read_test",
    "version": "0.1.0",
    "spec": "0.1.0",
    "description": "Created with Anchor"
  },
  "instructions": [
    {
      "name": "create_log",
      "discriminator": [
        215,
        95,
        248,
        114,
        153,
        204,
        208,
        48
      ],
      "accounts": [
        {
          "name": "authority",
          "signer": true
        },
        {
          "name": "system_program"
        }
      ],
      "args": [
        {
          "name": "value",
          "type": "u64"
        }
      ]
    },
    {
      "name": "create_truncated_log",
      "discriminator": [
        133,
        74,
        116,
        132,
        80,
        11,
        241,
        64
      ],
      "accounts": [
        {
          "name": "authority",
          "signer": true
        },
        {
          "name": "system_program"
        }
      ],
      "args": [
        {
          "name": "value",
          "type": "u64"
        }
      ]
    }
  ],
  "events": [
    {
      "name": "TestEvent",
      "discriminator": [
        28,
        52,
        39,
        105,
        8,
        210,
        91,
        9
      ]
    }
  ],
  "types": [
    {
      "name": "TestEvent",
      "type": {
        "kind": "struct",
        "fields": [
          {
            "name": "str_val",
            "type": "string"
          },
          {
            "name": "u64_value",
            "type": "u64"
          }
        ]
      }
    }
  ]
}`)
}

func RunSolLogTriggerWorkflow(cfg config.Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[config.Config], error) {
	logger.Info("RunSolLogTriggerWorkflow called")

	eventIdlJson := getTestEventIdlJson()

	expectedValueBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(expectedValueBytes, cfg.ExpectedU64Value)

	filterLogTriggerRequest := &solana.FilterLogTriggerRequest{
		Name:         "test-event-filter",
		Address:      cfg.LogReadTestProgramID[:],
		EventName:    "TestEvent",
		ContractIdlJson: eventIdlJson,
		Subkeys: []*solana.SubkeyConfig{
			{Path: []string{"U64Value"}, Comparers: []*solana.ValueComparator{
				{
					Value:    expectedValueBytes,
					Operator: solana.ComparisonOperator_COMPARISON_OPERATOR_EQ,
				},
			}},
		},
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
