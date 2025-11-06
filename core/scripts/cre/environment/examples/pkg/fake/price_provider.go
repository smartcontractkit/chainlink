package fake

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"time"
)

type PriceResponse struct {
	AccountName string  `json:"accountName"`
	TotalTrust  float64 `json:"totalTrust"`
	Ripcord     bool    `json:"ripcord"`
	UpdatedAt   string  `json:"updatedAt"`
}

type PriceServer struct {
	authKey   string
	priceData map[string][]float64
	server    *http.Server
}

// DeployPriceProvider starts a local HTTP server serving fake price data
func DeployPriceProvider(authKey string, port int, feedIDs []string, containerName string) (string, error) {
	// Generate prices for each feed ID
	priceData := make(map[string][]float64)
	for _, feedID := range feedIDs {
		// Clean feed ID (remove 0x prefix and limit to 32 chars)
		cleanFeedID := cleanFeedID(feedID)

		// Generate 3 random prices between 50.00 and 150.00
		prices := make([]float64, 3)
		for i := range 3 {
			prices[i] = math.Round((rand.Float64()*100+50)*100) / 100 //nolint:gosec // this is a fake price generator
		}
		priceData[cleanFeedID] = prices

		fmt.Printf("📊 Generated prices for feed %s: %.2f, %.2f, %.2f\n",
			cleanFeedID, prices[0], prices[1], prices[2])
	}

	// Create and start the server
	server := &PriceServer{
		authKey:   authKey,
		priceData: priceData,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/fake/api/price", server.handlePriceRequest)

	server.server = &http.Server{ //nolint:gosec // this is a fake price provider
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// Start server in background
	go func() {
		fmt.Printf("🚀 Starting fake price provider server on port %d...\n", port)
		if err := server.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ Server error: %v\n", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	url := fmt.Sprintf("http://localhost:%d/fake/api/price", port)
	fmt.Printf("✅ Fake price provider deployed successfully!\n")
	fmt.Printf("🌐 URL: %s\n", url)
	fmt.Printf("🔑 Auth Key: %s\n", authKey)
	fmt.Printf("📊 Feed IDs: %v\n", feedIDs)
	fmt.Printf("🔄 The service is now running locally.\n")
	fmt.Printf("⚠️  Press Ctrl+C in the terminal to stop the service.\n")

	return url, nil
}

func (ps *PriceServer) handlePriceRequest(w http.ResponseWriter, r *http.Request) {
	// Check auth header if provided
	if ps.authKey != "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader != ps.authKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			if err != nil {
				fmt.Println("Error encoding unauthorized response:", err)
			}
			return
		}
	}

	feedID := r.URL.Query().Get("feedID")
	if feedID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]string{"error": "no feedID provided"})
		if err != nil {
			fmt.Println("Error encoding bad request response:", err)
		}
		return
	}

	// Clean the feed ID
	cleanID := cleanFeedID(feedID)

	// Get price for this feed
	prices, exists := ps.priceData[cleanID]
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(map[string]string{"error": "feedID not found"})
		if err != nil {
			fmt.Println("Error encoding not found response:", err)
		}
		return
	}

	// Return the first price (you could cycle through them if needed)
	currentPrice := prices[0]
	response := PriceResponse{
		AccountName: "TrueUSD",
		TotalTrust:  currentPrice,
		Ripcord:     false,
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Error encoding price response:", err)
	}
}

func cleanFeedID(feedID string) string {
	cleanID := feedID
	if len(cleanID) > 2 && cleanID[:2] == "0x" {
		cleanID = cleanID[2:]
	}
	if len(cleanID) > 32 {
		cleanID = cleanID[:32]
	}
	return "0x" + cleanID
}
