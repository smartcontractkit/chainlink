package vrfkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVRFKeys_Models(t *testing.T) {
	kv2, err := NewV2()
	require.NoError(t, err)
	k := EncryptedVRFKey{
		PublicKey: kv2.PublicKey,
		VRFKey:    gethKeyStruct{},
	}

	v, err := k.JSON()
	assert.NotEmpty(t, v)
	require.NoError(t, err)

	vrfk := gethKeyStruct{}

	v2, err := vrfk.Value()
	require.NoError(t, err)
	err = vrfk.Scan(v2)
	require.NoError(t, err)
	err = vrfk.Scan("")
	require.Error(t, err)
	err = vrfk.Scan(1)
	require.Error(t, err)
}
