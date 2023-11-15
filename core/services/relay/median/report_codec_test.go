package median_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	libocr2median "github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/median"
)

func TestReportCodec(t *testing.T) {
	anyReports := []libocr2median.ParsedAttributedObservation{
		{
			Timestamp:       123,
			Value:           big.NewInt(300),
			JuelsPerFeeCoin: big.NewInt(100),
			Observer:        0,
		},
		{
			Timestamp:       125,
			Value:           big.NewInt(200),
			JuelsPerFeeCoin: big.NewInt(110),
			Observer:        1,
		},
		{
			Timestamp:       124,
			Value:           big.NewInt(250),
			JuelsPerFeeCoin: big.NewInt(90),
			Observer:        2,
		},
	}

	aggReports := median.AggregatedAttributedObservation{
		Timestamp: 124,
		Observers: [32]commontypes.OracleID{1, 2, 0},
		Observations: []*big.Int{
			big.NewInt(200),
			big.NewInt(250),
			big.NewInt(300),
		},
		JuelsPerFeeCoin: big.NewInt(100),
	}

	anyEncodedReport := ocrtypes.Report{5, 6, 7, 8}

	t.Run("BuildReport builds the type and delegates to relay", func(t *testing.T) {
		reportCodec, err := median.NewReportCodec(&testCodec{
			t:        t,
			expected: &aggReports,
			result:   anyEncodedReport,
		})
		require.NoError(t, err)

		encoded, err := reportCodec.BuildReport(anyReports)
		require.NoError(t, err)
		assert.Equal(t, anyEncodedReport, encoded)
	})

	t.Run("MedianFromReport delegates to relay and gets the median", func(t *testing.T) {
		reportCodec, err := median.NewReportCodec(&testCodec{
			t:        t,
			expected: []uint8(anyEncodedReport),
			result:   aggReports,
		})
		require.NoError(t, err)

		medianVal, err := reportCodec.MedianFromReport(anyEncodedReport)
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(250), medianVal)
	})

	t.Run("MaxReportLength delegates to relay", func(t *testing.T) {
		anyN := 10
		anyLen := 200
		reportCodec, err := median.NewReportCodec(&testCodec{
			t:        t,
			expected: anyN,
			result:   anyLen,
		})
		require.NoError(t, err)

		length, err := reportCodec.MaxReportLength(anyN)
		require.NoError(t, err)
		assert.Equal(t, anyLen, length)
	})
}

type testCodec struct {
	t        *testing.T
	expected any
	result   any
}

func (t testCodec) Encode(_ context.Context, item any, itemType string) (ocrtypes.Report, error) {
	assert.Equal(t.t, t.expected, item)
	assert.Equal(t.t, median.MedianTypeName, itemType)
	return t.result.(ocrtypes.Report), nil
}

func (t testCodec) GetMaxEncodingSize(_ context.Context, n int, itemType string) (int, error) {
	assert.Equal(t.t, t.expected, n)
	assert.Equal(t.t, median.MedianTypeName, itemType)
	return t.result.(int), nil
}

func (t testCodec) Decode(_ context.Context, raw []byte, into any, itemType string) error {
	assert.Equal(t.t, t.expected, raw)
	assert.Equal(t.t, median.MedianTypeName, itemType)
	set := into.(*median.AggregatedAttributedObservation)
	*set = t.result.(median.AggregatedAttributedObservation)
	return nil
}

func (t testCodec) GetMaxDecodingSize(_ context.Context, n int, itemType string) (int, error) {
	assert.Equal(t.t, t.expected, n)
	assert.Equal(t.t, median.MedianTypeName, itemType)
	return t.result.(int), nil
}
