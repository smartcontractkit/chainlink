package mercury

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type MockWSRPCClient struct {
	transmit     func(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error)
	latestReport func(ctx context.Context, req *pb.LatestReportRequest) (resp *pb.LatestReportResponse, err error)
}

func (m MockWSRPCClient) Name() string                   { return "" }
func (m MockWSRPCClient) Start(context.Context) error    { return nil }
func (m MockWSRPCClient) Close() error                   { return nil }
func (m MockWSRPCClient) HealthReport() map[string]error { return map[string]error{} }
func (m MockWSRPCClient) Ready() error                   { return nil }
func (m MockWSRPCClient) Transmit(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error) {
	return m.transmit(ctx, in)
}
func (m MockWSRPCClient) LatestReport(ctx context.Context, in *pb.LatestReportRequest) (*pb.LatestReportResponse, error) {
	return m.latestReport(ctx, in)
}

var _ wsrpc.Client = &MockWSRPCClient{}

type MockTracker struct {
	latestConfigDetails func(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error)
}

func (m MockTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	return m.latestConfigDetails(ctx)
}

var _ ConfigTracker = &MockTracker{}

func Test_MercuryTransmitter_Transmit(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)

	t.Run("successful transmit", func(t *testing.T) {
		c := MockWSRPCClient{
			transmit: func(ctx context.Context, in *pb.TransmitRequest) (out *pb.TransmitResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, samplePayloadHex, hexutil.Encode(in.Payload))
				out = new(pb.TransmitResponse)
				out.Code = 42
				out.Error = ""
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID)
		err := mt.Transmit(testutils.Context(t), sampleReportContext, sampleReport, sampleSigs)

		require.NoError(t, err)
	})

	t.Run("failing transmit", func(t *testing.T) {
		c := MockWSRPCClient{
			transmit: func(ctx context.Context, in *pb.TransmitRequest) (out *pb.TransmitResponse, err error) {
				return nil, errors.New("foo error")
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID)
		err := mt.Transmit(testutils.Context(t), sampleReportContext, sampleReport, sampleSigs)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "foo error")
	})
}

func Test_MercuryTransmitter_LatestConfigDigestAndEpoch(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)

	sampleConfigDigest := utils.NewHash().Bytes()
	wrongFeedID := []byte{1, 2, 3, 4}

	t.Run("successful query", func(t *testing.T) {
		c := MockWSRPCClient{
			latestReport: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.FeedId = sampleFeedID[:]
				out.Report.ConfigDigest = sampleConfigDigest
				out.Report.Epoch = 42
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID)
		cd, epoch, err := mt.LatestConfigDigestAndEpoch(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, hexutil.Encode(sampleConfigDigest[:]), hexutil.Encode(cd[:]))
		assert.Equal(t, 42, int(epoch))
	})
	t.Run("failing query", func(t *testing.T) {
		c := MockWSRPCClient{
			latestReport: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				return nil, errors.New("something exploded")
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID)
		_, _, err := mt.LatestConfigDigestAndEpoch(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")
	})
	t.Run("return feed ID is wrong", func(t *testing.T) {
		c := MockWSRPCClient{
			latestReport: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.FeedId = wrongFeedID
				out.Report.ConfigDigest = sampleConfigDigest
				out.Report.Epoch = 42
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID)
		_, _, err := mt.LatestConfigDigestAndEpoch(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "LatestConfigDigestAndEpoch failed; mismatched feed IDs, expected: 0x1c916b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472, got: 0x01020304")
	})
	t.Run("LatestReport returns nil response", func(t *testing.T) {
		c := MockWSRPCClient{
			latestReport: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				return
			},
		}
		tracker := &MockTracker{}
		mt := NewTransmitter(lggr, tracker, c, sampleClientPubKey, sampleFeedID)
		_, _, err := mt.LatestConfigDigestAndEpoch(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "LatestConfigDigestAndEpoch expected LatestReport to return non-nil response")
	})

	t.Run("falls back to latest config details if report is nil", func(t *testing.T) {
		c := MockWSRPCClient{
			latestReport: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				out = new(pb.LatestReportResponse)
				return
			},
		}
		t.Run("LatestConfigDetails succeeds", func(t *testing.T) {
			tracker := &MockTracker{
				latestConfigDetails: func(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
					return 123, ocrtypes.ConfigDigest(sampleConfigDigest), nil
				},
			}
			mt := NewTransmitter(lggr, tracker, c, sampleClientPubKey, sampleFeedID)
			cd, epoch, err := mt.LatestConfigDigestAndEpoch(testutils.Context(t))
			require.NoError(t, err)

			assert.Equal(t, ocrtypes.ConfigDigest(sampleConfigDigest), cd)
			assert.Equal(t, 0, int(epoch)) // always returns zero epoch if latest report is empty
		})
		t.Run("LatestConfigDetails fails", func(t *testing.T) {
			tracker := &MockTracker{
				latestConfigDetails: func(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
					return changedInBlock, configDigest, errors.New("something exploded")
				},
			}
			mt := NewTransmitter(lggr, tracker, c, sampleClientPubKey, sampleFeedID)
			_, _, err := mt.LatestConfigDigestAndEpoch(testutils.Context(t))
			require.Error(t, err)
			assert.Contains(t, err.Error(), "something exploded")
		})
	})
}

func Test_MercuryTransmitter_FetchInitialMaxFinalizedBlockNumber(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)

	t.Run("successful query", func(t *testing.T) {
		c := MockWSRPCClient{
			latestReport: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.FeedId = sampleFeedID[:]
				out.Report.CurrentBlockNumber = 42
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID)
		bn, err := mt.FetchInitialMaxFinalizedBlockNumber(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, 42, int(bn))
	})
	t.Run("failing query", func(t *testing.T) {
		c := MockWSRPCClient{
			latestReport: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				return nil, errors.New("something exploded")
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID)
		_, err := mt.FetchInitialMaxFinalizedBlockNumber(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")
	})
	t.Run("return feed ID is wrong", func(t *testing.T) {
		c := MockWSRPCClient{
			latestReport: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.CurrentBlockNumber = 42
				out.Report.FeedId = []byte{1, 2}
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, nil, c, sampleClientPubKey, sampleFeedID)
		_, err := mt.FetchInitialMaxFinalizedBlockNumber(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FetchInitialMaxFinalizedBlockNumber failed; mismatched feed IDs, expected: 0x1c916b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472, got: 0x")
	})
}
