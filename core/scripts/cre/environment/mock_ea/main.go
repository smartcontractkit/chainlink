// Package main implements a mock External Adapter for E2E testing
// of the Streams Trigger capability.
//
// This adapter returns hardcoded price values for various trading pairs,
// simulating a real price feed source for the LLO (Streams) DON.
//
// MAGIC_NUMBER (424242) is embedded in every response. If this number
// appears in workflow logs, it proves the full E2E pipeline is working:
// Mock EA → Stream Jobs → LLO Plugin → streams-trigger → Workflow
//
// Usage:
//
//	go run main.go
//
// Environment Variables:
//
//	PORT - HTTP server port (default: 8080)
//	BTC_USD_PRICE - Bitcoin price in USD (default: 65000.00)
//	ETH_USD_PRICE - Ethereum price in USD (default: 3500.00)
//	LINK_USD_PRICE - LINK price in USD (default: 15.00)
//	VOLATILITY - Price volatility percentage (default: 0.1)
//	LOG_LEVEL - Log level: debug, info, warn, error (default: info)
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Magic numbers embedded in EA responses to prove E2E connectivity
// Different streams return different magic numbers to verify both report formats
const (
	// MAGIC_NUMBER_FORMAT5 is returned for Stream 1 (TEST/USD) - ReportFormat 5 (CapabilityTrigger)
	MAGIC_NUMBER_FORMAT5 = 424242
	// MAGIC_NUMBER_FORMAT7 is returned for Stream 4 (DATA/USD) - ReportFormat 7 (EVMABIEncodeUnpackedExpr)
	MAGIC_NUMBER_FORMAT7 = 555555
)

// Config holds the adapter configuration
type Config struct {
	Port       string
	Prices     map[string]float64
	Volatility float64
	LogLevel   string
}

// EARequest represents a Chainlink External Adapter request
type EARequest struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data"`
}

// EAResponse represents a Chainlink External Adapter response
type EAResponse struct {
	JobRunID   string      `json:"jobRunID"`
	Data       interface{} `json:"data"`
	Result     interface{} `json:"result"`
	StatusCode int         `json:"statusCode"`
}

