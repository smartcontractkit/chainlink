package reportcodec

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func Test_ReportCodec_BuildReport(t *testing.T) {
	r := EVMReportCodec{}

	f := 1

	t.Run("BuildReport errors if observations are empty", func(t *testing.T) {
		paos := []relaymercury.ParsedAttributedObservation{}
		_, err := r.BuildReport(paos, f)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot build report from empty attributed observation")
	})

	t.Run("BuildReport constructs a report from observations", func(t *testing.T) {
		// only need to test happy path since validations are done in relaymercury

		hash := hexutil.MustDecode("0x552c2cea3ab43bae137d89ee6142a01db3ae2b5678bc3c9bd5f509f537bea57b")

		paos := []relaymercury.ParsedAttributedObservation{
			{
				Timestamp:         uint32(42),
				BenchmarkPrice:    big.NewInt(43),
				Bid:               big.NewInt(44),
				Ask:               big.NewInt(45),
				CurrentBlockNum:   48,
				CurrentBlockHash:  hash,
				ValidFromBlockNum: 46,
				Observer:          commontypes.OracleID(49),
			},
			{
				Timestamp:         uint32(142),
				BenchmarkPrice:    big.NewInt(143),
				Bid:               big.NewInt(144),
				Ask:               big.NewInt(145),
				CurrentBlockNum:   48,
				CurrentBlockHash:  hash,
				ValidFromBlockNum: 46,
				Observer:          commontypes.OracleID(149),
			},
			{
				Timestamp:         uint32(242),
				BenchmarkPrice:    big.NewInt(243),
				Bid:               big.NewInt(244),
				Ask:               big.NewInt(245),
				CurrentBlockNum:   248,
				CurrentBlockHash:  hash,
				ValidFromBlockNum: 246,
				Observer:          commontypes.OracleID(249),
			},
			{
				Timestamp:         uint32(342),
				BenchmarkPrice:    big.NewInt(343),
				Bid:               big.NewInt(344),
				Ask:               big.NewInt(345),
				CurrentBlockNum:   348,
				CurrentBlockHash:  hash,
				ValidFromBlockNum: 346,
				Observer:          commontypes.OracleID(250),
			},
		}
		rep, err := r.BuildReport(paos, f)
		require.NoError(t, err)

		assert.Equal(t, types.Report{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x30, 0x55, 0x2c, 0x2c, 0xea, 0x3a, 0xb4, 0x3b, 0xae, 0x13, 0x7d, 0x89, 0xee, 0x61, 0x42, 0xa0, 0x1d, 0xb3, 0xae, 0x2b, 0x56, 0x78, 0xbc, 0x3c, 0x9b, 0xd5, 0xf5, 0x9, 0xf5, 0x37, 0xbe, 0xa5, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2e}, rep)
		assert.LessOrEqual(t, len(rep), r.MaxReportLength(4))
	})
}

func Test_ReportCodec_MaxReportLength(t *testing.T) {
	r := EVMReportCodec{}
	n := 4

	t.Run("MaxReportLength returns correct length", func(t *testing.T) {
		assert.Equal(t, 1248, r.MaxReportLength(n))
	})
}

func Test_ReportCodec_CurrentBlockNumFromReport(t *testing.T) {
	r := EVMReportCodec{}

	var validBn int64 = 42
	var invalidBn int64 = -1

	t.Run("CurrentBlockNumFromReport extracts the current block number from a valid report", func(t *testing.T) {
		report := buildSampleReport(validBn)

		bn, err := r.CurrentBlockNumFromReport(report)
		require.NoError(t, err)

		assert.Equal(t, validBn, bn)
	})
	t.Run("CurrentBlockNumFromReport returns error if block num is too large", func(t *testing.T) {
		report := buildSampleReport(invalidBn)

		_, err := r.CurrentBlockNumFromReport(report)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "blockNum overflows max int64, got: 18446744073709551615")
	})
}

func buildSampleReport(bn int64) []byte {
	feedID := [32]byte{'f', 'o', 'o'}
	timestamp := uint32(42)
	bp := big.NewInt(242)
	bid := big.NewInt(243)
	ask := big.NewInt(244)
	currentBlockNumber := uint64(bn)
	currentBlockHash := utils.NewHash()
	validFromBlockNum := uint64(143)

	b, err := ReportTypes.Pack(feedID, timestamp, bp, bid, ask, currentBlockNumber, currentBlockHash, validFromBlockNum)
	if err != nil {
		panic(err)
	}
	return b
}
