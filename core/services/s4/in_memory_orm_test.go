package s4_test

import (
	"testing"
	"time"

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
	expiration := time.Now().Add(100 * time.Millisecond).UnixMilli()
	row := &s4.Row{
		Address:    address.String(),
		SlotId:     slotId,
		Payload:    payload[:],
		Version:    3,
		Expiration: expiration,
		Confirmed:  true,
		Signature:  signature[:],
	}

	orm := s4.NewInMemoryORM()

	t.Run("row not found", func(t *testing.T) {
		_, err := orm.Get(address, slotId)
		assert.ErrorIs(t, err, s4.ErrNotFound)
	})

	t.Run("insert and get", func(t *testing.T) {
		err := orm.Update(row)
		assert.NoError(t, err)

		e, err := orm.Get(address, slotId)
		assert.NoError(t, err)
		row.UpdatedAt = e.UpdatedAt
		assert.Equal(t, row, e)
	})

	t.Run("update and get", func(t *testing.T) {
		err := orm.Update(row)
		assert.ErrorIs(t, err, s4.ErrVersionTooLow)

		row.Version = 5
		err = orm.Update(row)
		assert.NoError(t, err)

		e, err := orm.Get(address, slotId)
		assert.NoError(t, err)
		row.UpdatedAt = e.UpdatedAt
		assert.Equal(t, row, e)
	})

	t.Run("delete expired", func(t *testing.T) {
		ms := row.Expiration - time.Now().UnixMilli() + 100
		time.Sleep(time.Duration(ms) * time.Millisecond)
		err := orm.DeleteExpired()
		assert.NoError(t, err)

		_, err = orm.Get(address, slotId)
		assert.ErrorIs(t, err, s4.ErrNotFound)
	})

	t.Run("snapshots", func(t *testing.T) {
		expiration := time.Now().Add(100 * time.Second).UnixMilli()

		for i := 0; i < 256; i++ {
			var thisAddress common.Address
			thisAddress[0] = byte(i)

			row := &s4.Row{
				Address:    thisAddress.String(),
				SlotId:     1,
				Payload:    []byte{},
				Version:    1,
				Expiration: expiration,
				Confirmed:  false,
				Signature:  []byte{},
			}
			err := orm.Update(row)
			assert.NoError(t, err)
		}

		rows, err := orm.GetSnapshot(s4.NewFullAddressRange())
		assert.NoError(t, err)
		assert.Len(t, rows, 256)

		ar, err := s4.NewInitialAddressRangeForIntervals(2)
		assert.NoError(t, err)

		rows, err = orm.GetSnapshot(ar)
		assert.NoError(t, err)
		assert.Len(t, rows, 128)
	})
}
