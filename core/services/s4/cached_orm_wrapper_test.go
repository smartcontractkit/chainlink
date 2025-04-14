package s4_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4/mocks"
)

func TestGetSnapshotEmpty(t *testing.T) {
	t.Run("OK-no_rows", func(t *testing.T) {
		ctx := testutils.Context(t)
		psqlORM := setupORM(t, "test")
		lggr := logger.TestLogger(t)
		orm := s4.NewCachedORMWrapper(psqlORM, lggr)

		rows, err := orm.GetSnapshot(ctx, s4.NewFullAddressRange())
		assert.NoError(t, err)
		assert.Empty(t, rows)
	})
}

func TestGetSnapshotCacheFilled(t *testing.T) {
	t.Run("OK_with_rows_already_cached", func(t *testing.T) {
		ctx := testutils.Context(t)
		rows := generateTestSnapshotRows(t, 100)

		fullAddressRange := s4.NewFullAddressRange()

		lggr := logger.TestLogger(t)
		underlayingORM := mocks.NewORM(t)
		underlayingORM.On("GetSnapshot", mock.Anything, fullAddressRange).Return(rows, nil).Once()

		orm := s4.NewCachedORMWrapper(underlayingORM, lggr)

		// first call will go to the underlaying orm implementation to fill the cache
		first_snapshot, err := orm.GetSnapshot(ctx, fullAddressRange)
		assert.NoError(t, err)
		assert.Equal(t, len(rows), len(first_snapshot))

		// on the second call, the results will come from the cache, if not the mock will return an error because of .Once()
		cache_snapshot, err := orm.GetSnapshot(ctx, fullAddressRange)
		assert.NoError(t, err)
		assert.Equal(t, len(rows), len(cache_snapshot))

		snapshotRowMap := make(map[string]*s4.SnapshotRow)
		for i, sr := range cache_snapshot {
			// assuming unique addresses
			snapshotRowMap[sr.Address.String()] = cache_snapshot[i]
		}

		for _, sr := range rows {
			snapshotRow, ok := snapshotRowMap[sr.Address.String()]
			assert.True(t, ok)
			assert.NotNil(t, snapshotRow)
			assert.Equal(t, snapshotRow.Address, sr.Address)
			assert.Equal(t, snapshotRow.SlotId, sr.SlotId)
			assert.Equal(t, snapshotRow.Version, sr.Version)
			assert.Equal(t, snapshotRow.Expiration, sr.Expiration)
			assert.Equal(t, snapshotRow.Confirmed, sr.Confirmed)
			assert.Equal(t, snapshotRow.PayloadSize, sr.PayloadSize)
		}
	})
}

