package solana

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
)

func TestContractReaderConfigRaw(t *testing.T) {
	cfg, err := DestContractReaderConfig()
	require.NoError(t, err)

	raw, err := json.Marshal(cfg)
	require.NoError(t, err)

	var result config.ContractReader
	require.NoError(t, json.Unmarshal(raw, &result))
	require.EqualValues(t, cfg, result)

	cfg, err = SourceContractReaderConfig()
	require.NoError(t, err)

	raw, err = json.Marshal(cfg)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &result))
	require.EqualValues(t, cfg, result)
}
