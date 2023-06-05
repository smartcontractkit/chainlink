package s4_test

import (
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/stretchr/testify/assert"
)

func setupORM(t *testing.T) s4.ORM {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	orm := s4.NewPostgresORM(db, lggr, pgtest.NewQConfig(true), "functions")

	t.Cleanup(func() {
		assert.NoError(t, db.Close())
	})

	return orm
}

func mustRandomBytes(t *testing.T, n int) []byte {
	b := make([]byte, n)
	k, err := rand.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, n, k)
	return b
}

func generateTestRows(t *testing.T, n int) []*s4.Row {
	t.Helper()

	rows := make([]*s4.Row, n)
	for i := 0; i < n; i++ {
		row := &s4.Row{
			Address:    utils.NewBig(testutils.NewAddress().Big()),
			SlotId:     1,
			Payload:    mustRandomBytes(t, 32),
			Version:    1 + uint64(i),
			Expiration: time.Now().Add(time.Hour).UnixMilli(),
			Confirmed:  i%2 == 0,
			Signature:  mustRandomBytes(t, 32),
		}
		rows[i] = row
	}

	return rows
}

func TestNewPostgresOrm(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	assert.NotNil(t, orm)
}

func TestPostgresORM_UpdateAndGet(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)
	rows := generateTestRows(t, 10)

	for _, row := range rows {
		err := orm.Update(row)
		assert.NoError(t, err)

		row.Version++
		err = orm.Update(row)
		assert.NoError(t, err)

		err = orm.Update(row)
		assert.ErrorIs(t, err, s4.ErrVersionTooLow)
	}

	for _, row := range rows {
		gotRow, err := orm.Get(row.Address, row.SlotId)
		assert.NoError(t, err)
		assert.Equal(t, row, gotRow)
	}

	rows = generateTestRows(t, 1)
	_, err := orm.Get(rows[0].Address, rows[0].SlotId)
	assert.ErrorIs(t, err, s4.ErrNotFound)
}

func TestPostgresORM_DeleteExpired(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)

	const total = 10
	const expired = 4
	rows := generateTestRows(t, total)

	for _, row := range rows {
		err := orm.Update(row)
		assert.NoError(t, err)
	}

	err := orm.DeleteExpired(expired, time.Now().Add(2*time.Hour).UTC())
	assert.NoError(t, err)

	count := 0
	for _, row := range rows {
		_, err := orm.Get(row.Address, row.SlotId)
		if !errors.Is(err, s4.ErrNotFound) {
			count++
		}
	}
	assert.Equal(t, total-expired, count)
}

func TestPostgresORM_GetSnapshot(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)

	t.Run("no rows", func(t *testing.T) {
		rows, err := orm.GetSnapshot(s4.NewFullAddressRange())
		assert.NoError(t, err)
		assert.Empty(t, rows)
	})

	t.Run("with rows", func(t *testing.T) {
		rows := generateTestRows(t, 100)

		for _, row := range rows {
			err := orm.Update(row)
			assert.NoError(t, err)
		}

		t.Run("full range", func(t *testing.T) {
			snapshot, err := orm.GetSnapshot(s4.NewFullAddressRange())
			assert.NoError(t, err)
			for i, sr := range snapshot {
				assert.Equal(t, rows[i].Address, sr.Address)
				assert.Equal(t, rows[i].SlotId, sr.SlotId)
				assert.Equal(t, rows[i].Version, sr.Version)
			}
		})

		t.Run("half range", func(t *testing.T) {
			ar, err := s4.NewInitialAddressRangeForIntervals(2)
			assert.NoError(t, err)
			snapshot, err := orm.GetSnapshot(ar)
			assert.NoError(t, err)
			for _, sr := range snapshot {
				assert.True(t, ar.Contains(sr.Address))
			}
		})
	})
}

func TestPostgresORM_GetUnconfirmedRows(t *testing.T) {
	t.Parallel()

	orm := setupORM(t)

	t.Run("no rows", func(t *testing.T) {
		rows, err := orm.GetUnconfirmedRows(5)
		assert.NoError(t, err)
		assert.Empty(t, rows)
	})

	t.Run("with rows", func(t *testing.T) {
		rows := generateTestRows(t, 10)

		for _, row := range rows {
			err := orm.Update(row)
			assert.NoError(t, err)
			time.Sleep(testutils.TestInterval / 10)
		}

		gotRows, err := orm.GetUnconfirmedRows(5)
		assert.NoError(t, err)
		assert.Len(t, gotRows, 5)

		for _, row := range gotRows {
			assert.False(t, row.Confirmed)
		}
	})
}
