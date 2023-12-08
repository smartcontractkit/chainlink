package wsrpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

var _ cache.CacheSet = &mockCacheSet{}

type mockCacheSet struct{}

func (m *mockCacheSet) Get(ctx context.Context, client cache.Client) (cache.Fetcher, error) {
	return nil, nil
}
func (m *mockCacheSet) Start(context.Context) error    { return nil }
func (m *mockCacheSet) Ready() error                   { return nil }
func (m *mockCacheSet) HealthReport() map[string]error { return nil }
func (m *mockCacheSet) Name() string                   { return "" }
func (m *mockCacheSet) Close() error                   { return nil }

var _ cache.Cache = &mockCache{}

type mockCache struct{}

func (m *mockCache) LatestReport(ctx context.Context, req *pb.LatestReportRequest) (resp *pb.LatestReportResponse, err error) {
	return nil, nil
}
func (m *mockCache) Start(context.Context) error    { return nil }
func (m *mockCache) Ready() error                   { return nil }
func (m *mockCache) HealthReport() map[string]error { return nil }
func (m *mockCache) Name() string                   { return "" }
func (m *mockCache) Close() error                   { return nil }

func newNoopCacheSet() cache.CacheSet {
	return &mockCacheSet{}
}

func Test_Client_Transmit(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	req := &pb.TransmitRequest{}

	noopCacheSet := newNoopCacheSet()

	t.Run("sends on reset channel after MaxConsecutiveTransmitFailures timed out transmits", func(t *testing.T) {
		calls := 0
		transmitErr := context.DeadlineExceeded
		wsrpcClient := &mocks.MockWSRPCClient{
			TransmitF: func(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error) {
				calls++
				return nil, transmitErr
			},
		}
		conn := &mocks.MockConn{
			Ready: true,
		}
		c := newClient(lggr, csakey.KeyV2{}, nil, "", noopCacheSet)
		c.conn = conn
		c.rawClient = wsrpcClient
		require.NoError(t, c.StartOnce("Mock WSRPC Client", func() error { return nil }))
		for i := 1; i < MaxConsecutiveTransmitFailures; i++ {
			_, err := c.Transmit(ctx, req)
			require.EqualError(t, err, "context deadline exceeded")
		}
		assert.Equal(t, 4, calls)
		select {
		case <-c.chResetTransport:
			t.Fatal("unexpected send on chResetTransport")
		default:
		}
		_, err := c.Transmit(ctx, req)
		require.EqualError(t, err, "context deadline exceeded")
		assert.Equal(t, 5, calls)
		select {
		case <-c.chResetTransport:
		default:
			t.Fatal("expected send on chResetTransport")
		}

		t.Run("successful transmit resets the counter", func(t *testing.T) {
			transmitErr = nil
			// working transmit to reset counter
			_, err = c.Transmit(ctx, req)
			require.NoError(t, err)
			assert.Equal(t, 6, calls)
			assert.Equal(t, 0, int(c.consecutiveTimeoutCnt.Load()))
		})

		t.Run("doesn't block in case channel is full", func(t *testing.T) {
			transmitErr = context.DeadlineExceeded
			c.chResetTransport = nil // simulate full channel
			for i := 0; i < MaxConsecutiveTransmitFailures; i++ {
				_, err := c.Transmit(ctx, req)
				require.EqualError(t, err, "context deadline exceeded")
			}
		})
	})
}

func Test_Client_LatestReport(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)

	t.Run("with nil cache", func(t *testing.T) {
		req := &pb.LatestReportRequest{}
		noopCacheSet := newNoopCacheSet()
		resp := &pb.LatestReportResponse{}

		wsrpcClient := &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (*pb.LatestReportResponse, error) {
				assert.Equal(t, req, in)
				return resp, nil
			},
		}

		conn := &mocks.MockConn{
			Ready: true,
		}
		c := newClient(lggr, csakey.KeyV2{}, nil, "", noopCacheSet)
		c.conn = conn
		c.rawClient = wsrpcClient
		require.NoError(t, c.StartOnce("Mock WSRPC Client", func() error { return nil }))

		r, err := c.LatestReport(ctx, req)

		require.NoError(t, err)
		assert.Equal(t, resp, r)
	})

	t.Run("with cache disabled", func(t *testing.T) {
		req := &pb.LatestReportRequest{}
		cacheSet := cache.NewCacheSet(lggr, cache.Config{LatestReportTTL: 0})
		resp := &pb.LatestReportResponse{}

		var calls int
		wsrpcClient := &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (*pb.LatestReportResponse, error) {
				calls++
				assert.Equal(t, req, in)
				return resp, nil
			},
		}

		conn := &mocks.MockConn{
			Ready: true,
		}
		c := newClient(lggr, csakey.KeyV2{}, nil, "", cacheSet)
		c.conn = conn
		c.rawClient = wsrpcClient

		// simulate start without dialling
		require.NoError(t, c.StartOnce("Mock WSRPC Client", func() error { return nil }))
		var err error
		servicetest.Run(t, cacheSet)
		c.cache, err = cacheSet.Get(ctx, c)
		require.NoError(t, err)

		for i := 0; i < 5; i++ {
			r, err := c.LatestReport(ctx, req)

			require.NoError(t, err)
			assert.Equal(t, resp, r)
		}
		assert.Equal(t, 5, calls, "expected 5 calls to LatestReport but it was called %d times", calls)
	})

	t.Run("with caching", func(t *testing.T) {
		req := &pb.LatestReportRequest{}
		const neverExpireTTL = 1000 * time.Hour // some massive value that will never expire during a test
		cacheSet := cache.NewCacheSet(lggr, cache.Config{LatestReportTTL: neverExpireTTL})
		resp := &pb.LatestReportResponse{}

		var calls int
		wsrpcClient := &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (*pb.LatestReportResponse, error) {
				calls++
				assert.Equal(t, req, in)
				return resp, nil
			},
		}

		conn := &mocks.MockConn{
			Ready: true,
		}
		c := newClient(lggr, csakey.KeyV2{}, nil, "", cacheSet)
		c.conn = conn
		c.rawClient = wsrpcClient

		// simulate start without dialling
		require.NoError(t, c.StartOnce("Mock WSRPC Client", func() error { return nil }))
		var err error
		servicetest.Run(t, cacheSet)
		c.cache, err = cacheSet.Get(ctx, c)
		require.NoError(t, err)

		for i := 0; i < 5; i++ {
			r, err := c.LatestReport(ctx, req)

			require.NoError(t, err)
			assert.Equal(t, resp, r)
		}
		assert.Equal(t, 1, calls, "expected only 1 call to LatestReport but it was called %d times", calls)
	})
}
