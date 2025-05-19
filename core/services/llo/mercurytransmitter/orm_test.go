package mercurytransmitter

import (
	"testing"
	"time"

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

	t.Run("Insert, Get, Delete, Prune", func(t *testing.T) {
		donID := uint32(654321)
		orm := NewORM(db, donID)

		t.Run("DonID", func(t *testing.T) {
			assert.Equal(t, donID, orm.DonID())
		})
		const n = 10
		transmissions := makeSampleTransmissions(n, sURL)
		// Insert
		err := orm.Insert(ctx, transmissions)
		require.NoError(t, err)
		// Get limits
		result, err := orm.Get(ctx, sURL, 0, 0)
		require.NoError(t, err)
		assert.Empty(t, result)

		// Get limits
		result, err = orm.Get(ctx, sURL, 1, 0)
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, transmissions[len(transmissions)-1], result[0])

		result, err = orm.Get(ctx, sURL, 100, 0)
		require.NoError(t, err)

		assert.ElementsMatch(t, transmissions, result)

		result, err = orm.Get(ctx, "other server url", 100, 0)
		require.NoError(t, err)
		assert.Empty(t, result)
		// Delete
		err = orm.Delete(ctx, [][32]byte{transmissions[0].Hash()})
		require.NoError(t, err)

		result, err = orm.Get(ctx, sURL, 100, 0)
		require.NoError(t, err)

		require.Len(t, result, n-1)
		assert.NotContains(t, result, transmissions[0])
		assert.Contains(t, result, transmissions[1])

		err = orm.Delete(ctx, [][32]byte{transmissions[1].Hash()})
		require.NoError(t, err)

		result, err = orm.Get(ctx, sURL, 100, 0)
		require.NoError(t, err)
		require.Len(t, result, n-2)
		// Prune
		// ensure that len(transmissions) exceeds batch size to test batching
		err = orm.Insert(ctx, transmissions)
		require.NoError(t, err)

		d, err := orm.Prune(ctx, sURL, 1, n/3)
		require.NoError(t, err)
		assert.Equal(t, int64(n-1), d)

		result, err = orm.Get(ctx, sURL, 100, 0)
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, transmissions[len(transmissions)-1], result[0])

		// Prune again, should not delete anything
		d, err = orm.Prune(ctx, sURL, 1, n/3)
		require.NoError(t, err)
		assert.Zero(t, d)
		result, err = orm.Get(ctx, sURL, 100, 0)
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, transmissions[len(transmissions)-1], result[0])

		// Pruning with max allowed records = 0 deletes everything
		d, err = orm.Prune(ctx, sURL, 0, 1)
		require.NoError(t, err)
		assert.Equal(t, int64(1), d)
		result, err = orm.Get(ctx, sURL, 100, 0)
		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("Prune trims to exactly the correct number of records", func(t *testing.T) {
		donID := uint32(100)
		orm := NewORM(db, donID)

		var transmissions []*Transmission
		// create 100 records (10 * 10 duplicate sequence numbers)
		for seqNr := uint64(0); seqNr < 10; seqNr++ {
			for i := 0; i < 10; i++ {
				transmissions = append(transmissions, makeSampleTransmission(seqNr, sURL, []byte{byte(i)}))
			}
		}
		err := orm.Insert(ctx, transmissions)
		require.NoError(t, err)

		d, err := orm.Prune(ctx, sURL, 43, 3)
		require.NoError(t, err)
		assert.Equal(t, int64(57), d)

		result, err := orm.Get(ctx, sURL, 100, 0)
		require.NoError(t, err)
		require.Len(t, result, 43)
		assert.Equal(t, uint64(9), result[0].SeqNr)
	})

	t.Run("Get respects maxAge argument and does not retrieve records older than this", func(t *testing.T) {
		donID := uint32(101)
		orm := NewORM(db, donID)

		transmissions := makeSampleTransmissions(10, sURL)
		err := orm.Insert(ctx, transmissions)
		require.NoError(t, err)

		pgtest.MustExec(t, db, `UPDATE llo_mercury_transmit_queue SET inserted_at = NOW() - INTERVAL '1 year' WHERE seq_nr < 5`)

		// Get with maxAge = 0 should return all records
		result, err := orm.Get(ctx, sURL, 100, 0)
		require.NoError(t, err)
		require.Len(t, result, 10)

		// Get with maxAge = 1 month should return only the records with seq_nr >= 5
		result, err = orm.Get(ctx, sURL, 100, 30*24*time.Hour)
		require.NoError(t, err)
		require.Len(t, result, 5)
	})
}
