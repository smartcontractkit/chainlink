package agentcommerce

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"
)

var ErrAttestationReplay = errors.New("attestation execution has already been used")

type ReplayGuard interface {
	Claim(context.Context, string) error
}

// InMemoryReplayGuard is for tests and single-process prototypes only.
// Production requires durable, shared, transactionally consistent storage.
type InMemoryReplayGuard struct {
	mu      sync.Mutex
	claimed map[string]struct{}
}

func NewInMemoryReplayGuard() *InMemoryReplayGuard {
	return &InMemoryReplayGuard{claimed: make(map[string]struct{})}
}

func (g *InMemoryReplayGuard) Claim(ctx context.Context, key string) error {
	if ctx == nil {
		return errors.New("context is required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if key == "" {
		return errors.New("replay key is required")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.claimed[key]; exists {
		return ErrAttestationReplay
	}
	g.claimed[key] = struct{}{}
	return nil
}

func AttestationReplayKey(a Attestation) (string, error) {
	input := struct {
		Domain       string `json:"domain"`
		Network      string `json:"network"`
		WorkflowID   string `json:"workflow_id"`
		CapabilityID string `json:"capability_id"`
		ExecutionID  string `json:"execution_id"`
	}{a.Domain, a.Network, a.WorkflowID, a.CapabilityID, a.ExecutionID}
	values := map[string]string{
		"domain": input.Domain, "network": input.Network,
		"workflow id": input.WorkflowID, "capability id": input.CapabilityID,
		"execution id": input.ExecutionID,
	}
	for name, value := range values {
		if err := requireCleanValue(name, value); err != nil {
			return "", err
		}
	}
	encoded, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}
