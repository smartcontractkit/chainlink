package s4_test

import (
	"errors"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/stretchr/testify/assert"
)

func setupORM(t *testing.T, namespace string) s4.ORM {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	orm := s4.NewPostgresORM(db, lggr, pgtest.NewQConfig(true), s4.SharedTableName, namespace)

	t.Cleanup(func() {
		assert.NoError(t, db.Close())
	})

	return orm
}

func generateTestRows(t *testing.T, n int) []*s4.Row {
	t.Helper()

	rows := make([]*s4.Row, n)
	for i := 0; i < n; i++ {
		row := &s4.Row{
			Address:    utils.NewBig(testutils.NewAddress().Big()),
			SlotId:     1,
			Payload:    cltest.MustRandomBytes(t, 32),
			Version:    1 + uint64(i),
			Expiration: time.Now().Add(time.Hour).UnixMilli(),
			Confirmed:  i%2 == 0,
			Signature:  cltest.MustRandomBytes(t, 32),
		}
		rows[i] = row
	}

	return rows
}

func TestNewPostgresOrm(t *testing.T) {
	t.Parallel()

	orm := setupORM(t, "test")
	assert.NotNil(t, orm)
}

func TestPostgresORM_UpdateAndGet(t *testing.T) {
	t.Parallel()

	orm := setupORM(t, "test")
	rows := generateTestRows(t, 10)

	for _, row := range rows {
		err := orm.Update(row)
		assert.NoError(t, err)

		row.Version++
		err = orm.Update(row)
		assert.NoError(t, err)

		err = orm.Update(row)
		if !row.Confirmed {
			assert.ErrorIs(t, err, s4.ErrVersionTooLow)
		}
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

	orm := setupORM(t, "test")

	const total = 10
	const expired = 4
	rows := generateTestRows(t, total)

	for _, row := range rows {
		err := orm.Update(row)
		assert.NoError(t, err)
	}

	deleted, err := orm.DeleteExpired(expired, time.Now().Add(2*time.Hour).UTC())
	assert.NoError(t, err)
	assert.Equal(t, int64(expired), deleted)

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

	orm := setupORM(t, "test")

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
			assert.Equal(t, len(rows), len(snapshot))

			snapshotRowMap := make(map[string]*s4.SnapshotRow)
			for i, sr := range snapshot {
				// assuming unique addresses
				snapshotRowMap[sr.Address.String()] = snapshot[i]
			}

			for _, sr := range rows {
				snapshotRow, ok := snapshotRowMap[sr.Address.String()]
				assert.True(t, ok)
				assert.Equal(t, snapshotRow.Address, sr.Address)
				assert.Equal(t, snapshotRow.SlotId, sr.SlotId)
				assert.Equal(t, snapshotRow.Version, sr.Version)
				assert.Equal(t, snapshotRow.Expiration, sr.Expiration)
				assert.Equal(t, snapshotRow.Confirmed, sr.Confirmed)
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

	orm := setupORM(t, "test")

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

func TestPostgresORM_Namespace(t *testing.T) {
	t.Parallel()

	ormA := setupORM(t, "a")
	ormB := setupORM(t, "b")

	const n = 10
	rowsA := generateTestRows(t, n)
	rowsB := generateTestRows(t, n)
	for i := 0; i < n; i++ {
		err := ormA.Update(rowsA[i])
		assert.NoError(t, err)

		err = ormB.Update(rowsB[i])
		assert.NoError(t, err)
	}

	urowsA, err := ormA.GetUnconfirmedRows(n)
	assert.NoError(t, err)
	assert.Len(t, urowsA, n/2)

	urowsB, err := ormB.GetUnconfirmedRows(n)
	assert.NoError(t, err)
	assert.Len(t, urowsB, n/2)

	_, err = ormB.DeleteExpired(n, time.Now().UTC())
	assert.NoError(t, err)

	snapshotA, err := ormA.GetSnapshot(s4.NewFullAddressRange())
	assert.NoError(t, err)
	assert.Len(t, snapshotA, n)
}