func TestUpdateInvalidatesSnapshotCache(t *testing.T) {
	t.Run("OK-GetSnapshot_cache_invalidated_after_update", func(t *testing.T) {
		ctx := testutils.Context(t)
		rows := generateTestSnapshotRows(t, 100)

		fullAddressRange := s4.NewFullAddressRange()

		lggr := logger.TestLogger(t)
		underlayingORM := mocks.NewORM(t)
		underlayingORM.On("GetSnapshot", mock.Anything, fullAddressRange).Return(rows, nil).Once()

		orm := s4.NewCachedORMWrapper(underlayingORM, lggr)

		// first call will go to the underlaying orm implementation to fill the cache
		first_snapshot, err := orm.GetSnapshot(ctx, fullAddressRange)
		assert.NoError(t, err)
		assert.Equal(t, len(rows), len(first_snapshot))

		// on the second call, the results will come from the cache, if not the mock will return an error because of .Once()
		cache_snapshot, err := orm.GetSnapshot(ctx, fullAddressRange)
		assert.NoError(t, err)
		assert.Equal(t, len(rows), len(cache_snapshot))

		// this update call will invalidate the cache
		row := &s4.Row{
			Address:    big.New(common.HexToAddress("0x0000000000000000000000000000000000000000000000000000000000000005").Big()),
			SlotId:     1,
			Payload:    cltest.MustRandomBytes(t, 32),
			Version:    1,
			Expiration: time.Now().Add(time.Hour).UnixMilli(),
			Confirmed:  true,
			Signature:  cltest.MustRandomBytes(t, 32),
		}
		underlayingORM.On("Update", mock.Anything, row).Return(nil).Once()
		err = orm.Update(ctx, row)
		assert.NoError(t, err)

		// given the cache was invalidated this request will reach the underlaying orm implementation
		underlayingORM.On("GetSnapshot", mock.Anything, fullAddressRange).Return(rows, nil).Once()
		third_snapshot, err := orm.GetSnapshot(ctx, fullAddressRange)
		assert.NoError(t, err)
		assert.Equal(t, len(rows), len(third_snapshot))
	})

	t.Run("OK-GetSnapshot_cache_not_invalidated_after_update", func(t *testing.T) {
		ctx := testutils.Context(t)
		rows := generateTestSnapshotRows(t, 5)

		addressRange := &s4.AddressRange{
			MinAddress: ubig.New(common.BytesToAddress(bytes.Repeat([]byte{0x00}, common.AddressLength)).Big()),
			MaxAddress: ubig.New(common.BytesToAddress(append(bytes.Repeat([]byte{0x00}, common.AddressLength-1), 3)).Big()),
		}

		lggr := logger.TestLogger(t)
		underlayingORM := mocks.NewORM(t)
		underlayingORM.On("GetSnapshot", mock.Anything, addressRange).Return(rows, nil).Once()

		orm := s4.NewCachedORMWrapper(underlayingORM, lggr)

		// first call will go to the underlaying orm implementation to fill the cache
		first_snapshot, err := orm.GetSnapshot(ctx, addressRange)
		assert.NoError(t, err)
		assert.Equal(t, len(rows), len(first_snapshot))

		// on the second call, the results will come from the cache, if not the mock will return an error because of .Once()
		cache_snapshot, err := orm.GetSnapshot(ctx, addressRange)
		assert.NoError(t, err)
		assert.Equal(t, len(rows), len(cache_snapshot))

		// this update call wont invalidate the cache because the address is out of the cache address range
		outOfCachedRangeAddress := ubig.New(common.BytesToAddress(append(bytes.Repeat([]byte{0x00}, common.AddressLength-1), 5)).Big())
		row := &s4.Row{
			Address:    outOfCachedRangeAddress,
			SlotId:     1,
			Payload:    cltest.MustRandomBytes(t, 32),
			Version:    1,
			Expiration: time.Now().Add(time.Hour).UnixMilli(),
			Confirmed:  true,
			Signature:  cltest.MustRandomBytes(t, 32),
		}
		underlayingORM.On("Update", mock.Anything, row).Return(nil).Once()
		err = orm.Update(ctx, row)
		assert.NoError(t, err)

		// given the cache was not invalidated this request wont reach the underlaying orm implementation
		third_snapshot, err := orm.GetSnapshot(ctx, addressRange)
		assert.NoError(t, err)
		assert.Equal(t, len(rows), len(third_snapshot))
	})
}

func TestGet(t *testing.T) {
	address := big.New(testutils.NewAddress().Big())
	var slotID uint = 1

	lggr := logger.TestLogger(t)

	t.Run("OK-Get_underlaying_ORM_returns_a_row", func(t *testing.T) {
		ctx := testutils.Context(t)
		underlayingORM := mocks.NewORM(t)
		expectedRow := &s4.Row{
			Address: address,
			SlotId:  slotID,
		}
		underlayingORM.On("Get", mock.Anything, address, slotID).Return(expectedRow, nil).Once()
		orm := s4.NewCachedORMWrapper(underlayingORM, lggr)

		row, err := orm.Get(ctx, address, slotID)
		require.NoError(t, err)
		require.Equal(t, expectedRow, row)
	})
	t.Run("NOK-Get_underlaying_ORM_returns_an_error", func(t *testing.T) {
		ctx := testutils.Context(t)
		underlayingORM := mocks.NewORM(t)
		underlayingORM.On("Get", mock.Anything, address, slotID).Return(nil, errors.New("some_error")).Once()
		orm := s4.NewCachedORMWrapper(underlayingORM, lggr)

		row, err := orm.Get(ctx, address, slotID)
		require.Nil(t, row)
		require.EqualError(t, err, "some_error")
	})
}

