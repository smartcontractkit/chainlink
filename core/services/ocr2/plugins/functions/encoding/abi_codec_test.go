package encoding_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/encoding"
)

func TestABICodec_EncodeDecodeV1Success(t *testing.T) {
	t.Parallel()
	codec, err := encoding.NewReportCodec(1)
	require.NoError(t, err)

	var report = []*encoding.ProcessedRequest{
		{
			RequestID:           []byte(fmt.Sprintf("%032d", 123)),
			Result:              []byte("abcd"),
			Error:               []byte("err string"),
			CoordinatorContract: []byte("contract_1"),
			OnchainMetadata:     []byte("commitment_1"),
		},
		{
			RequestID:           []byte(fmt.Sprintf("%032d", 4321)),
			Result:              []byte("0xababababab"),
			Error:               []byte(""),
			CoordinatorContract: []byte("contract_2"),
			OnchainMetadata:     []byte("commitment_2"),
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
		require.Equal(t, report[i].CoordinatorContract, decoded[i].CoordinatorContract, "Contracts not equal at index %d", i)
		require.Equal(t, report[i].OnchainMetadata, decoded[i].OnchainMetadata, "Metadata not equal at index %d", i)
	}
}

func TestABICodec_SliceToByte32(t *testing.T) {
	t.Parallel()

	_, err := encoding.SliceToByte32([]byte("abcd"))
	require.Error(t, err)
	_, err = encoding.SliceToByte32([]byte("0123456789012345678901234567890123456789"))
	require.Error(t, err)

	var expected [32]byte
	for i := 0; i < 32; i++ {
		expected[i] = byte(i)
	}
	res, err := encoding.SliceToByte32(expected[:])
	require.NoError(t, err)
	require.Equal(t, expected, res)
}
