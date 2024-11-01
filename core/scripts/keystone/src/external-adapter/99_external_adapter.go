package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"time"
)

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Price struct encapsulates bid, mid, ask values along with a mutex for synchronization
type Price struct {
	mu  sync.RWMutex
	Bid float64
	Mid float64
	Ask float64
}

// Update safely updates the price values within the specified bounds
func (p *Price) Update(step, floor, ceiling float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Mid = adjustValue(p.Mid, step, floor, ceiling)
	p.Bid = adjustValue(p.Mid, step, floor, p.Mid)
	p.Ask = adjustValue(p.Mid, step, p.Mid, ceiling)
}

// GetSnapshot safely retrieves a copy of the current price values
func (p *Price) GetSnapshot() (bid, mid, ask float64) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Bid, p.Mid, p.Ask
}

func main() {
	// Get initial values from environment variables or use defaults
	btcUsdInitialValue := getInitialValue("BTCUSD_INITIAL_VALUE", 1000.0)
	linkInitialValue := getInitialValue("LINK_INITIAL_VALUE", 11.0)
	nativeInitialValue := getInitialValue("NATIVE_INITIAL_VALUE", 2400.0)

	pctBounds := 0.3

	// Start external adapters on different ports
	externalAdapter(btcUsdInitialValue, "4001", pctBounds)
	externalAdapter(linkInitialValue, "4002", pctBounds)
	externalAdapter(nativeInitialValue, "4003", pctBounds)

	// Block main goroutine indefinitely
	select {}
}

// getInitialValue retrieves the initial value from the environment or returns a default
func getInitialValue(envVar string, defaultValue float64) float64 {
	valueEnv := os.Getenv(envVar)
	if valueEnv == "" {
		fmt.Printf("%s not set, using default value: %.4f\n", envVar, defaultValue)
		return defaultValue
	}
	fmt.Printf("%s set to %s\n", envVar, valueEnv)
	val, err := strconv.ParseFloat(valueEnv, 64)
	PanicErr(err)
	return val
}

// externalAdapter sets up a mock external adapter server for a specific asset
func externalAdapter(initialValue float64, port string, pctBounds float64) *httptest.Server {
	// Create a custom listener on the specified port
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		panic(err)
	}

	// Initialize the Price struct
	price := &Price{
		Bid: initialValue,
		Mid: initialValue,
		Ask: initialValue,
	}

	step := initialValue * pctBounds / 10
	ceiling := initialValue * (1 + pctBounds)
	floor := initialValue * (1 - pctBounds)

	// Perform initial adjustment to set bid and ask
	price.Update(step, floor, ceiling)

	// Start a goroutine to periodically update the price
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			price.Update(step, floor, ceiling)
			fmt.Printf("Updated prices on port %s: bid=%.4f, mid=%.4f, ask=%.4f\n", port, price.Bid, price.Mid, price.Ask)
		}
	}()

	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		bid, mid, ask := price.GetSnapshot()

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		resp := fmt.Sprintf(`{"result": {"bid": %.4f, "mid": %.4f, "ask": %.4f}}`, bid, mid, ask)
		if _, err := res.Write([]byte(resp)); err != nil {
			fmt.Printf("failed to write response: %v\n", err)
		}
	})

	// Create and start the test server
	ea := &httptest.Server{
		Listener: listener,
		Config: &http.Server{
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
	ea.Start()

	fmt.Printf("Mock external adapter started at %s\n", ea.URL)
	fmt.Printf("Initial value: %.4f, Floor: %.4f, Ceiling: %.4f\n", initialValue, floor, ceiling)
	return ea
}

// adjustValue takes a starting value and randomly shifts it up or down by a step.
// It ensures that the value stays within the specified bounds.
func adjustValue(start, step, floor, ceiling float64) float64 {
	// Randomly choose to increase or decrease the value
	// #nosec G404
	if rand.Intn(2) == 0 {
		step = -step
	}

	// Apply the step to the starting value
	newValue := start + step

	// Ensure the value is within the bounds
	if newValue < floor {
		newValue = floor
	} else if newValue > ceiling {
		newValue = ceiling
	}

	return newValue
}
