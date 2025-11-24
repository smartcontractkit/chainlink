package svr

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
)

// FlashbotsMockServer is a mock HTTP server that receives secondary transactions
// (Flashbots endpoint) and submits them to the chain.
type FlashbotsMockServer struct {
	*httptest.Server
	t            *testing.T
	backend      *simulated.Backend
	receivedTxs  sync.Map // map[string]bool - track received transactions by hash
	submittedTxs sync.Map // map[string]bool - track transaction hashes submitted to chain
	txCount      atomic.Int32
}

// setupFlashbotsMock creates and starts a new FlashbotsMockServer.
func setupFlashbotsMock(t *testing.T, backend *simulated.Backend) *FlashbotsMockServer {
	mock := &FlashbotsMockServer{
		t:       t,
		backend: backend,
	}
	mock.Server = httptest.NewServer(http.HandlerFunc(mock.handleRequest))
	return mock
}

// URL returns the URL of the mock server.
func (m *FlashbotsMockServer) URL() string {
	return m.Server.URL
}

// TransactionCount returns the number of transactions received by the mock server.
func (m *FlashbotsMockServer) TransactionCount() int32 {
	return m.txCount.Load()
}

// WasSubmitted checks if a transaction with the given hash was submitted to the chain.
func (m *FlashbotsMockServer) WasSubmitted(txHash string) bool {
	_, ok := m.submittedTxs.Load(txHash)
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

	// Submit transaction to chain
	m.t.Logf("Mock Flashbots server: submitting secondary transaction %s to chain", txHash)
	err = m.backend.Client().SendTransaction(context.Background(), &tx)
	if err != nil {
		m.t.Logf("Mock Flashbots server: failed to submit transaction to chain: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "failed to submit: %v"}`, err)))
		return
	}

	// Commit the block to include the transaction
	m.backend.Commit()

	m.txCount.Add(1)

	// Track that this transaction hash was submitted to chain by the mock server
	m.submittedTxs.Store(txHash, true)

	m.t.Logf("Mock Flashbots server: successfully submitted secondary transaction %s to chain (count: %d)", txHash, m.TransactionCount())

	// Return success response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(fmt.Appendf(nil, `{"status": "success", "txHash": "%s"}`, txHash))
	if err != nil {
		m.t.Logf("Mock Flashbots server: failed to write response: %v", err)
	}
}
