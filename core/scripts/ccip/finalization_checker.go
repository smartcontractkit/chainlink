package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

// Endpoint represents an RPC endpoint with a name and URL
type Endpoint struct {
	Name string
	URL  string
}

// BlockResponse is the structure for eth_getBlockByNumber response
// Only need the 'number' field
// The 'result' can be null if the node is not synced
// So we use a pointer

type Block struct {
	Number string `json:"number"`
}
type BlockResponse struct {
	Result *Block `json:"result"`
}

// MODIFY HERE THE ENDPOINTS YOU WANT TO TEST
var endpoints = []Endpoint{
	{"LinkPool", "https://rpc1"},
	{"Chainstack", "https://rpc2"},
	{"SimplyVC", "https://rpc3"},
}

func getFinalizedBlockNumber(url string) (uint64, error) {
	// JSON-RPC request
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{"finalized", false},
		"id":      1,
	}
	b, _ := json.Marshal(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var br BlockResponse
	err = json.Unmarshal(body, &br)
	if err != nil {
		return 0, err
	}
	if br.Result == nil || br.Result.Number == "" {
		return 0, fmt.Errorf("no block number in response: %s", string(body))
	}
	var blockNum uint64
	_, err = fmt.Sscanf(br.Result.Number, "0x%x", &blockNum)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block number: %w", err)
	}
	return blockNum, nil
}

func main() {
	interval := flag.Int("interval", 10, "Polling interval in seconds")
	flag.Parse()

	fmt.Printf("Monitoring finalized blocks every %d seconds. Press Ctrl+C to stop.\n", *interval)
	fmt.Println("Waiting for state changes (synchronized <-> unsynchronized)...")

	var lastSynced *bool // nil = unknown, true = synced (diff=0), false = not synced (diff>0)
	var timestamps []time.Time
	var states []string

	// Track last block number for each endpoint
	lastBlockNumbers := make(map[string]uint64)

	// Track last time we printed full status
	lastStatusPrint := time.Now()
	statusInterval := 120 * time.Second

	for {
		blockNumbers := make([]uint64, len(endpoints))
		minBlock, maxBlock := uint64(math.MaxUint64), uint64(0)
		hasValidData := false

		for i, ep := range endpoints {
			num, err := getFinalizedBlockNumber(ep.URL)
			if err != nil {
				blockNumbers[i] = 0
				continue
			}
			blockNumbers[i] = num
			hasValidData = true

			// Check if block number changed for this endpoint
			if lastBlock, exists := lastBlockNumbers[ep.Name]; !exists || lastBlock != num {
				fmt.Printf("[%s] %s block changed: %d -> %d\n",
					time.Now().Format(time.RFC3339), ep.Name, lastBlock, num)
				lastBlockNumbers[ep.Name] = num
			}

			if num < minBlock {
				minBlock = num
			}
			if num > maxBlock {
				maxBlock = num
			}
		}

		if !hasValidData {
			fmt.Printf("[%s] All endpoints failed, retrying...\n", time.Now().Format(time.RFC3339))
			time.Sleep(time.Duration(*interval) * time.Second)
			continue
		}

		// Print full status every 90 seconds
		if time.Since(lastStatusPrint) >= statusInterval {
			fmt.Printf("\n[%s] === PERIODIC STATUS CHECK ===\n", time.Now().Format(time.RFC3339))
			fmt.Println("Current finalized blocks:")
			for i, ep := range endpoints {
				if blockNumbers[i] > 0 {
					fmt.Printf("  %s: %d\n", ep.Name, blockNumbers[i])
				} else {
					fmt.Printf("  %s: ERROR\n", ep.Name)
				}
			}
			fmt.Printf("Block difference: %d (min: %d, max: %d)\n", maxBlock-minBlock, minBlock, maxBlock)
			fmt.Println("================================")
			lastStatusPrint = time.Now()
		}

		currentDiff := maxBlock - minBlock
		currentSynced := currentDiff == 0
		now := time.Now()

		// Check if state changed
		if lastSynced == nil || *lastSynced != currentSynced {
			stateStr := "UNSYNCHRONIZED"
			if currentSynced {
				stateStr = "SYNCHRONIZED"
			}

			fmt.Printf("\n[%s] STATE CHANGE: %s (block difference: %d)\n",
				now.Format(time.RFC3339), stateStr, currentDiff)

			// Print current block numbers for context
			fmt.Println("Current finalized blocks:")
			for i, ep := range endpoints {
				if blockNumbers[i] > 0 {
					fmt.Printf("  %s: %d\n", ep.Name, blockNumbers[i])
				}
			}

			// Store timestamp and state
			timestamps = append(timestamps, now)
			states = append(states, stateStr)

			// Print summary of recorded events
			fmt.Printf("\nRecorded events (%d total):\n", len(timestamps))
			for i, ts := range timestamps {
				fmt.Printf("  %d. [%s] %s\n", i+1, ts.Format(time.RFC3339), states[i])
			}

			lastSynced = &currentSynced
		}

		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
