package test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/ccip"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestStaticTokenPool(t *testing.T) {
	t.Parallel()

	// static test implementation is self consistent
	ctx := context.Background()
	assert.NoError(t, TokenPoolBatchedReader.Evaluate(ctx, TokenPoolBatchedReader))

	// error when the test implementation is evaluates something that differs from the static implementation
	botched := TokenPoolBatchedReader
	botched.getInboundTokenPoolRateLimitsRequest = []cciptypes.Address{"not the right request"}
	err := TokenPoolBatchedReader.Evaluate(ctx, botched)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not the right request")
}

func TestTokenPoolGRPC(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	// create a price registry server
	port := freeport.GetOne(t)
	addr := fmt.Sprintf("localhost:%d", port)
	lis, err := net.Listen("tcp", addr)
	require.NoError(t, err, "failed to listen on port %d", port)
	t.Cleanup(func() { lis.Close() })
	// we explicitly stop the server later, do not add a cleanup function here
	testServer := grpc.NewServer()
	// handle client close and server stop
	shutdown := make(chan struct{})
	closer := &serviceCloser{closeFn: func() error { close(shutdown); return nil }}

	tokenPool := ccip.NewTokenPoolBatchedReaderGRPCServer(TokenPoolBatchedReader)
	require.NoError(t, err)
	tokenPool = tokenPool.AddDep(closer)

	ccippb.RegisterTokenPoolBatcherReaderServer(testServer, tokenPool)
	// start the server and shutdown handler
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, testServer.Serve(lis))
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-shutdown
		t.Log("shutting down server")
		testServer.Stop()
	}()
	// create a token data client
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "failed to dial %s", addr)
	t.Cleanup(func() { conn.Close() })
	client := ccip.NewTokenPoolBatchedReaderGRPCClient(conn)

	// test the client
	roundTripTokenPoolTests(ctx, t, client)
	// closing the client executes the shutdown callback
	// which stops the server.  the wg.Wait() below ensures
	// that the server has stopped, which is what we care about.
	cerr := client.Close()
	require.NoError(t, cerr, "failed to close client %T, %v", cerr, status.Code(cerr))
	wg.Wait()
}

func roundTripTokenPoolTests(ctx context.Context, t *testing.T, client cciptypes.TokenPoolBatchedReader) {
	t.Helper()
	// test read token data
	limits, err := client.GetInboundTokenPoolRateLimits(ctx, TokenPoolBatchedReader.getInboundTokenPoolRateLimitsRequest)
	require.NoError(t, err)
	assert.Equal(t, TokenPoolBatchedReader.getInboundTokenPoolRateLimitsResponse, limits)
}
