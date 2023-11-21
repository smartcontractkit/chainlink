package mercury

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	mercurytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
	mocks "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

func Test_MercuryTransmitter_Transmit(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	var jobID int32
	pgtest.MustExec(t, db, `SET CONSTRAINTS mercury_transmit_requests_job_id_fkey DEFERRED`)
	pgtest.MustExec(t, db, `SET CONSTRAINTS feed_latest_reports_job_id_fkey DEFERRED`)
	q := NewTransmitQueue(lggr, "", 0, nil, nil)

	t.Run("v1 report transmission successfully enqueued", func(t *testing.T) {
		report := sampleV1Report
		c := mocks.MockWSRPCClient{
			TransmitF: func(ctx context.Context, in *pb.TransmitRequest) (out *pb.TransmitResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(buildSamplePayload(report)), hexutil.Encode(in.Payload))
				out = new(pb.TransmitResponse)
				out.Code = 42
				out.Error = ""
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, jobID, sampleFeedID, db, pgtest.NewQConfig(true), nil)
		mt.queue = q
		err := mt.Transmit(testutils.Context(t), sampleReportContext, report, sampleSigs)

		require.NoError(t, err)
	})
	t.Run("v2 report transmission successfully enqueued", func(t *testing.T) {
		report := sampleV2Report
		c := mocks.MockWSRPCClient{
			TransmitF: func(ctx context.Context, in *pb.TransmitRequest) (out *pb.TransmitResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(buildSamplePayload(report)), hexutil.Encode(in.Payload))
				out = new(pb.TransmitResponse)
				out.Code = 42
				out.Error = ""
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, jobID, sampleFeedID, db, pgtest.NewQConfig(true), nil)
		mt.queue = q
		err := mt.Transmit(testutils.Context(t), sampleReportContext, report, sampleSigs)

		require.NoError(t, err)
	})
	t.Run("v3 report transmission successfully enqueued", func(t *testing.T) {
		report := sampleV3Report
		c := mocks.MockWSRPCClient{
			TransmitF: func(ctx context.Context, in *pb.TransmitRequest) (out *pb.TransmitResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(buildSamplePayload(report)), hexutil.Encode(in.Payload))
				out = new(pb.TransmitResponse)
				out.Code = 42
				out.Error = ""
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, jobID, sampleFeedID, db, pgtest.NewQConfig(true), nil)
		mt.queue = q
		err := mt.Transmit(testutils.Context(t), sampleReportContext, report, sampleSigs)

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
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, 0, sampleFeedID, db, pgtest.NewQConfig(true), nil)
		ts, err := mt.LatestTimestamp(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, int64(42), ts)
	})

	t.Run("successful query returning nil report (new feed) gives latest timestamp = -1", func(t *testing.T) {
		c := mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				out = new(pb.LatestReportResponse)
				out.Report = nil
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, 0, sampleFeedID, db, pgtest.NewQConfig(true), nil)
		ts, err := mt.LatestTimestamp(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, int64(-1), ts)
	})

	t.Run("failing query", func(t *testing.T) {
		c := mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				return nil, errors.New("something exploded")
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, 0, sampleFeedID, db, pgtest.NewQConfig(true), nil)
		_, err := mt.LatestTimestamp(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")
	})
}

type mockCodec struct {
	val *big.Int
	err error
}

var _ mercurytypes.ReportCodec = &mockCodec{}

func (m *mockCodec) BenchmarkPriceFromReport(_ ocrtypes.Report) (*big.Int, error) {
	return m.val, m.err
}

func Test_MercuryTransmitter_LatestPrice(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)

	codec := new(mockCodec)

	t.Run("successful query", func(t *testing.T) {
		originalPrice := big.NewInt(123456789)
		c := mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.FeedId = sampleFeedID[:]
				out.Report.Payload = buildSamplePayload([]byte("doesn't matter"))
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, 0, sampleFeedID, db, pgtest.NewQConfig(true), codec)

		t.Run("BenchmarkPriceFromReport succeeds", func(t *testing.T) {
			codec.val = originalPrice
			codec.err = nil

			price, err := mt.LatestPrice(testutils.Context(t), sampleFeedID)
			require.NoError(t, err)

			assert.Equal(t, originalPrice, price)
		})
		t.Run("BenchmarkPriceFromReport fails", func(t *testing.T) {
			codec.val = nil
			codec.err = errors.New("something exploded")

			_, err := mt.LatestPrice(testutils.Context(t), sampleFeedID)
			require.Error(t, err)

			assert.EqualError(t, err, "something exploded")
		})
	})

	t.Run("successful query returning nil report (new feed)", func(t *testing.T) {
		c := mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				out = new(pb.LatestReportResponse)
				out.Report = nil
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, 0, sampleFeedID, db, pgtest.NewQConfig(true), nil)
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
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, 0, sampleFeedID, db, pgtest.NewQConfig(true), nil)
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
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, 0, sampleFeedID, db, pgtest.NewQConfig(true), nil)
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
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, 0, sampleFeedID, db, pgtest.NewQConfig(true), nil)
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
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, 0, sampleFeedID, db, pgtest.NewQConfig(true), nil)
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
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, 0, sampleFeedID, db, pgtest.NewQConfig(true), nil)
		_, err := mt.FetchInitialMaxFinalizedBlockNumber(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "latestReport failed; mismatched feed IDs, expected: 0x1c916b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472, got: 0x")
	})
}
