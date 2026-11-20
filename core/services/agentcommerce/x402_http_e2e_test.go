package agentcommerce

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type x402AuthorizerFunc func(context.Context, X402Challenge, string) (string, error)

func (f x402AuthorizerFunc) AuthorizeX402(ctx context.Context, challenge X402Challenge, requestHash string) (string, error) {
	return f(ctx, challenge, requestHash)
}

type settlementVerifierFunc func(context.Context, RailSettlementRequest, RailSettlementReceipt) error

func (f settlementVerifierFunc) VerifySettlementEvidence(ctx context.Context, request RailSettlementRequest, receipt RailSettlementReceipt) error {
	return f(ctx, request, receipt)
}

type deterministicX402Facilitator struct {
	mu           sync.Mutex
	request      RailSettlementRequest
	requestHash  string
	paymentID    string
	observations []SettlementPhase
	observation  int
	malformed    bool
	reject       bool
}

func (f *deterministicX402Facilitator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		f.mu.Lock()
		defer f.mu.Unlock()
		if r.URL.Query().Get("payment_id") != f.paymentID {
			http.Error(w, "unknown payment", http.StatusNotFound)
			return
		}
		if f.malformed {
			_, _ = w.Write([]byte("{"))
			return
		}
		phase := f.observations[f.observation]
		if f.observation < len(f.observations)-1 {
			f.observation++
		}
		finality := ""
		if phase == SettlementSettled {
			sum := sha256.Sum256([]byte(f.requestHash + ":" + f.paymentID))
			finality = hex.EncodeToString(sum[:])
		}
		_ = json.NewEncoder(w).Encode(RailSettlementReceipt{Version: "1", RequestHash: f.requestHash, Rail: "x402", Phase: phase, ExternalID: f.paymentID, FinalityHash: finality, ObservedAt: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)})
		return
	}
	var request RailSettlementRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	requestHash, err := SettlementRequestHash(request)
	if err != nil {
		http.Error(w, "bad identity", http.StatusBadRequest)
		return
	}
	if r.Header.Get("PAYMENT-SIGNATURE") == "" {
		challenge := X402Challenge{X402Version: 2}
		challenge.Accepted.Scheme, challenge.Accepted.Network = "exact", request.Network
		challenge.Accepted.Asset, challenge.Accepted.Amount, challenge.Accepted.PayTo = request.Asset, request.Amount, request.Destination
		payload, _ := json.Marshal(challenge)
		w.Header().Set("PAYMENT-REQUIRED", base64.StdEncoding.EncodeToString(payload))
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
	if f.reject || r.Header.Get("PAYMENT-SIGNATURE") != "payment:"+requestHash {
		http.Error(w, "rejected", http.StatusPaymentRequired)
		return
	}
	f.mu.Lock()
	if f.requestHash == "" {
		f.request, f.requestHash = request, requestHash
	}
	if f.requestHash != requestHash {
		f.mu.Unlock()
		http.Error(w, "conflict", http.StatusConflict)
		return
	}
	f.mu.Unlock()
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(X402Submission{RequestHash: requestHash, PaymentID: f.paymentID, Phase: SettlementSubmitted})
}

