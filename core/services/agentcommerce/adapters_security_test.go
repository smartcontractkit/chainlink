package agentcommerce

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type ccipClientFunc func(context.Context, RailSettlementRequest, ccipocr3.ChainSelector, ccipocr3.ChainSelector) (CCIPSubmission, error)

func (f ccipClientFunc) SubmitCCIP(ctx context.Context, req RailSettlementRequest, source, destination ccipocr3.ChainSelector) (CCIPSubmission, error) {
	return f(ctx, req, source, destination)
}

func (f ccipClientFunc) ObserveCCIP(context.Context, RailSettlementRequest, CCIPSubmission) (RailSettlementReceipt, error) {
	return RailSettlementReceipt{}, errors.New("observation not configured")
}

func TestCCIPAdapterPreservesTruthStatesAndIdentity(t *testing.T) {
	req := validRailRequest(t, "ccip")
	client := ccipClientFunc(func(_ context.Context, got RailSettlementRequest, source, destination ccipocr3.ChainSelector) (CCIPSubmission, error) {
		hash, err := SettlementRequestHash(got)
		messageID := sha256.Sum256([]byte("ccip-message-e2e"))
		sourceTx := sha256.Sum256([]byte("source-transaction-e2e"))
		return CCIPSubmission{
			RequestHash: hash, SourceChainSelector: source, DestinationChainSelector: destination,
			MessageID: ccipocr3.Bytes32(messageID), SourceTransactionHash: ccipocr3.Bytes32(sourceTx),
			Phase: SettlementBroadcast,
		}, err
	})
	adapter, err := NewCCIPSettlementAdapter(client)
	if err != nil {
		t.Fatal(err)
	}
	result, err := adapter.Submit(context.Background(), req, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if result.Phase != SettlementBroadcast {
		t.Fatalf("submission phase = %s", result.Phase)
	}
	if result.Phase == SettlementSettled {
		t.Fatal("CCIP source broadcast must not be treated as settled")
	}
}

func TestSettlementRejectsNonCanonicalAmountAndPrematureFinality(t *testing.T) {
	req := validRailRequest(t, "x402")
	req.Amount = "01"
	if _, err := SettlementRequestHash(req); err == nil {
		t.Fatal("non-canonical amount must fail")
	}
	req = validRailRequest(t, "x402")
	hash, _ := SettlementRequestHash(req)
	receipt, _ := SealSettlementReceipt(RailSettlementReceipt{
		Version: "1", RequestHash: hash, Rail: "x402", Phase: SettlementBroadcast,
		ExternalID: "payment", FinalityHash: hash, ObservedAt: time.Now().UTC(),
	})
	if err := VerifyFinalSettlement(context.Background(), req, receipt, settlementVerifierFunc(func(context.Context, RailSettlementRequest, RailSettlementReceipt) error {
		return nil
	})); err == nil {
		t.Fatal("BROADCAST must not verify as SETTLED")
	}
}

func TestFileCommerceStoreRejectsReplayConflictAndTampering(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "store.json")
	store, _ := NewFileCommerceStore(path)
	scope := SettlementScope{TenantID: "tenant-test", UserID: "user-test", ScopeID: "scope-test"}
	first := sha256.Sum256([]byte("first-evidence"))
	second := sha256.Sum256([]byte("second-evidence"))
	if err := store.ClaimReplay(ctx, scope, "replay-key", hex.EncodeToString(first[:])); err != nil {
		t.Fatal(err)
	}
	if err := store.ClaimReplay(ctx, scope, "replay-key", hex.EncodeToString(first[:])); err != nil {
		t.Fatalf("exact replay claim should be idempotent: %v", err)
	}
	if err := store.ClaimReplay(ctx, scope, "replay-key", hex.EncodeToString(second[:])); !errors.Is(err, ErrReplayClaimed) {
		t.Fatalf("conflicting replay error = %v", err)
	}
	if err := os.WriteFile(path, []byte(`{"records":{},"replay":{},"audit":[{"Sequence":1,"EntryHash":"tampered"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := store.VerifyIntegrity(ctx); err == nil {
		t.Fatal("tampered audit must fail integrity verification")
	}
}

func TestNativeOCRVerifierRejectsEnvelopeSubstitutionAndReplay(t *testing.T) {
	intent, proof, _, now := verifiedNativeE2EFixture(t)
	seed := sha256.Sum256([]byte("agentcommerce-e2e-ocr-seed"))
	key, err := crypto.ToECDSA(seed[:])
	if err != nil {
		t.Fatal(err)
	}
	policy := NativeOCRPolicy{Signers: []ocrtypes.OnchainPublicKey{crypto.PubkeyToAddress(key.PublicKey).Bytes()}, MinSignatures: 1, Issuer: "e2e-don", MaxLifetime: time.Minute}
	newVerifier := func() *NativeOCRProofVerifier {
		verifier, err := NewNativeOCRProofVerifier(NativeOCRProofVerifierConfig{Policy: policy, ReplayGuard: NewInMemoryReplayGuard(), Clock: ClockFunc(func() time.Time { return now }), Network: "chainlink-local-e2e", WorkflowID: "workflow-e2e", CapabilityID: "capability-e2e", MaxFutureSkew: time.Second, MaxAge: time.Minute})
		if err != nil {
			t.Fatal(err)
		}
		return verifier
	}
	verifier := newVerifier()
	result, err := verifier.Verify(context.Background(), intent, proof)
	if err != nil || !result.Verified {
		t.Fatalf("valid native proof rejected: %+v %v", result, err)
	}
	result, err = verifier.Verify(context.Background(), intent, proof)
	if err != nil || result.Verified {
		t.Fatalf("native replay accepted: %+v %v", result, err)
	}
	mutations := map[string]func(*SignedIntent, *Proof){
		"version":       func(_ *SignedIntent, p *Proof) { p.Metadata[metadataAttestationVersion] = "2" },
		"domain":        func(_ *SignedIntent, p *Proof) { p.Metadata[metadataAttestationDomain] = "wrong-domain" },
		"intent":        func(_ *SignedIntent, p *Proof) { p.Metadata[metadataIntentHash] = fmt.Sprintf("%064x", 1) },
		"execution":     func(_ *SignedIntent, p *Proof) { p.Metadata[metadataExecutionID] = "exec-substituted" },
		"output":        func(_ *SignedIntent, p *Proof) { p.Metadata[metadataOutputSHA256] = fmt.Sprintf("%064x", 2) },
		"report":        func(_ *SignedIntent, p *Proof) { p.Metadata[metadataReportSHA256] = fmt.Sprintf("%064x", 3) },
		"network":       func(_ *SignedIntent, p *Proof) { p.Metadata[metadataNetwork] = "wrong-network" },
		"workflow":      func(_ *SignedIntent, p *Proof) { p.Metadata[metadataWorkflowID] = "wrong-workflow" },
		"capability":    func(_ *SignedIntent, p *Proof) { p.Metadata[metadataCapabilityID] = "wrong-capability" },
		"issuer":        func(_ *SignedIntent, p *Proof) { p.Metadata[metadataIssuer] = "wrong-issuer" },
		"key_id":        func(_ *SignedIntent, p *Proof) { p.Metadata[metadataKeyID] = fmt.Sprintf("%064x", 4) },
		"evidence":      func(_ *SignedIntent, p *Proof) { p.Metadata[metadataAttestationSignature] = "00" },
		"signed_intent": func(i *SignedIntent, _ *Proof) { i.Hash = fmt.Sprintf("%064x", 5) },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			candidateIntent := intent
			candidateProof := proof
			candidateProof.Metadata = cloneStringMap(proof.Metadata)
			mutate(&candidateIntent, &candidateProof)
			result, err := newVerifier().Verify(context.Background(), candidateIntent, candidateProof)
			if err != nil || result.Verified {
				t.Fatalf("substitution accepted: %+v %v", result, err)
			}
		})
	}
}

func validRailRequest(t *testing.T, rail string) RailSettlementRequest {
	t.Helper()
	intent := sha256.Sum256([]byte("settlement-intent-" + rail))
	intentHash := hex.EncodeToString(intent[:])
	executionID, err := ExecutionIDForIntent(intentHash)
	if err != nil {
		t.Fatal(err)
	}
	authorization := sha256.Sum256([]byte("settlement-authorization-" + rail))
	return RailSettlementRequest{
		Version: "1", Scope: SettlementScope{TenantID: "tenant-test", UserID: "user-test", ScopeID: "scope-test"}, IntentHash: intentHash, ExecutionID: executionID,
		Rail: rail, Network: "local-e2e", Source: "buyer", Destination: "seller",
		Asset: "USDC", Amount: "10", AuthorizationHash: hex.EncodeToString(authorization[:]),
	}
}
