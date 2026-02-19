package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

func TestStartComponentReturnsErrorCodeForUnsupportedSchema(t *testing.T) {
	server := NewServer(zerolog.Nop(), nil)
	handler := server.Handler()

	req := httptest.NewRequest(http.MethodPost, "/v1/components/start", strings.NewReader(`{"schemaVersion":"v0","operation":"StartComponent","payload":{}}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), ErrCodeUnsupportedSchema) {
		t.Fatalf("expected response to include error code %q, got body: %s", ErrCodeUnsupportedSchema, rr.Body.String())
	}
}

func TestStartComponentReturnsErrorCodeForUnsupportedComponent(t *testing.T) {
	server := NewServer(zerolog.Nop(), nil)
	handler := server.Handler()

	body := bytes.NewBufferString(`{"schemaVersion":"v1","operation":"StartComponent","payload":{"componentType":"not-supported"}}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/components/start", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), ErrCodeUnsupportedComponent) {
		t.Fatalf("expected response to include error code %q, got body: %s", ErrCodeUnsupportedComponent, rr.Body.String())
	}
}

type fakeOutputDeployer struct {
	calls int
}

func (f *fakeOutputDeployer) Deploy(context.Context, *blockchain.Input) (blockchains.Blockchain, error) {
	return nil, nil
}

func (f *fakeOutputDeployer) DeployOutput(context.Context, *blockchain.Input) (*blockchain.Output, error) {
	f.calls++
	return &blockchain.Output{
		Type:    blockchain.TypeAnvil,
		ChainID: "1337",
		Nodes: []*blockchain.Node{
			{
				ExternalHTTPUrl: "http://127.0.0.1:8545",
				ExternalWSUrl:   "ws://127.0.0.1:8546",
			},
		},
	}, nil
}

func TestStartComponentReuseIfIdenticalPayload(t *testing.T) {
	deployer := &fakeOutputDeployer{}
	server := NewServer(zerolog.Nop(), map[blockchain.ChainFamily]blockchains.Deployer{
		blockchain.FamilyEVM: deployer,
	})
	handler := server.Handler()

	payload := `{"componentType":"blockchain","blockchain":{"type":"anvil","chain_id":"1337"}}`
	body := bytes.NewBufferString(`{"schemaVersion":"v1","operation":"StartComponent","payload":` + payload + `}`)

	req1 := httptest.NewRequest(http.MethodPost, "/v1/components/start", bytes.NewReader(body.Bytes()))
	req1.Header.Set("Content-Type", "application/json")
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Fatalf("expected first request OK, got %d: %s", rr1.Code, rr1.Body.String())
	}

	req2 := httptest.NewRequest(http.MethodPost, "/v1/components/start", bytes.NewReader(body.Bytes()))
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("expected second request OK, got %d: %s", rr2.Code, rr2.Body.String())
	}

	if deployer.calls != 1 {
		t.Fatalf("expected deployer to be called once with reuse mode, got %d", deployer.calls)
	}

	var resp StartComponentResponse
	if err := json.Unmarshal(rr2.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.AgentLogs) == 0 || !strings.Contains(strings.Join(resp.AgentLogs, " "), "reusing existing component") {
		t.Fatalf("expected reuse log in response, got: %v", resp.AgentLogs)
	}
}

func TestStartComponentAlwaysPolicyDisablesReuse(t *testing.T) {
	deployer := &fakeOutputDeployer{}
	server := NewServer(zerolog.Nop(), map[blockchain.ChainFamily]blockchains.Deployer{
		blockchain.FamilyEVM: deployer,
	})
	handler := server.Handler()

	payload := `{"componentType":"blockchain","reusePolicy":"always","blockchain":{"type":"anvil","chain_id":"1337"}}`
	body := bytes.NewBufferString(`{"schemaVersion":"v1","operation":"StartComponent","payload":` + payload + `}`)

	req1 := httptest.NewRequest(http.MethodPost, "/v1/components/start", bytes.NewReader(body.Bytes()))
	req1.Header.Set("Content-Type", "application/json")
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Fatalf("expected first request OK, got %d: %s", rr1.Code, rr1.Body.String())
	}

	req2 := httptest.NewRequest(http.MethodPost, "/v1/components/start", bytes.NewReader(body.Bytes()))
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("expected second request OK, got %d: %s", rr2.Code, rr2.Body.String())
	}

	if deployer.calls != 2 {
		t.Fatalf("expected deployer to be called twice with always policy, got %d", deployer.calls)
	}
}

