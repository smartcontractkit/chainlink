// Package main implements a simple HTTP echo server for testing CRE workflows.
// It listens on a configurable port, logs all incoming POST request bodies,
// and returns a JSON success response matching the OrderResponse struct
// expected by the http_simple workflow.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// OrderResponse fields are intentionally static for consensus compatibility.
// See: ConsensusIdenticalAggregation requires all nodes to observe the same value.

type OrderResponse struct {
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

var (
	requestCount atomic.Int64
	responseDelay time.Duration
)

func handler(w http.ResponseWriter, r *http.Request) {
	count := requestCount.Add(1)
	start := time.Now()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[#%d] ERROR reading body: %v", count, err)
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("[#%d] %s %s from %s (%d bytes)", count, r.Method, r.URL.Path, r.RemoteAddr, len(body))

	if responseDelay > 0 {
		log.Printf("[#%d] Delaying response by %s...", count, responseDelay)
		time.Sleep(responseDelay)
	}

	// IMPORTANT: Return static/deterministic response so that all DON nodes
	// observe identical values. This is required for ConsensusIdenticalAggregation
	// to succeed (f+1 nodes must agree on the same value).
	resp := OrderResponse{
		Status:  "success",
		Message: "Order placed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("[#%d] ERROR encoding response: %v", count, err)
	}

	log.Printf("[#%d] Response sent (200 OK, took %s)", count, time.Since(start).Round(time.Millisecond))
}

func main() {
	port := flag.Int("port", 9999, "Port to listen on")
	delay := flag.Duration("delay", 0, "Delay before responding (e.g. 77s)")
	flag.Parse()

	responseDelay = *delay

	addr := fmt.Sprintf(":%d", *port)

	http.HandleFunc("/", handler)

	log.Printf("Echo server starting on %s", addr)
	if responseDelay > 0 {
		log.Printf("Response delay: %s", responseDelay)
	}
	log.Printf("All routes handled - POST to any path to receive a JSON response")
	log.Printf("Press Ctrl+C to stop")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
