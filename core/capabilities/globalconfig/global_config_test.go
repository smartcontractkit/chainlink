package globalconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalConfig_Store(t *testing.T) {
	t.Parallel()

	t.Run("applies first version", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.NoError(t, g.Store(Update{Raw: `{"version":1}`, Hash: "h1"}))
		raw, v := g.Load()
		assert.Equal(t, `{"version":1}`, raw)
		assert.Equal(t, uint64(1), v)
	})

	t.Run("accepts strictly newer version", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.NoError(t, g.Store(Update{Raw: `{"version":1}`, Hash: "h1"}))
		require.NoError(t, g.Store(Update{Raw: `{"version":2}`, Hash: "h2"}))
		_, v := g.Load()
		assert.Equal(t, uint64(2), v)
	})

	t.Run("rejects non-increasing version", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.NoError(t, g.Store(Update{Raw: `{"version":5}`, Hash: "h5"}))
		err := g.Store(Update{Raw: `{"version":5}`, Hash: "hDifferent"})
		require.Error(t, err)
		err = g.Store(Update{Raw: `{"version":4}`, Hash: "h4"})
		require.Error(t, err)
		_, v := g.Load()
		assert.Equal(t, uint64(5), v)
	})

	t.Run("idempotent re-apply of same hash", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.NoError(t, g.Store(Update{Raw: `{"version":5}`, Hash: "h5"}))
		require.NoError(t, g.Store(Update{Raw: `{"version":5}`, Hash: "h5"}))
		_, v := g.Load()
		assert.Equal(t, uint64(5), v)
	})

	t.Run("proto-json string version", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.NoError(t, g.Store(Update{Raw: `{"version":"42"}`, Hash: "h"}))
		_, v := g.Load()
		assert.Equal(t, uint64(42), v)
	})

	t.Run("rejects invalid payload", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.Error(t, g.Store(Update{Raw: ``, Hash: "h"}))
		require.Error(t, g.Store(Update{Raw: `not json`, Hash: "h"}))
		require.Error(t, g.Store(Update{Raw: `{"version":"nope"}`, Hash: "h"}))
	})
}

func TestValidate(t *testing.T) {
	t.Parallel()
	require.NoError(t, Validate(`{"version":1,"dons":{}}`))
	require.NoError(t, Validate(`{"version":"1"}`))
	require.Error(t, Validate(``))
	require.Error(t, Validate(`{`))
}
