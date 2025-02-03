package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// Tx represents a transaction structure
type Tx struct {
	TxAttempts []TxAttempt
}

// TxAttempt represents a transaction attempt structure
type TxAttempt struct {
	Hash common.Hash
}

// zircuitResponse represents the response structure from zirc_isQuarantined
type zircuitResponse struct {
	IsQuarantined bool `json:"isQuarantined"`
}

// stuckTxDetector represents the structure for stuck transaction detection
type stuckTxDetector struct {
	chainClient *rpc.Client
	lggr        Logger
}

// Logger interface for logging (simplified here for demonstration)
type Logger struct{}

func (l Logger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// detectFraudTransactionsZircuit detects quarantined transactions
func (d *stuckTxDetector) detectFraudTransactionsZircuit(ctx context.Context, txs []Tx) ([]Tx, error) {
	txReqs := make([]rpc.BatchElem, len(txs))
	txHashMap := make(map[common.Hash]Tx)
	txRes := make([]*zircuitResponse, len(txs))

	// Build batch request elements
	for i, tx := range txs {
		latestAttemptHash := tx.TxAttempts[0].Hash
		var result zircuitResponse
		txReqs[i] = rpc.BatchElem{
			Method: "zirc_isQuarantined",
			Args:   []interface{}{latestAttemptHash},
			Result: &result,
		}
		txHashMap[latestAttemptHash] = tx
		txRes[i] = &result
	}

	// Send batch request using the Ethereum client
	err := d.chainClient.BatchCallContext(ctx, txReqs)
	if err != nil {
		return nil, fmt.Errorf("failed to check Quarantine transactions in batch: %w", err)
	}

	// Process the results and check if any transactions are quarantined
	var fraudTxs []Tx
	for i, req := range txReqs {
		txHash := req.Args[0].(common.Hash)
		if req.Error != nil {
			d.lggr.Debugf("failed to check fraud transaction by hash (%s): %v", txHash.String(), req.Error)
			continue
		}

		result := txRes[i]
		if result != nil && result.IsQuarantined {
			tx := txHashMap[txHash]
			fraudTxs = append(fraudTxs, tx)
		}
	}
	return fraudTxs, nil
}

//
//func main() {
//	// Create an Ethereum RPC client to communicate with the Zircuit endpoint
//	client, err := rpc.Dial("https://zircuit1-mainnet.p2pify.com")
//	if err != nil {
//		log.Fatalf("Failed to create Ethereum RPC client: %v", err)
//	}
//
//	// Create the detector with the Ethereum client
//	detector := &stuckTxDetector{
//		chainClient: client,
//		lggr:        Logger{},
//	}
//
//	// Example transactions (replace with real data)
//	txs := []Tx{
//		{TxAttempts: []TxAttempt{{Hash: common.HexToHash("0xbf7179cbeb760388d972b68dfeef0ee687dd1286d5b4cd534b6b75c42bbdf3c4")}}},
//		{TxAttempts: []TxAttempt{{Hash: common.HexToHash("0x3f9d8133af4ab1ca64c4a6fa572161995eb88fd942d71a61910a04e1639ee51e")}}},
//	}
//
//	// Detect fraudulent transactions
//	ctx := context.Background()
//	fraudTxs, err := detector.detectFraudTransactionsZircuit(ctx, txs)
//	if err != nil {
//		log.Fatalf("Error detecting fraud transactions: %v", err)
//	}
//
//	// Print out the fraudulent transactions
//	fmt.Printf("Fraudulent transactions: %+v\n", fraudTxs)
//}
