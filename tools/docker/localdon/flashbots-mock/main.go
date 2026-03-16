// flashbots-mock is a minimal JSON-RPC server that mimics a Flashbots/OFA endpoint.
// It responds to eth_getTransactionCount (needed by TXMv2 for nonce initialization)
// and eth_sendRawTransaction (accepts secondary txs, returns fake hash).
// Used for testing SVR dual transmission locally.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
)

var txCount atomic.Int64

func main() {
	http.HandleFunc("/", handleRequest)
	log.Println("Flashbots mock server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	defer r.Body.Close()

	var req struct {
		JSONRPC string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		ID      interface{} `json:"id"`
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.Unmarshal(body, &req); err == nil && req.Method != "" {
		var result interface{}

		switch req.Method {
		case "eth_getTransactionCount":
			result = "0x0"
		case "eth_sendRawTransaction":
			count := txCount.Add(1)
			log.Printf("Secondary tx received (total: %d)", count)
			result = fmt.Sprintf("0x%064x", count)
		case "eth_chainId":
			result = "0x7a69" // 31337
		default:
			log.Printf("Unhandled JSON-RPC method=%s", req.Method)
			result = "0x0"
		}

		resp := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result":  result,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	count := txCount.Load()
	log.Printf("Non-RPC request %s %s (secondary txs: %d)", r.Method, r.URL.Path, count)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","secondaryTxCount":%d}`, count)
}
