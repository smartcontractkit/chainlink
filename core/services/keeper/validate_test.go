package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidatedKeeperSpec(t *testing.T) {
	t.Parallel()
	toml := `
		type                = "keeper"
		schemaVersion       = 1
		name                = "example keeper spec"
		contractAddress     = "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba"
		fromAddress         = "0xa8037A20989AFcBC51798de9762b351D63ff462e"
	`

	s, err := ValidatedKeeperSpec(toml)
	require.NoError(t, err)

	require.Equal(t, int32(0), s.ID)
	require.Equal(t, "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba", s.KeeperSpec.ContractAddress.Hex())
	require.Equal(t, "0xa8037A20989AFcBC51798de9762b351D63ff462e", s.KeeperSpec.FromAddress.Hex())
	require.Equal(t, time.Time{}, s.KeeperSpec.CreatedAt)
	require.Equal(t, time.Time{}, s.KeeperSpec.UpdatedAt)
}
