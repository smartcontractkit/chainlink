package secp256k1

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestSuite(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	s := NewBlakeKeccackSecp256k1()
	emptyHashAsHex := hex.EncodeToString(s.Hash().Sum(nil))
	require.Equal(t, "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470", emptyHashAsHex)
	_ = s.RandomStream()
}
