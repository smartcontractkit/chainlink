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

func setupTestStorage(t *testing.T) s4.Storage {
	logger := logger.TestLogger(t)
	constraints := s4.Constraints{
		MaxSlotsPerUser:     5,
		MaxPayloadSizeBytes: 32,
	}
	storage := s4.NewInMemoryStorage(logger, constraints, time.Second)
	err := storage.Start(testutils.Context(t))
	assert.NoError(t, err)

	t.Cleanup(func() {
		err = storage.Close()
		assert.NoError(t, err)
	})
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

func TestStorage_StartStop(t *testing.T) {
	t.Parallel()

	setupTestStorage(t)
}

func TestStorage_PutAndGet(t *testing.T) {
	t.Parallel()

	storage := setupTestStorage(t)

	t.Run("ErrRecordNotFound", func(t *testing.T) {
		record, metadata, err := storage.Get(testutils.Context(t), testutils.NewAddress(), 2)
		assert.Nil(t, record)
		assert.Nil(t, metadata)
		assert.ErrorIs(t, err, s4.ErrRecordNotFound)
	})

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
		assert.Equal(t, s4.NewRecordState, metadata.State)
		assert.Equal(t, record.Expiration, metadata.HighestExpiration)
		assert.Equal(t, signature, metadata.Signature)
		assert.Equal(t, record.Version, rec.Version)
		assert.Equal(t, record.Expiration, rec.Expiration)
		assert.Equal(t, record.Payload, rec.Payload)
	})
}
