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

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/ccip"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestStaticPriceRegistry(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	// static test implementation is self consistent
	assert.NoError(t, PriceRegistryReader.Evaluate(ctx, PriceRegistryReader))

	// error when the test implementation is evaluates something that differs from the static implementation
	botched := PriceRegistryReader
	botched.addressResponse = "not what we expect"
	err := PriceRegistryReader.Evaluate(ctx, botched)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not what we expect")
}

func TestPriceRegistryGRPC(t *testing.T) {
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
	priceRegistry := ccip.NewPriceRegistryGRPCServer(PriceRegistryReader).AddDep(closer)

	ccippb.RegisterPriceRegistryReaderServer(testServer, priceRegistry)
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
	client := ccip.NewPriceRegistryGRPCClient(conn)

	// test the client
	roundTripPriceRegistryTests(ctx, t, client)
	// closing the client executes the shutdown callback
	// which stops the server.  the wg.Wait() below ensures
	// that the server has stopped, which is what we care about.
	cerr := client.Close()
	require.NoError(t, cerr, "failed to close client %T, %v", cerr, status.Code(cerr))
	wg.Wait()
}

// roundTripPriceRegistryTests tests the round trip of the client<->server.
// it should exercise all the methods of the client.
// do not add client.Close to this test, test that from the driver test
func roundTripPriceRegistryTests(ctx context.Context, t *testing.T, client *ccip.PriceRegistryGRPCClient) {
	t.Run("Address", func(t *testing.T) {
		address, err := client.Address(ctx)
		require.NoError(t, err)
		assert.Equal(t, PriceRegistryReader.addressResponse, address)
	})

	t.Run("GetFeeTokens", func(t *testing.T) {
		price, err := client.GetFeeTokens(ctx)
		require.NoError(t, err)
		assert.Equal(t, PriceRegistryReader.getFeeTokensResponse, price)
	})

	t.Run("GetGasPriceUpdatesCreatedAfter", func(t *testing.T) {
		price, err := client.GetGasPriceUpdatesCreatedAfter(ctx,
			PriceRegistryReader.getGasPriceUpdatesCreatedAfterRequest.chainSelector,
			PriceRegistryReader.getGasPriceUpdatesCreatedAfterRequest.ts,
			PriceRegistryReader.getGasPriceUpdatesCreatedAfterRequest.confirmations,
		)
		require.NoError(t, err)
		assert.Equal(t, PriceRegistryReader.getGasPriceUpdatesCreatedAfterResponse, price)
	})

	t.Run("GetTokenPriceUpdatesCreatedAfter", func(t *testing.T) {
		price, err := client.GetTokenPriceUpdatesCreatedAfter(ctx,
			PriceRegistryReader.getTokenPriceUpdatesCreatedAfterRequest.ts,
			PriceRegistryReader.getTokenPriceUpdatesCreatedAfterRequest.confirmations,
		)
		require.NoError(t, err)
		assert.Equal(t, PriceRegistryReader.getTokenPriceUpdatesCreatedAfterResponse, price)
	})

	t.Run("GetTokenPrices", func(t *testing.T) {
		price, err := client.GetTokenPrices(ctx, PriceRegistryReader.getTokenPricesRequest)
		require.NoError(t, err)
		assert.Equal(t, PriceRegistryReader.getTokenPricesResponse, price)
	})

	t.Run("GetTokensDecimals", func(t *testing.T) {
		price, err := client.GetTokensDecimals(ctx, PriceRegistryReader.getTokensDecimalsRequest)
		require.NoError(t, err)
		assert.Equal(t, PriceRegistryReader.getTokensDecimalsResponse, price)
	})
}
