package tunnel

import (
	"context"
	"testing"
)

type fakeProvider struct {
	openCount  int
	closeCount int
}

func (f *fakeProvider) Open(_ context.Context, ref EndpointRef) (TunnelBinding, error) {
	f.openCount++
	return TunnelBinding{
		EndpointRef: ref,
		LocalPort:   10000 + f.openCount,
		LocalURL:    "http://127.0.0.1:10000",
	}, nil
}

func (f *fakeProvider) Close(_ context.Context, _ TunnelBinding) error {
	f.closeCount++
	return nil
}

func (f *fakeProvider) Name() string { return "fake" }

func TestManagerStartDedupsAndStops(t *testing.T) {
	provider := &fakeProvider{}
	mgr := NewManager(provider)

	refs := []EndpointRef{
		{
			ComponentID:  "blockchain:0:anvil",
			EndpointName: "node-0-http",
			Scheme:       "http",
			Host:         "127.0.0.1",
			Port:         8545,
		},
		{
			ComponentID:  "blockchain:0:anvil",
			EndpointName: "node-0-ws",
			Scheme:       "ws",
			Host:         "127.0.0.1",
			Port:         8546,
		},
	}

	started, err := mgr.Start(context.Background(), refs)
	if err != nil {
		t.Fatalf("expected start to succeed: %v", err)
	}
	if len(started) != 2 {
		t.Fatalf("expected 2 bindings, got %d", len(started))
	}
	if provider.openCount != 2 {
		t.Fatalf("expected 2 opens, got %d", provider.openCount)
	}

	startedAgain, err := mgr.Start(context.Background(), refs)
	if err != nil {
		t.Fatalf("expected dedup start to succeed: %v", err)
	}
	if len(startedAgain) != 2 {
		t.Fatalf("expected 2 dedup bindings, got %d", len(startedAgain))
	}
	if provider.openCount != 2 {
		t.Fatalf("expected no extra open calls after dedup, got %d", provider.openCount)
	}

	if !mgr.IsStarted() {
		t.Fatalf("expected manager to report started")
	}

	if err := mgr.Stop(context.Background()); err != nil {
		t.Fatalf("expected idempotent stop to succeed: %v", err)
	}
	if provider.closeCount != 2 {
		t.Fatalf("expected 2 closes from stop, got %d", provider.closeCount)
	}
	if mgr.IsStarted() {
		t.Fatalf("expected manager to report no active tunnels after stop")
	}
}
