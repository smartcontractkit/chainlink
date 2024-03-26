package test

import (
	"context"
	"errors"
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

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	loopnet "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	loopnettest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net/test"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestStaticCommitProvider(t *testing.T) {
	ctx := tests.Context(t)
	t.Run("Self consistent Evaluate", func(t *testing.T) {
		t.Parallel()
		// static test implementation is self consistent
		assert.NoError(t, CommitProvider.Evaluate(ctx, CommitProvider))

		// error when the test implementation evaluates something that differs from form itself
		botched := CommitProvider
		botched.priceRegistryReader = staticPriceRegistryReader{}
		err := CommitProvider.Evaluate(ctx, botched)
		require.Error(t, err)
		var evalErr evaluationError
		require.True(t, errors.As(err, &evalErr), "expected error to be an evaluationError")
		assert.Equal(t, priceRegistryComponent, evalErr.component)
	})
	t.Run("Self consistent AssertEqual", func(t *testing.T) {
		// no parallel because the AssertEqual is parallel
		CommitProvider.AssertEqual(ctx, t, CommitProvider)
	})
}

func TestCommitProviderGRPC(t *testing.T) {
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

	/*
		// handle client close and server stop
		shutdown := make(chan struct{})
		closer := &serviceCloser{closeFn: func() error { close(shutdown); return nil }}
	*/
	lggr := logger.Test(t)
	broker := &loopnettest.Broker{T: t}
	brokerExt := &loopnet.BrokerExt{
		Broker:       broker,
		BrokerConfig: loopnet.BrokerConfig{Logger: lggr, StopCh: make(chan struct{})},
	}
	commitProvider := ccip.NewCommitProviderServer(CommitProvider, brokerExt)
	require.NoError(t, err)

	ccippb.RegisterCommitCustomHandlersServer(testServer, commitProvider)
	// start the server and shutdown handler
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, testServer.Serve(lis))
	}()
	// create a price registry client
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "failed to dial %s", addr)
	t.Cleanup(func() { conn.Close() })
	client := ccip.NewCommitProviderClient(brokerExt, conn)

	roundTripCommitProviderTests(ctx, t, client)
	// closing the client executes the shutdown callback
	// which stops the server.  the wg.Wait() below ensures
	// that the server has stopped, which is what we care about.
	cerr := client.Close()
	require.NoError(t, cerr, "failed to close client %T, %v", cerr, status.Code(cerr))
	testServer.Stop()
	wg.Wait()
}

func roundTripCommitProviderTests(ctx context.Context, t *testing.T, client types.CCIPCommitProvider) {
	t.Run("CommitStore", func(t *testing.T) {
		commitClient, err := client.NewCommitStoreReader(ctx, "ignored")
		require.NoError(t, err)
		roundTripCommitStoreTests(ctx, t, commitClient)
		require.NoError(t, commitClient.Close())
	})

	t.Run("OffRamp", func(t *testing.T) {
		offRampClient, err := client.NewOffRampReader(ctx, "ignored")
		require.NoError(t, err)
		roundTripOffRampTests(ctx, t, offRampClient)
		require.NoError(t, offRampClient.Close())
	})

	t.Run("OnRamp", func(t *testing.T) {
		onRampClient, err := client.NewOnRampReader(ctx, "ignored")
		require.NoError(t, err)
		roundTripOnRampTests(ctx, t, onRampClient)
		require.NoError(t, onRampClient.Close())
	})

	t.Run("PriceGetter", func(t *testing.T) {
		priceGetterClient, err := client.NewPriceGetter(ctx)
		require.NoError(t, err)
		roundTripPriceGetterTests(ctx, t, priceGetterClient)
		require.NoError(t, priceGetterClient.Close())
	})

	t.Run("PriceRegistry", func(t *testing.T) {
		priceRegistryClient, err := client.NewPriceRegistryReader(ctx, "ignored")
		require.NoError(t, err)
		roundTripPriceRegistryTests(ctx, t, priceRegistryClient)
		require.NoError(t, priceRegistryClient.Close())
	})
}
