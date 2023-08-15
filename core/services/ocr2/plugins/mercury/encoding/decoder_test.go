package encoding

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"

	mercuryv1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1"
	mercuryv2 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v2"
	mercuryv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	hash          = hexutil.MustDecode("0x552c2cea3ab43bae137d89ee6142a01db3ae2b5678bc3c9bd5f509f537bea57b")
	v1FeedId      = [32]uint8{00, 01, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
	v2FeedId      = [32]uint8{00, 02, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
	v3FeedId      = [32]uint8{00, 03, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
	invalidFeedId = [32]uint8{28, 145, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
)

func buildV1Report(fId [32]byte) []byte {
	feedID := fId
	timestamp := uint32(42)
	bp := big.NewInt(242)
	bid := big.NewInt(243)
	ask := big.NewInt(244)
	currentBlockNumber := uint64(143)
	currentBlockHash := common.BytesToHash(hash)
	currentBlockTimestamp := uint64(123)
	validFromBlockNum := uint64(142)

	report, err := mercuryv1.ReportTypes.Pack(feedID, timestamp, bp, bid, ask, currentBlockNumber, currentBlockHash, validFromBlockNum, currentBlockTimestamp)
	if err != nil {
		panic(err)
	}
	return report
}

func buildV2Report(fId [32]byte) []byte {
	feedID := fId
	timestamp := uint32(52)
	bp := big.NewInt(342)
	validFromTimestamp := uint32(343)
	expiresAt := uint32(344)
	linkFee := big.NewInt(345)
	nativeFee := big.NewInt(346)

	report, err := mercuryv2.ReportTypes.Pack(feedID, timestamp, bp, validFromTimestamp, expiresAt, linkFee, nativeFee)
	if err != nil {
		panic(err)
	}
	return report
}

func buildV3Report(fId [32]byte) []byte {
	feedID := fId
	timestamp := uint32(62)
	bp := big.NewInt(442)
	bid := big.NewInt(443)
	ask := big.NewInt(444)
	validFromTimestamp := uint32(445)
	expiresAt := uint32(446)
	linkFee := big.NewInt(447)
	nativeFee := big.NewInt(448)

	report, err := mercuryv3.ReportTypes.Pack(feedID, timestamp, bp, bid, ask, validFromTimestamp, expiresAt, linkFee, nativeFee)
	if err != nil {
		panic(err)
	}
	return report
}

func Test_ReportDecoder(t *testing.T) {
	lggr := logger.TestLogger(t)

	t.Run("invalid report length", func(t *testing.T) {
		_, err := DecodeV1([]byte{0x00}, lggr)
		assert.EqualError(t, err, "invalid length for report: 1")

		_, err = DecodeV2([]byte{0x00, 0x01}, lggr)
		assert.EqualError(t, err, "invalid length for report: 2")

		_, err = DecodeV3([]byte{0x00, 0x01, 0x02}, lggr)
		assert.EqualError(t, err, "invalid length for report: 3")
	})

	t.Run("invalid schema version", func(t *testing.T) {
		report := buildV1Report(invalidFeedId)
		_, err := DecodeV1(report, lggr)
		assert.EqualError(t, err, "invalid schema version: 7313")

		report = buildV2Report(invalidFeedId)
		_, err = DecodeV2(report, lggr)
		assert.EqualError(t, err, "invalid schema version: 7313")

		report = buildV3Report(invalidFeedId)
		_, err = DecodeV3(report, lggr)
		assert.EqualError(t, err, "invalid schema version: 7313")
	})
	t.Run("v1", func(t *testing.T) {
		v1Report := buildV1Report(v1FeedId)
		report, err := DecodeV1(v1Report, lggr)
		assert.NoError(t, err)
		assert.Equal(t, v1FeedId, report.FeedId)
		assert.Equal(t, uint32(42), report.ObservationsTimestamp)
		assert.Equal(t, big.NewInt(242), report.BenchmarkPrice)
		assert.Equal(t, big.NewInt(243), report.Bid)
		assert.Equal(t, big.NewInt(244), report.Ask)
		assert.Equal(t, uint64(143), report.CurrentBlockNum)
		assert.Equal(t, hash, report.CurrentBlockHash[:])
		assert.Equal(t, uint64(123), report.CurrentBlockTimestamp)
		assert.Equal(t, uint64(142), report.ValidFromBlockNum)
	})

	t.Run("v2", func(t *testing.T) {
		v2Report := buildV2Report(v2FeedId)
		report, err := DecodeV2(v2Report, lggr)
		assert.NoError(t, err)
		assert.Equal(t, v2FeedId, report.FeedId)
		assert.Equal(t, uint32(52), report.ObservationsTimestamp)
		assert.Equal(t, big.NewInt(342), report.BenchmarkPrice)
		assert.Equal(t, uint32(343), report.ValidFromTimestamp)
		assert.Equal(t, uint32(344), report.ExpiresAt)
		assert.Equal(t, big.NewInt(345), report.LinkFee)
		assert.Equal(t, big.NewInt(346), report.NativeFee)
	})

	t.Run("v3", func(t *testing.T) {
		v3Report := buildV3Report(v3FeedId)
		report, err := DecodeV3(v3Report, lggr)
		assert.NoError(t, err)
		assert.Equal(t, v3FeedId, report.FeedId)
		assert.Equal(t, uint32(62), report.ObservationsTimestamp)
		assert.Equal(t, big.NewInt(442), report.BenchmarkPrice)
		assert.Equal(t, big.NewInt(443), report.Bid)
		assert.Equal(t, big.NewInt(444), report.Ask)
		assert.Equal(t, uint32(445), report.ValidFromTimestamp)
		assert.Equal(t, uint32(446), report.ExpiresAt)
		assert.Equal(t, big.NewInt(447), report.LinkFee)
		assert.Equal(t, big.NewInt(448), report.NativeFee)
	})
}
