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

	body := bytes.NewBufferString(`{"schemaVersion":"v1","operation":"StartComponent","payload":{"componentType":"jd"}}`)
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
