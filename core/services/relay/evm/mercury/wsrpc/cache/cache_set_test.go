package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func Test_CacheSet(t *testing.T) {
	lggr := logger.Test(t)
	cs := newCacheSet(lggr, Config{LatestReportTTL: 1})
	disabledCs := newCacheSet(lggr, Config{LatestReportTTL: 0})
	ctx := testutils.Context(t)
	servicetest.Run(t, cs)

	t.Run("Get", func(t *testing.T) {
		c := &mockClient{}

		var err error
		var f Fetcher
		t.Run("with caching disabled, returns nil, nil", func(t *testing.T) {
			assert.Len(t, disabledCs.caches, 0)

			f, err = disabledCs.Get(ctx, c)
			require.NoError(t, err)

			assert.Nil(t, f)
			assert.Len(t, disabledCs.caches, 0)
		})

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
