package testhelpers

import (
	"math/rand"
	"testing"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// MakeConfigDigest makes config digest
func MakeConfigDigest(t *testing.T) ocrtypes.ConfigDigest {
	t.Helper()
	b := make([]byte, 32)
	/* #nosec G404 */
	_, err := rand.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	return MustBytesToConfigDigest(t, b)
}

// MustBytesToConfigDigest returns config digest from bytes
func MustBytesToConfigDigest(t *testing.T, b []byte) ocrtypes.ConfigDigest {
	t.Helper()
	configDigest, err := ocrtypes.BytesToConfigDigest(b)
	if err != nil {
		t.Fatal(err)
	}
	return configDigest
}
