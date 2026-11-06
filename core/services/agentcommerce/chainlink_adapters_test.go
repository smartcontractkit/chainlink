package agentcommerce

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"
)

const (
	testNetwork      = "ethereum-sepolia"
	testWorkflowID   = "workflow-agentcommerce-v1"
	testCapabilityID = "capability-agent-execution-v1"
	testIssuer       = "chainlink-test-issuer"
	testKeyID        = "test-key-v1"
)

var testNow = time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)

type testSigningBackend struct {
	privateKey ed25519.PrivateKey
	output     []byte
	report     []byte
	issuedAt   time.Time
	expiresAt  time.Time
	mutate     func(*Attestation)
	err        error
}

func (b *testSigningBackend) Execute(_ context.Context, request ChainlinkExecutionRequest) (ChainlinkExecutionReport, error) {
	if b.err != nil {
		return ChainlinkExecutionReport{}, b.err
	}
	outputSum, reportSum := sha256.Sum256(b.output), sha256.Sum256(b.report)
	issuedAt := b.issuedAt
	if issuedAt.IsZero() {
		issuedAt = testNow
	}
	expiresAt := b.expiresAt
	if expiresAt.IsZero() {
		expiresAt = issuedAt.Add(5 * time.Minute)
	}
	a := Attestation{
		Version: ChainlinkAttestationVersion, Domain: ChainlinkAttestationDomain,
		IntentHash: request.IntentHash, ExecutionID: request.ExecutionID,
		OutputHash: hex.EncodeToString(outputSum[:]), ReportHash: hex.EncodeToString(reportSum[:]),
		Network: request.Network, WorkflowID: request.WorkflowID, CapabilityID: request.CapabilityID,
		Issuer: testIssuer, KeyID: testKeyID, IssuedAt: issuedAt, ExpiresAt: expiresAt,
	}
	if b.mutate != nil {
		b.mutate(&a)
	}
	payload, err := CanonicalAttestationPayload(a)
	if err != nil {
		return ChainlinkExecutionReport{}, err
	}
	a.Signature = ed25519.Sign(b.privateKey, payload)
	return ChainlinkExecutionReport{
		Output: append([]byte(nil), b.output...), Report: append([]byte(nil), b.report...),
		Attestation: a, Metadata: map[string]string{"backend": "deterministic-test-backend"},
	}, nil
}

func TestChainlinkAdaptersAcceptValidAttestation(t *testing.T) {
	f := newAdapterFixture(t, nil)
	proof := f.executeProof(t)
	result, err := f.verifier.Verify(context.Background(), f.intent, proof)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !result.Verified || result.Method != "chainlink-attestation-v1" {
		t.Fatalf("unexpected result: %+v", result)
	}
	if !result.Time.Equal(testNow) {
		t.Fatalf("verification time = %s, want %s", result.Time, testNow)
	}
	var executeHook ServiceExecutor = f.executor.Execute
	var verifyHook ProofVerifier = f.verifier.Verify
	if executeHook == nil || verifyHook == nil {
		t.Fatal("adapter hooks must match orchestrator hooks")
	}
}

