package s4_test

import (
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	constraints = s4.Constraints{
		MaxSlotsPerUser:     5,
		MaxPayloadSizeBytes: 32,
	}
)

func setupTestStorage(t *testing.T) s4.Storage {
	logger := logger.TestLogger(t)
	orm := s4.NewInMemoryORM()
	storage := s4.NewStorage(logger, constraints, orm)
	return storage
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

	storage := setupTestStorage(t)
	c := storage.Constraints()
	assert.Equal(t, constraints, c)
}

func TestStorage_Errors(t *testing.T) {
	t.Parallel()

	storage := setupTestStorage(t)

	t.Run("ErrNotFound", func(t *testing.T) {
		_, _, err := storage.Get(testutils.Context(t), testutils.NewAddress(), 1)
		assert.ErrorIs(t, err, s4.ErrNotFound)
	})

	t.Run("ErrSlotIdTooBig", func(t *testing.T) {
		_, _, err := storage.Get(testutils.Context(t), testutils.NewAddress(), constraints.MaxSlotsPerUser+1)
		assert.ErrorIs(t, err, s4.ErrSlotIdTooBig)

		record := &s4.Record{
			Payload:    make([]byte, 10),
			Version:    0,
			Expiration: time.Now().UnixMilli() + 1,
		}
		err = storage.Put(testutils.Context(t), testutils.NewAddress(), constraints.MaxSlotsPerUser+1, record, []byte{})
		assert.ErrorIs(t, err, s4.ErrSlotIdTooBig)
	})

	t.Run("ErrPayloadTooBig", func(t *testing.T) {
		record := &s4.Record{
			Payload:    make([]byte, constraints.MaxPayloadSizeBytes+1),
			Version:    0,
			Expiration: time.Now().UnixMilli() + 1,
		}
		err := storage.Put(testutils.Context(t), testutils.NewAddress(), 1, record, []byte{})
		assert.ErrorIs(t, err, s4.ErrPayloadTooBig)
	})

	t.Run("ErrPastExpiration", func(t *testing.T) {
		record := &s4.Record{
			Payload:    make([]byte, 10),
			Version:    0,
			Expiration: time.Now().UnixMilli() - 1,
		}
		err := storage.Put(testutils.Context(t), testutils.NewAddress(), 1, record, []byte{})
		assert.ErrorIs(t, err, s4.ErrPastExpiration)
	})

	t.Run("ErrWrongSignature", func(t *testing.T) {
		privateKey, _, address := generateCryptoEntity(t)
		var slotid uint = 2
		record := &s4.Record{
			Payload:    []byte("foobar"),
			Version:    0,
			Expiration: time.Now().UnixMilli() + 1,
		}
		env := s4.NewEnvelopeFromRecord(address, slotid, record)
		signature, err := env.Sign(privateKey)
		assert.NoError(t, err)

		signature[0]++
		err = storage.Put(testutils.Context(t), address, slotid, record, signature)
		assert.ErrorIs(t, err, s4.ErrWrongSignature)
	})

	t.Run("ErrNotFound if expired", func(t *testing.T) {
		privateKey, _, address := generateCryptoEntity(t)
		var slotid uint = 2
		record := &s4.Record{
			Payload:    []byte("foobar"),
			Version:    0,
			Expiration: time.Now().UnixMilli() + 1,
		}
		env := s4.NewEnvelopeFromRecord(address, slotid, record)
		signature, err := env.Sign(privateKey)
		assert.NoError(t, err)

		err = storage.Put(testutils.Context(t), address, slotid, record, signature)
		assert.NoError(t, err)

		time.Sleep(testutils.TestInterval)
		_, _, err = storage.Get(testutils.Context(t), address, slotid)
		assert.ErrorIs(t, err, s4.ErrNotFound)
	})

	t.Run("ErrVersionTooLow", func(t *testing.T) {
		privateKey, _, address := generateCryptoEntity(t)
		var slotid uint = 2
		record := &s4.Record{
			Payload:    []byte("foobar"),
			Version:    5,
			Expiration: time.Now().Add(time.Hour).UnixMilli(),
		}
		env := s4.NewEnvelopeFromRecord(address, slotid, record)
		signature, err := env.Sign(privateKey)
		assert.NoError(t, err)

		err = storage.Put(testutils.Context(t), address, slotid, record, signature)
		assert.NoError(t, err)

		record.Version--
		env = s4.NewEnvelopeFromRecord(address, slotid, record)
		signature, err = env.Sign(privateKey)
		assert.NoError(t, err)

		err = storage.Put(testutils.Context(t), address, slotid, record, signature)
		assert.ErrorIs(t, err, s4.ErrVersionTooLow)
	})
}

func TestStorage_PutAndGet(t *testing.T) {
	t.Parallel()

	storage := setupTestStorage(t)

	t.Run("Happy Put then Get", func(t *testing.T) {
		privateKey, _, address := generateCryptoEntity(t)
		var slotid uint = 2
		record := &s4.Record{
			Payload:    []byte("foobar"),
			Version:    0,
			Expiration: time.Now().Add(time.Hour).UnixMilli(),
		}
		env := s4.NewEnvelopeFromRecord(address, slotid, record)
		signature, err := env.Sign(privateKey)
		assert.NoError(t, err)

		err = storage.Put(testutils.Context(t), address, slotid, record, signature)
		assert.NoError(t, err)

		rec, metadata, err := storage.Get(testutils.Context(t), address, slotid)
		assert.NoError(t, err)
		assert.Equal(t, false, metadata.Confirmed)
		assert.Equal(t, signature, metadata.Signature)
		assert.Equal(t, record.Version, rec.Version)
		assert.Equal(t, record.Expiration, rec.Expiration)
		assert.Equal(t, record.Payload, rec.Payload)
	})
}
