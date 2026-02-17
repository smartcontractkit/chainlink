package environment

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/tunnel"
)

type countingTunnelManager struct {
	stopCalls int
}

func (c *countingTunnelManager) Start(_ context.Context, _ []tunnel.EndpointRef) ([]tunnel.TunnelBinding, error) {
	return nil, nil
}
func (c *countingTunnelManager) Stop(_ context.Context) error {
	c.stopCalls++
	return nil
}
func (c *countingTunnelManager) IsStarted() bool { return false }

func TestSetupOutputCloseIsIdempotent(t *testing.T) {
	manager := &countingTunnelManager{}
	out := &SetupOutput{tunnelManager: manager}

	if err := out.Close(context.Background()); err != nil {
		t.Fatalf("expected first close to succeed: %v", err)
	}
	if err := out.Close(context.Background()); err != nil {
		t.Fatalf("expected second close to succeed: %v", err)
	}
	if manager.stopCalls != 1 {
		t.Fatalf("expected tunnel manager stop once, got %d", manager.stopCalls)
	}
}