func TestHTTPX402LifecycleE2E(t *testing.T) {
	req := validRailRequest(t, "x402")
	facilitator := &deterministicX402Facilitator{paymentID: "local-payment-1", observations: []SettlementPhase{SettlementBroadcast, SettlementUnknown, SettlementConfirmed, SettlementSettled}}
	server := httptest.NewServer(facilitator)
	defer server.Close()
	client, err := NewHTTPX402Client(server.Client(), server.URL, x402AuthorizerFunc(func(_ context.Context, _ X402Challenge, hash string) (string, error) { return "payment:" + hash, nil }))
	if err != nil {
		t.Fatal(err)
	}
	adapter, _ := NewX402SettlementAdapter(client)
	submission, err := adapter.Submit(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if submission.Phase != SettlementSubmitted {
		t.Fatalf("phase=%s", submission.Phase)
	}
	var receipt RailSettlementReceipt
	for _, expected := range []SettlementPhase{SettlementBroadcast, SettlementUnknown, SettlementConfirmed, SettlementSettled} {
		receipt, err = adapter.Observe(context.Background(), req, submission.PaymentID)
		if err != nil || receipt.Phase != expected {
			t.Fatalf("observation want=%s got=%s err=%v", expected, receipt.Phase, err)
		}
	}
	verifier := settlementVerifierFunc(func(_ context.Context, request RailSettlementRequest, got RailSettlementReceipt) error {
		want := sha256.Sum256([]byte(facilitator.requestHash + ":" + facilitator.paymentID))
		if got.FinalityHash != hex.EncodeToString(want[:]) || request.Scope != facilitator.request.Scope {
			return errors.New("invalid facilitator finality")
		}
		return nil
	})
	if err := VerifyFinalSettlement(context.Background(), req, receipt, verifier); err != nil {
		t.Fatal(err)
	}
	duplicate, err := adapter.Submit(context.Background(), req)
	if err != nil || duplicate.PaymentID != submission.PaymentID {
		t.Fatalf("duplicate was not idempotent: %+v %v", duplicate, err)
	}
}

func TestHTTPX402NegativeMatrix(t *testing.T) {
	base := validRailRequest(t, "x402")
	tests := []struct {
		name   string
		mutate func(*X402Challenge)
		raw    string
	}{
		{name: "malformed", raw: "not-base64"},
		{name: "wrong network", mutate: func(c *X402Challenge) { c.Accepted.Network = "wrong" }},
		{name: "wrong asset", mutate: func(c *X402Challenge) { c.Accepted.Asset = "WRONG" }},
		{name: "wrong amount", mutate: func(c *X402Challenge) { c.Accepted.Amount = "11" }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				header := test.raw
				if header == "" {
					challenge := X402Challenge{X402Version: 2}
					challenge.Accepted.Scheme = "exact"
					challenge.Accepted.Network = base.Network
					challenge.Accepted.Asset = base.Asset
					challenge.Accepted.Amount = base.Amount
					challenge.Accepted.PayTo = base.Destination
					test.mutate(&challenge)
					payload, _ := json.Marshal(challenge)
					header = base64.StdEncoding.EncodeToString(payload)
				}
				w.Header().Set("PAYMENT-REQUIRED", header)
				w.WriteHeader(http.StatusPaymentRequired)
			}))
			defer server.Close()
			client, _ := NewHTTPX402Client(server.Client(), server.URL, x402AuthorizerFunc(func(context.Context, X402Challenge, string) (string, error) { return "payment", nil }))
			adapter, _ := NewX402SettlementAdapter(client)
			if _, err := adapter.Submit(context.Background(), base); err == nil {
				t.Fatal("expected rejection")
			}
		})
	}
	wrongExecution := base
	wrongExecution.ExecutionID = "exec-wrong"
	if err := validateRailSettlementRequest(wrongExecution); err == nil {
		t.Fatal("wrong execution id accepted")
	}
	wrongIntent := base
	wrongIntent.IntentHash = fmt.Sprintf("%064x", 1)
	if err := validateRailSettlementRequest(wrongIntent); err == nil {
		t.Fatal("wrong intent binding accepted")
	}
}

func TestHTTPX402RejectionMalformedResponseTimeoutAndCancellation(t *testing.T) {
	req := validRailRequest(t, "x402")
	for _, tc := range []struct {
		name    string
		handler http.Handler
		timeout time.Duration
	}{
		{name: "rejected", handler: &deterministicX402Facilitator{paymentID: "p", observations: []SettlementPhase{SettlementFailed}, reject: true}},
		{name: "malformed response", handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("PAYMENT-SIGNATURE") == "" {
				c := X402Challenge{X402Version: 2}
				c.Accepted.Scheme = "exact"
				c.Accepted.Network = req.Network
				c.Accepted.Asset = req.Asset
				c.Accepted.Amount = req.Amount
				c.Accepted.PayTo = req.Destination
				b, _ := json.Marshal(c)
				w.Header().Set("PAYMENT-REQUIRED", base64.StdEncoding.EncodeToString(b))
				w.WriteHeader(402)
				return
			}
			w.WriteHeader(202)
			_, _ = w.Write([]byte("{"))
		})},
		{name: "timeout", handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-r.Context().Done():
				return
			case <-time.After(100 * time.Millisecond):
				http.Error(w, "late", http.StatusGatewayTimeout)
			}
		}), timeout: 10 * time.Millisecond},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()
			httpClient := server.Client()
			if tc.timeout > 0 {
				httpClient.Timeout = tc.timeout
			}
			client, _ := NewHTTPX402Client(httpClient, server.URL, x402AuthorizerFunc(func(context.Context, X402Challenge, string) (string, error) { return "payment:invalid", nil }))
			adapter, _ := NewX402SettlementAdapter(client)
			if _, err := adapter.Submit(context.Background(), req); err == nil {
				t.Fatal("expected error")
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client, _ := NewHTTPX402Client(http.DefaultClient, "http://127.0.0.1", x402AuthorizerFunc(func(context.Context, X402Challenge, string) (string, error) { return "x", nil }))
	adapter, _ := NewX402SettlementAdapter(client)
	if _, err := adapter.Submit(ctx, req); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation=%v", err)
	}
}
