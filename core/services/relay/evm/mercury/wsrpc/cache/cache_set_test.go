package cache

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CacheSet(t *testing.T) {
	lggr := logger.TestLogger(t)
	cs := newCacheSet(Config{Logger: lggr})
	ctx := testutils.Context(t)
	require.NoError(t, cs.Start(ctx))
	t.Cleanup(func() {
		assert.NoError(t, cs.Close())
	})

	t.Run("Get", func(t *testing.T) {
		c := &mockClient{}

		var err error
		var f Fetcher
		t.Run("with virgin cacheset, makes new entry and returns it", func(t *testing.T) {
			assert.Len(t, cs.caches, 0)

			f, err = cs.Get(ctx, c)
			require.NoError(t, err)

			assert.IsType(t, f, &memCache{})
			assert.Len(t, cs.caches, 1)
		})
		t.Run("with existing cache for value, returns that", func(t *testing.T) {
			var f2 Fetcher
			assert.Len(t, cs.caches, 1)

			f2, err = cs.Get(ctx, c)
			require.NoError(t, err)

			assert.IsType(t, f, &memCache{})
			assert.Equal(t, f, f2)
			assert.Len(t, cs.caches, 1)
		})
	})
}
