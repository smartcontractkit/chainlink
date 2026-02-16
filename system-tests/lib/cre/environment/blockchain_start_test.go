package environment

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

func TestValidatePhase2ARemoteBlockchainInput(t *testing.T) {
	if err := validatePhase2ARemoteBlockchainInput(nil); err == nil {
		t.Fatalf("expected nil input to fail validation")
	}

	if err := validatePhase2ARemoteBlockchainInput(&blockchain.Input{Type: blockchain.TypeGeth}); err == nil {
		t.Fatalf("expected non-anvil input to fail validation")
	}

	if err := validatePhase2ARemoteBlockchainInput(&blockchain.Input{Type: blockchain.TypeAnvil}); err != nil {
		t.Fatalf("expected anvil input to pass validation, got %v", err)
	}
}
