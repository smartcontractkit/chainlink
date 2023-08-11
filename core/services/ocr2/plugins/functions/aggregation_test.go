package functions_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/encoding"
)

func req(id int, result []byte, err []byte) *encoding.ProcessedRequest {
	return &encoding.ProcessedRequest{
		RequestID:           []byte(strconv.Itoa(id)),
		Result:              result,
		Error:               err,
		CallbackGasLimit:    0,
		CoordinatorContract: []byte{},
		OnchainMetadata:     []byte{},
	}
}

func reqS(id int, result string, err string) *encoding.ProcessedRequest {
	return req(id, []byte(result), []byte(err))
}

func reqMeta(id int, result []byte, err []byte, callbackGas uint32, coordinatorContract []byte, onchainMeta []byte) *encoding.ProcessedRequest {
	return &encoding.ProcessedRequest{
		RequestID:           []byte(strconv.Itoa(id)),
		Result:              result,
		Error:               err,
		CallbackGasLimit:    callbackGas,
		CoordinatorContract: coordinatorContract,
		OnchainMetadata:     onchainMeta,
	}
}

func TestCanAggregate(t *testing.T) {
	t.Parallel()
	obs := make([]*encoding.ProcessedRequest, 10)

	require.True(t, functions.CanAggregate(4, 1, obs[:4]))
	require.True(t, functions.CanAggregate(4, 1, obs[:3]))
	require.True(t, functions.CanAggregate(6, 1, obs[:3]))

	require.False(t, functions.CanAggregate(4, 1, obs[:5]))
	require.False(t, functions.CanAggregate(4, 1, obs[:2]))
	require.False(t, functions.CanAggregate(4, 1, obs[:0]))
	require.False(t, functions.CanAggregate(0, 0, obs[:0]))
}

func TestAggregate_Successful(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mode     config.AggregationMethod
		input    []*encoding.ProcessedRequest
		expected *encoding.ProcessedRequest
	}{
		{
			"Mode",
			config.AggregationMethod_AGGREGATION_MODE,
			[]*encoding.ProcessedRequest{
				reqS(21, "ab", ""),
				reqS(21, "abcd", ""),
				reqS(21, "cd", ""),
				reqS(21, "abcd", ""),
			},
			reqS(21, "abcd", ""),
		},
		{
			"Errors",
			config.AggregationMethod_AGGREGATION_MEDIAN,
			[]*encoding.ProcessedRequest{
				reqS(21, "", "bug"),
				reqS(21, "", "compile error"),
				reqS(21, "", "bug"),
			},
			reqS(21, "", "bug"),
		},
		{
			"Median Odd",
			config.AggregationMethod_AGGREGATION_MEDIAN,
			// NOTE: binary values of those strings represent different integers
			// but they still should be sorted correctly
			[]*encoding.ProcessedRequest{
				reqS(21, "7", ""),
				reqS(21, "101", ""),
				reqS(21, "8", ""),
				reqS(21, "19", ""),
				reqS(21, "10", ""),
			},
			reqS(21, "10", ""),
		},
		{
			"Median Even",
			config.AggregationMethod_AGGREGATION_MEDIAN,
			[]*encoding.ProcessedRequest{
				req(21, []byte{9, 200, 2}, []byte{}),
				req(21, []byte{9, 11}, []byte{}),
				req(21, []byte{5, 100}, []byte{}),
				req(21, []byte{12, 2}, []byte{}),
			},
			req(21, []byte{9, 11}, []byte{}),
		},
		{
			"Median Even Aligned",
			config.AggregationMethod_AGGREGATION_MEDIAN,
			[]*encoding.ProcessedRequest{
				req(21, []byte{0, 9, 200, 2}, []byte{}),
				req(21, []byte{0, 0, 9, 11}, []byte{}),
				req(21, []byte{0, 0, 5, 100}, []byte{}),
				req(21, []byte{0, 0, 12, 2}, []byte{}),
			},
			req(21, []byte{0, 0, 9, 11}, []byte{}),
		},
		{
			"Metadata With Results",
			config.AggregationMethod_AGGREGATION_MEDIAN,
			[]*encoding.ProcessedRequest{
				reqMeta(21, []byte{1}, []byte{}, 100, []byte{2}, []byte{4}),
				reqMeta(21, []byte{1}, []byte{}, 90, []byte{2}, []byte{4}),
				reqMeta(21, []byte{1}, []byte{}, 100, []byte{0}, []byte{4}),
				reqMeta(21, []byte{1}, []byte{}, 100, []byte{2}, []byte{1}),
			},
			reqMeta(21, []byte{1}, []byte{}, 100, []byte{2}, []byte{4}),
		},
		{
			"Metadata With Errors",
			config.AggregationMethod_AGGREGATION_MEDIAN,
			[]*encoding.ProcessedRequest{
				reqMeta(21, []byte{}, []byte{2}, 90, []byte{0}, []byte{4}),
				reqMeta(21, []byte{}, []byte{2}, 100, []byte{2}, []byte{4}),
				reqMeta(21, []byte{}, []byte{2}, 100, []byte{2}, []byte{1}),
				reqMeta(21, []byte{}, []byte{2}, 100, []byte{2}, []byte{4}),
			},
			reqMeta(21, []byte{}, []byte{2}, 100, []byte{2}, []byte{4}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := functions.Aggregate(test.mode, test.input)
			require.NoError(t, err)
			require.Equal(t, test.expected, result)
		})
	}
}
