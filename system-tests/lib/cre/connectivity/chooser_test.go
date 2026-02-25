package connectivity

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveSamePlacementUsesInternal(t *testing.T) {
	r, err := Resolve(PlacementLocal, PlacementLocal, EndpointPair{
		Name:     "evm-rpc",
		Internal: "http://anvil:8545",
		External: "http://10.0.0.1:8545",
	})
	require.NoError(t, err, "expected resolve to succeed")
	require.Equal(t, "http://anvil:8545", r.URL, "unexpected URL resolution")
	require.Equal(t, "internal", r.SelectedKind, "unexpected endpoint kind")
	require.False(t, r.RequiresBridge, "did not expect bridge requirement for same placement")
}

func TestResolveRemoteToLocalRequiresBridge(t *testing.T) {
	r, err := Resolve(PlacementRemote, PlacementLocal, EndpointPair{
		Name:     "jd-grpc",
		Internal: "jd:14231",
		External: "127.0.0.1:14231",
	})
	require.NoError(t, err, "expected resolve to succeed")
	require.True(t, r.RequiresBridge, "expected bridge requirement for remote caller to local target")
	require.Equal(t, 14231, r.BridgePort, "unexpected bridge port")
}

func TestResolveCrossPlacementLocalToRemoteUsesExternalWithoutBridge(t *testing.T) {
	r, err := Resolve(PlacementLocal, PlacementRemote, EndpointPair{
		Name:     "gateway",
		Internal: "ws://gateway-node:5003/node",
		External: "ws://203.0.113.10:5003/node",
	})
	require.NoError(t, err, "expected cross-placement resolve to succeed")
	require.Equal(t, "external", r.SelectedKind, "expected external URL for cross-placement")
	require.Equal(t, "ws://203.0.113.10:5003/node", r.URL, "unexpected cross-placement URL")
	require.False(t, r.RequiresBridge, "local caller to remote target should not require bridge")
}

func TestResolveAndEnsureReachableCallsEnsurer(t *testing.T) {
	called := false
	r, err := ResolveAndEnsureReachable(context.Background(), PlacementRemote, PlacementLocal, EndpointPair{
		Name:     "jd-grpc",
		Internal: "jd:14231",
		External: "127.0.0.1:14231",
	}, func(_ context.Context, endpoint EndpointPair, port int) error {
		called = true
		require.Equal(t, "jd-grpc", endpoint.Name, "unexpected endpoint name in bridge callback")
		require.Equal(t, 14231, port, "unexpected port in bridge callback")
		return nil
	})
	require.NoError(t, err, "expected resolve+ensure to succeed")
	require.True(t, called, "expected bridge ensurer to be called")
	require.Equal(t, "127.0.0.1:14231", r.URL, "unexpected resolution URL")
}

func TestResolveAndEnsureReachableFailsWithoutEnsurer(t *testing.T) {
	_, err := ResolveAndEnsureReachable(context.Background(), PlacementRemote, PlacementLocal, EndpointPair{
		Name:     "jd-grpc",
		Internal: "jd:14231",
		External: "127.0.0.1:14231",
	}, nil)
	require.Error(t, err, "expected missing bridge ensurer to fail")
}

func TestResolveAndEnsureReachablePropagatesEnsurerError(t *testing.T) {
	_, err := ResolveAndEnsureReachable(context.Background(), PlacementRemote, PlacementLocal, EndpointPair{
		Name:     "jd-grpc",
		Internal: "jd:14231",
		External: "127.0.0.1:14231",
	}, func(_ context.Context, _ EndpointPair, _ int) error {
		return errors.New("boom")
	})
	require.Error(t, err, "expected ensurer error to be returned")
}
