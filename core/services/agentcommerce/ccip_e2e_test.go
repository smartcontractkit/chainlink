package agentcommerce

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type deterministicCCIPHarness struct {
	mu          sync.Mutex
	submission  CCIPSubmission
	request     RailSettlementRequest
	phases      []SettlementPhase
	observation int
	delay       time.Duration
	fail        error
}

func (h *deterministicCCIPHarness) SubmitCCIP(ctx context.Context, req RailSettlementRequest, source, destination ccipocr3.ChainSelector) (CCIPSubmission, error) {
	if err := waitContext(ctx, h.delay); err != nil {
		return CCIPSubmission{}, err
	}
	requestHash, err := SettlementRequestHash(req)
	if err != nil {
		return CCIPSubmission{}, err
	}
	sourceTx := sha256.Sum256([]byte(fmt.Sprintf("source:%s:%d", requestHash, source)))
	message := sha256.Sum256([]byte(fmt.Sprintf("message:%s:%d:%d:%x", requestHash, source, destination, sourceTx)))
	submission := CCIPSubmission{RequestHash: requestHash, SourceChainSelector: source, DestinationChainSelector: destination, MessageID: ccipocr3.Bytes32(message), SourceTransactionHash: ccipocr3.Bytes32(sourceTx), Phase: SettlementBroadcast}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.submission.RequestHash == "" {
		h.submission, h.request = submission, req
	}
	if h.submission != submission || h.request != req {
		return CCIPSubmission{}, ErrCommerceConflict
	}
	return h.submission, nil
}

func (h *deterministicCCIPHarness) ObserveCCIP(ctx context.Context, req RailSettlementRequest, submission CCIPSubmission) (RailSettlementReceipt, error) {
	if err := waitContext(ctx, h.delay); err != nil {
		return RailSettlementReceipt{}, err
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.fail != nil {
		return RailSettlementReceipt{}, h.fail
	}
	if submission != h.submission || req != h.request {
		return RailSettlementReceipt{}, errors.New("CCIP observation identity mismatch")
	}
	phase := h.phases[h.observation]
	if h.observation < len(h.phases)-1 {
		h.observation++
	}
	finality := ""
	if phase == SettlementSettled {
		sum := sha256.Sum256([]byte(fmt.Sprintf("destination:%x:%d", submission.MessageID, submission.DestinationChainSelector)))
		finality = hex.EncodeToString(sum[:])
	}
	return RailSettlementReceipt{Version: "1", RequestHash: submission.RequestHash, Rail: "ccip", Phase: phase, ExternalID: hex.EncodeToString(submission.MessageID[:]), FinalityHash: finality, ObservedAt: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)}, nil
}

