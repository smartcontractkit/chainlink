package chainlink

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config"
)

type stubSharding struct {
	config.Sharding
	enabled bool
	index   uint16
}

func (s stubSharding) ShardingEnabled() bool { return s.enabled }
func (s stubSharding) ShardIndex() uint16    { return s.index }

func TestLocalShardIndex(t *testing.T) {
	t.Parallel()

	require.Equal(t, uint32(0), localShardIndex(nil))
	require.Equal(t, uint32(0), localShardIndex(stubSharding{enabled: false, index: 3}), "disabled sharding reports shard 0")
	require.Equal(t, uint32(3), localShardIndex(stubSharding{enabled: true, index: 3}))
}
