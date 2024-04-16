package test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	loopnet "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/ccip"
	looptest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestStaticCommitProvider(t *testing.T) {
	t.Run("Self consistent Evaluate", func(t *testing.T) {
		t.Parallel()
		ctx := tests.Context(t)
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
		CommitProvider.AssertEqual(tests.Context(t), t, CommitProvider)
	})
}

func TestCommitProviderGRPC(t *testing.T) {
	t.Parallel()

	grpcScaffold := looptest.NewGRPCScaffold(t, setupCommitProviderServer, ccip.NewCommitProviderClient)
	t.Cleanup(grpcScaffold.Close)
	roundTripCommitProviderTests(t, grpcScaffold.Client())
}

func roundTripCommitProviderTests(t *testing.T, client types.CCIPCommitProvider) {
	t.Run("CommitStore", func(t *testing.T) {
		commitClient, err := client.NewCommitStoreReader(tests.Context(t), "ignored")
		require.NoError(t, err)
		roundTripCommitStoreTests(t, commitClient)
		require.NoError(t, commitClient.Close())
	})

	t.Run("OffRamp", func(t *testing.T) {
		offRampClient, err := client.NewOffRampReader(tests.Context(t), "ignored")
		require.NoError(t, err)
		roundTripOffRampTests(t, offRampClient)
		require.NoError(t, offRampClient.Close())
	})

	t.Run("OnRamp", func(t *testing.T) {
		onRampClient, err := client.NewOnRampReader(tests.Context(t), "ignored")
		require.NoError(t, err)
		roundTripOnRampTests(t, onRampClient)
		require.NoError(t, onRampClient.Close())
	})

	t.Run("PriceGetter", func(t *testing.T) {
		priceGetterClient, err := client.NewPriceGetter(tests.Context(t))
		require.NoError(t, err)
		roundTripPriceGetterTests(t, priceGetterClient)
		require.NoError(t, priceGetterClient.Close())
	})

	t.Run("PriceRegistry", func(t *testing.T) {
		priceRegistryClient, err := client.NewPriceRegistryReader(tests.Context(t), "ignored")
		require.NoError(t, err)
		roundTripPriceRegistryTests(t, priceRegistryClient)
		require.NoError(t, priceRegistryClient.Close())
	})
}

func setupCommitProviderServer(t *testing.T, s *grpc.Server, b *loopnet.BrokerExt) *ccip.CommitProviderServer {
	commitProvider := ccip.NewCommitProviderServer(CommitProvider, b)
	ccippb.RegisterCommitCustomHandlersServer(s, commitProvider)
	return commitProvider
}

var _ looptest.SetupGRPCServer[*ccip.CommitProviderServer] = setupCommitProviderServer
var _ looptest.SetupGRPCClient[*ccip.CommitProviderClient] = ccip.NewCommitProviderClient
