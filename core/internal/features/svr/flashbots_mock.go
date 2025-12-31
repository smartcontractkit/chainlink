package svr

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

// FlashbotsMockServer is a mock HTTP server that receives secondary transactions
// (Flashbots endpoint) but does NOT submit them to the chain.
// This allows the test to distinguish between primary (on-chain) and secondary (Flashbots-only) transmissions.
type FlashbotsMockServer struct {
	*httptest.Server
	t           *testing.T
	receivedTxs sync.Map // map[string]bool - track received transactions by hash
	txCount     atomic.Int32
}

// setupFlashbotsMock creates and starts a new FlashbotsMockServer.
func setupFlashbotsMock(t *testing.T) *FlashbotsMockServer {
	mock := &FlashbotsMockServer{
		t: t,
	}
	mock.Server = httptest.NewServer(http.HandlerFunc(mock.handleRequest))
	return mock
}

// URL returns the URL of the mock server with "/flashbots" in the path.
// This ensures the Flashbots client is created (it only creates if URL contains "flashbots").
func (m *FlashbotsMockServer) URL() string {
	baseURL, err := url.Parse(m.Server.URL)
	if err != nil {
		// Fallback to original URL if parsing fails
		return m.Server.URL + "/flashbots"
	}
	baseURL.Path = "/flashbots"
	return baseURL.String()
}

// TransactionCount returns the number of transactions received by the mock server.
func (m *FlashbotsMockServer) TransactionCount() int32 {
	return m.txCount.Load()
}

// WasReceived checks if a transaction with the given hash was received by the mock server.
func (m *FlashbotsMockServer) WasReceived(txHash string) bool {
	_, ok := m.receivedTxs.Load(txHash)
	return ok
}

// handleRequest handles incoming HTTP requests to the mock Flashbots server.
func (m *FlashbotsMockServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		m.t.Logf("Mock Flashbots server: failed to read request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	m.t.Logf("Mock Flashbots server: received request: %s", string(body))

	// Try to parse as JSON-RPC format
	var jsonRPCReq struct {
		JSONRPC string          `json:"jsonrpc"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params"`
		ID      interface{}     `json:"id"`
	}

	// Try to parse as plain JSON with transaction data
	var txData struct {
		Tx       string `json:"tx"`       // hex encoded transaction
		TxHash   string `json:"txHash"`   // transaction hash
		RawTx    string `json:"rawTx"`    // raw transaction bytes
		SignedTx string `json:"signedTx"` // signed transaction
	}

	var txBytes []byte
	var txHash string

	// First try JSON-RPC format
	if err := json.Unmarshal(body, &jsonRPCReq); err == nil && jsonRPCReq.Method != "" {
		// Parse params - could be array or object
		var params []interface{}
		if err := json.Unmarshal(jsonRPCReq.Params, &params); err == nil && len(params) > 0 {
			// Params might be an array with transaction data
			if txStr, ok := params[0].(string); ok {
				txBytes, _ = hex.DecodeString(strings.TrimPrefix(txStr, "0x"))
			}
		} else {
			// Try as object
			var paramsObj map[string]interface{}
			if err := json.Unmarshal(jsonRPCReq.Params, &paramsObj); err == nil {
				if tx, ok := paramsObj["tx"].(string); ok {
					txBytes, _ = hex.DecodeString(strings.TrimPrefix(tx, "0x"))
				}
			}
		}
	} else if err := json.Unmarshal(body, &txData); err == nil {
		// Try plain JSON format
		if txData.RawTx != "" {
			txBytes, _ = hex.DecodeString(strings.TrimPrefix(txData.RawTx, "0x"))
			txHash = txData.TxHash
		} else if txData.Tx != "" {
			txBytes, _ = hex.DecodeString(strings.TrimPrefix(txData.Tx, "0x"))
		} else if txData.SignedTx != "" {
			txBytes, _ = hex.DecodeString(strings.TrimPrefix(txData.SignedTx, "0x"))
		}
	} else {
		// Try to decode as raw hex
		bodyStr := strings.TrimSpace(string(body))
		if strings.HasPrefix(bodyStr, "0x") || strings.HasPrefix(bodyStr, "\"0x") {
			bodyStr = strings.Trim(bodyStr, "\"")
			txBytes, _ = hex.DecodeString(strings.TrimPrefix(bodyStr, "0x"))
		}
	}

	if len(txBytes) == 0 {
		m.t.Logf("Mock Flashbots server: could not extract transaction from request body: %s", string(body))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "could not parse transaction"}`))
		return
	}

	// Decode the transaction
	var tx gethtypes.Transaction
	if err := tx.UnmarshalBinary(txBytes); err != nil {
		m.t.Logf("Mock Flashbots server: failed to unmarshal transaction: %v, body: %s", err, string(body))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid transaction format"}`))
		return
	}

	// Calculate hash if not provided
	if txHash == "" {
		txHash = tx.Hash().Hex()
	}

	// Check if we've already received this transaction
	if _, exists := m.receivedTxs.LoadOrStore(txHash, true); exists {
		m.t.Logf("Mock Flashbots server: duplicate transaction %s", txHash)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "duplicate"}`))
		return
	}

	// Track received transaction (but do NOT submit to chain)
	m.txCount.Add(1)
	m.t.Logf("Mock Flashbots server: received secondary transaction %s (count: %d) - NOT submitting to chain", txHash, m.TransactionCount())

	// Return success response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(fmt.Appendf(nil, `{"status": "success", "txHash": "%s"}`, txHash))
	if err != nil {
		m.t.Logf("Mock Flashbots server: failed to write response: %v", err)
	}
}
