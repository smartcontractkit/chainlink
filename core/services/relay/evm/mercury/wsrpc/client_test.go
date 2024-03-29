package wsrpc

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/wsrpc"
)

// simulate start without dialling
func simulateStart(ctx context.Context, t *testing.T, c *client) {
	require.NoError(t, c.StartOnce("Mock WSRPC Client", func() (err error) {
		c.cache, err = c.cacheSet.Get(ctx, c)
		return err
	}))
}

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

	t.Run("sends on reset channel after MaxConsecutiveRequestFailures timed out transmits", func(t *testing.T) {
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
		c := newClient(lggr, csakey.KeyV2{}, nil, "", noopCacheSet, nil)
		c.conn = conn
		c.rawClient = wsrpcClient
		require.NoError(t, c.StartOnce("Mock WSRPC Client", func() error { return nil }))
		for i := 1; i < MaxConsecutiveRequestFailures; i++ {
			_, err := c.Transmit(ctx, req)
			require.EqualError(t, err, "context deadline exceeded")
		}
		assert.Equal(t, MaxConsecutiveRequestFailures-1, calls)
		select {
		case <-c.chResetTransport:
			t.Fatal("unexpected send on chResetTransport")
		default:
		}
		_, err := c.Transmit(ctx, req)
		require.EqualError(t, err, "context deadline exceeded")
		assert.Equal(t, MaxConsecutiveRequestFailures, calls)
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
			assert.Equal(t, MaxConsecutiveRequestFailures+1, calls)
			assert.Equal(t, 0, int(c.consecutiveTimeoutCnt.Load()))
		})

		t.Run("doesn't block in case channel is full", func(t *testing.T) {
			transmitErr = context.DeadlineExceeded
			c.chResetTransport = nil // simulate full channel
			for i := 0; i < MaxConsecutiveRequestFailures; i++ {
				_, err := c.Transmit(ctx, req)
				require.EqualError(t, err, "context deadline exceeded")
			}
		})
	})
}

func Test_Client_LatestReport(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	cacheReads := 5

	tests := []struct {
		name          string
		ttl           time.Duration
		expectedCalls int
	}{
		{
			name:          "with cache disabled",
			ttl:           0,
			expectedCalls: 5,
		},
		{
			name:          "with cache enabled",
			ttl:           1000 * time.Hour, //some large value that will never expire during a test
			expectedCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.LatestReportRequest{}

			cacheSet := cache.NewCacheSet(lggr, cache.Config{LatestReportTTL: tt.ttl})

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
			c := newClient(lggr, csakey.KeyV2{}, nil, "", cacheSet, nil)
			c.conn = conn
			c.rawClient = wsrpcClient

			servicetest.Run(t, cacheSet)
			simulateStart(ctx, t, c)

			for i := 0; i < cacheReads; i++ {
				r, err := c.LatestReport(ctx, req)

				require.NoError(t, err)
				assert.Equal(t, resp, r)
			}
			assert.Equal(t, tt.expectedCalls, calls, "expected %d calls to LatestReport but it was called %d times", tt.expectedCalls, calls)
		})
	}
}

type TestServer interface {
	Serve(lis net.Listener) error
	Stop()
}

type WrappedWsrpcServer struct {
	*wsrpc.Server
}

func (s *WrappedWsrpcServer) Serve(lis net.Listener) error {
	s.Server.Serve(lis)
	return nil
}

func NewWrappedWsrpcServer() TestServer {
	return &WrappedWsrpcServer{wsrpc.NewServer()}
}

// Tests that when start is called, the appropriate type of connection is made
func Test_Start_Dial(t *testing.T) {
	wsrpcName := "WSRPC"
	grpcName := "GRPC"
	tests := []struct {
		name string
		tlsCertFile *string
		server TestServer
	}{
		{
			name: wsrpcName,
			tlsCertFile: nil,
			server: NewWrappedWsrpcServer(),
		},
		{
			name: grpcName,
			tlsCertFile: ptr("./fixtures/domain.pem"), // TODO: dialing with insecure creds for now
			server: grpc.NewServer(),
		},
	}
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)

	for _, tt := range tests {

		t.Logf("Running test for %s", tt.name)
		port := freeport.GetOne(t)
		addr := fmt.Sprintf("127.0.0.1:%v", port)

		// Set up client
		c := newClient(lggr, csakey.KeyV2{}, nil, addr, newNoopCacheSet(), nil)

		// Set up server
		lis, err := net.Listen("tcp",  addr)
		require.NoError(t, err)
		s := tt.server
		go s.Serve(lis)
		t.Cleanup(s.Stop)

		// Start client
		err = c.Start(ctx)
		require.NoError(t, err)
		t.Cleanup(func() {c.Close()})

		// Validate connection type
		switch tt.name {
		case wsrpcName:
			_, ok := c.conn.(*wsrpc.ClientConn); 
			if !ok {
				t.Fatalf("expected wsrpc.ClientConn, got %T", c.conn)
			}
		case grpcName:
			_, ok:= c.conn.(*AdapatedGrpcClientConn)
			if !ok { // will always fail for now while I'm hacking tlsCertFile to be null
				t.Fatalf("expected AdaptedGrpcClientConn, got %T", c.conn)
			}
		}

		}
}

// Tests that when start is called, a the appropriate type of connection is made
// spin up server, make dial, validate type of connection is correct
func Test_Start_Dial_GRPC(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)

	port := freeport.GetOne(t)
	addr := fmt.Sprintf("127.0.0.1:%v", port)

	// Set up client
	c := newClient(lggr, csakey.KeyV2{}, nil, addr, newNoopCacheSet(), nil)

	// Set up server
	lis, err := net.Listen("tcp",  addr)
	require.NoError(t, err)
	s := grpc.NewServer()
	err = s.Serve(lis)
	require.NoError(t, err)

	// Start client
	err = c.Start(ctx)
	require.NoError(t, err)

	// Validate connection type
	_, ok:= c.conn.(*AdapatedGrpcClientConn)
	if !ok {
		t.Fatalf("expected AdaptedGrpcClientConn, got %T", c.conn)
	}

}

// TODO:
// * figure out if I want a mode where the client can dial the grpc server without a cert bundle


func ptr[T any](t T) *T { return &t }