func TestStartComponentRequiresJDPayload(t *testing.T) {
	server := NewServer(zerolog.Nop(), nil)
	handler := server.Handler()

	body := bytes.NewBufferString(`{"schemaVersion":"v1","operation":"StartComponent","payload":{"componentType":"jd"}}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/components/start", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), ErrCodeMissingComponentInput) {
		t.Fatalf("expected missing input error code in response, got body: %s", rr.Body.String())
	}
}

func TestShouldReuseRemoteStartDisablesJDReuse(t *testing.T) {
	if shouldReuseRemoteStart(ComponentTypeJD, RemoteStartPolicyReuseIdentical) {
		t.Fatal("expected JD reuse to be hard disabled")
	}
	if !shouldReuseRemoteStart(ComponentTypeBlockchain, "") {
		t.Fatal("expected blockchain reuse to default to enabled")
	}
}

func TestStopComponentIdempotent(t *testing.T) {
	deployer := &fakeOutputDeployer{}
	server := NewServer(zerolog.Nop(), map[blockchain.ChainFamily]blockchains.Deployer{
		blockchain.FamilyEVM: deployer,
	})
	handler := server.Handler()

	startPayload := `{"componentType":"blockchain","blockchain":{"type":"anvil","chain_id":"1337"}}`
	startBody := bytes.NewBufferString(`{"schemaVersion":"v1","operation":"StartComponent","payload":` + startPayload + `}`)
	startReq := httptest.NewRequest(http.MethodPost, "/v1/components/start", bytes.NewReader(startBody.Bytes()))
	startReq.Header.Set("Content-Type", "application/json")
	startRR := httptest.NewRecorder()
	handler.ServeHTTP(startRR, startReq)
	if startRR.Code != http.StatusOK {
		t.Fatalf("expected start request OK, got %d: %s", startRR.Code, startRR.Body.String())
	}

	stopBody := bytes.NewBufferString(`{"schemaVersion":"v1","operation":"StopComponent","payload":` + startPayload + `}`)
	stopReq1 := httptest.NewRequest(http.MethodPost, "/v1/components/start", bytes.NewReader(stopBody.Bytes()))
	stopReq1.Header.Set("Content-Type", "application/json")
	stopRR1 := httptest.NewRecorder()
	handler.ServeHTTP(stopRR1, stopReq1)
	if stopRR1.Code != http.StatusOK {
		t.Fatalf("expected first stop request OK, got %d: %s", stopRR1.Code, stopRR1.Body.String())
	}

	var stopResp1 StartComponentResponse
	if err := json.Unmarshal(stopRR1.Body.Bytes(), &stopResp1); err != nil {
		t.Fatalf("failed to decode first stop response: %v", err)
	}
	if !stopResp1.Found || !stopResp1.Stopped {
		t.Fatalf("expected first stop to find and stop component, got found=%v stopped=%v", stopResp1.Found, stopResp1.Stopped)
	}

	stopReq2 := httptest.NewRequest(http.MethodPost, "/v1/components/start", bytes.NewReader(stopBody.Bytes()))
	stopReq2.Header.Set("Content-Type", "application/json")
	stopRR2 := httptest.NewRecorder()
	handler.ServeHTTP(stopRR2, stopReq2)
	if stopRR2.Code != http.StatusOK {
		t.Fatalf("expected second stop request OK, got %d: %s", stopRR2.Code, stopRR2.Body.String())
	}

	var stopResp2 StartComponentResponse
	if err := json.Unmarshal(stopRR2.Body.Bytes(), &stopResp2); err != nil {
		t.Fatalf("failed to decode second stop response: %v", err)
	}
	if stopResp2.Found || stopResp2.Stopped {
		t.Fatalf("expected second stop to be no-op, got found=%v stopped=%v", stopResp2.Found, stopResp2.Stopped)
	}
}

func TestShouldCleanupFailedContainersDefaultsToTrue(t *testing.T) {
	t.Setenv(EnvKeepFailedContainers, "")
	if !shouldCleanupFailedContainers() {
		t.Fatal("expected cleanup to be enabled by default")
	}
}

func TestShouldCleanupFailedContainersCanBeDisabled(t *testing.T) {
	for _, value := range []string{"1", "true", "yes", "on", "TRUE"} {
		t.Setenv(EnvKeepFailedContainers, value)
		if shouldCleanupFailedContainers() {
			t.Fatalf("expected cleanup to be disabled for value %q", value)
		}
	}
}
