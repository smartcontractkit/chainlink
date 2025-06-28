package solana

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types/solana"
)

func TestContractReaderConfigRaw(t *testing.T) {
	cfg, err := DestContractReaderConfig()
	require.NoError(t, err)

	raw, err := json.Marshal(cfg)
	require.NoError(t, err)

	var result solana.ContractReader
	require.NoError(t, json.Unmarshal(raw, &result))
	require.Equal(t, cfg, result)

	cfg, err = SourceContractReaderConfig()
	require.NoError(t, err)

	raw, err = json.Marshal(cfg)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &result))
	require.Equal(t, cfg, result)
}