func TestOrchestratorPropagatesInjectedClockThroughChainlinkPath(t *testing.T) {
	f := newAdapterFixture(t, nil)
	ctx := context.Background()
	keys := &sync.Map{}
	buyerWallet, err := NewEd25519Wallet(keys)
	if err != nil {
		t.Fatalf("NewEd25519Wallet buyer: %v", err)
	}
	sellerWallet, err := NewEd25519Wallet(keys)
	if err != nil {
		t.Fatalf("NewEd25519Wallet seller: %v", err)
	}

	reputation := NewInMemoryReputation()
	if err := reputation.Update(ctx, ReputationEvent{
		AgentID: sellerWallet.Address(),
		Delta:   10,
		Reason:  "deterministic-test-seed",
	}); err != nil {
		t.Fatalf("seed reputation: %v", err)
	}
	sellerScore, err := reputation.Query(ctx, sellerWallet.Address())
	if err != nil {
		t.Fatalf("query reputation: %v", err)
	}

	directory := NewInMemoryDirectory()
	directory.Register(AgentProfile{
		ID:           sellerWallet.Address(),
		Name:         "deterministic-chainlink-test-agent",
		Capabilities: []string{"agent-research"},
		Pricing: map[string]Amount{
			"agent-research": {Value: 10, Currency: "USDC", Rail: "x402"},
		},
		Reputation: sellerScore,
	})

	audit := &AuditLog{}
	orchestrator, err := NewOrchestrator(OrchestratorConfig{
		Directory:  directory,
		Escrow:     NewInMemoryEscrow(),
		Reputation: reputation,
		Audit:      audit,
		Wallet:     buyerWallet,
		Clock:      ClockFunc(func() time.Time { return testNow }),
		VerifyWith: f.verifier.Verify,
		ExecuteAs:  f.executor.Execute,
		Timeout:    10 * time.Second,
		MaxRetries: 1,
		Policy: Policy{
			MaxSpend:          20,
			AllowedCurrencies: []string{"USDC"},
			AllowedRails:      []string{"x402"},
		},
	})
	if err != nil {
		t.Fatalf("NewOrchestrator: %v", err)
	}

	receipt, err := orchestrator.RunTransaction(
		ctx,
		ServiceQuery{
			Capability:     "agent-research",
			MaxPrice:       Amount{Value: 20, Currency: "USDC", Rail: "x402"},
			SettlementRail: "x402",
		},
		ServiceRequest{
			Service:         "agent-research",
			Deliverables:    []string{"report"},
			SettlementChain: testNetwork,
			EscrowTerms: EscrowTerms{
				ReleaseConditions: []string{"chainlink-attestation"},
				Expiration:        time.Hour,
			},
		},
		sellerWallet,
	)
	if err != nil {
		t.Fatalf("RunTransaction: %v", err)
	}
	if receipt.Status != "released" {
		t.Fatalf("settlement status = %q", receipt.Status)
	}

	var intent SignedIntent
	var proof Proof
	var result VerificationResult
	for _, entry := range audit.Entries() {
		switch entry.Kind {
		case "intent":
			if err := json.Unmarshal(entry.Payload, &intent); err != nil {
				t.Fatalf("decode intent audit entry: %v", err)
			}
		case "proof":
			if err := json.Unmarshal(entry.Payload, &proof); err != nil {
				t.Fatalf("decode proof audit entry: %v", err)
			}
		case "verification":
			if err := json.Unmarshal(entry.Payload, &result); err != nil {
				t.Fatalf("decode verification audit entry: %v", err)
			}
		}
	}

	if !intent.Terms.Timestamp.Equal(testNow) {
		t.Fatalf("intent time = %s, want %s", intent.Terms.Timestamp, testNow)
	}
	if !proof.Timestamp.Equal(testNow) {
		t.Fatalf("proof time = %s, want %s", proof.Timestamp, testNow)
	}
	if !result.Time.Equal(testNow) {
		t.Fatalf("verification time = %s, want %s", result.Time, testNow)
	}
	executionID, err := ExecutionIDForIntent(intent.Hash)
	if err != nil {
		t.Fatalf("ExecutionIDForIntent: %v", err)
	}
	if proof.ExecutionID != executionID {
		t.Fatalf("execution id = %q, want %q", proof.ExecutionID, executionID)
	}
}

func TestChainlinkExecutorRejectsEmptyReport(t *testing.T) {
	f := newAdapterFixture(t, func(b *testSigningBackend) { b.report = nil })
	if _, _, err := f.executor.Execute(context.Background(), f.intent); err == nil {
		t.Fatal("expected empty report rejection")
	}
}

func TestChainlinkExecutorRejectsExecutionIDSubstitution(t *testing.T) {
	f := newAdapterFixture(t, func(b *testSigningBackend) {
		b.mutate = func(a *Attestation) { a.ExecutionID = "exec-substituted" }
	})
	if _, _, err := f.executor.Execute(context.Background(), f.intent); err == nil {
		t.Fatal("expected execution id mismatch")
	}
}

func TestChainlinkVerifierRejectsReplay(t *testing.T) {
	f := newAdapterFixture(t, nil)
	proof := f.executeProof(t)
	first, err := f.verifier.Verify(context.Background(), f.intent, proof)
	if err != nil || !first.Verified {
		t.Fatalf("first verification: result=%+v err=%v", first, err)
	}
	second, err := f.verifier.Verify(context.Background(), f.intent, proof)
	if err != nil {
		t.Fatalf("second verification: %v", err)
	}
	if second.Verified || second.Reason != "attestation replay rejected" {
		t.Fatalf("replayed result: %+v", second)
	}
}

