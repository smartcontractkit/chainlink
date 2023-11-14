package tokendata_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
)

// Black box test
func TestCachedReader_ReadTokenData(t *testing.T) {
	mockReader := tokendata.MockReader{}
	cachedReader := tokendata.NewCachedReader(&mockReader)

	msgData := []byte("msgData")
	mockReader.On("ReadTokenData", mock.Anything, mock.Anything).Return(msgData, nil)

	ctx := context.Background()
	msg := internal.EVM2EVMOnRampCCIPSendRequestedWithMeta{}

	// Call ReadTokenData twice, expect only one call to underlying reader
	data, err := cachedReader.ReadTokenData(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, msgData, data)

	// First time, calls the underlying reader
	mockReader.AssertNumberOfCalls(t, "ReadTokenData", 1)

	data, err = cachedReader.ReadTokenData(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, msgData, data)

	// Second time, get data from cache
	mockReader.AssertNumberOfCalls(t, "ReadTokenData", 1)
}
