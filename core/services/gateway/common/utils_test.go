package common_test

import (
	"testing"
	"time"

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

	val = "世界"
	data = common.StringToAlignedBytes(val, 40)
	require.Equal(t, val, common.AlignedBytesToString(data))
}