func TestChainlinkVerifierRejectsSubstitutedEvidence(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Proof)
	}{
		{"malformed output hash", func(p *Proof) { p.OutputHash = "abcd" }},
		{"different output hash", func(p *Proof) { p.OutputHash = hashText("other-output") }},
		{"different execution id", func(p *Proof) { p.ExecutionID = "exec-substituted" }},
		{"different network", func(p *Proof) { p.Metadata[metadataNetwork] = "ethereum-mainnet" }},
		{"different workflow", func(p *Proof) { p.Metadata[metadataWorkflowID] = "other-workflow" }},
		{"different issuer", func(p *Proof) { p.Metadata[metadataIssuer] = "untrusted-issuer" }},
		{"unknown key", func(p *Proof) { p.Metadata[metadataKeyID] = "unknown-key" }},
		{"different intent", func(p *Proof) { p.Metadata[metadataIntentHash] = hashText("other-intent") }},
		{"different report", func(p *Proof) { p.Metadata[metadataReportSHA256] = hashText("other-report") }},
		{"invalid signature", func(p *Proof) {
			p.Metadata[metadataAttestationSignature] = hex.EncodeToString(make([]byte, ed25519.SignatureSize))
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := newAdapterFixture(t, nil)
			proof := f.executeProof(t)
			proof.Metadata = cloneStringMap(proof.Metadata)
			tc.mutate(&proof)
			result, err := f.verifier.Verify(context.Background(), f.intent, proof)
			if err != nil {
				t.Fatalf("Verify: %v", err)
			}
			if result.Verified {
				t.Fatalf("substituted evidence verified: %+v", result)
			}
		})
	}
}

func TestChainlinkVerifierRejectsStaleExpiredAndFutureEvidence(t *testing.T) {
	tests := []struct {
		name      string
		issuedAt  time.Time
		expiresAt time.Time
	}{
		{"stale", testNow.Add(-20 * time.Minute), testNow.Add(5 * time.Minute)},
		{"expired", testNow.Add(-5 * time.Minute), testNow.Add(-time.Second)},
		{"future", testNow.Add(2 * time.Minute), testNow.Add(7 * time.Minute)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := newAdapterFixture(t, func(b *testSigningBackend) {
				b.issuedAt, b.expiresAt = tc.issuedAt, tc.expiresAt
			})
			result, err := f.verifier.Verify(context.Background(), f.intent, f.executeProof(t))
			if err != nil {
				t.Fatalf("Verify: %v", err)
			}
			if result.Verified {
				t.Fatalf("%s evidence verified", tc.name)
			}
		})
	}
}

func TestChainlinkVerifierRejectsWrongSigningKey(t *testing.T) {
	f := newAdapterFixture(t, nil)
	proof := f.executeProof(t)
	a, err := attestationFromProof(proof)
	if err != nil {
		t.Fatalf("attestationFromProof: %v", err)
	}
	payload, err := CanonicalAttestationPayload(a)
	if err != nil {
		t.Fatalf("CanonicalAttestationPayload: %v", err)
	}
	seed := sha256.Sum256([]byte("different deterministic test key"))
	wrongKey := ed25519.NewKeyFromSeed(seed[:])
	proof.Metadata = cloneStringMap(proof.Metadata)
	proof.Metadata[metadataAttestationSignature] = hex.EncodeToString(ed25519.Sign(wrongKey, payload))
	result, err := f.verifier.Verify(context.Background(), f.intent, proof)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if result.Verified {
		t.Fatal("wrong signing key verified")
	}
}

func TestCanonicalAttestationPayloadIsDeterministic(t *testing.T) {
	f := newAdapterFixture(t, nil)
	a, err := attestationFromProof(f.executeProof(t))
	if err != nil {
		t.Fatalf("attestationFromProof: %v", err)
	}
	first, err := CanonicalAttestationPayload(a)
	if err != nil {
		t.Fatalf("first payload: %v", err)
	}
	second, err := CanonicalAttestationPayload(a)
	if err != nil {
		t.Fatalf("second payload: %v", err)
	}
	if string(first) != string(second) {
		t.Fatalf("canonical payload changed:\n%s\n%s", first, second)
	}
	firstDigest, _ := AttestationDigest(a)
	secondDigest, _ := AttestationDigest(a)
	if firstDigest != secondDigest {
		t.Fatalf("digest changed: %s != %s", firstDigest, secondDigest)
	}
}

