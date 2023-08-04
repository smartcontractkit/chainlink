package mercury_v1

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	mocks "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

var (
	sampleFeedID        = [32]uint8{28, 145, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
	sampleClientPubKey  = hexutil.MustDecode("0x724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93")
	sig2                = ocrtypes.AttributedOnchainSignature{Signature: mercury.MustDecodeBase64("kbeuRczizOJCxBzj7MUAFpz3yl2WRM6K/f0ieEBvA+oTFUaKslbQey10krumVjzAvlvKxMfyZo0WkOgNyfF6xwE="), Signer: 2}
	sig3                = ocrtypes.AttributedOnchainSignature{Signature: mercury.MustDecodeBase64("9jz4b6Dh2WhXxQ97a6/S9UNjSfrEi9016XKTrfN0mLQFDiNuws23x7Z4n+6g0sqKH/hnxx1VukWUH/ohtw83/wE="), Signer: 3}
	sampleSigs          = []ocrtypes.AttributedOnchainSignature{sig2, sig3}
	sampleReportContext = ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: mercury.MustHexToConfigDigest("0x0006fc30092226b37f6924b464e16a54a7978a9a524519a73403af64d487dc45"),
			Epoch:        6,
			Round:        28,
		},
		ExtraHash: [32]uint8{27, 144, 106, 73, 166, 228, 123, 166, 179, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114},
	}
	sampleReport     = buildSampleReport(123)
	samplePayload    = mercury.BuildSamplePayload(sampleReport, sampleReportContext, sampleSigs)
	samplePayloadHex = hexutil.Encode(samplePayload)
)

func Test_MercuryTransmitter_Transmit(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)

	t.Run("transmission successfully enqueued", func(t *testing.T) {
		c := mocks.MockWSRPCClient{
			TransmitF: func(ctx context.Context, in *pb.TransmitRequest) (out *pb.TransmitResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, samplePayloadHex, hexutil.Encode(in.Payload))
				out = new(pb.TransmitResponse)
				out.Code = 42
				out.Error = ""
				return out, nil
			},
		}
		mt := mercury.NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
		err := mt.Transmit(testutils.Context(t), sampleReportContext, sampleReport, sampleSigs)

		require.NoError(t, err)
	})
}

func Test_MercuryTransmitter_LatestTimestamp(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)

	t.Run("successful query", func(t *testing.T) {
		c := mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.FeedId = sampleFeedID[:]
				out.Report.ObservationsTimestamp = 42
				return out, nil
			},
		}
		mt := mercury.NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
		ts, err := mt.LatestTimestamp(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, ts, uint32(42))
	})

	t.Run("successful query returning nil report (new feed)", func(t *testing.T) {
		c := mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				out = new(pb.LatestReportResponse)
				out.Report = nil
				return out, nil
			},
		}
		mt := mercury.NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
		ts, err := mt.LatestTimestamp(testutils.Context(t))
		require.NoError(t, err)

		assert.Zero(t, ts)
	})

	t.Run("failing query", func(t *testing.T) {
		c := mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				return nil, errors.New("something exploded")
			},
		}
		mt := mercury.NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
		_, err := mt.LatestTimestamp(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")
	})
}

func Test_MercuryTransmitter_LatestPrice(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)

	t.Run("successful query", func(t *testing.T) {
		originalPrice := big.NewInt(123456789)
		encodedPrice, _ := relaymercury.EncodeValueInt192(originalPrice)
		c := mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.FeedId = sampleFeedID[:]
				out.Report.Price = encodedPrice
				return out, nil
			},
		}
		mt := mercury.NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
		price, err := mt.LatestPrice(testutils.Context(t), sampleFeedID)
		require.NoError(t, err)

		assert.Equal(t, price, originalPrice)
	})

	t.Run("successful query returning nil report (new feed)", func(t *testing.T) {
		c := mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				out = new(pb.LatestReportResponse)
				out.Report = nil
				return out, nil
			},
		}
		mt := mercury.NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
		price, err := mt.LatestPrice(testutils.Context(t), sampleFeedID)
		require.NoError(t, err)

		assert.Nil(t, price)
	})

	t.Run("failing query", func(t *testing.T) {
		c := mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				return nil, errors.New("something exploded")
			},
		}
		mt := mercury.NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
		_, err := mt.LatestPrice(testutils.Context(t), sampleFeedID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")
	})
}
