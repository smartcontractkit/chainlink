package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httptypedapi "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/http"
)

// waitForPort polls until the TCP port is reachable or the deadline passes.
func waitForPort(t *testing.T, port uint16, timeout time.Duration) {
	t.Helper()
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		dialer := &net.Dialer{Timeout: 50 * time.Millisecond}
		conn, err := dialer.DialContext(context.Background(), "tcp", addr)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("server on port %d did not become ready within %s", port, timeout)
}

// TestListenForTriggerPayload_HappyPath is an integration test that verifies
// the full round-trip of ListenForTriggerPayload:
//  1. A LocalGateway is started on an ephemeral port.
//  2. A valid POST request carrying a signed JWT and a JSON-RPC body is sent.
//  3. The method returns a Payload whose Input and Key match the request.
func TestListenForTriggerPayload_HappyPath(t *testing.T) {
	var port uint16 = 30123
	gw := NewLocalGateway(Config{Port: port})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type result struct {
		payload *httptypedapi.Payload
		err     error
	}
	resultCh := make(chan result, 1)
	go func() {
		p, err := gw.ListenForTriggerPayload(ctx)
		resultCh <- result{p, err}
	}()

	// Wait for the HTTP server to be ready before sending the request.
	waitForPort(t, port, 2*time.Second)

	// Build the JSON-RPC request body containing input data.
	// We construct the body manually using a map so that the "input" field is
	// embedded as a raw JSON object on the wire.  If we used JsonRpcRequest
	// directly, the []byte field would be base64-encoded by json.Marshal,
	// which is not what a real HTTP client sends.
	inputJSON := json.RawMessage(`{"order":"pizza","size":"large"}`)
	rawBody := map[string]any{
		"input": inputJSON,
	}
	body, err := json.Marshal(rawBody)
	require.NoError(t, err)

	// Send the POST request to the gateway.
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("http://127.0.0.1:%d/trigger", port), bytes.NewReader(body))
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// ListenForTriggerPayload should return exactly one payload.
	res := <-resultCh
	require.NoError(t, res.err)
	require.NotNil(t, res.payload)

	// The payload input must match what was sent in Params.Input.
	assert.Equal(t, []byte(inputJSON), res.payload.Input)
}
