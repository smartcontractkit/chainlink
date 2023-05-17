package s4_test

import (
	"crypto/ecdsa"
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
	env := s4.NewEnvelopeFromRecord(testutils.NewAddress(), 3, &s4.Record{
		Payload:    payload[:],
		Version:    5,
		Expiration: expiration,
	})
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
}
