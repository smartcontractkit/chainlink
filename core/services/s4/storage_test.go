package s4_test

import (
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	constraints = s4.Constraints{
		MaxSlotsPerUser:     5,
		MaxPayloadSizeBytes: 32,
	}
)

func setupTestStorage(t *testing.T) (*mocks.ORM, s4.Storage) {
	logger := logger.TestLogger(t)
	orm := mocks.NewORM(t)
	storage := s4.NewStorage(logger, constraints, orm)
	return orm, storage
}

func generateCryptoEntity(t *testing.T) (*ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address) {
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	assert.True(t, ok)

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, publicKeyECDSA, address
}

func TestStorage_Constraints(t *testing.T) {
	t.Parallel()

	_, storage := setupTestStorage(t)
	c := storage.Constraints()
	assert.Equal(t, constraints, c)
}

func TestStorage_Errors(t *testing.T) {
	t.Parallel()

	ormMock, storage := setupTestStorage(t)

	t.Run("ErrNotFound", func(t *testing.T) {
		key := &s4.Key{
			Address: testutils.NewAddress(),
			SlotId:  1,
			Version: 0,
		}
		ormMock.On("Get", utils.NewBig(key.Address.Big()), uint(key.SlotId), mock.Anything).Return(nil, s4.ErrNotFound)
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
			Expiration: time.Now().UnixMilli() + 1,
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
			Expiration: time.Now().UnixMilli() + 1,
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
			Expiration: time.Now().UnixMilli() - 1,
		}
		err := storage.Put(testutils.Context(t), key, record, []byte{})
		assert.ErrorIs(t, err, s4.ErrPastExpiration)
	})

	t.Run("ErrWrongSignature", func(t *testing.T) {
		privateKey, _, address := generateCryptoEntity(t)
		key := &s4.Key{
			Address: address,
			SlotId:  2,
			Version: 0,
		}
		record := &s4.Record{
			Payload:    []byte("foobar"),
			Expiration: time.Now().UnixMilli() + 1,
		}
		env := s4.NewEnvelopeFromRecord(key, record)
		signature, err := env.Sign(privateKey)
		assert.NoError(t, err)

		signature[0]++
		err = storage.Put(testutils.Context(t), key, record, signature)
		assert.ErrorIs(t, err, s4.ErrWrongSignature)
	})

	/* to be reworked
	t.Run("ErrNotFound if expired", func(t *testing.T) {
		privateKey, _, address := generateCryptoEntity(t)
		key := &s4.Key{
			Address: address,
			SlotId:  2,
			Version: 0,
		}
		record := &s4.Record{
			Payload:    []byte("foobar"),
			Expiration: time.Now().UnixMilli() + 1,
		}
		env := s4.NewEnvelopeFromRecord(key, record)
		signature, err := env.Sign(privateKey)
		assert.NoError(t, err)

		ormMock.On("Update", mock.Anything, mock.Anything).Return(nil)
		ormMock.On("Get", utils.NewBig(key.Address.Big()), uint(key.SlotId), mock.Anything).Return(nil, s4.ErrNotFound)

		err = storage.Put(testutils.Context(t), key, record, signature)
		assert.NoError(t, err)

		_, _, err = storage.Get(testutils.Context(t), key)
		assert.ErrorIs(t, err, s4.ErrNotFound)
	})
	*/

	t.Run("ErrVersionTooLow", func(t *testing.T) {
		privateKey, _, address := generateCryptoEntity(t)
		key := &s4.Key{
			Address: address,
			SlotId:  2,
			Version: 5,
		}
		record := &s4.Record{
			Payload:    []byte("foobar"),
			Expiration: time.Now().Add(time.Hour).UnixMilli(),
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

	ormMock, storage := setupTestStorage(t)

	t.Run("Happy Put then Get", func(t *testing.T) {
		privateKey, _, address := generateCryptoEntity(t)
		key := &s4.Key{
			Address: address,
			SlotId:  2,
			Version: 0,
		}
		record := &s4.Record{
			Payload:    []byte("foobar"),
			Expiration: time.Now().Add(time.Hour).UnixMilli(),
		}
		env := s4.NewEnvelopeFromRecord(key, record)
		signature, err := env.Sign(privateKey)
		assert.NoError(t, err)

		ormMock.On("Update", mock.Anything, mock.Anything).Return(nil)
		ormMock.On("Get", utils.NewBig(key.Address.Big()), uint(2), mock.Anything).Return(&s4.Row{
			Address:    utils.NewBig(key.Address.Big()),
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
	})
}
