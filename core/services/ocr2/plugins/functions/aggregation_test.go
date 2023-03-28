package functions_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
)

func req(id int, result []byte, err []byte) *functions.ProcessedRequest {
	return &functions.ProcessedRequest{
		RequestID: []byte(strconv.Itoa(id)),
		Result:    result,
		Error:     err,
	}
}

func reqS(id int, result string, err string) *functions.ProcessedRequest {
	return req(id, []byte(result), []byte(err))
}

func TestCanAggregate(t *testing.T) {
	t.Parallel()
	obs := make([]*functions.ProcessedRequest, 10)

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
		input    []*functions.ProcessedRequest
		expected *functions.ProcessedRequest
	}{
		{
			"Mode",
			config.AggregationMethod_AGGREGATION_MODE,
			[]*functions.ProcessedRequest{
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
			[]*functions.ProcessedRequest{
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
			[]*functions.ProcessedRequest{
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
			[]*functions.ProcessedRequest{
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
			[]*functions.ProcessedRequest{
				req(21, []byte{0, 9, 200, 2}, []byte{}),
				req(21, []byte{0, 0, 9, 11}, []byte{}),
				req(21, []byte{0, 0, 5, 100}, []byte{}),
				req(21, []byte{0, 0, 12, 2}, []byte{}),
			},
			req(21, []byte{0, 0, 9, 11}, []byte{}),
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
