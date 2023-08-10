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

func Test_DecodeSchemaVersionFromFeedId(t *testing.T) {
	tests := []struct {
		name    string
		feedID  [32]byte
		want    uint16
		wantErr bool
	}{
		{
			name:    "schemaVersion v1",
			feedID:  [32]byte{0x00, 0x01},
			want:    1,
			wantErr: false,
		},
		{
			name:    "schemaVersion v2",
			feedID:  [32]byte{0x00, 0x02},
			want:    2,
			wantErr: false,
		},
		{
			name:    "schemaVersion v3",
			feedID:  [32]byte{0x00, 0x03},
			want:    3,
			wantErr: false,
		},
		{
			name:    "schemaVersion invalid",
			feedID:  [32]byte{0x00, 0x04},
			want:    0,
			wantErr: true,
		},
		{
			name:    "schemaVersion invalid",
			feedID:  [32]byte{0x00, 0x00},
			want:    0,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := DecodeSchemaVersionFromFeedId(test.feedID)
			if (err != nil) != test.wantErr {
				t.Errorf("DecodeSchemaVersionFromFeedId() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("DecodeSchemaVersionFromFeedId() = %v, want %v", got, test.want)
			}
		})
	}
}

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
		_, err := NewReportDecoder([]byte{0x00}, lggr)
		assert.EqualError(t, err, "invalid length for report: 1")
	})

	t.Run("invalid feed id", func(t *testing.T) {
		report := buildV1Report(invalidFeedId)
		_, err := NewReportDecoder(report, lggr)
		assert.EqualError(t, err, "invalid schema version: 7313")
	})
	t.Run("v1", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			v1Report := buildV1Report(v1FeedId)
			reportDecoder, err := NewReportDecoder(v1Report, lggr)
			assert.NoError(t, err)
			assert.Equal(t, reportDecoder.GetSchemaVersion(), REPORT_V1)
			report, err := reportDecoder.DecodeAsV1()
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

		t.Run("errors if invalid schema version", func(t *testing.T) {
			v1Report := buildV1Report(v2FeedId)
			reportDecoder, err := NewReportDecoder(v1Report, lggr)
			assert.NoError(t, err)
			assert.NotEqual(t, reportDecoder.GetSchemaVersion(), REPORT_V1)
			_, err = reportDecoder.DecodeAsV1()
			assert.EqualError(t, err, "invalid schema version: 2")
		})

	})

	t.Run("v2", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			v2Report := buildV2Report(v2FeedId)
			reportDecoder, err := NewReportDecoder(v2Report, lggr)
			assert.NoError(t, err)
			assert.Equal(t, reportDecoder.GetSchemaVersion(), REPORT_V2)
			report, err := reportDecoder.DecodeAsV2()
			assert.NoError(t, err)
			assert.Equal(t, v2FeedId, report.FeedId)
			assert.Equal(t, uint32(52), report.ObservationsTimestamp)
			assert.Equal(t, big.NewInt(342), report.BenchmarkPrice)
			assert.Equal(t, uint32(343), report.ValidFromTimestamp)
			assert.Equal(t, uint32(344), report.ExpiresAt)
			assert.Equal(t, big.NewInt(345), report.LinkFee)
			assert.Equal(t, big.NewInt(346), report.NativeFee)
		})

		t.Run("errors if invalid schema version", func(t *testing.T) {
			v2Report := buildV2Report(v3FeedId)
			reportDecoder, err := NewReportDecoder(v2Report, lggr)
			assert.NoError(t, err)
			assert.NotEqual(t, reportDecoder.GetSchemaVersion(), REPORT_V2)
			_, err = reportDecoder.DecodeAsV2()
			assert.EqualError(t, err, "invalid schema version: 3")
		})

		t.Run("errors when decoding wrong report of larger size", func(t *testing.T) {
			v2Report := buildV2Report(v1FeedId)
			reportDecoder, err := NewReportDecoder(v2Report, lggr)
			assert.NoError(t, err)
			_, err = reportDecoder.DecodeAsV1()
			assert.EqualError(t, err, "error during unpack: abi: cannot marshal in to go type: length insufficient 224 require 256")
		})
	})

	t.Run("v3", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			v3Report := buildV3Report(v3FeedId)
			reportDecoder, err := NewReportDecoder(v3Report, lggr)
			assert.NoError(t, err)
			assert.Equal(t, reportDecoder.GetSchemaVersion(), REPORT_V3)
			report, err := reportDecoder.DecodeAsV3()
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

		t.Run("errors if invalid schema version", func(t *testing.T) {
			v3Report := buildV3Report(v1FeedId)
			reportDecoder, err := NewReportDecoder(v3Report, lggr)
			assert.NoError(t, err)
			assert.NotEqual(t, reportDecoder.GetSchemaVersion(), REPORT_V3)
			_, err = reportDecoder.DecodeAsV3()
			assert.EqualError(t, err, "invalid schema version: 1")
		})
	})
}
