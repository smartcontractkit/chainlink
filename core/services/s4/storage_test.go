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

	t.Run("ErrEntryNotFound", func(t *testing.T) {
		_, _, err := storage.Get(testutils.Context(t), testutils.NewAddress(), 1)
		assert.ErrorIs(t, err, s4.ErrEntryNotFound)
	})

	t.Run("ErrTooBigSlotId", func(t *testing.T) {
		_, _, err := storage.Get(testutils.Context(t), testutils.NewAddress(), constraints.MaxSlotsPerUser+1)
		assert.ErrorIs(t, err, s4.ErrTooBigSlotId)

		record := &s4.Record{
			Payload:    make([]byte, 10),
			Version:    0,
			Expiration: time.Now().UnixMilli() + 1,
		}
		err = storage.Put(testutils.Context(t), testutils.NewAddress(), constraints.MaxSlotsPerUser+1, record, []byte{})
		assert.ErrorIs(t, err, s4.ErrTooBigSlotId)
	})

	t.Run("ErrTooBigPayload", func(t *testing.T) {
		record := &s4.Record{
			Payload:    make([]byte, constraints.MaxPayloadSizeBytes+1),
			Version:    0,
			Expiration: time.Now().UnixMilli() + 1,
		}
		err := storage.Put(testutils.Context(t), testutils.NewAddress(), 1, record, []byte{})
		assert.ErrorIs(t, err, s4.ErrTooBigPayload)
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
		slotid := 2
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

	t.Run("ErrRecordExpired", func(t *testing.T) {
		privateKey, _, address := generateCryptoEntity(t)
		slotid := 2
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
		assert.ErrorIs(t, err, s4.ErrRecordExpired)
	})

	t.Run("ErrOlderVersion", func(t *testing.T) {
		privateKey, _, address := generateCryptoEntity(t)
		slotid := 2
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
		assert.ErrorIs(t, err, s4.ErrOlderVersion)
	})
}

func TestStorage_PutAndGet(t *testing.T) {
	t.Parallel()

	storage := setupTestStorage(t)

	t.Run("Happy Put then Get", func(t *testing.T) {
		privateKey, _, address := generateCryptoEntity(t)
		slotid := 2
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
		assert.Equal(t, record.Expiration, metadata.HighestExpiration)
		assert.Equal(t, signature, metadata.Signature)
		assert.Equal(t, record.Version, rec.Version)
		assert.Equal(t, record.Expiration, rec.Expiration)
		assert.Equal(t, record.Payload, rec.Payload)
	})

	t.Run("HighestExpiration", func(t *testing.T) {
		privateKey, _, address := generateCryptoEntity(t)
		slotid := 1
		record2s := &s4.Record{
			Payload:    []byte("two-seconds"),
			Version:    1,
			Expiration: time.Now().Add(2 * time.Second).UnixMilli(),
		}
		env := s4.NewEnvelopeFromRecord(address, slotid, record2s)
		signature, err := env.Sign(privateKey)
		assert.NoError(t, err)

		err = storage.Put(testutils.Context(t), address, slotid, record2s, signature)
		assert.NoError(t, err)

		record1s := &s4.Record{
			Payload:    []byte("one-second"),
			Version:    2,
			Expiration: time.Now().Add(time.Second).UnixMilli(),
		}
		env = s4.NewEnvelopeFromRecord(address, slotid, record1s)
		signature, err = env.Sign(privateKey)
		assert.NoError(t, err)

		err = storage.Put(testutils.Context(t), address, slotid, record1s, signature)
		assert.NoError(t, err)

		_, metadata, err := storage.Get(testutils.Context(t), address, slotid)
		assert.NoError(t, err)
		assert.Equal(t, record2s.Expiration, metadata.HighestExpiration)
	})
}
