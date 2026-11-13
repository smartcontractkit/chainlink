package agentcommerce

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type e2eNativeExecutable struct {
	key    []byte
	digest ocrtypes.ConfigDigest
}

func verifiedNativeE2EFixture(t *testing.T) (SignedIntent, Proof, map[string]string, time.Time) {
	t.Helper()
	ctx := context.Background()
	now := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	seed := sha256.Sum256([]byte("agentcommerce-e2e-ocr-seed"))
	key, err := crypto.ToECDSA(seed[:])
	if err != nil {
		t.Fatal(err)
	}
	var digest ocrtypes.ConfigDigest
	digestSeed := sha256.Sum256([]byte("agentcommerce-e2e-config"))
	copy(digest[:], digestSeed[:])
	policy := NativeOCRPolicy{Signers: []ocrtypes.OnchainPublicKey{crypto.PubkeyToAddress(key.PublicKey).Bytes()}, MinSignatures: 1, Issuer: "e2e-don", MaxLifetime: time.Minute}
	terms := IntentTerms{ServiceRequest: ServiceRequest{Service: "agent-research", SettlementChain: "base-sepolia", Deliverables: []string{"report"}, EscrowTerms: EscrowTerms{ReleaseConditions: []string{"verified-attestation"}, Expiration: time.Hour}}, Price: Amount{Value: 10, Currency: "USDC", Rail: "x402"}, Buyer: "buyer-e2e", Seller: "seller-e2e", Timestamp: now}
	intentHash, err := IntentHash(terms)
	if err != nil {
		t.Fatal(err)
	}
	intent := SignedIntent{Terms: terms, Hash: intentHash}
	backend, err := NewNativeCapabilityBackend(NativeCapabilityBackendConfig{Executable: e2eNativeExecutable{key: crypto.FromECDSA(key), digest: digest}, Policy: policy, Network: "chainlink-local-e2e", WorkflowID: "workflow-e2e", WorkflowOwner: "owner-e2e", CapabilityID: "capability-e2e", Method: "execute"})
	if err != nil {
		t.Fatal(err)
	}
	executor, err := NewChainlinkExecutor(ChainlinkExecutorConfig{Backend: backend, Clock: ClockFunc(func() time.Time { return now }), Network: "chainlink-local-e2e", WorkflowID: "workflow-e2e", CapabilityID: "capability-e2e"})
	if err != nil {
		t.Fatal(err)
	}
	output, metadata, err := executor.Execute(ctx, intent)
	if err != nil {
		t.Fatal(err)
	}
	outputSum := sha256.Sum256(output)
	executionID, _ := ExecutionIDForIntent(intentHash)
	proof := Proof{OutputHash: hex.EncodeToString(outputSum[:]), ExecutionID: executionID, Timestamp: now, Metadata: metadata}
	verifier, err := NewNativeOCRProofVerifier(NativeOCRProofVerifierConfig{Policy: policy, ReplayGuard: NewInMemoryReplayGuard(), Clock: ClockFunc(func() time.Time { return now }), Network: "chainlink-local-e2e", WorkflowID: "workflow-e2e", CapabilityID: "capability-e2e", MaxFutureSkew: time.Second, MaxAge: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	result, err := verifier.Verify(ctx, intent, proof)
	if err != nil || !result.Verified {
		t.Fatalf("native verification: result=%+v err=%v", result, err)
	}
	return intent, proof, metadata, now
}

func (e e2eNativeExecutable) RegisterToWorkflow(context.Context, commoncap.RegisterToWorkflowRequest) error {
	return nil
}
func (e e2eNativeExecutable) UnregisterFromWorkflow(context.Context, commoncap.UnregisterFromWorkflowRequest) error {
	return nil
}
func (e e2eNativeExecutable) Execute(ctx context.Context, req commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	if err := ctx.Err(); err != nil {
		return commoncap.CapabilityResponse{}, err
	}
	output := []byte("deterministic-capability-output")
	metadata := commoncap.ResponseMetadata{
		CapDON_N: 1,
		Metering: []commoncap.MeteringNodeDetail{{SpendUnit: "COMPUTE", SpendValue: "1"}},
	}
	reportData, err := commoncap.ResponseToReportData(req.Metadata.WorkflowExecutionID, req.Metadata.ReferenceID, output, metadata)
	if err != nil {
		return commoncap.CapabilityResponse{}, err
	}
	key, err := crypto.ToECDSA(e.key)
	if err != nil {
		return commoncap.CapabilityResponse{}, err
	}
	signature, err := crypto.Sign(ocr2key.ReportToSigData3(e.digest, 1, reportData[:]), key)
	if err != nil {
		return commoncap.CapabilityResponse{}, err
	}
	return commoncap.CapabilityResponse{
		Payload:        &anypb.Any{TypeUrl: "type.googleapis.com/agentcommerce.v1.ExecutionResult", Value: output},
		Metadata:       metadata,
		OCRAttestation: &commoncap.OCRAttestation{ConfigDigest: e.digest, SequenceNumber: 1, Sigs: []commoncap.AttributedSignature{{Signer: 0, Signature: signature}}},
	}, nil
}

func TestAgentCommercePartsThreeThroughFiveE2E(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	seed := sha256.Sum256([]byte("agentcommerce-e2e-ocr-seed"))
	key, err := crypto.ToECDSA(seed[:])
	if err != nil {
		t.Fatal(err)
	}
	var digest ocrtypes.ConfigDigest
	digestSeed := sha256.Sum256([]byte("agentcommerce-e2e-config"))
	copy(digest[:], digestSeed[:])
	policy := NativeOCRPolicy{Signers: []ocrtypes.OnchainPublicKey{crypto.PubkeyToAddress(key.PublicKey).Bytes()}, MinSignatures: 1, Issuer: "e2e-don", MaxLifetime: time.Minute}

	terms := IntentTerms{ServiceRequest: ServiceRequest{Service: "agent-research", SettlementChain: "base-sepolia", Deliverables: []string{"report"}, EscrowTerms: EscrowTerms{ReleaseConditions: []string{"verified-attestation"}, Expiration: time.Hour}}, Price: Amount{Value: 10, Currency: "USDC", Rail: "x402"}, Buyer: "buyer-e2e", Seller: "seller-e2e", Timestamp: now}
	intentHash, err := IntentHash(terms)
	if err != nil {
		t.Fatal(err)
	}
	intent := SignedIntent{Terms: terms, Hash: intentHash}
	backend, err := NewNativeCapabilityBackend(NativeCapabilityBackendConfig{Executable: e2eNativeExecutable{key: crypto.FromECDSA(key), digest: digest}, Policy: policy, Network: "chainlink-local-e2e", WorkflowID: "workflow-e2e", WorkflowOwner: "owner-e2e", CapabilityID: "capability-e2e", Method: "execute"})
	if err != nil {
		t.Fatal(err)
	}
	executor, err := NewChainlinkExecutor(ChainlinkExecutorConfig{Backend: backend, Clock: ClockFunc(func() time.Time { return now }), Network: "chainlink-local-e2e", WorkflowID: "workflow-e2e", CapabilityID: "capability-e2e"})
	if err != nil {
		t.Fatal(err)
	}
	output, metadata, err := executor.Execute(ctx, intent)
	if err != nil {
		t.Fatal(err)
	}
	outputSum := sha256.Sum256(output)
	executionID, _ := ExecutionIDForIntent(intentHash)
	proof := Proof{OutputHash: hex.EncodeToString(outputSum[:]), ExecutionID: executionID, Timestamp: now, Metadata: metadata}
	proofVerifier, err := NewNativeOCRProofVerifier(NativeOCRProofVerifierConfig{Policy: policy, ReplayGuard: NewInMemoryReplayGuard(), Clock: ClockFunc(func() time.Time { return now }), Network: "chainlink-local-e2e", WorkflowID: "workflow-e2e", CapabilityID: "capability-e2e", MaxFutureSkew: time.Second, MaxAge: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	verification, err := proofVerifier.Verify(ctx, intent, proof)
	if err != nil || !verification.Verified {
		t.Fatalf("native verification: result=%+v err=%v", verification, err)
	}

	authorization := sha256.Sum256([]byte("deterministic-x402-authorization"))
	scope := SettlementScope{TenantID: "tenant-e2e", UserID: "user-e2e", ScopeID: "scope-e2e"}
	settlementRequest := RailSettlementRequest{Version: "1", Scope: scope, IntentHash: intentHash, ExecutionID: executionID, Rail: "x402", Network: "base-sepolia", Source: "buyer-e2e", Destination: "seller-e2e", Asset: "USDC", Amount: "10", AuthorizationHash: hex.EncodeToString(authorization[:])}
	facilitator := &deterministicX402Facilitator{paymentID: "payment-e2e-1", observations: []SettlementPhase{SettlementBroadcast, SettlementUnknown, SettlementConfirmed, SettlementSettled}}
	server := httptest.NewServer(facilitator)
	defer server.Close()
	httpClient, err := NewHTTPX402Client(server.Client(), server.URL, x402AuthorizerFunc(func(_ context.Context, _ X402Challenge, requestHash string) (string, error) {
		return "payment:" + requestHash, nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	x402, _ := NewX402SettlementAdapter(httpClient)
	submission, err := x402.Submit(ctx, settlementRequest)
	if err != nil {
		t.Fatal(err)
	}
	var receipt RailSettlementReceipt
	for _, expected := range []SettlementPhase{SettlementBroadcast, SettlementUnknown, SettlementConfirmed, SettlementSettled} {
		receipt, err = x402.Observe(ctx, settlementRequest, submission.PaymentID)
		if err != nil || receipt.Phase != expected {
			t.Fatalf("x402 observation want=%s got=%s err=%v", expected, receipt.Phase, err)
		}
	}
	verifier := settlementVerifierFunc(func(_ context.Context, request RailSettlementRequest, got RailSettlementReceipt) error {
		want := sha256.Sum256([]byte(facilitator.requestHash + ":" + facilitator.paymentID))
		if got.FinalityHash != hex.EncodeToString(want[:]) || request.Scope != facilitator.request.Scope {
			return errors.New("invalid x402 finality evidence")
		}
		return nil
	})
	if err := VerifyFinalSettlement(ctx, settlementRequest, receipt, verifier); err != nil {
		t.Fatal(err)
	}

	storePath := filepath.Join(t.TempDir(), "agentcommerce.json")
	store, _ := NewFileCommerceStore(storePath)
	record := CommerceRecord{Version: "1", Scope: scope, IntentHash: intentHash, ExecutionID: executionID, ExecutionPhase: ExecutionPrepared, OutputHash: proof.OutputHash, ReportHash: metadata[metadataReportSHA256], UpdatedAt: now}
	if _, created, err := store.Put(ctx, record); err != nil || !created {
		t.Fatalf("initial record: created=%v err=%v", created, err)
	}
	for _, phase := range []ExecutionPhase{ExecutionSubmitted, ExecutionExecuting, ExecutionExecuted, ExecutionVerified} {
		if phase == ExecutionVerified {
			evidence := sha256.Sum256([]byte(proof.OutputHash + metadata[metadataReportSHA256] + proof.ExecutionID))
			record.EvidenceHash = hex.EncodeToString(evidence[:])
		}
		record.ExecutionPhase, record.UpdatedAt = phase, record.UpdatedAt.Add(time.Second)
		if _, _, err := store.Put(ctx, record); err != nil {
			t.Fatal(err)
		}
	}
	for _, phase := range []SettlementPhase{SettlementPrepared, SettlementSubmitted, SettlementBroadcast, SettlementConfirmed, SettlementSettled} {
		record.SettlementPhase, record.UpdatedAt = phase, record.UpdatedAt.Add(time.Second)
		if phase == SettlementSettled {
			record.SettlementHash = receipt.ReceiptHash
		}
		if _, _, err := store.Put(ctx, record); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.ClaimReplay(ctx, scope, "settlement:"+submission.PaymentID, receipt.ReceiptHash); err != nil {
		t.Fatal(err)
	}
	reopened, _ := NewFileCommerceStore(storePath)
	if err := reopened.VerifyIntegrity(ctx); err != nil {
		t.Fatal(err)
	}
	loaded, err := reopened.Get(ctx, scope, intentHash, executionID)
	if err != nil || loaded.SettlementPhase != SettlementSettled || loaded.ExecutionPhase != ExecutionVerified {
		t.Fatalf("reopened record: %+v err=%v", loaded, err)
	}
	loaded.UpdatedAt = loaded.UpdatedAt.Add(time.Hour)
	if _, created, err := reopened.Put(ctx, loaded); err != nil || created {
		t.Fatalf("exact retry: created=%v err=%v", created, err)
	}
	conflict := sha256.Sum256([]byte("conflicting-evidence"))
	if err := reopened.ClaimReplay(ctx, scope, "settlement:"+submission.PaymentID, hex.EncodeToString(conflict[:])); !errors.Is(err, ErrReplayClaimed) {
		t.Fatalf("conflicting replay = %v", err)
	}
}
