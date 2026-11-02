package agentcommerce

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRunTransactionReleasesEscrowAndAuditsLifecycle(t *testing.T) {
	ctx := context.Background()
	keys := &sync.Map{}
	buyerWallet, err := NewEd25519Wallet(keys)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sellerWallet, err := NewEd25519Wallet(keys)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reputation := NewInMemoryReputation()
	if err := reputation.Update(ctx, ReputationEvent{AgentID: sellerWallet.Address(), Delta: 10, Reason: "seed"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sellerScore, err := reputation.Query(ctx, sellerWallet.Address())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	directory := NewInMemoryDirectory()
	directory.Register(AgentProfile{
		ID:           sellerWallet.Address(),
		Name:         "chainlink-research-agent",
		Capabilities: []string{"agent-research"},
		Pricing: map[string]Amount{
			"agent-research": {Value: 25, Currency: "LINK", Rail: "evm-escrow"},
		},
		Reputation: sellerScore,
		Metadata:   map[string]string{"oracle": "chainlink-functions"},
	})

	audit := &AuditLog{}

	cfg := OrchestratorConfig{
		Directory:  directory,
		Escrow:     NewInMemoryEscrow(),
		Reputation: reputation,
		Audit:      audit,
		Wallet:     buyerWallet,
		Policy: Policy{
			MaxSpend:            100,
			AllowedCurrencies:   []string{"LINK"},
			AllowedRails:        []string{"evm-escrow"},
			MinSellerReputation: 5,
		},
		Timeout:    10 * time.Second,
		MaxRetries: 2,
		ExecuteAs: func(_ context.Context, _ SignedIntent) ([]byte, map[string]string, error) {
			return []byte("verified research deliverable"), map[string]string{"content-type": "text/plain"}, nil
		},
	}

	orchestrator, err := NewOrchestrator(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	receipt, err := orchestrator.RunTransaction(ctx, ServiceQuery{
		Capability:     "agent-research",
		MaxPrice:       Amount{Value: 50, Currency: "LINK", Rail: "evm-escrow"},
		MinReputation:  5,
		SettlementRail: "evm-escrow",
	}, ServiceRequest{
		Service:         "agent-research",
		SLA:             "p95<30s",
		Deliverables:    []string{"summary", "citations", "output_hash"},
		SettlementChain: "ethereum-sepolia",
		EscrowTerms:     EscrowTerms{ReleaseConditions: []string{"deterministic-proof"}, Expiration: time.Hour},
	}, sellerWallet)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receipt.Status != "released" {
		t.Fatalf("status = %q", receipt.Status)
	}
	if receipt.To != sellerWallet.Address() {
		t.Fatalf("to = %q", receipt.To)
	}
	if receipt.Amount.Value != 25 {
		t.Fatalf("amount = %d", receipt.Amount.Value)
	}

	updatedSeller, err := reputation.Query(ctx, sellerWallet.Address())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updatedSeller.Score != 11 {
		t.Fatalf("seller score = %d", updatedSeller.Score)
	}
	entries := audit.Entries()
	if len(entries) != 5 {
		t.Fatalf("entries len = %d", len(entries))
	}
	if got := []string{entries[0].Kind, entries[1].Kind, entries[2].Kind, entries[3].Kind, entries[4].Kind}; got[0] != "intent" || got[1] != "escrow" || got[2] != "proof" || got[3] != "verification" || got[4] != "settlement" {
		t.Fatalf("audit kinds = %#v", got)
	}
}

func TestRunTransactionRefundsWhenVerificationFails(t *testing.T) {
	ctx := context.Background()
	keys := &sync.Map{}
	buyerWallet, err := NewEd25519Wallet(keys)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sellerWallet, err := NewEd25519Wallet(keys)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reputation := NewInMemoryReputation()
	directory := NewInMemoryDirectory()
	directory.Register(AgentProfile{ID: sellerWallet.Address(), Capabilities: []string{"build"}, Pricing: map[string]Amount{"build": {Value: 3, Currency: "USDC", Rail: "x402"}}, Reputation: ReputationScore{Score: 1}})

	cfg := OrchestratorConfig{
		Directory:  directory,
		Escrow:     NewInMemoryEscrow(),
		Reputation: reputation,
		Audit:      &AuditLog{},
		Wallet:     buyerWallet,
		Policy:     Policy{MaxSpend: 10, AllowedCurrencies: []string{"USDC"}, AllowedRails: []string{"x402"}},
		Timeout:    10 * time.Second,
		MaxRetries: 2,
		ExecuteAs: func(_ context.Context, _ SignedIntent) ([]byte, map[string]string, error) {
			return []byte("bad output"), nil, nil
		},
		VerifyWith: func(_ context.Context, _ SignedIntent, _ Proof) (VerificationResult, error) {
			return VerificationResult{Verified: false, Method: "oracle", Reason: "oracle rejected deliverable", Time: time.Now().UTC()}, nil
		},
	}

	orchestrator, err := NewOrchestrator(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	receipt, err := orchestrator.RunTransaction(ctx, ServiceQuery{Capability: "build", MaxPrice: Amount{Value: 10}, SettlementRail: "x402"}, ServiceRequest{Service: "build", EscrowTerms: EscrowTerms{ReleaseConditions: []string{"oracle"}, Expiration: time.Hour}}, sellerWallet)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receipt.Status != "refunded" {
		t.Fatalf("status = %q", receipt.Status)
	}
	if receipt.To != buyerWallet.Address() {
		t.Fatalf("to = %q", receipt.To)
	}
	sellerScore, err := reputation.Query(ctx, sellerWallet.Address())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sellerScore.Score != -1 {
		t.Fatalf("seller score = %d", sellerScore.Score)
	}
}

func TestPolicyRejectsOverspendBeforeEscrow(t *testing.T) {
	err := EvaluatePolicy(Policy{MaxSpend: 5}, AgentProfile{}, IntentTerms{
		ServiceRequest: ServiceRequest{
			Service:     "agent-research",
			EscrowTerms: EscrowTerms{ReleaseConditions: []string{"proof"}, Expiration: time.Hour},
		},
		Price:     Amount{Value: 6, Currency: "LINK", Rail: "evm-escrow"},
		Buyer:     "buyer",
		Seller:    "seller",
		Timestamp: time.Now().UTC(),
	})
	if err == nil || !strings.Contains(err.Error(), "exceeds max spend") {
		t.Fatalf("expected overspend error, got %v", err)
	}
}

func TestDirectoryFiltersRequiredMetadataAndAuditHashChain(t *testing.T) {
	directory := NewInMemoryDirectory()
	directory.Register(AgentProfile{
		ID:           "seller-a",
		Capabilities: []string{"proof"},
		Pricing:      map[string]Amount{"proof": {Value: 1, Currency: "LINK", Rail: "evm-escrow"}},
		Reputation:   ReputationScore{Score: 10},
		Metadata:     map[string]string{"oracle": "functions"},
	})
	directory.Register(AgentProfile{
		ID:           "seller-b",
		Capabilities: []string{"proof"},
		Pricing:      map[string]Amount{"proof": {Value: 1, Currency: "LINK", Rail: "evm-escrow"}},
		Reputation:   ReputationScore{Score: 20},
		Metadata:     map[string]string{"oracle": "manual"},
	})

	matches := directory.Discover(ServiceQuery{
		Capability:       "proof",
		SettlementRail:   "evm-escrow",
		RequiredMetadata: map[string]string{"oracle": "functions"},
	})
	if len(matches) != 1 || matches[0].ID != "seller-a" {
		t.Fatalf("matches = %#v", matches)
	}

	audit := &AuditLog{}
	if err := audit.Store("intent", "intent-1", map[string]string{"hash": "abc"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := audit.Store("settlement", "settle-1", map[string]string{"status": "released"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := audit.Entries()
	if len(entries) != 2 {
		t.Fatalf("entries len = %d", len(entries))
	}
	if entries[0].EntryHash == "" || entries[0].PayloadHash == "" {
		t.Fatalf("missing first audit hashes: %#v", entries[0])
	}
	if entries[1].PreviousHash != entries[0].EntryHash {
		t.Fatalf("previous hash = %q, want %q", entries[1].PreviousHash, entries[0].EntryHash)
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("nil wallet in SignIntent returns error", func(t *testing.T) {
		_, err := SignIntent(IntentTerms{}, nil, nil)
		if err == nil {
			t.Fatal("expected error for nil wallets")
		}
	})

	t.Run("empty hash in InMemoryEscrow.Lock returns error", func(t *testing.T) {
		e := NewInMemoryEscrow()
		intent := SignedIntent{Hash: "", Terms: IntentTerms{
			ServiceRequest: ServiceRequest{
				Service:     "test",
				EscrowTerms: EscrowTerms{ReleaseConditions: []string{"proof"}, Expiration: time.Hour},
			},
			Price:     Amount{Value: 1, Currency: "LINK", Rail: "evm-escrow"},
			Buyer:     "buyer",
			Seller:    "seller",
			Timestamp: time.Now().UTC(),
		}}
		_, err := e.Lock(context.Background(), intent)
		if err == nil {
			t.Fatal("expected error for empty hash")
		}
	})

	t.Run("orchestrator handles nil execute as", func(t *testing.T) {
		cfg := OrchestratorConfig{
			Directory:  NewInMemoryDirectory(),
			Escrow:     NewInMemoryEscrow(),
			Reputation: NewInMemoryReputation(),
			Audit:      &AuditLog{},
			Wallet:     &Ed25519Wallet{address: "buyer"},
			ExecuteAs:  nil,
		}
		_, err := NewOrchestrator(cfg)
		if err == nil {
			t.Fatal("expected error for nil ExecuteAs")
		}
	})

	t.Run("SafeWallet handles nil inner wallet", func(t *testing.T) {
		sw := NewSafeWallet(nil)
		if sw.Address() != "" {
			t.Error("expected empty address for nil wallet")
		}
		if _, err := sw.Sign([]byte("test")); err == nil {
			t.Error("expected error for nil wallet Sign")
		}
		if sw.Verify("test", []byte("test"), []byte("sig")) {
			t.Error("expected false for nil wallet Verify")
		}
	})
}
