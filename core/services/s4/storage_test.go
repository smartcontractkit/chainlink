package s4_test

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	constraints = s4.Constraints{
		MaxSlotsPerUser:        5,
		MaxPayloadSizeBytes:    32,
		MaxExpirationLengthSec: 3600,
	}
)

func setupTestStorage(t *testing.T, now time.Time) (*mocks.ORM, s4.Storage) {
	logger := logger.TestLogger(t)
	orm := mocks.NewORM(t)
	clock := clockwork.NewFakeClock()
	storage := s4.NewStorage(logger, constraints, orm, clock)
	return orm, storage
}

func TestStorage_Constraints(t *testing.T) {
	t.Parallel()

	_, storage := setupTestStorage(t, time.Now())
	c := storage.Constraints()
	assert.Equal(t, constraints, c)
}

func TestStorage_Errors(t *testing.T) {
	t.Parallel()

	now := time.Now()
	ormMock, storage := setupTestStorage(t, now)

	t.Run("ErrNotFound", func(t *testing.T) {
		key := &s4.Key{
			Address: testutils.NewAddress(),
			SlotId:  1,
			Version: 0,
		}
		ormMock.On("Get", mock.Anything, big.New(key.Address.Big()), key.SlotId).Return(nil, s4.ErrNotFound)
		_, _, err := storage.Get(testutils.Context(t), key)
		assert.ErrorIs(t, err, s4.ErrNotFound)
	})

	t.Run("ErrSlotIdTooBig", func(t *testing.T) {
		key := &s4.Key{
			Address: testutils.NewAddress(),
			SlotId:  constraints.MaxSlotsPerUser + 1,
			Version: 0,
		}
		_, _, err := storage.Get(testutils.Context(t), key)
		assert.ErrorIs(t, err, s4.ErrSlotIdTooBig)

		record := &s4.Record{
			Payload:    make([]byte, 10),
			Expiration: now.Add(time.Minute).UnixMilli(),
		}
		err = storage.Put(testutils.Context(t), key, record, []byte{})
		assert.ErrorIs(t, err, s4.ErrSlotIdTooBig)
	})

	t.Run("ErrPayloadTooBig", func(t *testing.T) {
		key := &s4.Key{
			Address: testutils.NewAddress(),
			SlotId:  1,
			Version: 0,
		}
		record := &s4.Record{
			Payload:    make([]byte, constraints.MaxPayloadSizeBytes+1),
			Expiration: now.Add(time.Minute).UnixMilli(),
		}
		err := storage.Put(testutils.Context(t), key, record, []byte{})
		assert.ErrorIs(t, err, s4.ErrPayloadTooBig)
	})

	t.Run("ErrPastExpiration", func(t *testing.T) {
		key := &s4.Key{
			Address: testutils.NewAddress(),
			SlotId:  1,
			Version: 0,
		}
		record := &s4.Record{
			Payload:    make([]byte, 10),
			Expiration: now.UnixMilli() - 1,
		}
		err := storage.Put(testutils.Context(t), key, record, []byte{})
		assert.ErrorIs(t, err, s4.ErrPastExpiration)
	})

	t.Run("ErrExpirationTooLong", func(t *testing.T) {
		key := &s4.Key{
			Address: testutils.NewAddress(),
			SlotId:  1,
			Version: 0,
		}
		record := &s4.Record{
			Payload:    make([]byte, 10),
			Expiration: now.UnixMilli() + 10000000,
		}
		err := storage.Put(testutils.Context(t), key, record, []byte{})
		assert.ErrorIs(t, err, s4.ErrExpirationTooLong)
	})

	t.Run("ErrWrongSignature", func(t *testing.T) {
		privateKey, address := testutils.NewPrivateKeyAndAddress(t)
		key := &s4.Key{
			Address: address,
			SlotId:  2,
			Version: 0,
		}
		record := &s4.Record{
			Payload:    []byte("foobar"),
			Expiration: now.Add(time.Minute).UnixMilli(),
		}
		env := s4.NewEnvelopeFromRecord(key, record)
		signature, err := env.Sign(privateKey)
		assert.NoError(t, err)

		signature[0]++
		err = storage.Put(testutils.Context(t), key, record, signature)
		assert.ErrorIs(t, err, s4.ErrWrongSignature)
	})

	t.Run("ErrVersionTooLow", func(t *testing.T) {
		privateKey, address := testutils.NewPrivateKeyAndAddress(t)
		key := &s4.Key{
			Address: address,
			SlotId:  2,
			Version: 5,
		}
		record := &s4.Record{
			Payload:    []byte("foobar"),
			Expiration: now.Add(time.Hour).UnixMilli(),
		}
		env := s4.NewEnvelopeFromRecord(key, record)
		signature, err := env.Sign(privateKey)
		assert.NoError(t, err)

		ormMock.ExpectedCalls = make([]*mock.Call, 0)
		ormMock.On("Update", mock.Anything, mock.Anything).Return(s4.ErrVersionTooLow).Once()

		err = storage.Put(testutils.Context(t), key, record, signature)
		assert.ErrorIs(t, err, s4.ErrVersionTooLow)
	})
}

func TestStorage_PutAndGet(t *testing.T) {
	t.Parallel()

	now := time.Now()
	ormMock, storage := setupTestStorage(t, now)

	privateKey, address := testutils.NewPrivateKeyAndAddress(t)
	key := &s4.Key{
		Address: address,
		SlotId:  2,
		Version: 0,
	}
	record := &s4.Record{
		Payload:    []byte("foobar"),
		Expiration: now.Add(time.Hour).UnixMilli(),
	}
	env := s4.NewEnvelopeFromRecord(key, record)
	signature, err := env.Sign(privateKey)
	assert.NoError(t, err)

	ormMock.On("Update", mock.Anything, mock.Anything).Return(nil)
	ormMock.On("Get", mock.Anything, big.New(key.Address.Big()), uint(2)).Return(&s4.Row{
		Address:    big.New(key.Address.Big()),
		SlotId:     key.SlotId,
		Version:    key.Version,
		Payload:    record.Payload,
		Expiration: record.Expiration,
		Signature:  signature,
	}, nil)

	err = storage.Put(testutils.Context(t), key, record, signature)
	assert.NoError(t, err)

	rec, metadata, err := storage.Get(testutils.Context(t), key)
	assert.NoError(t, err)
	assert.Equal(t, false, metadata.Confirmed)
	assert.Equal(t, signature, metadata.Signature)
	assert.Equal(t, record.Expiration, rec.Expiration)
	assert.Equal(t, record.Payload, rec.Payload)
}

func TestStorage_List(t *testing.T) {
	t.Parallel()

	ormMock, storage := setupTestStorage(t, time.Now())
	address := testutils.NewAddress()
	ormRows := []*s4.SnapshotRow{
		{
			SlotId:     1,
			Version:    1,
			Expiration: 1,
		},
		{
			SlotId:     5,
			Version:    5,
			Expiration: 5,
		},
	}

	addressRange, err := s4.NewSingleAddressRange(big.New(address.Big()))
	assert.NoError(t, err)
	ormMock.On("GetSnapshot", mock.Anything, addressRange).Return(ormRows, nil)

	rows, err := storage.List(testutils.Context(t), address)
	require.NoError(t, err)
	assert.Len(t, rows, 2)
	for _, row := range rows {
		if row.SlotId == ormRows[0].SlotId {
			assert.Equal(t, ormRows[0], row)
		}
		if row.SlotId == ormRows[1].SlotId {
			assert.Equal(t, ormRows[1], row)
		}
	}
}
