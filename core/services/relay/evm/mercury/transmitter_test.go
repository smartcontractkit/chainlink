package mercury

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	mocks "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

var (
	sampleFeedID       = [32]uint8{28, 145, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
	sampleClientPubKey = hexutil.MustDecode("0x724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93")
)

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
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
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
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
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
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
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
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
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
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
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
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
		_, err := mt.LatestPrice(testutils.Context(t), sampleFeedID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")
	})
}

func Test_MercuryTransmitter_FetchInitialMaxFinalizedBlockNumber(t *testing.T) {
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
				out.Report.CurrentBlockNumber = 42
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
		bn, err := mt.FetchInitialMaxFinalizedBlockNumber(testutils.Context(t))
		require.NoError(t, err)

		require.NotNil(t, bn)
		assert.Equal(t, 42, int(*bn))
	})
	t.Run("successful query returning nil report (new feed)", func(t *testing.T) {
		c := mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				out = new(pb.LatestReportResponse)
				out.Report = nil
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
		bn, err := mt.FetchInitialMaxFinalizedBlockNumber(testutils.Context(t))
		require.NoError(t, err)

		assert.Nil(t, bn)
	})
	t.Run("failing query", func(t *testing.T) {
		c := mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				return nil, errors.New("something exploded")
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
		_, err := mt.FetchInitialMaxFinalizedBlockNumber(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")
	})
	t.Run("return feed ID is wrong", func(t *testing.T) {
		c := mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.CurrentBlockNumber = 42
				out.Report.FeedId = []byte{1, 2}
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID, db, pgtest.NewQConfig(true))
		_, err := mt.FetchInitialMaxFinalizedBlockNumber(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "latestReport failed; mismatched feed IDs, expected: 0x1c916b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472, got: 0x")
	})
}
