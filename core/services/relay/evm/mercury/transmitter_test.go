package mercury

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type MockWSRPCClient struct {
	transmit     func(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error)
	latestReport func(ctx context.Context, req *pb.LatestReportRequest) (resp *pb.LatestReportResponse, err error)
}

func (m MockWSRPCClient) Name() string                   { return "" }
func (m MockWSRPCClient) Start(context.Context) error    { return nil }
func (m MockWSRPCClient) Close() error                   { return nil }
func (m MockWSRPCClient) Healthy() error                 { return nil }
func (m MockWSRPCClient) HealthReport() map[string]error { return map[string]error{} }
func (m MockWSRPCClient) Ready() error                   { return nil }
func (m MockWSRPCClient) Transmit(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error) {
	return m.transmit(ctx, in)
}
func (m MockWSRPCClient) LatestReport(ctx context.Context, in *pb.LatestReportRequest) (*pb.LatestReportResponse, error) {
	return m.latestReport(ctx, in)
}

var _ wsrpc.Client = &MockWSRPCClient{}

func Test_MercuryTransmitter_Transmit(t *testing.T) {
	lggr := logger.TestLogger(t)
	reportURL := "http://report.test/foo"

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
		mt := NewTransmitter(lggr, c, reportURL, sampleClientPubKey, sampleFeedID)
		err := mt.Transmit(testutils.Context(t), sampleReportContext, sampleReport, sampleSigs)

		require.NoError(t, err)
	})

	t.Run("failing transmit", func(t *testing.T) {
		c := MockWSRPCClient{
			transmit: func(ctx context.Context, in *pb.TransmitRequest) (out *pb.TransmitResponse, err error) {
				return nil, errors.New("foo error")
			},
		}
		mt := NewTransmitter(lggr, c, reportURL, sampleClientPubKey, sampleFeedID)
		err := mt.Transmit(testutils.Context(t), sampleReportContext, sampleReport, sampleSigs)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "foo error")
	})
}

func Test_MercuryTransmitter_LatestConfigDigestAndEpoch(t *testing.T) {
	lggr := logger.TestLogger(t)
	reportURL := "http://report.test/foo"

	sampleConfigDigest := utils.NewHash().Bytes()
	wrongFeedID := []byte{1, 2, 3, 4}

	t.Run("successful query", func(t *testing.T) {
		c := MockWSRPCClient{
			latestReport: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.FeedId = sampleFeedID[:]
				out.ConfigDigest = sampleConfigDigest
				out.Epoch = 42
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, c, reportURL, sampleClientPubKey, sampleFeedID)
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
		mt := NewTransmitter(lggr, c, reportURL, sampleClientPubKey, sampleFeedID)
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
				out.FeedId = wrongFeedID
				out.ConfigDigest = sampleConfigDigest
				out.Epoch = 42
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, c, reportURL, sampleClientPubKey, sampleFeedID)
		_, _, err := mt.LatestConfigDigestAndEpoch(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "LatestConfigDigestAndEpoch failed; mismatched feed IDs, expected: 0x1c916b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472, got: 0x01020304")
	})
}

func Test_MercuryTransmitter_FetchInitialMaxFinalizedBlockNumber(t *testing.T) {
	lggr := logger.TestLogger(t)
	reportURL := "http://report.test/foo"

	t.Run("successful query", func(t *testing.T) {
		c := MockWSRPCClient{
			latestReport: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.FeedId = sampleFeedID[:]
				out.BlockNumber = 42
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, c, reportURL, sampleClientPubKey, sampleFeedID)
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
		mt := NewTransmitter(lggr, c, reportURL, sampleClientPubKey, sampleFeedID)
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
				out.BlockNumber = 42
				out.FeedId = nil
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, c, reportURL, sampleClientPubKey, sampleFeedID)
		_, err := mt.FetchInitialMaxFinalizedBlockNumber(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FetchInitialMaxFinalizedBlockNumber failed; mismatched feed IDs, expected: 0x1c916b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472, got: 0x")
	})
}
