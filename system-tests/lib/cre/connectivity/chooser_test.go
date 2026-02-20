package connectivity

import (
	"context"
	"errors"
	"testing"
)

func TestResolveSamePlacementUsesInternal(t *testing.T) {
	r, err := Resolve(PlacementLocal, PlacementLocal, EndpointPair{
		Name:     "evm-rpc",
		Internal: "http://anvil:8545",
		External: "http://10.0.0.1:8545",
	})
	if err != nil {
		t.Fatalf("expected resolve to succeed: %v", err)
	}
	if r.URL != "http://anvil:8545" || r.SelectedKind != "internal" {
		t.Fatalf("unexpected resolution: %+v", r)
	}
	if r.RequiresBridge {
		t.Fatalf("did not expect bridge requirement for same placement")
	}
}

func TestResolveRemoteToLocalRequiresBridge(t *testing.T) {
	r, err := Resolve(PlacementRemote, PlacementLocal, EndpointPair{
		Name:     "jd-grpc",
		Internal: "jd:14231",
		External: "127.0.0.1:14231",
	})
	if err != nil {
		t.Fatalf("expected resolve to succeed: %v", err)
	}
	if !r.RequiresBridge || r.BridgePort != 14231 {
		t.Fatalf("expected bridge requirement with port 14231, got %+v", r)
	}
}

func TestResolveAndEnsureReachableCallsEnsurer(t *testing.T) {
	called := false
	r, err := ResolveAndEnsureReachable(context.Background(), PlacementRemote, PlacementLocal, EndpointPair{
		Name:     "jd-grpc",
		Internal: "jd:14231",
		External: "127.0.0.1:14231",
	}, func(_ context.Context, endpoint EndpointPair, port int) error {
		called = true
		if endpoint.Name != "jd-grpc" || port != 14231 {
			t.Fatalf("unexpected bridge args: endpoint=%s port=%d", endpoint.Name, port)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected resolve+ensure to succeed: %v", err)
	}
	if !called {
		t.Fatalf("expected bridge ensurer to be called")
	}
	if r.URL != "127.0.0.1:14231" {
		t.Fatalf("unexpected resolution URL: %s", r.URL)
	}
}

func TestResolveAndEnsureReachableFailsWithoutEnsurer(t *testing.T) {
	_, err := ResolveAndEnsureReachable(context.Background(), PlacementRemote, PlacementLocal, EndpointPair{
		Name:     "jd-grpc",
		Internal: "jd:14231",
		External: "127.0.0.1:14231",
	}, nil)
	if err == nil {
		t.Fatalf("expected missing bridge ensurer to fail")
	}
}

func TestResolveAndEnsureReachablePropagatesEnsurerError(t *testing.T) {
	_, err := ResolveAndEnsureReachable(context.Background(), PlacementRemote, PlacementLocal, EndpointPair{
		Name:     "jd-grpc",
		Internal: "jd:14231",
		External: "127.0.0.1:14231",
	}, func(_ context.Context, _ EndpointPair, _ int) error {
		return errors.New("boom")
	})
	if err == nil {
		t.Fatalf("expected ensurer error to be returned")
	}
}
