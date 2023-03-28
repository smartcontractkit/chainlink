package functions_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions"
)

func TestDRCodec_EncodeDecodeSuccess(t *testing.T) {
	t.Parallel()
	codec, err := functions.NewReportCodec()
	require.NoError(t, err)

	var report = []*functions.ProcessedRequest{
		{
			RequestID: []byte(fmt.Sprintf("%032d", 123)),
			Result:    []byte("abcd"),
			Error:     []byte("err string"),
		},
		{
			RequestID: []byte(fmt.Sprintf("%032d", 4321)),
			Result:    []byte("0xababababab"),
			Error:     []byte(""),
		},
	}

	encoded, err := codec.EncodeReport(report)
	require.NoError(t, err)
	decoded, err := codec.DecodeReport(encoded)
	require.NoError(t, err)

	require.Equal(t, len(report), len(decoded))
	for i := 0; i < len(report); i++ {
		require.Equal(t, report[i].RequestID, decoded[i].RequestID, "RequestIDs not equal at index %d", i)
		require.Equal(t, report[i].Result, decoded[i].Result, "Results not equal at index %d", i)
		require.Equal(t, report[i].Error, decoded[i].Error, "Errors not equal at index %d", i)
	}
}