func waitContext(ctx context.Context, delay time.Duration) error {
	if delay == 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func TestCCIPEndToEndSourceToDestination(t *testing.T) {
	intent, proof, metadata, now := verifiedNativeE2EFixture(t)
	executionID, _ := ExecutionIDForIntent(intent.Hash)
	req := validRailRequest(t, "ccip")
	req.IntentHash = intent.Hash
	req.ExecutionID = executionID
	harness := &deterministicCCIPHarness{phases: []SettlementPhase{SettlementBroadcast, SettlementUnknown, SettlementConfirmed, SettlementSettled}}
	adapter, _ := NewCCIPSettlementAdapter(harness)
	submission, err := adapter.Submit(context.Background(), req, 100, 200)
	if err != nil {
		t.Fatal(err)
	}
	duplicate, err := adapter.Submit(context.Background(), req, 100, 200)
	if err != nil || duplicate != submission {
		t.Fatalf("duplicate CCIP submission: %+v %v", duplicate, err)
	}
	var receipt RailSettlementReceipt
	for _, expected := range []SettlementPhase{SettlementBroadcast, SettlementUnknown, SettlementConfirmed, SettlementSettled} {
		receipt, err = adapter.Observe(context.Background(), req, submission)
		if err != nil || receipt.Phase != expected {
			t.Fatalf("want=%s got=%s err=%v", expected, receipt.Phase, err)
		}
	}
	verifier := settlementVerifierFunc(func(_ context.Context, _ RailSettlementRequest, got RailSettlementReceipt) error {
		sum := sha256.Sum256([]byte(fmt.Sprintf("destination:%x:%d", submission.MessageID, submission.DestinationChainSelector)))
		if got.FinalityHash != hex.EncodeToString(sum[:]) {
			return errors.New("invalid destination finality")
		}
		return nil
	})
	if err := VerifyFinalSettlement(context.Background(), req, receipt, verifier); err != nil {
		t.Fatal(err)
	}
	storePath := filepath.Join(t.TempDir(), "ccip-commerce.json")
	store, _ := NewFileCommerceStore(storePath)
	record := CommerceRecord{Version: "1", Scope: req.Scope, IntentHash: intent.Hash, ExecutionID: executionID, ExecutionPhase: ExecutionPrepared, OutputHash: proof.OutputHash, ReportHash: metadata[metadataReportSHA256], UpdatedAt: now}
	if _, created, err := store.Put(context.Background(), record); err != nil || !created {
		t.Fatalf("create record: created=%v err=%v", created, err)
	}
	for _, phase := range []ExecutionPhase{ExecutionSubmitted, ExecutionExecuting, ExecutionExecuted, ExecutionVerified} {
		if phase == ExecutionVerified {
			evidence := sha256.Sum256([]byte(proof.OutputHash + metadata[metadataReportSHA256] + proof.ExecutionID))
			record.EvidenceHash = hex.EncodeToString(evidence[:])
		}
		record.ExecutionPhase, record.UpdatedAt = phase, record.UpdatedAt.Add(time.Second)
		if _, _, err := store.Put(context.Background(), record); err != nil {
			t.Fatal(err)
		}
	}
	for _, phase := range []SettlementPhase{SettlementPrepared, SettlementSubmitted, SettlementBroadcast, SettlementConfirmed, SettlementSettled} {
		record.SettlementPhase, record.UpdatedAt = phase, record.UpdatedAt.Add(time.Second)
		if phase == SettlementSettled {
			record.SettlementHash = receipt.ReceiptHash
		}
		if _, _, err := store.Put(context.Background(), record); err != nil {
			t.Fatal(err)
		}
	}
	messageReplayKey := "ccip:" + hex.EncodeToString(submission.MessageID[:])
	if err := store.ClaimReplay(context.Background(), req.Scope, messageReplayKey, receipt.ReceiptHash); err != nil {
		t.Fatal(err)
	}
	reopened, _ := NewFileCommerceStore(storePath)
	if err := reopened.VerifyIntegrity(context.Background()); err != nil {
		t.Fatal(err)
	}
	loaded, err := reopened.Get(context.Background(), req.Scope, intent.Hash, executionID)
	if err != nil || loaded.SettlementPhase != SettlementSettled {
		t.Fatalf("reloaded CCIP record: %+v err=%v", loaded, err)
	}
	loaded.UpdatedAt = loaded.UpdatedAt.Add(time.Hour)
	if _, created, err := reopened.Put(context.Background(), loaded); err != nil || created {
		t.Fatalf("CCIP retry: created=%v err=%v", created, err)
	}
	conflict := sha256.Sum256([]byte("conflicting-ccip-evidence"))
	if err := reopened.ClaimReplay(context.Background(), req.Scope, messageReplayKey, hex.EncodeToString(conflict[:])); !errors.Is(err, ErrReplayClaimed) {
		t.Fatalf("CCIP replay conflict=%v", err)
	}
}

func TestCCIPNegativeMatrix(t *testing.T) {
	req := validRailRequest(t, "ccip")
	harness := &deterministicCCIPHarness{phases: []SettlementPhase{SettlementConfirmed}}
	adapter, _ := NewCCIPSettlementAdapter(harness)
	if _, err := adapter.Submit(context.Background(), req, 0, 2); err == nil {
		t.Fatal("zero selector accepted")
	}
	if _, err := adapter.Submit(context.Background(), req, 2, 2); err == nil {
		t.Fatal("identical selectors accepted")
	}
	submission, err := adapter.Submit(context.Background(), req, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	wrongMessage := submission
	wrongMessage.MessageID[0] ^= 1
	if _, err := adapter.Observe(context.Background(), req, wrongMessage); err == nil {
		t.Fatal("wrong message id accepted")
	}
	wrongTx := submission
	wrongTx.SourceTransactionHash[0] ^= 1
	if _, err := adapter.Observe(context.Background(), req, wrongTx); err == nil {
		t.Fatal("wrong source tx accepted")
	}
	wrongDestination := submission
	wrongDestination.DestinationChainSelector = 3
	if _, err := adapter.Observe(context.Background(), req, wrongDestination); err == nil {
		t.Fatal("wrong destination accepted")
	}
	wrongReq := req
	wrongReq.ExecutionID = "exec-wrong"
	if _, err := adapter.Submit(context.Background(), wrongReq, 1, 2); err == nil {
		t.Fatal("wrong execution binding accepted")
	}
	receipt, err := adapter.Observe(context.Background(), req, submission)
	if err != nil {
		t.Fatal(err)
	}
	if receipt.Phase != SettlementConfirmed {
		t.Fatalf("source/destination pending phase=%s", receipt.Phase)
	}
	if err := VerifyFinalSettlement(context.Background(), req, receipt, settlementVerifierFunc(func(context.Context, RailSettlementRequest, RailSettlementReceipt) error { return nil })); err == nil {
		t.Fatal("destination pending verified as settled")
	}
}

func TestCCIPFailureTimeoutAndMalformedFinality(t *testing.T) {
	req := validRailRequest(t, "ccip")
	failing := &deterministicCCIPHarness{phases: []SettlementPhase{SettlementFailed}, fail: errors.New("destination failed")}
	adapter, _ := NewCCIPSettlementAdapter(failing)
	submission, err := adapter.Submit(context.Background(), req, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := adapter.Observe(context.Background(), req, submission); err == nil {
		t.Fatal("destination failure accepted")
	}
	timed := &deterministicCCIPHarness{phases: []SettlementPhase{SettlementBroadcast}, delay: time.Second}
	timedAdapter, _ := NewCCIPSettlementAdapter(timed)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if _, err := timedAdapter.Submit(ctx, req, 1, 2); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("timeout=%v", err)
	}
	malformed := RailSettlementReceipt{Version: "1", RequestHash: submission.RequestHash, Rail: "ccip", Phase: SettlementSettled, ExternalID: hex.EncodeToString(submission.MessageID[:]), FinalityHash: "bad", ObservedAt: time.Now().UTC()}
	malformed, _ = SealSettlementReceipt(malformed)
	if err := VerifyFinalSettlement(context.Background(), req, malformed, settlementVerifierFunc(func(context.Context, RailSettlementRequest, RailSettlementReceipt) error {
		return errors.New("malformed finality")
	})); err == nil {
		t.Fatal("malformed finality accepted")
	}
}
