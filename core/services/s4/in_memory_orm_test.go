package s4_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryORM(t *testing.T) {
	t.Parallel()

	address := testutils.NewAddress()
	var slotId uint = 3
	payload := testutils.Random32Byte()
	signature := testutils.Random32Byte()
	expiration := time.Now().Add(time.Minute).UnixMilli()
	row := &s4.Row{
		Address:    big.New(address.Big()),
		SlotId:     slotId,
		Payload:    payload[:],
		Version:    3,
		Expiration: expiration,
		Confirmed:  false,
		Signature:  signature[:],
	}

	orm := s4.NewInMemoryORM()

	t.Run("row not found", func(t *testing.T) {
		ctx := testutils.Context(t)
		_, err := orm.Get(ctx, big.New(address.Big()), slotId)
		assert.ErrorIs(t, err, s4.ErrNotFound)
	})

	t.Run("insert and get", func(t *testing.T) {
		ctx := testutils.Context(t)
		err := orm.Update(ctx, row)
		assert.NoError(t, err)

		e, err := orm.Get(ctx, big.New(address.Big()), slotId)
		assert.NoError(t, err)
		assert.Equal(t, row, e)
	})

	t.Run("update and get", func(t *testing.T) {
		ctx := testutils.Context(t)
		row.Version = 5
		err := orm.Update(ctx, row)
		assert.NoError(t, err)

		// unconfirmed row requires greater version
		err = orm.Update(ctx, row)
		assert.ErrorIs(t, err, s4.ErrVersionTooLow)

		row.Confirmed = true
		err = orm.Update(ctx, row)
		assert.NoError(t, err)

		e, err := orm.Get(ctx, big.New(address.Big()), slotId)
		assert.NoError(t, err)
		assert.Equal(t, row, e)
	})
}

func TestInMemoryORM_DeleteExpired(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm := s4.NewInMemoryORM()
	baseTime := time.Now().Add(time.Minute).UTC()

	for i := 0; i < 256; i++ {
		var thisAddress common.Address
		thisAddress[0] = byte(i)

		row := &s4.Row{
			Address:    big.New(thisAddress.Big()),
			SlotId:     1,
			Payload:    []byte{},
			Version:    1,
			Expiration: baseTime.Add(time.Duration(i) * time.Second).UnixMilli(),
			Confirmed:  false,
			Signature:  []byte{},
		}
		err := orm.Update(ctx, row)
		assert.NoError(t, err)
	}

	deadline := baseTime.Add(100 * time.Second)
	count, err := orm.DeleteExpired(ctx, 200, deadline)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), count)

	rows, err := orm.GetUnconfirmedRows(ctx, 200)
	assert.NoError(t, err)
	assert.Len(t, rows, 156)
}

func TestInMemoryORM_GetUnconfirmedRows(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm := s4.NewInMemoryORM()
	expiration := time.Now().Add(100 * time.Second).UnixMilli()

	for i := 0; i < 256; i++ {
		var thisAddress common.Address
		thisAddress[0] = byte(i)

		row := &s4.Row{
			Address:    big.New(thisAddress.Big()),
			SlotId:     1,
			Payload:    []byte{},
			Version:    1,
			Expiration: expiration,
			Confirmed:  i >= 100,
			Signature:  []byte{},
		}
		err := orm.Update(ctx, row)
		assert.NoError(t, err)
		time.Sleep(time.Millisecond)
	}

	rows, err := orm.GetUnconfirmedRows(ctx, 100)
	assert.NoError(t, err)
	assert.Len(t, rows, 100)
}

func TestInMemoryORM_GetSnapshot(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	orm := s4.NewInMemoryORM()
	expiration := time.Now().Add(100 * time.Second).UnixMilli()

	const n = 256
	for i := 0; i < n; i++ {
		var thisAddress common.Address
		thisAddress[0] = byte(i)

		row := &s4.Row{
			Address:    big.New(thisAddress.Big()),
			SlotId:     1,
			Payload:    []byte{},
			Version:    uint64(i),
			Expiration: expiration,
			Confirmed:  i >= 100,
			Signature:  []byte{},
		}
		err := orm.Update(ctx, row)
		assert.NoError(t, err)
	}

	rows, err := orm.GetSnapshot(ctx, s4.NewFullAddressRange())
	assert.NoError(t, err)
	assert.Len(t, rows, n)

	testMap := make(map[uint64]int)
	for i := 0; i < n; i++ {
		testMap[rows[i].Version]++
	}
	assert.Len(t, testMap, n)
	for _, c := range testMap {
		assert.Equal(t, 1, c)
	}
}
