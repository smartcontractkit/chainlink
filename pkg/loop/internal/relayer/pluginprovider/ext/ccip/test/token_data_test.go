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

func TestStaticTokenData(t *testing.T) {
	t.Parallel()

	// static test implementation is self consistent
	ctx := context.Background()
	assert.NoError(t, TokenDataReader.Evaluate(ctx, TokenDataReader))

	// error when the test implementation is evaluates something that differs from the static implementation
	botched := TokenDataReader
	botched.readTokenDataResponse = []byte("not the right response")
	err := TokenDataReader.Evaluate(ctx, botched)
	require.Error(t, err)
}

func TestTokenDataGRPC(t *testing.T) {
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

	tokenData := ccip.NewTokenDataReaderGRPCServer(TokenDataReader)
	require.NoError(t, err)
	tokenData = tokenData.AddDep(closer)

	ccippb.RegisterTokenDataReaderServer(testServer, tokenData)
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
	client := ccip.NewTokenDataReaderGRPCClient(conn)

	// test the client
	roundTripTokenDataTests(ctx, t, client)
	// closing the client executes the shutdown callback
	// which stops the server.  the wg.Wait() below ensures
	// that the server has stopped, which is what we care about.
	cerr := client.Close()
	require.NoError(t, cerr, "failed to close client %T, %v", cerr, status.Code(cerr))
	wg.Wait()
}

func roundTripTokenDataTests(ctx context.Context, t *testing.T, client cciptypes.TokenDataReader) {
	t.Helper()
	// test read token data
	tokenData, err := client.ReadTokenData(ctx, TokenDataReader.readTokenDataRequest.msg, TokenDataReader.readTokenDataRequest.tokenIndex)
	require.NoError(t, err)
	assert.Equal(t, TokenDataReader.readTokenDataResponse, tokenData)
}
