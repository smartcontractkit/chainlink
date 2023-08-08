package mercury_v2

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercuryv2 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v2"
)

var paos = []relaymercuryv2.ParsedAttributedObservation{
	relaymercuryv2.NewParsedAttributedObservation(42, commontypes.OracleID(49), big.NewInt(43), true, 123, true, big.NewInt(143), true, big.NewInt(456), true),
	relaymercuryv2.NewParsedAttributedObservation(142, commontypes.OracleID(149), big.NewInt(143), true, 456, true, big.NewInt(456), true, big.NewInt(345), true),
	relaymercuryv2.NewParsedAttributedObservation(242, commontypes.OracleID(249), big.NewInt(243), true, 789, true, big.NewInt(764), true, big.NewInt(167), true),
	relaymercuryv2.NewParsedAttributedObservation(342, commontypes.OracleID(250), big.NewInt(343), true, 123, true, big.NewInt(378), true, big.NewInt(643), true),
}

func Test_ReportCodec_BuildReport(t *testing.T) {
	r := ReportCodec{}
	f := 1

	t.Run("BuildReport errors if observations are empty", func(t *testing.T) {
		paos := []relaymercuryv2.ParsedAttributedObservation{}
		_, err := r.BuildReport(paos, f, 123, 10)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot build report from empty attributed observation")
	})

	t.Run("BuildReport constructs a report from observations", func(t *testing.T) {
		// only need to test happy path since validations are done in relaymercury

		report, err := r.BuildReport(paos, f, 123, 20)
		require.NoError(t, err)

		reportElems := make(map[string]interface{})
		err = ReportTypes.UnpackIntoMap(reportElems, report)
		require.NoError(t, err)

		assert.Equal(t, int(reportElems["observationsTimestamp"].(uint32)), 242)
		assert.Equal(t, reportElems["benchmarkPrice"].(*big.Int).Int64(), int64(243))
		assert.Equal(t, reportElems["validFromTimestamp"].(uint32), uint32(123))
		assert.Equal(t, reportElems["expiresAt"].(uint32), uint32(20))
		assert.Equal(t, reportElems["linkFee"].(*big.Int).Int64(), int64(456))
		assert.Equal(t, reportElems["nativeFee"].(*big.Int).Int64(), int64(456))

		assert.Equal(t, types.Report{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0xc8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0xc8}, report)
		max, err := r.MaxReportLength(4)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(report), max)
	})
}

func buildSampleReport(ts int64) []byte {
	feedID := [32]byte{'f', 'o', 'o'}
	timestamp := uint32(ts)
	bp := big.NewInt(242)
	validFromTimestamp := uint32(123)
	expiresAt := uint32(456)
	linkFee := big.NewInt(3334455)
	nativeFee := big.NewInt(556677)

	b, err := ReportTypes.Pack(feedID, timestamp, bp, validFromTimestamp, expiresAt, linkFee, nativeFee)
	if err != nil {
		panic(err)
	}
	return b
}

func Test_ReportCodec_ObservationTimestampFromReport(t *testing.T) {
	r := ReportCodec{}

	t.Run("ObservationTimestampFromReport extracts observation timestamp from a valid report", func(t *testing.T) {
		report := buildSampleReport(123)

		ts, err := r.ObservationTimestampFromReport(report)
		require.NoError(t, err)

		assert.Equal(t, ts, uint32(123))
	})
	t.Run("ObservationTimestampFromReport returns error when timestamp is too big", func(t *testing.T) {
		report := buildSampleReport(math.MaxInt32 + 1)

		_, err := r.ObservationTimestampFromReport(report)
		require.Error(t, err)

		assert.EqualError(t, err, "timestamp overflows max uint32, got: 2147483648")
	})
}