func TestChainlinkExecutorPropagatesBackendFailure(t *testing.T) {
	f := newAdapterFixture(t, func(b *testSigningBackend) { b.err = errors.New("backend unavailable") })
	if _, _, err := f.executor.Execute(context.Background(), f.intent); err == nil {
		t.Fatal("expected backend failure")
	}
}

type adapterFixture struct {
	executor *ChainlinkExecutor
	verifier *ChainlinkProofVerifier
	intent   SignedIntent
}

func newAdapterFixture(t *testing.T, configure func(*testSigningBackend)) adapterFixture {
	t.Helper()
	seed := sha256.Sum256([]byte("agentcommerce-part2-deterministic-test-key"))
	privateKey := ed25519.NewKeyFromSeed(seed[:])
	publicKey := privateKey.Public().(ed25519.PublicKey)
	backend := &testSigningBackend{privateKey: privateKey, output: []byte("deterministic verified result"), report: []byte("deterministic chainlink report")}
	if configure != nil {
		configure(backend)
	}
	clock := ClockFunc(func() time.Time { return testNow })
	executor, err := NewChainlinkExecutor(ChainlinkExecutorConfig{Backend: backend, Clock: clock, Network: testNetwork, WorkflowID: testWorkflowID, CapabilityID: testCapabilityID})
	if err != nil {
		t.Fatalf("NewChainlinkExecutor: %v", err)
	}
	signatureVerifier, err := NewEd25519AttestationVerifier(map[AttestationKey]ed25519.PublicKey{{Issuer: testIssuer, KeyID: testKeyID}: publicKey})
	if err != nil {
		t.Fatalf("NewEd25519AttestationVerifier: %v", err)
	}
	verifier, err := NewChainlinkProofVerifier(ChainlinkProofVerifierConfig{
		SignatureVerifier: signatureVerifier, ReplayGuard: NewInMemoryReplayGuard(), Clock: clock,
		Policy: VerificationPolicy{
			Version: ChainlinkAttestationVersion, Domain: ChainlinkAttestationDomain,
			Network: testNetwork, WorkflowID: testWorkflowID, CapabilityID: testCapabilityID,
			TrustedIdentities: []TrustedAttestationIdentity{{Issuer: testIssuer, KeyID: testKeyID}},
			MaxFutureSkew:     30 * time.Second, MaxAge: 10 * time.Minute, MaxLifetime: 10 * time.Minute,
		},
	})
	if err != nil {
		t.Fatalf("NewChainlinkProofVerifier: %v", err)
	}
	return adapterFixture{executor, verifier, validChainlinkTestIntent(t)}
}

func (f adapterFixture) executeProof(t *testing.T) Proof {
	t.Helper()
	output, metadata, err := f.executor.Execute(context.Background(), f.intent)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	executionID, err := ExecutionIDForIntent(f.intent.Hash)
	if err != nil {
		t.Fatalf("ExecutionIDForIntent: %v", err)
	}
	outputSum := sha256.Sum256(output)
	return Proof{OutputHash: hex.EncodeToString(outputSum[:]), Timestamp: testNow, ExecutionID: executionID, Metadata: cloneStringMap(metadata)}
}

func validChainlinkTestIntent(t *testing.T) SignedIntent {
	t.Helper()
	terms := IntentTerms{
		ServiceRequest: ServiceRequest{
			Service: "agent-research", SLA: "p95<30s", Deliverables: []string{"report"},
			SettlementChain: testNetwork,
			EscrowTerms:     EscrowTerms{ReleaseConditions: []string{"chainlink-attestation"}, Expiration: time.Hour},
			Metadata:        map[string]string{"content-type": "text/plain"},
		},
		Price: Amount{Value: 10, Currency: "USDC", Rail: "x402"},
		Buyer: "buyer", Seller: "seller", Timestamp: testNow.Add(-time.Minute),
	}
	hash, err := IntentHash(terms)
	if err != nil {
		t.Fatalf("IntentHash: %v", err)
	}
	return SignedIntent{Terms: terms, Hash: hash}
}

func hashText(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
