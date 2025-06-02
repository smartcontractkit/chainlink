package common_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
)

func TestUtils_BytesUint32Conversions(t *testing.T) {
	t.Parallel()

	val := uint32(time.Now().Unix())
	data := common.Uint32ToBytes(val)
	require.Equal(t, val, common.BytesToUint32(data))
}

func TestUtils_StringAlignedBytesConversions(t *testing.T) {
	t.Parallel()

	val := "my_string"
	data := common.StringToAlignedBytes(val, 40)
	require.Equal(t, val, common.AlignedBytesToString(data))

	val = "0123456789"
	data = common.StringToAlignedBytes(val, 10)
	require.Equal(t, val, common.AlignedBytesToString(data))

	val = "世界"
	data = common.StringToAlignedBytes(val, 40)
	require.Equal(t, val, common.AlignedBytesToString(data))
}

func TestUtils_BytesSignAndValidate(t *testing.T) {
	t.Parallel()

	data := []byte("data_data")
	incorrectData := []byte("some_other_data")

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Bytes()

	signature, err := common.SignData(privateKey, data)
	require.NoError(t, err)
	require.Len(t, signature, 65)

	// valid
	signer, err := common.ExtractSigner(signature, data)
	require.NoError(t, err)
	require.True(t, bytes.Equal(signer, address))

	// invalid
	signer, err = common.ExtractSigner(signature, incorrectData)
	require.NoError(t, err)
	require.False(t, bytes.Equal(signer, address))

	// invalid format
	_, err = common.ExtractSigner([]byte{0xaa, 0xbb}, data)
	require.Error(t, err)
}
