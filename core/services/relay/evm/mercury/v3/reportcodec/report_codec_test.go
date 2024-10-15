package reportcodec

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func newValidReportFields() v3.ReportFields {
	return v3.ReportFields{
		Timestamp:          242,
		BenchmarkPrice:     big.NewInt(243),
		Bid:                big.NewInt(244),
		Ask:                big.NewInt(245),
		ValidFromTimestamp: 123,
		ExpiresAt:          20,
		LinkFee:            big.NewInt(456),
		NativeFee:          big.NewInt(457),
	}
}

func Test_ReportCodec_BuildReport(t *testing.T) {
	r := ReportCodec{}

	t.Run("BuildReport errors on zero values", func(t *testing.T) {
		ctx := testutils.Context(t)
		_, err := r.BuildReport(ctx, v3.ReportFields{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "benchmarkPrice may not be nil")
		assert.Contains(t, err.Error(), "linkFee may not be nil")
		assert.Contains(t, err.Error(), "nativeFee may not be nil")
	})

	t.Run("BuildReport constructs a report from observations", func(t *testing.T) {
		rf := newValidReportFields()
		// only need to test happy path since validations are done in relaymercury

		ctx := testutils.Context(t)
		report, err := r.BuildReport(ctx, rf)
		require.NoError(t, err)

		reportElems := make(map[string]interface{})
		err = ReportTypes.UnpackIntoMap(reportElems, report)
		require.NoError(t, err)

		assert.Equal(t, int(reportElems["observationsTimestamp"].(uint32)), 242)
		assert.Equal(t, reportElems["benchmarkPrice"].(*big.Int).Int64(), int64(243))
		assert.Equal(t, reportElems["bid"].(*big.Int).Int64(), int64(244))
		assert.Equal(t, reportElems["ask"].(*big.Int).Int64(), int64(245))
		assert.Equal(t, reportElems["validFromTimestamp"].(uint32), uint32(123))
		assert.Equal(t, reportElems["expiresAt"].(uint32), uint32(20))
		assert.Equal(t, reportElems["linkFee"].(*big.Int).Int64(), int64(456))
		assert.Equal(t, reportElems["nativeFee"].(*big.Int).Int64(), int64(457))

		assert.Equal(t, types.Report{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0xc9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0xc8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf5}, report)
		max, err := r.MaxReportLength(ctx, 4)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(report), max)

		t.Run("Decode decodes the report", func(t *testing.T) {
			decoded, err := r.Decode(report)
			require.NoError(t, err)

			require.NotNil(t, decoded)

			assert.Equal(t, uint32(242), decoded.ObservationsTimestamp)
			assert.Equal(t, big.NewInt(243), decoded.BenchmarkPrice)
			assert.Equal(t, big.NewInt(244), decoded.Bid)
			assert.Equal(t, big.NewInt(245), decoded.Ask)
			assert.Equal(t, uint32(123), decoded.ValidFromTimestamp)
			assert.Equal(t, uint32(20), decoded.ExpiresAt)
			assert.Equal(t, big.NewInt(456), decoded.LinkFee)
			assert.Equal(t, big.NewInt(457), decoded.NativeFee)
		})
	})

	t.Run("errors on negative fee", func(t *testing.T) {
		rf := newValidReportFields()
		rf.LinkFee = big.NewInt(-1)
		rf.NativeFee = big.NewInt(-1)
		ctx := testutils.Context(t)
		_, err := r.BuildReport(ctx, rf)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "linkFee may not be negative (got: -1)")
		assert.Contains(t, err.Error(), "nativeFee may not be negative (got: -1)")
	})

	t.Run("Decode errors on invalid report", func(t *testing.T) {
		_, err := r.Decode([]byte{1, 2, 3})
		assert.EqualError(t, err, "failed to decode report: abi: cannot marshal in to go type: length insufficient 3 require 32")

		longBad := make([]byte, 64)
		for i := 0; i < len(longBad); i++ {
			longBad[i] = byte(i)
		}
		_, err = r.Decode(longBad)
		assert.EqualError(t, err, "failed to decode report: abi: improperly encoded uint32 value")
	})
}

func buildSampleReport(ts int64) []byte {
	feedID := [32]byte{'f', 'o', 'o'}
	timestamp := uint32(ts)
	bp := big.NewInt(242)
	bid := big.NewInt(243)
	ask := big.NewInt(244)
	validFromTimestamp := uint32(123)
	expiresAt := uint32(456)
	linkFee := big.NewInt(3334455)
	nativeFee := big.NewInt(556677)

	b, err := ReportTypes.Pack(feedID, validFromTimestamp, timestamp, nativeFee, linkFee, expiresAt, bp, bid, ask)
	if err != nil {
		panic(err)
	}
	return b
}

func Test_ReportCodec_ObservationTimestampFromReport(t *testing.T) {
	r := ReportCodec{}

	t.Run("ObservationTimestampFromReport extracts observation timestamp from a valid report", func(t *testing.T) {
		report := buildSampleReport(123)

		ctx := testutils.Context(t)
		ts, err := r.ObservationTimestampFromReport(ctx, report)
		require.NoError(t, err)

		assert.Equal(t, ts, uint32(123))
	})
	t.Run("ObservationTimestampFromReport returns error when report is invalid", func(t *testing.T) {
		report := []byte{1, 2, 3}

		ctx := testutils.Context(t)
		_, err := r.ObservationTimestampFromReport(ctx, report)
		require.Error(t, err)

		assert.EqualError(t, err, "failed to decode report: abi: cannot marshal in to go type: length insufficient 3 require 32")
	})
}

func Test_ReportCodec_BenchmarkPriceFromReport(t *testing.T) {
	r := ReportCodec{}

	t.Run("BenchmarkPriceFromReport extracts the benchmark price from valid report", func(t *testing.T) {
		ctx := testutils.Context(t)
		report := buildSampleReport(123)

		bp, err := r.BenchmarkPriceFromReport(ctx, report)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(242), bp)
	})
	t.Run("BenchmarkPriceFromReport errors on invalid report", func(t *testing.T) {
		ctx := testutils.Context(t)
		_, err := r.BenchmarkPriceFromReport(ctx, []byte{1, 2, 3})
		require.Error(t, err)
		assert.EqualError(t, err, "failed to decode report: abi: cannot marshal in to go type: length insufficient 3 require 32")
	})
}
