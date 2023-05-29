package s4_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

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
		Address:    utils.NewBig(address.Big()),
		SlotId:     slotId,
		Payload:    payload[:],
		Version:    3,
		Expiration: expiration,
		Confirmed:  true,
		Signature:  signature[:],
	}

	orm := s4.NewInMemoryORM()

	t.Run("row not found", func(t *testing.T) {
		_, err := orm.Get(utils.NewBig(address.Big()), slotId)
		assert.ErrorIs(t, err, s4.ErrNotFound)
	})

	t.Run("insert and get", func(t *testing.T) {
		err := orm.Update(row)
		assert.NoError(t, err)

		e, err := orm.Get(utils.NewBig(address.Big()), slotId)
		assert.NoError(t, err)
		row.UpdatedAt = e.UpdatedAt
		assert.Equal(t, row, e)
	})

	t.Run("update and get", func(t *testing.T) {
		err := orm.Update(row)
		assert.NoError(t, err)

		row.Version = 5
		err = orm.Update(row)
		assert.NoError(t, err)

		e, err := orm.Get(utils.NewBig(address.Big()), slotId)
		assert.NoError(t, err)
		row.UpdatedAt = e.UpdatedAt
		assert.Equal(t, row, e)
	})

	t.Run("delete expired", func(t *testing.T) {
		ms := row.Expiration - time.Now().UnixMilli() + 100
		time.Sleep(time.Duration(ms) * time.Millisecond)
		err := orm.DeleteExpired()
		assert.NoError(t, err)

		_, err = orm.Get(utils.NewBig(address.Big()), slotId)
		assert.ErrorIs(t, err, s4.ErrNotFound)
	})
}

func TestInMemoryORM_GetUnconfirmedRows(t *testing.T) {
	t.Parallel()

	orm := s4.NewInMemoryORM()
	expiration := time.Now().Add(100 * time.Second).UnixMilli()

	for i := 0; i < 256; i++ {
		var thisAddress common.Address
		thisAddress[0] = byte(i)

		row := &s4.Row{
			Address:    utils.NewBig(thisAddress.Big()),
			SlotId:     1,
			Payload:    []byte{},
			Version:    1,
			Expiration: expiration,
			Confirmed:  i >= 100,
			Signature:  []byte{},
		}
		err := orm.Update(row)
		assert.NoError(t, err)
		time.Sleep(time.Millisecond)
	}

	rows, err := orm.GetUnconfirmedRows(100)
	assert.NoError(t, err)
	assert.Len(t, rows, 100)
	assert.Less(t, rows[0].UpdatedAt, rows[99].UpdatedAt)
}

func TestInMemoryORM_GetVersions(t *testing.T) {
	t.Parallel()

	orm := s4.NewInMemoryORM()
	expiration := time.Now().Add(100 * time.Second).UnixMilli()

	const n = 256
	for i := 0; i < n; i++ {
		var thisAddress common.Address
		thisAddress[0] = byte(i)

		row := &s4.Row{
			Address:    utils.NewBig(thisAddress.Big()),
			SlotId:     1,
			Payload:    []byte{},
			Version:    uint64(i),
			Expiration: expiration,
			Confirmed:  i >= 100,
			Signature:  []byte{},
		}
		err := orm.Update(row)
		assert.NoError(t, err)
	}

	versions, err := orm.GetVersions(s4.NewFullAddressRange())
	assert.NoError(t, err)
	assert.Len(t, versions, n)

	testMap := make(map[uint64]int)
	for i := 0; i < n; i++ {
		testMap[versions[i].Version]++
	}
	assert.Len(t, testMap, n)
	for _, c := range testMap {
		assert.Equal(t, 1, c)
	}
}
