package suikey

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

func TestSuiKey(t *testing.T) {
	bytes, err := hex.DecodeString("f0d07ab448018b2754475f9a3b580218b0675a1456aad96ad607c7bbd7d9237b")
	require.NoError(t, err)
	k := KeyFor(internal.NewRaw(bytes))
	assert.Equal(t, "2acd605efc181e2af8a0b8c0686a5e12578efa1253d15a235fa5e5ad970c4b29", k.PublicKeyStr())
	assert.Equal(t, "28ebc047741ac19f2ff459d4ddb8f0c688415650edb103a22e4c66350dbcaff9", k.Account())
}
