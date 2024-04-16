package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	loopnet "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/ccip"
	looptest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
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
	scaffold := looptest.NewGRPCScaffold(t, setupTokenPoolServer, setupTokenPoolClient)
	roundTripTokenPoolTests(t, scaffold.Client())
	// token pool implements dependency management, test that it closes properly
	t.Run("Dependency management", func(t *testing.T) {
		d := &looptest.MockDep{}
		scaffold.Server().AddDep(d)
		assert.False(t, d.IsClosed())
		scaffold.Client().Close()
		assert.True(t, d.IsClosed())
	})
}

func roundTripTokenPoolTests(t *testing.T, client cciptypes.TokenPoolBatchedReader) {
	t.Helper()
	// test read token data
	limits, err := client.GetInboundTokenPoolRateLimits(tests.Context(t), TokenPoolBatchedReader.getInboundTokenPoolRateLimitsRequest)
	require.NoError(t, err)
	assert.Equal(t, TokenPoolBatchedReader.getInboundTokenPoolRateLimitsResponse, limits)
}

func setupTokenPoolServer(t *testing.T, s *grpc.Server, b *loopnet.BrokerExt) *ccip.TokenPoolBatchedReaderGRPCServer {
	tokenPool := ccip.NewTokenPoolBatchedReaderGRPCServer(TokenPoolBatchedReader)
	ccippb.RegisterTokenPoolBatcherReaderServer(s, tokenPool)
	return tokenPool
}

func setupTokenPoolClient(b *loopnet.BrokerExt, conn grpc.ClientConnInterface) *ccip.TokenPoolBatchedReaderGRPCClient {
	return ccip.NewTokenPoolBatchedReaderGRPCClient(conn)
}

var _ looptest.SetupGRPCServer[*ccip.TokenPoolBatchedReaderGRPCServer] = setupTokenPoolServer
var _ looptest.SetupGRPCClient[*ccip.TokenPoolBatchedReaderGRPCClient] = setupTokenPoolClient