// EAErrorResponse represents an error response
type EAErrorResponse struct {
	JobRunID   string `json:"jobRunID"`
	Status     string `json:"status"`
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

// PriceServer manages price data with optional volatility
type PriceServer struct {
	config     Config
	prices     map[string]float64
	pricesMu   sync.RWMutex
	requestLog []RequestLogEntry
	logMu      sync.Mutex
}

// RequestLogEntry records each request for debugging
type RequestLogEntry struct {
	Time      time.Time
	Pair      string
	Price     float64
	RequestID string
}

func main() {
	config := loadConfig()

	server := &PriceServer{
		config:     config,
		prices:     config.Prices,
		requestLog: make([]RequestLogEntry, 0),
	}

	http.HandleFunc("/", server.handleRequest)
	http.HandleFunc("/health", server.handleHealth)
	http.HandleFunc("/prices", server.handlePrices)
	http.HandleFunc("/logs", server.handleLogs)

	addr := ":" + config.Port
	log.Printf("Mock External Adapter starting on %s", addr)
	log.Printf("Configured prices: %+v", config.Prices)
	log.Printf("Volatility: %.2f%%", config.Volatility*100)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loadConfig() Config {
	config := Config{
		Port: getEnv("PORT", "8080"),
		Prices: map[string]float64{
			// ============================================================
			// Format 5 (CapabilityTrigger) - Stream 1
			// ============================================================
			// TEST/USD returns MAGIC_NUMBER_FORMAT5 (424242)
			"TEST/USD": float64(MAGIC_NUMBER_FORMAT5),

			// ============================================================
			// Format 7 (EVMABIEncodeUnpackedExpr) - Streams 2, 3, 4
			// ============================================================
			// Format 7 requires fee streams as first two streams:
			// - Stream 2: NATIVE/USD (for native fee calculation)
			// - Stream 3: LINK/USD (for LINK fee calculation)
			// - Stream 4+: Data streams (DATA/USD returns MAGIC_NUMBER_FORMAT7)
			"NATIVE/USD": 3000.00,                       // Native token price (e.g., ETH)
			"LINK/USD":   15.00,                         // LINK price for fees
			"DATA/USD":   float64(MAGIC_NUMBER_FORMAT7), // Data stream with magic number (555555)
		},
		Volatility: 0, // No volatility - exact match required
		LogLevel:   getEnv("LOG_LEVEL", "info"),
	}
	return config
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvFloat(key string, defaultVal float64) float64 {
	if val := os.Getenv(key); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return defaultVal
}

func (s *PriceServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req EARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, req.ID, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract the trading pair from the request
	pair := s.extractPair(req.Data)
	if pair == "" {
		s.sendError(w, req.ID, "Missing or invalid trading pair", http.StatusBadRequest)
		return
	}

	// Get the price (with optional volatility)
	price, ok := s.getPrice(pair)
	if !ok {
		s.sendError(w, req.ID, fmt.Sprintf("Unknown trading pair: %s", pair), http.StatusNotFound)
		return
	}

	// Log the request
	s.logRequest(pair, price, req.ID)

	// Send successful response with magic number for E2E verification
	// Different pairs return different magic numbers to verify both report formats
	magicNumber := 0
	switch pair {
	case "TEST/USD":
		magicNumber = MAGIC_NUMBER_FORMAT5 // 424242 - for Format 5
	case "DATA/USD":
		magicNumber = MAGIC_NUMBER_FORMAT7 // 555555 - for Format 7
	}

	response := EAResponse{
		JobRunID: req.ID,
		Data: map[string]interface{}{
			"result":       price,
			"pair":         pair,
			"magic_number": magicNumber, // This proves E2E connectivity when seen in workflow logs
		},
		Result:     price,
		StatusCode: http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	if s.config.LogLevel == "debug" {
		log.Printf("Request: pair=%s, price=%.6f, id=%s", pair, price, req.ID)
	}
}

func (s *PriceServer) extractPair(data map[string]interface{}) string {
	// Try various common formats
	if base, ok := data["base"].(string); ok {
		if quote, ok := data["quote"].(string); ok {
			return base + "/" + quote
		}
	}

	if pair, ok := data["pair"].(string); ok {
		return pair
	}

	if from, ok := data["from"].(string); ok {
		if to, ok := data["to"].(string); ok {
			return from + "/" + to
		}
	}

	// Check for stream ID (used by LLO)
	if streamID, ok := data["streamId"]; ok {
		return s.streamIDToPair(streamID)
	}

	return ""
}

func (s *PriceServer) streamIDToPair(streamID interface{}) string {
	// Map stream IDs to trading pairs
	// This matches the channel definitions in the LLO job spec
	//
	// Channel 1 (Format 5 - CapabilityTrigger):
	//   Stream 1: TEST/USD → MAGIC_NUMBER_FORMAT5 (424242)
	//
	// Channel 2 (Format 7 - EVMABIEncodeUnpackedExpr):
	//   Stream 2: NATIVE/USD (fee stream)
	//   Stream 3: LINK/USD (fee stream)
	//   Stream 4: DATA/USD → MAGIC_NUMBER_FORMAT7 (555555)
	switch id := streamID.(type) {
	case float64:
		switch int(id) {
		case 1:
			return "TEST/USD" // Format 5 data stream
		case 2:
			return "NATIVE/USD" // Format 7 fee stream (native)
		case 3:
			return "LINK/USD" // Format 7 fee stream (LINK)
		case 4:
			return "DATA/USD" // Format 7 data stream
		}
	case string:
		switch id {
		case "1":
			return "TEST/USD"
		case "2":
			return "NATIVE/USD"
		case "3":
			return "LINK/USD"
		case "4":
			return "DATA/USD"
		}
	}
	return ""
}

func (s *PriceServer) getPrice(pair string) (float64, bool) {
	s.pricesMu.RLock()
	basePrice, ok := s.prices[pair]
	s.pricesMu.RUnlock()

	if !ok {
		return 0, false
	}

	// Apply volatility if configured
	if s.config.Volatility > 0 {
		change := (rand.Float64()*2 - 1) * s.config.Volatility * basePrice
		return basePrice + change, true
	}

	return basePrice, true
}

func (s *PriceServer) logRequest(pair string, price float64, requestID string) {
	s.logMu.Lock()
	defer s.logMu.Unlock()

	entry := RequestLogEntry{
		Time:      time.Now(),
		Pair:      pair,
		Price:     price,
		RequestID: requestID,
	}
	s.requestLog = append(s.requestLog, entry)

	// Keep only last 1000 entries
	if len(s.requestLog) > 1000 {
		s.requestLog = s.requestLog[len(s.requestLog)-1000:]
	}
}

func (s *PriceServer) sendError(w http.ResponseWriter, jobRunID, errorMsg string, statusCode int) {
	response := EAErrorResponse{
		JobRunID:   jobRunID,
		Status:     "errored",
		Error:      errorMsg,
		StatusCode: statusCode,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	log.Printf("Error: %s (id=%s, status=%d)", errorMsg, jobRunID, statusCode)
}

func (s *PriceServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (s *PriceServer) handlePrices(w http.ResponseWriter, r *http.Request) {
	s.pricesMu.RLock()
	defer s.pricesMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"prices":     s.prices,
		"volatility": s.config.Volatility,
		"time":       time.Now().Format(time.RFC3339),
	})
}

func (s *PriceServer) handleLogs(w http.ResponseWriter, r *http.Request) {
	s.logMu.Lock()
	defer s.logMu.Unlock()

	// Return last 100 log entries
	start := 0
	if len(s.requestLog) > 100 {
		start = len(s.requestLog) - 100
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"logs":  s.requestLog[start:],
		"total": len(s.requestLog),
	})
}
