package ocr3

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOracleConfig_JSON(t *testing.T) {
	t.Run("Legacy format general OCR config", func(t *testing.T) {
		var cfg OracleConfig
		err := json.Unmarshal([]byte(ocr3Cfg), &cfg)
		require.NoError(t, err)
		// ensure the values were correctly unmarshalled
		require.Equal(t, cfg.MaxFaultyOracles, 3)
		require.IsType(t, &ConsensusCapOffchainConfig{}, cfg.OffchainConfig)
		consensusCapCfg := cfg.OffchainConfig.(*ConsensusCapOffchainConfig)
		require.Equal(t, consensusCapCfg.RequestTimeout, 30*time.Second)
		require.Equal(t, consensusCapCfg.MaxBatchSize, uint32(20))
		// ensure that marshalling back to JSON works
		asJSON, err := json.Marshal(cfg)
		require.NoError(t, err)
		var cfg2 OracleConfig
		err = json.Unmarshal(asJSON, &cfg2)
		require.NoError(t, err)
		require.Equal(t, cfg, cfg2)
	})
	t.Run("Chain Capability OCR Config", func(t *testing.T) {
		cfg := OracleConfig{
			OffchainConfigType:  OffchainConfigTypeChainCap,
			DeltaProgressMillis: 5000,
			OffchainConfig: &ChainCapOffchainConfig{
				MaxBatchSize:        100,
				MaxQueryLengthBytes: 1000000,
			},
		}
		asJSON, err := json.Marshal(cfg)
		require.NoError(t, err)

		var fromJSON OracleConfig
		err = json.Unmarshal(asJSON, &fromJSON)
		require.NoError(t, err)
		require.Equal(t, cfg, fromJSON)
	})
}
