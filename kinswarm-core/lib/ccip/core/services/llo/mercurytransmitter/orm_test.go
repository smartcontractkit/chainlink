package mercurytransmitter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

var (
	sURL  = "wss://example.com/mercury"
	sURL2 = "wss://mercuryserver.test"
	sURL3 = "wss://mercuryserver.example/foo"
)

func TestORM(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)

	donID := uint32(654321)
	orm := NewORM(db, donID)

	t.Run("DonID", func(t *testing.T) {
		assert.Equal(t, donID, orm.DonID())
	})

	transmissions := makeSampleTransmissions()[:2]

	t.Run("Insert", func(t *testing.T) {
		err := orm.Insert(ctx, transmissions)
		require.NoError(t, err)
	})
	t.Run("Get", func(t *testing.T) {
		result, err := orm.Get(ctx, sURL)
		require.NoError(t, err)

		assert.ElementsMatch(t, transmissions, result)

		result, err = orm.Get(ctx, "other server url")
		require.NoError(t, err)

		assert.Empty(t, result)
	})
	t.Run("Delete", func(t *testing.T) {
		err := orm.Delete(ctx, [][32]byte{transmissions[0].Hash()})
		require.NoError(t, err)

		result, err := orm.Get(ctx, sURL)
		require.NoError(t, err)

		require.Len(t, result, 1)
		assert.Equal(t, transmissions[1], result[0])

		err = orm.Delete(ctx, [][32]byte{transmissions[1].Hash()})
		require.NoError(t, err)

		result, err = orm.Get(ctx, sURL)
		require.NoError(t, err)
		require.Len(t, result, 0)
	})
	t.Run("Prune", func(t *testing.T) {
		err := orm.Insert(ctx, transmissions)
		require.NoError(t, err)

		err = orm.Prune(ctx, sURL, 1)
		require.NoError(t, err)

		result, err := orm.Get(ctx, sURL)
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, transmissions[1], result[0])

		err = orm.Prune(ctx, sURL, 1)
		require.NoError(t, err)
		result, err = orm.Get(ctx, sURL)
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, transmissions[1], result[0])

		err = orm.Prune(ctx, sURL, 0)
		require.NoError(t, err)
		result, err = orm.Get(ctx, sURL)
		require.NoError(t, err)
		require.Len(t, result, 0)
	})
}
