package reportcodec

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	v1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

var hash = hexutil.MustDecode("0x552c2cea3ab43bae137d89ee6142a01db3ae2b5678bc3c9bd5f509f537bea57b")

func newValidReportFields() v1.ReportFields {
	return v1.ReportFields{
		Timestamp:             242,
		BenchmarkPrice:        big.NewInt(243),
		Bid:                   big.NewInt(244),
		Ask:                   big.NewInt(245),
		CurrentBlockNum:       248,
		CurrentBlockHash:      hash,
		ValidFromBlockNum:     46,
		CurrentBlockTimestamp: 123,
	}
}

func Test_ReportCodec(t *testing.T) {
	r := ReportCodec{}

	t.Run("BuildReport errors on zero fields", func(t *testing.T) {
		_, err := r.BuildReport(v1.ReportFields{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "benchmarkPrice may not be nil")
		assert.Contains(t, err.Error(), "bid may not be nil")
		assert.Contains(t, err.Error(), "ask may not be nil")
		assert.Contains(t, err.Error(), "invalid length for currentBlockHash, expected: 32, got: 0")
	})

	t.Run("BuildReport constructs a report from observations", func(t *testing.T) {
		rf := newValidReportFields()
		// only need to test happy path since validations are done in relaymercury

		report, err := r.BuildReport(rf)
		require.NoError(t, err)

		reportElems := make(map[string]interface{})
		err = ReportTypes.UnpackIntoMap(reportElems, report)
		require.NoError(t, err)

		assert.Equal(t, int(reportElems["observationsTimestamp"].(uint32)), 242)
		assert.Equal(t, reportElems["benchmarkPrice"].(*big.Int).Int64(), int64(243))
		assert.Equal(t, reportElems["bid"].(*big.Int).Int64(), int64(244))
		assert.Equal(t, reportElems["ask"].(*big.Int).Int64(), int64(245))
		assert.Equal(t, reportElems["currentBlockNum"].(uint64), uint64(248))
		assert.Equal(t, common.Hash(reportElems["currentBlockHash"].([32]byte)), common.BytesToHash(hash))
		assert.Equal(t, reportElems["currentBlockTimestamp"].(uint64), uint64(123))
		assert.Equal(t, reportElems["validFromBlockNum"].(uint64), uint64(46))

		assert.Equal(t, types.Report{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf8, 0x55, 0x2c, 0x2c, 0xea, 0x3a, 0xb4, 0x3b, 0xae, 0x13, 0x7d, 0x89, 0xee, 0x61, 0x42, 0xa0, 0x1d, 0xb3, 0xae, 0x2b, 0x56, 0x78, 0xbc, 0x3c, 0x9b, 0xd5, 0xf5, 0x9, 0xf5, 0x37, 0xbe, 0xa5, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b}, report)

		max, err := r.MaxReportLength(4)
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
			assert.Equal(t, uint64(248), decoded.CurrentBlockNum)
			assert.Equal(t, [32]byte(common.BytesToHash(hash)), decoded.CurrentBlockHash)
			assert.Equal(t, uint64(123), decoded.CurrentBlockTimestamp)
			assert.Equal(t, uint64(46), decoded.ValidFromBlockNum)
		})
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

func buildSampleReport(bn, validFromBn int64, feedID [32]byte) []byte {
	timestamp := uint32(42)
	bp := big.NewInt(242)
	bid := big.NewInt(243)
	ask := big.NewInt(244)
	currentBlockNumber := uint64(bn)
	currentBlockHash := utils.NewHash()
	currentBlockTimestamp := uint64(123)
	validFromBlockNum := uint64(validFromBn)

	b, err := ReportTypes.Pack(feedID, timestamp, bp, bid, ask, currentBlockNumber, currentBlockHash, validFromBlockNum, currentBlockTimestamp)
	if err != nil {
		panic(err)
	}
	return b
}

func Test_ReportCodec_CurrentBlockNumFromReport(t *testing.T) {
	r := ReportCodec{}
	feedID := utils.NewHash()

	var validBn int64 = 42
	var invalidBn int64 = -1

	t.Run("CurrentBlockNumFromReport extracts the current block number from a valid report", func(t *testing.T) {
		report := buildSampleReport(validBn, 143, feedID)

		bn, err := r.CurrentBlockNumFromReport(report)
		require.NoError(t, err)

		assert.Equal(t, validBn, bn)
	})
	t.Run("CurrentBlockNumFromReport returns error if block num is too large", func(t *testing.T) {
		report := buildSampleReport(invalidBn, 143, feedID)

		_, err := r.CurrentBlockNumFromReport(report)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "CurrentBlockNum=18446744073709551615 overflows max int64")
	})
}
func Test_ReportCodec_ValidFromBlockNumFromReport(t *testing.T) {
	r := ReportCodec{}
	feedID := utils.NewHash()

	t.Run("ValidFromBlockNumFromReport extracts the valid from block number from a valid report", func(t *testing.T) {
		report := buildSampleReport(42, 999, feedID)

		bn, err := r.ValidFromBlockNumFromReport(report)
		require.NoError(t, err)

		assert.Equal(t, int64(999), bn)
	})
	t.Run("ValidFromBlockNumFromReport returns error if valid from block number is too large", func(t *testing.T) {
		report := buildSampleReport(42, -1, feedID)

		_, err := r.ValidFromBlockNumFromReport(report)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "ValidFromBlockNum=18446744073709551615 overflows max int64")
	})
}

func Test_ReportCodec_BenchmarkPriceFromReport(t *testing.T) {
	r := ReportCodec{}
	feedID := utils.NewHash()

	t.Run("BenchmarkPriceFromReport extracts the benchmark price from valid report", func(t *testing.T) {
		report := buildSampleReport(42, 999, feedID)

		bp, err := r.BenchmarkPriceFromReport(report)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(242), bp)
	})
	t.Run("BenchmarkPriceFromReport errors on invalid report", func(t *testing.T) {
		_, err := r.BenchmarkPriceFromReport([]byte{1, 2, 3})
		require.Error(t, err)
		assert.EqualError(t, err, "failed to decode report: abi: cannot marshal in to go type: length insufficient 3 require 32")
	})
}
