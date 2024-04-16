package test

import (
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	loopnet "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	loopnettest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net/test"
)

type Client interface {
	io.Closer
	goplugin.GRPCClientConn
}

type Server any

// GRPCScaffold is a test scaffold for grpc tests
type GRPCScaffold[C Client, S Server] struct {
	t *testing.T

	server S
	client C

	grpcServer *grpc.Server
}

// Caller must either call Close() on the returned GRPCScaffold
// or manually close the Client. some implementations in our suite of LOOPPs release server resources on client.Close()
// and in that case, the caller should call Client().Close() and test that the server resources are released
func (t *GRPCScaffold[T, S]) Close() {
	require.NoError(t.t, t.client.Close(), "failed to close client")
}

func (t *GRPCScaffold[T, S]) Client() T {
	return t.client
}

func (t *GRPCScaffold[T, S]) Server() S {
	return t.server
}

func NewGRPCScaffold[T Client, S any](t *testing.T, serverFn SetupGRPCServer[S], clientFn SetupGRPCClient[T]) *GRPCScaffold[T, S] {
	lis := tcpListener(t)
	grpcServer := grpc.NewServer()
	t.Cleanup(grpcServer.Stop)

	lggr := logger.Test(t)
	broker := &loopnettest.Broker{T: t}
	brokerExt := &loopnet.BrokerExt{
		Broker:       broker,
		BrokerConfig: loopnet.BrokerConfig{Logger: lggr, StopCh: make(chan struct{})},
	}

	s := serverFn(t, grpcServer, brokerExt)
	go func() {
		// the cleanup call to grpcServer.Stop will unblock this
		require.NoError(t, grpcServer.Serve(lis))
	}()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "failed to dial %s", lis.Addr().String())
	t.Cleanup(func() { require.NoError(t, conn.Close(), "failed to close connection") })

	client := clientFn(brokerExt, conn)

	return &GRPCScaffold[T, S]{t: t, server: s, client: client, grpcServer: grpcServer}
}

func tcpListener(t *testing.T) net.Listener {
	port := freeport.GetOne(t)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	require.NoError(t, err, "failed to listen on port %d", port)
	t.Cleanup(func() { lis.Close() })
	return lis
}

// SetupGRPCServer is a function that sets up a grpc server with a given broker
// typical and expected usage is to instantiate a grpc server implementation with a static test interface implementation
// and then register that grpc server
// e.g.
// ```
//
//	func setupCCIPCommitProviderGRPCServer(t *testing.T, s *grpc.Server, b *loopnet.BrokerExt) *grpc.Server {
//	  commitProvider := ccip.NewCommitProviderServer(CommitProvider, b)
//	  ccippb.RegisterCommitCustomHandlersServer(s, commitProvider)
//	  return s
//	}
//
// ```
type SetupGRPCServer[S any] func(t *testing.T, s *grpc.Server, b *loopnet.BrokerExt) S

// SetupGRPCClient is a function that sets up a grpc client with a given broker and connection
// analogous to SetupGRPCServer. Typically it is implemented as a light wrapper around the grpc client constructor
type SetupGRPCClient[T Client] func(b *loopnet.BrokerExt, conn grpc.ClientConnInterface) T

// MockDep is a mock dependency that can be used to test that a grpc client closes its dependencies
// to be used in tests that require a grpc client to close its dependencies
type MockDep struct {
	closeCalled bool
}

func (m *MockDep) Close() error {
	m.closeCalled = true
	return nil
}

func (m *MockDep) IsClosed() bool {
	return m.closeCalled
}
