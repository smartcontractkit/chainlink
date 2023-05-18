package s4_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryORM(t *testing.T) {
	t.Parallel()

	address := testutils.NewAddress()
	slotId := 3
	payload := testutils.Random32Byte()
	signature := testutils.Random32Byte()
	expiration := time.Now().Add(100 * time.Millisecond).UnixMilli()
	entry := &s4.Entry{
		Payload:           payload[:],
		Version:           3,
		Expiration:        expiration,
		Confirmed:         true,
		HighestExpiration: expiration,
		Signature:         signature[:],
	}

	orm := s4.NewInMemoryORM()

	t.Run("entry not found", func(t *testing.T) {
		_, err := orm.Get(address, slotId)
		assert.ErrorIs(t, err, s4.ErrEntryNotFound)
	})

	t.Run("insert and get", func(t *testing.T) {
		err := orm.Upsert(address, slotId, entry)
		assert.NoError(t, err)

		e, err := orm.Get(address, slotId)
		assert.NoError(t, err)
		assert.Equal(t, entry, e)
	})

	t.Run("update and get", func(t *testing.T) {
		err := orm.Upsert(address, slotId, entry)
		assert.NoError(t, err)

		entry.Version = 4
		err = orm.Upsert(address, slotId, entry)
		assert.NoError(t, err)

		e, err := orm.Get(address, slotId)
		assert.NoError(t, err)
		assert.Equal(t, entry, e)
	})

	t.Run("delete expired", func(t *testing.T) {
		ms := entry.Expiration - time.Now().UnixMilli() + 100
		time.Sleep(time.Duration(ms) * time.Millisecond)
		err := orm.DeleteExpired()
		assert.NoError(t, err)

		_, err = orm.Get(address, slotId)
		assert.ErrorIs(t, err, s4.ErrEntryNotFound)
	})
}
