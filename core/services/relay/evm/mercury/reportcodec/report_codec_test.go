package reportcodec

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var hash = hexutil.MustDecode("0x552c2cea3ab43bae137d89ee6142a01db3ae2b5678bc3c9bd5f509f537bea57b")
var paos = []relaymercury.ParsedAttributedObservation{
	{
		Timestamp: uint32(42),
		Observer:  commontypes.OracleID(49),

		BenchmarkPrice: big.NewInt(43),
		Bid:            big.NewInt(44),
		Ask:            big.NewInt(45),
		PricesValid:    true,

		CurrentBlockNum:       48,
		CurrentBlockHash:      hash,
		CurrentBlockTimestamp: uint64(123),
		CurrentBlockValid:     true,
	},
	{
		Timestamp: uint32(142),
		Observer:  commontypes.OracleID(149),

		BenchmarkPrice: big.NewInt(143),
		Bid:            big.NewInt(144),
		Ask:            big.NewInt(145),
		PricesValid:    true,

		CurrentBlockNum:       48,
		CurrentBlockHash:      hash,
		CurrentBlockTimestamp: uint64(123),
		CurrentBlockValid:     true,
	},
	{
		Timestamp: uint32(242),
		Observer:  commontypes.OracleID(249),

		BenchmarkPrice: big.NewInt(243),
		Bid:            big.NewInt(244),
		Ask:            big.NewInt(245),
		PricesValid:    true,

		CurrentBlockNum:       248,
		CurrentBlockHash:      hash,
		CurrentBlockTimestamp: uint64(789),
		CurrentBlockValid:     true,
	},
	{
		Timestamp: uint32(342),
		Observer:  commontypes.OracleID(250),

		BenchmarkPrice: big.NewInt(343),
		Bid:            big.NewInt(344),
		Ask:            big.NewInt(345),
		PricesValid:    true,

		CurrentBlockNum:       348,
		CurrentBlockHash:      hash,
		CurrentBlockTimestamp: uint64(123456),
		CurrentBlockValid:     true,
	},
}

func Test_ReportCodec_BuildReport(t *testing.T) {
	r := EVMReportCodec{}

	f := 1

	t.Run("BuildReport errors if observations are empty", func(t *testing.T) {
		paos := []relaymercury.ParsedAttributedObservation{}
		_, err := r.BuildReport(paos, f, 1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot build report from empty attributed observation")
	})

	t.Run("BuildReport constructs a report from observations", func(t *testing.T) {
		// only need to test happy path since validations are done in relaymercury

		report, err := r.BuildReport(paos, f, 46)
		require.NoError(t, err)

		reportElems := make(map[string]interface{})
		err = ReportTypes.UnpackIntoMap(reportElems, report)
		require.NoError(t, err)

		assert.Equal(t, int(reportElems["observationsTimestamp"].(uint32)), 242)
		assert.Equal(t, reportElems["benchmarkPrice"].(*big.Int).Int64(), int64(243))
		assert.Equal(t, reportElems["bid"].(*big.Int).Int64(), int64(244))
		assert.Equal(t, reportElems["ask"].(*big.Int).Int64(), int64(245))
		assert.Equal(t, reportElems["currentBlockNum"].(uint64), uint64(48))
		assert.Equal(t, common.Hash(reportElems["currentBlockHash"].([32]byte)), common.BytesToHash(hash))
		assert.Equal(t, reportElems["currentBlockTimestamp"].(uint64), uint64(123))
		assert.Equal(t, reportElems["validFromBlockNum"].(uint64), uint64(46))

		assert.Equal(t, types.Report{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x30, 0x55, 0x2c, 0x2c, 0xea, 0x3a, 0xb4, 0x3b, 0xae, 0x13, 0x7d, 0x89, 0xee, 0x61, 0x42, 0xa0, 0x1d, 0xb3, 0xae, 0x2b, 0x56, 0x78, 0xbc, 0x3c, 0x9b, 0xd5, 0xf5, 0x9, 0xf5, 0x37, 0xbe, 0xa5, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b}, report)
		max, err := r.MaxReportLength(4)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(report), max)
	})
}

func Test_ReportCodec_MaxReportLength(t *testing.T) {
	r := EVMReportCodec{}
	n := len(paos)
	report, err := r.BuildReport(paos, 1, 46)
	require.NoError(t, err)

	t.Run("MaxReportLength returns length of report", func(t *testing.T) {
		max, err := r.MaxReportLength(n)
		require.NoError(t, err)
		assert.Equal(t, len(report), max)
	})
}

func Test_ReportCodec_CurrentBlockNumFromReport(t *testing.T) {
	r := EVMReportCodec{}
	feedID := utils.NewHash()

	var validBn int64 = 42
	var invalidBn int64 = -1

	t.Run("CurrentBlockNumFromReport extracts the current block number from a valid report", func(t *testing.T) {
		report := buildSampleReport(validBn, feedID)

		bn, err := r.CurrentBlockNumFromReport(report)
		require.NoError(t, err)

		assert.Equal(t, validBn, bn)
	})
	t.Run("CurrentBlockNumFromReport returns error if block num is too large", func(t *testing.T) {
		report := buildSampleReport(invalidBn, feedID)

		_, err := r.CurrentBlockNumFromReport(report)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "blockNum overflows max int64, got: 18446744073709551615")
	})
}

func Test_ReportCodec_FeedIDFromReport(t *testing.T) {
	r := EVMReportCodec{}

	feedID := utils.NewHash()
	var validBn int64 = 42

	t.Run("FeedIDFromReport extracts the current block number from a valid report", func(t *testing.T) {
		report := buildSampleReport(validBn, feedID)

		f, err := r.FeedIDFromReport(report)
		require.NoError(t, err)

		assert.Equal(t, feedID[:], f[:])
	})
	t.Run("FeedIDFromReport returns error if report is invalid", func(t *testing.T) {
		report := []byte{1}

		_, err := r.FeedIDFromReport(report)
		assert.EqualError(t, err, "invalid length for report: 1")
	})
}

func buildSampleReport(bn int64, feedID [32]byte) []byte {
	timestamp := uint32(42)
	bp := big.NewInt(242)
	bid := big.NewInt(243)
	ask := big.NewInt(244)
	currentBlockNumber := uint64(bn)
	currentBlockHash := utils.NewHash()
	currentBlockTimestamp := uint64(123)
	validFromBlockNum := uint64(143)

	b, err := ReportTypes.Pack(feedID, timestamp, bp, bid, ask, currentBlockNumber, currentBlockHash, currentBlockTimestamp, validFromBlockNum)
	if err != nil {
		panic(err)
	}
	return b
}
