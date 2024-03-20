package test

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"sync"
	"testing"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/ccip"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestStaticOnRamp(t *testing.T) {
	t.Parallel()

	// static test implementation is self consistent
	ctx := context.Background()
	assert.NoError(t, OnRampReader.Evaluate(ctx, OnRampReader))

	// error when the test implementation is evaluates something that differs from the static implementation
	botched := OnRampReader
	botched.addressResponse = "not the right address"
	err := OnRampReader.Evaluate(ctx, botched)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not the right address")
}

func TestOnRampGRPC(t *testing.T) {
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

	onRamp := ccip.NewOnRampReaderGRPCServer(OnRampReader)
	require.NoError(t, err)
	onRamp = onRamp.AddDep(closer)

	ccippb.RegisterOnRampReaderServer(testServer, onRamp)
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
	// create a price registry client
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "failed to dial %s", addr)
	t.Cleanup(func() { conn.Close() })
	client := ccip.NewOnRampReaderGRPCClient(conn)

	// test the client
	roundTripOnRampTests(ctx, t, client)
	// closing the client executes the shutdown callback
	// which stops the server.  the wg.Wait() below ensures
	// that the server has stopped, which is what we care about.
	cerr := client.Close()
	require.NoError(t, cerr, "failed to close client %T, %v", cerr, status.Code(cerr))
	wg.Wait()
}

func roundTripOnRampTests(ctx context.Context, t *testing.T, client *ccip.OnRampReaderGRPCClient) {
	// test the client

	t.Run("Address", func(t *testing.T) {
		got, err := client.Address(ctx)
		require.NoError(t, err)
		assert.Equal(t, OnRampReader.addressResponse, got)
	})

	t.Run("GetDynamicConfig", func(t *testing.T) {
		got, err := client.GetDynamicConfig(ctx)
		require.NoError(t, err)
		assert.Equal(t, OnRampReader.dynamicConfigResponse, got)
	})

	t.Run("GetSendRequestsBetweenSeqNums", func(t *testing.T) {
		got, err := client.GetSendRequestsBetweenSeqNums(ctx, OnRampReader.getSendRequestsBetweenSeqNums.SeqNumMin, OnRampReader.getSendRequestsBetweenSeqNums.SeqNumMax, OnRampReader.getSendRequestsBetweenSeqNums.Finalized)
		require.NoError(t, err)
		if !reflect.DeepEqual(OnRampReader.getSendRequestsBetweenSeqNumsResponse.EVM2EVMMessageWithTxMeta, got) {
			t.Errorf("expected %v, got %v", OnRampReader.getSendRequestsBetweenSeqNumsResponse.EVM2EVMMessageWithTxMeta, got)
		}
	})

	t.Run("RouterAddress", func(t *testing.T) {
		got, err := client.RouterAddress(ctx)
		require.NoError(t, err)
		assert.Equal(t, OnRampReader.routerResponse, got)
	})

	// TODO: BCF-3106 implement the new methods
}
