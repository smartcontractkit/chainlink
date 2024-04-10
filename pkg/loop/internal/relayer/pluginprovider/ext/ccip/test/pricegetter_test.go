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

func Test_staticPriceGetter_Evaluate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	assert.NoError(t, PriceGetter.Evaluate(ctx, PriceGetter))

	botched := PriceGetter
	botched.config.Addresses = []cciptypes.Address{"wrong token"}
	assert.Error(t, PriceGetter.Evaluate(ctx, botched))
}

func TestPriceGetterGRPC(t *testing.T) {
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

	priceGetter := ccip.NewPriceGetterGRPCServer(PriceGetter)
	require.NoError(t, err)
	priceGetter = priceGetter.AddDep(closer)

	ccippb.RegisterPriceGetterServer(testServer, priceGetter)
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
	client := ccip.NewPriceGetterGRPCClient(conn)

	// test the client
	roundTripPriceGetterTests(ctx, t, client)
	// closing the client executes the shutdown callback
	// which stops the server.  the wg.Wait() below ensures
	// that the server has stopped, which is what we care about.
	cerr := client.Close()
	require.NoError(t, cerr, "failed to close client %T, %v", cerr, status.Code(cerr))
	wg.Wait()
}

func roundTripPriceGetterTests(ctx context.Context, t *testing.T, client cciptypes.PriceGetter) {
	t.Run("FilterConfiguredTokens", func(t *testing.T) {
		// test token is configured
		configuredTokens, unconfiguredTokens, err := client.FilterConfiguredTokens(ctx, PriceGetter.config.Addresses)
		require.NoError(t, err)
		assert.Equal(t, PriceGetter.config.Addresses, configuredTokens)
		assert.Equal(t, []cciptypes.Address{}, unconfiguredTokens)

		var unconfTk cciptypes.Address = "JK"
		unconfTks := []cciptypes.Address{unconfTk}
		configuredTokens2, unconfiguredTokens2, err := client.FilterConfiguredTokens(ctx, unconfTks)
		require.NoError(t, err)
		assert.Equal(t, []cciptypes.Address{}, configuredTokens2)
		assert.Equal(t, unconfTks, unconfiguredTokens2)
	})
	t.Run("TokenPricesUSD", func(t *testing.T) {
		// test token prices
		prices, err := client.TokenPricesUSD(ctx, PriceGetter.config.Addresses)
		require.NoError(t, err)
		assert.Equal(t, PriceGetter.config.Prices, prices)
	})
}
