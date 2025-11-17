package ocr3

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestOracleConfig_JSON(t *testing.T) {
	t.Run("Legacy config returns an error", func(t *testing.T) {
		var cfg OracleConfig
		err := json.Unmarshal([]byte(legacyOcr3Cfg), &cfg)
		require.ErrorContains(t, err, "not supported config format detected: field MaxQueryLengthBytes is not supported")
	})
	t.Run("Legacy yaml config returns an error", func(t *testing.T) {
		var cfg OracleConfig
		err := yaml.Unmarshal([]byte(legacyOCR3YamlCfg), &cfg)
		require.ErrorContains(t, err, "not supported config format detected: field maxQueryLengthBytes is not supported")
	})
	t.Run("Consensus Capability OCR config", func(t *testing.T) {
		var cfg OracleConfig
		err := json.Unmarshal([]byte(ocr3Cfg), &cfg)
		require.NoError(t, err)
		// ensure the values were correctly unmarshalled
		require.Equal(t, 3, cfg.MaxFaultyOracles)
		consensusCapCfg := cfg.ConsensusCapOffchainConfig
		require.Equal(t, 30*time.Second, consensusCapCfg.RequestTimeout)
		require.Equal(t, uint32(20), consensusCapCfg.MaxBatchSize)
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
			DeltaProgressMillis: 5000,
			ChainCapOffchainConfig: &ChainCapOffchainConfig{
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

var legacyOCR3YamlCfg = `
deltaProgressMillis: 5000
deltaResendMillis: 5000
deltaInitialMillis: 5000
deltaRoundMillis: 2000
deltaGraceMillis: 500
deltaCertifiedCommitRequestMillis: 1000
deltaStageMillis: 30000
maxRoundsPerEpoch: 10
transmissionSchedule:
- 7
maxDurationQueryMillis: 1000
maxDurationObservationMillis: 1000
maxDurationShouldAcceptMillis: 1000
maxDurationShouldTransmitMillis: 1000
maxFaultyOracles: 2
maxQueryLengthBytes: 1000000
maxObservationLengthBytes: 1000000
maxReportLengthBytes: 1000000
maxBatchSize: 1000
uniqueReports: true
`

var legacyOcr3Cfg = `
{
  "MaxQueryLengthBytes": 1000000,
  "MaxObservationLengthBytes": 1000000,
  "MaxReportLengthBytes": 1000000,
  "MaxOutcomeLengthBytes": 1000000,
  "MaxReportCount": 20,
  "MaxBatchSize": 20,
  "OutcomePruningThreshold": 3600,
  "UniqueReports": true,
  "RequestTimeout": "30s",
  "DeltaProgressMillis": 5000,
  "DeltaResendMillis": 5000,
  "DeltaInitialMillis": 5000,
  "DeltaRoundMillis": 2000,
  "DeltaGraceMillis": 500,
  "DeltaCertifiedCommitRequestMillis": 1000,
  "DeltaStageMillis": 30000,
  "MaxRoundsPerEpoch": 10,
  "TransmissionSchedule": [
    10
  ],
  "MaxDurationQueryMillis": 1000,
  "MaxDurationObservationMillis": 1000,
  "MaxDurationReportMillis": 1000,
  "MaxDurationShouldAcceptMillis": 1000,
  "MaxDurationShouldTransmitMillis": 1000,
  "MaxFaultyOracles": 3
}`
