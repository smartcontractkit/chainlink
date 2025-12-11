package solana

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
)

func TestContractReaderConfigRaw(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	cfg, err := DestContractReaderConfig()
	require.NoError(t, err)

	raw, err := json.Marshal(cfg)
	require.NoError(t, err)

	var result config.ContractReader
	require.NoError(t, json.Unmarshal(raw, &result))
	require.Equal(t, cfg, result)

	cfg, err = SourceContractReaderConfig()
	require.NoError(t, err)

	raw, err = json.Marshal(cfg)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &result))
	require.Equal(t, cfg, result)
}
