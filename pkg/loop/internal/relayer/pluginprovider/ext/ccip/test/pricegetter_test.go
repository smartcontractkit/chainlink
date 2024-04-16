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

	scaffold := looptest.NewGRPCScaffold(t, setupPriceGetterServer, setupPriceGetterClient)
	roundTripPriceGetterTests(t, scaffold.Client())
	// price getter implements dependency management, test that it closes properly
	t.Run("Dependency management", func(t *testing.T) {
		d := &looptest.MockDep{}
		scaffold.Server().AddDep(d)
		assert.False(t, d.IsClosed())
		scaffold.Client().Close()
		assert.True(t, d.IsClosed())
	})
}

func roundTripPriceGetterTests(t *testing.T, client cciptypes.PriceGetter) {
	t.Run("FilterConfiguredTokens", func(t *testing.T) {
		// test token is configured
		configuredTokens, unconfiguredTokens, err := client.FilterConfiguredTokens(tests.Context(t), PriceGetter.config.Addresses)
		require.NoError(t, err)
		assert.Equal(t, PriceGetter.config.Addresses, configuredTokens)
		assert.Equal(t, []cciptypes.Address{}, unconfiguredTokens)

		var unconfTk cciptypes.Address = "JK"
		unconfTks := []cciptypes.Address{unconfTk}
		configuredTokens2, unconfiguredTokens2, err := client.FilterConfiguredTokens(tests.Context(t), unconfTks)
		require.NoError(t, err)
		assert.Equal(t, []cciptypes.Address{}, configuredTokens2)
		assert.Equal(t, unconfTks, unconfiguredTokens2)
	})
	t.Run("TokenPricesUSD", func(t *testing.T) {
		// test token prices
		prices, err := client.TokenPricesUSD(tests.Context(t), PriceGetter.config.Addresses)
		require.NoError(t, err)
		assert.Equal(t, PriceGetter.config.Prices, prices)
	})
}

func setupPriceGetterServer(t *testing.T, s *grpc.Server, b *loopnet.BrokerExt) *ccip.PriceGetterGRPCServer {
	priceGetter := ccip.NewPriceGetterGRPCServer(PriceGetter)
	ccippb.RegisterPriceGetterServer(s, priceGetter)
	return priceGetter
}

func setupPriceGetterClient(b *loopnet.BrokerExt, conn grpc.ClientConnInterface) *ccip.PriceGetterGRPCClient {
	return ccip.NewPriceGetterGRPCClient(conn)
}

var _ looptest.SetupGRPCServer[*ccip.PriceGetterGRPCServer] = setupPriceGetterServer
var _ looptest.SetupGRPCClient[*ccip.PriceGetterGRPCClient] = setupPriceGetterClient
