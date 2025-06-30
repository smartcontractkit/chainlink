package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
// var endpoints = []Endpoint{
// 	{"LinkPool", "https://rpcs.cldev.sh/base/sepolia/linkpool1"},
// 	{"Chainstack", "https://rpcs.cldev.sh/base/sepolia/chainstack1"},
// 	{"SimplyVC", "https://rpcs.cldev.sh/base/sepolia/simplyvc1"},
// }

// var endpoints = []Endpoint{
// 	{"Chainstack", "https://rpcs.cldev.sh/base/mainnet/chainstack1"},
// 	{"LinkPool", "https://rpcs.cldev.sh/base/mainnet/linkpool1"},
// 	{"SimplyVC", "https://rpcs.cldev.sh/base/mainnet/simplyvc1"},
// }

var endpoints = []Endpoint{
	{"SimplyVC", "https://rpcs.cldev.sh/optimism/sepolia/simplyvc1"},
	{"LinkPool", "https://rpcs.cldev.sh/optimism/sepolia/linkpool1"},
	{"Chainstack", "https://rpcs.cldev.sh/optimism/sepolia/chainstack1"},
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
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var br BlockResponse
	if err := json.Unmarshal(body, &br); err != nil {
		return 0, err
	}
	if br.Result == nil || br.Result.Number == "" {
		return 0, fmt.Errorf("no block number in response: %s", string(body))
	}
	var blockNum uint64
	_, err = fmt.Sscanf(br.Result.Number, "0x%x", &blockNum)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block number: %v", err)
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

	for {
		blockNumbers := make([]uint64, len(endpoints))
		min, max := uint64(^uint64(0)), uint64(0)
		hasValidData := false

		for i, ep := range endpoints {
			num, err := getFinalizedBlockNumber(ep.URL)
			if err != nil {
				blockNumbers[i] = 0
				continue
			}
			blockNumbers[i] = num
			hasValidData = true
			if num < min {
				min = num
			}
			if num > max {
				max = num
			}
		}

		if !hasValidData {
			fmt.Printf("[%s] All endpoints failed, retrying...\n", time.Now().Format(time.RFC3339))
			time.Sleep(time.Duration(*interval) * time.Second)
			continue
		}

		currentDiff := max - min
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
