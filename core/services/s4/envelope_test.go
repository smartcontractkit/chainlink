package s4_test

import (
	"crypto/ecdsa"
	"encoding/json"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestEnvelope(t *testing.T) {
	t.Parallel()

	payload := testutils.Random32Byte()
	expiration := time.Now().Add(time.Hour).UnixMilli()
	key := &s4.Key{
		Address: testutils.NewAddress(),
		SlotId:  3,
		Version: 5,
	}
	env := s4.NewEnvelopeFromRecord(key, &s4.Record{
		Payload:    payload[:],
		Expiration: expiration,
	})

	t.Run("signing", func(t *testing.T) {
		privateKey, err := crypto.GenerateKey()
		assert.NoError(t, err)

		sig, err := env.Sign(privateKey)
		assert.NoError(t, err)

		addr, err := env.GetSignerAddress(sig)
		assert.NoError(t, err)

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		assert.True(t, ok)
		assert.Equal(t, crypto.PubkeyToAddress(*publicKeyECDSA), addr)
	})

	t.Run("json", func(t *testing.T) {
		js, err := env.ToJson()
		assert.NoError(t, err)

		var decoded s4.Envelope
		err = json.Unmarshal(js, &decoded)
		assert.NoError(t, err)

		js2, err := decoded.ToJson()
		assert.NoError(t, err)
		assert.Equal(t, js, js2)

		assert.Equal(t, *env, decoded)
	})
}