func TestDeletedExpired(t *testing.T) {
	var limit uint = 1
	now := time.Now()

	lggr := logger.TestLogger(t)

	t.Run("OK-DeletedExpired_underlaying_ORM_returns_a_row", func(t *testing.T) {
		ctx := testutils.Context(t)
		var expectedDeleted int64 = 10
		underlayingORM := mocks.NewORM(t)
		underlayingORM.On("DeleteExpired", mock.Anything, limit, now).Return(expectedDeleted, nil).Once()
		orm := s4.NewCachedORMWrapper(underlayingORM, lggr)

		actualDeleted, err := orm.DeleteExpired(ctx, limit, now)
		require.NoError(t, err)
		require.Equal(t, expectedDeleted, actualDeleted)
	})
	t.Run("NOK-DeletedExpired_underlaying_ORM_returns_an_error", func(t *testing.T) {
		ctx := testutils.Context(t)
		var expectedDeleted int64
		underlayingORM := mocks.NewORM(t)
		underlayingORM.On("DeleteExpired", mock.Anything, limit, now).Return(expectedDeleted, errors.New("some_error")).Once()
		orm := s4.NewCachedORMWrapper(underlayingORM, lggr)

		actualDeleted, err := orm.DeleteExpired(ctx, limit, now)
		require.EqualError(t, err, "some_error")
		require.Equal(t, expectedDeleted, actualDeleted)
	})
}

// GetUnconfirmedRows(limit uint, qopts ...pg.QOpt) ([]*Row, error)
func TestGetUnconfirmedRows(t *testing.T) {
	var limit uint = 1
	lggr := logger.TestLogger(t)

	t.Run("OK-GetUnconfirmedRows_underlaying_ORM_returns_a_row", func(t *testing.T) {
		ctx := testutils.Context(t)
		address := big.New(testutils.NewAddress().Big())
		var slotID uint = 1

		expectedRow := []*s4.Row{{
			Address: address,
			SlotId:  slotID,
		}}
		underlayingORM := mocks.NewORM(t)
		underlayingORM.On("GetUnconfirmedRows", mock.Anything, limit).Return(expectedRow, nil).Once()
		orm := s4.NewCachedORMWrapper(underlayingORM, lggr)

		actualRow, err := orm.GetUnconfirmedRows(ctx, limit)
		require.NoError(t, err)
		require.Equal(t, expectedRow, actualRow)
	})
	t.Run("NOK-GetUnconfirmedRows_underlaying_ORM_returns_an_error", func(t *testing.T) {
		ctx := testutils.Context(t)
		underlayingORM := mocks.NewORM(t)
		underlayingORM.On("GetUnconfirmedRows", mock.Anything, limit).Return(nil, errors.New("some_error")).Once()
		orm := s4.NewCachedORMWrapper(underlayingORM, lggr)

		actualRow, err := orm.GetUnconfirmedRows(ctx, limit)
		require.Nil(t, actualRow)
		require.EqualError(t, err, "some_error")
	})
}

func generateTestSnapshotRows(t *testing.T, n int) []*s4.SnapshotRow {
	t.Helper()

	rows := make([]*s4.SnapshotRow, n)
	for i := 0; i < n; i++ {
		row := &s4.SnapshotRow{
			Address:     big.New(testutils.NewAddress().Big()),
			SlotId:      1,
			PayloadSize: 32,
			Version:     1 + uint64(i),
			Expiration:  time.Now().Add(time.Hour).UnixMilli(),
			Confirmed:   i%2 == 0,
		}
		rows[i] = row
	}

	return rows
}
