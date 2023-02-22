package ocr2key

import (
	"crypto/ed25519"
	"testing"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
)

func TestLegacyEdDSAKeyring(t *testing.T) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	assert.NoError(t, err)
	lk := legacyEdDSAKeyring{privKey}

	assert.EqualValues(t, pubKey, lk.PublicKey())

	_, err = lk.Sign(ocrtypes.ReportContext{}, nil)
	assert.ErrorContains(t, err, "cannot use a legacy key to sign")

	verify := lk.Verify(nil, ocrtypes.ReportContext{}, nil, nil)
	assert.False(t, verify)

	maxSigLength := lk.MaxSignatureLength()
	assert.Equal(t, maxSigLength, ed25519.PublicKeySize+ed25519.SignatureSize)

	seed, err := lk.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, seed, privKey.Seed())

	err = lk.Unmarshal([]byte{0, 0})
	assert.ErrorContains(t, err, "unexpected seed size")

	_lk := legacyEdDSAKeyring{}
	err = _lk.Unmarshal(seed)
	assert.Equal(t, _lk, lk)
}
