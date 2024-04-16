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

	scaffold := looptest.NewGRPCScaffold(t, setupTokenDataServer, setupTokenDataClient)
	roundTripTokenDataTests(t, scaffold.Client())
	// token data reader implements dependency management, test that it closes properly
	t.Run("Dependency management", func(t *testing.T) {
		d := &looptest.MockDep{}
		scaffold.Server().AddDep(d)
		assert.False(t, d.IsClosed())
		scaffold.Client().Close()
		assert.True(t, d.IsClosed())
	})
}

func roundTripTokenDataTests(t *testing.T, client cciptypes.TokenDataReader) {
	t.Helper()
	// test read token data
	tokenData, err := client.ReadTokenData(tests.Context(t), TokenDataReader.readTokenDataRequest.msg, TokenDataReader.readTokenDataRequest.tokenIndex)
	require.NoError(t, err)
	assert.Equal(t, TokenDataReader.readTokenDataResponse, tokenData)
}

func setupTokenDataServer(t *testing.T, s *grpc.Server, b *loopnet.BrokerExt) *ccip.TokenDataReaderGRPCServer {
	tokenData := ccip.NewTokenDataReaderGRPCServer(TokenDataReader)
	ccippb.RegisterTokenDataReaderServer(s, tokenData)
	return tokenData
}

func setupTokenDataClient(b *loopnet.BrokerExt, conn grpc.ClientConnInterface) *ccip.TokenDataReaderGRPCClient {
	return ccip.NewTokenDataReaderGRPCClient(conn)
}

var _ looptest.SetupGRPCServer[*ccip.TokenDataReaderGRPCServer] = setupTokenDataServer
var _ looptest.SetupGRPCClient[*ccip.TokenDataReaderGRPCClient] = setupTokenDataClient
